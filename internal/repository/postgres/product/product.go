package product

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
	"restu-backend/internal/repository/postgres"
	"restu-backend/internal/service/product"
	"time"
)

type Repository struct {
	*postgresql.Database
}

// @admin

func (r Repository) AdminGetList(ctx context.Context, filter product.Filter) ([]product.AdminGetList, int, error) {
	claims, err := r.CheckClaims(ctx, auth.RoleAdmin)
	if err != nil {
		return nil, 0, err
	}

	whereQuery := fmt.Sprintf(`WHERE p.deleted_at IS NULL AND p.restaurant_id = %d`, *claims.RestaurantID)
	countWhereQuery := whereQuery

	if filter.Name != nil {
		whereQuery += fmt.Sprintf(" AND p.name ilike '%s'", "%"+*filter.Name+"%")
	}

	query := fmt.Sprintf(`
					SELECT 
					    p.id,
					    p.name,
					    p.measure_unit_id,
					    m.name as measure_unit,
					    p.barcode
					FROM 
					    products as p
					LEFT OUTER JOIN measure_unit as m ON m.id = p.measure_unit_id
					%s
	`, whereQuery)

	list := make([]product.AdminGetList, 0)

	rows, err := r.QueryContext(ctx, query)
	if err != nil {
		return nil, 0, web.NewRequestError(errors.Wrap(err, "select product"), http.StatusInternalServerError)
	}

	err = r.ScanRows(ctx, rows, &list)
	if err != nil {
		return nil, 0, web.NewRequestError(errors.Wrap(err, "scanning products"), http.StatusBadRequest)
	}

	countQuery := fmt.Sprintf(`
		SELECT
			count(p.id)
		FROM
		    products as p
		LEFT OUTER JOIN measure_unit as m ON m.id = p.measure_unit_id
		%s
	`, countWhereQuery)

	countRows, err := r.QueryContext(ctx, countQuery)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, 0, web.NewRequestError(postgres.ErrNotFound, http.StatusNotFound)
	}
	if err != nil {
		return nil, 0, web.NewRequestError(errors.Wrap(err, "selecting products"), http.StatusBadRequest)
	}

	count := 0

	for countRows.Next() {
		if err = countRows.Scan(&count); err != nil {
			return nil, 0, web.NewRequestError(errors.Wrap(err, "scanning products count"), http.StatusBadRequest)
		}
	}

	return list, count, nil
}

func (r Repository) AdminGetDetail(ctx context.Context, id int64) (entity.Product, error) {
	claims, err := r.CheckClaims(ctx, auth.RoleAdmin)
	if err != nil {
		return entity.Product{}, err
	}

	var detail entity.Product

	err = r.NewSelect().Model(&detail).Where("id = ? AND deleted_at IS NULL AND restaurant_id = ?", id, claims.RestaurantID).Scan(ctx)
	if err != nil {
		return entity.Product{}, err
	}

	return detail, nil
}

func (r Repository) AdminCreate(ctx context.Context, request product.AdminCreateRequest) (product.AdminCreateResponse, error) {
	claims, err := r.CheckClaims(ctx, auth.RoleAdmin)
	if err != nil {
		return product.AdminCreateResponse{}, err
	}

	err = r.ValidateStruct(&request, "Name", "MeasureUnitID")
	if err != nil {
		return product.AdminCreateResponse{}, err
	}

	response := product.AdminCreateResponse{
		Name:          request.Name,
		MeasureUnitID: request.MeasureUnitID,
		Barcode:       request.Barcode,
		CreatedAt:     time.Now(),
		CreatedBy:     claims.UserId,
		RestaurantID:  *claims.RestaurantID,
	}

	_, err = r.NewInsert().Model(&response).Exec(ctx)
	if err != nil {
		return product.AdminCreateResponse{}, web.NewRequestError(errors.Wrap(err, "creating product"), http.StatusBadRequest)
	}

	return response, nil
}

