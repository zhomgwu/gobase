package acinbase

import (
	"sync"
	//"fmt"
)

type ServiceManager struct {
	serviceID uint64							//自增id
	sessionID uint64							//请求id
	syncMutex *sync.Mutex 						//服务映射表同步锁
	syncWaiter *sync.WaitGroup					//同步所有service
	mapService map[uint64]*Service 				//服务映射表
	mapSyncService map[uint64]*sync.Mutex 		//服务锁，当服务被其它协程调用时锁住，保证同一时刻只有一条调用协程占用
}

func (self *ServiceManager)NewService(handler IServiceHandler) *Service {
	self.syncMutex.Lock()
	defer self.syncMutex.Unlock()

	self.serviceID++
	service := NewService(self.serviceID, handler)
	service.Waiter(self.syncWaiter)
	
	go service.GoRunner()

	self.mapService[self.serviceID] = service
	self.mapSyncService[self.serviceID] = new(sync.Mutex)
	return service
}

func (self *ServiceManager)Run() {
	self.syncWaiter.Wait()
}

func (self *ServiceManager)CallAsync(service *Service, event ServiceEvent) bool {
	return self.CallAsyncByID(service.GetSID(), event)
}

func (self *ServiceManager)CallAsyncByID(sid uint64, event interface{}) bool {
	locker := self.mapSyncService[sid]
	if locker == nil {
		return false
	}

	locker.Lock()
	defer locker.Unlock()

	service := self.mapService[sid]
	if service == nil {
		return false
	}

	go func() {
		service.RecvEvent(event)
	}()
	return true
}

func (self *ServiceManager)CallSync(service *Service, event ServiceEvent) interface{} {
	return self.CallSyncByID(service.GetSID(), event)
}

func (self *ServiceManager)CallSyncByID(sid uint64, event ServiceEvent) interface{} {
	locker := self.mapSyncService[sid]
	if locker == nil {
		return nil
	}

	locker.Lock()
	defer locker.Unlock()

	service := self.mapService[sid]
	if service == nil {
		return nil
	}
	
	event.SynCall = true
	service.RecvEventWait(event)
	return <-service.writec
}

func (self *ServiceManager)RemoveService(sid uint64) {
	self.syncMutex.Lock()
	defer self.syncMutex.Unlock()

	delete(self.mapService, sid)
	delete(self.mapSyncService, sid)
}

func NewServiceManager() *ServiceManager {
	self := &ServiceManager {
		serviceID: 0,
		syncMutex: new(sync.Mutex),
		syncWaiter: new(sync.WaitGroup),
		mapService : make(map[uint64]*Service),
		mapSyncService : make(map[uint64]*sync.Mutex),
	}
	return self
}
