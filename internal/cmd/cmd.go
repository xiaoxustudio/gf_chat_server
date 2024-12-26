/*
 * _______________#########_______________________
 * ______________############_____________________
 * ______________#############____________________
 * _____________##__###########___________________
 * ____________###__######_#####__________________
 * ____________###_#######___####_________________
 * ___________###__##########_####________________
 * __________####__###########_####_______________
 * ________#####___###########__#####_____________
 * _______######___###_########___#####___________
 * _______#####___###___########___######_________
 * ______######___###__###########___######_______
 * _____######___####_##############__######______
 * ____#######__#####################_#######_____
 * ____#######__##############################____
 * ___#######__######_#################_#######___
 * ___#######__######_######_#########___######___
 * ___#######____##__######___######_____######___
 * ___#######________######____#####_____#####____
 * ____######________#####_____#####_____####_____
 * _____#####________####______#####_____###______
 * ______#####______;###________###______#________
 * ________##_______####________####______________
 */

package cmd

import (
	"context"
	"gf_chat_server/internal/controller/group"
	"gf_chat_server/internal/controller/home"
	"gf_chat_server/internal/controller/user"
	websocketunit "gf_chat_server/internal/controller/websocketUnit"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
	"github.com/gogf/gf/v2/os/gcmd"
)

// 允许跨域请求中间件
func Middleware(r *ghttp.Request) {
	r.Response.CORSDefault()
	r.Middleware.Next()
}

var (
	Main = gcmd.Command{
		Name:  "main",
		Usage: "main",
		Brief: "start http server",
		Func: func(ctx context.Context, parser *gcmd.Parser) (err error) {
			s := g.Server()
			imc := websocketunit.New()

			s.BindHandler("/imc", func(r *ghttp.Request) { // 聊天websocket
				imc.HandleWebSocket(r.Response.Writer, r.Request)
			})
			s.Group("/user", func(group *ghttp.RouterGroup) { // 用户相关接口
				group.Middleware(ghttp.MiddlewareHandlerResponse)
				group.Middleware(Middleware)
				group.Bind(user.New())
			})
			s.Group("/group", func(g *ghttp.RouterGroup) { // 群组相关接口
				g.Middleware(ghttp.MiddlewareHandlerResponse)
				g.Middleware(Middleware)
				g.Bind(group.New())
			})
			s.Group("/home", func(group *ghttp.RouterGroup) {
				group.Middleware(ghttp.MiddlewareHandlerResponse)
				group.Bind(home.New())
			})
			s.SetServerRoot("/resource/public")
			s.AddStaticPath("/temp", "./temp")
			s.AddStaticPath("/resource", "./resource")
			s.Run()
			return nil
		},
	}
)
