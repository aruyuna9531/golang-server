package main

import (
	"fmt"

	"main/myhttp"
	"main/mytcp"
)

func main() {
	fmt.Println("Hello!")

	mytcp.CreateTcpServer()
	myhttp.CreateHttpServer()
}