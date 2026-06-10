package miniapp

import "github.com/gogf/gf/v2/frame/g"

// CreateOrderReq 小程序下单请求。
type CreateOrderReq struct {
	g.Meta `path:"/api/miniapp/orders" method:"post" tags:"小程序订单" summary:"创建订单"`
	Items  []CreateOrderItemReq `json:"items"`
	Remark string               `json:"remark"`
}

// CreateOrderItemReq 小程序下单明细请求。
type CreateOrderItemReq struct {
	DishID   int64 `json:"dish_id"`
	Quantity int   `json:"quantity"`
}

// OrderListReq 小程序订单列表请求。
type OrderListReq struct {
	g.Meta `path:"/api/miniapp/orders" method:"get" tags:"小程序订单" summary:"我的订单"`
}

// OrderDetailReq 小程序订单详情请求。
type OrderDetailReq struct {
	g.Meta `path:"/api/miniapp/orders/{id}" method:"get" tags:"小程序订单" summary:"订单详情"`
	ID     int64 `json:"id"`
}

// CancelOrderReq 小程序取消订单请求。
type CancelOrderReq struct {
	g.Meta `path:"/api/miniapp/orders/{id}/cancel" method:"post" tags:"小程序订单" summary:"取消订单"`
	ID     int64 `json:"id"`
}
