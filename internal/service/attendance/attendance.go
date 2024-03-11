package attendance

import (
	"context"
)

type Service struct {
	repo Repository
}

func NewService(repo Repository) *Service {
	return &Service{repo: repo}
}

// @waiter

func (s Service) WaiterCameCreate(ctx context.Context) (WaiterCreateResponse, error) {
	return s.repo.WaiterCameCreate(ctx)
}

func (s Service) WaiterGoneCreate(ctx context.Context) (WaiterCreateResponse, error) {
	return s.repo.WaiterGoneCreate(ctx)
}
