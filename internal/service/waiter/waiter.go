package waiter

import (
	"context"
)

type Service struct {
	repo Repository
}

// @cashier

func (s Service) CashierGetList(ctx context.Context, filter Filter) ([]CashierGetList, int, error) {
	return s.repo.CashierGetList(ctx, filter)
}

// @admin

func (s Service) AdminGetList(ctx context.Context, filter Filter) ([]AdminGetList, int, error) {
	return s.repo.AdminGetList(ctx, filter)
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

func (s Service) WaiterUpdatePhoto(ctx context.Context, request WaiterPhotoUpdateRequest) error {
	return s.repo.WaiterUpdatePhoto(ctx, request)
}

// others

func (s Service) UpdatePassword(ctx context.Context, request BranchUpdatePassword) error {
	return s.repo.UpdatePassword(ctx, request)
}

func (s Service) UpdatePhone(ctx context.Context, request BranchUpdatePhone) error {
	return s.repo.UpdatePhone(ctx, request)
}

func (s Service) WaiterGetMe(ctx context.Context) (*GetMeResponse, error) {
	return s.repo.WaiterGetMe(ctx)
}

func (s Service) WaiterGetPersonalInfo(ctx context.Context) (*GetPersonalInfoResponse, error) {
	return s.repo.WaiterGetPersonalInfo(ctx)
}

func (s Service) CalculateWaitersKPI(ctx context.Context) error {
	return s.repo.CalculateWaitersKPI(ctx)
}

func NewService(repo Repository) *Service {
	return &Service{repo: repo}
}
