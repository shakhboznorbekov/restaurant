package food_recipe

import (
	"github.com/uptrace/bun"
	"time"
)

type Filter struct {
	Limit  *int
	Offset *int
	Name   *string
	FoodId *int
}

// @admin

type AdminGetList struct {
	ID          int64    `json:"id" bun:"id"`
	Food        *string  `json:"food" bun:"food"`
	Product     *string  `json:"product" bun:"product"`
	Amount      *float64 `json:"amount" bun:"amount"`
	MeasureUnit *string  `json:"measure_unit" bun:"measure_unit"`
}

type AdminGetDetail struct {
	ID          int64    `json:"id"`
	Food        *string  `json:"food"`
	FoodId      *int64   `json:"food_id"`
	Amount      *float64 `json:"amount"`
	Product     *string  `json:"product"`
	ProductId   *int64   `json:"product_id"`
	MeasureUnit *string  `json:"measure_unit"`
	CreatedAt   *string  `json:"created_at"`
}

type AdminCreateRequest struct {
	Amount    float64 `json:"amount"`
	ProductId int64   `json:"product_id"`
	FoodId    int64   `json:"food_id"`
}

type AdminCreateResponse struct {
	bun.BaseModel `bun:"table:food_recipe"`

	ID        int64     `json:"id" bun:"id,pk,autoincrement"`
	Amount    float64   `json:"amount" bun:"amount"`
	ProductId int64     `json:"product_id" bun:"product_id"`
	FoodId    int64     `json:"food_id" bun:"food_id"`
	CreatedAt time.Time `json:"-" bun:"created_at"`
	CreatedBy int64     `json:"-" bun:"created_by"`
}

type AdminUpdateRequest struct {
	ID        int64    `json:"id" form:"id"`
	Amount    *float64 `json:"amount" form:"amount"`
	ProductId *int64   `json:"product_id" form:"product_id"`
	FoodId    *int64   `json:"food_id" form:"food_id"`
}

// @branch

type BranchGetList struct {
	ID          int64    `json:"id" bun:"id"`
	Food        *string  `json:"food" bun:"food"`
	Product     *string  `json:"product" bun:"product"`
	Amount      *float64 `json:"amount" bun:"amount"`
	MeasureUnit *string  `json:"measure_unit" bun:"measure_unit"`
}

type BranchGetDetail struct {
	ID          int64    `json:"id"`
	Food        *string  `json:"food"`
	FoodId      *int64   `json:"food_id"`
	Amount      *float64 `json:"amount"`
	Product     *string  `json:"product"`
	ProductId   *int64   `json:"product_id"`
	MeasureUnit *string  `json:"measure_unit"`
	CreatedAt   *string  `json:"created_at"`
}

type BranchCreateRequest struct {
	Amount    float64 `json:"amount"`
	ProductId int64   `json:"product_id"`
	FoodId    int64   `json:"food_id"`
}

type BranchCreateResponse struct {
	bun.BaseModel `bun:"table:food_recipe"`

	ID        int64     `json:"id" bun:"id,pk,autoincrement"`
	Amount    float64   `json:"amount" bun:"amount"`
	ProductId int64     `json:"product_id" bun:"product_id"`
	FoodId    int64     `json:"food_id" bun:"food_id"`
	CreatedAt time.Time `json:"-" bun:"created_at"`
	CreatedBy int64     `json:"-" bun:"created_by"`
}

type BranchUpdateRequest struct {
	ID        int64    `json:"id" form:"id"`
	Amount    *float64 `json:"amount" form:"amount"`
	ProductId *int64   `json:"product_id" form:"product_id"`
	FoodId    *int64   `json:"food_id" form:"food_id"`
}
