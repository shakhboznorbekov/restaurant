package food_price

import (
	"github.com/uptrace/bun"
	"time"
)

type Filter struct {
	Limit  *int
	Offset *int
}

// @admin

type AdminGetList struct {
	ID      int64    `json:"id"`
	Price   *float64 `json:"price"`
	SetDate *string  `json:"set_date"`
	FoodID  *int64   `json:"food_id"`
}

type AdminGetDetail struct {
	ID      int64    `json:"id"`
	Price   *float64 `json:"price"`
	SetDate *string  `json:"set_date"`
	FoodID  *int64   `json:"food_id"`
}

type AdminCreateRequest struct {
	Price   *float64 `json:"price" form:"price"`
	SetDate *string  `json:"set_date" form:"set_date"`
	FoodID  *int64   `json:"food_id" form:"food_id"`
}

type AdminCreateResponse struct {
	bun.BaseModel `bun:"table:food_price"`

	ID        int64      `json:"id" bun:"id,pk,autoincrement"`
	Price     *float64   `json:"price" form:"price"`
	SetDate   *time.Time `json:"set_date" form:"set_date"`
	FoodID    *int64     `json:"food_id" form:"food_id"`
	CreatedAt time.Time  `json:"-" bun:"created_at"`
	CreatedBy int64      `json:"-" bun:"created_by"`
}

type AdminUpdateRequest struct {
	ID      int64    `json:"id" form:"id"`
	Price   *float64 `json:"price" form:"price"`
	SetDate *string  `json:"set_date" form:"set_date"`
	FoodID  *int64   `json:"food_id" form:"food_id"`
}
