package main

import (
	"fmt"

	"main/myhttp"
	"main/mytcp"
	"main/zookeeper"
	"main/db"
)

func main() {
	fmt.Println("Hello!")

	go mytcp.CreateTcpServer(9000)
	go myhttp.CreateHttpServer(9001)
	go zookeeper.LinkZookeeper()

	conf := db.MysqlConf{
		Username: "root",
		Password: "mysql",
		RemoteIp: "192.168.2.181",
		RemotePort: 3306,
		DbName: "test_db",
		OpenConns: 1,
		IdleConns: 1,
	}
	db.LinkMysql(conf)

	for {

	}
}