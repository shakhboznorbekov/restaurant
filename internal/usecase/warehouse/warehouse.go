package warehouse

import (
	"context"
	"fmt"
	"github.com/pkg/errors"
	"github.com/restaurant/foundation/web"
	"github.com/restaurant/internal/service/warehouse"
	"github.com/restaurant/internal/service/warehouse_state"
	"github.com/restaurant/internal/service/warehouse_state_history"
	"github.com/restaurant/internal/service/warehouse_transaction"
	"github.com/restaurant/internal/service/warehouse_transaction_product"
	"net/http"
)

type UseCase struct {
	warehouses                  Warehouses
	warehouseSate               WarehouseSate
	warehouseTransaction        WarehouseTransaction
	warehouseSateHistory        WarehouseSateHistory
	warehouseTransactionProduct WarehouseTransactionProduct
}

func NewUseCase(
	warehouses Warehouses,
	warehouseSate WarehouseSate,
	warehouseTransaction WarehouseTransaction,
	warehouseSateHistory WarehouseSateHistory,
	warehouseTransactionProduct WarehouseTransactionProduct,
) *UseCase {
	return &UseCase{
		warehouses,
		warehouseSate,
		warehouseTransaction,
		warehouseSateHistory,
		warehouseTransactionProduct,
	}
}

//warehouse

// @admin

func (s UseCase) AdminGetWarehouseList(ctx context.Context, filter warehouse.Filter) ([]warehouse.AdminGetList, int, error) {
	return s.warehouses.AdminGetList(ctx, filter)
}

func (s UseCase) AdminGetWarehouseDetail(ctx context.Context, id int64) (warehouse.AdminGetDetail, error) {
	return s.warehouses.AdminGetDetail(ctx, id)
}

func (s UseCase) AdminCreateWarehouse(ctx context.Context, request warehouse.AdminCreateRequest) (warehouse.AdminCreateResponse, error) {
	return s.warehouses.AdminCreate(ctx, request)
}

func (s UseCase) AdminUpdateWarehouseAll(ctx context.Context, request warehouse.AdminUpdateRequest) error {
	return s.warehouses.AdminUpdateAll(ctx, request)
}

func (s UseCase) AdminUpdateWarehouseColumns(ctx context.Context, request warehouse.AdminUpdateRequest) error {
	return s.warehouses.AdminUpdateColumns(ctx, request)
}

func (s UseCase) AdminDeleteWarehouse(ctx context.Context, id int64) error {
	return s.warehouses.AdminDelete(ctx, id)
}

// @branch

func (s UseCase) BranchGetWarehouseList(ctx context.Context, filter warehouse.Filter) ([]warehouse.BranchGetList, int, error) {
	return s.warehouses.BranchGetList(ctx, filter)
}

func (s UseCase) BranchGetWarehouseDetail(ctx context.Context, id int64) (warehouse.BranchGetDetail, error) {
	return s.warehouses.BranchGetDetail(ctx, id)
}

func (s UseCase) BranchCreateWarehouse(ctx context.Context, request warehouse.BranchCreateRequest) (warehouse.BranchCreateResponse, error) {
	return s.warehouses.BranchCreate(ctx, request)
}

func (s UseCase) BranchUpdateWarehouseAll(ctx context.Context, request warehouse.BranchUpdateRequest) error {
	return s.warehouses.BranchUpdateAll(ctx, request)
}

func (s UseCase) BranchUpdateWarehouseColumns(ctx context.Context, request warehouse.BranchUpdateRequest) error {
	return s.warehouses.BranchUpdateColumns(ctx, request)
}

func (s UseCase) BranchDeleteWarehouse(ctx context.Context, id int64) error {
	return s.warehouses.BranchDelete(ctx, id)
}

// warehouse-state

// @admin

func (s UseCase) AdminGetWarehouseStateByWarehouseIDList(ctx context.Context, warehouseId int64, filter warehouse_state.Filter) ([]warehouse_state.AdminGetByWarehouseIDList, int, error) {
	return s.warehouseSate.AdminGetListByWarehouseID(ctx, warehouseId, filter)
}

