package repository

import (
	"context"
	"dozenChairs/internal/models"
	"encoding/json"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type ChairRepository struct {
	db *pgxpool.Pool
}

func NewChairRepository(db *pgxpool.Pool) *ChairRepository {
	return &ChairRepository{db: db}
}

func (r *ChairRepository) Create(ctx context.Context, chair *models.Chair) error {
	imagesJSON, _ := json.Marshal(chair.Images)
	attributesJSON, _ := json.Marshal(chair.Attributes)
	tagsJSON, _ := json.Marshal(chair.Tags)

	query := `
		INSERT INTO chairs (id, type, category, title, slug, description, price, old_price, in_stock, unit_count, images, attributes, tags, created_at, updated_at)
		VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14,$15)
	`
	_, err := r.db.Exec(ctx, query,
		chair.ID,
		chair.Type,
		chair.Category,
		chair.Title,
		chair.Slug,
		chair.Description,
		chair.Price,
		chair.OldPrice,
		chair.InStock,
		chair.UnitCount,
		imagesJSON,
		attributesJSON,
		tagsJSON,
		chair.CreatedAt,
		chair.UpdatedAt,
	)
	return err
}

func (r *ChairRepository) GetBySlug(ctx context.Context, slug string) (*models.Chair, error) {
	query := `SELECT id, type, category, title, slug, description, price, old_price, in_stock, unit_count, images, attributes, tags, created_at, updated_at FROM chairs WHERE slug = $1`
	row := r.db.QueryRow(ctx, query, slug)

	var chair models.Chair
	var imagesJSON, attributesJSON, tagsJSON []byte
	err := row.Scan(
		&chair.ID,
		&chair.Type,
		&chair.Category,
		&chair.Title,
		&chair.Slug,
		&chair.Description,
		&chair.Price,
		&chair.OldPrice,
		&chair.InStock,
		&chair.UnitCount,
		&imagesJSON,
		&attributesJSON,
		&tagsJSON,
		&chair.CreatedAt,
		&chair.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	_ = json.Unmarshal(imagesJSON, &chair.Images)
	_ = json.Unmarshal(attributesJSON, &chair.Attributes)
	_ = json.Unmarshal(tagsJSON, &chair.Tags)
	return &chair, nil
}

func (r *ChairRepository) GetAll(ctx context.Context) ([]*models.Chair, error) {
	query := `SELECT id, type, category, title, slug, description, price, old_price, in_stock, unit_count, images, attributes, tags, created_at, updated_at FROM chairs`
	rows, err := r.db.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var chairs []*models.Chair
	for rows.Next() {
		var chair models.Chair
		var imagesJSON, attributesJSON, tagsJSON []byte
		err := rows.Scan(
			&chair.ID,
			&chair.Type,
			&chair.Category,
			&chair.Title,
			&chair.Slug,
			&chair.Description,
			&chair.Price,
			&chair.OldPrice,
			&chair.InStock,
			&chair.UnitCount,
			&imagesJSON,
			&attributesJSON,
			&tagsJSON,
			&chair.CreatedAt,
			&chair.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		_ = json.Unmarshal(imagesJSON, &chair.Images)
		_ = json.Unmarshal(attributesJSON, &chair.Attributes)
		_ = json.Unmarshal(tagsJSON, &chair.Tags)
		chairs = append(chairs, &chair)
	}
	return chairs, nil
}

func (r *ChairRepository) UpdateBySlug(ctx context.Context, slug string, chair *models.Chair) error {
	imagesJSON, _ := json.Marshal(chair.Images)
	attributesJSON, _ := json.Marshal(chair.Attributes)
	tagsJSON, _ := json.Marshal(chair.Tags)

	query := `
		UPDATE chairs SET
			title = $1,
			description = $2,
			price = $3,
			old_price = $4,
			in_stock = $5,
			unit_count = $6,
			images = $7,
			attributes = $8,
			tags = $9,
			updated_at = $10
		WHERE slug = $11
	`
	_, err := r.db.Exec(ctx, query,
		chair.Title,
		chair.Description,
		chair.Price,
		chair.OldPrice,
		chair.InStock,
		chair.UnitCount,
		imagesJSON,
		attributesJSON,
		tagsJSON,
		time.Now().UTC(),
		slug,
	)
	return err
}

func (r *ChairRepository) DeleteBySlug(ctx context.Context, slug string) error {
	query := `DELETE FROM chairs WHERE slug = $1`
	_, err := r.db.Exec(ctx, query, slug)
	return err
}
