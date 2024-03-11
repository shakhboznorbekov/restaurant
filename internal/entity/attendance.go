package entity

import (
	"github.com/uptrace/bun"
)

type Attendance struct {
	bun.BaseModel `bun:"table:attendances"`

	ID     int64   `json:"id" bun:"id,pk,autoincrement"`
	UserId *int64  `json:"user_id" bun:"user_id"`
	CameAt *string `json:"came_at" bun:"came_at"`
	GoneAt *string `json:"gone_at" bun:"gone_at"`
}

//for---> waiterlarni keldi ketdisini nazorat qilish uchun
