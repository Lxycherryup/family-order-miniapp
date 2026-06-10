package auth

import (
	"context"
	"fmt"
	"time"

	"family-order/server/internal/consts"

	"github.com/gogf/gf/v2/frame/g"
	"golang.org/x/crypto/bcrypt"
)

const adminRole = "admin"

// AdminLoginInput 管理员登录输入参数。
type AdminLoginInput struct {
	Username  string
	Password  string
	Secret    string
	ExpiresIn time.Duration
}

// AdminLoginOutput 管理员登录输出结果。
type AdminLoginOutput struct {
	Token    string `json:"token"`
	AdminID  int64  `json:"admin_id"`
	Username string `json:"username"`
}

// adminUser 管理员数据库记录。
type adminUser struct {
	ID           int64  `json:"id"`
	Username     string `json:"username"`
	PasswordHash string `json:"password_hash"`
	Status       int    `json:"status"`
}

// LoginAdmin 校验管理员账号密码并生成 JWT。
func LoginAdmin(ctx context.Context, in AdminLoginInput) (*AdminLoginOutput, error) {
	if in.Username == "" || in.Password == "" {
		return nil, fmt.Errorf("管理员账号或密码不能为空")
	}

	var admin adminUser
	if err := g.DB().Model("admin_users").Ctx(ctx).Where("username", in.Username).Scan(&admin); err != nil {
		return nil, fmt.Errorf("查询管理员失败: %w", err)
	}
	if admin.ID == 0 {
		return nil, fmt.Errorf("管理员账号或密码错误")
	}
	if admin.Status != consts.StatusEnabled {
		return nil, fmt.Errorf("管理员账号已禁用")
	}
	if err := bcrypt.CompareHashAndPassword([]byte(admin.PasswordHash), []byte(in.Password)); err != nil {
		return nil, fmt.Errorf("管理员账号或密码错误")
	}

	secret := in.Secret
	if secret == "" {
		secret = g.Cfg().MustGet(ctx, "jwt.secret").String()
	}
	expiresIn := in.ExpiresIn
	if expiresIn <= 0 {
		expireHours := g.Cfg().MustGet(ctx, "jwt.adminExpireHours", 24).Int()
		expiresIn = time.Duration(expireHours) * time.Hour
	}
	token, err := GenerateToken(TokenClaimsInput{
		SubjectID: admin.ID,
		Role:      adminRole,
		Secret:    secret,
		ExpiresIn: expiresIn,
	})
	if err != nil {
		return nil, fmt.Errorf("生成管理员Token失败: %w", err)
	}

	return &AdminLoginOutput{
		Token:    token,
		AdminID:  admin.ID,
		Username: admin.Username,
	}, nil
}
