package order_payment

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
	"restu-backend/internal/service/order_payment"
	"time"
)

type Repository struct {
	*postgresql.Database
}

// @client

func (r Repository) ClientGetList(ctx context.Context, filter order_payment.Filter) ([]order_payment.ClientGetList, int, error) {
	_, err := r.CheckClaims(ctx, auth.RoleClient)
	if err != nil {
		return nil, 0, err
	}

	whereQuery := fmt.Sprintf(`WHERE op.deleted_at IS NULL `)
	countWhereQuery := whereQuery

	var limitQuery, offsetQuery string
	if filter.Limit != nil {
		limitQuery = fmt.Sprintf(" LIMIT '%d'", *filter.Limit)
	}
	if filter.Offset != nil {
		offsetQuery = fmt.Sprintf(" OFFSET '%d'", *filter.Offset)
	}

	whereQuery += fmt.Sprintf(" %s %s", limitQuery, offsetQuery)

	query := fmt.Sprintf(`
		SELECT
			op.id,
			op.status,
			op.order_id
		FROM
		    order_payment as op
		LEFT OUTER JOIN orders o on o.id = op.order_id
		%s
	`, whereQuery)

	list := make([]order_payment.ClientGetList, 0)

	rows, err := r.QueryContext(ctx, query)
	if err != nil {
		return nil, 0, web.NewRequestError(errors.Wrap(err, "select order_payments"), http.StatusInternalServerError)
	}

	err = r.ScanRows(ctx, rows, &list)
	if err != nil {
		return nil, 0, web.NewRequestError(errors.Wrap(err, "scanning order_payments"), http.StatusBadRequest)
	}

	countQuery := fmt.Sprintf(`
		SELECT
			count(op.id)
		FROM
		    order_payment as op
		%s
	`, countWhereQuery)
	countRows, err := r.QueryContext(ctx, countQuery)
	if err == sql.ErrNoRows {
		return nil, 0, web.NewRequestError(postgres.ErrNotFound, http.StatusNotFound)
	}
	if err != nil {
		return nil, 0, web.NewRequestError(errors.Wrap(err, "selecting order_payment"), http.StatusBadRequest)
	}

	count := 0
	for countRows.Next() {
		if err = countRows.Scan(&count); err != nil {
			return nil, 0, web.NewRequestError(errors.Wrap(err, "scanning order_payment count"), http.StatusBadRequest)
		}
	}

	return list, count, nil
}

func (r Repository) ClientGetDetail(ctx context.Context, id int64) (entity.OrderPayment, error) {
	var detail entity.OrderPayment

	err := r.NewSelect().Model(&detail).Where("id = ? AND deleted_at IS NULL", id).Scan(ctx)
	if err != nil {
		return entity.OrderPayment{}, err
	}
	return detail, nil
}

func (r Repository) ClientCreate(ctx context.Context, request order_payment.ClientCreateRequest) (order_payment.ClientCreateResponse, error) {
	claims, err := r.CheckClaims(ctx, auth.RoleClient)
	if err != nil {
		return order_payment.ClientCreateResponse{}, err
	}

	err = r.ValidateStruct(&request, "OrderID")
	if err != nil {
		return order_payment.ClientCreateResponse{}, err
	}

	response := order_payment.ClientCreateResponse{
		OrderID:   request.OrderID,
		CreatedAt: time.Now(),
		CreatedBy: claims.UserId,
	}

	_, err = r.NewInsert().Model(&response).Exec(ctx)
	if err != nil {
		return order_payment.ClientCreateResponse{}, web.NewRequestError(errors.Wrap(err, "creating order_payment"), http.StatusBadRequest)
	}

	return response, nil
}

func (r Repository) ClientUpdateAll(ctx context.Context, request order_payment.ClientUpdateRequest) error {
	claims, err := r.CheckClaims(ctx, auth.RoleClient)
	if err != nil {
		return err
	}

	if err = r.ValidateStruct(&request, "ID", "Status", "OrderID"); err != nil {
		return err
	}

	q := r.NewUpdate().Table("order_payment").Where("deleted_at IS NULL AND id = ?", request.ID)

	q.Set("status = ?", request.Status)
	q.Set("order_id = ?", request.OrderID)
	q.Set("updated_at = ?", time.Now())
	q.Set("updated_by = ?", claims.UserId)

	_, err = q.Exec(ctx)
	if err != nil {
		return web.NewRequestError(errors.Wrap(err, "updating order_payment"), http.StatusBadRequest)
	}

	return nil
}

