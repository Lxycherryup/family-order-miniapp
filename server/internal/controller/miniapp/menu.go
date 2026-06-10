package miniapp

import (
	"family-order/server/internal/consts"
	menulogic "family-order/server/internal/logic/menu"
	"family-order/server/internal/middleware"

	"github.com/gogf/gf/v2/net/ghttp"
)

// Menu 小程序菜单控制器。
type Menu struct{}

// NewMenu 创建小程序菜单控制器。
func NewMenu() *Menu {
	return &Menu{}
}

// Categories 查询启用分类列表。
func (c *Menu) Categories(r *ghttp.Request) {
	list, err := menulogic.ListEnabledCategories(r.Context())
	if err != nil {
		middleware.WriteError(r, consts.CodeSystemError, err.Error())
		return
	}
	middleware.WriteSuccess(r, list)
}

// Dishes 查询上架菜品列表。
func (c *Menu) Dishes(r *ghttp.Request) {
	list, err := menulogic.ListOnlineDishes(r.Context(), r.Get("category_id").Int64())
	if err != nil {
		middleware.WriteError(r, consts.CodeSystemError, err.Error())
		return
	}
	middleware.WriteSuccess(r, list)
}

// DishDetail 查询上架菜品详情。
func (c *Menu) DishDetail(r *ghttp.Request) {
	dish, err := menulogic.GetOnlineDish(r.Context(), r.Get("id").Int64())
	if err != nil {
		middleware.WriteError(r, consts.CodeNotFound, err.Error())
		return
	}
	middleware.WriteSuccess(r, dish)
}
