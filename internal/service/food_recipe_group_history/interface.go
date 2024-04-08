package food_recipe_group_history

import "context"

type Repository interface {

	// @admin

	AdminGetList(ctx context.Context, filter Filter) ([]AdminGetListByFoodID, int, error)
	AdminGetDetail(ctx context.Context, id int64) (*AdminGetDetail, error)
	AdminCreate(ctx context.Context, request AdminCreateRequest) (*AdminCreateResponse, error)
	AdminDelete(ctx context.Context, id int64) error

	// @branch

	BranchGetList(ctx context.Context, filter Filter) ([]BranchGetListByFoodID, int, error)
	BranchGetDetail(ctx context.Context, id int64) (*BranchGetDetail, error)
	BranchCreate(ctx context.Context, request BranchCreateRequest) (*BranchCreateResponse, error)
	BranchDelete(ctx context.Context, id int64) error
}
