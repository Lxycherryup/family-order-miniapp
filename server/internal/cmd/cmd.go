package cmd

import (
	"context"

	"family-order/server/internal/controller"
	admincontroller "family-order/server/internal/controller/admin"
	miniappcontroller "family-order/server/internal/controller/miniapp"
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
			group.Group("/api/admin", func(adminGroup *ghttp.RouterGroup) {
				adminGroup.Middleware(middleware.AdminAuth)

				category := admincontroller.NewCategory()
				adminGroup.GET("/categories", category.List)
				adminGroup.POST("/categories", category.Create)
				adminGroup.PUT("/categories/{id}", category.Update)
				adminGroup.DELETE("/categories/{id}", category.Delete)
				adminGroup.PUT("/categories/{id}/status", category.UpdateStatus)

				dish := admincontroller.NewDish()
				adminGroup.GET("/dishes", dish.List)
				adminGroup.POST("/dishes", dish.Create)
				adminGroup.GET("/dishes/{id}", dish.Detail)
				adminGroup.PUT("/dishes/{id}", dish.Update)
				adminGroup.DELETE("/dishes/{id}", dish.Delete)
				adminGroup.PUT("/dishes/{id}/status", dish.UpdateStatus)

				order := admincontroller.NewOrder()
				adminGroup.GET("/orders", order.List)
				adminGroup.GET("/orders/{id}", order.Detail)
				adminGroup.PUT("/orders/{id}/status", order.UpdateStatus)
			})

			miniappAuth := miniappcontroller.NewAuth()
			group.Group("/api/miniapp/auth", func(authGroup *ghttp.RouterGroup) {
				authGroup.POST("/login", miniappAuth.Login)
			})
			group.Group("/api/miniapp/menu", func(menuGroup *ghttp.RouterGroup) {
				menuGroup.Middleware(middleware.MiniappAuth)

				menu := miniappcontroller.NewMenu()
				menuGroup.GET("/categories", menu.Categories)
				menuGroup.GET("/dishes", menu.Dishes)
				menuGroup.GET("/dishes/{id}", menu.DishDetail)
			})
			group.Group("/api/miniapp", func(miniappGroup *ghttp.RouterGroup) {
				miniappGroup.Middleware(middleware.MiniappAuth)

				order := miniappcontroller.NewOrder()
				miniappGroup.POST("/orders", order.Create)
				miniappGroup.GET("/orders", order.List)
				miniappGroup.GET("/orders/{id}", order.Detail)
				miniappGroup.POST("/orders/{id}/cancel", order.Cancel)
			})
		})
		s.Run()
		return nil
	},
}
