package order_menu

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

func (s Service) ClientGetDetail(ctx context.Context, id int64) (entity.OrderMenu, error) {
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

func (s Service) ClientGetOftenList(ctx context.Context, branchID int) ([]ClientGetOftenList, error) {
	return s.repo.ClientGetOftenList(ctx, branchID)
}

// @waiter

func (s Service) WaiterUpdateStatus(ctx context.Context, ids []int64, status string) error {
	return s.repo.WaiterUpdateStatus(ctx, ids, status)
}

func NewService(repo Repository) *Service {
	return &Service{repo: repo}
}
