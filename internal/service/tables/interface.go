package tables

import (
	"context"
	"github.com/restaurant/internal/entity"
)

type Repository interface {
	// @admin
	AdminGetList(ctx context.Context, filter Filter) ([]AdminGetList, int, error)
	AdminGetDetail(ctx context.Context, id int64) (entity.Table, error)
	AdminCreate(ctx context.Context, request AdminCreateRequest) (AdminCreateResponse, error)
	AdminUpdateAll(ctx context.Context, request AdminUpdateRequest) error
	AdminUpdateColumns(ctx context.Context, request AdminUpdateRequest) error
	AdminDelete(ctx context.Context, id int64) error

	// @branch
	BranchGetList(ctx context.Context, filter Filter) ([]BranchGetList, int, error)
	BranchGetDetail(ctx context.Context, id int64) (entity.Table, error)
	BranchCreate(ctx context.Context, request BranchCreateRequest) (BranchCreateResponse, error)
	BranchUpdateAll(ctx context.Context, request BranchUpdateRequest) error
	BranchUpdateColumns(ctx context.Context, request BranchUpdateRequest) error
	BranchDelete(ctx context.Context, id int64) error

	// @waiter

	WaiterGetList(ctx context.Context, filter Filter) ([]WaiterGetListResponse, int, error)
}
