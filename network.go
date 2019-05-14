package acinbase

import (
	"net"
	"fmt"
)

type Network struct {
	listener net.Listener 				// listener
	port string 						// 监听端口
	svrDone chan struct{}				// 服务器停止信号
	handler INetworkHandler				// 网络处理器
	connMgr *ConnectionManager  		//连接管理器
}

type INetworkHandler interface {
	OnNewConnection(*Connection)
	OnRemoveConnection(*Connection)
}

func (self * Network)Listen(port string) bool{
	self.port = port

	listener, err := net.Listen("tcp", self.port)
	if err != nil {
		fmt.Println("server listen error", self.port, err)
		return false
	}
	self.listener = listener
	return true
}

func (self * Network)Loop() {
	fmt.Println("network loop!", self.svrDone)
	if self.listener == nil {
		fmt.Println("server listener does not ready!")
		return
	}

	defer self.listener.Close()

	// 放到独立goroutine中，否则Accept会阻塞退出通道svrDone，导致svrDone的消息无法读取
	go func () {
		for {
			conn, err := self.listener.Accept()
			if err != nil {
				fmt.Println("accept", err)
				return
			}

			connection := NewConnection(&conn, self.connMgr)
			self.handler.OnNewConnection(connection)

			go self.GoHandleConnection(connection)
		}
	}()

	// 监听是否要退出
	select {
	case <-self.svrDone:
		self.connMgr.CloseAll()
		return
	//default:
	}
}

func (self * Network)GoHandleConnection(conn * Connection) {
	// 开始处理网络消息
	conn.Start()
	// 等待退出
	conn.WaitExit()
	// 到这里时退出
	self.handler.OnRemoveConnection(conn)
}

func NewNetwork(handler INetworkHandler, svrDone chan struct{}) *Network{
	self := &Network {
		handler : handler,
		svrDone : svrDone,
		connMgr : NewConnectionManager(),
	}

	return self
}
