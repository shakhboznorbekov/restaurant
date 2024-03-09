package restaurant

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
	"github.com/restaurant/internal/service/hashing"
	"github.com/restaurant/internal/service/restaurant"
	"net/http"
	"time"
)

type Repository struct {
	*postgresql.Database
}

// @super-admin

func (r Repository) SuperAdminGetList(ctx context.Context, filter restaurant.Filter) ([]restaurant.SuperAdminGetList, int, error) {
	_, err := r.CheckClaims(ctx, auth.RoleSuperAdmin)
	if err != nil {
		return nil, 0, err
	}

	table := "restaurants"
	whereQuery := fmt.Sprintf(`WHERE %s.deleted_at IS NULL`, table)

	if filter.Name != nil {
		whereQuery += fmt.Sprintf(" AND %s.name ilike '%s'", table, "%"+*filter.Name+"%")
	}
	countWhereQuery := whereQuery

	whereQuery += fmt.Sprintf(` ORDER BY %s.created_at DESC`, table)

	if filter.Limit != nil {
		whereQuery += fmt.Sprintf(" Limit %d", *filter.Limit)
	}
	if filter.Offset != nil {
		whereQuery += fmt.Sprintf(" Offset %d", *filter.Offset)
	}

	query, err := utils.SelectQuery(filter.Fields, filter.Joins, &table, &whereQuery)
	if err != nil {
		return nil, 0, errors.Wrap(err, "select query")
	}

	list := make([]restaurant.SuperAdminGetList, 0)

	rows, err := r.QueryContext(ctx, query)
	if err != nil {
		return nil, 0, web.NewRequestError(errors.Wrap(err, "select user"), http.StatusInternalServerError)
	}

	err = r.ScanRows(ctx, rows, &list)
	if err != nil {
		return nil, 0, web.NewRequestError(errors.Wrap(err, "scanning restaurants"), http.StatusBadRequest)
	}

	for i := range list {
		if list[i].Logo != nil {
			link := hashing.GenerateHash(r.ServerBaseUrl, *list[i].Logo)
			list[i].Logo = &link
		}
		if list[i].MiniLogo != nil {
			link := hashing.GenerateHash(r.ServerBaseUrl, *list[i].MiniLogo)
			list[i].MiniLogo = &link
		}
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
		return nil, 0, web.NewRequestError(errors.Wrap(err, "selecting restaurants"), http.StatusBadRequest)
	}

	count := 0

	for countRows.Next() {
		if err = countRows.Scan(&count); err != nil {
			return nil, 0, web.NewRequestError(errors.Wrap(err, "scanning user count"), http.StatusBadRequest)
		}
	}

	return list, count, nil
}

func (r Repository) SuperAdminGetDetail(ctx context.Context, id int64) (entity.Restaurant, error) {
	_, err := r.CheckClaims(ctx, auth.RoleSuperAdmin)
	if err != nil {
		return entity.Restaurant{}, err
	}

	var detail entity.Restaurant

	err = r.NewSelect().Model(&detail).Where("id = ? AND deleted_at IS NULL", id).Scan(ctx)
	if err != nil {
		return entity.Restaurant{}, err
	}

	if detail.Logo != nil {
		link := hashing.GenerateHash(r.ServerBaseUrl, *detail.Logo)
		detail.Logo = &link
	}
	if detail.MiniLogo != nil {
		link := hashing.GenerateHash(r.ServerBaseUrl, *detail.MiniLogo)
		detail.MiniLogo = &link
	}

	return detail, nil
}

func (r Repository) SuperAdminCreate(ctx context.Context, request restaurant.SuperAdminCreateRequest) (restaurant.SuperAdminCreateResponse, error) {
	claims, err := r.CheckClaims(ctx, auth.RoleSuperAdmin)
	if err != nil {
		return restaurant.SuperAdminCreateResponse{}, err
	}

	err = r.ValidateStruct(&request, "Name")
	if err != nil {
		return restaurant.SuperAdminCreateResponse{}, err
	}

	response := restaurant.SuperAdminCreateResponse{
		Name:       request.Name,
		Logo:       request.LogoLink,
		MiniLogo:   request.MiniLogoLink,
		WebsiteUrl: request.WebsiteUrl,
		CreatedAt:  time.Now(),
		CreatedBy:  claims.UserId,
	}

	_, err = r.NewInsert().Model(&response).Exec(ctx)
	if err != nil {
		return restaurant.SuperAdminCreateResponse{}, web.NewRequestError(errors.Wrap(err, "creating restaurant"), http.StatusBadRequest)
	}

	return response, nil
}

