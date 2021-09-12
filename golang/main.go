package main

import (
	"fmt"
	"net/http"

	"main/myhttp"
	"main/mytcp"
)

func main() {
	fmt.Println("Hello!")

	mytcp.CreateTcpServer()

	http.HandleFunc("/", myhttp.MainHandler)
	http.ListenAndServe("0.0.0.0:8000", nil)
}