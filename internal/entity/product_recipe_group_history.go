package entity

import (
	"github.com/uptrace/bun"
	"time"
)

type ProductRecipeGroupHistory struct {
	bun.BaseModel `bun:"table:product_recipe_group_histories"`

	ID        int64      `json:"id" bun:"id,pk,autoincrement"`
	From      *time.Time `json:"from" bun:"from"`
	To        *time.Time `json:"to" bun:"to"`
	ProductID *int64     `json:"product_id" bun:"product_id"`
	GroupID   *int64     `json:"group_id" bun:"group_id"`
	CreatedAt *time.Time `json:"created_at" bun:"created_at"`
	CreatedBy *int64     `json:"created_by" bun:"created_by"`
	UpdatedAt *time.Time `json:"updated_at" bun:"updated_at"`
	UpdatedBy *int64     `json:"updated_by" bun:"updated_by"`
	DeletedAt *time.Time `json:"deleted_at" bun:"deleted_at"`
	DeletedBy *int64     `json:"deleted_by" bun:"deleted_by"`
}