func (r Repository) AdminUpdateAll(ctx context.Context, request product.AdminUpdateRequest) error {
	claims, err := r.CheckClaims(ctx, auth.RoleAdmin)
	if err != nil {
		return err
	}

	if err = r.ValidateStruct(&request, "ID", "Name", "MeasureUnitID"); err != nil {
		return err
	}

	q := r.NewUpdate().Table("products").Where("deleted_at IS NULL "+
		"AND id = ? "+
		"AND restaurant_id = ?", request.ID, claims.RestaurantID)

	q.Set("name = ?", request.Name)
	q.Set("measure_unit_id =?", request.MeasureUnitID)
	q.Set("barcode = ?", request.Barcode)
	q.Set("updated_at = ?", time.Now())
	q.Set("updated_by = ?", claims.UserId)

	_, err = q.Exec(ctx)
	if err != nil {
		return web.NewRequestError(errors.Wrap(err, "updating product"), http.StatusBadRequest)
	}

	return nil
}

func (r Repository) AdminUpdateColumns(ctx context.Context, request product.AdminUpdateRequest) error {
	claims, err := r.CheckClaims(ctx, auth.RoleAdmin)
	if err != nil {
		return err
	}

	if err = r.ValidateStruct(&request, "ID"); err != nil {
		return err
	}

	q := r.NewUpdate().Table("products").Where("deleted_at IS NULL "+
		"AND id = ? "+
		"AND restaurant_id = ?", request.ID, claims.RestaurantID)

	if request.Name != nil {
		q.Set("name = ?", request.Name)
	}

	if request.MeasureUnitID != nil {
		q.Set("measure_unit_id = ?", request.MeasureUnitID)
	}

	if request.Barcode != nil {
		q.Set("barcode = ?", request.Barcode)
	}

	q.Set("updated_at = ?", time.Now())
	q.Set("updated_by = ?", claims.UserId)

	_, err = q.Exec(ctx)
	if err != nil {
		return web.NewRequestError(errors.Wrap(err, "updating product"), http.StatusBadRequest)
	}

	return nil
}

func (r Repository) AdminDelete(ctx context.Context, id int64) error {
	return r.DeleteRow(ctx, "products", id, auth.RoleAdmin)
}

func (r Repository) AdminGetSpendingByBranch(ctx context.Context, filter product.SpendingFilter) ([]product.AdminGetSpendingByBranchResponse, error) {
	claims, err := r.CheckClaims(ctx, auth.RoleAdmin)
	if err != nil {
		return nil, err
	}

	var (
		response  []product.AdminGetSpendingByBranchResponse
		whereDate string
	)

	if filter.FromDate != nil {
		from := filter.FromDate.Format("02.01.2006")
		whereDate += fmt.Sprintf(` and to_char(or.created_at, 'DD.MM.YYYY') >= '%s' and case when frgh.from isnull then true else frgh.from <= '%s' end`, from, from)
	}
	if filter.ToDate != nil {
		to := filter.ToDate.Format("02.01.2006")
		whereDate += fmt.Sprintf(` and to_char(or.created_at, 'DD.MM.YYYY') <= '%s' and case when frgh.to isnull then true else frgh.to >= '%s' end`, to, to)
	}

	if filter.BranchId != nil {
		// building product query for every branch of restaurant...
		productQuery := fmt.Sprintf(`select
												p.id as id,
												p.name as name,
												sum(fr.amount*om.count) as amount,
												mu.name as measure_unit
											from foods f 
											    	 join food_recipe_group_histories frgh 
											    	     on f.id = frgh.food_id 
											    	 join food_recipe_groups frg 
											    	     on f.id = frg.food_id
													 join menus m
														  on f.id = any(m.food_ids)
													 join order_menu om
														  on m.id = om.menu_id 
											         join food_recipe fr 
											              on fr.id = any (frg.recipe_ids)
													 join products p
														  on fr.product_id = p.id
													 join measure_unit mu
														  on p.measure_unit_id = mu.id
													 join orders o
														 on om.order_id = o.id
													 join tables t
														 on o.table_id = t.id
													 join restaurants r 
													     on p.restaurant_id = r.id
											where t.branch_id='%d' and p.deleted_at isnull and om.deleted_at isnull and o.status = 'PAID' and r.id='%d' and om.status = 'PAID' %s
											group by p.id, mu.name`, *filter.BranchId, *claims.RestaurantID, whereDate)
		// scanning products spending [heart of the api]...
		rows, err := r.QueryContext(ctx, productQuery)
		if err != nil {
			return nil, err
		}
		if err = r.ScanRows(ctx, rows, &response); err != nil {
			return nil, err
		}
	} else {
		err = errors.New("branch_id not specified")
		return nil, web.NewRequestError(err, http.StatusBadRequest)
	}

	return response, nil
}

