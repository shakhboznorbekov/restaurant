package food_recipe_group_history

import (
	"github.com/uptrace/bun"
	"time"
)

type Filter struct {
	Limit  *int
	Offset *int
	FoodID *int
}

// @admin

type AdminGetListByFoodID struct {
	ID    int64   `json:"id"`
	From  *string `json:"from"`
	To    *string `json:"to"`
	Food  *string `json:"food"`
	Group *string `json:"group"`
}

type AdminGetDetail struct {
	ID      int64   `json:"id"`
	From    *string `json:"from"`
	To      *string `json:"to"`
	Food    *string `json:"food"`
	Group   *string `json:"group"`
	FoodID  *int64  `json:"food_id"`
	GroupID *int64  `json:"group_id"`
}

type AdminCreateRequest struct {
	Date    *string `json:"date" form:"date"`
	FoodID  *int64  `json:"food_id" form:"food_id"`
	GroupID *int64  `json:"group_id" form:"group_id"`
}

type AdminCreateResponse struct {
	bun.BaseModel `bun:"table:food_recipe_group_histories"`

	ID        int64      `json:"id" bun:"id,pk,autoincrement"`
	From      *time.Time `json:"from" bun:"from"`
	To        *time.Time `json:"to" bun:"to"`
	FoodID    *int64     `json:"food_id" bun:"food_id"`
	GroupID   *int64     `json:"group_id" bun:"group_id"`
	CreatedAt *time.Time `json:"created_at" bun:"created_at"`
	CreatedBy *int64     `json:"created_by" bun:"created_by"`
}

// @branch

type BranchGetListByFoodID struct {
	ID    int64   `json:"id"`
	From  *string `json:"from"`
	To    *string `json:"to"`
	Food  *string `json:"food"`
	Group *string `json:"group"`
}

type BranchGetDetail struct {
	ID      int64   `json:"id"`
	From    *string `json:"from"`
	To      *string `json:"to"`
	Food    *string `json:"food"`
	Group   *string `json:"group"`
	FoodID  *int64  `json:"food_id"`
	GroupID *int64  `json:"group_id"`
}

type BranchCreateRequest struct {
	Date    *string `json:"date" form:"date"`
	FoodID  *int64  `json:"food_id" form:"food_id"`
	GroupID *int64  `json:"group_id" form:"group_id"`
}

type BranchCreateResponse struct {
	bun.BaseModel `bun:"table:food_recipe_group_histories"`

	ID        int64      `json:"id" bun:"id,pk,autoincrement"`
	From      *time.Time `json:"from" bun:"from"`
	To        *time.Time `json:"to" bun:"to"`
	FoodID    *int64     `json:"food_id" bun:"food_id"`
	GroupID   *int64     `json:"group_id" bun:"group_id"`
	CreatedAt *time.Time `json:"created_at" bun:"created_at"`
	CreatedBy *int64     `json:"created_by" bun:"created_by"`
}
