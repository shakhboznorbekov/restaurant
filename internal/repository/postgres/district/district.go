package district

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
	"github.com/restaurant/internal/service/district"
	"net/http"
	"time"
)

type Repository struct {
	*postgresql.Database
}

func (r Repository) SuperAdminGetList(ctx context.Context, filter district.Filter) ([]district.SuperAdminGetList, int, error) {
	_, err := r.CheckClaims(ctx, auth.RoleSuperAdmin)
	if err != nil {
		return nil, 0, err
	}

	lang := r.DefaultLang

	whereQuery := fmt.Sprintf(`WHERE d.deleted_at IS NULL`)
	countWhereQuery := whereQuery

	if filter.Name != nil {
		whereQuery += fmt.Sprintf(" AND d.name->>'%s' ilike '%s'", lang, "%"+*filter.Name+"%")
	}

	query := fmt.Sprintf(`
						SELECT 
						    d.id, 
						    d.name->>'%s' as name, 
						    d.code,
						    d.region_id,
						    r.name->>'%s' as region
						FROM districts as d
						LEFT OUTER JOIN regions as r ON r.id = d.region_id
						%s`, lang, lang, whereQuery)
	list := make([]district.SuperAdminGetList, 0)

	rows, err := r.QueryContext(ctx, query)
	if err != nil {
		return nil, 0, web.NewRequestError(errors.Wrap(err, "select district"), http.StatusInternalServerError)
	}

	err = r.ScanRows(ctx, rows, &list)
	if err != nil {
		return nil, 0, web.NewRequestError(errors.Wrap(err, "scanning districts"), http.StatusBadRequest)
	}

	countQuery := fmt.Sprintf(`
		SELECT
			count(d.id)
		FROM
		    districts as d
		%s
	`, countWhereQuery)

	countRows, err := r.QueryContext(ctx, countQuery)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, 0, web.NewRequestError(postgres.ErrNotFound, http.StatusNotFound)
	}
	if err != nil {
		return nil, 0, web.NewRequestError(errors.Wrap(err, "selecting districts"), http.StatusBadRequest)
	}

	count := 0

	for countRows.Next() {
		if err = countRows.Scan(&count); err != nil {
			return nil, 0, web.NewRequestError(errors.Wrap(err, "scanning district count"), http.StatusBadRequest)
		}
	}

	return list, count, nil
}

func (r Repository) SuperAdminGetDetail(ctx context.Context, id int64) (entity.District, error) {
	_, err := r.CheckClaims(ctx, auth.RoleSuperAdmin)
	if err != nil {
		return entity.District{}, err
	}

	var detail entity.District

	err = r.NewSelect().Model(&detail).Where("id = ? AND deleted_at IS NULL", id).Scan(ctx)
	if err != nil {
		return entity.District{}, err
	}

	return detail, nil
}

func (r Repository) SuperAdminCreate(ctx context.Context, request district.SuperAdminCreateRequest) (district.SuperAdminCreateResponse, error) {
	claims, err := r.CheckClaims(ctx, auth.RoleSuperAdmin)
	if err != nil {
		return district.SuperAdminCreateResponse{}, err
	}

	err = r.ValidateStruct(&request, "Name", "Code", "RegionID")
	if err != nil {
		return district.SuperAdminCreateResponse{}, err
	}

	response := district.SuperAdminCreateResponse{
		Name:      request.Name,
		RegionID:  request.RegionID,
		Code:      request.Code,
		CreatedAt: time.Now(),
		CreatedBy: claims.UserId,
	}

	_, err = r.NewInsert().Model(&response).Exec(ctx)
	if err != nil {
		return district.SuperAdminCreateResponse{}, web.NewRequestError(errors.Wrap(err, "creating district"), http.StatusBadRequest)
	}

	return response, nil
}

func (r Repository) SuperAdminUpdateAll(ctx context.Context, request district.SuperAdminUpdateRequest) error {
	_, err := r.CheckClaims(ctx, auth.RoleSuperAdmin)
	if err != nil {
		return err
	}

	if err = r.ValidateStruct(&request, "ID", "Name", "Code", "RegionID"); err != nil {
		return err
	}

	q := r.NewUpdate().Table("districts").Where("deleted_at IS NULL AND id = ?", request.ID)

	q.Set("name = ?", request.Name)
	q.Set("code = ?", request.Code)
	q.Set("region_id = ?", request.RegionID)

	_, err = q.Exec(ctx)
	if err != nil {
		return web.NewRequestError(errors.Wrap(err, "updating districts"), http.StatusBadRequest)
	}

	return nil
}

func (r Repository) SuperAdminUpdateColumns(ctx context.Context, request district.SuperAdminUpdateRequest) error {
	_, err := r.CheckClaims(ctx, auth.RoleSuperAdmin)
	if err != nil {
		return err
	}

	if err = r.ValidateStruct(&request, "ID"); err != nil {
		return err
	}

	q := r.NewUpdate().Table("districts").Where("deleted_at IS NULL AND id = ?", request.ID)

	if request.Name != nil {
		q.Set("name = ?", request.Name)
	}
	if request.Code != nil {
		q.Set("code = ?", request.Code)
	}
	if request.RegionID != nil {
		q.Set("region_id = ?", request.RegionID)
	}

	_, err = q.Exec(ctx)
	if err != nil {
		return web.NewRequestError(errors.Wrap(err, "updating districts"), http.StatusBadRequest)
	}

	return nil
}

func (r Repository) SuperAdminDelete(ctx context.Context, id int64) error {
	return r.DeleteRow(ctx, "districts", id, auth.RoleSuperAdmin)
}

func NewRepository(DB *postgresql.Database) *Repository {
	return &Repository{DB}
}
