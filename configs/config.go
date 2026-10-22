package configs

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	API_PREFIX                    string
	CF_ACCOUNT_ID                 string
	CF_API_TOKEN                  string
	CF_KV_CONFIG_KEY              string
	CF_KV_NAMESPACE_ID            string
	CONFIG_API_KEY                string
	CONFIG_PUBLIC_URL             string
	DATABASE_URL                  string
	DEFAULT_SALARY_AMOUNT         string
	DEFAULT_SALARY_DAY            string
	DEFAULT_TIMEZONE              string
	FIREBASE_SERVICE_ACCOUNT_JSON string
	FIREBASE_SERVICE_ACCOUNT_PATH string
	GEMINI_API_KEY                string
	GEMINI_MODEL                  string
	GOOGLE_CLIENT_ID              string
	JWT_EXPIRES_IN                string
	JWT_REFRESH_EXPIRES_IN        string
	JWT_REFRESH_SECRET            string
	JWT_SECRET                    string
	NODE_ENV                      string
	PORT                          string
	REDIS_HOST                    string
	REDIS_PASSWORD                string
	REDIS_PORT                    string
	THROTTLE_LIMIT                string
	THROTTLE_TTL                  string
}

var AppConfig Config

func LoadConfig() {

	if err := godotenv.Load(); err != nil {
		log.Printf("No .env file found: %v", err)
	}

	AppConfig = Config{
		API_PREFIX:                    getEnv("API_PREFIX", "api/v1"),
		CF_ACCOUNT_ID:                 getEnv("CF_ACCOUNT_ID", ""),
		CF_API_TOKEN:                  getEnv("CF_API_TOKEN", ""),
		CF_KV_CONFIG_KEY:              getEnv("CF_KV_CONFIG_KEY", "remote_config"),
		CF_KV_NAMESPACE_ID:            getEnv("CF_KV_NAMESPACE_ID", ""),
		CONFIG_API_KEY:                getEnv("CONFIG_API_KEY", ""),
		CONFIG_PUBLIC_URL:             getEnv("CONFIG_PUBLIC_URL", "https://expent-config.adityaworkhbss.workers.dev/"),
		DATABASE_URL:                  getEnv("DATABASE_URL", ""),
		DEFAULT_SALARY_AMOUNT:         getEnv("DEFAULT_SALARY_AMOUNT", "87500"),
		DEFAULT_SALARY_DAY:            getEnv("DEFAULT_SALARY_DAY", "10"),
		DEFAULT_TIMEZONE:              getEnv("DEFAULT_TIMEZONE", "Asia/Kolkata"),
		FIREBASE_SERVICE_ACCOUNT_JSON: getEnv("FIREBASE_SERVICE_ACCOUNT_JSON", ""),
		FIREBASE_SERVICE_ACCOUNT_PATH: getEnv("FIREBASE_SERVICE_ACCOUNT_PATH", ""),
		GEMINI_API_KEY:                getEnv("GEMINI_API_KEY", ""),
		GEMINI_MODEL:                  getEnv("GEMINI_MODEL", "gemini-3.5-flash"),
		GOOGLE_CLIENT_ID:              getEnv("GOOGLE_CLIENT_ID", ""),
		JWT_EXPIRES_IN:                getEnv("JWT_EXPIRES_IN", "15m"),
		JWT_REFRESH_EXPIRES_IN:        getEnv("JWT_REFRESH_EXPIRES_IN", "7d"),
		JWT_REFRESH_SECRET:            getEnv("JWT_REFRESH_SECRET", "change-in-production"),
		JWT_SECRET:                    getEnv("JWT_SECRET", "change-in-production"),
		NODE_ENV:                      getEnv("NODE_ENV", "development"),
		PORT:                          getEnv("PORT", "3000"),
		REDIS_HOST:                    getEnv("REDIS_HOST", "localhost"),
		REDIS_PASSWORD:                getEnv("REDIS_PASSWORD", ""),
		REDIS_PORT:                    getEnv("REDIS_PORT", "6379"),
		THROTTLE_LIMIT:                getEnv("THROTTLE_LIMIT", "100"),
		THROTTLE_TTL:                  getEnv("THROTTLE_TTL", "60000"),
	}
}

func getEnv(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return fallback
}
