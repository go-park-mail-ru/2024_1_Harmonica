package config

import (
	"os"
	"strconv"
)

func getEnv(key string, defaultVal string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}

	return defaultVal
}

func getEnvAsInt(name string, defaultVal int) int {
	valueStr := getEnv(name, "")
	if value, err := strconv.Atoi(valueStr); err == nil {
		return value
	}

	return defaultVal
}

func getEnvAsBool(name string, defaultVal bool) bool {
	valStr := getEnv(name, "")
	if val, err := strconv.ParseBool(valStr); err == nil {
		return val
	}

	return defaultVal
}

type DBConf struct {
	Host     string
	Port     int
	User     string
	Password string
	DBname   string
}

type Config struct {
	DB DBConf
}

func New() *Config {
	return &Config{
		DB: DBConf{
			Host:     getEnv("DBHost", ""),
			Port:     getEnvAsInt("DBPort", 0),
			User:     getEnv("DBUser", ""),
			Password: getEnv("DBPassword", ""),
			DBname:   getEnv("DBname", ""),
		},
	}
}
