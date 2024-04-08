package foodCategory

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
	"restu-backend/internal/service/foodCategory"
	"restu-backend/internal/service/hashing"
	"time"
)

type Repository struct {
	*postgresql.Database
}

// @admin

func (r Repository) AdminGetList(ctx context.Context, filter foodCategory.Filter) ([]foodCategory.AdminGetList, int, error) {
	claims, err := r.CheckClaims(ctx, auth.RoleAdmin)
	if err != nil {
		return nil, 0, err
	}

	table := "food_category"
	whereQuery := fmt.Sprintf(`WHERE %s.deleted_at IS NULL AND %s.restaurant_id = %d`, table, table, *claims.RestaurantID)
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
			baseLink := hashing.GenerateHash(r.ServerBaseUrl, *v.Logo)
			list[k].Logo = &baseLink
		}
	}

	return list, count, nil
}

func (r Repository) AdminGetDetail(ctx context.Context, id int64) (entity.FoodCategory, error) {
	claims, err := r.CheckClaims(ctx, auth.RoleAdmin)
	if err != nil {
		return entity.FoodCategory{}, err
	}

	var detail entity.FoodCategory

	err = r.NewSelect().Model(&detail).Where("id = ? AND deleted_at IS NULL AND restaurant_id = ?", id, *claims.RestaurantID).Scan(ctx)
	if err != nil {
		return entity.FoodCategory{}, err
	}

	if detail.Logo != nil {
		baseLink := hashing.GenerateHash(r.ServerBaseUrl, *detail.Logo)
		detail.Logo = &baseLink
	}

	return detail, nil
}

func (r Repository) AdminCreate(ctx context.Context, request foodCategory.AdminCreateRequest) (foodCategory.AdminCreateResponse, error) {
	claims, err := r.CheckClaims(ctx, auth.RoleAdmin)
	if err != nil {
		return foodCategory.AdminCreateResponse{}, err
	}

	err = r.ValidateStruct(&request, "Name", "LogoLink")
	if err != nil {
		return foodCategory.AdminCreateResponse{}, err
	}

	response := foodCategory.AdminCreateResponse{
		Name:         request.Name,
		Logo:         request.LogoLink,
		Main:         request.Main,
		RestaurantID: claims.RestaurantID,
		CreatedAt:    time.Now(),
		CreatedBy:    claims.UserId,
	}

	_, err = r.NewInsert().Model(&response).Exec(ctx)
	if err != nil {
		return foodCategory.AdminCreateResponse{}, err
	}

	return response, nil
}

func (r Repository) AdminUpdateAll(ctx context.Context, request foodCategory.AdminUpdateRequest) error {
	claims, err := r.CheckClaims(ctx, auth.RoleAdmin)
	if err != nil {
		return err
	}

	if err = r.ValidateStruct(&request, "ID", "Name", "LogoLink", "Main"); err != nil {
		return err
	}

	q := r.NewUpdate().Table("food_category").Where("deleted_at IS NULL AND id = ? AND restaurant_id = ?", request.ID, *claims.RestaurantID)

	q.Set("name = ?", request.Name)
	q.Set("main = ?", request.Main)
	q.Set("logo = ?", request.LogoLink)

	_, err = q.Exec(ctx)
	if err != nil {
		return web.NewRequestError(errors.Wrap(err, "updating food_category"), http.StatusBadRequest)
	}

	return nil
}

func (r Repository) AdminUpdateColumns(ctx context.Context, request foodCategory.AdminUpdateRequest) error {
	claims, err := r.CheckClaims(ctx, auth.RoleAdmin)
	if err != nil {
		return err
	}

	if err = r.ValidateStruct(&request, "ID"); err != nil {
		return err
	}

	q := r.NewUpdate().Table("food_category").Where("deleted_at IS NULL AND id = ? AND restaurant_id = ?", request.ID, *claims.RestaurantID)

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

func (r Repository) AdminDelete(ctx context.Context, id int64) error {
	claims, err := r.CheckClaims(ctx, auth.RoleAdmin)
	if err != nil {
		return err
	}

	_, err = r.NewUpdate().Table("food_category").
		Set("deleted_at = ?", time.Now()).
		Where("restaurant_id = ? and id = ?", *claims.RestaurantID, id).Exec(ctx)
	if err != nil {
		return err
	}

	return nil
}

// @branch

func (r Repository) BranchGetList(ctx context.Context, filter foodCategory.Filter) ([]foodCategory.BranchGetList, int, error) {
	claims, err := r.CheckClaims(ctx, auth.RoleBranch)
	if err != nil {
		return nil, 0, err
	}

	table := "food_category"
	whereQuery := fmt.Sprintf(`WHERE %s.deleted_at IS NULL AND %s.restaurant_id = (select restaurant_id from branches where id = %d)`, table, table, *claims.BranchID)
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
			baseLink := hashing.GenerateHash(r.ServerBaseUrl, *v.Logo)
			list[k].Logo = &baseLink
		}
	}

	return list, count, nil
}

