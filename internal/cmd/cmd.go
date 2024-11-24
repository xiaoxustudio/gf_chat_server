package cmd

import (
	"context"
	"gf_server/internal/controller/home"
	"gf_server/internal/controller/user"
	websocketunit "gf_server/internal/controller/websocketUnit"
	"net/http"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
	"github.com/gogf/gf/v2/os/gcmd"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

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
