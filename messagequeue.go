package acinbase

import (
	"sync"
)

type MessageQueue struct {
	queue []interface{}
	lockguard sync.Mutex
}


func (self *MessageQueue)Push(msg interface{}) {
	self.lockguard.Lock()
	defer self.lockguard.Unlock()

	self.queue = append(self.queue, msg)
}

func (self *MessageQueue)Pop() interface{} {
	self.lockguard.Lock()
	defer self.lockguard.Unlock()

	var msg interface{}
	if len(self.queue) > 0 {
		msg = self.queue[0]
		self.queue = self.queue[1:]
	} else {
		return nil
	}
	return msg
}

func (self *MessageQueue)Count() int {
	return len(self.queue)
}

// 可用来查看队列内容，不建议对队列进行内容修改
func (self *MessageQueue)VisitQueue(callback func(int, interface{})) {
	for i, msg := range self.queue {
		callback(i, msg)
	}
}

func NewMessageQueue() *MessageQueue {
	self := &MessageQueue {
		
	}
	return self
}
