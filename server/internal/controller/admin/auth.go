package admin

import (
	adminapi "family-order/server/api/admin"
	"family-order/server/internal/consts"
	authlogic "family-order/server/internal/logic/auth"
	"family-order/server/internal/middleware"

	"github.com/gogf/gf/v2/net/ghttp"
)

// Auth 管理端认证控制器。
type Auth struct{}

// NewAuth 创建管理端认证控制器。
func NewAuth() *Auth {
	return &Auth{}
}

// Login 管理员登录。
func (c *Auth) Login(r *ghttp.Request) {
	var req adminapi.LoginReq
	if err := r.Parse(&req); err != nil {
		middleware.WriteError(r, consts.CodeInvalidParams, "登录参数错误")
		return
	}

	out, err := authlogic.LoginAdmin(r.Context(), authlogic.AdminLoginInput{
		Username: req.Username,
		Password: req.Password,
	})
	if err != nil {
		middleware.WriteError(r, consts.CodeUnauthorized, err.Error())
		return
	}

	middleware.WriteSuccess(r, adminapi.LoginRes{
		Token:    out.Token,
		AdminID:  out.AdminID,
		Username: out.Username,
	})
}

// Logout 管理员退出登录。
func (c *Auth) Logout(r *ghttp.Request) {
	middleware.WriteSuccess(r, adminapi.LogoutRes{Success: true})
}

// Profile 返回当前管理员资料。
func (c *Auth) Profile(r *ghttp.Request) {
	middleware.WriteSuccess(r, adminapi.ProfileRes{
		AdminID: r.GetCtxVar(middleware.CtxKeySubjectID).Int64(),
		Role:    r.GetCtxVar(middleware.CtxKeyRole).String(),
	})
}
