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

// Create godoc
// @Summary      Создать товар
// @Description  Только для админов. Создаёт новый товар или набор. У набора нужно указать поле `includes`, а у обычного товара — `unitCount` и `attributes`.
// @Tags         Products
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Param        product  body      models.Product  true  "Данные нового товара или набора"
// @Success      201      {object}  models.Product
// @Failure      400      {object}  httphelper.APIResponse
// @Failure      500      {object}  httphelper.APIResponse
// @Router       /api/v1/products [post]
func (h *ProductHandler) Create(w http.ResponseWriter, r *http.Request) {
	if r.Header.Get("Content-Type") != "application/json" {
		httphelper.WriteError(w, http.StatusUnsupportedMediaType, "Content-Type must be application/json")
		return
	}

	var p models.Product
	if err := json.NewDecoder(r.Body).Decode(&p); err != nil {
		h.logger.Error("failed to decode product", zap.Error(err))
		httphelper.WriteError(w, http.StatusBadRequest, "Invalid JSON body")
		return
	}

	if err := validation.ValidateStruct(p); err != nil {
		h.logger.Warn("product validation failed", zap.Error(err))
		httphelper.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}

	if err := h.service.Create(&p); err != nil {
		h.logger.Error("product creation failed", zap.String("slug", p.Slug), zap.Error(err))
		httphelper.WriteError(w, http.StatusInternalServerError, "Failed to create product")
		return
	}

	h.logger.Info("product created", zap.String("slug", p.Slug))
	httphelper.WriteSuccess(w, http.StatusCreated, p)
}

// GetBySlug godoc
// @Summary      Получить товар по slug
// @Description  Возвращает один товар по его уникальному slug. Включает изображения, атрибуты и, при типе set — включённые товары.
// @Tags         Products
// @Produce      json
// @Param        slug  path      string  true  "Slug товара"
// @Success      200   {object}  models.Product
// @Failure      404   {object}  httphelper.APIResponse
// @Failure      500   {object}  httphelper.APIResponse
// @Router       /api/v1/products/{slug} [get]
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

// GetAll godoc
// @Summary      Получить список товаров
// @Description  Возвращает список товаров или наборов. Доступна фильтрация по типу, категории и наличию, а также сортировка по цене и дате создания.
// @Tags         Products
// @Produce      json
// @Param        type     query    string  false  "Тип товара (product или set)"
// @Param        category query    string  false  "Категория"
// @Param        inStock  query    boolean false  "Есть в наличии"
// @Param        sort     query    string  false  "Сортировка (price или createdAt)"
// @Param        limit    query    int     false  "Лимит на страницу (по умолчанию 20)"
// @Param        offset   query    int     false  "Смещение (по умолчанию 0)"
// @Success      200      {array}  models.Product
// @Failure      500      {object} httphelper.APIResponse
// @Router       /api/v1/products [get]
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

// GetSets godoc
// @Summary      Получить список наборов
// @Description  Возвращает все товары типа set. Наборы включают список вложенных товаров (`includes`).
// @Tags         Sets
// @Produce      json
// @Param        inStock query    boolean false  "Есть в наличии"
// @Param        sort     query    string  false  "Сортировка (price или createdAt)"
// @Param        limit    query    int     false  "Лимит на страницу"
// @Param        offset   query    int     false  "Смещение"
// @Success      200      {array}  models.Product
// @Failure      500      {object} httphelper.APIResponse
// @Router       /api/v1/sets [get]
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

// GetSetBySlug godoc
// @Summary      Получить набор по slug
// @Description  Возвращает товар типа set по его slug. Если товар не является набором, будет возвращена ошибка.
// @Tags         Sets
// @Produce      json
// @Param        slug  path      string  true  "Slug набора"
// @Success      200   {object}  models.Product
// @Failure      400   {object}  httphelper.APIResponse
// @Failure      404   {object}  httphelper.APIResponse
// @Failure      500   {object}  httphelper.APIResponse
// @Router       /api/v1/sets/{slug} [get]
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

// GetCategories godoc
// @Summary      Получить список категорий
// @Description  Возвращает все уникальные категории товаров
// @Tags         Products
// @Produce      json
// @Success      200  {array}  string
// @Failure      500  {object} httphelper.APIResponse
// @Router       /api/v1/categories [get]
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

// Update godoc
// @Summary      Обновить товар
// @Description  Только для админов. Обновляет данные товара по slug. Все поля можно изменить, включая изображения, состав набора и атрибуты.
// @Tags         Products
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Param        slug     path      string          true  "Slug товара"
// @Param        product  body      models.Product  true  "Обновлённые данные товара"
// @Success      200      {object}  models.Product
// @Failure      400      {object}  httphelper.APIResponse
// @Failure      500      {object}  httphelper.APIResponse
// @Router       /api/v1/products/{slug} [put]
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

// Delete godoc
// @Summary      Удалить товар
// @Description  Только для админов. Удаляет товар по slug. При удалении набора удаляется только сам набор, не включённые в него товары.
// @Tags         Products
// @Security     BearerAuth
// @Produce      json
// @Param        slug  path  string  true  "Slug товара"
// @Success      204   "No Content"
// @Failure      500   {object} httphelper.APIResponse
// @Router       /api/v1/products/{slug} [delete]
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
