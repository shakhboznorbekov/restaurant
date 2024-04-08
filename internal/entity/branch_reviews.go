package entity

import (
	"github.com/uptrace/bun"
	"time"
)

type BranchReview struct {
	bun.BaseModel `bun:"table:branch_reviews"`

	ID       int64    `json:"id" bun:"id,pk,autoincrement"`
	Point    *int     `json:"point" bun:"point"`
	Comment  *string  `json:"comment" bun:"comment"`
	Rate     *float32 `json:"rate" bun:"rate"`
	UserID   int64    `json:"user_id" bun:"user_id"`
	BranchID int64    `json:"branch_id" bun:"branch_id"`

	CreatedAt *time.Time `json:"created_at" bun:"created_at"`
	CreatedBy *int64     `json:"created_by" bun:"created_by"`
	UpdatedAt *time.Time `json:"updated_at" bun:"updated_at"`
	UpdatedBy *int64     `json:"updated_by" bun:"updated_by"`
	DeletedAt *time.Time `json:"deleted_at" bun:"deleted_at"`
	DeletedBy *int64     `json:"deleted_by" bun:"deleted_by"`
}
