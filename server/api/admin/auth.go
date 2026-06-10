package admin

import "github.com/gogf/gf/v2/frame/g"

// LoginReq 管理员登录请求。
type LoginReq struct {
	g.Meta   `path:"/api/admin/auth/login" method:"post" tags:"管理端认证" summary:"管理员登录"`
	Username string `json:"username" v:"required#请输入管理员账号"`
	Password string `json:"password" v:"required#请输入管理员密码"`
}

// LoginRes 管理员登录响应。
type LoginRes struct {
	Token    string `json:"token"`
	AdminID  int64  `json:"admin_id"`
	Username string `json:"username"`
}

// LogoutReq 管理员退出登录请求。
type LogoutReq struct {
	g.Meta `path:"/api/admin/auth/logout" method:"post" tags:"管理端认证" summary:"管理员退出登录"`
}

// LogoutRes 管理员退出登录响应。
type LogoutRes struct {
	Success bool `json:"success"`
}

// ProfileReq 管理员资料请求。
type ProfileReq struct {
	g.Meta `path:"/api/admin/auth/profile" method:"get" tags:"管理端认证" summary:"管理员资料"`
}

// ProfileRes 管理员资料响应。
type ProfileRes struct {
	AdminID int64  `json:"admin_id"`
	Role    string `json:"role"`
}
