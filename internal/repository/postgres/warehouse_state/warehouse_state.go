package warehouse_state

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/pkg/errors"
	"net/http"
	"restu-backend/foundation/web"
	"restu-backend/internal/auth"
	"restu-backend/internal/entity"
	"restu-backend/internal/pkg/repository/postgresql"
	"restu-backend/internal/repository/postgres"
	"restu-backend/internal/service/warehouse_state"
	"restu-backend/internal/service/warehouse_transaction_product"
	"time"
)

type Repository struct {
	*postgresql.Database
}

// #admin

func (r Repository) AdminGetListByWarehouseID(ctx context.Context, warehouseId int64, filter warehouse_state.Filter) ([]warehouse_state.AdminGetByWarehouseIDList, int, error) {
	_, err := r.CheckClaims(ctx, auth.RoleAdmin)
	if err != nil {
		return nil, 0, err
	}

	whereQuery := fmt.Sprintf(`WHERE ws.deleted_at IS NULL AND p.deleted_at IS NULL AND ws.warehouse_id = %d`, warehouseId)

	var limitQuery, offsetQuery string
	if filter.Limit != nil {
		limitQuery = fmt.Sprintf(" LIMIT '%d'", *filter.Limit)
	}
	if filter.Offset != nil {
		offsetQuery = fmt.Sprintf(" OFFSET '%d'", *filter.Offset)
	}

	query := fmt.Sprintf(`
		SELECT
			ws.id,
			ws.amount,
			ws.average_price,
			ws.product_id,
			p.name
		FROM
		    warehouse_state as ws
		LEFT JOIN products p ON p.id = ws.product_id
		%s %s %s
	`, whereQuery, limitQuery, offsetQuery)

	rows, err := r.QueryContext(ctx, query)
	if err != nil {
		return nil, 0, web.NewRequestError(errors.Wrap(err, "select user"), http.StatusInternalServerError)
	}

	var list []warehouse_state.AdminGetByWarehouseIDList

	for rows.Next() {
		var detail warehouse_state.AdminGetByWarehouseIDList
		if err = rows.Scan(&detail.ID, &detail.Amount, &detail.AveragePrice, &detail.ProductID, &detail.Product); err != nil {
			return nil, 0, web.NewRequestError(errors.Wrap(err, "scanning warehouse"), http.StatusBadRequest)
		}

		list = append(list, detail)
	}

	countQuery := fmt.Sprintf(`
		SELECT
			count(ws.id)
		FROM
		    warehouse_state as ws
		LEFT JOIN products p ON p.id = ws.product_id
		%s
	`, whereQuery)

	countRows, err := r.QueryContext(ctx, countQuery)
	if err == sql.ErrNoRows {
		return nil, 0, web.NewRequestError(postgres.ErrNotFound, http.StatusNotFound)
	}
	if err != nil {
		return nil, 0, web.NewRequestError(errors.Wrap(err, "selecting users"), http.StatusBadRequest)
	}

	count := 0

	for countRows.Next() {
		if err = countRows.Scan(&count); err != nil {
			return nil, 0, web.NewRequestError(errors.Wrap(err, "scanning user count"), http.StatusBadRequest)
		}
	}

	return list, count, nil
}

