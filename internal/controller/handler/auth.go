package controller

import (
	"net/http"

	"project-api/internal/core/common/utils"
	"project-api/internal/core/model/request"
	"project-api/internal/core/model/response"
	In "project-api/internal/core/port/service"
	"project-api/internal/infra/logger"

	"github.com/RichardKnop/machinery/v2"
	"github.com/RichardKnop/machinery/v2/tasks"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
)

type AuthHandler struct {
	service In.IUserService
	server  *machinery.Server
}

func NewAuthHandler(service In.IUserService, machineryServer *machinery.Server) *AuthHandler {
	return &AuthHandler{
		service: service,
		server:  machineryServer,
	}
}

func (l *AuthHandler) LoginHandle(c *fiber.Ctx) error {
	var req request.LoginRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(http.StatusBadRequest).JSON(response.ErrParser)
	}
	if err := req.Validate(); err != nil {
		return c.Status(http.StatusBadRequest).JSON(response.ErrorResponse{
			Code: http.StatusBadRequest,
			Msg:  "Bad request, please check the request body",
			Data: err.Error(),
		})
	}
	user, err := l.service.GetUserByName(c.Context(), req.UserName)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(response.ErrAuth)
	}
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		return c.Status(http.StatusUnauthorized).JSON(response.ErrorResponse{
			Code: http.StatusUnauthorized,
			Msg:  "Password or username is incorrect",
			Data: err.Error(),
		})
	}
	token, err := utils.GenerateJWT(user)
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
	if err := req.Validate(); err != nil {
		return c.Status(http.StatusBadRequest).JSON(response.ErrorResponse{
			Code: http.StatusBadRequest,
			Msg:  "Bad request, please check the request body",
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
	token := uuid.New().String()
	user := request.UserRequest{
		FirstName:    req.FirstName,
		LastName:     req.LastName,
		Username:     req.UserName,
		Email:        req.Email,
		Password:     string(hashed),
		IsActive:     false,
		ConfirmToken: token,
	}
	userEntity, err := user.ToEntity()
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	if err := l.service.Create(c.Context(), userEntity); err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	host := "http://localhost:8000"
	signature := &tasks.Signature{
		Name: "send_confirmation_email",
		Args: []tasks.Arg{
			{Type: "string", Value: user.Email},
			{Type: "string", Value: token},
			{Type: "string", Value: user.FirstName},
			{Type: "string", Value: host},
		},
	}
	_, err = l.server.SendTask(signature)
	if err != nil {
		logger.Error("Failed to queue confirmation email task", zap.Error(err))
	} else {
		logger.Info("Successfully queued confirmation email task", zap.String("email", user.Email))
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

func (h *AuthHandler) ConfirmEmailHandler(c *fiber.Ctx) error {
	token := c.Params("token")
	if token == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Token is required"})
	}

	if err := h.service.ConfirmEmail(c.UserContext(), token); err != nil {
		logger.Error("Email confirmation failed", zap.Error(err))
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{"message": "Email confirmed successfully"})
}
