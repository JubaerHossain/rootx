package app

import (
	"context"
	"database/sql"
	"fmt"
	"io/fs"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/JubaerHossain/rootx/pkg/core/cache"
	"github.com/JubaerHossain/rootx/pkg/core/config"
	"github.com/JubaerHossain/rootx/pkg/core/database"
	"github.com/JubaerHossain/rootx/pkg/core/filesystem"
	"github.com/JubaerHossain/rootx/pkg/core/logger"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"
	"golang.org/x/sync/errgroup"
)

type App struct {
	Config       *config.Config
	HttpServer   *http.Server
	BuildVersion string
	HttpPort     int
	PublicFS     fs.FS
	Cache        cache.CacheService
	DB           *pgxpool.Pool
	MDB          *sql.DB
	Logger       *zap.Logger
	FileUpload   *filesystem.FileUploadService
}

func NewApp(cfg *config.Config) *App {
	return &App{Config: cfg}
}

func StartApp() (*App, error) {
	if err := logger.Init(); err != nil {
		return nil, fmt.Errorf("error initializing logger: %w", err)
	}

	cfg, err := config.LoadConfig()
	if err != nil {
		return nil, fmt.Errorf("error loading config: %w", err)
	}

	app := &App{
		Config:       cfg,
		HttpPort:     cfg.AppPort,
		BuildVersion: cfg.BuildVersion,
		Logger:       logger.Logger,
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := app.initializeResources(ctx); err != nil {
		return nil, fmt.Errorf("error initializing resources: %w", err)
	}

	return app, nil
}

func (app *App) initializeResources(ctx context.Context) error {
	g, ctx := errgroup.WithContext(ctx)

	g.Go(func() error {
		var err error
		app.Cache, err = InitCache(ctx)
		return err
	})

	g.Go(func() error {
		var err error
		if strings.TrimSpace(app.Config.DBType) == "postgres" {
			app.DB, err = InitPqDatabase(app.Config)
		} else if app.Config.DBType == "mysql" {
			app.MDB, err = InitMySQLDatabase(app.Config)
		}
		return err
	})

	g.Go(func() error {
		app.FileUpload = filesystem.NewFileUploadService(app.Config)
		return nil
	})

	return g.Wait()
}

func InitPqDatabase(cfg *config.Config) (*pgxpool.Pool, error) {
	dbService, err := database.NewPgxDatabaseService(cfg)
	if err != nil {
		return nil, err
	}
	return dbService.GetDB(), nil
}

func InitMySQLDatabase(cfg *config.Config) (*sql.DB, error) {
	dbService, err := database.NewMySQLService(cfg)
	if err != nil {
		return nil, err
	}
	return dbService.GetDB(), nil
}

func InitCache(ctx context.Context) (cache.CacheService, error) {
	return cache.NewRedisCacheService(ctx)
}

func (app *App) CloseResources() error {
	var wg sync.WaitGroup
	errChan := make(chan error, 3)

	wg.Add(1)
	go func() {
		defer wg.Done()
		if app.Cache != nil {
			if err := app.Cache.Close(); err != nil {
				errChan <- fmt.Errorf("failed to close cache: %w", err)
			}
		}
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		if app.DB != nil {
			app.DB.Close()
		}
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		if app.MDB != nil {
			if err := app.MDB.Close(); err != nil {
				errChan <- fmt.Errorf("failed to close MySQL database: %w", err)
			}
		}
	}()

	wg.Wait()
	close(errChan)

	var errs []error
	for err := range errChan {
		errs = append(errs, err)
	}

	if len(errs) > 0 {
		return fmt.Errorf("errors closing resources: %v", errs)
	}

	return nil
}

func (app *App) SetupHTTPServer(handler http.Handler) {
	app.HttpServer = &http.Server{
		Addr:         fmt.Sprintf(":%d", app.Config.AppPort),
		Handler:      handler,
		ReadTimeout:  time.Duration(app.Config.ReadTimeout) * time.Second,
		WriteTimeout: time.Duration(app.Config.WriteTimeout) * time.Second,
		IdleTimeout:  time.Duration(app.Config.IdleTimeout) * time.Second,
		MaxHeaderBytes: app.Config.MaxHeaderBytes,
	}
}


func (app *App) CheckDatabaseHealth() error {
	if app.Config.DBType == "" {
		return fmt.Errorf("database type is not configured")
	}

	var err error
	switch app.Config.DBType {
	case "postgres":
		if app.DB == nil {
			return fmt.Errorf("PostgreSQL database connection is not initialized")
		}
		// Set a timeout for the health check
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		err = app.DB.Ping(ctx)
		if err != nil {
			return fmt.Errorf("PostgreSQL database is unhealthy: %w", err)
		}
	case "mysql":
		if app.MDB == nil {
			return fmt.Errorf("MySQL database connection is not initialized")
		}
		// Set a timeout for the health check
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		err = app.MDB.PingContext(ctx)
		if err != nil {
			return fmt.Errorf("MySQL database is unhealthy: %w", err)
		}
	default:
		return fmt.Errorf("unsupported database type: %s", app.Config.DBType)
	}

	app.Logger.Info("Database health check passed", zap.String("dbType", app.Config.DBType))
	return nil
}