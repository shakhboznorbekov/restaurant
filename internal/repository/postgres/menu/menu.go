package menu

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/lib/pq"
	"github.com/pkg/errors"
	"github.com/restaurant/foundation/web"
	"github.com/restaurant/internal/auth"
	"github.com/restaurant/internal/entity"
	"github.com/restaurant/internal/pkg/repository/postgresql"
	"github.com/restaurant/internal/repository/postgres"
	"github.com/restaurant/internal/service/hashing"
	"github.com/restaurant/internal/service/menu"
	"net/http"
	"strings"
	"time"
)

type Repository struct {
	*postgresql.Database
}

// @admin

func (r Repository) AdminGetList(ctx context.Context, filter menu.Filter) ([]menu.AdminGetList, int, error) {
	claims, err := r.CheckClaims(ctx, auth.RoleAdmin)
	if err != nil {
		return nil, 0, err
	}

	whereQuery := fmt.Sprintf(`WHERE fc.deleted_at IS NULL AND m.deleted_at IS NULL AND m.branch_id in (select id from branches where restaurant_id = %d)`, *claims.RestaurantID)
	countWhereQuery := whereQuery

	var limitQuery, offsetQuery string
	if filter.Limit != nil {
		limitQuery = fmt.Sprintf(" LIMIT '%d'", *filter.Limit)
	}
	if filter.Offset != nil {
		offsetQuery = fmt.Sprintf(" OFFSET '%d'", *filter.Offset)
	}

	//whereQuery += fmt.Sprintf(" %s %s", limitQuery, offsetQuery)
	//query := fmt.Sprintf(`
	//	SELECT
	//		m.id,
	//		m.branch_id,
	//		m.food_id,
	//		m.status,
	//		m.new_price,
	//		m.old_price,
	//		b.name as branch_name,
	//		f.name as food_name
	//	FROM
	//	    menus as m
	//	LEFT JOIN public.branches b on b.id = m.branch_id
	//	LEFT OUTER JOIN public.foods f on f.id = m.food_id
	//	%s
	//`, whereQuery)

	query := fmt.Sprintf(`SELECT 
					fc.id AS category_id, 
					fc.name AS category_name, 
					json_agg(json_build_object('id', m.id, 'name', f.name, 'photos', f.photos, 'price', m.new_price, 'status', m.status)) AS menus
				FROM 
					food_category fc
				JOIN foods f ON f.category_id = fc.id
				JOIN menus m ON m.food_id = f.id
				%s 
				GROUP BY fc.id, fc.name ORDER BY fc.name %s %s;`, whereQuery, limitQuery, offsetQuery)

	list := make([]menu.AdminGetList, 0)

	rows, err := r.QueryContext(ctx, query)
	if err != nil {
		return nil, 0, web.NewRequestError(errors.Wrap(err, "select menus"), http.StatusInternalServerError)
	}

	err = r.ScanRows(ctx, rows, &list)
	if err != nil {
		return nil, 0, web.NewRequestError(errors.Wrap(err, "scanning menus"), http.StatusBadRequest)
	}

	countQuery := fmt.Sprintf(`SELECT
			count(fc.id)
		FROM
		    food_category fc
		JOIN foods f ON f.category_id = fc.id
		JOIN menus m ON m.food_id = f.id
		%s`, countWhereQuery)
	countRows, err := r.QueryContext(ctx, countQuery)

	if err == sql.ErrNoRows {
		return nil, 0, web.NewRequestError(postgres.ErrNotFound, http.StatusNotFound)
	}
	if err != nil {
		return nil, 0, web.NewRequestError(errors.Wrap(err, "selecting menu"), http.StatusBadRequest)
	}

	count := 0
	for countRows.Next() {
		if err = countRows.Scan(&count); err != nil {
			return nil, 0, web.NewRequestError(errors.Wrap(err, "scanning menu count"), http.StatusBadRequest)
		}
	}

	for k, v := range list {
		list[k].UserID = &claims.UserId
		for k1, v1 := range v.Menus {
			var photoLink pq.StringArray
			if v1.Photos != nil {
				for _, v2 := range *v1.Photos {
					baseLink := hashing.GenerateHash(r.ServerBaseUrl, v2)
					photoLink = append(photoLink, baseLink)
				}
				v.Menus[k1].Photos = &photoLink
			}
		}
	}

	return list, count, nil
}

func (r Repository) AdminGetDetail(ctx context.Context, id int64) (entity.Menu, error) {
	claims, err := r.CheckClaims(ctx, auth.RoleAdmin)
	if err != nil {
		return entity.Menu{}, err
	}

	var detail entity.Menu

	err = r.NewSelect().Model(&detail).Where("id = ? AND deleted_at IS NULL AND branch_id in (select id from branches where restaurant_id = ?)", id, claims.RestaurantID).Scan(ctx)
	if err != nil {
		return entity.Menu{}, err
	}
	return detail, nil
}

