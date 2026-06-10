package consts

const (
	// CodeSuccess 成功。
	CodeSuccess = 0
	// CodeInvalidParams 参数错误。
	CodeInvalidParams = 40001
	// CodeUnauthorized 未登录或 Token 失效。
	CodeUnauthorized = 40101
	// CodeForbidden 无访问权限。
	CodeForbidden = 40301
	// CodeNotFound 资源不存在。
	CodeNotFound = 40401
	// CodeConflict 状态冲突。
	CodeConflict = 40901
	// CodeSystemError 系统错误。
	CodeSystemError = 50001
)
