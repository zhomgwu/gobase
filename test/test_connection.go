package main 

import (
	. "acinbase"
	"fmt"
)


type ServerHandler struct {

}

func (self *ServerHandler) OnServerBegin(* Server){
	fmt.Println("OnServerBegin")
}

func (self *ServerHandler) OnServerEnd(* Server) {
	fmt.Println("OnServerEnd")
}

func (self *ServerHandler) OnNewConnection(* Connection) {
	fmt.Println("OnNewConnection")
}

func (self *ServerHandler) OnRemoveConnection(* Connection) {
	fmt.Println("OnRemoveConnection")
}

type PlayerHandler struct {

}

func (self *PlayerHandler) OnConnected(conn *Connection) {
	fmt.Println("new connection connected")
}

func (self *PlayerHandler) OnDisconnected(conn *Connection) {
	fmt.Println("connection disconnected")
}

func (self *PlayerHandler) OnMessage(conn *Connection, data []byte) {
	/*append (databuffer, data )
	if len(databuffer) >= HEADER_LENGTH {

	}*/

	fmt.Println("connection OnMessage")
}

func main() {
	conf := ServerConfig {
		Port : ":10123",
	}

	svrHandler := &ServerHandler {}
	if svr := NewServer(conf, svrHandler, svrHandler); svr != nil {
		svr.Run()	
	}
}
