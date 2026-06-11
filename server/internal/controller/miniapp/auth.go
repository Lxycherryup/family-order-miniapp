package miniapp

import (
	miniappapi "family-order/server/api/miniapp"
	"family-order/server/internal/consts"
	authlogic "family-order/server/internal/logic/auth"
	userlogic "family-order/server/internal/logic/user"
	"family-order/server/internal/middleware"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
)

// Auth 小程序认证控制器。
type Auth struct{}

// NewAuth 创建小程序认证控制器。
func NewAuth() *Auth {
	return &Auth{}
}

// Login 小程序微信登录。
func (c *Auth) Login(r *ghttp.Request) {
	var req miniappapi.LoginReq
	if err := r.Parse(&req); err != nil {
		middleware.WriteError(r, consts.CodeInvalidParams, "登录参数错误")
		return
	}

	client := newWechatClient(r)
	out, err := userlogic.LoginByWechatCode(r.Context(), userlogic.WechatLoginInput{
		Code:      req.Code,
		Nickname:  req.Nickname,
		AvatarURL: req.AvatarURL,
		Client:    client,
	})
	if err != nil {
		middleware.WriteError(r, consts.CodeUnauthorized, err.Error())
		return
	}

	middleware.WriteSuccess(r, miniappapi.LoginRes{
		Token:       out.Token,
		IsWhitelist: out.IsWhitelist,
	})
}

// newWechatClient 根据配置创建微信登录客户端。
func newWechatClient(r *ghttp.Request) authlogic.WechatClient {
	if g.Cfg().MustGet(r.Context(), "wechat.mockEnabled", false).Bool() {
		return authlogic.NewMockWechatClient(
			g.Cfg().MustGet(r.Context(), "wechat.mockOpenID", "dev-family-user").String(),
		)
	}
	return &authlogic.HTTPWechatClient{
		AppID:     g.Cfg().MustGet(r.Context(), "wechat.appId").String(),
		AppSecret: g.Cfg().MustGet(r.Context(), "wechat.appSecret").String(),
	}
}
