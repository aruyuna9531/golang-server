package modules

import (
	"go_svr/define"
	"go_svr/dependency"
	"go_svr/game/user"
	"go_svr/log"
	"go_svr/proto_codes/rpc"
)

func ReqLogin(clientInfo *define.ClientInfo, msg *rpc.CS_Login) error {
	log.Debug("ReqLogin called, open id %s", clientInfo.OpenId)
	u := user.GetMgr().GetUserByOpen(clientInfo.OpenId)
	if u == nil {
		u = user.GetMgr().CreateUser(clientInfo.OpenId)
	}
	clientInfo.UserId = u.UserId
	return dependency.SendToClient(clientInfo, rpc.MessageId_Msg_SC_Login, &rpc.SC_Login{
		ErrCode:   0,
		SessionId: clientInfo.SessionId,
		OpenId:    clientInfo.OpenId,
		UserId:    clientInfo.UserId,
	})
}

func ReqMessage(clientInfo *define.ClientInfo, msg *rpc.CS_Message) error {
	log.Debug("ReqMessage called, userid = %d, message = %s", clientInfo.UserId, msg.Message)
	// todo 业务代码
	return nil
}