// @branch

func (s UseCase) BranchGetWarehouseStateByWarehouseIDList(ctx context.Context, warehouseId int64, filter warehouse_state.Filter) ([]warehouse_state.BranchGetByWarehouseIDList, int, error) {
	return s.warehouseSate.BranchGetListByWarehouseID(ctx, warehouseId, filter)
}

//warehouse-transaction

// @admin

func (s UseCase) AdminGetWarehouseTransactionList(ctx context.Context, filter warehouse_transaction.Filter) ([]warehouse_transaction.AdminGetListResponse, int, error) {
	return s.warehouseTransaction.AdminGetList(ctx, filter)
}

func (s UseCase) AdminGetWarehouseTransactionDetail(ctx context.Context, id int64) (warehouse_transaction.AdminGetDetailByIdResponse, error) {
	return s.warehouseTransaction.AdminGetDetailByID(ctx, id)
}

func (s UseCase) AdminCreateWarehouseTransaction(ctx context.Context, request warehouse_transaction.AdminCreateRequest) (warehouse_transaction.AdminCreateResponse, error) {
	return s.warehouseTransaction.AdminCreate(ctx, request)
}

func (s UseCase) AdminUpdateWarehouseTransactionColumns(ctx context.Context, request warehouse_transaction.AdminUpdateRequest) error {
	return s.warehouseTransaction.AdminUpdateColumn(ctx, request)
}

func (s UseCase) AdminDeleteWarehouseTransaction(ctx context.Context, id int64) error {
	return s.warehouseTransaction.AdminDelete(ctx, id)
}

// @branch

func (s UseCase) BranchGetWarehouseTransactionList(ctx context.Context, filter warehouse_transaction.Filter) ([]warehouse_transaction.BranchGetListResponse, int, error) {
	return s.warehouseTransaction.BranchGetList(ctx, filter)
}

func (s UseCase) BranchGetWarehouseTransactionDetail(ctx context.Context, id int64) (warehouse_transaction.BranchGetDetailByIdResponse, error) {
	return s.warehouseTransaction.BranchGetDetailByID(ctx, id)
}

func (s UseCase) BranchCreateWarehouseTransaction(ctx context.Context, request warehouse_transaction.BranchCreateRequest) (warehouse_transaction.BranchCreateResponse, error) {
	return s.warehouseTransaction.BranchCreate(ctx, request)
}

func (s UseCase) BranchUpdateWarehouseTransactionColumns(ctx context.Context, request warehouse_transaction.BranchUpdateRequest) error {
	return s.warehouseTransaction.BranchUpdateColumn(ctx, request)
}

func (s UseCase) BranchDeleteWarehouseTransaction(ctx context.Context, id int64) error {
	return s.warehouseTransaction.BranchDelete(ctx, id)
}

//warehouse-transaction-product

// @admin

func (s UseCase) AdminGetWarehouseTransactionProductList(ctx context.Context, filter warehouse_transaction_product.Filter, transactionID int64) ([]warehouse_transaction_product.AdminGetListResponse, int, error) {
	return s.warehouseTransactionProduct.AdminGetList(ctx, filter, transactionID)
}

func (s UseCase) AdminGetWarehouseTransactionProductByID(ctx context.Context, id int64) (warehouse_transaction_product.AdminGetDetailByIdResponse, error) {
	return s.warehouseTransactionProduct.AdminGetDetailByID(ctx, id)
}

