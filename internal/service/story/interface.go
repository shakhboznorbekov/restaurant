package story

import (
	"context"
	"github.com/restaurant/internal/entity"
)

type Repository interface {
	// @admin
	AdminGetList(ctx context.Context, filter Filter) ([]AdminGetList, int, error)
	AdminGetDetail(ctx context.Context, id int64) (entity.Story, error)
	AdminCreate(ctx context.Context, request AdminCreateRequest) (AdminCreateResponse, error)
	AdminDelete(ctx context.Context, id int64) error
	AdminUpdateStatus(ctx context.Context, id int64) error

	// @client
	ClientGetList(ctx context.Context, filter Filter) ([]ClientGetList, int, error)
	ClientSetViewed(ctx context.Context, id int64) error

	// @super-admin

	SuperAdminGetList(ctx context.Context, filter Filter) ([]SuperAdminGetListResponse, int, error)
	SuperAdminUpdateStatus(ctx context.Context, id int64, status string) error
}
