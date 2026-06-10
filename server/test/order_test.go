package test

import (
	"testing"

	"family-order/server/internal/consts"
	orderlogic "family-order/server/internal/logic/order"
)

// TestCanChangeStatus 验证订单状态流转规则。
func TestCanChangeStatus(t *testing.T) {
	if !orderlogic.CanChangeStatus(consts.OrderStatusPending, consts.OrderStatusCooking) {
		t.Fatalf("待处理应该允许变更为制作中")
	}
	if !orderlogic.CanChangeStatus(consts.OrderStatusCooking, consts.OrderStatusCompleted) {
		t.Fatalf("制作中应该允许变更为已完成")
	}
	if orderlogic.CanChangeStatus(consts.OrderStatusCompleted, consts.OrderStatusCanceled) {
		t.Fatalf("已完成不允许变更为已取消")
	}
}