// @branch

func (r Repository) BranchGetList(ctx context.Context, filter product.Filter) ([]product.BranchGetList, int, error) {
	claims, err := r.CheckClaims(ctx, auth.RoleBranch)
	if err != nil {
		return nil, 0, err
	}

	whereQuery := fmt.Sprintf(`
WHERE p.deleted_at IS NULL AND p.restaurant_id in (select restaurant_id from branches where id = '%d')`,
		*claims.BranchID)
	countWhereQuery := whereQuery

	if filter.Name != nil {
		whereQuery += fmt.Sprintf(" AND p.name ilike '%s'", "%"+*filter.Name+"%")
	}

	query := fmt.Sprintf(`
					SELECT 
					    p.id,
					    p.name,
					    p.measure_unit_id,
					    m.name as measure_unit,
					    p.barcode
					FROM 
					    products as p
					LEFT OUTER JOIN measure_unit as m ON m.id = p.measure_unit_id
					%s
	`, whereQuery)

	list := make([]product.BranchGetList, 0)

	rows, err := r.QueryContext(ctx, query)
	if err != nil {
		return nil, 0, web.NewRequestError(errors.Wrap(err, "select product"), http.StatusInternalServerError)
	}

	err = r.ScanRows(ctx, rows, &list)
	if err != nil {
		return nil, 0, web.NewRequestError(errors.Wrap(err, "scanning products"), http.StatusBadRequest)
	}

	countQuery := fmt.Sprintf(`
		SELECT
			count(p.id)
		FROM
		    products as p
		LEFT OUTER JOIN measure_unit as m ON m.id = p.measure_unit_id
		%s
	`, countWhereQuery)

	countRows, err := r.QueryContext(ctx, countQuery)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, 0, web.NewRequestError(postgres.ErrNotFound, http.StatusNotFound)
	}
	if err != nil {
		return nil, 0, web.NewRequestError(errors.Wrap(err, "selecting products"), http.StatusBadRequest)
	}

	count := 0

	for countRows.Next() {
		if err = countRows.Scan(&count); err != nil {
			return nil, 0, web.NewRequestError(errors.Wrap(err, "scanning products count"), http.StatusBadRequest)
		}
	}

	return list, count, nil
}

func (r Repository) BranchGetDetail(ctx context.Context, id int64) (entity.Product, error) {
	claims, err := r.CheckClaims(ctx, auth.RoleBranch)
	if err != nil {
		return entity.Product{}, err
	}

	var restaurantID *int64

	err = r.QueryRowContext(ctx, fmt.Sprintf("SELECT restaurant_id FROM branches WHERE id = '%d'", *claims.BranchID)).Scan(&restaurantID)
	if err != nil {
		return entity.Product{}, web.NewRequestError(errors.Wrap(err, "restaurant not found in "), http.StatusBadRequest)
	}

	if restaurantID == nil {
		return entity.Product{}, web.NewRequestError(errors.New("not found restaurant"), http.StatusBadRequest)
	}

	var detail entity.Product

	err = r.NewSelect().Model(&detail).Where("id = ? AND deleted_at IS NULL AND restaurant_id = ?", id, restaurantID).Scan(ctx)
	if errors.Is(err, sql.ErrNoRows) {
		return entity.Product{}, web.NewRequestError(postgres.ErrNotFound, http.StatusBadRequest)
	}
	if err != nil {
		return entity.Product{}, web.NewRequestError(errors.Wrap(err, "selecting products"), http.StatusBadRequest)
	}

	return detail, nil
}