func (r Repository) AdminCreate(ctx context.Context, request warehouse_transaction_product.AdminCreateRequest) (warehouse_state.AdminCreateResponse, error) {
	claims, err := r.CheckClaims(ctx, auth.RoleAdmin)
	if err != nil {
		return warehouse_state.AdminCreateResponse{}, err
	}

	err = r.ValidateStruct(&request, "Amount", "ProductID")
	if err != nil {
		return warehouse_state.AdminCreateResponse{}, err
	}

	response := warehouse_state.AdminCreateResponse{}

	if request.FromWarehouseID != nil {
		res := warehouse_state.AdminCreate{}
		state, err := r.getByProductIDAndWarehouseID(ctx, *request.ProductID, *request.FromWarehouseID)
		if err != nil {
			if !errors.Is(err, sql.ErrNoRows) {
				return warehouse_state.AdminCreateResponse{}, err
			} else {
				return warehouse_state.AdminCreateResponse{}, web.NewRequestError(errors.New("the product is not available in the warehouse"), http.StatusBadRequest)
			}
		}

		if *state.Amount < *request.Amount {
			return warehouse_state.AdminCreateResponse{}, web.NewRequestError(errors.New("there is not enough product in the warehouse"), http.StatusBadRequest)
		}
		q := r.NewUpdate().Table("warehouse_state").Where("product_id = ? AND warehouse_id = ?", *request.ProductID, *request.FromWarehouseID)
		q.Set("amount = amount - ?", request.Amount)
		q.Set("updated_at = ?", time.Now())
		q.Set("updated_by = ?", claims.UserId)

		_, err = q.Exec(ctx)
		if err != nil {
			return warehouse_state.AdminCreateResponse{}, web.NewRequestError(errors.Wrap(err, "update warehouse_state"), http.StatusInternalServerError)
		}

		totalPrice := *state.AveragePrice * (*state.Amount)
		request.TotalPrice = &totalPrice

		res = warehouse_state.AdminCreate{
			ID:           state.ID,
			Amount:       state.Amount,
			ProductID:    state.ProductID,
			AveragePrice: state.AveragePrice,
			WarehouseID:  state.WarehouseID,
			CreatedAt:    *state.CreatedAt,
			CreatedBy:    *state.CreatedBy,
		}

		response.OutcomeSate = &res
	}

	err = r.ValidateStruct(&request, "TotalPrice")
	if err != nil {
		return warehouse_state.AdminCreateResponse{}, err
	}

	if request.ToWarehouseID != nil {
		res := warehouse_state.AdminCreate{}
		state, err := r.getByProductIDAndWarehouseID(ctx, *request.ProductID, *request.ToWarehouseID)
		if err != nil {
			if !errors.Is(err, sql.ErrNoRows) {
				return warehouse_state.AdminCreateResponse{}, err
			} else {
				averagePrice := *request.TotalPrice / (*request.Amount)
				resQ := warehouse_state.AdminCreate{
					Amount:       request.Amount,
					AveragePrice: &averagePrice,
					ProductID:    request.ProductID,
					WarehouseID:  request.ToWarehouseID,
					CreatedBy:    claims.UserId,
					CreatedAt:    time.Now(),
				}
				_, err = r.NewInsert().Model(&resQ).Exec(ctx)
				if err != nil {
					return warehouse_state.AdminCreateResponse{}, web.NewRequestError(errors.Wrap(err, "creating warehouse_state"), http.StatusBadRequest)
				}
				res = resQ
				var (
					amount, price float64
				)
				res.Amount = &amount
				res.AveragePrice = &price
				response.IncomeSate = &res
				return response, nil
			}
		}

		q := r.NewUpdate().Table("warehouse_state").Where("product_id = ? AND warehouse_id = ?", *request.ProductID, *request.ToWarehouseID)
		q.Set("amount = amount + ?", request.Amount)
		q.Set("average_price = (average_price * amount + ?)/(amount + ?)", request.TotalPrice, request.Amount)
		q.Set("updated_at = ?", time.Now())
		q.Set("updated_by = ?", claims.UserId)

		_, err = q.Exec(ctx)
		if err != nil {
			return warehouse_state.AdminCreateResponse{}, web.NewRequestError(errors.Wrap(err, "update warehouse_state"), http.StatusInternalServerError)
		}

		res = warehouse_state.AdminCreate{
			ID:           state.ID,
			Amount:       state.Amount,
			ProductID:    state.ProductID,
			AveragePrice: state.AveragePrice,
			WarehouseID:  state.WarehouseID,
			CreatedAt:    *state.CreatedAt,
			CreatedBy:    *state.CreatedBy,
		}

		response.IncomeSate = &res
	}

	return response, nil
}

