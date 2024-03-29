// 一个测试TCP服务的Demo，go run client.go
// 可以拿这个去测试其他语言写的TCP服务器的连接，端口改成那个服务器开放的端口就行
package main

import (
	"encoding/json"
	"errors"
	"flag"
	"go_svr/log"
	"go_svr/proto_codes/rpc"
	"google.golang.org/protobuf/proto"
	"io"
	"math/rand"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"
)

type Msg struct {
	SessionId   uint64 `json:"session_id"`
	OpenId      string `json:"open_id"`
	MessageId   int32  `json:"message_id"`
	MessageBody []byte `json:"message_body"`
}

type client struct {
	SessionId uint64
	UserId    uint64
	OpenId    string
	Alive     bool
	conn      net.Conn
}

var cl = &client{SessionId: 0}

func (c *client) Send(id int32, msg proto.Message, forceSend bool) {
	// 这里id类型定义为int32——因为服务器不能假设客户端传过来的值一定是从枚举来，实践中往往会使用通用数值类型
	if c.SessionId == 0 && !forceSend {
		log.Error("client not ready login")
		return
	}
	//m := &Msg{
	//	SessionId: c.SessionId,
	//	Message:   utils.CopySlice(msg),
	//}
	b, err := proto.Marshal(msg)
	if err != nil {
		log.Error("send marshal error: %s", err.Error())
		return
	}
	m := &Msg{
		SessionId:   c.SessionId,
		OpenId:      c.OpenId,
		MessageId:   id,
		MessageBody: b,
	}
	bb, err := json.Marshal(m)
	if err != nil {
		log.Error("send marshal error: %s", err.Error())
		return
	}
	_, err = c.conn.Write(bb)
	if err != nil {
		log.Error("send error: %s", err.Error())
		return
	}
}

var open = flag.String("openid", "default", "定义openid")

func main() {
	flag.Parse()
	conn, err := net.Dial("tcp", "127.0.0.1:9000")
	if err != nil {
		log.Fatal("client dial err: " + err.Error())
	}
	cl.OpenId = *open
	cl.conn = conn
	log.Debug("open id = %s", cl.OpenId)

	cl.Send(int32(rpc.MessageId_Msg_CS_Login), &rpc.CS_Login{}, true)

	hbSig := make(chan struct{}, 1)
	go func() {
		// 每3秒给服务器发一个长为10的随机字符串
		for {
			select {
			case <-hbSig:
				return
			default:
				b := make([]byte, 10)
				for i := 0; i < 9; i++ {
					b[i] = byte(rand.Int()%10) + 'a'
				}
				cl.Send(int32(rpc.MessageId_Msg_CS_Message), &rpc.CS_Message{Message: string(b)}, false)
				time.Sleep(3 * time.Second)
			}
		}
	}()

	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM)
	gorExit := make(chan struct{})
	go func() {
		buf := make([]byte, 1024)
		for {
			n, err := conn.Read(buf)
			if err != nil {
				if errors.Is(err, io.EOF) {
					log.Warn("server closed")
					gorExit <- struct{}{}
					return
				}
				var netErr *net.OpError
				if errors.As(err, &netErr) && errors.Is(netErr.Err, net.ErrClosed) {
					log.Info("connection closed")
					gorExit <- struct{}{}
					return
				}
				log.Error("read error: %s", err.Error())
				gorExit <- struct{}{}
				return
			}
			msg := &Msg{}
			err = json.Unmarshal(buf[:n], msg)
			if err != nil {
				log.Error("unmarshal error: %s", err.Error())
				continue
			}
			// ↓懒得做解析器了
			switch rpc.MessageId(msg.MessageId) {
			case rpc.MessageId_Msg_SC_Login:
				loginP := &rpc.SC_Login{}
				err = proto.Unmarshal(msg.MessageBody, loginP)
				if err != nil {
					log.Error("unmarshal error: %s", err.Error())
					continue
				}
				cl.SessionId = loginP.SessionId
				cl.UserId = loginP.UserId
				cl.Alive = true
				log.Info("Login OK, session id: %d, user id %d", cl.SessionId, loginP.UserId)
			case rpc.MessageId_Msg_SC_Message:
				mp := &rpc.SC_Message{}
				err = proto.Unmarshal(msg.MessageBody, mp)
				if err != nil {
					log.Error("unmarshal error: %s", err.Error())
					continue
				}
				log.Info("receive message from server: %s", mp.Message)
			}
		}
	}()
	sig := <-sc
	hbSig <- struct{}{}
	log.Info("client terminated by signal %v, exit", sig)
	conn.Close()
	<-gorExit
	// 断线重连先不做了
}
