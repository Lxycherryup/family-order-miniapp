const { request } = require('../../utils/request')

const STATUS_TEXT = {
  1: '待处理',
  2: '制作中',
  3: '已完成',
  4: '已取消',
}

Page({
  data: {
    id: 0,
    loading: false,
    order: null,
  },

  // onLoad 页面加载时读取订单ID。
  onLoad(options) {
    this.setData({ id: Number(options.id || 0) })
    this.loadDetail()
  },

  // loadDetail 加载订单详情。
  async loadDetail() {
    if (!this.data.id) {
      return
    }

    this.setData({ loading: true })
    try {
      const order = await request({ url: `/orders/${this.data.id}` })
      this.setData({
        order: {
          ...order,
          status_text: STATUS_TEXT[order.status] || '未知',
          amount_text: Number(order.total_amount).toFixed(2),
          items: order.items.map((item) => ({
            ...item,
            unit_price_text: Number(item.unit_price).toFixed(2),
            subtotal_text: Number(item.subtotal_amount).toFixed(2),
          })),
        },
        loading: false,
      })
    } catch (err) {
      this.setData({ loading: false })
      wx.showToast({ title: err.message || '加载订单详情失败', icon: 'none' })
    }
  },

  // cancelOrder 取消待处理订单。
  cancelOrder() {
    wx.showModal({
      title: '取消订单',
      content: '确认取消这个订单？',
      success: async (res) => {
        if (!res.confirm) {
          return
        }
        try {
          await request({
            url: `/orders/${this.data.id}/cancel`,
            method: 'POST',
          })
          wx.showToast({ title: '订单已取消', icon: 'success' })
          this.loadDetail()
        } catch (err) {
          wx.showToast({ title: err.message || '取消订单失败', icon: 'none' })
        }
      },
    })
  },
})
