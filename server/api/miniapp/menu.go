package miniapp

import "github.com/gogf/gf/v2/frame/g"

// CategoryListReq 小程序分类列表请求。
type CategoryListReq struct {
	g.Meta `path:"/api/miniapp/menu/categories" method:"get" tags:"小程序菜单" summary:"分类列表"`
}

// DishListReq 小程序菜品列表请求。
type DishListReq struct {
	g.Meta     `path:"/api/miniapp/menu/dishes" method:"get" tags:"小程序菜单" summary:"菜品列表"`
	CategoryID int64 `json:"category_id"`
}

// DishDetailReq 小程序菜品详情请求。
type DishDetailReq struct {
	g.Meta `path:"/api/miniapp/menu/dishes/{id}" method:"get" tags:"小程序菜单" summary:"菜品详情"`
	ID     int64 `json:"id"`
}
