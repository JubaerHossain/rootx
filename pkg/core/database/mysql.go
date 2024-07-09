package database

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	"github.com/JubaerHossain/rootx/pkg/core/config"
	_ "github.com/go-sql-driver/mysql"
)

// MySQLService manages MySQL database connections
type MySQLService struct {
	db *sql.DB
}

// NewMySQLService creates a new instance of MySQLService with advanced configurations
func NewMySQLService(cfg *config.Config) (*MySQLService, error) {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s",
		cfg.DBUser, cfg.DBPassword, cfg.DBHost, cfg.DBPort, cfg.DBName)

	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to open database connection: %w", err)
	}

	// Configure connection pooling
	db.SetMaxOpenConns(cfg.MaxConns)                                          // Maximum number of open connections to the database
	db.SetMaxIdleConns(cfg.DBMaxIdleConns)                                    // Maximum number of idle connections in the pool
	db.SetConnMaxLifetime(time.Duration(cfg.DBMaxConnLifetime) * time.Minute) // Maximum amount of time a connection may be reused

	// Test the database connection
	if err := testConnection(db); err != nil {
		return nil, fmt.Errorf("database connection test failed: %w", err)
	}

	return &MySQLService{db: db}, nil
}

// testConnection attempts to ping the database with retry logic
func testConnection(db *sql.DB) error {
	for i := 0; i < 3; i++ { // Retry logic
		if err := db.Ping(); err == nil {
			log.Println("connected to database")
			return nil
		}
		log.Printf("failed to ping database"+" (attempt %d)", i+1)
		time.Sleep(2 * time.Second) // Wait before retrying
	}

	return fmt.Errorf("failed to ping database after multiple attempts")
}

// GetDB returns the underlying *sql.DB instance
func (dbService *MySQLService) GetDB() *sql.DB {
	return dbService.db
}

// Close closes the database connection gracefully
func (dbService *MySQLService) Close() {
	if dbService.db != nil {
		dbService.db.Close()
		log.Println("database connection closed")
	}
}

// HealthCheck performs a health check by pinging the database
func (dbService *MySQLService) HealthCheck() error {
	if err := dbService.db.Ping(); err != nil {
		return fmt.Errorf("database health check failed: %w", err)
	}
	return nil
}

// LogStats logs the current connection pool statistics
func (dbService *MySQLService) LogStats() {
	stats := dbService.db.Stats()
	log.Printf("MySQL connection pool stats: %+v", stats)
}
