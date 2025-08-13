package repository

import (
	"context"
	"dozenChairs/internal/models"
	"encoding/json"
	"fmt"
	"github.com/jackc/pgx/v5/pgxpool"
	"strings"
	"time"
)

type ProductRepository interface {
	Create(p *models.Product) error
	GetBySlug(slug string) (*models.Product, error)
	GetAll(filter ProductFilter) ([]*models.Product, error)
	GetCategories() ([]string, error)
	Update(slug string, p *models.Product) error
	Delete(slug string) error
}

type productRepo struct {
	db *pgxpool.Pool
}

type ProductFilter struct {
	Type     string
	Category string
	InStock  *bool
	Sort     string
	Limit    int
	Offset   int
	FromDate time.Time
}

func NewProductRepo(db *pgxpool.Pool) ProductRepository {
	return &productRepo{db: db}
}

func (r *productRepo) Create(p *models.Product) error {
	attrJson, _ := json.Marshal(p.Attributes)
	includesJson, _ := json.Marshal(p.Includes)
	tagsJson, _ := json.Marshal(p.Tags)

	query := `
	INSERT INTO products (
		id, type, category, title, slug, description,
		price, old_price, in_stock, unit_count,
		attributes, includes, tags, created_at, updated_at
	) VALUES (
		$1, $2, $3, $4, $5, $6,
		$7, $8, $9, $10,
		$11, $12, $13, $14, $15
	)`

	_, err := r.db.Exec(
		context.Background(),
		query,
		p.ID, p.Type, p.Category, p.Title, p.Slug, p.Description,
		p.Price, p.OldPrice, p.InStock, p.UnitCount,
		string(attrJson), string(includesJson), string(tagsJson),
		p.CreatedAt, p.UpdatedAt,
	)

	return err
}

func (r *productRepo) GetBySlug(slug string) (*models.Product, error) {
	query := `SELECT id, type, category, title, slug, description, price, old_price,
	                 in_stock, unit_count, attributes, includes, tags, created_at, updated_at
	          FROM products WHERE slug = $1`

	var p models.Product
	var attributes, includes, tags []byte

	err := r.db.QueryRow(context.Background(), query, slug).Scan(
		&p.ID, &p.Type, &p.Category, &p.Title, &p.Slug, &p.Description,
		&p.Price, &p.OldPrice, &p.InStock, &p.UnitCount,
		&attributes, &includes, &tags,
		&p.CreatedAt, &p.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	// Распаковываем JSON-поля
	_ = json.Unmarshal(attributes, &p.Attributes)
	_ = json.Unmarshal(includes, &p.Includes)
	_ = json.Unmarshal(tags, &p.Tags)

	// Загружаем изображения
	imageQuery := `SELECT id, product_id, url, filename FROM images WHERE product_id = $1`

	rows, err := r.db.Query(context.Background(), imageQuery, p.ID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var images []models.Image
	for rows.Next() {
		var img models.Image
		if err := rows.Scan(&img.ID, &img.ProductID, &img.URL, &img.Filename); err != nil {
			return nil, err
		}
		images = append(images, img)
	}
	p.Images = images

	return &p, nil
}

func (r *productRepo) GetAll(f ProductFilter) ([]*models.Product, error) {
	query := `SELECT id, type, category, title, slug, description, price, old_price, in_stock, unit_count,
	                 attributes, includes, tags, created_at, updated_at
	          FROM products`
	args := []interface{}{}
	where := []string{}
	i := 1

	if f.Type != "" {
		where = append(where, fmt.Sprintf("type = $%d", i))
		args = append(args, f.Type)
		i++
	}
	if f.Category != "" {
		where = append(where, fmt.Sprintf("category = $%d", i))
		args = append(args, f.Category)
		i++
	}
	if f.InStock != nil {
		where = append(where, fmt.Sprintf("in_stock = $%d", i))
		args = append(args, *f.InStock)
		i++
	}
	if !f.FromDate.IsZero() {
		where = append(where, fmt.Sprintf("created_at >= $%d", i))
		args = append(args, f.FromDate)
		i++
	}

	if len(where) > 0 {
		query += " WHERE " + strings.Join(where, " AND ")
	}

	if f.Sort != "" {
		switch f.Sort {
		case "price", "createdAt":
			query += fmt.Sprintf(" ORDER BY %s", map[string]string{"price": "price", "createdAt": "created_at"}[f.Sort])
		}
	} else {
		query += " ORDER BY created_at DESC"
	}

	if f.Limit > 0 {
		query += fmt.Sprintf(" LIMIT $%d", i)
		args = append(args, f.Limit)
		i++
	}
	if f.Offset > 0 {
		query += fmt.Sprintf(" OFFSET $%d", i)
		args = append(args, f.Offset)
	}

	rows, err := r.db.Query(context.Background(), query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var products []*models.Product
	for rows.Next() {
		var p models.Product
		var attributes, includes, tags []byte

		err := rows.Scan(
			&p.ID, &p.Type, &p.Category, &p.Title, &p.Slug, &p.Description,
			&p.Price, &p.OldPrice, &p.InStock, &p.UnitCount,
			&attributes, &includes, &tags,
			&p.CreatedAt, &p.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}

		_ = json.Unmarshal(attributes, &p.Attributes)
		_ = json.Unmarshal(includes, &p.Includes)
		_ = json.Unmarshal(tags, &p.Tags)

		// загрузка изображений отдельно
		imgRows, err := r.db.Query(context.Background(), `SELECT id, product_id, url, filename FROM images WHERE product_id = $1`, p.ID)
		if err != nil {
			return nil, err
		}
		defer imgRows.Close()

		for imgRows.Next() {
			var img models.Image
			if err := imgRows.Scan(&img.ID, &img.ProductID, &img.URL, &img.Filename); err != nil {
				return nil, err
			}
			p.Images = append(p.Images, img)
		}

		products = append(products, &p)
	}

	return products, nil
}

func (r *productRepo) GetCategories() ([]string, error) {
	query := `SELECT DISTINCT category FROM products ORDER BY category`
	rows, err := r.db.Query(context.Background(), query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var categories []string
	for rows.Next() {
		var cat string
		if err := rows.Scan(&cat); err != nil {
			return nil, err
		}
		categories = append(categories, cat)
	}

	return categories, nil
}

func (r *productRepo) Update(slug string, p *models.Product) error {
	query := `
	UPDATE products SET
		id = $1,
		type = $2,
		category = $3,
		title = $4,
		description = $5,
		price = $6,
		old_price = $7,
		in_stock = $8,
		unit_count = $9,
		attributes = $10,
		includes = $11,
		tags = $12,
		updated_at = $13
	WHERE slug = $14
	`

	attrs, _ := json.Marshal(p.Attributes)
	includes, _ := json.Marshal(p.Includes)
	tags, _ := json.Marshal(p.Tags)

	_, err := r.db.Exec(context.Background(), query,
		p.ID, p.Type, p.Category, p.Title, p.Description,
		p.Price, p.OldPrice, p.InStock, p.UnitCount,
		attrs, includes, tags,
		p.UpdatedAt, slug,
	)

	return err
}

func (r *productRepo) Delete(slug string) error {
	_, err := r.db.Exec(context.Background(), `DELETE FROM products WHERE slug = $1`, slug)
	return err
}
