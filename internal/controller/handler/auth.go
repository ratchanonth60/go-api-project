package controller

import (
	"fmt"
	"net/http"

	"project-api/internal/core/common/utils"
	"project-api/internal/core/model/request"
	"project-api/internal/core/model/response"
	In "project-api/internal/core/port/service"
	"project-api/internal/infra/config"
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
		return c.Status(fiber.StatusOK).JSON(response.ErrParser)
	}
	if err := req.Validate(); err != nil {
		return c.Status(fiber.StatusOK).JSON(response.ErrorResponse{
			Code: http.StatusBadRequest,
			Msg:  "Bad request, please check the request body",
			Data: err.Error(),
		})
	}
	user, err := l.service.GetUserByName(c.Context(), req.UserName)
	if err != nil {
		return c.Status(fiber.StatusOK).JSON(response.ErrAuth)
	}
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		return c.Status(fiber.StatusOK).JSON(response.ErrorResponse{
			Code: http.StatusUnauthorized,
			Msg:  "Password or username is incorrect",
			Data: err.Error(),
		})
	}
	token, err := utils.GenerateJWT(user)
	if err != nil {
		return err
	}
	return c.Status(fiber.StatusOK).JSON(
		response.SuccResponse{
			Msg:  "successfully logged in",
			Data: token,
		})
}

func (l *AuthHandler) RegisterHandler(c *fiber.Ctx) error {
	var req request.RegisterRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusOK).JSON(response.ErrorResponse{
			Code: http.StatusBadRequest,
			Msg:  "bad request, please check the request body",
			Data: err.Error(),
		})
	}
	if err := req.Validate(); err != nil {
		return c.Status(fiber.StatusOK).JSON(response.ErrorResponse{
			Code: http.StatusBadRequest,
			Msg:  "Bad request, please check the request body",
			Data: err.Error(),
		})
	}
	if !req.ConfirmPassword() {
		return c.Status(fiber.StatusOK).JSON(response.ErrCofirmPassword)
	}
	hashed, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return c.Status(fiber.StatusOK).JSON(response.ErrorResponse{
			Code: fiber.StatusUnauthorized,
			Msg:  "Error hashing password",
		})
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
		return c.Status(fiber.StatusOK).JSON(response.ErrorResponse{
			Code: fiber.StatusUnauthorized,
			Msg:  "Error to convert to entity",
		})
	}
	if err := l.service.Create(c.Context(), userEntity); err != nil {
		return c.Status(fiber.StatusOK).JSON(response.ErrorResponse{
			Code: fiber.StatusInternalServerError,
			Msg:  "Error to create user",
		})
	}
	host := fmt.Sprintf("http://%s:%s", config.Config.Server.Host, config.Config.Server.Port)
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
		return c.Status(fiber.StatusOK).JSON(fiber.Map{"error": err.Error()})
	}
	jsonData, err := toEntity.ToJson()
	if err != nil {
		return c.Status(fiber.StatusOK).JSON(response.ErrorResponse{
			Code: fiber.ErrBadGateway.Code,
			Msg:  "Can't convert user to json",
			Data: err,
		})
	}
	return c.Status(fiber.StatusOK).JSON(response.SuccResponse{
		Msg:  "successfully created user",
		Data: jsonData,
	})
}

