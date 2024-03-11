package entity

import (
	"github.com/uptrace/bun"
	"time"
)

type ProductRecipe struct {
	bun.BaseModel `bun:"table:product_recipe"`

	ID        int64      `json:"id" bun:"id,pk,autoincrement"`
	Amount    *float64   `json:"amount" bun:"amount"`
	ProductID *int64     `json:"product_id" bun:"product_id"`
	RecipeID  *int64     `json:"recipe_id" bun:"recipe_id"`
	CreatedAt *time.Time `json:"created_at" bun:"created_at"`
	CreatedBy *int64     `json:"created_by" bun:"created_by"`
	UpdatedAt *time.Time `json:"updated_at" bun:"updated_at"`
	UpdatedBy *int64     `json:"updated_by" bun:"updated_by"`
	DeletedAt *time.Time `json:"deleted_at" bun:"deleted_at"`
	DeletedBy *int64     `json:"deleted_by" bun:"deleted_by"`
}

//for----> bazi mahsulotlar bir nerchta mahsulotlardan tayyorlanadi shuning uchun
