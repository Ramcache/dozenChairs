package handlers

import (
	"dozenChairs/internal/auth"
	"dozenChairs/internal/dto"
	"dozenChairs/internal/middlewares"
	"dozenChairs/internal/services"
	"dozenChairs/pkg/httphelper"
	"dozenChairs/pkg/logger"
	security "dozenChairs/pkg/secutiry"
	"dozenChairs/pkg/validation"
	"encoding/json"
	"go.uber.org/zap"
	"net/http"
	"time"
)

type AuthHandler struct {
	service    services.AuthService
	logger     logger.Logger
	jwtManager *auth.JWTManager
}

func NewAuthHandler(s services.AuthService, l logger.Logger, jwtManager *auth.JWTManager) *AuthHandler {
	return &AuthHandler{
		service:    s,
		logger:     l,
		jwtManager: jwtManager,
	}
}

// Register godoc
// @Summary      Регистрация пользователя
// @Description  Создаёт нового пользователя и возвращает access и refresh токены
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        input  body      dto.RegisterRequest  true  "Данные регистрации"
// @Success      201    {object}  dto.AuthResponse
// @Failure      400    {object}  dto.ErrorResponse
// @Failure      409    {object}  dto.ErrorResponse
// @Router       /api/v1/auth/register [post]
func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	var input dto.RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		h.logger.Error("invalid register payload", zap.Error(err))
		httphelper.WriteError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	if err := validation.ValidateStruct(input); err != nil {
		h.logger.Error("register validation failed", zap.Error(err))
		httphelper.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}

	user, err := h.service.Register(input)
	if err != nil {
		h.logger.Error("register failed", zap.Error(err))
		httphelper.WriteError(w, http.StatusInternalServerError, "Failed to register")
		return
	}

	h.logger.Info("user registered", zap.String("id", user.ID), zap.String("email", user.Email))
	httphelper.WriteSuccess(w, http.StatusCreated, map[string]string{
		"id": user.ID,
	})
}

// Login godoc
// @Summary      Авторизация пользователя
// @Description  Принимает email и пароль, возвращает access и refresh токены
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        input  body      dto.LoginRequest  true  "Данные авторизации"
// @Success      200    {object}  dto.AuthResponse
// @Failure      400    {object}  dto.ErrorResponse
// @Failure      401    {object}  dto.ErrorResponse
// @Router       /api/v1/auth/login [post]
func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req dto.LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httphelper.WriteError(w, http.StatusBadRequest, "Invalid JSON")
		return
	}
	if err := validation.ValidateStruct(req); err != nil {
		httphelper.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}

	user, refreshToken, accessToken, err := h.service.Login(req, h.jwtManager, r.RemoteAddr, r.UserAgent())
	if err != nil {
		h.logger.Error("login failed", zap.Error(err))
		httphelper.WriteError(w, http.StatusUnauthorized, "Invalid credentials")
		return
	}

	// Set refresh token in cookie
	http.SetCookie(w, &http.Cookie{
		Name:     "refresh_token",
		Value:    refreshToken,
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteStrictMode,
		Path:     "/",
		Expires:  time.Now().Add(h.jwtManager.RefreshTTL),
	})

	// return access token
	h.logger.Info("user logged in", zap.String("id", user.ID))
	httphelper.WriteSuccess(w, http.StatusOK, map[string]string{
		"access_token": accessToken,
	})
}

// Refresh godoc
// @Summary      Обновление access токена
// @Description  Обновляет access токен по refresh токену из куки
// @Tags         auth
// @Accept       json
// @Produce      json
// @Success      200  {object}  dto.AccessTokenResponse
// @Failure      401  {object}  dto.ErrorResponse
// @Router       /api/v1/auth/refresh [post]
func (h *AuthHandler) Refresh(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("refresh_token")
	if err != nil || cookie.Value == "" {
		httphelper.WriteError(w, http.StatusUnauthorized, "No refresh token")
		return
	}

	userID, err := h.jwtManager.ValidateRefresh(cookie.Value)
	if err != nil {
		httphelper.WriteError(w, http.StatusUnauthorized, "Invalid refresh token")
		return
	}

	hash := security.SHA256Sum(cookie.Value)
	if err := h.service.ValidateSession(userID, hash); err != nil {
		httphelper.WriteError(w, http.StatusUnauthorized, "Session not found or expired")
		return
	}

	user, err := h.service.Me(userID)
	if err != nil {
		httphelper.WriteError(w, http.StatusInternalServerError, "Failed to load user")
		return
	}

	accessToken, err := h.jwtManager.GenerateAccess(userID, user.Role)
	if err != nil {
		httphelper.WriteError(w, http.StatusInternalServerError, "Failed to generate access token")
		return
	}

	httphelper.WriteSuccess(w, http.StatusOK, map[string]string{
		"access_token": accessToken,
	})
}

// Logout godoc
// @Summary      Выход пользователя
// @Description  Удаляет refresh токен из хранилища и куки
// @Tags         auth
// @Accept       json
// @Produce      json
// @Success      204  "No Content"
// @Failure      401  {object}  dto.ErrorResponse
// @Router       /api/v1/auth/logout [post]
func (h *AuthHandler) Logout(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("refresh_token")
	if err != nil || cookie.Value == "" {
		httphelper.WriteError(w, http.StatusBadRequest, "No refresh token")
		return
	}

	hash := security.SHA256Sum(cookie.Value)
	if err := h.service.Logout(hash); err != nil {
		h.logger.Error("logout failed", zap.Error(err))
		httphelper.WriteError(w, http.StatusInternalServerError, "Failed to logout")
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "refresh_token",
		Value:    "",
		Path:     "/",
		Expires:  time.Unix(0, 0),
		MaxAge:   -1,
		HttpOnly: true,
	})

	httphelper.WriteSuccess(w, http.StatusOK, map[string]string{"message": "Logged out"})
}

// Me godoc
// @Summary      Получить информацию о текущем пользователе
// @Description  Возвращает профиль пользователя по access токену
// @Tags         auth
// @Security     BearerAuth
// @Produce      json
// @Success      200  {object}  dto.UserResponse
// @Failure      401  {object}  dto.ErrorResponse
// @Router       /api/v1/auth/me [get]
func (h *AuthHandler) Me(w http.ResponseWriter, r *http.Request) {
	val := r.Context().Value(middlewares.UserID())
	if val == nil {
		httphelper.WriteError(w, http.StatusUnauthorized, "Missing user context")
		return
	}

	userID, ok := val.(string)
	if !ok {
		httphelper.WriteError(w, http.StatusInternalServerError, "Invalid user ID in context")
		return
	}

	user, err := h.service.Me(userID)
	if err != nil {
		h.logger.Error("failed to fetch user", zap.Error(err))
		httphelper.WriteError(w, http.StatusInternalServerError, "Failed to load user")
		return
	}

	httphelper.WriteSuccess(w, http.StatusOK, map[string]interface{}{
		"id":       user.ID,
		"username": user.Username,
		"email":    user.Email,
		"role":     user.Role,
	})
}
