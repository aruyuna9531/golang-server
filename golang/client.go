// 一个测试TCP服务的Demo，go run client.go
// 可以拿这个去测试其他语言写的TCP服务器的连接，端口改成那个服务器开放的端口就行
package main

import (
	"fmt"
	"net"
	"time"
	"bytes"
)

func main() {
	conn, err := net.Dial("tcp", "127.0.0.1:8003")
	if err != nil {
		fmt.Println("client dial err: " + err.Error())
		return
	}
	defer conn.Close()
	conn.Write([]byte("hello?"))
	fmt.Println("program will exit in 10 seconds, return message will only available in this time");

	result := bytes.NewBuffer(nil)
    var buf [512]byte
	go func() {
		for {
			n,err := conn.Read(buf[0:])
			result.Write(buf[0:n])
			if err != nil {
				return
			}
			fmt.Println("client dial return: " + string(result.Bytes()))
		}
	}()
	time.Sleep(10 * time.Second)
}