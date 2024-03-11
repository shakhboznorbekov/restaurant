package user_work_time

import (
	"context"
	"github.com/restaurant/internal/service/attendance"
)

type Service struct {
	repo Repository
}

func NewService(repo Repository) *Service {
	return &Service{repo: repo}
}

// waiter

func (s Service) WaiterCreate(ctx context.Context, request attendance.WaiterCreateResponse) error {
	return s.repo.WaiterCreate(ctx, request)
}

func (s Service) GetDetailByWaiterIDAndDate(ctx context.Context, request Filter) (GetDetailByWaiterIDAndDateResponse, error) {
	return s.repo.GetDetailByWaiterIDAndDate(ctx, request)
}

func (s Service) WaiterGetListWorkTime(ctx context.Context, filter ListFilter) ([]GetListResponse, int, error) {
	return s.repo.WaiterGetListWorkTime(ctx, filter)
}

func (s Service) BranchGetListWaiterWorkTime(ctx context.Context, filter BranchFilter) ([]BranchGetListResponse, int, error) {
	return s.repo.BranchGetListWaiterWorkTime(ctx, filter)
}

func (s Service) CashierGetListWaiterWorkTime(ctx context.Context, filter BranchFilter) ([]BranchGetListResponse, int, error) {
	return s.repo.CashierGetListWaiterWorkTime(ctx, filter)
}

func (s Service) BranchGetDetailWaiterWorkTime(ctx context.Context, filter ListFilter) ([]GetListResponse, int, error) {
	return s.repo.BranchGetDetailWaiterWorkTime(ctx, filter)
}

func (s Service) CashierGetDetailWaiterWorkTime(ctx context.Context, filter ListFilter) ([]GetListResponse, int, error) {
	return s.repo.CashierGetDetailWaiterWorkTime(ctx, filter)
}