func (r Repository) BranchCreate(ctx context.Context, request product.BranchCreateRequest) (product.BranchCreateResponse, error) {
	claims, err := r.CheckClaims(ctx, auth.RoleBranch)
	if err != nil {
		return product.BranchCreateResponse{}, err
	}

	err = r.ValidateStruct(&request, "Name", "MeasureUnitID")
	if err != nil {
		return product.BranchCreateResponse{}, err
	}

	var restaurantID *int64

	err = r.QueryRowContext(ctx, fmt.Sprintf("SELECT restaurant_id FROM branches WHERE id = '%d'", *claims.BranchID)).Scan(&restaurantID)
	if err != nil {
		return product.BranchCreateResponse{}, web.NewRequestError(errors.Wrap(err, "restaurant not found in "), http.StatusBadRequest)
	}

	if restaurantID == nil {
		return product.BranchCreateResponse{}, web.NewRequestError(errors.New("not found restaurant"), http.StatusBadRequest)
	}

	response := product.BranchCreateResponse{
		Name:          request.Name,
		MeasureUnitID: request.MeasureUnitID,
		Barcode:       request.Barcode,
		CreatedAt:     time.Now(),
		CreatedBy:     claims.UserId,
		RestaurantID:  *restaurantID,
	}

	_, err = r.NewInsert().Model(&response).Exec(ctx)
	if err != nil {
		return product.BranchCreateResponse{}, web.NewRequestError(errors.Wrap(err, "creating product"), http.StatusBadRequest)
	}

	return response, nil
}

func (r Repository) BranchUpdateAll(ctx context.Context, request product.BranchUpdateRequest) error {
	claims, err := r.CheckClaims(ctx, auth.RoleBranch)
	if err != nil {
		return err
	}

	if err = r.ValidateStruct(&request, "ID", "Name", "MeasureUnitID"); err != nil {
		return err
	}

	var restaurantID *int64

	err = r.QueryRowContext(ctx, fmt.Sprintf("SELECT restaurant_id FROM branches WHERE id = '%d'", *claims.BranchID)).Scan(&restaurantID)
	if err != nil {
		return web.NewRequestError(errors.Wrap(err, "restaurant not found in "), http.StatusBadRequest)
	}

	if restaurantID == nil {
		return web.NewRequestError(errors.New("not found restaurant"), http.StatusBadRequest)
	}

	q := r.NewUpdate().Table("products").Where("deleted_at IS NULL "+
		"AND id = ? "+
		"AND restaurant_id = ?", request.ID, restaurantID)

	q.Set("name = ?", request.Name)
	q.Set("measure_unit_id =?", request.MeasureUnitID)
	q.Set("barcode = ?", request.Barcode)
	q.Set("updated_at = ?", time.Now())
	q.Set("updated_by = ?", claims.UserId)

	_, err = q.Exec(ctx)
	if err != nil {
		return web.NewRequestError(errors.Wrap(err, "updating product"), http.StatusBadRequest)
	}

	return nil
}

