package table

import (
	"context"
	"database/sql"
	"fmt"
	"net/http"
	"restu-backend/foundation/web"
	"restu-backend/internal/auth"
	"restu-backend/internal/entity"
	"restu-backend/internal/pkg/repository/postgresql"
	"restu-backend/internal/pkg/utils"
	"restu-backend/internal/repository/postgres"
	"restu-backend/internal/service/tables"
	"time"

	"github.com/pkg/errors"
)

type Repository struct {
	*postgresql.Database
}

// @admin

func (r Repository) AdminGetList(ctx context.Context, filter tables.Filter) ([]tables.AdminGetList, int, error) {
	claims, err := r.CheckClaims(ctx, auth.RoleAdmin)
	if err != nil {
		return nil, 0, err
	}

	table := "tables"
	whereQuery := fmt.Sprintf(`WHERE %s.deleted_at IS NULL and %s.branch_id IN (select id from branches where restaurant_id = %d)`, table, table, *claims.RestaurantID)
	if filter.BranchID != nil {
		whereQuery += fmt.Sprintf(" AND %s.branch_id = %d", table, *filter.BranchID)
	}
	if filter.HallID != nil {
		whereQuery += fmt.Sprintf(" AND %s.hall_id = %d", table, *filter.HallID)
	}
	countWhereQuery := whereQuery

	query, err := utils.SelectQuery(filter.Fields, filter.Joins, &table, &whereQuery)
	if err != nil {
		return nil, 0, errors.Wrap(err, "select query")
	}

	list := make([]tables.AdminGetList, 0)

	rows, err := r.QueryContext(ctx, query)
	if err != nil {
		return nil, 0, web.NewRequestError(errors.Wrap(err, "select table"), http.StatusInternalServerError)
	}

	err = r.ScanRows(ctx, rows, &list)
	if err != nil {
		return nil, 0, web.NewRequestError(errors.Wrap(err, "scanning tables"), http.StatusBadRequest)
	}

	countQuery := fmt.Sprintf(`
		SELECT
			count(tables.id)
		FROM
		    %s
			JOIN branches on branches.id=tables.branch_id
		%s
	`, table, countWhereQuery)

	countRows, err := r.QueryContext(ctx, countQuery)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, 0, web.NewRequestError(postgres.ErrNotFound, http.StatusNotFound)
	}
	if err != nil {
		return nil, 0, web.NewRequestError(errors.Wrap(err, "selecting tables"), http.StatusBadRequest)
	}

	count := 0

	for countRows.Next() {
		if err = countRows.Scan(&count); err != nil {
			return nil, 0, web.NewRequestError(errors.Wrap(err, "scanning user count"), http.StatusBadRequest)
		}
	}

	return list, count, nil
}

func (r Repository) AdminGetDetail(ctx context.Context, id int64) (entity.Table, error) {
	_, err := r.CheckClaims(ctx, auth.RoleAdmin)
	if err != nil {
		return entity.Table{}, err
	}

	var detail entity.Table

	err = r.NewSelect().Model(&detail).Where("id = ? AND deleted_at IS NULL", id).Scan(ctx)
	if err != nil {
		return entity.Table{}, err
	}

	return detail, nil
}

func (r Repository) AdminCreate(ctx context.Context, request tables.AdminCreateRequest) (tables.AdminCreateResponse, error) {
	claims, err := r.CheckClaims(ctx, auth.RoleAdmin)
	if err != nil {
		return tables.AdminCreateResponse{}, err
	}

	err = r.ValidateStruct(&request, "Number", "Capacity", "BranchID")
	if err != nil {
		return tables.AdminCreateResponse{}, err
	}

	response := tables.AdminCreateResponse{
		Number:    request.Number,
		Capacity:  request.Capacity,
		BranchID:  request.BranchID,
		HallID:    request.HallID,
		CreatedAt: time.Now(),
		CreatedBy: claims.UserId,
	}

	_, err = r.NewInsert().Model(&response).Exec(ctx)
	if err != nil {
		return tables.AdminCreateResponse{}, web.NewRequestError(errors.Wrap(err, "creating table"), http.StatusBadRequest)
	}

	return response, nil
}

