package config

import (
	"os"
	"strconv"
)

func GetEnv(key string, defaultVal string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}

	return defaultVal
}

func GetEnvAsInt(name string, defaultVal int) int {
	valueStr := GetEnv(name, "")
	if value, err := strconv.Atoi(valueStr); err == nil {
		return value
	}

	return defaultVal
}

func GetEnvAsBool(name string, defaultVal bool) bool {
	valStr := GetEnv(name, "")
	if val, err := strconv.ParseBool(valStr); err == nil {
		return val
	}

	return defaultVal
}

type Config struct {
	DB    DBConf
	Minio MinioConf
}

func New() *Config {
	return &Config{
		DB: DBConf{
			Host:     GetEnv("DBHost", ""),
			Port:     GetEnvAsInt("DBPort", 0),
			User:     GetEnv("DBUser", ""),
			Password: GetEnv("DBPassword", ""),
			DBname:   GetEnv("DBname", ""),
		},
		Minio: MinioConf{
			Endpoint:        GetEnv("MinioEndpoint", ""),
			AccessKeyID:     GetEnv("MinioAccessKeyID", ""),
			SecretAccessKey: GetEnv("MinioSecretAccessKey", ""),
			UseSSL:          GetEnvAsBool("MinioUseSSL", false),
		},
	}
}