func (r Repository) AdminCreate(ctx context.Context, request menu.AdminCreateRequest) ([]menu.AdminCreateResponse, error) {
	claims, err := r.CheckClaims(ctx, auth.RoleAdmin)
	if err != nil {
		return nil, err
	}

	err = r.ValidateStruct(&request, "FoodID", "BranchID")
	if err != nil {
		return nil, err
	}

	var response []menu.AdminCreateResponse
	if len(request.BranchID) > 0 {
		for i := range request.BranchID {
			detail := menu.AdminCreateResponse{
				FoodID:      request.FoodID,
				NewPrice:    request.NewPrice,
				CreatedAt:   time.Now(),
				CreatedBy:   claims.UserId,
				Description: request.Description,
			}
			detail.BranchID = &request.BranchID[i]

			_, err = r.NewInsert().Model(&detail).Exec(ctx)
			if err != nil {
				return nil, web.NewRequestError(errors.Wrap(err, "creating menu"), http.StatusBadRequest)
			}

			response = append(response, detail)
		}
		return response, nil
	}

	err = errors.New("branch_id required field")
	return nil, err
}

func (r Repository) AdminUpdateAll(ctx context.Context, request menu.AdminUpdateRequest) error {
	claims, err := r.CheckClaims(ctx, auth.RoleAdmin)
	if err != nil {
		return err
	}

	if err = r.ValidateStruct(&request, "ID", "FoodID", "Status"); err != nil {
		return err
	}

	q := r.NewUpdate().Table("menus").Where(`deleted_at IS NULL AND id = ? AND branch_id in (select id from branches where restaurant_id = ?)`, request.ID, claims.RestaurantID)

	q.Set("food_id = ?", request.FoodID)
	q.Set("old_price = new_price")
	q.Set("new_price = ?", request.NewPrice)
	q.Set("branch_id = ?", request.BranchID)
	q.Set("status = ?", request.Status)
	q.Set("updated_at = ?", time.Now())
	q.Set("updated_by = ?", claims.UserId)
	q.Set("description = ?", request.Description)

	_, err = q.Exec(ctx)
	if err != nil {
		return web.NewRequestError(errors.Wrap(err, "updating menu"), http.StatusBadRequest)
	}

	return nil
}

func (r Repository) AdminUpdateColumns(ctx context.Context, request menu.AdminUpdateRequest) error {
	claims, err := r.CheckClaims(ctx, auth.RoleAdmin)
	if err != nil {
		return err
	}

	if err = r.ValidateStruct(&request, "ID"); err != nil {
		return err
	}

	q := r.NewUpdate().Table("menus").Where("deleted_at IS NULL AND id = ? "+
		"AND branch_id in (select id from branches where restaurant_id = ?)", request.ID, claims.RestaurantID)

	if request.FoodID != nil {
		q.Set("food_id = ?", request.FoodID)
	}
	if request.BranchID != nil {
		q.Set("branch_id = ?", request.BranchID)
	}
	if request.Status != nil {
		q.Set("status = ?", request.Status)
	}
	if request.NewPrice != nil {
		q.Set("old_price = new_price")
		q.Set("new_price = ?", request.NewPrice)
	}
	if request.Description != nil {
		q.Set("description = ?", request.Description)
	}

	q.Set("updated_at = ?", time.Now())
	q.Set("updated_by = ?", claims.UserId)

	_, err = q.Exec(ctx)
	if err != nil {
		return web.NewRequestError(errors.Wrap(err, "updating menu"), http.StatusBadRequest)
	}

	if request.Status != nil || request.BranchID != nil {
		_, err = r.ExecContext(ctx, fmt.Sprintf(`
			UPDATE branches
    		SET menu_names = (SELECT string_agg( ' {' || text(m.id) || '} ' || f.name, ' ') AS aggregated_names FROM menus m JOIN foods f ON m.food_id = f.id WHERE m.branch_id = branches.id AND m.deleted_at IS NULL AND m.status = 'active')
    		WHERE restaurant_id = %d`, *claims.RestaurantID))
		if err != nil {
			return web.NewRequestError(errors.Wrap(err, "updating branch menu_names"), http.StatusBadRequest)
		}
	}

	return nil
}

func (r Repository) AdminDelete(ctx context.Context, id int64) error {
	claims, err := r.CheckClaims(ctx, auth.RoleAdmin)
	if err != nil {
		return err
	}

	err = r.DeleteRow(ctx, "menus", id, auth.RoleAdmin)
	if err != nil {
		return err
	}

	_, err = r.ExecContext(ctx, fmt.Sprintf(`
			UPDATE branches
    		SET menu_names = (SELECT string_agg( ' {' || text(m.id) || '} ' || f.name, ' ') AS aggregated_names FROM menus m JOIN foods f ON m.food_id = f.id WHERE m.branch_id = branches.id AND m.deleted_at IS NULL AND m.status = 'active')
    		WHERE restaurant_id = %d`, *claims.RestaurantID))
	if err != nil {
		return web.NewRequestError(errors.Wrap(err, "updating branch menu_names"), http.StatusBadRequest)
	}
	return nil
}

// @branch

