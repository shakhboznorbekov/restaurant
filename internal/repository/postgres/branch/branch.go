package branch

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/google/uuid"
	"github.com/lib/pq"
	"github.com/pkg/errors"
	"github.com/restaurant/foundation/web"
	"github.com/restaurant/internal/auth"
	"github.com/restaurant/internal/entity"
	"github.com/restaurant/internal/pkg/repository/postgresql"
	"github.com/restaurant/internal/repository/postgres"
	"github.com/restaurant/internal/service/branch"
	"github.com/restaurant/internal/service/hashing"
	"github.com/uptrace/bun/dialect/pgdialect"
	"net/http"
	"slices"
	"time"
)

type Repository struct {
	*postgresql.Database
}

// @admin

func (r Repository) AdminGetList(ctx context.Context, filter branch.Filter) ([]branch.AdminGetList, int, error) {
	claims, err := r.CheckClaims(ctx, auth.RoleAdmin)
	if err != nil {
		return nil, 0, err
	}

	table := "branches"
	whereQuery := fmt.Sprintf(`WHERE b.deleted_at IS NULL AND b.restaurant_id = '%d'`, *claims.RestaurantID)
	countWhereQuery := whereQuery

	var limitQuery, offsetQuery string
	if filter.Limit != nil {
		limitQuery = fmt.Sprintf(" LIMIT '%d'", *filter.Limit)
	}
	if filter.Offset != nil {
		offsetQuery = fmt.Sprintf(" OFFSET '%d'", *filter.Offset)
	}

	whereQuery += fmt.Sprintf(" %s %s", limitQuery, offsetQuery)

	query := fmt.Sprintf(`
					SELECT 
					    b.id, 
					    b.status, 
					    b.location, 
					    b.photos, 
					    b.work_time,
					    b.name,
					    b.category_id,
					    rc.name as category_name
					FROM 
					    branches as b
					LEFT OUTER JOIN restaurant_category as rc ON rc.id = b.category_id
					%s`, whereQuery)

	list := make([]branch.AdminGetList, 0)

	rows, err := r.QueryContext(ctx, query)
	if err != nil {
		return nil, 0, web.NewRequestError(errors.Wrap(err, "select branches"), http.StatusInternalServerError)
	}

	err = r.ScanRows(ctx, rows, &list)
	if err != nil {
		return nil, 0, web.NewRequestError(errors.Wrap(err, "scanning branches"), http.StatusBadRequest)
	}

	countQuery := fmt.Sprintf(`
		SELECT
			count(id)
		FROM
		    %s as b
		%s
	`, table, countWhereQuery)
	countRows, err := r.QueryContext(ctx, countQuery)
	if err == sql.ErrNoRows {
		return nil, 0, web.NewRequestError(postgres.ErrNotFound, http.StatusNotFound)
	}
	if err != nil {
		return nil, 0, web.NewRequestError(errors.Wrap(err, "selecting branch"), http.StatusBadRequest)
	}

	count := 0

	for countRows.Next() {
		if err = countRows.Scan(&count); err != nil {
			return nil, 0, web.NewRequestError(errors.Wrap(err, "scanning branch count"), http.StatusBadRequest)
		}
	}

	for k, v := range list {
		var basePhotos pq.StringArray
		if v.Photos != nil {
			for _, v1 := range *v.Photos {
				baseLink := hashing.GenerateHash(r.ServerBaseUrl, v1)
				basePhotos = append(basePhotos, baseLink)
			}
		}
		list[k].Photos = &basePhotos
	}

	return list, count, nil
}

func (r Repository) AdminGetDetail(ctx context.Context, id int64) (entity.Branch, error) {
	claims, err := r.CheckClaims(ctx, auth.RoleAdmin)
	if err != nil {
		return entity.Branch{}, err
	}

	var detail entity.Branch

	err = r.NewSelect().Model(&detail).Where("id = ? AND restaurant_id = ? AND deleted_at IS NULL", id, claims.RestaurantID).Scan(ctx)
	if err != nil {
		return entity.Branch{}, err
	}

	if detail.Photos != nil {
		var photos pq.StringArray
		for _, v := range *detail.Photos {
			baseLink := hashing.GenerateHash(r.ServerBaseUrl, v)
			photos = append(photos, baseLink)
		}
		detail.Photos = &photos
	}

	return detail, nil
}

