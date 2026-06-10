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
	tokenString := extractBearerToken(r.GetHeader("Authorization"))
	if tokenString == "" {
		WriteError(r, consts.CodeUnauthorized, "未登录或Token失效")
		r.ExitAll()
		return
	}

	secret := g.Cfg().MustGet(r.Context(), "jwt.secret").String()
	claims, err := authlogic.ParseToken(tokenString, secret)
	if err != nil || claims.Role != "admin" {
		WriteError(r, consts.CodeUnauthorized, "未登录或Token失效")
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
