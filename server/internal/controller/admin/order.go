package admin

import (
	adminapi "family-order/server/api/admin"
	"family-order/server/internal/consts"
	orderlogic "family-order/server/internal/logic/order"
	"family-order/server/internal/middleware"

	"github.com/gogf/gf/v2/net/ghttp"
)

// Order 管理端订单控制器。
type Order struct{}

// NewOrder 创建管理端订单控制器。
func NewOrder() *Order {
	return &Order{}
}

// List 查询订单列表。
func (c *Order) List(r *ghttp.Request) {
	list, err := orderlogic.ListAdminOrders(r.Context(), orderlogic.AdminOrderListInput{
		Status:   r.Get("status").Int(),
		Page:     r.Get("page").Int(),
		PageSize: r.Get("page_size").Int(),
	})
	if err != nil {
		middleware.WriteError(r, consts.CodeSystemError, err.Error())
		return
	}
	middleware.WriteSuccess(r, list)
}

// Detail 查询订单详情。
func (c *Order) Detail(r *ghttp.Request) {
	detail, err := orderlogic.GetOrderDetail(r.Context(), r.Get("id").Int64(), 0)
	if err != nil {
		middleware.WriteError(r, mapAdminOrderErrorCode(err), err.Error())
		return
	}
	middleware.WriteSuccess(r, detail)
}

// UpdateStatus 更新订单状态。
func (c *Order) UpdateStatus(r *ghttp.Request) {
	var req adminapi.OrderStatusReq
	if err := r.Parse(&req); err != nil {
		middleware.WriteError(r, consts.CodeInvalidParams, "订单状态参数错误")
		return
	}

	if err := orderlogic.ChangeOrderStatus(r.Context(), r.Get("id").Int64(), req.Status); err != nil {
		middleware.WriteError(r, mapAdminOrderErrorCode(err), err.Error())
		return
	}
	middleware.WriteSuccess(r, true)
}
