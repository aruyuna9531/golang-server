package mytcp

import (
	"encoding/json"
	"errors"
	"go_svr/log"
	"io"
	"net"
	"strconv"
)

// 假设这是客户端发过来的东西
type ClientLogin struct {
	UserId uint64 `json:"user_id"`
}

type ClientPack struct {
	UserId uint64
	Msg    []byte
}

type ClientInfo struct {
	UserId uint64
	conn   net.Conn
}

type TcpServer struct {
	Port       int
	listener   net.Listener
	conns      map[string]*ClientInfo
	clientMsgs chan *ClientPack
}

var tcpSvr = &TcpServer{}

func GetTcpSvr() *TcpServer {
	return tcpSvr
}

func (ts *TcpServer) GetMsgChan() <-chan *ClientPack {
	return ts.clientMsgs
}

func (ts *TcpServer) handleConnection(conn net.Conn) {
	defer conn.Close()
	defer delete(ts.conns, conn.RemoteAddr().String())
	readbytes := make([]byte, 1024)
	for {
		readSize, err := conn.Read(readbytes)
		if err != nil {
			if errors.Is(err, io.EOF) {
				// client主动发起关闭
				log.Printf("client actively closed")
				return
			}
			log.Printf("client read error: " + err.Error())
			return
		}
		js := readbytes[:readSize]
		addr := conn.RemoteAddr()
		if _, e := ts.conns[addr.String()]; !e {
			// 当然，这里还得判断一下是不是login data。这里先从略
			ld := &ClientLogin{}
			err := json.Unmarshal(js, ld)
			if err != nil {
				log.Printf("unmarshal login data error: %s", err.Error())
				return
			}
			ts.conns[addr.String()] = &ClientInfo{
				UserId: ld.UserId,
				conn:   conn,
			}
			log.Printf("remote login %s success, user id %d", addr.String(), ld.UserId)
			continue
		}
		if u, e := ts.conns[addr.String()]; e {
			ts.clientMsgs <- &ClientPack{
				UserId: u.UserId,
				Msg:    js,
			}
		} else {
			log.Printf("illegal connection source: %s", addr.String())
			return
		}
	}
}

func (ts *TcpServer) Create(port int) {
	ts.Port = port
	log.Printf("creating tcp server at port %d...", ts.Port)
	var err error

	ts.listener, err = net.Listen("tcp", "0.0.0.0:"+strconv.Itoa(ts.Port))
	if err != nil {
		// handle error
		log.Printf("Error: " + err.Error())
		return
	}

	if ts.listener == nil {
		log.Printf("Error: socket is nil")
		return
	}

	ts.conns = make(map[string]*ClientInfo)
	ts.clientMsgs = make(chan *ClientPack, 1000)
}

func (ts *TcpServer) OnClose() {
	defer close(ts.clientMsgs)
	for addr, conn := range ts.conns {
		err := conn.conn.Close()
		if err != nil {
			log.Printf("connection close to %s error: %s", addr, err.Error())
			continue
		}
		log.Printf("connection close to %s success", addr)
	}
	err := ts.listener.Close()
	if err != nil {
		log.Printf("TcpServer close error: %s\n", err.Error())
		return
	}
	log.Printf("tcp server closed")
}

func (ts *TcpServer) OnLoop() {
	for {
		conn, err := ts.listener.Accept()
		if err != nil {
			var netErr *net.OpError
			if errors.As(err, &netErr) && errors.Is(netErr.Err, net.ErrClosed) {
				// 因为listener被关闭，中断了Accept过程（已经停止服务）
				log.Printf("tcp server Accept terminated because listener is closed")
				return
			}
			// 其他错误
			log.Printf("tcp server Accept error: %s\n", err.Error())
			continue
		}
		go ts.handleConnection(conn)
	}
}
