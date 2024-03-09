package restaurant_category

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/pkg/errors"
	"github.com/restaurant/foundation/web"
	"github.com/restaurant/internal/auth"
	"github.com/restaurant/internal/entity"
	"github.com/restaurant/internal/pkg/repository/postgresql"
	"github.com/restaurant/internal/pkg/utils"
	"github.com/restaurant/internal/repository/postgres"
	"github.com/restaurant/internal/service/restaurant_category"
	"net/http"
	"time"
)

type Repository struct {
	*postgresql.Database
}

// @super-admin

func (r Repository) SuperAdminGetList(ctx context.Context, filter restaurant_category.Filter) ([]restaurant_category.SuperAdminGetList, int, error) {
	_, err := r.CheckClaims(ctx, auth.RoleSuperAdmin)
	if err != nil {
		return nil, 0, err
	}

	table := "restaurant_category"
	whereQuery := fmt.Sprintf(` WHERE %s.deleted_at IS NULL`, table)
	countWhereQuery := whereQuery

	if filter.Name != nil {
		whereQuery += fmt.Sprintf(" AND %s.name ILIKE '%s'", table, "%"+*filter.Name+"%")
	}

	whereQuery += fmt.Sprintf(` ORDER BY %s.created_at DESC`, table)

	query, err := utils.SelectQuery(filter.Fields, filter.Joins, &table, &whereQuery)
	if err != nil {
		return nil, 0, errors.Wrap(err, "select query")
	}

	list := make([]restaurant_category.SuperAdminGetList, 0)

	rows, err := r.QueryContext(ctx, query)
	if err != nil {
		return nil, 0, web.NewRequestError(errors.Wrap(err, "select restaurantCategories"), http.StatusInternalServerError)
	}

	err = r.ScanRows(ctx, rows, &list)
	if err != nil {
		return nil, 0, web.NewRequestError(errors.Wrap(err, "scanning restaurantCategories"), http.StatusBadRequest)
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
		return nil, 0, web.NewRequestError(errors.Wrap(err, "selecting restaurantCategories"), http.StatusBadRequest)
	}

	count := 0

	for countRows.Next() {
		if err = countRows.Scan(&count); err != nil {
			return nil, 0, web.NewRequestError(errors.Wrap(err, "scanning restaurantCategory count"), http.StatusBadRequest)
		}
	}
	return list, count, nil
}

func (r Repository) SuperAdminGetDetail(ctx context.Context, id int64) (entity.RestaurantCategory, error) {
	_, err := r.CheckClaims(ctx, auth.RoleSuperAdmin)
	if err != nil {
		return entity.RestaurantCategory{}, err
	}

	var detail entity.RestaurantCategory

	err = r.NewSelect().Model(&detail).Where("id = ? AND deleted_at IS NULL", id).Scan(ctx)
	if err != nil {
		return entity.RestaurantCategory{}, err
	}

	return detail, nil
}

func (r Repository) SuperAdminCreate(ctx context.Context, request restaurant_category.SuperAdminCreateRequest) (restaurant_category.SuperAdminCreateResponse, error) {
	claims, err := r.CheckClaims(ctx, auth.RoleSuperAdmin)
	if err != nil {
		return restaurant_category.SuperAdminCreateResponse{}, err
	}

	err = r.ValidateStruct(&request, "Name")
	if err != nil {
		return restaurant_category.SuperAdminCreateResponse{}, err
	}

	response := restaurant_category.SuperAdminCreateResponse{
		Name:      request.Name,
		CreatedAt: time.Now(),
		CreatedBy: claims.UserId,
	}

	_, err = r.NewInsert().Model(&response).Exec(ctx)
	if err != nil {
		return restaurant_category.SuperAdminCreateResponse{}, web.NewRequestError(errors.Wrap(err, "creating restaurant_category"), http.StatusBadRequest)
	}

	return response, nil
}

func (r Repository) SuperAdminUpdateAll(ctx context.Context, request restaurant_category.SuperAdminUpdateRequest) error {
	_, err := r.CheckClaims(ctx, auth.RoleSuperAdmin)
	if err != nil {
		return err
	}

	if err = r.ValidateStruct(&request, "ID", "Name"); err != nil {
		return nil
	}

	q := r.NewUpdate().Table("restaurant_category").Where(" deleted_at IS NULL AND id = ?", request.ID)

	q.Set(" name = ?", request.Name)

	_, err = q.Exec(ctx)
	if err != nil {
		return web.NewRequestError(errors.Wrap(err, "updating restaurant_category"), http.StatusBadRequest)
	}

	return nil
}