func (r Repository) AdminCreate(ctx context.Context, request branch.AdminCreateRequest) (branch.AdminCreateResponse, error) {
	claims, err := r.CheckClaims(ctx, auth.RoleAdmin)
	if err != nil {
		return branch.AdminCreateResponse{}, err
	}

	err = r.ValidateStruct(&request, "Location", "WorkTime", "Name", "CategoryID")
	if err != nil {
		return branch.AdminCreateResponse{}, err
	}

	if request.WorkTime != nil {
		workTime := []string{"Monday", "Tuesday", "Wednesday", "Thursday", "Friday", "Saturday", "Sunday"}
		//workTime := []string{"Monday", "Tuesday"}
		st := 0
		for k, _ := range request.WorkTime {
			if slices.Contains(workTime, k) {
				st += 1
			}
		}
		if st != len(workTime) {
			return branch.AdminCreateResponse{}, errors.New("invalid WorkTime format")
		}
	}

	response := branch.AdminCreateResponse{
		Name:         request.Name,
		CategoryID:   request.CategoryID,
		Location:     request.Location,
		Photos:       request.PhotosLink,
		Status:       request.Status,
		WorkTime:     request.WorkTime,
		CreatedAt:    time.Now(),
		CreatedBy:    claims.UserId,
		RestaurantID: *claims.RestaurantID,
	}

	_, err = r.NewInsert().Model(&response).Exec(ctx)
	if err != nil {
		return branch.AdminCreateResponse{}, web.NewRequestError(errors.Wrap(err, "creating user"), http.StatusBadRequest)
	}

	return response, nil
}

func (r Repository) AdminUpdateAll(ctx context.Context, request branch.AdminUpdateRequest) error {
	claims, err := r.CheckClaims(ctx, auth.RoleAdmin)
	if err != nil {
		return err
	}

	if err = r.ValidateStruct(&request, "ID", "Location", "Photos",
		"WorkTime", "Status", "Name", "CategoryID"); err != nil {
		return err
	}

	if request.WorkTime != nil {
		workTime := []string{"Monday", "Tuesday", "Wednesday", "Thursday", "Friday", "Saturday", "Sunday"}
		//workTime := []string{"Monday", "Tuesday"}
		st := 0
		for k, _ := range request.WorkTime {
			if slices.Contains(workTime, k) {
				st += 1
			}
		}
		if st != len(workTime) {
			return errors.New("invalid WorkTime format")
		}
	}

	q := r.NewUpdate().Table("branches").Where("deleted_at IS NULL AND id = ? AND restaurant_id = ?",
		request.ID, claims.RestaurantID)

	q.Set("name = ?", request.Name)
	q.Set("category_id = ?", request.CategoryID)
	q.Set("location = ?", request.Location)
	q.Set("photos = array_cat(photos, ?)", request.PhotosLink)
	q.Set("status = ?", request.Status)
	q.Set("work_time = ?", request.WorkTime)
	q.Set("updated_at = ?", time.Now())
	q.Set("updated_by = ?", claims.UserId)

	_, err = q.Exec(ctx)
	if err != nil {
		return web.NewRequestError(errors.Wrap(err, "updating branch"), http.StatusBadRequest)
	}

	return nil
}

