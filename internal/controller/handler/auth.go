package controller

import (
	"net/http"

	"project-api/internal/core/common/utils"
	"project-api/internal/core/model/request"
	"project-api/internal/core/model/response"
	In "project-api/internal/core/port/service"

	"github.com/gofiber/fiber/v2"
	"golang.org/x/crypto/bcrypt"
)

type AuthHandler struct {
	service In.IUserService
}

func NewAuthHandler(service In.IUserService) *AuthHandler {
	return &AuthHandler{
		service: service,
	}
}

func (l *AuthHandler) LoginHandle(c *fiber.Ctx) error {
	var req request.LoginRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(http.StatusBadRequest).JSON(response.ErrParser)
	}
	user, err := l.service.GetUserByUserName(c.Context(), req.UserName)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(response.ErrAuth)
	}
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		return c.Status(http.StatusBadRequest).JSON(response.ErrorResponse{
			Code: http.StatusBadRequest,
			Msg:  "Password can't hash",
			Data: err.Error(),
		})
	}
	token, err := utils.GenerateJWT(user.Email)
	if err != nil {
		return err
	}
	return c.Status(http.StatusOK).JSON(
		response.SuccResponse{
			Code: http.StatusOK,
			Msg:  "successfully logged in",
			Data: token,
		})
}

func (l *AuthHandler) RegisterHandler(c *fiber.Ctx) error {
	var req request.RegisterRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(http.StatusBadRequest).JSON(response.ErrorResponse{
			Code: http.StatusBadRequest,
			Msg:  "bad request, please check the request body",
			Data: err.Error(),
		})
	}
	if !req.ConfirmPassword() {
		return c.Status(http.StatusBadRequest).JSON(response.ErrCofirmPassword)
	}
	hashed, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	user := request.UserRequest{
		FirstName: req.FirstName,
		LastName:  req.LastName,
		Username:  req.UserName,
		Email:     req.Email,
		Password:  string(hashed),
	}
	if err := l.service.Create(c.Context(), &user); err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	toEntity, err := user.ToEntity()
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	jsonData, err := toEntity.ToJson()
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(response.ErrorResponse{
			Code: fiber.ErrBadGateway.Code,
			Msg:  "Can't convert user to json",
			Data: err,
		})
	}
	return c.Status(http.StatusCreated).JSON(response.SuccResponse{
		Code: http.StatusCreated,
		Msg:  "successfully created user",
		Data: jsonData,
	})
}
