package entity

import (
	"github.com/uptrace/bun"
	"time"
)

type OrderPayment struct {
	bun.BaseModel `bun:"table:order_payment"`

	ID     int64   `json:"id" bun:"id,pk,autoincrement"`
	Status *string `json:"status" bun:"status"`

	OrderID   *int64     `json:"order_id" bun:"order_id"`
	CreatedAt *time.Time `json:"created_at" bun:"created_at"`
	CreatedBy *int64     `json:"created_by" bun:"created_by"`
	UpdatedAt *time.Time `json:"updated_at" bun:"updated_at"`
	UpdatedBy *int64     `json:"updated_by" bun:"updated_by"`
	DeletedAt *time.Time `json:"deleted_at" bun:"deleted_at"`
	DeletedBy *int64     `json:"deleted_by" bun:"deleted_by"`
}
