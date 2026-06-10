# 家庭微信小程序点餐系统 Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use `superpowers:subagent-driven-development`（推荐）或 `superpowers:executing-plans` 按任务逐项实现。步骤使用 checkbox（`- [ ]`）语法追踪进度。

**Goal:** 构建一个仅供自己和家人使用的微信小程序点餐系统，包含 GoFrame 后端、PostgreSQL 数据库、Vue 管理端、小程序端、Docker Compose 部署和 K8s 预留配置。

**Architecture:** 项目采用单仓库多应用结构：`server/` 提供 GoFrame RESTful API，`admin-web/` 提供 Vue 管理后台，`miniapp/` 提供微信小程序端，`deploy/` 保存部署配置。后端采用单体分层架构，所有点餐接口通过微信 OpenID 白名单控制访问，管理端通过账号密码和 JWT 鉴权。

**Tech Stack:** Golang 1.23、GoFrame v2、PostgreSQL、Vue 3、Vite、TypeScript、Pinia、Vue Router、Axios、Element Plus、微信小程序原生框架、Docker、Docker Compose、K8s。

---

## 文件结构规划

### 后端 `server/`

```text
server/
  ├─ api/
  │  ├─ admin/
  │  │  ├─ auth.go
  │  │  ├─ category.go
  │  │  ├─ dish.go
  │  │  ├─ order.go
  │  │  └─ user.go
  │  └─ miniapp/
  │     ├─ auth.go
  │     ├─ menu.go
  │     └─ order.go
  ├─ internal/
  │  ├─ cmd/
  │  │  └─ cmd.go
  │  ├─ consts/
  │  │  ├─ error_code.go
  │  │  └─ status.go
  │  ├─ controller/
  │  │  ├─ admin/
  │  │  └─ miniapp/
  │  ├─ dao/
  │  ├─ logic/
  │  │  ├─ auth/
  │  │  ├─ menu/
  │  │  ├─ order/
  │  │  └─ user/
  │  ├─ middleware/
  │  │  ├─ auth.go
  │  │  ├─ response.go
  │  │  └─ request_log.go
  │  ├─ model/
  │  │  ├─ entity/
  │  │  └─ input/
  │  └─ service/
  ├─ manifest/
  │  ├─ config/
  │  │  ├─ config.example.yaml
  │  │  └─ config.yaml
  │  ├─ deploy/
  │  │  ├─ migrations/
  │  │  │  └─ 001_init.sql
  │  │  ├─ docker-compose.yaml
  │  │  └─ k8s/
  │  └─ docker/
  │     └─ Dockerfile
  ├─ test/
  │  ├─ auth_test.go
  │  ├─ menu_test.go
  │  └─ order_test.go
  ├─ go.mod
  ├─ go.sum
  └─ main.go
```

### 管理端 `admin-web/`

```text
admin-web/
  ├─ src/
  │  ├─ api/
  │  │  ├─ auth.ts
  │  │  ├─ category.ts
  │  │  ├─ dish.ts
  │  │  ├─ order.ts
  │  │  └─ user.ts
  │  ├─ components/
  │  │  └─ PageHeader.vue
  │  ├─ router/
  │  │  └─ index.ts
  │  ├─ stores/
  │  │  └─ auth.ts
  │  ├─ views/
  │  │  ├─ Login.vue
  │  │  ├─ Dashboard.vue
  │  │  ├─ CategoryList.vue
  │  │  ├─ DishList.vue
  │  │  ├─ DishForm.vue
  │  │  ├─ OrderList.vue
  │  │  ├─ OrderDetail.vue
  │  │  └─ UserWhitelist.vue
  │  ├─ App.vue
  │  └─ main.ts
  ├─ index.html
  ├─ package.json
  ├─ tsconfig.json
  └─ vite.config.ts
```

### 小程序 `miniapp/`

```text
miniapp/
  ├─ app.js
  ├─ app.json
  ├─ app.wxss
  ├─ project.config.json
  ├─ utils/
  │  ├─ request.js
  │  └─ cart.js
  └─ pages/
     ├─ login/
     ├─ menu/
     ├─ dish-detail/
     ├─ cart/
     ├─ order-submit/
     ├─ orders/
     └─ order-detail/
```

### 部署 `deploy/`

```text
deploy/
  ├─ nginx/
  │  └─ family-order.conf
  ├─ compose/
  │  └─ docker-compose.yaml
  └─ k8s/
     ├─ namespace.yaml
     ├─ backend-deployment.yaml
     ├─ admin-web-deployment.yaml
     ├─ configmap.yaml
     ├─ secret.example.yaml
     └─ ingress.yaml
```

---

## Task 1: 创建项目骨架和基础说明

**Files:**
- Create: `server/.gitkeep`
- Create: `admin-web/.gitkeep`
- Create: `miniapp/.gitkeep`
- Create: `deploy/.gitkeep`
- Create: `docs/family-order/README.md`
- Modify: `.gitignore`

- [x] **Step 1: 创建目录**

Run:

```bash
mkdir -p server admin-web miniapp deploy docs/family-order
touch server/.gitkeep admin-web/.gitkeep miniapp/.gitkeep deploy/.gitkeep
```

Expected: 四个应用目录创建成功。

- [x] **Step 2: 更新 `.gitignore`**

在 `.gitignore` 增加以下内容：

```gitignore
# 家庭点餐系统
server/manifest/config/config.yaml
server/tmp/
server/log/
admin-web/node_modules/
admin-web/dist/
miniapp/miniprogram_npm/
deploy/**/*.local.yaml
```

- [x] **Step 3: 写项目说明**

创建 `docs/family-order/README.md`：

```markdown
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
```

- [x] **Step 4: 提交骨架**

Run:

```bash
git add .gitignore server admin-web miniapp deploy docs/family-order/README.md
git commit -m "chore: 初始化家庭点餐项目骨架"
```

Expected: 提交成功，只包含骨架和说明文件。

---

## Task 2: 初始化 GoFrame 后端和健康检查

**Files:**
- Create: `server/go.mod`
- Create: `server/main.go`
- Create: `server/internal/cmd/cmd.go`
- Create: `server/internal/controller/health.go`
- Create: `server/manifest/config/config.example.yaml`
- Create: `server/test/health_test.go`

- [x] **Step 1: 初始化 Go 模块**

Run:

```bash
cd server
go mod init family-order/server
go get github.com/gogf/gf/v2@latest
```

Expected: `server/go.mod` 和 `server/go.sum` 生成成功，`go.mod` 中 Go 版本为本机 `go env GOVERSION` 对应主版本，不手动升级。

- [x] **Step 2: 编写健康检查测试**

创建 `server/test/health_test.go`：

