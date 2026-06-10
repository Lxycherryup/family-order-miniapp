package auth

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// TokenClaimsInput 生成 Token 的输入参数。
type TokenClaimsInput struct {
	SubjectID int64
	Role      string
	Secret    string
	ExpiresIn time.Duration
}

// TokenClaims Token 载荷。
type TokenClaims struct {
	SubjectID int64  `json:"subject_id"`
	Role      string `json:"role"`
	jwt.RegisteredClaims
}

// GenerateToken 生成 JWT。
func GenerateToken(in TokenClaimsInput) (string, error) {
	if in.Secret == "" {
		return "", fmt.Errorf("JWT密钥不能为空")
	}
	if in.SubjectID <= 0 {
		return "", fmt.Errorf("Token主体ID不能为空")
	}
	if in.Role == "" {
		return "", fmt.Errorf("Token角色不能为空")
	}
	if in.ExpiresIn <= 0 {
		return "", fmt.Errorf("Token有效期必须大于0")
	}

	now := time.Now()
	claims := TokenClaims{
		SubjectID: in.SubjectID,
		Role:      in.Role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(now.Add(in.ExpiresIn)),
			IssuedAt:  jwt.NewNumericDate(now),
		},
	}
	return jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString([]byte(in.Secret))
}

// ParseToken 解析 JWT。
func ParseToken(tokenString string, secret string) (*TokenClaims, error) {
	if tokenString == "" {
		return nil, fmt.Errorf("Token不能为空")
	}
	if secret == "" {
		return nil, fmt.Errorf("JWT密钥不能为空")
	}

	token, err := jwt.ParseWithClaims(tokenString, &TokenClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Token签名方法无效")
		}
		return []byte(secret), nil
	})
	if err != nil {
		return nil, fmt.Errorf("解析Token失败: %w", err)
	}
	claims, ok := token.Claims.(*TokenClaims)
	if !ok || !token.Valid {
		return nil, fmt.Errorf("Token无效")
	}
	return claims, nil
}
