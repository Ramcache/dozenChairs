package repository

import (
	"context"
	"dozenChairs/internal/models"
	"github.com/jackc/pgx/v5/pgxpool"
)

type SessionRepository interface {
	Create(s *models.Session) error
	DeleteByTokenHash(hash string) error
	DeleteAllForUser(userID string) error
	FindByUserID(userID string) ([]*models.Session, error)
}

type sessionRepo struct {
	db *pgxpool.Pool
}

func NewSessionRepo(db *pgxpool.Pool) SessionRepository {
	return &sessionRepo{db: db}
}

func (r *sessionRepo) Create(s *models.Session) error {
	_, err := r.db.Exec(context.Background(), `
		INSERT INTO user_sessions (id, user_id, token_hash, user_agent, ip_address, expires_at, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
	`, s.ID, s.UserID, s.TokenHash, s.UserAgent, s.IPAddress, s.ExpiresAt, s.CreatedAt)
	return err
}

func (r *sessionRepo) DeleteByTokenHash(hash string) error {
	_, err := r.db.Exec(context.Background(), `DELETE FROM user_sessions WHERE token_hash = $1`, hash)
	return err
}

func (r *sessionRepo) DeleteAllForUser(userID string) error {
	_, err := r.db.Exec(context.Background(), `DELETE FROM user_sessions WHERE user_id = $1`, userID)
	return err
}

func (r *sessionRepo) FindByUserID(userID string) ([]*models.Session, error) {
	rows, err := r.db.Query(context.Background(), `
		SELECT id, user_id, token_hash, user_agent, ip_address, expires_at, created_at
		FROM user_sessions
		WHERE user_id = $1
	`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var sessions []*models.Session
	for rows.Next() {
		var s models.Session
		if err := rows.Scan(
			&s.ID,
			&s.UserID,
			&s.TokenHash,
			&s.UserAgent,
			&s.IPAddress,
			&s.ExpiresAt,
			&s.CreatedAt,
		); err != nil {
			return nil, err
		}
		sessions = append(sessions, &s)
	}
	return sessions, nil
}
