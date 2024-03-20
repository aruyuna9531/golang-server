package mytcp

import (
	"encoding/json"
	"errors"
	"go_svr/log"
	"io"
	"net"
	"strconv"
	"sync/atomic"
)

// 假设这是客户端发过来的东西（先不用protobuf）
type ClientPack struct {
	SessionId uint64 `json:"session_id"`
	Msg       []byte `json:"message"`
}

type ClientInfo struct {
	SessionId uint64
	RemoteIp  string
	conn      net.Conn
}

type TcpServer struct {
	Port          int
	listener      net.Listener
	conns         map[uint64]*ClientInfo // key - sessionid
	clientMsgs    chan *ClientPack
	SessionIdUsed atomic.Uint64
}

var tcpSvr = &TcpServer{}

func GetTcpSvr() *TcpServer {
	return tcpSvr
}

// 服务器的主动推送
func (ts *TcpServer) Broadcast(msg []byte) {
	log.Debug("pushing notify to all clients, msg: %s, receive clients: %d", msg, len(ts.conns))
	b, err := json.Marshal(&ClientPack{
		SessionId: 0,
		Msg:       msg,
	})
	if err != nil {
		log.Error("broadcast marshal error: %s", err.Error())
		return
	}
	for sId, conn := range ts.conns {
		_, err := conn.conn.Write(b)
		if err != nil {
			log.Error("broadcast to session id %d error: %s", sId, err.Error())
			delete(ts.conns, sId)
			return
		}
	}
}

func (ts *TcpServer) PushSessionResponse(sessionId uint64, msg []byte) {
	cInfo, ok := ts.conns[sessionId]
	if !ok {
		log.Error("Session Id %d not found or removed", sessionId)
		return
	}

	b, err := json.Marshal(&ClientPack{
		SessionId: sessionId,
		Msg:       msg,
	})
	if err != nil {
		log.Error("Session Id %d marshal error: %s", sessionId, err.Error())
		return
	}

	_, err = cInfo.conn.Write(b) // 这里可能已经不可写
	if err != nil {
		log.Error("Session Id %d write error: %s", sessionId, err.Error())
		return
	}
}

func (ts *TcpServer) GetMsgChan() <-chan *ClientPack {
	return ts.clientMsgs
}

func (ts *TcpServer) handleConnection(conn net.Conn) {
	connCloseAcq := true
	sessionId := uint64(0)
	defer func() {
		if connCloseAcq {
			err := conn.Close()
			if err != nil {
				log.Error("handleConnection close client error: %s", err.Error())
			}
			delete(ts.conns, sessionId)
		}
	}()
	readbytes := make([]byte, 1024)
	for {
		readSize, err := conn.Read(readbytes)
		if err != nil {
			if isEof(err) {
				// client主动发起的关闭
				log.Info("client actively closed")
				return
			}
			if isNetClosedErr(err) {
				// 服务器在其他地方主动关闭了connection（比如OnClose）导致Read阻塞解除并返回error
				log.Error("connection already closed")
				connCloseAcq = false
				return
			}
			log.Error("client read error: " + err.Error())
			return
		}
		js := readbytes[:readSize]
		addr := conn.RemoteAddr()
		// 当然，这里还得判断一下是不是login data。这里先从略
		ld := &ClientPack{}
		err = json.Unmarshal(js, ld)
		if err != nil {
			log.Error("unmarshal login data error: %s", err.Error())
			continue
		}
		if string(ld.Msg) == "LoginReq" {
			ts.SessionIdUsed.Add(1)
			newSessionId := ts.SessionIdUsed.Load()
			sessionId = newSessionId
			newC := &ClientInfo{
				SessionId: newSessionId,
				RemoteIp:  addr.String(),
				conn:      conn,
			}
			ts.conns[newSessionId] = newC
			log.Info("remote login %s success, session id %d", addr.String(), newSessionId)
			ts.PushSessionResponse(newSessionId, []byte("LoginResp"))
			continue
		}
		if u, e := ts.conns[sessionId]; e {
			ts.clientMsgs <- &ClientPack{
				SessionId: u.SessionId,
				Msg:       js,
			}
		} else {
			log.Error("illegal connection source: %s", addr.String())
			return
		}
	}
}

func (ts *TcpServer) Create(port int) {
	ts.Port = port
	log.Error("creating tcp server at port %d...", ts.Port)
	var err error

	ts.listener, err = net.Listen("tcp", "0.0.0.0:"+strconv.Itoa(ts.Port))
	if err != nil {
		// handle error
		log.Error("Error: " + err.Error())
		return
	}

	if ts.listener == nil {
		log.Error("Error: socket is nil")
		return
	}

	ts.conns = make(map[uint64]*ClientInfo)
	ts.clientMsgs = make(chan *ClientPack, 1000)
}

func (ts *TcpServer) OnClose() {
	defer close(ts.clientMsgs)
	for sessId, conn := range ts.conns {
		err := conn.conn.Close()
		if err != nil {
			if !isNetClosedErr(err) {
				log.Error("connection close to %s error: %s", sessId, err.Error())
			}
			continue
		}
		log.Info("connection close to %s success", sessId)
		delete(ts.conns, sessId)
	}
	err := ts.listener.Close()
	if err != nil {
		log.Error("TcpServer close error: %s\n", err.Error())
		return
	}
	log.Info("tcp server closed")
}

func (ts *TcpServer) OnLoop() {
	for {
		conn, err := ts.listener.Accept()
		if err != nil {
			if isNetClosedErr(err) {
				// 因为listener被关闭，中断了Accept过程（已经停止服务）
				log.Info("tcp server Accept terminated because listener is closed")
				return
			}
			// 其他错误
			log.Error("tcp server Accept error: %s\n", err.Error())
			continue
		}
		go ts.handleConnection(conn)
	}
}

func isNetClosedErr(err error) bool {
	var netErr *net.OpError
	return errors.As(err, &netErr) && errors.Is(netErr.Err, net.ErrClosed)
}

func isEof(err error) bool {
	return errors.Is(err, io.EOF)
}
