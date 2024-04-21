package rudp

import (
	"fmt"
	"go_svr/log"
	"net"
)

type RudpInst struct {
	listener   *RudpListener
	rconn      map[string]*RudpConn
	listenPort int
}

var inst = &RudpInst{
	listenPort: 8888,
}

func GetInst() *RudpInst {
	return inst
}

func (ri *RudpInst) Init() {
	addr := &net.UDPAddr{IP: net.ParseIP("0.0.0.0"), Port: ri.listenPort}
	conn, err := net.ListenUDP("udp", addr)
	if err != nil {
		fmt.Println(err)
		return
	}
	ri.listener = NewListener(conn)
	ri.rconn = make(map[string]*RudpConn)
	go ri.Loop()
}

func (ri *RudpInst) Loop() {
	for {
		rconn, err := ri.listener.AcceptRudp()
		if err != nil {
			fmt.Printf("accept err %v\n", err)
			return
		}
		ri.rconn[rconn.remoteAddr.String()] = rconn
		data := make([]byte, MAX_PACKAGE)
		n, err := rconn.Read(data)
		if err != nil {
			fmt.Printf("read err %s\n", err)
			return
		}
		log.Debug("rudp data: %s", data[:n])
		rconn.Write([]byte(fmt.Sprintf("received message: %s", data[:n])))
	}
}
