package hall

import (
	"context"
	"database/sql"
	"fmt"
	"net/http"
	"restu-backend/foundation/web"
	"restu-backend/internal/auth"
	"restu-backend/internal/entity"
	"restu-backend/internal/pkg/repository/postgresql"
	"restu-backend/internal/repository/postgres"
	"restu-backend/internal/service/hall"
	"time"

	"github.com/pkg/errors"
)

type Repository struct {
	*postgresql.Database
}

func NewRepository(DB *postgresql.Database) *Repository {
	return &Repository{DB}
}

// @admin

func (r Repository) AdminGetList(ctx context.Context, filter halls.Filter) ([]halls.AdminGetList, int, error) {
	claims, err := r.CheckClaims(ctx, auth.RoleAdmin)
	if err != nil {
		return nil, 0, err
	}

	whereQuery := fmt.Sprintf(`WHERE h.deleted_at IS NULL 
	AND h.branch_id IN (select id from branches where restaurant_id = %d)`, *claims.RestaurantID)

	if filter.Search != nil {
		whereQuery += fmt.Sprintf(" AND h.name ilike '%s'", "%"+*filter.Search+"%")
	}

	if filter.BranchID != nil {
		whereQuery += fmt.Sprintf(" AND h.branch_id = %d", *filter.BranchID)
	}

	countWhereQuery := whereQuery

	var limitQuery, offsetQuery string
	if filter.Page != nil && filter.Limit != nil {
		offset := (*filter.Page - 1) * (*filter.Limit)
		filter.Offset = &offset
	}
	if filter.Limit != nil {
		limitQuery = fmt.Sprintf(" LIMIT '%d'", *filter.Limit)
	}
	if filter.Offset != nil {
		offsetQuery = fmt.Sprintf(" OFFSET '%d'", *filter.Offset)
	}
	orderQuery := " ORDER BY h.created_at desc"
	whereQuery += fmt.Sprintf("%s %s %s", orderQuery, limitQuery, offsetQuery)

	query := fmt.Sprintf(`
					SELECT 
					    h.id,
					    h.name,
					    h.branch_id,
					    b.name as branch
					FROM 
					    halls as h
					    left join branches as b on b.id = h.branch_id
					%s
	`, whereQuery)

	list := make([]halls.AdminGetList, 0)

	rows, err := r.QueryContext(ctx, query)
	if err != nil {
		return nil, 0, web.NewRequestError(errors.Wrap(err, "select hall"), http.StatusInternalServerError)
	}

	err = r.ScanRows(ctx, rows, &list)
	if err != nil {
		return nil, 0, web.NewRequestError(errors.Wrap(err, "scanning halls"), http.StatusBadRequest)
	}

	countQuery := fmt.Sprintf(`
		SELECT
			count(h.id)
		FROM
		    halls as h
		    left join branches as b on b.id = h.branch_id
		%s
	`, countWhereQuery)

	countRows, err := r.QueryContext(ctx, countQuery)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, 0, web.NewRequestError(postgres.ErrNotFound, http.StatusNotFound)
	}
	if err != nil {
		return nil, 0, web.NewRequestError(errors.Wrap(err, "selecting halls"), http.StatusBadRequest)
	}

	count := 0

	for countRows.Next() {
		if err = countRows.Scan(&count); err != nil {
			return nil, 0, web.NewRequestError(errors.Wrap(err, "scanning halls count"), http.StatusBadRequest)
		}
	}

	return list, count, nil
}

func (r Repository) AdminGetDetail(ctx context.Context, id int64) (entity.Hall, error) {
	_, err := r.CheckClaims(ctx, auth.RoleAdmin)
	if err != nil {
		return entity.Hall{}, err
	}

	var detail entity.Hall

	err = r.NewSelect().Model(&detail).Where("id = ? AND deleted_at IS NULL", id).Scan(ctx)
	if err != nil {
		if err == sql.ErrNoRows {
			return entity.Hall{}, web.NewRequestError(postgres.ErrNotFound, http.StatusBadRequest)
		}
		return entity.Hall{}, err
	}

	return detail, nil
}