func (r Repository) BranchGetDetail(ctx context.Context, id int64) (entity.FoodCategory, error) {
	claims, err := r.CheckClaims(ctx, auth.RoleBranch)
	if err != nil {
		return entity.FoodCategory{}, err
	}

	var detail entity.FoodCategory

	err = r.NewSelect().Model(&detail).Where("id = ? AND deleted_at IS NULL AND restaurant_id = (select restaurant_id from branches where id = ?)", id, *claims.BranchID).Scan(ctx)
	if err != nil {
		return entity.FoodCategory{}, err
	}

	if detail.Logo != nil {
		baseLink := hashing.GenerateHash(r.ServerBaseUrl, *detail.Logo)
		detail.Logo = &baseLink
	}

	return detail, nil
}

func (r Repository) BranchCreate(ctx context.Context, request foodCategory.BranchCreateRequest) (foodCategory.BranchCreateResponse, error) {
	claims, err := r.CheckClaims(ctx, auth.RoleBranch)
	if err != nil {
		return foodCategory.BranchCreateResponse{}, err
	}

	err = r.ValidateStruct(&request, "Name", "LogoLink")
	if err != nil {
		return foodCategory.BranchCreateResponse{}, err
	}

	var restaurantID *int64
	query := fmt.Sprintf(`SELECT restaurant_id FROM branches WHERE id='%d'`, *claims.BranchID)
	if err = r.QueryRowContext(ctx, query).Scan(&restaurantID); err != nil || restaurantID == nil {
		return foodCategory.BranchCreateResponse{}, err
	}

	response := foodCategory.BranchCreateResponse{
		Name:         request.Name,
		Logo:         request.LogoLink,
		Main:         request.Main,
		RestaurantID: restaurantID,
		CreatedAt:    time.Now(),
		CreatedBy:    claims.UserId,
	}

	_, err = r.NewInsert().Model(&response).Exec(ctx)
	if err != nil {
		return foodCategory.BranchCreateResponse{}, err
	}

	return response, nil
}

func (r Repository) BranchUpdateAll(ctx context.Context, request foodCategory.BranchUpdateRequest) error {
	claims, err := r.CheckClaims(ctx, auth.RoleBranch)
	if err != nil {
		return err
	}

	if err = r.ValidateStruct(&request, "ID", "Name", "LogoLink", "Main"); err != nil {
		return err
	}

	q := r.NewUpdate().Table("food_category").Where("deleted_at IS NULL AND id = ? AND restaurant_id = (select restaurant_id from branches where id = ?)", request.ID, *claims.BranchID)

	q.Set("name = ?", request.Name)
	q.Set("main = ?", request.Main)
	q.Set("logo = ?", request.LogoLink)

	_, err = q.Exec(ctx)
	if err != nil {
		return web.NewRequestError(errors.Wrap(err, "updating food_category"), http.StatusBadRequest)
	}

	return nil
}

