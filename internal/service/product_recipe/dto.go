package product_recipe

import (
	"github.com/uptrace/bun"
	"time"
)

type Filter struct {
	Limit     *int
	Offset    *int
	Name      *string
	ProductId *int
}

// @admin

type AdminGetList struct {
	ID          int64    `json:"id" bun:"id"`
	Recipe      *string  `json:"recipe" bun:"recipe"`
	Product     *string  `json:"product" bun:"product"`
	Amount      *float64 `json:"amount" bun:"amount"`
	MeasureUnit *string  `json:"measure_unit" bun:"measure_unit"`
}

type AdminGetDetail struct {
	ID          int64    `json:"id"`
	Recipe      *string  `json:"recipe"`
	RecipeId    *int64   `json:"recipe_id"`
	Amount      *float64 `json:"amount"`
	Product     *string  `json:"product"`
	ProductId   *int64   `json:"product_id"`
	MeasureUnit *string  `json:"measure_unit"`
	CreatedAt   *string  `json:"created_at"`
}

type AdminCreateRequest struct {
	Amount    float64 `json:"amount"`
	ProductId int64   `json:"product_id"`
	RecipeId  int64   `json:"recipe_id"`
}

type AdminCreateResponse struct {
	bun.BaseModel `bun:"table:product_recipe"`

	ID        int64     `json:"id" bun:"id,pk,autoincrement"`
	Amount    float64   `json:"amount" bun:"amount"`
	ProductId int64     `json:"product_id" bun:"product_id"`
	RecipeId  int64     `json:"recipe_id" bun:"recipe_id"`
	CreatedAt time.Time `json:"-" bun:"created_at"`
	CreatedBy int64     `json:"-" bun:"created_by"`
}

type AdminUpdateRequest struct {
	ID        int64    `json:"id" form:"id"`
	Amount    *float64 `json:"amount" form:"amount"`
	ProductId *int64   `json:"product_id" form:"product_id"`
	RecipeId  *int64   `json:"recipe_id" form:"recipe_id"`
}

// @branch

type BranchGetList struct {
	ID          int64    `json:"id" bun:"id"`
	Recipe      *string  `json:"recipe" bun:"recipe"`
	Product     *string  `json:"product" bun:"product"`
	Amount      *float64 `json:"amount" bun:"amount"`
	MeasureUnit *string  `json:"measure_unit" bun:"measure_unit"`
}

type BranchGetDetail struct {
	ID          int64    `json:"id"`
	Recipe      *string  `json:"recipe"`
	RecipeId    *int64   `json:"recipe_id"`
	Amount      *float64 `json:"amount"`
	Product     *string  `json:"product"`
	ProductId   *int64   `json:"product_id"`
	MeasureUnit *string  `json:"measure_unit"`
	CreatedAt   *string  `json:"created_at"`
}

type BranchCreateRequest struct {
	Amount    float64 `json:"amount"`
	ProductId int64   `json:"product_id"`
	RecipeId  int64   `json:"recipe_id"`
}

type BranchCreateResponse struct {
	bun.BaseModel `bun:"table:product_recipe"`

	ID        int64     `json:"id" bun:"id,pk,autoincrement"`
	Amount    float64   `json:"amount" bun:"amount"`
	ProductId int64     `json:"product_id" bun:"product_id"`
	RecipeId  int64     `json:"recipe_id" bun:"recipe_id"`
	CreatedAt time.Time `json:"-" bun:"created_at"`
	CreatedBy int64     `json:"-" bun:"created_by"`
}

type BranchUpdateRequest struct {
	ID        int64    `json:"id" form:"id"`
	Amount    *float64 `json:"amount" form:"amount"`
	ProductId *int64   `json:"product_id" form:"product_id"`
	RecipeId  *int64   `json:"recipe_id" form:"recipe_id"`
}