func (r Repository) AdminUpdate(ctx context.Context, request *warehouse_transaction_product.AdminUpdateRequest) (response warehouse_state.AdminUpdateResponse, err error) {
	claims, err := r.CheckClaims(ctx, auth.RoleAdmin)
	if err != nil {
		return warehouse_state.AdminUpdateResponse{}, err
	}
	response = warehouse_state.AdminUpdateResponse{}
	err = r.ValidateStruct(request, "Amount", "ProductID", "TotalPrice")
	if err != nil {
		return warehouse_state.AdminUpdateResponse{}, err
	}

	tx, err := r.BeginTx(ctx, nil)
	if err != nil {
		return warehouse_state.AdminUpdateResponse{}, err
	}

	defer func() {
		if err != nil {
			_ = tx.Rollback()
			return
		}

		_ = tx.Commit()
	}()

	if request.FromWarehouseID != nil {
		stateHistory, er := r.getOldStateByTransactionIDAndWarehouseID(ctx, *request.ID, *request.FromWarehouseID)
		if er != nil {
			if !errors.Is(er, sql.ErrNoRows) {
				err = er
				return warehouse_state.AdminUpdateResponse{}, err
			} else {
				err = web.NewRequestError(errors.New("the product is not available in the warehouse"), http.StatusNotImplemented)
				return warehouse_state.AdminUpdateResponse{}, err
			}
		}
		if *stateHistory.ProductID == *request.ProductID {
			if *stateHistory.Amount < *request.Amount {
				err = web.NewRequestError(errors.New("there is not enough product in the warehouse"), http.StatusBadRequest)
				return warehouse_state.AdminUpdateResponse{}, err
			}

			q := tx.NewUpdate().Table("warehouse_state").Where("product_id = ? AND warehouse_id = ?", *request.ProductID, *request.FromWarehouseID)
			q.Set("amount = ? - ?", stateHistory.Amount, request.Amount)
			q.Set("updated_at = ?", time.Now())
			q.Set("updated_by = ?", claims.UserId)

			_, err = q.Exec(ctx)
			if err != nil {
				return warehouse_state.AdminUpdateResponse{}, web.NewRequestError(errors.Wrap(err, "update warehouse_state"), http.StatusInternalServerError)
			}
		} else {
			q1 := tx.NewUpdate().Table("warehouse_state").Where("product_id = ? AND warehouse_id = ?", *stateHistory.ProductID, *request.FromWarehouseID)
			q1.Set("amount = ?", stateHistory.Amount)
			q1.Set("average_price = ?", stateHistory.AveragePrice)
			q1.Set("updated_at = ?", time.Now())
			q1.Set("updated_by = ?", claims.UserId)

			_, err = q1.Exec(ctx)
			if err != nil {
				return warehouse_state.AdminUpdateResponse{}, web.NewRequestError(errors.Wrap(err, "update warehouse_state"), http.StatusInternalServerError)
			}

			state, er := r.getByProductIDAndWarehouseID(ctx, *request.ProductID, *request.FromWarehouseID)
			if er != nil {
				if !errors.Is(er, sql.ErrNoRows) {
					err = er
					return warehouse_state.AdminUpdateResponse{}, err
				} else {
					err = web.NewRequestError(errors.New("the product is not available in the warehouse"), http.StatusBadRequest)
					return warehouse_state.AdminUpdateResponse{}, err
				}
			}

			if *state.Amount < *request.Amount {
				err = web.NewRequestError(errors.New("there is not enough product in the warehouse"), http.StatusBadRequest)
				return warehouse_state.AdminUpdateResponse{}, err
			}
			q := r.NewUpdate().Table("warehouse_state").Where("product_id = ? AND warehouse_id = ?", *request.ProductID, *request.FromWarehouseID)
			q.Set("amount = amount - ?", request.Amount)
			q.Set("updated_at = ?", time.Now())
			q.Set("updated_by = ?", claims.UserId)

			_, err = q.Exec(ctx)
			if err != nil {
				return warehouse_state.AdminUpdateResponse{}, web.NewRequestError(errors.Wrap(err, "update warehouse_state"), http.StatusInternalServerError)
			}

			res := warehouse_state.AdminUpdate{
				ID:                      state.ID,
				Amount:                  state.Amount,
				ProductID:               state.ProductID,
				AveragePrice:            state.AveragePrice,
				WarehouseID:             state.WarehouseID,
				WarehouseStateHistoryID: &stateHistory.ID,
			}

			response.OutcomeSate = &res
			totalPrice := *state.Amount * (*state.AveragePrice)
			request.TotalPrice = &totalPrice
		}
	}

	if request.ToWarehouseID != nil {
		stateHistory, er := r.getOldStateByTransactionIDAndWarehouseID(ctx, *request.ID, *request.ToWarehouseID)
		if er != nil {
			if !errors.Is(er, sql.ErrNoRows) {
				err = er
				return warehouse_state.AdminUpdateResponse{}, err
			} else {
				err = web.NewRequestError(errors.New("the product is not available in the warehouse"), http.StatusNotImplemented)
				return warehouse_state.AdminUpdateResponse{}, err
			}
		}
		if *stateHistory.ProductID == *request.ProductID {
			q := tx.NewUpdate().Table("warehouse_state").Where("product_id = ? AND warehouse_id = ?", *request.ProductID, *request.ToWarehouseID)
			q.Set("amount = ? + ?", stateHistory.Amount, request.Amount)
			q.Set("average_price = (? * ? + ?)/(? + ?)", stateHistory.AveragePrice, stateHistory.Amount, request.TotalPrice, stateHistory.Amount, request.Amount)
			q.Set("updated_at = ?", time.Now())
			q.Set("updated_by = ?", claims.UserId)

			_, err = q.Exec(ctx)
			if err != nil {
				return warehouse_state.AdminUpdateResponse{}, web.NewRequestError(errors.Wrap(err, "update warehouse_state"), http.StatusInternalServerError)
			}
		} else {
			fmt.Println(stateHistory.Amount)
			q1 := tx.NewUpdate().Table("warehouse_state").Where("product_id = ? AND warehouse_id = ?", *stateHistory.ProductID, *request.ToWarehouseID)
			q1.Set("amount = ?", stateHistory.Amount)
			q1.Set("average_price = ?", stateHistory.AveragePrice)
			q1.Set("updated_at = ?", time.Now())
			q1.Set("updated_by = ?", claims.UserId)
			_, err = q1.Exec(ctx)
			if err != nil {
				return warehouse_state.AdminUpdateResponse{}, web.NewRequestError(errors.Wrap(err, "update warehouse_state"), http.StatusInternalServerError)
			}

			fmt.Println(333)

			state, er := r.getByProductIDAndWarehouseID(ctx, *request.ProductID, *request.ToWarehouseID)
			if er != nil {
				if !errors.Is(er, sql.ErrNoRows) {
					err = er
					return warehouse_state.AdminUpdateResponse{}, err
				} else {
					averagePrice := *request.TotalPrice / (*request.Amount)
					resQ := warehouse_state.AdminCreate{
						Amount:       request.Amount,
						AveragePrice: &averagePrice,
						ProductID:    request.ProductID,
						WarehouseID:  request.ToWarehouseID,
						CreatedBy:    claims.UserId,
						CreatedAt:    time.Now(),
					}
					_, err = r.NewInsert().Model(&resQ).Exec(ctx)
					if err != nil {
						return warehouse_state.AdminUpdateResponse{}, web.NewRequestError(errors.Wrap(err, "creating warehouse_state"), http.StatusBadRequest)
					}
					var (
						amount, price float64
					)

					res := warehouse_state.AdminUpdate{
						ID:                      resQ.ID,
						Amount:                  &amount,
						ProductID:               resQ.ProductID,
						AveragePrice:            &price,
						WarehouseID:             resQ.WarehouseID,
						WarehouseStateHistoryID: &stateHistory.ID,
					}
					response.IncomeSate = &res
					return response, nil
				}
			}

			q := r.NewUpdate().Table("warehouse_state").Where("product_id = ? AND warehouse_id = ?", *request.ProductID, *request.ToWarehouseID)
			q.Set("amount = amount + ?", request.Amount)
			q.Set("average_price = (average_price * amount + ?)/(amount + ?)", request.TotalPrice, request.Amount)
			q.Set("updated_at = ?", time.Now())
			q.Set("updated_by = ?", claims.UserId)

			_, err = q.Exec(ctx)
			if err != nil {
				return warehouse_state.AdminUpdateResponse{}, web.NewRequestError(errors.Wrap(err, "update warehouse_state"), http.StatusInternalServerError)
			}

			res := warehouse_state.AdminUpdate{
				ID:                      state.ID,
				Amount:                  state.Amount,
				ProductID:               state.ProductID,
				AveragePrice:            state.AveragePrice,
				WarehouseID:             state.WarehouseID,
				WarehouseStateHistoryID: &stateHistory.ID,
			}

			response.IncomeSate = &res
		}
	}

	return response, nil
}

