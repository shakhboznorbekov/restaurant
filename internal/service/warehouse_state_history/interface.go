package warehouse_state_history

import "context"

type Repository interface {
	//admin

	AdminCreate(ctx context.Context, request AdminCreateRequest) (AdminCreateResponse, error)
	AdminUpdate(ctx context.Context, request AdminUpdateRequest) error
	AdminDeleteTransaction(ctx context.Context, warehouseTransactionID int64) error

	//admin

	BranchCreate(ctx context.Context, request BranchCreateRequest) (BranchCreateResponse, error)
	BranchUpdate(ctx context.Context, request BranchUpdateRequest) error
	BranchDeleteTransaction(ctx context.Context, warehouseTransactionID int64) error
}
