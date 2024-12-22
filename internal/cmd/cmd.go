package cmd

import (
	"context"
	"gf_chat_server/internal/controller/home"
	"gf_chat_server/internal/controller/user"
	websocketunit "gf_chat_server/internal/controller/websocketUnit"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
	"github.com/gogf/gf/v2/os/gcmd"
)

var (
	Main = gcmd.Command{
		Name:  "main",
		Usage: "main",
		Brief: "start http server",
		Func: func(ctx context.Context, parser *gcmd.Parser) (err error) {
			s := g.Server()
			imc := websocketunit.New()
			s.BindHandler("/imc", func(r *ghttp.Request) {
				imc.HandleWebSocket(r.Response.Writer, r.Request)
			})
			s.Group("/user", func(group *ghttp.RouterGroup) {
				group.Middleware(ghttp.MiddlewareHandlerResponse)
				group.Bind(user.New())
			})
			s.Group("/home", func(group *ghttp.RouterGroup) {
				group.Middleware(ghttp.MiddlewareHandlerResponse)
				group.Bind(home.New())
			})
			s.SetServerRoot("/resource/public")
			s.Run()
			return nil
		},
	}
)
