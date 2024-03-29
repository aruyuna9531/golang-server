package dependency

import (
	"go_svr/define"
	"go_svr/proto_codes/rpc"
	"google.golang.org/protobuf/proto"
)

var SendToClient func(info *define.ClientInfo, messageId rpc.MessageId, message proto.Message) error
