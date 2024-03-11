package foodCategory

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
	"github.com/restaurant/internal/service/foodCategory"
	"github.com/restaurant/internal/service/hashing"
	"net/http"
	"time"
)

type Repository struct {
	*postgresql.Database
}

// @super-admin

func (r Repository) SuperAdminGetList(ctx context.Context, filter foodCategory.Filter) ([]foodCategory.SuperAdminGetList, int, error) {
	_, err := r.CheckClaims(ctx, auth.RoleSuperAdmin)
	if err != nil {
		return nil, 0, err
	}

	table := "food_category"
	whereQuery := fmt.Sprintf(`WHERE %s.deleted_at IS NULL`, table)
	countWhereQuery := whereQuery

	if filter.Name != nil {
		whereQuery += fmt.Sprintf(" AND %s.name ilike '%s'", table, "%"+*filter.Name+"%")
	}

	whereQuery += fmt.Sprintf(` ORDER BY %s.created_at DESC`, table)

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

	list := make([]foodCategory.SuperAdminGetList, 0)

	rows, err := r.QueryContext(ctx, query)
	if err != nil {
		return nil, 0, web.NewRequestError(errors.Wrap(err, "select foodCategory"), http.StatusInternalServerError)
	}

	err = r.ScanRows(ctx, rows, &list)
	if err != nil {
		return nil, 0, web.NewRequestError(errors.Wrap(err, "scanning foodCategory"), http.StatusBadRequest)
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
		return nil, 0, web.NewRequestError(errors.Wrap(err, "selecting foodCategory"), http.StatusBadRequest)
	}

	count := 0

	for countRows.Next() {
		if err = countRows.Scan(&count); err != nil {
			return nil, 0, web.NewRequestError(errors.Wrap(err, "scanning foodCategory count"), http.StatusBadRequest)
		}
	}

	for k, v := range list {
		if v.Logo != nil {
			baseLink := hashing.GenerateHash(r.ServerBaseUrl, *v.Logo)
			list[k].Logo = &baseLink
		}
	}

	return list, count, nil
}

func (r Repository) SuperAdminGetDetail(ctx context.Context, id int64) (entity.FoodCategory, error) {
	_, err := r.CheckClaims(ctx, auth.RoleSuperAdmin)
	if err != nil {
		return entity.FoodCategory{}, err
	}

	var detail entity.FoodCategory

	err = r.NewSelect().Model(&detail).Where("id = ? AND deleted_at IS NULL", id).Scan(ctx)
	if err != nil {
		return entity.FoodCategory{}, err
	}

	if detail.Logo != nil {
		baseLink := hashing.GenerateHash(r.ServerBaseUrl, *detail.Logo)
		detail.Logo = &baseLink
	}

	return detail, nil
}

func (r Repository) SuperAdminCreate(ctx context.Context, request foodCategory.SuperAdminCreateRequest) (foodCategory.SuperAdminCreateResponse, error) {
	claims, err := r.CheckClaims(ctx, auth.RoleSuperAdmin)
	if err != nil {
		return foodCategory.SuperAdminCreateResponse{}, err
	}

	err = r.ValidateStruct(&request, "Name", "LogoLink")
	if err != nil {
		return foodCategory.SuperAdminCreateResponse{}, err
	}

	response := foodCategory.SuperAdminCreateResponse{
		Name:      request.Name,
		Logo:      request.LogoLink,
		Main:      request.Main,
		CreatedAt: time.Now(),
		CreatedBy: claims.UserId,
	}

	_, err = r.NewInsert().Model(&response).Exec(ctx)
	if err != nil {
		return foodCategory.SuperAdminCreateResponse{}, web.NewRequestError(errors.Wrap(err, "creating foodCategory"), http.StatusBadRequest)
	}

	return response, nil
}

func (r Repository) SuperAdminUpdateAll(ctx context.Context, request foodCategory.SuperAdminUpdateRequest) error {
	_, err := r.CheckClaims(ctx, auth.RoleSuperAdmin)
	if err != nil {
		return err
	}

	if err = r.ValidateStruct(&request, "ID", "Name", "LogoLink", "Main"); err != nil {
		return err
	}

	q := r.NewUpdate().Table("food_category").Where("deleted_at IS NULL AND id = ?", request.ID)

	q.Set("name = ?", request.Name)
	q.Set("main = ?", request.Main)
	q.Set("logo = ?", request.LogoLink)

	_, err = q.Exec(ctx)
	if err != nil {
		return web.NewRequestError(errors.Wrap(err, "updating food_category"), http.StatusBadRequest)
	}

	return nil
}

