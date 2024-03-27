package entity

import (
	"github.com/uptrace/bun"
	"time"
)

type FoodRecipeGroup struct {
	bun.BaseModel `bun:"table:food_recipe_groups"`

	ID        int64      `json:"id" bun:"id,pk,autoincrement"`
	Name      *string    `json:"name" bun:"name"`
	FoodID    *int64     `json:"food_id" bun:"food_id"`
	RecipeIds []int64    `json:"recipe_ids" bun:"recipe_ids"`
	CreatedAt *time.Time `json:"created_at" bun:"created_at"`
	CreatedBy *int64     `json:"created_by" bun:"created_by"`
	UpdatedAt *time.Time `json:"updated_at" bun:"updated_at"`
	UpdatedBy *int64     `json:"updated_by" bun:"updated_by"`
	DeletedAt *time.Time `json:"deleted_at" bun:"deleted_at"`
	DeletedBy *int64     `json:"deleted_by" bun:"deleted_by"`
}