func (h *AuthHandler) ConfirmEmailHandler(c *fiber.Ctx) error {
	token := c.Params("token")
	if token == "" {
		return c.Status(fiber.StatusOK).JSON(response.ErrorResponse{
			Code: http.StatusBadRequest,
			Msg:  "Missing token",
		})
	}

	if err := h.service.ConfirmEmail(c.UserContext(), token); err != nil {
		logger.Error("Email confirmation failed", zap.Error(err))
		return c.Status(fiber.StatusOK).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(fiber.StatusOK).JSON(response.SuccResponse{
		Msg: "Email confirmed successfully",
	})
}

func (h *AuthHandler) ResendConfirmationEmailHandler(c *fiber.Ctx) error {
	var req request.EmailRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusOK).JSON(response.ErrorResponse{
			Code: http.StatusBadRequest,
			Msg:  "Bad request, please check the request body",
			Data: err.Error(),
		})
	}

	// ตรวจสอบว่า email ถูกต้อง
	if req.Email == "" {
		return c.Status(fiber.StatusOK).JSON(response.ErrorResponse{
			Code: http.StatusBadRequest,
			Msg:  "Email is required",
		})
	}

	// ส่งคำขอไปยัง service เพื่อสร้าง token ใหม่และอัปเดต user
	user, err := h.service.ResendConfirmationEmail(c.UserContext(), req.Email)
	if err != nil {
		logger.Error("Failed to resend confirmation email", zap.String("email", req.Email), zap.Error(err))
		return c.Status(fiber.StatusOK).JSON(response.ErrorResponse{
			Code: http.StatusInternalServerError,
			Msg:  "Failed to resend confirmation email",
			Data: err.Error(),
		})
	}

	// ส่ง email confirmation ผ่าน Machinery
	host := fmt.Sprintf("http://%s:%s", config.Config.Server.Host, config.Config.Server.Port)
	signature := &tasks.Signature{
		Name: "send_confirmation_email",
		Args: []tasks.Arg{
			{Type: "string", Value: user.Email},
			{Type: "string", Value: user.ConfirmToken},
			{Type: "string", Value: user.FirstName},
			{Type: "string", Value: host},
		},
	}
	_, err = h.server.SendTask(signature)
	if err != nil {
		logger.Error("Failed to queue resend confirmation email task", zap.String("email", user.Email), zap.Error(err))
	} else {
		logger.Info("Successfully queued resend confirmation email task", zap.String("email", user.Email))
	}

	return c.Status(fiber.StatusOK).JSON(response.SuccResponse{
		Msg: "Email confirmation re-sent successfully",
	})
}

func (h *AuthHandler) ResetPasswordHandler(c *fiber.Ctx) error {
	req := request.EmailRequest{}
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusOK).JSON(response.ErrorResponse{
			Code: http.StatusBadRequest,
			Msg:  "Bad request, please check the request body",
			Data: err.Error(),
		})
	}

	// ตรวจสอบว่า email ถูกต้อง
	if req.Email == "" {
		return c.Status(fiber.StatusOK).JSON(response.ErrorResponse{
			Code: http.StatusBadRequest,
			Msg:  "Email is required",
		})
	}

	// ส่งคำขอไปยัง service เพื่อสร้าง reset password token
	user, err := h.service.ResetPassword(c.UserContext(), req.Email)
	if err != nil {
		logger.Error("Failed to request reset password", zap.String("email", req.Email), zap.Error(err))
		return c.Status(fiber.StatusOK).JSON(response.ErrorResponse{
			Code: http.StatusInternalServerError,
			Msg:  "Failed to request reset password",
			Data: err.Error(),
		})
	}

	// ส่ง email reset password ผ่าน Machinery
	host := fmt.Sprintf("http://%s:%s", config.Config.Server.Host, config.Config.Server.Port)
	signature := &tasks.Signature{
		Name: "send_reset_password_email",
		Args: []tasks.Arg{
			{Type: "string", Value: user.Email},
			{Type: "string", Value: user.ResetPasswordToken},
			{Type: "string", Value: user.FirstName},
			{Type: "string", Value: host},
		},
	}
	_, err = h.server.SendTask(signature)
	if err != nil {
		logger.Error("Failed to queue reset password email task", zap.String("email", user.Email), zap.Error(err))
	} else {
		logger.Info("Successfully queued reset password email task", zap.String("email", user.Email))
	}

	return c.Status(fiber.StatusOK).JSON(response.SuccResponse{
		Msg: "Reset password email sent successfully",
	})
}

func (h *AuthHandler) ConfirmResetPasswordHandler(c *fiber.Ctx) error {
	var req request.ConfirmResetPassword
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusOK).JSON(response.ErrorResponse{
			Code: http.StatusBadRequest,
			Msg:  "Bad request, please check the request body",
			Data: err.Error(),
		})
	}

	// Validate struct
	if err := req.IsValid(); err {
		logger.Warn("Validation failed for reset password")
		return c.Status(fiber.StatusOK).JSON(response.ErrorResponse{
			Code: http.StatusBadRequest,
			Msg:  "Validation failed",
		})
	}

	// เรียก service เพื่อยืนยันรหัสผ่านใหม่
	if err := h.service.ConfirmResetPassword(c.UserContext(), req.Token, req.NewPassword); err != nil {
		logger.Error("Failed to confirm reset password", zap.String("token", req.Token), zap.Error(err))
		return c.Status(fiber.StatusOK).JSON(response.ErrorResponse{
			Code: http.StatusInternalServerError,
			Msg:  "Failed to reset password",
			Data: err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(response.SuccResponse{
		Msg: "Password reset successfully",
	})
}