func (r Repository) AdminUpdateAll(ctx context.Context, request tables.AdminUpdateRequest) error {
	claims, err := r.CheckClaims(ctx, auth.RoleAdmin)
	if err != nil {
		return err
	}

	if err = r.ValidateStruct(&request, "ID", "Number", "Status", "Capacity", "BranchID", "HallID"); err != nil {
		return err
	}

	q := r.NewUpdate().Table("tables").Where("deleted_at IS NULL AND id = ?", request.ID)

	q.Set("number = ?", request.Number)
	q.Set("status =?", request.Status)
	q.Set("capacity =?", request.Capacity)
	q.Set("branch_id =?", request.BranchID)
	q.Set("hall_id =?", request.HallID)
	q.Set("updated_at = ?", time.Now())
	q.Set("updated_by = ?", claims.UserId)

	_, err = q.Exec(ctx)
	if err != nil {
		return web.NewRequestError(errors.Wrap(err, "updating table"), http.StatusBadRequest)
	}

	return nil
}

func (r Repository) AdminUpdateColumns(ctx context.Context, request tables.AdminUpdateRequest) error {
	claims, err := r.CheckClaims(ctx, auth.RoleAdmin)
	if err != nil {
		return err
	}

	if err = r.ValidateStruct(&request, "ID"); err != nil {
		return err
	}

	q := r.NewUpdate().Table("tables").Where("deleted_at IS NULL AND id = ?", request.ID)

	if request.Number != nil {
		q.Set("number = ?", request.Number)
	}
	if request.Status != nil {
		q.Set("status = ?", request.Status)
	}
	if request.Capacity != nil {
		q.Set("capacity = ?", request.Capacity)
	}
	if request.BranchID != nil {
		q.Set("branch_id = ?", request.BranchID)
	}
	if request.HallID != nil {
		q.Set("hall_id = ?", request.HallID)
	}

	q.Set("updated_at = ?", time.Now())
	q.Set("updated_by = ?", claims.UserId)

	_, err = q.Exec(ctx)
	if err != nil {
		return web.NewRequestError(errors.Wrap(err, "updating table"), http.StatusBadRequest)
	}

	return nil
}

func (r Repository) AdminDelete(ctx context.Context, id int64) error {
	return r.DeleteRow(ctx, "tables", id, auth.RoleAdmin)
}

// @branch

func (r Repository) BranchGetList(ctx context.Context, filter tables.Filter) ([]tables.BranchGetList, int, error) {
	claims, err := r.CheckClaims(ctx, auth.RoleBranch)
	if err != nil {
		return nil, 0, err
	}

	table := "tables"
	whereQuery := fmt.Sprintf(`WHERE %s.deleted_at IS NULL and %s.branch_id = %d`, table, table, *claims.BranchID)

	if filter.HallID != nil {
		whereQuery += fmt.Sprintf(" AND %s.hall_id = %d", table, *filter.HallID)
	}

	countWhereQuery := whereQuery

	whereQuery += fmt.Sprintf(" ORDER BY %s.number asc", table)

	query, err := utils.SelectQuery(filter.Fields, filter.Joins, &table, &whereQuery)
	if err != nil {
		return nil, 0, errors.Wrap(err, "select query")
	}

	list := make([]tables.BranchGetList, 0)

	rows, err := r.QueryContext(ctx, query)
	if err != nil {
		return nil, 0, web.NewRequestError(errors.Wrap(err, "select table"), http.StatusInternalServerError)
	}

	err = r.ScanRows(ctx, rows, &list)
	if err != nil {
		return nil, 0, web.NewRequestError(errors.Wrap(err, "scanning tables"), http.StatusBadRequest)
	}

	countQuery := fmt.Sprintf(`
		SELECT
			count(id)
		FROM
		    %s
		%s
	`, table, countWhereQuery)

	countRows, err := r.QueryContext(ctx, countQuery)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, 0, web.NewRequestError(postgres.ErrNotFound, http.StatusNotFound)
	}
	if err != nil {
		return nil, 0, web.NewRequestError(errors.Wrap(err, "selecting tables"), http.StatusBadRequest)
	}

	count := 0

	for countRows.Next() {
		if err = countRows.Scan(&count); err != nil {
			return nil, 0, web.NewRequestError(errors.Wrap(err, "scanning user count"), http.StatusBadRequest)
		}
	}

	return list, count, nil
}