func (s UseCase) AdminCreateWarehouseTransactionProduct(ctx context.Context, request warehouse_transaction_product.AdminCreateRequest) (warehouse_transaction_product.AdminCreateResponse, error) {
	transaction, err := s.warehouseTransaction.AdminGetDetailByID(ctx, *request.TransactionId)
	if err != nil {
		return warehouse_transaction_product.AdminCreateResponse{}, err
	}

	request.FromWarehouseID = transaction.FromWarehouseID
	request.FromPartnerID = transaction.FromPartnerID
	request.ToWarehouseID = transaction.ToWarehouseID
	request.ToPartnerID = transaction.ToPartnerID

	StateResponse, err := s.warehouseSate.AdminCreate(ctx, request)
	if err != nil {
		return warehouse_transaction_product.AdminCreateResponse{}, err
	}
	response, err := s.warehouseTransactionProduct.AdminCreate(ctx, request)
	if err != nil {
		return warehouse_transaction_product.AdminCreateResponse{}, err
	}
	if StateResponse.OutcomeSate != nil {
		_, err := s.warehouseSateHistory.AdminCreate(ctx, warehouse_state_history.AdminCreateRequest{
			Amount:                        StateResponse.OutcomeSate.Amount,
			AveragePrice:                  StateResponse.OutcomeSate.AveragePrice,
			WarehouseStateID:              &StateResponse.OutcomeSate.ID,
			WarehouseTransactionProductID: &response.ID,
		})
		if err != nil {
			return warehouse_transaction_product.AdminCreateResponse{}, err
		}
	}

	if StateResponse.IncomeSate != nil {
		_, err := s.warehouseSateHistory.AdminCreate(ctx, warehouse_state_history.AdminCreateRequest{
			Amount:                        StateResponse.IncomeSate.Amount,
			AveragePrice:                  StateResponse.IncomeSate.AveragePrice,
			WarehouseStateID:              &StateResponse.IncomeSate.ID,
			WarehouseTransactionProductID: &response.ID,
		})
		if err != nil {
			return warehouse_transaction_product.AdminCreateResponse{}, err
		}
	}

	return response, nil
}

func (s UseCase) AdminUpdateWarehouseTransactionProductColumn(ctx context.Context, request warehouse_transaction_product.AdminUpdateRequest) error {
	transactionProduct, err := s.warehouseTransactionProduct.AdminGetDetailByID(ctx, *request.ID)
	if err != nil {
		return err
	}

	OldTransaction, err := s.warehouseTransaction.AdminGetDetailByID(ctx, *transactionProduct.TransactionId)
	if err != nil {
		return err
	}

	if OldTransaction.FromWarehouseID != nil {
		limit := 1
		LastTransaction, _, err := s.warehouseTransaction.AdminGetList(ctx, warehouse_transaction.Filter{
			WarehouseID: OldTransaction.FromWarehouseID,
			Limit:       &limit,
		})
		if err != nil {
			return err
		}
		if len(LastTransaction) > 0 {
			if LastTransaction[0].ID != OldTransaction.ID {
				return web.NewRequestError(errors.New("this is not the last transaction from warehouse"), http.StatusConflict)
			}
		} else {
			return web.NewRequestError(errors.New("the last transaction from this warehouse was not found"), http.StatusConflict)
		}

	}

	if OldTransaction.ToWarehouseID != nil {
		limit := 1
		LastTransaction, _, err := s.warehouseTransaction.AdminGetList(ctx, warehouse_transaction.Filter{
			WarehouseID: OldTransaction.ToWarehouseID,
			Limit:       &limit,
		})
		if err != nil {
			return err
		}
		if len(LastTransaction) > 0 {
			if LastTransaction[0].ID != OldTransaction.ID {
				return web.NewRequestError(errors.New("this is not the last transaction to warehouse"), http.StatusConflict)
			}
		} else {
			return web.NewRequestError(errors.New("the last transaction to this warehouse was not found"), http.StatusConflict)
		}
	}

	if request.Amount == nil {
		request.Amount = transactionProduct.Amount
	}
	if request.ProductID == nil {
		request.ProductID = transactionProduct.ProductID
	}
	if request.TotalPrice == nil {
		request.TotalPrice = transactionProduct.TotalPrice
	}
	request.FromWarehouseID = OldTransaction.FromWarehouseID
	request.ToWarehouseID = OldTransaction.ToWarehouseID
	request.FromPartnerID = OldTransaction.FromPartnerID
	request.ToPartnerID = OldTransaction.ToPartnerID

	response, err := s.warehouseSate.AdminUpdate(ctx, &request)
	if err != nil {
		return err
	}

	err = s.warehouseTransactionProduct.AdminUpdateColumn(ctx, request)
	if err != nil {
		return err
	}

	fmt.Println(response)

	if response.IncomeSate != nil {
		err = s.warehouseSateHistory.AdminUpdate(ctx, warehouse_state_history.AdminUpdateRequest{
			ID:                            response.IncomeSate.WarehouseStateHistoryID,
			Amount:                        response.IncomeSate.Amount,
			AveragePrice:                  response.IncomeSate.AveragePrice,
			WarehouseStateID:              &response.IncomeSate.ID,
			WarehouseTransactionProductID: request.ID,
		})
		if err != nil {
			return err
		}
	}

	if response.OutcomeSate != nil {
		err = s.warehouseSateHistory.AdminUpdate(ctx, warehouse_state_history.AdminUpdateRequest{
			ID:                            response.OutcomeSate.WarehouseStateHistoryID,
			Amount:                        response.OutcomeSate.Amount,
			AveragePrice:                  response.OutcomeSate.AveragePrice,
			WarehouseStateID:              &response.OutcomeSate.ID,
			WarehouseTransactionProductID: request.ID,
		})
		if err != nil {
			return err
		}
	}

	return nil
}

