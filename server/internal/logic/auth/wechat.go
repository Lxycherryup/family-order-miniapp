package auth

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
)

const (
	// wechatCode2SessionURL 微信 code 换 session 接口地址。
	wechatCode2SessionURL = "https://api.weixin.qq.com/sns/jscode2session"
)

// WechatClient 微信登录客户端接口。
type WechatClient interface {
	Code2Session(ctx context.Context, code string) (*WechatSession, error)
}

// WechatSession 微信登录会话。
type WechatSession struct {
	OpenID     string `json:"openid"`
	SessionKey string `json:"session_key"`
	ErrCode    int    `json:"errcode"`
	ErrMsg     string `json:"errmsg"`
}

// HTTPWechatClient 基于 HTTP 的微信客户端。
type HTTPWechatClient struct {
	AppID     string
	AppSecret string
}

// MockWechatClient 本地开发使用的微信客户端。
type MockWechatClient struct {
	OpenID string
}

// NewMockWechatClient 创建本地开发微信 mock 客户端。
func NewMockWechatClient(openID string) *MockWechatClient {
	if openID == "" {
		openID = "dev-family-user"
	}
	return &MockWechatClient{OpenID: openID}
}

// Code2Session 返回本地开发固定 OpenID。
func (c *MockWechatClient) Code2Session(ctx context.Context, code string) (*WechatSession, error) {
	return &WechatSession{
		OpenID:     c.OpenID,
		SessionKey: "mock-session-key",
	}, nil
}

// Code2Session 使用微信 code 换取 OpenID。
func (c *HTTPWechatClient) Code2Session(ctx context.Context, code string) (*WechatSession, error) {
	if c.AppID == "" || c.AppSecret == "" {
		return nil, fmt.Errorf("微信小程序配置不能为空")
	}
	if code == "" {
		return nil, fmt.Errorf("微信登录code不能为空")
	}

	values := url.Values{}
	values.Set("appid", c.AppID)
	values.Set("secret", c.AppSecret)
	values.Set("js_code", code)
	values.Set("grant_type", "authorization_code")
	apiURL := wechatCode2SessionURL + "?" + values.Encode()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, apiURL, nil)
	if err != nil {
		return nil, fmt.Errorf("创建微信登录请求失败: %w", err)
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("请求微信登录接口失败: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("微信登录接口状态异常: %d", resp.StatusCode)
	}

	var session WechatSession
	if err := json.NewDecoder(resp.Body).Decode(&session); err != nil {
		return nil, fmt.Errorf("解析微信登录响应失败: %w", err)
	}
	if session.ErrCode != 0 {
		return nil, fmt.Errorf("微信登录失败: %s", session.ErrMsg)
	}
	if session.OpenID == "" {
		return nil, fmt.Errorf("微信登录未返回OpenID")
	}
	return &session, nil
}
