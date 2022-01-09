package netdemo

import (
	"container/list"
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"math/rand"
	"net"
	"sync"
	"testing"
	"time"

	"github.com/golang/protobuf/proto"

	"github.com/Sogues/SG/netdemo/proto/proto_csmsg"
)

var (
	MsgChan     = make(chan *MsgSt, 1024)
	playerGroup = sync.Map{}
	colorGroup  = []color4{
		{1, 0, 0, 1},
		{1, 1, 0, 1},
		{1, 1, 1, 1},
		{0, 1, 0, 1},
		{0, 1, 1, 1},
		{0, 0, 1, 1},
		{0, 0, 0, 1},
	}
	colorIdx int
)

func sendFn(id uint32, msg proto.Message) (out []byte) {
	fmt.Println("send", id, msg)
	by, _ := proto.Marshal(msg)
	size := len(by)
	totalSize := 12 + size
	out = make([]byte, totalSize)
	binary.BigEndian.PutUint32(out, uint32(totalSize))
	binary.BigEndian.PutUint32(out[4:], id)
	copy(out[12:], by)
	return
}

func Login(p *playerUnit) {
	scLogin := &proto_csmsg.SC_Login{
		Mine:   nil,
		Theirs: nil,
	}
	scLogin.Mine = p.DetailPlayerInfo()
	playerGroup.Range(func(key, value interface{}) bool {
		player := value.(*playerUnit)
		scLogin.Theirs = append(scLogin.Theirs, player.DetailPlayerInfo())
		player.conn.Write(
			sendFn(
				uint32(proto_csmsg.MSG_ID_MSG_ID_SC_NtfLogin),
				&proto_csmsg.SC_NtfLogin{Player: scLogin.GetMine()},
			),
		)
		return true
	})
	p.conn.Write(
		sendFn(uint32(proto_csmsg.MSG_ID_MSG_ID_SC_LogIn), scLogin),
	)
	playerGroup.Store(p.addr, p)
}

func Logoff(addr string) {
	playerGroup.Delete(addr)
	msg := &proto_csmsg.SC_LogoutNtf{Addr: addr}
	playerGroup.Range(func(key, value interface{}) bool {
		player := value.(*playerUnit)
		player.conn.Write(
			sendFn(uint32(proto_csmsg.MSG_ID_MSG_ID_SC_LogoutNtf), msg),
		)
		return true
	})
}

func UpdateMove() {
	playerGroup.Range(func(key, value interface{}) bool {
		player := value.(*playerUnit)
		player.UpdateMove()
		return true
	})
}

func SendDiff() {
	syncMove := &proto_csmsg.SC_SyncMove{}
	playerGroup.Range(func(key, value interface{}) bool {
		player := value.(*playerUnit)
		diff := player.fillDiff()
		if nil != diff {
			syncMove.Diffs = append(syncMove.Diffs, diff)
		}
		return true
	})
	if 0 == len(syncMove.Diffs) {
		return
	}
	playerGroup.Range(func(key, value interface{}) bool {
		player := value.(*playerUnit)
		player.conn.Write(
			sendFn(
				uint32(proto_csmsg.MSG_ID_MSG_ID_SC_SyncMove),
				syncMove),
		)
		return true
	})
}

const (
	statePos = 1
)

type (
	vec3   [3]float32
	color4 [4]float32
	MsgSt  struct {
		msg  proto.Message
		addr string
	}
	MoveList struct {
		moves      list.List
		lastMoveTm float32
	}
	playerUnit struct {
		addr string
		conn net.Conn

		color color4
		loc   vec3

		moveList MoveList

		state uint8
	}
)

func (v vec3) EQ(t vec3) bool {
	return v[0] == t[0] && v[1] == t[1] && v[2] == t[2]
}

func (v vec3) Add(t vec3) vec3 {
	v[0] += t[0]
	v[1] += t[1]
	v[2] += t[2]
	return v
}

func (v vec3) Mul(f float32) vec3 {
	v[0] *= f
	v[1] *= f
	v[2] *= f
	return v
}

func (m *MoveList) Add(move *proto_csmsg.CS_SyncMove_Move) bool {
	if move.GetTm() <= m.lastMoveTm {
		return false
	}
	if 0 == m.lastMoveTm {
		move.Delta = 0
	} else {
		move.Delta = move.GetTm() - m.lastMoveTm
	}
	m.lastMoveTm = move.GetTm()
	m.moves.PushBack(move)
	return true
}

func (p *playerUnit) SyncMove(move *proto_csmsg.CS_SyncMove) {
	for _, v := range move.GetMoves() {
		p.moveList.Add(v)
	}
}

