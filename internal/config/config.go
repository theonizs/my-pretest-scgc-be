package config

import (
	"os"
)

type Config struct {
	ServerPort string
}

// LoadConfig อ่านค่าจาก Env หรือใช้ค่า Default
func LoadConfig() *Config {
	return &Config{
		ServerPort: getEnv("PORT", ":8080"),
	}
}

// Helper function: อ่าน Env ถ้าไม่มีให้ใช้ default
func getEnv(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return fallback
}