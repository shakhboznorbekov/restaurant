package measureUnit

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
	"restu-backend/internal/pkg/utils"
	"restu-backend/internal/repository/postgres"
	"restu-backend/internal/service/measureUnit"
	"time"
)

type Repository struct {
	*postgresql.Database
}

// @super-admin

func (r Repository) SuperAdminGetList(ctx context.Context, filter measureUnit.Filter) ([]measureUnit.SuperAdminGetList, int, error) {
	_, err := r.CheckClaims(ctx, auth.RoleSuperAdmin)
	if err != nil {
		return nil, 0, err
	}

	table := "measure_unit"
	whereQuery := fmt.Sprintf(`WHERE %s.deleted_at IS NULL`, table)
	countWhereQuery := whereQuery

	if filter.Name != nil {
		whereQuery += fmt.Sprintf(" AND %s.name ilike '%s'", table, "%"+*filter.Name+"%")
	}

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

	whereQuery += fmt.Sprintf(" %s %s", limitQuery, offsetQuery)

	query, err := utils.SelectQuery(filter.Fields, filter.Joins, &table, &whereQuery)
	if err != nil {
		return nil, 0, errors.Wrap(err, "select query")
	}

	list := make([]measureUnit.SuperAdminGetList, 0)

	rows, err := r.QueryContext(ctx, query)
	if err != nil {
		return nil, 0, web.NewRequestError(errors.Wrap(err, "select measure_unit"), http.StatusInternalServerError)
	}

	err = r.ScanRows(ctx, rows, &list)
	if err != nil {
		return nil, 0, web.NewRequestError(errors.Wrap(err, "scanning measure_unit"), http.StatusBadRequest)
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
		return nil, 0, web.NewRequestError(errors.Wrap(err, "selecting measure_unit"), http.StatusBadRequest)
	}

	count := 0

	for countRows.Next() {
		if err = countRows.Scan(&count); err != nil {
			return nil, 0, web.NewRequestError(errors.Wrap(err, "scanning measure_unit count"), http.StatusBadRequest)
		}
	}

	return list, count, nil
}

func (r Repository) SuperAdminGetDetail(ctx context.Context, id int64) (entity.MeasureUnit, error) {
	_, err := r.CheckClaims(ctx, auth.RoleSuperAdmin)
	if err != nil {
		return entity.MeasureUnit{}, err
	}

	var detail entity.MeasureUnit

	err = r.NewSelect().Model(&detail).Where("id = ? AND deleted_at IS NULL", id).Scan(ctx)
	if err != nil {
		return entity.MeasureUnit{}, err
	}

	return detail, nil
}

func (r Repository) SuperAdminCreate(ctx context.Context, request measureUnit.SuperAdminCreateRequest) (measureUnit.SuperAdminCreateResponse, error) {
	claims, err := r.CheckClaims(ctx, auth.RoleSuperAdmin)
	if err != nil {
		return measureUnit.SuperAdminCreateResponse{}, err
	}

	err = r.ValidateStruct(&request, "Name")
	if err != nil {
		return measureUnit.SuperAdminCreateResponse{}, err
	}

	response := measureUnit.SuperAdminCreateResponse{
		Name:      request.Name,
		CreatedAt: time.Now(),
		CreatedBy: claims.UserId,
	}

	_, err = r.NewInsert().Model(&response).Exec(ctx)
	if err != nil {
		return measureUnit.SuperAdminCreateResponse{}, web.NewRequestError(errors.Wrap(err, "creating measure_unit"), http.StatusBadRequest)
	}

	return response, nil
}

func (r Repository) SuperAdminUpdateAll(ctx context.Context, request measureUnit.SuperAdminUpdateRequest) error {
	_, err := r.CheckClaims(ctx, auth.RoleSuperAdmin)
	if err != nil {
		return err
	}

	if err = r.ValidateStruct(&request, "ID", "Name"); err != nil {
		return err
	}

	q := r.NewUpdate().Table("measure_unit").Where("deleted_at IS NULL AND id = ?", request.ID)

	q.Set("name = ?", request.Name)

	_, err = q.Exec(ctx)
	if err != nil {
		return web.NewRequestError(errors.Wrap(err, "updating measure_unit"), http.StatusBadRequest)
	}

	return nil
}

