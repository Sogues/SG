package sync_demo

import (
	"encoding/binary"
	"fmt"
	"math"
	"net"
)

var (
	HELO = binary.BigEndian.Uint32([]byte("HELO"))
	WLCM = binary.BigEndian.Uint32([]byte("WLCM"))
	STAT = binary.BigEndian.Uint32([]byte("STAT"))
	INPT = binary.BigEndian.Uint32([]byte("INPT"))

	GOBJ = binary.BigEndian.Uint32([]byte("GOBJ"))
	RCAT = binary.BigEndian.Uint32([]byte("RCAT"))
	RPLM = binary.BigEndian.Uint32([]byte("RPLM"))
)

var (
	CMDG = map[uint32]string{
		HELO: `HELO`,
		WLCM: `WLCM`,
		STAT: `STAT`,
		INPT: `INPT`,

		GOBJ: `GOBJ`,
		RCAT: `RCAT`,
		RPLM: `RPLM`,
	}

	ng = &NetMgr{}
)

type (
	netBuffer struct {
		data    []byte
		bitHead int
	}
	netMsg struct {
		receivedTm float32
		data       *netBuffer
		addr       *net.UDPAddr
	}
	NetMgr struct {
		ln *net.UDPConn

		msgChan chan *netMsg

		addrMap map[string]*ClientProxy
		idMap   map[uint64]*ClientProxy

		id uint64
	}
)

func (n *netBuffer) readBits(bitCount int) uint8 {
	byteOffset := n.bitHead >> 3
	bitOffset := n.bitHead & 0x7
	out := n.data[byteOffset] >> bitOffset

	bitsFreeThisByte := 8 - bitOffset
	if bitsFreeThisByte < bitCount {
		out |= (n.data[byteOffset+1] << bitsFreeThisByte) & 0xff
	}
	out &= ^(0xff << bitCount)
	n.bitHead += bitCount
	return out
}

func (n *netBuffer) writeBits(data uint8, bitCount int) {
	nextBitHead := n.bitHead + bitCount

	byteOffset := n.bitHead >> 3
	bitOffset := n.bitHead & 0x7

	var currentMask uint8 = ^(0xff << bitOffset)
	n.data[byteOffset] = (n.data[byteOffset] & currentMask) | data<<bitOffset

	bitsFreeThisByte := 8 - bitOffset

	if bitsFreeThisByte < bitCount {
		n.data[byteOffset+1] = data >> bitsFreeThisByte
	}
	n.bitHead = nextBitHead
}

func (n *netBuffer) readInputState() InputState {
	out := InputState{}
	if n.readBool() {
		if n.readBool() {
			out.r = 1
		} else {
			out.r = -1
		}
	}
	if n.readBool() {
		if n.readBool() {
			out.f = 1
		} else {
			out.b = -1
		}
	}
	if n.readBool() {
		out.s = 1
	}
	return out
}

func (n *netBuffer) writeInputState(input InputState) {
	h := input.GetH()
	n.writeBool(0 != h)
	if 0 != h {
		n.writeBool(h > 0)
	}

	v := input.GetV()
	n.writeBool(0 != v)
	if 0 != v {
		n.writeBool(v > 0)
	}
	n.writeBool(0 != input.s)
}

func (n *netBuffer) writeBool(b bool) {
	if b {
		n.writeBits(1, 1)

	} else {
		n.writeBits(0, 1)

	}
}

func (n *netBuffer) writeFloat(f float32) {
	u := math.Float32bits(f)
	n.writeUint32(u)
}

func (n *netBuffer) readFloat() float32 {
	u := n.readUint32()
	f := math.Float32frombits(u)
	return f
}

func (n *netBuffer) readBool() bool {
	return 0 != n.readBits(1)
}

func (n *netBuffer) writeInteger(val interface{}) {
	var by []byte
	switch tv := val.(type) {
	case uint8:
		by = append(by, byte(tv))
	case uint16:
		by = append(by, byte(tv>>8), byte(tv))
	case packetSeqNum:
		by = append(by, byte(tv>>8), byte(tv))
	case uint32:
		by = append(by, byte(tv>>24), byte(tv>>16), byte(tv>>8), byte(tv))
	default:
		panic(tv)
	}
	for i := 0; i < len(by); i++ {
		n.writeBits(by[len(by)-1-i], 8)
	}
}

