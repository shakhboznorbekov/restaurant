package entity

import (
	"github.com/uptrace/bun"
	"time"
)

type Restaurant struct {
	bun.BaseModel `bun:"table:restaurants"`

	ID         int64   `json:"id" bun:"id,pk,autoincrement"`
	Name       *string `json:"name" bun:"name"`
	Logo       *string `json:"logo" bun:"logo"`
	MiniLogo   *string `json:"mini_logo" bun:"mini_logo"`
	WebsiteUrl *string `json:"website_url" bun:"website_url"`

	CreatedAt *time.Time `json:"created_at" bun:"created_at"`
	CreatedBy *int64     `json:"created_by" bun:"created_by"`
	UpdatedAt *time.Time `json:"updated_at" bun:"updated_at"`
	UpdatedBy *int64     `json:"updated_by" bun:"updated_by"`
	DeletedAt *time.Time `json:"deleted_at" bun:"deleted_at"`
	DeletedBy *int64     `json:"deleted_by" bun:"deleted_by"`
}