func (r Repository) SuperAdminUpdateColumns(ctx context.Context, request measureUnit.SuperAdminUpdateRequest) error {
	_, err := r.CheckClaims(ctx, auth.RoleSuperAdmin)
	if err != nil {
		return err
	}

	if err = r.ValidateStruct(&request, "ID"); err != nil {
		return err
	}

	q := r.NewUpdate().Table("measure_unit").Where("deleted_at IS NULL AND id = ?", request.ID)

	if request.Name != nil {
		q.Set("name = ?", request.Name)
	}

	_, err = q.Exec(ctx)
	if err != nil {
		return web.NewRequestError(errors.Wrap(err, "updating measure_unit"), http.StatusBadRequest)
	}

	return nil
}

func (r Repository) SuperAdminDelete(ctx context.Context, id int64) error {
	return r.DeleteRow(ctx, "measure_unit", id, auth.RoleSuperAdmin)
}

// @admin

func (r Repository) AdminGetList(ctx context.Context, filter measureUnit.Filter) ([]measureUnit.AdminGetList, int, error) {
	_, err := r.CheckClaims(ctx, auth.RoleAdmin)
	if err != nil {
		return nil, 0, err
	}

	table := "measure_unit"
	whereQuery := fmt.Sprintf(`WHERE %s.deleted_at IS NULL`, table)
	countWhereQuery := whereQuery

	if filter.Name != nil {
		whereQuery += fmt.Sprintf(" AND %s.name ilike '%s'", table, "%"+*filter.Name+"%")
	}

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

	whereQuery += fmt.Sprintf(" %s %s", limitQuery, offsetQuery)

	query, err := utils.SelectQuery(filter.Fields, filter.Joins, &table, &whereQuery)
	if err != nil {
		return nil, 0, errors.Wrap(err, "select query")
	}

	list := make([]measureUnit.AdminGetList, 0)

	rows, err := r.QueryContext(ctx, query)
	if err != nil {
		return nil, 0, web.NewRequestError(errors.Wrap(err, "select measure_unit"), http.StatusInternalServerError)
	}

	err = r.ScanRows(ctx, rows, &list)
	if err != nil {
		return nil, 0, web.NewRequestError(errors.Wrap(err, "scanning measure_unit"), http.StatusBadRequest)
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
		return nil, 0, web.NewRequestError(errors.Wrap(err, "selecting measure_unit"), http.StatusBadRequest)
	}

	count := 0

	for countRows.Next() {
		if err = countRows.Scan(&count); err != nil {
			return nil, 0, web.NewRequestError(errors.Wrap(err, "scanning measure_unit count"), http.StatusBadRequest)
		}
	}

	return list, count, nil
}

// @branch

func (r Repository) BranchGetList(ctx context.Context, filter measureUnit.Filter) ([]measureUnit.BranchGetList, int, error) {
	_, err := r.CheckClaims(ctx, auth.RoleBranch)
	if err != nil {
		return nil, 0, err
	}

	table := "measure_unit"
	whereQuery := fmt.Sprintf(`WHERE %s.deleted_at IS NULL`, table)
	countWhereQuery := whereQuery

	if filter.Name != nil {
		whereQuery += fmt.Sprintf(" AND %s.name ilike '%s'", table, "%"+*filter.Name+"%")
	}

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

	whereQuery += fmt.Sprintf(" %s %s", limitQuery, offsetQuery)

	query, err := utils.SelectQuery(filter.Fields, filter.Joins, &table, &whereQuery)
	if err != nil {
		return nil, 0, errors.Wrap(err, "select query")
	}

	list := make([]measureUnit.BranchGetList, 0)

	rows, err := r.QueryContext(ctx, query)
	if err != nil {
		return nil, 0, web.NewRequestError(errors.Wrap(err, "select measure_unit"), http.StatusInternalServerError)
	}

	err = r.ScanRows(ctx, rows, &list)
	if err != nil {
		return nil, 0, web.NewRequestError(errors.Wrap(err, "scanning measure_unit"), http.StatusBadRequest)
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
		return nil, 0, web.NewRequestError(errors.Wrap(err, "selecting measure_unit"), http.StatusBadRequest)
	}

	count := 0

	for countRows.Next() {
		if err = countRows.Scan(&count); err != nil {
			return nil, 0, web.NewRequestError(errors.Wrap(err, "scanning measure_unit count"), http.StatusBadRequest)
		}
	}

	return list, count, nil
}

func NewRepository(DB *postgresql.Database) *Repository {
	return &Repository{DB}
}
