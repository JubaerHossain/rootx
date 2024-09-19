package database

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/JubaerHossain/rootx/pkg/core/config"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PgxDatabaseService struct {
	pool *pgxpool.Pool
}

// NewPgxDatabaseService initializes a new database service using pgxpool
func NewPgxDatabaseService(cfg *config.Config) (*PgxDatabaseService, error) {
	dsn := fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=disable",
		cfg.DBUser, cfg.DBPassword, cfg.DBHost, cfg.DBPort, cfg.DBName)

	config, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to parse database config: %w", err)
	}

	config.MaxConnIdleTime = time.Duration(cfg.DBMaxIdleConns) * time.Minute
	config.MaxConnLifetime = cfg.DBMaxConnLifetime
	config.MaxConns = int32(cfg.MaxConns) // Adjust based on your environment
	config.MinConns = int32(cfg.MinConns) // Adjust based on your environment

	pool, err := pgxpool.NewWithConfig(context.Background(), config)
	if err != nil {
		return nil, fmt.Errorf("failed to create database connection pool: %w", err)
	}

	// Test the database connection
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	for i := 0; i < 3; i++ { // Retry logic
		err := pool.Ping(ctx)
		if err == nil {
			log.Println("connected to database")
			break
		}
		log.Printf("failed to ping database: %v (attempt %d)", err, i+1)
		time.Sleep(2 * time.Second) // Wait before retrying
		if i == 2 {                 // After the last attempt
			return nil, fmt.Errorf("failed to ping database after multiple attempts: %w", err)
		}
	}

	return &PgxDatabaseService{pool: pool}, nil
}

func (db *PgxDatabaseService) GetDB() *pgxpool.Pool {
	return db.pool
}

func (db *PgxDatabaseService) Close() {
	db.pool.Close()
	log.Println("database connection closed")
}

func (db *PgxDatabaseService) LogStats() {
	stats := db.pool.Stat()
	log.Printf("Total Connections: %d\n", stats.TotalConns())
	log.Printf("Idle Connections: %d\n", stats.IdleConns())
	log.Printf("Max Connections: %d\n", stats.MaxConns())
}

func (db *PgxDatabaseService) HealthCheck() error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := db.pool.Ping(ctx); err != nil {
		return fmt.Errorf("database health check failed: %w", err)
	}
	return nil
}
