package mytcp

import (
	"fmt"
	"net"
)

func handleConnection(conn net.Conn) {
	defer conn.Close()

	readbytes := make([]byte, 1024)
	readSize, err := conn.Read(readbytes)
	if err != nil {
		fmt.Println("client exit, err: " + err.Error())
		return
	}
	addr := conn.RemoteAddr()
	fmt.Printf("string: %s, bytes: %s\n", addr.String(), string(readbytes[:readSize]))
}

func CreateTcpServer() {
	ln, err := net.Listen("tcp", "0.0.0.0:8080")
	if err != nil {
		// handle error
	}
	for {
		conn, err := ln.Accept()
		if err != nil {
			// handle error
			continue
		}
		go handleConnection(conn)
	}
}