func (r Repository) AdminCreate(ctx context.Context, request halls.AdminCreateRequest) (halls.AdminCreateResponse, error) {
	claims, err := r.CheckClaims(ctx, auth.RoleAdmin)
	if err != nil {
		return halls.AdminCreateResponse{}, err
	}

	err = r.ValidateStruct(&request, "Name", "BranchID")
	if err != nil {
		return halls.AdminCreateResponse{}, err
	}

	response := halls.AdminCreateResponse{
		Name:      request.Name,
		BranchID:  request.BranchID,
		CreatedAt: time.Now(),
		CreatedBy: claims.UserId,
	}

	_, err = r.NewInsert().Model(&response).Exec(ctx)
	if err != nil {
		return halls.AdminCreateResponse{}, web.NewRequestError(errors.Wrap(err, "creating hall"), http.StatusBadRequest)
	}

	return response, nil
}

func (r Repository) AdminUpdateAll(ctx context.Context, request halls.AdminUpdateRequest) error {
	claims, err := r.CheckClaims(ctx, auth.RoleAdmin)
	if err != nil {
		return err
	}

	if err = r.ValidateStruct(&request, "ID", "Name", "BranchID"); err != nil {
		return err
	}

	q := r.NewUpdate().Table("halls").Where("deleted_at IS NULL AND id = ?", request.ID)

	q.Set("name = ?", request.Name)
	q.Set("branch_id =?", request.BranchID)
	q.Set("updated_at = ?", time.Now())
	q.Set("updated_by = ?", claims.UserId)

	_, err = q.Exec(ctx)
	if err != nil {
		return web.NewRequestError(errors.Wrap(err, "updating hall"), http.StatusBadRequest)
	}

	return nil
}

func (r Repository) AdminUpdateColumns(ctx context.Context, request halls.AdminUpdateRequest) error {
	claims, err := r.CheckClaims(ctx, auth.RoleAdmin)
	if err != nil {
		return err
	}

	if err = r.ValidateStruct(&request, "ID"); err != nil {
		return err
	}

	q := r.NewUpdate().Table("halls").Where("deleted_at IS NULL AND id = ?", request.ID)

	if request.Name != nil {
		q.Set("name = ?", request.Name)
	}

	if request.BranchID != nil {
		q.Set("branch_id = ?", request.BranchID)
	}

	q.Set("updated_at = ?", time.Now())
	q.Set("updated_by = ?", claims.UserId)

	_, err = q.Exec(ctx)
	if err != nil {
		return web.NewRequestError(errors.Wrap(err, "updating hall"), http.StatusBadRequest)
	}

	return nil
}

func (r Repository) AdminDelete(ctx context.Context, id int64) error {
	return r.DeleteRow(ctx, "halls", id, auth.RoleAdmin)
}

// @branch

func (r Repository) BranchGetList(ctx context.Context, filter halls.Filter) ([]halls.BranchGetList, int, error) {
	claims, err := r.CheckClaims(ctx, auth.RoleBranch)
	if err != nil {
		return nil, 0, err
	}

	whereQuery := fmt.Sprintf(`WHERE h.deleted_at IS NULL AND h.branch_id = %d`, *claims.BranchID)

	if filter.Search != nil {
		whereQuery += fmt.Sprintf(" AND h.name ilike '%s'", "%"+*filter.Search+"%")
	}

	countWhereQuery := whereQuery

	var limitQuery, offsetQuery string
	if filter.Page != nil && filter.Limit != nil {
		offset := (*filter.Page - 1) * (*filter.Limit)
		filter.Offset = &offset
	}
	if filter.Limit != nil {
		limitQuery = fmt.Sprintf(" LIMIT '%d'", *filter.Limit)
	}
	if filter.Offset != nil {
		offsetQuery = fmt.Sprintf(" OFFSET '%d'", *filter.Offset)
	}

	orderQuery := " ORDER BY h.created_at desc"
	whereQuery += fmt.Sprintf("%s %s %s", orderQuery, limitQuery, offsetQuery)

	query := fmt.Sprintf(`
					SELECT 
					    h.id,
					    h.name
					FROM 
					    halls as h
					%s
	`, whereQuery)

	list := make([]halls.BranchGetList, 0)

	rows, err := r.QueryContext(ctx, query)
	if err != nil {
		return nil, 0, web.NewRequestError(errors.Wrap(err, "select halls"), http.StatusInternalServerError)
	}

	err = r.ScanRows(ctx, rows, &list)
	if err != nil {
		return nil, 0, web.NewRequestError(errors.Wrap(err, "scanning halls"), http.StatusInternalServerError)
	}

	countQuery := fmt.Sprintf(`
		SELECT
			count(h.id)
		FROM
		    halls as h
		%s
	`, countWhereQuery)

	countRows, err := r.QueryContext(ctx, countQuery)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, 0, web.NewRequestError(postgres.ErrNotFound, http.StatusNotFound)
	}
	if err != nil {
		return nil, 0, web.NewRequestError(errors.Wrap(err, "selecting halls"), http.StatusBadRequest)
	}

	count := 0

	for countRows.Next() {
		if err = countRows.Scan(&count); err != nil {
			return nil, 0, web.NewRequestError(errors.Wrap(err, "scanning halls count"), http.StatusBadRequest)
		}
	}

	return list, count, nil
}

