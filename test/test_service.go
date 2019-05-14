package main 

import (
	"fmt"
	"acinbase"
	"time"
)


type ServiceA struct {
	id int
}

func (self *ServiceA)OnCallHandler(event acinbase.ServiceEvent) interface{} {
	fmt.Println("OnCallHandler", event)
	return nil
}
func (self *ServiceA)OnEventHandler(event acinbase.ServiceEvent){
	fmt.Println("OnEventHandler", event)
}
func (self *ServiceA)OnBegin(){
	fmt.Println("OnBegin", self.id)
}
func (self *ServiceA)OnStep() {
	//fmt.Println("OnStep", self.id)
}
func (self *ServiceA)OnStop(){
	fmt.Println("OnStop", self.id)
}

func main() {
	serviceManager := acinbase.NewServiceManager()

	s1 := serviceManager.NewService(&ServiceA{id:123131313131} )
	//s2 := serviceManager.NewService(&ServiceA{id:2} )
	//s2.Stop()

	sval := acinbase.ServiceEvent{
		MsgType: 123,
		Msg: "hello",
	}

	time.AfterFunc(time.Duration(3*time.Second), func() {
		//val := serviceManager.CallSyncByID(s1.GetSID(), sval)
		//va2 := serviceManager.CallSyncByID(s1.GetSID(), sval)
		va3 := serviceManager.CallAsyncByID(s1.GetSID(), sval)
		//fmt.Println("main", val, va2, va3)
		fmt.Println("main 11111", va3)
	})


	for i:= 0; i < 10000;i++ {
		go func(index int) {
			sval := acinbase.ServiceEvent{
				MsgType: index,
				Msg: "hello",
			}
			va3 := serviceManager.CallSyncByID(s1.GetSID(), sval)
			fmt.Printf("main i=%d v=%v\n", index, va3)
		}(i)	
	}
	
	serviceManager.Run()
}
