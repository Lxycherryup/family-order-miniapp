package admin

import (
	"strings"

	"family-order/server/internal/consts"
)

// mapMenuErrorCode 映射菜单模块错误码。
func mapMenuErrorCode(err error) int {
	if err == nil {
		return consts.CodeSuccess
	}
	msg := err.Error()
	if strings.Contains(msg, "不存在") {
		return consts.CodeNotFound
	}
	if strings.Contains(msg, "不能为空") || strings.Contains(msg, "无效") || strings.Contains(msg, "不能小于") {
		return consts.CodeInvalidParams
	}
	return consts.CodeSystemError
}

// mapCategoryDeleteErrorCode 映射分类删除错误码。
func mapCategoryDeleteErrorCode(err error) int {
	if err == nil {
		return consts.CodeSuccess
	}
	if strings.Contains(err.Error(), "存在菜品") {
		return consts.CodeConflict
	}
	return mapMenuErrorCode(err)
}
