package branch

import (
	"github.com/lib/pq"
	"github.com/restaurant/internal/service/foodCategory"
	"github.com/restaurant/internal/service/order"
	"github.com/restaurant/internal/service/user"
	"github.com/uptrace/bun"
	"mime/multipart"
	"time"
)

type Filter struct {
	Limit   *int
	Offset  *int
	Page    *int
	Lon     *float64
	Lat     *float64
	IsLiked *bool
	HasMenu *bool
}

type DetailFilter struct {
	ID  int64
	Lon *float64
	Lat *float64
}

// @admin

type AdminGetList struct {
	ID           int64              `json:"id" bun:"id"`
	Location     map[string]float32 `json:"location" bun:"location"`
	Photos       *pq.StringArray    `json:"photos" bun:"photos"`
	Status       *string            `json:"status" bun:"status"`
	WorkTime     map[string]string  `json:"work_time" bun:"work_time"`
	Name         *string            `json:"name" bun:"name"`
	CategoryID   *int64             `json:"category_id" bun:"category_id"`
	CategoryName *string            `json:"category_name" bun:"category_name"`
}

type AdminGetDetail struct {
	ID         int64              `json:"id"`
	Location   map[string]float32 `json:"location"`
	Photos     *pq.StringArray    `json:"photos"`
	Status     *string            `json:"status"`
	WorkTime   map[string]string  `json:"work_time"`
	Name       *string            `json:"name"`
	CategoryID *int64             `json:"category_id"`
	User       User               `json:"user"`
}

type AdminCreateRequest struct {
	Location                   map[string]float32       `json:"location" form:"location"`
	Photos                     []*multipart.FileHeader  `json:"photos" form:"photos"`
	PhotosLink                 *string                  `json:"-" form:"-"`
	Status                     *string                  `json:"status" form:"status"`
	WorkTime                   map[string]string        `json:"work_time" form:"work_time"`
	Name                       *string                  `json:"name" form:"name"`
	CategoryID                 *int64                   `json:"category_id" form:"category_id"`
	User                       user.BranchCreateRequest `json:"user" form:"user"`
	DefaultServicePercentage   *float64                 `json:"default_service_percentage" form:"default_service_percentage"`
	DefaultServicePercentageID *int64                   `json:"-" form:"-"`
}

type AdminCreateResponse struct {
	bun.BaseModel `bun:"table:branches"`

	ID         int64              `json:"id" bun:"id,pk,autoincrement"`
	Location   map[string]float32 `json:"location" bun:"location"`
	Photos     *string            `json:"photos" bun:"photos"`
	Status     *string            `json:"status" bun:"status"`
	WorkTime   map[string]string  `json:"work_time" bun:"work_time"`
	Name       *string            `json:"name" bun:"name"`
	CategoryID *int64             `json:"category_id" bun:"category_id"`

	CreatedAt                time.Time `json:"created_at" bun:"created_at"`
	CreatedBy                int64     `json:"created_by" bun:"created_by"`
	RestaurantID             int64     `json:"restaurant_id" bun:"restaurant_id"`
	DefaultServicePercentage *int64    `json:"default_service_percentage" bun:"default_service_percentage"`
}

type AdminUpdateRequest struct {
	ID         int64                   `json:"id" form:"id"`
	Location   map[string]float32      `json:"location" form:"location"`
	Photos     []*multipart.FileHeader `json:"photos" form:"photos"`
	PhotosLink *string                 `json:"photo_links" form:"photo_links"`
	Status     *string                 `json:"status" form:"status"`
	WorkTime   map[string]string       `json:"work_time" form:"work_time"`
	Name       *string                 `json:"name" form:"name"`
	CategoryID *int64                  `json:"category_id" form:"category_id"`
}

type AdminDeleteImageRequest struct {
	ID         int64 `json:"id" form:"id"`
	ImageIndex *int  `json:"image_index" form:"image_index"`
}

type AdminUpdateBranchAdmin struct {
	ID       int64   `json:"id" form:"id"`
	Password *string `json:"password" form:"password"`
}

