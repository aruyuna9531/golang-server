package main

import (
	"fmt"
	"net"
)

func main() {
	conn, err := net.Dial("tcp", "127.0.0.1:8080")
	if err != nil {
		fmt.Println("client dial err: " + err.Error())
		return
	}
	conn.Write([]byte("hello?"))
}