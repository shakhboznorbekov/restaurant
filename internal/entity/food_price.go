package entity

import (
	"github.com/uptrace/bun"
	"time"
)

type FoodPrice struct {
	bun.BaseModel `bun:"table:food_price"`

	ID      int64      `json:"id" bun:"id,pk,autoincrement"`
	Price   *float64   `json:"price" bun:"price"`
	SetDate *time.Time `json:"set_date" bun:"set_date"`

	FoodID    *int64     `json:"food_id" bun:"food_id"`
	CreatedAt *time.Time `json:"created_at" bun:"created_at"`
	CreatedBy *int64     `json:"created_by" bun:"created_by"`
	UpdatedAt *time.Time `json:"updated_at" bun:"updated_at"`
	UpdatedBy *int64     `json:"updated_by" bun:"updated_by"`
	DeletedAt *time.Time `json:"deleted_at" bun:"deleted_at"`
	DeletedBy *int64     `json:"deleted_by" bun:"deleted_by"`
}

//for ---->  taomlarni narxlari branchlar uchun
