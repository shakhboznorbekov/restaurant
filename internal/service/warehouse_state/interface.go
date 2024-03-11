package warehouse_state

import (
	"context"
	"github.com/restaurant/internal/service/warehouse_transaction_product"
)

type Repository interface {

	//admin

	AdminGetListByWarehouseID(ctx context.Context, warehouseId int64, filter Filter) ([]AdminGetByWarehouseIDList, int, error)
	AdminCreate(ctx context.Context, request warehouse_transaction_product.AdminCreateRequest) (AdminCreateResponse, error)
	AdminUpdate(ctx context.Context, request *warehouse_transaction_product.AdminUpdateRequest) (AdminUpdateResponse, error)
	AdminDeleteTransaction(ctx context.Context, request AdminDeleteTransactionRequest) error

	//branch

	BranchGetListByWarehouseID(ctx context.Context, warehouseId int64, filter Filter) ([]BranchGetByWarehouseIDList, int, error)
	BranchCreate(ctx context.Context, request warehouse_transaction_product.BranchCreateRequest) (BranchCreateResponse, error)
	BranchUpdate(ctx context.Context, request *warehouse_transaction_product.BranchUpdateRequest) (BranchUpdateResponse, error)
	BranchDeleteTransaction(ctx context.Context, request BranchDeleteTransactionRequest) error
}
