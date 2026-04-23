package config

import (
	"os"
	"strconv"

	"gopkg.in/yaml.v2"
)

type config struct {
	Server        server        `yaml:"server"`
	Db            db            `yaml:"db"`
	Redis         redis         `yaml:"redis"`
	Integrations  Integrations  `yaml:"integrations"`
	ImageSettings imageSettings `yaml:"imageSettings"`
	Log           log           `yaml:"log"`
	Monitor       monitor       `yaml:"monitor"`
}

type Integrations struct {
	Systems []ExternalSystem `yaml:"systems"`
}

type ExternalSystem struct {
	Key          string            `yaml:"key"`
	DisplayName  string            `yaml:"displayName"`
	Category     string            `yaml:"category"`
	Provider     string            `yaml:"provider"`
	Mode         string            `yaml:"mode"`
	BaseURL      string            `yaml:"baseUrl"`
	Enabled      bool              `yaml:"enabled"`
	Capabilities []string          `yaml:"capabilities"`
	Metadata     map[string]string `yaml:"metadata"`
}

type monitor struct {
	Prometheus  prometheus  `yaml:"prometheus"`
	Pushgateway pushgateway `yaml:"pushgateway"`
	Agent       agent       `yaml:"agent"`
	Webhook     webhook     `yaml:"webhook"`
}

type pushgateway struct {
	URL string `yaml:"url"`
}

type agent struct {
	HeartbeatServerURL string `yaml:"heartbeat_server_url"`
	HeartbeatToken     string `yaml:"heartbeat_token"`
}

type webhook struct {
	Token string `yaml:"token"`
}

type prometheus struct {
	URL string `yaml:"url"`
}

type server struct {
	Address       string `yaml:"address"`
	Model         string `yaml:"model"`
	EnableSwagger bool   `yaml:"enableSwagger"`
	PublicUrl     string `yaml:"publicUrl"`
}

type Db = db

type db struct {
	Dialects      string `yaml:"dialects"`
	Host          string `yaml:"host"`
	Port          int    `yaml:"port"`
	Db            string `yaml:"db"`
	Username      string `yaml:"username"`
	Password      string `yaml:"password"`
	Charset       string `yaml:"charset"`
	SSLMode       string `yaml:"sslMode"`
	MigrationPath string `yaml:"migrationPath"`
	AutoMigrate   bool   `yaml:"autoMigrate"`
	MaxIdle       int    `yaml:"maxIdle"`
	MaxOpen       int    `yaml:"maxOpen"`
}

type redis struct {
	Address  string `yaml:"address"`
	Password string `yaml:"password"`
}

type imageSettings struct {
	UploadDir string `yaml:"uploadDir"`
	ImageHost string `yaml:"imageHost"`
}

type log struct {
	Path  string `yaml:"path"`
	Name  string `yaml:"name"`
	Model string `yaml:"model"`
}

var Config *config

func init() {}

func LoadConfig(configPath string) error {
	if configPath == "" {
		configPath = "./config.yaml"
	}

	yamlFile, err := os.ReadFile(configPath)
	if err != nil {
		return err
	}

	cfg := &config{}
	if err := yaml.Unmarshal(yamlFile, cfg); err != nil {
		return err
	}

	applyEnvOverrides(cfg)
	Config = cfg
	return nil
}

func GetConfig() *db {
	if Config == nil {
		panic("Config is not initialized")
	}
	return &Config.Db
}

func GetRedisConfig() *redis {
	if Config == nil {
		panic("Config is not initialized")
	}
	return &Config.Redis
}

func Setup() {
	if Config == nil {
		panic("Config initialization failed")
	}
}

func applyEnvOverrides(cfg *config) {
	if cfg == nil {
		return
	}

	applyStringEnv("SERVER_ADDRESS", &cfg.Server.Address)
	applyStringEnv("SERVER_MODEL", &cfg.Server.Model)
	applyBoolEnv("SERVER_ENABLE_SWAGGER", &cfg.Server.EnableSwagger)
	applyStringEnv("SERVER_PUBLIC_URL", &cfg.Server.PublicUrl)

	applyStringEnv("DB_DIALECTS", &cfg.Db.Dialects)
	applyStringEnv("DB_HOST", &cfg.Db.Host)
	applyIntEnv("DB_PORT", &cfg.Db.Port)
	applyStringEnv("DB_NAME", &cfg.Db.Db)
	applyStringEnv("DB_USER", &cfg.Db.Username)
	applyStringEnv("DB_PASSWORD", &cfg.Db.Password)
	applyStringEnv("DB_CHARSET", &cfg.Db.Charset)
	applyStringEnv("DB_SSLMODE", &cfg.Db.SSLMode)
	applyStringEnv("DB_MIGRATION_PATH", &cfg.Db.MigrationPath)
	applyBoolEnv("DB_AUTO_MIGRATE", &cfg.Db.AutoMigrate)
	applyIntEnv("DB_MAX_IDLE", &cfg.Db.MaxIdle)
	applyIntEnv("DB_MAX_OPEN", &cfg.Db.MaxOpen)

	applyStringEnv("REDIS_ADDR", &cfg.Redis.Address)
	applyStringEnv("REDIS_PASSWORD", &cfg.Redis.Password)

	applyStringEnv("IMAGE_UPLOAD_DIR", &cfg.ImageSettings.UploadDir)
	applyStringEnv("IMAGE_HOST", &cfg.ImageSettings.ImageHost)

	applyStringEnv("LOG_PATH", &cfg.Log.Path)
	applyStringEnv("LOG_NAME", &cfg.Log.Name)
	applyStringEnv("LOG_MODEL", &cfg.Log.Model)

	applyStringEnv("PROMETHEUS_URL", &cfg.Monitor.Prometheus.URL)
	applyStringEnv("PUSHGATEWAY_URL", &cfg.Monitor.Pushgateway.URL)
	applyStringEnv("HEARTBEAT_SERVER_URL", &cfg.Monitor.Agent.HeartbeatServerURL)
	applyStringEnv("HEARTBEAT_TOKEN", &cfg.Monitor.Agent.HeartbeatToken)
	applyStringEnv("WEBHOOK_TOKEN", &cfg.Monitor.Webhook.Token)
}

func applyStringEnv(key string, target *string) {
	if value := os.Getenv(key); value != "" {
		*target = value
	}
}

func applyIntEnv(key string, target *int) {
	value := os.Getenv(key)
	if value == "" {
		return
	}

	if parsed, err := strconv.Atoi(value); err == nil {
		*target = parsed
	}
}

func applyBoolEnv(key string, target *bool) {
	value := os.Getenv(key)
	if value == "" {
		return
	}

	if parsed, err := strconv.ParseBool(value); err == nil {
		*target = parsed
	}
}