```go
package test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
)

// TestHealthHandler 验证健康检查接口返回成功。
func TestHealthHandler(t *testing.T) {
	s := g.Server()
	s.BindHandler("/health", func(r *ghttp.Request) {
		r.Response.WriteJson(g.Map{
			"code": 0,
			"msg":  "success",
			"data": g.Map{"status": "ok"},
		})
	})

	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	w := httptest.NewRecorder()
	s.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("期望状态码 200，实际为 %d", w.Code)
	}
}
```

- [x] **Step 3: 运行测试并确认失败或编译提示**

Run:

```bash
cd server
go test ./...
```

Expected: 如果缺少导入或初始化代码，测试失败；记录失败信息后进入实现。

- [x] **Step 4: 实现后端入口**

`server/main.go`：

```go
package main

import (
	_ "family-order/server/internal/cmd"

	"github.com/gogf/gf/v2/os/gctx"
	"github.com/gogf/gf/v2/os/gcmd"
)

// main 后端服务启动入口。
func main() {
	gcmd.Main.Run(gctx.GetInitCtx())
}
```

`server/internal/cmd/cmd.go`：

```go
package cmd

import (
	"context"

	"family-order/server/internal/controller"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
	"github.com/gogf/gf/v2/os/gcmd"
)

// Main 后端主命令。
var Main = gcmd.Command{
	Name:  "main",
	Usage: "main",
	Brief: "启动家庭点餐后端服务",
	Func: func(ctx context.Context, parser *gcmd.Parser) error {
		s := g.Server()
		s.Group("/", func(group *ghttp.RouterGroup) {
			group.Bind(controller.NewHealth())
		})
		s.Run()
		return nil
	},
}
```

`server/internal/controller/health.go`：

```go
package controller

import (
	"context"

	"github.com/gogf/gf/v2/frame/g"
)

// Health 健康检查控制器。
type Health struct{}

// NewHealth 创建健康检查控制器。
func NewHealth() *Health {
	return &Health{}
}

// HealthReq 健康检查请求。
type HealthReq struct {
	g.Meta `path:"/health" method:"get" tags:"系统" summary:"健康检查"`
}

// HealthRes 健康检查响应。
type HealthRes struct {
	Status string `json:"status"`
}

// Health 返回服务健康状态。
func (c *Health) Health(ctx context.Context, req *HealthReq) (res *HealthRes, err error) {
	return &HealthRes{Status: "ok"}, nil
}
```

- [x] **Step 5: 运行测试**

Run:

```bash
cd server
gofmt -w main.go internal test
go test ./...
```

Expected: `go test ./...` 通过。

- [x] **Step 6: 提交后端骨架**

Run:

```bash
git add server
git commit -m "feat: 初始化GoFrame后端服务"
```

Expected: 提交成功。

---

## Task 3: 增加数据库迁移和基础配置

**Files:**
- Create: `server/manifest/deploy/migrations/001_init.sql`
- Create: `server/manifest/config/config.example.yaml`
- Create: `server/internal/consts/status.go`
- Create: `server/test/migration_test.go`

- [x] **Step 1: 编写迁移文件测试**

创建 `server/test/migration_test.go`：

```go
package test

import (
	"os"
	"strings"
	"testing"
)

// TestInitMigrationContainsCoreTables 验证初始化迁移包含核心表。
func TestInitMigrationContainsCoreTables(t *testing.T) {
	content, err := os.ReadFile("../manifest/deploy/migrations/001_init.sql")
	if err != nil {
		t.Fatalf("读取迁移文件失败: %v", err)
	}
	sql := string(content)
	tables := []string{
		"CREATE TABLE IF NOT EXISTS users",
		"CREATE TABLE IF NOT EXISTS admin_users",
		"CREATE TABLE IF NOT EXISTS menu_categories",
		"CREATE TABLE IF NOT EXISTS dishes",
		"CREATE TABLE IF NOT EXISTS orders",
		"CREATE TABLE IF NOT EXISTS order_items",
		"CREATE TABLE IF NOT EXISTS system_configs",
	}
	for _, table := range tables {
		if !strings.Contains(sql, table) {
			t.Fatalf("迁移文件缺少表定义: %s", table)
		}
	}
}
```

- [x] **Step 2: 运行测试确认失败**

Run:

```bash
cd server
go test ./test -run TestInitMigrationContainsCoreTables -v
```

Expected: 因迁移文件不存在而失败。

- [x] **Step 3: 编写初始化 SQL**

创建 `server/manifest/deploy/migrations/001_init.sql`，包含 7 张表、索引和中文注释：

