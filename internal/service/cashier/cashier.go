package cashier

import (
	"context"
)

type Service struct {
	repo Repository
}

// @admin

func (s Service) AdminGetList(ctx context.Context, filter Filter) ([]AdminGetList, int, error) {
	return s.repo.AdminGetList(ctx, filter)
}

func (s Service) AdminGetDetail(ctx context.Context, id int64) (AdminGetDetail, error) {
	return s.repo.AdminGetDetail(ctx, id)
}

func (s Service) AdminCreate(ctx context.Context, request AdminCreateRequest) (AdminCreateResponse, error) {
	return s.repo.AdminCreate(ctx, request)
}

func (s Service) AdminUpdateAll(ctx context.Context, request AdminUpdateRequest) error {
	return s.repo.AdminUpdateAll(ctx, request)
}

func (s Service) AdminUpdateColumns(ctx context.Context, request AdminUpdateRequest) error {
	return s.repo.AdminUpdateColumns(ctx, request)
}

func (s Service) AdminUpdateStatus(ctx context.Context, id int64, status string) error {
	return s.repo.AdminUpdateStatus(ctx, id, status)
}

func (s Service) AdminDelete(ctx context.Context, id int64) error {
	return s.repo.AdminDelete(ctx, id)
}

// others

func (s Service) AdminUpdatePassword(ctx context.Context, request AdminUpdatePassword) error {
	return s.repo.AdminUpdatePassword(ctx, request)
}

func (s Service) AdminUpdatePhone(ctx context.Context, request AdminUpdatePhone) error {
	return s.repo.AdminUpdatePhone(ctx, request)
}

// @branch

func (s Service) BranchGetList(ctx context.Context, filter Filter) ([]BranchGetList, int, error) {
	return s.repo.BranchGetList(ctx, filter)
}

func (s Service) BranchGetDetail(ctx context.Context, id int64) (BranchGetDetail, error) {
	return s.repo.BranchGetDetail(ctx, id)
}

func (s Service) BranchCreate(ctx context.Context, request BranchCreateRequest) (BranchCreateResponse, error) {
	return s.repo.BranchCreate(ctx, request)
}

func (s Service) BranchUpdateAll(ctx context.Context, request BranchUpdateRequest) error {
	return s.repo.BranchUpdateAll(ctx, request)
}

func (s Service) BranchUpdateColumns(ctx context.Context, request BranchUpdateRequest) error {
	return s.repo.BranchUpdateColumns(ctx, request)
}

func (s Service) BranchUpdateStatus(ctx context.Context, id int64, status string) error {
	return s.repo.BranchUpdateStatus(ctx, id, status)
}

func (s Service) BranchDelete(ctx context.Context, id int64) error {
	return s.repo.BranchDelete(ctx, id)
}

// others

func (s Service) UpdatePassword(ctx context.Context, request BranchUpdatePassword) error {
	return s.repo.UpdatePassword(ctx, request)
}

func (s Service) UpdatePhone(ctx context.Context, request BranchUpdatePhone) error {
	return s.repo.UpdatePhone(ctx, request)
}

func NewService(repo Repository) *Service {
	return &Service{repo: repo}
}