func (r Repository) AdminDeleteTransaction(ctx context.Context, request warehouse_state.AdminDeleteTransactionRequest) error {
	claims, err := r.CheckClaims(ctx, auth.RoleAdmin)
	if err != nil {
		return err
	}
	err = r.ValidateStruct(request, "Amount", "ProductID", "TotalPrice")
	if err != nil {
		return err
	}

	tx, err := r.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	defer func() {
		if err != nil {
			_ = tx.Rollback()
			return
		}

		_ = tx.Commit()
	}()

	if request.FromWarehouseID != nil {
		stateHistory, er := r.getOldStateByTransactionIDAndWarehouseID(ctx, *request.TransactionProductID, *request.FromWarehouseID)
		if er != nil {
			if !errors.Is(er, sql.ErrNoRows) {
				err = er
				return err
			} else {
				err = web.NewRequestError(errors.New("the product is not available in the warehouse"), http.StatusNotImplemented)
				return err
			}
		}

		q1 := tx.NewUpdate().Table("warehouse_state").Where("product_id = ? AND warehouse_id = ?", *stateHistory.ProductID, *request.FromWarehouseID)
		q1.Set("amount = ?", stateHistory.Amount)
		q1.Set("average_price = ?", stateHistory.AveragePrice)
		q1.Set("updated_at = ?", time.Now())
		q1.Set("updated_by = ?", claims.UserId)

		_, err = q1.Exec(ctx)
		if err != nil {
			return web.NewRequestError(errors.Wrap(err, "update warehouse_state"), http.StatusInternalServerError)
		}
	}

	if request.ToWarehouseID != nil {
		stateHistory, er := r.getOldStateByTransactionIDAndWarehouseID(ctx, *request.TransactionProductID, *request.ToWarehouseID)
		if er != nil {
			if !errors.Is(er, sql.ErrNoRows) {
				err = er
				return err
			} else {
				err = web.NewRequestError(errors.New("the product is not available in the warehouse"), http.StatusNotImplemented)
				return err
			}
		}
		fmt.Println(stateHistory.Amount)
		q1 := tx.NewUpdate().Table("warehouse_state").Where("product_id = ? AND warehouse_id = ?", *stateHistory.ProductID, *request.ToWarehouseID)
		q1.Set("amount = ?", stateHistory.Amount)
		q1.Set("average_price = ?", stateHistory.AveragePrice)
		q1.Set("updated_at = ?", time.Now())
		q1.Set("updated_by = ?", claims.UserId)
		_, err = q1.Exec(ctx)
		if err != nil {
			return web.NewRequestError(errors.Wrap(err, "update warehouse_state"), http.StatusInternalServerError)
		}
	}

	return nil
}

// #branch