func (r Repository) SuperAdminUpdateColumns(ctx context.Context, request restaurant_category.SuperAdminUpdateRequest) error {
	_, err := r.CheckClaims(ctx, auth.RoleSuperAdmin)
	if err != nil {
		return err
	}

	if err = r.ValidateStruct(&request, "ID"); err != nil {
		return err
	}

	q := r.NewUpdate().Table("restaurant_category").Where("deleted_at IS NULL AND id = ?", request.ID)

	if request.Name != nil {
		q.Set("name = ?", request.Name)
	}

	_, err = q.Exec(ctx)
	if err != nil {
		return web.NewRequestError(errors.Wrap(err, "updating restaurant_category"), http.StatusBadRequest)
	}

	return nil
}

func (r Repository) SuperAdminDelete(ctx context.Context, id int64) error {
	return r.DeleteRow(ctx, "restaurant_category", id, auth.RoleSuperAdmin)
}

// @admin

func (r Repository) AdminGetList(ctx context.Context, filter restaurant_category.Filter) ([]restaurant_category.AdminGetList, int, error) {
	_, err := r.CheckClaims(ctx, auth.RoleAdmin)
	if err != nil {
		return nil, 0, err
	}

	table := "restaurant_category"
	whereQuery := fmt.Sprintf(`WHERE %s.deleted_at IS NULL`, table)
	countWhereQuery := whereQuery

	if filter.Name != nil {
		whereQuery += fmt.Sprintf(" AND %s.name ilike '%s'", table, "%"+*filter.Name+"%")
	}

	query, err := utils.SelectQuery(filter.Fields, filter.Joins, &table, &whereQuery)
	if err != nil {
		return nil, 0, errors.Wrap(err, "select query")
	}

	list := make([]restaurant_category.AdminGetList, 0)

	rows, err := r.QueryContext(ctx, query)
	if err != nil {
		return nil, 0, web.NewRequestError(errors.Wrap(err, "select restaurantCategories"), http.StatusInternalServerError)
	}

	err = r.ScanRows(ctx, rows, &list)
	if err != nil {
		return nil, 0, web.NewRequestError(errors.Wrap(err, "scanning restaurantCategories"), http.StatusBadRequest)
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
		return nil, 0, web.NewRequestError(errors.Wrap(err, "selecting restaurantCategories"), http.StatusBadRequest)
	}

	count := 0

	for countRows.Next() {
		if err = countRows.Scan(&count); err != nil {
			return nil, 0, web.NewRequestError(errors.Wrap(err, "scanning restaurantCategory count"), http.StatusBadRequest)
		}
	}

	return list, count, nil
}

// @site

func (r Repository) SiteGetList(ctx context.Context) ([]restaurant_category.SiteGetListResponse, int, error) {
	list := make([]restaurant_category.SiteGetListResponse, 0)
	where := fmt.Sprintf(`WHERE deleted_at isnull`)
	query := fmt.Sprintf(`SELECT name, generate_single_hash(photo) as photo FROM restaurant_category %s`, where)

	rows, err := r.QueryContext(ctx, query)
	if err != nil {
		return nil, 0, web.NewRequestError(errors.Wrap(err, "selecting restaurants"), http.StatusInternalServerError)
	}

	err = r.ScanRows(ctx, rows, &list)
	if err != nil {
		return nil, 0, web.NewRequestError(errors.Wrap(err, "scanning restaurants"), http.StatusInternalServerError)
	}

	var count int
	countQuery := fmt.Sprintf(`SELECT count(id) FROM restaurant_category %s`, where)
	if err = r.QueryRowContext(ctx, countQuery).Scan(&count); err != nil {
		return nil, 0, web.NewRequestError(errors.Wrap(err, "scanning restaurants count"), http.StatusInternalServerError)
	}

	return list, count, nil
}

func NewRepository(DB *postgresql.Database) *Repository {
	return &Repository{DB}
}
