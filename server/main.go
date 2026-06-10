package main

import (
	"family-order/server/internal/cmd"

	"github.com/gogf/gf/v2/os/gctx"
)

// main 后端服务启动入口。
func main() {
	cmd.Main.Run(gctx.GetInitCtx())
}
