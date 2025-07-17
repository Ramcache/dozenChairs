package repository

import (
	"context"
	"database/sql"
	"dozenChairs/internal/models"
	"errors"

	"github.com/jackc/pgx/v5/pgxpool"
)

type UserRepository struct {
	db *pgxpool.Pool
}

func NewUserRepository(db *pgxpool.Pool) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) Create(ctx context.Context, user *models.User) error {
	query := `
        INSERT INTO users (username, password_hash, email, full_name, phone, address, role, created_at, email_verified)
        VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
        RETURNING id
    `
	err := r.db.QueryRow(ctx, query,
		user.Username,
		user.PasswordHash,
		user.Email,
		user.FullName,
		user.Phone,
		user.Address,
		user.Role,
		user.CreatedAt,
		user.EmailVerified,
	).Scan(&user.ID)
	return err
}

func (r *UserRepository) GetByUsername(ctx context.Context, username string) (*models.User, error) {
	user := &models.User{}
	query := `SELECT id, username, password_hash, email, full_name, phone, address, role, created_at, email_verified
              FROM users WHERE username = $1`
	err := r.db.QueryRow(ctx, query, username).Scan(
		&user.ID,
		&user.Username,
		&user.PasswordHash,
		&user.Email,
		&user.FullName,
		&user.Phone,
		&user.Address,
		&user.Role,
		&user.CreatedAt,
		&user.EmailVerified,
	)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (r *UserRepository) GetByEmail(ctx context.Context, email string) (*models.User, error) {
	user := &models.User{}
	query := `SELECT id, username, password_hash, email, full_name, phone, address, role, created_at, email_verified
              FROM users WHERE email = $1`
	err := r.db.QueryRow(ctx, query, email).Scan(
		&user.ID,
		&user.Username,
		&user.PasswordHash,
		&user.Email,
		&user.FullName,
		&user.Phone,
		&user.Address,
		&user.Role,
		&user.CreatedAt,
		&user.EmailVerified,
	)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (r *UserRepository) GetByID(ctx context.Context, id int) (*models.User, error) {
	user := &models.User{}
	query := `SELECT id, username, password_hash, email, full_name, phone, address, role, created_at, email_verified
              FROM users WHERE id = $1`
	err := r.db.QueryRow(ctx, query, id).Scan(
		&user.ID,
		&user.Username,
		&user.PasswordHash,
		&user.Email,
		&user.FullName,
		&user.Phone,
		&user.Address,
		&user.Role,
		&user.CreatedAt,
		&user.EmailVerified,
	)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return user, nil
}