func (r Repository) AdminUpdateColumns(ctx context.Context, request branch.AdminUpdateRequest) error {
	claims, err := r.CheckClaims(ctx, auth.RoleAdmin)
	if err != nil {
		return err
	}

	if err = r.ValidateStruct(&request, "ID"); err != nil {
		return err
	}

	if request.WorkTime != nil {
		workTime := []string{"Monday", "Tuesday", "Wednesday", "Thursday", "Friday", "Saturday", "Sunday"}
		//workTime := []string{"Monday", "Tuesday"}
		st := 0
		for k, _ := range request.WorkTime {
			if slices.Contains(workTime, k) {
				st += 1
			}
		}
		if st != len(workTime) {
			return errors.New("invalid WorkTime format")
		}
	}

	q := r.NewUpdate().Table("branches").Where("deleted_at IS NULL AND id = ? AND restaurant_id = ?",
		request.ID, claims.RestaurantID)

	if request.Name != nil {
		q.Set("name = ?", request.Name)
	}
	if request.CategoryID != nil {
		q.Set("category_id = ?", request.CategoryID)
	}
	if request.Location != nil {
		q.Set("location = ?", request.Location)
	}
	if request.PhotosLink != nil {
		q.Set("photos = array_cat(photos, ?)", request.PhotosLink)
	}
	if request.Status != nil {
		q.Set("status = ?", request.Status)
	}
	if request.WorkTime != nil {
		q.Set("work_time = ?", request.WorkTime)
	}

	q.Set("updated_at = ?", time.Now())
	q.Set("updated_by = ?", claims.UserId)

	_, err = q.Exec(ctx)
	if err != nil {
		return web.NewRequestError(errors.Wrap(err, "updating branch"), http.StatusBadRequest)
	}

	return nil
}

func (r Repository) AdminDelete(ctx context.Context, id int64) error {
	return r.DeleteRow(ctx, "branches", id, auth.RoleAdmin)
}

func (r Repository) AdminDeleteImage(ctx context.Context, request branch.AdminDeleteImageRequest) error {
	_, err := r.CheckClaims(ctx, auth.RoleAdmin)
	if err != nil {
		return err
	}

	if err := r.ValidateStruct(&request,
		"ID",
	); err != nil {
		return err
	}

	var images []string
	err = r.NewSelect().Table("branches").
		Column("photos").
		Where("deleted_at IS NULL AND id = ?", request.ID).
		Scan(ctx,
			pgdialect.Array(&images))

	if err != nil {
		return web.NewRequestError(errors.Wrap(err, "scan branches"), http.StatusBadRequest)
	}

	if request.ImageIndex != nil {
		err = r.DeleteImageIndex(ctx, request.ID, "branches", "photos", images, *request.ImageIndex, auth.RoleAdmin)
		if err != nil {
			return web.NewRequestError(err, http.StatusBadRequest)
		}
	}

	return nil
}

// @client

