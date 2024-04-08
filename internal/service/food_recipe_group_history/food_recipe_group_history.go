package food_recipe_group_history

import "context"

type Service struct {
	repo Repository
}

// @branch

func (s Service) BranchGetList(ctx context.Context, filter Filter) ([]BranchGetListByFoodID, int, error) {
	return s.repo.BranchGetList(ctx, filter)
}

func (s Service) BranchGetDetail(ctx context.Context, id int64) (*BranchGetDetail, error) {
	return s.repo.BranchGetDetail(ctx, id)
}

func (s Service) BranchCreate(ctx context.Context, request BranchCreateRequest) (*BranchCreateResponse, error) {
	return s.repo.BranchCreate(ctx, request)
}

func (s Service) BranchDelete(ctx context.Context, id int64) error {
	return s.repo.BranchDelete(ctx, id)
}

// @admin

func (s Service) AdminGetList(ctx context.Context, filter Filter) ([]AdminGetListByFoodID, int, error) {
	return s.repo.AdminGetList(ctx, filter)
}

func (s Service) AdminGetDetail(ctx context.Context, id int64) (*AdminGetDetail, error) {
	return s.repo.AdminGetDetail(ctx, id)
}

func (s Service) AdminCreate(ctx context.Context, request AdminCreateRequest) (*AdminCreateResponse, error) {
	return s.repo.AdminCreate(ctx, request)
}

func (s Service) AdminDelete(ctx context.Context, id int64) error {
	return s.repo.AdminDelete(ctx, id)
}

func NewService(repo Repository) *Service {
	return &Service{repo}
}
