package framebase

import (
	"encoding/json"
	"fmt"
	"go_svr/define"
	"go_svr/dependency"
	"go_svr/log"
	"go_svr/proto_codes/rpc"
	"go_svr/utils"
	"google.golang.org/protobuf/proto"
	"net"
	"strconv"
	"sync/atomic"
)

// 假设这是客户端发过来的东西（先不用protobuf）
type ClientPack struct {
	SessionId uint64 `json:"session_id"`
	OpenId    string `json:"open_id"`
	MsgId     int32  `json:"message_id"`
	Msg       []byte `json:"message_body"`
}

func (cp *ClientPack) Exec() {
	alzer, ok := protoMaps[rpc.MessageId(cp.MsgId)]
	if !ok {
		log.Error("Exec error, message id %d not exist", cp.MsgId)
		return
	}
	cl, ok := tcpSvr.conns[cp.SessionId]
	if !ok {
		log.Error("Exec error, conn session id %d not exist", cp.SessionId)
		return
	}
	err := alzer.Exec(cl, cp.Msg)
	if err != nil {
		log.Error("exec error: %s", err.Error())
	}
}

type TcpServer struct {
	Port          int
	listener      net.Listener
	conns         map[uint64]*define.ClientInfo // key - sessionid
	open2conn     map[string]*define.ClientInfo
	clientMsgs    chan *ClientPack
	SessionIdUsed atomic.Uint64
}

var tcpSvr = &TcpServer{}

func GetTcpSvr() *TcpServer {
	return tcpSvr
}

func init() {
	dependency.SendToClient = tcpSvr.SendToClient
	dependency.Disconnect = tcpSvr.ForceDisconnect
}

func (ts *TcpServer) SendToClient(cInfo *define.ClientInfo, messageId rpc.MessageId, message proto.Message) error {
	b, err := proto.Marshal(message)
	if err != nil {
		return err
	}

	b2, err := json.Marshal(&ClientPack{
		SessionId: cInfo.SessionId,
		OpenId:    cInfo.OpenId,
		MsgId:     int32(messageId),
		Msg:       utils.CopySlice(b),
	})
	if err != nil {
		log.Error("broadcast marshal error: %s", err.Error())
		return err
	}

	_, err = cInfo.Conn.Write(b2) // 这里可能已经不可写
	if err != nil {
		return err
	}
	return nil
}

// 服务器的主动推送
func (ts *TcpServer) Broadcast(msgId rpc.MessageId, msg proto.Message) {
	log.Debug("pushing notify to all clients, msg: %s, receive clients: %d", msg, len(ts.conns))
	bs, err := proto.Marshal(msg)
	if err != nil {
		log.Error("broadcast marshal error: %s", err.Error())
		return
	}

	b, err := json.Marshal(&ClientPack{
		SessionId: 0,
		OpenId:    "",
		MsgId:     int32(msgId),
		Msg:       utils.CopySlice(bs),
	})
	if err != nil {
		log.Error("broadcast marshal error: %s", err.Error())
		return
	}
	for sId, conn := range ts.conns {
		_, err := conn.Conn.Write(b)
		if err != nil {
			log.Error("broadcast to session id %d error: %s", sId, err.Error())
			delete(ts.conns, sId)
			return
		}
	}
}

func (ts *TcpServer) GetMsgChan() <-chan *ClientPack {
	return ts.clientMsgs
}

func (ts *TcpServer) handleConnection(conn net.Conn) {
	connCloseAcq := true
	sessionId := ts.SessionIdUsed.Add(1)
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
	cl := &define.ClientInfo{
		SessionId: sessionId,
		RemoteIp:  conn.RemoteAddr().String(),
		Conn:      conn,
	}
	log.Debug("new session %d established", sessionId)
	ts.conns[sessionId] = cl
	for {
		// todo ↓ 注意包大小校验
		readSize, err := conn.Read(readbytes)
		if err != nil {
			if utils.IsEof(err) {
				// client主动发起的关闭
				log.Info("client actively closed")
				return
			}
			if utils.IsNetClosedErr(err) {
				// 服务器在其他地方主动关闭了connection（比如OnClose）导致Read阻塞解除并返回error
				log.Error("connection already closed")
				connCloseAcq = false
				return
			}
			log.Error("client read error: " + err.Error())
			return
		}
		js := readbytes[:readSize]
		ld := &ClientPack{}
		err = json.Unmarshal(js, ld)
		if err != nil {
			log.Error("unmarshal data error: %s", err.Error())
			continue
		}
		if cl.OpenId == "" {
			cl.OpenId = ld.OpenId
			ts.open2conn[cl.OpenId] = cl
		}
		if ld.SessionId == 0 {
			ld.SessionId = sessionId
		}
		ts.clientMsgs <- ld
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

	ts.conns = make(map[uint64]*define.ClientInfo)
	ts.open2conn = make(map[string]*define.ClientInfo)
	ts.clientMsgs = make(chan *ClientPack, 1000)
}

func (ts *TcpServer) OnClose() {
	defer close(ts.clientMsgs)
	for sessId, conn := range ts.conns {
		err := conn.Conn.Close()
		if err != nil {
			if !utils.IsNetClosedErr(err) {
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
			if utils.IsNetClosedErr(err) {
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

func (ts *TcpServer) ForceDisconnect(openId string) error {
	cl, ok := ts.open2conn[openId]
	if !ok {
		return fmt.Errorf("open id %s is not connected", openId)
	}
	defer delete(ts.open2conn, openId)
	err := ts.SendToClient(cl, rpc.MessageId_Msg_SC_DisconnectNotify, &rpc.SC_DisconnectNotify{Reason: 1})
	if err != nil {
		return err
	}
	return cl.Conn.Close()
}
