package services

import (
	"context"
	"dozenChairs/internal/models"
	"dozenChairs/internal/repository"
	"errors"
	"strings"
	"time"

	"golang.org/x/crypto/bcrypt"
)

type UserService struct {
	userRepo *repository.UserRepository
}

func NewUserService(userRepo *repository.UserRepository) *UserService {
	return &UserService{userRepo: userRepo}
}

var (
	ErrEmailTaken         = errors.New("email уже используется")
	ErrUsernameTaken      = errors.New("username уже используется")
	ErrInvalidCredentials = errors.New("неверный email/username или пароль")
)

func (s *UserService) Register(ctx context.Context, input *models.User) (*models.User, error) {
	// Минимальная валидация
	input.Email = strings.TrimSpace(input.Email)
	input.Username = strings.TrimSpace(input.Username)

	if input.Email == "" || input.Username == "" || input.PasswordHash == "" {
		return nil, errors.New("email, username и password обязательны")
	}

	// Проверка уникальности email
	if user, _ := s.userRepo.GetByEmail(ctx, input.Email); user != nil {
		return nil, ErrEmailTaken
	}
	// Проверка уникальности username
	if user, _ := s.userRepo.GetByUsername(ctx, input.Username); user != nil {
		return nil, ErrUsernameTaken
	}

	// Хеширование пароля
	hash, err := bcrypt.GenerateFromPassword([]byte(input.PasswordHash), bcrypt.DefaultCost)
	if err != nil {
		return nil, errors.New("ошибка при хешировании пароля")
	}
	input.PasswordHash = string(hash)

	// Прочие поля
	input.CreatedAt = time.Now().UTC()
	input.Role = "user"
	input.EmailVerified = false

	// Сохраняем пользователя
	err = s.userRepo.Create(ctx, input)
	if err != nil {
		return nil, err
	}

	// Не возвращаем пароль наружу
	input.PasswordHash = ""
	return input, nil
}

func (s *UserService) Login(ctx context.Context, emailOrUsername, password string) (*models.User, error) {
	var user *models.User
	var err error

	// Можно определить, это email или username, или искать в обоих случаях
	if strings.Contains(emailOrUsername, "@") {
		user, err = s.userRepo.GetByEmail(ctx, emailOrUsername)
	} else {
		user, err = s.userRepo.GetByUsername(ctx, emailOrUsername)
	}

	if err != nil || user == nil {
		return nil, ErrInvalidCredentials
	}

	// Сравниваем пароль
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)); err != nil {
		return nil, ErrInvalidCredentials
	}

	// Всё ок, возвращаем пользователя
	return user, nil
}

func (s *UserService) GetByID(ctx context.Context, id int) (*models.User, error) {
	return s.userRepo.GetByID(ctx, id)
}

func (s *UserService) GetByUsername(ctx context.Context, username string) (*models.User, error) {
	return s.userRepo.GetByUsername(ctx, username)
}

func (s *UserService) GetByEmail(ctx context.Context, email string) (*models.User, error) {
	return s.userRepo.GetByEmail(ctx, email)
}

func (s *UserService) Create(ctx context.Context, user *models.User) error {
	return s.userRepo.Create(ctx, user)
}
