// 一个测试TCP服务的Demo，go run client.go
// 可以拿这个去测试其他语言写的TCP服务器的连接，端口改成那个服务器开放的端口就行
package main

import (
	"encoding/json"
	"errors"
	"go_svr/log"
	"go_svr/utils"
	"io"
	"math/rand"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"
)

type Msg struct {
	SessionId uint64 `json:"session_id"`
	Message   []byte `json:"message"`
}

type client struct {
	SessionId uint64
	Alive     bool
	conn      net.Conn
}

var cl = &client{SessionId: 0}

func (c *client) Send(msg []byte) {
	if c.SessionId == 0 {
		log.Error("client not ready login")
		return
	}
	m := &Msg{
		SessionId: c.SessionId,
		Message:   utils.CopySlice(msg),
	}
	b, err := json.Marshal(m)
	if err != nil {
		log.Error("send marshal error: %s", err.Error())
		return
	}
	_, err = c.conn.Write(b)
	if err != nil {
		log.Error("send error: %s", err.Error())
		return
	}
}

func main() {
	conn, err := net.Dial("tcp", "127.0.0.1:9000")
	if err != nil {
		log.Fatal("client dial err: " + err.Error())
	}
	cl.conn = conn

	loginData := &Msg{
		SessionId: 0, // 登录的时候还没有这个
		Message:   []byte("LoginReq"),
	}

	b, err := json.Marshal(loginData)
	if err != nil {
		log.Error("marshal error: %v\n", err)
		return
	}

	_, err = conn.Write(b)
	if err != nil {
		log.Error("error: %v", err)
		return
	}

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
				cl.Send(b)
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
			if string(msg.Message) == "LoginResp" {
				cl.SessionId = msg.SessionId
				cl.Alive = true
				log.Info("Login OK, session id: %d", cl.SessionId)
			} else {
				if msg.SessionId == 0 {
					log.Debug("Server pushed broadcast message: %s", msg.Message)
				} else if msg.SessionId == cl.SessionId {
					log.Debug("Server response message: %s", msg.Message)
				} else {
					log.Error("Server response session id illegal: %d (self %d)", msg.SessionId, cl.SessionId)
				}
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
