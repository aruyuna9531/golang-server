package myhttp

import (
	"fmt"
	"net/http"
	"strconv"
)



func MainHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "hello!!")
}

func CreateHttpServer(port int) {
	fmt.Println("creating http server at port " + strconv.Itoa(port) + "...")

	// 注册url及其执行的回调函数
	http.HandleFunc("/", MainHandler)

	// 启动http服务
	http.ListenAndServe("0.0.0.0:" + strconv.Itoa(port), nil)
}