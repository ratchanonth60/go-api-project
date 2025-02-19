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
	"project-api/internal/core/common/utils"
	"project-api/internal/core/service"
	"project-api/internal/infra/config"
	"project-api/internal/infra/repository"

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
	utils.Logger.Printf("Server is running on port %s", port)

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	go func() {
		if err := app.Serve(port); err != nil {
			utils.Logger.Fatalf("Failed to start server: %v", err)
		}
	}()

	<-ctx.Done()
	utils.Logger.Println("Shutting down server...")
	defer sqlDB.Close()
	if err := app.Shutdown(); err != nil {
		utils.Logger.Fatalf("Server shutdown failed: %v", err)
	}

	utils.Logger.Println("Server gracefully stopped")
}

func Init(db *config.GormDB) *controller.Services {
	userRepository := repository.NewUserRepository(db.DB)
	userService := service.NewUserService(userRepository)
	return &controller.Services{
		UserService: userService,
	}
}
