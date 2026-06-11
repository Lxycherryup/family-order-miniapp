package menu

import (
	"context"
	"fmt"
	"time"

	"family-order/server/internal/consts"

	"github.com/gogf/gf/v2/frame/g"
)

// Dish 菜品。
type Dish struct {
	ID          int64     `json:"id"`
	CategoryID  int64     `json:"category_id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	ImageURL    string    `json:"image_url"`
	Price       float64   `json:"price"`
	Unit        string    `json:"unit"`
	Status      int       `json:"status"`
	Sort        int       `json:"sort"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// DishInput 菜品保存输入参数。
type DishInput struct {
	ID          int64
	CategoryID  int64
	Name        string
	Description string
	ImageURL    string
	Price       float64
	Unit        string
	Status      int
	Sort        int
}

// ValidateDishQuantity 校验点餐数量。
func ValidateDishQuantity(quantity int) error {
	if quantity <= 0 {
		return fmt.Errorf("菜品数量必须大于0")
	}
	if quantity > 99 {
		return fmt.Errorf("单个菜品数量不能超过99")
	}
	return nil
}

// ListAdminDishes 查询管理端菜品列表。
func ListAdminDishes(ctx context.Context, categoryID int64, status int) ([]Dish, error) {
	model := g.DB().Model("dishes").Ctx(ctx).
		Fields("id,category_id,name,description,image_url,price,unit,status,sort,created_at,updated_at")
	if categoryID > 0 {
		model = model.Where("category_id", categoryID)
	}
	if status > 0 {
		model = model.Where("status", status)
	}

	list := make([]Dish, 0)
	if err := model.Order("sort ASC,id ASC").Scan(&list); err != nil {
		return nil, fmt.Errorf("查询菜品列表失败: %w", err)
	}
	return list, nil
}

// ListOnlineDishes 查询小程序上架菜品列表。
func ListOnlineDishes(ctx context.Context, categoryID int64) ([]Dish, error) {
	model := g.DB().Model("dishes").Ctx(ctx).
		Fields("id,category_id,name,description,image_url,price,unit,status,sort,created_at,updated_at").
		Where("status", consts.DishStatusOn)
	if categoryID > 0 {
		model = model.Where("category_id", categoryID)
	}

	list := make([]Dish, 0)
	if err := model.Order("sort ASC,id ASC").Scan(&list); err != nil {
		return nil, fmt.Errorf("查询上架菜品列表失败: %w", err)
	}
	return list, nil
}

// CreateDish 创建菜品。
func CreateDish(ctx context.Context, in DishInput) (*Dish, error) {
	if err := validateDishInput(ctx, in); err != nil {
		return nil, err
	}
	in = fillDishDefault(in)

	now := time.Now()
	id, err := g.DB().Model("dishes").Ctx(ctx).Data(g.Map{
		"category_id": in.CategoryID,
		"name":        in.Name,
		"description": in.Description,
		"image_url":   in.ImageURL,
		"price":       in.Price,
		"unit":        in.Unit,
		"status":      in.Status,
		"sort":        in.Sort,
		"created_at":  now,
		"updated_at":  now,
	}).InsertAndGetId()
	if err != nil {
		return nil, fmt.Errorf("创建菜品失败: %w", err)
	}
	return GetAdminDish(ctx, id)
}

// UpdateDish 更新菜品。
func UpdateDish(ctx context.Context, in DishInput) (*Dish, error) {
	if in.ID <= 0 {
		return nil, fmt.Errorf("菜品ID不能为空")
	}
	if err := validateDishInput(ctx, in); err != nil {
		return nil, err
	}
	in = fillDishDefault(in)

	result, err := g.DB().Model("dishes").Ctx(ctx).
		Where("id", in.ID).
		Data(g.Map{
			"category_id": in.CategoryID,
			"name":        in.Name,
			"description": in.Description,
			"image_url":   in.ImageURL,
			"price":       in.Price,
			"unit":        in.Unit,
			"sort":        in.Sort,
			"updated_at":  time.Now(),
		}).Update()
	if err != nil {
		return nil, fmt.Errorf("更新菜品失败: %w", err)
	}
	rows, _ := result.RowsAffected()
	if rows == 0 {
		return nil, fmt.Errorf("菜品不存在")
	}
	return GetAdminDish(ctx, in.ID)
}

// UpdateDishStatus 更新菜品状态。
func UpdateDishStatus(ctx context.Context, id int64, status int) error {
	if id <= 0 {
		return fmt.Errorf("菜品ID不能为空")
	}
	if status != consts.DishStatusOn && status != consts.DishStatusOff {
		return fmt.Errorf("菜品状态无效")
	}

	result, err := g.DB().Model("dishes").Ctx(ctx).
		Where("id", id).
		Data(g.Map{"status": status, "updated_at": time.Now()}).
		Update()
	if err != nil {
		return fmt.Errorf("更新菜品状态失败: %w", err)
	}
	rows, _ := result.RowsAffected()
	if rows == 0 {
		return fmt.Errorf("菜品不存在")
	}
	return nil
}

// DeleteDish 下架菜品，不物理删除历史数据。
func DeleteDish(ctx context.Context, id int64) error {
	return UpdateDishStatus(ctx, id, consts.DishStatusOff)
}

// GetAdminDish 查询管理端菜品详情。
func GetAdminDish(ctx context.Context, id int64) (*Dish, error) {
	return getDish(ctx, id, false)
}

// GetOnlineDish 查询小程序上架菜品详情。
func GetOnlineDish(ctx context.Context, id int64) (*Dish, error) {
	return getDish(ctx, id, true)
}

// getDish 查询菜品详情。
func getDish(ctx context.Context, id int64, onlyOnline bool) (*Dish, error) {
	if id <= 0 {
		return nil, fmt.Errorf("菜品ID不能为空")
	}

	model := g.DB().Model("dishes").Ctx(ctx).
		Fields("id,category_id,name,description,image_url,price,unit,status,sort,created_at,updated_at").
		Where("id", id)
	if onlyOnline {
		model = model.Where("status", consts.DishStatusOn)
	}

	var dish Dish
	if err := model.Scan(&dish); err != nil {
		return nil, fmt.Errorf("查询菜品失败: %w", err)
	}
	if dish.ID == 0 {
		return nil, fmt.Errorf("菜品不存在")
	}
	return &dish, nil
}

// validateDishInput 校验菜品保存输入。
func validateDishInput(ctx context.Context, in DishInput) error {
	if in.Name == "" {
		return fmt.Errorf("菜品名称不能为空")
	}
	if in.Price < 0 {
		return fmt.Errorf("菜品价格不能小于0")
	}
	if in.CategoryID <= 0 {
		return fmt.Errorf("菜品分类不能为空")
	}

	count, err := g.DB().Model("menu_categories").Ctx(ctx).
		Where("id", in.CategoryID).
		Where("status", consts.StatusEnabled).
		Count()
	if err != nil {
		return fmt.Errorf("检查菜品分类失败: %w", err)
	}
	if count == 0 {
		return fmt.Errorf("菜品分类不存在或已禁用")
	}
	return nil
}

// fillDishDefault 填充菜品默认值。
func fillDishDefault(in DishInput) DishInput {
	if in.Unit == "" {
		in.Unit = "份"
	}
	if in.Sort == 0 {
		in.Sort = 100
	}
	if in.Status == 0 {
		in.Status = consts.DishStatusOn
	}
	return in
}