func (r Repository) BranchGetList(ctx context.Context, filter menu.Filter) ([]menu.BranchGetList, int, error) {
	claims, err := r.CheckClaims(ctx, auth.RoleBranch)
	if err != nil {
		return nil, 0, err
	}

	whereQuery := fmt.Sprintf(`WHERE fc.deleted_at IS NULL AND m.deleted_at IS NULL and m.branch_id = '%d'`, *claims.BranchID)
	if filter.PrinterID != nil {
		whereQuery += fmt.Sprintf(" AND m.printer_id = %d", *filter.PrinterID)
	}
	if filter.Printer != nil {
		if !*filter.Printer {
			whereQuery += fmt.Sprintf(" AND m.printer_id IS NULL")
		}
	}
	countWhereQuery := whereQuery

	var limitQuery, offsetQuery string
	if filter.Limit != nil {
		limitQuery = fmt.Sprintf(" LIMIT '%d'", *filter.Limit)
	}
	if filter.Offset != nil {
		offsetQuery = fmt.Sprintf(" OFFSET '%d'", *filter.Offset)
	}

	//whereQuery += fmt.Sprintf(" %s %s", limitQuery, offsetQuery)
	//query := fmt.Sprintf(`
	//	SELECT
	//		m.id,
	//		m.branch_id,
	//		m.food_id,
	//		m.status,
	//		m.new_price,
	//		m.old_price,
	//		b.name as branch_name,
	//		f.name as food_name,
	//		f.photos[1] as photo
	//	FROM
	//	    menus as m
	//	LEFT JOIN public.branches b on b.id = m.branch_id
	//	LEFT OUTER JOIN public.foods f on f.id = m.food_id
	//	%s
	//`, whereQuery)

	query := fmt.Sprintf(`SELECT 
					fc.id AS category_id, 
					fc.name AS category_name, 
					json_agg(json_build_object('id', m.id, 'name', f.name, 'photos', f.photos, 'price', m.new_price, 'status', m.status, 'printer', CASE WHEN m.printer_id is null THEN false ELSE true END )) AS menus
				FROM 
					food_category fc
				JOIN foods f ON f.category_id = fc.id
				JOIN menus m ON m.food_id = f.id
				%s 
				GROUP BY fc.id, fc.name ORDER BY fc.name %s %s;`, whereQuery, limitQuery, offsetQuery)

	list := make([]menu.BranchGetList, 0)

	rows, err := r.QueryContext(ctx, query)
	if err != nil {
		return nil, 0, web.NewRequestError(errors.Wrap(err, "select menus"), http.StatusInternalServerError)
	}

	err = r.ScanRows(ctx, rows, &list)
	if err != nil {
		return nil, 0, web.NewRequestError(errors.Wrap(err, "scanning menus"), http.StatusBadRequest)
	}

	for k, v := range list {
		list[k].UserID = &claims.UserId
		for k1, v1 := range v.Menus {
			var photoLink pq.StringArray
			if v1.Photos != nil {
				for _, v2 := range *v1.Photos {
					baseLink := hashing.GenerateHash(r.ServerBaseUrl, v2)
					photoLink = append(photoLink, baseLink)
				}
				v.Menus[k1].Photos = &photoLink
			}
		}
	}

	countQuery := fmt.Sprintf(`
		SELECT
			count(m.id)
		FROM
		    food_category fc
		JOIN foods f ON f.category_id = fc.id
		JOIN menus m ON m.food_id = f.id
		%s
	`, countWhereQuery)
	countRows, err := r.QueryContext(ctx, countQuery)
	if err == sql.ErrNoRows {
		return nil, 0, web.NewRequestError(postgres.ErrNotFound, http.StatusNotFound)
	}
	if err != nil {
		return nil, 0, web.NewRequestError(errors.Wrap(err, "selecting menu"), http.StatusBadRequest)
	}

	count := 0
	for countRows.Next() {
		if err = countRows.Scan(&count); err != nil {
			return nil, 0, web.NewRequestError(errors.Wrap(err, "scanning menu count"), http.StatusBadRequest)
		}
	}

	return list, count, nil
}

func (r Repository) BranchGetDetail(ctx context.Context, id int64) (entity.Menu, error) {
	claims, err := r.CheckClaims(ctx, auth.RoleBranch)
	if err != nil {
		return entity.Menu{}, err
	}

	var detail entity.Menu

	err = r.NewSelect().Model(&detail).Where("id = ? AND branch_id = ? AND deleted_at IS NULL", id, claims.BranchID).Scan(ctx)
	if err != nil {
		return entity.Menu{}, err
	}
	return detail, nil
}

func (r Repository) BranchCreate(ctx context.Context, request menu.BranchCreateRequest) (menu.BranchCreateResponse, error) {
	claims, err := r.CheckClaims(ctx, auth.RoleBranch)
	if err != nil {
		return menu.BranchCreateResponse{}, err
	}

	err = r.ValidateStruct(&request, "FoodID", "Status", "New_Price")
	if err != nil {
		return menu.BranchCreateResponse{}, err
	}

	response := menu.BranchCreateResponse{
		FoodID:      request.FoodID,
		NewPrice:    request.NewPrice,
		CreatedAt:   time.Now(),
		CreatedBy:   claims.UserId,
		BranchID:    claims.BranchID,
		Description: request.Description,
	}

	_, err = r.NewInsert().Model(&response).Exec(ctx)
	if err != nil {
		return menu.BranchCreateResponse{}, web.NewRequestError(errors.Wrap(err, "creating menu"), http.StatusBadRequest)
	}

	return response, nil
}

