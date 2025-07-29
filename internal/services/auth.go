package services

import (
	"dozenChairs/internal/auth"
	"dozenChairs/internal/dto"
	"dozenChairs/internal/models"
	"dozenChairs/internal/repository"
	security "dozenChairs/pkg/security"
	"fmt"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"time"
)

type AuthService interface {
	Register(input dto.RegisterRequest) (*models.User, error)
	Login(input dto.LoginRequest, jwt *auth.JWTManager, ip, ua string) (*models.User, string, string, error)
	ValidateSession(userID, tokenHash string) error
	Logout(tokenHash string) error
	Me(userID string) (*models.User, error)
}

type authService struct {
	userRepo    repository.UserRepository
	sessionRepo repository.SessionRepository
}

func NewAuthService(r repository.UserRepository, sR repository.SessionRepository) AuthService {
	return &authService{userRepo: r,
		sessionRepo: sR}
}

func (s *authService) Register(input dto.RegisterRequest) (*models.User, error) {
	hashed, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	user := &models.User{
		ID:           uuid.NewString(),
		Username:     input.Username,
		Email:        input.Email,
		PasswordHash: string(hashed),
		Role:         "user",
		CreatedAt:    time.Now(),
	}

	if err := s.userRepo.Create(user); err != nil {
		return nil, err
	}

	return user, nil
}

func (s *authService) Login(input dto.LoginRequest, jwt *auth.JWTManager, ip, ua string) (*models.User, string, string, error) {
	user, err := s.userRepo.GetByEmail(input.Email)
	if err != nil {
		return nil, "", "", fmt.Errorf("user not found")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(input.Password)); err != nil {
		return nil, "", "", fmt.Errorf("invalid credentials")
	}

	accessToken, err := jwt.GenerateAccess(user.ID, user.Role)
	if err != nil {
		return nil, "", "", err
	}

	refreshToken, err := jwt.GenerateRefresh(user.ID)
	if err != nil {
		return nil, "", "", err
	}

	hash := security.SHA256Sum(refreshToken)

	session := &models.Session{
		ID:        uuid.NewString(),
		UserID:    user.ID,
		TokenHash: hash,
		UserAgent: ua,
		IPAddress: ip,
		ExpiresAt: time.Now().Add(jwt.RefreshTTL),
		CreatedAt: time.Now(),
	}

	if err := s.sessionRepo.Create(session); err != nil {
		return nil, "", "", err
	}

	return user, refreshToken, accessToken, nil
}

func (s *authService) ValidateSession(userID, tokenHash string) error {
	// Простая проверка — есть ли сессия в БД по userID и хешу
	sessions, err := s.sessionRepo.FindByUserID(userID)
	if err != nil {
		return err
	}

	for _, sess := range sessions {
		if sess.TokenHash == tokenHash && sess.ExpiresAt.After(time.Now()) {
			return nil // валидная сессия найдена
		}
	}

	return fmt.Errorf("session not found or expired")
}

func (s *authService) Logout(tokenHash string) error {
	return s.sessionRepo.DeleteByTokenHash(tokenHash)
}

func (s *authService) Me(userID string) (*models.User, error) {
	return s.userRepo.GetByID(userID)
}
