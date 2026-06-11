package main

import (
	"family-order/server/internal/cmd"

	// 注册 PostgreSQL 数据库驱动。
	_ "github.com/gogf/gf/contrib/drivers/pgsql/v2"
	"github.com/gogf/gf/v2/os/gctx"
)

// main 后端服务启动入口。
func main() {
	cmd.Main.Run(gctx.GetInitCtx())
}
