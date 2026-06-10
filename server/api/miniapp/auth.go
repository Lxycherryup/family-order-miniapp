package miniapp

import "github.com/gogf/gf/v2/frame/g"

// LoginReq 小程序微信登录请求。
type LoginReq struct {
	g.Meta    `path:"/api/miniapp/auth/login" method:"post" tags:"小程序认证" summary:"微信登录"`
	Code      string `json:"code" v:"required#请输入微信登录code"`
	Nickname  string `json:"nickname"`
	AvatarURL string `json:"avatar_url"`
}

// LoginRes 小程序微信登录响应。
type LoginRes struct {
	Token       string `json:"token"`
	IsWhitelist bool   `json:"is_whitelist"`
}