func (r Repository) ClientUpdateColumns(ctx context.Context, request order_payment.ClientUpdateRequest) error {
	claims, err := r.CheckClaims(ctx, auth.RoleClient)
	if err != nil {
		return err
	}

	if err = r.ValidateStruct(&request, "ID"); err != nil {
		return err
	}

	q := r.NewUpdate().Table("order_payment").Where("deleted_at IS NULL AND id = ?", request.ID)

	if request.Status != nil {
		q.Set("status = ?", request.Status)
	}
	if request.OrderID != nil {
		q.Set("order_id = ?", request.OrderID)
	}

	q.Set("updated_at = ?", time.Now())
	q.Set("updated_by = ?", claims.UserId)

	_, err = q.Exec(ctx)
	if err != nil {
		return web.NewRequestError(errors.Wrap(err, "updating order_payment"), http.StatusBadRequest)
	}

	return nil
}

func (r Repository) ClientDelete(ctx context.Context, id int64) error {
	return r.DeleteRow(ctx, "order_payment", id, auth.RoleClient)
}

// @cashier

func (r Repository) CashierGetList(ctx context.Context, filter order_payment.Filter) ([]order_payment.CashierGetList, int, error) {
	_, err := r.CheckClaims(ctx, auth.RoleCashier)
	if err != nil {
		return nil, 0, err
	}

	whereQuery := fmt.Sprintf(`WHERE op.deleted_at IS NULL`)
	countWhereQuery := whereQuery

	var limitQuery, offsetQuery string
	if filter.Limit != nil {
		limitQuery = fmt.Sprintf(" LIMIT '%d'", *filter.Limit)
	}
	if filter.Offset != nil {
		offsetQuery = fmt.Sprintf(" OFFSET '%d'", *filter.Offset)
	}

	whereQuery += fmt.Sprintf(" %s %s", limitQuery, offsetQuery)

	query := fmt.Sprintf(`
		SELECT
			op.id,
			op.status,
			op.order_id
		FROM
		    order_payment as op
		LEFT OUTER JOIN orders o on o.id = op.order_id
		%s
	`, whereQuery)

	list := make([]order_payment.CashierGetList, 0)

	rows, err := r.QueryContext(ctx, query)
	if err != nil {
		return nil, 0, web.NewRequestError(errors.Wrap(err, "select order_payments"), http.StatusInternalServerError)
	}

	err = r.ScanRows(ctx, rows, &list)
	if err != nil {
		return nil, 0, web.NewRequestError(errors.Wrap(err, "scanning order_payments"), http.StatusBadRequest)
	}

	countQuery := fmt.Sprintf(`
		SELECT
			count(op.id)
		FROM
		    order_payment as op
		%s
	`, countWhereQuery)
	countRows, err := r.QueryContext(ctx, countQuery)
	if err == sql.ErrNoRows {
		return nil, 0, web.NewRequestError(postgres.ErrNotFound, http.StatusNotFound)
	}
	if err != nil {
		return nil, 0, web.NewRequestError(errors.Wrap(err, "selecting order_payment"), http.StatusBadRequest)
	}

	count := 0
	for countRows.Next() {
		if err = countRows.Scan(&count); err != nil {
			return nil, 0, web.NewRequestError(errors.Wrap(err, "scanning order_payment count"), http.StatusBadRequest)
		}
	}

	return list, count, nil
}

func (r Repository) CashierGetDetail(ctx context.Context, id int64) (entity.OrderPayment, error) {
	var detail entity.OrderPayment

	err := r.NewSelect().Model(&detail).Where("id = ? AND deleted_at IS NULL", id).Scan(ctx)
	if err != nil {
		return entity.OrderPayment{}, err
	}
	return detail, nil
}

func (r Repository) CashierCreate(ctx context.Context, request order_payment.CashierCreateRequest) (order_payment.CashierCreateResponse, error) {
	claims, err := r.CheckClaims(ctx, auth.RoleCashier)
	if err != nil {
		return order_payment.CashierCreateResponse{}, err
	}

	err = r.ValidateStruct(&request, "OrderID")
	if err != nil {
		return order_payment.CashierCreateResponse{}, err
	}

	query := fmt.Sprintf(`
			SELECT
					SUM(m.new_price * om.count) AS price
			FROM
			orders o
			JOIN order_menu om ON o.id = om.order_id
			JOIN menus m ON om.menu_id = m.id
			WHERE
			om.order_id = '%d'
			AND o.deleted_at IS NULL
			AND o.status != 'CANCELLED'
-- 			AND NOT EXISTS (
-- 				SELECT 1
-- 					FROM order_menu
-- 			WHERE order_id = '%d'
-- 			AND status = 'NEW'
-- 			AND deleted_at IS NULL
-- 				)
				GROUP BY o.id
					`, *request.OrderID, *request.OrderID)

	var price *float64

	err = r.QueryRowContext(ctx, query).Scan(
		&price,
	)
	if errors.Is(err, sql.ErrNoRows) {
		return order_payment.CashierCreateResponse{}, web.NewRequestError(errors.New("order_menu is new"), http.StatusBadRequest)
	}
	if err != nil {
		return order_payment.CashierCreateResponse{}, web.NewRequestError(errors.Wrap(err, "select order price"), http.StatusInternalServerError)
	}

	status := "PAID"

	response := order_payment.CashierCreateResponse{
		OrderID:   request.OrderID,
		Status:    &status,
		Price:     price,
		CreatedAt: time.Now(),
		CreatedBy: claims.UserId,
	}

	_, err = r.NewInsert().Model(&response).Exec(ctx)
	if err != nil {
		return order_payment.CashierCreateResponse{}, web.NewRequestError(errors.Wrap(err, "creating order_payment"), http.StatusBadRequest)
	}

	return response, nil
}

