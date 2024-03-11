package entity

import (
	"github.com/uptrace/bun"
	"time"
)

type Partner struct {
	bun.BaseModel `bun:"table:partners"`

	ID   int64   `json:"id" bun:"id,pk,autoincrement"`
	Name *string `json:"name" bun:"name"`
	Type *string `json:"type" bun:"type"`

	WarehouseID *int64     `json:"warehouse_id" bun:"warehouse_id"`
	CreatedAt   *time.Time `json:"created_at" bun:"created_at"`
	CreatedBy   *int64     `json:"created_by" bun:"created_by"`
	UpdatedAt   *time.Time `json:"updated_at" bun:"updated_at"`
	UpdatedBy   *int64     `json:"updated_by" bun:"updated_by"`
	DeletedAt   *time.Time `json:"deleted_at" bun:"deleted_at"`
	DeletedBy   *int64     `json:"deleted_by" bun:"deleted_by"`
}

//for ---->partniyorlar uchun
