package database

import (
	"context"
	"fmt"
	"log"
	"os"
	"path/filepath"
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

	fmt.Println("MaxIdleConns:", cfg.DBMaxIdleConns)
	fmt.Println("MaxConnLifetime:", cfg.DBMaxConnLifetime)
	fmt.Println("MaxConns:", cfg.MaxConns)
	fmt.Println("MinConns:", cfg.MinConns)

	config.MaxConnIdleTime = time.Duration(cfg.DBMaxIdleConns) * time.Minute
	config.MaxConnLifetime = time.Duration(cfg.DBMaxConnLifetime) * time.Minute
	// config.MaxConnIdleTime = 10 * time.Minute
	// config.MaxConnLifetime = 60 * time.Minute // Set to 1 hour
	config.MaxConns = int32(cfg.MaxConns)                  // Adjust based on your environment
	config.MinConns = int32(cfg.MinConns)                   // Adjust based on your environment

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

func (db *PgxDatabaseService) GetPool() *pgxpool.Pool {
	return db.pool
}

func (db *PgxDatabaseService) Close() {
	db.pool.Close()
	log.Println("database connection closed")
}

func (db *PgxDatabaseService) Migrate() error {
	return db.executeSQLFiles("migrations")
}

func (db *PgxDatabaseService) executeSQLFiles(directory string) error {
	entries, err := os.ReadDir(directory)
	if err != nil {
		return fmt.Errorf("failed to read directory: %w", err)
	}

	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}

		filePath := filepath.Join(directory, entry.Name())
		content, err := os.ReadFile(filePath)
		if err != nil {
			return fmt.Errorf("failed to read file %s: %w", filePath, err)
		}

		_, err = db.pool.Exec(context.Background(), string(content))
		if err != nil {
			return fmt.Errorf("failed to execute file %s: %w", filePath, err)
		}
	}

	log.Printf("%s files executed successfully", directory)
	return nil
}

func (db *PgxDatabaseService) Seed() error {
	return db.ExecuteSeeders("seeds")
}

func (db *PgxDatabaseService) ExecuteSeeders(directory string) error {
	entries, err := os.ReadDir(directory)
	if err != nil {
		return fmt.Errorf("failed to read directory: %w", err)
	}

	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}

		filePath := filepath.Join(directory, entry.Name())
		content, err := os.ReadFile(filePath)
		if err != nil {
			return fmt.Errorf("failed to read file %s: %w", filePath, err)
		}

		tx, err := db.pool.Begin(context.Background())
		if err != nil {
			return fmt.Errorf("failed to begin transaction: %w", err)
		}
		defer func() {
			if err != nil {
				tx.Rollback(context.Background())
				return
			}
			err = tx.Commit(context.Background())
			if err != nil {
				fmt.Println("Error committing transaction for file", filePath, ":", err)
			}
		}()

		_, err = tx.Exec(context.Background(), string(content))
		if err != nil {
			return fmt.Errorf("failed to execute file %s: %w", filePath, err)
		}
	}

	log.Printf("%s files executed successfully", directory)
	return nil
}

// PoolStats returns the statistics of the connection pool
func (db *PgxDatabaseService) PoolStats() *pgxpool.Stat {
	return db.pool.Stat()
}
