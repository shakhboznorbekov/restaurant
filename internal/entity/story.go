package entity

import (
	"github.com/uptrace/bun"
	"time"
)

type Story struct {
	bun.BaseModel `bun:"table:stories"`

	ID           int64      `json:"id" bun:"id,pk,autoincrement"`
	Name         *string    `json:"name" bun:"name"`
	File         *string    `json:"file" bun:"file"`
	Type         *string    `json:"type" bun:"type"`
	Duration     *string    `json:"duration" bun:"duration"`
	ExpiredAt    *time.Time `json:"expired_at" bun:"expired_at"`
	RestaurantID *int64     `json:"restaurant_id" bun:"restaurant_id"`
	CreatedAt    *time.Time `json:"created_at" bun:"created_at"`
	CreatedBy    *int64     `json:"created_by" bun:"created_by"`
	DeletedAt    *time.Time `json:"deleted_at" bun:"deleted_at"`
	DeletedBy    *int64     `json:"deleted_by" bun:"deleted_by"`
}

//istoriya foylasj uchun
