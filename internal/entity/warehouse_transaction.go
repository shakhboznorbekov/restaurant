package entity

import (
	"github.com/uptrace/bun"
	"time"
)

type WarehouseTransaction struct {
	bun.BaseModel `bun:"table:warehouse_transactions"`

	ID              int64      `json:"id" bun:"id,pk,autoincrement"`
	FromWarehouseID *int64     `json:"from_warehouse_id" bun:"from_warehouse_id"`
	FromPartnerID   *int64     `json:"from_partner_id" bun:"from_partner_id"`
	ToWarehouseID   *int64     `json:"to_warehouse_id" bun:"to_warehouse_id"`
	ToPartnerID     *int64     `json:"to_partner_id" bun:"to_partner_id"`
	CreatedAt       *time.Time `json:"created_at" bun:"created_at"`
	CreatedBy       *int64     `json:"created_by" bun:"created_by"`
	UpdatedAt       *time.Time `json:"updated_at" bun:"updated_at"`
	UpdatedBy       *int64     `json:"updated_by" bun:"updated_by"`
	DeletedAt       *time.Time `json:"deleted_at" bun:"deleted_at"`
	DeletedBy       *int64     `json:"deleted_by" bun:"deleted_by"`
}