func (r Repository) BranchUpdateAll(ctx context.Context, request menu.BranchUpdateRequest) error {
	claims, err := r.CheckClaims(ctx, auth.RoleBranch)
	if err != nil {
		return err
	}

	if err = r.ValidateStruct(&request,
		"ID", "FoodID", "Status", "NewPrice",
	); err != nil {
		return err
	}

	q := r.NewUpdate().Table("menus").Where("deleted_at IS NULL AND id = ? AND branch_id = ?",
		request.ID, claims.BranchID)

	q.Set("food_id = ?", request.FoodID)
	q.Set("old_price = new_price")
	q.Set("new_price = ?", request.NewPrice)
	q.Set("status = ?", request.Status)
	q.Set("updated_at = ?", time.Now())
	q.Set("updated_by = ?", claims.UserId)
	q.Set("description = ?", request.Description)

	_, err = q.Exec(ctx)
	if err != nil {
		return web.NewRequestError(errors.Wrap(err, "updating menu"), http.StatusBadRequest)
	}

	return nil
}

func (r Repository) BranchUpdateColumns(ctx context.Context, request menu.BranchUpdateRequest) error {
	claims, err := r.CheckClaims(ctx, auth.RoleBranch)
	if err != nil {
		return err
	}

	if err = r.ValidateStruct(&request, "ID"); err != nil {
		return err
	}

	q := r.NewUpdate().Table("menus").Where("deleted_at IS NULL AND id = ? AND branch_id = ?",
		request.ID, claims.BranchID)

	if request.FoodID != nil {
		q.Set("food_id = ?", request.FoodID)
	}
	if request.NewPrice != nil {
		q.Set("old_price = new_price")
		q.Set("new_price = ?", request.NewPrice)
	}
	if request.Status != nil {
		q.Set("status = ?", request.Status)
	}
	if request.Description != nil {
		q.Set("description = ?", request.Description)
	}

	q.Set("updated_at = ?", time.Now())
	q.Set("updated_by = ?", claims.UserId)

	_, err = q.Exec(ctx)
	if err != nil {
		return web.NewRequestError(errors.Wrap(err, "updating menu"), http.StatusBadRequest)
	}

	if request.Status != nil {
		_, err = r.ExecContext(ctx, fmt.Sprintf(`
			UPDATE branches
    		SET menu_names = (SELECT string_agg( ' {' || text(m.id) || '} ' || f.name, ' ') AS aggregated_names FROM menus m JOIN foods f ON m.food_id = f.id WHERE m.branch_id = branches.id AND m.deleted_at IS NULL AND m.status = 'active')
    		WHERE id = %d`, *claims.BranchID))
		if err != nil {
			return web.NewRequestError(errors.Wrap(err, "updating branch menu_names"), http.StatusBadRequest)
		}
	}
	return nil
}

func (r Repository) BranchDelete(ctx context.Context, id int64) error {
	claims, err := r.CheckClaims(ctx, auth.RoleBranch)
	if err != nil {
		return err
	}

	err = r.DeleteRow(ctx, "menus", id, auth.RoleBranch)
	if err != nil {
		return err
	}
	_, err = r.ExecContext(ctx, fmt.Sprintf(`
			UPDATE branches
    		SET menu_names = (SELECT string_agg( ' {' || text(m.id) || '} ' || f.name, ' ') AS aggregated_names FROM menus m JOIN foods f ON m.food_id = f.id WHERE m.branch_id = branches.id AND m.deleted_at IS NULL AND m.status = 'active')
    		WHERE id = %d`, *claims.BranchID))
	if err != nil {
		return web.NewRequestError(errors.Wrap(err, "updating branch menu_names"), http.StatusBadRequest)
	}
	return nil
}

func (r Repository) BranchUpdatePrinterID(ctx context.Context, request menu.BranchUpdatePrinterIDRequest) error {
	claims, err := r.CheckClaims(ctx, auth.RoleBranch)
	if err != nil {
		return err
	}
	if err = r.ValidateStruct(&request, "PrinterID", "MenuIds"); err != nil {
		return err
	}

	for _, v := range request.MenuIds {
		_, err = r.NewUpdate().Table("menus").Where("id = ? AND branch_id = ?", v, claims.BranchID).Set("printer_id = ?", *request.PrinterID).Exec(ctx)
		if err != nil {
			return web.NewRequestError(errors.Wrap(err, "updating menu token"), http.StatusBadRequest)
		}
	}
	return nil
}

func (r Repository) BranchDeletePrinterID(ctx context.Context, menuID int64) error {
	claims, err := r.CheckClaims(ctx, auth.RoleBranch)
	if err != nil {
		return err
	}

	_, err = r.NewUpdate().Table("menus").Where("id = ? AND branch_id = ?", menuID, claims.BranchID).Set("printer_id = null").Exec(ctx)
	if err != nil {
		return web.NewRequestError(errors.Wrap(err, "updating menu token"), http.StatusBadRequest)
	}
	return nil
}

// @client

