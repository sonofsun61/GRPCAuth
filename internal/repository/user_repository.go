package repository

import (
	"authservice/internal/entity"
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
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

func (r *UserRepository) CreateUser(ctx context.Context, user entity.User) error {
	query := `
			insert into users (
				email, 
				name, 
				surname, 
				phone_number, 
				password_hash, 
				password_salt)
			values ($1, $2, $3, $4, $5, $6)`
	tag, err := r.pool.Exec(ctx, query, 
		user.Email, user.Name, 
		user.Surname, user.PhoneNumber, 
		user.PasswordHash, user.PasswordSalt)
	if err != nil {
		return fmt.Errorf("failed to insert new user information: %w", err)
	}
	if tag.RowsAffected() == 0 {
		return pgx.ErrNoRows
	}
	return nil
}