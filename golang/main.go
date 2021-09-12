package main

import (
	"fmt"
	"net"
	// "net/url"
	"net/http"
)

func MainHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "hello!!")
}

func handleConnection(conn net.Conn) {
	var readbytes []byte
	conn.Read(readbytes)
	addr := conn.RemoteAddr()
	fmt.Printf("network: %s, string: %s, bytes: %s\n", addr.Network(), addr.String(), string(readbytes))
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

func main() {
	fmt.Println("Hello!")

	http.HandleFunc("/", MainHandler)
	http.ListenAndServe("0.0.0.0:8000", nil)
}