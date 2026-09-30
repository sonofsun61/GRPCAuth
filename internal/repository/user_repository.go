package repository

import (
	"authservice/internal/entity"
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type UserRepository struct {
	pool *pgxpool.Pool
}

func NewUserRepository(pool *pgxpool.Pool) *UserRepository {
	return &UserRepository{
		pool: pool,
	}
}

func (r *UserRepository) CreateUser(ctx context.Context, user entity.User) (uuid.UUID, error) {
	query := `
			insert into users (
				email, 
				name, 
				surname, 
				phone_number, 
				password_hash, 
				password_salt)
			values ($1, $2, $3, $4, $5, $6)
			returning id`
	var newUserID uuid.UUID
	err := r.pool.QueryRow(ctx, query,
		user.Email, user.Name,
		user.Surname, user.PhoneNumber,
		user.PasswordHash, user.PasswordSalt).Scan(&newUserID)
	if err != nil {
		return uuid.Nil, fmt.Errorf("failed to insert new user information: %w", err)
	}
	return newUserID, nil
}

