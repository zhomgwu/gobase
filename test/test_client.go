
package main 

import (
	"fmt"
	"net"
	"strconv"
	"sync"
	"time"
	//"bytes"
	//"encoding/binary"
)

func main() {
	//connectOnce()	
	connectForever()
}

func connectOnce() {
	conn, err := net.Dial("tcp", "127.0.0.1:10123")	
	if err != nil {
		fmt.Println("connection error!")
		return
	}
	
	str := "hello go net!! "
	data := []byte(str)
	conn.Write(data)

	for {
		readbuf := make([]byte, 65535)
		n,err := conn.Read(readbuf)		
		fmt.Println("read buffer", readbuf, n,err)
	}
	conn.Close()
}

func connectForever(){
	var i int64 = 1
	var waiter sync.WaitGroup
	var count int = 10000
	waiter.Add(count)
	for idx:=0; idx<count; idx++{
		go func() {
			defer waiter.Done()
			
			conn, err := net.Dial("tcp", "127.0.0.1:10123")	
			if err != nil {
				fmt.Println("connection error!", err)
				return
			}
			defer conn.Close()

			str := "hello go net!!"+strconv.FormatInt(i, 10)
			data := []byte(str)
			start := time.Now()
			conn.Write(data)
			recv := make([]byte, 1024) 
			n, werr := conn.Read(recv)
			if werr != nil {
				return
			}

			n = n
			//var then int32
			//binary.Read(bytes.NewBuffer(recv), binary.BigEndian, &then)  

			fmt.Println("recv", string(recv))
			fmt.Println("spend time", time.Now().Sub(start).Nanoseconds())
			i++	
		
			//time.Sleep(5 * time.Minute)
		}()
		// 1秒100次，让发送平滑
		time.Sleep(10 * time.Millisecond)
	}

	waiter.Wait()
}