func (r Repository) BranchUpdateColumns(ctx context.Context, request product.BranchUpdateRequest) error {
	claims, err := r.CheckClaims(ctx, auth.RoleBranch)
	if err != nil {
		return err
	}

	if err = r.ValidateStruct(&request, "ID"); err != nil {
		return err
	}

	var restaurantID *int64

	err = r.QueryRowContext(ctx, fmt.Sprintf("SELECT restaurant_id FROM branches WHERE id = '%d'", *claims.BranchID)).Scan(&restaurantID)
	if err != nil {
		return web.NewRequestError(errors.Wrap(err, "restaurant not found in "), http.StatusBadRequest)
	}

	if restaurantID == nil {
		return web.NewRequestError(errors.New("not found restaurant"), http.StatusBadRequest)
	}

	q := r.NewUpdate().Table("products").Where("deleted_at IS NULL "+
		"AND id = ? "+
		"AND restaurant_id = ?", request.ID, restaurantID)

	if request.Name != nil {
		q.Set("name = ?", request.Name)
	}

	if request.MeasureUnitID != nil {
		q.Set("measure_unit_id = ?", request.MeasureUnitID)
	}

	if request.Barcode != nil {
		q.Set("barcode = ?", request.Barcode)
	}

	q.Set("updated_at = ?", time.Now())
	q.Set("updated_by = ?", claims.UserId)

	_, err = q.Exec(ctx)
	if err != nil {
		return web.NewRequestError(errors.Wrap(err, "updating product"), http.StatusBadRequest)
	}

	return nil
}

func (r Repository) BranchDelete(ctx context.Context, id int64) error {
	return r.DeleteRow(ctx, "products", id, auth.RoleBranch)
}

// @cashier

func (r Repository) CashierGetList(ctx context.Context, filter product.Filter) ([]product.CashierGetList, int, error) {
	claims, err := r.CheckClaims(ctx, auth.RoleCashier)
	if err != nil {
		return nil, 0, err
	}

	whereQuery := fmt.Sprintf(`
WHERE p.deleted_at IS NULL AND p.restaurant_id in (select restaurant_id from branches where id = '%d')`,
		*claims.BranchID)
	countWhereQuery := whereQuery

	if filter.Name != nil {
		whereQuery += fmt.Sprintf(" AND p.name ilike '%s'", "%"+*filter.Name+"%")
	}

	query := fmt.Sprintf(`
					SELECT 
					    p.id,
					    p.name,
					    p.measure_unit_id,
					    m.name as measure_unit
					FROM 
					    products as p
					LEFT OUTER JOIN measure_unit as m ON m.id = p.measure_unit_id
					%s
	`, whereQuery)

	list := make([]product.CashierGetList, 0)

	rows, err := r.QueryContext(ctx, query)
	if err != nil {
		return nil, 0, web.NewRequestError(errors.Wrap(err, "select product"), http.StatusInternalServerError)
	}

	err = r.ScanRows(ctx, rows, &list)
	if err != nil {
		return nil, 0, web.NewRequestError(errors.Wrap(err, "scanning products"), http.StatusBadRequest)
	}

	countQuery := fmt.Sprintf(`
		SELECT
			count(p.id)
		FROM
		    products as p
		LEFT OUTER JOIN measure_unit as m ON m.id = p.measure_unit_id
		%s
	`, countWhereQuery)

	countRows, err := r.QueryContext(ctx, countQuery)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, 0, web.NewRequestError(postgres.ErrNotFound, http.StatusNotFound)
	}
	if err != nil {
		return nil, 0, web.NewRequestError(errors.Wrap(err, "selecting products"), http.StatusBadRequest)
	}

	count := 0

	for countRows.Next() {
		if err = countRows.Scan(&count); err != nil {
			return nil, 0, web.NewRequestError(errors.Wrap(err, "scanning products count"), http.StatusBadRequest)
		}
	}

	return list, count, nil
}