```sql
CREATE TABLE IF NOT EXISTS users (
    id BIGSERIAL PRIMARY KEY,
    openid VARCHAR(64) NOT NULL,
    nickname VARCHAR(64) NOT NULL DEFAULT '',
    avatar_url VARCHAR(255) NOT NULL DEFAULT '',
    status SMALLINT NOT NULL DEFAULT 1,
    is_whitelist BOOLEAN NOT NULL DEFAULT FALSE,
    last_login_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT uk_users_openid UNIQUE (openid)
);
COMMENT ON TABLE users IS '小程序用户表';

CREATE INDEX IF NOT EXISTS idx_users_whitelist ON users (is_whitelist, status);

CREATE TABLE IF NOT EXISTS admin_users (
    id BIGSERIAL PRIMARY KEY,
    username VARCHAR(64) NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    status SMALLINT NOT NULL DEFAULT 1,
    last_login_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT uk_admin_users_username UNIQUE (username)
);
COMMENT ON TABLE admin_users IS '管理端用户表';

CREATE TABLE IF NOT EXISTS menu_categories (
    id BIGSERIAL PRIMARY KEY,
    name VARCHAR(64) NOT NULL,
    sort INT NOT NULL DEFAULT 100,
    status SMALLINT NOT NULL DEFAULT 1,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
COMMENT ON TABLE menu_categories IS '菜品分类表';
CREATE INDEX IF NOT EXISTS idx_menu_categories_status_sort ON menu_categories (status, sort);

CREATE TABLE IF NOT EXISTS dishes (
    id BIGSERIAL PRIMARY KEY,
    category_id BIGINT NOT NULL,
    name VARCHAR(100) NOT NULL,
    description VARCHAR(500) NOT NULL DEFAULT '',
    image_url VARCHAR(255) NOT NULL DEFAULT '',
    price NUMERIC(10,2) NOT NULL,
    unit VARCHAR(16) NOT NULL DEFAULT '份',
    status SMALLINT NOT NULL DEFAULT 1,
    sort INT NOT NULL DEFAULT 100,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT fk_dishes_category_id FOREIGN KEY (category_id) REFERENCES menu_categories(id)
);
COMMENT ON TABLE dishes IS '菜品表';
CREATE INDEX IF NOT EXISTS idx_dishes_category_status_sort ON dishes (category_id, status, sort);
CREATE INDEX IF NOT EXISTS idx_dishes_status ON dishes (status);

CREATE TABLE IF NOT EXISTS orders (
    id BIGSERIAL PRIMARY KEY,
    order_no VARCHAR(32) NOT NULL,
    user_id BIGINT NOT NULL,
    total_amount NUMERIC(10,2) NOT NULL,
    status SMALLINT NOT NULL DEFAULT 1,
    remark VARCHAR(500) NOT NULL DEFAULT '',
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT uk_orders_order_no UNIQUE (order_no),
    CONSTRAINT fk_orders_user_id FOREIGN KEY (user_id) REFERENCES users(id)
);
COMMENT ON TABLE orders IS '订单表';
CREATE INDEX IF NOT EXISTS idx_orders_user_created ON orders (user_id, created_at DESC);
CREATE INDEX IF NOT EXISTS idx_orders_status_created ON orders (status, created_at DESC);

CREATE TABLE IF NOT EXISTS order_items (
    id BIGSERIAL PRIMARY KEY,
    order_id BIGINT NOT NULL,
    dish_id BIGINT NOT NULL,
    dish_name VARCHAR(100) NOT NULL,
    dish_image_url VARCHAR(255) NOT NULL DEFAULT '',
    unit_price NUMERIC(10,2) NOT NULL,
    quantity INT NOT NULL,
    subtotal_amount NUMERIC(10,2) NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT fk_order_items_order_id FOREIGN KEY (order_id) REFERENCES orders(id)
);
COMMENT ON TABLE order_items IS '订单明细表';
CREATE INDEX IF NOT EXISTS idx_order_items_order_id ON order_items (order_id);

CREATE TABLE IF NOT EXISTS system_configs (
    id BIGSERIAL PRIMARY KEY,
    config_key VARCHAR(100) NOT NULL,
    config_value TEXT NOT NULL,
    remark VARCHAR(255) NOT NULL DEFAULT '',
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT uk_system_configs_key UNIQUE (config_key)
);
COMMENT ON TABLE system_configs IS '系统配置表';
```

- [x] **Step 4: 增加状态常量**

创建 `server/internal/consts/status.go`：

```go
package consts

const (
	// StatusEnabled 启用状态。
	StatusEnabled = 1
	// StatusDisabled 禁用状态。
	StatusDisabled = 2
)

const (
	// DishStatusOn 菜品上架状态。
	DishStatusOn = 1
	// DishStatusOff 菜品下架状态。
	DishStatusOff = 2
)

const (
	// OrderStatusPending 订单待处理。
	OrderStatusPending = 1
	// OrderStatusCooking 订单制作中。
	OrderStatusCooking = 2
	// OrderStatusCompleted 订单已完成。
	OrderStatusCompleted = 3
	// OrderStatusCanceled 订单已取消。
	OrderStatusCanceled = 4
)
```

- [x] **Step 5: 增加配置示例**

创建 `server/manifest/config/config.example.yaml`：

```yaml
server:
  address: ":8000"
  openapiPath: "/api.json"
  swaggerPath: "/swagger"

database:
  default:
    link: "pgsql:host=127.0.0.1 port=5432 user=family_order password=family_order dbname=family_order sslmode=disable"

jwt:
  secret: "please-change-me"
  miniappExpireHours: 168
  adminExpireHours: 24

wechat:
  appId: "your-miniapp-app-id"
  appSecret: "your-miniapp-app-secret"
```

- [x] **Step 6: 运行测试**

Run:

```bash
cd server
gofmt -w internal test
go test ./...
```

Expected: 测试通过。

- [x] **Step 7: 提交数据库设计**

Run:

```bash
git add server/manifest server/internal/consts server/test/migration_test.go
git commit -m "feat: 新增点餐系统数据库迁移"
```

Expected: 提交成功。

---

## Task 4: 实现统一响应、错误码和请求日志

**Files:**
- Create: `server/internal/consts/error_code.go`
- Create: `server/internal/model/response.go`
- Create: `server/internal/middleware/response.go`
- Create: `server/internal/middleware/request_log.go`
- Modify: `server/internal/cmd/cmd.go`
- Create: `server/test/error_code_test.go`

- [x] **Step 1: 编写错误码测试**

创建 `server/test/error_code_test.go`：

```go
package test

import (
	"testing"

	"family-order/server/internal/consts"
)

// TestErrorCodeValues 验证核心错误码稳定。
func TestErrorCodeValues(t *testing.T) {
	cases := map[string]int{
		"成功":        consts.CodeSuccess,
		"参数错误":      consts.CodeInvalidParams,
		"未登录":       consts.CodeUnauthorized,
		"不在白名单":     consts.CodeForbidden,
		"资源不存在":     consts.CodeNotFound,
		"状态冲突":      consts.CodeConflict,
		"系统错误":      consts.CodeSystemError,
	}
	expected := map[string]int{
		"成功":        0,
		"参数错误":      40001,
		"未登录":       40101,
		"不在白名单":     40301,
		"资源不存在":     40401,
		"状态冲突":      40901,
		"系统错误":      50001,
	}
	for name, actual := range cases {
		if actual != expected[name] {
			t.Fatalf("%s 错误码期望 %d，实际 %d", name, expected[name], actual)
		}
	}
}
```

- [x] **Step 2: 运行测试确认失败**

Run:

```bash
cd server
go test ./test -run TestErrorCodeValues -v
```

Expected: 因常量不存在而失败。

- [x] **Step 3: 实现错误码和响应结构**

创建 `server/internal/consts/error_code.go`：

```go
package consts

const (
	// CodeSuccess 成功。
	CodeSuccess = 0
	// CodeInvalidParams 参数错误。
	CodeInvalidParams = 40001
	// CodeUnauthorized 未登录或 Token 失效。
	CodeUnauthorized = 40101
	// CodeForbidden 无访问权限。
	CodeForbidden = 40301
	// CodeNotFound 资源不存在。
	CodeNotFound = 40401
	// CodeConflict 状态冲突。
	CodeConflict = 40901
	// CodeSystemError 系统错误。
	CodeSystemError = 50001
)
```

创建 `server/internal/model/response.go`：

```go
package model

// Response 统一接口响应。
type Response struct {
	Code int         `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data"`
}

// NewSuccessResponse 创建成功响应。
func NewSuccessResponse(data interface{}) Response {
	return Response{Code: 0, Msg: "success", Data: data}
}

