package config

import (
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

var AppConfig *Config

// Config holds all configuration for the application
type Config struct {
	Server ServerConfig
	MySql  MySqlConfig
	Jwt    JwtConfig
	Log    LogConfig
}

type ServerConfig struct {
	Port string
}

type MySqlConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	Name     string
}

type JwtConfig struct {
	Secret string
}

type LogConfig struct {
	Level          string
	Type           string
	LogFileEnabled bool
	LogFilePath    string
}

// LoadConfig loads configuration from environment variables
func LoadConfig() {
	// Load .env file if it exists
	if err := godotenv.Load(); err != nil {
		log.Println("Warning: .env file not found")
	}

	AppConfig = &Config{
		Server: ServerConfig{
			Port: getEnv("PORT", "8080"),
		},
		MySql: MySqlConfig{
			Host:     getEnv("DB_HOST", "localhost"),
			Port:     getEnv("DB_PORT", "3306"),
			User:     getEnv("DB_USER", "root"),
			Password: getEnv("DB_PASS", "rootpassword"),
			Name:     getEnv("DB_NAME", "blog-db"),
		},
		Jwt: JwtConfig{
			Secret: getEnv("JWT_SECRET_KEY", "secret"),
		},
		Log: LogConfig{
			Level:       getEnv("LOG_LEVEL", "debug"),
			Type:        getEnv("LOG_TYPE", "json"),
			LogFilePath: getEnv("LOG_FILE_PATH", "logs/app.log"),
		},
	}

	AppConfig.Log.LogFileEnabled, _ = strconv.ParseBool(getEnv("LOG_FILE_ENABLED", "true"))
}

// Helper function to get environment variable with a default value
func getEnv(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return fallback
}
