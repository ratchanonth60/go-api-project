package controller

import (
	"context"
	"fmt"
	"project-api/internal/controller/handler"
	"project-api/internal/core/middleware"
	In "project-api/internal/core/port/service"
	"project-api/internal/infra/logger"
	"github.com/gofiber/template/html/v2"

	"github.com/RichardKnop/machinery/v2"
	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
)

// RouterConfig holds configuration for the router
type RouterConfig struct {
	JWTSecret      string
	ExcludedRoutes []string
}

// Services holds all required services
type Services struct {
	UserService In.IUserService
	FileService In.IS3Service
	Server      *machinery.Server
}

// Router encapsulates the Fiber app and its configuration
type Router struct {
	app *fiber.App
}

// New creates a new Router instance with optimized configuration
func New(services *Services) (*Router, error) {
	if services == nil || services.UserService == nil || services.FileService == nil {
		return nil, fmt.Errorf("services cannot be nil")
	}

	// Initialize Fiber with custom configuration
	app := fiber.New(fiber.Config{
		ErrorHandler: customErrorHandler,
		Views:        html.New("./templates", ".html"),
	})

	// Configure router
	router := &Router{app}
	router.setupRoutes(services)
	return router, nil
}

// setupRoutes configures all API routes
func (r *Router) setupRoutes(services *Services) {

	// health check
	r.app.Get("/", func(c *fiber.Ctx) error {
		return c.Status(fiber.StatusOK).JSON(fiber.Map{"message": "Welcome to the API"})

	})
	// Public routes (no authentication)
	auth := r.app.Group("/api/v1/auth")
	r.setupAuthRoutes(auth, services.UserService, services.Server)

	// Protected routes
	v1 := r.app.Group("/api/v1", middleware.JWTAuthMiddleware)
	r.setupProtectedRoutes(v1, services)
}

// setupAuthRoutes configures authentication routes
func (r *Router) setupAuthRoutes(group fiber.Router, userService In.IUserService, machineryServer *machinery.Server) {
	authHandler := controller.NewAuthHandler(userService, machineryServer)
	group.Post("/login", authHandler.LoginHandle)
	group.Post("/register", authHandler.RegisterHandler)
	group.Get("/confirm/:token", authHandler.ConfirmEmailHandler)
	group.Post("/resend", authHandler.ResendConfirmationEmailHandler)
	group.Post("/reset-password", authHandler.ResetPasswordHandler)
	group.Post("/reset-password/confirm/:token", authHandler.ConfirmResetPasswordHandler)
}

// setupProtectedRoutes configures authenticated routes
func (r *Router) setupProtectedRoutes(group fiber.Router, services *Services) {
	// User routes
	userGroup := group.Group("/users")
	userHandler := controller.NewUserHandler(services.UserService)
	userGroup.Post("/", userHandler.CreateUser)
	userGroup.Get("/:email", userHandler.GetUserByEmail)

	// File routes
	fileGroup := group.Group("/files")
	fileHandler := controller.NewFileHandler(services.UserService, services.FileService)
	fileGroup.Post("/upload", fileHandler.UploadFile)
	fileGroup.Delete("/delete/:key", fileHandler.DeleteFile)
	fileGroup.Get("/download/:key", fileHandler.DownloadFile)
	fileGroup.Use(func(c *fiber.Ctx) error {
		logger.Warn("Unhandled file route", zap.String("path", c.Path()))
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"code":  404,
			"error": fmt.Sprintf("Cannot %s %s", c.Method(), c.Path()),
		})
	})
}

// customErrorHandler handles Fiber errors
func customErrorHandler(c *fiber.Ctx, err error) error {
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
			"code":  fiber.StatusInternalServerError,
		})
	}
	return nil
}

// jwtErrorHandler handles JWT-specific errors
func jwtErrorHandler(c *fiber.Ctx, err error) error {
	return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
		"error": err.Error(),
		"code":  fiber.StatusUnauthorized,
	})
}

func (r *Router) ShutdownWithContext(ctx context.Context) error {
	return r.app.ShutdownWithContext(ctx)
}
func (r *Router) Serve(port string) error {
	return r.app.Listen(port)
}
