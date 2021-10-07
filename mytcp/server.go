package mytcp

import (
	"fmt"
	"net"
	"strconv"
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

func CreateTcpServer(port int) {
	fmt.Println("creating tcp server at port " + strconv.Itoa(port) + "...")
	ln, err := net.Listen("tcp", "0.0.0.0:" + strconv.Itoa(port))
	if err != nil {
		// handle error
		fmt.Println("Error: " + err.Error())
		return
	}

	if ln == nil {
		fmt.Println("Error: socket is nil")
		return
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