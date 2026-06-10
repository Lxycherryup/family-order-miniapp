# 家庭微信小程序点餐系统设计文档

## 问题分析

本项目目标是实现一个仅供自己和家人使用的微信小程序点餐系统，不对外开放，不接入支付。系统需要支持家人在小程序端浏览菜单、加入购物车、提交订单，也需要管理员在 Vue 后台维护菜单、处理订单和管理白名单用户。

由于使用场景是家庭内部小工具，第一版应优先保证功能可用、权限清晰、部署简单，避免引入支付、配送、多门店、库存、优惠券等商业点餐系统复杂能力。

## 设计目标

- 支持微信小程序端点餐。
- 支持 Vue 管理后台维护菜单和订单。
- 后端使用 Golang 和 GoFrame 实现。
- 数据持久化使用 PostgreSQL。
- 第一阶段使用 Docker Compose 部署到腾讯云服务器。
- 第二阶段预留 K8s 部署能力。
- 系统仅允许白名单微信用户访问点餐能力。
- 不接入微信支付，不处理支付、退款和配送。

## 非目标范围

第一版暂不实现以下能力：

- 微信支付。
- 配送和自提排队。
- 打印机。
- 菜品库存。
- 多规格 SKU。
- 优惠券。
- 多门店。
- 商业报表。
- 复杂会员体系。
- 消息订阅通知。

## 总体方案

系统采用轻量生产版架构：

- 微信小程序：家人点餐、查看自己的订单。
- Vue 管理端：管理员维护分类、菜品、白名单和订单。
- GoFrame 后端：统一提供 RESTful API、登录鉴权和业务逻辑。
- PostgreSQL：持久化菜单、订单、用户和配置数据。
- Docker Compose：第一阶段部署。
- K8s：第二阶段迁移部署。

总体架构如下：

```text
微信小程序
  ├─ 家人登录
  ├─ 浏览菜单
  ├─ 加入购物车
  └─ 提交订单

Vue 管理后台
  ├─ 管理员登录
  ├─ 菜品分类管理
  ├─ 菜品管理
  ├─ 订单列表
  └─ 订单状态处理

GoFrame 后端服务
  ├─ 小程序认证模块
  ├─ 管理端认证模块
  ├─ 用户白名单模块
  ├─ 菜单分类模块
  ├─ 菜品模块
  ├─ 订单模块
  └─ 系统配置模块

PostgreSQL
  ├─ users
  ├─ admin_users
  ├─ menu_categories
  ├─ dishes
  ├─ orders
  ├─ order_items
  └─ system_configs
```

## 功能设计

### 小程序端

小程序端按常见微信点餐风格设计，核心是打开后能快速点餐。

页面：

```text
pages/
  ├─ login/               # 微信登录和白名单提示
  ├─ menu/                # 点餐首页
  ├─ dish-detail/         # 菜品详情
  ├─ cart/                # 购物车确认
  ├─ order-submit/        # 提交订单结果
  ├─ orders/              # 我的订单
  └─ order-detail/        # 订单详情
```

核心流程：

```text
打开小程序
  -> 微信登录
  -> 后端校验 OpenID 白名单
  -> 进入菜单页
  -> 选择分类和菜品
  -> 加入购物车
  -> 提交订单
  -> 查看订单状态
```

菜单页布局：

```text
顶部：家庭点餐 / 当前问候语
左侧：分类列表
右侧：菜品列表
底部：购物车栏 + 总价 + 去下单
```

状态处理：

- 未登录时跳转登录页。
- 未加入白名单时展示无权限提示。
- 菜品下架时由后端在提交订单时拦截。
- 订单提交失败时保留购物车内容，方便重试。
- 订单提交成功后清空购物车，并跳转订单详情。

### 管理端

管理端采用清晰、紧凑、易维护的后台系统风格。

推荐技术栈：

- Vue 3。
- Vite。
- TypeScript。
- Pinia。
- Vue Router。
- Axios。
- Element Plus。

页面：

```text
src/
  ├─ views/
  │  ├─ Login.vue
  │  ├─ Dashboard.vue
  │  ├─ CategoryList.vue
  │  ├─ DishList.vue
  │  ├─ DishForm.vue
  │  ├─ OrderList.vue
  │  ├─ OrderDetail.vue
  │  └─ UserWhitelist.vue
  ├─ router/
  ├─ stores/
  ├─ api/
  └─ components/
```

