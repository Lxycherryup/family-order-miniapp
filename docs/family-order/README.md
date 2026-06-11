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

## 本地部署

启动后端、PostgreSQL 和本地 Nginx：

```bash
docker compose -f deploy/compose/docker-compose.yaml up -d
```

常用访问地址：

- 后端健康检查：`http://127.0.0.1:8000/health`
- 管理后台开发服务：`http://localhost:5174/`
- PostgreSQL：`127.0.0.1:5432`

Navicat 连接信息：

- 数据库：`family_order`
- 用户名：`family_order`
- 密码：`family_order`

管理后台默认开发账号：

```json
{
  "username": "admin",
  "password": "admin123456"
}
```

## 后端验证

```bash
cd server
go test -count=1 ./...
curl http://127.0.0.1:8000/health
```

## K8s 说明

`deploy/k8s/` 仅作为第二阶段部署模板。生产使用前必须替换以下内容：

- 镜像地址：`your-registry.example.com`
- 域名：`api.yourdomain.com`、`admin.yourdomain.com`
- TLS Secret：`family-order-tls`
- `secret.example.yaml` 中的数据库、JWT 和微信密钥

K8s 模板默认生产环境关闭小程序 mock 登录。

## 最终验收清单

- 后端 `go test -count=1 ./...` 通过。
- 后端 `/health` 可访问。
- PostgreSQL 初始化表结构成功。
- 管理员可以登录管理后台。
- 管理员可以新增分类和菜品。
- 管理员可以把小程序用户加入白名单。
- 白名单用户可以在小程序端看到菜单并提交订单。
- 后台可以查看订单详情并修改订单状态。
- 非白名单用户无法访问菜单和下单。
- Docker Compose 配置可解析。
- Nginx 可以代理 `/health` 和 `/api/`。
