package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestLoadConfigAppliesEnvironmentOverrides(t *testing.T) {
	t.Setenv("SERVER_ADDRESS", "0.0.0.0:9000")
	t.Setenv("SERVER_ENABLE_SWAGGER", "false")
	t.Setenv("SERVER_PUBLIC_URL", "http://docker.local:8080/")
	t.Setenv("DB_DIALECTS", "postgres")
	t.Setenv("DB_HOST", "postgres")
	t.Setenv("DB_PORT", "5433")
	t.Setenv("DB_NAME", "phase1")
	t.Setenv("DB_USER", "phase1_user")
	t.Setenv("DB_PASSWORD", "phase1_password")
	t.Setenv("DB_SSLMODE", "disable")
	t.Setenv("DB_MIGRATION_PATH", "/app/migrations")
	t.Setenv("DB_AUTO_MIGRATE", "false")
	t.Setenv("REDIS_ADDR", "redis:6379")
	t.Setenv("REDIS_PASSWORD", "redis-pass")
	t.Setenv("IMAGE_HOST", "http://docker.local:8080")
	t.Setenv("PUSHGATEWAY_URL", "http://pushgateway:9091")

	Config = nil
	configPath := writeTempConfig(t)

	if err := LoadConfig(configPath); err != nil {
		t.Fatalf("LoadConfig() error = %v", err)
	}

	if Config.Server.Address != "0.0.0.0:9000" {
		t.Fatalf("Server.Address = %q, want %q", Config.Server.Address, "0.0.0.0:9000")
	}
	if Config.Server.EnableSwagger {
		t.Fatalf("Server.EnableSwagger = true, want false")
	}
	if Config.Server.PublicUrl != "http://docker.local:8080/" {
		t.Fatalf("Server.PublicUrl = %q, want %q", Config.Server.PublicUrl, "http://docker.local:8080/")
	}
	if Config.Db.Host != "postgres" || Config.Db.Port != 5433 || Config.Db.Db != "phase1" {
		t.Fatalf("DB override not applied: %+v", Config.Db)
	}
	if Config.Db.Username != "phase1_user" || Config.Db.Password != "phase1_password" {
		t.Fatalf("DB credentials override not applied: %+v", Config.Db)
	}
	if Config.Db.MigrationPath != "/app/migrations" || Config.Db.AutoMigrate {
		t.Fatalf("DB migration settings override not applied: %+v", Config.Db)
	}
	if Config.Redis.Address != "redis:6379" || Config.Redis.Password != "redis-pass" {
		t.Fatalf("Redis override not applied: %+v", Config.Redis)
	}
	if Config.ImageSettings.ImageHost != "http://docker.local:8080" {
		t.Fatalf("Image host override not applied: %+v", Config.ImageSettings)
	}
	if Config.Monitor.Pushgateway.URL != "http://pushgateway:9091" {
		t.Fatalf("Pushgateway override not applied: %+v", Config.Monitor.Pushgateway)
	}
}

func TestLoadConfigKeepsYamlValuesWhenEnvInvalid(t *testing.T) {
	t.Setenv("DB_PORT", "not-a-number")
	t.Setenv("DB_AUTO_MIGRATE", "not-a-bool")

	Config = nil
	configPath := writeTempConfig(t)

	if err := LoadConfig(configPath); err != nil {
		t.Fatalf("LoadConfig() error = %v", err)
	}

	if Config.Db.Port != 5432 {
		t.Fatalf("DB.Port = %d, want %d", Config.Db.Port, 5432)
	}
	if !Config.Db.AutoMigrate {
		t.Fatalf("DB.AutoMigrate = false, want true")
	}
}

func writeTempConfig(t *testing.T) string {
	t.Helper()

	dir := t.TempDir()
	path := filepath.Join(dir, "config.yaml")
	content := []byte(`server:
  address: 127.0.0.1:8000
  model: debug
  enableSwagger: true
  publicUrl: "http://127.0.0.1:8080/"
db:
  dialects: mysql
  host: 127.0.0.1
  port: 5432
  db: autoops
  username: postgres
  password: postgres123456
  charset: utf8
  sslMode: disable
  migrationPath: ./migrations
  autoMigrate: true
  maxIdle: 50
  maxOpen: 150
redis:
  address: 127.0.0.1:6379
  password: "123456"
imageSettings:
  uploadDir: ./upload/
  imageHost: http://127.0.0.1:8080
log:
  path: ./log
  name: sys
  model: console
monitor:
  prometheus:
    url: http://127.0.0.1:9090
  pushgateway:
    url: http://127.0.0.1:9091
  agent:
    heartbeat_server_url: http://127.0.0.1:8000/api/v1/monitor/agent/heartbeat
    heartbeat_token: agent-heartbeat-token-2024
  webhook:
    token: webhook-token
`)
	if err := os.WriteFile(path, content, 0o600); err != nil {
		t.Fatalf("WriteFile() error = %v", err)
	}

	return path
}
