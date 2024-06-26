package product

import (
	"context"
	"github.com/restaurant/internal/entity"
)

type Service struct {
	repo Repository
}

// @admin

func (s Service) AdminGetList(ctx context.Context, filter Filter) ([]AdminGetList, int, error) {
	return s.repo.AdminGetList(ctx, filter)
}

func (s Service) AdminGetDetail(ctx context.Context, id int64) (entity.Product, error) {
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

func (s Service) AdminDelete(ctx context.Context, id int64) error {
	return s.repo.AdminDelete(ctx, id)
}

func (s Service) AdminGetSpendingByBranch(ctx context.Context, filter SpendingFilter) ([]AdminGetSpendingByBranchResponse, error) {
	return s.repo.AdminGetSpendingByBranch(ctx, filter)
}

// @branch

func (s Service) BranchGetList(ctx context.Context, filter Filter) ([]BranchGetList, int, error) {
	return s.repo.BranchGetList(ctx, filter)
}

func (s Service) BranchGetDetail(ctx context.Context, id int64) (entity.Product, error) {
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

func (s Service) BranchDelete(ctx context.Context, id int64) error {
	return s.repo.BranchDelete(ctx, id)
}

// @cashier

func (s Service) CashierGetList(ctx context.Context, filter Filter) ([]CashierGetList, int, error) {
	return s.repo.CashierGetList(ctx, filter)
}

func (s Service) CashierGetDetail(ctx context.Context, id int64) (entity.Product, error) {
	return s.repo.CashierGetDetail(ctx, id)
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

func (s Service) CashierDelete(ctx context.Context, id int64) error {
	return s.repo.CashierDelete(ctx, id)
}

func (s Service) CashierGetSpending(ctx context.Context, filter SpendingFilter) ([]CashierGetSpendingResponse, error) {
	return s.repo.CashierGetSpending(ctx, filter)
}

func NewService(repo Repository) *Service {
	return &Service{repo: repo}
}
