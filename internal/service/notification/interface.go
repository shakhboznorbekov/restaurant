package notification

import (
	"context"
)

type Repository interface {
	// @admin

	AdminGetList(ctx context.Context, filter Filter) ([]AdminGetList, int, error)
	AdminGetDetail(ctx context.Context, id int64) (*AdminGetDetail, error)
	AdminCreate(ctx context.Context, request AdminCreateRequest) (*AdminCreateResponse, error)
	AdminUpdateAll(ctx context.Context, request AdminUpdateRequest) error
	AdminUpdateColumn(ctx context.Context, request AdminUpdateRequest) error
	AdminDelete(ctx context.Context, id int64) error

	//@super-admin

	SuperAdminUpdateStatus(ctx context.Context, id int64, status string) error
	SuperAdminGetList(ctx context.Context, filter Filter) ([]SuperAdminGetList, int, error)
	SuperAdminGetDetail(ctx context.Context, id int64) (*SuperAdminGetDetail, error)
	SuperAdminSend(ctx context.Context, request SuperAdminSendRequest) ([]SuperAdminSendResponse, error)

	// @client

	ClientGetList(ctx context.Context, filter Filter) ([]ClientGetListResponse, int, error)
	ClientGetUnseenCount(ctx context.Context) (int, error)
	ClientSetAsViewed(ctx context.Context, id int64) error
	ClientSetAllAsViewed(ctx context.Context) error
}
