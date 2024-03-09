package entity

import (
	"github.com/uptrace/bun"
)

type BranchLikes struct {
	bun.BaseModel `bun:"table:branch_likes"`

	ID       int64 `json:"id" bun:"id,pk,autoincrement"`
	UserID   int64 `json:"user_id" bun:"user_id"`
	BranchID int64 `json:"branch_id" bun:"branch_id"`
}
