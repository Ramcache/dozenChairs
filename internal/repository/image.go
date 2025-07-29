package repository

import (
	"context"
	"dozenChairs/internal/models"
	"github.com/jackc/pgx/v5/pgxpool"
	"time"
)

type ImageRepository interface {
	Save(ctx context.Context, img *models.Image) error
	Delete(ctx context.Context, id string) error
	GetByProductID(ctx context.Context, productID string) ([]models.Image, error)
}

type imageRepo struct {
	db *pgxpool.Pool
}

func NewImageRepo(db *pgxpool.Pool) ImageRepository {
	return &imageRepo{db: db}
}

func (r *imageRepo) Save(ctx context.Context, img *models.Image) error {
	query := `
		INSERT INTO images (id, product_id, url, filename, created_at)
		VALUES ($1, $2, $3, $4, $5)`
	_, err := r.db.Exec(ctx, query,
		img.ID,
		img.ProductID,
		img.URL,
		img.Filename,
		time.Now(),
	)
	return err
}

func (r *imageRepo) Delete(ctx context.Context, id string) error {
	_, err := r.db.Exec(ctx, `DELETE FROM images WHERE id = $1`, id)
	return err
}

func (r *imageRepo) GetByProductID(ctx context.Context, productID string) ([]models.Image, error) {
	query := `
		SELECT id, product_id, url, filename, created_at
		FROM images
		WHERE product_id = $1`

	rows, err := r.db.Query(ctx, query, productID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var images []models.Image
	for rows.Next() {
		var img models.Image
		if err := rows.Scan(&img.ID, &img.ProductID, &img.URL, &img.Filename, &img.CreatedAt); err != nil {
			return nil, err
		}
		images = append(images, img)
	}
	return images, nil
}
