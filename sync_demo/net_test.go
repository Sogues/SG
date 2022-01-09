package sync_demo

import (
	"fmt"
	"net"
	"testing"
	"time"
)

func TestNetMgr_ReadIntoQueue(t *testing.T) {
	for k, v := range CMDG {
		fmt.Println(k, v)
	}
	ng.Init(45000)
	tk := time.NewTicker(time.Millisecond)
	ng.ReadIntoQueue()
	for {
		select {
		case <-tk.C:
			ng.ProcessPackets()
			cg.Update()
			cg.SendOutgoingPackets()
		}
	}
}

func TestNetBuffer(t *testing.T) {
	b := netBuffer{
		data:    make([]byte, 1500),
		bitHead: 0,
	}
	b.writeUint32(WLCM)
	b.writeUint64(10)
	fmt.Println(b.bitHead)

	b.bitHead = 0
	fmt.Println(WLCM, b.readUint32())
	fmt.Println(b.readUint64())
}

func TestDial(t *testing.T) {
	ln, err := net.DialUDP(`udp`, nil, &net.UDPAddr{
		IP:   net.IPv4(42, 192, 146, 173),
		Port: 45000,
		Zone: "",
	})
	if nil != err {
		panic(err)
	}
	fmt.Println(ln.WriteToUDP([]byte("xxx"), &net.UDPAddr{
		IP:   net.IPv4(42, 192, 146, 173),
		Port: 45000,
		Zone: "",
	}))
}
