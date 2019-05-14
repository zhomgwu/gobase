
package acinbase

import (
	"sync"
	"sync/atomic"
)

type ConnectionManager struct {
	connectionByID sync.Map		//[connectionId] *Connection
	genID int64
	count int64
}

func (self *ConnectionManager)Add(conn *Connection) {
	self.count = atomic.AddInt64(&self.count, 1)
	ID := atomic.AddInt64(&self.genID, 1)

	conn.SetConnectionID(self.genID)
	self.connectionByID.Store(ID, conn)
}

func (self *ConnectionManager)Count() int64 {
	return self.count
}

func (self *ConnectionManager)Remove(id int64) {
	if self.GetConnection(id) == nil {
		//连接不存在
		return
	}

	self.count = atomic.AddInt64(&self.count, -1)
	self.connectionByID.Delete(id)
}

func (self *ConnectionManager)GetConnection(id int64) *Connection{
	if v, ok := self.connectionByID.Load(id); ok {
		return v.(*Connection)
	}
	return nil
}

func (self *ConnectionManager)Kick(id int64) {
	if v, ok := self.connectionByID.Load(id); ok {
		if conn, cok := v.(*Connection); cok {
			conn.Close()
		}
		self.connectionByID.Delete(id)
	}
}

func (self *ConnectionManager)CloseAll() {
	self.connectionByID.Range(func(key, value interface{})bool {
		conn := value.(*Connection)
		conn.Close()
		// 测试遍历中删除是安全的
		self.connectionByID.Delete(key)
		return true
	})
	
	self.count = 0
}

func NewConnectionManager() *ConnectionManager{
	self := &ConnectionManager {
		genID: 0,
		count: 0,
	}
	return self
}
