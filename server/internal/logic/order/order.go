package order

import (
	"context"
	"fmt"
	"time"

	"family-order/server/internal/consts"
	menulogic "family-order/server/internal/logic/menu"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/util/grand"
)

// CreateOrderInput 创建订单输入参数。
type CreateOrderInput struct {
	UserID int64
	Items  []CreateOrderItemInput
	Remark string
}

// CreateOrderItemInput 创建订单明细输入参数。
type CreateOrderItemInput struct {
	DishID   int64
	Quantity int
}

// Order 订单主表。
type Order struct {
	ID          int64     `json:"id"`
	OrderNo     string    `json:"order_no"`
	UserID      int64     `json:"user_id"`
	TotalAmount float64   `json:"total_amount"`
	Status      int       `json:"status"`
	Remark      string    `json:"remark"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// OrderItem 订单明细。
type OrderItem struct {
	ID             int64     `json:"id"`
	OrderID        int64     `json:"order_id"`
	DishID         int64     `json:"dish_id"`
	DishName       string    `json:"dish_name"`
	DishImageURL   string    `json:"dish_image_url"`
	UnitPrice      float64   `json:"unit_price"`
	Quantity       int       `json:"quantity"`
	SubtotalAmount float64   `json:"subtotal_amount"`
	CreatedAt      time.Time `json:"created_at"`
}

// OrderDetail 订单详情。
type OrderDetail struct {
	Order
	Items []OrderItem `json:"items"`
}

// AdminOrderListInput 管理端订单列表输入参数。
type AdminOrderListInput struct {
	Status   int
	Page     int
	PageSize int
}

// CanChangeStatus 判断订单状态是否允许流转。
func CanChangeStatus(from int, to int) bool {
	switch from {
	case consts.OrderStatusPending:
		return to == consts.OrderStatusCooking || to == consts.OrderStatusCanceled
	case consts.OrderStatusCooking:
		return to == consts.OrderStatusCompleted || to == consts.OrderStatusCanceled
	default:
		return false
	}
}

// CreateOrder 创建订单并保存菜品快照。
func CreateOrder(ctx context.Context, in CreateOrderInput) (*OrderDetail, error) {
	if in.UserID <= 0 {
		return nil, fmt.Errorf("用户ID不能为空")
	}
	if len(in.Items) == 0 {
		return nil, fmt.Errorf("订单明细不能为空")
	}

	var orderID int64
	err := g.DB().Transaction(ctx, func(ctx context.Context, tx gdb.TX) error {
		dishIDs := make([]int64, 0, len(in.Items))
		quantityMap := make(map[int64]int, len(in.Items))
		for _, item := range in.Items {
			if item.DishID <= 0 {
				return fmt.Errorf("菜品ID不能为空")
			}
			if err := menulogic.ValidateDishQuantity(item.Quantity); err != nil {
				return err
			}
			if _, exists := quantityMap[item.DishID]; !exists {
				dishIDs = append(dishIDs, item.DishID)
			}
			quantityMap[item.DishID] += item.Quantity
		}

		dishes, err := loadOrderDishes(ctx, tx, dishIDs)
		if err != nil {
			return err
		}
		if len(dishes) != len(dishIDs) {
			return fmt.Errorf("存在菜品不存在或已下架")
		}

		now := time.Now()
		orderNo := GenerateOrderNo(now)
		totalAmount := 0.0
		orderItems := make([]g.Map, 0, len(dishes))
		for _, dish := range dishes {
			quantity := quantityMap[dish.ID]
			subtotal := dish.Price * float64(quantity)
			totalAmount += subtotal
			orderItems = append(orderItems, g.Map{
				"dish_id":         dish.ID,
				"dish_name":       dish.Name,
				"dish_image_url":  dish.ImageURL,
				"unit_price":      dish.Price,
				"quantity":        quantity,
				"subtotal_amount": subtotal,
				"created_at":      now,
			})
		}

		id, err := tx.Model("orders").Ctx(ctx).Data(g.Map{
			"order_no":     orderNo,
			"user_id":      in.UserID,
			"total_amount": totalAmount,
			"status":       consts.OrderStatusPending,
			"remark":       in.Remark,
			"created_at":   now,
			"updated_at":   now,
		}).InsertAndGetId()
		if err != nil {
			return fmt.Errorf("创建订单失败: %w", err)
		}
		orderID = id

		for _, item := range orderItems {
			item["order_id"] = orderID
		}
		if _, err := tx.Model("order_items").Ctx(ctx).Data(orderItems).Insert(); err != nil {
			return fmt.Errorf("创建订单明细失败: %w", err)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return GetOrderDetail(ctx, orderID, in.UserID)
}

// ListUserOrders 查询用户订单列表。
func ListUserOrders(ctx context.Context, userID int64) ([]Order, error) {
	if userID <= 0 {
		return nil, fmt.Errorf("用户ID不能为空")
	}

	var list []Order
	err := g.DB().Model("orders").Ctx(ctx).
		Fields("id,order_no,user_id,total_amount,status,remark,created_at,updated_at").
		Where("user_id", userID).
		Order("created_at DESC,id DESC").
		Scan(&list)
	if err != nil {
		return nil, fmt.Errorf("查询我的订单失败: %w", err)
	}
	return list, nil
}

// ListAdminOrders 查询管理端订单列表。
func ListAdminOrders(ctx context.Context, in AdminOrderListInput) ([]Order, error) {
	if in.Page <= 0 {
		in.Page = 1
	}
	if in.PageSize <= 0 || in.PageSize > 100 {
		in.PageSize = 20
	}

	model := g.DB().Model("orders").Ctx(ctx).
		Fields("id,order_no,user_id,total_amount,status,remark,created_at,updated_at")
	if in.Status > 0 {
		model = model.Where("status", in.Status)
	}

	var list []Order
	err := model.Page(in.Page, in.PageSize).
		Order("created_at DESC,id DESC").
		Scan(&list)
	if err != nil {
		return nil, fmt.Errorf("查询订单列表失败: %w", err)
	}
	return list, nil
}

// GetOrderDetail 查询订单详情，可传 userID 限定只能查询自己的订单。
func GetOrderDetail(ctx context.Context, id int64, userID int64) (*OrderDetail, error) {
	orderInfo, err := getOrder(ctx, id, userID)
	if err != nil {
		return nil, err
	}

	items, err := listOrderItems(ctx, id)
	if err != nil {
		return nil, err
	}
	return &OrderDetail{Order: *orderInfo, Items: items}, nil
}

// CancelUserOrder 取消用户自己的待处理订单。
func CancelUserOrder(ctx context.Context, id int64, userID int64) error {
	orderInfo, err := getOrder(ctx, id, userID)
	if err != nil {
		return err
	}
	if orderInfo.Status != consts.OrderStatusPending {
		return fmt.Errorf("只能取消待处理订单")
	}
	return updateOrderStatus(ctx, id, consts.OrderStatusCanceled)
}

// ChangeOrderStatus 修改订单状态。
func ChangeOrderStatus(ctx context.Context, id int64, status int) error {
	orderInfo, err := getOrder(ctx, id, 0)
	if err != nil {
		return err
	}
	if !CanChangeStatus(orderInfo.Status, status) {
		return fmt.Errorf("订单状态流转冲突")
	}
	return updateOrderStatus(ctx, id, status)
}

// GenerateOrderNo 生成订单号。
func GenerateOrderNo(now time.Time) string {
	return fmt.Sprintf("%s%06d", now.Format("20060102150405"), grand.N(100000, 999999))
}

// orderDish 下单使用的菜品快照来源。
type orderDish struct {
	ID       int64   `json:"id"`
	Name     string  `json:"name"`
	ImageURL string  `json:"image_url"`
	Price    float64 `json:"price"`
}

// loadOrderDishes 查询上架菜品用于下单快照。
func loadOrderDishes(ctx context.Context, tx gdb.TX, dishIDs []int64) ([]orderDish, error) {
	var dishes []orderDish
	err := tx.Model("dishes").Ctx(ctx).
		Fields("id,name,image_url,price").
		WhereIn("id", dishIDs).
		Where("status", consts.DishStatusOn).
		Scan(&dishes)
	if err != nil {
		return nil, fmt.Errorf("查询下单菜品失败: %w", err)
	}
	return dishes, nil
}

// getOrder 查询订单主表。
func getOrder(ctx context.Context, id int64, userID int64) (*Order, error) {
	if id <= 0 {
		return nil, fmt.Errorf("订单ID不能为空")
	}
	model := g.DB().Model("orders").Ctx(ctx).
		Fields("id,order_no,user_id,total_amount,status,remark,created_at,updated_at").
		Where("id", id)
	if userID > 0 {
		model = model.Where("user_id", userID)
	}

	var orderInfo Order
	if err := model.Scan(&orderInfo); err != nil {
		return nil, fmt.Errorf("查询订单失败: %w", err)
	}
	if orderInfo.ID == 0 {
		return nil, fmt.Errorf("订单不存在")
	}
	return &orderInfo, nil
}

// listOrderItems 查询订单明细。
func listOrderItems(ctx context.Context, orderID int64) ([]OrderItem, error) {
	var items []OrderItem
	err := g.DB().Model("order_items").Ctx(ctx).
		Fields("id,order_id,dish_id,dish_name,dish_image_url,unit_price,quantity,subtotal_amount,created_at").
		Where("order_id", orderID).
		Order("id ASC").
		Scan(&items)
	if err != nil {
		return nil, fmt.Errorf("查询订单明细失败: %w", err)
	}
	return items, nil
}

// updateOrderStatus 更新订单状态。
func updateOrderStatus(ctx context.Context, id int64, status int) error {
	result, err := g.DB().Model("orders").Ctx(ctx).
		Where("id", id).
		Data(g.Map{"status": status, "updated_at": time.Now()}).
		Update()
	if err != nil {
		return fmt.Errorf("更新订单状态失败: %w", err)
	}
	rows, _ := result.RowsAffected()
	if rows == 0 {
		return fmt.Errorf("订单不存在")
	}
	return nil
}
