package repository

import (
	"context"
	"dozenChairs/internal/models"
	"encoding/json"
	"fmt"
	"github.com/jackc/pgx/v5"
	"strings"
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
	db *pgx.Conn
}

type ProductFilter struct {
	Type     string
	Category string
	InStock  *bool
	Sort     string
	Limit    int
	Offset   int
}

func NewProductRepo(db *pgx.Conn) ProductRepository {
	return &productRepo{db: db}
}

func (r *productRepo) Create(p *models.Product) error {
	attrJson, _ := json.Marshal(p.Attributes)
	imagesJson, _ := json.Marshal(p.Images)
	includesJson, _ := json.Marshal(p.Includes)
	tagsJson, _ := json.Marshal(p.Tags)

	query := `
	INSERT INTO products (
		id, type, category, title, slug, description,
		price, old_price, in_stock, unit_count,
		images, attributes, includes, tags, created_at, updated_at
	) VALUES (
		$1, $2, $3, $4, $5, $6,
		$7, $8, $9, $10,
		$11, $12, $13, $14, $15, $16
	)`

	_, err := r.db.Exec(
		context.Background(),
		query,
		p.ID, p.Type, p.Category, p.Title, p.Slug, p.Description,
		p.Price, p.OldPrice, p.InStock, p.UnitCount,
		string(imagesJson), string(attrJson), string(includesJson), string(tagsJson),
		p.CreatedAt, p.UpdatedAt,
	)

	return err
}

func (r *productRepo) GetBySlug(slug string) (*models.Product, error) {
	query := `SELECT id, type, category, title, slug, description, price, old_price, in_stock, unit_count,
	                 images, attributes, includes, tags, created_at, updated_at
	          FROM products WHERE slug = $1`

	var p models.Product
	var images, attributes, includes, tags []byte

	err := r.db.QueryRow(context.Background(), query, slug).Scan(
		&p.ID, &p.Type, &p.Category, &p.Title, &p.Slug, &p.Description,
		&p.Price, &p.OldPrice, &p.InStock, &p.UnitCount,
		&images, &attributes, &includes, &tags,
		&p.CreatedAt, &p.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	// Разбираем JSON-поля
	_ = json.Unmarshal(images, &p.Images)
	_ = json.Unmarshal(attributes, &p.Attributes)
	_ = json.Unmarshal(includes, &p.Includes)
	_ = json.Unmarshal(tags, &p.Tags)

	return &p, nil
}

func (r *productRepo) GetAll(f ProductFilter) ([]*models.Product, error) {
	query := `SELECT id, type, category, title, slug, description, price, old_price, in_stock, unit_count,
	                 images, attributes, includes, tags, created_at, updated_at
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
		var images, attributes, includes, tags []byte

		err := rows.Scan(
			&p.ID, &p.Type, &p.Category, &p.Title, &p.Slug, &p.Description,
			&p.Price, &p.OldPrice, &p.InStock, &p.UnitCount,
			&images, &attributes, &includes, &tags,
			&p.CreatedAt, &p.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}

		_ = json.Unmarshal(images, &p.Images)
		_ = json.Unmarshal(attributes, &p.Attributes)
		_ = json.Unmarshal(includes, &p.Includes)
		_ = json.Unmarshal(tags, &p.Tags)

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
		images = $10,
		attributes = $11,
		includes = $12,
		tags = $13,
		updated_at = $14
	WHERE slug = $15
	`

	images, _ := json.Marshal(p.Images)
	attrs, _ := json.Marshal(p.Attributes)
	includes, _ := json.Marshal(p.Includes)
	tags, _ := json.Marshal(p.Tags)

	_, err := r.db.Exec(context.Background(), query,
		p.ID, p.Type, p.Category, p.Title, p.Description,
		p.Price, p.OldPrice, p.InStock, p.UnitCount,
		images, attrs, includes, tags,
		p.UpdatedAt, slug,
	)

	return err
}

func (r *productRepo) Delete(slug string) error {
	_, err := r.db.Exec(context.Background(), `DELETE FROM products WHERE slug = $1`, slug)
	return err
}