func (r Repository) BranchGetDetail(ctx context.Context, id int64) (entity.Hall, error) {
	_, err := r.CheckClaims(ctx, auth.RoleBranch)
	if err != nil {
		return entity.Hall{}, err
	}

	var detail entity.Hall

	err = r.NewSelect().Model(&detail).Where("id = ? AND deleted_at IS NULL ", id).Scan(ctx)
	if err != nil {
		if err == sql.ErrNoRows {
			return entity.Hall{}, web.NewRequestError(postgres.ErrNotFound, http.StatusBadRequest)
		}
		return entity.Hall{}, err
	}

	return detail, nil
}

func (r Repository) BranchCreate(ctx context.Context, request halls.BranchCreateRequest) (halls.BranchCreateResponse, error) {
	claims, err := r.CheckClaims(ctx, auth.RoleBranch)
	if err != nil {
		return halls.BranchCreateResponse{}, err
	}

	err = r.ValidateStruct(&request, "Name")
	if err != nil {
		return halls.BranchCreateResponse{}, err
	}

	response := halls.BranchCreateResponse{
		Name:      request.Name,
		BranchID:  claims.BranchID,
		CreatedAt: time.Now(),
		CreatedBy: claims.UserId,
	}

	_, err = r.NewInsert().Model(&response).Exec(ctx)
	if err != nil {
		return halls.BranchCreateResponse{}, web.NewRequestError(errors.Wrap(err, "creating hall"), http.StatusBadRequest)
	}

	return response, nil
}

func (r Repository) BranchUpdateAll(ctx context.Context, request halls.BranchUpdateRequest) error {
	claims, err := r.CheckClaims(ctx, auth.RoleBranch)
	if err != nil {
		return err
	}

	if err = r.ValidateStruct(&request, "ID", "Name"); err != nil {
		return err
	}

	q := r.NewUpdate().Table("halls").Where("deleted_at IS NULL AND id = ?", request.ID)

	q.Set("name = ?", request.Name)
	q.Set("updated_at = ?", time.Now())
	q.Set("updated_by = ?", claims.UserId)

	_, err = q.Exec(ctx)
	if err != nil {
		return web.NewRequestError(errors.Wrap(err, "updating hall"), http.StatusBadRequest)
	}

	return nil
}

func (r Repository) BranchUpdateColumns(ctx context.Context, request halls.BranchUpdateRequest) error {
	claims, err := r.CheckClaims(ctx, auth.RoleBranch)
	if err != nil {
		return err
	}

	if err = r.ValidateStruct(&request, "ID"); err != nil {
		return err
	}

	q := r.NewUpdate().Table("halls").Where("deleted_at IS NULL AND id = ?", request.ID)

	if request.Name != nil {
		q.Set("name = ?", request.Name)
	}
	q.Set("updated_at = ?", time.Now())
	q.Set("updated_by = ?", claims.UserId)

	_, err = q.Exec(ctx)
	if err != nil {
		return web.NewRequestError(errors.Wrap(err, "updating hall"), http.StatusBadRequest)
	}

	return nil
}

func (r Repository) BranchDelete(ctx context.Context, id int64) error {
	return r.DeleteRow(ctx, "halls", id, auth.RoleBranch)
}

