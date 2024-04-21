package framebase

import (
	"fmt"
	"go_svr/define"
	"go_svr/game/modules"
	"go_svr/proto_codes/rpc"
	"google.golang.org/protobuf/proto"
)

type analyzer interface {
	Exec(clientInfo *define.ClientInfo, message []byte) error
}

var protoMaps = map[rpc.MessageId]analyzer{}

func init() {
	protoMaps[rpc.MessageId_Msg_CS_Login] = &P_CS_Login{}
	protoMaps[rpc.MessageId_Msg_CS_Message] = &P_CS_Message{}
}

type P_CS_Login struct{}

func (*P_CS_Login) Exec(clientInfo *define.ClientInfo, message []byte) error {
	p := &rpc.CS_Login{}
	err := proto.Unmarshal(message, p)
	if err != nil {
		return fmt.Errorf("cannot parse to CS_Login")
	}
	// other checkers
	return modules.ReqLogin(clientInfo, p)
}

type P_CS_Message struct{}

func (*P_CS_Message) Exec(clientInfo *define.ClientInfo, message []byte) error {
	p := &rpc.CS_Message{}
	err := proto.Unmarshal(message, p)
	if err != nil {
		return fmt.Errorf("cannot parse to CS_Message")
	}
	// other checkers
	return modules.ReqMessage(clientInfo, p)
}
