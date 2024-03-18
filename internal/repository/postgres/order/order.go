package order

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/dariubs/percent"
	"github.com/pkg/errors"
	"github.com/restaurant/foundation/web"
	"github.com/restaurant/internal/auth"
	"github.com/restaurant/internal/pkg/repository/postgresql"
	"github.com/restaurant/internal/repository/postgres"
	"github.com/restaurant/internal/service/hashing"
	"github.com/restaurant/internal/service/order"
	waiter2 "github.com/restaurant/internal/service/waiter"
	"net/http"
	"time"
)

type Repository struct {
	*postgresql.Database
}

func NewRepository(db *postgresql.Database) *Repository {
	return &Repository{db}
}

// @client

func (r *Repository) MobileListOrder(ctx context.Context, filter order.Filter) ([]order.MobileGetList, int, error) {
	claims, err := r.CheckClaims(ctx, auth.RoleClient)
	if err != nil {
		return nil, 0, err
	}

	whereQuery := fmt.Sprintf(
		`WHERE o.deleted_at IS NULL AND o.user_id = %d`,
		claims.UserId,
	)
	countWhereQuery := whereQuery

	//countWhereQuery := fmt.Sprintf(
	//AND ord.status='PAID' AND o.status='NEW'
	//	`WHERE ord.deleted_at IS NULL AND ord.user_id = %d AND ord.status='PAID'`,
	//	claims.UserId,
	//)

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

	//o.id,
	//	o.status,
	//	TO_CHAR(o.created_at, 'DD.MM.YYYY HH24:MI') AS created_at,
	//	b.name as branch_name,
	//	m.new_price as price

	//	query := fmt.Sprintf(`
	//			SELECT
	//			    json_agg(
	//			    	json_build_object(
	//			    		'id', o.id,
	//			    		'status', o.status,
	//			    		'created_at', TO_CHAR(o.created_at, 'DD.MM.YYYY HH24:MI'),
	//			    		'branch_name', b.name,
	//			    		'price', m.new_price)) FILTER ( WHERE o.status = 'NEW' ) as active_order,
	//				json_agg(
	//					json_build_object(
	//						'id', o.id,
	//						'status', o.status,
	//						'created_at', TO_CHAR(o.created_at, 'DD.MM.YYYY HH24:MI'),
	//						'branch_name', b.name,
	//						'price', m.new_price)) FILTER ( WHERE o.status = 'PAID' ) as in_active_order
	//			FROM
	//				orders o
	//			LEFT OUTER JOIN tables as t ON t.id = o.table_id
	//			LEFT OUTER JOIN menus as m ON m.branch_id = t.branch_id
	//			LEFT OUTER JOIN	branches as b ON b.id = t.branch_id
	//
	//-- 				LEFT OUTER JOIN orders as ord ON ord.id = o.id
	//-- 				LEFT OUTER JOIN tables as tbl ON t.id = ord.table_id
	//-- 				LEFT OUTER JOIN menus as mnu ON m.branch_id = tbl.branch_id
	//-- 				LEFT OUTER JOIN	branches as brn ON b.id = tbl.branch_id
	//				%s
	//	`, whereQuery)

	query := fmt.Sprintf(`
				SELECT
					o.id,
					o.status,
					TO_CHAR(o.created_at, 'DD.MM.YYYY HH24:MI') AS created_at,
					b.name as branch_name,
					(select sum(m.new_price*om.count) from menus m join order_menu om on m.id = om.menu_id where om.order_id=o.id and om.deleted_at isnull ) as price,
					t.branch_id
				FROM
					orders o
				LEFT OUTER JOIN tables as t ON t.id = o.table_id
				LEFT OUTER JOIN	branches as b ON b.id = t.branch_id
				%s
				GROUP BY o.id, o.status, TO_CHAR(o.created_at, 'DD.MM.YYYY HH24:MI'), b.name, t.branch_id
				ORDER BY o.status, o.created_at desc %s %s`, whereQuery, limitQuery, offsetQuery)

	list := make([]order.MobileGetList, 0)

	rows, err := r.QueryContext(ctx, query)
	if err != nil {
		return nil, 0, web.NewRequestError(errors.Wrap(err, "select orders"), http.StatusInternalServerError)
	}

	err = r.ScanRows(ctx, rows, &list)
	if err != nil {
		return nil, 0, web.NewRequestError(errors.Wrap(err, "scanning orders"), http.StatusBadRequest)
	}

	countQuery := fmt.Sprintf(`
		SELECT
			count(o.id)
		FROM
		    orders o
		LEFT OUTER JOIN tables as t ON t.id = o.table_id
		LEFT OUTER JOIN	branches as b ON b.id = t.branch_id
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

func (r *Repository) MobileDetailOrder(ctx context.Context, id int64) (order.MobileGetDetail, error) {
	claims, err := r.CheckClaims(ctx, auth.RoleClient)
	if err != nil {
		return order.MobileGetDetail{}, err
	}

	whereQuery := fmt.Sprintf(`WHERE o.deleted_at IS NULL AND o.id = '%d' AND o.user_id = '%d'`, id, claims.UserId)

	query := fmt.Sprintf(`
					SELECT 
					    o.id,
					    o.status,
					    TO_CHAR(o.created_at, 'DD.MM.YYYY HH24:MI') as created_at,
					    b.name as branch_name,
					    t.number,
					    o.number
					FROM 
					    orders as o
					LEFT OUTER JOIN tables as t ON t.id = o.table_id
					LEFT OUTER JOIN	branches as b ON b.id = t.branch_id
					%s`, whereQuery)

	var detail order.MobileGetDetail

	var status string
	err = r.QueryRowContext(ctx, query).Scan(
		&detail.ID,
		&status,
		&detail.CreatedAt,
		&detail.BranchName,
		&detail.TableNumber,
		&detail.OrderNumber,
	)
	if err != nil {
		return order.MobileGetDetail{}, web.NewRequestError(errors.Wrap(err, "select branches"), http.StatusInternalServerError)
	}

	// ----------------------foods------------------------------------------------------------------

	foodQuery := fmt.Sprintf(`
					SELECT
					    m.id,
					    of.count,
					    f.name,
					    m.new_price * of.count as price
					FROM
					    order_menu as of
					LEFT OUTER JOIN menus as m ON m.id = of.menu_id
					LEFT OUTER JOIN foods as f ON f.id = m.food_id
					WHERE of.deleted_at IS NULL AND of.order_id = '%d'`, detail.ID)

	menuLists := make([]order.MenuList, 0)
	rows, err := r.QueryContext(ctx, foodQuery)
	if err != nil {
		return order.MobileGetDetail{}, web.NewRequestError(errors.Wrap(err, "select food_category"), http.StatusInternalServerError)
	}

	err = r.ScanRows(ctx, rows, &menuLists)
	if err != nil {
		return order.MobileGetDetail{}, web.NewRequestError(errors.Wrap(err, "scanning food"), http.StatusBadRequest)
	}

	percentage := 10
	if menuLists != nil && len(menuLists) > 0 {
		detail.Menus = menuLists

		var sum float32
		for k, v := range menuLists {
			if v.Price != nil {
				sum += float32(v.Count) * *v.Price
			} else {
				continue
			}
			menuLists[k].Status = &status
		}

		detail.Sum = &sum
		detail.ServicePercentage = &percentage

		service := float32(percent.PercentFloat(float64(percentage), float64(sum)))
		detail.Service = &service

		overAll := service + sum
		detail.OverAll = &overAll
		detail.Menus = menuLists
	}

	// ----------------------end_of_process--------------------------------------------------------

	return detail, nil
}

func (r *Repository) MobileCreateOrder(ctx context.Context, data order.MobileCreateRequest) (*order.MobileCreateResponse, error) {
	var (
		orderNumber *int
		response    order.MobileCreateResponse
		price       float64
		name        string
	)
	claims, err := r.CheckClaims(ctx, auth.RoleClient)
	if err != nil {
		return nil, web.NewRequestError(errors.Wrap(err, "extracting claims"), http.StatusInternalServerError)
	}

	var exists bool
	checkOrderExistence := fmt.Sprintf(`SELECT exists(select id from orders where user_id='%d' and status='NEW' and deleted_at isnull )`, claims.UserId)
	if err = r.QueryRowContext(ctx, checkOrderExistence).Scan(&exists); err != nil {
		return nil, err
	}
	if exists {
		err = errors.New("order already exists")
		return nil, err
	}

	today := time.Now().Format("02.01.2006")

	var branchID, restaurantID int64
	var tableStatus string
	tableNumberQuery := fmt.Sprintf(`SELECT t.branch_id,t.status,b.restaurant_id FROM tables t LEFT JOIN branches b ON b.id = t.branch_id WHERE t.id = '%d' AND t.deleted_at IS NULL `, data.TableID)
	if err = r.QueryRowContext(ctx, tableNumberQuery).Scan(&branchID, &tableStatus, &restaurantID); err != nil {
		return nil, web.NewRequestError(errors.Wrap(err, "selecting table number"), http.StatusInternalServerError)
	}

	//AND t.status='active'

	orderNumberQuery := fmt.Sprintf(`SELECT max(o.number)
											FROM orders o
											JOIN tables t ON t.id = '%d'
											JOIN branches b ON b.id = t.branch_id
											WHERE TO_CHAR(o.created_at, 'DD.MM.YYYY')='%s' AND o.deleted_at IS NULL `, data.TableID, today)
	if err = r.QueryRowContext(ctx, orderNumberQuery).Scan(&orderNumber); err != nil {
		if !errors.Is(err, sql.ErrNoRows) {
			return nil, web.NewRequestError(errors.Wrap(err, "selecting order number"), http.StatusInternalServerError)
		}
	}

	if orderNumber == nil {
		n := 1
		orderNumber = &n
	} else {
		n := *orderNumber + 1
		orderNumber = &n
	}

	tx, err := r.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}

	defer func() {
		if err != nil {
			_ = tx.Rollback()
			return
		}

		_ = tx.Commit()
	}()

	insertOrderQuery := fmt.Sprintf(`INSERT INTO orders (table_id, user_id, created_by, number, client_count) VALUES ('%d', '%d', '%d', '%d', '%d') RETURNING id, TO_CHAR(created_at, 'YYYY-MM-DD"T"HH24:MI:SS"Z"')`, data.TableID, claims.UserId, claims.UserId, *orderNumber, data.ClientCount)

	if err = tx.QueryRowContext(ctx, insertOrderQuery).Scan(&response.ID, &response.CreatedAt); err != nil {
		return nil, web.NewRequestError(errors.Wrap(err, "inserting order"), http.StatusInternalServerError)
	}

	var menus []order.ResponseMenu

	for _, v := range data.Menus {
		var menu order.ResponseMenu
		insertOrderFoodQuery := fmt.Sprintf(`INSERT INTO order_menu (count, order_id, menu_id, created_by) VALUES ('%d', '%d', '%d', '%d')`, v.Count, response.ID, v.ID, claims.UserId)
		if _, err = tx.ExecContext(ctx, insertOrderFoodQuery); err != nil {
			return nil, web.NewRequestError(errors.Wrap(err, "inserting order_menu"), http.StatusInternalServerError)
		}

		name, price, err = r.GetFoodPrice(ctx, v.ID, branchID)
		if err != nil {
			return nil, web.NewRequestError(errors.Wrap(err, "selecting price"), http.StatusInternalServerError)
		}

		response.Price += price * float64(v.Count)

		menu.Price = price * float64(v.Count)
		menu.ID = v.ID
		menu.Count = v.Count
		menu.Name = name

		menus = append(menus, menu)
	}
	response.BranchID = branchID
	response.RestaurantID = restaurantID
	response.UserID = claims.UserId
	response.Menus = menus

	return &response, nil
}

func (r *Repository) MobileUpdateOrder(ctx context.Context, data order.MobileUpdateRequest) error {
	claims, err := r.CheckClaims(ctx, auth.RoleClient)
	if err != nil {
		return web.NewRequestError(errors.Wrap(err, "extracting claims"), http.StatusInternalServerError)
	}

	// inserting foods...
	for _, v := range data.Menus {
		foodQuery := fmt.Sprintf(`INSERT INTO
											order_menu
												(order_id, menu_id, count)
											VALUES ('%d', '%d', '%d')`, data.Id, v.ID, v.Count)
		if _, err = r.ExecContext(ctx, foodQuery); err != nil {
			return err
		}
	}

	// updating orders...
	q := r.NewUpdate().Table("orders").Where("user_id=? and deleted_at isnull and status!='PAID'", claims.UserId)

	// fake update for cheating bun...
	q.Set("id=id")

	// actual value...
	if data.Status != nil && *data.Status == "CANCELLED" {
		q.Set("status=?", *data.Status)
	}

	_, err = q.Exec(ctx)
	if err != nil {
		return web.NewRequestError(errors.Wrap(err, "updating order"), http.StatusBadRequest)
	}

	return nil
}

func (r *Repository) MobileDeleteOrder(ctx context.Context, id int64) error {
	return r.DeleteRow(ctx, "orders", id, auth.RoleClient)
}

func (r *Repository) ClientReview(ctx context.Context, request order.ClientReviewRequest) (*order.ClientReviewResponse, error) {
	claims, err := r.CheckClaims(ctx, auth.RoleClient)
	if err != nil {
		return nil, err
	}

	err = r.ValidateStruct(&request, "Star", "OrderId")
	if err != nil {
		return nil, err
	}

	if request.Star > 5 || request.Star < 1 {
		err = errors.New("star must be between 1 and 5")
		return nil, web.NewRequestError(err, http.StatusBadRequest)
	}

	var score int
	switch request.Star {
	case 5:
		score = 5
	case 4:
		score = 3
	case 3:
		score = 0
	case 2:
		score = -3
	case 1:
		score = -5
	}

	response := order.ClientReviewResponse{
		Star:        request.Star,
		Score:       score,
		Description: request.Description,
		OrderId:     request.OrderId,
		CreatedAt:   time.Now(),
		CreatedBy:   claims.UserId,
	}

	_, err = r.NewInsert().Model(&response).Exec(ctx)
	if err != nil {
		return nil, web.NewRequestError(errors.Wrap(err, "applying review"), http.StatusInternalServerError)
	}

	return &response, nil
}

//@cashier

func (r *Repository) CashierListOrder(ctx context.Context, filter order.Filter) ([]order.CashierGetList, int, error) {
	claims, err := r.CheckClaims(ctx, auth.RoleCashier)
	if err != nil {
		return nil, 0, err
	}

	whereQuery := fmt.Sprintf(`WHERE o.deleted_at IS NULL AND t.branch_id = %d`, *claims.BranchID)
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
	orderQuery := " ORDER BY o.created_at desc"
	whereQuery += fmt.Sprintf("%s %s %s", orderQuery, limitQuery, offsetQuery)

	query := fmt.Sprintf(`
				SELECT 
					o.id,
					o.number,
					o.status,
					t.id AS table_id,
				    t.number  AS table_number,
					TO_CHAR(o.created_at, 'HH24:MI | DD.MM.YYYY') as created_at,
					(select sum(m.new_price*om.count) from menus m join order_menu om on m.id = om.menu_id where om.order_id=o.id and om.deleted_at isnull ) as price
				FROM 
					orders o
				LEFT OUTER JOIN tables as t ON t.id = o.table_id
				LEFT OUTER JOIN	branches as b ON b.id = t.branch_id
				%s
	`, whereQuery)

	list := make([]order.CashierGetList, 0)

	rows, err := r.QueryContext(ctx, query)
	if err != nil {
		return nil, 0, web.NewRequestError(errors.Wrap(err, "select orders"), http.StatusInternalServerError)
	}

	err = r.ScanRows(ctx, rows, &list)
	if err != nil {
		return nil, 0, web.NewRequestError(errors.Wrap(err, "scanning orders"), http.StatusBadRequest)
	}

	countQuery := fmt.Sprintf(`
		SELECT
			count(o.id)
		FROM
		    orders o
		LEFT OUTER JOIN tables as t ON t.id = o.table_id
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

func (r *Repository) CashierDetailOrder(ctx context.Context, id int64) (order.CashierGetDetail, error) {
	claims, err := r.CheckClaims(ctx, auth.RoleCashier)
	if err != nil {
		return order.CashierGetDetail{}, err
	}

	whereQuery := fmt.Sprintf(`WHERE o.deleted_at IS NULL AND o.id = '%d' AND t.branch_id = '%d'`, id, *claims.BranchID)

	query := fmt.Sprintf(`
					SELECT 
					    o.id,
						o.number,
						o.status,
						o.waiter_id,
						t.id,
				    	t.number,
						TO_CHAR(o.created_at, 'DD.MM.YYYY HH24:MI') as created_at,
						(select sum(m.new_price*om.count) from menus m join order_menu om on m.id = om.menu_id where om.order_id=o.id and om.deleted_at isnull ) as price
					FROM 
					    orders as o
					LEFT OUTER JOIN tables as t ON t.id = o.table_id
					%s`, whereQuery)

	var detail order.CashierGetDetail

	err = r.QueryRowContext(ctx, query).Scan(
		&detail.ID,
		&detail.Number,
		&detail.Status,
		&detail.WaiterID,
		&detail.TableID,
		&detail.TableNumber,
		&detail.CreatedAt,
		&detail.Price,
	)
	if errors.Is(err, sql.ErrNoRows) {
		return order.CashierGetDetail{}, web.NewRequestError(postgres.ErrNotFound, http.StatusBadRequest)
	}
	if err != nil {
		return order.CashierGetDetail{}, web.NewRequestError(errors.Wrap(err, "select branches"), http.StatusInternalServerError)
	}

	// ----------------------menus------------------------------------------------------------------

	foodQuery := fmt.Sprintf(`
					SELECT
					    m.id,
					    of.count,
					    f.name,
					    m.new_price * of.count as price,
					    f.photos
					FROM
					    order_menu as of
					LEFT OUTER JOIN menus as m ON m.id = of.menu_id
					LEFT OUTER JOIN foods as f ON f.id = m.food_id
					WHERE of.deleted_at IS NULL AND of.order_id = '%d'`, detail.ID)

	menuList := make([]order.MenuList, 0)
	rows, err := r.QueryContext(ctx, foodQuery)
	if err != nil {
		return order.CashierGetDetail{}, web.NewRequestError(errors.Wrap(err, "select food_category"), http.StatusInternalServerError)
	}

	err = r.ScanRows(ctx, rows, &menuList)
	if err != nil {
		return order.CashierGetDetail{}, web.NewRequestError(errors.Wrap(err, "scanning food"), http.StatusBadRequest)
	}

	if menuList != nil && len(menuList) > 0 {
		detail.Menus = menuList
	}

	if detail.Menus != nil {
		for _, menu := range detail.Menus {
			if menu.Photos != nil {
				for i, v := range menu.Photos {
					hashedValue := hashing.GenerateHash(r.ServerBaseUrl, v)
					menu.Photos[i] = hashedValue
				}
			}
		}
	}

	// ----------------------waiter----------------------------------------------------------------
	if detail.WaiterID != nil {
		waiterQuery := fmt.Sprintf(`SELECT
   									u.id,
   									u.name,
   									u.photo
								 FROM users u
								 WHERE u.id = '%d' AND u.deleted_at ISNULL AND u.role='WAITER'`, *detail.WaiterID)

		var waiterDet waiter2.CashierGetDetail
		err = r.QueryRowContext(ctx, waiterQuery).Scan(
			&waiterDet.ID,
			&waiterDet.Name,
			&waiterDet.Avatar,
		)
		if waiterDet.Avatar != nil {
			baseLink := r.ServerBaseUrl + *waiterDet.Avatar
			waiterDet.Avatar = &baseLink
		}
		detail.Waiter = waiterDet

		if err != nil {
			return order.CashierGetDetail{}, web.NewRequestError(errors.Wrap(err, "select waiter"), http.StatusInternalServerError)
		}
	}

	// ----------------------end_of_process--------------------------------------------------------

	return detail, nil
}

func (r *Repository) CashierUpdateStatus(ctx context.Context, data order.CashierUpdateStatusRequest) error {
	claims, err := r.CheckClaims(ctx, auth.RoleCashier)
	if err != nil {
		return web.NewRequestError(errors.Wrap(err, "extracting claims"), http.StatusInternalServerError)
	}

	err = r.ValidateStruct(&data, "Status", "Id")
	if err != nil {
		return web.NewRequestError(errors.Wrap(err, "data error"), http.StatusBadRequest)
	}

	// updating orders...
	q := r.NewUpdate().Table("orders").Where("table_id = (SELECT t.id FROM branches as b LEFT JOIN tables as t ON t.branch_id = b.id WHERE b.id = ? AND t.id = orders.table_id) AND deleted_at isnull and status!='PAID'", claims.BranchID)

	// actual value...
	if data.Status != nil {
		q.Set("status=?", *data.Status)
	}

	_, err = q.Exec(ctx)
	if err != nil {
		return web.NewRequestError(errors.Wrap(err, "updating order"), http.StatusBadRequest)
	}

	return nil
}

// @admin

func (r *Repository) AdminGetListOrder(ctx context.Context, filter order.Filter) ([]order.AdminList, int, error) {
	_, err := r.CheckClaims(ctx, auth.RoleAdmin)
	if err != nil {
		return nil, 0, err
	}

	table := "orders"
	whereQuery := fmt.Sprintf(`WHERE o.deleted_at IS NULL`)
	countWhereQuery := whereQuery

	var limitQuery, offsetQuery string
	if filter.Limit != nil {
		limitQuery = fmt.Sprintf(" LIMIT '%d'", *filter.Limit)
	}
	if filter.Offset != nil {
		offsetQuery = fmt.Sprintf(" OFFSET '%d'", *filter.Offset)
	}

	whereQuery += fmt.Sprintf("%s %s", limitQuery, offsetQuery)

	query := fmt.Sprintf(`
					SELECT 
					    o.id, 
					    o.status,
					    o.table_id,
					    o.user_id,
					    o.number,
					    t.number as table_number,
					    t.status as table_status,
					    u.name as user_name,
					    u.phone as user_phone
					FROM 
					    orders as o
					LEFT OUTER JOIN tables as t ON t.id = o.table_id
					LEFT OUTER JOIN users as u on u.id = o.user_id
					%s`, whereQuery)

	list := make([]order.AdminList, 0)

	rows, err := r.QueryContext(ctx, query)
	if err != nil {
		return nil, 0, web.NewRequestError(errors.Wrap(err, "select branches"), http.StatusInternalServerError)
	}

	err = r.ScanRows(ctx, rows, &list)
	if err != nil {
		return nil, 0, web.NewRequestError(errors.Wrap(err, "scanning order"), http.StatusBadRequest)
	}

	countQuery := fmt.Sprintf(`
		SELECT
			count(id)
		FROM
		    %s as o
		%s
	`, table, countWhereQuery)
	countRows, err := r.QueryContext(ctx, countQuery)
	if err == sql.ErrNoRows {
		return nil, 0, web.NewRequestError(postgres.ErrNotFound, http.StatusNotFound)
	}
	if err != nil {
		return nil, 0, web.NewRequestError(errors.Wrap(err, "selecting order"), http.StatusBadRequest)
	}

	count := 0

	for countRows.Next() {
		if err = countRows.Scan(&count); err != nil {
			return nil, 0, web.NewRequestError(errors.Wrap(err, "scanning branch count"), http.StatusBadRequest)
		}
	}

	return list, count, nil
}

// @waiter

func (r *Repository) WaiterGetList(ctx context.Context, filter order.Filter) ([]order.WaiterGetListResponse, int, error) {
	claims, err := r.CheckClaims(ctx, auth.RoleWaiter)
	if err != nil {
		return nil, 0, err
	}

	where := fmt.Sprintf(` WHERE o.deleted_at isnull AND t.deleted_at isnull AND t.branch_id='%d'`, *claims.BranchID)

	if filter.Whose != nil {
		if *filter.Whose == "MY" {
			where += fmt.Sprintf(` AND o.waiter_id='%d'`, claims.UserId)
		} else if *filter.Whose == "OTHERS" {
			where += fmt.Sprintf(` AND o.waiter_id!='%d'`, claims.UserId)
		}
	}
	if filter.Archived != nil && *filter.Archived {
		where += fmt.Sprintf(` AND (o.status='PAID' or o.status='CANCELLED')`)
	} else {
		where += fmt.Sprintf(` AND o.status!='PAID' AND o.status!='CANCELLED'`) // NOTE: current_date have not added because, restaurants can work 24/7, and in some cases order could not be completed
	}
	if filter.Search != nil {
		where += fmt.Sprintf(` AND 
											(o.number::text ilike '%s'
										OR 
											t.number::text ilike '%s'
										OR 
											(SELECT EXISTS(SELECT 
												m.id
											FROM menus m
												join order_menu om 
													on m.id = om.menu_id
												join foods f
													on m.food_id = f.id
											WHERE f.deleted_at isnull and om.order_id=o.id and f.name ilike '%s' and m.status='active')))`, "%"+*filter.Search+"%", "%"+*filter.Search+"%", "%"+*filter.Search+"%")
	}
	var limit, offset string
	if filter.Limit != nil {
		limit = fmt.Sprintf(` LIMIT %d`, *filter.Limit)
	}
	if filter.Page != nil {
		page := (*filter.Page - 1) * (*filter.Limit)
		offset = fmt.Sprintf(` OFFSET %d`, page)
	}
	orderQuery := ` ORDER BY o.created_at desc `

	query := fmt.Sprintf(`SELECT 
										o.id, 
										o.number, 
										t.number, 
										o.status, 
										to_char(o.created_at, 'DD.MM.YYYY') as created_date,
										to_char(o.created_at, 'HH24:MI') as created_at,
										case when u.role='WAITER' then u.id end as waiter_id,
										case when u.role='WAITER' then u.name end as waiter_name,
										case when o.waiter_id is not null then true else false end as accepted
									 FROM orders o 
										 join tables t 
											 on o.table_id = t.id
										 join users u
											 on o.user_id = u.id
									 %s %s %s %s`, where, orderQuery, limit, offset)
	rows, err := r.QueryContext(ctx, query)
	if err != nil {
		return nil, 0, err
	}

	response := []order.WaiterGetListResponse{}
	for rows.Next() {
		var row order.WaiterGetListResponse
		if err = rows.Scan(&row.Id, &row.Number, &row.TableNumber, &row.Status, &row.CreatedDate, &row.CreatedAt, &row.WaiterId, &row.WaiterName, &row.Accepted); err != nil {
			return nil, 0, err
		}

		query = fmt.Sprintf(`SELECT 
    									m.id, 
    									om.count, 
    									f.name, 
    									m.new_price as price,
    									om.id as order_menu_id,
    									om.status as order_menu_status
									FROM menus m
									    join order_menu om 
									        on m.id = om.menu_id 
									    join foods f 
									        on m.food_id = f.id
									WHERE f.deleted_at isnull and om.order_id='%d' 
									ORDER BY om.status!='NEW', om.status!='SERVED'`, row.Id)

		mRows, mErr := r.QueryContext(ctx, query)
		if mErr != nil {
			return nil, 0, err
		}

		var menus []order.WaiterMenu
		if err = r.ScanRows(ctx, mRows, &menus); err != nil {
			return nil, 0, err
		}

		row.Menus = menus

		response = append(response, row)
	}

	var count int
	query = fmt.Sprintf(`SELECT 
    								count(o.id) 
								FROM orders o 
								    join tables t 
								         on o.table_id = t.id
								%s`, where)
	if err = r.QueryRowContext(ctx, query).Scan(&count); err != nil {
		return nil, 0, err
	}

	return response, count, nil
}

func (r *Repository) WaiterCreate(ctx context.Context, data order.WaiterCreateRequest) (*order.WaiterCreateResponse, error) {
	var (
		orderNumber *int
		response    order.WaiterCreateResponse
		price       float64
	)
	claims, err := r.CheckClaims(ctx, auth.RoleWaiter)
	if err != nil {
		return nil, web.NewRequestError(errors.Wrap(err, "extracting claims"), http.StatusInternalServerError)
	}

	today := time.Now().Format("02.01.2006")

	var branchID, restaurantID int64
	var tableStatus, userStatus string
	tableNumberQuery := fmt.Sprintf(`SELECT t.branch_id,t.status,b.restaurant_id, (SELECT status FROM users WHERE id = '%d') FROM tables t LEFT JOIN branches b ON b.id = t.branch_id WHERE t.id = '%d' AND t.deleted_at IS NULL`, claims.UserId, data.TableID)
	if err = r.QueryRowContext(ctx, tableNumberQuery).Scan(&branchID, &tableStatus, &restaurantID, &userStatus); err != nil {
		return nil, web.NewRequestError(errors.Wrap(err, "selecting table number"), http.StatusInternalServerError)
	}

	orderNumberQuery := fmt.Sprintf(`SELECT max(o.number)
											FROM orders o
											JOIN tables t ON t.id = '%d'
											JOIN branches b ON b.id = t.branch_id
											WHERE TO_CHAR(o.created_at, 'DD.MM.YYYY')='%s' AND o.deleted_at IS NULL`, data.TableID, today)
	if err = r.QueryRowContext(ctx, orderNumberQuery).Scan(&orderNumber); err != nil {
		if !errors.Is(err, sql.ErrNoRows) {
			return nil, web.NewRequestError(errors.Wrap(err, "selecting order number"), http.StatusInternalServerError)
		}
	}

	if orderNumber == nil {
		n := 1
		orderNumber = &n
	} else {
		n := *orderNumber + 1
		orderNumber = &n
	}

	tx, err := r.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}

	defer func() {
		if err != nil {
			_ = tx.Rollback()
			return
		}

		_ = tx.Commit()
	}()

	insertOrderQuery := fmt.Sprintf(`INSERT INTO orders (table_id, user_id, created_by, number, client_count, waiter_id) VALUES ('%d', '%d', '%d', '%d', '%d', '%d') RETURNING id, TO_CHAR(created_at, 'YYYY-MM-DD"T"HH24:MI:SS"Z"')`, data.TableID, claims.UserId, claims.UserId, *orderNumber, data.ClientCount, claims.UserId)

	if err = tx.QueryRowContext(ctx, insertOrderQuery).Scan(&response.ID, &response.CreatedAt); err != nil {
		return nil, web.NewRequestError(errors.Wrap(err, "inserting order"), http.StatusInternalServerError)
	}

	for _, v := range data.Menus {
		insertOrderFoodQuery := fmt.Sprintf(`INSERT INTO order_menu (count, order_id, menu_id, created_by) VALUES ('%d', '%d', '%d', '%d')`, v.Count, response.ID, v.ID, claims.UserId)
		if _, err = tx.ExecContext(ctx, insertOrderFoodQuery); err != nil {
			return nil, web.NewRequestError(errors.Wrap(err, "inserting order_menu"), http.StatusInternalServerError)
		}

		_, price, err = r.GetFoodPrice(ctx, v.ID, branchID)
		if err != nil {
			return nil, web.NewRequestError(errors.Wrap(err, "selecting price"), http.StatusInternalServerError)
		}

		response.Price += price * float64(v.Count)
	}
	response.BranchID = branchID
	response.RestaurantID = restaurantID
	response.UserID = claims.UserId
	response.ClientCount = data.ClientCount

	return &response, nil
}

func (r *Repository) WaiterGetDetail(ctx context.Context, id int64) (*order.WaiterGetDetailResponse, error) {
	claims, err := r.CheckClaims(ctx, auth.RoleWaiter)
	if err != nil {
		return nil, err
	}

	where := fmt.Sprintf(` WHERE o.deleted_at isnull AND t.deleted_at isnull AND t.branch_id='%d' AND u.status='active' AND o.id='%d'`, *claims.BranchID, id)

	var response order.WaiterGetDetailResponse
	query := fmt.Sprintf(`SELECT 
										o.id, 
										o.number, 
										t.number, 
										o.status, 
										to_char(o.created_at, 'DD.MM.YYYY') as created_date,
										to_char(o.created_at, 'HH24:MI') as created_at,
										case when u.role='WAITER' then u.id end as waiter_id,
										case when u.role='WAITER' then u.name end as waiter_name,
										o.client_count
									 FROM orders o 
										 join tables t 
											 on o.table_id = t.id
										 join users u 
											 on o.user_id = u.id
									 %s`, where)
	if err = r.QueryRowContext(ctx, query).Scan(&response.Id, &response.Number, &response.TableNumber, &response.Status, &response.CreatedDate, &response.CreatedAt, &response.WaiterId, &response.WaiterName, &response.ClientCount); err != nil {
		return nil, err
	}

	query = fmt.Sprintf(`SELECT 
    									m.id, 
    									om.count, 
    									f.name, 
    									m.new_price as price,
    									om.id as order_menu_id,
    									om.status as order_menu_status
									FROM menus m
									    join order_menu om 
									        on m.id = om.menu_id 
									    join foods f 
									        on m.food_id = f.id
									WHERE f.deleted_at isnull and om.order_id='%d'
									ORDER BY om.status!='NEW', om.status!='SERVED'`, id)

	mRows, mErr := r.QueryContext(ctx, query)
	if mErr != nil {
		return nil, err
	}

	var menus []order.WaiterMenu
	if err = r.ScanRows(ctx, mRows, &menus); err != nil {
		return nil, err
	}

	var total float64
	for _, v := range menus {
		total += v.Price * float64(v.Count)
	}
	response.Price = &total
	response.Menus = menus

	return &response, nil
}

func (r *Repository) WaiterUpdate(ctx context.Context, data order.WaiterUpdateRequest) error {
	claims, err := r.CheckClaims(ctx, auth.RoleWaiter)
	if err != nil {
		return web.NewRequestError(errors.Wrap(err, "extracting claims"), http.StatusInternalServerError)
	}

	for _, v := range data.Menus {
		checkQuery := fmt.Sprintf(`SELECT EXISTS(SELECT m.id FROM menus m WHERE m.deleted_at isnull AND m.branch_id='%d' AND m.id='%d')`, *claims.BranchID, v.ID)

		var exists bool
		if err = r.QueryRowContext(ctx, checkQuery).Scan(&exists); err != nil {
			return err
		}

		if !exists {
			err = errors.New("there is not rule to update an order with the credentials")
			return err
		}
	}

	// updating foods...
	for _, v := range data.Menus {
		foodQuery := fmt.Sprintf(`INSERT INTO
											order_menu
												(order_id, menu_id, count)
											VALUES ('%d', '%d', '%d')`, data.Id, v.ID, v.Count)
		if _, err = r.ExecContext(ctx, foodQuery); err != nil {
			return err
		}
	}

	return nil
}

func (r *Repository) WaiterUpdateStatus(ctx context.Context, id int64, status string) error {
	claims, err := r.CheckClaims(ctx, auth.RoleWaiter)
	if err != nil {
		return web.NewRequestError(errors.Wrap(err, "extracting claims"), http.StatusInternalServerError)
	}

	// updating orders...
	q := r.NewUpdate().Table("orders").Where("waiter_id=? and deleted_at isnull and status!='PAID' and id=?", claims.UserId, id)

	// fake update for cheating bun...
	q.Set("id=id")

	// actual value...
	if status == "CANCELLED" || status == "SERVED" {
		q.Set("status=?", status)
	}

	_, err = q.Exec(ctx)
	if err != nil {
		return web.NewRequestError(errors.Wrap(err, "updating order"), http.StatusBadRequest)
	}

	return nil
}

func (r *Repository) WaiterAccept(ctx context.Context, id int64) (*order.WaiterAcceptOrderResponse, error) {
	claims, err := r.CheckClaims(ctx, auth.RoleWaiter)
	if err != nil {
		return nil, web.NewRequestError(errors.Wrap(err, "extracting claims"), http.StatusInternalServerError)
	}

	// checking before update...
	exists, err := r.NewSelect().Table("orders").Where("id = ? and accepted_at isnull and waiter_id isnull", id).Exists(ctx)
	if err != nil {
		return nil, err
	}

	if !exists {
		err = errors.New("order already taken")
		return nil, web.NewRequestError(err, http.StatusGone)
	}

	// updating orders...
	q := r.NewUpdate().Table("orders").Where("deleted_at isnull and status!='PAID' and id=? and waiter_id isnull and accepted_at isnull", id)

	q.Set("accepted_at = now()")
	q.Set("waiter_id = ?", claims.UserId)

	_, err = q.Exec(ctx)
	if err != nil {
		return nil, web.NewRequestError(errors.Wrap(err, "updating order"), http.StatusBadRequest)
	}

	var response order.WaiterAcceptOrderResponse
	query := fmt.Sprintf(`SELECT 
    									o.number, 
    									o.table_id, 
    									t.number, 
    									o.user_id, 
    									u.id, 
    									u.name, 
    									u.photo, 
    									o.accepted_at 
								 FROM orders o 
								     JOIN tables t 
								         ON o.table_id = t.id 
								     JOIN users u 
								         ON o.waiter_id = u.id 
								 WHERE 
								     o.id = '%d' 
								   and 
								     u.role = 'WAITER' 
								   and 
								     o.deleted_at is null
								   and 
								     u.id='%d'`, id, claims.UserId)
	if err = r.QueryRowContext(ctx, query).Scan(&response.OrderNumber, &response.TableID, &response.TableNumber, &response.ClientID, &response.WaiterID, &response.WaiterName, &response.WaiterPhoto, &response.AcceptedAt); err != nil {
		return nil, err
	}

	if response.WaiterPhoto != nil {
		link := hashing.GenerateHash(r.ServerBaseUrl, *response.WaiterPhoto)
		response.WaiterPhoto = &link
	}
	return &response, nil
}

// ----------others---------------------------------------------------------------------------------------------------------

func (r *Repository) OrderExists(ctx context.Context, tableID int64, userID int64) (*order.MobileCreateResponse, error) {
	var (
		response *order.MobileCreateResponse
	)

	_, err := r.NewSelect().Table("orders").Column("id", "TO_CHAR(created_at, 'YYYY-MM-DD\"T\"HH24:MI:SS\"Z\"')").Where("status='NEW' AND table_id=? AND deleted_at ISNULL AND user_id=?", tableID, userID).Exec(ctx, &response)
	if err != nil {
		return nil, err
	}

	if response != nil {
		priceQuery := fmt.Sprintf(`SELECT m.new_price * of.count
										  FROM orders o
										      JOIN order_menu of ON o.id = of.order_id
										      JOIN tables t ON t.id = o.table_id
										      JOIN menus m ON m.id=of.menu_id AND m.branch_id= t.branch_id
										  WHERE o.id = '%d';`, response.ID)
		rows, err := r.QueryContext(ctx, priceQuery)
		if err != nil {
			return nil, err
		}

		for rows.Next() {
			var (
				price float64
			)

			if err = rows.Scan(&price); err != nil {
				return nil, err
			}

			response.Price += price
		}
	}

	return response, nil
}

func (r *Repository) GetFoodPrice(ctx context.Context, menuID, branchID int64) (name string, price float64, err error) {
	query := fmt.Sprintf(`SELECT m.new_price, f.name FROM menus m JOIN foods f ON m.food_id = f.id WHERE m.deleted_at isnull  AND m.id = '%d' AND m.branch_id = '%d'`, menuID, branchID)
	if err = r.QueryRowContext(ctx, query).Scan(&price, &name); err != nil {
		return
	}
	return
}

func (r *Repository) GetWsMessage(ctx context.Context, orderId int64) (order.GetWsMessageResponse, error) {
	response := order.GetWsMessageResponse{}
	priceQuery := fmt.Sprintf(`SELECT 
    										t.id,
    										t.number,
    										o.id,
    										o.number,
    										b.id,
    										b.restaurant_id,
    										TO_CHAR(o.created_at, 'DD.MM.YYYY HH24:MI') as created_at,
    										o.user_id,
    										(select sum(m.new_price*om.count) from menus m join order_menu om on m.id = om.menu_id where om.order_id=o.id and om.deleted_at isnull ) as price,
											w.id
										  FROM orders o
										      JOIN order_menu of ON o.id = of.order_id
										      JOIN tables t ON t.id = o.table_id
										      JOIN branches b ON t.branch_id = b.id
										      JOIN menus m ON m.id=of.menu_id AND m.branch_id= t.branch_id
										      LEFT JOIN users w ON w.id = o.waiter_id
										  WHERE o.id = '%d';`, orderId)
	err := r.QueryRowContext(ctx, priceQuery).Scan(&response.TableID, &response.TableNumber, &response.OrderID, &response.OrderNumber, &response.BranchId, &response.RestaurantId, &response.CreatedAt, &response.UserId, &response.Price, &response.WaiterID)
	if err != nil {
		return order.GetWsMessageResponse{}, web.NewRequestError(errors.Wrap(err, "get ws massage"), http.StatusInternalServerError)
	}

	return response, nil
}

func (r *Repository) CheckOrderIfAccepted(id int64) error {
	var exists bool

	query := fmt.Sprintf(`SELECT EXISTS (SELECT 
    									o.id
								 FROM orders o 
								     JOIN users w 
								         ON o.waiter_id = w.id 
								 WHERE o.id='%d' AND w.role='WAITER' AND o.accepted_at is not null AND o.waiter_id is not null)`, id)
	if err := r.QueryRowContext(context.TODO(), query).Scan(&exists); err != nil {
		return err
	}

	if exists {
		err := errors.New("order already accepted")
		return err
	}

	return nil
}

func (r *Repository) CancelOrder(id int64) error {
	_, err := r.NewUpdate().Table("orders").Where("id = ?", id).Set("status = 'CANCELLED'").Exec(context.TODO())
	if err != nil {
		return err
	}

	return nil
}

func (r *Repository) GetWsOrderMenus(ctx context.Context, orderId int64, menus []order.Menu) ([]order.GetWsOrderMenusResponse, int64, error) {
	whereQuery := fmt.Sprintf("WHERE o.id = '%d'", orderId)

	menusMap := make(map[int64]int)
	if menus != nil {
		where := " AND m.id in ("
		for k, v := range menus {
			menusMap[v.ID] = v.Count
			if len(menus)-1 != k {
				where += fmt.Sprintf(" %d,", v.ID)
			} else {
				where += fmt.Sprintf(" %d", v.ID)
			}
		}
		whereQuery += where + ")"
	}

	query := fmt.Sprintf(`SELECT 
    										o.number,
    										t.branch_id,
    										t.number,
    										p.ip,
    										f.name,
    										om.count,
    										u.name,
    										m.id
    									 FROM orders o
											  LEFT JOIN users u ON u.id = o.waiter_id
											  LEFT JOIN tables t ON t.id = o.table_id
											  LEFT JOIN order_menu om ON o.id = om.order_id
										      LEFT JOIN menus m ON m.id = om.menu_id
											  LEFT JOIN foods f ON f.id = m.food_id
											  LEFT JOIN printers p ON p.id = m.printer_id                  
										  %s`, whereQuery)
	rows, err := r.QueryContext(ctx, query)
	if err != nil {
		return nil, 0, web.NewRequestError(errors.Wrap(err, "get ws order menu"), http.StatusInternalServerError)
	}

	list := make([]order.GetWsOrderMenusResponse, 0)
	listMap := make(map[string]int)
	var branchID int64
	for rows.Next() {
		var (
			food        = order.WsFood{}
			orderNumber *int64
			tableNumber *int64
			ip          *string
			waiter      *string
			menuID      int64
		)
		err = rows.Scan(&orderNumber, &branchID, &tableNumber, &ip, &food.Name, &food.Count, &waiter, &menuID)
		if err != nil {
			return nil, 0, web.NewRequestError(errors.Wrap(err, "scan ws order menu"), http.StatusInternalServerError)
		}
		if menus != nil {
			if v, ok := menusMap[menuID]; ok {
				food.Count = &v
			}
		}

		if ip != nil {
			if k, ok := listMap[*ip]; ok {
				list[k].Foods = append(list[k].Foods, food)
			} else {
				listMap[*ip] = len(list)
				list = append(list, order.GetWsOrderMenusResponse{
					OrderNumber: orderNumber,
					Ip:          ip,
					Waiter:      waiter,
					TableNumber: tableNumber,
					Foods: []order.WsFood{
						food,
					},
				})
			}
		}
	}

	return list, branchID, nil
}

func (r *Repository) GetWsWaiter(waiterID int64) (order.GetWsWaiterResponse, error) {
	whereQuery := fmt.Sprintf("WHERE w.id = '%d'", waiterID)
	response := order.GetWsWaiterResponse{}

	query := fmt.Sprintf(`SELECT
								     w.id,
								     CASE WHEN oc.count IS NOT NULL THEN oc.count ELSE 0 END
								 FROM users AS w
								          LEFT JOIN (
								     SELECT
								         waiter_id,
								         count(id) as count
								     FROM orders
								     WHERE
								             status = 'NEW' AND
								         waiter_id IS NOT NULL AND
								         deleted_at IS NULL
								     GROUP BY waiter_id) as oc ON oc.waiter_id = w.id
    				               
										  %s`, whereQuery)
	err := r.QueryRowContext(context.Background(), query).Scan(&response.ID, &response.OrderCount)
	if err != nil {
		return order.GetWsWaiterResponse{}, web.NewRequestError(errors.Wrap(err, "get ws order menu"), http.StatusInternalServerError)
	}

	return response, nil
}

//(
//SELECT
//EXTRACT(WEEK FROM o.created_at) AS week_number,
//COALESCE(SUM(o.id), 0) AS order_count,
//COALESCE(SUM(EXTRACT(EPOCH FROM (wwt.periods->>'gone_time')::timestamp - (wwt.periods->>'came_time')::timestamp) / 3600), 0) AS weekly_work_hours
//FROM
//waiter_work_time wwt
//LEFT JOIN orders o ON wwt.waiter_id = o.user_id
//WHERE
//wwt.waiter_id = 1363
//AND EXTRACT(WEEK FROM o.created_at) = EXTRACT(WEEK FROM CURRENT_DATE)
//GROUP BY
//week_number;
//	)
