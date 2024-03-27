package service_percentage

import (
	"github.com/uptrace/bun"
	"time"
)

type Filter struct {
	Limit  *int
	Offset *int
}

// @branch-admin

type AdminGetList struct {
	ID      int64    `json:"id"`
	Percent *float64 `json:"percent"`
}

type AdminGetDetail struct {
	ID      int64    `json:"id"`
	Percent *float64 `json:"percent"`
}

type AdminCreateRequest struct {
	Percent *float64 `json:"percent" form:"percent"`
}

type AdminCreateResponse struct {
	bun.BaseModel `bun:"table:service_percentage"`

	ID        int64     `json:"id" bun:"id,pk,autoincrement"`
	Percent   *float64  `json:"percent" bun:"percent"`
	BranchID  *int64    `json:"branch_id" bun:"branch_id"`
	CreatedAt time.Time `json:"-" bun:"created_at"`
	CreatedBy int64     `json:"-" bun:"created_by"`
}

type AdminUpdateRequest struct {
	ID      int64    `json:"id" form:"id"`
	Percent *float64 `json:"percent" form:"percent"`
}

type AdminUpdateBranchRequest struct {
	ID       int64  `json:"id" form:"id"`
	BranchID *int64 `json:"branch_id" form:"branch_id"`
}