func (r Repository) SuperAdminUpdateColumns(ctx context.Context, request foodCategory.SuperAdminUpdateRequest) error {
	_, err := r.CheckClaims(ctx, auth.RoleSuperAdmin)
	if err != nil {
		return err
	}

	if err = r.ValidateStruct(&request, "ID"); err != nil {
		return err
	}

	q := r.NewUpdate().Table("food_category").Where("deleted_at IS NULL AND id = ?", request.ID)

	if request.Name != nil {
		q.Set("name = ?", request.Name)
	}

	if request.LogoLink != nil {
		q.Set("logo = ?", request.LogoLink)
	}

	if request.Main != nil {
		q.Set("main = ?", request.Main)
	}

	_, err = q.Exec(ctx)
	if err != nil {
		return web.NewRequestError(errors.Wrap(err, "updating food_category"), http.StatusBadRequest)
	}

	return nil
}

func (r Repository) SuperAdminDelete(ctx context.Context, id int64) error {
	return r.DeleteRow(ctx, "food_category", id, auth.RoleSuperAdmin)
}

// @client

//func (r Repository) ClientGetList(ctx context.Context, filter foodCategory.Filter) ([]foodCategory.ClientGetList, int, error) {
//	//_, err := r.CheckClaims(ctx, auth.RoleClient)
//	//if err != nil {
//	//	return nil, 0, err
//	//}
//
//	//table := "food_category"
//	whereQuery := fmt.Sprintf(`WHERE fc.deleted_at IS NULL AND fc.main = true`)
//	if filter.Name != nil {
//		whereQuery += fmt.Sprintf(" AND fc.name ilike '%s'", "%"+*filter.Name+"%")
//	}
//
//	countWhereQuery := whereQuery
//
//	var limitQuery, offsetQuery string
//	if filter.Page != nil && filter.Limit != nil {
//		offset := (*filter.Page - 1) * (*filter.Limit)
//		filter.Offset = &offset
//	}
//	if filter.Limit != nil {
//		limitQuery = fmt.Sprintf(" LIMIT '%d'", *filter.Limit)
//	}
//	if filter.Offset != nil {
//		offsetQuery = fmt.Sprintf(" OFFSET '%d'", *filter.Offset)
//	}
//
//	whereQuery += fmt.Sprintf(" GROUP BY fc.id, fc.name ORDER BY fc.name %s %s", limitQuery, offsetQuery)
//
//	query := fmt.Sprintf(`
//					SELECT
//					    fc.id,
//					    fc.name,
//					    fc.logo
//					FROM
//					    food_category as fc
//					LEFT OUTER JOIN foods as f ON f.category_id = fc.id
//					%s`, whereQuery)
//
//	//query, err := utils.SelectQuery(filter.Fields, filter.Joins, &table, &whereQuery)
//	//if err != nil {
//	//	return nil, 0, errors.Wrap(err, "select query")
//	//}
//
//	list := make([]foodCategory.ClientGetList, 0)
//
//	rows, err := r.QueryContext(ctx, query)
//	if err != nil {
//		return nil, 0, web.NewRequestError(errors.Wrap(err, "select foodCategory"), http.StatusInternalServerError)
//	}
//
//	err = r.ScanRows(ctx, rows, &list)
//	if err != nil {
//		return nil, 0, web.NewRequestError(errors.Wrap(err, "scanning foodCategory"), http.StatusBadRequest)
//	}
//
//	countQuery := fmt.Sprintf(`
//		SELECT
//			count(fc.id)
//		FROM
//		    food_category as fc
//		%s
//	`, countWhereQuery)
//
//	countRows, err := r.QueryContext(ctx, countQuery)
//	if errors.Is(err, sql.ErrNoRows) {
//		return nil, 0, web.NewRequestError(postgres.ErrNotFound, http.StatusNotFound)
//	}
//	if err != nil {
//		return nil, 0, web.NewRequestError(errors.Wrap(err, "selecting foodCategory"), http.StatusBadRequest)
//	}
//
//	count := 0
//
//	for countRows.Next() {
//		if err = countRows.Scan(&count); err != nil {
//			return nil, 0, web.NewRequestError(errors.Wrap(err, "scanning foodCategory count"), http.StatusBadRequest)
//		}
//	}
//
//	for k, v := range list {
//		if v.Logo != nil {
//			baseLink := hashing.GenerateHash(r.ServerBaseUrl, *v.Logo)
//			list[k].Logo = &baseLink
//		}
//	}
//
//	return list, count, nil
//}

// @branch

