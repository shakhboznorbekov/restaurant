package warehouse_state_history

import (
	"context"
	"github.com/pkg/errors"
	"net/http"
	"restu-backend/foundation/web"
	"restu-backend/internal/auth"
	"restu-backend/internal/pkg/repository/postgresql"
	"restu-backend/internal/service/warehouse_state_history"
	"time"
)

type Repository struct {
	*postgresql.Database
}

func (r Repository) AdminCreate(ctx context.Context, request warehouse_state_history.AdminCreateRequest) (warehouse_state_history.AdminCreateResponse, error) {
	claims, err := r.CheckClaims(ctx, auth.RoleAdmin)
	if err != nil {
		return warehouse_state_history.AdminCreateResponse{}, err
	}

	err = r.ValidateStruct(&request)
	if err != nil {
		return warehouse_state_history.AdminCreateResponse{}, err
	}

	response := warehouse_state_history.AdminCreateResponse{
		Amount:                        request.Amount,
		AveragePrice:                  request.AveragePrice,
		WarehouseStateID:              request.WarehouseStateID,
		WarehouseTransactionProductID: request.WarehouseTransactionProductID,
		CreatedAt:                     time.Now(),
		CreatedBy:                     claims.UserId,
	}

	_, err = r.NewInsert().Model(&response).Exec(ctx)
	if err != nil {
		return warehouse_state_history.AdminCreateResponse{}, web.NewRequestError(errors.Wrap(err, "creating user"), http.StatusBadRequest)
	}

	return response, nil
}

func (r Repository) AdminUpdate(ctx context.Context, request warehouse_state_history.AdminUpdateRequest) error {
	claims, err := r.CheckClaims(ctx, auth.RoleAdmin)
	if err != nil {
		return err
	}

	err = r.ValidateStruct(&request, "ID")
	if err != nil {
		return err
	}
	q := r.NewUpdate().Table("warehouse_state_history").Where("deleted_at IS NULL AND id = ?",
		request.ID)

	if request.Amount != nil {
		q.Set("amount = ?", request.Amount)
	}
	if request.AveragePrice != nil {
		q.Set("average_price = ?", request.AveragePrice)
	}
	if request.WarehouseStateID != nil {
		q.Set("warehouse_state_id = ?", request.WarehouseStateID)
	}
	if request.Amount != nil {
		q.Set("warehouse_transaction_product_id = ?", request.WarehouseTransactionProductID)
	}

	q.Set("updated_at = ?", time.Now())
	q.Set("updated_by = ?", claims.UserId)

	_, err = q.Exec(ctx)
	if err != nil {
		return web.NewRequestError(errors.Wrap(err, "updating food"), http.StatusBadRequest)
	}

	return nil
}

func (r Repository) AdminDeleteTransaction(ctx context.Context, warehouseTransactionProductID int64) error {
	claims, err := r.CheckClaims(ctx, auth.RoleAdmin)
	if err != nil {
		return err
	}

	q := r.NewUpdate().Table("warehouse_state_history").Where("deleted_at IS NULL AND warehouse_transaction_product_id = ?",
		warehouseTransactionProductID)

	q.Set("deleted_at = ?", time.Now())
	q.Set("deleted_by = ?", claims.UserId)

	_, err = q.Exec(ctx)
	if err != nil {
		return web.NewRequestError(errors.Wrap(err, "updating food"), http.StatusBadRequest)
	}

	return nil
}

func (r Repository) BranchCreate(ctx context.Context, request warehouse_state_history.BranchCreateRequest) (warehouse_state_history.BranchCreateResponse, error) {
	claims, err := r.CheckClaims(ctx, auth.RoleBranch)
	if err != nil {
		return warehouse_state_history.BranchCreateResponse{}, err
	}

	err = r.ValidateStruct(&request, "Amount", "AveragePrice", "WarehouseStateID", "WarehouseTransactionID")
	if err != nil {
		return warehouse_state_history.BranchCreateResponse{}, err
	}

	response := warehouse_state_history.BranchCreateResponse{
		Amount:                        request.Amount,
		AveragePrice:                  request.AveragePrice,
		WarehouseStateID:              request.WarehouseStateID,
		WarehouseTransactionProductID: request.WarehouseTransactionProductID,
		CreatedAt:                     time.Now(),
		CreatedBy:                     claims.UserId,
	}

	_, err = r.NewInsert().Model(&response).Exec(ctx)
	if err != nil {
		return warehouse_state_history.BranchCreateResponse{}, web.NewRequestError(errors.Wrap(err, "creating user"), http.StatusNotImplemented)
	}

	return response, nil
}

func (r Repository) BranchUpdate(ctx context.Context, request warehouse_state_history.BranchUpdateRequest) error {
	claims, err := r.CheckClaims(ctx, auth.RoleBranch)
	if err != nil {
		return err
	}

	err = r.ValidateStruct(&request, "ID")
	if err != nil {
		return err
	}
	q := r.NewUpdate().Table("warehouse_state_history").Where("deleted_at IS NULL AND id = ?",
		request.ID)

	if request.Amount != nil {
		q.Set("amount = ?", request.Amount)
	}
	if request.AveragePrice != nil {
		q.Set("average_price = ?", request.AveragePrice)
	}
	if request.WarehouseStateID != nil {
		q.Set("warehouse_state_id = ?", request.WarehouseStateID)
	}
	if request.Amount != nil {
		q.Set("warehouse_transaction_product_id = ?", request.WarehouseTransactionProductID)
	}

	q.Set("updated_at = ?", time.Now())
	q.Set("updated_by = ?", claims.UserId)

	_, err = q.Exec(ctx)
	if err != nil {
		return web.NewRequestError(errors.Wrap(err, "updating food"), http.StatusBadRequest)
	}

	return nil
}

func (r Repository) BranchDeleteTransaction(ctx context.Context, warehouseTransactionProductID int64) error {
	claims, err := r.CheckClaims(ctx, auth.RoleBranch)
	if err != nil {
		return err
	}

	q := r.NewUpdate().Table("warehouse_state_history").Where("deleted_at IS NULL AND warehouse_transaction_product_id = ?",
		warehouseTransactionProductID)

	q.Set("deleted_at = ?", time.Now())
	q.Set("deleted_by = ?", claims.UserId)

	_, err = q.Exec(ctx)
	if err != nil {
		return web.NewRequestError(errors.Wrap(err, "updating food"), http.StatusBadRequest)
	}

	return nil
}

func NewRepository(DB *postgresql.Database) *Repository {
	return &Repository{DB}
}
