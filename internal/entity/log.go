package entity

import (
	"github.com/uptrace/bun"
	"time"
)

type Log struct {
	bun.BaseModel `bun:"table:logs"`

	ID        int64                  `json:"id"              bun:"id,pk,autoincrement"`
	Data      map[string]interface{} `json:"data"            bun:"data"`
	Action    *string                `json:"action"          bun:"action"`
	CreatedAt *time.Time             `json:"created_at"      bun:"created_at"`
	CreatedBy *int64                 `json:"created_by"      bun:"created_by"`
	UpdatedAt *time.Time             `json:"updated_at"      bun:"updated_at"`
	UpdatedBy *int64                 `json:"updated_by"      bun:"updated_by"`
	DeletedAt *time.Time             `json:"deleted_at"      bun:"deleted_at"`
	DeletedBy *int64                 `json:"deleted_by"      bun:"deleted_by"`
}
