package banner

import (
	"context"
)

type Service struct {
	repo Repository
}

// @admin

func (s Service) BranchGetList(ctx context.Context, filter Filter) ([]BranchGetList, int, error) {
	return s.repo.BranchGetList(ctx, filter)
}

func (s Service) BranchGetDetail(ctx context.Context, id int64) (*BranchGetDetail, error) {
	return s.repo.BranchGetDetail(ctx, id)
}

func (s Service) BranchUpdateAll(ctx context.Context, request BranchUpdateRequest) error {
	return s.repo.BranchUpdateAll(ctx, request)
}

func (s Service) BranchUpdateColumn(ctx context.Context, request BranchUpdateRequest) error {
	return s.repo.BranchUpdateColumn(ctx, request)
}

func (s Service) BranchUpdateStatus(ctx context.Context, id int64, expireAt string) error {
	return s.repo.BranchUpdateStatus(ctx, id, expireAt)
}

func (s Service) BranchCreate(ctx context.Context, request BranchCreateRequest) (*BranchCreateResponse, error) {
	return s.repo.BranchCreate(ctx, request)
}

func (s Service) BranchDelete(ctx context.Context, id int64) error {
	return s.repo.BranchDelete(ctx, id)
}

// @client

func (s Service) ClientGetList(ctx context.Context, filter Filter) ([]ClientGetList, int, error) {
	return s.repo.ClientGetList(ctx, filter)
}

func (s Service) ClientGetDetail(ctx context.Context, id int64) (*ClientGetDetail, error) {
	return s.repo.ClientGetDetail(ctx, id)
}

// @super-admin

func (s Service) SuperAdminGetList(ctx context.Context, filter Filter) ([]SuperAdminGetListResponse, int, error) {
	return s.repo.SuperAdminGetList(ctx, filter)
}

func (s Service) SuperAdminGetDetail(ctx context.Context, id int64) (*SuperAdminGetDetailResponse, error) {
	return s.repo.SuperAdminGetDetail(ctx, id)
}

func (s Service) SuperAdminUpdateStatus(ctx context.Context, id int64, status string) error {
	return s.repo.SuperAdminUpdateStatus(ctx, id, status)
}

func NewService(repo Repository) *Service {
	return &Service{repo: repo}
}
