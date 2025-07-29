package services

import (
	"context"
	"dozenChairs/internal/models"
	"dozenChairs/internal/repository"
)

type ImageService interface {
	SaveImage(ctx context.Context, img *models.Image) error
	DeleteImage(ctx context.Context, id string) error
	GetImagesByProductID(ctx context.Context, productID string) ([]models.Image, error)
}

type imageService struct {
	repo repository.ImageRepository
}

func NewImageService(r repository.ImageRepository) ImageService {
	return &imageService{repo: r}
}

func (s *imageService) SaveImage(ctx context.Context, img *models.Image) error {
	return s.repo.Save(ctx, img)
}

func (s *imageService) DeleteImage(ctx context.Context, id string) error {
	return s.repo.Delete(ctx, id)
}

func (s *imageService) GetImagesByProductID(ctx context.Context, productID string) ([]models.Image, error) {
	return s.repo.GetByProductID(ctx, productID)
}
