package cashier

import (
	"context"
)

type Repository interface {

	// admin

	AdminGetList(ctx context.Context, filter Filter) ([]AdminGetList, int, error)
	AdminGetDetail(ctx context.Context, id int64) (AdminGetDetail, error)
	AdminCreate(ctx context.Context, request AdminCreateRequest) (AdminCreateResponse, error)
	AdminUpdateAll(ctx context.Context, request AdminUpdateRequest) error
	AdminUpdateColumns(ctx context.Context, request AdminUpdateRequest) error
	AdminDelete(ctx context.Context, id int64) error
	AdminUpdateStatus(ctx context.Context, id int64, status string) error

	// others

	AdminUpdatePassword(ctx context.Context, request AdminUpdatePassword) error
	AdminUpdatePhone(ctx context.Context, request AdminUpdatePhone) error

	// @branch

	BranchGetList(ctx context.Context, filter Filter) ([]BranchGetList, int, error)
	BranchGetDetail(ctx context.Context, id int64) (BranchGetDetail, error)
	BranchCreate(ctx context.Context, request BranchCreateRequest) (BranchCreateResponse, error)
	BranchUpdateAll(ctx context.Context, request BranchUpdateRequest) error
	BranchUpdateColumns(ctx context.Context, request BranchUpdateRequest) error
	BranchDelete(ctx context.Context, id int64) error
	BranchUpdateStatus(ctx context.Context, id int64, status string) error

	// others

	UpdatePassword(ctx context.Context, request BranchUpdatePassword) error
	UpdatePhone(ctx context.Context, request BranchUpdatePhone) error
}
