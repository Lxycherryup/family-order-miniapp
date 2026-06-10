package miniapp

import (
	miniappapi "family-order/server/api/miniapp"
	"family-order/server/internal/consts"
	orderlogic "family-order/server/internal/logic/order"
	"family-order/server/internal/middleware"

	"github.com/gogf/gf/v2/net/ghttp"
)

// Order 小程序订单控制器。
type Order struct{}

// NewOrder 创建小程序订单控制器。
func NewOrder() *Order {
	return &Order{}
}

// Create 创建订单。
func (c *Order) Create(r *ghttp.Request) {
	var req miniappapi.CreateOrderReq
	if err := r.Parse(&req); err != nil {
		middleware.WriteError(r, consts.CodeInvalidParams, "下单参数错误")
		return
	}

	detail, err := orderlogic.CreateOrder(r.Context(), orderlogic.CreateOrderInput{
		UserID: r.GetCtxVar(middleware.CtxKeySubjectID).Int64(),
		Items:  toCreateOrderItems(req.Items),
		Remark: req.Remark,
	})
	if err != nil {
		middleware.WriteError(r, mapOrderErrorCode(err), err.Error())
		return
	}
	middleware.WriteSuccess(r, detail)
}

// List 查询我的订单列表。
func (c *Order) List(r *ghttp.Request) {
	list, err := orderlogic.ListUserOrders(r.Context(), r.GetCtxVar(middleware.CtxKeySubjectID).Int64())
	if err != nil {
		middleware.WriteError(r, mapOrderErrorCode(err), err.Error())
		return
	}
	middleware.WriteSuccess(r, list)
}

// Detail 查询我的订单详情。
func (c *Order) Detail(r *ghttp.Request) {
	detail, err := orderlogic.GetOrderDetail(
		r.Context(),
		r.Get("id").Int64(),
		r.GetCtxVar(middleware.CtxKeySubjectID).Int64(),
	)
	if err != nil {
		middleware.WriteError(r, mapOrderErrorCode(err), err.Error())
		return
	}
	middleware.WriteSuccess(r, detail)
}

// Cancel 取消我的订单。
func (c *Order) Cancel(r *ghttp.Request) {
	err := orderlogic.CancelUserOrder(
		r.Context(),
		r.Get("id").Int64(),
		r.GetCtxVar(middleware.CtxKeySubjectID).Int64(),
	)
	if err != nil {
		middleware.WriteError(r, mapOrderErrorCode(err), err.Error())
		return
	}
	middleware.WriteSuccess(r, true)
}

// toCreateOrderItems 转换下单明细请求。
func toCreateOrderItems(items []miniappapi.CreateOrderItemReq) []orderlogic.CreateOrderItemInput {
	result := make([]orderlogic.CreateOrderItemInput, 0, len(items))
	for _, item := range items {
		result = append(result, orderlogic.CreateOrderItemInput{
			DishID:   item.DishID,
			Quantity: item.Quantity,
		})
	}
	return result
}