// @client

type ClientGetList struct {
	ID            int64              `json:"id" bun:"id"`
	Location      map[string]float32 `json:"location" bun:"location"`
	Photos        *pq.StringArray    `json:"photos" bun:"photos"`
	Status        *string            `json:"status" bun:"status"`
	OpenTime      *string            `json:"open_time" bun:"-"`
	CloseTime     *string            `json:"close_time" bun:"-"`
	Name          *string            `json:"name" bun:"name"`
	CategoryID    *int64             `json:"category_id" bun:"category_id"`
	CategoryName  *string            `json:"category_name" bun:"category_name"`
	Point         *int               `json:"point" bun:"point"`
	Rate          *float32           `json:"rate" bun:"rate"`
	Distance      *string            `json:"distance" bun:"distance"`
	WorkTimeToday *string            `json:"-" bun:"work_time_today"`
	IsLiked       *bool              `json:"is_liked" bun:"is_liked"`
	IsClosed      *bool              `json:"is_closed" bun:"-"`
}

type ClientGetMapList struct {
	ID         int64    `json:"id" bun:"id"`
	Lat        *float32 `json:"lat" bun:"lat"`
	Lon        *float32 `json:"lon" bun:"lon"`
	Logo       *string  `json:"logo" bun:"logo"`
	Status     *string  `json:"status" bun:"status"`
	Name       *string  `json:"name" bun:"name"`
	CategoryID *int64   `json:"category_id" bun:"category_id"`
}

type ClientGetDetail struct {
	ID             int64                        `json:"id" bun:"id"`
	Status         *string                      `json:"status" bun:"status"`
	Location       map[string]float32           `json:"location" bun:"location"`
	Photos         *pq.StringArray              `json:"photos" bun:"photos"`
	WorkTimeToday  *string                      `json:"-" bun:"work_time_today"`
	Name           *string                      `json:"name" bun:"name"`
	Rate           *float32                     `json:"rate" bun:"rate"`
	Category       []foodCategory.ClientGetList `json:"category" bun:"-"`
	Orders         []order.ClientGetDetail      `json:"orders" bun:"-"`
	NewOrders      []order.ClientGetDetail      `json:"new_orders" bun:"-"`
	RestaurantName *string                      `json:"restaurant_name" bun:"restaurant_name"`
	RestaurantLogo *string                      `json:"restaurant_logo" bun:"restaurant_logo"`
	CanOrder       *bool                        `json:"can_order" bun:"-"`
	Distance       *string                      `json:"distance" bun:"-"`

	//CategoryID    *int64                       `json:"category_id" bun:"category_id"`
	//OpenTime     *string `json:"open_time" bun:"-"`
	//CloseTime    *string `json:"close_time" bun:"-"`
	//CategoryName *string `json:"category_name" bun:"category_name"`
	//Point        *int    `json:"point" bun:"point"`
	//IsLiked      *bool   `json:"is_liked" bun:"is_liked"`
	//IsClosed     *bool   `json:"is_closed" bun:"-"`
}

type ClientUpdateRequest struct {
	ID      int64 `json:"id" form:"id"`
	IsLiked *bool `json:"is_liked" form:"is_liked"`
}

type User struct {
	ID        int64   `json:"id"`
	Name      *string `json:"name"`
	Phone     *string `json:"phone"`
	Role      *string `json:"role"`
	BirthDate *string `json:"birth_date"`
	Gender    *string `json:"gender"`
}

// @branch

type BranchGetDetail struct {
	ID   int64   `json:"id"`
	Logo *string `json:"logo"`
}

type CashierGetDetail struct {
	ID   int64   `json:"id"`
	Logo *string `json:"logo"`
}

// @token

type BranchGetToken struct {
	Token          *string `json:"token"`
	TokenExpiredAt *string `json:"token_expired_at"`
}

type WsGetByTokenResponse struct {
	ID             *int64     `json:"id"`
	TokenExpiredAt *time.Time `json:"token_expired_at"`
}
