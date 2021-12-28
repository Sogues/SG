package netdemo

import (
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"net"
	"testing"

	"github.com/golang/protobuf/proto"

	"github.com/Sogues/SG/netdemo/proto/proto_csmsg"
)

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
	for {
		select {
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
							msg := &proto_csmsg.CS_Login{}
							err = proto.UnmarshalMerge(arr[8:], msg)
							if nil != err {
								return err
							}
							fmt.Println(conn.RemoteAddr(), "receive", msg)
							conn.Write(append(msgLen[:], arr...))
						}
						return nil
					}()
					if nil != err {
						fmt.Println(conn.RemoteAddr(), "close", err)
						conn.Close()
						return
					}
				}
			}()
		}
	}
}
