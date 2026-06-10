package cmd

import (
	"context"

	"family-order/server/internal/controller"
	admincontroller "family-order/server/internal/controller/admin"
	"family-order/server/internal/middleware"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
	"github.com/gogf/gf/v2/os/gcmd"
)

// Main 后端主命令。
var Main = gcmd.Command{
	Name:  "main",
	Usage: "main",
	Brief: "启动家庭点餐后端服务",
	Func: func(ctx context.Context, parser *gcmd.Parser) error {
		s := g.Server()
		s.SetAddr(":8000")
		s.Group("/", func(group *ghttp.RouterGroup) {
			group.Middleware(middleware.RequestLog)
			health := controller.NewHealth()
			group.GET("/health", health.Handler)

			adminAuth := admincontroller.NewAuth()
			group.Group("/api/admin/auth", func(authGroup *ghttp.RouterGroup) {
				authGroup.POST("/login", adminAuth.Login)
				authGroup.Group("/", func(protectedGroup *ghttp.RouterGroup) {
					protectedGroup.Middleware(middleware.AdminAuth)
					protectedGroup.POST("/logout", adminAuth.Logout)
					protectedGroup.GET("/profile", adminAuth.Profile)
				})
			})
		})
		s.Run()
		return nil
	},
}
