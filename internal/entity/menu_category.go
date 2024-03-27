package entity

import (
	"github.com/uptrace/bun"
	"time"
)

type MenuCategory struct {
	bun.BaseModel `bun:"table:menu_categories"`

	ID           int64   `json:"id" bun:"id,pk,autoincrement"`
	Name         *string `json:"name" bun:"name"`
	RestaurantId *int64  `json:"restaurant_id" bun:"restaurant_id"`
	Logo         *string `json:"logo" bun:"logo"`

	CreatedAt *time.Time `json:"created_at" bun:"created_at"`
	CreatedBy *int64     `json:"created_by" bun:"created_by"`
	DeletedAt *time.Time `json:"deleted_at" bun:"deleted_at"`
	DeletedBy *int64     `json:"deleted_by" bun:"deleted_by"`
}