后台菜单：

- 首页概览：今日订单数、待处理订单数、菜品数量。
- 分类管理：新增、编辑、启用、禁用、排序。
- 菜品管理：新增、编辑、上下架、图片、价格、分类。
- 订单管理：查看订单、筛选状态、修改状态。
- 白名单管理：查看小程序用户、加入白名单、移出白名单、禁用用户。

## 后端设计

后端采用 GoFrame 单体服务。当前场景不拆微服务，避免增加部署、日志、链路追踪和数据库一致性成本。

目录结构：

```text
server/
  ├─ api/                 # 接口请求和响应定义
  │  ├─ admin/            # 管理端接口定义
  │  └─ miniapp/          # 小程序接口定义
  ├─ internal/
  │  ├─ cmd/              # 启动入口
  │  ├─ consts/           # 常量定义
  │  ├─ controller/       # 控制器层
  │  │  ├─ admin/
  │  │  └─ miniapp/
  │  ├─ dao/              # GoFrame 生成的 DAO
  │  ├─ logic/            # 业务逻辑层
  │  │  ├─ auth/
  │  │  ├─ menu/
  │  │  ├─ order/
  │  │  └─ user/
  │  ├─ middleware/       # 鉴权、日志、限流中间件
  │  ├─ model/            # 业务模型
  │  └─ service/          # 服务接口定义
  ├─ manifest/
  │  ├─ config/           # 配置文件
  │  ├─ docker/           # Dockerfile
  │  └─ deploy/           # Compose 和 K8s 部署文件
  ├─ resource/
  │  └─ public/           # 可选静态资源
  ├─ utility/             # 工具函数
  ├─ go.mod
  └─ main.go
```

分层原则：

- `controller` 只处理参数、调用服务和返回结果。
- `logic` 承担主要业务规则，例如下单校验、订单状态流转。
- `dao` 只负责数据库访问，不写复杂业务判断。
- `middleware` 统一处理 JWT、请求日志和错误响应。
- 所有新增 Go 类型和函数都需要中文注释。

## 数据库设计

### 用户表 `users`

用于记录家人微信登录身份和白名单状态。

| 字段 | 类型 | 说明 |
|---|---|---|
| id | bigserial | 主键 |
| openid | varchar(64) | 微信 OpenID，唯一 |
| nickname | varchar(64) | 昵称，可为空 |
| avatar_url | varchar(255) | 头像，可为空 |
| status | smallint | 状态：1 启用，2 禁用 |
| is_whitelist | boolean | 是否白名单 |
| last_login_at | timestamptz | 最近登录时间 |
| created_at | timestamptz | 创建时间 |
| updated_at | timestamptz | 更新时间 |

索引：

- `uk_users_openid(openid)`：微信登录查询使用。
- `idx_users_whitelist(is_whitelist, status)`：白名单管理筛选使用。

### 管理员表 `admin_users`

用于 Vue 管理端登录。

| 字段 | 类型 | 说明 |
|---|---|---|
| id | bigserial | 主键 |
| username | varchar(64) | 管理员用户名，唯一 |
| password_hash | varchar(255) | 密码哈希 |
| status | smallint | 状态：1 启用，2 禁用 |
| last_login_at | timestamptz | 最近登录时间 |
| created_at | timestamptz | 创建时间 |
| updated_at | timestamptz | 更新时间 |

索引：

- `uk_admin_users_username(username)`：管理员登录查询使用。

### 菜品分类表 `menu_categories`

| 字段 | 类型 | 说明 |
|---|---|---|
| id | bigserial | 主键 |
| name | varchar(64) | 分类名称 |
| sort | int | 排序值，越小越靠前 |
| status | smallint | 状态：1 启用，2 禁用 |
| created_at | timestamptz | 创建时间 |
| updated_at | timestamptz | 更新时间 |

索引：

- `idx_menu_categories_status_sort(status, sort)`：小程序菜单分类查询使用。

### 菜品表 `dishes`

