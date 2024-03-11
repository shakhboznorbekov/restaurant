package warehouse_transaction_product

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

func (s Service) AdminGetList(ctx context.Context, filter Filter, transactionID int64) ([]AdminGetListResponse, int, error) {
	return s.repo.AdminGetList(ctx, filter, transactionID)
}

func (s Service) AdminGetDetailByID(ctx context.Context, id int64) (AdminGetDetailByIdResponse, error) {
	return s.repo.AdminGetDetailByID(ctx, id)
}

func (s Service) AdminUpdateColumn(ctx context.Context, request AdminUpdateRequest) error {
	return s.repo.AdminUpdateColumn(ctx, request)
}

func (s Service) AdminDelete(ctx context.Context, id int64) error {
	return s.repo.AdminDelete(ctx, id)
}

// branch

func (s Service) BranchCreate(ctx context.Context, request BranchCreateRequest) (BranchCreateResponse, error) {
	return s.repo.BranchCreate(ctx, request)
}

func (s Service) BranchGetList(ctx context.Context, filter Filter, transactionID int64) ([]BranchGetListResponse, int, error) {
	return s.repo.BranchGetList(ctx, filter, transactionID)
}

func (s Service) BranchGetDetailByID(ctx context.Context, id int64) (BranchGetDetailByIdResponse, error) {
	return s.repo.BranchGetDetailByID(ctx, id)
}

func (s Service) BranchUpdateColumn(ctx context.Context, request BranchUpdateRequest) error {
	return s.repo.BranchUpdateColumn(ctx, request)
}

func (s Service) BranchDelete(ctx context.Context, id int64) error {
	return s.repo.AdminDelete(ctx, id)
}

func NewService(repo Repository) *Service {
	return &Service{repo: repo}
}
