package home

import (
	"gf_server/internal/consts"
	"gf_server/utility/msgtoken"
	"gf_server/utility/token"

	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
)

type Home struct{}

func New() *Home {
	return &Home{}
}

func (c *Home) GetHome(req *ghttp.Request) {
	tok := req.Header.Get("Authorization")
	if len(tok) == 0 {
		req.Response.WriteJsonExit(msgtoken.ToGMap(msgtoken.MsgToken(consts.TokenInValid, "token校验失败！", nil)))
	}
	_, err := token.ValidToken(tok)
	if err == nil {
		req.Response.WriteJsonExit(msgtoken.ToGMap(msgtoken.MsgToken(1, "OK", g.Map{
			"num":   1,
			"token": tok,
		})))
	} else {
		req.Response.WriteJsonExit(msgtoken.ToGMap(msgtoken.MsgToken(consts.TokenExpired, "token校验失败:"+err.Error(), nil)))
	}
}
