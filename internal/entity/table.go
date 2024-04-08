package entity

import (
	"github.com/uptrace/bun"
	"time"
)

type Table struct {
	bun.BaseModel `bun:"table:tables"`

	ID       int64   `json:"id" bun:"id,pk,autoincrement"`
	Number   *int    `json:"number" bun:"number"`
	Status   *string `json:"status" bun:"status"`
	Capacity *int    `json:"capacity" bun:"capacity"`

	BranchID  *int64     `json:"branch_id" bun:"branch_id"`
	HallID    *int64     `json:"hall_id" bun:"hall_id"`
	CreatedAt *time.Time `json:"created_at" bun:"created_at"`
	CreatedBy *int64     `json:"created_by" bun:"created_by"`
	UpdatedAt *time.Time `json:"updated_at" bun:"updated_at"`
	UpdatedBy *int64     `json:"updated_by" bun:"updated_by"`
	DeletedAt *time.Time `json:"deleted_at" bun:"deleted_at"`
	DeletedBy *int64     `json:"deleted_by" bun:"deleted_by"`
}
