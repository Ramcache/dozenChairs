package models

import "time"

type Chair struct {
	ID          string          `json:"id"`
	Type        string          `json:"type"`     // всегда "product"
	Category    string          `json:"category"` // всегда "chair"
	Title       string          `json:"title"`
	Slug        string          `json:"slug"`
	Description string          `json:"description"`
	Price       int             `json:"price"`
	OldPrice    *int            `json:"oldPrice,omitempty"`
	InStock     bool            `json:"inStock"`
	UnitCount   *int            `json:"unitCount,omitempty"`
	Images      []string        `json:"images"`
	Attributes  ChairAttributes `json:"attributes"`
	Tags        []string        `json:"tags"`
	CreatedAt   time.Time       `json:"createdAt"`
	UpdatedAt   time.Time       `json:"updatedAt"`
}

type ChairAttributes struct {
	Color          string `json:"color"`
	Material       string `json:"material"`
	MaterialPillow string `json:"materialPillow"`
	MaterialFrame  string `json:"materialFrame"`
	ColorPillow    string `json:"colorPillow"`
	ColorFrame     string `json:"colorFrame"`
	TotalHeight    int    `json:"totalHeight"`
	Width          int    `json:"width"`
}
