const { request } = require('../../utils/request')
const { addDish, getCartSummary } = require('../../utils/cart')

Page({
  data: {
    loading: false,
    categories: [],
    dishes: [],
    activeCategoryID: 0,
    cartCount: 0,
    cartAmount: '0.00',
  },

  // onLoad 页面加载时初始化菜单。
  onLoad() {
    this.loadMenu()
  },

  // onShow 页面显示时刷新购物车汇总。
  onShow() {
    this.refreshCartSummary()
  },

  // onPullDownRefresh 下拉刷新菜单。
  async onPullDownRefresh() {
    await this.loadMenu()
    wx.stopPullDownRefresh()
  },

  // loadMenu 加载分类和菜品。
  async loadMenu() {
    this.setData({ loading: true })
    try {
      const categories = await request({ url: '/menu/categories' })
      const activeCategoryID = this.data.activeCategoryID || (categories[0] && categories[0].id) || 0
      const dishes = await this.loadDishes(activeCategoryID)
      this.setData({
        categories,
        activeCategoryID,
        dishes,
        loading: false,
      })
      this.refreshCartSummary()
    } catch (err) {
      this.setData({ loading: false })
      wx.showToast({ title: err.message || '加载菜单失败', icon: 'none' })
      if ((err.message || '').includes('登录')) {
        wx.redirectTo({ url: '/pages/login/login' })
      }
    }
  },

  // loadDishes 加载指定分类菜品。
  loadDishes(categoryID) {
    return request({
      url: '/menu/dishes',
      data: categoryID ? { category_id: categoryID } : {},
    })
  },

  // switchCategory 切换分类。
  async switchCategory(event) {
    const categoryID = Number(event.currentTarget.dataset.id)
    if (categoryID === this.data.activeCategoryID) {
      return
    }
    this.setData({ activeCategoryID: categoryID, loading: true })
    try {
      const dishes = await this.loadDishes(categoryID)
      this.setData({ dishes, loading: false })
    } catch (err) {
      this.setData({ loading: false })
      wx.showToast({ title: err.message || '加载菜品失败', icon: 'none' })
    }
  },

  // addToCart 添加菜品到购物车。
  addToCart(event) {
    const index = Number(event.currentTarget.dataset.index)
    const dish = this.data.dishes[index]
    if (!dish) {
      return
    }
    addDish(dish)
    this.refreshCartSummary()
    wx.showToast({ title: '已加入购物车', icon: 'success' })
  },

  // refreshCartSummary 刷新购物车汇总。
  refreshCartSummary() {
    const summary = getCartSummary()
    this.setData({
      cartCount: summary.count,
      cartAmount: summary.amount.toFixed(2),
    })
  },

  // goCart 跳转购物车。
  goCart() {
    wx.navigateTo({ url: '/pages/cart/cart' })
  },

  // goOrders 跳转我的订单。
  goOrders() {
    wx.navigateTo({ url: '/pages/orders/orders' })
  },
})
