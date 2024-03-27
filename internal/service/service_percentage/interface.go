package service_percentage

import "context"

type Repository interface {
	// @admin
	AdminGetList(ctx context.Context, filter Filter) ([]AdminGetList, int, error)
	AdminGetDetail(ctx context.Context, id int64) (*AdminGetDetail, error)
	AdminCreate(ctx context.Context, request AdminCreateRequest) (AdminCreateResponse, error)
	AdminUpdateAll(ctx context.Context, request AdminUpdateRequest) error
	AdminDelete(ctx context.Context, id int64) error
	AdminUpdateBranchID(ctx context.Context, request AdminUpdateBranchRequest) error

	//	@branch
	BranchCreate(ctx context.Context, request AdminCreateRequest) (AdminCreateResponse, error)
}
