package database

import (
	"fmt"
	"log"      // Standard log for initial setup errors before slog is ready
	"log/slog" // Use slog for operational logging
	"os"
	"sync"
	"time"

	"github.com/saengepch-sys/iple-backend/internal/config"
	"github.com/saengepch-sys/iple-backend/internal/models"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var (
	DB *gorm.DB
	once sync.Once
)

// ConnectDB establishes a connection to the database using the loaded config.
// It ensures the connection is attempted only once using sync.Once.
func ConnectDB() {
	once.Do(func() {
		var err error
		dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=%s TimeZone=%s",
			config.Env.DBHost,
			config.Env.DBUser,
			config.Env.DBPassword,
			config.Env.DBName,
			config.Env.DBPort,
			config.Env.DBSslMode,
			config.Env.DBTimezone,
		)

		// Configure GORM Logger
		gormLogLevel := logger.Warn
		if config.Env.LogLevel == "debug" || config.Env.LogLevel == "info" { // Show Info level logs in debug/info mode
			gormLogLevel = logger.Info
		}

		newLogger := logger.New(
			log.New(os.Stdout, "\r\n", log.LstdFlags), // Use Go's standard logger for GORM's output destination
			logger.Config{
				SlowThreshold:             250 * time.Millisecond, // Log queries slower than 250ms
				LogLevel:                  gormLogLevel,
				IgnoreRecordNotFoundError: true,  // Don't log ErrRecordNotFound as Error
				Colorful:                  false, // Disable color in logs for consistency
				ParameterizedQueries:      true, // Log queries with parameters hidden (safer)
			},
		)

		// Attempt to connect
		retryCount := 0
		maxRetries := 5
		retryDelay := 2 * time.Second

		for retryCount < maxRetries {
			DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{
				Logger: newLogger,
				// PrepareStmt: true, // Can improve performance but check compatibility
				// DisableForeignKeyConstraintWhenMigrating: true, // Use only if necessary during migrations
			})

			if err == nil {
                // Connection successful, verify with Ping
                sqlDB, pingErr := DB.DB()
                if pingErr != nil {
                    slog.Error("Failed to get underlying sql.DB", "error", pingErr)
                     os.Exit(1) // Exit if cannot get DB instance
                }
                if pingErr = sqlDB.Ping(); pingErr == nil {
                    slog.Info("Database connection successful and verified.")
                    break // Exit retry loop on successful ping
                } else {
                     slog.Warn("Database connection opened but ping failed, retrying...", "error", pingErr, "attempt", retryCount+1)
                     DB = nil // Reset DB instance
                }

			}

            // Connection failed or ping failed
			slog.Warn("Failed to connect to database, retrying...", "attempt", retryCount+1, "error", err)
			retryCount++
			if retryCount < maxRetries {
				time.Sleep(retryDelay)
                retryDelay *= 2 // Exponential backoff
			} else {
                 slog.Error("Failed to connect to database after multiple retries", "attempts", maxRetries, "last_error", err)
                 os.Exit(1) // Exit application if connection ultimately fails
            }

		}


		// Configure Connection Pool
		sqlDB, err := DB.DB()
		if err != nil {
			slog.Error("Failed to get generic database object after connection", "error", err)
			os.Exit(1)
		}
		sqlDB.SetMaxIdleConns(10)          // Adjust based on expected idle load
		sqlDB.SetMaxOpenConns(50)         // Adjust based on expected peak load & DB limits
		sqlDB.SetConnMaxLifetime(time.Hour) // Reuse connections for up to an hour
        sqlDB.SetConnMaxIdleTime(10 * time.Minute) // Close connections idle for 10 mins

		slog.Info("Database connection pool configured",
            "max_idle", 10,
            "max_open", 50,
            "max_lifetime", time.Hour,
            "max_idletime", 10*time.Minute,
        )

		// --- Auto Migration (Development ONLY!) ---
		// Run migrations using a dedicated tool in production (Step 1.5b)
		if config.Env.Env == "development" {
			runAutoMigrateDev()
		}
	})
}

// GetDB returns the global database instance. Exits if not connected.
func GetDB() *gorm.DB {
	if DB == nil {
		// Should not happen if main calls ConnectDB and exits on failure
		slog.Error("FATAL: Database instance is nil. ConnectDB must be called successfully first.")
        os.Exit(1)
	}
	return DB
}

// runAutoMigrateDev runs GORM's AutoMigrate in development environments.
// WARNING: DO NOT USE IN PRODUCTION. Use a proper migration tool.
func runAutoMigrateDev() {
	slog.Warn("Running Auto Migration (Development Mode ONLY!)...")
	err := DB.AutoMigrate(
        // Add ALL your GORM models here for auto-migration in dev
		&models.User{}, &models.Course{}, &models.Section{}, &models.Lesson{},
		&models.Enrollment{}, &models.Submission{}, &models.Grade{},
		&models.ParentChildLink{}, &models.PDPAConsentLog{},
        // ... other models ...
	)
	if err != nil {
		slog.Error("Failed to auto migrate database", "error", err)
		os.Exit(1) // Exit if auto-migration fails in dev, as it indicates a model issue
	}
	slog.Info("Auto Migration Completed (Development Mode)")
}