func (r Repository) BranchGetListByWarehouseID(ctx context.Context, warehouseId int64, filter warehouse_state.Filter) ([]warehouse_state.BranchGetByWarehouseIDList, int, error) {
	_, err := r.CheckClaims(ctx, auth.RoleBranch)
	if err != nil {
		return nil, 0, err
	}

	whereQuery := fmt.Sprintf(`WHERE ws.deleted_at IS NULL AND p.deleted_at IS NULL AND ws.warehouse_id = %d`, warehouseId)

	var limitQuery, offsetQuery string
	if filter.Limit != nil {
		limitQuery = fmt.Sprintf(" LIMIT '%d'", *filter.Limit)
	}
	if filter.Offset != nil {
		offsetQuery = fmt.Sprintf(" OFFSET '%d'", *filter.Offset)
	}

	query := fmt.Sprintf(`
		SELECT
			ws.id,
			ws.amount,
			ws.average_price,
			ws.product_id,
			p.name
		FROM
		    warehouse_state as ws
		LEFT JOIN products p ON p.id = ws.product_id
		%s %s %s
	`, whereQuery, limitQuery, offsetQuery)

	rows, err := r.QueryContext(ctx, query)
	if err != nil {
		return nil, 0, web.NewRequestError(errors.Wrap(err, "select user"), http.StatusInternalServerError)
	}

	var list []warehouse_state.BranchGetByWarehouseIDList

	for rows.Next() {
		var detail warehouse_state.BranchGetByWarehouseIDList
		if err = rows.Scan(&detail.ID, &detail.Amount, &detail.AveragePrice, &detail.ProductID, &detail.Product); err != nil {
			return nil, 0, web.NewRequestError(errors.Wrap(err, "scanning warehouse"), http.StatusBadRequest)
		}

		list = append(list, detail)
	}

	countQuery := fmt.Sprintf(`
		SELECT
			count(ws.id)
		FROM
		    warehouse_state as ws
		LEFT JOIN products p ON p.id = ws.product_id
		%s
	`, whereQuery)

	countRows, err := r.QueryContext(ctx, countQuery)
	if err == sql.ErrNoRows {
		return nil, 0, web.NewRequestError(postgres.ErrNotFound, http.StatusNotFound)
	}
	if err != nil {
		return nil, 0, web.NewRequestError(errors.Wrap(err, "selecting users"), http.StatusBadRequest)
	}

	count := 0

	for countRows.Next() {
		if err = countRows.Scan(&count); err != nil {
			return nil, 0, web.NewRequestError(errors.Wrap(err, "scanning user count"), http.StatusBadRequest)
		}
	}

	return list, count, nil
}

func (r Repository) BranchCreate(ctx context.Context, request warehouse_transaction_product.BranchCreateRequest) (warehouse_state.BranchCreateResponse, error) {
	claims, err := r.CheckClaims(ctx, auth.RoleBranch)
	if err != nil {
		return warehouse_state.BranchCreateResponse{}, err
	}

	err = r.ValidateStruct(&request, "Amount", "ProductID")
	if err != nil {
		return warehouse_state.BranchCreateResponse{}, err
	}

	response := warehouse_state.BranchCreateResponse{}

	if request.FromWarehouseID != nil {
		res := warehouse_state.BranchCreate{}
		state, err := r.getByProductIDAndWarehouseID(ctx, *request.ProductID, *request.FromWarehouseID)
		if err != nil {
			if !errors.Is(err, sql.ErrNoRows) {
				return warehouse_state.BranchCreateResponse{}, err
			} else {
				return warehouse_state.BranchCreateResponse{}, web.NewRequestError(errors.New("the product is not available in the warehouse"), http.StatusBadRequest)
			}
		}

		if *state.Amount < *request.Amount {
			return warehouse_state.BranchCreateResponse{}, web.NewRequestError(errors.New("there is not enough product in the warehouse"), http.StatusBadRequest)
		}
		q := r.NewUpdate().Table("warehouse_state").Where("product_id = ? AND warehouse_id = ?", *request.ProductID, *request.FromWarehouseID)
		q.Set("amount = amount - ?", request.Amount)
		q.Set("updated_at = ?", time.Now())
		q.Set("updated_by = ?", claims.UserId)

		_, err = q.Exec(ctx)
		if err != nil {
			return warehouse_state.BranchCreateResponse{}, web.NewRequestError(errors.Wrap(err, "update warehouse_state"), http.StatusInternalServerError)
		}

		totalPrice := *state.AveragePrice * (*state.Amount)
		request.TotalPrice = &totalPrice

		res = warehouse_state.BranchCreate{
			ID:           state.ID,
			Amount:       state.Amount,
			ProductID:    state.ProductID,
			AveragePrice: state.AveragePrice,
			WarehouseID:  state.WarehouseID,
			CreatedAt:    *state.CreatedAt,
			CreatedBy:    *state.CreatedBy,
		}

		response.OutcomeSate = &res
	}

	err = r.ValidateStruct(&request, "TotalPrice")
	if err != nil {
		return warehouse_state.BranchCreateResponse{}, err
	}

	if request.ToWarehouseID != nil {
		res := warehouse_state.BranchCreate{}
		state, err := r.getByProductIDAndWarehouseID(ctx, *request.ProductID, *request.ToWarehouseID)
		if err != nil {
			if !errors.Is(err, sql.ErrNoRows) {
				return warehouse_state.BranchCreateResponse{}, err
			} else {
				averagePrice := *request.TotalPrice / (*request.Amount)
				resQ := warehouse_state.BranchCreate{
					Amount:       request.Amount,
					AveragePrice: &averagePrice,
					ProductID:    request.ProductID,
					WarehouseID:  request.ToWarehouseID,
					CreatedBy:    claims.UserId,
					CreatedAt:    time.Now(),
				}
				_, err = r.NewInsert().Model(&resQ).Exec(ctx)
				if err != nil {
					return warehouse_state.BranchCreateResponse{}, web.NewRequestError(errors.Wrap(err, "creating warehouse_state"), http.StatusBadRequest)
				}
				res = resQ
				var (
					amount, price float64
				)
				res.Amount = &amount
				res.AveragePrice = &price
				response.IncomeSate = &res
				return response, nil
			}
		}

		q := r.NewUpdate().Table("warehouse_state").Where("product_id = ? AND warehouse_id = ?", *request.ProductID, *request.ToWarehouseID)
		q.Set("amount = amount + ?", request.Amount)
		q.Set("average_price = (average_price * amount + ?)/(amount + ?)", request.TotalPrice, request.Amount)
		q.Set("updated_at = ?", time.Now())
		q.Set("updated_by = ?", claims.UserId)

		_, err = q.Exec(ctx)
		if err != nil {
			return warehouse_state.BranchCreateResponse{}, web.NewRequestError(errors.Wrap(err, "update warehouse_state"), http.StatusInternalServerError)
		}

		res = warehouse_state.BranchCreate{
			ID:           state.ID,
			Amount:       state.Amount,
			ProductID:    state.ProductID,
			AveragePrice: state.AveragePrice,
			WarehouseID:  state.WarehouseID,
			CreatedAt:    *state.CreatedAt,
			CreatedBy:    *state.CreatedBy,
		}

		response.OutcomeSate = &res
	}

	return response, nil
}

