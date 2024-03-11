package banner

import (
	"github.com/lib/pq"
	"github.com/restaurant/internal/pkg/utils"
	"github.com/uptrace/bun"
	"mime/multipart"
	"time"
)

type Filter struct {
	Limit          *int
	Offset         *int
	Page           *int
	FoodCategoryID *int64
	Expired        *bool
	Lat            *float64
	Lon            *float64
	Status         *string
	Fields         map[string][]string
	Joins          map[string]utils.Joins
}

// @admin

type BranchGetList struct {
	ID          int64    `json:"id"`
	Title       *string  `json:"title"`
	Description *string  `json:"description"`
	Photo       *string  `json:"photo"`
	Price       *float64 `json:"price"`
	OldPrice    *float64 `json:"old_price"`
	ExpiredAt   *string  `json:"expired_at"`
	Expired     *bool    `json:"expired"`
	Discount    *int     `json:"discount"`
	Status      *string  `json:"status"`
}

type BranchGetDetail struct {
	ID          int64              `json:"id"`
	Title       *map[string]string `json:"title"`
	Description *map[string]string `json:"description"`
	Photo       *string            `json:"photo"`
	Price       *float64           `json:"price"`
	OldPrice    *float64           `json:"old_price"`
	ExpiredAt   string             `json:"expired_at"`
	Discount    *int               `json:"discount"`
	Status      *string            `json:"status"`
	MenuIds     *pq.Int64Array     `json:"menu_ids"`
	Menus       []Menu             `json:"menus"`
}

type BranchCreateRequest struct {
	Title       map[string]string     `json:"title" form:"title"`
	Description map[string]string     `json:"description" form:"description"`
	Photo       *multipart.FileHeader `json:"photo" form:"photo"`
	PhotoLink   *string               `json:"-"`
	ExpiredAt   *string               `json:"expired_at" form:"expired_at"`
	Price       *float64              `json:"price" form:"price"`
	MenuIds     *pq.Int64Array        `json:"menu_ids" form:"menu_ids"`
}

type BranchCreateResponse struct {
	bun.BaseModel `bun:"table:banners"`

	ID          int64             `json:"id" bun:"id,pk,autoincrement"`
	Title       map[string]string `json:"title" bun:"title"`
	Description map[string]string `json:"description" bun:"description"`
	MenuIds     *pq.Int64Array    `json:"menu_ids" bun:"menu_ids"`
	Photo       *string           `json:"photo" bun:"photo"`
	Price       *float64          `json:"price" bun:"price"`
	OldPrice    *float64          `json:"old_price" bun:"old_price"`
	ExpiredAt   time.Time         `json:"expired_at" bun:"expired_at"`
	BranchID    int64             `json:"branch_id" bun:"branch_id"`
	CreatedAt   time.Time         `json:"-" bun:"created_at"`
	CreatedBy   int64             `json:"-" bun:"created_by"`
}

type BranchUpdateRequest struct {
	ID          int64                 `json:"id" form:"id"`
	Title       map[string]string     `json:"title" form:"title"`
	Description map[string]string     `json:"description" form:"description"`
	Photo       *multipart.FileHeader `json:"photo" form:"photo"`
	PhotoLink   *string               `json:"-"`
	Price       *float64              `json:"price" form:"price"`
	MenuIds     *pq.Int64Array        `json:"menu_ids" form:"menu_ids"`
}

type BranchUpdateMenuRequest struct {
	Id     int64  `json:"id" form:"id"`
	MenuId int64  `json:"menu_id" form:"menu_id"`
	Index  *int   `json:"index" form:"index"`
	Type   string `json:"type" form:"type"`
}

// @client

type ClientGetList struct {
	Id              int64    `json:"id"`
	Title           *string  `json:"title"`
	Description     *string  `json:"description"`
	Photo           *string  `json:"photo"`
	Price           *float64 `json:"price"`
	OldPrice        *float64 `json:"old_price"`
	RestaurantID    *int     `json:"restaurant_id"`
	Restaurant      *string  `json:"restaurant"`
	RestaurantPhoto *string  `json:"restaurant_photo"`
	Discount        *int     `json:"discount"`
	Distance        *string  `json:"distance"`
}

type ClientGetDetail struct {
	Title           *string        `json:"title"`
	Description     *string        `json:"description"`
	Photo           *string        `json:"photo"`
	Price           *float64       `json:"price"`
	OldPrice        *float64       `json:"old_price"`
	Discount        *int           `json:"discount"`
	RestaurantID    *int           `json:"restaurant_id"`
	Restaurant      *string        `json:"restaurant"`
	RestaurantPhoto *string        `json:"restaurant_photo"`
	MenuIds         *pq.Int64Array `json:"-"`
	Menus           []Menu         `json:"menus"`
}

type Menu struct {
	Id          int64    `json:"id"`
	Title       *string  `json:"title"`
	Description *string  `json:"description"`
	Photo       *string  `json:"photo"`
	Price       *float64 `json:"price"`
}

// @super-admin

type SuperAdminGetListResponse struct {
	Id              int64    `json:"id"`
	Title           *string  `json:"title"`
	Description     *string  `json:"description"`
	Photo           *string  `json:"photo"`
	Price           *float64 `json:"price"`
	OldPrice        *float64 `json:"old_price"`
	RestaurantID    *int     `json:"restaurant_id"`
	Restaurant      *string  `json:"restaurant"`
	RestaurantPhoto *string  `json:"restaurant_photo"`
	Discount        *int     `json:"discount"`
	Status          *string  `json:"status"`
}

type SuperAdminGetDetailResponse struct {
	Title           map[string]interface{} `json:"title"`
	Description     map[string]interface{} `json:"description"`
	Photo           *string                `json:"photo"`
	Price           *float64               `json:"price"`
	OldPrice        *float64               `json:"old_price"`
	Discount        *int                   `json:"discount"`
	RestaurantID    *int                   `json:"restaurant_id"`
	Restaurant      *string                `json:"restaurant"`
	RestaurantPhoto *string                `json:"restaurant_photo"`
	MenuIds         *pq.Int64Array         `json:"-"`
	Menus           []Menu                 `json:"menus"`
}
