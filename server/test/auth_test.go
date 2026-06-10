package test

import (
	"testing"
	"time"

	authlogic "family-order/server/internal/logic/auth"
)

// TestJWTGenerateAndParse 验证 JWT 可以生成和解析。
func TestJWTGenerateAndParse(t *testing.T) {
	token, err := authlogic.GenerateToken(authlogic.TokenClaimsInput{
		SubjectID: 1,
		Role:      "admin",
		Secret:    "test-secret",
		ExpiresIn: time.Hour,
	})
	if err != nil {
		t.Fatalf("生成 Token 失败: %v", err)
	}

	claims, err := authlogic.ParseToken(token, "test-secret")
	if err != nil {
		t.Fatalf("解析 Token 失败: %v", err)
	}
	if claims.SubjectID != 1 || claims.Role != "admin" {
		t.Fatalf("Token 内容不符合预期: %+v", claims)
	}
}
