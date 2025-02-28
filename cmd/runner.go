package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"project-api/internal/controller"
	"project-api/internal/core/service"
	"project-api/internal/infra/aws"
	"project-api/internal/infra/config"
	"project-api/internal/infra/logger"
	"project-api/internal/infra/repository"
	"project-api/internal/task"

	"github.com/RichardKnop/machinery/v2"
	"github.com/gofiber/storage/s3/v2"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

// Application holds the main application components
type Application struct {
	config   *config.AppConfig
	db       *config.GormDB
	router   *controller.Router
	services *controller.Services
}

// Config holds runtime configuration
type Config struct {
	configType string
	port       string
}

func main() {
	if err := run(); err != nil {
		log.Fatalf("Application failed: %v", err)
	}
}

func run() error {
	// Parse flags and load configuration first
	configType := parseFlags()
	if err := loadConfig(configType); err != nil {
		return fmt.Errorf("failed to load configuration: %w", err)
	}

	// Set port after config is loaded
	cfg := Config{
		configType: configType,
		port:       config.Config.Server.Port, // Use the port directly without adding extra colon
	}

	// Initialize application
	app, err := initializeApp()
	if err != nil {
		return fmt.Errorf("failed to initialize application: %w", err)
	}

	// Start server with graceful shutdown
	return runServer(app, cfg.port)
}

func parseFlags() string {
	configType := flag.String("config", "yaml", "Configuration type (env or yaml)")
	flag.Parse()
	return *configType
}

func loadConfig(configType string) error {
	if configType == "env" {
		config.IsYaml = false
		return config.LoadConfig("")
	}
	return config.LoadConfig("conf/app.yaml")
}

func initializeApp() (*Application, error) {
	// Initialize database
	db := &config.GormDB{Config: &gorm.Config{}}
	if err := db.Connect(); err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	// Get SQL DB instance for cleanup
	sqlDB, err := db.DB.DB()
	if err != nil {
		return nil, fmt.Errorf("failed to get SQL DB: %w", err)
	}

	machineryServer, err := task.NewMachineryServer()
	if err != nil {
		return nil, fmt.Errorf("failed to initialize Machinery server: %w", err)
	}
	// redisClient := redis.NewRedisClient()
	// Initialize services
	services := initializeServices(db, machineryServer)

	// Create router
	router, err := controller.New(services)
	if err != nil {
		sqlDB.Close()
		return nil, fmt.Errorf("failed to create router: %w", err)
	}

	return &Application{
		config:   config.Config,
		db:       db,
		router:   router,
		services: services,
	}, nil
}

func initializeServices(db *config.GormDB, machineryServer *machinery.Server) *controller.Services {
	s3Config := config.Config.GetS3Config()
	credential := config.Config.GetCredentials()

	awsConfig := s3.Config{
		Bucket:      s3Config.Bucket,
		Region:      s3Config.Region,
		Endpoint:    s3Config.Endpoint,
		Credentials: credential,
	}

	fileRepo := repository.NewFileRepository(db.DB)
	userRepo := repository.NewUserRepository(db.DB)
	userService := service.NewUserService(userRepo)
	s3Repo := aws.New(awsConfig)
	fileService := service.NewS3Service(fileRepo, s3Repo)

	return &controller.Services{
		UserService: userService,
		FileService: fileService,
		Server:      machineryServer,
	}
}

func runServer(app *Application, port string) error {
	// Ensure port has a leading colon if it doesn't already
	logger.Info("Parsed port value", zap.String("port", port))
	if port != "" && port[0] != ':' {
		port = ":" + port
	}

	// Setup signal handling
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	// Start server in goroutine
	errChan := make(chan error, 1)
	go func() {
		logger.Info("Server starting", zap.String("port", port))
		if err := app.router.Serve(port); err != nil {
			errChan <- fmt.Errorf("server failed: %w", err)
		}
	}()

	// Wait for shutdown signal or server error
	select {
	case <-ctx.Done():
		logger.Info("Shutdown signal received")
	case err := <-errChan:
		return err
	}

	// Graceful shutdown
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := app.router.ShutdownWithContext(shutdownCtx); err != nil {
		logger.Error("Failed to shutdown server gracefully", zap.Error(err))
	}
	// Close database connection
	if sqlDB, err := app.db.DB.DB(); err == nil {
		if err := sqlDB.Close(); err != nil {
			logger.Error("Failed to close database", zap.Error(err))
		}
	}

	logger.Info("Server shut down gracefully")
	return nil
}
