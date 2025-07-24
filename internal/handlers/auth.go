package handlers

import (
	"context"
	middleware "dozenChairs/internal/middlewares"
	"dozenChairs/internal/models"
	"dozenChairs/internal/services"
	"dozenChairs/pkg/utils"
	"encoding/json"
	"net/http"
	"time"
)

type AuthHandler struct {
	userService *services.UserService
}

func NewAuthHandler(userService *services.UserService) *AuthHandler {
	return &AuthHandler{userService: userService}
}

type RegisterRequest struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
	FullName string `json:"full_name"`
	Phone    string `json:"phone,omitempty"`
	Address  string `json:"address,omitempty"`
}

type LoginRequest struct {
	EmailOrUsername string `json:"email_or_username"`
	Password        string `json:"password"`
}

type LoginResponse struct {
	Token string      `json:"token"`
	User  interface{} `json:"user"`
}

// Register godoc
// @Summary      Регистрация пользователя
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        data body RegisterRequest true "Данные для регистрации"
// @Success      201  {object}  models.User
// @Failure      400  {string}  string  "Некорректный запрос"
// @Router       /register [post]
func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	var req RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Некорректный запрос", http.StatusBadRequest)
		return
	}
	//
	// Собираем модель для передачи в сервис
	user := &models.User{
		Username:     req.Username,
		Email:        req.Email,
		PasswordHash: req.Password, // Передаем пароль — сервис его захеширует
		FullName:     req.FullName,
		Phone:        req.Phone,
		Address:      req.Address,
	}

	createdUser, err := h.userService.Register(ctx, user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Не показываем пароль
	createdUser.PasswordHash = ""

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(createdUser)
}

// Login godoc
// @Summary      Логин пользователя
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        data body LoginRequest true "Данные для логина"
// @Success      200  {object}  LoginResponse
// @Failure      401  {string}  string  "Неверные учетные данные"
// @Router       /login [post]
func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	var req LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Некорректный запрос", http.StatusBadRequest)
		return
	}

	user, err := h.userService.Login(ctx, req.EmailOrUsername, req.Password)
	if err != nil {
		http.Error(w, "Неверные учетные данные", http.StatusUnauthorized)
		return
	}

	token, err := utils.GenerateJWT(user)
	if err != nil {
		http.Error(w, "Ошибка генерации токена", http.StatusInternalServerError)
		return
	}

	user.PasswordHash = "" // не возвращаем хеш пароля!

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(LoginResponse{
		Token: token,
		User:  user,
	})
}

// Profile godoc
// @Summary      Получение профиля пользователя
// @Tags         auth
// @Security     BearerAuth
// @Produce      json
// @Success      200  {object}  models.User
// @Failure      401  {string}  string  "Не авторизован"
// @Router       /profile [get]
func (h *AuthHandler) Profile(w http.ResponseWriter, r *http.Request) {
	user, ok := r.Context().Value(middleware.UserCtxKey).(*middleware.JWTUser)
	if !ok {
		http.Error(w, "Пользователь не найден", http.StatusUnauthorized)
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	dbUser, err := h.userService.GetByID(ctx, user.ID)
	if err != nil {
		http.Error(w, "Ошибка при получении профиля", http.StatusInternalServerError)
		return
	}
	if dbUser == nil {
		http.Error(w, "Пользователь не найден", http.StatusNotFound)
		return
	}
	dbUser.PasswordHash = ""

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(dbUser)
}
