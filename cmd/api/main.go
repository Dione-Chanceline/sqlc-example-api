package main

import (
	"context"
	"errors"
	"fmt"
	"net/url"
	"os"

	"github.com/ardanlabs/conf/v3"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"

	"github.com/Iknite-Space/sqlc-example-api/api"
	"github.com/Iknite-Space/sqlc-example-api/db/repo"
)

// DBConfig holds the database configuration.
type DBConfig struct {
	DBUser      string `conf:"env:DB_USER,required"`
	DBPassword  string `conf:"env:DB_PASSWORD,required,mask"`
	DBHost      string `conf:"env:DB_HOST,required"`
	DBPort      uint16 `conf:"env:DB_PORT,required"`
	DBName      string `conf:"env:DB_Name,required"`
	TLSDisabled bool   `conf:"env:DB_TLS_DISABLED"`
}

// Config holds the application config.
type Config struct {
	ListenPort     uint16 `conf:"env:LISTEN_PORT,required"`
	MigrationsPath string `conf:"env:MIGRATIONS_PATH,required"`
	DB             DBConfig
}

func main() {
	if err := run(); err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}
}

func run() error {
	ctx := context.Background()
	config := Config{}

	// Load .env + config
	if err := LoadConfig(&config); err != nil {
		fmt.Println("Error loading config:", err)
		return err
	}

	// Database setup
	dbURL := getPostgresConnectionURL(config.DB)
	db, err := pgxpool.New(ctx, dbURL)
	if err != nil {
		return fmt.Errorf("failed to connect to DB: %w", err)
	}
	defer db.Close()

	// Run migrations
	if err := repo.Migrate(dbURL, config.MigrationsPath); err != nil {
		return fmt.Errorf("failed running migrations: %w", err)
	}

	// Initialize SQLC querier
	queries := repo.New(db)

	// Initialize Gin router
	router := gin.Default()

	// Register message routes
	messageHandler := api.NewMessageHandler(queries)
	messageHandler.RegisterRoutes(router)

	// Register attachment routes
	attachmentHandler := api.NewAttachmentHandler(queries)
	attachmentHandler.RegisterRoutes(router)

	// Start server
	return router.Run(fmt.Sprintf(":%d", config.ListenPort))
}

// LoadConfig reads configuration from env file.
func LoadConfig(cfg *Config) error {
	if _, err := os.Stat(".env"); err == nil {
		if err := godotenv.Load(); err != nil {
			return fmt.Errorf("failed to load .env: %w", err)
		}
	}

	_, err := conf.Parse("", cfg)
	if errors.Is(err, conf.ErrHelpWanted) {
		return err
	}
	return err
}

// Generate DB URL
func getPostgresConnectionURL(config DBConfig) string {
	values := url.Values{}
	if config.TLSDisabled {
		values.Add("sslmode", "disable")
	} else {
		values.Add("sslmode", "require")
	}

	return (&url.URL{
		Scheme:   "postgres",
		User:     url.UserPassword(config.DBUser, config.DBPassword),
		Host:     fmt.Sprintf("%s:%d", config.DBHost, config.DBPort),
		Path:     config.DBName,
		RawQuery: values.Encode(),
	}).String()
}
