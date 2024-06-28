package config

import (
	"fmt"
	"sync"

	"github.com/spf13/viper"
)

type Config struct {
	BuildVersion      string `mapstructure:"VERSION"`
	AppEnv            string `mapstructure:"APP_ENV"`
	AppPort           int    `mapstructure:"APP_PORT"`
	Domain            string `mapstructure:"DOMAIN"`
	DBType            string `mapstructure:"DB_TYPE"`
	DBHost            string `mapstructure:"DB_HOST"`
	DBPort            int    `mapstructure:"DB_PORT"`
	DBName            string `mapstructure:"DB_NAME"`
	DBUser            string `mapstructure:"DB_USER"`
	DBPassword        string `mapstructure:"DB_PASSWORD"`
	DBSSLMode         string `mapstructure:"DB_SSLMODE"`
	DBMaxIdleConns    int    `mapstructure:"DB_MAX_IDLE_CONNS"`
	DBMaxConnLifetime int    `mapstructure:"DB_MAX_CONN_LIFETIME"`
	MaxConns          int    `mapstructure:"MAX_CONNS"`
	MinConns          int    `mapstructure:"MIN_CONNS"`
	Migrate           bool   `mapstructure:"MIGRATE"`
	Seed              bool   `mapstructure:"SEED"`
	RedisExp          int    `mapstructure:"REDIS_EXP"`
	RedisURI          string `mapstructure:"REDIS_URI"`
	RedisPassword     string `mapstructure:"REDIS_PASSWORD"`
	RedisDB           int    `mapstructure:"REDIS_DB"`
	IsRedis           bool   `mapstructure:"IS_REDIS"`
	RateLimitEnabled  bool   `mapstructure:"RATE_LIMIT_ENABLED"`
	RateLimit         int    `mapstructure:"RATE_LIMIT"`
	RateLimitDuration string `mapstructure:"RATE_LIMIT_DURATION"`
	JwtSecretKey      string `mapstructure:"JWT_SECRET_KEY"`
	JwtExpiration     string `mapstructure:"JWT_EXPIRATION"`
	StorageDisk       string `mapstructure:"STORAGE_DISK"`
	StoragePath       string `mapstructure:"STORAGE_PATH"`
	AwsRegion         string `mapstructure:"AWS_REGION"`
	AwsAccessKey      string `mapstructure:"AWS_ACCESS_KEY"`
	AwsSecretKey      string `mapstructure:"AWS_SECRET_KEY"`
	AwsBucket         string `mapstructure:"AWS_BUCKET"`
	AwsEndpoint       string `mapstructure:"AWS_ENDPOINT"`
}

var (
	GlobalConfig *Config
	configMutex  sync.Mutex
)

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

	// Set default values for configuration fields
	setDefaultValues(&cfg)

	// Set the global configuration variable
	GlobalConfig = &cfg

	return &cfg, nil
}

// setDefaultValues sets default values for configuration fields
func setDefaultValues(cfg *Config) {
	if cfg.RedisDB == 0 {
		cfg.RedisDB = 0 // Default Redis database
	}
	// Add default values for other configuration fields as needed
}
