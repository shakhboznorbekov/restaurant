package foodCategory

import (
	"github.com/restaurant/internal/pkg/utils"
	"github.com/uptrace/bun"
	"mime/multipart"
	"time"
)

type Filter struct {
	Limit  *int
	Offset *int
	Name   *string
	Page   *int
	Fields map[string][]string
	Joins  map[string]utils.Joins
}

// @super-admin

type SuperAdminGetList struct {
	ID   int64   `json:"id"`
	Name *string `json:"name"`
	Logo *string `json:"logo"`
	Main *bool   `json:"main"`
}

type SuperAdminGetDetail struct {
	ID   int64   `json:"id"`
	Name *string `json:"name"`
	Logo *string `json:"logo"`
	Main *bool   `json:"main"`
}

type SuperAdminCreateRequest struct {
	Name     *string               `json:"name" form:"name"`
	Logo     *multipart.FileHeader `json:"logo" form:"logo"`
	LogoLink *string               `json:"-" form:"-"`
	Main     *bool                 `json:"main" form:"main"`
}

type SuperAdminCreateResponse struct {
	bun.BaseModel `bun:"table:food_category"`

	ID        int64     `json:"id" bun:"id,pk,autoincrement"`
	Name      *string   `json:"name" bun:"name"`
	Logo      *string   `json:"logo" bun:"logo"`
	Main      *bool     `json:"main" bun:"main"`
	CreatedAt time.Time `json:"-" bun:"created_at"`
	CreatedBy int64     `json:"-" bun:"created_by"`
}

type SuperAdminUpdateRequest struct {
	ID       int64                 `json:"id" form:"id"`
	Name     *string               `json:"name" form:"name"`
	Logo     *multipart.FileHeader `json:"-" form:"logo"`
	LogoLink *string               `json:"-" form:"-"`
	Main     *bool                 `json:"main" form:"main"`
}

// @client

type ClientGetList struct {
	ID   int64   `json:"id"`
	Name *string `json:"name"`
	Logo *string `json:"logo"`
}

type BranchGetList struct {
	ID   int64   `json:"id"`
	Name *string `json:"name"`
	Logo *string `json:"logo"`
}

type AdminGetList struct {
	ID   int64   `json:"id"`
	Name *string `json:"name"`
	Logo *string `json:"logo"`
}

// @waiter

type WaiterGetList struct {
	Id   int64   `json:"id" bun:"id"`
	Name *string `json:"name" bun:"name"`
}
