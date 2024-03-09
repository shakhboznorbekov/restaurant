package entity

import (
	"github.com/uptrace/bun"
	"time"
)

type CreditCard struct {
	bun.BaseModel `bun:"table:credit_cards"`

	ID        int64   `json:"id" bun:"id,pk,autoincrement"`
	Number    *string `json:"number" bun:"number"`
	ExpiresAt *string `json:"expires_at" bun:"number"`
	Phone     *string `json:"phone" bun:"phone"`
	IsMain    bool    `json:"is_main" bun:"is_main"`

	CreatedAt *time.Time `json:"created_at" bun:"created_at"`
	CreatedBy *int64     `json:"created_by" bun:"created_by"`
	UpdatedAt *time.Time `json:"updated_at" bun:"updated_at"`
	UpdatedBy *int64     `json:"updated_by" bun:"updated_by"`
	DeletedAt *time.Time `json:"deleted_at" bun:"deleted_at"`
	DeletedBy *int64     `json:"deleted_by" bun:"deleted_by"`
	UserID    *int64     `json:"user_id" bun:"user_id"`
}