// NewErrorResponse 创建失败响应。
func NewErrorResponse(code int, msg string) Response {
	return Response{Code: code, Msg: msg, Data: nil}
}
```

- [x] **Step 4: 实现请求日志中间件**

创建 `server/internal/middleware/request_log.go`：

```go
package middleware

import (
	"time"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
)

// RequestLog 记录请求日志。
func RequestLog(r *ghttp.Request) {
	start := time.Now()
	r.Middleware.Next()
	g.Log().Infof(r.Context(), "请求完成 method=%s path=%s status=%d cost=%s",
		r.Method, r.URL.Path, r.Response.Status, time.Since(start))
}
```

- [x] **Step 5: 接入中间件**

修改 `server/internal/cmd/cmd.go`，在路由组内增加：

```go
group.Middleware(middleware.RequestLog)
```

并导入：

```go
"family-order/server/internal/middleware"
```

- [x] **Step 6: 运行测试**

Run:

```bash
cd server
gofmt -w internal test
go test ./...
```

Expected: 测试通过。

- [x] **Step 7: 提交统一响应基础**

Run:

```bash
git add server
git commit -m "feat: 新增统一响应和请求日志"
```

Expected: 提交成功。

---

## Task 5: 实现 JWT、管理员登录和管理端鉴权

**Files:**
- Create: `server/internal/logic/auth/jwt.go`
- Create: `server/internal/logic/auth/admin.go`
- Create: `server/internal/middleware/auth.go`
- Create: `server/api/admin/auth.go`
- Create: `server/internal/controller/admin/auth.go`
- Create: `server/test/auth_test.go`

- [ ] **Step 1: 安装依赖**

Run:

```bash
cd server
go get github.com/golang-jwt/jwt/v5@latest
go get golang.org/x/crypto/bcrypt@latest
```

Expected: JWT 和 bcrypt 依赖写入 `go.mod`。

- [ ] **Step 2: 编写 JWT 测试**

创建 `server/test/auth_test.go`：

```go
package test

import (
	"testing"
	"time"

	authlogic "family-order/server/internal/logic/auth"
)

// TestJWTGenerateAndParse 验证 JWT 可以生成和解析。
func TestJWTGenerateAndParse(t *testing.T) {
	token, err := authlogic.GenerateToken(authlogic.TokenClaimsInput{
		SubjectID: 1,
		Role:      "admin",
		Secret:    "test-secret",
		ExpiresIn: time.Hour,
	})
	if err != nil {
		t.Fatalf("生成 Token 失败: %v", err)
	}
	claims, err := authlogic.ParseToken(token, "test-secret")
	if err != nil {
		t.Fatalf("解析 Token 失败: %v", err)
	}
	if claims.SubjectID != 1 || claims.Role != "admin" {
		t.Fatalf("Token 内容不符合预期: %+v", claims)
	}
}
```

- [ ] **Step 3: 运行测试确认失败**

Run:

```bash
cd server
go test ./test -run TestJWTGenerateAndParse -v
```

Expected: 因 JWT 逻辑不存在而失败。

- [ ] **Step 4: 实现 JWT 逻辑**

创建 `server/internal/logic/auth/jwt.go`：

```go
package auth

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// TokenClaimsInput 生成 Token 的输入参数。
type TokenClaimsInput struct {
	SubjectID int64
	Role      string
	Secret    string
	ExpiresIn time.Duration
}

// TokenClaims Token 载荷。
type TokenClaims struct {
	SubjectID int64  `json:"subject_id"`
	Role      string `json:"role"`
	jwt.RegisteredClaims
}

// GenerateToken 生成 JWT。
func GenerateToken(in TokenClaimsInput) (string, error) {
	now := time.Now()
	claims := TokenClaims{
		SubjectID: in.SubjectID,
		Role:      in.Role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(now.Add(in.ExpiresIn)),
			IssuedAt:  jwt.NewNumericDate(now),
		},
	}
	return jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString([]byte(in.Secret))
}

// ParseToken 解析 JWT。
func ParseToken(tokenString string, secret string) (*TokenClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &TokenClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(secret), nil
	})
	if err != nil {
		return nil, fmt.Errorf("解析Token失败: %w", err)
	}
	claims, ok := token.Claims.(*TokenClaims)
	if !ok || !token.Valid {
		return nil, fmt.Errorf("Token无效")
	}
	return claims, nil
}
```

- [ ] **Step 5: 实现管理员登录接口**

实现以下文件：

`server/api/admin/auth.go` 定义登录请求、响应、个人资料响应。

`server/internal/controller/admin/auth.go` 负责参数接收和调用登录逻辑。

`server/internal/logic/auth/admin.go` 负责根据用户名查询管理员、bcrypt 校验密码、生成管理员 JWT。

管理员登录错误必须返回明确中文错误，例如：

```go
return nil, fmt.Errorf("管理员账号或密码错误")
```

- [ ] **Step 6: 运行测试**

Run:

```bash
cd server
gofmt -w api internal test
go test ./...
```

Expected: 测试通过。

- [ ] **Step 7: 提交管理员认证**

Run:

```bash
git add server
git commit -m "feat: 新增管理员登录鉴权"
```

Expected: 提交成功。

---

## Task 6: 实现微信登录和白名单校验

**Files:**
- Create: `server/internal/logic/auth/wechat.go`
- Create: `server/internal/logic/user/user.go`
- Create: `server/api/miniapp/auth.go`
- Create: `server/internal/controller/miniapp/auth.go`
- Modify: `server/internal/middleware/auth.go`
- Create: `server/test/wechat_auth_test.go`

- [ ] **Step 1: 编写微信客户端接口测试**

创建 `server/test/wechat_auth_test.go`：

```go
package test

import (
	"context"
	"testing"

	authlogic "family-order/server/internal/logic/auth"
)

type fakeWechatClient struct{}

// Code2Session 模拟微信 code 换 OpenID。
func (fakeWechatClient) Code2Session(ctx context.Context, code string) (*authlogic.WechatSession, error) {
	return &authlogic.WechatSession{OpenID: "openid-family-001", SessionKey: "session-key"}, nil
}

// TestWechatClientInterface 验证微信登录依赖可以被替换为 mock。
func TestWechatClientInterface(t *testing.T) {
	client := fakeWechatClient{}
	session, err := client.Code2Session(context.Background(), "test-code")
	if err != nil {
		t.Fatalf("模拟微信登录失败: %v", err)
	}
	if session.OpenID != "openid-family-001" {
		t.Fatalf("OpenID 不符合预期: %s", session.OpenID)
	}
}
```

- [ ] **Step 2: 运行测试确认失败**

Run:

```bash
cd server
go test ./test -run TestWechatClientInterface -v
```

Expected: 因 `WechatSession` 未定义而失败。

- [ ] **Step 3: 实现微信客户端接口**

创建 `server/internal/logic/auth/wechat.go`：

```go
package auth

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
)

