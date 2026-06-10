package menu

import (
	"context"
	"fmt"
	"time"

	"family-order/server/internal/consts"

	"github.com/gogf/gf/v2/frame/g"
)

// Category 菜品分类。
type Category struct {
	ID        int64     `json:"id"`
	Name      string    `json:"name"`
	Sort      int       `json:"sort"`
	Status    int       `json:"status"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// CategoryInput 分类保存输入参数。
type CategoryInput struct {
	ID     int64
	Name   string
	Sort   int
	Status int
}

// ListAdminCategories 查询管理端分类列表。
func ListAdminCategories(ctx context.Context) ([]Category, error) {
	var list []Category
	err := g.DB().Model("menu_categories").Ctx(ctx).
		Fields("id,name,sort,status,created_at,updated_at").
		Order("sort ASC,id ASC").
		Scan(&list)
	if err != nil {
		return nil, fmt.Errorf("查询分类列表失败: %w", err)
	}
	return list, nil
}

// ListEnabledCategories 查询启用分类列表。
func ListEnabledCategories(ctx context.Context) ([]Category, error) {
	var list []Category
	err := g.DB().Model("menu_categories").Ctx(ctx).
		Fields("id,name,sort,status,created_at,updated_at").
		Where("status", consts.StatusEnabled).
		Order("sort ASC,id ASC").
		Scan(&list)
	if err != nil {
		return nil, fmt.Errorf("查询启用分类列表失败: %w", err)
	}
	return list, nil
}

// CreateCategory 创建菜品分类。
func CreateCategory(ctx context.Context, in CategoryInput) (*Category, error) {
	if in.Name == "" {
		return nil, fmt.Errorf("分类名称不能为空")
	}
	if in.Sort == 0 {
		in.Sort = 100
	}
	if in.Status == 0 {
		in.Status = consts.StatusEnabled
	}

	now := time.Now()
	id, err := g.DB().Model("menu_categories").Ctx(ctx).Data(g.Map{
		"name":       in.Name,
		"sort":       in.Sort,
		"status":     in.Status,
		"created_at": now,
		"updated_at": now,
	}).InsertAndGetId()
	if err != nil {
		return nil, fmt.Errorf("创建分类失败: %w", err)
	}
	return GetCategory(ctx, id)
}

// UpdateCategory 更新菜品分类。
func UpdateCategory(ctx context.Context, in CategoryInput) (*Category, error) {
	if in.ID <= 0 {
		return nil, fmt.Errorf("分类ID不能为空")
	}
	if in.Name == "" {
		return nil, fmt.Errorf("分类名称不能为空")
	}
	if in.Sort == 0 {
		in.Sort = 100
	}

	result, err := g.DB().Model("menu_categories").Ctx(ctx).
		Where("id", in.ID).
		Data(g.Map{
			"name":       in.Name,
			"sort":       in.Sort,
			"updated_at": time.Now(),
		}).Update()
	if err != nil {
		return nil, fmt.Errorf("更新分类失败: %w", err)
	}
	rows, _ := result.RowsAffected()
	if rows == 0 {
		return nil, fmt.Errorf("分类不存在")
	}
	return GetCategory(ctx, in.ID)
}

// UpdateCategoryStatus 更新分类状态。
func UpdateCategoryStatus(ctx context.Context, id int64, status int) error {
	if id <= 0 {
		return fmt.Errorf("分类ID不能为空")
	}
	if status != consts.StatusEnabled && status != consts.StatusDisabled {
		return fmt.Errorf("分类状态无效")
	}

	result, err := g.DB().Model("menu_categories").Ctx(ctx).
		Where("id", id).
		Data(g.Map{"status": status, "updated_at": time.Now()}).
		Update()
	if err != nil {
		return fmt.Errorf("更新分类状态失败: %w", err)
	}
	rows, _ := result.RowsAffected()
	if rows == 0 {
		return fmt.Errorf("分类不存在")
	}
	return nil
}

// DeleteCategory 删除菜品分类。
func DeleteCategory(ctx context.Context, id int64) error {
	if id <= 0 {
		return fmt.Errorf("分类ID不能为空")
	}

	count, err := g.DB().Model("dishes").Ctx(ctx).Where("category_id", id).Count()
	if err != nil {
		return fmt.Errorf("检查分类菜品失败: %w", err)
	}
	if count > 0 {
		return fmt.Errorf("分类下存在菜品，不能删除")
	}

	result, err := g.DB().Model("menu_categories").Ctx(ctx).Where("id", id).Delete()
	if err != nil {
		return fmt.Errorf("删除分类失败: %w", err)
	}
	rows, _ := result.RowsAffected()
	if rows == 0 {
		return fmt.Errorf("分类不存在")
	}
	return nil
}

// GetCategory 查询单个分类。
func GetCategory(ctx context.Context, id int64) (*Category, error) {
	var category Category
	err := g.DB().Model("menu_categories").Ctx(ctx).
		Fields("id,name,sort,status,created_at,updated_at").
		Where("id", id).
		Scan(&category)
	if err != nil {
		return nil, fmt.Errorf("查询分类失败: %w", err)
	}
	if category.ID == 0 {
		return nil, fmt.Errorf("分类不存在")
	}
	return &category, nil
}