| 字段 | 类型 | 说明 |
|---|---|---|
| id | bigserial | 主键 |
| category_id | bigint | 分类 ID |
| name | varchar(100) | 菜品名称 |
| description | varchar(500) | 菜品描述 |
| image_url | varchar(255) | 图片地址 |
| price | numeric(10,2) | 单价 |
| unit | varchar(16) | 单位，例如份、碗、杯 |
| status | smallint | 状态：1 上架，2 下架 |
| sort | int | 排序值 |
| created_at | timestamptz | 创建时间 |
| updated_at | timestamptz | 更新时间 |

索引：

- `idx_dishes_category_status_sort(category_id, status, sort)`：按分类查上架菜品使用。
- `idx_dishes_status(status)`：后台按上下架筛选使用。

### 订单表 `orders`

| 字段 | 类型 | 说明 |
|---|---|---|
| id | bigserial | 主键 |
| order_no | varchar(32) | 订单号，唯一 |
| user_id | bigint | 下单用户 ID |
| total_amount | numeric(10,2) | 订单总金额 |
| status | smallint | 状态：1 待处理，2 制作中，3 已完成，4 已取消 |
| remark | varchar(500) | 下单备注 |
| created_at | timestamptz | 创建时间 |
| updated_at | timestamptz | 更新时间 |

索引：

- `uk_orders_order_no(order_no)`：按订单号查询使用。
- `idx_orders_user_created(user_id, created_at desc)`：我的订单列表使用。
- `idx_orders_status_created(status, created_at desc)`：后台订单处理列表使用。

### 订单明细表 `order_items`

下单时保存菜品快照，避免后续菜品改价影响历史订单。

| 字段 | 类型 | 说明 |
|---|---|---|
| id | bigserial | 主键 |
| order_id | bigint | 订单 ID |
| dish_id | bigint | 菜品 ID |
| dish_name | varchar(100) | 下单时菜品名称 |
| dish_image_url | varchar(255) | 下单时菜品图片 |
| unit_price | numeric(10,2) | 下单时单价 |
| quantity | int | 数量 |
| subtotal_amount | numeric(10,2) | 小计金额 |
| created_at | timestamptz | 创建时间 |

索引：

- `idx_order_items_order_id(order_id)`：订单详情查询使用。

### 系统配置表 `system_configs`

用于保存低频配置，例如业务开关。敏感密钥不直接明文放数据库，优先使用环境变量或 K8s Secret。

| 字段 | 类型 | 说明 |
|---|---|---|
| id | bigserial | 主键 |
| config_key | varchar(100) | 配置键，唯一 |
| config_value | text | 配置值 |
| remark | varchar(255) | 备注 |
| created_at | timestamptz | 创建时间 |
| updated_at | timestamptz | 更新时间 |

索引：

- `uk_system_configs_key(config_key)`：按配置键读取使用。

### SQL 风险说明

- 菜单查询使用 `status + sort` 索引，第一版数据量较小，全表扫描风险低。
- 菜品列表使用 `category_id + status + sort` 索引，适合小程序按分类查询。
- 订单列表使用 `user_id + created_at` 和 `status + created_at` 索引，能支撑后续订单增长。
- 下单必须使用事务，订单主表和订单明细需要同时成功或同时失败。
- 订单明细保存菜品快照，避免菜单修改影响历史订单。
- 删除分类前需要检查是否存在菜品，避免产生孤儿数据。
- 菜品建议优先下架或软删除，不建议物理删除。

## API 设计

接口统一前缀：

```text
/api/miniapp/*
/api/admin/*
```

统一响应结构：

```json
{
  "code": 0,
  "msg": "success",
  "data": {}
}
```

### 小程序认证

```http
POST /api/miniapp/auth/login
```

请求：

```json
{
  "code": "微信临时登录凭证",
  "nickname": "昵称",
  "avatar_url": "头像地址"
}
```

响应：

```json
{
  "token": "JWT",
  "is_whitelist": true
}
```

逻辑：

- 后端使用 `code` 调用微信 `jscode2session` 获取 OpenID。
- 用户不存在时自动创建。
- OpenID 不在白名单时允许登录，但禁止访问菜单和下单能力。
- 白名单用户返回 JWT。

### 小程序菜单

