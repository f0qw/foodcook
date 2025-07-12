package config

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/spf13/viper"
)

type Config struct {
	App      AppConfig      `mapstructure:"app"`
	Database DatabaseConfig `mapstructure:"database"`
	Redis    RedisConfig    `mapstructure:"redis"`
	JWT      JWTConfig      `mapstructure:"jwt"`
	Upload   UploadConfig   `mapstructure:"upload"`
	CORS     CORSConfig     `mapstructure:"cors"`
}

type AppConfig struct {
	Name    string `mapstructure:"name"`
	Version string `mapstructure:"version"`
	Port    int    `mapstructure:"port"`
	Mode    string `mapstructure:"mode"`
}

type DatabaseConfig struct {
	Host            string `mapstructure:"host"`
	Port            int    `mapstructure:"port"`
	User            string `mapstructure:"user"`
	Password        string `mapstructure:"password"`
	Name            string `mapstructure:"name"`
	Charset         string `mapstructure:"charset"`
	ParseTime       bool   `mapstructure:"parse_time"`
	Loc             string `mapstructure:"loc"`
	MaxIdleConns    int    `mapstructure:"max_idle_conns"`
	MaxOpenConns    int    `mapstructure:"max_open_conns"`
	ConnMaxLifetime int    `mapstructure:"conn_max_lifetime"`
}

type RedisConfig struct {
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	Password string `mapstructure:"password"`
	DB       int    `mapstructure:"db"`
	PoolSize int    `mapstructure:"pool_size"`
}

type JWTConfig struct {
	Secret      string `mapstructure:"secret"`
	ExpireHours int    `mapstructure:"expire_hours"`
}

type UploadConfig struct {
	MaxSize      int64    `mapstructure:"max_size"`
	AllowedTypes []string `mapstructure:"allowed_types"`
	UploadPath   string   `mapstructure:"upload_path"`
}

type CORSConfig struct {
	AllowedOrigins   []string `mapstructure:"allowed_origins"`
	AllowedMethods   []string `mapstructure:"allowed_methods"`
	AllowedHeaders   []string `mapstructure:"allowed_headers"`
	AllowCredentials bool     `mapstructure:"allow_credentials"`
}

var GlobalConfig *Config

func LoadConfig() error {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("./configs")
	viper.AddConfigPath(".")

	// 从环境变量读取配置
	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	// 设置默认值
	setDefaults()

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return fmt.Errorf("failed to read config file: %w", err)
		}
		log.Println("No config file found, using environment variables and defaults")
	}

	GlobalConfig = &Config{}
	if err := viper.Unmarshal(GlobalConfig); err != nil {
		return fmt.Errorf("failed to unmarshal config: %w", err)
	}

	// 从环境变量覆盖数据库配置
	if host := os.Getenv("DB_HOST"); host != "" {
		GlobalConfig.Database.Host = host
	}
	if port := os.Getenv("DB_PORT"); port != "" {
		// 这里可以添加端口解析逻辑
	}
	if user := os.Getenv("DB_USER"); user != "" {
		GlobalConfig.Database.User = user
	}
	if password := os.Getenv("DB_PASSWORD"); password != "" {
		GlobalConfig.Database.Password = password
	}
	if name := os.Getenv("DB_NAME"); name != "" {
		GlobalConfig.Database.Name = name
	}

	// 从环境变量覆盖Redis配置
	if redisHost := os.Getenv("REDIS_HOST"); redisHost != "" {
		GlobalConfig.Redis.Host = redisHost
	}
	if redisPort := os.Getenv("REDIS_PORT"); redisPort != "" {
		// 这里可以添加端口解析逻辑
	}
	if redisPassword := os.Getenv("REDIS_PASSWORD"); redisPassword != "" {
		GlobalConfig.Redis.Password = redisPassword
	}

	// 从环境变量覆盖JWT配置
	if jwtSecret := os.Getenv("JWT_SECRET"); jwtSecret != "" {
		GlobalConfig.JWT.Secret = jwtSecret
	}

	return nil
}

func setDefaults() {
	viper.SetDefault("app.name", "foodcook")
	viper.SetDefault("app.version", "1.0.0")
	viper.SetDefault("app.port", 8080)
	viper.SetDefault("app.mode", "debug")

	viper.SetDefault("database.host", "127.0.0.1")
	viper.SetDefault("database.port", 3306)
	viper.SetDefault("database.user", "root")
	viper.SetDefault("database.password", "root")
	viper.SetDefault("database.name", "foodcook")
	viper.SetDefault("database.charset", "utf8mb4")
	viper.SetDefault("database.parse_time", true)
	viper.SetDefault("database.loc", "Local")
	viper.SetDefault("database.max_idle_conns", 10)
	viper.SetDefault("database.max_open_conns", 100)
	viper.SetDefault("database.conn_max_lifetime", 3600)

	viper.SetDefault("redis.host", "127.0.0.1")
	viper.SetDefault("redis.port", 6379)
	viper.SetDefault("redis.password", "")
	viper.SetDefault("redis.db", 0)
	viper.SetDefault("redis.pool_size", 10)

	viper.SetDefault("jwt.secret", "your-secret-key-change-in-production")
	viper.SetDefault("jwt.expire_hours", 24)

	viper.SetDefault("upload.max_size", 10485760)
	viper.SetDefault("upload.allowed_types", []string{"image/jpeg", "image/png", "image/gif"})
	viper.SetDefault("upload.upload_path", "./uploads")

	viper.SetDefault("cors.allowed_origins", []string{"http://localhost:3000", "http://localhost:5173"})
	viper.SetDefault("cors.allowed_methods", []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"})
	viper.SetDefault("cors.allowed_headers", []string{"*"})
	viper.SetDefault("cors.allow_credentials", true)
}

func GetConfig() *Config {
	return GlobalConfig
}