func (r Repository) SuperAdminUpdateAll(ctx context.Context, request restaurant.SuperAdminUpdateRequest) error {
	claims, err := r.CheckClaims(ctx, auth.RoleSuperAdmin)
	if err != nil {
		return err
	}

	if err = r.ValidateStruct(&request, "ID", "Name", "LogoLink"); err != nil {
		return err
	}

	q := r.NewUpdate().Table("restaurants").Where("deleted_at IS NULL AND id = ?", request.ID)

	q.Set("name = ?", request.Name)
	q.Set("logo = ?", request.LogoLink)
	q.Set("website_url = ?", request.WebsiteUrl)
	q.Set("mini_logo = ?", request.MiniLogoLink)
	q.Set("updated_at = ?", time.Now())
	q.Set("updated_by = ?", claims.UserId)

	_, err = q.Exec(ctx)
	if err != nil {
		return web.NewRequestError(errors.Wrap(err, "updating user"), http.StatusBadRequest)
	}

	return nil
}

func (r Repository) SuperAdminUpdateColumns(ctx context.Context, request restaurant.SuperAdminUpdateRequest) error {
	claims, err := r.CheckClaims(ctx, auth.RoleSuperAdmin)
	if err != nil {
		return err
	}

	if err = r.ValidateStruct(&request, "ID"); err != nil {
		return err
	}

	q := r.NewUpdate().Table("restaurants").Where("deleted_at IS NULL AND id = ?", request.ID)

	if request.Name != nil {
		q.Set("name = ?", request.Name)
	}
	if request.LogoLink != nil {
		q.Set("logo = ?", request.LogoLink)
	}
	if request.WebsiteUrl != nil {
		q.Set("website_url = ?", request.WebsiteUrl)
	}
	if request.MiniLogoLink != nil {
		q.Set("mini_logo = ?", request.MiniLogoLink)
	}

	q.Set("updated_at = ?", time.Now())
	q.Set("updated_by = ?", claims.UserId)

	_, err = q.Exec(ctx)
	if err != nil {
		return web.NewRequestError(errors.Wrap(err, "updating user"), http.StatusBadRequest)
	}

	return nil
}

func (r Repository) SuperAdminDelete(ctx context.Context, id int64) error {
	return r.DeleteRow(ctx, "restaurants", id, auth.RoleSuperAdmin)
}

// @site

func (r Repository) SiteGetList(ctx context.Context) ([]restaurant.SiteGetListResponse, int, error) {
	list := make([]restaurant.SiteGetListResponse, 0)
	where := fmt.Sprintf(`WHERE deleted_at isnull`)
	query := fmt.Sprintf(`SELECT generate_single_hash(logo), generate_single_hash(mini_logo), website_url FROM restaurants %s`, where)

	rows, err := r.QueryContext(ctx, query)
	if err != nil {
		return nil, 0, web.NewRequestError(errors.Wrap(err, "selecting restaurants"), http.StatusInternalServerError)
	}

	err = r.ScanRows(ctx, rows, &list)
	if err != nil {
		return nil, 0, web.NewRequestError(errors.Wrap(err, "scanning restaurants"), http.StatusInternalServerError)
	}

	var count int
	countQuery := fmt.Sprintf(`SELECT count(id) FROM restaurants %s`, where)
	if err = r.QueryRowContext(ctx, countQuery).Scan(&count); err != nil {
		return nil, 0, web.NewRequestError(errors.Wrap(err, "scanning restaurants count"), http.StatusInternalServerError)
	}

	return list, count, nil
}

func NewRepository(DB *postgresql.Database) *Repository {
	return &Repository{DB}
}
