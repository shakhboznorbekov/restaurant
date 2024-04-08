package menu_category

import (
	"context"
	"github.com/restaurant/internal/entity"
)

type Service struct {
	repo Repository
}

// @admin

func (s Service) AdminCreate(ctx context.Context, data AdminCreateRequest) (*entity.MenuCategory, error) {
	return s.repo.AdminCreate(ctx, data)
}

func (s Service) AdminUpdate(ctx context.Context, data AdminUpdateRequest) error {
	return s.repo.AdminUpdate(ctx, data)
}

func (s Service) AdminGetList(ctx context.Context, filter Filter) ([]entity.MenuCategory, int, error) {
	return s.repo.AdminGetList(ctx, filter)
}

func (s Service) AdminGetDetail(ctx context.Context, id int64) (*entity.MenuCategory, error) {
	return s.repo.AdminGetDetail(ctx, id)
}

func (s Service) AdminDelete(ctx context.Context, id int64) error {
	return s.repo.AdminDelete(ctx, id)
}

// @branch

func (s Service) BranchCreate(ctx context.Context, data BranchCreateRequest) (*entity.MenuCategory, error) {
	return s.repo.BranchCreate(ctx, data)
}

func (s Service) BranchUpdate(ctx context.Context, data BranchUpdateRequest) error {
	return s.repo.BranchUpdate(ctx, data)
}

func (s Service) BranchGetList(ctx context.Context, filter Filter) ([]entity.MenuCategory, int, error) {
	return s.repo.BranchGetList(ctx, filter)
}

func (s Service) BranchGetDetail(ctx context.Context, id int64) (*entity.MenuCategory, error) {
	return s.repo.BranchGetDetail(ctx, id)
}

func (s Service) BranchDelete(ctx context.Context, id int64) error {
	return s.repo.BranchDelete(ctx, id)
}

// @cashier

func (s Service) CashierCreate(ctx context.Context, data CashierCreateRequest) (*entity.MenuCategory, error) {
	return s.repo.CashierCreate(ctx, data)
}

func (s Service) CashierUpdate(ctx context.Context, data CashierUpdateRequest) error {
	return s.repo.CashierUpdate(ctx, data)
}

func (s Service) CashierGetList(ctx context.Context, filter Filter) ([]entity.MenuCategory, int, error) {
	return s.repo.CashierGetList(ctx, filter)
}

func (s Service) CashierGetDetail(ctx context.Context, id int64) (*entity.MenuCategory, error) {
	return s.repo.CashierGetDetail(ctx, id)
}

func (s Service) CashierDelete(ctx context.Context, id int64) error {
	return s.repo.CashierDelete(ctx, id)
}

func NewService(repo Repository) *Service {
	return &Service{repo}
}
