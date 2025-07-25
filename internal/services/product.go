package services

import (
	"dozenChairs/internal/models"
	"dozenChairs/internal/repository"
	"time"
)

type ProductService interface {
	Create(p *models.Product) error
	GetBySlug(slug string) (*models.Product, error)
	GetAll(filter repository.ProductFilter) ([]*models.Product, error)
	GetCategories() ([]string, error)
	Update(slug string, p *models.Product) error
	Delete(slug string) error
}

type productService struct {
	repo repository.ProductRepository
}

func NewProductService(r repository.ProductRepository) ProductService {
	return &productService{repo: r}
}

func (s *productService) Create(p *models.Product) error {
	return s.repo.Create(p)
}

func (s *productService) GetBySlug(slug string) (*models.Product, error) {
	return s.repo.GetBySlug(slug)
}

func (s *productService) GetAll(filter repository.ProductFilter) ([]*models.Product, error) {
	return s.repo.GetAll(filter)
}

func (s *productService) GetCategories() ([]string, error) {
	return s.repo.GetCategories()
}

func (s *productService) Update(slug string, p *models.Product) error {
	p.UpdatedAt = time.Now().UTC()
	return s.repo.Update(slug, p)
}

func (s *productService) Delete(slug string) error {
	return s.repo.Delete(slug)
}
