package rudp

import (
	"go_svr/log"
	"net"
	"sync"
)

func NewListener(conn *net.UDPConn) *RudpListener {
	listen := &RudpListener{conn: conn,
		newRudpConn: make(chan *RudpConn, 1024),
		newRudpErr:  make(chan error, 12),
		rudpConnMap: make(map[string]*RudpConn)}
	go listen.run()
	return listen
}

type RudpListener struct {
	conn *net.UDPConn
	lock sync.RWMutex

	newRudpConn chan *RudpConn
	newRudpErr  chan error
	rudpConnMap map[string]*RudpConn
}

// net listener interface
func (this *RudpListener) Accept() (net.Conn, error) { return this.AcceptRudp() }
func (this *RudpListener) Close() error {
	this.CloseAllRudp()
	return this.conn.Close()
}
func (this *RudpListener) Addr() net.Addr { return this.conn.LocalAddr() }

func (this *RudpListener) CloseRudp(addr string) {
	this.lock.Lock()
	defer this.lock.Unlock()
	if conn, ok := this.rudpConnMap[addr]; ok {
		conn.Close()
		delete(this.rudpConnMap, addr)
		log.Debug("来自%s的rudp连接已被删除", addr)
	}
	this.lock.Unlock()
}

func (this *RudpListener) CloseAllRudp() {
	this.lock.Lock()
	for _, rconn := range this.rudpConnMap {
		rconn.Close()
		log.Debug("来自%s的rudp连接已被删除", rconn.remoteAddr.String())
	}
	this.rudpConnMap = make(map[string]*RudpConn)
	this.lock.Unlock()
}
func (this *RudpListener) AcceptRudp() (*RudpConn, error) {
	select {
	case c := <-this.newRudpConn:
		return c, nil
	case e := <-this.newRudpErr:
		return nil, e
	}
}
func (this *RudpListener) run() {
	data := make([]byte, MAX_PACKAGE)
	for {
		n, remoteAddr, err := this.conn.ReadFromUDP(data)
		if err != nil {
			this.CloseAllRudp()
			this.newRudpErr <- err
			return
		}
		this.lock.RLock()
		rudpConn, ok := this.rudpConnMap[remoteAddr.String()]
		this.lock.RUnlock()
		if !ok {
			rudpConn = NewUnConn(this.conn, remoteAddr, New(), this.CloseRudp)
			this.lock.Lock()
			this.rudpConnMap[remoteAddr.String()] = rudpConn
			this.lock.Unlock()
			this.newRudpConn <- rudpConn
			log.Debug("新的UDP连接，来自%s", remoteAddr.String())
		}
		bts := make([]byte, n)
		copy(bts, data[:n])
		rudpConn.in <- bts
	}
}
