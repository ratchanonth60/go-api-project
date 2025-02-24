package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"project-api/internal/controller"
	"project-api/internal/core/service"
	"project-api/internal/infra/aws"
	"project-api/internal/infra/config"
	"project-api/internal/infra/logger"
	"project-api/internal/infra/repository"

	"github.com/gofiber/storage/s3/v2"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

func main() {
	configType := flag.String("config", "yaml", "Configuration type (env or yaml)") // Default เป็น yaml
	flag.Parse()

	var err error
	if *configType == "env" {
		config.IsYaml = false
		err = config.LoadConfig("") // โหลดจาก YAML
	} else {
		err = config.LoadConfig("conf/app.yaml")
	}

	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}
	db := config.GormDB{
		Config: &gorm.Config{},
	}
	if db.Connect() != nil {
		log.Fatal("Failed to connect to database")
	}
	sqlDB, err := db.DB.DB()
	if err != nil {
		log.Fatalln(err)
	}
	defer sqlDB.Close()
	app, err := controller.New(Init(&db))
	if err != nil {
		log.Fatal("Failed to create server: ", err)
	}

	port := fmt.Sprintf(":%s", config.Config.Server.Port)
	logger.Info("Info server running at:", zap.Any("port", port))

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	go func() {
		if err := app.Serve(port); err != nil {
			logger.Fatal("Failed to start server:", zap.Error(err))
		}
	}()

	<-ctx.Done()
	logger.Info("Shutting down server...")
	defer sqlDB.Close()
	if err := app.Shutdown(); err != nil {
		logger.Fatal("Server shutdown failed:", zap.Error(err))
	}

	logger.Info("Server shut down gracefully")
}

func Init(db *config.GormDB) *controller.Services {
	s3Congfig := config.Config.GetS3Config()
	credential := config.Config.GetCredentials()
	awsConfig := s3.Config{
		Bucket:      s3Congfig.Bucket,
		Region:      s3Congfig.Region,
		Endpoint:    s3Congfig.Endpoint,
		Credentials: credential,
	}
	userRepository := repository.NewUserRepository(db.DB)
	userService := service.NewUserService(userRepository)
	aws3Repo := aws.New(awsConfig)
	fileService := service.NewS3Service(userService, aws3Repo)
	return &controller.Services{
		UserService: userService,
		FileService: fileService,
	}
}
