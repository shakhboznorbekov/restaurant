package entity

import (
	"github.com/uptrace/bun"
	"time"
)

type Banner struct {
	bun.BaseModel `bun:"table:banners"`

	ID           int64             `json:"id" bun:"id,pk,autoincrement"`
	Title        map[string]string `json:"title" bun:"title"`
	Description  map[string]string `json:"description" bun:"description"`
	Photo        *string           `json:"photo" bun:"photo"`
	Price        *float64          `json:"price" bun:"price"`
	OldPrice     *float64          `json:"old_price" bun:"old_price"`
	ExpiredAt    *time.Time        `json:"expired_at" bun:"expired_at"`
	RestaurantID *int64            `json:"restaurant_id" bun:"restaurant_id"`
	CreatedAt    *time.Time        `json:"created_at" bun:"created_at"`
	CreatedBy    *int64            `json:"created_by" bun:"created_by"`
	DeletedAt    *time.Time        `json:"deleted_at" bun:"deleted_at"`
	DeletedBy    *int64            `json:"deleted_by" bun:"deleted_by"`
}
