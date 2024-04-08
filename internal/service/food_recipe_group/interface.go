package food_recipe_group

import (
	"golang.org/x/net/context"
)

type Repository interface {

	// @admin

	AdminGetList(ctx context.Context, filter Filter, foodID int64) ([]AdminGetListByFoodID, int, error)
	AdminGetDetail(ctx context.Context, id int64) (*AdminGetDetail, error)
	AdminCreate(ctx context.Context, request AdminCreateRequest) (*AdminCreateResponse, error)
	AdminUpdateAll(ctx context.Context, request AdminUpdateRequest) error
	AdminUpdateColumns(ctx context.Context, request AdminUpdateRequest) error
	AdminDeleteRecipe(ctx context.Context, request AdminDeleteRecipeRequest) error
	AdminDelete(ctx context.Context, id int64) error

	// @branch

	BranchGetList(ctx context.Context, filter Filter, foodID int64) ([]BranchGetListByFoodID, int, error)
	BranchGetDetail(ctx context.Context, id int64) (*BranchGetDetail, error)
	BranchCreate(ctx context.Context, request BranchCreateRequest) (*BranchCreateResponse, error)
	BranchUpdateAll(ctx context.Context, request BranchUpdateRequest) error
	BranchUpdateColumns(ctx context.Context, request BranchUpdateRequest) error
	BranchDeleteRecipe(ctx context.Context, request BranchDeleteRecipeRequest) error
	BranchDelete(ctx context.Context, id int64) error
}
