package netdemo

import (
	"fmt"
	"net"
	"testing"
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
				var bt [1024]byte
				for {
					l, err := conn.Read(bt[:])
					if nil != err {
						fmt.Println(conn.RemoteAddr(), err)
						conn.Close()
						return
					}
					fmt.Println("receive from", conn.RemoteAddr(), string(bt[:l]))
					conn.Write(bt[:l])
				}
			}()
		}
	}
}
