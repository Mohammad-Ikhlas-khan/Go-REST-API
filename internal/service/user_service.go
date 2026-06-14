package service

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	db "github.com/example/go-user-api/db/sqlc"
	"github.com/example/go-user-api/internal/models"
	"github.com/example/go-user-api/internal/repository"
	"github.com/example/go-user-api/internal/logger"
	"go.uber.org/zap"
)

const dobLayout = "2006-01-02"

// UserService holds business logic for user operations.
type UserService interface {
	CreateUser(ctx context.Context, req models.CreateUserRequest) (models.UserResponse, error)
	GetUser(ctx context.Context, id int32) (models.UserWithAgeResponse, error)
	UpdateUser(ctx context.Context, id int32, req models.UpdateUserRequest) (models.UserResponse, error)
	DeleteUser(ctx context.Context, id int32) error
	ListUsers(ctx context.Context, page, limit int) (models.PaginatedUsersResponse, error)
}

type userService struct {
	repo repository.UserRepository
	log  *zap.Logger
}

// NewUserService creates a new UserService.
func NewUserService(repo repository.UserRepository) UserService {
	return &userService{
		repo: repo,
		log:  logger.Get(),
	}
}

func (s *userService) CreateUser(ctx context.Context, req models.CreateUserRequest) (models.UserResponse, error) {
	dob, err := time.Parse(dobLayout, req.DOB)
	if err != nil {
		return models.UserResponse{}, fmt.Errorf("invalid dob format: %w", err)
	}

	user, err := s.repo.Create(ctx, db.CreateUserParams{
		Name: req.Name,
		Dob:  dob,
	})
	if err != nil {
		s.log.Error("failed to create user", zap.Error(err))
		return models.UserResponse{}, err
	}

	s.log.Info("user created", zap.Int32("id", user.ID), zap.String("name", user.Name))
	return toUserResponse(user), nil
}

func (s *userService) GetUser(ctx context.Context, id int32) (models.UserWithAgeResponse, error) {
	user, err := s.repo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return models.UserWithAgeResponse{}, ErrNotFound
		}
		s.log.Error("failed to get user", zap.Int32("id", id), zap.Error(err))
		return models.UserWithAgeResponse{}, err
	}

	resp := toUserWithAgeResponse(user)
	s.log.Info("user fetched", zap.Int32("id", id))
	return resp, nil
}

func (s *userService) UpdateUser(ctx context.Context, id int32, req models.UpdateUserRequest) (models.UserResponse, error) {
	// Ensure the user exists first
	_, err := s.repo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return models.UserResponse{}, ErrNotFound
		}
		return models.UserResponse{}, err
	}

	dob, err := time.Parse(dobLayout, req.DOB)
	if err != nil {
		return models.UserResponse{}, fmt.Errorf("invalid dob format: %w", err)
	}

	user, err := s.repo.Update(ctx, db.UpdateUserParams{
		ID:   id,
		Name: req.Name,
		Dob:  dob,
	})
	if err != nil {
		s.log.Error("failed to update user", zap.Int32("id", id), zap.Error(err))
		return models.UserResponse{}, err
	}

	s.log.Info("user updated", zap.Int32("id", id))
	return toUserResponse(user), nil
}

func (s *userService) DeleteUser(ctx context.Context, id int32) error {
	_, err := s.repo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return ErrNotFound
		}
		return err
	}

	if err := s.repo.Delete(ctx, id); err != nil {
		s.log.Error("failed to delete user", zap.Int32("id", id), zap.Error(err))
		return err
	}

	s.log.Info("user deleted", zap.Int32("id", id))
	return nil
}

func (s *userService) ListUsers(ctx context.Context, page, limit int) (models.PaginatedUsersResponse, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 10
	}

	offset := (page - 1) * limit

	users, err := s.repo.List(ctx, db.ListUsersParams{
		Limit:  int32(limit),
		Offset: int32(offset),
	})
	if err != nil {
		s.log.Error("failed to list users", zap.Error(err))
		return models.PaginatedUsersResponse{}, err
	}

	total, err := s.repo.Count(ctx)
	if err != nil {
		return models.PaginatedUsersResponse{}, err
	}

	data := make([]models.UserWithAgeResponse, 0, len(users))
	for _, u := range users {
		data = append(data, toUserWithAgeResponse(u))
	}

	totalPages := int(total) / limit
	if int(total)%limit != 0 {
		totalPages++
	}

	return models.PaginatedUsersResponse{
		Data:       data,
		Total:      total,
		Page:       page,
		Limit:      limit,
		TotalPages: totalPages,
	}, nil
}

// --- helpers ---

func toUserResponse(u db.User) models.UserResponse {
	return models.UserResponse{
		ID:   u.ID,
		Name: u.Name,
		DOB:  u.Dob.Format(dobLayout),
	}
}

func toUserWithAgeResponse(u db.User) models.UserWithAgeResponse {
	return models.UserWithAgeResponse{
		ID:   u.ID,
		Name: u.Name,
		DOB:  u.Dob.Format(dobLayout),
		Age:  models.CalculateAge(u.Dob),
	}
}

// ErrNotFound is returned when a requested resource does not exist.
var ErrNotFound = errors.New("user not found")
