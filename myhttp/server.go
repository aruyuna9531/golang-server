package myhttp

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"go_svr/myhttp/handlers"
	"log"
	"net/http"
	"strconv"
	"strings"
	"sync/atomic"
	"time"
)

type FlowControlData struct {
	count   int32
	last    int64
	lasturl string
}

func (d *FlowControlData) Allow(url string) bool {
	if url != d.lasturl {
		d.Clear()
		return true
	}

	if d.count < handlers.FlowControlMax {
		return true
	}

	if time.Now().UnixMilli() > d.last+handlers.FlowControlTimeGap*1000 {
		return true
	}

	return false
}

func (d *FlowControlData) Clear() {
	d.count = 0
	d.last = 0
	d.lasturl = ""
}

func (d *FlowControlData) Add() {
	if time.Now().UnixMilli() > d.last+handlers.FlowControlTimeGap*int64(time.Millisecond) {
		d.Clear()
	}
	d.count++
	d.last = time.Now().UnixMilli()
}

func (d *FlowControlData) String() string {
	return fmt.Sprintf("FlowControlData: count %d, last %s", d.count, time.UnixMilli(d.last).Format("2006-01-02 15:04:05"))
}

var reqSolving atomic.Int32
var reqIpCount = map[string]*FlowControlData{}

func CreateHttpServer(port int) {
	log.Println("creating http server at port " + strconv.Itoa(port) + "...")

	//// 注册url及其执行的回调函数
	//RegHttpHandler("/mainpage", handlers.MainHandler)
	//
	//// 启动http服务
	//http.ListenAndServe("0.0.0.0:"+strconv.Itoa(port), nil)
	r := gin.Default()
	r.GET("/mainpage", func(context *gin.Context) {
		if !check(context.Writer, context.Request) {
			return
		}
		handlers.MainHandlerByGin(context)
	})
	r.GET("/AMoney", func(context *gin.Context) {
		if !check(context.Writer, context.Request) {
			return
		}
		handlers.AMoney(context)
	})
	r.POST("/ktvaacalc", func(context *gin.Context) {
		if !check(context.Writer, context.Request) {
			return
		}
		handlers.AACalcResult(context)
	})
	err := r.Run("0.0.0.0:9001")
	if err != nil {
		panic(err)
	}
}

func check(writer http.ResponseWriter, request *http.Request) bool {
	// 事前准备
	log.Printf("request url: %s, request host: %s, source ip: %s", request.URL, request.Host, request.RemoteAddr)
	// 检查代理
	proxyForward := request.Header.Get("X-Forwarded-For")
	if proxyForward != "" {
		log.Printf("request proxy detected, route = %s", proxyForward)
	}
	// 限制同时只能处理5条请求
	if reqSolving.Load() >= 5 {
		log.Printf("request error, now solving request count >= 5")
		handlers.ToClient(writer, http.StatusNotAcceptable, "现在访问量过大，请稍后再试")
		return false
	}
	reqSolving.Add(1)
	defer reqSolving.Add(-1)

	realIp := request.Header.Get("X-Real-Ip")
	if realIp == "" {
		realIp = strings.Split(request.RemoteAddr, ":")[0]
	}
	d, ok := reqIpCount[realIp]
	if !ok {
		reqIpCount[realIp] = &FlowControlData{}
		d = reqIpCount[realIp]
	}
	log.Printf("request url path=%s", request.URL.Path)
	if !d.Allow(request.URL.Path) {
		log.Printf("request error, ip %s flow control failed, %s", realIp, d.String())
		handlers.ToClient(writer, http.StatusForbidden, "访问过于频繁，请稍后再试")
		return false
	}
	d.Add()
	return true
}

func RegHttpHandler(url string, handler func(http.ResponseWriter, *http.Request)) {
	http.HandleFunc(url, func(writer http.ResponseWriter, request *http.Request) {
		if !check(writer, request) {
			return
		}
		handler(writer, request)
	})
}
