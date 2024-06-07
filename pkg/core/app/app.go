package app

import (
	"context"
	"fmt"
	"io/fs"
	"net/http"
	"os"        // Add this import
	"os/signal" // Add this import
	"syscall"   // Add this import
	"time"

	"github.com/JubaerHossain/rootx/pkg/core/cache"
	"github.com/JubaerHossain/rootx/pkg/core/config"
	"github.com/JubaerHossain/rootx/pkg/core/database"
	"github.com/JubaerHossain/rootx/pkg/core/logger"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"
)

// App represents the application struct
type App struct {
	Config       *config.Config
	HttpServer   *http.Server
	BuildVersion string
	HttpPort     int
	PublicFS     fs.FS
	Cache        cache.CacheService
	DB           *pgxpool.Pool
	Logger       *zap.Logger
}

// NewApp creates a new instance of the App struct
func NewApp(cfg *config.Config) *App {
	return &App{Config: cfg}
}

// StartApp initializes and starts the application
func StartApp() (*App, error) {
	// Initialize logger
	if err := logger.Init(); err != nil {
		return nil, fmt.Errorf("error initializing logger: %w", err)
	}

	// Load configuration
	// Load configuration
	cfg, err := config.LoadConfig()
	if err != nil {
		return nil, fmt.Errorf("error loading config: %w", err)
	}

	// Initialize database and cache asynchronously
	dbPool, err := initDatabase(cfg)
	if err != nil {
		return nil, err
	}

	cacheService, err := initCache()
	if err != nil {
		return nil, err
	}
	app := &App{
		Config:       cfg,
		HttpPort:     cfg.AppPort,
		BuildVersion: cfg.BuildVersion,
		Cache:        cacheService,
		DB:           dbPool,
		Logger:       logger.Logger,
	}

	return app, nil
}

// initDatabase initializes the database
func initDatabase(cfg *config.Config) (*pgxpool.Pool, error) {
	dbService, err := database.NewPgxDatabaseService()
	if err != nil {
		return nil, err
	}

	if cfg.Migrate {
		if err := dbService.Migrate(); err != nil {
			logger.Logger.Error("Failed to migrate database", zap.Error(err))
			return nil, err
		}
	}

	if cfg.Seed {
		if err := dbService.Seed(); err != nil {
			logger.Logger.Error("Failed to seed database", zap.Error(err))
			return nil, err
		}
	}

	return dbService.GetPool(), nil
}

// initCache initializes the cache
func initCache() (cache.CacheService, error) {
	ctx := context.Background()
	cacheService, err := cache.NewRedisCacheService(ctx)
	if err != nil {
		return nil, err
	}
	return cacheService, nil
}

// StartServer starts the HTTP server
func (app *App) StartServer() error {
	// Start HTTP server in a goroutine
	go func() {
		if err := app.HttpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			app.Logger.Error("Could not start server", zap.Error(err))
		}
	}()

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	<-quit

	app.Logger.Info("Server is shutting down")
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := app.HttpServer.Shutdown(ctx); err != nil {
		app.Logger.Error("Could not gracefully shutdown server", zap.Error(err))
	}

	if err := app.closeResources(); err != nil {
		app.Logger.Error("Failed to close resources", zap.Error(err))
	}

	app.Logger.Info("Server stopped")
	return nil
}

// closeResources closes resources like database connections, cache, etc.
func (app *App) closeResources() error {
	if app.Cache != nil {
		if err := app.Cache.Close(); err != nil {
			return fmt.Errorf("failed to close cache: %w", err)
		}
	}
	if app.DB != nil {
		app.DB.Close()
	}
	return nil
}
