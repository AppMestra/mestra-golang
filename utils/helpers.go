package utils

import (
	"log"
	"os"
	"time"
)

func LogInfo(message string) {
	log.Printf("[INFO] %s: %s\n", time.Now().Format("2006-01-02 15:04:05"), message)
}

func LogError(message string, err error) {
	log.Printf("[ERROR] %s: %s - %v\n", time.Now().Format("2006-01-02 15:04:05"), message, err)
}

func ValidatePhoneNumber(phone string) bool {
	return len(phone) >= 10 && len(phone) <= 15
}

func GetEnvOrDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
