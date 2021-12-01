package main

import (
	"fmt"

	"go_svr/myhttp"
	"go_svr/mytcp"
	"go_svr/zookeeper"
	"go_svr/db"
	"io/ioutil"
	"encoding/json"
	"os"
	"os/signal"
	
)

var ServerConf struct {
	ServerId	int			`json:"server_id"`
	TcpPort		int			`json:"tcp_port"`
	HttpPort	int			`json:"http_port"`
	ZkAddr		[]string	`json:"zookeeper_addr"`
}

func main() {
	fmt.Println("Hello!")

	// 注册os.Signal
	os_channel := make(chan os.Signal, 1)
	signal.Notify(os_channel, os.Interrupt, os.Kill)

	b, err := ioutil.ReadFile("conf/server.conf") // just pass the file name
    if err != nil {
        fmt.Print(err)
		return
    }
	if err = json.Unmarshal(b, &ServerConf); err != nil{
        fmt.Print(err)
		return
	}

	go mytcp.CreateTcpServer(ServerConf.TcpPort)
	go myhttp.CreateHttpServer(ServerConf.HttpPort)
	go zookeeper.LinkZookeeper(&zookeeper.ZkConf{
		Addr:ServerConf.ZkAddr,
		ServerId:ServerConf.ServerId, 
		HttpListen:ServerConf.HttpPort,
		TcpListen:ServerConf.TcpPort,
	})

	conf := db.MysqlConf{
		Username: "root",
		Password: "mysql",
		RemoteIp: "localhost",
		RemotePort: 3306,
		DbName: "test_db",
		OpenConns: 1,
		IdleConns: 1,
	}
	db.LinkMysql(conf)

	OVERBREAK:
	for {
		signal := <-os_channel
		fmt.Printf("receive Signal %v\n", signal)
		if signal == os.Interrupt || signal == os.Kill {
			fmt.Println("Signal is interrupt or kill, exit")
			//善后工作
			zookeeper.DeleteSelfZnode()
			break OVERBREAK
		}
	}
}