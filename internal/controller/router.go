package controller

import (
	handler "project-api/internal/controller/handler"
	"project-api/internal/core/middleware"
	In "project-api/internal/core/port/service"
	"project-api/internal/infra/config"

	jwtware "github.com/gofiber/contrib/jwt"
	"github.com/gofiber/fiber/v2"
)

type Router struct {
	*fiber.App
}
type Services struct {
	UserService In.IUserService
	FileService In.IS3Service
}

func New(services *Services) (*Router, error) {
	app := fiber.New()
	auth := app.Group("api/v1/auth")
	{
		auth.Post("/login", handler.NewAuthHandler(services.UserService).LoginHandle)
		auth.Post("/register", handler.NewAuthHandler(services.UserService).RegisterHandler)

	}
	app.Use(jwtware.New(jwtware.Config{
		SigningKey: jwtware.SigningKey{Key: []byte(config.Config.JWT.Signed)},
		Filter: func(c *fiber.Ctx) bool {
			excludedRoutes := []string{
				"/api/v1/auth/login",
				"/api/v1/auth/register",
				"/api/v1/auth/logout",
			}

			for _, route := range excludedRoutes {
				if c.Path() == route {
					return true
				}
			}
			return false
		},
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": err.Error(),
				"code":  fiber.StatusUnauthorized,
			})
		},
	}))

	v1 := app.Group("/api/v1", middleware.JWTAuthMiddleware)
	{
		user := v1.Group("users")
		{
			user.Post("/", handler.NewUserHandler(services.UserService).CreateUser)
			user.Get("/:email", handler.NewUserHandler(services.UserService).GetUserByEmail)
		}
		file := v1.Group("files")
		{
			file.Post("/upload", handler.NewFileHandler(services.UserService, services.FileService).UploadFile)
		}
	}
	return &Router{
		app,
	}, nil
}

func (r *Router) Serve(port string) error {
	return r.Listen(port)
}
