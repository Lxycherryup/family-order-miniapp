package admin

import "github.com/gogf/gf/v2/frame/g"

// CategoryListReq 分类列表请求。
type CategoryListReq struct {
	g.Meta `path:"/api/admin/categories" method:"get" tags:"管理端分类" summary:"分类列表"`
}

// CategorySaveReq 分类保存请求。
type CategorySaveReq struct {
	Name string `json:"name"`
	Sort int    `json:"sort"`
}

// CategoryStatusReq 分类状态请求。
type CategoryStatusReq struct {
	Status int `json:"status"`
}
