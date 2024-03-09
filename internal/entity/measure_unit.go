package entity

import (
	"github.com/uptrace/bun"
	"time"
)

type MeasureUnit struct {
	bun.BaseModel `bun:"table:measure_unit"`

	ID   int64   `json:"id" bun:"id,pk,autoincrement"`
	Name *string `json:"name" bun:"name"`

	CreatedAt *time.Time `json:"created_at" bun:"created_at"`
	CreatedBy *int64     `json:"created_by" bun:"created_by"`
	DeletedAt *time.Time `json:"deleted_at" bun:"deleted_at"`
	DeletedBy *int64     `json:"deleted_by" bun:"deleted_by"`
}
