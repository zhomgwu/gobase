package acinbase

import (
	"net"
	"fmt"
	"sync"
	"time"
	//"strconv"
)

const (
	MAX_PACKAGE_SIZE = 65535
	MAX_SEND_BUFFER = 65535
)

type IConnectionHandlder interface {
	OnConnected(*Connection)
	OnDisconnected(*Connection)
	OnMessage(*Connection, []byte)
}

type Connection struct {
	connMgr *ConnectionManager  //连接管理器
	conn *net.Conn 				//连接实体
	connID int64				//connection id, 由ConnectionManager分配， 非fd
	validity bool				//是否通过验证
	close bool					//是否关闭
	userdata interface{}		//用户对象
	handler IConnectionHandlder //消息处理器
	recvBuffer []byte			//接收缓存区
	recvLength int 				//接收数据长度
	sendBuffer []byte			//待发送缓存区
	sendLength int 				//待发送长度

	exitSync sync.WaitGroup		//同步变量，等待Read， Write协程退出
	sendMutex *sync.Mutex		//发送缓存区锁
}

func (self *Connection)Send(data []byte) {
	// copy入发送缓存区，会有协程自动发送
	datalen := len(data)
	if datalen > 1400 {
		fmt.Printf("WARNING: this package larger than 1400b, maybe you can optimize it!")
	}

	if datalen > MAX_PACKAGE_SIZE {
		fmt.Printf("ERROR: this package larger than MAX_PACKAGE_SIZE=%d, it will be discarded!!\n", MAX_PACKAGE_SIZE)
		return
	}

	// copy
	self.sendMutex.Lock()
	defer self.sendMutex.Unlock()
	copy(self.sendBuffer[self.sendLength:], data)
	self.sendLength += datalen
}


func (self *Connection)SetUserData(userdata interface{}) {
	self.userdata = userdata
}

func (self *Connection)UserData() interface{} {
	return self.userdata
}

func (self *Connection)SetHandler(handler IConnectionHandlder) {
	self.handler = handler
}

func (self *Connection)SetConnectionID(ID int64) {
	self.connID = ID
}

func (self *Connection)ConnectionID() int64 {
	return self.connID
}

func (self *Connection)OnConnected() {
	if self.handler != nil {
		self.handler.OnConnected(self)
	}
}

func (self *Connection)OnDisconnected() {
	if self.handler != nil {
		self.handler.OnDisconnected(self)
	}
}

func (self *Connection)OnMessage(data []byte) {
	if self.handler != nil {
		self.handler.OnMessage(self, data)
	}
}

func (self *Connection)Start() {
	// 开始的时候，为“读/写”协程增加同步等待
	self.exitSync.Add(2)
	go self.GoRead()
	go self.GoWrite()
}

func (self *Connection)Close() {
	if self.conn != nil && self.close == false {
		self.close = true
		(*self.conn).Close()
		self.OnDisconnected()

		self.connMgr.Remove(self.connID)
	}
}

func (self *Connection)GoRead() {
	defer self.exitSync.Done()

	readbuf := make([]byte, 1024)
	for {
		n, err := (*self.conn).Read(readbuf)
		if err != nil {
			// fmt.Println("connection error occur when read socket")
			fmt.Println("connection error", err)
			self.Close()
			break
		}
		
		self.OnMessage(readbuf[:n])

/*		// test send
		fmt.Println("recv", string(readbuf), n)
		now := time.Now().Nanosecond()
		self.Send([]byte(strconv.Itoa(now)))
*/

/*
		copy(self.recvBuffer[:self.recvLength], readbuf[:n])
		self.recvLength = self.recvLength + n

		for {
			//数据小于2字节, 前2个字节代表单个包长度，长度包括这2个字节
			if self.recvLength < 2 {
				break
			}

			pkglen := BytesToUint16(self.recvBuffer[:2])
			if self.recvLength < int(pkglen) {
				break
			}
			// 处理消息
			self.OnMessage(self.recvBuffer[2:pkglen])
			// 删除已处理包体
			self.recvBuffer = self.recvBuffer[pkglen:]
			self.recvLength = self.recvLength - int(pkglen)
		}
*/
		// 如果连接关闭，则直接退出，缓存区内的数据会被丢弃
		if self.close {
			break
		}
	}
}

func (self *Connection)GoWrite() {
	defer self.exitSync.Done()
	for {
		if self.close {
			break
		}	
		
		if self.sendLength == 0 {
			// 将自己的时间片让出1毫秒让其它线程有机会执行 
			time.Sleep(time.Nanosecond * 1000)
			continue
		}
		
		//缓存区加锁防止数据清除出错
		self.sendMutex.Lock()
		n, err := (*self.conn).Write(self.sendBuffer[:self.sendLength])
		if err != nil {
			fmt.Println("error occur when send data", err)
			self.Close()
			self.sendMutex.Unlock()
			break
		}

		self.sendBuffer = self.sendBuffer[:n]
		self.sendLength = self.sendLength - n
		self.sendMutex.Unlock()
	}
}

func (self *Connection)WaitExit() {
	self.exitSync.Wait()
}

func NewConnection(conn *net.Conn, connMgr *ConnectionManager) *Connection {
	self := &Connection {
		connMgr : connMgr,
		conn : conn,
		validity : false,
		close : false,
		recvBuffer : make([]byte, MAX_SEND_BUFFER),
		sendBuffer : make([]byte, MAX_SEND_BUFFER),
		sendMutex : new(sync.Mutex),
	}

	connMgr.Add(self)
	return self
}
