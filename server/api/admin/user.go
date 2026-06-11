package admin

import "github.com/gogf/gf/v2/frame/g"

// UserListReq 管理端用户列表请求。
type UserListReq struct {
	g.Meta      `path:"/api/admin/users" method:"get" tags:"管理端用户" summary:"用户列表"`
	IsWhitelist string `json:"is_whitelist"`
	Status      int    `json:"status"`
	Page        int    `json:"page"`
	PageSize    int    `json:"page_size"`
}

// UserWhitelistReq 用户白名单请求。
type UserWhitelistReq struct {
	IsWhitelist bool `json:"is_whitelist"`
}

// UserStatusReq 用户状态请求。
type UserStatusReq struct {
	Status int `json:"status"`
}
