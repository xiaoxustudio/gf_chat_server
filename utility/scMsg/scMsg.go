package scmsg

import (
	"gf_server/internal/consts"

	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
)

type SCMsgO struct {
	Type    consts.WsCode `json:"type"`
	Message string        `json:"message"`
	Data    any           `json:"data"`
}
type SCMsgD struct {
	Data g.Map `json:"data"`
}

func New(d g.Map) *SCMsgD {
	return &SCMsgD{Data: d}
}

func (r *SCMsgD) ConvertToSCMsg() (SCMsgO, error) {
	var msg SCMsgO
	if typeVal, ok := r.Data["Type"].(consts.WsCode); ok {
		msg.Type = typeVal
	} else {
		return msg, gerror.New("Type field is missing or not a string")
	}

	if messageVal, ok := r.Data["Message"].(string); ok {
		msg.Message = messageVal
	} else {
		return msg, gerror.New("Message field is missing or not a string")
	}

	// 直接赋值Data字段
	msg.Data = r.Data

	return msg, nil
}