// WechatClient 微信登录客户端接口。
type WechatClient interface {
	Code2Session(ctx context.Context, code string) (*WechatSession, error)
}

// WechatSession 微信登录会话。
type WechatSession struct {
	OpenID     string `json:"openid"`
	SessionKey string `json:"session_key"`
}

// HTTPWechatClient 基于 HTTP 的微信客户端。
type HTTPWechatClient struct {
	AppID     string
	AppSecret string
}

// Code2Session 使用微信 code 换取 OpenID。
func (c *HTTPWechatClient) Code2Session(ctx context.Context, code string) (*WechatSession, error) {
	values := url.Values{}
	values.Set("appid", c.AppID)
	values.Set("secret", c.AppSecret)
	values.Set("js_code", code)
	values.Set("grant_type", "authorization_code")
	apiURL := "https://api.weixin.qq.com/sns/jscode2session?" + values.Encode()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, apiURL, nil)
	if err != nil {
		return nil, fmt.Errorf("创建微信登录请求失败: %w", err)
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("请求微信登录接口失败: %w", err)
	}
	defer resp.Body.Close()

	var session WechatSession
	if err := json.NewDecoder(resp.Body).Decode(&session); err != nil {
		return nil, fmt.Errorf("解析微信登录响应失败: %w", err)
	}
	if session.OpenID == "" {
		return nil, fmt.Errorf("微信登录未返回OpenID")
	}
	return &session, nil
}
```

- [ ] **Step 4: 实现用户白名单逻辑**

在 `server/internal/logic/user/user.go` 实现：

- 根据 OpenID 查询用户。
- 用户不存在时创建用户。
- 更新昵称、头像和最近登录时间。
- 判断 `is_whitelist` 和 `status`。
- 非白名单用户允许登录，但返回 `is_whitelist=false`。

- [ ] **Step 5: 实现小程序登录接口**

实现：

```http
POST /api/miniapp/auth/login
```

响应结构：

```json
{
  "token": "JWT",
  "is_whitelist": true
}
```

白名单用户返回可访问点餐接口的 JWT；非白名单用户返回空 Token 或受限 Token，并在访问菜单和下单接口时由中间件拦截。

- [ ] **Step 6: 运行测试**

Run:

```bash
cd server
gofmt -w api internal test
go test ./...
```

Expected: 测试通过。

- [ ] **Step 7: 提交微信登录**

Run:

```bash
git add server
git commit -m "feat: 新增小程序微信登录"
```

Expected: 提交成功。

---

## Task 7: 实现分类和菜品接口

**Files:**
- Create: `server/api/admin/category.go`
- Create: `server/api/admin/dish.go`
- Create: `server/api/miniapp/menu.go`
- Create: `server/internal/logic/menu/category.go`
- Create: `server/internal/logic/menu/dish.go`
- Create: `server/internal/controller/admin/category.go`
- Create: `server/internal/controller/admin/dish.go`
- Create: `server/internal/controller/miniapp/menu.go`
- Create: `server/test/menu_test.go`

- [ ] **Step 1: 编写菜单排序测试**

创建 `server/test/menu_test.go`：

```go
package test

import (
	"testing"

	menulogic "family-order/server/internal/logic/menu"
)

// TestValidateDishQuantity 验证菜品数量校验。
func TestValidateDishQuantity(t *testing.T) {
	if err := menulogic.ValidateDishQuantity(0); err == nil {
		t.Fatalf("数量为0时应该返回错误")
	}
	if err := menulogic.ValidateDishQuantity(1); err != nil {
		t.Fatalf("数量为1时不应该返回错误: %v", err)
	}
}
```

- [ ] **Step 2: 运行测试确认失败**

Run:

```bash
cd server
go test ./test -run TestValidateDishQuantity -v
```

Expected: 因校验函数不存在而失败。

- [ ] **Step 3: 实现菜品基础校验**

创建 `server/internal/logic/menu/dish.go`：

```go
package menu

import "fmt"

// ValidateDishQuantity 校验点餐数量。
func ValidateDishQuantity(quantity int) error {
	if quantity <= 0 {
		return fmt.Errorf("菜品数量必须大于0")
	}
	if quantity > 99 {
		return fmt.Errorf("单个菜品数量不能超过99")
	}
	return nil
}
```

- [ ] **Step 4: 实现管理端分类接口**

实现以下接口：

```http
GET /api/admin/categories
POST /api/admin/categories
PUT /api/admin/categories/{id}
DELETE /api/admin/categories/{id}
PUT /api/admin/categories/{id}/status
```

规则：

- 分类名称不能为空。
- 排序值默认为 `100`。
- 分类下存在菜品时不允许删除。
- 删除失败返回 `40901`。

- [ ] **Step 5: 实现管理端菜品接口**

实现以下接口：

```http
GET /api/admin/dishes
POST /api/admin/dishes
GET /api/admin/dishes/{id}
PUT /api/admin/dishes/{id}
DELETE /api/admin/dishes/{id}
PUT /api/admin/dishes/{id}/status
```

规则：

- 菜品名称不能为空。
- 价格必须大于等于 `0`。
- 分类必须存在且启用。
- 删除优先做下架处理，不物理删除历史数据。

- [ ] **Step 6: 实现小程序菜单接口**

实现以下接口：

```http
GET /api/miniapp/menu/categories
GET /api/miniapp/menu/dishes?category_id=1
GET /api/miniapp/menu/dishes/{id}
```

规则：

- 必须通过小程序 JWT 和白名单校验。
- 只返回启用分类和上架菜品。
- 按 `sort ASC, id ASC` 排序。

- [ ] **Step 7: 运行测试**

Run:

```bash
cd server
gofmt -w api internal test
go test ./...
```

Expected: 测试通过。

- [ ] **Step 8: 提交菜单接口**

Run:

```bash
git add server
git commit -m "feat: 新增分类和菜品接口"
```

Expected: 提交成功。

---

## Task 8: 实现订单创建、查询和状态流转

**Files:**
- Create: `server/api/admin/order.go`
- Create: `server/api/miniapp/order.go`
- Create: `server/internal/logic/order/order.go`
- Create: `server/internal/controller/admin/order.go`
- Create: `server/internal/controller/miniapp/order.go`
- Create: `server/test/order_test.go`

- [ ] **Step 1: 编写订单状态流转测试**

创建 `server/test/order_test.go`：

```go
package test

import (
	"testing"

	"family-order/server/internal/consts"
	orderlogic "family-order/server/internal/logic/order"
)

