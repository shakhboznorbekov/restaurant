package product

import (
	"context"
	"github.com/restaurant/internal/entity"
	"github.com/restaurant/internal/service/product"
	"github.com/restaurant/internal/service/product_recipe"
	"github.com/restaurant/internal/service/product_recipe_group"
	"github.com/restaurant/internal/service/product_recipe_group_history"
)

type Product interface {

	// @admin

	AdminGetList(ctx context.Context, filter product.Filter) ([]product.AdminGetList, int, error)
	AdminGetDetail(ctx context.Context, id int64) (entity.Product, error)
	AdminCreate(ctx context.Context, request product.AdminCreateRequest) (product.AdminCreateResponse, error)
	AdminUpdateAll(ctx context.Context, request product.AdminUpdateRequest) error
	AdminUpdateColumns(ctx context.Context, request product.AdminUpdateRequest) error
	AdminDelete(ctx context.Context, id int64) error
	AdminGetSpendingByBranch(ctx context.Context, filter product.SpendingFilter) ([]product.AdminGetSpendingByBranchResponse, error)

	//	@branch

	BranchGetList(ctx context.Context, filter product.Filter) ([]product.BranchGetList, int, error)
	BranchGetDetail(ctx context.Context, id int64) (entity.Product, error)
	BranchCreate(ctx context.Context, request product.BranchCreateRequest) (product.BranchCreateResponse, error)
	BranchUpdateAll(ctx context.Context, request product.BranchUpdateRequest) error
	BranchUpdateColumns(ctx context.Context, request product.BranchUpdateRequest) error
	BranchDelete(ctx context.Context, id int64) error

	//	@cashier

	CashierGetList(ctx context.Context, filter product.Filter) ([]product.CashierGetList, int, error)
	CashierGetDetail(ctx context.Context, id int64) (entity.Product, error)
	CashierCreate(ctx context.Context, request product.CashierCreateRequest) (product.CashierCreateResponse, error)
	CashierUpdateAll(ctx context.Context, request product.CashierUpdateRequest) error
	CashierUpdateColumns(ctx context.Context, request product.CashierUpdateRequest) error
	CashierDelete(ctx context.Context, id int64) error
	CashierGetSpending(ctx context.Context, filter product.SpendingFilter) ([]product.CashierGetSpendingResponse, error)
}

type ProductRecipe interface {
	AdminGetList(ctx context.Context, filter product_recipe.Filter) ([]product_recipe.AdminGetList, int, error)
	AdminGetDetail(ctx context.Context, id int64) (*product_recipe.AdminGetDetail, error)
	AdminCreate(ctx context.Context, request product_recipe.AdminCreateRequest) (*product_recipe.AdminCreateResponse, error)
	AdminUpdateAll(ctx context.Context, request product_recipe.AdminUpdateRequest) error
	AdminUpdateColumns(ctx context.Context, request product_recipe.AdminUpdateRequest) error
	AdminDelete(ctx context.Context, id int64) error

	//	 @branch

	BranchGetList(ctx context.Context, filter product_recipe.Filter) ([]product_recipe.BranchGetList, int, error)
	BranchGetDetail(ctx context.Context, id int64) (*product_recipe.BranchGetDetail, error)
	BranchCreate(ctx context.Context, request product_recipe.BranchCreateRequest) (*product_recipe.BranchCreateResponse, error)
	BranchUpdateAll(ctx context.Context, request product_recipe.BranchUpdateRequest) error
	BranchUpdateColumns(ctx context.Context, request product_recipe.BranchUpdateRequest) error
	BranchDelete(ctx context.Context, id int64) error
}

type ProductRecipeGroup interface {

	// @admin

	AdminGetList(ctx context.Context, filter product_recipe_group.Filter, productID int64) ([]product_recipe_group.AdminGetListByProductID, int, error)
	AdminGetDetail(ctx context.Context, id int64) (*product_recipe_group.AdminGetDetail, error)
	AdminCreate(ctx context.Context, request product_recipe_group.AdminCreateRequest) (*product_recipe_group.AdminCreateResponse, error)
	AdminUpdateAll(ctx context.Context, request product_recipe_group.AdminUpdateRequest) error
	AdminUpdateColumns(ctx context.Context, request product_recipe_group.AdminUpdateRequest) error
	AdminDeleteRecipe(ctx context.Context, request product_recipe_group.AdminDeleteRecipeRequest) error
	AdminDelete(ctx context.Context, id int64) error

	// @branch

	BranchGetList(ctx context.Context, filter product_recipe_group.Filter, productID int64) ([]product_recipe_group.BranchGetListByProductID, int, error)
	BranchGetDetail(ctx context.Context, id int64) (*product_recipe_group.BranchGetDetail, error)
	BranchCreate(ctx context.Context, request product_recipe_group.BranchCreateRequest) (*product_recipe_group.BranchCreateResponse, error)
	BranchUpdateAll(ctx context.Context, request product_recipe_group.BranchUpdateRequest) error
	BranchUpdateColumns(ctx context.Context, request product_recipe_group.BranchUpdateRequest) error
	BranchDeleteRecipe(ctx context.Context, request product_recipe_group.BranchDeleteRecipeRequest) error
	BranchDelete(ctx context.Context, id int64) error
}

type ProductRecipeGroupHistory interface {

	// @admin

	AdminGetList(ctx context.Context, filter product_recipe_group_history.Filter) ([]product_recipe_group_history.AdminGetListByProductID, int, error)
	AdminGetDetail(ctx context.Context, id int64) (*product_recipe_group_history.AdminGetDetail, error)
	AdminCreate(ctx context.Context, request product_recipe_group_history.AdminCreateRequest) (*product_recipe_group_history.AdminCreateResponse, error)
	AdminDelete(ctx context.Context, id int64) error

	// @branch

	BranchGetList(ctx context.Context, filter product_recipe_group_history.Filter) ([]product_recipe_group_history.BranchGetListByProductID, int, error)
	BranchGetDetail(ctx context.Context, id int64) (*product_recipe_group_history.BranchGetDetail, error)
	BranchCreate(ctx context.Context, request product_recipe_group_history.BranchCreateRequest) (*product_recipe_group_history.BranchCreateResponse, error)
	BranchDelete(ctx context.Context, id int64) error
}