func (r Repository) BranchGetDetail(ctx context.Context, id int64) (entity.Table, error) {
	_, err := r.CheckClaims(ctx, auth.RoleBranch)
	if err != nil {
		return entity.Table{}, err
	}

	var detail entity.Table

	err = r.NewSelect().Model(&detail).Where("id = ? AND deleted_at IS NULL", id).Scan(ctx)
	if err != nil {
		return entity.Table{}, err
	}

	return detail, nil
}

func (r Repository) BranchCreate(ctx context.Context, request tables.BranchCreateRequest) (tables.BranchCreateResponse, error) {
	claims, err := r.CheckClaims(ctx, auth.RoleBranch)
	if err != nil {
		return tables.BranchCreateResponse{}, err
	}

	err = r.ValidateStruct(&request, "Number", "Capacity")
	if err != nil {
		return tables.BranchCreateResponse{}, err
	}

	response := tables.BranchCreateResponse{
		Number:    request.Number,
		Capacity:  request.Capacity,
		BranchID:  claims.BranchID,
		HallID:    request.HallID,
		CreatedAt: time.Now(),
		CreatedBy: claims.UserId,
	}

	_, err = r.NewInsert().Model(&response).Exec(ctx)
	if err != nil {
		return tables.BranchCreateResponse{}, web.NewRequestError(errors.Wrap(err, "creating table"), http.StatusBadRequest)
	}

	return response, nil
}

func (r Repository) BranchUpdateAll(ctx context.Context, request tables.BranchUpdateRequest) error {
	claims, err := r.CheckClaims(ctx, auth.RoleBranch)
	if err != nil {
		return err
	}

	if err = r.ValidateStruct(&request, "ID", "Number", "Status", "Capacity"); err != nil {
		return err
	}

	q := r.NewUpdate().Table("tables").Where("deleted_at IS NULL AND id = ?", request.ID)

	q.Set("number = ?", request.Number)
	q.Set("status =?", request.Status)
	q.Set("capacity =?", request.Capacity)
	q.Set("hall_id =?", request.HallID)
	q.Set("updated_at = ?", time.Now())
	q.Set("updated_by = ?", claims.BranchID)

	_, err = q.Exec(ctx)
	if err != nil {
		return web.NewRequestError(errors.Wrap(err, "updating table"), http.StatusBadRequest)
	}

	return nil
}

func (r Repository) BranchUpdateColumns(ctx context.Context, request tables.BranchUpdateRequest) error {
	_, err := r.CheckClaims(ctx, auth.RoleBranch)
	if err != nil {
		return err
	}

	if err = r.ValidateStruct(&request, "ID"); err != nil {
		return err
	}

	q := r.NewUpdate().Table("tables").Where("deleted_at IS NULL AND id = ?", request.ID)

	if request.Number != nil {
		q.Set("number = ?", request.Number)
	}
	if request.Status != nil {
		q.Set("status", request.Status)
	}
	if request.Capacity != nil {
		q.Set("capacity = ?", request.Capacity)
	}
	if request.HallID != nil {
		q.Set("hall_id = ?", request.HallID)
	}

	_, err = q.Exec(ctx)
	if err != nil {
		return web.NewRequestError(errors.Wrap(err, "updating table"), http.StatusBadRequest)
	}

	return nil
}

func (r Repository) BranchDelete(ctx context.Context, id int64) error {
	return r.DeleteRow(ctx, "tables", id, auth.RoleBranch)
}

// @admin

