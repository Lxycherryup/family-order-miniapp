package test

import (
	"context"
	"testing"

	"family-order/server/internal/controller"
)

// TestHealthController 验证健康检查控制器返回正常状态。
func TestHealthController(t *testing.T) {
	res, err := controller.NewHealth().Health(context.Background(), &controller.HealthReq{})
	if err != nil {
		t.Fatalf("健康检查不应该返回错误: %v", err)
	}
	if res.Status != "ok" {
		t.Fatalf("期望健康状态为 ok，实际为 %s", res.Status)
	}
}

// TestHealthPayload 验证健康检查接口响应符合统一结构。
func TestHealthPayload(t *testing.T) {
	payload := controller.NewHealthPayload()
	if payload.Code != 0 {
		t.Fatalf("期望响应码为 0，实际为 %d", payload.Code)
	}
	if payload.Msg != "success" {
		t.Fatalf("期望响应消息为 success，实际为 %s", payload.Msg)
	}
	if payload.Data.Status != "ok" {
		t.Fatalf("期望健康状态为 ok，实际为 %s", payload.Data.Status)
	}
}
