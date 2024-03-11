package banner

import (
	"context"
)

type Repository interface {
	// @admin
	BranchGetList(ctx context.Context, filter Filter) ([]BranchGetList, int, error)
	BranchGetDetail(ctx context.Context, id int64) (*BranchGetDetail, error)
	BranchCreate(ctx context.Context, request BranchCreateRequest) (*BranchCreateResponse, error)
	BranchUpdateAll(ctx context.Context, request BranchUpdateRequest) error
	BranchUpdateColumn(ctx context.Context, request BranchUpdateRequest) error
	BranchDelete(ctx context.Context, id int64) error
	BranchUpdateStatus(ctx context.Context, id int64, expireAt string) error

	// @client
	ClientGetList(ctx context.Context, filter Filter) ([]ClientGetList, int, error)
	ClientGetDetail(ctx context.Context, id int64) (*ClientGetDetail, error)

	// @super-admin
	SuperAdminGetList(ctx context.Context, filter Filter) ([]SuperAdminGetListResponse, int, error)
	SuperAdminGetDetail(ctx context.Context, id int64) (*SuperAdminGetDetailResponse, error)
	SuperAdminUpdateStatus(ctx context.Context, id int64, status string) error
}
