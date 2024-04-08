package entity

import (
	"github.com/uptrace/bun"
	"time"
)

type Product struct {
	bun.BaseModel `bun:"table:products"`

	ID      int64   `json:"id" bun:"id,pk,autoincrement"`
	Name    *string `json:"name" bun:"name"`
	Barcode *string `json:"barcode" bun:"barcode"`

	MeasureUnitID *int64     `json:"measure_unit_id" bun:"measure_unit_id"`
	CreatedAt     *time.Time `json:"created_at" bun:"created_at"`
	CreatedBy     *int64     `json:"created_by" bun:"created_by"`
	UpdatedAt     *time.Time `json:"updated_at" bun:"updated_at"`
	UpdatedBy     *int64     `json:"updated_by" bun:"updated_by"`
	DeletedAt     *time.Time `json:"deleted_at" bun:"deleted_at"`
	DeletedBy     *int64     `json:"deleted_by" bun:"deleted_by"`
	RestaurantID  *int64     `json:"restaurant_id" bun:"restaurant_id"`
}
