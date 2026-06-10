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
