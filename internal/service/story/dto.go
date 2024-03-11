package story

import (
	"github.com/restaurant/internal/pkg/utils"
	"github.com/uptrace/bun"
	"mime/multipart"
	"time"
)

type Filter struct {
	Limit          *int
	Offset         *int
	FoodCategoryID *int64
	Expired        *bool
	Status         *string
	Fields         map[string][]string
	Joins          map[string]utils.Joins
}

// @admin

type AdminGetList struct {
	ID        int64    `json:"id"`
	Name      string   `json:"name"`
	File      *string  `json:"file"`
	Type      *string  `json:"type"`
	ExpiredAt *string  `json:"expired_at"`
	Expired   *bool    `json:"expired"`
	Duration  *float32 `json:"duration"`
	Status    *string  `json:"status"`
}

type AdminCreateRequest struct {
	Name     *string               `json:"name" form:"name"`
	File     *multipart.FileHeader `json:"file" form:"file"`
	FileLink *string               `json:"-" form:"-"`
	Type     *string               `json:"type" form:"type"`
	Duration *float32              `json:"duration" form:"duration"`
}

type AdminCreateResponse struct {
	bun.BaseModel `bun:"table:stories"`

	ID           int64     `json:"id" bun:"id,pk,autoincrement"`
	Name         *string   `json:"name" bun:"name"`
	File         *string   `json:"file" bun:"file"`
	Type         *string   `json:"type" bun:"type"`
	Duration     *float32  `json:"duration" bun:"duration"`
	ExpiredAt    time.Time `json:"expired_at" bun:"expired_at"`
	RestaurantID int64     `json:"restaurant_id" bun:"restaurant_id"`
	CreatedAt    time.Time `json:"-" bun:"created_at"`
	CreatedBy    int64     `json:"-" bun:"created_by"`
}

// @client

type ClientGetList struct {
	RestaurantId   int64   `json:"restaurant_id"`
	RestaurantName *string `json:"restaurant_name"`
	RestaurantLogo *string `json:"restaurant_logo"`
	Seen           bool    `json:"seen"`
	Stories        []Story `json:"stories"`
}

type Story struct {
	Id       int64    `json:"id"`
	File     *string  `json:"file"`
	Type     *string  `json:"type"`
	Duration *float32 `json:"duration"`
	Seen     *bool    `json:"seen"`
	Status   *string  `json:"status"`
}

// @super-admin

type SuperAdminGetListResponse struct {
	RestaurantId   int64   `json:"restaurant_id"`
	RestaurantName *string `json:"restaurant_name"`
	RestaurantLogo *string `json:"restaurant_logo"`
	Seen           bool    `json:"seen"`
	Stories        []Story `json:"stories"`
}