func (r Repository) CashierGetList(ctx context.Context, filter tables.Filter) ([]tables.CashierGetList, int, error) {
	claims, err := r.CheckClaims(ctx, auth.RoleCashier)
	if err != nil {
		return nil, 0, err
	}

	table := "tables"
	whereQuery := fmt.Sprintf(`WHERE %s.deleted_at IS NULL and %s.branch_id = %d`, table, table, *claims.BranchID)

	if filter.HallID != nil {
		whereQuery += fmt.Sprintf(" AND %s.hall_id = %d", table, *filter.HallID)
	}
	countWhereQuery := whereQuery

	whereQuery += fmt.Sprintf(" ORDER BY %s.number asc", table)

	query, err := utils.SelectQuery(filter.Fields, filter.Joins, &table, &whereQuery)
	if err != nil {
		return nil, 0, errors.Wrap(err, "select query")
	}

	list := make([]tables.CashierGetList, 0)

	rows, err := r.QueryContext(ctx, query)
	if err != nil {
		return nil, 0, web.NewRequestError(errors.Wrap(err, "select table"), http.StatusInternalServerError)
	}

	err = r.ScanRows(ctx, rows, &list)
	if err != nil {
		return nil, 0, web.NewRequestError(errors.Wrap(err, "scanning tables"), http.StatusBadRequest)
	}

	countQuery := fmt.Sprintf(`
		SELECT
			count(id)
		FROM
		    %s
		%s
	`, table, countWhereQuery)

	countRows, err := r.QueryContext(ctx, countQuery)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, 0, web.NewRequestError(postgres.ErrNotFound, http.StatusBadRequest)
	}
	if err != nil {
		return nil, 0, web.NewRequestError(errors.Wrap(err, "selecting tables"), http.StatusBadRequest)
	}

	count := 0

	for countRows.Next() {
		if err = countRows.Scan(&count); err != nil {
			return nil, 0, web.NewRequestError(errors.Wrap(err, "scanning user count"), http.StatusBadRequest)
		}
	}

	return list, count, nil
}

func (r Repository) CashierGetDetail(ctx context.Context, id int64) (entity.Table, error) {
	claims, err := r.CheckClaims(ctx, auth.RoleCashier)
	if err != nil {
		return entity.Table{}, err
	}

	var detail entity.Table

	err = r.NewSelect().Model(&detail).Where("id = ? AND deleted_at IS NULL AND branch_id = ?", id, *claims.BranchID).Scan(ctx)
	if err != nil {
		return entity.Table{}, err
	}

	return detail, nil
}

func (r Repository) CashierCreate(ctx context.Context, request tables.CashierCreateRequest) (tables.CashierCreateResponse, error) {
	claims, err := r.CheckClaims(ctx, auth.RoleCashier)
	if err != nil {
		return tables.CashierCreateResponse{}, err
	}

	err = r.ValidateStruct(&request, "Number", "Capacity")
	if err != nil {
		return tables.CashierCreateResponse{}, err
	}

	response := tables.CashierCreateResponse{
		Number:    request.Number,
		Capacity:  request.Capacity,
		BranchID:  claims.BranchID,
		HallID:    request.HallID,
		CreatedAt: time.Now(),
		CreatedBy: claims.UserId,
	}

	_, err = r.NewInsert().Model(&response).Exec(ctx)
	if err != nil {
		return tables.CashierCreateResponse{}, web.NewRequestError(errors.Wrap(err, "creating table"), http.StatusBadRequest)
	}

	return response, nil
}

func (r Repository) CashierUpdateAll(ctx context.Context, request tables.CashierUpdateRequest) error {
	claims, err := r.CheckClaims(ctx, auth.RoleCashier)
	if err != nil {
		return err
	}

	if err = r.ValidateStruct(&request, "ID", "Number", "Status", "Capacity"); err != nil {
		return err
	}

	q := r.NewUpdate().Table("tables").Where("deleted_at IS NULL AND id = ? AND branch_id = ? ", request.ID, *claims.BranchID)

	q.Set("number = ?", request.Number)
	q.Set("status =?", request.Status)
	q.Set("capacity =?", request.Capacity)
	q.Set("hall_id =?", request.HallID)
	q.Set("updated_at = ?", time.Now())
	q.Set("updated_by = ?", claims.UserId)

	_, err = q.Exec(ctx)
	if err != nil {
		return web.NewRequestError(errors.Wrap(err, "updating table"), http.StatusBadRequest)
	}

	return nil
}

