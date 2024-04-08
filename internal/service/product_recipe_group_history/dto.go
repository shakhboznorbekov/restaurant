package product_recipe_group_history

import (
	"github.com/uptrace/bun"
	"time"
)

type Filter struct {
	Limit     *int
	Offset    *int
	ProductID *int
}

// @admin

type AdminGetListByProductID struct {
	ID      int64   `json:"id"`
	From    *string `json:"from"`
	To      *string `json:"to"`
	Product *string `json:"product"`
	Group   *string `json:"group"`
}

type AdminGetDetail struct {
	ID        int64   `json:"id"`
	From      *string `json:"from"`
	To        *string `json:"to"`
	Product   *string `json:"product"`
	Group     *string `json:"group"`
	ProductID *int64  `json:"product_id"`
	GroupID   *int64  `json:"group_id"`
}

type AdminCreateRequest struct {
	Date      *string `json:"date" form:"date"`
	ProductID *int64  `json:"product_id" form:"product_id"`
	GroupID   *int64  `json:"group_id" form:"group_id"`
}

type AdminCreateResponse struct {
	bun.BaseModel `bun:"table:product_recipe_group_histories"`

	ID        int64      `json:"id" bun:"id,pk,autoincrement"`
	From      *time.Time `json:"from" bun:"from"`
	To        *time.Time `json:"to" bun:"to"`
	ProductID *int64     `json:"product_id" bun:"product_id"`
	GroupID   *int64     `json:"group_id" bun:"group_id"`
	CreatedAt *time.Time `json:"created_at" bun:"created_at"`
	CreatedBy *int64     `json:"created_by" bun:"created_by"`
}

// @branch

type BranchGetListByProductID struct {
	ID      int64   `json:"id"`
	From    *string `json:"from"`
	To      *string `json:"to"`
	Product *string `json:"product"`
	Group   *string `json:"group"`
}

type BranchGetDetail struct {
	ID        int64   `json:"id"`
	From      *string `json:"from"`
	To        *string `json:"to"`
	Product   *string `json:"product"`
	Group     *string `json:"group"`
	ProductID *int64  `json:"product_id"`
	GroupID   *int64  `json:"group_id"`
}

type BranchCreateRequest struct {
	Date      *string `json:"date" form:"date"`
	ProductID *int64  `json:"product_id" form:"product_id"`
	GroupID   *int64  `json:"group_id" form:"group_id"`
}

type BranchCreateResponse struct {
	bun.BaseModel `bun:"table:product_recipe_group_histories"`

	ID        int64      `json:"id" bun:"id,pk,autoincrement"`
	From      *time.Time `json:"from" bun:"from"`
	To        *time.Time `json:"to" bun:"to"`
	ProductID *int64     `json:"product_id" bun:"product_id"`
	GroupID   *int64     `json:"group_id" bun:"group_id"`
	CreatedAt *time.Time `json:"created_at" bun:"created_at"`
	CreatedBy *int64     `json:"created_by" bun:"created_by"`
}
