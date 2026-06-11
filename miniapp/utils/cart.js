const CART_KEY = 'cart_items'

// getCart 获取购物车列表。
function getCart() {
  return wx.getStorageSync(CART_KEY) || []
}

// saveCart 保存购物车列表。
function saveCart(items) {
  wx.setStorageSync(CART_KEY, items)
}

// clearCart 清空购物车。
function clearCart() {
  wx.removeStorageSync(CART_KEY)
}

// addDish 添加菜品到购物车。
function addDish(dish) {
  const items = getCart()
  const current = items.find((item) => item.dish_id === dish.id)
  if (current) {
    current.quantity += 1
  } else {
    items.push({
      dish_id: dish.id,
      name: dish.name,
      image_url: dish.image_url,
      price: Number(dish.price),
      quantity: 1,
    })
  }
  saveCart(items)
  return items
}

// updateQuantity 更新购物车菜品数量。
function updateQuantity(dishID, quantity) {
  const items = getCart()
    .map((item) => (item.dish_id === dishID ? { ...item, quantity } : item))
    .filter((item) => item.quantity > 0)
  saveCart(items)
  return items
}

// getCartSummary 计算购物车汇总信息。
function getCartSummary() {
  const items = getCart()
  return items.reduce(
    (summary, item) => ({
      count: summary.count + item.quantity,
      amount: summary.amount + Number(item.price) * item.quantity,
    }),
    { count: 0, amount: 0 },
  )
}

module.exports = {
  addDish,
  clearCart,
  getCart,
  getCartSummary,
  saveCart,
  updateQuantity,
}
