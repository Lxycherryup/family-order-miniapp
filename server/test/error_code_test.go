package test

import (
	"testing"

	"family-order/server/internal/consts"
	"family-order/server/internal/model"
)

// TestErrorCodeValues 验证核心错误码稳定。
func TestErrorCodeValues(t *testing.T) {
	cases := map[string]int{
		"成功":    consts.CodeSuccess,
		"参数错误":  consts.CodeInvalidParams,
		"未登录":   consts.CodeUnauthorized,
		"不在白名单": consts.CodeForbidden,
		"资源不存在": consts.CodeNotFound,
		"状态冲突":  consts.CodeConflict,
		"系统错误":  consts.CodeSystemError,
	}
	expected := map[string]int{
		"成功":    0,
		"参数错误":  40001,
		"未登录":   40101,
		"不在白名单": 40301,
		"资源不存在": 40401,
		"状态冲突":  40901,
		"系统错误":  50001,
	}
	for name, actual := range cases {
		if actual != expected[name] {
			t.Fatalf("%s 错误码期望 %d，实际 %d", name, expected[name], actual)
		}
	}
}

// TestResponseBuilders 验证统一响应构造函数返回稳定结构。
func TestResponseBuilders(t *testing.T) {
	success := model.NewSuccessResponse(map[string]string{"status": "ok"})
	if success.Code != consts.CodeSuccess || success.Msg != "success" || success.Data == nil {
		t.Fatalf("成功响应不符合预期: %+v", success)
	}

	failed := model.NewErrorResponse(consts.CodeInvalidParams, "参数错误")
	if failed.Code != consts.CodeInvalidParams || failed.Msg != "参数错误" || failed.Data != nil {
		t.Fatalf("失败响应不符合预期: %+v", failed)
	}
}
