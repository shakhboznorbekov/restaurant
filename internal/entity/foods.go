package entity

import (
	"github.com/lib/pq"
	"github.com/uptrace/bun"
	"time"
)

type Foods struct {
	bun.BaseModel `bun:"table:foods"`

	ID           int64           `json:"id" bun:"id,pk,autoincrement"`
	Name         *string         `json:"name" bun:"name"`
	Photos       *pq.StringArray `json:"photos" bun:"photos"`
	CategoryID   *int64          `json:"category_id" bun:"category_id"`
	BranchID     *int64          `json:"branch_id" bun:"branch_id"`
	Price        *float64        `json:"price" bun:"price"`
	RestaurantID *int64          `json:"restaurant_id" bun:"restaurant_id"`
	CreatedAt    *time.Time      `json:"created_at" bun:"created_at"`
	CreatedBy    *int64          `json:"created_by" bun:"created_by"`
	UpdatedAt    *time.Time      `json:"updated_at" bun:"updated_at"`
	UpdatedBy    *int64          `json:"updated_by" bun:"updated_by"`
	DeletedAt    *time.Time      `json:"deleted_at" bun:"deleted_at"`
	DeletedBy    *int64          `json:"deleted_by" bun:"deleted_by"`
}
