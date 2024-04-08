package product_recipe_group

import (
	"github.com/lib/pq"
	"github.com/uptrace/bun"
	"time"
)

type Filter struct {
	Limit     *int
	Offset    *int
	Name      *string
	ProductID *int
}

// @admin

type AdminGetListByProductID struct {
	ID   int64   `json:"id"`
	Name *string `json:"name"`
}

type AdminGetDetail struct {
	ID        int64          `json:"id"`
	Name      *string        `json:"name"`
	ProductID *int64         `json:"product_id"`
	RecipeIds *pq.Int64Array `json:"recipe_ids"`
}

type AdminCreateRequest struct {
	Name      *string `json:"name" form:"name"`
	ProductID *int64  `json:"product_id" form:"product_id"`
	RecipeIds []int64 `json:"recipe_ids" form:"recipe_ids"`
}

type AdminCreateResponse struct {
	bun.BaseModel `bun:"table:product_recipe_groups"`

	ID        int64     `json:"id" bun:"id,pk,autoincrement"`
	Name      *string   `json:"name" bun:"name"`
	ProductID *int64    `json:"product_id" bun:"product_id"`
	RecipeIds []int64   `json:"recipe_ids" bun:"recipe_ids"`
	CreatedAt time.Time `json:"created_at" bun:"created_at"`
	CreatedBy int64     `json:"created_by" bun:"created_by"`
}

type AdminUpdateRequest struct {
	ID        int64   `json:"id" form:"id"`
	Name      *string `json:"name" form:"name"`
	ProductID *int64  `json:"product_id" form:"product_id"`
	RecipeIds []int64 `json:"recipe_ids" form:"recipe_ids"`
}

type AdminDeleteRecipeRequest struct {
	ID       int64  `json:"id" form:"id"`
	RecipeId *int64 `json:"recipe_id" form:"recipe_id"`
}

// @branch

type BranchGetListByProductID struct {
	ID   int64   `json:"id"`
	Name *string `json:"name"`
}

type BranchGetDetail struct {
	ID        int64   `json:"id"`
	Name      *string `json:"name"`
	ProductID *int64  `json:"product_id"`
	RecipeIds []int64 `json:"recipe_ids"`
}

type BranchCreateRequest struct {
	Name      *string `json:"name" form:"name"`
	ProductID *int64  `json:"product_id" form:"product_id"`
	RecipeIds []int64 `json:"recipe_ids" form:"recipe_ids"`
}

type BranchCreateResponse struct {
	bun.BaseModel `bun:"table:product_recipe_groups"`

	ID        int64     `json:"id" bun:"id,pk,autoincrement"`
	Name      *string   `json:"name" bun:"name"`
	ProductID *int64    `json:"product_id" bun:"product_id"`
	RecipeIds []int64   `json:"recipe_ids" bun:"recipe_ids"`
	CreatedAt time.Time `json:"created_at" bun:"created_at"`
	CreatedBy int64     `json:"created_by" bun:"created_by"`
}

type BranchUpdateRequest struct {
	ID        int64   `json:"id" form:"id"`
	Name      *string `json:"name" form:"name"`
	ProductID *int64  `json:"product_id" form:"product_id"`
	RecipeIds []int64 `json:"recipe_ids" form:"recipe_ids"`
}

type BranchDeleteRecipeRequest struct {
	ID       int64  `json:"id" form:"id"`
	RecipeId *int64 `json:"recipe_id" form:"recipe_id"`
}
