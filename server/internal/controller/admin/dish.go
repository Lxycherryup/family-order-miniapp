package admin

import (
	adminapi "family-order/server/api/admin"
	"family-order/server/internal/consts"
	menulogic "family-order/server/internal/logic/menu"
	"family-order/server/internal/middleware"

	"github.com/gogf/gf/v2/net/ghttp"
)

// Dish 管理端菜品控制器。
type Dish struct{}

// NewDish 创建管理端菜品控制器。
func NewDish() *Dish {
	return &Dish{}
}

// List 查询菜品列表。
func (c *Dish) List(r *ghttp.Request) {
	list, err := menulogic.ListAdminDishes(r.Context(), r.Get("category_id").Int64(), r.Get("status").Int())
	if err != nil {
		middleware.WriteError(r, consts.CodeSystemError, err.Error())
		return
	}
	middleware.WriteSuccess(r, list)
}

// Create 创建菜品。
func (c *Dish) Create(r *ghttp.Request) {
	var req adminapi.DishSaveReq
	if err := r.Parse(&req); err != nil {
		middleware.WriteError(r, consts.CodeInvalidParams, "菜品参数错误")
		return
	}

	dish, err := menulogic.CreateDish(r.Context(), toDishInput(0, req))
	if err != nil {
		middleware.WriteError(r, consts.CodeInvalidParams, err.Error())
		return
	}
	middleware.WriteSuccess(r, dish)
}

// Detail 查询菜品详情。
func (c *Dish) Detail(r *ghttp.Request) {
	dish, err := menulogic.GetAdminDish(r.Context(), r.Get("id").Int64())
	if err != nil {
		middleware.WriteError(r, mapMenuErrorCode(err), err.Error())
		return
	}
	middleware.WriteSuccess(r, dish)
}

// Update 更新菜品。
func (c *Dish) Update(r *ghttp.Request) {
	var req adminapi.DishSaveReq
	if err := r.Parse(&req); err != nil {
		middleware.WriteError(r, consts.CodeInvalidParams, "菜品参数错误")
		return
	}

	dish, err := menulogic.UpdateDish(r.Context(), toDishInput(r.Get("id").Int64(), req))
	if err != nil {
		middleware.WriteError(r, mapMenuErrorCode(err), err.Error())
		return
	}
	middleware.WriteSuccess(r, dish)
}

// Delete 下架菜品。
func (c *Dish) Delete(r *ghttp.Request) {
	if err := menulogic.DeleteDish(r.Context(), r.Get("id").Int64()); err != nil {
		middleware.WriteError(r, mapMenuErrorCode(err), err.Error())
		return
	}
	middleware.WriteSuccess(r, true)
}

// UpdateStatus 更新菜品状态。
func (c *Dish) UpdateStatus(r *ghttp.Request) {
	var req adminapi.DishStatusReq
	if err := r.Parse(&req); err != nil {
		middleware.WriteError(r, consts.CodeInvalidParams, "菜品状态参数错误")
		return
	}

	if err := menulogic.UpdateDishStatus(r.Context(), r.Get("id").Int64(), req.Status); err != nil {
		middleware.WriteError(r, mapMenuErrorCode(err), err.Error())
		return
	}
	middleware.WriteSuccess(r, true)
}

// toDishInput 转换菜品保存参数。
func toDishInput(id int64, req adminapi.DishSaveReq) menulogic.DishInput {
	return menulogic.DishInput{
		ID:          id,
		CategoryID:  req.CategoryID,
		Name:        req.Name,
		Description: req.Description,
		ImageURL:    req.ImageURL,
		Price:       req.Price,
		Unit:        req.Unit,
		Sort:        req.Sort,
	}
}
