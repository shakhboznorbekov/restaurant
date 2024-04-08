package order_menu

import (
	"github.com/lib/pq"
	"github.com/uptrace/bun"
	"time"
)

type Filter struct {
	Limit  *int
	Offset *int
}

// @client

type ClientGetList struct {
	ID      int64   `json:"id" bun:"id"`
	Count   *int    `json:"count" bun:"count"`
	MenuID  *int64  `json:"menu_id" bun:"menu_id"`
	OrderID *int64  `json:"order_id" bun:"order_id"`
	Name    *string `json:"name" bun:"name"`
}

type ClientGetDetail struct {
	ID      int64  `json:"id"`
	Count   *int   `json:"count"`
	MenuID  *int64 `json:"menu_id"`
	OrderID *int64 `json:"order_id"`
}

type ClientCreateRequest struct {
	Count   *int   `json:"count" form:"count"`
	MenuID  *int64 `json:"menu_id" form:"menu_id"`
	OrderID *int64 `json:"order_id" form:"order_id"`
}

type ClientGetOftenList struct {
	ID     int64           `json:"id"`
	Name   *string         `json:"name"`
	Photos *pq.StringArray `json:"photos"`
	Price  *float32        `json:"price"`
}

type ClientCreateResponse struct {
	bun.BaseModel `bun:"table:order_menu"`

	ID      int64  `json:"id" bun:"id,pk,autoincrement"`
	Count   *int   `json:"count" bun:"count"`
	MenuID  *int64 `json:"menu_id" bun:"menu_id"`
	OrderID *int64 `json:"order_id" bun:"order_id"`

	CreatedAt time.Time `json:"created_at" bun:"created_at"`
	CreatedBy int64     `json:"created_by" bun:"created_by"`
}

type ClientUpdateRequest struct {
	ID      int64  `json:"id" form:"id"`
	Count   *int   `json:"count" form:"count"`
	MenuID  *int64 `json:"menu_id" form:"menu_id"`
	OrderID *int64 `json:"order_id" form:"order_id"`
}
