package controller

import (
	"net/http"

	"project-api/internal/core/model/request"
	"project-api/internal/core/model/response"

	In "project-api/internal/core/port/service"

	"github.com/gofiber/fiber/v2"
)

type UserHandler struct {
	service In.IUserService
}

func NewUserHandler(service In.IUserService) *UserHandler {
	return &UserHandler{
		service: service,
	}
}

func (u *UserHandler) CreateUser(ctx *fiber.Ctx) error {
	user := &request.UserRequest{}
	if err := ctx.BodyParser(user); err != nil {
		return ctx.Status(http.StatusBadRequest).JSON(response.ErrParser)
	}
	userEntity, err := user.ToEntity()
	if err != nil {
		return ctx.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	if err := u.service.Create(ctx.Context(), userEntity); err != nil {
		return ctx.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return ctx.Status(http.StatusCreated).JSON(user)
}

func (u *UserHandler) GetUserByEmail(ctx *fiber.Ctx) error {
	email := ctx.Params("email")
	user, err := u.service.GetUserByEmail(ctx.Context(), email)
	if err != nil {
		return ctx.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return ctx.Status(http.StatusOK).JSON(user)
}
