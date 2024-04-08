package warehouse_transaction

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
	"restu-backend/internal/service/warehouse_transaction"
	"time"
)

type Repository struct {
	*postgresql.Database
}

func (r Repository) AdminGetList(ctx context.Context, filter warehouse_transaction.Filter) ([]warehouse_transaction.AdminGetListResponse, int, error) {
	claims, err := r.CheckClaims(ctx, auth.RoleAdmin)
	if err != nil {
		return nil, 0, err
	}

	whereQuery := fmt.Sprintf(`WHERE wt.deleted_at IS NULL AND (fb.restaurant_id = '%d' OR tb.restaurant_id  = %d OR fp.restaurant_id  = %d OR tp.restaurant_id  = %d)`, *claims.RestaurantID, *claims.RestaurantID, *claims.RestaurantID, *claims.RestaurantID)

	var limitQuery, offsetQuery string
	if filter.Limit != nil {
		limitQuery = fmt.Sprintf(" LIMIT '%d'", *filter.Limit)
	}
	if filter.Offset != nil {
		offsetQuery = fmt.Sprintf(" OFFSET '%d'", *filter.Offset)
	}

	OrderQuery := " ORDER BY wt.created_at desc"

	if filter.WarehouseID != nil {
		whereQuery += fmt.Sprintf(" AND (CASE WHEN wt.from_warehouse_id IS NOT NULL THEN wt.from_warehouse_id = %d else false end  OR CASE WHEN wt.to_warehouse_id IS NOT NULL THEN wt.to_warehouse_id = %d else false end)", *filter.WarehouseID, *filter.WarehouseID)
	}

	query := fmt.Sprintf(`
		SELECT
			wt.id,
			wt.from_warehouse_id,
			fw.name,
			wt.from_partner_id,
			fp.name,
			wt.to_warehouse_id,
			tw.name,
			wt.to_partner_id,
			tp.name
		FROM
		    warehouse_transactions as wt
		LEFT JOIN warehouses fw ON fw.id = wt.from_warehouse_id
		LEFT JOIN partners fp ON fp.id = wt.from_partner_id
		LEFT JOIN warehouses tw ON tw.id = wt.to_warehouse_id
		LEFT JOIN partners tp ON tp.id = wt.to_partner_id
		LEFT JOIN branches fb ON fw.branch_id = fb.id
		LEFT JOIN branches tb ON tw.branch_id = tb.id
		%s %s %s %s
	`, whereQuery, OrderQuery, limitQuery, offsetQuery)

	rows, err := r.QueryContext(ctx, query)
	if err != nil {
		return nil, 0, web.NewRequestError(errors.Wrap(err, "select user"), http.StatusInternalServerError)
	}

	var list []warehouse_transaction.AdminGetListResponse

	for rows.Next() {
		var detail warehouse_transaction.AdminGetListResponse
		if err = rows.Scan(
			&detail.ID,
			&detail.FromWarehouseID,
			&detail.FromWarehouse,
			&detail.FromPartnerID,
			&detail.FromPartner,
			&detail.ToWarehouseID,
			&detail.ToWarehouse,
			&detail.ToPartnerID,
			&detail.ToPartner,
		); err != nil {
			return nil, 0, web.NewRequestError(errors.Wrap(err, "scanning warehouse"), http.StatusBadRequest)
		}

		list = append(list, detail)
	}

	countQuery := fmt.Sprintf(`
		SELECT
			count(wt.id)
		FROM
		    warehouse_transactions as wt
		LEFT JOIN warehouses fw ON fw.id = wt.from_warehouse_id
		LEFT JOIN partners fp ON fp.id = wt.from_partner_id
		LEFT JOIN warehouses tw ON tw.id = wt.to_warehouse_id
		LEFT JOIN partners tp ON tp.id = wt.to_partner_id
		LEFT JOIN branches fb ON fw.branch_id = fb.id
		LEFT JOIN branches tb ON tw.branch_id = tb.id
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

func (r Repository) AdminCreate(ctx context.Context, request warehouse_transaction.AdminCreateRequest) (warehouse_transaction.AdminCreateResponse, error) {
	claims, err := r.CheckClaims(ctx, auth.RoleAdmin)
	if err != nil {
		return warehouse_transaction.AdminCreateResponse{}, err
	}

	err = r.ValidateStruct(&request)
	if err != nil {
		return warehouse_transaction.AdminCreateResponse{}, err
	}

	response := warehouse_transaction.AdminCreateResponse{
		FromWarehouseID: request.FromWarehouseID,
		FromPartnerID:   request.FromPartnerID,
		ToWarehouseID:   request.ToWarehouseID,
		ToPartnerID:     request.ToPartnerID,
		CreatedAt:       time.Now(),
		CreatedBy:       claims.UserId,
	}

	_, err = r.NewInsert().Model(&response).Exec(ctx)
	if err != nil {
		return warehouse_transaction.AdminCreateResponse{}, web.NewRequestError(errors.Wrap(err, "creating user"), http.StatusBadRequest)
	}

	return response, nil
}

func (r Repository) AdminUpdateColumn(ctx context.Context, request warehouse_transaction.AdminUpdateRequest) error {
	claims, err := r.CheckClaims(ctx, auth.RoleAdmin)
	if err != nil {
		return err
	}

	err = r.ValidateStruct(&request, "ID")
	if err != nil {
		return err
	}
	q := r.NewUpdate().Table("warehouse_transactions").Where("deleted_at IS NULL AND id = ?",
		request.ID, claims.RestaurantID)

	if request.FromWarehouseID != nil {
		q.Set("from_warehouse_id = ?, from_partner_id = null ", request.FromWarehouseID)
	}
	if request.FromPartnerID != nil {
		q.Set("from_warehouse_id = null, from_partner_id = ? ", request.FromPartnerID)
	}
	if request.ToWarehouseID != nil {
		q.Set("to_warehouse_id = ?, to_partner_id = null ", request.ToWarehouseID)
	}
	if request.ToPartnerID != nil {
		q.Set("to_warehouse_id = null, to_partner_id = ? ", request.ToPartnerID)
	}

	q.Set("updated_at = ?", time.Now())
	q.Set("updated_by = ?", claims.UserId)

	_, err = q.Exec(ctx)
	if err != nil {
		return web.NewRequestError(errors.Wrap(err, "updating food"), http.StatusBadRequest)
	}

	return nil
}

func (r Repository) AdminGetDetailByID(ctx context.Context, id int64) (warehouse_transaction.AdminGetDetailByIdResponse, error) {
	claims, err := r.CheckClaims(ctx, auth.RoleAdmin)
	if err != nil {
		return warehouse_transaction.AdminGetDetailByIdResponse{}, err
	}

	whereQuery := fmt.Sprintf(`WHERE wt.deleted_at IS NULL AND (fb.restaurant_id = %d OR tb.restaurant_id  = %d OR fp.restaurant_id  = %d OR tp.restaurant_id  = %d) AND wt.id = %d`, *claims.RestaurantID, *claims.RestaurantID, *claims.RestaurantID, *claims.RestaurantID, id)

	query := fmt.Sprintf(`
		SELECT
			wt.id,
			wt.from_warehouse_id,
			wt.from_partner_id,
			wt.to_warehouse_id,
			wt.to_partner_id
		FROM
		    warehouse_transactions as wt
		LEFT JOIN warehouses fw ON fw.id = wt.from_warehouse_id
		LEFT JOIN partners fp ON fp.id = wt.from_partner_id
		LEFT JOIN warehouses tw ON tw.id = wt.to_warehouse_id
		LEFT JOIN partners tp ON tp.id = wt.to_partner_id
		LEFT JOIN branches fb ON fw.branch_id = fb.id
		LEFT JOIN branches tb ON tw.branch_id = tb.id
		%s
	`, whereQuery)
	var detail warehouse_transaction.AdminGetDetailByIdResponse

	err = r.QueryRowContext(ctx, query).Scan(
		&detail.ID,
		&detail.FromWarehouseID,
		&detail.FromPartnerID,
		&detail.ToWarehouseID,
		&detail.ToPartnerID,
	)
	if err != nil {
		return warehouse_transaction.AdminGetDetailByIdResponse{}, web.NewRequestError(errors.Wrap(err, "select user"), http.StatusInternalServerError)
	}

	return detail, nil
}

func (r Repository) AdminDelete(ctx context.Context, id int64) error {
	return r.DeleteRow(ctx, "warehouse_transactions", id, auth.RoleAdmin)
}

func (r Repository) BranchGetList(ctx context.Context, filter warehouse_transaction.Filter) ([]warehouse_transaction.BranchGetListResponse, int, error) {
	claims, err := r.CheckClaims(ctx, auth.RoleBranch)
	if err != nil {
		return nil, 0, err
	}

	whereQuery := fmt.Sprintf(`WHERE wt.deleted_at IS NULL AND (fb.id = %d OR tb.id  = %d)`, *claims.BranchID, *claims.BranchID)

	var limitQuery, offsetQuery string
	if filter.Limit != nil {
		limitQuery = fmt.Sprintf(" LIMIT '%d'", *filter.Limit)
	}
	if filter.Offset != nil {
		offsetQuery = fmt.Sprintf(" OFFSET '%d'", *filter.Offset)
	}

	if filter.WarehouseID != nil {
		whereQuery += fmt.Sprintf(" AND (wt.from_warehouse_id = %d OR wt.to_warehouse_id = %d)", *filter.WarehouseID, *filter.WarehouseID)
	}

	query := fmt.Sprintf(`
		SELECT
			wt.id,
			wt.from_warehouse_id,
			fw.name,
			wt.from_partner_id,
			fp.name,
			wt.to_warehouse_id,
			tw.name,
			wt.to_partner_id,
			tp.name
		FROM
		    warehouse_transactions as wt
		LEFT JOIN warehouses fw ON fw.id = wt.from_warehouse_id
		LEFT JOIN partners fp ON fp.id = wt.from_partner_id
		LEFT JOIN warehouses tw ON tw.id = wt.to_warehouse_id
		LEFT JOIN partners tp ON tp.id = wt.to_partner_id
		LEFT JOIN branches fb ON fw.branch_id = fb.id
		LEFT JOIN branches tb ON tw.branch_id = tb.id
		%s %s %s
	`, whereQuery, limitQuery, offsetQuery)

	rows, err := r.QueryContext(ctx, query)
	if err != nil {
		return nil, 0, web.NewRequestError(errors.Wrap(err, "select user"), http.StatusInternalServerError)
	}

	var list []warehouse_transaction.BranchGetListResponse

	for rows.Next() {
		var detail warehouse_transaction.BranchGetListResponse
		if err = rows.Scan(
			&detail.ID,
			&detail.FromWarehouseID,
			&detail.FromWarehouse,
			&detail.FromPartnerID,
			&detail.FromPartner,
			&detail.ToWarehouseID,
			&detail.ToWarehouse,
			&detail.ToPartnerID,
			&detail.ToPartner,
		); err != nil {
			return nil, 0, web.NewRequestError(errors.Wrap(err, "scanning warehouse"), http.StatusBadRequest)
		}

		list = append(list, detail)
	}

	countQuery := fmt.Sprintf(`
		SELECT
			count(wt.id)
		FROM
		    warehouse_transactions as wt
		LEFT JOIN warehouses fw ON fw.id = wt.from_warehouse_id
		LEFT JOIN partners fp ON fp.id = wt.from_partner_id
		LEFT JOIN warehouses tw ON tw.id = wt.to_warehouse_id
		LEFT JOIN partners tp ON tp.id = wt.to_partner_id
		LEFT JOIN branches fb ON fw.branch_id = fb.id
		LEFT JOIN branches tb ON tw.branch_id = tb.id
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

func (r Repository) BranchCreate(ctx context.Context, request warehouse_transaction.BranchCreateRequest) (warehouse_transaction.BranchCreateResponse, error) {
	claims, err := r.CheckClaims(ctx, auth.RoleBranch)
	if err != nil {
		return warehouse_transaction.BranchCreateResponse{}, err
	}

	err = r.ValidateStruct(&request)
	if err != nil {
		return warehouse_transaction.BranchCreateResponse{}, err
	}

	response := warehouse_transaction.BranchCreateResponse{
		FromWarehouseID: request.FromWarehouseID,
		FromPartnerID:   request.FromPartnerID,
		ToWarehouseID:   request.ToWarehouseID,
		ToPartnerID:     request.ToPartnerID,
		CreatedAt:       time.Now(),
		CreatedBy:       claims.UserId,
	}

	_, err = r.NewInsert().Model(&response).Exec(ctx)
	if err != nil {
		return warehouse_transaction.BranchCreateResponse{}, web.NewRequestError(errors.Wrap(err, "creating warehouse transaction"), http.StatusNotImplemented)
	}

	return response, nil
}

func (r Repository) BranchUpdateColumn(ctx context.Context, request warehouse_transaction.BranchUpdateRequest) error {
	claims, err := r.CheckClaims(ctx, auth.RoleBranch)
	if err != nil {
		return err
	}

	err = r.ValidateStruct(&request, "ID")
	if err != nil {
		return err
	}
	q := r.NewUpdate().Table("warehouse_transactions").Where("deleted_at IS NULL AND id = ?",
		request.ID, claims.RestaurantID)

	if request.FromWarehouseID != nil {
		q.Set("from_warehouse_id = ?, from_partner_id = null ", request.FromWarehouseID)
	}
	if request.FromPartnerID != nil {
		q.Set("from_warehouse_id = null, from_partner_id = ? ", request.FromPartnerID)
	}
	if request.ToWarehouseID != nil {
		q.Set("to_warehouse_id = ?, to_partner_id = null ", request.ToWarehouseID)
	}
	if request.ToPartnerID != nil {
		q.Set("to_warehouse_id = null, to_partner_id = ? ", request.ToPartnerID)
	}

	q.Set("updated_at = ?", time.Now())
	q.Set("updated_by = ?", claims.UserId)

	_, err = q.Exec(ctx)
	if err != nil {
		return web.NewRequestError(errors.Wrap(err, "updating food"), http.StatusBadRequest)
	}

	return nil
}

func (r Repository) BranchGetDetailByID(ctx context.Context, id int64) (warehouse_transaction.BranchGetDetailByIdResponse, error) {
	claims, err := r.CheckClaims(ctx, auth.RoleBranch)
	if err != nil {
		return warehouse_transaction.BranchGetDetailByIdResponse{}, err
	}

	whereQuery := fmt.Sprintf(`WHERE wt.deleted_at IS NULL AND (fb.id = %d OR tb.id = %d) AND wt.id = %d`, *claims.BranchID, *claims.BranchID, id)

	query := fmt.Sprintf(`
		SELECT
			wt.id,
			wt.from_warehouse_id,
			wt.from_partner_id,
			wt.to_warehouse_id,
			wt.to_partner_id
		FROM
		    warehouse_transactions as wt
		LEFT JOIN warehouses fw ON fw.id = wt.from_warehouse_id
		LEFT JOIN partners fp ON fp.id = wt.from_partner_id
		LEFT JOIN warehouses tw ON tw.id = wt.to_warehouse_id
		LEFT JOIN partners tp ON tp.id = wt.to_partner_id
		LEFT JOIN branches fb ON fw.branch_id = fb.id
		LEFT JOIN branches tb ON tw.branch_id = tb.id
		%s
	`, whereQuery)
	var detail warehouse_transaction.BranchGetDetailByIdResponse

	err = r.QueryRowContext(ctx, query).Scan(
		&detail.ID,
		&detail.FromWarehouseID,
		&detail.FromPartnerID,
		&detail.ToWarehouseID,
		&detail.ToPartnerID,
	)
	if err != nil {
		return warehouse_transaction.BranchGetDetailByIdResponse{}, web.NewRequestError(errors.Wrap(err, "select user"), http.StatusInternalServerError)
	}

	return detail, nil
}

func (r Repository) BranchDelete(ctx context.Context, id int64) error {
	return r.DeleteRow(ctx, "warehouse_transactions", id, auth.RoleBranch)
}

func NewRepository(DB *postgresql.Database) *Repository {
	return &Repository{DB}
}
