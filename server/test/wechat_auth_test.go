package test

import (
	"context"
	"testing"

	authlogic "family-order/server/internal/logic/auth"
)

// fakeWechatClient 模拟微信登录客户端。
type fakeWechatClient struct{}

// Code2Session 模拟微信 code 换 OpenID。
func (fakeWechatClient) Code2Session(ctx context.Context, code string) (*authlogic.WechatSession, error) {
	return &authlogic.WechatSession{OpenID: "openid-family-001", SessionKey: "session-key"}, nil
}

// TestWechatClientInterface 验证微信登录依赖可以被替换为 mock。
func TestWechatClientInterface(t *testing.T) {
	client := fakeWechatClient{}
	session, err := client.Code2Session(context.Background(), "test-code")
	if err != nil {
		t.Fatalf("模拟微信登录失败: %v", err)
	}
	if session.OpenID != "openid-family-001" {
		t.Fatalf("OpenID 不符合预期: %s", session.OpenID)
	}
}

// TestWechatMockClient 验证开发环境 mock 微信客户端返回固定 OpenID。
func TestWechatMockClient(t *testing.T) {
	client := authlogic.NewMockWechatClient("dev-family-user")
	session, err := client.Code2Session(context.Background(), "any-code")
	if err != nil {
		t.Fatalf("mock 微信登录不应该失败: %v", err)
	}
	if session.OpenID != "dev-family-user" {
		t.Fatalf("期望 mock OpenID 为 dev-family-user，实际为 %s", session.OpenID)
	}
}
