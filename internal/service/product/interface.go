package product

import (
	"context"
	"github.com/restaurant/internal/entity"
)

type Repository interface {

	// @admin

	AdminGetList(ctx context.Context, filter Filter) ([]AdminGetList, int, error)
	AdminGetDetail(ctx context.Context, id int64) (entity.Product, error)
	AdminCreate(ctx context.Context, request AdminCreateRequest) (AdminCreateResponse, error)
	AdminUpdateAll(ctx context.Context, request AdminUpdateRequest) error
	AdminUpdateColumns(ctx context.Context, request AdminUpdateRequest) error
	AdminDelete(ctx context.Context, id int64) error

	// @branch

	BranchGetList(ctx context.Context, filter Filter) ([]BranchGetList, int, error)
	BranchGetDetail(ctx context.Context, id int64) (entity.Product, error)
	BranchCreate(ctx context.Context, request BranchCreateRequest) (BranchCreateResponse, error)
	BranchUpdateAll(ctx context.Context, request BranchUpdateRequest) error
	BranchUpdateColumns(ctx context.Context, request BranchUpdateRequest) error
	BranchDelete(ctx context.Context, id int64) error
	AdminGetSpendingByBranch(ctx context.Context, filter SpendingFilter) ([]AdminGetSpendingByBranchResponse, error)

	// @cashier

	CashierGetList(ctx context.Context, filter Filter) ([]CashierGetList, int, error)
	CashierGetDetail(ctx context.Context, id int64) (entity.Product, error)
	CashierCreate(ctx context.Context, request CashierCreateRequest) (CashierCreateResponse, error)
	CashierUpdateAll(ctx context.Context, request CashierUpdateRequest) error
	CashierUpdateColumns(ctx context.Context, request CashierUpdateRequest) error
	CashierDelete(ctx context.Context, id int64) error
	CashierGetSpending(ctx context.Context, filter SpendingFilter) ([]CashierGetSpendingResponse, error)
}
