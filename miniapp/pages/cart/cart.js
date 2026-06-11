const { request } = require('../../utils/request')
const { clearCart, getCart, getCartSummary, updateQuantity } = require('../../utils/cart')

Page({
  data: {
    items: [],
    count: 0,
    amount: '0.00',
    remark: '',
    submitting: false,
  },

  // onShow 页面显示时加载购物车。
  onShow() {
    this.loadCart()
  },

  // loadCart 加载购物车数据。
  loadCart() {
    const summary = getCartSummary()
    this.setData({
      items: getCart(),
      count: summary.count,
      amount: summary.amount.toFixed(2),
    })
  },

  // increase 增加菜品数量。
  increase(event) {
    const dishID = Number(event.currentTarget.dataset.id)
    const item = this.data.items.find((current) => current.dish_id === dishID)
    if (!item) {
      return
    }
    updateQuantity(dishID, item.quantity + 1)
    this.loadCart()
  },

  // decrease 减少菜品数量。
  decrease(event) {
    const dishID = Number(event.currentTarget.dataset.id)
    const item = this.data.items.find((current) => current.dish_id === dishID)
    if (!item) {
      return
    }
    updateQuantity(dishID, item.quantity - 1)
    this.loadCart()
  },

  // onRemarkInput 更新订单备注。
  onRemarkInput(event) {
    this.setData({ remark: event.detail.value })
  },

  // submitOrder 提交订单。
  async submitOrder() {
    if (this.data.submitting || this.data.items.length === 0) {
      return
    }

    this.setData({ submitting: true })
    try {
      const order = await request({
        url: '/orders',
        method: 'POST',
        data: {
          remark: this.data.remark,
          items: this.data.items.map((item) => ({
            dish_id: item.dish_id,
            quantity: item.quantity,
          })),
        },
      })
      clearCart()
      wx.redirectTo({ url: `/pages/order-detail/order-detail?id=${order.id}` })
    } catch (err) {
      wx.showToast({ title: err.message || '提交订单失败', icon: 'none' })
      this.setData({ submitting: false })
    }
  },

  // goMenu 返回菜单页。
  goMenu() {
    wx.redirectTo({ url: '/pages/menu/menu' })
  },
})
