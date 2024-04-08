package foodCategory

import (
	"github.com/restaurant/internal/pkg/utils"
	"mime/multipart"
	"time"

	"github.com/uptrace/bun"
)

type Filter struct {
	Limit  *int
	Offset *int
	Name   *string
	Page   *int
	Fields map[string][]string
	Joins  map[string]utils.Joins
}

// @admin

type AdminGetList struct {
	ID   int64   `json:"id"`
	Name *string `json:"name"`
	Logo *string `json:"logo"`
	Main *bool   `json:"main"`
}

type AdminGetDetail struct {
	ID   int64   `json:"id"`
	Name *string `json:"name"`
	Logo *string `json:"logo"`
	Main *bool   `json:"main"`
}

type AdminCreateRequest struct {
	Name     *string               `json:"name" form:"name"`
	Logo     *multipart.FileHeader `json:"logo" form:"logo"`
	LogoLink *string               `json:"-" form:"-"`
	Main     *bool                 `json:"main" form:"main"`
}

type AdminCreateResponse struct {
	bun.BaseModel `bun:"table:food_category"`

	ID           int64     `json:"id" bun:"id,pk,autoincrement"`
	Name         *string   `json:"name" bun:"name"`
	Logo         *string   `json:"-" bun:"logo"`
	Main         *bool     `json:"main" bun:"main"`
	RestaurantID *int64    `json:"restaurant_id" bun:"restaurant_id"`
	CreatedAt    time.Time `json:"-" bun:"created_at"`
	CreatedBy    int64     `json:"-" bun:"created_by"`
}

type AdminUpdateRequest struct {
	ID       int64                 `json:"id" form:"id"`
	Name     *string               `json:"name" form:"name"`
	Logo     *multipart.FileHeader `json:"-" form:"logo"`
	LogoLink *string               `json:"-" form:"-"`
	Main     *bool                 `json:"main" form:"main"`
}

// @branch

type BranchGetList struct {
	ID   int64   `json:"id"`
	Name *string `json:"name"`
	Logo *string `json:"logo"`
	Main *bool   `json:"main"`
}

type BranchGetDetail struct {
	ID   int64   `json:"id"`
	Name *string `json:"name"`
	Logo *string `json:"logo"`
	Main *bool   `json:"main"`
}

type BranchCreateRequest struct {
	Name     *string               `json:"name" form:"name"`
	Logo     *multipart.FileHeader `json:"logo" form:"logo"`
	LogoLink *string               `json:"-" form:"-"`
	Main     *bool                 `json:"main" form:"main"`
}

type BranchCreateResponse struct {
	bun.BaseModel `bun:"table:food_category"`

	ID           int64     `json:"id" bun:"id,pk,autoincrement"`
	Name         *string   `json:"name" bun:"name"`
	Logo         *string   `json:"-" bun:"logo"`
	Main         *bool     `json:"main" bun:"main"`
	RestaurantID *int64    `json:"restaurant_id" bun:"restaurant_id"`
	CreatedAt    time.Time `json:"-" bun:"created_at"`
	CreatedBy    int64     `json:"-" bun:"created_by"`
}

type BranchUpdateRequest struct {
	ID       int64                 `json:"id" form:"id"`
	Name     *string               `json:"name" form:"name"`
	Logo     *multipart.FileHeader `json:"-" form:"logo"`
	LogoLink *string               `json:"-" form:"-"`
	Main     *bool                 `json:"main" form:"main"`
}

// @cashier

type CashierGetList struct {
	ID   int64   `json:"id"`
	Name *string `json:"name"`
	Logo *string `json:"logo"`
	Main *bool   `json:"main"`
}

type CashierGetDetail struct {
	ID   int64   `json:"id"`
	Name *string `json:"name"`
	Logo *string `json:"logo"`
	Main *bool   `json:"main"`
}

type CashierCreateRequest struct {
	Name     *string               `json:"name" form:"name"`
	Logo     *multipart.FileHeader `json:"logo" form:"logo"`
	LogoLink *string               `json:"-" form:"-"`
	Main     *bool                 `json:"main" form:"main"`
}

type CashierCreateResponse struct {
	bun.BaseModel `bun:"table:food_category"`

	ID           int64     `json:"id" bun:"id,pk,autoincrement"`
	Name         *string   `json:"name" bun:"name"`
	Logo         *string   `json:"-" bun:"logo"`
	Main         *bool     `json:"main" bun:"main"`
	RestaurantID *int64    `json:"restaurant_id" bun:"restaurant_id"`
	CreatedAt    time.Time `json:"-" bun:"created_at"`
	CreatedBy    int64     `json:"-" bun:"created_by"`
}

type CashierUpdateRequest struct {
	ID       int64                 `json:"id" form:"id"`
	Name     *string               `json:"name" form:"name"`
	Logo     *multipart.FileHeader `json:"-" form:"logo"`
	LogoLink *string               `json:"-" form:"-"`
	Main     *bool                 `json:"main" form:"main"`
}

// @client

type ClientGetList struct {
	ID   int64   `json:"id"`
	Name *string `json:"name"`
	Logo *string `json:"logo"`
}

// @waiter

type WaiterGetList struct {
	Id   int64   `json:"id" bun:"id"`
	Name *string `json:"name" bun:"name"`
}
