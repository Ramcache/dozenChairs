package models

import (
	"time"
)

type ProductType string

const (
	TypeProduct ProductType = "product"
	TypeSet     ProductType = "set"
)

type Product struct {
	ID          string                 `json:"id" validate:"required"` // можно добавить `uuid4` при необходимости
	Type        ProductType            `json:"type" validate:"required,oneof=product set"`
	Category    string                 `json:"category" validate:"required"`
	Title       string                 `json:"title" validate:"required"`
	Slug        string                 `json:"slug" validate:"required"` // можно добавить custom slug-валидацию
	Description string                 `json:"description,omitempty"`
	Price       int                    `json:"price" validate:"gte=0"`
	OldPrice    *int                   `json:"oldPrice,omitempty" validate:"omitempty,gte=0"`
	InStock     bool                   `json:"inStock"`
	UnitCount   *int                   `json:"unitCount,omitempty" validate:"omitempty,gte=0"`
	Images      []Image                `json:"images"`
	Attributes  map[string]interface{} `json:"attributes,omitempty"`
	Includes    []IncludeItem          `json:"includes,omitempty" validate:"omitempty,dive"` // только для sets
	Tags        []string               `json:"tags,omitempty"`
	CreatedAt   time.Time              `json:"createdAt"`
	UpdatedAt   time.Time              `json:"updatedAt"`
}

type IncludeItem struct {
	ProductID string `json:"productId" validate:"required"`
	Quantity  int    `json:"quantity"  validate:"required,gt=0"`
}
