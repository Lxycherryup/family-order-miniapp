package middleware

import (
	"family-order/server/internal/model"

	"github.com/gogf/gf/v2/net/ghttp"
)

// WriteSuccess 写入统一成功响应。
func WriteSuccess(r *ghttp.Request, data interface{}) {
	r.Response.WriteJson(model.NewSuccessResponse(data))
}

// WriteError 写入统一失败响应。
func WriteError(r *ghttp.Request, code int, msg string) {
	r.Response.WriteJson(model.NewErrorResponse(code, msg))
}