func (n *netBuffer) readInteger(i interface{}) uint32 {
	var bitCount int
	switch i.(type) {
	case uint8:
		bitCount = 8
	case uint16, packetSeqNum:
		bitCount = 16
	case uint32:
		bitCount = 32
	}
	b := n.readSpBits(bitCount)
	switch len(b) {
	case 1:
		return uint32(b[0])
	case 2:
		return uint32(b[1]) | uint32(b[0])<<8
	case 3:
		return uint32(b[2]) | uint32(b[1])<<8 | uint32(b[0])<<16
	case 4:
		return uint32(b[3]) | uint32(b[2])<<8 | uint32(b[1])<<16 | uint32(b[0])<<24
	default:
		panic(len(b))
	}
}

func (n *netBuffer) writeUint32(val uint32) {
	var data [4]byte
	binary.BigEndian.PutUint32(data[:], val)
	for i := 0; i < 4; i++ {
		n.writeBits(data[3-i], 8)
	}
}

func (n *netBuffer) writeUint64(val uint64) {
	n.writeUint32(uint32(val))
	return
	var data [8]byte
	binary.BigEndian.PutUint64(data[:], val)
	for i := 0; i < 8; i++ {
		n.writeBits(data[7-i], 8)
	}
}

func (n *netBuffer) writeSpBits(b []byte, bitCount int) {
	size := bitCount >> 3
	for i := 0; i < size; i++ {
		n.writeBits(b[len(b)-1-i], 8)
		bitCount -= 8
	}
	if 0 != bitCount {
		n.writeBits(b[0], bitCount)
	}
}

func (n *netBuffer) readSpBits(bitCount int) []byte {
	size := bitCount >> 3
	data := make([]byte, size)
	for i := 0; i < size; i++ {
		data[size-1-i] = n.readBits(8)
		bitCount -= 8
	}
	if 0 != bitCount {
		tmp := append([]byte{}, n.readBits(bitCount))
		data = append(tmp, data...)
	}
	return data
}

func (n *netBuffer) readUint32() uint32 {
	var out [4]byte
	for i := 0; i < 4; i++ {
		out[3-i] = n.readBits(8)
	}
	return binary.BigEndian.Uint32(out[:])
}

func (n *netBuffer) readUint64() uint64 {
	var out [8]byte
	for i := 0; i < 8; i++ {
		out[7-i] = n.readBits(8)
	}
	return binary.BigEndian.Uint64(out[:])
}

func (n *netBuffer) writeString(s string) {
	size := len(s)
	n.writeUint32(uint32(size))
	for _, v := range s {
		n.writeBits(uint8(v), 8)
	}
}

func (n *netBuffer) readString() string {
	size := n.readUint32()
	out := make([]byte, size)
	for i := range out {
		out[i] = n.readBits(8)
	}
	return string(out)
}

func (n *NetMgr) Init(port int) {
	ln, err := net.ListenUDP(`udp`, &net.UDPAddr{
		IP:   net.IPv4zero,
		Port: port,
		Zone: "",
	})
	if nil != err {
		panic(err)
	}
	n.msgChan = make(chan *netMsg, 1024)
	n.ln = ln
	n.ReadIntoQueue()
	n.id = 100000
}

func (n *NetMgr) ProcessPackets() {
	n.ProcessQueue()
}

func (n *NetMgr) ReadIntoQueue() {
	go func() {
		for {
			var packet [1500]byte
			size, addr, err := n.ln.ReadFromUDP(packet[:])
			if nil != err {
				fmt.Println(err)
				continue
			}
			tm := TimingIst.GetTimeF()
			fmt.Println(tm, addr, size)
			n.msgChan <- &netMsg{
				receivedTm: tm,
				data: &netBuffer{
					data: packet[:size],
				},
				addr: addr,
			}
		}
	}()
}

