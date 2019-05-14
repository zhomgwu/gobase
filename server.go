package acinbase

import (
	"fmt"
	"os"
	"os/signal"
	"sync"
	"runtime/pprof"
	//"time"
)


type ServerConfig struct {
	Port string
}

type Server struct {
	ServerConfig						// 服务器配置
	maxClient int						// 最大用户数
	stop bool							// 是否停止服务器
	handler IServerHandler				// 服务器处理器
	syncWaiter sync.WaitGroup			// 同步信号
	network * Network 					// 网络服务
	svrDone chan struct{}				// 退出信号
}

type IServerHandler interface {
	OnServerBegin(* Server)
	OnServerEnd(* Server)
}

func (self * Server)SetMaxClient(client int) {
	self.maxClient = client
}

func (self * Server)GoSignalHandler() {
	c := make(chan os.Signal, 1)
	// os.Kill 是捕获不到的
	signal.Notify(c, os.Interrupt)

	s:= <-c
	if s == os.Interrupt {
		fmt.Println("capture SIGINT, server will EXIT ...")
		self.StopCPUProfile()
		self.Stop()
	}
}

func (self *Server)StartCPUProfile() {
	f, err := os.Create("cpu.profile")
    	if err != nil {
        fmt.Fprintf(os.Stderr, "Can not create cpu profile output file: %s",
            err)
        return
    }
	pprof.StartCPUProfile(f)
}

func (self *Server)StopCPUProfile() {
	pprof.StopCPUProfile()
}

func (self *Server)GoNetworkHandler() {
	defer self.syncWaiter.Done()
	self.network.Loop()
}

func (self * Server)Run() {
	self.StartCPUProfile()

	self.handler.OnServerBegin(self)
	// 不等待signal
	self.syncWaiter.Add(1)

	go self.GoSignalHandler()
	go self.GoNetworkHandler()

	self.syncWaiter.Wait()
	self.handler.OnServerEnd(self)
}

func (self *Server)Stop() {
	self.stop = true
	// 通知所有模块协程，退出工作
	fmt.Println("close svrdone!", self.svrDone)
	close(self.svrDone)
}

func NewServer(conf ServerConfig, handler IServerHandler, nwHandler INetworkHandler) *Server{
	self := &Server {
		svrDone : make(chan struct{}),
		stop : false,
		handler : handler,
	}

	self.network = NewNetwork(nwHandler, self.svrDone)
	if !self.network.Listen(conf.Port) {
		fmt.Println("network listen errror")
		return nil
	}
	return self
}

