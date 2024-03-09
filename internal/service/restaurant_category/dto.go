package restaurant_category

import (
	"github.com/restaurant/internal/pkg/utils"
	"github.com/uptrace/bun"
	"time"
)

type Filter struct {
	Limit  *int
	Offset *int
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
	bun.BaseModel `bun:"table:restaurant_category"`

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

// @site

type SiteGetListResponse struct {
	Name  *string `json:"name"`
	Photo *string `json:"photo"`
}