//func (r Repository) ClientGetList(ctx context.Context, filter branch.Filter) ([]branch.ClientGetList, int, error) {
//	claims, err := r.CheckClaims(ctx, auth.RoleClient)
//	if err != nil {
//		claims.UserId = 0
//	}
//
//	table := "branches"
//	whereQuery := fmt.Sprintf(`WHERE b.deleted_at IS NULL`)
//
//	if filter.IsLiked != nil {
//		if *filter.IsLiked {
//			whereQuery += fmt.Sprintf(` and b.id in (SELECT branch_id FROM branch_likes WHERE user_id = %d)`, claims.UserId)
//		}
//	}
//
//	countWhereQuery := whereQuery
//
//	todayInt := time.Now().Weekday()
//	var today string
//	switch todayInt {
//	case 1:
//		today = "Monday"
//	case 2:
//		today = "Tuesday"
//	case 3:
//		today = "Wednesday"
//	case 4:
//		today = "Thursday"
//	case 5:
//		today = "Friday"
//	case 6:
//		today = "Saturday"
//	case 7:
//		today = "Sunday"
//	}
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
//	whereQuery += fmt.Sprintf(" ORDER BY b.rate desc %s %s", limitQuery, offsetQuery)
//
//	query := fmt.Sprintf(`
//					SELECT
//					    b.id,
//					    b.status,
//					    b.location,
//					    b.photos,
//					    b.work_time->>'%s' as work_time_today,
//					    b.name,
//					    b.category_id,
//					    rc.name as category_name,
//					    br.point,
//					    b.rate,
//					    CASE WHEN (SELECT id FROM branch_likes WHERE branch_id = b.id AND user_id = %d) IS NOT NULL THEN true ELSE false END AS is_liked
//					FROM
//					    branches as b
//					LEFT OUTER JOIN restaurant_category as rc ON rc.id = b.category_id
//					LEFT OUTER JOIN branch_reviews br on b.id = br.branch_id
//					%s`, today, claims.UserId, whereQuery)
//
//	list := make([]branch.ClientGetList, 0)
//
//	rows, err := r.QueryContext(ctx, query)
//	if err != nil {
//		return nil, 0, web.NewRequestError(errors.Wrap(err, "select branches"), http.StatusInternalServerError)
//	}
//
//	err = r.ScanRows(ctx, rows, &list)
//	if err != nil {
//		return nil, 0, web.NewRequestError(errors.Wrap(err, "scanning branches"), http.StatusBadRequest)
//	}
//
//	countQuery := fmt.Sprintf(`
//		SELECT
//			count(id)
//		FROM
//		    %s as b
//		%s
//	`, table, countWhereQuery)
//	countRows, err := r.QueryContext(ctx, countQuery)
//	if err == sql.ErrNoRows {
//		return nil, 0, web.NewRequestError(postgres.ErrNotFound, http.StatusNotFound)
//	}
//	if err != nil {
//		return nil, 0, web.NewRequestError(errors.Wrap(err, "selecting branch"), http.StatusBadRequest)
//	}
//
//	count := 0
//
//	for countRows.Next() {
//		if err = countRows.Scan(&count); err != nil {
//			return nil, 0, web.NewRequestError(errors.Wrap(err, "scanning branch count"), http.StatusBadRequest)
//		}
//	}
//
//	for k, v := range list {
//		var basePhotos pq.StringArray
//		for _, v1 := range *v.Photos {
//			baseLink := hashing.GenerateHash(r.ServerBaseUrl, v1)
//			basePhotos = append(basePhotos, baseLink)
//		}
//		list[k].Photos = &basePhotos
//	}
//
//	return list, count, nil
//}
//
//func (r Repository) ClientGetMapList(ctx context.Context, filter branch.Filter) ([]branch.ClientGetMapList, int, error) {
//	//_, err := r.CheckClaims(ctx, auth.RoleClient)
//	//if err != nil {
//	//	return nil, 0, err
//	//}
//
//	table := "branches"
//	whereQuery := fmt.Sprintf(`WHERE b.deleted_at IS NULL`)
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
//	whereQuery += fmt.Sprintf(" ORDER BY b.rate desc %s %s", limitQuery, offsetQuery)
//
//	query := fmt.Sprintf(`
//					SELECT
//					    b.id,
//					    b.status,
//					    b.location->>'lat' as lat,
//					    b.location->>'lon' as lon,
//					    b.name,
//					    b.category_id,
//					    r.logo as logo
//					FROM
//					    branches as b
//					LEFT OUTER JOIN restaurant_category as rc ON rc.id = b.category_id
//					LEFT OUTER JOIN branch_reviews br on b.id = br.branch_id
//					LEFT OUTER JOIN restaurants r ON r.id = b.restaurant_id
//					%s`, whereQuery)
//
//	list := make([]branch.ClientGetMapList, 0)
//
//	rows, err := r.QueryContext(ctx, query)
//	if err != nil {
//		return nil, 0, web.NewRequestError(errors.Wrap(err, "select branches"), http.StatusInternalServerError)
//	}
//
//	err = r.ScanRows(ctx, rows, &list)
//	if err != nil {
//		return nil, 0, web.NewRequestError(errors.Wrap(err, "scanning branches"), http.StatusBadRequest)
//	}
//
//	countQuery := fmt.Sprintf(`
//		SELECT
//			count(id)
//		FROM
//		    %s as b
//		%s
//	`, table, countWhereQuery)
//	countRows, err := r.QueryContext(ctx, countQuery)
//	if err == sql.ErrNoRows {
//		return nil, 0, web.NewRequestError(postgres.ErrNotFound, http.StatusNotFound)
//	}
//	if err != nil {
//		return nil, 0, web.NewRequestError(errors.Wrap(err, "selecting branch"), http.StatusBadRequest)
//	}
//
//	count := 0
//
//	for countRows.Next() {
//		if err = countRows.Scan(&count); err != nil {
//			return nil, 0, web.NewRequestError(errors.Wrap(err, "scanning branch count"), http.StatusBadRequest)
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
//
//func (r Repository) ClientGetDetail(ctx context.Context, id int64) (branch.ClientGetDetail, error) {
//	claims, err := r.CheckClaims(ctx, auth.RoleClient)
//	if err != nil {
//		claims.UserId = 0
//	}
//
//	whereQuery := fmt.Sprintf(`WHERE b.deleted_at IS NULL AND b.id = '%d'`, id)
//
//	todayInt := time.Now().Weekday()
//	var today string
//	switch todayInt {
//	case 1:
//		today = "Monday"
//	case 2:
//		today = "Tuesday"
//	case 3:
//		today = "Wednesday"
//	case 4:
//		today = "Thursday"
//	case 5:
//		today = "Friday"
//	case 6:
//		today = "Saturday"
//	case 7:
//		today = "Sunday"
//	}
//
//	query := fmt.Sprintf(`
//					SELECT
//					    b.id,
//					    b.status,
//					    b.location,
//					    b.photos,
//					    b.work_time->>'%s' as work_time_today,
//					    b.name,
//					    b.rate,
//					    r.name as restaurant_name,
//					    r.logo as restaurant_logo
//					FROM
//					    branches as b
//					Left Outer Join restaurants as r on r.id = b.restaurant_id
//					%s`, today, whereQuery)
//
//	var detail branch.ClientGetDetail
//
//	location := make([]byte, 0)
//	err = r.QueryRowContext(ctx, query).Scan(
//		&detail.ID, &detail.Status, &location, &detail.Photos,
//		&detail.WorkTimeToday, &detail.Name, &detail.Rate,
//		&detail.RestaurantName, &detail.RestaurantLogo,
//	)
//	if err != nil {
//		return branch.ClientGetDetail{}, web.NewRequestError(errors.Wrap(err, "select branches"), http.StatusInternalServerError)
//	}
//
//	if location != nil {
//		err = json.Unmarshal(location, &detail.Location)
//		if err != nil {
//			return branch.ClientGetDetail{}, errors.Wrap(err, "location unmarshal")
//		}
//	}
//	if detail.Photos != nil && len(*detail.Photos) > 0 {
//		var photoLinks pq.StringArray
//		for _, v := range *detail.Photos {
//			linkP := hashing.GenerateHash(r.ServerBaseUrl, v)
//			photoLinks = append(photoLinks, linkP)
//		}
//		detail.Photos = &photoLinks
//	}
//	if detail.RestaurantLogo != nil {
//		logoR := hashing.GenerateHash(r.ServerBaseUrl, *detail.RestaurantLogo)
//		detail.RestaurantLogo = &logoR
//	}
//
//	// ----------------------------food-category-----------------------------------------------------
//
//	//categoryQuery := fmt.Sprintf(`
//	//				SELECT
//	//				    fc.id,
//	//				    fc.name
//	//				FROM
//	//				    food_category as fc
//	//				LEFT OUTER JOIN foods as f ON f.category_id = fc.id
//	//				WHERE fc.deleted_at IS NULL AND f.branch_id = '%d'
//	//				GROUP BY fc.id, fc.name`, detail.ID)
//	//
//	//categoryDetail := make([]foodCategory.ClientGetList, 0)
//	//rows, err := r.QueryContext(ctx, categoryQuery)
//	//if err != nil {
//	//	return branch.ClientGetDetail{}, web.NewRequestError(errors.Wrap(err, "select food_category"), http.StatusInternalServerError)
//	//}
//	//
//	//err = r.ScanRows(ctx, rows, &categoryDetail)
//	//if err != nil {
//	//	return branch.ClientGetDetail{}, web.NewRequestError(errors.Wrap(err, "scanning food_category"), http.StatusBadRequest)
//	//}
//	//
//	//if categoryDetail != nil && len(categoryDetail) > 0 {
//	//	detail.Category = categoryDetail
//	//}
//
//	// --------------------------order----------------------------------------------------------------
//
//	orderQuery := fmt.Sprintf(`
//					SELECT
//					    o.id,
//					    o.number as order_number,
//					    o.status,
//					    o.table_id,
//					    t.number as table_number,
//                        t.branch_id as branch_id,
//                        (
//                        	select
//                        	    sum(m.new_price*om.count)
//                        	from menus m
//                        	join order_menu om on m.id = om.menu_id
//                        	where om.order_id=o.id
//                        	  and om.deleted_at isnull
//                        ) as price,
//					    CASE WHEN accepted_at IS NULL THEN false ELSE true END AS accept
//					FROM
//					    orders AS o
//					LEFT OUTER JOIN tables AS t ON t.id = o.table_id
//					WHERE o.user_id='%d' AND o.status='NEW'`, claims.UserId)
//
//	orderDetail := make([]order2.ClientGetDetail, 0)
//	orderRows, err := r.QueryContext(ctx, orderQuery)
//	if err != nil {
//		return branch.ClientGetDetail{}, web.NewRequestError(errors.Wrap(err, "select order"), http.StatusInternalServerError)
//	}
//
//	err = r.ScanRows(ctx, orderRows, &orderDetail)
//	if err != nil {
//		return branch.ClientGetDetail{}, web.NewRequestError(errors.Wrap(err, "scanning order"), http.StatusBadRequest)
//	}
//	if orderDetail != nil && len(orderDetail) > 0 {
//		detail.Orders = make([]order2.ClientGetDetail, 0)
//		detail.NewOrders = make([]order2.ClientGetDetail, 0)
//		for _, v := range orderDetail {
//			if *v.BranchID == id {
//				if v.Accept {
//					detail.Orders = append(detail.Orders, v)
//				} else {
//					detail.NewOrders = append(detail.NewOrders, v)
//				}
//			}
//		}
//		canOrder := false
//		detail.CanOrder = &canOrder
//		if len(detail.Orders) == 0 {
//			detail.Orders = nil
//		}
//	} else {
//		detail.Orders = nil
//		canOrder := true
//		detail.CanOrder = &canOrder
//	}
//
//	// --------------------------end_of_process--------------------------------------------------------
//
//	if detail.Name != nil && detail.RestaurantName != nil {
//		if strings.Compare(*detail.Name, *detail.RestaurantName) == 0 {
//			detail.Name = nil
//		}
//	}
//
//	return detail, nil
//}
//
//func (r Repository) ClientNearlyBranchGetList(ctx context.Context, filter branch.Filter) ([]branch.ClientGetList, int, error) {
//	claims, err := r.CheckClaims(ctx, auth.RoleClient)
//	if err != nil {
//		claims.UserId = 0
//	}
//
//	table := "branches"
//	whereQuery := fmt.Sprintf(`WHERE b.deleted_at IS NULL`)
//	countWhereQuery := whereQuery
//
//	todayInt := time.Now().Weekday()
//	var today string
//	switch todayInt {
//	case 1:
//		today = "Monday"
//	case 2:
//		today = "Tuesday"
//	case 3:
//		today = "Wednesday"
//	case 4:
//		today = "Thursday"
//	case 5:
//		today = "Friday"
//	case 6:
//		today = "Saturday"
//	case 7:
//		today = "Sunday"
//	}
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
//	whereQuery += fmt.Sprintf(" ORDER BY distance %s %s", limitQuery, offsetQuery)
//
//	query := fmt.Sprintf(`
//					SELECT
//					    b.id,
//					    b.status,
//					    b.location,
//					    b.photos,
//					    b.work_time->>'%s' as work_time_today,
//					    b.name,
//					    b.category_id,
//					    rc.name as category_name,
//					    b.rate,
//					    CASE WHEN (SELECT id FROM branch_likes WHERE branch_id = b.id AND user_id = %d) IS NOT NULL THEN true ELSE false END AS is_liked,
//					    ST_DistanceSphere(
//						   ST_SetSRID(ST_MakePoint(CAST(b.location->>'lon' AS float), CAST(b.location->>'lat' AS float)), 4326),
//						   ST_SetSRID(ST_MakePoint(%v, %v), 4326)
//					    ) as distance
//					FROM
//					    branches as b
//					LEFT OUTER JOIN restaurant_category as rc ON rc.id = b.category_id
//					%s`, today, claims.UserId, *filter.Lon, *filter.Lat, whereQuery)
//
//	list := make([]branch.ClientGetList, 0)
//
//	rows, err := r.QueryContext(ctx, query)
//	if err != nil {
//		return nil, 0, web.NewRequestError(errors.Wrap(err, "select branches"), http.StatusInternalServerError)
//	}
//
//	err = r.ScanRows(ctx, rows, &list)
//	if err != nil {
//		return nil, 0, web.NewRequestError(errors.Wrap(err, "scanning branches"), http.StatusBadRequest)
//	}
//
//	countQuery := fmt.Sprintf(`
//		SELECT
//			count(id)
//		FROM
//		    %s as b
//		%s
//	`, table, countWhereQuery)
//	countRows, err := r.QueryContext(ctx, countQuery)
//	if err == sql.ErrNoRows {
//		return nil, 0, web.NewRequestError(postgres.ErrNotFound, http.StatusNotFound)
//	}
//	if err != nil {
//		return nil, 0, web.NewRequestError(errors.Wrap(err, "selecting branch"), http.StatusBadRequest)
//	}
//
//	count := 0
//	for countRows.Next() {
//		if err = countRows.Scan(&count); err != nil {
//			return nil, 0, web.NewRequestError(errors.Wrap(err, "scanning branch count"), http.StatusBadRequest)
//		}
//	}
//
//	for k, v := range list {
//		var basePhotos pq.StringArray
//		for _, v1 := range *v.Photos {
//			baseLink := hashing.GenerateHash(r.ServerBaseUrl, v1)
//			basePhotos = append(basePhotos, baseLink)
//		}
//		list[k].Photos = &basePhotos
//	}
//
//	return list, count, nil
//}
//
//func (r Repository) ClientUpdateColumns(ctx context.Context, request branch.ClientUpdateRequest) error {
//	claims, err := r.CheckClaims(ctx, auth.RoleClient)
//	if err != nil {
//		return nil
//	}
//
//	if err = r.ValidateStruct(&request, "ID, IsLiked"); err != nil {
//		return err
//	}
//	if *request.IsLiked {
//		branchLikes := entity.BranchLikes{
//			BranchID: request.ID,
//			UserID:   claims.UserId,
//		}
//		_, err = r.NewInsert().Model(&branchLikes).Exec(ctx)
//		if err != nil {
//			return web.NewRequestError(errors.Wrap(err, "updating branch_likes"), http.StatusBadRequest)
//		}
//	} else {
//		_, err = r.NewDelete().Table("branch_likes").Where("branch_id = ? AND user_id = ?", request.ID, claims.UserId).Exec(ctx)
//		if err != nil {
//			return web.NewRequestError(errors.Wrap(err, "updating branch_likes"), http.StatusBadRequest)
//		}
//	}
//
//	return nil
//}

