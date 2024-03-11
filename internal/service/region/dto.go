package region

import (
	"github.com/restaurant/internal/pkg/utils"
	"github.com/uptrace/bun"
	"time"
)

type Filter struct {
	Limit  *int
	Offset *int
	Name   *string
	Lang   *string
	Fields map[string][]string
	Joins  map[string]utils.Joins
}

// @super-admin

type SuperAdminGetList struct {
	ID   int64   `json:"id"`
	Name *string `json:"name"`
	Code *int    `json:"code"`
}

type SuperAdminGetDetail struct {
	ID   int64              `json:"id"`
	Name *map[string]string `json:"name"`
	Code *int               `json:"code"`
}

type SuperAdminCreateRequest struct {
	Name *map[string]string `json:"name" form:"name"`
	Code *int               `json:"code" form:"code"`
}

type SuperAdminCreateResponse struct {
	bun.BaseModel `bun:"table:regions"`

	ID        int64              `json:"id" bun:"id,pk,autoincrement"`
	Name      *map[string]string `json:"name" bun:"name"`
	Code      *int               `json:"code" bun:"code"`
	CreatedAt time.Time          `json:"-" bun:"created_at"`
	CreatedBy int64              `json:"-" bun:"created_by"`
}

type SuperAdminUpdateRequest struct {
	ID   int64              `json:"id" form:"id"`
	Name *map[string]string `json:"name" form:"name"`
	Code *int               `json:"code" form:"code"`
}
