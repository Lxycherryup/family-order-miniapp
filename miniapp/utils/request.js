const DEFAULT_BASE_URL = 'http://127.0.0.1:8000/api/miniapp'

// getBaseURL 获取接口基础地址。
function getBaseURL() {
  const app = getApp()
  return (app.globalData && app.globalData.apiBaseUrl) || DEFAULT_BASE_URL
}

// request 发起统一接口请求。
function request(options) {
  const token = wx.getStorageSync('token')
  return new Promise((resolve, reject) => {
    wx.request({
      url: `${getBaseURL()}${options.url}`,
      method: options.method || 'GET',
      data: options.data || {},
      header: {
        Authorization: token ? `Bearer ${token}` : '',
        'Content-Type': 'application/json',
      },
      success(res) {
        const body = res.data || {}
        if (body.code !== 0) {
          reject(new Error(body.msg || '请求失败'))
          return
        }
        resolve(body.data)
      },
      fail(err) {
        reject(err)
      },
    })
  })
}

module.exports = { request }