func (s UseCase) AdminDeleteWarehouseTransactionProduct(ctx context.Context, id int64) error {
	transactionProduct, err := s.warehouseTransactionProduct.AdminGetDetailByID(ctx, id)
	if err != nil {
		return err
	}

	transaction, err := s.warehouseTransaction.AdminGetDetailByID(ctx, *transactionProduct.TransactionId)
	if err != nil {
		return err
	}

	if transaction.FromWarehouseID != nil {
		limit := 1
		LastTransaction, _, err := s.warehouseTransaction.AdminGetList(ctx, warehouse_transaction.Filter{
			WarehouseID: transaction.FromWarehouseID,
			Limit:       &limit,
		})
		if err != nil {
			return err
		}
		if len(LastTransaction) > 0 {
			if LastTransaction[0].ID != transaction.ID {
				return web.NewRequestError(errors.New("this is not the last transaction from warehouse"), http.StatusConflict)
			}
		} else {
			return web.NewRequestError(errors.New("the last transaction from this warehouse was not found"), http.StatusConflict)
		}

	}

	if transaction.ToWarehouseID != nil {
		limit := 1
		LastTransaction, _, err := s.warehouseTransaction.AdminGetList(ctx, warehouse_transaction.Filter{
			WarehouseID: transaction.ToWarehouseID,
			Limit:       &limit,
		})
		if err != nil {
			return err
		}
		if len(LastTransaction) > 0 {
			if LastTransaction[0].ID != transaction.ID {
				return web.NewRequestError(errors.New("this is not the last transaction to warehouse"), http.StatusConflict)
			}
		} else {
			return web.NewRequestError(errors.New("the last transaction to this warehouse was not found"), http.StatusConflict)
		}
	}

	err = s.warehouseSate.AdminDeleteTransaction(ctx, warehouse_state.AdminDeleteTransactionRequest{
		TransactionProductID: &id,
		FromWarehouseID:      transaction.FromWarehouseID,
		ToWarehouseID:        transaction.ToWarehouseID,
	})
	if err != nil {
		return err
	}

	err = s.warehouseTransactionProduct.AdminDelete(ctx, id)
	if err != nil {
		return err
	}

	err = s.warehouseSateHistory.AdminDeleteTransaction(ctx, id)
	if err != nil {
		return err
	}

	return nil
}

// @branch

func (s UseCase) BranchGetWarehouseTransactionProductList(ctx context.Context, filter warehouse_transaction_product.Filter, transactionID int64) ([]warehouse_transaction_product.BranchGetListResponse, int, error) {
	return s.warehouseTransactionProduct.BranchGetList(ctx, filter, transactionID)
}

func (s UseCase) BranchGetWarehouseTransactionProductByID(ctx context.Context, id int64) (warehouse_transaction_product.BranchGetDetailByIdResponse, error) {
	return s.warehouseTransactionProduct.BranchGetDetailByID(ctx, id)
}

