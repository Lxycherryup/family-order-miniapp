package test

import (
	"os"
	"strings"
	"testing"
)

// TestDatabaseLinkUsesGoFrameFormat 验证数据库连接串使用 GoFrame 支持的格式。
func TestDatabaseLinkUsesGoFrameFormat(t *testing.T) {
	content, err := os.ReadFile("../manifest/config/config.example.yaml")
	if err != nil {
		t.Fatalf("读取配置文件失败: %v", err)
	}

	config := string(content)
	want := "pgsql:family_order:family_order@tcp(127.0.0.1:5432)/family_order?sslmode=disable"
	if !strings.Contains(config, want) {
		t.Fatalf("数据库连接串格式错误，期望包含: %s", want)
	}
	if strings.Contains(config, "pgsql:host=") {
		t.Fatal("数据库连接串不能使用 PostgreSQL 原生 DSN 风格")
	}
}