// TestCanChangeStatus 验证订单状态流转规则。
func TestCanChangeStatus(t *testing.T) {
	if !orderlogic.CanChangeStatus(consts.OrderStatusPending, consts.OrderStatusCooking) {
		t.Fatalf("待处理应该允许变更为制作中")
	}
	if !orderlogic.CanChangeStatus(consts.OrderStatusCooking, consts.OrderStatusCompleted) {
		t.Fatalf("制作中应该允许变更为已完成")
	}
	if orderlogic.CanChangeStatus(consts.OrderStatusCompleted, consts.OrderStatusCanceled) {
		t.Fatalf("已完成不允许变更为已取消")
	}
}
```

- [ ] **Step 2: 运行测试确认失败**

Run:

```bash
cd server
go test ./test -run TestCanChangeStatus -v
```

Expected: 因订单逻辑不存在而失败。

- [ ] **Step 3: 实现状态流转函数**

创建 `server/internal/logic/order/order.go`：

```go
package order

import "family-order/server/internal/consts"

// CanChangeStatus 判断订单状态是否允许流转。
func CanChangeStatus(from int, to int) bool {
	switch from {
	case consts.OrderStatusPending:
		return to == consts.OrderStatusCooking || to == consts.OrderStatusCanceled
	case consts.OrderStatusCooking:
		return to == consts.OrderStatusCompleted || to == consts.OrderStatusCanceled
	default:
		return false
	}
}
```

- [ ] **Step 4: 实现下单接口**

实现：

```http
POST /api/miniapp/orders
```

规则：

- 必须校验白名单。
- `items` 不能为空。
- 每个 `quantity` 必须在 `1` 到 `99`。
- 查询数据库中的菜品价格，不信任前端价格。
- 菜品必须存在且上架。
- 使用事务创建 `orders` 和 `order_items`。
- `order_items` 保存菜品名称、图片、单价快照。
- 订单号格式使用 `yyyyMMddHHmmss` 加随机后缀，保证唯一。

- [ ] **Step 5: 实现我的订单接口**

实现：

```http
GET /api/miniapp/orders
GET /api/miniapp/orders/{id}
POST /api/miniapp/orders/{id}/cancel
```

规则：

- 只能查询自己的订单。
- 订单详情包含订单主表和明细。
- 只能取消自己的待处理订单。

- [ ] **Step 6: 实现管理端订单接口**

实现：

```http
GET /api/admin/orders
GET /api/admin/orders/{id}
PUT /api/admin/orders/{id}/status
```

规则：

- 支持按状态筛选。
- 支持分页。
- 修改状态前调用 `CanChangeStatus`。
- 状态冲突返回 `40901`。

- [ ] **Step 7: 运行测试**

Run:

```bash
cd server
gofmt -w api internal test
go test ./...
```

Expected: 测试通过。

- [ ] **Step 8: 提交订单模块**

Run:

```bash
git add server
git commit -m "feat: 新增订单下单和状态流转"
```

Expected: 提交成功。

---

## Task 9: 实现 Docker Compose 本地部署

**Files:**
- Create: `server/manifest/docker/Dockerfile`
- Create: `deploy/compose/docker-compose.yaml`
- Create: `deploy/nginx/family-order.conf`
- Create: `.env.example`

- [ ] **Step 1: 编写后端 Dockerfile**

创建 `server/manifest/docker/Dockerfile`：

```dockerfile
FROM golang:1.23-alpine AS builder

WORKDIR /app
COPY server/go.mod server/go.sum ./
RUN go mod download
COPY server/ ./
RUN go build -o family-order-server main.go

FROM alpine:3.20

WORKDIR /app
COPY --from=builder /app/family-order-server /app/family-order-server
COPY server/manifest/config/config.example.yaml /app/manifest/config/config.yaml
EXPOSE 8000
CMD ["/app/family-order-server"]
```

- [ ] **Step 2: 编写 Compose 配置**

创建 `deploy/compose/docker-compose.yaml`：

```yaml
services:
  postgres:
    image: postgres:16-alpine
    container_name: family-order-postgres
    environment:
      POSTGRES_DB: family_order
      POSTGRES_USER: family_order
      POSTGRES_PASSWORD: family_order
    volumes:
      - postgres_data:/var/lib/postgresql/data
      - ../../server/manifest/deploy/migrations:/docker-entrypoint-initdb.d:ro
    ports:
      - "5432:5432"

  backend:
    build:
      context: ../..
      dockerfile: server/manifest/docker/Dockerfile
    container_name: family-order-backend
    depends_on:
      - postgres
    ports:
      - "8000:8000"

  nginx:
    image: nginx:1.27-alpine
    container_name: family-order-nginx
    depends_on:
      - backend
    volumes:
      - ../nginx/family-order.conf:/etc/nginx/conf.d/default.conf:ro
    ports:
      - "80:80"

volumes:
  postgres_data:
```

- [ ] **Step 3: 编写 Nginx 配置**

创建 `deploy/nginx/family-order.conf`：

```nginx
server {
    listen 80;
    server_name _;

    client_max_body_size 10m;

    location /api/ {
        proxy_pass http://backend:8000/api/;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
    }

    location /health {
        proxy_pass http://backend:8000/health;
    }
}
```

- [ ] **Step 4: 验证 Compose 配置**

Run:

```bash
docker compose -f deploy/compose/docker-compose.yaml config
```

Expected: 输出解析后的 Compose 配置，无错误。

- [ ] **Step 5: 构建后端镜像**

Run:

```bash
docker compose -f deploy/compose/docker-compose.yaml build backend
```

Expected: 后端镜像构建成功。

- [ ] **Step 6: 提交部署配置**

Run:

```bash
git add server/manifest/docker deploy .env.example
git commit -m "feat: 新增Docker Compose部署配置"
```

Expected: 提交成功。

---

## Task 10: 初始化 Vue 管理端

**Files:**
- Create: `admin-web/package.json`
- Create: `admin-web/index.html`
- Create: `admin-web/src/main.ts`
- Create: `admin-web/src/App.vue`
- Create: `admin-web/src/router/index.ts`
- Create: `admin-web/src/stores/auth.ts`
- Create: `admin-web/src/api/request.ts`

- [ ] **Step 1: 创建 Vue 项目依赖**

在 `admin-web/package.json` 写入：

```json
{
  "scripts": {
    "dev": "vite --host 0.0.0.0",
    "build": "vue-tsc && vite build",
    "preview": "vite preview --host 0.0.0.0"
  },
  "dependencies": {
    "@vitejs/plugin-vue": "latest",
    "axios": "latest",
    "element-plus": "latest",
    "pinia": "latest",
    "vue": "latest",
    "vue-router": "latest"
  },
  "devDependencies": {
    "typescript": "latest",
    "vite": "latest",
    "vue-tsc": "latest"
  }
}
```

- [ ] **Step 2: 安装依赖**

Run:

```bash
cd admin-web
npm install
```

Expected: `node_modules` 和 `package-lock.json` 生成成功。

- [ ] **Step 3: 实现请求封装**

创建 `admin-web/src/api/request.ts`：

```ts
import axios from 'axios'

