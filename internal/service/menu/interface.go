package menu

import (
	"context"
	"github.com/restaurant/internal/entity"
)

type Repository interface {

	// @admin

	AdminGetList(ctx context.Context, filter Filter) ([]AdminGetList, int, error)
	AdminGetDetail(ctx context.Context, id int64) (entity.Menu, error)
	AdminCreate(ctx context.Context, request AdminCreateRequest) ([]AdminCreateResponse, error)
	AdminUpdateAll(ctx context.Context, request AdminUpdateRequest) error
	AdminUpdateColumns(ctx context.Context, request AdminUpdateRequest) error
	AdminDelete(ctx context.Context, id int64) error

	// @client

	ClientGetList(ctx context.Context, filter Filter) ([]ClientGetList, error)
	ClientGetDetail(ctx context.Context, id int64) (ClientGetDetail, error)
	ClientGetListByCategoryID(ctx context.Context, foodCategoryID int, filter Filter) ([]ClientGetListByCategoryID, error)

	// @branch

	BranchGetList(ctx context.Context, filter Filter) ([]BranchGetList, int, error)
	BranchGetDetail(ctx context.Context, id int64) (entity.Menu, error)
	BranchCreate(ctx context.Context, request BranchCreateRequest) (BranchCreateResponse, error)
	BranchUpdateAll(ctx context.Context, request BranchUpdateRequest) error
	BranchUpdateColumns(ctx context.Context, request BranchUpdateRequest) error
	BranchDelete(ctx context.Context, id int64) error
	BranchUpdatePrinterID(ctx context.Context, request BranchUpdatePrinterIDRequest) error
	BranchDeletePrinterID(ctx context.Context, menuID int64) error

	// @cashier

	CashierUpdateColumns(ctx context.Context, request CashierUpdateMenuStatus) error
	CashierGetList(ctx context.Context, filter Filter) ([]CashierGetList, int, error)
	CashierGetDetail(ctx context.Context, id int64) (entity.Menu, error)
	CashierCreate(ctx context.Context, request CashierCreateRequest) (CashierCreateResponse, error)
	CashierUpdateAll(ctx context.Context, request CashierUpdateRequest) error
	CashierUpdateColumn(ctx context.Context, request CashierUpdateRequest) error
	CashierDelete(ctx context.Context, id int64) error
	CashierUpdatePrinterID(ctx context.Context, request CashierUpdatePrinterIDRequest) error
	CashierDeletePrinterID(ctx context.Context, menuID int64) error
	// @waiter

	WaiterGetMenuList(ctx context.Context, filter Filter) ([]WaiterGetMenuListResponse, error)
}
