package warehouse_state

import (
	"context"
	"github.com/restaurant/internal/service/warehouse_transaction_product"
)

type Service struct {
	repo Repository
}

// admin

func (s Service) AdminGetListByWarehouseID(ctx context.Context, warehouseId int64, filter Filter) ([]AdminGetByWarehouseIDList, int, error) {
	return s.repo.AdminGetListByWarehouseID(ctx, warehouseId, filter)
}

func (s Service) AdminCreate(ctx context.Context, request warehouse_transaction_product.AdminCreateRequest) (AdminCreateResponse, error) {
	return s.repo.AdminCreate(ctx, request)
}

func (s Service) AdminUpdate(ctx context.Context, request *warehouse_transaction_product.AdminUpdateRequest) (AdminUpdateResponse, error) {
	return s.repo.AdminUpdate(ctx, request)
}

func (s Service) AdminDeleteTransaction(ctx context.Context, request AdminDeleteTransactionRequest) error {
	return s.repo.AdminDeleteTransaction(ctx, request)
}

// branch

func (s Service) BranchGetListByWarehouseID(ctx context.Context, warehouseId int64, filter Filter) ([]BranchGetByWarehouseIDList, int, error) {
	return s.repo.BranchGetListByWarehouseID(ctx, warehouseId, filter)
}

func (s Service) BranchCreate(ctx context.Context, request warehouse_transaction_product.BranchCreateRequest) (BranchCreateResponse, error) {
	return s.repo.BranchCreate(ctx, request)
}

func (s Service) BranchUpdate(ctx context.Context, request *warehouse_transaction_product.BranchUpdateRequest) (BranchUpdateResponse, error) {
	return s.repo.BranchUpdate(ctx, request)
}

func (s Service) BranchDeleteTransaction(ctx context.Context, request BranchDeleteTransactionRequest) error {
	return s.repo.BranchDeleteTransaction(ctx, request)
}
func NewService(repo Repository) *Service {
	return &Service{repo: repo}
}
