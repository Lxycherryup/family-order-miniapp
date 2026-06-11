package user

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"family-order/server/internal/consts"
	authlogic "family-order/server/internal/logic/auth"

	"github.com/gogf/gf/v2/frame/g"
)

const (
	// miniappRole 小程序用户 JWT 角色。
	miniappRole = "miniapp"
)

// WechatLoginInput 微信登录输入参数。
type WechatLoginInput struct {
	Code      string
	Nickname  string
	AvatarURL string
	Client    authlogic.WechatClient
	Secret    string
	ExpiresIn time.Duration
}

// WechatLoginOutput 微信登录输出结果。
type WechatLoginOutput struct {
	Token       string `json:"token"`
	IsWhitelist bool   `json:"is_whitelist"`
	UserID      int64  `json:"user_id"`
	OpenID      string `json:"openid"`
}

// miniappUser 小程序用户数据库记录。
type miniappUser struct {
	ID          int64  `json:"id"`
	OpenID      string `json:"openid"`
	Nickname    string `json:"nickname"`
	AvatarURL   string `json:"avatar_url"`
	Status      int    `json:"status"`
	IsWhitelist bool   `json:"is_whitelist"`
}

// LoginByWechatCode 使用微信 code 登录并返回白名单状态。
func LoginByWechatCode(ctx context.Context, in WechatLoginInput) (*WechatLoginOutput, error) {
	if in.Client == nil {
		return nil, fmt.Errorf("微信客户端不能为空")
	}
	session, err := in.Client.Code2Session(ctx, in.Code)
	if err != nil {
		return nil, fmt.Errorf("微信登录失败: %w", err)
	}

	u, err := findOrCreateWechatUser(ctx, session.OpenID, in.Nickname, in.AvatarURL)
	if err != nil {
		return nil, err
	}

	out := &WechatLoginOutput{
		IsWhitelist: u.IsWhitelist && u.Status == consts.StatusEnabled,
		UserID:      u.ID,
		OpenID:      u.OpenID,
	}
	if !out.IsWhitelist {
		return out, nil
	}

	secret := in.Secret
	if secret == "" {
		secret = g.Cfg().MustGet(ctx, "jwt.secret").String()
	}
	expiresIn := in.ExpiresIn
	if expiresIn <= 0 {
		expireHours := g.Cfg().MustGet(ctx, "jwt.miniappExpireHours", 168).Int()
		expiresIn = time.Duration(expireHours) * time.Hour
	}
	token, err := authlogic.GenerateToken(authlogic.TokenClaimsInput{
		SubjectID: u.ID,
		Role:      miniappRole,
		Secret:    secret,
		ExpiresIn: expiresIn,
	})
	if err != nil {
		return nil, fmt.Errorf("生成小程序Token失败: %w", err)
	}
	out.Token = token
	return out, nil
}

// findOrCreateWechatUser 查询或创建微信用户，并更新登录信息。
func findOrCreateWechatUser(ctx context.Context, openID string, nickname string, avatarURL string) (*miniappUser, error) {
	if openID == "" {
		return nil, fmt.Errorf("微信OpenID不能为空")
	}

	var u miniappUser
	if err := g.DB().Model("users").Ctx(ctx).Where("openid", openID).Scan(&u); err != nil {
		if !errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("查询微信用户失败: %w", err)
		}
	}

	now := time.Now()
	if u.ID == 0 {
		id, err := g.DB().Model("users").Ctx(ctx).Data(g.Map{
			"openid":        openID,
			"nickname":      nickname,
			"avatar_url":    avatarURL,
			"status":        consts.StatusEnabled,
			"is_whitelist":  false,
			"last_login_at": now,
			"created_at":    now,
			"updated_at":    now,
		}).InsertAndGetId()
		if err != nil {
			return nil, fmt.Errorf("创建微信用户失败: %w", err)
		}
		return &miniappUser{
			ID:          id,
			OpenID:      openID,
			Nickname:    nickname,
			AvatarURL:   avatarURL,
			Status:      consts.StatusEnabled,
			IsWhitelist: false,
		}, nil
	}

	if _, err := g.DB().Model("users").Ctx(ctx).Where("id", u.ID).Data(g.Map{
		"nickname":      nickname,
		"avatar_url":    avatarURL,
		"last_login_at": now,
		"updated_at":    now,
	}).Update(); err != nil {
		return nil, fmt.Errorf("更新微信用户登录信息失败: %w", err)
	}
	u.Nickname = nickname
	u.AvatarURL = avatarURL
	return &u, nil
}
