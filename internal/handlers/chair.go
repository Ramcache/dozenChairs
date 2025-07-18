package handlers

import (
	"context"
	"dozenChairs/internal/models"
	"dozenChairs/internal/services"
	"encoding/json"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
)

type ChairHandler struct {
	chairService *services.ChairService
}

func NewChairHandler(service *services.ChairService) *ChairHandler {
	return &ChairHandler{chairService: service}
}

// Create
// @Summary Создать стул
// @Tags chairs
// @Accept json
// @Produce json
// @Param chair body models.Chair true "Стул"
// @Success 201 {object} models.Chair
// @Router /chairs [post]
func (h *ChairHandler) Create(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	var chair models.Chair
	if err := json.NewDecoder(r.Body).Decode(&chair); err != nil {
		http.Error(w, "Некорректные данные", http.StatusBadRequest)
		return
	}

	if err := h.chairService.Create(ctx, &chair); err != nil {
		http.Error(w, "Ошибка создания", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(chair)
}

// GetBySlug
// @Summary Получить стул по slug
// @Tags chairs
// @Produce json
// @Param slug path string true "Slug"
// @Success 200 {object} models.Chair
// @Failure 404 {string} string "Не найден"
// @Router /chairs/{slug} [get]
func (h *ChairHandler) GetBySlug(w http.ResponseWriter, r *http.Request) {
	slug := chi.URLParam(r, "slug")
	if slug == "" {
		http.Error(w, "slug обязателен", http.StatusBadRequest)
		return
	}
	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()
	chair, err := h.chairService.GetBySlug(ctx, slug)
	if err != nil || chair == nil {
		http.Error(w, "Товар не найден", http.StatusNotFound)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(chair)
}

// GetAll
// @Summary Получить все стулья
// @Tags chairs
// @Produce json
// @Success 200 {array} models.Chair
// @Router /chairs [get]
func (h *ChairHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()
	chairs, err := h.chairService.GetAll(ctx)
	if err != nil {
		http.Error(w, "Ошибка получения списка", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(chairs)
}

// UpdateBySlug
// @Summary Обновить стул по slug
// @Tags chairs
// @Accept json
// @Param slug path string true "Slug"
// @Param chair body models.Chair true "Обновляемые данные"
// @Success 204 {string} string ""
// @Failure 400 {string} string "Некорректные данные"
// @Failure 404 {string} string "Не найден"
// @Router /chairs/{slug} [patch]
func (h *ChairHandler) UpdateBySlug(w http.ResponseWriter, r *http.Request) {
	slug := chi.URLParam(r, "slug")
	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	var chair models.Chair
	if err := json.NewDecoder(r.Body).Decode(&chair); err != nil {
		http.Error(w, "Некорректные данные", http.StatusBadRequest)
		return
	}

	err := h.chairService.UpdateBySlug(ctx, slug, &chair)
	if err != nil {
		http.Error(w, "Ошибка обновления", http.StatusInternalServerError)
		return
	}

	// Получаем обновлённый объект (уже после изменений)
	updated, err := h.chairService.GetBySlug(ctx, slug)
	if err != nil || updated == nil {
		http.Error(w, "Ошибка получения обновленного объекта", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(updated)
}

// DeleteBySlug
// @Summary Удалить стул по slug
// @Tags chairs
// @Param slug path string true "Slug"
// @Success 204 {string} string ""
// @Failure 404 {string} string "Не найден"
// @Router /chairs/{slug} [delete]
func (h *ChairHandler) DeleteBySlug(w http.ResponseWriter, r *http.Request) {
	slug := chi.URLParam(r, "slug")
	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	// Получаем объект ДО удаления
	deleted, err := h.chairService.GetBySlug(ctx, slug)
	if err != nil || deleted == nil {
		http.Error(w, "Объект не найден", http.StatusNotFound)
		return
	}

	// Удаляем
	err = h.chairService.DeleteBySlug(ctx, slug)
	if err != nil {
		http.Error(w, "Ошибка удаления", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(deleted)
}
