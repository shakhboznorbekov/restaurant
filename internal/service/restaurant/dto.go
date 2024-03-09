package restaurant

import (
	"github.com/restaurant/internal/pkg/utils"
	"github.com/restaurant/internal/service/user"
	"github.com/uptrace/bun"
	"mime/multipart"
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
	ID       int64   `json:"id"`
	Name     *string `json:"name"`
	Logo     *string `json:"logo"`
	MiniLogo *string `json:"mini_logo"`
}

type SuperAdminGetDetail struct {
	ID         int64   `json:"id"`
	Name       *string `json:"name"`
	Logo       *string `json:"logo"`
	MiniLogo   *string `json:"mini_logo"`
	WebsiteUrl *string `json:"website_url"`
	User       User    `json:"user"`
}

type SuperAdminCreateRequest struct {
	Name         *string                 `json:"name" form:"name"`
	Logo         *multipart.FileHeader   `json:"logo" form:"logo"`
	LogoLink     *string                 `json:"-" form:"-"`
	User         user.AdminCreateRequest `json:"user" form:"user"`
	MiniLogo     *multipart.FileHeader   `json:"mini_logo" form:"mini_logo"`
	MiniLogoLink *string                 `json:"-" form:"-"`
	WebsiteUrl   *string                 `json:"website_url" form:"website_url"`
}

type SuperAdminCreateResponse struct {
	bun.BaseModel `bun:"table:restaurants"`

	ID         int64     `json:"id" bun:"id,pk,autoincrement"`
	Name       *string   `json:"name" bun:"name"`
	Logo       *string   `json:"logo" bun:"logo"`
	MiniLogo   *string   `json:"mini_logo" bun:"mini_logo"`
	WebsiteUrl *string   `json:"website_url" bun:"website_url"`
	CreatedAt  time.Time `json:"-" bun:"created_at"`
	CreatedBy  int64     `json:"-" bun:"created_by"`
}

type SuperAdminUpdateRequest struct {
	ID           int64                    `json:"id" form:"id"`
	Name         *string                  `json:"name" form:"name"`
	Logo         *multipart.FileHeader    `json:"logo" form:"logo"`
	LogoLink     *string                  `json:"logo_link" form:"logo_link"`
	User         *user.AdminUpdateRequest `json:"user" form:"user"`
	MiniLogo     *multipart.FileHeader    `json:"mini_logo" form:"mini_logo"`
	MiniLogoLink *string                  `json:"-" form:"-"`
	WebsiteUrl   *string                  `json:"website_url" form:"website_url"`
}

// @site

type SiteGetListResponse struct {
	Logo       *string `json:"logo"`
	MiniLogo   *string `json:"mini_logo"`
	WebsiteURL *string `json:"website_url"`
}

// @others

type User struct {
	ID        int64   `json:"id"`
	Name      *string `json:"name"`
	Phone     *string `json:"phone"`
	Role      *string `json:"role"`
	BirthDate *string `json:"birth_date"`
	Gender    *string `json:"gender"`
}
