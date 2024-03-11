package food

import (
	"github.com/lib/pq"
	"github.com/restaurant/internal/pkg/utils"
	"github.com/restaurant/internal/service/user"
	"github.com/uptrace/bun"
	"mime/multipart"
	"time"
)

type Filter struct {
	Limit          *int
	Offset         *int
	Page           *int
	FoodCategoryID *int64
	Search         *string
	Lat            *float64
	Lon            *float64
	BranchID       *int64
	Fields         map[string][]string
	Joins          map[string]utils.Joins
}

// @admin

type AdminGetList struct {
	ID         int64           `json:"id"`
	Name       string          `json:"name"`
	Photos     *pq.StringArray `json:"photos"`
	Price      *float64        `json:"price"`
	CategoryID *int64          `json:"category_id"`
	Category   *string         `json:"category"`
}

type AdminGetDetail struct {
	ID         int64           `json:"id"`
	Name       *string         `json:"name"`
	Photos     *pq.StringArray `json:"photos"`
	Price      *float64        `json:"price"`
	CategoryID *int64          `json:"category_id"`
}

type AdminCreateRequest struct {
	Name       string                   `json:"name" form:"name"`
	Photos     []*multipart.FileHeader  `json:"photos" form:"photos"`
	PhotosLink *string                  `json:"-" form:"-"`
	CategoryID *int64                   `json:"category_id" form:"category_id"`
	User       user.BranchCreateRequest `json:"user" form:"user"`
	Price      *float64                 `json:"price" form:"price"`
}

type AdminCreateResponse struct {
	bun.BaseModel `bun:"table:foods"`

	ID           int64     `json:"id" bun:"id,pk,autoincrement"`
	Name         string    `json:"name" bun:"name"`
	Photos       *string   `json:"photos" bun:"photos"`
	Price        *float64  `json:"price" bun:"price"`
	CategoryID   *int64    `json:"category_id" bun:"category_id"`
	CreatedAt    time.Time `json:"created_at" bun:"created_at"`
	CreatedBy    int64     `json:"created_by" bun:"created_by"`
	RestaurantID *int64    `json:"restaurant_id" bun:"restaurant_id"`
}

type AdminUpdateRequest struct {
	ID         int64                   `json:"id" form:"id"`
	Name       *string                 `json:"name" form:"name"`
	Photos     []*multipart.FileHeader `json:"photos" form:"photos"`
	PhotosLink *string                 `json:"-" form:"-"`
	CategoryID *int64                  `json:"category_id" form:"category_id"`
	Price      *float64                `json:"price" form:"price"`
}

type AdminDeleteImageRequest struct {
	ID         int64 `json:"id" form:"id"`
	ImageIndex *int  `json:"image_index" form:"image_index"`
}

//@branch

type BranchGetList struct {
	ID         int64           `json:"id"`
	Name       string          `json:"name"`
	Photos     *pq.StringArray `json:"photos"`
	Price      *float64        `json:"price"`
	CategoryID *int64          `json:"category_id"`
	BranchID   *int64          `json:"branch_id"`
}

type BranchGetDetail struct {
	ID         int64           `json:"id"`
	Name       *string         `json:"name"`
	Photos     *pq.StringArray `json:"photos"`
	CategoryID *int64          `json:"category_id"`
	Price      *float64        `json:"price"`
}

type BranchCreateRequest struct {
	Name       string                   `json:"name" form:"name"`
	Photos     []*multipart.FileHeader  `json:"photos" form:"photos"`
	PhotosLink *string                  `json:"-" form:"-"`
	CategoryID *int64                   `json:"category_id" form:"category_id"`
	BranchID   *int64                   `json:"branch_id" form:"branch_id"`
	User       user.BranchCreateRequest `json:"user" form:"user"`
	Price      *float64                 `json:"price" form:"price"`
}

type BranchCreateResponse struct {
	bun.BaseModel `bun:"table:foods"`

	ID           int64     `json:"id" bun:"id,pk,autoincrement"`
	Name         string    `json:"name" bun:"name"`
	Photos       *string   `json:"photos" bun:"photos"`
	CategoryID   *int64    `json:"category_id" bun:"category_id"`
	Price        *float64  `json:"price" bun:"price"`
	CreatedAt    time.Time `json:"created_at" bun:"created_at"`
	CreatedBy    int64     `json:"created_by" bun:"created_by"`
	RestaurantID int64     `json:"restaurant_id" bun:"restaurant_id"`
}

type BranchUpdateRequest struct {
	ID         *int64                  `json:"id" form:"id"`
	Name       *string                 `json:"name" form:"name"`
	Photos     []*multipart.FileHeader `json:"photos" form:"photos"`
	PhotosLink *string                 `json:"-" form:"-"`
	CategoryID *int64                  `json:"category_id" form:"category_id"`
	Price      *float64                `json:"price" form:"price"`
}

//@cashier

type CashierGetList struct {
	ID         int64           `json:"id"`
	Name       string          `json:"name"`
	Photos     *pq.StringArray `json:"photos"`
	Price      *float64        `json:"price"`
	CategoryID *int64          `json:"category_id"`
	BranchID   *int64          `json:"branch_id"`
}
