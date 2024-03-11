package warehouse

import (
	"github.com/uptrace/bun"
	"time"
)

type Filter struct {
	Limit    *int
	Offset   *int
	BranchID *int64
	Type     *string
}

// @admin
type AdminGetList struct {
	ID       int64               `json:"id"`
	Name     *string             `json:"name"`
	Location *map[string]float32 `json:"location"`
	Type     *string             `json:"type"`
	BranchID *int64              `json:"branch_id"`
	Branch   *string             `json:"branch"`
}

type AdminGetDetail struct {
	ID       int64               `json:"id"`
	Name     *string             `json:"name"`
	Location *map[string]float32 `json:"location"`
	Type     *string             `json:"type"`
	BranchID *int64              `json:"branch_id"`
}

type AdminCreateRequest struct {
	Name     *string             `json:"name" form:"name"`
	Location *map[string]float32 `json:"location" form:"location"`
	Type     *string             `json:"type" form:"type"`
	BranchID *int64              `json:"branch_id" form:"branch_id"`
}

type AdminCreateResponse struct {
	bun.BaseModel `bun:"table:warehouses"`

	ID        int64               `json:"id" bun:"id,pk,autoincrement"`
	Name      *string             `json:"name" bun:"name"`
	Location  *map[string]float32 `json:"location" bun:"location"`
	Type      *string             `json:"type" bun:"type"`
	BranchID  *int64              `json:"branch_id" bun:"branch_id"`
	CreatedAt time.Time           `json:"created_at" bun:"created_at"`
	CreatedBy int64               `json:"created_by" bun:"created_by"`
}

type AdminUpdateRequest struct {
	ID       int64               `json:"id" form:"id"`
	Name     *string             `json:"name" form:"name"`
	Location *map[string]float32 `json:"location" form:"location"`
	Type     *string             `json:"type" form:"type"`
	BranchID *int64              `json:"branch_id" form:"branch_id"`
}

// @admin
type BranchGetList struct {
	ID       int64               `json:"id"`
	Name     *string             `json:"name"`
	Location *map[string]float32 `json:"location"`
	Type     *string             `json:"type"`
}

type BranchGetDetail struct {
	ID       int64               `json:"id"`
	Name     *string             `json:"name"`
	Location *map[string]float32 `json:"location"`
	Type     *string             `json:"type"`
}

type BranchCreateRequest struct {
	Name     *string             `json:"name" form:"name"`
	Location *map[string]float32 `json:"location" form:"location"`
	Type     *string             `json:"type" form:"type"`
}

type BranchCreateResponse struct {
	bun.BaseModel `bun:"table:warehouses"`

	ID        int64               `json:"id" bun:"id,pk,autoincrement"`
	Name      *string             `json:"name" bun:"name"`
	Location  *map[string]float32 `json:"location" bun:"location"`
	Type      *string             `json:"type" bun:"type"`
	BranchID  *int64              `json:"-" bun:"branch_id"`
	CreatedAt time.Time           `json:"created_at" bun:"created_at"`
	CreatedBy int64               `json:"created_by" bun:"created_by"`
}

type BranchUpdateRequest struct {
	ID       int64               `json:"id" form:"id"`
	Name     *string             `json:"name" form:"name"`
	Location *map[string]float32 `json:"location" form:"location"`
	Type     *string             `json:"type" form:"type"`
}
