package handler

import (
	"log/slog"
	"net/http"

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

type RegisterRequest struct {
	FullName string `json:"full_name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (h *FiberHandler) Register(c *fiber.Ctx) error {
	var req RegisterRequest
	if err := c.BodyParser(&req); err!= nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "Cannot parse JSON"})
	}

	user, err := h.uc.Register(c.Context(), req.FullName, req.Email, req.Password)
	if err!= nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "Could not register user"})
	}

	return c.Status(http.StatusCreated).JSON(fiber.Map{"id": user.ID, "email": user.Email})
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (h *FiberHandler) Login(c *fiber.Ctx) error {
	var req LoginRequest
	if err := c.BodyParser(&req); err!= nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "Cannot parse JSON"})
	}

	accessToken, refreshToken, err := h.uc.Login(c.Context(), req.Email, req.Password)
	if err!= nil {
		return c.Status(http.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid credentials"})
	}

	return c.JSON(fiber.Map{
		"access_token":  accessToken,
		"refresh_token": refreshToken,
	})
}