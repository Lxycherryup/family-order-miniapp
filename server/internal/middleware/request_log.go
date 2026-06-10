package middleware

import (
	"fmt"
	"time"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
)

// RequestLog 记录请求日志。
func RequestLog(r *ghttp.Request) {
	start := time.Now()
	requestID := r.GetHeader("X-Request-ID")
	if requestID == "" {
		requestID = fmt.Sprintf("%d", start.UnixNano())
	}
	r.Response.Header().Set("X-Request-ID", requestID)

	r.Middleware.Next()

	errMsg := ""
	if err := r.GetError(); err != nil {
		errMsg = err.Error()
	}
	userID := r.GetHeader("X-User-ID", "0")
	g.Log().Infof(r.Context(), "请求完成 request_id=%s user_id=%s method=%s path=%s status=%d cost=%s error=%s",
		requestID, userID, r.Method, r.URL.Path, r.Response.Status, time.Since(start), errMsg)
}
