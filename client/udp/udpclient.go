package main

import (
	"fmt"
	"go_svr/log"
	"go_svr/share/rudp"
	"net"
	"os"
	"os/signal"
	"syscall"
)

// UDP 客户端
func main() {
	conn, err := net.DialUDP("udp", nil, &net.UDPAddr{IP: net.IPv4zero, Port: 8888})
	if err != nil {
		fmt.Println(err)
		return
	}
	rconn := rudp.NewConn(conn, rudp.New())
	for !rconn.Connected() {
	}
	rconn.Write([]byte("rudp hello"))
	go func() {
		var b [1024]byte
		for {
			n, err := rconn.Read(b[:])
			if err != nil {
				log.Error("eof")
				return
			}
			if n == 0 {
				continue
			}
			log.Debug("rudp msg: %s", b[:n])
		}
	}()
	sc := make(chan os.Signal)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM)
	<-sc
}
