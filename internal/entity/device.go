package entity

import (
	"github.com/uptrace/bun"
	"time"
)

type Device struct {
	bun.BaseModel `bun:"table:devices"`

	ID          int64   `json:"id" bun:"id,pk,autoincrement"`
	Name        *string `json:"name" bun:"name"`
	UserID      *int64  `json:"user_id" bun:"user_id"`
	DeviceLang  *string `json:"device_lang" bun:"device_lang"`
	IsLogOut    *bool   `json:"is_log_out" bun:"is_log_out"`
	DeviceID    *string `json:"device_id" bun:"device_id"`
	DeviceToken *string `json:"device_token" bun:"device_token"`

	CreatedAt *time.Time `json:"created_at" bun:"created_at"`
	UpdatedAt *time.Time `json:"updated_at" bun:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at" bun:"deleted_at"`
}

//for----> mobile dasturchilar uchun malumotlarni saqlashga
