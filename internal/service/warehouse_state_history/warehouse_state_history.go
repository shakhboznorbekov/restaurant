package warehouse_state_history

import (
	"context"
)

type Service struct {
	repo Repository
}

// admin

func (s Service) AdminCreate(ctx context.Context, request AdminCreateRequest) (AdminCreateResponse, error) {
	return s.repo.AdminCreate(ctx, request)
}

func (s Service) AdminUpdate(ctx context.Context, request AdminUpdateRequest) error {
	return s.repo.AdminUpdate(ctx, request)
}

func (s Service) AdminDeleteTransaction(ctx context.Context, warehouseTransactionID int64) error {
	return s.repo.AdminDeleteTransaction(ctx, warehouseTransactionID)
}

// branch

func (s Service) BranchCreate(ctx context.Context, request BranchCreateRequest) (BranchCreateResponse, error) {
	return s.repo.BranchCreate(ctx, request)
}

func (s Service) BranchUpdate(ctx context.Context, request BranchUpdateRequest) error {
	return s.repo.BranchUpdate(ctx, request)
}

func (s Service) BranchDeleteTransaction(ctx context.Context, warehouseTransactionID int64) error {
	return s.repo.BranchDeleteTransaction(ctx, warehouseTransactionID)
}

func NewService(repo Repository) *Service {
	return &Service{repo: repo}
}
