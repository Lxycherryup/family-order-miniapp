const { request } = require('../../utils/request')

const STATUS_TEXT = {
  1: '待处理',
  2: '制作中',
  3: '已完成',
  4: '已取消',
}

Page({
  data: {
    loading: false,
    orders: [],
  },

  // onShow 页面显示时加载订单。
  onShow() {
    this.loadOrders()
  },

  // onPullDownRefresh 下拉刷新订单。
  async onPullDownRefresh() {
    await this.loadOrders()
    wx.stopPullDownRefresh()
  },

  // loadOrders 加载我的订单列表。
  async loadOrders() {
    this.setData({ loading: true })
    try {
      const orders = await request({ url: '/orders' })
      this.setData({
        orders: orders.map((order) => ({
          ...order,
          status_text: STATUS_TEXT[order.status] || '未知',
          amount_text: Number(order.total_amount).toFixed(2),
        })),
        loading: false,
      })
    } catch (err) {
      this.setData({ loading: false })
      wx.showToast({ title: err.message || '加载订单失败', icon: 'none' })
    }
  },

  // goDetail 跳转订单详情。
  goDetail(event) {
    const id = event.currentTarget.dataset.id
    wx.navigateTo({ url: `/pages/order-detail/order-detail?id=${id}` })
  },
})
