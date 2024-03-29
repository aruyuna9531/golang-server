package define

import (
	"net"
)

type ClientInfo struct {
	OpenId    string
	UserId    uint64
	SessionId uint64
	RemoteIp  string
	Conn      net.Conn
}
