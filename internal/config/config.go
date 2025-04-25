package config

import (
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/joho/godotenv"
)

type Config struct {
	Port                  string
	LogLevel              string // debug, info, warn, error
	Env                   string // development, staging, production
	AllowedOrigins        string // Comma-separated list of allowed CORS origins
	GracefulShutdownTimeout time.Duration // เพิ่ม Timeout สำหรับ Shutdown

	DBHost     string
	DBPort     string
	DBUser     string
	DBPassword string // Loaded from secret manager in prod ideally
	DBName     string
	DBSslMode  string
	DBTimezone string

	JWTSecretKey  string // Loaded from secret manager in prod ideally
	JWTExpiresIn time.Duration
}

var Env *Config

// LoadConfig loads configuration from environment variables or .env file.
func LoadConfig() {
	// Load .env file only if ENV is not 'production' or not set
	appEnv := getEnv("ENV", "development") // Default to development
	if appEnv != "production" {
		if err := godotenv.Load(); err != nil {
			log.Println("WARN: No .env file found or error loading .env file")
		}
	}

	port := getEnv("PORT", "8080")
	logLevel := getEnv("LOG_LEVEL", "info")
	allowedOrigins := getEnv("ALLOWED_ORIGINS", "http://localhost:3000") // Default for local Next.js dev
    shutdownTimeoutStr := getEnv("GRACEFUL_SHUTDOWN_TIMEOUT_SECONDS", "30")

	dbHost := getEnv("DB_HOST", "localhost")
	dbPort := getEnv("DB_PORT", "5432")
	dbUser := getEnv("DB_USER", "postgres")
	dbPassword := getEnv("DB_PASSWORD", "123456") // *** MUST be set via ENV or Secret Manager in Prod ***
	dbName := getEnv("DB_NAME", "iple_lms_db")
	dbSslMode := getEnv("DB_SSLMODE", "disable")
	dbTimezone := getEnv("DB_TIMEZONE", "Asia/Bangkok")

	jwtSecret := getEnv("JWT_SECRET_KEY", "!!!change_this_very_secret_key_in_production!!!") // *** MUST be set via ENV or Secret Manager in Prod ***
	jwtExpiresMinutesStr := getEnv("JWT_EXPIRES_IN_MINUTES", "60")

	jwtExpiresMinutes, err := strconv.Atoi(jwtExpiresMinutesStr)
	if err != nil || jwtExpiresMinutes <= 0 {
		jwtExpiresMinutes = 60
	}
    shutdownTimeoutSec, err := strconv.Atoi(shutdownTimeoutStr)
    if err != nil || shutdownTimeoutSec <= 0 {
        shutdownTimeoutSec = 30
    }

	// TODO: Implement logic to load DBPassword and JWTSecretKey from GCP Secret Manager if appEnv == "production"
	// This would typically involve using the secretmanager client library and the service account credentials.
	if appEnv == "production" {
		if dbPassword == "" {
			log.Fatal("FATAL: DB_PASSWORD environment variable not set in production!")
		}
		if jwtSecret == "!!!change_this_very_secret_key_in_production!!!" {
			log.Fatal("FATAL: JWT_SECRET_KEY environment variable not set or using default in production!")
		}
		// Example placeholder for fetching from Secret Manager (needs implementation)
		// dbPassword = fetchSecret("projects/YOUR_PROJECT_ID/secrets/db-password/versions/latest")
		// jwtSecret = fetchSecret("projects/YOUR_PROJECT_ID/secrets/jwt-secret/versions/latest")
	}

	Env = &Config{
		Port:                  port,
		LogLevel:              strings.ToLower(logLevel),
		Env:                   appEnv,
        AllowedOrigins:        allowedOrigins,
        GracefulShutdownTimeout: time.Duration(shutdownTimeoutSec) * time.Second,
		DBHost:                dbHost,
		DBPort:                dbPort,
		DBUser:                dbUser,
		DBPassword:            dbPassword,
		DBName:                dbName,
		DBSslMode:             dbSslMode,
		DBTimezone:            dbTimezone,
		JWTSecretKey:          jwtSecret,
		JWTExpiresIn:          time.Duration(jwtExpiresMinutes) * time.Minute,
	}

	// Log only non-sensitive config values on startup
	log.Printf("Configuration loaded: Port=%s, LogLevel=%s, Env=%s", Env.Port, Env.LogLevel, Env.Env)
}

// getEnv retrieves an environment variable or returns a fallback value.
// Logs a warning if fallback is used (except for DB_PASSWORD/JWT_SECRET).
func getEnv(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists && value != "" {
		return value
	}
	// Avoid logging fallback for sensitive keys
	if key != "DB_PASSWORD" && key != "JWT_SECRET_KEY" {
		log.Printf("WARN: Environment variable %s not set, using fallback: %s", key, fallback)
	} else if fallback != "" {
        log.Printf("WARN: Environment variable %s not set, using configured fallback (should not happen in prod!)", key)
    } else {
         log.Printf("WARN: Environment variable %s not set, no fallback provided.", key)
    }

	return fallback
}

// GetAllowedOrigins parses the comma-separated ALLOWED_ORIGINS string.
func GetAllowedOrigins() string {
    // In a real app, you might want more robust parsing or return a slice.
    // Fiber's CORS middleware might handle comma-separated string directly.
    return Env.AllowedOrigins
}

/*
// Example function placeholder for fetching secrets (needs implementation)
func fetchSecret(secretName string) string {
    ctx := context.Background()
    client, err := secretmanager.NewClient(ctx)
    if err != nil {
        log.Fatalf("Failed to create secretmanager client: %v", err)
    }
    defer client.Close()

    req := &secretmanagerpb.AccessSecretVersionRequest{
        Name: secretName,
    }

    result, err := client.AccessSecretVersion(ctx, req)
    if err != nil {
        log.Fatalf("Failed to access secret version %s: %v", secretName, err)
    }

    // WARNING: Do not log the secret payload!
    // log.Printf("Payload for secret %s: %s", secretName, string(result.Payload.Data))
    return string(result.Payload.Data)
}
*/