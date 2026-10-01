package repository

import (
	"authservice/internal/entity"
	"context"
	"fmt"

	"github.com/google/uuid"
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

func (r *UserRepository) GetUserByEmail(ctx context.Context, email string) (entity.User, error) {
	query := `
			select id, email, name, surname, phone_number, 
				password_hash, password_salt, created_at
			from users
			where email = $1
	`
	rows, err := r.pool.Query(ctx, query, email)
	if err != nil {
		return entity.User{}, fmt.Errorf("failed to get user by email: %w", err)
	}
	userInfo, err := pgx.CollectExactlyOneRow(rows, pgx.RowToStructByName[entity.User])
	if err != nil {
		return entity.User{}, fmt.Errorf("failed to place user data into struct: %w", err)
	}
	return userInfo, nil
}

func (r *UserRepository) UpdatePassword(ctx context.Context, email, newPasswordHash, newPasswordSalt string) error {
	query := `update users set password_hash = $1, password_salt = $2 where email = $3`
	tag, err := r.pool.Exec(ctx, query, newPasswordHash, newPasswordSalt, email)
	if err != nil {
		return fmt.Errorf("failed to update user's password: %w", err)
	}
	if tag.RowsAffected() == 0 {
		return pgx.ErrNoRows
	}
	return nil
}