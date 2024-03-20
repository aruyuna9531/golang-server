// 一个测试TCP服务的Demo，go run client.go
// 可以拿这个去测试其他语言写的TCP服务器的连接，端口改成那个服务器开放的端口就行
package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"
)

type LoginData struct {
	UserId uint64 `json:"user_id"` // 暂且先上传这个做例子——真实登录肯定不是这么传
}

func main() {
	conn, err := net.Dial("tcp", "127.0.0.1:9000")
	if err != nil {
		fmt.Println("client dial err: " + err.Error())
		return
	}

	loginData := &LoginData{
		UserId: 1,
	}

	b, err := json.Marshal(loginData)
	if err != nil {
		log.Printf("marshal error: %v\n", err)
		return
	}

	_, err = conn.Write(b)
	if err != nil {
		log.Printf("error: %v", err)
		return
	}

	time.Sleep(1 * time.Second)
	conn.Write([]byte("hello, i am client"))

	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM)
	gorExit := make(chan struct{})
	go func() {
		buf := make([]byte, 1024)
		for {
			n, err := conn.Read(buf)
			if err != nil {
				if errors.Is(err, io.EOF) {
					log.Println("server closed")
					gorExit <- struct{}{}
					return
				}
				var netErr *net.OpError
				if errors.As(err, &netErr) && errors.Is(netErr.Err, net.ErrClosed) {
					log.Println("connection closed")
					gorExit <- struct{}{}
					return
				}
				log.Printf("read error: %s", err.Error())
				gorExit <- struct{}{}
				return
			}
			fmt.Println("client dial return: " + string(buf[:n]))
		}
	}()
	sig := <-sc
	log.Printf("client terminated by signal %v, exit", sig)
	conn.Close()
	<-gorExit
	// 断线重连先不做了
}
