package warehouse_transaction_product

import (
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
	ID         int64    `json:"id"`
	Amount     *float64 `json:"amount"`
	ProductID  *int64   `json:"product_id"`
	Product    *string  `json:"product"`
	TotalPrice *float64 `json:"total_price"`
}

type AdminCreateRequest struct {
	Amount          *float64 `json:"amount" form:"amount"`
	ProductID       *int64   `json:"product_id" form:"product_id"`
	TotalPrice      *float64 `json:"total_price" form:"total_price"`
	TransactionId   *int64   `json:"transaction_id" form:"transaction_id"`
	FromWarehouseID *int64   `json:"-" bun:"-"`
	FromPartnerID   *int64   `json:"-" bun:"-"`
	ToWarehouseID   *int64   `json:"-" bun:"-"`
	ToPartnerID     *int64   `json:"-" bun:"-"`
}

type AdminCreateResponse struct {
	bun.BaseModel `bun:"table:warehouse_transaction_products"`

	ID            int64     `json:"-" bun:"id,pk,autoincrement"`
	Amount        *float64  `json:"amount" bun:"amount"`
	ProductID     *int64    `json:"product_id" bun:"product_id"`
	TotalPrice    *float64  `json:"total_price" bun:"total_price"`
	TransactionId *int64    `json:"transaction_id" bun:"transaction_id"`
	CreatedAt     time.Time `json:"-" bun:"created_at"`
	CreatedBy     int64     `json:"-" bun:"created_by"`
}

type AdminUpdateRequest struct {
	ID              *int64   `json:"id" form:"id"`
	Amount          *float64 `json:"amount" form:"amount"`
	ProductID       *int64   `json:"product_id" form:"product_id"`
	TotalPrice      *float64 `json:"total_price" form:"total_price"`
	FromWarehouseID *int64   `json:"-" bun:"-"`
	FromPartnerID   *int64   `json:"-" bun:"-"`
	ToWarehouseID   *int64   `json:"-" bun:"-"`
	ToPartnerID     *int64   `json:"-" bun:"-"`
}

type AdminGetDetailByIdResponse struct {
	ID            int64      `json:"id"`
	Amount        *float64   `json:"amount"`
	ProductID     *int64     `json:"product_id"`
	TotalPrice    *float64   `json:"total_price"`
	TransactionId *int64     `json:"transaction_id"`
	CreatedAt     *time.Time `json:"-"`
	CreatedBy     *int64     `json:"-"`
}

type BranchGetListResponse struct {
	ID         int64    `json:"id"`
	Amount     *float64 `json:"amount"`
	ProductID  *int64   `json:"product_id"`
	Product    *string  `json:"product"`
	TotalPrice *float64 `json:"total_price"`
}

type BranchCreateRequest struct {
	Amount          *float64 `json:"amount" form:"amount"`
	ProductID       *int64   `json:"product_id" form:"product_id"`
	TotalPrice      *float64 `json:"total_price" form:"total_price"`
	TransactionId   *int64   `json:"transaction_id" form:"transaction_id"`
	FromWarehouseID *int64   `json:"-" bun:"-"`
	FromPartnerID   *int64   `json:"-" bun:"-"`
	ToWarehouseID   *int64   `json:"-" bun:"-"`
	ToPartnerID     *int64   `json:"-" bun:"-"`
}

type BranchCreateResponse struct {
	bun.BaseModel `bun:"table:warehouse_transaction_products"`

	ID            int64     `json:"-" bun:"id,pk,autoincrement"`
	Amount        *float64  `json:"amount" bun:"amount"`
	ProductID     *int64    `json:"product_id" bun:"product_id"`
	TotalPrice    *float64  `json:"total_price" bun:"total_price"`
	TransactionId *int64    `json:"transaction_id" bun:"transaction_id"`
	CreatedAt     time.Time `json:"-" bun:"created_at"`
	CreatedBy     int64     `json:"-" bun:"created_by"`
}

type BranchUpdateRequest struct {
	ID              *int64   `json:"id" form:"id"`
	Amount          *float64 `json:"amount" form:"amount"`
	ProductID       *int64   `json:"product_id" form:"product_id"`
	TotalPrice      *float64 `json:"total_price" form:"total_price"`
	FromWarehouseID *int64   `json:"-" bun:"-"`
	FromPartnerID   *int64   `json:"-" bun:"-"`
	ToWarehouseID   *int64   `json:"-" bun:"-"`
	ToPartnerID     *int64   `json:"-" bun:"-"`
}

type BranchGetDetailByIdResponse struct {
	ID            int64      `json:"id"`
	Amount        *float64   `json:"amount"`
	ProductID     *int64     `json:"product_id"`
	TotalPrice    *float64   `json:"total_price"`
	TransactionId *int64     `json:"transaction_id"`
	CreatedAt     *time.Time `json:"-"`
	CreatedBy     *int64     `json:"-"`
}