func (r Repository) BranchUpdateColumns(ctx context.Context, request foodCategory.BranchUpdateRequest) error {
	claims, err := r.CheckClaims(ctx, auth.RoleBranch)
	if err != nil {
		return err
	}

	if err = r.ValidateStruct(&request, "ID"); err != nil {
		return err
	}

	q := r.NewUpdate().Table("food_category").Where("deleted_at IS NULL AND id = ? AND restaurant_id = (select restaurant_id from branches where id = ?)", request.ID, *claims.BranchID)

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

func (r Repository) BranchDelete(ctx context.Context, id int64) error {
	claims, err := r.CheckClaims(ctx, auth.RoleBranch)
	if err != nil {
		return err
	}

	_, err = r.NewUpdate().Table("food_category").
		Set("deleted_at = ?", time.Now()).
		Where("restaurant_id = (select restaurant_id from branches where id = ?) and id = ?", *claims.BranchID, id).Exec(ctx)
	if err != nil {
		return err
	}

	return nil
}

// @cashier

func (r Repository) CashierGetList(ctx context.Context, filter foodCategory.Filter) ([]foodCategory.CashierGetList, int, error) {
	claims, err := r.CheckClaims(ctx, auth.RoleCashier)
	if err != nil {
		return nil, 0, err
	}

	table := "food_category"
	whereQuery := fmt.Sprintf(`WHERE %s.deleted_at IS NULL AND %s.restaurant_id = (select restaurant_id from branches where id = %d)`, table, table, *claims.BranchID)
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

	list := make([]foodCategory.CashierGetList, 0)

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

func (r Repository) CashierGetDetail(ctx context.Context, id int64) (entity.FoodCategory, error) {
	claims, err := r.CheckClaims(ctx, auth.RoleCashier)
	if err != nil {
		return entity.FoodCategory{}, err
	}

	var detail entity.FoodCategory

	err = r.NewSelect().Model(&detail).Where("id = ? AND deleted_at IS NULL AND restaurant_id = (select restaurant_id from branches where id = ?)", id, *claims.BranchID).Scan(ctx)
	if err != nil {
		return entity.FoodCategory{}, err
	}

	if detail.Logo != nil {
		baseLink := hashing.GenerateHash(r.ServerBaseUrl, *detail.Logo)
		detail.Logo = &baseLink
	}

	return detail, nil
}

func (r Repository) CashierCreate(ctx context.Context, request foodCategory.CashierCreateRequest) (foodCategory.CashierCreateResponse, error) {
	claims, err := r.CheckClaims(ctx, auth.RoleCashier)
	if err != nil {
		return foodCategory.CashierCreateResponse{}, err
	}

	err = r.ValidateStruct(&request, "Name", "LogoLink")
	if err != nil {
		return foodCategory.CashierCreateResponse{}, err
	}

	var restaurantID *int64
	query := fmt.Sprintf(`SELECT restaurant_id FROM branches WHERE id='%d'`, *claims.BranchID)
	if err = r.QueryRowContext(ctx, query).Scan(&restaurantID); err != nil || restaurantID == nil {
		return foodCategory.CashierCreateResponse{}, err
	}

	response := foodCategory.CashierCreateResponse{
		Name:         request.Name,
		Logo:         request.LogoLink,
		Main:         request.Main,
		RestaurantID: restaurantID,
		CreatedAt:    time.Now(),
		CreatedBy:    claims.UserId,
	}

	_, err = r.NewInsert().Model(&response).Exec(ctx)
	if err != nil {
		return foodCategory.CashierCreateResponse{}, err
	}

	return response, nil
}

func (r Repository) CashierUpdateAll(ctx context.Context, request foodCategory.CashierUpdateRequest) error {
	claims, err := r.CheckClaims(ctx, auth.RoleCashier)
	if err != nil {
		return err
	}

	if err = r.ValidateStruct(&request, "ID", "Name", "LogoLink", "Main"); err != nil {
		return err
	}

	q := r.NewUpdate().Table("food_category").Where("deleted_at IS NULL AND id = ? AND restaurant_id = (select restaurant_id from branches where id = ?)", request.ID, *claims.BranchID)

	q.Set("name = ?", request.Name)
	q.Set("main = ?", request.Main)
	q.Set("logo = ?", request.LogoLink)

	_, err = q.Exec(ctx)
	if err != nil {
		return web.NewRequestError(errors.Wrap(err, "updating food_category"), http.StatusBadRequest)
	}

	return nil
}

func (r Repository) CashierUpdateColumns(ctx context.Context, request foodCategory.CashierUpdateRequest) error {
	claims, err := r.CheckClaims(ctx, auth.RoleCashier)
	if err != nil {
		return err
	}

	if err = r.ValidateStruct(&request, "ID"); err != nil {
		return err
	}

	q := r.NewUpdate().Table("food_category").Where("deleted_at IS NULL AND id = ? AND restaurant_id = (select restaurant_id from branches where id = ?)", request.ID, *claims.BranchID)

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

func (r Repository) CashierDelete(ctx context.Context, id int64) error {
	claims, err := r.CheckClaims(ctx, auth.RoleCashier)
	if err != nil {
		return err
	}

	_, err = r.NewUpdate().Table("food_category").
		Set("deleted_at = ?", time.Now()).
		Where("restaurant_id = (select restaurant_id from branches where id = ?) and id = ?", *claims.BranchID, id).Exec(ctx)
	if err != nil {
		return err
	}

	return nil
}

// @client

func (r Repository) ClientGetList(ctx context.Context, filter foodCategory.Filter) ([]foodCategory.ClientGetList, int, error) {
	//_, err := r.CheckClaims(ctx, auth.RoleClient)
	//if err != nil {
	//	return nil, 0, err
	//}

	//table := "food_category"
	whereQuery := fmt.Sprintf(`WHERE fc.deleted_at IS NULL AND fc.main = true`)
	if filter.Name != nil {
		whereQuery += fmt.Sprintf(" AND fc.name ilike '%s'", "%"+*filter.Name+"%")
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

	whereQuery += fmt.Sprintf(" GROUP BY fc.id, fc.name ORDER BY fc.name %s %s", limitQuery, offsetQuery)

	query := fmt.Sprintf(`
					SELECT 
					    fc.id,
					    fc.name,
					    fc.logo
					FROM 
					    food_category as fc
					LEFT OUTER JOIN foods as f ON f.category_id = fc.id
					%s`, whereQuery)

	//query, err := utils.SelectQuery(filter.Fields, filter.Joins, &table, &whereQuery)
	//if err != nil {
	//	return nil, 0, errors.Wrap(err, "select query")
	//}

	list := make([]foodCategory.ClientGetList, 0)

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
			count(fc.id)
		FROM
		    food_category as fc
		%s
	`, countWhereQuery)

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

func (r Repository) WaiterGetList(ctx context.Context) ([]foodCategory.WaiterGetList, error) {
	claims, err := r.CheckClaims(ctx, auth.RoleWaiter)
	if err != nil {
		return nil, err
	}

	query := fmt.Sprintf(`SELECT 
								fc.id, 
								fc.name 
						  FROM menu_categories fc
						  WHERE fc.deleted_at isnull AND (select count(m.id) from menus m where m.menu_category_id=fc.id and m.branch_id='%d') != 0`, *claims.BranchID)

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