func (r Repository) ClientGetList(ctx context.Context, filter menu.Filter) ([]menu.ClientGetList, error) {
	//_, err := r.CheckClaims(ctx, auth.RoleClient)
	//if err != nil {
	//	return nil, err
	//}

	whereQuery := fmt.Sprintf(`WHERE f.deleted_at IS NULL AND fc.deleted_at IS NULL AND m.status = 'active' AND m.deleted_at IS NULL`)

	if filter.Search != nil {
		search := strings.Replace(*filter.Search, " ", "", -1)
		whereQuery += fmt.Sprintf(` AND 
			(
				REPLACE(f.name, ' ', '') ilike '%s' OR
				REPLACE(fc.name, ' ', '') ilike '%s'
				
			)`,
			"%"+search+"%",
			"%"+search+"%",
		)
	}

	if filter.BranchID != nil {
		whereQuery += fmt.Sprintf(` and f.branch_id = '%d'`, *filter.BranchID)
	}

	query := fmt.Sprintf(`
				SELECT 
					fc.id AS category_id, 
					fc.name AS category_name, 
					json_agg(json_build_object('id', m.id, 'name', f.name, 'photos', f.photos, 'price', m.new_price)) AS menus
				FROM 
					food_category fc
				JOIN 
					foods f ON f.category_id = fc.id
				JOIN
					menus m ON m.food_id = f.id
				%s
				GROUP BY 
					fc.id, fc.name;
	`, whereQuery)

	list := make([]menu.ClientGetList, 0)

	rows, err := r.QueryContext(ctx, query)
	if err != nil {
		return nil, web.NewRequestError(errors.Wrap(err, "select foods"), http.StatusInternalServerError)
	}

	err = r.ScanRows(ctx, rows, &list)
	if err != nil {
		return nil, web.NewRequestError(errors.Wrap(err, "scanning foods"), http.StatusBadRequest)
	}

	for _, v := range list {
		for k1, v1 := range v.Menus {
			var photoLink pq.StringArray
			if v1.Photos != nil {
				for _, v2 := range *v1.Photos {
					baseLink := hashing.GenerateHash(r.ServerBaseUrl, v2)
					photoLink = append(photoLink, baseLink)
				}
				v.Menus[k1].Photos = &photoLink
			}
		}
	}

	//countQuery := fmt.Sprintf(`
	//	SELECT
	//		count(f.id)
	//	FROM
	//	    foods as f
	//	JOIN
	//				foods f ON f.category_id = fc.id
	//	%s
	//`, countWhereQuery)
	//
	//countRows, err := r.QueryContext(ctx, countQuery)
	//if err == sql.ErrNoRows {
	//	return nil, 0, web.NewRequestError(postgres.ErrNotFound, http.StatusNotFound)
	//}
	//if err != nil {
	//	return nil, 0, web.NewRequestError(errors.Wrap(err, "selecting foods"), http.StatusBadRequest)
	//}
	//
	//count := 0
	//
	//for countRows.Next() {
	//	if err = countRows.Scan(&count); err != nil {
	//		return nil, 0, web.NewRequestError(errors.Wrap(err, "scanning user count"), http.StatusBadRequest)
	//	}
	//}

	//for k, v := range list {
	//	if v.Photos != nil {
	//		var photoLink pq.StringArray
	//		for _, v1 := range *v.Photos {
	//			baseLink := hashing.GenerateHash(r.ServerBaseUrl, v1)
	//			photoLink = append(photoLink, baseLink)
	//		}
	//		list[k].Photos = &photoLink
	//	}
	//}

	return list, nil
}

func (r Repository) ClientGetDetail(ctx context.Context, id int64) (menu.ClientGetDetail, error) {
	table := "foods"
	whereQuery := fmt.Sprintf(`WHERE %s.deleted_at IS NULL AND %s.id = %d`, table, table, id)

	query := fmt.Sprintf(`
			SELECT
				menus.id,
				foods.name,
				foods.photos,
				menus.new_price
			FROM foods
			LEFT JOIN menus ON foods.id = menus.food_id
			%s
	`, whereQuery)

	var detail menu.ClientGetDetail

	rows, err := r.QueryContext(ctx, query)
	if err != nil {
		return menu.ClientGetDetail{},
			web.NewRequestError(errors.Wrap(err, "select foods"), http.StatusInternalServerError)
	}

	for rows.Next() {
		err = rows.Scan(&detail.ID, &detail.Name, &detail.Photos, &detail.Price)
		if err != nil {
			return menu.ClientGetDetail{},
				web.NewRequestError(errors.Wrap(err, "scanning foods"), http.StatusBadRequest)
		}
	}

	if detail.Photos != nil {
		var basePhoto pq.StringArray
		for _, v := range *detail.Photos {
			baseLink := hashing.GenerateHash(r.ServerBaseUrl, v)
			basePhoto = append(basePhoto, baseLink)
		}
		detail.Photos = &basePhoto
	}

	return detail, nil
}