func (p *playerUnit) UpdateMove() {
	oldLoc := p.loc
	for it := p.moveList.moves.Front(); nil != it; it = it.Next() {
		move := it.Value.(*proto_csmsg.CS_SyncMove_Move)
		v, h := float32(0.0), float32(0)
		if move.InputState.GetUp() {
			v = 1
		} else if move.InputState.GetDown() {
			v = -1
		}
		if move.InputState.GetLeft() {
			h = -1
		} else if move.InputState.GetRight() {
			h = 1
		}
		p.loc = p.loc.Add(vec3{h, 0, v}.Mul(move.Delta).Mul(10))
	}
	p.moveList.moves.Init()
	if !oldLoc.EQ(p.loc) {
		p.state |= statePos
	}
}

func (p *playerUnit) fillDiff() *proto_csmsg.SC_SyncMove_PlayerDiff {
	if 0 == p.state {
		return nil
	}
	out := &proto_csmsg.SC_SyncMove_PlayerDiff{}
	out.Addr = p.addr
	if 0 != p.state&statePos {
		out.PosDiff = &proto_csmsg.SC_SyncMove_PlayerDiff_PositionDiff{
			Pos: &proto_csmsg.Position{
				X: p.loc[0],
				Y: p.loc[1],
				Z: p.loc[2],
			},
			Tm: p.moveList.lastMoveTm,
		}
	}
	p.state = 0
	return out
}

func (p *playerUnit) DetailPlayerInfo() *proto_csmsg.PlayerInfo {
	out := &proto_csmsg.PlayerInfo{
		Addr: p.addr,
		Position: &proto_csmsg.Position{
			X: p.loc[0],
			Y: p.loc[1],
			Z: p.loc[2],
		},
		PlayerColor: &proto_csmsg.PlayerColor{
			A: p.color[0],
			B: p.color[1],
			C: p.color[2],
			D: p.color[3],
		},
	}
	return out
}

func TestEchoServer(t *testing.T) {
	ln, err := net.Listen(`tcp`, `:7777`)
	if nil != err {
		panic(err)
	}
	connChan := make(chan net.Conn, 1024)
	go func() {
		for {
			conn, err := ln.Accept()
			if nil != err {
				fmt.Println(err)
				continue
			}
			fmt.Println("start from", conn.RemoteAddr())
			connChan <- conn
		}
	}()
	// 30帧
	tk := time.NewTicker(time.Second / 33)

	for {
		select {
		case <-tk.C:
			var cnt int

		outFor:
			for cnt < 10 {
				cnt++
				select {
				case msg := <-MsgChan:
					p, ok := playerGroup.Load(msg.addr)
					if !ok {
						continue
					}
					player := p.(*playerUnit)
					switch sp := msg.msg.(type) {
					case *proto_csmsg.CS_SyncMove:
						player.SyncMove(sp)
					}
				default:
					break outFor
				}
			}

			UpdateMove()

			SendDiff()

		case conn := <-connChan:
			go func() {
				for {
					err := func() error {
						var msgLen [4]byte
						_, err := io.ReadFull(conn, msgLen[:])
						if nil != err {
							return err
						}
						size := binary.BigEndian.Uint32(msgLen[:])
						if size < 12 || size > 1<<16 {
							return errors.New(fmt.Sprintf("size %v", size))
						}
						arr := make([]byte, size-4)
						_, err = io.ReadFull(conn, arr[:])
						if nil != err {
							return err
						}

						cmdId := binary.BigEndian.Uint32(arr[:])
						switch cmdId {
						case uint32(proto_csmsg.MSG_ID_MSG_ID_CS_Login):
							addr := conn.RemoteAddr().String()
							if _, ok := playerGroup.Load(addr); ok {
								fmt.Println(addr, "duplicate login")
								return nil
							}
							p := &playerUnit{}
							p.addr = conn.RemoteAddr().String()
							p.conn = conn
							p.color = colorGroup[colorIdx%len(colorGroup)]
							colorIdx++
							p.loc = vec3{rand.Float32() * 10, 0, rand.Float32() * 10}
							p.moveList.moves.Init()
							Login(p)
						case uint32(proto_csmsg.MSG_ID_MSG_ID_CS_SyncMove):
							msg := &proto_csmsg.CS_SyncMove{}
							err = proto.UnmarshalMerge(arr[8:], msg)
							if nil != err {
								return err
							}
							fmt.Println(conn.RemoteAddr(), "sync move", msg)
							MsgChan <- &MsgSt{
								msg:  msg,
								addr: conn.RemoteAddr().String(),
							}
						}
						return nil
					}()
					if nil != err {
						fmt.Println(conn.RemoteAddr(), "close", err)
						conn.Close()
						Logoff(conn.RemoteAddr().String())
						return
					}
				}
			}()
		}
	}
}