func (r Repository) BranchGetList(ctx context.Context, filter foodCategory.Filter) ([]foodCategory.BranchGetList, int, error) {
	_, err := r.CheckClaims(ctx, auth.RoleBranch)
	if err != nil {
		return nil, 0, err
	}

	table := "food_category"
	whereQuery := fmt.Sprintf(`WHERE %s.deleted_at IS NULL`, table)

	if filter.Name != nil {
		whereQuery += fmt.Sprintf(" AND %s.name ilike '%s'", table, "%"+*filter.Name+"%")
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

	whereQuery += fmt.Sprintf(" %s %s", limitQuery, offsetQuery)

	query, err := utils.SelectQuery(filter.Fields, filter.Joins, &table, &whereQuery)
	if err != nil {
		return nil, 0, errors.Wrap(err, "select query")
	}

	list := make([]foodCategory.BranchGetList, 0)

	rows, err := r.QueryContext(ctx, query)
	if err != nil {
		return nil, 0, web.NewRequestError(errors.Wrap(err, "select foodCategory"), http.StatusInternalServerError)
	}

	err = r.ScanRows(ctx, rows, &list)
	if err != nil {
		return nil, 0, web.NewRequestError(errors.Wrap(err, "scanning foodCategory"), http.StatusBadRequest)
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
		return nil, 0, web.NewRequestError(errors.Wrap(err, "selecting foodCategory"), http.StatusBadRequest)
	}

	count := 0

	for countRows.Next() {
		if err = countRows.Scan(&count); err != nil {
			return nil, 0, web.NewRequestError(errors.Wrap(err, "scanning foodCategory count"), http.StatusBadRequest)
		}
	}

	for k, v := range list {
		if v.Logo != nil {
			baseLink := r.ServerBaseUrl + *v.Logo
			list[k].Logo = &baseLink
		}
	}

	return list, count, nil
}

// @admin

func (r Repository) AdminGetList(ctx context.Context, filter foodCategory.Filter) ([]foodCategory.AdminGetList, int, error) {
	_, err := r.CheckClaims(ctx, auth.RoleAdmin)
	if err != nil {
		return nil, 0, err
	}

	table := "food_category"
	whereQuery := fmt.Sprintf(`WHERE %s.deleted_at IS NULL`, table)

	if filter.Name != nil {
		whereQuery += fmt.Sprintf(" AND %s.name ilike '%s'", table, "%"+*filter.Name+"%")
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

	whereQuery += fmt.Sprintf(" %s %s", limitQuery, offsetQuery)

	query, err := utils.SelectQuery(filter.Fields, filter.Joins, &table, &whereQuery)
	if err != nil {
		return nil, 0, errors.Wrap(err, "select query")
	}

	list := make([]foodCategory.AdminGetList, 0)

	rows, err := r.QueryContext(ctx, query)
	if err != nil {
		return nil, 0, web.NewRequestError(errors.Wrap(err, "select foodCategory"), http.StatusInternalServerError)
	}

	err = r.ScanRows(ctx, rows, &list)
	if err != nil {
		return nil, 0, web.NewRequestError(errors.Wrap(err, "scanning foodCategory"), http.StatusBadRequest)
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
		return nil, 0, web.NewRequestError(errors.Wrap(err, "selecting foodCategory"), http.StatusBadRequest)
	}

	count := 0

	for countRows.Next() {
		if err = countRows.Scan(&count); err != nil {
			return nil, 0, web.NewRequestError(errors.Wrap(err, "scanning foodCategory count"), http.StatusBadRequest)
		}
	}

	for k, v := range list {
		if v.Logo != nil {
			baseLink := r.ServerBaseUrl + *v.Logo
			list[k].Logo = &baseLink
		}
	}

	return list, count, nil
}

func (r Repository) WaiterGetList(ctx context.Context) ([]foodCategory.WaiterGetList, error) {
	claims, err := r.CheckClaims(ctx, auth.RoleWaiter)
	if err != nil {
		return nil, err
	}

	query := fmt.Sprintf(`SELECT 
								fc.id, 
								fc.name 
						  FROM food_category fc
						  WHERE fc.deleted_at isnull AND (select count(m.id) from menus m join foods f on f.id = m.food_id where f.category_id=fc.id and m.branch_id='%d') != 0`, *claims.BranchID)

	rows, err := r.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}

	var response []foodCategory.WaiterGetList
	if err = r.ScanRows(ctx, rows, &response); err != nil {
		return nil, err
	}

	return response, nil
}

func NewRepository(DB *postgresql.Database) *Repository {
	return &Repository{DB}
}