func (r Repository) ClientGetListByCategoryID(ctx context.Context, foodCategoryID int, filter menu.Filter) ([]menu.ClientGetListByCategoryID, error) {
	claims, err := r.CheckClaims(ctx, auth.RoleClient)
	if err != nil {
		claims.UserId = 0
	}
	whereQuery := fmt.Sprintf(`WHERE b.deleted_at IS NULL AND m.deleted_at IS NULL AND fc.id = %d`, foodCategoryID)
	location := true
	if filter.Lat == nil || filter.Lon == nil {
		var l float64
		filter.Lat = &l
		filter.Lon = &l
		location = false
	}

	todayInt := time.Now().Weekday()
	var today string
	switch todayInt {
	case 1:
		today = "Monday"
	case 2:
		today = "Tuesday"
	case 3:
		today = "Wednesday"
	case 4:
		today = "Thursday"
	case 5:
		today = "Friday"
	case 6:
		today = "Saturday"
	case 7:
		today = "Sunday"
	}

	queryOrder := ""
	if location {
		queryOrder += " ORDER BY distance"
	} else {
		queryOrder += " ORDER BY b.rate"
	}

	query := fmt.Sprintf(`
				SELECT
				    b.id,
				    b.status,
				    b.location,
				    b.photos as photos,
				    b.work_time->>'%s' as work_time_today,
				    b.name,
				    b.category_id,
				    rc.name as category_name,
				   	CASE WHEN %v THEN ST_DistanceSphere(
				            ST_SetSRID(ST_MakePoint(CAST(b.location->>'lon' AS float), CAST(b.location->>'lat' AS float)), 4326),
				            ST_SetSRID(ST_MakePoint(%v, %v), 4326)
				        ) END as distance,
				    b.rate,
					CASE WHEN (SELECT id FROM branch_likes WHERE branch_id = b.id AND user_id = %d) IS NOT NULL THEN true ELSE false END AS is_liked
				FROM
				    branches AS b
				LEFT OUTER JOIN restaurant_category AS rc ON rc.id = b.category_id
				LEFT OUTER JOIN menus m ON m.branch_id = b.id
				LEFT OUTER JOIN foods f ON f.id = m.food_id
				LEFT OUTER JOIN food_category fc ON fc.id = f.category_id
				%s
				GROUP BY b.id, b.status, b.location, rc.name,b.photos, b.name, b.category_id, b.rate
				%s
	`, today, location, *filter.Lon, *filter.Lat, claims.UserId, whereQuery, queryOrder)

	list := make([]menu.ClientGetListByCategoryID, 0)

	rows, err := r.QueryContext(ctx, query)
	if err != nil {
		return nil, web.NewRequestError(errors.Wrap(err, "select branches"), http.StatusInternalServerError)
	}

	err = r.ScanRows(ctx, rows, &list)
	if err != nil {
		return nil, web.NewRequestError(errors.Wrap(err, "scanning branches"), http.StatusBadRequest)
	}

	for k, v := range list {
		var photoLink pq.StringArray
		for _, v2 := range *v.Photos {
			baseLink := hashing.GenerateHash(r.ServerBaseUrl, v2)
			photoLink = append(photoLink, baseLink)
		}
		list[k].Photos = &photoLink
	}
	return list, nil
}

//  @cashier

func (r Repository) CashierUpdateColumns(ctx context.Context, request menu.CashierUpdateMenuStatus) error {
	claims, err := r.CheckClaims(ctx, auth.RoleCashier)
	if err != nil {
		return err
	}

	if err = r.ValidateStruct(&request, "ID"); err != nil {
		return err
	}

	q := r.NewUpdate().Table("menus").Where("deleted_at IS NULL AND id = ? AND branch_id = ?",
		request.ID, claims.BranchID)

	if request.Status != nil {
		q.Set("status = ?", request.Status)
	}

	q.Set("updated_at = ?", time.Now())
	q.Set("updated_by = ?", claims.UserId)

	_, err = q.Exec(ctx)
	if err != nil {
		return web.NewRequestError(errors.Wrap(err, "updating menu"), http.StatusBadRequest)
	}

	return nil
}

func (r Repository) CashierGetList(ctx context.Context, filter menu.Filter) ([]menu.CashierGetList, int, error) {
	claims, err := r.CheckClaims(ctx, auth.RoleCashier)
	if err != nil {
		return nil, 0, err
	}

	whereQuery := fmt.Sprintf(`WHERE fc.deleted_at IS NULL AND m.deleted_at IS NULL and m.branch_id = '%d'`, *claims.BranchID)
	if filter.PrinterID != nil {
		whereQuery += fmt.Sprintf(" AND m.printer_id = %d", *filter.PrinterID)
	}
	if filter.Printer != nil {
		if !*filter.Printer {
			whereQuery += fmt.Sprintf(" AND m.printer_id IS NULL")
		}
	}
	countWhereQuery := whereQuery

	var limitQuery, offsetQuery string
	if filter.Limit != nil {
		limitQuery = fmt.Sprintf(" LIMIT '%d'", *filter.Limit)
	}
	if filter.Offset != nil {
		offsetQuery = fmt.Sprintf(" OFFSET '%d'", *filter.Offset)
	}

	//whereQuery += fmt.Sprintf(" %s %s", limitQuery, offsetQuery)
	//query := fmt.Sprintf(`
	//	SELECT
	//		m.id,
	//		m.branch_id,
	//		m.food_id,
	//		m.status,
	//		m.new_price,
	//		m.old_price,
	//		b.name as branch_name,
	//		f.name as food_name,
	//		f.photos[1] as photo
	//	FROM
	//	    menus as m
	//	LEFT JOIN public.branches b on b.id = m.branch_id
	//	LEFT OUTER JOIN public.foods f on f.id = m.food_id
	//	%s
	//`, whereQuery)

	query := fmt.Sprintf(`SELECT 
					fc.id AS category_id, 
					fc.name AS category_name, 
					json_agg(json_build_object('id', m.id, 'name', f.name, 'photos', f.photos, 'price', m.new_price, 'status', m.status, 'printer', CASE WHEN m.printer_id is null THEN false ELSE true END )) AS menus
				FROM 
					food_category fc
				JOIN foods f ON f.category_id = fc.id
				JOIN menus m ON m.food_id = f.id
				%s 
				GROUP BY fc.id, fc.name ORDER BY fc.name %s %s;`, whereQuery, limitQuery, offsetQuery)

	list := make([]menu.CashierGetList, 0)

	rows, err := r.QueryContext(ctx, query)
	if err != nil {
		return nil, 0, web.NewRequestError(errors.Wrap(err, "select menus"), http.StatusBadRequest)
	}

	err = r.ScanRows(ctx, rows, &list)
	if err != nil {
		return nil, 0, web.NewRequestError(errors.Wrap(err, "scanning menus"), http.StatusInternalServerError)
	}

	for k, v := range list {
		list[k].UserID = &claims.UserId
		for k1, v1 := range v.Menus {
			var photoLink pq.StringArray
			if v1.Photos != nil {
				for _, v2 := range *v1.Photos {
					baseLink := hashing.GenerateHash(r.ServerBaseUrl, v2)
					photoLink = append(photoLink, baseLink)
				}
				v.Menus[k1].Photos = &photoLink
			}
		}
	}

	countQuery := fmt.Sprintf(`
		SELECT
			count(m.id)
		FROM
		    food_category fc
		JOIN foods f ON f.category_id = fc.id
		JOIN menus m ON m.food_id = f.id
		%s
	`, countWhereQuery)
	countRows, err := r.QueryContext(ctx, countQuery)
	if err == sql.ErrNoRows {
		return nil, 0, web.NewRequestError(postgres.ErrNotFound, http.StatusBadRequest)
	}
	if err != nil {
		return nil, 0, web.NewRequestError(errors.Wrap(err, "selecting menu"), http.StatusBadRequest)
	}

	count := 0
	for countRows.Next() {
		if err = countRows.Scan(&count); err != nil {
			return nil, 0, web.NewRequestError(errors.Wrap(err, "scanning menu count"), http.StatusBadRequest)
		}
	}

	return list, count, nil
}

