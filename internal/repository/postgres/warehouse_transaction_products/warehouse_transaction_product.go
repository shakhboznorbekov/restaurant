package warehouse_transaction_product

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/pkg/errors"
	"net/http"
	"restu-backend/foundation/web"
	"restu-backend/internal/auth"
	"restu-backend/internal/pkg/repository/postgresql"
	"restu-backend/internal/repository/postgres"
	"restu-backend/internal/service/warehouse_transaction_product"
	"time"
)

type Repository struct {
	*postgresql.Database
}

func (r Repository) AdminGetList(ctx context.Context, filter warehouse_transaction_product.Filter, transactionId int64) ([]warehouse_transaction_product.AdminGetListResponse, int, error) {
	_, err := r.CheckClaims(ctx, auth.RoleAdmin)
	if err != nil {
		return nil, 0, err
	}

	whereQuery := fmt.Sprintf(`WHERE wtp.deleted_at IS NULL AND wtp.transaction_id = %d`, transactionId)

	var limitQuery, offsetQuery string
	if filter.Limit != nil {
		limitQuery = fmt.Sprintf(" LIMIT '%d'", *filter.Limit)
	}
	if filter.Offset != nil {
		offsetQuery = fmt.Sprintf(" OFFSET '%d'", *filter.Offset)
	}

	OrderQuery := " ORDER BY wtp.created_at desc"

	if filter.WarehouseID != nil {
		whereQuery += fmt.Sprintf(" AND (CASE WHEN wtp.from_warehouse_id IS NOT NULL THEN wtp.from_warehouse_id = %d else false end  OR CASE WHEN wtp.to_warehouse_id IS NOT NULL THEN wtp.to_warehouse_id = %d else false end)", *filter.WarehouseID, *filter.WarehouseID)
	}

	query := fmt.Sprintf(`
		SELECT
			wtp.id,
			wtp.amount,
			wtp.total_price,
			wtp.product_id,
			p.name
		FROM
		    warehouse_transaction_products as wtp
		LEFT JOIN products p ON p.id = wtp.product_id
		%s %s %s %s
	`, whereQuery, OrderQuery, limitQuery, offsetQuery)

	rows, err := r.QueryContext(ctx, query)
	if err != nil {
		return nil, 0, web.NewRequestError(errors.Wrap(err, "select user"), http.StatusInternalServerError)
	}

	var list []warehouse_transaction_product.AdminGetListResponse

	for rows.Next() {
		var detail warehouse_transaction_product.AdminGetListResponse
		if err = rows.Scan(
			&detail.ID,
			&detail.Amount,
			&detail.TotalPrice,
			&detail.ProductID,
			&detail.Product,
		); err != nil {
			return nil, 0, web.NewRequestError(errors.Wrap(err, "scanning warehouse"), http.StatusBadRequest)
		}

		list = append(list, detail)
	}

	countQuery := fmt.Sprintf(`
		SELECT
			count(wtp.id)
		FROM
		    warehouse_transaction_products as wtp
		LEFT JOIN products p ON p.id = wtp.product_id
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

func (r Repository) AdminCreate(ctx context.Context, request warehouse_transaction_product.AdminCreateRequest) (warehouse_transaction_product.AdminCreateResponse, error) {
	claims, err := r.CheckClaims(ctx, auth.RoleAdmin)
	if err != nil {
		return warehouse_transaction_product.AdminCreateResponse{}, err
	}

	err = r.ValidateStruct(&request, "Amount", "ProductID", "TotalPrice", "TransactionId")
	if err != nil {
		return warehouse_transaction_product.AdminCreateResponse{}, err
	}

	response := warehouse_transaction_product.AdminCreateResponse{
		Amount:        request.Amount,
		TotalPrice:    request.TotalPrice,
		ProductID:     request.ProductID,
		TransactionId: request.TransactionId,
		CreatedAt:     time.Now(),
		CreatedBy:     claims.UserId,
	}

	_, err = r.NewInsert().Model(&response).Exec(ctx)
	if err != nil {
		return warehouse_transaction_product.AdminCreateResponse{}, web.NewRequestError(errors.Wrap(err, "creating user"), http.StatusBadRequest)
	}

	return response, nil
}

func (r Repository) AdminUpdateColumn(ctx context.Context, request warehouse_transaction_product.AdminUpdateRequest) error {
	claims, err := r.CheckClaims(ctx, auth.RoleAdmin)
	if err != nil {
		return err
	}

	err = r.ValidateStruct(&request, "ID")
	if err != nil {
		return err
	}
	q := r.NewUpdate().Table("warehouse_transaction_products").Where("deleted_at IS NULL AND id = ?",
		request.ID, claims.RestaurantID)

	if request.Amount != nil {
		q.Set("amount = ?", request.Amount)
	}
	if request.ProductID != nil {
		q.Set("product_id = ?", request.ProductID)
	}
	if request.TotalPrice != nil {
		q.Set("total_price = ?", request.TotalPrice)
	}

	q.Set("updated_at = ?", time.Now())
	q.Set("updated_by = ?", claims.UserId)

	_, err = q.Exec(ctx)
	if err != nil {
		return web.NewRequestError(errors.Wrap(err, "updating food"), http.StatusBadRequest)
	}

	return nil
}

func (r Repository) AdminGetDetailByID(ctx context.Context, id int64) (warehouse_transaction_product.AdminGetDetailByIdResponse, error) {
	_, err := r.CheckClaims(ctx, auth.RoleAdmin)
	if err != nil {
		return warehouse_transaction_product.AdminGetDetailByIdResponse{}, err
	}

	whereQuery := fmt.Sprintf(`WHERE wtp.deleted_at IS NULL AND wtp.id = %d`, id)

	query := fmt.Sprintf(`
		SELECT
			wtp.id,
			wtp.amount,
			wtp.total_price,
			wtp.product_id,
			wtp.created_at,
		    wtp.created_by,
		    wtp.transaction_id
		FROM
		    warehouse_transaction_products as wtp
		LEFT JOIN products p ON p.id = wtp.product_id
		%s
	`, whereQuery)
	var detail warehouse_transaction_product.AdminGetDetailByIdResponse

	err = r.QueryRowContext(ctx, query).Scan(
		&detail.ID,
		&detail.Amount,
		&detail.TotalPrice,
		&detail.ProductID,
		&detail.CreatedAt,
		&detail.CreatedBy,
		&detail.TransactionId,
	)
	if err != nil {
		return warehouse_transaction_product.AdminGetDetailByIdResponse{}, web.NewRequestError(errors.Wrap(err, "select user"), http.StatusInternalServerError)
	}

	return detail, nil
}

func (r Repository) AdminDelete(ctx context.Context, id int64) error {
	return r.DeleteRow(ctx, "warehouse_transaction_products", id, auth.RoleAdmin)
}

func (r Repository) BranchGetList(ctx context.Context, filter warehouse_transaction_product.Filter, transactionId int64) ([]warehouse_transaction_product.BranchGetListResponse, int, error) {
	_, err := r.CheckClaims(ctx, auth.RoleBranch)
	if err != nil {
		return nil, 0, err
	}

	whereQuery := fmt.Sprintf(`WHERE wtp.deleted_at IS NULL AND wtp.transaction_id = %d`, transactionId)

	var limitQuery, offsetQuery string
	if filter.Limit != nil {
		limitQuery = fmt.Sprintf(" LIMIT '%d'", *filter.Limit)
	}
	if filter.Offset != nil {
		offsetQuery = fmt.Sprintf(" OFFSET '%d'", *filter.Offset)
	}

	OrderQuery := " ORDER BY wtp.created_at desc"

	if filter.WarehouseID != nil {
		whereQuery += fmt.Sprintf(" AND (CASE WHEN wtp.from_warehouse_id IS NOT NULL THEN wtp.from_warehouse_id = %d else false end  OR CASE WHEN wtp.to_warehouse_id IS NOT NULL THEN wtp.to_warehouse_id = %d else false end)", *filter.WarehouseID, *filter.WarehouseID)
	}

	query := fmt.Sprintf(`
		SELECT
			wtp.id,
			wtp.amount,
			wtp.total_price,
			wtp.product_id,
			p.name
		FROM
		    warehouse_transaction_products as wtp
		LEFT JOIN products p ON p.id = wtp.product_id
		%s %s %s %s
	`, whereQuery, OrderQuery, limitQuery, offsetQuery)

	rows, err := r.QueryContext(ctx, query)
	if err != nil {
		return nil, 0, web.NewRequestError(errors.Wrap(err, "select user"), http.StatusInternalServerError)
	}

	var list []warehouse_transaction_product.BranchGetListResponse

	for rows.Next() {
		var detail warehouse_transaction_product.BranchGetListResponse
		if err = rows.Scan(
			&detail.ID,
			&detail.Amount,
			&detail.TotalPrice,
			&detail.ProductID,
			&detail.Product,
		); err != nil {
			return nil, 0, web.NewRequestError(errors.Wrap(err, "scanning warehouse"), http.StatusBadRequest)
		}

		list = append(list, detail)
	}

	countQuery := fmt.Sprintf(`
		SELECT
			count(wtp.id)
		FROM
		    warehouse_transaction_products as wtp
		LEFT JOIN products p ON p.id = wtp.product_id
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

func (r Repository) BranchCreate(ctx context.Context, request warehouse_transaction_product.BranchCreateRequest) (warehouse_transaction_product.BranchCreateResponse, error) {
	claims, err := r.CheckClaims(ctx, auth.RoleBranch)
	if err != nil {
		return warehouse_transaction_product.BranchCreateResponse{}, err
	}

	err = r.ValidateStruct(&request, "Amount", "ProductID", "TotalPrice", "TransactionId")
	if err != nil {
		return warehouse_transaction_product.BranchCreateResponse{}, err
	}

	response := warehouse_transaction_product.BranchCreateResponse{
		Amount:        request.Amount,
		TotalPrice:    request.TotalPrice,
		ProductID:     request.ProductID,
		TransactionId: request.TransactionId,
		CreatedAt:     time.Now(),
		CreatedBy:     claims.UserId,
	}

	_, err = r.NewInsert().Model(&response).Exec(ctx)
	if err != nil {
		return warehouse_transaction_product.BranchCreateResponse{}, web.NewRequestError(errors.Wrap(err, "creating user"), http.StatusBadRequest)
	}

	return response, nil
}

func (r Repository) BranchUpdateColumn(ctx context.Context, request warehouse_transaction_product.BranchUpdateRequest) error {
	claims, err := r.CheckClaims(ctx, auth.RoleBranch)
	if err != nil {
		return err
	}

	err = r.ValidateStruct(&request, "ID")
	if err != nil {
		return err
	}
	q := r.NewUpdate().Table("warehouse_transaction_products").Where("deleted_at IS NULL AND id = ?",
		request.ID, claims.RestaurantID)

	if request.Amount != nil {
		q.Set("amount = ?", request.Amount)
	}
	if request.ProductID != nil {
		q.Set("product_id = ?", request.ProductID)
	}
	if request.TotalPrice != nil {
		q.Set("total_price = ?", request.TotalPrice)
	}

	q.Set("updated_at = ?", time.Now())
	q.Set("updated_by = ?", claims.UserId)

	_, err = q.Exec(ctx)
	if err != nil {
		return web.NewRequestError(errors.Wrap(err, "updating food"), http.StatusBadRequest)
	}

	return nil
}

func (r Repository) BranchGetDetailByID(ctx context.Context, id int64) (warehouse_transaction_product.BranchGetDetailByIdResponse, error) {
	_, err := r.CheckClaims(ctx, auth.RoleBranch)
	if err != nil {
		return warehouse_transaction_product.BranchGetDetailByIdResponse{}, err
	}

	whereQuery := fmt.Sprintf(`WHERE wtp.deleted_at IS NULL AND wtp.id = %d`, id)

	query := fmt.Sprintf(`
		SELECT
			wtp.id,
			wtp.amount,
			wtp.total_price,
			wtp.product_id,
			wtp.created_at,
		    wtp.created_by,
		    wtp.transaction_id
		FROM
		    warehouse_transaction_products as wtp
		LEFT JOIN products p ON p.id = wtp.product_id
		%s
	`, whereQuery)
	var detail warehouse_transaction_product.BranchGetDetailByIdResponse

	err = r.QueryRowContext(ctx, query).Scan(
		&detail.ID,
		&detail.Amount,
		&detail.TotalPrice,
		&detail.ProductID,
		&detail.CreatedAt,
		&detail.CreatedBy,
		&detail.TransactionId,
	)
	if err != nil {
		return warehouse_transaction_product.BranchGetDetailByIdResponse{}, web.NewRequestError(errors.Wrap(err, "select user"), http.StatusInternalServerError)
	}

	return detail, nil
}

func (r Repository) BranchDelete(ctx context.Context, id int64) error {
	return r.DeleteRow(ctx, "warehouse_transaction_products", id, auth.RoleBranch)
}
func NewRepository(DB *postgresql.Database) *Repository {
	return &Repository{DB}
}