// @branch

func (r Repository) BranchGetDetail(ctx context.Context, id int64) (branch.BranchGetDetail, error) {
	_, err := r.CheckClaims(ctx, auth.RoleBranch)
	if err != nil {
		return branch.BranchGetDetail{}, err
	}

	whereQuery := fmt.Sprintf(`WHERE b.deleted_at IS NULL AND b.id = '%d'`, id)

	query := fmt.Sprintf(`
					SELECT 
					    b.id, 
					    r.logo as restaurant_logo
					FROM 
					    branches as b
					Left Outer Join restaurants as r on r.id = b.restaurant_id
					%s`, whereQuery)

	var detail branch.BranchGetDetail

	err = r.QueryRowContext(ctx, query).Scan(
		&detail.ID, &detail.Logo,
	)
	if err != nil {
		return branch.BranchGetDetail{}, web.NewRequestError(errors.Wrap(err, "select branches"), http.StatusInternalServerError)
	}

	if detail.Logo != nil {
		logoR := hashing.GenerateHash(r.ServerBaseUrl, *detail.Logo)
		detail.Logo = &logoR
	}

	return detail, nil
}

// @token

func (r Repository) BranchGetToken(ctx context.Context) (branch.BranchGetToken, error) {
	claims, err := r.CheckClaims(ctx, auth.RoleBranch)
	if err != nil {
		return branch.BranchGetToken{}, err
	}

	whereQuery := fmt.Sprintf(`WHERE deleted_at IS NULL AND id = '%d'`, *claims.BranchID)

	query := fmt.Sprintf(`
					SELECT 
					    token,
					    TO_CHAR(token_expired_at,'DD.MM.YYYY HH24:MI'),
					    CASE WHEN 
					        	CASE WHEN token_expired_at is not null THEN now()>token_expired_at ELSE true END 
					        THEN TRUE 
					        ELSE FALSE 
					    END
					FROM 
					    branches
	%s`, whereQuery)
	response := branch.BranchGetToken{}
	var expired bool
	err = r.QueryRowContext(ctx, query).Scan(response.Token, response.TokenExpiredAt, &expired)
	if err != nil {
		return branch.BranchGetToken{}, web.NewRequestError(errors.Wrap(err, "scan branch"), http.StatusNotImplemented)
	}

	if response.Token != nil {
		if !expired {
			return response, nil
		} else {
			token := uuid.NewString()
			_, err = r.NewUpdate().Table("branches").Where("id = ?", claims.BranchID).Set("token_expired_at = ?", time.Now().Add(30*24*time.Hour)).Set("token = ?", token).Exec(ctx)
			if err != nil {
				return branch.BranchGetToken{}, web.NewRequestError(errors.Wrap(err, "scan branch"), http.StatusNotImplemented)
			}
			response.Token = &token
			expiredAt := time.Now().Add(30 * 24 * time.Hour).Format("02.01.2006 15:04")
			response.TokenExpiredAt = &expiredAt
			return response, nil
		}
	}

	token := uuid.NewString()
	_, err = r.NewUpdate().Table("branches").Where("id = ?", claims.BranchID).Set("token_expired_at = ?", time.Now().Add(30*24*time.Hour)).Set("token = ?", token).Exec(ctx)
	if err != nil {
		return branch.BranchGetToken{}, web.NewRequestError(errors.Wrap(err, "scan branch"), http.StatusNotImplemented)
	}
	response.Token = &token
	expiredAt := time.Now().Add(30 * 24 * time.Hour).Format("02.01.2006 15:04")
	response.TokenExpiredAt = &expiredAt
	return response, nil
}