func (r Repository) CashierGetDetail(ctx context.Context, id int64) (entity.Menu, error) {
	claims, err := r.CheckClaims(ctx, auth.RoleCashier)
	if err != nil {
		return entity.Menu{}, err
	}

	var detail entity.Menu

	err = r.NewSelect().Model(&detail).Where("id = ? AND branch_id = ? AND deleted_at IS NULL", id, claims.BranchID).Scan(ctx)
	if err != nil {
		return entity.Menu{}, err
	}
	return detail, nil
}

func (r Repository) CashierCreate(ctx context.Context, request menu.CashierCreateRequest) (menu.CashierCreateResponse, error) {
	claims, err := r.CheckClaims(ctx, auth.RoleCashier)
	if err != nil {
		return menu.CashierCreateResponse{}, err
	}

	err = r.ValidateStruct(&request, "FoodID", "Status", "New_Price")
	if err != nil {
		return menu.CashierCreateResponse{}, err
	}

	response := menu.CashierCreateResponse{
		FoodID:      request.FoodID,
		NewPrice:    request.NewPrice,
		CreatedAt:   time.Now(),
		CreatedBy:   claims.UserId,
		BranchID:    claims.BranchID,
		Description: request.Description,
	}

	_, err = r.NewInsert().Model(&response).Exec(ctx)
	if err != nil {
		return menu.CashierCreateResponse{}, web.NewRequestError(errors.Wrap(err, "creating menu"), http.StatusBadRequest)
	}

	return response, nil
}

func (r Repository) CashierUpdateAll(ctx context.Context, request menu.CashierUpdateRequest) error {
	claims, err := r.CheckClaims(ctx, auth.RoleCashier)
	if err != nil {
		return err
	}

	if err = r.ValidateStruct(&request,
		"ID", "FoodID", "Status", "NewPrice",
	); err != nil {
		return err
	}

	q := r.NewUpdate().Table("menus").Where("deleted_at IS NULL AND id = ? AND branch_id = ?",
		request.ID, claims.BranchID)

	q.Set("food_id = ?", request.FoodID)
	q.Set("old_price = new_price")
	q.Set("new_price = ?", request.NewPrice)
	q.Set("status = ?", request.Status)
	q.Set("updated_at = ?", time.Now())
	q.Set("updated_by = ?", claims.UserId)
	q.Set("description = ?", request.Description)

	_, err = q.Exec(ctx)
	if err != nil {
		return web.NewRequestError(errors.Wrap(err, "updating menu"), http.StatusBadRequest)
	}

	return nil
}

