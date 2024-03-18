package user

import (
	"context"
	"github.com/restaurant/internal/entity"
)

type Repository interface {

	// @super-admin

	SuperAdminGetList(ctx context.Context, filter Filter) ([]SuperAdminGetList, int, error)
	SuperAdminGetDetail(ctx context.Context, id int64) (entity.User, error)
	SuperAdminCreate(ctx context.Context, request SuperAdminCreateRequest) (SuperAdminCreateResponse, error)
	SuperAdminUpdateAll(ctx context.Context, request SuperAdminUpdateRequest) error
	SuperAdminUpdateColumns(ctx context.Context, request SuperAdminUpdateRequest) error
	SuperAdminDelete(ctx context.Context, id int64) error

	// @client

	ClientCreate(ctx context.Context, request ClientCreateRequest) (ClientCreateResponse, error)
	ClientUpdateAll(ctx context.Context, request ClientUpdateRequest) error
	ClientUpdateColumn(ctx context.Context, request ClientUpdateRequest) error
	GetByPhone(ctx context.Context, phone string) (entity.User, error)
	ClientGetMe(ctx context.Context, id int64) (entity.User, error)
	ClientDeleteMe(ctx context.Context) error
	ClientUpdateMePhone(ctx context.Context, newPhone string) error

	// @admin

	AdminGetList(ctx context.Context, filter Filter) ([]AdminGetList, int, error)
	AdminCreate(ctx context.Context, request AdminCreateRequest) (AdminCreateResponse, error)
	AdminGetDetailByRestaurantID(ctx context.Context, restaurantID int64) (entity.User, error)
	AdminUpdateColumnsByRestaurantID(ctx context.Context, request AdminUpdateByRestaurantIDRequest) error
	AdminUpdateColumns(ctx context.Context, request AdminUpdateRequest) error

	// @branch

	BranchCreate(ctx context.Context, request BranchCreateRequest) (BranchCreateResponse, error)

	// @waiter

	WaiterUpdateMePhone(ctx context.Context, newPhone string) error
	WaiterUpdatePassword(ctx context.Context, password string, waiterId int64) error

	// others

	IsPhoneExists(ctx context.Context, phone string) (bool, error)
	IsSABCPhoneExists(ctx context.Context, phone string) (bool, error)
	IsWaiterPhoneExists(ctx context.Context, phone string) (bool, error)

	//	@cashier

	CashierGetMe(ctx context.Context) (*CashierGetMeResponse, error)
	GetMe(ctx context.Context, userID int64) (*GetMeResponse, error)
}
