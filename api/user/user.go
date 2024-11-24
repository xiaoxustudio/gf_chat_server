package user

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
)

type UserReq struct {
	g.Meta   `method:"post"`
	Response *ghttp.Response
}
type UserRes struct {
	Msg  string `json:"msg"`
	Code int    `json:"code"`
}
type UserResultReq struct {
	g.Meta `method:"post"`
}
type UserResultRes struct {
	Token string
}
