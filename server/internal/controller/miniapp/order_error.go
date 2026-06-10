package miniapp

import (
	"strings"

	"family-order/server/internal/consts"
)

// mapOrderErrorCode 映射小程序订单错误码。
func mapOrderErrorCode(err error) int {
	if err == nil {
		return consts.CodeSuccess
	}
	msg := err.Error()
	if strings.Contains(msg, "冲突") || strings.Contains(msg, "只能取消") {
		return consts.CodeConflict
	}
	if strings.Contains(msg, "不存在") || strings.Contains(msg, "下架") {
		return consts.CodeNotFound
	}
	if strings.Contains(msg, "不能为空") || strings.Contains(msg, "必须") || strings.Contains(msg, "超过") {
		return consts.CodeInvalidParams
	}
	return consts.CodeSystemError
}
