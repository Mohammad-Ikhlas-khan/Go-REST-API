package handler

import (
	"errors"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/go-playground/validator/v10"
	"go.uber.org/zap"

	"github.com/example/go-user-api/internal/logger"
	"github.com/example/go-user-api/internal/models"
	"github.com/example/go-user-api/internal/service"
)

// UserHandler holds all HTTP handlers for the /users resource.
type UserHandler struct {
	svc      service.UserService
	validate *validator.Validate
	log      *zap.Logger
}

// NewUserHandler creates a new UserHandler.
func NewUserHandler(svc service.UserService) *UserHandler {
	return &UserHandler{
		svc:      svc,
		validate: validator.New(),
		log:      logger.Get(),
	}
}

// CreateUser handles POST /users
func (h *UserHandler) CreateUser(c *fiber.Ctx) error {
	var req models.CreateUserRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid request body",
		})
	}

	if err := h.validate.Struct(req); err != nil {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{
			"error":  "validation failed",
			"detail": err.Error(),
		})
	}

	resp, err := h.svc.CreateUser(c.Context(), req)
	if err != nil {
		h.log.Error("CreateUser handler error", zap.Error(err))
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "failed to create user",
		})
	}

	return c.Status(fiber.StatusCreated).JSON(resp)
}

// GetUser handles GET /users/:id
func (h *UserHandler) GetUser(c *fiber.Ctx) error {
	id, err := parseID(c)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid user id"})
	}

	resp, err := h.svc.GetUser(c.Context(), id)
	if err != nil {
		if errors.Is(err, service.ErrNotFound) {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "user not found"})
		}
		h.log.Error("GetUser handler error", zap.Error(err))
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "internal server error"})
	}

	return c.JSON(resp)
}

// UpdateUser handles PUT /users/:id
func (h *UserHandler) UpdateUser(c *fiber.Ctx) error {
	id, err := parseID(c)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid user id"})
	}

	var req models.UpdateUserRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid request body"})
	}

	if err := h.validate.Struct(req); err != nil {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{
			"error":  "validation failed",
			"detail": err.Error(),
		})
	}

	resp, err := h.svc.UpdateUser(c.Context(), id, req)
	if err != nil {
		if errors.Is(err, service.ErrNotFound) {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "user not found"})
		}
		h.log.Error("UpdateUser handler error", zap.Error(err))
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "internal server error"})
	}

	return c.JSON(resp)
}

// DeleteUser handles DELETE /users/:id
func (h *UserHandler) DeleteUser(c *fiber.Ctx) error {
	id, err := parseID(c)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid user id"})
	}

	if err := h.svc.DeleteUser(c.Context(), id); err != nil {
		if errors.Is(err, service.ErrNotFound) {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "user not found"})
		}
		h.log.Error("DeleteUser handler error", zap.Error(err))
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "internal server error"})
	}

	return c.SendStatus(fiber.StatusNoContent)
}

// ListUsers handles GET /users
func (h *UserHandler) ListUsers(c *fiber.Ctx) error {
	page, _ := strconv.Atoi(c.Query("page", "1"))
	limit, _ := strconv.Atoi(c.Query("limit", "10"))

	resp, err := h.svc.ListUsers(c.Context(), page, limit)
	if err != nil {
		h.log.Error("ListUsers handler error", zap.Error(err))
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "internal server error"})
	}

	return c.JSON(resp)
}

func parseID(c *fiber.Ctx) (int32, error) {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil || id <= 0 {
		return 0, errors.New("invalid id")
	}
	return int32(id), nil
}