// @cashier

func (r Repository) CashierGetList(ctx context.Context, filter halls.Filter) ([]halls.CashierGetList, int, error) {
	claims, err := r.CheckClaims(ctx, auth.RoleCashier)
	if err != nil {
		return nil, 0, err
	}

	whereQuery := fmt.Sprintf(`WHERE h.deleted_at IS NULL AND h.branch_id = %d`, *claims.BranchID)

	if filter.Search != nil {
		whereQuery += fmt.Sprintf(" AND h.name ilike '%s'", "%"+*filter.Search+"%")
	}

	countWhereQuery := whereQuery

	var limitQuery, offsetQuery string
	if filter.Page != nil && filter.Limit != nil {
		offset := (*filter.Page - 1) * (*filter.Limit)
		filter.Offset = &offset
	}
	if filter.Limit != nil {
		limitQuery = fmt.Sprintf(" LIMIT '%d'", *filter.Limit)
	}
	if filter.Offset != nil {
		offsetQuery = fmt.Sprintf(" OFFSET '%d'", *filter.Offset)
	}

	orderQuery := " ORDER BY h.created_at desc"
	whereQuery += fmt.Sprintf("%s %s %s", orderQuery, limitQuery, offsetQuery)

	query := fmt.Sprintf(`
					SELECT 
					    h.id,
					    h.name
					FROM 
					    halls as h
					%s
	`, whereQuery)

	list := make([]halls.CashierGetList, 0)

	rows, err := r.QueryContext(ctx, query)
	if err != nil {
		return nil, 0, web.NewRequestError(errors.Wrap(err, "select halls"), http.StatusInternalServerError)
	}

	err = r.ScanRows(ctx, rows, &list)
	if err != nil {
		return nil, 0, web.NewRequestError(errors.Wrap(err, "scanning halls"), http.StatusInternalServerError)
	}

	countQuery := fmt.Sprintf(`
		SELECT
			count(h.id)
		FROM
		    halls as h
		%s
	`, countWhereQuery)

	countRows, err := r.QueryContext(ctx, countQuery)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, 0, web.NewRequestError(postgres.ErrNotFound, http.StatusNotFound)
	}
	if err != nil {
		return nil, 0, web.NewRequestError(errors.Wrap(err, "selecting halls"), http.StatusBadRequest)
	}

	count := 0

	for countRows.Next() {
		if err = countRows.Scan(&count); err != nil {
			return nil, 0, web.NewRequestError(errors.Wrap(err, "scanning halls count"), http.StatusBadRequest)
		}
	}

	return list, count, nil
}

func (r Repository) CashierGetDetail(ctx context.Context, id int64) (entity.Hall, error) {
	_, err := r.CheckClaims(ctx, auth.RoleCashier)
	if err != nil {
		return entity.Hall{}, err
	}

	var detail entity.Hall

	err = r.NewSelect().Model(&detail).Where("id = ? AND deleted_at IS NULL ", id).Scan(ctx)
	if err != nil {
		if err == sql.ErrNoRows {
			return entity.Hall{}, web.NewRequestError(postgres.ErrNotFound, http.StatusBadRequest)
		}
		return entity.Hall{}, err
	}

	return detail, nil
}

func (r Repository) CashierCreate(ctx context.Context, request halls.CashierCreateRequest) (halls.CashierCreateResponse, error) {
	claims, err := r.CheckClaims(ctx, auth.RoleCashier)
	if err != nil {
		return halls.CashierCreateResponse{}, err
	}

	err = r.ValidateStruct(&request, "Name")
	if err != nil {
		return halls.CashierCreateResponse{}, err
	}

	response := halls.CashierCreateResponse{
		Name:      request.Name,
		BranchID:  claims.BranchID,
		CreatedAt: time.Now(),
		CreatedBy: claims.UserId,
	}

	_, err = r.NewInsert().Model(&response).Exec(ctx)
	if err != nil {
		return halls.CashierCreateResponse{}, web.NewRequestError(errors.Wrap(err, "creating hall"), http.StatusBadRequest)
	}

	return response, nil
}

