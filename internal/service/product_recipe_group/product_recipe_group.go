package product_recipe_group

import (
	"golang.org/x/net/context"
)

type Service struct {
	repo Repository
}

// @admin

func (s Service) AdminGetList(ctx context.Context, filter Filter, foodID int64) ([]AdminGetListByProductID, int, error) {
	return s.repo.AdminGetList(ctx, filter, foodID)
}

func (s Service) AdminGetDetail(ctx context.Context, restaurantID int64) (*AdminGetDetail, error) {
	return s.repo.AdminGetDetail(ctx, restaurantID)
}

func (s Service) AdminCreate(ctx context.Context, request AdminCreateRequest) (*AdminCreateResponse, error) {
	return s.repo.AdminCreate(ctx, request)
}

func (s Service) AdminUpdateColumns(ctx context.Context, request AdminUpdateRequest) error {
	return s.repo.AdminUpdateColumns(ctx, request)
}

func (s Service) AdminUpdateAll(ctx context.Context, request AdminUpdateRequest) error {
	return s.repo.AdminUpdateAll(ctx, request)
}

func (s Service) AdminDeleteRecipe(ctx context.Context, request AdminDeleteRecipeRequest) error {
	return s.repo.AdminDeleteRecipe(ctx, request)
}

func (s Service) AdminDelete(ctx context.Context, id int64) error {
	return s.repo.AdminDelete(ctx, id)
}

// @branch

func (s Service) BranchGetList(ctx context.Context, filter Filter, foodID int64) ([]BranchGetListByProductID, int, error) {
	return s.repo.BranchGetList(ctx, filter, foodID)
}

func (s Service) BranchGetDetail(ctx context.Context, restaurantID int64) (*BranchGetDetail, error) {
	return s.repo.BranchGetDetail(ctx, restaurantID)
}

func (s Service) BranchCreate(ctx context.Context, request BranchCreateRequest) (*BranchCreateResponse, error) {
	return s.repo.BranchCreate(ctx, request)
}

func (s Service) BranchUpdateColumns(ctx context.Context, request BranchUpdateRequest) error {
	return s.repo.BranchUpdateColumns(ctx, request)
}

func (s Service) BranchUpdateAll(ctx context.Context, request BranchUpdateRequest) error {
	return s.repo.BranchUpdateAll(ctx, request)
}

func (s Service) BranchDeleteRecipe(ctx context.Context, request BranchDeleteRecipeRequest) error {
	return s.repo.BranchDeleteRecipe(ctx, request)
}

func (s Service) BranchDelete(ctx context.Context, id int64) error {
	return s.repo.BranchDelete(ctx, id)
}

func NewService(repo Repository) *Service {
	return &Service{repo}
}
