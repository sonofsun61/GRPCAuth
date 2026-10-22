//go:build integration

package repository_test

import (
	"authservice/internal/entity"
	"authservice/internal/repository"
	"context"
	"errors"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"
)

const uniqueViolationCode = "23505"

// pool is shared by all tests in the package; the container is started once in TestMain.
var pool *pgxpool.Pool

func TestMain(m *testing.M) {
	os.Exit(runTests(m))
}

func runTests(m *testing.M) int {
	ctx := context.Background()

	container, err := postgres.Run(ctx, "postgres:16",
		postgres.WithDatabase("auth_test"),
		postgres.WithUsername("test"),
		postgres.WithPassword("test"),
		postgres.WithInitScripts(filepath.Join("..", "..", "migrations", "000001_create_users.up.sql")),
		testcontainers.WithWaitStrategy(
			wait.ForLog("database system is ready to accept connections").
				WithOccurrence(2).
				WithStartupTimeout(time.Minute),
		),
	)
	if err != nil {
		panic(err)
	}
	defer func() {
		_ = testcontainers.TerminateContainer(container)
	}()

	connString, err := container.ConnectionString(ctx, "sslmode=disable")
	if err != nil {
		panic(err)
	}
	pool, err = pgxpool.New(ctx, connString)
	if err != nil {
		panic(err)
	}
	defer pool.Close()

	return m.Run()
}

func newRepo(t *testing.T) *repository.UserRepository {
	t.Helper()
	_, err := pool.Exec(context.Background(), "truncate table users")
	require.NoError(t, err)
	return repository.NewUserRepository(pool)
}

func newUser(email, phone string) entity.User {
	return entity.User{
		Email:        email,
		Name:         "Name",
		Surname:      "Surname",
		PhoneNumber:  phone,
		PasswordHash: "hash",
		PasswordSalt: "salt",
	}
}

func TestCreateUser_GetUserByEmail(t *testing.T) {
	repo := newRepo(t)
	ctx := context.Background()
	user := newUser("example@mail.com", "+79999999999")

	id, err := repo.CreateUser(ctx, user)
	require.NoError(t, err)
	require.NotEqual(t, uuid.Nil, id)

	got, err := repo.GetUserByEmail(ctx, user.Email)
	require.NoError(t, err)
	assert.Equal(t, id, got.ID)
	assert.Equal(t, user.Email, got.Email)
	assert.Equal(t, user.Name, got.Name)
	assert.Equal(t, user.Surname, got.Surname)
	assert.Equal(t, user.PhoneNumber, got.PhoneNumber)
	assert.Equal(t, user.PasswordHash, got.PasswordHash)
	assert.Equal(t, user.PasswordSalt, got.PasswordSalt)
	assert.WithinDuration(t, time.Now(), got.CreatedAt, time.Minute)
}

func TestCreateUser_Duplicate(t *testing.T) {
	tests := []struct {
		name   string
		second entity.User
	}{
		{"same email", newUser("example@mail.com", "+70000000000")},
		{"same phone", newUser("other@mail.com", "+79999999999")},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := newRepo(t)
			ctx := context.Background()
			_, err := repo.CreateUser(ctx, newUser("example@mail.com", "+79999999999"))
			require.NoError(t, err)

			id, err := repo.CreateUser(ctx, tt.second)
			var pgErr *pgconn.PgError
			require.True(t, errors.As(err, &pgErr), "expected *pgconn.PgError, got %v", err)
			assert.Equal(t, uniqueViolationCode, pgErr.Code)
			assert.Equal(t, uuid.Nil, id)
		})
	}
}

func TestGetUserByEmail_NotFound(t *testing.T) {
	repo := newRepo(t)

	got, err := repo.GetUserByEmail(context.Background(), "missing@mail.com")
	require.ErrorIs(t, err, pgx.ErrNoRows)
	assert.Equal(t, entity.User{}, got)
}

func TestUpdatePassword_Success(t *testing.T) {
	repo := newRepo(t)
	ctx := context.Background()
	user := newUser("example@mail.com", "+79999999999")
	_, err := repo.CreateUser(ctx, user)
	require.NoError(t, err)

	err = repo.UpdatePassword(ctx, user.Email, "new-hash", "new-salt")
	require.NoError(t, err)

	got, err := repo.GetUserByEmail(ctx, user.Email)
	require.NoError(t, err)
	assert.Equal(t, "new-hash", got.PasswordHash)
	assert.Equal(t, "new-salt", got.PasswordSalt)
}

func TestUpdatePassword_NotFound(t *testing.T) {
	repo := newRepo(t)

	err := repo.UpdatePassword(context.Background(), "missing@mail.com", "new-hash", "new-salt")
	require.ErrorIs(t, err, pgx.ErrNoRows)
}
