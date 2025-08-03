package handlers

import (
	"dozenChairs/internal/auth"
	"dozenChairs/internal/dto"
	"dozenChairs/internal/middlewares"
	"dozenChairs/internal/services"
	"dozenChairs/pkg/config"
	"dozenChairs/pkg/httphelper"
	"dozenChairs/pkg/logger"
	security "dozenChairs/pkg/security"
	"dozenChairs/pkg/validation"
	"encoding/json"
	"fmt"
	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
	"golang.org/x/oauth2"
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

	// Регистрируем пользователя
	user, err := h.service.Register(input)
	if err != nil {
		h.logger.Error("register failed", zap.Error(err))
		httphelper.WriteError(w, http.StatusInternalServerError, "Failed to register")
		return
	}

	// Генерируем токены (как при логине)
	refreshToken, accessToken, err := h.jwtManager.GenerateTokens(user.ID, user.Role)
	if err != nil {
		h.logger.Error("token generation failed", zap.Error(err))
		httphelper.WriteError(w, http.StatusInternalServerError, "Failed to generate tokens")
		return
	}

	// Устанавливаем refresh токен в cookie
	http.SetCookie(w, &http.Cookie{
		Name:     "refresh_token",
		Value:    refreshToken,
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteStrictMode,
		Path:     "/",
		Expires:  time.Now().Add(h.jwtManager.RefreshTTL),
	})

	// Лог и ответ
	h.logger.Info("user registered", zap.String("id", user.ID), zap.String("email", user.Email))
	httphelper.WriteSuccess(w, http.StatusOK, map[string]interface{}{
		"access_token": accessToken,
		"user": map[string]interface{}{
			"id":    user.ID,
			"email": user.Email,
			"role":  user.Role,
			"name":  user.Username,
		},
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
	httphelper.WriteSuccess(w, http.StatusOK, map[string]interface{}{
		"access_token": accessToken,
		"user": map[string]interface{}{
			"id":    user.ID,
			"email": user.Email,
			"role":  user.Role,
			"name":  user.Username,
		},
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

// OAuthCallback godoc
// @Summary      Callback от OAuth-провайдера
// @Description  Обрабатывает код, полученный от VK, Google или Yandex, и возвращает JWT токены
// @Tags         auth
// @Produce      json
// @Param        provider  path      string  true  "OAuth-провайдер"  Enums(google, yandex, vk)
// @Param        code      query     string  true  "Код от OAuth-провайдера"
// @Success      200       {object}  dto.AuthResponse
// @Failure      400       {object}  dto.ErrorResponse
// @Failure      500       {object}  dto.ErrorResponse
// @Router       /api/v1/auth/callback/{provider} [get]
func (h *AuthHandler) OAuthCallback(w http.ResponseWriter, r *http.Request) {
	provider := chi.URLParam(r, "provider")
	cfg, ok := config.OauthProviders[provider]
	if !ok {
		httphelper.WriteError(w, http.StatusBadRequest, "Unknown provider")
		return
	}

	code := r.URL.Query().Get("code")
	if code == "" {
		httphelper.WriteError(w, http.StatusBadRequest, "Missing code")
		return
	}

	ctx := r.Context()
	var email, username string

	switch provider {
	case "vk":
		// 1. Обмен кода на access_token
		resp, err := http.Get(fmt.Sprintf("https://oauth.vk.com/access_token?client_id=%s&client_secret=%s&redirect_uri=%s&code=%s",
			cfg.Config.ClientID,
			cfg.Config.ClientSecret,
			cfg.Config.RedirectURL,
			code,
		))
		if err != nil {
			h.logger.Error("VK token exchange failed", zap.Error(err))
			httphelper.WriteError(w, http.StatusInternalServerError, "VK token exchange failed")
			return
		}
		defer resp.Body.Close()

		var tokenData struct {
			AccessToken string `json:"access_token"`
			Email       string `json:"email"`
			UserID      int    `json:"user_id"`
		}
		if err := json.NewDecoder(resp.Body).Decode(&tokenData); err != nil {
			httphelper.WriteError(w, http.StatusInternalServerError, "Invalid VK token response")
			return
		}
		email = tokenData.Email

		// 2. Запрашиваем имя пользователя
		userInfoResp, err := http.Get(fmt.Sprintf("https://api.vk.com/method/users.get?user_ids=%d&fields=first_name,last_name&access_token=%s&v=5.131",
			tokenData.UserID, tokenData.AccessToken))
		if err != nil {
			httphelper.WriteError(w, http.StatusInternalServerError, "VK user info fetch failed")
			return
		}
		defer userInfoResp.Body.Close()

		var userInfo struct {
			Response []struct {
				FirstName string `json:"first_name"`
				LastName  string `json:"last_name"`
			} `json:"response"`
		}
		_ = json.NewDecoder(userInfoResp.Body).Decode(&userInfo)

		if len(userInfo.Response) > 0 {
			username = userInfo.Response[0].FirstName + " " + userInfo.Response[0].LastName
		}

	case "google", "yandex":
		// 1. Обмен кода на токен
		token, err := cfg.Config.Exchange(ctx, code)
		if err != nil {
			h.logger.Error("token exchange failed", zap.Error(err))
			httphelper.WriteError(w, http.StatusBadRequest, "Token exchange failed")
			return
		}

		client := cfg.Config.Client(ctx, token)
		resp, err := client.Get(cfg.UserInfoURL)
		if err != nil {
			h.logger.Error("user info request failed", zap.Error(err))
			httphelper.WriteError(w, http.StatusInternalServerError, "User info fetch failed")
			return
		}
		defer resp.Body.Close()

		var data map[string]interface{}
		if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
			httphelper.WriteError(w, http.StatusInternalServerError, "Invalid user info")
			return
		}

		switch provider {
		case "google":
			email, _ = data["email"].(string)
			username, _ = data["name"].(string)
		case "yandex":
			email, _ = data["default_email"].(string)
			username, _ = data["real_name"].(string)
		}
	}

	if email == "" || username == "" {
		httphelper.WriteError(w, http.StatusInternalServerError, "Failed to get user info")
		return
	}

	// 3. Авторизация через сервис
	user, refreshToken, accessToken, err := h.service.OAuthLogin(email, username, provider, h.jwtManager)
	if err != nil {
		h.logger.Error("oauth login failed", zap.Error(err))
		httphelper.WriteError(w, http.StatusInternalServerError, "OAuth login failed")
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "refresh_token",
		Value:    refreshToken,
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteStrictMode,
		Path:     "/",
		Expires:  time.Now().Add(h.jwtManager.RefreshTTL),
	})

	httphelper.WriteSuccess(w, http.StatusOK, map[string]interface{}{
		"access_token": accessToken,
		"user": map[string]interface{}{
			"id":    user.ID,
			"email": user.Email,
			"role":  user.Role,
			"name":  user.Username,
		},
	})
}

// BeginOAuth godoc
// @Summary      Перенаправление на OAuth-провайдера
// @Description  Редиректит пользователя на страницу авторизации Google, VK или Yandex
// @Tags         auth
// @Produce      json
// @Param        provider  path      string  true  "OAuth-провайдер"  Enums(google, yandex, vk)
// @Success      307       {string}  string  "Redirect"
// @Failure      400       {object}  dto.ErrorResponse
// @Router       /api/v1/auth/oauth/{provider} [get]
func (h *AuthHandler) BeginOAuth(w http.ResponseWriter, r *http.Request) {
	provider := chi.URLParam(r, "provider")

	cfg, ok := config.OauthProviders[provider]
	if !ok {
		httphelper.WriteError(w, http.StatusBadRequest, "Unknown provider")
		return
	}

	authURL := cfg.Config.AuthCodeURL("state-token", oauth2.AccessTypeOffline)
	http.Redirect(w, r, authURL, http.StatusTemporaryRedirect)
}
