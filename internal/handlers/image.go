package handlers

import (
	"dozenChairs/internal/metrics"
	"dozenChairs/internal/models"
	"dozenChairs/internal/services"
	"dozenChairs/pkg/httphelper"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"io"
	"net/http"
	"os"
	"path/filepath"
)

type ImageHandler struct {
	service services.ImageService
}

func NewImageHandler(s services.ImageService) *ImageHandler {
	return &ImageHandler{service: s}
}

// Upload
// @Summary Загрузить изображения для продукта
// @Tags Images
// @Accept multipart/form-data
// @Produce json
// @Param product_id formData string true "ID товара"
// @Param images formData file true "Файлы изображений (можно несколько)"
// @Success 201 {array} models.Image
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/v1//upload [post]
func (h *ImageHandler) Upload(w http.ResponseWriter, r *http.Request) {
	err := r.ParseMultipartForm(10 << 20) // 10 MB
	if err != nil {
		httphelper.WriteError(w, http.StatusBadRequest, "Invalid multipart form")
		return
	}

	productID := r.FormValue("product_id")
	files := r.MultipartForm.File["images"]
	if len(files) == 0 {
		httphelper.WriteError(w, http.StatusBadRequest, "No images provided")
		return
	}

	var uploaded []models.Image

	for _, fileHeader := range files {
		file, err := fileHeader.Open()
		if err != nil {
			continue // пропустим файл и логируем ошибку, если нужно
		}
		defer file.Close()

		filename := uuid.NewString() + filepath.Ext(fileHeader.Filename)
		dstPath := filepath.Join("uploads", filename)

		uploadDir := "uploads"

		// создаём папку, если не существует
		if _, err := os.Stat(uploadDir); os.IsNotExist(err) {
			err = os.MkdirAll(uploadDir, os.ModePerm)
			if err != nil {
				httphelper.WriteError(w, http.StatusInternalServerError, "Failed to create upload directory")
				return
			}
		}

		dst, err := os.Create(dstPath)
		if err != nil {
			continue
		}
		defer dst.Close()

		if _, err := io.Copy(dst, file); err != nil {
			continue
		}

		image := models.Image{
			ID:        uuid.NewString(),
			ProductID: productID,
			URL:       "/uploads/" + filename,
			Filename:  filename,
		}

		if err := h.service.SaveImage(r.Context(), &image); err != nil {
			continue
		}

		uploaded = append(uploaded, image)
	}

	httphelper.WriteSuccess(w, http.StatusCreated, uploaded)
	metrics.ImagesUploaded.Inc()
}

// Delete
// @Summary Удалить изображение по ID
// @Tags Images
// @Param id path string true "Image ID"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/v1//images/{id} [delete]
func (h *ImageHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		httphelper.WriteError(w, http.StatusBadRequest, "Image ID required")
		return
	}

	err := h.service.DeleteImage(r.Context(), id)
	if err != nil {
		httphelper.WriteError(w, http.StatusInternalServerError, "Failed to delete image")
		return
	}

	httphelper.WriteSuccess(w, http.StatusOK, map[string]string{"message": "Image deleted"})
}

// GetByProductID
// @Summary Получить изображения по товару
// @Tags Images
// @Param product_id path string true "ID или slug товара"
// @Success 200 {array} models.Image
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/v1//products/{product_id}/images [get]
func (h *ImageHandler) GetByProductID(w http.ResponseWriter, r *http.Request) {
	productID := chi.URLParam(r, "product_id")
	if productID == "" {
		httphelper.WriteError(w, http.StatusBadRequest, "Product ID required")
		return
	}

	images, err := h.service.GetImagesByProductID(r.Context(), productID)
	if err != nil {
		httphelper.WriteError(w, http.StatusInternalServerError, "Failed to get images")
		return
	}

	httphelper.WriteSuccess(w, http.StatusOK, images)
}
