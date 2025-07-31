package infrastructure

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/hydr0g3nz/mini_bank/internal/adapter/repository/gorm/model"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/tracelog"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type SimpleLogger struct{}

func (sl SimpleLogger) Log(ctx context.Context, level tracelog.LogLevel, msg string, data map[string]interface{}) {
	log.Printf("[PGX] %s: %s, Data: %v\n", level.String(), msg, data)
}

// DBConfig holds database connection configuration
type DBConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	DBName   string
	SSLMode  string
}

// ConnectDB creates a database connection pool
func ConnectDB(config *DBConfig) (*gorm.DB, error) {
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=%s",
		config.Host,
		config.User,
		config.Password,
		config.DBName,
		config.Port,
		config.SSLMode,
	)

	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
		logger.Config{
			SlowThreshold:             time.Second, // Slow SQL threshold
			LogLevel:                  logger.Info, // Log level
			IgnoreRecordNotFoundError: true,        // Ignore ErrRecordNotFound error for logger
			Colorful:                  true,        // Disable color
		},
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: newLogger,
	})
	if err != nil {
		return nil, err
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}

	// SetMaxIdleConns sets the maximum number of connections in the idle connection pool
	sqlDB.SetMaxIdleConns(10)
	// SetMaxOpenConns sets the maximum number of open connections to the database
	sqlDB.SetMaxOpenConns(100)
	// SetConnMaxLifetime sets the maximum amount of time a connection may be reused
	sqlDB.SetConnMaxLifetime(time.Hour)

	return db, nil
}

// CloseDB closes the database connection pool
func CloseDB(pool *pgxpool.Pool) {
	if pool != nil {
		pool.Close()
	}
}
func MigrateDB(db *gorm.DB) error {
	log.Println("Running database migrations...")

	// Auto migrate all model
	err := db.AutoMigrate(
		// &model.Hospital{},
		&model.Account{},
		&model.Transaction{},
	)

	if err != nil {
		log.Fatalf("Failed to run migrations: %v", err)
		return err
	}

	log.Println("Database migrations completed successfully")
	return nil
}
