package bootstrap

import (
	"log"
	"os"
	"strconv"

	"github.com/spf13/viper"
)

type Env struct {
	AppEnv                 string `mapstructure:"APP_ENV"`
	ServerAddress          string `mapstructure:"SERVER_ADDRESS"`
	ContextTimeout         int    `mapstructure:"CONTEXT_TIMEOUT"`
	DBHost                 string `mapstructure:"DB_HOST"`
	DBPort                 string `mapstructure:"DB_PORT"`
	DBUser                 string `mapstructure:"DB_USER"`
	DBPass                 string `mapstructure:"DB_PASS"`
	DBName                 string `mapstructure:"DB_NAME"`
	AccessTokenExpiryHour  int    `mapstructure:"ACCESS_TOKEN_EXPIRY_HOUR"`
	RefreshTokenExpiryHour int    `mapstructure:"REFRESH_TOKEN_EXPIRY_HOUR"`
	AccessTokenSecret      string `mapstructure:"ACCESS_TOKEN_SECRET"`
	RefreshTokenSecret     string `mapstructure:"REFRESH_TOKEN_SECRET"`
}

func NewEnv() *Env {
	env := Env{}

	// 尝试读取 .env 文件（用于本地开发）
	viper.SetConfigFile(".env")
	if err := viper.ReadInConfig(); err != nil {
		log.Println("No .env file found, using system environment variables")
		// 从系统环境变量读取配置
		env = loadFromSystemEnv()
	} else {
		// 从 .env 文件读取配置
		if err := viper.Unmarshal(&env); err != nil {
			log.Fatal("Environment can't be loaded: ", err)
		}
	}

	if env.AppEnv == "development" {
		log.Println("The App is running in development env")
	} else if env.AppEnv == "production" {
		log.Println("The App is running in production env")
	}

	return &env
}

// loadFromSystemEnv 从系统环境变量加载配置（用于 Docker 部署）
func loadFromSystemEnv() Env {
	return Env{
		AppEnv:                 getEnv("APP_ENV", "production"),
		ServerAddress:          getEnv("SERVER_ADDRESS", "0.0.0.0:8080"),
		ContextTimeout:         getEnvAsInt("CONTEXT_TIMEOUT", 30),
		DBHost:                 getEnv("DB_HOST", "localhost"),
		DBPort:                 getEnv("DB_PORT", "27017"),
		DBUser:                 getEnv("DB_USER", ""),
		DBPass:                 getEnv("DB_PASS", ""),
		DBName:                 getEnv("DB_NAME", "flow_link"),
		AccessTokenExpiryHour:  getEnvAsInt("ACCESS_TOKEN_EXPIRY_HOUR", 24),
		RefreshTokenExpiryHour: getEnvAsInt("REFRESH_TOKEN_EXPIRY_HOUR", 168),
		AccessTokenSecret:      getEnv("ACCESS_TOKEN_SECRET", ""),
		RefreshTokenSecret:     getEnv("REFRESH_TOKEN_SECRET", ""),
	}
}

// getEnv 获取环境变量，如果不存在则返回默认值
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// getEnvAsInt 获取环境变量并转换为 int
func getEnvAsInt(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if intVal, err := strconv.Atoi(value); err == nil {
			return intVal
		}
	}
	return defaultValue
}