func (r Repository) BranchUpdate(ctx context.Context, request *warehouse_transaction_product.BranchUpdateRequest) (response warehouse_state.BranchUpdateResponse, err error) {
	claims, err := r.CheckClaims(ctx, auth.RoleBranch)
	if err != nil {
		return warehouse_state.BranchUpdateResponse{}, err
	}
	response = warehouse_state.BranchUpdateResponse{}
	err = r.ValidateStruct(request, "Amount", "ProductID", "TotalPrice")
	if err != nil {
		return warehouse_state.BranchUpdateResponse{}, err
	}

	tx, err := r.BeginTx(ctx, nil)
	if err != nil {
		return warehouse_state.BranchUpdateResponse{}, err
	}

	defer func() {
		if err != nil {
			_ = tx.Rollback()
			return
		}

		_ = tx.Commit()
	}()

	if request.FromWarehouseID != nil {
		stateHistory, er := r.getOldStateByTransactionIDAndWarehouseID(ctx, *request.ID, *request.FromWarehouseID)
		if er != nil {
			if !errors.Is(er, sql.ErrNoRows) {
				err = er
				return warehouse_state.BranchUpdateResponse{}, err
			} else {
				err = web.NewRequestError(errors.New("the product is not available in the warehouse"), http.StatusNotImplemented)
				return warehouse_state.BranchUpdateResponse{}, err
			}
		}
		if *stateHistory.ProductID == *request.ProductID {
			if *stateHistory.Amount < *request.Amount {
				err = web.NewRequestError(errors.New("there is not enough product in the warehouse"), http.StatusBadRequest)
				return warehouse_state.BranchUpdateResponse{}, err
			}

			q := tx.NewUpdate().Table("warehouse_state").Where("product_id = ? AND warehouse_id = ?", *request.ProductID, *request.FromWarehouseID)
			q.Set("amount = ? - ?", stateHistory.Amount, request.Amount)
			q.Set("updated_at = ?", time.Now())
			q.Set("updated_by = ?", claims.UserId)

			_, err = q.Exec(ctx)
			if err != nil {
				return warehouse_state.BranchUpdateResponse{}, web.NewRequestError(errors.Wrap(err, "update warehouse_state"), http.StatusInternalServerError)
			}
		} else {
			q1 := tx.NewUpdate().Table("warehouse_state").Where("product_id = ? AND warehouse_id = ?", *stateHistory.ProductID, *request.FromWarehouseID)
			q1.Set("amount = ?", stateHistory.Amount)
			q1.Set("average_price = ?", stateHistory.AveragePrice)
			q1.Set("updated_at = ?", time.Now())
			q1.Set("updated_by = ?", claims.UserId)

			_, err = q1.Exec(ctx)
			if err != nil {
				return warehouse_state.BranchUpdateResponse{}, web.NewRequestError(errors.Wrap(err, "update warehouse_state"), http.StatusInternalServerError)
			}

			state, er := r.getByProductIDAndWarehouseID(ctx, *request.ProductID, *request.FromWarehouseID)
			if er != nil {
				if !errors.Is(er, sql.ErrNoRows) {
					err = er
					return warehouse_state.BranchUpdateResponse{}, err
				} else {
					err = web.NewRequestError(errors.New("the product is not available in the warehouse"), http.StatusBadRequest)
					return warehouse_state.BranchUpdateResponse{}, err
				}
			}

			if *state.Amount < *request.Amount {
				err = web.NewRequestError(errors.New("there is not enough product in the warehouse"), http.StatusBadRequest)
				return warehouse_state.BranchUpdateResponse{}, err
			}
			q := r.NewUpdate().Table("warehouse_state").Where("product_id = ? AND warehouse_id = ?", *request.ProductID, *request.FromWarehouseID)
			q.Set("amount = amount - ?", request.Amount)
			q.Set("updated_at = ?", time.Now())
			q.Set("updated_by = ?", claims.UserId)

			_, err = q.Exec(ctx)
			if err != nil {
				return warehouse_state.BranchUpdateResponse{}, web.NewRequestError(errors.Wrap(err, "update warehouse_state"), http.StatusInternalServerError)
			}

			res := warehouse_state.BranchUpdate{
				ID:                      state.ID,
				Amount:                  state.Amount,
				ProductID:               state.ProductID,
				AveragePrice:            state.AveragePrice,
				WarehouseID:             state.WarehouseID,
				WarehouseStateHistoryID: &stateHistory.ID,
			}

			response.OutcomeSate = &res
			totalPrice := *state.Amount * (*state.AveragePrice)
			request.TotalPrice = &totalPrice
		}
	}

	if request.ToWarehouseID != nil {
		stateHistory, er := r.getOldStateByTransactionIDAndWarehouseID(ctx, *request.ID, *request.ToWarehouseID)
		if er != nil {
			if !errors.Is(er, sql.ErrNoRows) {
				err = er
				return warehouse_state.BranchUpdateResponse{}, err
			} else {
				err = web.NewRequestError(errors.New("the product is not available in the warehouse"), http.StatusNotImplemented)
				return warehouse_state.BranchUpdateResponse{}, err
			}
		}
		if *stateHistory.ProductID == *request.ProductID {
			q := tx.NewUpdate().Table("warehouse_state").Where("product_id = ? AND warehouse_id = ?", *request.ProductID, *request.ToWarehouseID)
			q.Set("amount = ? + ?", stateHistory.Amount, request.Amount)
			q.Set("average_price = (? * ? + ?)/(? + ?)", stateHistory.AveragePrice, stateHistory.Amount, request.TotalPrice, stateHistory.Amount, request.Amount)
			q.Set("updated_at = ?", time.Now())
			q.Set("updated_by = ?", claims.UserId)

			_, err = q.Exec(ctx)
			if err != nil {
				return warehouse_state.BranchUpdateResponse{}, web.NewRequestError(errors.Wrap(err, "update warehouse_state"), http.StatusInternalServerError)
			}
		} else {
			fmt.Println(stateHistory.Amount)
			q1 := tx.NewUpdate().Table("warehouse_state").Where("product_id = ? AND warehouse_id = ?", *stateHistory.ProductID, *request.ToWarehouseID)
			q1.Set("amount = ?", stateHistory.Amount)
			q1.Set("average_price = ?", stateHistory.AveragePrice)
			q1.Set("updated_at = ?", time.Now())
			q1.Set("updated_by = ?", claims.UserId)
			_, err = q1.Exec(ctx)
			if err != nil {
				return warehouse_state.BranchUpdateResponse{}, web.NewRequestError(errors.Wrap(err, "update warehouse_state"), http.StatusInternalServerError)
			}

			fmt.Println(333)

			state, er := r.getByProductIDAndWarehouseID(ctx, *request.ProductID, *request.ToWarehouseID)
			if er != nil {
				if !errors.Is(er, sql.ErrNoRows) {
					err = er
					return warehouse_state.BranchUpdateResponse{}, err
				} else {
					averagePrice := *request.TotalPrice / (*request.Amount)
					resQ := warehouse_state.AdminCreate{
						Amount:       request.Amount,
						AveragePrice: &averagePrice,
						ProductID:    request.ProductID,
						WarehouseID:  request.ToWarehouseID,
						CreatedBy:    claims.UserId,
						CreatedAt:    time.Now(),
					}
					_, err = r.NewInsert().Model(&resQ).Exec(ctx)
					if err != nil {
						return warehouse_state.BranchUpdateResponse{}, web.NewRequestError(errors.Wrap(err, "creating warehouse_state"), http.StatusBadRequest)
					}
					var (
						amount, price float64
					)

					res := warehouse_state.BranchUpdate{
						ID:                      resQ.ID,
						Amount:                  &amount,
						ProductID:               resQ.ProductID,
						AveragePrice:            &price,
						WarehouseID:             resQ.WarehouseID,
						WarehouseStateHistoryID: &stateHistory.ID,
					}
					response.IncomeSate = &res
					return response, nil
				}
			}

			q := r.NewUpdate().Table("warehouse_state").Where("product_id = ? AND warehouse_id = ?", *request.ProductID, *request.ToWarehouseID)
			q.Set("amount = amount + ?", request.Amount)
			q.Set("average_price = (average_price * amount + ?)/(amount + ?)", request.TotalPrice, request.Amount)
			q.Set("updated_at = ?", time.Now())
			q.Set("updated_by = ?", claims.UserId)

			_, err = q.Exec(ctx)
			if err != nil {
				return warehouse_state.BranchUpdateResponse{}, web.NewRequestError(errors.Wrap(err, "update warehouse_state"), http.StatusInternalServerError)
			}

			res := warehouse_state.BranchUpdate{
				ID:                      state.ID,
				Amount:                  state.Amount,
				ProductID:               state.ProductID,
				AveragePrice:            state.AveragePrice,
				WarehouseID:             state.WarehouseID,
				WarehouseStateHistoryID: &stateHistory.ID,
			}

			response.IncomeSate = &res
		}
	}

	return response, nil
}

