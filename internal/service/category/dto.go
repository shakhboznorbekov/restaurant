package category

import (
	"github.com/restaurant/internal/pkg/utils"
	"mime/multipart"
	"time"

	"github.com/uptrace/bun"
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
	ID     int64   `json:"id"`
	Name   *string `json:"name"`
	Logo   *string `json:"logo"`
	Status *bool   `json:"status"`
}

type SuperAdminGetDetail struct {
	ID     int64   `json:"id"`
	Name   *string `json:"name"`
	Logo   *string `json:"logo"`
	Status *bool   `json:"status"`
}

type SuperAdminCreateRequest struct {
	Name     *string               `json:"name" form:"name"`
	Logo     *multipart.FileHeader `json:"logo" form:"logo"`
	LogoLink *string               `json:"-" form:"-"`
	Status   *bool                 `json:"status" form:"status"`
}

type SuperAdminCreateResponse struct {
	bun.BaseModel `bun:"table:categories"`

	ID        int64     `json:"id" bun:"id,pk,autoincrement"`
	Name      *string   `json:"name" bun:"name"`
	Logo      *string   `json:"-" bun:"logo"`
	Status    *bool     `json:"status" bun:"status"`
	CreatedAt time.Time `json:"-" bun:"created_at"`
	CreatedBy int64     `json:"-" bun:"created_by"`
}

type SuperAdminUpdateRequest struct {
	ID       int64                 `json:"id" form:"id"`
	Name     *string               `json:"name" form:"name"`
	Logo     *multipart.FileHeader `json:"-" form:"logo"`
	LogoLink *string               `json:"-" form:"-"`
	Status   *bool                 `json:"status" form:"status"`
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

type CashierGetList struct {
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
	Logo *string `json:"logo" bun:"logo"`
}
