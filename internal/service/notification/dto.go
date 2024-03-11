package notification

import (
	"mime/multipart"
	"time"

	"github.com/lib/pq"
	"github.com/uptrace/bun"
)

type Filter struct {
	Limit    *int
	Offset   *int
	Page     *int
	Status   *string
	DeviceId *string
	Whose    *string
	Fields   map[string][]string
}

// @admin

type AdminGetList struct {
	ID          int64   `json:"id"`
	Title       string  `json:"title"`
	Description string  `json:"description"`
	Photo       *string `json:"photo"`
	Status      *string `json:"status"`
}

type AdminGetDetail struct {
	ID          int64                  `json:"id"`
	Title       map[string]interface{} `json:"title"`
	Description map[string]interface{} `json:"description"`
	Photo       *string                `json:"photo"`
	Status      *string                `json:"status"`
}

type AdminCreateRequest struct {
	Title       map[string]interface{} `json:"title" form:"title"`
	Description map[string]interface{} `json:"description" form:"description"`
	Photo       *multipart.FileHeader  `json:"photo" form:"photo"`
	PhotoLink   *string                `json:"-"`
}

type AdminCreateResponse struct {
	bun.BaseModel `bun:"table:notifications"`

	ID           int64                  `json:"id" bun:"id,pk,autoincrement"`
	Title        map[string]interface{} `json:"title" bun:"title"`
	Description  map[string]interface{} `json:"description" bun:"description"`
	Photo        *string                `json:"photo" bun:"photo"`
	Status       *string                `json:"status" bun:"status"`
	RestaurantId *int64                 `json:"restaurant_id" bun:"restaurant_id"`
	CreatedAt    time.Time              `json:"-" bun:"created_at"`
	CreatedBy    int64                  `json:"-" bun:"created_by"`
}

type AdminUpdateRequest struct {
	ID          int64                  `json:"id" form:"id"`
	Title       map[string]interface{} `json:"title" form:"title"`
	Description map[string]interface{} `json:"description" form:"description"`
	Photo       *multipart.FileHeader  `json:"photo" form:"photo"`
	PhotoLink   *string                `json:"-"`
}

// @super-admin

type SuperAdminGetList struct {
	ID          int64   `json:"id" bun:"id"`
	Title       *string `json:"title" bun:"title"`
	Description *string `json:"description" bun:"description"`
	Photo       *string `json:"photo" bun:"photo"`
	Status      *string `json:"status" bun:"status"`
}

type SuperAdminGetDetail struct {
	Title        *string         `json:"title" bun:"title"`
	Description  *string         `json:"description" bun:"description"`
	Photo        *string         `json:"photo"`
	DeviceTokens *pq.StringArray `json:"device_tokens"`
}

type SuperAdminSendRequest struct {
	Title       map[string]interface{} `json:"title" bun:"title" form:"title"`
	Description map[string]interface{} `json:"description" bun:"description" form:"description"`
	Photo       *multipart.FileHeader  `json:"photo" bun:"-" form:"photo"`
	Status      *string                `json:"status" bun:"status" form:"status"`
	PhotoLink   *string                `json:"-" bun:"photo" form:"-"`
	UserId      *int64                 `json:"user_id" bun:"user_id" form:"user_id"`
}

type SuperAdminSendResponse struct {
	DeviceToken string `json:"device_token" bun:"device_token"`
	DeviceLang  string `json:"device_lang" bun:":device_lang"`
}

// @client

type ClientGetListResponse struct {
	ID          int64   `json:"id" bun:"id"`
	Title       *string `json:"title" bun:"title"`
	Description *string `json:"description" bun:"description"`
	Photo       *string `json:"photo" bun:"photo"`
	CreatedAt   *string `json:"created_at" bun:"created_at"`
	Seen        *bool   `json:"seen" bun:"seen"`
}
