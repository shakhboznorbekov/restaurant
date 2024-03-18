package region

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/pkg/errors"
	"github.com/restaurant/foundation/web"
	"github.com/restaurant/internal/auth"
	"github.com/restaurant/internal/entity"
	"github.com/restaurant/internal/pkg/repository/postgresql"
	"github.com/restaurant/internal/repository/postgres"
	"github.com/restaurant/internal/service/region"
	"net/http"
	"time"
)

type Repository struct {
	*postgresql.Database
}

func (r Repository) SuperAdminGetList(ctx context.Context, filter region.Filter) ([]region.SuperAdminGetList, int, error) {
	_, err := r.CheckClaims(ctx, auth.RoleSuperAdmin)
	if err != nil {
		return nil, 0, err
	}

	lang := r.DefaultLang

	table := "regions"
	whereQuery := fmt.Sprintf(`WHERE %s.deleted_at IS NULL`, table)
	countWhereQuery := whereQuery

	if filter.Name != nil {
		whereQuery += fmt.Sprintf(" AND %s.name->>'%s' ilike '%s'", table, lang, "%"+*filter.Name+"%")
	}

	query := fmt.Sprintf(`
						SELECT 
						    id, 
						    name->>'%s' as name, 
						    code 
						FROM regions 
						%s`, lang, whereQuery)

	//query, err := utils.SelectQuery(filter.Fields, filter.Joins, &table, &whereQuery)
	//if err != nil {
	//	return nil, 0, errors.Wrap(err, "select query")
	//}

	list := make([]region.SuperAdminGetList, 0)

	rows, err := r.QueryContext(ctx, query)
	if err != nil {
		return nil, 0, web.NewRequestError(errors.Wrap(err, "select region"), http.StatusInternalServerError)
	}

	err = r.ScanRows(ctx, rows, &list)
	if err != nil {
		return nil, 0, web.NewRequestError(errors.Wrap(err, "scanning regions"), http.StatusBadRequest)
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
		return nil, 0, web.NewRequestError(errors.Wrap(err, "selecting regions"), http.StatusBadRequest)
	}

	count := 0

	for countRows.Next() {
		if err = countRows.Scan(&count); err != nil {
			return nil, 0, web.NewRequestError(errors.Wrap(err, "scanning user count"), http.StatusBadRequest)
		}
	}

	return list, count, nil
}

func (r Repository) SuperAdminGetDetail(ctx context.Context, id int64) (entity.Region, error) {
	_, err := r.CheckClaims(ctx, auth.RoleSuperAdmin)
	if err != nil {
		return entity.Region{}, err
	}

	var detail entity.Region

	err = r.NewSelect().Model(&detail).Where("id = ? AND deleted_at IS NULL", id).Scan(ctx)
	if err != nil {
		return entity.Region{}, err
	}

	return detail, nil
}

func (r Repository) SuperAdminCreate(ctx context.Context, request region.SuperAdminCreateRequest) (region.SuperAdminCreateResponse, error) {
	claims, err := r.CheckClaims(ctx, auth.RoleSuperAdmin)
	if err != nil {
		return region.SuperAdminCreateResponse{}, err
	}

	err = r.ValidateStruct(&request, "Name", "Code")
	if err != nil {
		return region.SuperAdminCreateResponse{}, err
	}

	response := region.SuperAdminCreateResponse{
		Name:      request.Name,
		Code:      request.Code,
		CreatedAt: time.Now(),
		CreatedBy: claims.UserId,
	}

	_, err = r.NewInsert().Model(&response).Exec(ctx)
	if err != nil {
		return region.SuperAdminCreateResponse{}, web.NewRequestError(errors.Wrap(err, "creating region"), http.StatusBadRequest)
	}

	return response, nil
}

func (r Repository) SuperAdminUpdateAll(ctx context.Context, request region.SuperAdminUpdateRequest) error {
	_, err := r.CheckClaims(ctx, auth.RoleSuperAdmin)
	if err != nil {
		return err
	}

	if err = r.ValidateStruct(&request, "ID", "Name", "Code"); err != nil {
		return err
	}

	q := r.NewUpdate().Table("regions").Where("deleted_at IS NULL AND id = ?", request.ID)

	q.Set("name = ?", request.Name)
	q.Set("code = ?", request.Code)

	_, err = q.Exec(ctx)
	if err != nil {
		return web.NewRequestError(errors.Wrap(err, "updating regions"), http.StatusBadRequest)
	}

	return nil
}

func (r Repository) SuperAdminUpdateColumns(ctx context.Context, request region.SuperAdminUpdateRequest) error {
	_, err := r.CheckClaims(ctx, auth.RoleSuperAdmin)
	if err != nil {
		return err
	}

	if err = r.ValidateStruct(&request, "ID"); err != nil {
		return err
	}

	q := r.NewUpdate().Table("regions").Where("deleted_at IS NULL AND id = ?", request.ID)

	if request.Name != nil {
		q.Set("name = ?", request.Name)
	}
	if request.Code != nil {
		q.Set("code = ?", request.Code)
	}

	_, err = q.Exec(ctx)
	if err != nil {
		return web.NewRequestError(errors.Wrap(err, "updating regions"), http.StatusBadRequest)
	}

	return nil
}

func (r Repository) SuperAdminDelete(ctx context.Context, id int64) error {
	return r.DeleteRow(ctx, "regions", id, auth.RoleSuperAdmin)
}

func NewRepository(DB *postgresql.Database) *Repository {
	return &Repository{DB}
}