func (r Repository) CashierGetDetail(ctx context.Context, id int64) (entity.Product, error) {
	claims, err := r.CheckClaims(ctx, auth.RoleCashier)
	if err != nil {
		return entity.Product{}, err
	}

	var restaurantID *int64

	err = r.QueryRowContext(ctx, fmt.Sprintf("SELECT restaurant_id FROM branches WHERE id = '%d'", *claims.BranchID)).Scan(&restaurantID)
	if err != nil {
		return entity.Product{}, web.NewRequestError(errors.Wrap(err, "restaurant not found in "), http.StatusBadRequest)
	}

	if restaurantID == nil {
		return entity.Product{}, web.NewRequestError(errors.New("not found restaurant"), http.StatusBadRequest)
	}

	var detail entity.Product

	err = r.NewSelect().Model(&detail).Where("id = ? AND deleted_at IS NULL AND restaurant_id = ?", id, restaurantID).Scan(ctx)
	if errors.Is(err, sql.ErrNoRows) {
		return entity.Product{}, web.NewRequestError(postgres.ErrNotFound, http.StatusBadRequest)
	}
	if err != nil {
		return entity.Product{}, web.NewRequestError(errors.Wrap(err, "selecting products"), http.StatusBadRequest)
	}

	return detail, nil
}

func (r Repository) CashierCreate(ctx context.Context, request product.CashierCreateRequest) (product.CashierCreateResponse, error) {
	claims, err := r.CheckClaims(ctx, auth.RoleCashier)
	if err != nil {
		return product.CashierCreateResponse{}, err
	}

	err = r.ValidateStruct(&request, "Name", "MeasureUnitID")
	if err != nil {
		return product.CashierCreateResponse{}, err
	}

	var restaurantID *int64

	err = r.QueryRowContext(ctx, fmt.Sprintf("SELECT restaurant_id FROM branches WHERE id = '%d'", *claims.BranchID)).Scan(&restaurantID)
	if err != nil {
		return product.CashierCreateResponse{}, web.NewRequestError(errors.Wrap(err, "restaurant not found in "), http.StatusBadRequest)
	}

	if restaurantID == nil {
		return product.CashierCreateResponse{}, web.NewRequestError(errors.New("not found restaurant"), http.StatusBadRequest)
	}

	response := product.CashierCreateResponse{
		Name:          request.Name,
		MeasureUnitID: request.MeasureUnitID,
		CreatedAt:     time.Now(),
		CreatedBy:     claims.UserId,
		RestaurantID:  *restaurantID,
	}

	_, err = r.NewInsert().Model(&response).Exec(ctx)
	if err != nil {
		return product.CashierCreateResponse{}, web.NewRequestError(errors.Wrap(err, "creating product"), http.StatusBadRequest)
	}

	return response, nil
}

func (r Repository) CashierUpdateAll(ctx context.Context, request product.CashierUpdateRequest) error {
	claims, err := r.CheckClaims(ctx, auth.RoleCashier)
	if err != nil {
		return err
	}

	if err = r.ValidateStruct(&request, "ID", "Name", "MeasureUnitID"); err != nil {
		return err
	}

	var restaurantID *int64

	err = r.QueryRowContext(ctx, fmt.Sprintf("SELECT restaurant_id FROM branches WHERE id = '%d'", *claims.BranchID)).Scan(&restaurantID)
	if err != nil {
		return web.NewRequestError(errors.Wrap(err, "restaurant not found in "), http.StatusBadRequest)
	}

	if restaurantID == nil {
		return web.NewRequestError(errors.New("not found restaurant"), http.StatusBadRequest)
	}

	q := r.NewUpdate().Table("products").Where("deleted_at IS NULL "+
		"AND id = ? "+
		"AND restaurant_id = ?", request.ID, restaurantID)

	q.Set("name = ?", request.Name)
	q.Set("measure_unit_id =?", request.MeasureUnitID)
	q.Set("updated_at = ?", time.Now())
	q.Set("updated_by = ?", claims.UserId)

	_, err = q.Exec(ctx)
	if err != nil {
		return web.NewRequestError(errors.Wrap(err, "updating product"), http.StatusBadRequest)
	}

	return nil
}

