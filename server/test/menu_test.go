package test

import (
	"testing"

	menulogic "family-order/server/internal/logic/menu"
)

// TestValidateDishQuantity 验证菜品数量校验。
func TestValidateDishQuantity(t *testing.T) {
	if err := menulogic.ValidateDishQuantity(0); err == nil {
		t.Fatalf("数量为0时应该返回错误")
	}
	if err := menulogic.ValidateDishQuantity(1); err != nil {
		t.Fatalf("数量为1时不应该返回错误: %v", err)
	}
}
