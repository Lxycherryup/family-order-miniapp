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