```http
GET /api/miniapp/menu/categories
GET /api/miniapp/menu/dishes?category_id=1
GET /api/miniapp/menu/dishes/{id}
```

逻辑：

- 只返回启用分类和上架菜品。
- 菜品列表按 `sort` 排序。
- 不返回后台内部字段。

### 小程序订单

```http
POST /api/miniapp/orders
GET /api/miniapp/orders
GET /api/miniapp/orders/{id}
POST /api/miniapp/orders/{id}/cancel
```

下单请求：

```json
{
  "items": [
    {
      "dish_id": 1,
      "quantity": 2
    }
  ],
  "remark": "少辣"
}
```

下单逻辑：

- 校验用户是否在白名单中。
- 校验菜品存在且已上架。
- 以数据库当前价格计算金额，不信任前端传价。
- 订单和订单明细写入同一个事务。
- 保存菜品名称、图片、单价快照。
- 用户只允许取消自己的待处理订单。

### 管理员认证

```http
POST /api/admin/auth/login
POST /api/admin/auth/logout
GET /api/admin/auth/profile
```

### 白名单用户

```http
GET /api/admin/users
PUT /api/admin/users/{id}/whitelist
PUT /api/admin/users/{id}/status
```

用途：

- 查看登录过的小程序用户。
- 设置是否加入白名单。
- 禁用异常用户。

### 分类管理

```http
GET /api/admin/categories
POST /api/admin/categories
PUT /api/admin/categories/{id}
DELETE /api/admin/categories/{id}
PUT /api/admin/categories/{id}/status
```

规则：

- 分类下存在菜品时不允许删除。
- 删除前需要先移动或下架菜品。

### 菜品管理

```http
GET /api/admin/dishes
POST /api/admin/dishes
GET /api/admin/dishes/{id}
PUT /api/admin/dishes/{id}
DELETE /api/admin/dishes/{id}
PUT /api/admin/dishes/{id}/status
```

规则：

- 第一版建议使用下架或软删除。
- 历史订单继续展示下单快照。

### 订单管理

```http
GET /api/admin/orders
GET /api/admin/orders/{id}
PUT /api/admin/orders/{id}/status
```

订单状态：

```text
1 待处理
2 制作中
3 已完成
4 已取消
```

允许状态流转：

```text
待处理 -> 制作中
待处理 -> 已取消
制作中 -> 已完成
制作中 -> 已取消
```

不允许状态流转：

```text
已完成 -> 其他状态
已取消 -> 其他状态
```

## 错误处理设计

错误码：

| code | 说明 |
|---|---|
| 0 | 成功 |
| 40001 | 参数错误 |
| 40101 | 未登录或 Token 失效 |
| 40301 | 不在白名单 |
| 40401 | 资源不存在 |
| 40901 | 状态冲突 |
| 50001 | 系统错误 |

错误响应示例：

```json
{
  "code": 40301,
  "msg": "当前微信用户未加入点餐白名单",
  "data": null
}
```

Go 错误处理要求：

- 不使用 `panic(err)` 处理业务错误。
- 错误信息需要明确、可定位、可追踪。
- 底层错误使用 `%w` 包装，保留错误链路。

## 日志设计

每次请求记录：

- 请求 ID。
- 用户 ID 或管理员 ID。
- 请求路径。
- 请求方法。
- 响应状态。
- 耗时。
- 错误信息。
- 关键参数。

禁止记录：

- JWT。
- 微信 `session_key`。
- 密码。
- Token。
- 身份证号等敏感信息。

## 认证与权限设计

小程序端：

- 使用微信登录。
- 后端通过 `jscode2session` 获取 OpenID。
- OpenID 写入用户表。
- 点餐接口必须校验白名单。
- 小程序 JWT 建议有效期为 7 天。

管理端：

- 使用账号密码登录。
- 密码使用 bcrypt 哈希存储。
- 管理端 JWT 建议有效期为 24 小时。
- 管理接口必须校验管理员登录态。

## Redis 设计

第一版 Redis 不作为强依赖，仅预留能力。

未来可用于：

- 登录限流。
- 热门菜单缓存。
- 管理端接口防刷。

Key 设计示例：

```text
family_order:rate_limit:login:{ip}
family_order:menu:categories
family_order:menu:dishes:{category_id}
```

