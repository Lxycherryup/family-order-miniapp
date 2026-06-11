const { request } = require('../../utils/request')

Page({
  data: {
    loading: false,
    message: '',
  },

  // onLoad 页面加载时检查本地登录态。
  onLoad() {
    const token = wx.getStorageSync('token')
    if (token) {
      wx.redirectTo({ url: '/pages/menu/menu' })
    }
  },

  // handleLogin 发起微信登录。
  handleLogin() {
    if (this.data.loading) {
      return
    }

    this.setData({ loading: true, message: '' })
    wx.login({
      success: async (loginRes) => {
        try {
          const data = await request({
            url: '/auth/login',
            method: 'POST',
            data: {
              code: loginRes.code,
            },
          })
          if (!data.is_whitelist) {
            this.setData({
              message: '当前账号暂未开通点餐权限，请联系管理员',
              loading: false,
            })
            return
          }
          wx.setStorageSync('token', data.token)
          wx.redirectTo({ url: '/pages/menu/menu' })
        } catch (err) {
          this.setData({
            message: err.message || '登录失败，请稍后重试',
            loading: false,
          })
        }
      },
      fail: () => {
        this.setData({
          message: '获取微信登录凭证失败',
          loading: false,
        })
      },
    })
  },
})
