package menu_category

import (
	"context"
	"github.com/restaurant/internal/entity"
)

type Repository interface {

	// @admin

	AdminCreate(ctx context.Context, data AdminCreateRequest) (*entity.MenuCategory, error)
	AdminUpdate(ctx context.Context, data AdminUpdateRequest) error
	AdminGetList(ctx context.Context, filter Filter) ([]entity.MenuCategory, int, error)
	AdminGetDetail(ctx context.Context, id int64) (*entity.MenuCategory, error)
	AdminDelete(ctx context.Context, id int64) error

	// @branch

	BranchCreate(ctx context.Context, data BranchCreateRequest) (*entity.MenuCategory, error)
	BranchUpdate(ctx context.Context, data BranchUpdateRequest) error
	BranchGetList(ctx context.Context, filter Filter) ([]entity.MenuCategory, int, error)
	BranchGetDetail(ctx context.Context, id int64) (*entity.MenuCategory, error)
	BranchDelete(ctx context.Context, id int64) error

	// @cashier

	CashierCreate(ctx context.Context, data CashierCreateRequest) (*entity.MenuCategory, error)
	CashierUpdate(ctx context.Context, data CashierUpdateRequest) error
	CashierGetList(ctx context.Context, filter Filter) ([]entity.MenuCategory, int, error)
	CashierGetDetail(ctx context.Context, id int64) (*entity.MenuCategory, error)
	CashierDelete(ctx context.Context, id int64) error
}
