
package acinbase

import (
	"sync"
	"fmt"
)

type ServiceEvent struct {
	SynCall bool
	SenderSid int
	MsgType int
	Msg interface{}
}

type IServiceHandler interface {
	OnCallHandler(ServiceEvent) interface{}
	OnEventHandler(ServiceEvent)
	OnBegin()
	OnStep()
	OnStop()
}

type Service struct {
	stop bool
	sid uint64
	readc chan interface{}				//读
	writec chan interface{}				//写，阻塞回调时用， 由调用方提供关闭
	done chan interface{}				//完成， 阻塞回调时用， 由调用方提供关闭

	serviceMgr *ServiceManager
	msgQueue *MessageQueue
	handler IServiceHandler
	syncWaiter *sync.WaitGroup
}

func (self * Service)GetSID() uint64 {
	return self.sid
}

func (self * Service)Stop() {
	self.stop = true
}

func (self * Service)Waiter(w *sync.WaitGroup) {
	self.syncWaiter = w
	self.syncWaiter.Add(1)
}

func (self * Service)GoRunner() {
	self.handler.OnBegin()
	for {
		select {
		case recvmsg := <-self.readc:
			event := recvmsg.(ServiceEvent)
			if event.SynCall {		// 如果是同步调用
				ret := self.handler.OnCallHandler(event)
				self.writec<-ret
			} else {				// 异步调用只入队列
				self.msgQueue.Push(recvmsg)	
			}

		case _ = <-self.done:
			
			fmt.Println("end of file!!")
		default:

			for self.msgQueue.Count() > 0 {
				event := self.msgQueue.Pop().(ServiceEvent)
				self.handler.OnEventHandler(event)
			}
			self.handler.OnStep()
		}

		if self.stop && self.msgQueue.Count() == 0 {
			break
		}
	}
	self.handler.OnStop()

	if self.syncWaiter != nil {
		self.syncWaiter.Done()
	}
}

func (self * Service)RecvEvent(event interface{}) {
	self.readc<-event
}

func (self * Service)RecvEventWait(event ServiceEvent) {
	self.readc<-event
}

func NewService(sid uint64, handler IServiceHandler) *Service{
	self := &Service {
		sid : sid,
		stop : false,
		readc : make(chan interface{}),
		writec : make(chan interface{}),
		done : make(chan interface{}),
		msgQueue : NewMessageQueue(),
		syncWaiter : nil,
		handler : handler,
	}
	return self
}


