package warehouse_state_history

import (
	"github.com/uptrace/bun"
	"time"
)

type AdminCreateRequest struct {
	Amount                        *float64 `json:"amount" form:"amount"`
	AveragePrice                  *float64 `json:"average_price" form:"average_price"`
	WarehouseStateID              *int64   `json:"warehouse_state_id" form:"warehouse_state_id"`
	WarehouseTransactionProductID *int64   `json:"warehouse_transaction_product_id" form:"warehouse_transaction_product_id"`
}

type AdminCreateResponse struct {
	bun.BaseModel `bun:"table:warehouse_state_history"`

	ID                            int64     `json:"id" bun:"id,pk,autoincrement"`
	Amount                        *float64  `json:"amount" bun:"amount"`
	AveragePrice                  *float64  `json:"average_price" bun:"average_price"`
	WarehouseStateID              *int64    `json:"warehouse_state_id" bun:"warehouse_state_id"`
	WarehouseTransactionProductID *int64    `json:"warehouse_transaction_product_id" bun:"warehouse_transaction_product_id"`
	CreatedAt                     time.Time `json:"-" bun:"created_at"`
	CreatedBy                     int64     `json:"-" bun:"created_by"`
}

type AdminUpdateRequest struct {
	ID                            *int64   `json:"id" from:"id"`
	Amount                        *float64 `json:"amount" form:"amount"`
	AveragePrice                  *float64 `json:"average_price" form:"average_price"`
	WarehouseStateID              *int64   `json:"warehouse_state_id" form:"warehouse_state_id"`
	WarehouseTransactionProductID *int64   `json:"warehouse_transaction_product_id" form:"warehouse_transaction_product_id"`
}

type BranchCreateRequest struct {
	Amount                        *float64 `json:"amount" form:"amount"`
	AveragePrice                  *float64 `json:"average_price" form:"average_price"`
	WarehouseStateID              *int64   `json:"warehouse_state_id" form:"warehouse_state_id"`
	WarehouseTransactionProductID *int64   `json:"warehouse_transaction_product_id" form:"warehouse_transaction_product_id"`
}

type BranchCreateResponse struct {
	bun.BaseModel `bun:"table:warehouse_state_history"`

	ID                            int64     `json:"id" bun:"id,pk,autoincrement"`
	Amount                        *float64  `json:"amount" bun:"amount"`
	AveragePrice                  *float64  `json:"average_price" bun:"average_price"`
	WarehouseStateID              *int64    `json:"warehouse_state_id" bun:"warehouse_state_id"`
	WarehouseTransactionProductID *int64    `json:"warehouse_transaction_product_id" bun:"warehouse_transaction_product_id"`
	CreatedAt                     time.Time `json:"-" bun:"created_at"`
	CreatedBy                     int64     `json:"-" bun:"created_by"`
}

type BranchUpdateRequest struct {
	ID                            *int64   `json:"id" from:"id"`
	Amount                        *float64 `json:"amount" form:"amount"`
	AveragePrice                  *float64 `json:"average_price" form:"average_price"`
	WarehouseStateID              *int64   `json:"warehouse_state_id" form:"warehouse_state_id"`
	WarehouseTransactionProductID *int64   `json:"warehouse_transaction_product_id" form:"warehouse_transaction_product_id"`
}
