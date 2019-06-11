package serve

import (
	"api"

	"wholeally.com/common/logs"
	"wholeally.com/share/v4/msqclient"
	"wholeally.com/share/v4/protocol"
)

// 消息通知
func msqNotify(msg *protocol.Message) {
	// 回复消息不处理(设置平台在线线状态的回复消息)
	if protocol.MT_RESPONSE == msg.GetType() {
		return
	}

	go func() {
		ctx := protocol.Context{
			SetSession: msg.GetMessageTo(),
			MsgFrom:    msg.GetMessageFrom(),
		}
		// 调用注册函数
		rmsg := api.Router.Call(msg, &ctx)
		// 输出回应消息
		if nil != rmsg {
			for i := 0; i < 3; i++ {
				err := msqclient.Post(rmsg)
				if nil != err {
					logs.Waring(err)
				} else {
					return
				}
			}
		}
	}()
}
