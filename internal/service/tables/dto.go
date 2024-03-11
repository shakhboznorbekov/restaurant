package tables

import (
	"github.com/restaurant/internal/pkg/utils"
	"time"

	"github.com/uptrace/bun"
)

type Filter struct {
	Limit  *int
	Offset *int
	Search *string
	Fields map[string][]string
	Joins  map[string]utils.Joins
}

// @admin

type AdminGetList struct {
	ID       int64   `json:"id"`
	Number   *int    `json:"number"`
	Status   *string `json:"status"`
	Capacity *int    `json:"capacity"`
	BranchID *int64  `json:"branch_id"`
}

type AdminGetDetail struct {
	ID       int64   `json:"id"`
	Number   *int    `json:"number"`
	Status   *string `json:"status"`
	Capacity *int    `json:"capacity"`
	BranchID *int64  `json:"branch_id"`
}

type AdminCreateRequest struct {
	Number   *int   `json:"number" form:"number"`
	Capacity *int   `json:"capacity" form:"capacity"`
	BranchID *int64 `json:"branch_id" form:"branch_id"`
}

type AdminCreateResponse struct {
	bun.BaseModel `bun:"table:tables"`

	ID        int64     `json:"id" bun:"id,pk,autoincrement"`
	Number    *int      `json:"number" bun:"number"`
	Status    *string   `json:"status" bun:"status"`
	Capacity  *int      `json:"capacity" bun:"capacity"`
	BranchID  *int64    `json:"branch_id" bun:"branch_id"`
	CreatedAt time.Time `json:"-" bun:"created_at"`
	CreatedBy int64     `json:"-" bun:"created_by"`
}

type AdminUpdateRequest struct {
	ID       int64   `json:"id" form:"id"`
	Number   *int    `json:"number" bun:"number"`
	Status   *string `json:"status" bun:"status"`
	Capacity *int    `json:"capacity" bun:"capacity"`
	BranchID *int64  `json:"branch_id" bun:"branch_id"`
}

// @branch

type BranchGetList struct {
	ID       int64   `json:"id"`
	Number   *int    `json:"number"`
	Status   *string `json:"status"`
	Capacity *int    `json:"capacity"`
	BranchID *int64  `json:"branch_id"`
}

type BranchGetDetail struct {
	ID       int64   `json:"id"`
	Number   *int    `json:"number"`
	Status   *string `json:"status"`
	Capacity *int    `json:"capacity"`
	BranchID *int64  `json:"branch_id"`
}

type BranchCreateRequest struct {
	From     *int `json:"from" form:"from"`
	To       *int `json:"to" form:"to"`
	Number   *int `json:"number" form:"number"`
	Capacity *int `json:"capacity" form:"capacity"`
}

type BranchCreateResponse struct {
	bun.BaseModel `bun:"table:tables"`

	ID        int64     `json:"id" bun:"id,pk,autoincrement"`
	Number    *int      `json:"number" bun:"number"`
	Status    *string   `json:"status" bun:"status"`
	Capacity  *int      `json:"capacity" bun:"capacity"`
	BranchID  *int64    `json:"branch_id" bun:"branch_id"`
	CreatedAt time.Time `json:"-" bun:"created_at"`
	CreatedBy int64     `json:"-" bun:"created_by"`
}

type BranchUpdateRequest struct {
	ID       int64   `json:"id" form:"id"`
	Number   *int    `json:"number" bun:"number"`
	Status   *string `json:"status" bun:"status"`
	Capacity *int    `json:"capacity" bun:"capacity"`
}

type BranchGenerateQRTable struct {
	Tables []int64 `json:"tables" form:"tables"`
}

// @waiter

type WaiterGetListResponse struct {
	ID          int64         `json:"id" bun:"id"`
	Number      *int          `json:"number" bun:"number"`
	Capacity    *int          `json:"capacity" bun:"capacity"`
	ClientCount *int          `json:"client_count" bun:"-"`
	Orders      []WaiterOrder `json:"orders" bun:"-"`
}

type WaiterOrder struct {
	Id          int64 `json:"id" bun:"id"`
	ClientCount *int  `json:"client_count" bun:"client_count"`
	Number      *int  `json:"number" bun:"number"`
}