func (r Repository) CashierUpdateAll(ctx context.Context, request halls.CashierUpdateRequest) error {
	claims, err := r.CheckClaims(ctx, auth.RoleCashier)
	if err != nil {
		return err
	}

	if err = r.ValidateStruct(&request, "ID", "Name"); err != nil {
		return err
	}

	q := r.NewUpdate().Table("halls").Where("deleted_at IS NULL AND id = ?", request.ID)

	q.Set("name = ?", request.Name)
	q.Set("updated_at = ?", time.Now())
	q.Set("updated_by = ?", claims.UserId)

	_, err = q.Exec(ctx)
	if err != nil {
		return web.NewRequestError(errors.Wrap(err, "updating hall"), http.StatusBadRequest)
	}

	return nil
}

func (r Repository) CashierUpdateColumns(ctx context.Context, request halls.CashierUpdateRequest) error {
	claims, err := r.CheckClaims(ctx, auth.RoleCashier)
	if err != nil {
		return err
	}

	if err = r.ValidateStruct(&request, "ID"); err != nil {
		return err
	}

	q := r.NewUpdate().Table("halls").Where("deleted_at IS NULL AND id = ?", request.ID)

	if request.Name != nil {
		q.Set("name = ?", request.Name)
	}
	q.Set("updated_at = ?", time.Now())
	q.Set("updated_by = ?", claims.UserId)

	_, err = q.Exec(ctx)
	if err != nil {
		return web.NewRequestError(errors.Wrap(err, "updating hall"), http.StatusBadRequest)
	}

	return nil
}

func (r Repository) CashierDelete(ctx context.Context, id int64) error {
	return r.DeleteRow(ctx, "halls", id, auth.RoleCashier)
}

// @waiter

func (r Repository) WaiterGetList(ctx context.Context, filter halls.Filter) ([]halls.WaiterGetList, int, error) {
	claims, err := r.CheckClaims(ctx, auth.RoleWaiter)
	if err != nil {
		return nil, 0, err
	}

	whereQuery := fmt.Sprintf(`WHERE h.deleted_at IS NULL AND h.branch_id = %d`, *claims.BranchID)

	if filter.Search != nil {
		whereQuery += fmt.Sprintf(" AND h.name ilike '%s'", "%"+*filter.Search+"%")
	}

	countWhereQuery := whereQuery

	var limitQuery, offsetQuery string
	if filter.Page != nil && filter.Limit != nil {
		offset := (*filter.Page - 1) * (*filter.Limit)
		filter.Offset = &offset
	}
	if filter.Limit != nil {
		limitQuery = fmt.Sprintf(" LIMIT '%d'", *filter.Limit)
	}
	if filter.Offset != nil {
		offsetQuery = fmt.Sprintf(" OFFSET '%d'", *filter.Offset)
	}

	orderQuery := " ORDER BY h.created_at desc"
	whereQuery += fmt.Sprintf("%s %s %s", orderQuery, limitQuery, offsetQuery)

	query := fmt.Sprintf(`
					SELECT 
					    h.id,
					    h.name
					FROM 
					    halls as h
					%s
	`, whereQuery)

	list := make([]halls.WaiterGetList, 0)

	rows, err := r.QueryContext(ctx, query)
	if err != nil {
		return nil, 0, web.NewRequestError(errors.Wrap(err, "select halls"), http.StatusInternalServerError)
	}

	err = r.ScanRows(ctx, rows, &list)
	if err != nil {
		return nil, 0, web.NewRequestError(errors.Wrap(err, "scanning halls"), http.StatusInternalServerError)
	}

	countQuery := fmt.Sprintf(`
		SELECT
			count(h.id)
		FROM
		    halls as h
		%s
	`, countWhereQuery)

	countRows, err := r.QueryContext(ctx, countQuery)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, 0, web.NewRequestError(postgres.ErrNotFound, http.StatusNotFound)
	}
	if err != nil {
		return nil, 0, web.NewRequestError(errors.Wrap(err, "selecting halls"), http.StatusBadRequest)
	}

	count := 0

	for countRows.Next() {
		if err = countRows.Scan(&count); err != nil {
			return nil, 0, web.NewRequestError(errors.Wrap(err, "scanning halls count"), http.StatusBadRequest)
		}
	}

	return list, count, nil
}
