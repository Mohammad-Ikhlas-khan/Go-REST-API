package repository

import (
	"context"
	"database/sql"

	db "github.com/example/go-user-api/db/sqlc"
)

//go:generate mockgen -source=user_repository.go -destination=mock_user_repository.go -package=repository

// UserRepository defines all DB operations for users.
type UserRepository interface {
	Create(ctx context.Context, arg db.CreateUserParams) (db.User, error)
	GetByID(ctx context.Context, id int32) (db.User, error)
	Update(ctx context.Context, arg db.UpdateUserParams) (db.User, error)
	Delete(ctx context.Context, id int32) error
	List(ctx context.Context, arg db.ListUsersParams) ([]db.User, error)
	Count(ctx context.Context) (int64, error)
}

type userRepository struct {
	q db.Querier
}

// NewUserRepository creates a new UserRepository backed by *db.Queries.
func NewUserRepository(sqlDB *sql.DB) UserRepository {
	return &userRepository{q: db.New(sqlDB)}
}

func (r *userRepository) Create(ctx context.Context, arg db.CreateUserParams) (db.User, error) {
	return r.q.CreateUser(ctx, arg)
}

func (r *userRepository) GetByID(ctx context.Context, id int32) (db.User, error) {
	return r.q.GetUserByID(ctx, id)
}

func (r *userRepository) Update(ctx context.Context, arg db.UpdateUserParams) (db.User, error) {
	return r.q.UpdateUser(ctx, arg)
}

func (r *userRepository) Delete(ctx context.Context, id int32) error {
	return r.q.DeleteUser(ctx, id)
}

func (r *userRepository) List(ctx context.Context, arg db.ListUsersParams) ([]db.User, error) {
	return r.q.ListUsers(ctx, arg)
}

func (r *userRepository) Count(ctx context.Context) (int64, error) {
	return r.q.CountUsers(ctx)
}
