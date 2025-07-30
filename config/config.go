package config

import (
	"os"
	"strconv"
)

type Config struct {
	Database DatabaseConfig
	API      APIConfig
	Server   ServerConfig
}

type DatabaseConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	Name     string
	SSLMode  string
}

type APIConfig struct {
	DotNetURL string
}

type ServerConfig struct {
	Port     string
	LogLevel string
}

var AppConfig *Config

func LoadConfig() {
	AppConfig = &Config{
		Database: DatabaseConfig{
			Host:     getEnvOrDefault("DB_HOST", "localhost"),
			Port:     getEnvOrDefault("DB_PORT", "5432"),
			User:     getEnvOrDefault("DB_USER", "mestra_user"),
			Password: getEnvOrDefault("DB_PASSWORD", "mestra_pass"),
			Name:     getEnvOrDefault("DB_NAME", "mestra_db"),
			SSLMode:  getEnvOrDefault("DB_SSLMODE", "disable"),
		},
		API: APIConfig{
			DotNetURL: getEnvOrDefault("DOTNET_API_URL", "http://localhost:5000/api/messages"),
		},
		Server: ServerConfig{
			Port:     getEnvOrDefault("PORT", "8080"),
			LogLevel: getEnvOrDefault("LOG_LEVEL", "info"),
		},
	}
}

// GetPort retorna a porta do servidor
func GetPort() string {
	return ":" + AppConfig.Server.Port
}

// IsProductionMode verifica se está em produção
func IsProductionMode() bool {
	env := os.Getenv("GIN_MODE")
	return env == "release"
}

// GetIntEnv retorna uma variável de ambiente como int
func GetIntEnv(key string, defaultValue int) int {
	str := os.Getenv(key)
	if str == "" {
		return defaultValue
	}

	value, err := strconv.Atoi(str)
	if err != nil {
		return defaultValue
	}

	return value
}
