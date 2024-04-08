package warehouse

import (
	"context"
	"github.com/restaurant/internal/service/warehouse"
	"github.com/restaurant/internal/service/warehouse_state"
	"github.com/restaurant/internal/service/warehouse_state_history"
	"github.com/restaurant/internal/service/warehouse_transaction"
	"github.com/restaurant/internal/service/warehouse_transaction_product"
)

type Warehouses interface {
	// admin

	AdminGetList(ctx context.Context, filter warehouse.Filter) ([]warehouse.AdminGetList, int, error)
	AdminGetDetail(ctx context.Context, id int64) (warehouse.AdminGetDetail, error)
	AdminCreate(ctx context.Context, request warehouse.AdminCreateRequest) (warehouse.AdminCreateResponse, error)
	AdminUpdateAll(ctx context.Context, request warehouse.AdminUpdateRequest) error
	AdminUpdateColumns(ctx context.Context, request warehouse.AdminUpdateRequest) error
	AdminDelete(ctx context.Context, id int64) error

	// branch

	BranchGetList(ctx context.Context, filter warehouse.Filter) ([]warehouse.BranchGetList, int, error)
	BranchGetDetail(ctx context.Context, id int64) (warehouse.BranchGetDetail, error)
	BranchCreate(ctx context.Context, request warehouse.BranchCreateRequest) (warehouse.BranchCreateResponse, error)
	BranchUpdateAll(ctx context.Context, request warehouse.BranchUpdateRequest) error
	BranchUpdateColumns(ctx context.Context, request warehouse.BranchUpdateRequest) error
	BranchDelete(ctx context.Context, id int64) error
}

type WarehouseSate interface {
	AdminCreate(ctx context.Context, request warehouse_transaction_product.AdminCreateRequest) (warehouse_state.AdminCreateResponse, error)
	AdminGetListByWarehouseID(ctx context.Context, warehouseId int64, filter warehouse_state.Filter) ([]warehouse_state.AdminGetByWarehouseIDList, int, error)
	AdminUpdate(ctx context.Context, request *warehouse_transaction_product.AdminUpdateRequest) (warehouse_state.AdminUpdateResponse, error)
	AdminDeleteTransaction(ctx context.Context, request warehouse_state.AdminDeleteTransactionRequest) error

	BranchCreate(ctx context.Context, request warehouse_transaction_product.BranchCreateRequest) (warehouse_state.BranchCreateResponse, error)
	BranchGetListByWarehouseID(ctx context.Context, warehouseId int64, filter warehouse_state.Filter) ([]warehouse_state.BranchGetByWarehouseIDList, int, error)
	BranchUpdate(ctx context.Context, request *warehouse_transaction_product.BranchUpdateRequest) (warehouse_state.BranchUpdateResponse, error)
	BranchDeleteTransaction(ctx context.Context, request warehouse_state.BranchDeleteTransactionRequest) error
}

type WarehouseSateHistory interface {
	AdminCreate(ctx context.Context, request warehouse_state_history.AdminCreateRequest) (warehouse_state_history.AdminCreateResponse, error)
	AdminUpdate(ctx context.Context, request warehouse_state_history.AdminUpdateRequest) error
	AdminDeleteTransaction(ctx context.Context, warehouseTransactionID int64) error

	BranchCreate(ctx context.Context, request warehouse_state_history.BranchCreateRequest) (warehouse_state_history.BranchCreateResponse, error)
	BranchUpdate(ctx context.Context, request warehouse_state_history.BranchUpdateRequest) error
	BranchDeleteTransaction(ctx context.Context, warehouseTransactionID int64) error
}

type WarehouseTransaction interface {
	AdminCreate(ctx context.Context, request warehouse_transaction.AdminCreateRequest) (warehouse_transaction.AdminCreateResponse, error)
	AdminGetList(ctx context.Context, filter warehouse_transaction.Filter) ([]warehouse_transaction.AdminGetListResponse, int, error)
	AdminUpdateColumn(ctx context.Context, request warehouse_transaction.AdminUpdateRequest) error
	AdminGetDetailByID(ctx context.Context, id int64) (warehouse_transaction.AdminGetDetailByIdResponse, error)
	AdminDelete(ctx context.Context, id int64) error

	BranchCreate(ctx context.Context, request warehouse_transaction.BranchCreateRequest) (warehouse_transaction.BranchCreateResponse, error)
	BranchGetList(ctx context.Context, filter warehouse_transaction.Filter) ([]warehouse_transaction.BranchGetListResponse, int, error)
	BranchUpdateColumn(ctx context.Context, request warehouse_transaction.BranchUpdateRequest) error
	BranchGetDetailByID(ctx context.Context, id int64) (warehouse_transaction.BranchGetDetailByIdResponse, error)
	BranchDelete(ctx context.Context, id int64) error
}

type WarehouseTransactionProduct interface {
	AdminCreate(ctx context.Context, request warehouse_transaction_product.AdminCreateRequest) (warehouse_transaction_product.AdminCreateResponse, error)
	AdminGetList(ctx context.Context, filter warehouse_transaction_product.Filter, transactionID int64) ([]warehouse_transaction_product.AdminGetListResponse, int, error)
	AdminUpdateColumn(ctx context.Context, request warehouse_transaction_product.AdminUpdateRequest) error
	AdminGetDetailByID(ctx context.Context, id int64) (warehouse_transaction_product.AdminGetDetailByIdResponse, error)
	AdminDelete(ctx context.Context, id int64) error

	BranchCreate(ctx context.Context, request warehouse_transaction_product.BranchCreateRequest) (warehouse_transaction_product.BranchCreateResponse, error)
	BranchGetList(ctx context.Context, filter warehouse_transaction_product.Filter, transactionID int64) ([]warehouse_transaction_product.BranchGetListResponse, int, error)
	BranchUpdateColumn(ctx context.Context, request warehouse_transaction_product.BranchUpdateRequest) error
	BranchGetDetailByID(ctx context.Context, id int64) (warehouse_transaction_product.BranchGetDetailByIdResponse, error)
	BranchDelete(ctx context.Context, id int64) error
}
