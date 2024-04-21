package main

import (
	"flag"
	"fmt"
	"go_svr/log"
	"go_svr/proto_codes/rpc"
	"go_svr/share/rudp"
	"google.golang.org/protobuf/proto"
	"math/rand"
	"net"
	"os"
	"os/signal"
	"sync/atomic"
	"syscall"
	"time"
)

var seat = flag.Int64("seat", 0, "")

// UDP 客户端
func main() {
	var a atomic.Int32
	conn, err := net.DialUDP("udp", nil, &net.UDPAddr{IP: net.IPv4zero, Port: 8888})
	if err != nil {
		fmt.Println(err)
		return
	}
	rconn := rudp.NewConn(conn, rudp.New())
	for !rconn.Connected() {
	}
	go func() {
		for {
			time.Sleep(1*time.Second + time.Duration(rand.Int()%1000)*time.Millisecond) // 随机睡1到2秒
			r := &rpc.InputData{
				Id:         uint64(a.Add(1)),
				Sid:        rand.Int31() % 5,
				X:          rand.Int31() % 100,
				Y:          rand.Int31() % 100,
				Roomseatid: int32(*seat),
			}
			b, err := proto.Marshal(r)
			if err != nil {
				log.Error("op error: %s", err.Error())
				continue
			}
			rconn.Write(b)
			log.Debug("write input to server, id %d, sid %d, x %d, y %d, seat %d", r.Id, r.Sid, r.X, r.Y, r.Roomseatid)
		}
	}()
	go func() {
		var b [1024]byte
		for {
			n, err := rconn.Read(b[:])
			if err != nil {
				log.Error("eof")
				return
			}
			if n == 0 {
				continue
			}
			f := &rpc.FrameData{}
			err = proto.Unmarshal(b[:n], f)
			if err != nil {
				continue
			}
			if len(f.Input) > 0 {
				log.Debug("Receive Frame data: frame id %d, input: \n", f.FrameID)
				for _, ip := range f.Input {
					log.Debug("id %d, sid %d, X %d, Y %d, seatid %d", ip.Id, ip.Sid, ip.X, ip.Y, ip.Roomseatid)
				}
			}
		}
	}()
	sc := make(chan os.Signal)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM)
	<-sc
}
