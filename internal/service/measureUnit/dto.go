package measureUnit

import (
	"github.com/restaurant/internal/pkg/utils"
	"github.com/uptrace/bun"
	"time"
)

type Filter struct {
	Limit  *int
	Offset *int
	Page   *int
	Name   *string
	Fields map[string][]string
	Joins  map[string]utils.Joins
}

// @super-admin

type SuperAdminGetList struct {
	ID   int64   `json:"id"`
	Name *string `json:"name"`
}

type SuperAdminGetDetail struct {
	ID   int64   `json:"id"`
	Name *string `json:"name"`
}

type SuperAdminCreateRequest struct {
	Name *string `json:"name" form:"name"`
}

type SuperAdminCreateResponse struct {
	bun.BaseModel `bun:"measure_unit"`

	ID        int64     `json:"id" bun:"id,pk,autoincrement"`
	Name      *string   `json:"name" bun:"name"`
	CreatedAt time.Time `json:"-" bun:"created_at"`
	CreatedBy int64     `json:"-" bun:"created_by"`
}

type SuperAdminUpdateRequest struct {
	ID   int64   `json:"id" form:"id"`
	Name *string `json:"name" form:"name"`
}

// @admin

type AdminGetList struct {
	ID   int64   `json:"id"`
	Name *string `json:"name"`
}

// @branch

type BranchGetList struct {
	ID   int64   `json:"id"`
	Name *string `json:"name"`
}