缓存风险：

- 缓存击穿：菜单缓存失效时可能集中访问数据库，后续可用互斥锁或短期本地缓存缓解。
- 缓存穿透：对不存在的分类 ID 可缓存空结果，并设置较短过期时间。
- 缓存雪崩：不同缓存 Key 设置随机过期时间，避免同一时间大量失效。

第一版菜单数据量很小，可以先直接查询数据库，避免过早增加缓存复杂度。

## 部署设计

### Docker Compose 部署

第一阶段在腾讯云服务器使用 Docker Compose 部署，结构简单、排障方便。

```text
腾讯云服务器
  ├─ Nginx
  │  ├─ admin.example.com       -> Vue 管理端静态资源
  │  └─ api.example.com         -> GoFrame API
  ├─ GoFrame 后端容器
  ├─ PostgreSQL 容器或腾讯云数据库
  └─ 可选 Redis 容器
```

推荐域名规划：

```text
api.yourdomain.com      # 后端 API
admin.yourdomain.com    # Vue 管理端
```

小程序后台配置：

```text
request 合法域名：https://api.yourdomain.com
```

Docker Compose 服务：

```text
services:
  backend       # GoFrame API
  admin-web     # Vue 静态资源，也可直接由 Nginx 挂载
  postgres      # PostgreSQL
  redis         # 可选，第一版可以不启用
  nginx         # 反向代理和 HTTPS
```

### K8s 部署

第二阶段迁移到 K8s，第一版仅预留部署文件结构。

```text
namespace: family-order
  ├─ Deployment/backend
  ├─ Service/backend
  ├─ Deployment/admin-web
  ├─ Service/admin-web
  ├─ StatefulSet/postgres 或外部云数据库
  ├─ Secret
  ├─ ConfigMap
  └─ Ingress
```

敏感配置：

- 微信 AppSecret 使用环境变量或 K8s Secret。
- 数据库密码使用环境变量或 K8s Secret。
- JWT 密钥使用环境变量或 K8s Secret。

## 安全设计

- 小程序端使用微信登录，后端以 OpenID 识别用户。
- 点餐接口必须校验白名单。
- 管理端密码使用 bcrypt 哈希。
- HTTPS 必须开启，小程序正式环境要求 HTTPS。
- 微信 `session_key` 不写日志。
- 后端日志不打印密码、Token 和敏感信息。
- 数据库不暴露公网，只允许容器网络或服务器本机访问。
- Nginx 配置请求体大小限制和超时时间。
- 管理端登录接口需要预留限流能力。

## 测试设计

后端优先覆盖：

- 微信登录 OpenID 处理逻辑，可使用接口 mock。
- 白名单校验。
- 菜单查询。
- 下单事务。
- 订单金额计算。
- 订单状态流转。
- 管理员登录。

前端优先覆盖：

- 登录态过期后跳转。
- 菜品加入购物车。
- 下单失败不清空购物车。
- 后台订单状态切换。

部署验收：

- 后端健康检查接口可访问。
- 管理端静态页面可访问。
- 小程序合法域名配置后可请求 API。
- HTTPS 证书有效。
- 数据库连接不暴露公网。

## 风险评估

- 微信登录依赖微信接口，开发阶段需要准备小程序 AppID 和 AppSecret。
- 小程序正式环境要求 HTTPS，需要提前配置域名证书。
- 如果使用服务器内 PostgreSQL 容器，需要做好数据卷挂载和备份。
- 如果后续切换 K8s，需要重新处理 Secret、Ingress 和数据持久化。
- 菜品图片存储方案尚未最终确定，第一版可先使用对象存储或静态资源地址。
- 白名单 OpenID 需要用户至少登录一次后才能在后台加入，初始化流程需要在管理端提示清楚。

## 后续建议

- 第一阶段先完成后端、数据库、管理端和小程序核心闭环。
- 菜品图片建议使用腾讯云 COS，避免后端容器本地存储带来的迁移问题。
- PostgreSQL 生产数据需要定时备份。
- 管理端后续可增加简单统计，例如今日订单数和常点菜品。
- 如果家庭使用频率较高，再考虑 Redis 缓存和消息通知。
