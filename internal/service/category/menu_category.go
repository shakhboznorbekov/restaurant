package category

import (
	"context"
	"github.com/restaurant/internal/entity"
)

type Service struct {
	repo Repository
}

// @client

func (s Service) ClientGetList(ctx context.Context, filter Filter) ([]ClientGetList, int, error) {
	return s.repo.ClientGetList(ctx, filter)
}

// @super-admin

func (s Service) SuperAdminGetList(ctx context.Context, filter Filter) ([]SuperAdminGetList, int, error) {
	return s.repo.SuperAdminGetList(ctx, filter)
}

func (s Service) SuperAdminGetDetail(ctx context.Context, id int64) (entity.Category, error) {
	return s.repo.SuperAdminGetDetail(ctx, id)
}

func (s Service) SuperAdminCreate(ctx context.Context, request SuperAdminCreateRequest) (SuperAdminCreateResponse, error) {
	return s.repo.SuperAdminCreate(ctx, request)
}

func (s Service) SuperAdminUpdateAll(ctx context.Context, request SuperAdminUpdateRequest) error {
	return s.repo.SuperAdminUpdateAll(ctx, request)
}

func (s Service) SuperAdminUpdateColumns(ctx context.Context, request SuperAdminUpdateRequest) error {
	return s.repo.SuperAdminUpdateColumns(ctx, request)
}

func (s Service) SuperAdminDelete(ctx context.Context, id int64) error {
	return s.repo.SuperAdminDelete(ctx, id)
}

// @branch

func (s Service) BranchGetList(ctx context.Context, filter Filter) ([]BranchGetList, int, error) {
	return s.repo.BranchGetList(ctx, filter)
}

// @cashier

func (s Service) CashierGetList(ctx context.Context, filter Filter) ([]CashierGetList, int, error) {
	return s.repo.CashierGetList(ctx, filter)
}

// @admin

func (s Service) AdminGetList(ctx context.Context, filter Filter) ([]AdminGetList, int, error) {
	return s.repo.AdminGetList(ctx, filter)
}

// @waiter

func (s Service) WaiterGetList(ctx context.Context) ([]WaiterGetList, error) {
	return s.repo.WaiterGetList(ctx)
}

func NewService(repo Repository) *Service {
	return &Service{repo: repo}
}
