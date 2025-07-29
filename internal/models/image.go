package models

import "time"

type Image struct {
	ID        string    `json:"id"`
	ProductID string    `json:"product_id"`
	URL       string    `json:"url"`
	Filename  string    `json:"filename"`
	CreatedAt time.Time `json:"created_at"`
}
