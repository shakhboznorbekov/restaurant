package entity

import (
	"github.com/uptrace/bun"
	"time"
)

type ServicePercentage struct {
	bun.BaseModel `bun:"table:service-percentage"`

	ID        int64      `json:"id" bun:"id,pk,autoincrement"`
	Percent   *float64   `json:"percent" bun:"percent"`
	CreatedAt *time.Time `json:"created_at" bun:"created_at"`
	CreatedBy *int64     `json:"created_by" bun:"created_by"`
	DeletedAt *time.Time `json:"deleted_at" bun:"deleted_at"`
	DeletedBy *int64     `json:"deleted_by" bun:"deleted_by"`
}