func (r Repository) CashierUpdateColumns(ctx context.Context, request product.CashierUpdateRequest) error {
	claims, err := r.CheckClaims(ctx, auth.RoleCashier)
	if err != nil {
		return err
	}

	if err = r.ValidateStruct(&request, "ID"); err != nil {
		return err
	}

	var restaurantID *int64

	err = r.QueryRowContext(ctx, fmt.Sprintf("SELECT restaurant_id FROM branches WHERE id = '%d'", *claims.BranchID)).Scan(&restaurantID)
	if err != nil {
		return web.NewRequestError(errors.Wrap(err, "restaurant not found in "), http.StatusBadRequest)
	}

	if restaurantID == nil {
		return web.NewRequestError(errors.New("not found restaurant"), http.StatusBadRequest)
	}

	q := r.NewUpdate().Table("products").Where("deleted_at IS NULL "+
		"AND id = ? "+
		"AND restaurant_id = ?", request.ID, restaurantID)

	if request.Name != nil {
		q.Set("name = ?", request.Name)
	}

	if request.MeasureUnitID != nil {
		q.Set("measure_unit_id = ?", request.MeasureUnitID)
	}

	q.Set("updated_at = ?", time.Now())
	q.Set("updated_by = ?", claims.UserId)

	_, err = q.Exec(ctx)
	if err != nil {
		return web.NewRequestError(errors.Wrap(err, "updating product"), http.StatusBadRequest)
	}

	return nil
}

func (r Repository) CashierDelete(ctx context.Context, id int64) error {
	return r.DeleteRow(ctx, "products", id, auth.RoleCashier)
}

func (r Repository) CashierGetSpending(ctx context.Context, filter product.SpendingFilter) ([]product.CashierGetSpendingResponse, error) {
	claims, err := r.CheckClaims(ctx, auth.RoleCashier)
	if err != nil {
		return nil, err
	}

	var (
		response  []product.CashierGetSpendingResponse
		whereDate string
	)

	if filter.FromDate != nil {
		from := filter.FromDate.Format("02.01.2006")
		whereDate += fmt.Sprintf(` and to_char(or.created_at, 'DD.MM.YYYY') >= '%s' and case when frgh.from isnull then true else frgh.from <= '%s' end`, from, from)
	}
	if filter.ToDate != nil {
		to := filter.ToDate.Format("02.01.2006")
		whereDate += fmt.Sprintf(` and to_char(or.created_at, 'DD.MM.YYYY') <= '%s' and case when frgh.to isnull then true else frgh.to >= '%s' end`, to, to)
	}

	if claims.BranchID != nil {
		productQuery := fmt.Sprintf(`select
												p.id as id,
												p.name as name,
												sum(fr.amount*om.count) as amount,
												mu.name as measure_unit
											from foods f 
											    	 join food_recipe_group_histories frgh 
											    	     on f.id = frgh.food_id 
											    	 join food_recipe_groups frg 
											    	     on f.id = frg.food_id
													 join menus m
														  on f.id = any(m.food_ids)
													 join order_menu om
														  on m.id = om.menu_id 
											         join food_recipe fr 
											              on fr.id = any (frg.recipe_ids)
													 join products p
														  on fr.product_id = p.id
													 join measure_unit mu
														  on p.measure_unit_id = mu.id
													 join orders o
														 on om.order_id = o.id
													 join tables t
														 on o.table_id = t.id
													 join restaurants r 
													     on p.restaurant_id = r.id
											where t.branch_id='%d' and p.deleted_at isnull and om.deleted_at isnull and o.status = 'PAID' and om.status = 'PAID' %s
											group by p.id, mu.name`, *claims.BranchID, whereDate)

		rows, err := r.QueryContext(ctx, productQuery)
		if err != nil {
			return nil, err
		}

		if err = r.ScanRows(ctx, rows, &response); err != nil {
			return nil, err
		}
	} else {
		err = errors.New("cashier does not exists branch_id")
		return nil, web.NewRequestError(err, http.StatusBadRequest)
	}

	return response, nil
}

func NewRepository(DB *postgresql.Database) *Repository {
	return &Repository{DB}
}