func (r Repository) CashierUpdateColumns(ctx context.Context, request tables.CashierUpdateRequest) error {
	claims, err := r.CheckClaims(ctx, auth.RoleCashier)
	if err != nil {
		return err
	}

	if err = r.ValidateStruct(&request, "ID"); err != nil {
		return err
	}

	q := r.NewUpdate().Table("tables").Where("deleted_at IS NULL AND id = ? AND branch_id = ?", request.ID, *claims.BranchID)

	if request.Number != nil {
		q.Set("number = ?", request.Number)
	}
	if request.Status != nil {
		q.Set("status", request.Status)
	}
	if request.Capacity != nil {
		q.Set("capacity = ?", request.Capacity)
	}
	if request.HallID != nil {
		q.Set("hall_id = ?", request.HallID)
	}

	_, err = q.Exec(ctx)
	if err != nil {
		return web.NewRequestError(errors.Wrap(err, "updating table"), http.StatusBadRequest)
	}

	return nil
}

func (r Repository) CashierDelete(ctx context.Context, id int64) error {
	return r.DeleteRow(ctx, "tables", id, auth.RoleCashier)
}

// @waiter

func (r Repository) WaiterGetList(ctx context.Context, filter tables.Filter) ([]tables.WaiterGetListResponse, int, error) {
	claims, err := r.CheckClaims(ctx, auth.RoleWaiter)
	if err != nil {
		return nil, 0, err
	}

	where := fmt.Sprintf(`WHERE t.branch_id = '%d' AND t.deleted_at ISNULL`, *claims.BranchID)
	if filter.Search != nil {
		where += fmt.Sprintf(` AND (t.number::text ilike '%s')`, "%"+*filter.Search+"%")
	}

	if filter.HallID != nil {
		where += fmt.Sprintf(" AND t.hall_id = %d", *filter.HallID)
	}

	query := fmt.Sprintf(`SELECT 
								t.id, 
								t.number, 
								t.capacity
						  FROM tables t
							%s`, where)

	list := make([]tables.WaiterGetListResponse, 0)

	rows, err := r.QueryContext(ctx, query)
	if err != nil {
		return nil, 0, web.NewRequestError(errors.Wrap(err, "select table"), http.StatusInternalServerError)
	}

	err = r.ScanRows(ctx, rows, &list)
	if err != nil {
		return nil, 0, web.NewRequestError(errors.Wrap(err, "scanning tables"), http.StatusBadRequest)
	}

	for i := range list {
		query = fmt.Sprintf(`SELECT 
									o.id, 
									o.number, 
									o.client_count
							  FROM orders o 
							  WHERE o.status != 'CANCELLED' 
							     AND o.status != 'PAID' 
								 AND o.deleted_at isnull 
								 AND o.table_id='%d' ORDER BY o.created_at DESC`, list[i].ID)
		rows, err = r.QueryContext(ctx, query)
		if err != nil {
			return nil, 0, err
		}

		for rows.Next() {
			var row tables.WaiterOrder
			if err = rows.Scan(&row.Id, &row.Number, &row.ClientCount); err != nil {
				return nil, 0, err
			}

			if list[i].ClientCount == nil {
				list[i].ClientCount = row.ClientCount
			}

			list[i].Orders = append(list[i].Orders, row)
		}
	}

	countQuery := fmt.Sprintf(`
		SELECT
			count(t.id)
		FROM
		    tables t
		JOIN orders o ON t.id = o.table_id
		 %s
	`, where)

	countRows, err := r.QueryContext(ctx, countQuery)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, 0, web.NewRequestError(postgres.ErrNotFound, http.StatusNotFound)
	}
	if err != nil {
		return nil, 0, web.NewRequestError(errors.Wrap(err, "selecting tables"), http.StatusBadRequest)
	}

	count := 0

	for countRows.Next() {
		if err = countRows.Scan(&count); err != nil {
			return nil, 0, web.NewRequestError(errors.Wrap(err, "scanning user count"), http.StatusBadRequest)
		}
	}

	return list, count, nil
}

func NewRepository(DB *postgresql.Database) *Repository {
	return &Repository{DB}
}
