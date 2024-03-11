package warehouse_state

import (
	"github.com/uptrace/bun"
	"time"
)

type Filter struct {
	Limit    *int
	Offset   *int
	BranchID *int64
	Type     *string
}

// @admin

type AdminGetList struct {
	ID           int64    `json:"id"`
	Amount       *float64 `json:"amount"`
	ProductID    *int64   `json:"product_id"`
	Product      *string  `json:"product"`
	AveragePrice *float64 `json:"average_price"`
	WarehouseID  *int64   `json:"warehouse_id"`
	Warehouse    *string  `json:"warehouse"`
	BranchID     *int64   `json:"branch_id"`
	Branch       *string  `json:"branch"`
}

type AdminGetByWarehouseIDList struct {
	ID           int64    `json:"id"`
	Amount       *float64 `json:"amount"`
	ProductID    *int64   `json:"product_id"`
	Product      *string  `json:"product"`
	AveragePrice *float64 `json:"average_price"`
}

type AdminGetDetail struct {
	ID           int64    `json:"id"`
	Amount       *float64 `json:"amount"`
	ProductID    *int64   `json:"product_id"`
	AveragePrice *float64 `json:"average_price"`
	WarehouseID  *int64   `json:"warehouse_id"`
	BranchID     *int64   `json:"branch_id"`
	Branch       *string  `json:"branch"`
}

type AdminCreateRequest struct {
	Name     *string             `json:"name" form:"name"`
	Location *map[string]float32 `json:"location" form:"location"`
	Type     *string             `json:"type" form:"type"`
	BranchID *int64              `json:"branch_id" form:"branch_id"`
}

type AdminCreate struct {
	bun.BaseModel `bun:"table:warehouse_state"`

	ID           int64     `json:"id" bun:"id,pk,autoincrement"`
	Amount       *float64  `json:"amount" bun:"amount"`
	ProductID    *int64    `json:"product_id" bun:"product_id"`
	AveragePrice *float64  `json:"average_price" bun:"average_price"`
	WarehouseID  *int64    `json:"warehouse_id" bun:"warehouse_id"`
	CreatedAt    time.Time `json:"created_at" bun:"created_at"`
	CreatedBy    int64     `json:"created_by" bun:"created_by"`
}

type AdminCreateResponse struct {
	IncomeSate  *AdminCreate
	OutcomeSate *AdminCreate
}

type AdminUpdate struct {
	ID                      int64    `json:"id" bun:"id,pk,autoincrement"`
	Amount                  *float64 `json:"amount" bun:"amount"`
	ProductID               *int64   `json:"product_id" bun:"product_id"`
	AveragePrice            *float64 `json:"average_price" bun:"average_price"`
	WarehouseID             *int64   `json:"warehouse_id" bun:"warehouse_id"`
	WarehouseStateHistoryID *int64   `json:"warehouse_state_history_id" bun:"warehouse_state_history_id"`
}

type AdminUpdateResponse struct {
	IncomeSate  *AdminUpdate
	OutcomeSate *AdminUpdate
}

type AdminDeleteTransactionRequest struct {
	TransactionProductID *int64 `json:"transaction_product_id"`
	FromWarehouseID      *int64 `json:"from_warehouse_id"`
	ToWarehouseID        *int64 `json:"to_warehouse_id"`
}

// @branch

type BranchGetByWarehouseIDList struct {
	ID           int64    `json:"id"`
	Amount       *float64 `json:"amount"`
	ProductID    *int64   `json:"product_id"`
	Product      *string  `json:"product"`
	AveragePrice *float64 `json:"average_price"`
}

type BranchGetDetail struct {
	ID           int64    `json:"id"`
	Amount       *float64 `json:"amount"`
	ProductID    *int64   `json:"product_id"`
	AveragePrice *float64 `json:"average_price"`
	WarehouseID  *int64   `json:"warehouse_id"`
	BranchID     *int64   `json:"branch_id"`
	Branch       *string  `json:"branch"`
}

type BranchCreateRequest struct {
	Name     *string             `json:"name" form:"name"`
	Location *map[string]float32 `json:"location" form:"location"`
	Type     *string             `json:"type" form:"type"`
	BranchID *int64              `json:"branch_id" form:"branch_id"`
}

type BranchCreate struct {
	bun.BaseModel `bun:"table:warehouse_state"`

	ID           int64     `json:"id" bun:"id,pk,autoincrement"`
	Amount       *float64  `json:"amount" bun:"amount"`
	ProductID    *int64    `json:"product_id" bun:"product_id"`
	AveragePrice *float64  `json:"average_price" bun:"average_price"`
	WarehouseID  *int64    `json:"warehouse_id" bun:"warehouse_id"`
	CreatedAt    time.Time `json:"created_at" bun:"created_at"`
	CreatedBy    int64     `json:"created_by" bun:"created_by"`
}

type BranchCreateResponse struct {
	IncomeSate  *BranchCreate
	OutcomeSate *BranchCreate
}

type BranchUpdate struct {
	ID                      int64    `json:"id" bun:"id,pk,autoincrement"`
	Amount                  *float64 `json:"amount" bun:"amount"`
	ProductID               *int64   `json:"product_id" bun:"product_id"`
	AveragePrice            *float64 `json:"average_price" bun:"average_price"`
	WarehouseID             *int64   `json:"warehouse_id" bun:"warehouse_id"`
	WarehouseStateHistoryID *int64   `json:"warehouse_state_history_id" bun:"warehouse_state_history_id"`
}

type BranchUpdateResponse struct {
	IncomeSate  *BranchUpdate
	OutcomeSate *BranchUpdate
}

type BranchDeleteTransactionRequest struct {
	TransactionProductID *int64 `json:"transaction_product_id"`
	FromWarehouseID      *int64 `json:"from_warehouse_id"`
	ToWarehouseID        *int64 `json:"to_warehouse_id"`
}

type WarehouseStateHistory struct {
	ID           int64    `json:"id" bun:"id"`
	Amount       *float64 `json:"amount" bun:"amount"`
	ProductID    *int64   `json:"product_id" bun:"product_id"`
	AveragePrice *float64 `json:"average_price" bun:"average_price"`
}
