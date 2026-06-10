package admin

import "github.com/gogf/gf/v2/frame/g"

// OrderListReq 管理端订单列表请求。
type OrderListReq struct {
	g.Meta   `path:"/api/admin/orders" method:"get" tags:"管理端订单" summary:"订单列表"`
	Status   int `json:"status"`
	Page     int `json:"page"`
	PageSize int `json:"page_size"`
}

// OrderDetailReq 管理端订单详情请求。
type OrderDetailReq struct {
	g.Meta `path:"/api/admin/orders/{id}" method:"get" tags:"管理端订单" summary:"订单详情"`
	ID     int64 `json:"id"`
}

// OrderStatusReq 管理端订单状态请求。
type OrderStatusReq struct {
	Status int `json:"status"`
}