func (r Repository) BranchDeleteTransaction(ctx context.Context, request warehouse_state.BranchDeleteTransactionRequest) error {
	claims, err := r.CheckClaims(ctx, auth.RoleBranch)
	if err != nil {
		return err
	}
	err = r.ValidateStruct(request, "Amount", "ProductID", "TotalPrice")
	if err != nil {
		return err
	}

	tx, err := r.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	defer func() {
		if err != nil {
			_ = tx.Rollback()
			return
		}

		_ = tx.Commit()
	}()

	if request.FromWarehouseID != nil {
		stateHistory, er := r.getOldStateByTransactionIDAndWarehouseID(ctx, *request.TransactionProductID, *request.FromWarehouseID)
		if er != nil {
			if !errors.Is(er, sql.ErrNoRows) {
				err = er
				return err
			} else {
				err = web.NewRequestError(errors.New("the product is not available in the warehouse"), http.StatusNotImplemented)
				return err
			}
		}

		q1 := tx.NewUpdate().Table("warehouse_state").Where("product_id = ? AND warehouse_id = ?", *stateHistory.ProductID, *request.FromWarehouseID)
		q1.Set("amount = ?", stateHistory.Amount)
		q1.Set("average_price = ?", stateHistory.AveragePrice)
		q1.Set("updated_at = ?", time.Now())
		q1.Set("updated_by = ?", claims.UserId)

		_, err = q1.Exec(ctx)
		if err != nil {
			return web.NewRequestError(errors.Wrap(err, "update warehouse_state"), http.StatusInternalServerError)
		}
	}

	if request.ToWarehouseID != nil {
		stateHistory, er := r.getOldStateByTransactionIDAndWarehouseID(ctx, *request.TransactionProductID, *request.ToWarehouseID)
		if er != nil {
			if !errors.Is(er, sql.ErrNoRows) {
				err = er
				return err
			} else {
				err = web.NewRequestError(errors.New("the product is not available in the warehouse"), http.StatusNotImplemented)
				return err
			}
		}
		fmt.Println(stateHistory.Amount)
		q1 := tx.NewUpdate().Table("warehouse_state").Where("product_id = ? AND warehouse_id = ?", *stateHistory.ProductID, *request.ToWarehouseID)
		q1.Set("amount = ?", stateHistory.Amount)
		q1.Set("average_price = ?", stateHistory.AveragePrice)
		q1.Set("updated_at = ?", time.Now())
		q1.Set("updated_by = ?", claims.UserId)
		_, err = q1.Exec(ctx)
		if err != nil {
			return web.NewRequestError(errors.Wrap(err, "update warehouse_state"), http.StatusInternalServerError)
		}
	}

	return nil
}

