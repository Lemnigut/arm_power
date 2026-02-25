package repository

import (
	"context"

	"arm_back/internal/model"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type UserRepository interface {
	Create(ctx context.Context, user *model.User) error
	GetByUsername(ctx context.Context, username string) (*model.User, error)
	GetByID(ctx context.Context, id uuid.UUID) (*model.User, error)
}

type pgUserRepo struct {
	pool *pgxpool.Pool
}

func NewUserRepository(pool *pgxpool.Pool) UserRepository {
	return &pgUserRepo{pool: pool}
}

func (r *pgUserRepo) Create(ctx context.Context, user *model.User) error {
	_, err := r.pool.Exec(ctx,
		`INSERT INTO users (id, username, password_hash) VALUES ($1, $2, $3)`,
		user.ID, user.Username, user.PasswordHash,
	)
	return err
}

func (r *pgUserRepo) GetByUsername(ctx context.Context, username string) (*model.User, error) {
	var u model.User
	err := r.pool.QueryRow(ctx,
		`SELECT id, username, password_hash, created_at FROM users WHERE username = $1`,
		username,
	).Scan(&u.ID, &u.Username, &u.PasswordHash, &u.CreatedAt)
	if err == pgx.ErrNoRows {
		return nil, model.ErrNotFound
	}
	return &u, err
}

func (r *pgUserRepo) GetByID(ctx context.Context, id uuid.UUID) (*model.User, error) {
	var u model.User
	err := r.pool.QueryRow(ctx,
		`SELECT id, username, password_hash, created_at FROM users WHERE id = $1`,
		id,
	).Scan(&u.ID, &u.Username, &u.PasswordHash, &u.CreatedAt)
	if err == pgx.ErrNoRows {
		return nil, model.ErrNotFound
	}
	return &u, err
}
