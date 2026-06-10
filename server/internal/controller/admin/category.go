package admin

import (
	adminapi "family-order/server/api/admin"
	"family-order/server/internal/consts"
	menulogic "family-order/server/internal/logic/menu"
	"family-order/server/internal/middleware"

	"github.com/gogf/gf/v2/net/ghttp"
)

// Category 管理端分类控制器。
type Category struct{}

// NewCategory 创建管理端分类控制器。
func NewCategory() *Category {
	return &Category{}
}

// List 查询分类列表。
func (c *Category) List(r *ghttp.Request) {
	list, err := menulogic.ListAdminCategories(r.Context())
	if err != nil {
		middleware.WriteError(r, consts.CodeSystemError, err.Error())
		return
	}
	middleware.WriteSuccess(r, list)
}

// Create 创建分类。
func (c *Category) Create(r *ghttp.Request) {
	var req adminapi.CategorySaveReq
	if err := r.Parse(&req); err != nil {
		middleware.WriteError(r, consts.CodeInvalidParams, "分类参数错误")
		return
	}

	category, err := menulogic.CreateCategory(r.Context(), menulogic.CategoryInput{
		Name: req.Name,
		Sort: req.Sort,
	})
	if err != nil {
		middleware.WriteError(r, consts.CodeInvalidParams, err.Error())
		return
	}
	middleware.WriteSuccess(r, category)
}

// Update 更新分类。
func (c *Category) Update(r *ghttp.Request) {
	var req adminapi.CategorySaveReq
	if err := r.Parse(&req); err != nil {
		middleware.WriteError(r, consts.CodeInvalidParams, "分类参数错误")
		return
	}

	category, err := menulogic.UpdateCategory(r.Context(), menulogic.CategoryInput{
		ID:   r.Get("id").Int64(),
		Name: req.Name,
		Sort: req.Sort,
	})
	if err != nil {
		middleware.WriteError(r, mapMenuErrorCode(err), err.Error())
		return
	}
	middleware.WriteSuccess(r, category)
}

// Delete 删除分类。
func (c *Category) Delete(r *ghttp.Request) {
	if err := menulogic.DeleteCategory(r.Context(), r.Get("id").Int64()); err != nil {
		middleware.WriteError(r, mapCategoryDeleteErrorCode(err), err.Error())
		return
	}
	middleware.WriteSuccess(r, true)
}

// UpdateStatus 更新分类状态。
func (c *Category) UpdateStatus(r *ghttp.Request) {
	var req adminapi.CategoryStatusReq
	if err := r.Parse(&req); err != nil {
		middleware.WriteError(r, consts.CodeInvalidParams, "分类状态参数错误")
		return
	}

	if err := menulogic.UpdateCategoryStatus(r.Context(), r.Get("id").Int64(), req.Status); err != nil {
		middleware.WriteError(r, mapMenuErrorCode(err), err.Error())
		return
	}
	middleware.WriteSuccess(r, true)
}
