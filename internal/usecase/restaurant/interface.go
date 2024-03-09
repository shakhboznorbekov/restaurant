package restaurant

import (
	"context"
	"github.com/restaurant/internal/entity"
	"github.com/restaurant/internal/service/restaurant"
	"github.com/restaurant/internal/service/restaurant_category"
	"github.com/restaurant/internal/service/user"
)

type Restaurant interface {
	SuperAdminGetList(ctx context.Context, filter restaurant.Filter) ([]restaurant.SuperAdminGetList, int, error)
	SuperAdminGetDetail(ctx context.Context, id int64) (entity.Restaurant, error)
	SuperAdminCreate(ctx context.Context, request restaurant.SuperAdminCreateRequest) (restaurant.SuperAdminCreateResponse, error)
	SuperAdminUpdateAll(ctx context.Context, request restaurant.SuperAdminUpdateRequest) error
	SuperAdminUpdateColumns(ctx context.Context, request restaurant.SuperAdminUpdateRequest) error
	SuperAdminDelete(ctx context.Context, id int64) error
	SiteGetList(ctx context.Context) ([]restaurant.SiteGetListResponse, int, error)
}

type User interface {

	// @super-admin

	SuperAdminGetList(ctx context.Context, filter user.Filter) ([]user.SuperAdminGetList, int, error)
	SuperAdminGetDetail(ctx context.Context, id int64) (entity.User, error)
	SuperAdminCreate(ctx context.Context, request user.SuperAdminCreateRequest) (user.SuperAdminCreateResponse, error)
	SuperAdminUpdateAll(ctx context.Context, request user.SuperAdminUpdateRequest) error
	SuperAdminUpdateColumns(ctx context.Context, request user.SuperAdminUpdateRequest) error
	SuperAdminDelete(ctx context.Context, id int64) error

	// @admin

	AdminGetList(ctx context.Context, filter user.Filter) ([]user.AdminGetList, int, error)
	AdminCreate(ctx context.Context, request user.AdminCreateRequest) (user.AdminCreateResponse, error)
	AdminUpdateColumns(ctx context.Context, request user.AdminUpdateRequest) error
	AdminGetDetailByRestaurantID(ctx context.Context, restaurantID int64) (entity.User, error)
	AdminUpdateColumnsByRestaurantID(ctx context.Context, request user.AdminUpdateByRestaurantIDRequest) error

	// @branch

	BranchCreate(ctx context.Context, request user.BranchCreateRequest) (user.BranchCreateResponse, error)

	// others

	IsPhoneExists(ctx context.Context, phone string) (bool, error)
}

//type Branch interface {
//
//	// @admin
//
//	AdminGetList(ctx context.Context, filter branch.Filter) ([]branch.AdminGetList, int, error)
//	AdminGetDetail(ctx context.Context, id int64) (entity.Branch, error)
//	AdminCreate(ctx context.Context, request branch.AdminCreateRequest) (branch.AdminCreateResponse, error)
//	AdminUpdateAll(ctx context.Context, request branch.AdminUpdateRequest) error
//	AdminUpdateColumns(ctx context.Context, request branch.AdminUpdateRequest) error
//	AdminDelete(ctx context.Context, id int64) error
//	AdminDeleteImage(ctx context.Context, request branch.AdminDeleteImageRequest) error
//
//	// @client
//
//	ClientGetList(ctx context.Context, filter branch.Filter) ([]branch.ClientGetList, int, error)
//	ClientGetMapList(ctx context.Context, filter branch.Filter) ([]branch.ClientGetMapList, int, error)
//	ClientGetDetail(ctx context.Context, id int64) (branch.ClientGetDetail, error)
//	ClientNearlyBranchGetList(ctx context.Context, filter branch.Filter) ([]branch.ClientGetList, int, error)
//	ClientUpdateColumns(ctx context.Context, request branch.ClientUpdateRequest) error
//
//	// @token
//
//	BranchGetToken(ctx context.Context) (branch.BranchGetToken, error)
//	WsGetByToken(ctx context.Context, token string) (branch.WsGetByTokenResponse, error)
//	WsUpdateTokenExpiredAt(ctx context.Context, id int64) (string, error)
//}

