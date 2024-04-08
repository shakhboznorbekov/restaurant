package order_menu

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/lib/pq"
	"github.com/pkg/errors"
	"net/http"
	"restu-backend/foundation/web"
	"restu-backend/internal/auth"
	"restu-backend/internal/entity"
	"restu-backend/internal/pkg/repository/postgresql"
	"restu-backend/internal/repository/postgres"
	"restu-backend/internal/service/hashing"
	"restu-backend/internal/service/order_menu"
	"strings"
	"time"
)

type Repository struct {
	*postgresql.Database
}

// @client

func (r Repository) ClientGetList(ctx context.Context, filter order_menu.Filter) ([]order_menu.ClientGetList, int, error) {
	_, err := r.CheckClaims(ctx, auth.RoleClient)
	if err != nil {
		return nil, 0, err
	}

	whereQuery := fmt.Sprintf(`WHERE of.deleted_at IS NULL `)
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
			of.id,
			of.count,
			of.menu_id,
			of.order_id,
			m.name
		FROM
		    order_menu as of
		LEFT OUTER JOIN menus m on m.id = of.menu_id
		LEFT OUTER JOIN orders o on o.id = of.order_id
		%s
	`, whereQuery)

	list := make([]order_menu.ClientGetList, 0)

	rows, err := r.QueryContext(ctx, query)
	if err != nil {
		return nil, 0, web.NewRequestError(errors.Wrap(err, "select order_foods"), http.StatusInternalServerError)
	}

	err = r.ScanRows(ctx, rows, &list)
	if err != nil {
		return nil, 0, web.NewRequestError(errors.Wrap(err, "scanning order_foods"), http.StatusBadRequest)
	}

	countQuery := fmt.Sprintf(`
		SELECT
			count(of.id)
		FROM
		    order_menu as of
		%s
	`, countWhereQuery)
	countRows, err := r.QueryContext(ctx, countQuery)
	if err == sql.ErrNoRows {
		return nil, 0, web.NewRequestError(postgres.ErrNotFound, http.StatusNotFound)
	}
	if err != nil {
		return nil, 0, web.NewRequestError(errors.Wrap(err, "selecting order_food"), http.StatusBadRequest)
	}

	count := 0
	for countRows.Next() {
		if err = countRows.Scan(&count); err != nil {
			return nil, 0, web.NewRequestError(errors.Wrap(err, "scanning order_food count"), http.StatusBadRequest)
		}
	}

	return list, count, nil
}

func (r Repository) ClientGetDetail(ctx context.Context, id int64) (entity.OrderMenu, error) {
	var detail entity.OrderMenu

	err := r.NewSelect().Model(&detail).Where("id = ? AND deleted_at IS NULL", id).Scan(ctx)
	if err != nil {
		return entity.OrderMenu{}, err
	}
	return detail, nil
}

func (r Repository) ClientCreate(ctx context.Context, request order_menu.ClientCreateRequest) (order_menu.ClientCreateResponse, error) {
	claims, err := r.CheckClaims(ctx, auth.RoleClient)
	if err != nil {
		return order_menu.ClientCreateResponse{}, err
	}

	err = r.ValidateStruct(&request, "Count", "MenuID", "OrderID")
	if err != nil {
		return order_menu.ClientCreateResponse{}, err
	}

	response := order_menu.ClientCreateResponse{
		Count:     request.Count,
		MenuID:    request.MenuID,
		OrderID:   request.OrderID,
		CreatedAt: time.Now(),
		CreatedBy: claims.UserId,
	}

	_, err = r.NewInsert().Model(&response).Exec(ctx)
	if err != nil {
		return order_menu.ClientCreateResponse{}, web.NewRequestError(errors.Wrap(err, "creating order_food"), http.StatusBadRequest)
	}

	return response, nil
}

func (r Repository) ClientUpdateAll(ctx context.Context, request order_menu.ClientUpdateRequest) error {
	claims, err := r.CheckClaims(ctx, auth.RoleClient)
	if err != nil {
		return err
	}

	if err = r.ValidateStruct(&request, "ID", "Count", "MenuID", "OrderID"); err != nil {
		return err
	}

	q := r.NewUpdate().Table("order_menu").Where("deleted_at IS NULL AND id = ?", request.ID)

	q.Set("count = ?", request.Count)
	q.Set("menu_id = ?", request.MenuID)
	q.Set("order_id = ?", request.OrderID)
	q.Set("updated_at = ?", time.Now())
	q.Set("updated_by = ?", claims.UserId)

	_, err = q.Exec(ctx)
	if err != nil {
		return web.NewRequestError(errors.Wrap(err, "updating order_menu"), http.StatusBadRequest)
	}

	return nil
}

func (r Repository) ClientUpdateColumns(ctx context.Context, request order_menu.ClientUpdateRequest) error {
	claims, err := r.CheckClaims(ctx, auth.RoleClient)
	if err != nil {
		return err
	}

	if err = r.ValidateStruct(&request, "ID"); err != nil {
		return err
	}

	q := r.NewUpdate().Table("order_menu").Where("deleted_at IS NULL AND id = ?", request.ID)

	if request.Count != nil {
		q.Set("count = ?", request.Count)
	}
	if request.MenuID != nil {
		q.Set("menu_id = ?", request.MenuID)
	}
	if request.OrderID != nil {
		q.Set("order_id = ?", request.OrderID)
	}

	q.Set("updated_at = ?", time.Now())
	q.Set("updated_by = ?", claims.UserId)

	_, err = q.Exec(ctx)
	if err != nil {
		return web.NewRequestError(errors.Wrap(err, "updating order_food"), http.StatusBadRequest)
	}

	return nil
}

func (r Repository) ClientDelete(ctx context.Context, id int64) error {
	return r.DeleteRow(ctx, "order_food", id, auth.RoleClient)
}

func (r Repository) ClientGetOftenList(ctx context.Context, branchID int) ([]order_menu.ClientGetOftenList, error) {
	//_, err := r.CheckClaims(ctx, auth.RoleClient)
	//if err != nil {
	//	return nil, err
	//}

	query := fmt.Sprintf(`
		WITH often_menus AS (
			SELECT orm.menu_id as foo, count(menu_id) AS count
			FROM order_menu as orm
					 LEFT JOIN orders as o ON orm.order_id = o.id
					 LEFT JOIN tables t on o.table_id = t.id
			WHERE t.branch_id='%d'
			GROUP BY orm.menu_id
			ORDER BY count(orm.menu_id) desc )
		SELECT
			m.id AS id,
			m.name AS name,
			m.photos AS photos,
			m.new_price AS price
		FROM menus AS m
				 LEFT JOIN often_menus ON m.id = often_menus.foo
		WHERE often_menus.foo IS NOT NULL;
	`, branchID)

	list := make([]order_menu.ClientGetOftenList, 0)

	rows, err := r.QueryContext(ctx, query)
	if err != nil {
		return nil, web.NewRequestError(errors.Wrap(err, "select order_foods"), http.StatusInternalServerError)
	}

	err = r.ScanRows(ctx, rows, &list)
	if err != nil {
		return nil, web.NewRequestError(errors.Wrap(err, "scanning order_foods"), http.StatusBadRequest)
	}

	for k, v := range list {
		var photoFoodLink pq.StringArray
		if v.Photos != nil {
			for _, v2 := range *v.Photos {
				baseLink := hashing.GenerateHash(r.ServerBaseUrl, v2)
				photoFoodLink = append(photoFoodLink, baseLink)
			}
			list[k].Photos = &photoFoodLink
		}
	}

	return list, nil
}

// @waiter

func (r Repository) WaiterUpdateStatus(ctx context.Context, ids []int64, status string) error {
	claims, err := r.CheckClaims(ctx, auth.RoleWaiter)
	if err != nil {
		return web.NewRequestError(errors.Wrap(err, "extracting claims"), http.StatusInternalServerError)
	}

	// checking waiter policy...
	checked, err := r.NewSelect().Table("orders").Column("id").Where("waiter_id=? and deleted_at isnull and status!='PAID' and status!='CANCELLED'", claims.UserId).Exists(ctx)
	if err != nil {
		return err
	}

	if !checked {
		err = errors.New("cannot update the order")
		return web.NewRequestError(err, http.StatusBadRequest)
	}

	if len(ids) != 0 {
		idRange := "("
		for i := range ids {
			if i != len(ids)-1 {
				idRange += fmt.Sprintf("%d, ", ids[i])
				continue
			}

			idRange += fmt.Sprintf("%d", ids[i])
		}
		idRange += ")"

		// updating orders...
		q := r.NewUpdate().Table("order_menu").Where(fmt.Sprintf("deleted_at isnull and status!='PAID' and status!='CANCELLED' and id in %s", idRange))

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
	}

	return nil
}

// @cashier

func (r Repository) CashierUpdateStatus(ctx context.Context, id int64, status string) error {
	_, err := r.CheckClaims(ctx, auth.RoleCashier)
	if err != nil {
		return web.NewRequestError(errors.Wrap(err, "extracting claims"), http.StatusInternalServerError)
	}

	if strings.ToUpper(status) == "CANCELLED" || strings.ToUpper(status) == "SERVED" {
		// updating orders...
		q := r.NewUpdate().Table("order_menu").Where(fmt.Sprintf("deleted_at isnull and status!='PAID' and status!='CANCELLED' and id = %d", id))

		q.Set("status = ?", strings.ToUpper(status))

		_, err = q.Exec(ctx)
		if err != nil {
			return web.NewRequestError(errors.Wrap(err, "updating order"), http.StatusBadRequest)
		}
		//} else if strings.ToUpper(status) == "DELETE" {
		//	// updating orders...
		//	q := r.NewUpdate().Table("order_menu").Where(fmt.Sprintf("deleted_at isnull and status!='PAID' and id = %d", id))
		//
		//	q.Set("deleted_at = ?", time.Now())
		//	q.Set("deleted_by = ?", claims.UserId)
		//
		//	_, err = q.Exec(ctx)
		//	if err != nil {
		//		return web.NewRequestError(errors.Wrap(err, "updating order"), http.StatusBadRequest)
		//	}
	} else {
		return web.NewRequestError(errors.New("invalid status"), http.StatusBadRequest)
	}

	return nil
}

func (r Repository) CashierUpdateStatusByOrderID(ctx context.Context, orderId int64, status string) error {
	claims, err := r.CheckClaims(ctx, auth.RoleCashier)
	if err != nil {
		return web.NewRequestError(errors.Wrap(err, "extracting claims"), http.StatusInternalServerError)
	}

	// updating orders...
	q := r.NewUpdate().Table("order_menu").Where(fmt.Sprintf("deleted_at isnull and status!='PAID' and status!='CANCELLED' and order_id = %d AND (SELECT t.branch_id FROM tables as t Where t.id = (SELECT o.table_id FROM orders as o Where o.id = order_menu.order_id)) = %d", orderId, *claims.BranchID))

	q.Set("status = ?", strings.ToUpper(status))

	_, err = q.Exec(ctx)
	if err != nil {
		return web.NewRequestError(errors.Wrap(err, "updating order"), http.StatusBadRequest)
	}

	return nil
}

func NewRepository(DB *postgresql.Database) *Repository {
	return &Repository{DB}
}
