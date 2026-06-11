package admin

import (
	adminapi "family-order/server/api/admin"
	"family-order/server/internal/consts"
	userlogic "family-order/server/internal/logic/user"
	"family-order/server/internal/middleware"

	"github.com/gogf/gf/v2/net/ghttp"
)

// User 管理端用户控制器。
type User struct{}

// NewUser 创建管理端用户控制器。
func NewUser() *User {
	return &User{}
}

// List 查询用户列表。
func (c *User) List(r *ghttp.Request) {
	list, err := userlogic.ListAdminUsers(r.Context(), userlogic.AdminUserListInput{
		IsWhitelist: r.Get("is_whitelist").String(),
		Status:      r.Get("status").Int(),
		Page:        r.Get("page").Int(),
		PageSize:    r.Get("page_size").Int(),
	})
	if err != nil {
		middleware.WriteError(r, consts.CodeSystemError, err.Error())
		return
	}
	middleware.WriteSuccess(r, list)
}

// UpdateWhitelist 更新用户白名单状态。
func (c *User) UpdateWhitelist(r *ghttp.Request) {
	var req adminapi.UserWhitelistReq
	if err := r.Parse(&req); err != nil {
		middleware.WriteError(r, consts.CodeInvalidParams, "用户白名单参数错误")
		return
	}

	if err := userlogic.UpdateUserWhitelist(r.Context(), r.Get("id").Int64(), req.IsWhitelist); err != nil {
		middleware.WriteError(r, consts.CodeInvalidParams, err.Error())
		return
	}
	middleware.WriteSuccess(r, true)
}

// UpdateStatus 更新用户状态。
func (c *User) UpdateStatus(r *ghttp.Request) {
	var req adminapi.UserStatusReq
	if err := r.Parse(&req); err != nil {
		middleware.WriteError(r, consts.CodeInvalidParams, "用户状态参数错误")
		return
	}

	if err := userlogic.UpdateUserStatus(r.Context(), r.Get("id").Int64(), req.Status); err != nil {
		middleware.WriteError(r, consts.CodeInvalidParams, err.Error())
		return
	}
	middleware.WriteSuccess(r, true)
}
