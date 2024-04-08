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

// #waiter

func (s Service) WaiterGetMe(ctx context.Context) (*GetMeResponse, error) {
	return s.repo.WaiterGetMe(ctx)
}

func (s Service) WaiterGetPersonalInfo(ctx context.Context) (*GetPersonalInfoResponse, error) {
	return s.repo.WaiterGetPersonalInfo(ctx)
}

func (s Service) WaiterGetActivityStatistics(ctx context.Context) (*GetActivityStatistics, error) {
	return s.repo.WaiterGetActivityStatistics(ctx)
}

func (s Service) WaiterGetWeeklyActivityStatistics(ctx context.Context, filter EarnedMoneyFilter) (*GetEarnedMoneyStatistics, error) {
	return s.repo.WaiterGetWeeklyActivityStatistics(ctx, filter)
}

func (s Service) WaiterGetWeeklyAcceptedOrdersStatistics(ctx context.Context, filter EarnedMoneyFilter) (*GetAcceptedOrdersStatistics, error) {
	return s.repo.WaiterGetWeeklyAcceptedOrdersStatistics(ctx, filter)
}

func (s Service) WaiterGetWeeklyRatingStatistics(ctx context.Context, filter Filter) ([]GetWeeklyRating, error) {
	return s.repo.WaiterGetWeeklyRatingStatistics(ctx, filter)
}

func (s Service) CalculateWaitersKPI(ctx context.Context) error {
	return s.repo.CalculateWaitersKPI(ctx)
}

// cashier

// @branch

func (s Service) CashierGetLists(ctx context.Context, filter Filter) ([]CashierGetLists, int, error) {
	return s.repo.CashierGetLists(ctx, filter)
}

func (s Service) CashierGetDetails(ctx context.Context, id int64) (CashierGetDetails, error) {
	return s.repo.CashierGetDetails(ctx, id)
}

func (s Service) CashierCreate(ctx context.Context, request CashierCreateRequest) (CashierCreateResponse, error) {
	return s.repo.CashierCreate(ctx, request)
}

func (s Service) CashierUpdateAll(ctx context.Context, request CashierUpdateRequest) error {
	return s.repo.CashierUpdateAll(ctx, request)
}

func (s Service) CashierUpdateColumns(ctx context.Context, request CashierUpdateRequest) error {
	return s.repo.CashierUpdateColumns(ctx, request)
}

func (s Service) CashierUpdateStatus(ctx context.Context, id int64, status string) error {
	return s.repo.CashierUpdateStatus(ctx, id, status)
}

func (s Service) CashierDelete(ctx context.Context, id int64) error {
	return s.repo.CashierDelete(ctx, id)
}

func (s Service) CashierUpdatePassword(ctx context.Context, request CashierUpdatePassword) error {
	return s.repo.CashierUpdatePassword(ctx, request)
}

func (s Service) CashierUpdatePhone(ctx context.Context, request CashierUpdatePhone) error {
	return s.repo.CashierUpdatePhone(ctx, request)
}

func NewService(repo Repository) *Service {
	return &Service{repo: repo}
}
