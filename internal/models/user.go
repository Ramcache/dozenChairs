package models

import "time"

type User struct {
	ID            int       `json:"id"`
	Username      string    `json:"username"`
	PasswordHash  string    `json:"-"`
	Email         string    `json:"email"`
	FullName      string    `json:"full_name"`
	Phone         string    `json:"phone,omitempty"`
	Address       string    `json:"address,omitempty"`
	Role          string    `json:"role"`
	CreatedAt     time.Time `json:"created_at"`
	EmailVerified bool      `json:"email_verified"`
}
