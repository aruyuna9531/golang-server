package rudp

import (
	"fmt"
	"go_svr/dependency"
	"go_svr/log"
	"go_svr/proto_codes/rpc"
	"google.golang.org/protobuf/proto"
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
		go ri.Conn(rconn)
	}
}

func (ri *RudpInst) Conn(rconn *RudpConn) {
	for {
		data := make([]byte, MAX_PACKAGE)
		n, err := rconn.Read(data)
		if err != nil {
			fmt.Printf("read err %s\n", err)
			return
		}
		r := &rpc.InputData{}
		err = proto.Unmarshal(data[:n], r) // todo 看看能不能省掉这一步
		if err != nil {
			log.Error("unmarshal error")
			continue
		}
		dependency.PushCmd(r)
	}
}

func (ri *RudpInst) Broadcast(contents []byte) {
	ri.listener.RudpBroadcast(contents)
}