export const request = axios.create({
  baseURL: import.meta.env.VITE_API_BASE_URL || '/api/admin',
  timeout: 10000,
})

request.interceptors.request.use((config) => {
  const token = localStorage.getItem('admin_token')
  if (token) {
    config.headers.Authorization = `Bearer ${token}`
  }
  return config
})

request.interceptors.response.use((response) => {
  const body = response.data
  if (body.code !== 0) {
    return Promise.reject(new Error(body.msg || '请求失败'))
  }
  return body.data
})
```

- [ ] **Step 4: 实现路由和登录态**

创建 `admin-web/src/router/index.ts` 和 `admin-web/src/stores/auth.ts`，实现：

- `/login` 登录页。
- `/` 首页。
- 未登录访问后台页面时跳转 `/login`。
- 登录成功保存 `admin_token`。

- [ ] **Step 5: 构建验证**

Run:

```bash
cd admin-web
npm run build
```

Expected: 构建通过，生成 `dist/`。

- [ ] **Step 6: 提交管理端骨架**

Run:

```bash
git add admin-web
git commit -m "feat: 初始化Vue管理后台"
```

Expected: 提交成功。

---

## Task 11: 实现管理端页面

**Files:**
- Create: `admin-web/src/api/auth.ts`
- Create: `admin-web/src/api/category.ts`
- Create: `admin-web/src/api/dish.ts`
- Create: `admin-web/src/api/order.ts`
- Create: `admin-web/src/api/user.ts`
- Create: `admin-web/src/components/PageHeader.vue`
- Create: `admin-web/src/views/Login.vue`
- Create: `admin-web/src/views/Dashboard.vue`
- Create: `admin-web/src/views/CategoryList.vue`
- Create: `admin-web/src/views/DishList.vue`
- Create: `admin-web/src/views/DishForm.vue`
- Create: `admin-web/src/views/OrderList.vue`
- Create: `admin-web/src/views/OrderDetail.vue`
- Create: `admin-web/src/views/UserWhitelist.vue`

- [ ] **Step 1: 实现 API 文件**

每个 API 文件只封装对应资源：

```ts
import { request } from './request'

export interface Category {
  id: number
  name: string
  sort: number
  status: number
}

export function listCategories() {
  return request.get<Category[]>('/categories')
}
```

`dish.ts`、`order.ts`、`user.ts` 按同样方式定义类型和函数。

- [ ] **Step 2: 实现登录页**

`Login.vue` 包含：

- 用户名输入框。
- 密码输入框。
- 登录按钮。
- 登录失败中文提示。

登录成功后跳转 `/`。

- [ ] **Step 3: 实现分类管理**

`CategoryList.vue` 包含：

- 表格展示分类名称、排序、状态。
- 新增和编辑弹窗。
- 启用、禁用按钮。
- 删除前二次确认。

- [ ] **Step 4: 实现菜品管理**

`DishList.vue` 和 `DishForm.vue` 包含：

- 按分类筛选。
- 菜品名称、图片、价格、单位、排序、状态。
- 上架和下架操作。
- 表单校验：名称必填、价格不小于 0。

- [ ] **Step 5: 实现订单管理**

`OrderList.vue` 和 `OrderDetail.vue` 包含：

- 按状态筛选。
- 展示订单号、用户、金额、状态、创建时间。
- 详情展示订单明细。
- 支持待处理改制作中、待处理改已取消、制作中改已完成、制作中改已取消。

- [ ] **Step 6: 实现白名单管理**

`UserWhitelist.vue` 包含：

- 用户列表。
- 白名单开关。
- 启用和禁用操作。
- OpenID 脱敏展示。

- [ ] **Step 7: 构建验证**

Run:

```bash
cd admin-web
npm run build
```

Expected: 构建通过。

- [ ] **Step 8: 提交管理端页面**

Run:

```bash
git add admin-web
git commit -m "feat: 新增管理后台核心页面"
```

Expected: 提交成功。

---

## Task 12: 初始化微信小程序端

**Files:**
- Create: `miniapp/app.js`
- Create: `miniapp/app.json`
- Create: `miniapp/app.wxss`
- Create: `miniapp/project.config.json`
- Create: `miniapp/utils/request.js`
- Create: `miniapp/utils/cart.js`
- Create: `miniapp/pages/login/login.js`
- Create: `miniapp/pages/login/login.json`
- Create: `miniapp/pages/login/login.wxml`
- Create: `miniapp/pages/login/login.wxss`

- [ ] **Step 1: 创建小程序配置**

`miniapp/app.json`：

```json
{
  "pages": [
    "pages/login/login",
    "pages/menu/menu",
    "pages/cart/cart",
    "pages/orders/orders",
    "pages/order-detail/order-detail"
  ],
  "window": {
    "navigationBarTitleText": "家庭点餐",
    "navigationBarBackgroundColor": "#ffffff",
    "navigationBarTextStyle": "black"
  }
}
```

- [ ] **Step 2: 实现请求封装**

`miniapp/utils/request.js`：

```js
const BASE_URL = 'https://api.yourdomain.com/api/miniapp'

