package user

import (
	"context"
	"fmt"
	"time"

	"family-order/server/internal/consts"

	"github.com/gogf/gf/v2/frame/g"
)

// AdminUser 管理端用户信息。
type AdminUser struct {
	ID          int64      `json:"id"`
	OpenID      string     `json:"openid"`
	Nickname    string     `json:"nickname"`
	AvatarURL   string     `json:"avatar_url"`
	Status      int        `json:"status"`
	IsWhitelist bool       `json:"is_whitelist"`
	LastLoginAt *time.Time `json:"last_login_at"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
}

// AdminUserListInput 管理端用户列表输入参数。
type AdminUserListInput struct {
	IsWhitelist string
	Status      int
	Page        int
	PageSize    int
}

// ListAdminUsers 查询管理端用户列表。
func ListAdminUsers(ctx context.Context, in AdminUserListInput) ([]AdminUser, error) {
	if in.Page <= 0 {
		in.Page = 1
	}
	if in.PageSize <= 0 || in.PageSize > 100 {
		in.PageSize = 20
	}

	model := g.DB().Model("users").Ctx(ctx).
		Fields("id,openid,nickname,avatar_url,status,is_whitelist,last_login_at,created_at,updated_at")
	if in.IsWhitelist == "true" {
		model = model.Where("is_whitelist", true)
	}
	if in.IsWhitelist == "false" {
		model = model.Where("is_whitelist", false)
	}
	if in.Status > 0 {
		model = model.Where("status", in.Status)
	}

	list := make([]AdminUser, 0)
	err := model.Page(in.Page, in.PageSize).
		Order("created_at DESC,id DESC").
		Scan(&list)
	if err != nil {
		return nil, fmt.Errorf("查询用户列表失败: %w", err)
	}
	return list, nil
}

// UpdateUserWhitelist 更新用户白名单状态。
func UpdateUserWhitelist(ctx context.Context, id int64, isWhitelist bool) error {
	if id <= 0 {
		return fmt.Errorf("用户ID不能为空")
	}

	result, err := g.DB().Model("users").Ctx(ctx).
		Where("id", id).
		Data(g.Map{"is_whitelist": isWhitelist, "updated_at": time.Now()}).
		Update()
	if err != nil {
		return fmt.Errorf("更新用户白名单失败: %w", err)
	}
	rows, _ := result.RowsAffected()
	if rows == 0 {
		return fmt.Errorf("用户不存在")
	}
	return nil
}

// UpdateUserStatus 更新用户状态。
func UpdateUserStatus(ctx context.Context, id int64, status int) error {
	if id <= 0 {
		return fmt.Errorf("用户ID不能为空")
	}
	if status != consts.StatusEnabled && status != consts.StatusDisabled {
		return fmt.Errorf("用户状态无效")
	}

	result, err := g.DB().Model("users").Ctx(ctx).
		Where("id", id).
		Data(g.Map{"status": status, "updated_at": time.Now()}).
		Update()
	if err != nil {
		return fmt.Errorf("更新用户状态失败: %w", err)
	}
	rows, _ := result.RowsAffected()
	if rows == 0 {
		return fmt.Errorf("用户不存在")
	}
	return nil
}
