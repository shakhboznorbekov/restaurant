package restaurant_category

import (
	"context"
	"github.com/restaurant/internal/entity"
)

type Repository interface {

	// @super-admin

	SuperAdminGetList(ctx context.Context, filter Filter) ([]SuperAdminGetList, int, error)
	SuperAdminGetDetail(ctx context.Context, id int64) (entity.RestaurantCategory, error)
	SuperAdminCreate(ctx context.Context, request SuperAdminCreateRequest) (SuperAdminCreateResponse, error)
	SuperAdminUpdateAll(ctx context.Context, request SuperAdminUpdateRequest) error
	SuperAdminUpdateColumns(ctx context.Context, request SuperAdminUpdateRequest) error
	SuperAdminDelete(ctx context.Context, id int64) error

	// @admin

	AdminGetList(ctx context.Context, filter Filter) ([]AdminGetList, int, error)

	// @site

	SiteGetList(ctx context.Context) ([]SiteGetListResponse, int, error)
}
