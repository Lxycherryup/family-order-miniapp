package controller

import (
	"context"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
)

// Health 健康检查控制器。
type Health struct{}

// NewHealth 创建健康检查控制器。
func NewHealth() *Health {
	return &Health{}
}

// HealthReq 健康检查请求。
type HealthReq struct {
	g.Meta `path:"/health" method:"get" tags:"系统" summary:"健康检查"`
}

// HealthRes 健康检查响应。
type HealthRes struct {
	Status string `json:"status"`
}

// HealthPayload 健康检查统一响应。
type HealthPayload struct {
	Code int       `json:"code"`
	Msg  string    `json:"msg"`
	Data HealthRes `json:"data"`
}

// NewHealthPayload 创建健康检查统一响应。
func NewHealthPayload() HealthPayload {
	return HealthPayload{
		Code: 0,
		Msg:  "success",
		Data: HealthRes{Status: "ok"},
	}
}

// Health 返回服务健康状态。
func (c *Health) Health(ctx context.Context, req *HealthReq) (res *HealthRes, err error) {
	return &HealthRes{Status: "ok"}, nil
}

// Handler 显式输出健康检查 JSON 响应。
func (c *Health) Handler(r *ghttp.Request) {
	r.Response.WriteJson(NewHealthPayload())
}
