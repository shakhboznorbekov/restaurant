package entity

import (
	"github.com/uptrace/bun"
	"time"
)

type Region struct {
	bun.BaseModel `bun:"table:regions"`

	ID   int64              `json:"id" bun:"id,pk,autoincrement"`
	Name *map[string]string `json:"name" bun:"name"`
	Code *int               `json:"code" bun:"code"`

	CreatedAt *time.Time `json:"created_at" bun:"created_at"`
	CreatedBy *int64     `json:"created_by" bun:"created_by"`
	DeletedAt *time.Time `json:"deleted_at" bun:"deleted_at"`
	DeletedBy *int64     `json:"deleted_by" bun:"deleted_by"`
}
