package config

import (
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/spf13/viper"
)

type Config struct {
	BuildVersion      string        `mapstructure:"VERSION"`
	AppEnv            string        `mapstructure:"APP_ENV"`
	AppPort           int           `mapstructure:"APP_PORT"`
	Domain            string        `mapstructure:"DOMAIN"`
	DBType            string        `mapstructure:"DB_TYPE"`
	DBHost            string        `mapstructure:"DB_HOST"`
	DBPort            int           `mapstructure:"DB_PORT"`
	DBName            string        `mapstructure:"DB_NAME"`
	DBUser            string        `mapstructure:"DB_USER"`
	DBPassword        string        `mapstructure:"DB_PASSWORD"`
	DBSSLMode         string        `mapstructure:"DB_SSLMODE"`
	DBMaxIdleConns    int           `mapstructure:"DB_MAX_IDLE_CONNS"`
	DBMaxConnLifetime time.Duration `mapstructure:"DB_MAX_CONN_LIFETIME"`
	MaxConns          int           `mapstructure:"MAX_CONNS"`
	MinConns          int           `mapstructure:"MIN_CONNS"`
	Migrate           bool          `mapstructure:"MIGRATE"`
	Seed              bool          `mapstructure:"SEED"`
	RedisExp          int           `mapstructure:"REDIS_EXP"`
	RedisURI          string        `mapstructure:"REDIS_URI"`
	RedisPassword     string        `mapstructure:"REDIS_PASSWORD"`
	RedisDB           int           `mapstructure:"REDIS_DB"`
	IsRedis           bool          `mapstructure:"IS_REDIS"`
	RateLimitEnabled  bool          `mapstructure:"RATE_LIMIT_ENABLED"`
	RateLimit         int           `mapstructure:"RATE_LIMIT"`
	RateLimitDuration time.Duration `mapstructure:"RATE_LIMIT_DURATION"`
	JwtSecretKey      string        `mapstructure:"JWT_SECRET_KEY"`
	JwtExpiration     time.Duration `mapstructure:"JWT_EXPIRATION"`
	StorageDisk       string        `mapstructure:"STORAGE_DISK"`
	StoragePath       string        `mapstructure:"STORAGE_PATH"`
	AwsRegion         string        `mapstructure:"AWS_REGION"`
	AwsAccessKey      string        `mapstructure:"AWS_ACCESS_KEY"`
	AwsSecretKey      string        `mapstructure:"AWS_SECRET_KEY"`
	AwsBucket         string        `mapstructure:"AWS_BUCKET"`
	AwsEndpoint       string        `mapstructure:"AWS_ENDPOINT"`
	OtpExpiration     int           `mapstructure:"OTP_EXPIRATION"`
	OtpLength         int           `mapstructure:"OTP_LENGTH"`
	OtpResendDuration int           `mapstructure:"OTP_RESEND_DURATION"`
	ReadTimeout       int           `mapstructure:"READ_TIMEOUT"`
	WriteTimeout      int           `mapstructure:"WRITE_TIMEOUT"`
	IdleTimeout       int           `mapstructure:"IDLE_TIMEOUT"`
	MaxHeaderBytes    int           `mapstructure:"MAX_HEADER_BYTES"`
}

var (
	GlobalConfig *Config
	configMutex  sync.Mutex
)

// LoadConfig loads the configuration from the environment and config file
func LoadConfig() (*Config, error) {
	// Lock the mutex to ensure thread safety during configuration loading
	configMutex.Lock()
	defer configMutex.Unlock()

	// Initialize viper to read from .env file and environment variables
	viper.SetConfigFile(".env")
	viper.AutomaticEnv()

	// Read the configuration file
	if err := viper.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	// Unmarshal the configuration into a Config struct
	var cfg Config
	if err := viper.Unmarshal(&cfg); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config: %w", err)
	}

	// Validate configuration values
	if err := stringTrim(&cfg); err != nil {
		return nil, fmt.Errorf("failed to trim config values: %w", err)
	}

	// Set the global configuration variable
	GlobalConfig = &cfg

	return &cfg, nil
}

func stringTrim(cfg *Config) error {
	cfg.Domain = strings.TrimSpace(cfg.Domain)
	cfg.DBType = strings.TrimSpace(cfg.DBType)
	cfg.DBHost = strings.TrimSpace(cfg.DBHost)
	cfg.DBName = strings.TrimSpace(cfg.DBName)
	cfg.DBUser = strings.TrimSpace(cfg.DBUser)
	cfg.DBPassword = strings.TrimSpace(cfg.DBPassword)
	cfg.DBSSLMode = strings.TrimSpace(cfg.DBSSLMode)
	cfg.JwtSecretKey = strings.TrimSpace(cfg.JwtSecretKey)
	cfg.StorageDisk = strings.TrimSpace(cfg.StorageDisk)
	cfg.StoragePath = strings.TrimSpace(cfg.StoragePath)
	cfg.AwsRegion = strings.TrimSpace(cfg.AwsRegion)
	cfg.AwsAccessKey = strings.TrimSpace(cfg.AwsAccessKey)
	cfg.AwsSecretKey = strings.TrimSpace(cfg.AwsSecretKey)
	cfg.AwsBucket = strings.TrimSpace(cfg.AwsBucket)
	cfg.AwsEndpoint = strings.TrimSpace(cfg.AwsEndpoint)

	return nil
}

func GetConfig() *Config {
	return GlobalConfig
}