func (s UseCase) BranchCreateWarehouseTransactionProduct(ctx context.Context, request warehouse_transaction_product.BranchCreateRequest) (warehouse_transaction_product.BranchCreateResponse, error) {
	transaction, err := s.warehouseTransaction.BranchGetDetailByID(ctx, *request.TransactionId)
	if err != nil {
		return warehouse_transaction_product.BranchCreateResponse{}, err
	}

	request.FromWarehouseID = transaction.FromWarehouseID
	request.FromPartnerID = transaction.FromPartnerID
	request.ToWarehouseID = transaction.ToWarehouseID
	request.ToPartnerID = transaction.ToPartnerID

	StateResponse, err := s.warehouseSate.BranchCreate(ctx, request)
	if err != nil {
		return warehouse_transaction_product.BranchCreateResponse{}, err
	}
	response, err := s.warehouseTransactionProduct.BranchCreate(ctx, request)
	if err != nil {
		return warehouse_transaction_product.BranchCreateResponse{}, err
	}
	if StateResponse.OutcomeSate != nil {
		_, err := s.warehouseSateHistory.BranchCreate(ctx, warehouse_state_history.BranchCreateRequest{
			Amount:                        StateResponse.OutcomeSate.Amount,
			AveragePrice:                  StateResponse.OutcomeSate.AveragePrice,
			WarehouseStateID:              &StateResponse.OutcomeSate.ID,
			WarehouseTransactionProductID: &response.ID,
		})
		if err != nil {
			return warehouse_transaction_product.BranchCreateResponse{}, err
		}
	}

	if StateResponse.IncomeSate != nil {
		_, err := s.warehouseSateHistory.BranchCreate(ctx, warehouse_state_history.BranchCreateRequest{
			Amount:                        StateResponse.IncomeSate.Amount,
			AveragePrice:                  StateResponse.IncomeSate.AveragePrice,
			WarehouseStateID:              &StateResponse.IncomeSate.ID,
			WarehouseTransactionProductID: &response.ID,
		})
		if err != nil {
			return warehouse_transaction_product.BranchCreateResponse{}, err
		}
	}

	return response, nil
}

func (s UseCase) BranchUpdateWarehouseTransactionProductColumn(ctx context.Context, request warehouse_transaction_product.BranchUpdateRequest) error {
	transactionProduct, err := s.warehouseTransactionProduct.BranchGetDetailByID(ctx, *request.ID)
	if err != nil {
		return err
	}

	OldTransaction, err := s.warehouseTransaction.BranchGetDetailByID(ctx, *transactionProduct.TransactionId)
	if err != nil {
		return err
	}

	if OldTransaction.FromWarehouseID != nil {
		limit := 1
		LastTransaction, _, err := s.warehouseTransaction.BranchGetList(ctx, warehouse_transaction.Filter{
			WarehouseID: OldTransaction.FromWarehouseID,
			Limit:       &limit,
		})
		if err != nil {
			return err
		}
		if len(LastTransaction) > 0 {
			if LastTransaction[0].ID != OldTransaction.ID {
				return web.NewRequestError(errors.New("this is not the last transaction from warehouse"), http.StatusConflict)
			}
		} else {
			return web.NewRequestError(errors.New("the last transaction from this warehouse was not found"), http.StatusConflict)
		}

	}

	if OldTransaction.ToWarehouseID != nil {
		limit := 1
		LastTransaction, _, err := s.warehouseTransaction.BranchGetList(ctx, warehouse_transaction.Filter{
			WarehouseID: OldTransaction.ToWarehouseID,
			Limit:       &limit,
		})
		if err != nil {
			return err
		}
		if len(LastTransaction) > 0 {
			if LastTransaction[0].ID != OldTransaction.ID {
				return web.NewRequestError(errors.New("this is not the last transaction to warehouse"), http.StatusConflict)
			}
		} else {
			return web.NewRequestError(errors.New("the last transaction to this warehouse was not found"), http.StatusConflict)
		}
	}

	if request.Amount == nil {
		request.Amount = transactionProduct.Amount
	}
	if request.ProductID == nil {
		request.ProductID = transactionProduct.ProductID
	}
	if request.TotalPrice == nil {
		request.TotalPrice = transactionProduct.TotalPrice
	}
	request.FromWarehouseID = OldTransaction.FromWarehouseID
	request.ToWarehouseID = OldTransaction.ToWarehouseID
	request.FromPartnerID = OldTransaction.FromPartnerID
	request.ToPartnerID = OldTransaction.ToPartnerID

	response, err := s.warehouseSate.BranchUpdate(ctx, &request)
	if err != nil {
		return err
	}

	err = s.warehouseTransactionProduct.BranchUpdateColumn(ctx, request)
	if err != nil {
		return err
	}

	fmt.Println(response)

	if response.IncomeSate != nil {
		err = s.warehouseSateHistory.BranchUpdate(ctx, warehouse_state_history.BranchUpdateRequest{
			ID:                            response.IncomeSate.WarehouseStateHistoryID,
			Amount:                        response.IncomeSate.Amount,
			AveragePrice:                  response.IncomeSate.AveragePrice,
			WarehouseStateID:              &response.IncomeSate.ID,
			WarehouseTransactionProductID: request.ID,
		})
		if err != nil {
			return err
		}
	}

	if response.OutcomeSate != nil {
		err = s.warehouseSateHistory.BranchUpdate(ctx, warehouse_state_history.BranchUpdateRequest{
			ID:                            response.OutcomeSate.WarehouseStateHistoryID,
			Amount:                        response.OutcomeSate.Amount,
			AveragePrice:                  response.OutcomeSate.AveragePrice,
			WarehouseStateID:              &response.OutcomeSate.ID,
			WarehouseTransactionProductID: request.ID,
		})
		if err != nil {
			return err
		}
	}

	return nil
}

