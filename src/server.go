package src

import (
	"fmt"
	"net"
)

func init() {
	listen, err := net.Listen("tcp",":48889")
	if(err != nil) {
		sendToTextarea("Couldn't start BattleTerminal server: "+err.Error())
	}
	go func() {
		for {
			conn, err := listen.Accept()
			if(err != nil) {
				sendToTextarea("Error accepting data: "+err.Error())
			}
			go connHandle(conn)
		}
	}()
}

func connHandle(conn net.Conn) {
	// try to read what we got
	ch := make(chan []byte)
	eCh := make(chan error)
	go func() {
		for {
			data := make([]byte,512)
			_, err := conn.Read(data)
			if(err != nil) {
				eCh<- err
				return
			}
			ch<- data
		}
	}()
	for {
		select {
		case data := <-ch:
			fmt.Println(string(data))
		case err := <-eCh:
			fmt.Println(err)
			break;
		}
	}
}