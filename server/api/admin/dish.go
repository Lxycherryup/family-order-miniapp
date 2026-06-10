package admin

import "github.com/gogf/gf/v2/frame/g"

// DishListReq 菜品列表请求。
type DishListReq struct {
	g.Meta     `path:"/api/admin/dishes" method:"get" tags:"管理端菜品" summary:"菜品列表"`
	CategoryID int64 `json:"category_id"`
	Status     int   `json:"status"`
}

// DishSaveReq 菜品保存请求。
type DishSaveReq struct {
	CategoryID  int64   `json:"category_id"`
	Name        string  `json:"name"`
	Description string  `json:"description"`
	ImageURL    string  `json:"image_url"`
	Price       float64 `json:"price"`
	Unit        string  `json:"unit"`
	Sort        int     `json:"sort"`
}

// DishStatusReq 菜品状态请求。
type DishStatusReq struct {
	Status int `json:"status"`
}
