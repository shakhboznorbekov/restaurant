package entity

import (
	"github.com/uptrace/bun"
	"time"
)

type District struct {
	bun.BaseModel `bun:"table:districts"`

	ID   int64              `json:"id" bun:"id,pk,autoincrement"`
	Name *map[string]string `json:"name" bun:"name"`
	Code *int               `json:"code" bun:"code"`

	RegionId  *int64     `json:"region_id" bun:"region_id"`
	CreatedAt *time.Time `json:"created_at" bun:"created_at"`
	CreatedBy *int64     `json:"created_by" bun:"created_by"`
	DeletedAt *time.Time `json:"deleted_at" bun:"deleted_at"`
	DeletedBy *int64     `json:"deleted_by" bun:"deleted_by"`
}

//for =----> qayerda joylashganini belgilash uchun
