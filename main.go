package main

import (
	"fmt"
	"go_svr/log"
	"go_svr/mytcp"
	"go_svr/panic_recover"
	"go_svr/timer"
	"syscall"
	"time"

	"encoding/json"
	"io/ioutil"
	"os"
	"os/signal"
)

var ServerConf struct {
	ServerId int `json:"server_id"`
	TcpPort  int `json:"tcp_port"`
	HttpPort int `json:"http_port"`
	//ZkAddr   []string `json:"zookeeper_addr"`
}

var (
	StopChan chan bool
)

func main() {
	defer panic_recover.PanicRecoverTrace()
	osChannel := make(chan os.Signal, 1)
	signal.Notify(osChannel, syscall.SIGINT, syscall.SIGKILL, syscall.SIGTERM)
	// SIGINT: kill -2 ctrl+C属于此列。
	// SIGKILL: kill -9 没有遗言的强杀（捕捉不到的信号，可能进程没了还会遗留一些现象，比如打点计时器还在跑）。不要乱用。Goland的停止按钮疑似SIGKILL（debug没抓到）
	// SIGTERM: kill -15 有遗言的退出

	b, err := ioutil.ReadFile("conf/server.conf") // just pass the file name
	if err != nil {
		fmt.Print(err)
		return
	}
	if err = json.Unmarshal(b, &ServerConf); err != nil {
		fmt.Print(err)
		return
	}

	//conf := db.MysqlConf{
	//	Username:   "root",
	//	Password:   "123456",
	//	RemoteIp:   "localhost",
	//	RemotePort: 3306,
	//	DbName:     "test",
	//}
	//db.GetDbPool().InitMysqlPool(conf)
	//defer db.GetDbPool().OnClose()

	mytcp.GetTcpSvr().Create(ServerConf.TcpPort)
	defer mytcp.GetTcpSvr().OnClose()

	go mytcp.GetTcpSvr().OnLoop()

	//go myhttp.CreateHttpServer(ServerConf.HttpPort)
	timer.TimerTestCode()

	tk := time.NewTicker(1 * time.Second)
	defer tk.Stop()
	startTime := time.Now()
	log.Printf("server start at %s\n", startTime.Format("2006-01-02 15:04:05"))
	for {
		select {
		case t := <-tk.C:
			//fmt.Printf("now: %s\n", t.Format("2006-01-02 15:04:05"))
			timer.GetInst().Trigger(t.Format("2006-01-02 15:04:05"))
		case msg := <-mytcp.GetTcpSvr().GetMsgChan():
			log.Printf("get msg from user %d, message: %s", msg.UserId, string(msg.Msg))
		case s := <-osChannel:
			fmt.Printf("receive signal %v, exit\n", s)
			return
		}
	}
}