func (r Repository) CashierUpdateColumn(ctx context.Context, request menu.CashierUpdateRequest) error {
	claims, err := r.CheckClaims(ctx, auth.RoleCashier)
	if err != nil {
		return err
	}

	if err = r.ValidateStruct(&request, "ID"); err != nil {
		return err
	}

	q := r.NewUpdate().Table("menus").Where("deleted_at IS NULL AND id = ? AND branch_id = ?",
		request.ID, claims.BranchID)

	if request.FoodID != nil {
		q.Set("food_id = ?", request.FoodID)
	}
	if request.NewPrice != nil {
		q.Set("old_price = new_price")
		q.Set("new_price = ?", request.NewPrice)
	}
	if request.Status != nil {
		q.Set("status = ?", request.Status)
	}
	if request.Description != nil {
		q.Set("description = ?", request.Description)
	}

	q.Set("updated_at = ?", time.Now())
	q.Set("updated_by = ?", claims.UserId)

	_, err = q.Exec(ctx)
	if err != nil {
		return web.NewRequestError(errors.Wrap(err, "updating menu"), http.StatusBadRequest)
	}

	if request.Status != nil {
		_, err = r.ExecContext(ctx, fmt.Sprintf(`
			UPDATE branches
    		SET menu_names = (SELECT string_agg( ' {' || text(m.id) || '} ' || f.name, ' ') AS aggregated_names FROM menus m JOIN foods f ON m.food_id = f.id WHERE m.branch_id = branches.id AND m.deleted_at IS NULL AND m.status = 'active')
    		WHERE id = %d`, *claims.BranchID))
		if err != nil {
			return web.NewRequestError(errors.Wrap(err, "updating branch menu_names"), http.StatusBadRequest)
		}
	}
	return nil
}

func (r Repository) CashierDelete(ctx context.Context, id int64) error {
	claims, err := r.CheckClaims(ctx, auth.RoleCashier)
	if err != nil {
		return err
	}

	err = r.DeleteRow(ctx, "menus", id, auth.RoleCashier)
	if err != nil {
		return err
	}
	_, err = r.ExecContext(ctx, fmt.Sprintf(`
			UPDATE branches
    		SET menu_names = (SELECT string_agg( ' {' || text(m.id) || '} ' || f.name, ' ') AS aggregated_names FROM menus m JOIN foods f ON m.food_id = f.id WHERE m.branch_id = branches.id AND m.deleted_at IS NULL AND m.status = 'active')
    		WHERE id = %d`, *claims.BranchID))
	if err != nil {
		return web.NewRequestError(errors.Wrap(err, "updating branch menu_names"), http.StatusBadRequest)
	}
	return nil
}

func (r Repository) CashierUpdatePrinterID(ctx context.Context, request menu.CashierUpdatePrinterIDRequest) error {
	claims, err := r.CheckClaims(ctx, auth.RoleCashier)
	if err != nil {
		return err
	}
	if err = r.ValidateStruct(&request, "PrinterID", "MenuIds"); err != nil {
		return err
	}

	for _, v := range request.MenuIds {
		_, err = r.NewUpdate().Table("menus").Where("id = ? AND branch_id = ?", v, claims.BranchID).Set("printer_id = ?", *request.PrinterID).Exec(ctx)
		if err != nil {
			return web.NewRequestError(errors.Wrap(err, "updating menu token"), http.StatusBadRequest)
		}
	}
	return nil
}

func (r Repository) CashierDeletePrinterID(ctx context.Context, menuID int64) error {
	claims, err := r.CheckClaims(ctx, auth.RoleCashier)
	if err != nil {
		return err
	}

	_, err = r.NewUpdate().Table("menus").Where("id = ? AND branch_id = ?", menuID, claims.BranchID).Set("printer_id = null").Exec(ctx)
	if err != nil {
		return web.NewRequestError(errors.Wrap(err, "updating menu token"), http.StatusBadRequest)
	}
	return nil
}

// @waiter

func (r Repository) WaiterGetMenuList(ctx context.Context, filter menu.Filter) ([]menu.WaiterGetMenuListResponse, error) {
	claims, err := r.CheckClaims(ctx, auth.RoleWaiter)
	if err != nil {
		return nil, err
	}

	fltr := fmt.Sprintf(` AND m.deleted_at isnull AND m.branch_id='%d' AND f.deleted_at isnull `, *claims.BranchID)
	if filter.Search != nil {
		fltr += fmt.Sprintf(` AND f.name ilike '%s'`, "%"+*filter.Search+"%")
	}
	if filter.CategoryId != nil {
		fltr += fmt.Sprintf(` AND f.category_id='%d'`, *filter.CategoryId)
	}

	response := make([]menu.WaiterGetMenuListResponse, 0)

	query := fmt.Sprintf(`SELECT 
								fc.id,
								fc.name
						  FROM 
						  		food_category fc
						  WHERE 
						  		fc.deleted_at isnull 
									AND
								(select count(m.id) from menus m join foods f on f.id = m.food_id where f.category_id=fc.id %s) != 0`, fltr)
	rows, err := r.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var row menu.WaiterGetMenuListResponse

		if err = rows.Scan(&row.Id, &row.Name); err != nil {
			return nil, err
		}

		query = fmt.Sprintf(`SELECT 
									m.id as id, 
									f.name as name, 
									m.status as status, 
									f.photos[1] as photo,
									m.new_price as price
							 FROM menus m 
							 JOIN foods f 
							 		ON f.id = m.food_id 
							 WHERE f.category_id='%d' %s`, row.Id, fltr)
		rows, err := r.QueryContext(ctx, query)
		if err != nil {
			return nil, err
		}

		var menus []menu.WaiterMenu
		if err = r.ScanRows(ctx, rows, &menus); err != nil {
			return nil, err
		}

		for i := range menus {
			if menus[i].Photo != nil {
				link := hashing.GenerateHash(r.ServerBaseUrl, *menus[i].Photo)

				menus[i].Photo = &link
			}
		}

		row.Menus = menus

		response = append(response, row)
	}

	return response, nil
}

func NewRepository(DB *postgresql.Database) *Repository {
	return &Repository{DB}
}
