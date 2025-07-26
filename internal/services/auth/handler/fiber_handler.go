package handler

import (
	"log/slog"
	"net/http"

	"github.com/dpalhz/microservice-exp-with-go/internal/pkg/apiresponse"
	"github.com/dpalhz/microservice-exp-with-go/internal/services/auth/dto"
	"github.com/dpalhz/microservice-exp-with-go/internal/services/auth/usecase"

	"github.com/gofiber/fiber/v2"
)

type FiberHandler struct {
	uc  *usecase.UserUsecase
	log *slog.Logger
}

func NewFiberHandler(uc *usecase.UserUsecase, log *slog.Logger) *FiberHandler {
	return &FiberHandler{uc: uc, log: log}
}

func (h *FiberHandler) RegisterRoutes(app *fiber.App) {
	api := app.Group("/api/v1/auth")
	api.Post("/register", h.Register)
	api.Post("/login", h.Login)
}

func (h *FiberHandler) Register(c *fiber.Ctx) error {
	var req dto.RegisterRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(http.StatusBadRequest).JSON(apiresponse.Error("cannot parse json"))
	}

	user, err := h.uc.Register(c.Context(), req.FullName, req.Email, req.Password)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(apiresponse.Error("could not register user"))
	}

	resp := dto.RegisterResponse{ID: user.ID, Email: user.Email}
	return c.Status(http.StatusCreated).JSON(apiresponse.Success(resp))
}

func (h *FiberHandler) Login(c *fiber.Ctx) error {
	var req dto.LoginRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(http.StatusBadRequest).JSON(apiresponse.Error("cannot parse json"))
	}

	accessToken, refreshToken, err := h.uc.Login(c.Context(), req.Email, req.Password)
	if err != nil {
		return c.Status(http.StatusUnauthorized).JSON(apiresponse.Error("invalid credentials"))
	}

	resp := dto.LoginResponse{AccessToken: accessToken, RefreshToken: refreshToken}
	return c.JSON(apiresponse.Success(resp))
}
