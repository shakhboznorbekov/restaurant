package branchReview

import (
	"github.com/uptrace/bun"
	"time"
)

type Filter struct {
	Limit    *int
	Offset   *int
	UserID   *int64
	BranchID *int64
}

// @client

type ClientGetList struct {
	ID         int64    `json:"id" bun:"id"`
	Point      *int     `json:"point" bun:"point"`
	Comment    *string  `json:"comment" bun:"comment"`
	Rate       *float32 `json:"rate" bun:"rate"`
	UserID     *int64   `json:"user_id" bun:"user_id"`
	UserName   *string  `json:"user_name" bun:"user_name"`
	BranchID   *int64   `json:"branch_id" bun:"branch_id"`
	BranchName *string  `json:"branch_name" bun:"branch_name"`
}

type ClientGetDetail struct {
	ID         int64    `json:"id" bun:"id"`
	Point      *int     `json:"point" bun:"point"`
	Comment    *string  `json:"comment" bun:"comment"`
	Rate       *float32 `json:"rate" bun:"rate"`
	UserID     *int64   `json:"user_id" bun:"user_id"`
	UserName   *string  `json:"user_name" bun:"user_name"`
	BranchID   *int64   `json:"branch_id" bun:"branch_id"`
	BranchName *string  `json:"branch_name" bun:"branch_name"`
}

type ClientCreateRequest struct {
	Point    *int     `json:"point" form:"point"`
	Comment  *string  `json:"comment" form:"comment"`
	BranchID *int64   `json:"branch_id" form:"branch_id"`
	Rate     *float32 `json:"rate" form:"rate"`
}

type ClientCreateResponse struct {
	bun.BaseModel `bun:"table:branch_reviews"`

	ID       int64    `json:"id" bun:"id,pk,autoincrement"`
	Point    *int     `json:"point" bun:"point"`
	Comment  *string  `json:"comment" bun:"comment"`
	Rate     *float32 `json:"rate" bun:"rate"`
	UserID   *int64   `json:"user_id" bun:"user_id"`
	BranchID *int64   `json:"branch_id" bun:"branch_id"`

	CreatedAt time.Time `json:"created_at" bun:"created_at"`
	CreatedBy int64     `json:"created_by" bun:"created_by"`
}

type ClientUpdateRequest struct {
	ID       int64    `json:"id" form:"id"`
	Point    *int     `json:"point" form:"point"`
	Comment  *string  `json:"comment" form:"comment"`
	Rate     *float32 `json:"rate" form:"rate"`
	UserID   *int64   `json:"user_id" form:"user_id"`
	BranchID *int64   `json:"branch_id" form:"branch_id"`
}
