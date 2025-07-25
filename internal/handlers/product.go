package handlers

import (
	"dozenChairs/internal/models"
	"dozenChairs/internal/repository"
	"dozenChairs/internal/services"
	"dozenChairs/pkg/httphelper"
	"dozenChairs/pkg/logger"
	"dozenChairs/pkg/validation"
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
)

type ProductHandler struct {
	service services.ProductService
	logger  logger.Logger
}

func NewProductHandler(s services.ProductService, l logger.Logger) *ProductHandler {
	return &ProductHandler{
		service: s,
		logger:  l,
	}
}

// Create
// @Summary Создать товар
// @Tags Products
// @Accept json
// @Produce json
// @Param product body models.Product true "Product JSON"
// @Success 201 {object} httphelper.APIResponse
// @Failure 400 {object} httphelper.APIResponse
// @Failure 500 {object} httphelper.APIResponse
// @Router /products [post]
func (h *ProductHandler) Create(w http.ResponseWriter, r *http.Request) {
	var p models.Product
	if err := json.NewDecoder(r.Body).Decode(&p); err != nil {
		h.logger.Error("failed to decode product", zap.Error(err))
		httphelper.WriteError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	if err := validation.ValidateStruct(p); err != nil {
		h.logger.Error("validation failed", zap.Error(err))
		httphelper.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}

	if err := h.service.Create(&p); err != nil {
		h.logger.Error("failed to create product", zap.String("id", p.ID), zap.Error(err))
		httphelper.WriteError(w, http.StatusInternalServerError, "Failed to create product")
		return
	}

	h.logger.Info("product created successfully", zap.String("id", p.ID), zap.String("slug", p.Slug))
	httphelper.WriteSuccess(w, http.StatusCreated, p)
}

// GetBySlug
// @Summary Получить товар по slug
// @Tags Products
// @Produce json
// @Param slug path string true "Slug товара"
// @Success 200 {object} httphelper.APIResponse
// @Failure 404 {object} httphelper.APIResponse
// @Failure 500 {object} httphelper.APIResponse
// @Router /products/{slug} [get]
func (h *ProductHandler) GetBySlug(w http.ResponseWriter, r *http.Request) {
	slug := chi.URLParam(r, "slug")

	p, err := h.service.GetBySlug(slug)
	if err != nil {
		h.logger.Error("product not found", zap.String("slug", slug), zap.Error(err))
		httphelper.WriteError(w, http.StatusNotFound, "Product not found")
		return
	}

	h.logger.Info("product fetched", zap.String("id", p.ID), zap.String("slug", slug))
	httphelper.WriteSuccess(w, http.StatusOK, p)
}

// GetAll
// @Summary Получить список товаров
// @Tags Products
// @Produce json
// @Param type query string false "Тип товара (product или set)"
// @Param category query string false "Категория"
// @Param inStock query boolean false "Есть в наличии"
// @Param sort query string false "Сортировка (price или createdAt)"
// @Param limit query int false "Лимит на страницу"
// @Param offset query int false "Смещение"
// @Success 200 {object} httphelper.APIResponse
// @Failure 500 {object} httphelper.APIResponse
// @Router /products [get]
func (h *ProductHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()

	filter := repository.ProductFilter{
		Type:     q.Get("type"),
		Category: q.Get("category"),
		Sort:     q.Get("sort"),
		Limit:    httphelper.ParseInt(q.Get("limit"), 20),
		Offset:   httphelper.ParseInt(q.Get("offset"), 0),
	}

	if inStockStr := q.Get("inStock"); inStockStr != "" {
		b := inStockStr == "true"
		filter.InStock = &b
	}

	products, err := h.service.GetAll(filter)
	if err != nil {
		h.logger.Error("failed to get products", zap.Error(err))
		httphelper.WriteError(w, http.StatusInternalServerError, "Failed to load products")
		return
	}

	h.logger.Info("products fetched", zap.Int("count", len(products)))
	httphelper.WriteSuccess(w, http.StatusOK, products)
}

// GetSets
// @Summary Получить список наборов
// @Tags Sets
// @Produce json
// @Param inStock query boolean false "Есть в наличии"
// @Param sort query string false "Сортировка (price или createdAt)"
// @Param limit query int false "Лимит на страницу"
// @Param offset query int false "Смещение"
// @Success 200 {object} httphelper.APIResponse
// @Failure 500 {object} httphelper.APIResponse
// @Router /sets [get]
func (h *ProductHandler) GetSets(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()

	filter := repository.ProductFilter{
		Type:   "set",
		Sort:   q.Get("sort"),
		Limit:  httphelper.ParseInt(q.Get("limit"), 10),
		Offset: httphelper.ParseInt(q.Get("offset"), 0),
	}

	if inStockStr := q.Get("inStock"); inStockStr != "" {
		b := inStockStr == "true"
		filter.InStock = &b
	}

	sets, err := h.service.GetAll(filter)
	if err != nil {
		h.logger.Error("failed to get sets", zap.Error(err))
		httphelper.WriteError(w, http.StatusInternalServerError, "Failed to load sets")
		return
	}

	h.logger.Info("sets fetched", zap.Int("count", len(sets)))
	httphelper.WriteSuccess(w, http.StatusOK, sets)
}

// GetSetBySlug
// @Summary Получить набор по slug
// @Tags Sets
// @Produce json
// @Param slug path string true "Slug набора"
// @Success 200 {object} httphelper.APIResponse
// @Failure 404 {object} httphelper.APIResponse
// @Failure 400 {object} httphelper.APIResponse
// @Failure 500 {object} httphelper.APIResponse
// @Router /sets/{slug} [get]
func (h *ProductHandler) GetSetBySlug(w http.ResponseWriter, r *http.Request) {
	slug := chi.URLParam(r, "slug")

	p, err := h.service.GetBySlug(slug)
	if err != nil {
		h.logger.Error("set not found", zap.String("slug", slug), zap.Error(err))
		httphelper.WriteError(w, http.StatusNotFound, "Set not found")
		return
	}

	if p.Type != "set" {
		h.logger.Warn("not a set", zap.String("slug", slug))
		httphelper.WriteError(w, http.StatusBadRequest, "Item is not a set")
		return
	}

	h.logger.Info("set fetched", zap.String("slug", slug))
	httphelper.WriteSuccess(w, http.StatusOK, p)
}

// GetCategories
// @Summary Получить список категорий
// @Tags Products
// @Produce json
// @Success 200 {object} httphelper.APIResponse
// @Failure 500 {object} httphelper.APIResponse
// @Router /categories [get]
func (h *ProductHandler) GetCategories(w http.ResponseWriter, r *http.Request) {
	categories, err := h.service.GetCategories()
	if err != nil {
		h.logger.Error("failed to get categories", zap.Error(err))
		httphelper.WriteError(w, http.StatusInternalServerError, "Failed to load categories")
		return
	}

	h.logger.Info("categories fetched", zap.Int("count", len(categories)))
	httphelper.WriteSuccess(w, http.StatusOK, categories)
}

// Update
// @Summary Обновить товар по slug
// @Tags Products
// @Accept json
// @Produce json
// @Param slug path string true "Slug товара"
// @Param product body models.Product true "Обновленные данные товара"
// @Success 200 {object} httphelper.APIResponse
// @Failure 400 {object} httphelper.APIResponse
// @Failure 500 {object} httphelper.APIResponse
// @Router /products/{slug} [put]
func (h *ProductHandler) Update(w http.ResponseWriter, r *http.Request) {
	slug := chi.URLParam(r, "slug")
	var p models.Product

	if err := json.NewDecoder(r.Body).Decode(&p); err != nil {
		h.logger.Error("failed to decode product on update", zap.Error(err))
		httphelper.WriteError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	if err := validation.ValidateStruct(p); err != nil {
		h.logger.Error("validation failed on update", zap.Error(err))
		httphelper.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}

	if err := h.service.Update(slug, &p); err != nil {
		h.logger.Error("failed to update product", zap.String("slug", slug), zap.Error(err))
		httphelper.WriteError(w, http.StatusInternalServerError, "Failed to update product")
		return
	}

	h.logger.Info("product updated", zap.String("slug", slug))
	httphelper.WriteSuccess(w, http.StatusOK, p)
}

// Delete
// @Summary Удалить товар по slug
// @Tags Products
// @Produce json
// @Param slug path string true "Slug товара"
// @Success 204 "No Content"
// @Failure 500 {object} httphelper.APIResponse
// @Router /products/{slug} [delete]
func (h *ProductHandler) Delete(w http.ResponseWriter, r *http.Request) {
	slug := chi.URLParam(r, "slug")

	if err := h.service.Delete(slug); err != nil {
		h.logger.Error("failed to delete product", zap.String("slug", slug), zap.Error(err))
		httphelper.WriteError(w, http.StatusInternalServerError, "Failed to delete product")
		return
	}

	h.logger.Info("product deleted", zap.String("slug", slug))
	w.WriteHeader(http.StatusNoContent)
}