func (n *NetMgr) ProcessQueue() {
	var cnt int = 10
	for cnt > 0 {
		select {
		case msg := <-n.msgChan:
			n.ProcessMsg(msg)
			cnt--
		default:
			return
		}
	}
}

func (n *NetMgr) ProcessMsg(msg *netMsg) {
	cp, ok := n.addrMap[msg.addr.String()]
	if !ok {
		n.HandleNewClient(msg)
	} else {
		tp := msg.data.readUint32()
		msg.data.bitHead = 0
		fmt.Println(tp, CMDG[tp])
		n.HandleMsgWithClient(msg, cp)
	}
}

func (n *NetMgr) HandleNewClient(msg *netMsg) {
	packetType := msg.data.readUint32()
	switch packetType {
	case HELO:
		name := msg.data.readString()
		fmt.Println(msg.addr, name)
		n.id++
		cp := &ClientProxy{
			updAddr: msg.addr,
			addr:    msg.addr.String(),
			name:    name,
			id:      n.id,
		}
		if nil == n.addrMap {
			n.addrMap = map[string]*ClientProxy{}
		}
		n.addrMap[cp.addr] = cp
		if nil == n.idMap {
			n.idMap = map[uint64]*ClientProxy{}
		}
		n.idMap[cp.id] = cp

		cg.AddCP(cp)

		n.SendWelcome(cp)

		for _, v := range cg.entries {
			v.rpl.ReplicateCreate(cp.id, dirtyAll)
		}
	default:
		fmt.Println("unknown packet", packetType, " from", msg.addr)
	}
}

func (n *NetMgr) HandleMsgWithClient(msg *netMsg, cp *ClientProxy) {
	packetType := msg.data.readUint32()
	switch packetType {
	case HELO:
		n.SendWelcome(cp)
	case INPT:
		if cp.dn.ReadAndProcessState(msg.data) {
			n.handleInput(msg.data, cp)
		}
	default:
		fmt.Println("unknown packet", packetType, " from", msg.addr)
	}
}

func (n *NetMgr) Gen(b []byte) uint32 {
	switch len(b) {
	case 1:
		return uint32(b[0])
	case 2:
		return uint32(b[1]) | uint32(b[0])<<8
	case 3:
		return uint32(b[2]) | uint32(b[1])<<8 | uint32(b[0])<<16
	case 4:
		return uint32(b[3]) | uint32(b[2])<<8 | uint32(b[1])<<16 | uint32(b[0])<<24
	default:
		panic(len(b))
	}
	return 0
}

func (n *NetMgr) handleInput(data *netBuffer, cp *ClientProxy) {
	// todo
	//sq := n.Gen(data.readSpBits(16))
	//if data.readBool() {
	//	start := n.Gen(data.readSpBits(16))
	//	hasCount := data.readBool()
	//	var ack uint32
	//	if hasCount {
	//		ack = n.Gen(data.readSpBits(8))
	//	}
	//	fmt.Println(cp.addr, "queue info", sq, start, hasCount, ack)
	//}
	moveCount := n.Gen(data.readSpBits(2))
	fmt.Println(cp.addr, "input move count", moveCount)
	for ; moveCount > 0; moveCount-- {
		is := data.readInputState()
		ts := math.Float32frombits(data.readUint32())
		fmt.Println(cp.addr, "input", is, ts)
		if cp.moveList.AddMoveIfNew(&Move{
			inputState: is,
			timestamp:  ts,
			deltaTime:  0,
		}) {
			cp.mIsLastMoveTimestampDirty = true
		}
	}
}

func (n *NetMgr) SendWelcome(cp *ClientProxy) {
	b := netBuffer{
		data:    make([]byte, 1500),
		bitHead: 0,
	}
	b.writeUint32(WLCM)
	b.writeUint64(cp.id)
	_, err := n.ln.WriteToUDP(b.data[:(b.bitHead+7)>>3], cp.updAddr)
	if nil != err {
		fmt.Println(cp.updAddr, "send error", err)
	}
}
