package main

import (
	"fmt"

	"main/myhttp"
	"main/mytcp"
	"main/zookeeper"
)

func main() {
	fmt.Println("Hello!")

	go mytcp.CreateTcpServer(9000)
	go myhttp.CreateHttpServer(9001)
	go zookeeper.LinkZookeeper()

	for {

	}
}