function request(options) {
  const token = wx.getStorageSync('token')
  return new Promise((resolve, reject) => {
    wx.request({
      url: `${BASE_URL}${options.url}`,
      method: options.method || 'GET',
      data: options.data || {},
      header: {
        Authorization: token ? `Bearer ${token}` : '',
      },
      success(res) {
        const body = res.data
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
```

- [ ] **Step 3: 实现登录页**

登录流程：

- 调用 `wx.login()` 获取 code。
- 请求 `/auth/login`。
- 保存 `token`。
- `is_whitelist=false` 时展示“当前账号暂未开通点餐权限，请联系管理员”。
- `is_whitelist=true` 时跳转菜单页。

- [ ] **Step 4: 使用微信开发者工具打开项目**

Run:

```bash
open -a "微信开发者工具" miniapp
```

Expected: 小程序项目能被微信开发者工具打开。

- [ ] **Step 5: 提交小程序骨架**

Run:

```bash
git add miniapp
git commit -m "feat: 初始化微信小程序端"
```

Expected: 提交成功。

---

## Task 13: 实现小程序点餐页面

**Files:**
- Create: `miniapp/pages/menu/menu.js`
- Create: `miniapp/pages/menu/menu.json`
- Create: `miniapp/pages/menu/menu.wxml`
- Create: `miniapp/pages/menu/menu.wxss`
- Create: `miniapp/pages/cart/cart.js`
- Create: `miniapp/pages/cart/cart.json`
- Create: `miniapp/pages/cart/cart.wxml`
- Create: `miniapp/pages/cart/cart.wxss`
- Create: `miniapp/pages/orders/orders.js`
- Create: `miniapp/pages/orders/orders.json`
- Create: `miniapp/pages/orders/orders.wxml`
- Create: `miniapp/pages/orders/orders.wxss`
- Create: `miniapp/pages/order-detail/order-detail.js`
- Create: `miniapp/pages/order-detail/order-detail.json`
- Create: `miniapp/pages/order-detail/order-detail.wxml`
- Create: `miniapp/pages/order-detail/order-detail.wxss`

- [ ] **Step 1: 实现购物车工具**

`miniapp/utils/cart.js`：

```js
const CART_KEY = 'cart_items'

function getCart() {
  return wx.getStorageSync(CART_KEY) || []
}

function saveCart(items) {
  wx.setStorageSync(CART_KEY, items)
}

function clearCart() {
  wx.removeStorageSync(CART_KEY)
}

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
      price: dish.price,
      quantity: 1,
    })
  }
  saveCart(items)
}

module.exports = { getCart, saveCart, clearCart, addDish }
```

- [ ] **Step 2: 实现菜单页**

菜单页要求：

- 左侧展示分类。
- 右侧展示菜品。
- 菜品支持加入购物车。
- 底部展示购物车数量、总价和“去下单”按钮。

- [ ] **Step 3: 实现购物车页**

购物车页要求：

- 展示已选菜品。
- 支持增加、减少数量。
- 数量减少到 0 时移除菜品。
- 支持填写备注。
- 提交订单调用 `POST /orders`。
- 提交失败保留购物车。
- 提交成功清空购物车并跳转订单详情。

- [ ] **Step 4: 实现我的订单页**

订单页要求：

- 调用 `GET /orders`。
- 展示订单号、金额、状态和时间。
- 点击进入订单详情。

- [ ] **Step 5: 实现订单详情页**

详情页要求：

- 调用 `GET /orders/{id}`。
- 展示订单明细。
- 待处理订单展示取消按钮。
- 取消订单调用 `POST /orders/{id}/cancel`。

- [ ] **Step 6: 手工验收**

使用微信开发者工具验证：

- 未登录时进入登录页。
- 非白名单用户不能进入菜单。
- 白名单用户能浏览菜单。
- 下单成功后购物车清空。
- 下单失败时购物车保留。

- [ ] **Step 7: 提交小程序点餐功能**

Run:

```bash
git add miniapp
git commit -m "feat: 新增小程序点餐流程"
```

Expected: 提交成功。

---

## Task 14: 增加 K8s 预留配置和最终验收

**Files:**
- Create: `deploy/k8s/namespace.yaml`
- Create: `deploy/k8s/configmap.yaml`
- Create: `deploy/k8s/secret.example.yaml`
- Create: `deploy/k8s/backend-deployment.yaml`
- Create: `deploy/k8s/admin-web-deployment.yaml`
- Create: `deploy/k8s/ingress.yaml`
- Modify: `docs/family-order/README.md`

- [ ] **Step 1: 创建命名空间**

`deploy/k8s/namespace.yaml`：

```yaml
apiVersion: v1
kind: Namespace
metadata:
  name: family-order
```

- [ ] **Step 2: 创建配置示例**

`deploy/k8s/secret.example.yaml`：

```yaml
apiVersion: v1
kind: Secret
metadata:
  name: family-order-secret
  namespace: family-order
type: Opaque
stringData:
  DATABASE_DSN: "pgsql:host=postgres port=5432 user=family_order password=family_order dbname=family_order sslmode=disable"
  JWT_SECRET: "please-change-me"
  WECHAT_APP_SECRET: "please-change-me"
```

- [ ] **Step 3: 编写后端 Deployment**

`deploy/k8s/backend-deployment.yaml` 包含：

- `Deployment/family-order-backend`。
- `Service/family-order-backend`。
- 容器端口 `8000`。
- 健康检查路径 `/health`。
- 从 Secret 读取数据库、JWT、微信密钥。

- [ ] **Step 4: 编写管理端 Deployment**

`deploy/k8s/admin-web-deployment.yaml` 包含：

- `Deployment/family-order-admin-web`。
- `Service/family-order-admin-web`。
- 容器端口 `80`。

- [ ] **Step 5: 编写 Ingress**

`deploy/k8s/ingress.yaml` 支持：

- `api.yourdomain.com` 转发到后端。
- `admin.yourdomain.com` 转发到管理端。
- TLS 证书由集群环境提供。

- [ ] **Step 6: 更新 README**

在 `docs/family-order/README.md` 增加：

```markdown
## 本地部署

```bash
docker compose -f deploy/compose/docker-compose.yaml up -d
```

## 后端验证

```bash
curl http://127.0.0.1/health
```

## K8s 说明

`deploy/k8s/` 仅作为第二阶段部署模板。生产使用前需要替换域名、镜像地址和 Secret。
```

- [ ] **Step 7: 执行最终验证**

Run:

```bash
cd server
go test ./...
```

Expected: 后端测试通过。

Run:

```bash
cd admin-web
npm run build
```

Expected: 管理端构建通过。

Run:

```bash
docker compose -f deploy/compose/docker-compose.yaml config
```

Expected: Compose 配置通过。

- [ ] **Step 8: 提交部署模板**

Run:

```bash
git add deploy docs/family-order/README.md
git commit -m "feat: 新增K8s部署模板"
```

Expected: 提交成功。

---

## 总体验收清单

- [ ] 后端 `go test ./...` 通过。
- [ ] 后端 `/health` 可访问。
- [ ] PostgreSQL 初始化表结构成功。
- [ ] 管理员可以登录管理后台。
- [ ] 管理员可以新增分类。
- [ ] 管理员可以新增菜品。
- [ ] 管理员可以把小程序用户加入白名单。
- [ ] 白名单用户可以在小程序端看到菜单。
- [ ] 白名单用户可以提交订单。
- [ ] 后台可以查看订单详情。
- [ ] 后台可以修改订单状态。
- [ ] 非白名单用户无法访问菜单和下单。
- [ ] Docker Compose 配置可解析。
- [ ] Nginx 可以代理 `/health` 和 `/api/`。

## 实施顺序建议

推荐先执行 Task 1 到 Task 9，形成后端和 Docker 的可运行闭环；再执行 Task 10 到 Task 11 完成管理后台；最后执行 Task 12 到 Task 13 完成小程序端。Task 14 作为部署模板和最终验收任务，在核心业务跑通后执行。