func (r Repository) CashierUpdateAll(ctx context.Context, request order_payment.CashierUpdateRequest) error {
	claims, err := r.CheckClaims(ctx, auth.RoleCashier)
	if err != nil {
		return err
	}

	if err = r.ValidateStruct(&request, "ID", "Status", "OrderID"); err != nil {
		return err
	}

	q := r.NewUpdate().Table("order_payment").Where("deleted_at IS NULL AND id = ?",
		request.ID)

	q.Set("status = ?", request.Status)
	q.Set("order_id = ?", request.OrderID)
	q.Set("updated_at = ?", time.Now())
	q.Set("updated_by = ?", claims.UserId)

	_, err = q.Exec(ctx)
	if err != nil {
		return web.NewRequestError(errors.Wrap(err, "updating order_payment"), http.StatusBadRequest)
	}

	return nil
}

func (r Repository) CashierUpdateColumns(ctx context.Context, request order_payment.CashierUpdateRequest) error {
	claims, err := r.CheckClaims(ctx, auth.RoleCashier)
	if err != nil {
		return err
	}

	if err = r.ValidateStruct(&request, "ID"); err != nil {
		return err
	}

	q := r.NewUpdate().Table("order_payment").Where("deleted_at IS NULL AND id = ?",
		request.ID)

	if request.Status != nil {
		q.Set("status = ?", request.Status)
	}
	if request.OrderID != nil {
		q.Set("order_id = ?", request.OrderID)
	}

	q.Set("updated_at = ?", time.Now())
	q.Set("updated_by = ?", claims.UserId)

	_, err = q.Exec(ctx)
	if err != nil {
		return web.NewRequestError(errors.Wrap(err, "updating order_payment"), http.StatusBadRequest)
	}

	return nil
}

func (r Repository) CashierDelete(ctx context.Context, id int64) error {
	return r.DeleteRow(ctx, "order_payment", id, auth.RoleCashier)
}

func NewRepository(DB *postgresql.Database) *Repository {
	return &Repository{DB}
}
func (r *Repository) GetFoodPrice(ctx context.Context, menuID, branchID int64) (price float64, err error) {
	_, err = r.NewSelect().Table("menus").Column("new_price").Where("deleted_at IS NULL AND id = ? AND branch_id = ?", menuID, branchID).Limit(1).Exec(ctx, &price)
	return
}

//
//if data.Status != nil && *data.Status == "PAID" {
//var response orderPayment.CashierCreateResponse
//
//query := fmt.Sprintf(`
//					SELECT
//					    SUM(m.new_price * om.count) AS price
//					FROM
//					    orders o
//					    JOIN order_menu om ON o.id = om.order_id
//					    JOIN menus m ON om.menu_id = m.id
//					WHERE om.order_id = '%d' AND o.deleted_at IS NULL AND o.status!='CANCELLED' AND o.status !='PAID' AND om.deleted_at IS NULL
//					GROUP BY
//					    o.id
//					`, data.Id)
//
//var detail float64
//
//err = r.QueryRowContext(ctx, query).Scan(
//&detail,
//)
//fmt.Println(detail)
//if err != nil {
//return web.NewRequestError(errors.Wrap(err, "select order price"), http.StatusInternalServerError)
//}
//// inserting order_payment...
//
//insertOrderPaymentQuery := fmt.Sprintf(`INSERT INTO order_payment (order_id, price, created_by, status) VALUES (%d, %f, %d, '%s') RETURNING id, created_at`, data.Id, detail, claims.UserId, *data.Status)
//
//if err = r.QueryRowContext(ctx, insertOrderPaymentQuery).Scan(&response.ID, &response.CreatedAt); err != nil {
//return web.NewRequestError(errors.Wrap(err, "inserting order_payment"), http.StatusInternalServerError)
//}
//}