type RestaurantCategory interface {

	// @super-admin

	SuperAdminGetList(ctx context.Context, filter restaurant_category.Filter) ([]restaurant_category.SuperAdminGetList, int, error)
	SuperAdminGetDetail(ctx context.Context, id int64) (entity.RestaurantCategory, error)
	SuperAdminCreate(ctx context.Context, request restaurant_category.SuperAdminCreateRequest) (restaurant_category.SuperAdminCreateResponse, error)
	SuperAdminUpdateAll(ctx context.Context, request restaurant_category.SuperAdminUpdateRequest) error
	SuperAdminUpdateColumns(ctx context.Context, request restaurant_category.SuperAdminUpdateRequest) error
	SuperAdminDelete(ctx context.Context, id int64) error

	// @admin

	AdminGetList(ctx context.Context, filter restaurant_category.Filter) ([]restaurant_category.AdminGetList, int, error)

	// @site

	SiteGetList(ctx context.Context) ([]restaurant_category.SiteGetListResponse, int, error)
}

//
//type Table interface {
//
//	// @admin
//
//	AdminGetList(ctx context.Context, filter tables.Filter) ([]tables.AdminGetList, int, error)
//	AdminGetDetail(ctx context.Context, id int64) (entity.Table, error)
//	AdminCreate(ctx context.Context, request tables.AdminCreateRequest) (tables.AdminCreateResponse, error)
//	AdminUpdateAll(ctx context.Context, request tables.AdminUpdateRequest) error
//	AdminUpdateColumns(ctx context.Context, request tables.AdminUpdateRequest) error
//	AdminDelete(ctx context.Context, id int64) error
//
//	// @branch
//
//	BranchGetList(ctx context.Context, filter tables.Filter) ([]tables.BranchGetList, int, error)
//	BranchGetDetail(ctx context.Context, id int64) (entity.Table, error)
//	BranchCreate(ctx context.Context, request tables.BranchCreateRequest) (tables.BranchCreateResponse, error)
//	BranchUpdateAll(ctx context.Context, request tables.BranchUpdateRequest) error
//	BranchUpdateColumns(ctx context.Context, request tables.BranchUpdateRequest) error
//	BranchDelete(ctx context.Context, id int64) error
//
//	// @waiter
//
//	WaiterGetList(ctx context.Context, filter tables.Filter) ([]tables.WaiterGetListResponse, int, error)
//}
//
//type BranchReview interface {
//	// @client
//
//	ClientGetList(ctx context.Context, filter branchReview.Filter) ([]branchReview.ClientGetList, int, error)
//	ClientGetDetail(ctx context.Context, id int64) (branchReview.ClientGetDetail, error)
//	ClientCreate(ctx context.Context, request branchReview.ClientCreateRequest) (branchReview.ClientCreateResponse, error)
//	ClientUpdateAll(ctx context.Context, request branchReview.ClientUpdateRequest) error
//	ClientUpdateColumns(ctx context.Context, request branchReview.ClientUpdateRequest) error
//	ClientDelete(ctx context.Context, id int64) error
//}
//
//type Printers interface {
//	// @branch
//
//	BranchGetList(ctx context.Context, filter printers.Filter) ([]printers.BranchGetList, int, error)
//	BranchGetDetail(ctx context.Context, id int64) (printers.BranchGetDetail, error)
//	BranchCreate(ctx context.Context, request printers.BranchCreateRequest) (printers.BranchCreateResponse, error)
//	BranchUpdateAll(ctx context.Context, request printers.BranchUpdateRequest) error
//	BranchUpdateColumns(ctx context.Context, request printers.BranchUpdateRequest) error
//	BranchDelete(ctx context.Context, id int64) error
//}
