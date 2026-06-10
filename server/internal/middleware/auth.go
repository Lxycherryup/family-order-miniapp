package middleware

import (
	"strings"

	"family-order/server/internal/consts"
	authlogic "family-order/server/internal/logic/auth"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
)

const (
	// CtxKeySubjectID 请求上下文中的主体 ID。
	CtxKeySubjectID = "subject_id"
	// CtxKeyRole 请求上下文中的角色。
	CtxKeyRole = "role"
)

// AdminAuth 校验管理端 JWT。
func AdminAuth(r *ghttp.Request) {
	authByRole(r, "admin", consts.CodeUnauthorized, "未登录或Token失效")
}

// MiniappAuth 校验小程序白名单用户 JWT。
func MiniappAuth(r *ghttp.Request) {
	authByRole(r, "miniapp", consts.CodeForbidden, "当前微信用户未加入点餐白名单")
}

// authByRole 按角色校验 JWT 并写入请求上下文。
func authByRole(r *ghttp.Request, role string, errorCode int, errorMsg string) {
	tokenString := extractBearerToken(r.GetHeader("Authorization"))
	if tokenString == "" {
		WriteError(r, errorCode, errorMsg)
		r.ExitAll()
		return
	}

	secret := g.Cfg().MustGet(r.Context(), "jwt.secret").String()
	claims, err := authlogic.ParseToken(tokenString, secret)
	if err != nil || claims.Role != role {
		WriteError(r, errorCode, errorMsg)
		r.ExitAll()
		return
	}

	r.SetCtxVar(CtxKeySubjectID, claims.SubjectID)
	r.SetCtxVar(CtxKeyRole, claims.Role)
	r.Middleware.Next()
}

// extractBearerToken 提取 Authorization 头中的 Bearer Token。
func extractBearerToken(header string) string {
	const prefix = "Bearer "
	if !strings.HasPrefix(header, prefix) {
		return ""
	}
	return strings.TrimSpace(strings.TrimPrefix(header, prefix))
}
