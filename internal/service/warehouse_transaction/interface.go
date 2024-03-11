package warehouse_transaction

import "context"

type Repository interface {
	//admin

	AdminCreate(ctx context.Context, request AdminCreateRequest) (AdminCreateResponse, error)
	AdminGetList(ctx context.Context, filter Filter) ([]AdminGetListResponse, int, error)
	AdminUpdateColumn(ctx context.Context, request AdminUpdateRequest) error
	AdminGetDetailByID(ctx context.Context, id int64) (AdminGetDetailByIdResponse, error)
	AdminDelete(ctx context.Context, id int64) error

	//branch

	BranchCreate(ctx context.Context, request BranchCreateRequest) (BranchCreateResponse, error)
	BranchGetList(ctx context.Context, filter Filter) ([]BranchGetListResponse, int, error)
	BranchUpdateColumn(ctx context.Context, request BranchUpdateRequest) error
	BranchGetDetailByID(ctx context.Context, id int64) (BranchGetDetailByIdResponse, error)
	BranchDelete(ctx context.Context, id int64) error
}
