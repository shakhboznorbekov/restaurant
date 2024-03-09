package entity

import (
	"github.com/uptrace/bun"
	"time"
)

type Warehouse struct {
	bun.BaseModel `bun:"table:warehouses"`

	ID       int64               `json:"id" bun:"id,pk,autoincrement"`
	Name     *string             `json:"name" bun:"name"`
	Location *map[string]float32 `json:"location" bun:"location"`
	Type     *string             `json:"type" bun:"type"`

	BranchID  *int64     `json:"branch_id" bun:"branch_id"`
	CreatedAt *time.Time `json:"created_at" bun:"created_at"`
	CreatedBy *int64     `json:"created_by" bun:"created_by"`
	UpdatedAt *time.Time `json:"updated_at" bun:"updated_at"`
	UpdatedBy *int64     `json:"updated_by" bun:"updated_by"`
	DeletedAt *time.Time `json:"deleted_at" bun:"deleted_at"`
	DeletedBy *int64     `json:"deleted_by" bun:"deleted_by"`
}
