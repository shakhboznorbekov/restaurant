package order_payment

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

func (s Service) ClientGetDetail(ctx context.Context, id int64) (entity.OrderPayment, error) {
	return s.repo.ClientGetDetail(ctx, id)
}

func (s Service) ClientCreate(ctx context.Context, request ClientCreateRequest) (ClientCreateResponse, error) {
	return s.repo.ClientCreate(ctx, request)
}

func (s Service) ClientUpdateAll(ctx context.Context, request ClientUpdateRequest) error {
	return s.repo.ClientUpdateAll(ctx, request)
}

func (s Service) ClientUpdateColumns(ctx context.Context, request ClientUpdateRequest) error {
	return s.repo.ClientUpdateColumns(ctx, request)
}

func (s Service) ClientDelete(ctx context.Context, id int64) error {
	return s.repo.ClientDelete(ctx, id)
}

// @cashier
func (s Service) CashierGetList(ctx context.Context, filter Filter) ([]CashierGetList, int, error) {
	return s.repo.CashierGetList(ctx, filter)
}

func (s Service) CashierGetDetail(ctx context.Context, id int64) (entity.OrderPayment, error) {
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

func NewService(repo Repository) *Service {
	return &Service{repo: repo}
}
