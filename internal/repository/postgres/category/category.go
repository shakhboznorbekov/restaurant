package category

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
	"restu-backend/internal/service/category"
	"restu-backend/internal/service/hashing"
	"time"
)

type Repository struct {
	*postgresql.Database
}

// @super-admin

func (r Repository) SuperAdminGetList(ctx context.Context, filter category.Filter) ([]category.SuperAdminGetList, int, error) {
	_, err := r.CheckClaims(ctx, auth.RoleSuperAdmin)
	if err != nil {
		return nil, 0, err
	}

	table := "categories"
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

	list := make([]category.SuperAdminGetList, 0)

	rows, err := r.QueryContext(ctx, query)
	if err != nil {
		return nil, 0, web.NewRequestError(errors.Wrap(err, "select menu_category"), http.StatusInternalServerError)
	}

	err = r.ScanRows(ctx, rows, &list)
	if err != nil {
		return nil, 0, web.NewRequestError(errors.Wrap(err, "scanning menu_category"), http.StatusBadRequest)
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
		return nil, 0, web.NewRequestError(errors.Wrap(err, "selecting menu_category"), http.StatusBadRequest)
	}

	count := 0

	for countRows.Next() {
		if err = countRows.Scan(&count); err != nil {
			return nil, 0, web.NewRequestError(errors.Wrap(err, "scanning menu_category count"), http.StatusBadRequest)
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

func (r Repository) SuperAdminGetDetail(ctx context.Context, id int64) (entity.Category, error) {
	_, err := r.CheckClaims(ctx, auth.RoleSuperAdmin)
	if err != nil {
		return entity.Category{}, err
	}

	var detail entity.Category

	err = r.NewSelect().Model(&detail).Where("id = ? AND deleted_at IS NULL", id).Scan(ctx)
	if err != nil {
		return entity.Category{}, err
	}

	if detail.Logo != nil {
		baseLink := hashing.GenerateHash(r.ServerBaseUrl, *detail.Logo)
		detail.Logo = &baseLink
	}

	return detail, nil
}

func (r Repository) SuperAdminCreate(ctx context.Context, request category.SuperAdminCreateRequest) (category.SuperAdminCreateResponse, error) {
	claims, err := r.CheckClaims(ctx, auth.RoleSuperAdmin)
	if err != nil {
		return category.SuperAdminCreateResponse{}, err
	}

	err = r.ValidateStruct(&request, "Name", "LogoLink")
	if err != nil {
		return category.SuperAdminCreateResponse{}, err
	}

	response := category.SuperAdminCreateResponse{
		Name:      request.Name,
		Logo:      request.LogoLink,
		Status:    request.Status,
		CreatedAt: time.Now(),
		CreatedBy: claims.UserId,
	}

	_, err = r.NewInsert().Model(&response).Exec(ctx)
	if err != nil {
		return category.SuperAdminCreateResponse{}, web.NewRequestError(errors.Wrap(err, "creating category"), http.StatusBadRequest)
	}

	return response, nil
}

func (r Repository) SuperAdminUpdateAll(ctx context.Context, request category.SuperAdminUpdateRequest) error {
	_, err := r.CheckClaims(ctx, auth.RoleSuperAdmin)
	if err != nil {
		return err
	}

	if err = r.ValidateStruct(&request, "ID", "Name", "LogoLink", "Main"); err != nil {
		return err
	}

	q := r.NewUpdate().Table("categories").Where("deleted_at IS NULL AND id = ?", request.ID)

	q.Set("name = ?", request.Name)
	q.Set("status = ?", request.Status)
	q.Set("logo = ?", request.LogoLink)

	_, err = q.Exec(ctx)
	if err != nil {
		return web.NewRequestError(errors.Wrap(err, "updating categories"), http.StatusBadRequest)
	}

	return nil
}

func (r Repository) SuperAdminUpdateColumns(ctx context.Context, request category.SuperAdminUpdateRequest) error {
	_, err := r.CheckClaims(ctx, auth.RoleSuperAdmin)
	if err != nil {
		return err
	}

	if err = r.ValidateStruct(&request, "ID"); err != nil {
		return err
	}

	q := r.NewUpdate().Table("categories").Where("deleted_at IS NULL AND id = ?", request.ID)

	if request.Name != nil {
		q.Set("name = ?", request.Name)
	}

	if request.LogoLink != nil {
		q.Set("logo = ?", request.LogoLink)
	}

	if request.Status != nil {
		q.Set("status = ?", request.Status)
	}

	_, err = q.Exec(ctx)
	if err != nil {
		return web.NewRequestError(errors.Wrap(err, "updating menu_categories"), http.StatusBadRequest)
	}

	return nil
}

func (r Repository) SuperAdminDelete(ctx context.Context, id int64) error {
	return r.DeleteRow(ctx, "categories", id, auth.RoleSuperAdmin)
}

// @client

func (r Repository) ClientGetList(ctx context.Context, filter category.Filter) ([]category.ClientGetList, int, error) {
	//_, err := r.CheckClaims(ctx, auth.RoleClient)
	//if err != nil {
	//	return nil, 0, err
	//}

	//table := "menu_categories"
	whereQuery := fmt.Sprintf(`WHERE deleted_at IS NULL AND status = true`)
	if filter.Name != nil {
		whereQuery += fmt.Sprintf(" AND name ilike '%s'", "%"+*filter.Name+"%")
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

	whereQuery += fmt.Sprintf(" ORDER BY name %s %s", limitQuery, offsetQuery)

	query := fmt.Sprintf(`
					SELECT 
					    id,
					    name,
					    logo
					FROM 
					    categories
					%s`, whereQuery)

	//query, err := utils.SelectQuery(filter.Fields, filter.Joins, &table, &whereQuery)
	//if err != nil {
	//	return nil, 0, errors.Wrap(err, "select query")
	//}

	list := make([]category.ClientGetList, 0)

	rows, err := r.QueryContext(ctx, query)
	if err != nil {
		return nil, 0, web.NewRequestError(errors.Wrap(err, "select menu_category"), http.StatusInternalServerError)
	}

	err = r.ScanRows(ctx, rows, &list)
	if err != nil {
		return nil, 0, web.NewRequestError(errors.Wrap(err, "scanning menu_category"), http.StatusBadRequest)
	}

	countQuery := fmt.Sprintf(`
		SELECT
			count(id)
		FROM
		    categories
		%s
	`, countWhereQuery)

	countRows, err := r.QueryContext(ctx, countQuery)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, 0, web.NewRequestError(postgres.ErrNotFound, http.StatusNotFound)
	}
	if err != nil {
		return nil, 0, web.NewRequestError(errors.Wrap(err, "selecting menu_category"), http.StatusBadRequest)
	}

	count := 0

	for countRows.Next() {
		if err = countRows.Scan(&count); err != nil {
			return nil, 0, web.NewRequestError(errors.Wrap(err, "scanning menu_category count"), http.StatusBadRequest)
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

// @branch

func (r Repository) BranchGetList(ctx context.Context, filter category.Filter) ([]category.BranchGetList, int, error) {
	_, err := r.CheckClaims(ctx, auth.RoleBranch)
	if err != nil {
		return nil, 0, err
	}

	table := "categories"
	whereQuery := fmt.Sprintf(`WHERE %s.deleted_at IS NULL AND %s.status = true`, table, table)

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

	list := make([]category.BranchGetList, 0)

	rows, err := r.QueryContext(ctx, query)
	if err != nil {
		return nil, 0, web.NewRequestError(errors.Wrap(err, "select category"), http.StatusInternalServerError)
	}

	err = r.ScanRows(ctx, rows, &list)
	if err != nil {
		return nil, 0, web.NewRequestError(errors.Wrap(err, "scanning category"), http.StatusBadRequest)
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
		return nil, 0, web.NewRequestError(errors.Wrap(err, "selecting category"), http.StatusBadRequest)
	}

	count := 0

	for countRows.Next() {
		if err = countRows.Scan(&count); err != nil {
			return nil, 0, web.NewRequestError(errors.Wrap(err, "scanning category count"), http.StatusBadRequest)
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

// @cashier

func (r Repository) CashierGetList(ctx context.Context, filter category.Filter) ([]category.CashierGetList, int, error) {
	_, err := r.CheckClaims(ctx, auth.RoleCashier)
	if err != nil {
		return nil, 0, err
	}

	table := "categories"
	whereQuery := fmt.Sprintf(`WHERE %s.deleted_at IS NULL AND %s.status = true`, table, table)

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

	list := make([]category.CashierGetList, 0)

	rows, err := r.QueryContext(ctx, query)
	if err != nil {
		return nil, 0, web.NewRequestError(errors.Wrap(err, "select category"), http.StatusInternalServerError)
	}

	err = r.ScanRows(ctx, rows, &list)
	if err != nil {
		return nil, 0, web.NewRequestError(errors.Wrap(err, "scanning category"), http.StatusBadRequest)
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
		return nil, 0, web.NewRequestError(errors.Wrap(err, "selecting category"), http.StatusBadRequest)
	}

	count := 0

	for countRows.Next() {
		if err = countRows.Scan(&count); err != nil {
			return nil, 0, web.NewRequestError(errors.Wrap(err, "scanning category count"), http.StatusBadRequest)
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

// @admin

func (r Repository) AdminGetList(ctx context.Context, filter category.Filter) ([]category.AdminGetList, int, error) {
	_, err := r.CheckClaims(ctx, auth.RoleAdmin)
	if err != nil {
		return nil, 0, err
	}

	table := "categories"
	whereQuery := fmt.Sprintf(`WHERE %s.deleted_at IS NULL AND %s.status = true`, table, table)

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

	list := make([]category.AdminGetList, 0)

	rows, err := r.QueryContext(ctx, query)
	if err != nil {
		return nil, 0, web.NewRequestError(errors.Wrap(err, "select menu_category"), http.StatusInternalServerError)
	}

	err = r.ScanRows(ctx, rows, &list)
	if err != nil {
		return nil, 0, web.NewRequestError(errors.Wrap(err, "scanning menu_category"), http.StatusBadRequest)
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
		return nil, 0, web.NewRequestError(errors.Wrap(err, "selecting menu_category"), http.StatusBadRequest)
	}

	count := 0

	for countRows.Next() {
		if err = countRows.Scan(&count); err != nil {
			return nil, 0, web.NewRequestError(errors.Wrap(err, "scanning menu_category count"), http.StatusBadRequest)
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

func (r Repository) WaiterGetList(ctx context.Context) ([]category.WaiterGetList, error) {
	_, err := r.CheckClaims(ctx, auth.RoleWaiter)
	if err != nil {
		return nil, err
	}

	query := fmt.Sprintf(`SELECT 
								fc.id, 
								fc.name,
								fc.logo
						  FROM categories fc
						  WHERE fc.deleted_at isnull AND fc.status = true`)

	rows, err := r.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}

	var response []category.WaiterGetList
	if err = r.ScanRows(ctx, rows, &response); err != nil {
		return nil, err
	}

	for k, v := range response {
		if v.Logo != nil {
			baseLink := hashing.GenerateHash(r.ServerBaseUrl, *v.Logo)
			response[k].Logo = &baseLink
		}
	}

	return response, nil
}

func NewRepository(DB *postgresql.Database) *Repository {
	return &Repository{DB}
}
