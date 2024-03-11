package printers

import (
	"github.com/uptrace/bun"
	"time"
)

type Filter struct {
	Limit       *int
	Offset      *int
	WarehouseID *int
}

// @branch

type BranchGetList struct {
	ID   int64   `json:"id" bun:"id"`
	Name *string `json:"name" bun:"name"`
	IP   *string `json:"ip" bun:"ip"`
	Port *string `json:"port" bun:"port"`
}

type BranchGetDetail struct {
	ID   int64   `json:"id" bun:"id"`
	Name *string `json:"name" bun:"name"`
	IP   *string `json:"ip" bun:"ip"`
	Port *string `json:"port" bun:"port"`
}

type BranchCreateRequest struct {
	Name *string `json:"name" form:"name"`
	IP   *string `json:"ip" form:"ip"`
	Port *string `json:"port" form:"port"`
}

type BranchCreateResponse struct {
	bun.BaseModel `bun:"table:printers"`

	ID       int64   `json:"id" bun:"id,pk,autoincrement"`
	Name     *string `json:"name" bun:"name"`
	IP       *string `json:"ip" bun:"ip"`
	Port     *string `json:"port" bun:"port"`
	BranchID *int64  `json:"branch_id" bun:"branch_id"`

	CreatedAt time.Time `json:"created_at" bun:"created_at"`
	CreatedBy int64     `json:"created_by" bun:"created_by"`
}

type BranchUpdateRequest struct {
	ID   int64   `json:"id" form:"id"`
	Name *string `json:"name" form:"name"`
	IP   *string `json:"ip" form:"ip"`
	Port *string `json:"port" form:"port"`
}