func (s UseCase) BranchDeleteWarehouseTransactionProduct(ctx context.Context, id int64) error {
	transactionProduct, err := s.warehouseTransactionProduct.BranchGetDetailByID(ctx, id)
	if err != nil {
		return err
	}

	transaction, err := s.warehouseTransaction.BranchGetDetailByID(ctx, *transactionProduct.TransactionId)
	if err != nil {
		return err
	}

	if transaction.FromWarehouseID != nil {
		limit := 1
		LastTransaction, _, err := s.warehouseTransaction.BranchGetList(ctx, warehouse_transaction.Filter{
			WarehouseID: transaction.FromWarehouseID,
			Limit:       &limit,
		})
		if err != nil {
			return err
		}
		if len(LastTransaction) > 0 {
			if LastTransaction[0].ID != transaction.ID {
				return web.NewRequestError(errors.New("this is not the last transaction from warehouse"), http.StatusConflict)
			}
		} else {
			return web.NewRequestError(errors.New("the last transaction from this warehouse was not found"), http.StatusConflict)
		}

	}

	if transaction.ToWarehouseID != nil {
		limit := 1
		LastTransaction, _, err := s.warehouseTransaction.BranchGetList(ctx, warehouse_transaction.Filter{
			WarehouseID: transaction.ToWarehouseID,
			Limit:       &limit,
		})
		if err != nil {
			return err
		}
		if len(LastTransaction) > 0 {
			if LastTransaction[0].ID != transaction.ID {
				return web.NewRequestError(errors.New("this is not the last transaction to warehouse"), http.StatusConflict)
			}
		} else {
			return web.NewRequestError(errors.New("the last transaction to this warehouse was not found"), http.StatusConflict)
		}
	}

	err = s.warehouseSate.BranchDeleteTransaction(ctx, warehouse_state.BranchDeleteTransactionRequest{
		TransactionProductID: &id,
		FromWarehouseID:      transaction.FromWarehouseID,
		ToWarehouseID:        transaction.ToWarehouseID,
	})
	if err != nil {
		return err
	}

	err = s.warehouseTransactionProduct.BranchDelete(ctx, id)
	if err != nil {
		return err
	}

	err = s.warehouseSateHistory.BranchDeleteTransaction(ctx, id)
	if err != nil {
		return err
	}

	return nil
}
