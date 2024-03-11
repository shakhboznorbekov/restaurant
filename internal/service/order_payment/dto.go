package order_payment

import (
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
	Status  *string `json:"status" bun:"status"`
	OrderID *int64  `json:"order_id" bun:"order_id"`
}

type ClientGetDetail struct {
	ID      int64   `json:"id"`
	Status  *string `json:"status"`
	OrderID *int64  `json:"order_id"`
}

type ClientCreateRequest struct {
	OrderID *int64 `json:"order_id" form:"order_id"`
}

type ClientCreateResponse struct {
	bun.BaseModel `bun:"table:order_payment"`

	ID      int64   `json:"id" bun:"id,pk,autoincrement"`
	Status  *string `json:"status" bun:"status"`
	OrderID *int64  `json:"order_id" bun:"order_id"`
	Price   float64 `json:"price" form:"price"`

	CreatedAt time.Time `json:"created_at" bun:"created_at"`
	CreatedBy int64     `json:"created_by" bun:"created_by"`
}

type ClientUpdateRequest struct {
	ID      int64   `json:"id" form:"id"`
	Status  *string `json:"status" bun:"status"`
	OrderID *int64  `json:"order_id" bun:"order_id"`
}

// @cashier
type CashierGetList struct {
	ID      int64   `json:"id" bun:"id"`
	Status  *string `json:"status" bun:"status"`
	OrderID *int64  `json:"order_id" bun:"order_id"`
}

type CashierGetDetail struct {
	ID      int64   `json:"id"`
	Status  *string `json:"status"`
	OrderID *int64  `json:"order_id"`
}

type CashierCreateRequest struct {
	OrderID *int64 `json:"order_id" form:"order_id"`
}

type CashierCreateResponse struct {
	bun.BaseModel `bun:"table:order_payment"`

	ID      int64    `json:"id" bun:"id,pk,autoincrement"`
	Status  *string  `json:"status" bun:"status"`
	OrderID *int64   `json:"order_id" bun:"order_id"`
	Price   *float64 `json:"price" form:"price"`

	CreatedAt time.Time `json:"created_at" bun:"created_at"`
	CreatedBy int64     `json:"created_by" bun:"created_by"`
}

type CashierUpdateRequest struct {
	ID      int64   `json:"id" form:"id"`
	Status  *string `json:"status" bun:"status"`
	OrderID *int64  `json:"order_id" bun:"order_id"`
}
