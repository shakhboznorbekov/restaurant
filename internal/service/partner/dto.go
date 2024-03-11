package partner

import (
	"github.com/restaurant/internal/pkg/utils"
	"github.com/uptrace/bun"
	"time"
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
	ID   int64   `json:"id"`
	Name *string `json:"name"`
	Type *string `json:"type"`
}

type AdminGetDetail struct {
	ID   int64   `json:"id"`
	Name *string `json:"name"`
	Type *string `json:"type"`
}

type AdminCreateRequest struct {
	Name *string `json:"name" form:"name"`
	Type *string `json:"type" form:"type"`
}

type AdminCreateResponse struct {
	bun.BaseModel `bun:"table:partners"`

	ID           int64     `json:"id" bun:"id,pk,autoincrement"`
	Name         *string   `json:"name" bun:"name"`
	Type         *string   `json:"type" bun:"type"`
	RestaurantID *int64    `json:"restaurant_id"`
	CreatedAt    time.Time `json:"-" bun:"created_at"`
	CreatedBy    int64     `json:"-" bun:"created_by"`
}

type AdminUpdateRequest struct {
	ID   int64   `json:"id" form:"id"`
	Name *string `json:"name" bun:"name"`
	Type *string `json:"type" bun:"type"`
}

// @branch

type BranchGetList struct {
	ID   int64   `json:"id"`
	Name *string `json:"name"`
	Type *string `json:"type"`
}
