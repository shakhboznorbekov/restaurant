package warehouse_transaction

import (
	"github.com/restaurant/internal/service/warehouse_transaction_product"
	"github.com/uptrace/bun"
	"time"
)

type Filter struct {
	Limit       *int
	Offset      *int
	Page        *int
	WarehouseID *int64
}

type AdminGetListResponse struct {
	ID              int64   `json:"id"`
	FromWarehouseID *int64  `json:"from_warehouse_id"`
	FromWarehouse   *string `json:"from_warehouse"`
	FromPartnerID   *int64  `json:"from_partner_id"`
	FromPartner     *string `json:"from_partner"`
	ToWarehouseID   *int64  `json:"to_warehouse_id"`
	ToWarehouse     *string `json:"to_warehouse"`
	ToPartnerID     *int64  `json:"to_partner_id"`
	ToPartner       *string `json:"to_partner"`
}

type AdminCreateRequest struct {
	FromWarehouseID *int64 `json:"from_warehouse_id" from:"from_warehouse_id"`
	FromPartnerID   *int64 `json:"from_partner_id" from:"from_partner_id"`
	ToWarehouseID   *int64 `json:"to_warehouse_id" from:"to_warehouse_id"`
	ToPartnerID     *int64 `json:"to_partner_id" from:"to_partner_id"`
}

type AdminCreateResponse struct {
	bun.BaseModel `bun:"table:warehouse_transactions"`

	ID              int64     `json:"-" bun:"id,pk,autoincrement"`
	FromWarehouseID *int64    `json:"from_warehouse_id" bun:"from_warehouse_id"`
	FromPartnerID   *int64    `json:"from_partner_id" bun:"from_partner_id"`
	ToWarehouseID   *int64    `json:"to_warehouse_id" bun:"to_warehouse_id"`
	ToPartnerID     *int64    `json:"to_partner_id" bun:"to_partner_id"`
	CreatedAt       time.Time `json:"-" bun:"created_at"`
	CreatedBy       int64     `json:"-" bun:"created_by"`
}

type AdminUpdateRequest struct {
	ID              *int64 `json:"id" form:"id"`
	FromWarehouseID *int64 `json:"from_warehouse_id" from:"from_warehouse_id"`
	FromPartnerID   *int64 `json:"from_partner_id" from:"from_partner_id"`
	ToWarehouseID   *int64 `json:"to_warehouse_id" from:"to_warehouse_id"`
	ToPartnerID     *int64 `json:"to_partner_id" from:"to_partner_id"`
}

type AdminGetDetailByIdResponse struct {
	ID              int64                                                `json:"id"`
	FromWarehouseID *int64                                               `json:"from_warehouse_id"`
	FromPartnerID   *int64                                               `json:"from_partner_id"`
	ToWarehouseID   *int64                                               `json:"to_warehouse_id"`
	ToPartnerID     *int64                                               `json:"to_partner_id"`
	Products        []warehouse_transaction_product.AdminGetListResponse `json:"products"`
}

type BranchGetListResponse struct {
	ID              int64   `json:"id"`
	FromWarehouseID *int64  `json:"from_warehouse_id"`
	FromWarehouse   *string `json:"from_warehouse"`
	FromPartnerID   *int64  `json:"from_partner_id"`
	FromPartner     *string `json:"from_partner"`
	ToWarehouseID   *int64  `json:"to_warehouse_id"`
	ToWarehouse     *string `json:"to_warehouse"`
	ToPartnerID     *int64  `json:"to_partner_id"`
	ToPartner       *string `json:"to_partner"`
}

type BranchCreateRequest struct {
	FromWarehouseID *int64 `json:"from_warehouse_id" from:"from_warehouse_id"`
	FromPartnerID   *int64 `json:"from_partner_id" from:"from_partner_id"`
	ToWarehouseID   *int64 `json:"to_warehouse_id" from:"to_warehouse_id"`
	ToPartnerID     *int64 `json:"to_partner_id" from:"to_partner_id"`
}

type BranchCreateResponse struct {
	bun.BaseModel `bun:"table:warehouse_transactions"`

	ID              int64     `json:"-" bun:"id,pk,autoincrement"`
	FromWarehouseID *int64    `json:"from_warehouse_id" bun:"from_warehouse_id"`
	FromPartnerID   *int64    `json:"from_partner_id" bun:"from_partner_id"`
	ToWarehouseID   *int64    `json:"to_warehouse_id" bun:"to_warehouse_id"`
	ToPartnerID     *int64    `json:"to_partner_id" bun:"to_partner_id"`
	CreatedAt       time.Time `json:"-" bun:"created_at"`
	CreatedBy       int64     `json:"-" bun:"created_by"`
}

type BranchUpdateRequest struct {
	ID              *int64 `json:"id" form:"id"`
	FromWarehouseID *int64 `json:"from_warehouse_id" from:"from_warehouse_id"`
	FromPartnerID   *int64 `json:"from_partner_id" from:"from_partner_id"`
	ToWarehouseID   *int64 `json:"to_warehouse_id" from:"to_warehouse_id"`
	ToPartnerID     *int64 `json:"to_partner_id" from:"to_partner_id"`
}

type BranchGetDetailByIdResponse struct {
	ID              int64                                                 `json:"id"`
	FromWarehouseID *int64                                                `json:"from_warehouse_id"`
	FromWarehouse   *string                                               `json:"from_warehouse"`
	FromPartnerID   *int64                                                `json:"from_partner_id"`
	FromPartner     *string                                               `json:"from_partner"`
	ToWarehouseID   *int64                                                `json:"to_warehouse_id"`
	ToWarehouse     *string                                               `json:"to_warehouse"`
	ToPartnerID     *int64                                                `json:"to_partner_id"`
	Products        []warehouse_transaction_product.BranchGetListResponse `json:"products"`
}
