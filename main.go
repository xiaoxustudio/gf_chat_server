package main

import (
	_ "gf_chat_server/internal/packed"

	"gf_chat_server/internal/cmd"

	_ "github.com/gogf/gf/contrib/drivers/mysql/v2"
	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/v2/os/gctx"
	"github.com/joho/godotenv"
)

func main() {
	res := godotenv.Load(".local.env") // 加载本地环境文件
	if res != nil {
		return
	}
	g.Cfg().SetFileName("./manifest/config/config.yaml")
	cmd.Main.Run(gctx.GetInitCtx())
}