func (r Repository) getByProductIDAndWarehouseID(ctx context.Context, productID, warehouseID int64) (entity.WarehouseState, error) {
	response := entity.WarehouseState{}
	return response, r.NewSelect().Model(&response).Where("product_id = ? AND warehouse_id = ?", productID, warehouseID).Scan(ctx)
}

func (r Repository) getOldStateByTransactionIDAndWarehouseID(ctx context.Context, transactionID, warehouseID int64) (warehouse_state.WarehouseStateHistory, error) {
	response := warehouse_state.WarehouseStateHistory{}
	return response, r.QueryRowContext(ctx, fmt.Sprintf(`
		SELECT 
			wsh.id,
			wsh.amount,
			wsh.average_price,
			ws.product_id
		FROM
			warehouse_state_history AS wsh 
		LEFT JOIN warehouse_state AS ws ON ws.id = wsh.warehouse_state_id
		WHERE 
		    wsh.deleted_at IS NULL AND 
		    wsh.warehouse_transaction_product_id = %d AND
		    ws.warehouse_id = %d 
		ORDER BY wsh.created_at desc`, transactionID, warehouseID)).Scan(&response.ID, &response.Amount, &response.AveragePrice, &response.ProductID)
}

func NewRepository(DB *postgresql.Database) *Repository {
	return &Repository{DB}
}