func (r Repository) WsGetByToken(ctx context.Context, token string) (branch.WsGetByTokenResponse, error) {

	whereQuery := fmt.Sprintf(`WHERE deleted_at IS NULL AND token = '%s'`, token)

	query := fmt.Sprintf(`
					SELECT 
					    id,
					    token_expired_at
					FROM 
					    branches
	%s`, whereQuery)
	response := branch.WsGetByTokenResponse{}
	err := r.QueryRowContext(ctx, query).Scan(&response.ID, &response.TokenExpiredAt)
	if err != nil {
		return branch.WsGetByTokenResponse{}, web.NewRequestError(errors.Wrap(err, "scan branch"), http.StatusNotImplemented)
	}

	return response, nil
}

func (r Repository) WsUpdateTokenExpiredAt(ctx context.Context, id int64) (string, error) {
	token := uuid.NewString()
	_, err := r.NewUpdate().Table("branches").Where("id = ?", id).Set("token_expired_at = ?", time.Now().Add(30*24*time.Hour)).Set("token = ?", token).Exec(ctx)
	if err != nil {
		return "", web.NewRequestError(errors.Wrap(err, "scan branch"), http.StatusNotImplemented)
	}
	return token, nil
}

func NewRepository(DB *postgresql.Database) *Repository {
	return &Repository{DB}
}
