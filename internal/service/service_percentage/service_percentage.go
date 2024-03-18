package service_percentage

import "context"

type Service struct {
	repo Repository
}

func (s Service) AdminGetList(ctx context.Context, filter Filter) ([]AdminGetList, int, error) {
	return s.repo.AdminGetList(ctx, filter)
}

func (s Service) AdminGetDetail(ctx context.Context, id int64) (*AdminGetDetail, error) {
	return s.repo.AdminGetDetail(ctx, id)
}

func (s Service) AdminCreate(ctx context.Context, request AdminCreateRequest) (AdminCreateResponse, error) {
	return s.repo.AdminCreate(ctx, request)
}

func (s Service) AdminUpdateAll(ctx context.Context, request AdminUpdateRequest) error {
	return s.repo.AdminUpdateAll(ctx, request)
}

func (s Service) AdminDelete(ctx context.Context, id int64) error {
	return s.repo.AdminDelete(ctx, id)
}

func NewService(repo Repository) *Service {
	return &Service{repo: repo}
}
