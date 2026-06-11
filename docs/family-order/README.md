# 家庭微信小程序点餐系统

本项目用于实现仅供自己和家人使用的微信小程序点餐系统。

## 应用组成

- `server/`：GoFrame 后端服务。
- `admin-web/`：Vue 管理后台。
- `miniapp/`：微信小程序端。
- `deploy/`：Nginx、Docker Compose 和 K8s 部署配置。

## 第一版范围

- 菜品分类管理。
- 菜品管理。
- 微信 OpenID 白名单。
- 小程序点餐。
- 订单状态处理。
- Docker Compose 部署。

## 非目标范围

第一版不支持支付、配送、打印机、多门店、库存、优惠券和商业报表。

## 本地开发登录

本地 Docker Compose 镜像会自动开启小程序 mock 登录：

```yaml
wechat:
  mockEnabled: true
  mockOpenID: "dev-family-user"
```

小程序调用 `POST /api/miniapp/auth/login` 时，后端不会请求微信服务器，而是使用固定 OpenID `dev-family-user` 创建或查询用户。管理员在后台白名单页面把该用户加入白名单后，小程序即可进入点餐流程。

## 生产微信登录

生产环境必须关闭 mock 登录，并配置真实微信小程序信息：

```yaml
wechat:
  mockEnabled: false
  appId: "真实小程序AppID"
  appSecret: "真实小程序AppSecret"
```

上线后家人不需要输入账号密码。小程序通过 `wx.login()` 获取 code，后端使用真实 `appId/appSecret` 换取 OpenID，再根据 OpenID 判断是否在白名单中：

- 白名单用户：返回小程序 Token，可以点餐。
- 非白名单用户：返回未开通提示，管理员可在后台加入白名单。

禁止在生产环境开启 `mockEnabled`，否则所有用户都会被识别为同一个 mock 用户，订单归属和白名单控制都会失效。
