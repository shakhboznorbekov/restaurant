package halls

import (
	"github.com/restaurant/internal/pkg/utils"
	"time"

	"github.com/uptrace/bun"
)

type Filter struct {
	Limit    *int
	Offset   *int
	Page     *int
	Search   *string
	BranchID *int
	Fields   map[string][]string
	Joins    map[string]utils.Joins
}

// @admin

type AdminGetList struct {
	ID       int64   `json:"id"`
	Name     *string `json:"name"`
	BranchID *int64  `json:"branch_id"`
	Branch   *string `json:"branch"`
}

type AdminGetDetail struct {
	ID       int64   `json:"id"`
	Name     *string `json:"name"`
	BranchID *int64  `json:"branch_id"`
}

type AdminCreateRequest struct {
	Name     *string `json:"name" form:"name"`
	BranchID *int64  `json:"branch_id" form:"branch_id"`
}

type AdminCreateResponse struct {
	bun.BaseModel `bun:"table:halls"`

	ID        int64     `json:"id" bun:"id,pk,autoincrement"`
	Name      *string   `json:"name" bun:"name"`
	BranchID  *int64    `json:"branch_id" bun:"branch_id"`
	CreatedAt time.Time `json:"-" bun:"created_at"`
	CreatedBy int64     `json:"-" bun:"created_by"`
}

type AdminUpdateRequest struct {
	ID       int64   `json:"id" form:"id"`
	Name     *string `json:"name" form:"name"`
	BranchID *int64  `json:"branch_id" form:"branch_id"`
}

// @branch

type BranchGetList struct {
	ID   int64   `json:"id"`
	Name *string `json:"name"`
}

type BranchGetDetail struct {
	ID   int64   `json:"id"`
	Name *string `json:"name"`
}

type BranchCreateRequest struct {
	Name *string `json:"name" form:"name"`
}

type BranchCreateResponse struct {
	bun.BaseModel `bun:"table:halls"`

	ID        int64     `json:"id" bun:"id,pk,autoincrement"`
	Name      *string   `json:"name" bun:"name"`
	BranchID  *int64    `json:"-" bun:"branch_id"`
	CreatedAt time.Time `json:"-" bun:"created_at"`
	CreatedBy int64     `json:"-" bun:"created_by"`
}

type BranchUpdateRequest struct {
	ID   int64   `json:"id" form:"id"`
	Name *string `json:"name" form:"name"`
}

// @cashier

type CashierGetList struct {
	ID   int64   `json:"id"`
	Name *string `json:"name"`
}

type CashierGetDetail struct {
	ID   int64   `json:"id"`
	Name *string `json:"name"`
}

type CashierCreateRequest struct {
	Name *string `json:"name" form:"name"`
}

type CashierCreateResponse struct {
	bun.BaseModel `bun:"table:halls"`

	ID        int64     `json:"id" bun:"id,pk,autoincrement"`
	Name      *string   `json:"name" bun:"name"`
	BranchID  *int64    `json:"-" bun:"branch_id"`
	CreatedAt time.Time `json:"-" bun:"created_at"`
	CreatedBy int64     `json:"-" bun:"created_by"`
}

type CashierUpdateRequest struct {
	ID   int64   `json:"id" form:"id"`
	Name *string `json:"name" form:"name"`
}

// @waiter

type WaiterGetList struct {
	ID   int64   `json:"id"`
	Name *string `json:"name"`
}
