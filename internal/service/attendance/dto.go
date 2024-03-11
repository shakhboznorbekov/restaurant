package attendance

import (
	"github.com/uptrace/bun"
	"time"
)

type WaiterFilter struct {
	Limit  *int
	Offset *int
	Page   *int
}

type WaiterCreateResponse struct {
	bun.BaseModel `bun:"table:attendances"`

	ID         int64      `json:"id"          bun:"-"`
	UserID     *int64     `json:"user_id"     bun:"user_id"`
	CameAt     *time.Time `json:"came_at"      bun:"came_at"`
	GoneAt     *time.Time `json:"gone_at" bun:"gone_at"`
	Action     *string    `json:"-" bun:"-"`
	ActionTime *time.Time `json:"-" bun:"-"`
}
