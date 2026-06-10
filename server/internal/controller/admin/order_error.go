package admin

import (
	"strings"

	"family-order/server/internal/consts"
)

// mapAdminOrderErrorCode 映射管理端订单错误码。
func mapAdminOrderErrorCode(err error) int {
	if err == nil {
		return consts.CodeSuccess
	}
	msg := err.Error()
	if strings.Contains(msg, "冲突") {
		return consts.CodeConflict
	}
	if strings.Contains(msg, "不存在") {
		return consts.CodeNotFound
	}
	if strings.Contains(msg, "不能为空") || strings.Contains(msg, "必须") || strings.Contains(msg, "超过") {
		return consts.CodeInvalidParams
	}
	return consts.CodeSystemError
}
