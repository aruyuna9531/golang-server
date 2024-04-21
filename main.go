package main

import (
	"fmt"
	"go_svr/framebase"
	"go_svr/lock_step"
	"go_svr/log"
	"go_svr/panic_recover"
	"go_svr/share/rudp"
	"go_svr/timer"
	"net/http"
	_ "net/http/pprof"
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

func pprof() {
	err := http.ListenAndServe(":9090", nil) //设置监听的端口
	if err != nil {
		fmt.Printf("ListenAndServe: %s", err)
	}
}

func main() {
	defer panic_recover.PanicRecoverTrace() // TODO 放这里没用 换个地方
	go pprof()

	osChannel := make(chan os.Signal, 1)
	signal.Notify(osChannel, syscall.SIGINT, syscall.SIGKILL, syscall.SIGTERM, syscall.SIGSEGV)
	// SIGINT: kill -2 ctrl+C属于此列。只能在非daemon模式的程序运行中使用，外部给这个进程发SIGINT不会被响应
	// SIGKILL: kill -9 没有遗言的强杀（捕捉不到的信号）。不要乱用。Goland的停止按钮疑似SIGKILL（debug没抓到）
	// SIGSEGV: kill -10 segmentation fault
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

	//redis.GetRedisCli().Init()

	framebase.GetTcpSvr().Create(ServerConf.TcpPort)
	defer framebase.GetTcpSvr().OnClose()
	go framebase.GetTcpSvr().OnLoop()

	rudp.GetInst().Init()
	lock_step.GetInst().Init()

	//go myhttp.CreateHttpServer(ServerConf.HttpPort)
	timer.TimerTestCode()

	tk := time.NewTicker(1 * time.Second)
	defer tk.Stop()
	startTime := time.Now()
	log.Debug("server start at %s\n", startTime.Format("2006-01-02 15:04:05"))
	for {
		select {
		case t := <-tk.C:
			//fmt.Printf("now: %s\n", t.Format("2006-01-02 15:04:05"))
			timer.GetInst().Trigger(t.UnixMilli())
		case msg := <-framebase.GetTcpSvr().GetMsgChan():
			msg.Exec()
		case s := <-osChannel:
			log.Info("receive signal %v, exit\n", s)
			return
		}
	}
}
