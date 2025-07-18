package services

import (
	"context"
	"dozenChairs/internal/models"
	"dozenChairs/internal/repository"
	"time"

	"github.com/gosimple/slug"
)

type ChairService struct {
	chairRepo *repository.ChairRepository
}

func NewChairService(repo *repository.ChairRepository) *ChairService {
	return &ChairService{chairRepo: repo}
}

func (s *ChairService) Create(ctx context.Context, chair *models.Chair) error {
	chair.Slug = slug.Make(chair.Title) // автоматом генерируем слаг
	chair.Type = "product"
	chair.Category = "chair"
	now := time.Now().UTC()
	chair.CreatedAt = now
	chair.UpdatedAt = now
	return s.chairRepo.Create(ctx, chair)
}

func (s *ChairService) GetBySlug(ctx context.Context, slug string) (*models.Chair, error) {
	return s.chairRepo.GetBySlug(ctx, slug)
}

func (s *ChairService) GetAll(ctx context.Context) ([]*models.Chair, error) {
	return s.chairRepo.GetAll(ctx)
}

func (s *ChairService) UpdateBySlug(ctx context.Context, slug string, chair *models.Chair) error {
	chair.UpdatedAt = time.Now().UTC()
	return s.chairRepo.UpdateBySlug(ctx, slug, chair)
}

func (s *ChairService) DeleteBySlug(ctx context.Context, slug string) error {
	return s.chairRepo.DeleteBySlug(ctx, slug)
}
