package repository

import (
	"context"
	"dozenChairs/internal/models"
	"github.com/jackc/pgx/v5/pgxpool"
)

type UserRepository interface {
	Create(u *models.User) error
	GetByEmail(email string) (*models.User, error)
	GetByID(userID string) (*models.User, error)
}

type userRepo struct {
	db *pgxpool.Pool
}

func NewUserRepo(db *pgxpool.Pool) UserRepository {
	return &userRepo{db: db}
}

func (r *userRepo) Create(u *models.User) error {
	_, err := r.db.Exec(
		context.Background(),
		`INSERT INTO users (id, email, username, password_hash, role, created_at) VALUES ($1, $2, $3, $4, $5, $6)`,
		u.ID, u.Email, u.Username, u.PasswordHash, u.Role, u.CreatedAt,
	)
	return err
}

func (r *userRepo) GetByEmail(email string) (*models.User, error) {
	var u models.User
	err := r.db.QueryRow(
		context.Background(),
		`SELECT id, email, username, password_hash, role, created_at FROM users WHERE email = $1`,
		email,
	).Scan(&u.ID, &u.Email, &u.Username, &u.PasswordHash, &u.Role, &u.CreatedAt)

	if err != nil {
		return nil, err
	}
	return &u, nil
}

func (r *userRepo) GetByID(userID string) (*models.User, error) {
	row := r.db.QueryRow(context.Background(), `
		SELECT id, username, email, password_hash, role, created_at
		FROM users
		WHERE id = $1
	`, userID)

	var u models.User
	if err := row.Scan(&u.ID, &u.Username, &u.Email, &u.PasswordHash, &u.Role, &u.CreatedAt); err != nil {
		return nil, err
	}
	return &u, nil
}
