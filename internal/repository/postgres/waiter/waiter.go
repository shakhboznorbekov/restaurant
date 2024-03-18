package waiter

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/pkg/errors"
	"github.com/restaurant/foundation/web"
	"github.com/restaurant/internal/auth"
	"github.com/restaurant/internal/entity"
	"github.com/restaurant/internal/pkg/repository/postgresql"
	"github.com/restaurant/internal/repository/postgres"
	"github.com/restaurant/internal/service/hashing"
	"github.com/restaurant/internal/service/waiter"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type Repository struct {
	*postgresql.Database
}

// @admin

func (r Repository) AdminGetList(ctx context.Context, filter waiter.Filter) ([]waiter.AdminGetList, int, error) {
	claims, err := r.CheckClaims(ctx, auth.RoleAdmin)
	if err != nil {
		return nil, 0, err
	}

	whereQuery := fmt.Sprintf(`WHERE w.deleted_at IS NULL and b.restaurant_id = %d and w.role='WAITER'`, *claims.RestaurantID)

	if filter.BranchID != nil {
		whereQuery += fmt.Sprintf(" AND w.branch_id = %d", *filter.BranchID)
	}

	countWhereQuery := whereQuery

	if filter.Limit != nil {
		whereQuery += fmt.Sprintf(" LIMIT '%d'", *filter.Limit)
	}
	if filter.Offset != nil {
		whereQuery += fmt.Sprintf(" OFFSET '%d'", *filter.Offset)
	}

	query := fmt.Sprintf(`
		SELECT 
		    w.id,
		    w.name,
		    w.phone,
		    w.role,
		    b.name as branch_name
		FROM 
		    users as w
		LEFT OUTER JOIN branches as b ON b.id = w.branch_id
		LEFT OUTER JOIN restaurants as r ON r.id = b.restaurant_id
		%s
	`, whereQuery)

	list := make([]waiter.AdminGetList, 0)

	rows, err := r.QueryContext(ctx, query)
	if err != nil {
		return nil, 0, web.NewRequestError(errors.Wrap(err, "select user"), http.StatusInternalServerError)
	}

	err = r.ScanRows(ctx, rows, &list)
	if err != nil {
		return nil, 0, web.NewRequestError(errors.Wrap(err, "scanning waiters"), http.StatusBadRequest)
	}

	countQuery := fmt.Sprintf(`
		SELECT
			count(w.id)
		FROM
		    users as w
		LEFT OUTER JOIN branches as b ON b.id = w.branch_id
		%s
	`, countWhereQuery)

	countRows, err := r.QueryContext(ctx, countQuery)
	if err == sql.ErrNoRows {
		return nil, 0, web.NewRequestError(postgres.ErrNotFound, http.StatusNotFound)
	}
	if err != nil {
		return nil, 0, web.NewRequestError(errors.Wrap(err, "selecting waiters"), http.StatusBadRequest)
	}

	count := 0

	for countRows.Next() {
		if err = countRows.Scan(&count); err != nil {
			return nil, 0, web.NewRequestError(errors.Wrap(err, "scanning user count"), http.StatusBadRequest)
		}
	}

	return list, count, nil
}

// @branch

func (r Repository) BranchGetList(ctx context.Context, filter waiter.Filter) ([]waiter.BranchGetList, int, error) {
	claims, err := r.CheckClaims(ctx, auth.RoleBranch)
	if err != nil {
		return nil, 0, err
	}

	whereQuery := fmt.Sprintf(`WHERE w.deleted_at IS NULL AND w.branch_id = '%d' AND w.role='WAITER'`, *claims.BranchID)

	countWhereQuery := whereQuery

	if filter.Limit != nil {
		whereQuery += fmt.Sprintf(" LIMIT '%d'", *filter.Limit)
	}
	if filter.Offset != nil {
		whereQuery += fmt.Sprintf(" OFFSET '%d'", *filter.Offset)
	}

	query := fmt.Sprintf(`
		SELECT 
		    w.id,
		    w.name,
		    w.phone,
		    w.role,
		    b.name as branch_name,
		    w.status,
		    w.photo,
		    w.address
		FROM 
		    users as w
		LEFT OUTER JOIN branches as b ON b.id = w.branch_id
		%s
	`, whereQuery)

	list := make([]waiter.BranchGetList, 0)

	rows, err := r.QueryContext(ctx, query)
	if err != nil {
		return nil, 0, web.NewRequestError(errors.Wrap(err, "select user"), http.StatusInternalServerError)
	}

	err = r.ScanRows(ctx, rows, &list)
	if err != nil {
		return nil, 0, web.NewRequestError(errors.Wrap(err, "scanning waiters"), http.StatusBadRequest)
	}

	countQuery := fmt.Sprintf(`
		SELECT
			count(w.id)
		FROM
		    users as w
		LEFT OUTER JOIN branches as b ON b.id = w.branch_id
		%s
	`, countWhereQuery)

	countRows, err := r.QueryContext(ctx, countQuery)
	if err == sql.ErrNoRows {
		return nil, 0, web.NewRequestError(postgres.ErrNotFound, http.StatusNotFound)
	}
	if err != nil {
		return nil, 0, web.NewRequestError(errors.Wrap(err, "selecting waiters"), http.StatusBadRequest)
	}

	for i, v := range list {
		var link string
		if v.Photo != nil {
			link = hashing.GenerateHash(r.ServerBaseUrl, *v.Photo)
		}
		list[i].Photo = &link
	}

	count := 0

	for countRows.Next() {
		if err = countRows.Scan(&count); err != nil {
			return nil, 0, web.NewRequestError(errors.Wrap(err, "scanning user count"), http.StatusBadRequest)
		}
	}

	return list, count, nil
}

func (r Repository) BranchGetDetail(ctx context.Context, id int64) (waiter.BranchGetDetail, error) {
	claims, err := r.CheckClaims(ctx, auth.RoleBranch)
	if err != nil {
		return waiter.BranchGetDetail{}, err
	}

	var data entity.User
	err = r.NewSelect().Model(&data).
		Where("id = ? AND branch_id = ? AND deleted_at IS NULL AND role='WAITER'", id, claims.BranchID).
		Scan(ctx)
	if err != nil {
		return waiter.BranchGetDetail{}, err
	}
	var detail waiter.BranchGetDetail

	if data.Rating != nil {
		query := fmt.Sprintf(`
		SELECT
			percent
		FROM
		    service_percentage 
		WHERE deleted_at IS NULL AND branch_id = %d AND id = %d
		
	`, *claims.BranchID, *data.Rating)
		err = r.QueryRowContext(ctx, query).Scan(
			&detail.Rating,
		)
	}

	detail.ID = data.ID
	detail.Name = data.Name
	detail.Phone = data.Phone
	detail.Gender = data.Gender
	detail.Role = data.Role
	detail.Address = data.Address

	if data.Photo != nil {
		photo := hashing.GenerateHash(r.ServerBaseUrl, *data.Photo)
		detail.Photo = &photo
	}

	birthDate := data.BirthDate.Format("02.01.2006")
	detail.BirthDate = &birthDate

	return detail, nil
}

func (r Repository) BranchCreate(ctx context.Context, request waiter.BranchCreateRequest) (waiter.BranchCreateResponse, error) {
	claims, err := r.CheckClaims(ctx, auth.RoleBranch)
	if err != nil {
		return waiter.BranchCreateResponse{}, err
	}

	err = r.ValidateStruct(&request, "Phone", "Name", "BirthDate", "Gender", "Password")
	if err != nil {
		return waiter.BranchCreateResponse{}, err
	}

	role := auth.RoleWaiter

	var birthDate time.Time
	if request.BirthDate != nil {
		birthDate, err = time.Parse("02.01.2006", *request.BirthDate)

		if err != nil {
			return waiter.BranchCreateResponse{}, web.NewRequestError(fmt.Errorf("incorrect birth-date format: '%v'", err), http.StatusBadRequest)
		}
	}

	var gender string
	if request.Gender != nil {
		if (strings.ToUpper(*request.Gender) == "M") && (strings.ToUpper(*request.Gender) == "F") {
			return waiter.BranchCreateResponse{}, web.NewRequestError(errors.New("incorrect gender. gender should be M (male) or F (female)"), http.StatusBadRequest)
		}
		gender = strings.ToUpper(*request.Gender)
	}

	response := waiter.BranchCreateResponse{
		Name:           request.Name,
		Password:       request.Password,
		Phone:          request.Phone,
		BirthDate:      &birthDate,
		Gender:         &gender,
		Role:           &role,
		CreatedAt:      time.Now(),
		BranchID:       claims.BranchID,
		CreatedBy:      claims.UserId,
		ServicePercent: request.ServicePercentageID,
		Photo:          request.PhotoLink,
		Address:        request.Address,
	}

	_, err = r.NewInsert().Model(&response).Exec(ctx)
	if err != nil {
		return waiter.BranchCreateResponse{}, web.NewRequestError(errors.Wrap(err, "creating user"), http.StatusBadRequest)
	}

	if response.Photo != nil {
		photo := hashing.GenerateHash(r.ServerBaseUrl, *response.Photo)
		response.Photo = &photo
	}

	return response, nil
}

func (r Repository) BranchUpdateAll(ctx context.Context, request waiter.BranchUpdateRequest) error {
	claims, err := r.CheckClaims(ctx, auth.RoleBranch)
	if err != nil {
		return err
	}

	if err = r.ValidateStruct(&request, "ID", "Name", "Phone", "BirthDate", "Gender", "Password"); err != nil {
		return err
	}

	if (strings.ToUpper(*request.Gender) == "M") && (strings.ToUpper(*request.Gender) == "F") {
		return web.NewRequestError(errors.New("incorrect gender. gender should be M (male) or F (female)"), http.StatusBadRequest)
	}

	var birthDate time.Time
	if request.BirthDate != nil {
		birthDate, err = time.Parse("02.01.2006", *request.BirthDate)
		if err != nil {
			return web.NewRequestError(errors.New("incorrect birth_date format"), http.StatusBadRequest)
		}
	}

	q := r.NewUpdate().Table("users").Where("deleted_at IS NULL AND id = ? AND w.role='WAITER'", request.ID)

	q.Set("name = ?", request.Name)
	q.Set("birth_date = ?", birthDate)
	q.Set("gender = ?", strings.ToUpper(*request.Gender))
	q.Set("updated_at = ?", time.Now())
	q.Set("updated_by = ?", claims.UserId)
	q.Set("photo = ?", request.PhotoLink)
	q.Set("rating = ?", request.ServicePercentageID)
	q.Set("address = ?", request.Address)

	_, err = q.Exec(ctx)
	if err != nil {
		return web.NewRequestError(errors.Wrap(err, "updating waiter"), http.StatusBadRequest)
	}

	return nil
}

func (r Repository) BranchUpdateColumns(ctx context.Context, request waiter.BranchUpdateRequest) error {
	claims, err := r.CheckClaims(ctx, auth.RoleBranch)
	if err != nil {
		return err
	}

	if err = r.ValidateStruct(&request, "ID"); err != nil {
		return err
	}

	if (request.Gender != nil) && (strings.ToUpper(*request.Gender) == "M") && (strings.ToUpper(*request.Gender) == "F") {
		return web.NewRequestError(errors.New("incorrect gender. gender should be M (male) or F (female)"), http.StatusBadRequest)
	}

	var birthDate time.Time
	if request.BirthDate != nil {
		birthDate, err = time.Parse("02.01.2006", *request.BirthDate)
		if err != nil {
			return web.NewRequestError(errors.New("incorrect birth_date format"), http.StatusBadRequest)
		}
	}

	q := r.NewUpdate().Table("users").Where("deleted_at IS NULL AND id = ? AND role='WAITER'", request.ID)

	if request.Name != nil {
		q.Set("name = ?", request.Name)
	}
	if request.BirthDate != nil {
		q.Set("birth_date = ?", birthDate)
	}
	if request.Gender != nil {
		q.Set("gender = ?", strings.ToUpper(*request.Gender))
	}

	if request.PhotoLink != nil {
		q.Set("photo = ?", request.PhotoLink)
	}
	if request.ServicePercentageID != nil {
		q.Set("rating = ?", request.ServicePercentageID)
	}
	if request.Address != nil {
		q.Set("address = ?", request.Address)
	}

	q.Set("updated_at = ?", time.Now())
	q.Set("updated_by = ?", claims.UserId)

	_, err = q.Exec(ctx)
	if err != nil {
		return web.NewRequestError(errors.Wrap(err, "updating waiter"), http.StatusBadRequest)
	}

	return nil
}

func (r Repository) BranchDelete(ctx context.Context, id int64) error {
	return r.DeleteRow(ctx, "waiters", id, auth.RoleBranch)
}

func (r Repository) BranchUpdateStatus(ctx context.Context, id int64, status string) error {
	claims, err := r.CheckClaims(ctx, auth.RoleBranch)
	if err != nil {
		return err
	}

	q := r.NewUpdate().Table("users").Where("deleted_at isnull and role='WAITER' and branch_id = ? and id = ?", *claims.BranchID, id)

	q.Set("updated_at = ?", time.Now())
	q.Set("updated_by = ?", claims.UserId)
	if status == "active" || status == "inactive" {
		q.Set("status = ?", status)
	}

	_, err = q.Exec(ctx)
	if err != nil {
		return web.NewRequestError(errors.Wrap(err, "updating waiter status"), http.StatusBadRequest)
	}

	return nil
}

// others

func (r Repository) UpdatePassword(ctx context.Context, request waiter.BranchUpdatePassword) error {
	q := r.NewUpdate().Table("users").Where("deleted_at IS NULL AND id = ? AND role='WAITER'", request.ID)

	q.Set("password = ?", request.Password)

	_, err := q.Exec(ctx)
	if err != nil {
		return web.NewRequestError(errors.Wrap(err, "updating waiter password"), http.StatusBadRequest)
	}

	return nil
}

func (r Repository) UpdatePhone(ctx context.Context, request waiter.BranchUpdatePhone) error {
	q := r.NewUpdate().Table("users").Where("deleted_at IS NULL AND id = ? AND role='WAITER'", request.ID)

	q.Set("phone = ?", request.Phone)

	_, err := q.Exec(ctx)
	if err != nil {
		return web.NewRequestError(errors.Wrap(err, "updating waiter phone"), http.StatusBadRequest)
	}

	return nil
}

// @waiter

func (r Repository) WaiterGetMe(ctx context.Context) (*waiter.GetMeResponse, error) {
	claims, err := r.CheckClaims(ctx, auth.RoleWaiter)
	if err != nil {
		return nil, err
	}

	query := fmt.Sprintf(`SELECT
   									u.id,
   									u.name,
   									u.photo,
   									coalesce(u.rating, 0),
   									coalesce((select sum(op.price) from order_payment op join orders o on op.order_id = o.id where o.deleted_at isnull and op.deleted_at isnull and o.status='PAID' and o.waiter_id=u.id), 0) as profit,
   									(select count(o.id) from orders o where o.status='PAID' and o.deleted_at isnull and o.waiter_id = u.id) as order_count,
   									to_char(u.birth_date, 'DD.MM.YYYY'),
   									u.phone,
   									u.address,
   									u.attendance_status
								 FROM users u
								 WHERE
								     u.id = '%d'
								   AND
								     u.deleted_at ISNULL
								   AND
								     u.role='WAITER'`, claims.UserId)

	var response waiter.GetMeResponse
	if err = r.QueryRowContext(ctx, query).Scan(&response.Id, &response.Name, &response.Photo, &response.Rating, &response.Profit, &response.OrderCount, &response.BirthDate, &response.Phone, &response.Address, &response.AttendanceStatus); err != nil {
		return nil, err
	}

	if response.Photo != nil {
		photo := hashing.GenerateHash(r.ServerBaseUrl, *response.Photo)
		response.Photo = &photo
	}

	return &response, nil
}

func (r Repository) WaiterGetPersonalInfo(ctx context.Context) (*waiter.GetPersonalInfoResponse, error) {
	claims, err := r.CheckClaims(ctx, auth.RoleWaiter)
	if err != nil {
		return nil, err
	}

	query := fmt.Sprintf(`SELECT
   									u.id,
   									u.name,
   									u.birth_date,
   									u.phone,
   									u.address
								 FROM users u
								 WHERE
								     u.id = '%d'
								   AND
								     u.deleted_at ISNULL
								   AND
								     u.role='WAITER'`, claims.UserId)

	var response waiter.GetPersonalInfoResponse
	if err = r.QueryRowContext(ctx, query).Scan(&response.Id, &response.Name, &response.BirthDate, &response.Phone, &response.Address); err != nil {
		return nil, err
	}

	return &response, nil
}

func (r Repository) WaiterUpdatePhoto(ctx context.Context, request waiter.WaiterPhotoUpdateRequest) error {
	claims, err := r.CheckClaims(ctx, auth.RoleWaiter)
	if err != nil {
		return err
	}

	q := r.NewUpdate().Table("users").Where("deleted_at is null and role='WAITER' and branch_id = ? and id = ? ", *claims.BranchID, claims.UserId)

	if request.PhotoLink != nil {
		q.Set("photo = ?", request.PhotoLink)
	}
	q.Set("updated_at = ?", time.Now())
	q.Set("updated_by = ?", claims.UserId)

	_, err = q.Exec(ctx)
	if err != nil {
		return web.NewRequestError(errors.Wrap(err, "updating waiter status"), http.StatusBadRequest)
	}

	return nil
}

// @cashier

func (r Repository) CashierGetList(ctx context.Context, filter waiter.Filter) ([]waiter.CashierGetList, int, error) {
	claims, err := r.CheckClaims(ctx, auth.RoleCashier)
	if err != nil {
		return nil, 0, err
	}

	whereQuery := fmt.Sprintf(`WHERE w.deleted_at IS NULL AND w.branch_id = %d and w.role='WAITER'`, *claims.BranchID)

	//if filter.BranchID != nil {
	//	whereQuery += fmt.Sprintf(" AND w.branch_id = %d", *filter.BranchID)
	//}

	countWhereQuery := whereQuery

	if filter.Limit != nil {
		whereQuery += fmt.Sprintf(" LIMIT '%d'", *filter.Limit)
	}
	if filter.Offset != nil {
		whereQuery += fmt.Sprintf(" OFFSET '%d'", *filter.Offset)
	}

	query := fmt.Sprintf(`
		SELECT 
		    w.id,
		    w.name,
		    w.phone,
		    w.role,
		    b.name as branch_name
		FROM 
		    users as w
		LEFT OUTER JOIN branches as b ON b.id = w.branch_id
		%s
	`, whereQuery)

	list := make([]waiter.CashierGetList, 0)

	rows, err := r.QueryContext(ctx, query)
	if err != nil {
		return nil, 0, web.NewRequestError(errors.Wrap(err, "select user"), http.StatusInternalServerError)
	}

	err = r.ScanRows(ctx, rows, &list)
	if err != nil {
		return nil, 0, web.NewRequestError(errors.Wrap(err, "scanning waiters"), http.StatusBadRequest)
	}

	countQuery := fmt.Sprintf(`
		SELECT
			count(w.id)
		FROM
		    users as w
		LEFT OUTER JOIN branches as b ON b.id = w.branch_id
		%s
	`, countWhereQuery)

	countRows, err := r.QueryContext(ctx, countQuery)
	if err == sql.ErrNoRows {
		return nil, 0, web.NewRequestError(postgres.ErrNotFound, http.StatusNotFound)
	}
	if err != nil {
		return nil, 0, web.NewRequestError(errors.Wrap(err, "selecting waiters"), http.StatusBadRequest)
	}

	count := 0

	for countRows.Next() {
		if err = countRows.Scan(&count); err != nil {
			return nil, 0, web.NewRequestError(errors.Wrap(err, "scanning user count"), http.StatusBadRequest)
		}
	}

	return list, count, nil
}

func NewRepository(DB *postgresql.Database) *Repository {
	return &Repository{DB}
}

// CalculateWaitersKPI : the func automatically executes every 30 days
func (r Repository) CalculateWaitersKPI(ctx context.Context) error {
	query := fmt.Sprintf(`SELECT 
    									id 
								 FROM 
								     branches b 
								 WHERE 
								     deleted_at isnull 
								   and 
								     (select COUNT(w.id) from users w where w.branch_id=b.id and w.role='WAITER' and w.deleted_at isnull and w.status='active') != 0`)
	rows, err := r.QueryContext(ctx, query)
	if err != nil {
		return err
	}

	for rows.Next() {
		var id int64
		if err = rows.Scan(&id); err != nil {
			return err
		}

		query = fmt.Sprintf(`SELECT 
    									w.id 
									FROM users w 
									WHERE 
									    w.status = 'active' 
									  and 
									    w.deleted_at isnull 
									  and 
									    w.role='WAITER' 
									  and 
									    w.branch_id='%d' 
-- 									  and 
-- 									    TO_CHAR(w.created_at, 'YYYY.MM.DD')::date <= TO_CHAR(CURRENT_DATE - INTERVAL '20 days', 'YYYY.MM.DD')::date 
									  and
									    (select count(id) from orders where status='PAID' and waiter_id=w.id and created_at >= now() - interval '24 hours') != 0
									  and 
									    (select count(id) from attendances where came_at >= now() - interval '24 hours' and w.id = user_id) != 0`, id)
		rowWaiters, err := r.QueryContext(ctx, query)
		if err != nil {
			return err
		}

		var waiters []int64

		if err = r.ScanRows(ctx, rowWaiters, &waiters); err != nil {
			return err
		}

		if err = r.calculateKPI(ctx, waiters); err != nil {
			return err
		}
	}

	return nil
}

// this field is for only calculating KPI for waiters. PLS! DO NOT TOUCH.
// this api is used every 30 days...
// after some time we can add none fixed time

// CalculateClientReview : waiters arg is for id of waiters, we retrieve count of reviews and given scores for orders
// max(score)/given(score) = waiter score(x) [1 <= x <= -1]
func (r Repository) CalculateClientReview(ctx context.Context, waiters []int64) (map[int64]float64, error) {
	response := make(map[int64]float64)

	for i := range waiters {
		var count, sum float64
		query := fmt.Sprintf(`SELECT 
    										COUNT(wr.id)::double precision, 
    										COALESCE(SUM(wr.score)::double precision, 0) 
									 FROM waiter_reviews wr
									 	JOIN orders o ON wr.order_id = o.id
									 	JOIN users w ON w.id = o.waiter_id
									 WHERE 
									     w.id = '%d'
									   AND 
									     TO_CHAR(wr.created_at, 'YYYY.MM.DD')::date >= TO_CHAR(CURRENT_DATE - INTERVAL '1 day', 'YYYY.MM.DD')::date`, waiters[i])
		if err := r.QueryRowContext(ctx, query).Scan(&count, &sum); err != nil && !errors.Is(err, sql.ErrNoRows) {
			return nil, err
		}
		if count == 0 {
			response[waiters[i]] = 0.5
			continue
		}

		score := sum / (count * 5)

		response[waiters[i]] = score
	}

	return response, nil
}

// CalculateOrderScore : waiters arg is for id of waiters, we retrieve waiter id with the approximate score for order's count
// 0.5 + (waiter.order.count - order.count.minimum) * (1.5 / (order.count.maximum - order.count.minimum))
func (r Repository) CalculateOrderScore(ctx context.Context, waiters []int64) (map[int64]float64, error) {
	var (
		maxCount, minCount       float64
		minWaiterId, maxWaiterId int64
	)
	response := make(map[int64]float64)

	for i := range waiters {
		var count float64
		query := fmt.Sprintf(`SELECT 
    										count(o.id)::double precision
									 FROM orders o 
									 WHERE 
									     o.waiter_id = '%d' 
									   AND 
									     o.status = 'PAID' 
									   AND 
									     TO_CHAR(o.created_at, 'YYYY.MM.DD')::date >= TO_CHAR(CURRENT_DATE - INTERVAL '1 day', 'YYYY.MM.DD')::date`, waiters[i])
		if err := r.QueryRowContext(ctx, query).Scan(&count); err != nil && !errors.Is(err, sql.ErrNoRows) {
			return nil, err
		}
		if i == 0 || minCount > count {
			minCount = count
			minWaiterId = waiters[i]
		}
		if maxCount < count {
			maxCount = count
			maxWaiterId = waiters[i]
		}

		response[waiters[i]] = count
	}

	if maxCount == minCount {
		for i := range response {
			response[i] = 0.5
		}
		return response, nil
	}

	point := approximate(1.5 / (maxCount - minCount))

	for i := range response {
		if minWaiterId != i && maxWaiterId != i && response[i] != response[maxWaiterId] && response[i] != response[minWaiterId] {
			response[i] = 0.5 + (response[i]-response[minWaiterId])*point
		} else if maxWaiterId != i && response[i] == response[maxWaiterId] {
			response[i] = 2.0
		} else if minWaiterId != i && response[i] == response[minWaiterId] {
			response[i] = 0.5
		}
	}
	response[minWaiterId] = 0.5
	response[maxWaiterId] = 2.0

	return response, nil
}

// CalculateProfitScore : waiters arg is for id of waiters, we retrieve waiter id with the approximate score for branches profit
// 0.5 + (waiter.profit - waiter.profit.minimum) * (2.5 / (waiter.profit.maximum - waiter.profit.minimum))
func (r Repository) CalculateProfitScore(ctx context.Context, waiters []int64) (map[int64]float64, error) {
	var (
		minWaiterId, maxWaiterId int64
		minProfit, maxProfit     float64
	)
	response := make(map[int64]float64)

	for i := range waiters {
		var profit float64
		query := fmt.Sprintf(`SELECT 
    										SUM(op.price)::double precision as profit
									 FROM order_payment op 
									     JOIN orders o 
									         ON op.order_id = o.id 
									 WHERE 
									     o.status = 'PAID' 
									   AND 
									     op.deleted_at isnull 
									   AND 
									     o.waiter_id = '%d'
									   AND 
									     TO_CHAR(op.created_at, 'YYYY.MM.DD')::date >= TO_CHAR(CURRENT_DATE - INTERVAL '1 day', 'YYYY.MM.DD')::date`, waiters[i])
		if err := r.QueryRowContext(ctx, query).Scan(&profit); err != nil && !errors.Is(err, sql.ErrNoRows) {
			return nil, err
		}

		if i == 0 || minProfit > profit {
			minProfit = profit
			minWaiterId = waiters[i]
		}
		if maxProfit < profit {
			maxProfit = profit
			maxWaiterId = waiters[i]
		}

		response[waiters[i]] = profit
	}
	if maxProfit == minProfit {
		for i := range response {
			response[i] = 3.0
		}
		return response, nil
	}

	point := approximate(2.5 / (maxProfit - minProfit))

	for i := range response {
		if minWaiterId != i && maxWaiterId != i && response[i] != response[maxWaiterId] && response[i] != response[minWaiterId] {
			response[i] = 0.5 + (response[i]-response[minWaiterId])*point
		} else if maxWaiterId != i && response[i] == response[maxWaiterId] {
			response[i] = 3.0
		} else if minWaiterId != i && response[i] == response[minWaiterId] {
			response[i] = 0.5
		}
	}
	response[minWaiterId] = 0.5
	response[maxWaiterId] = 3.0

	return response, nil
}

// CalculateOrderServeLikelihood : waiters arg is for id of waiters, we retrieve waiter id with the approximate score for branches order acceptance
// waiter.score.total * 2 / order.total * 3; score.minimum = 0.5
func (r Repository) CalculateOrderServeLikelihood(ctx context.Context, waiters []int64) (map[int64]float64, error) {
	response := make(map[int64]float64)

	for i := range waiters {
		var (
			cancel float64
			accept float64
			total  float64
		)
		cancelledOrdersCountQuery := fmt.Sprintf(`SELECT (SELECT 
    														count(o.id)::double precision
														 FROM orders o
															 JOIN tables t
																 ON o.table_id = t.id
															 JOIN users w
																 ON t.branch_id = w.branch_id
														 WHERE 
															 o.status = 'CANCELLED' 
														   AND 
															 o.waiter_id isnull 
														   AND 
															 o.accepted_at isnull 
														   AND 
															 TO_CHAR(o.created_at, 'YYYY.MM.DD')::date >= TO_CHAR(CURRENT_DATE - INTERVAL '1 day', 'YYYY.MM.DD')::date 
														   AND 
														     w.id = '%d') AS cancelled, 
    													 (SELECT 
    														count(o.id)::double precision
														 FROM orders o
														 WHERE 
															 o.status = 'PAID' 
														   AND 
															 TO_CHAR(o.created_at, 'YYYY.MM.DD')::date >= TO_CHAR(CURRENT_DATE - INTERVAL '1 day', 'YYYY.MM.DD')::date 
														   AND 
														     o.waiter_id = '%d') as accepted`, waiters[i], waiters[i])
		if err := r.QueryRowContext(ctx, cancelledOrdersCountQuery).Scan(&cancel, &accept); err != nil && !errors.Is(err, sql.ErrNoRows) {
			return nil, err
		}
		total = cancel + accept

		if total == 0 {
			response[waiters[i]] = 0.5
			continue
		}

		score := (((accept * 3) + (cancel * (-2))) * 2) / (total * 3)
		if score < 0.5 {
			score = 0.5
		}

		response[waiters[i]] = score
	}

	return response, nil
}

// CalculateProfitAndCountRelation : waiters arg is for id of waiters, we retrieve waiter id with the approximate score for branches order and profit relation
// waiter.order.total / waiter.profit.total = relation
func (r Repository) CalculateProfitAndCountRelation(ctx context.Context, waiters []int64) (map[int64]float64, error) {
	var (
		maxRelation, minRelation float64
		maxWaiterId, minWaiterId int64
	)
	response := make(map[int64]float64)

	for i := range waiters {
		var (
			count  float64
			profit float64
		)
		query := fmt.Sprintf(`SELECT 
    									(SELECT 
											 count(o.id)::double precision
										 FROM orders o 
										 WHERE 
											 o.waiter_id = '%d' 
										   AND 
											 o.status = 'PAID' 
										   AND 
											 TO_CHAR(o.created_at, 'YYYY.MM.DD')::date >= TO_CHAR(CURRENT_DATE - INTERVAL '1 day', 'YYYY.MM.DD')::date) AS count,
    									 (SELECT 
    									     COALESCE(sum(op.price)::double precision, 0)
    									  FROM 
    									      order_payment op 
    									          JOIN orders o 
    									              ON op.order_id = o.id 
    									  WHERE 
    									      o.waiter_id = '%d' 
    									    AND 
    									      o.status = 'PAID' 
    									    AND 
    									      op.deleted_at isnull
										    AND 
											  TO_CHAR(o.created_at, 'YYYY.MM.DD')::date >= TO_CHAR(CURRENT_DATE - INTERVAL '1 day', 'YYYY.MM.DD')::date) AS profit`, waiters[i], waiters[i])
		if err := r.QueryRowContext(ctx, query).Scan(&count, &profit); err != nil && !errors.Is(err, sql.ErrNoRows) {
			return nil, err
		}
		if count == 0 {
			response[waiters[i]] = 0.5
			continue
		}

		relation := profit / count

		if i == 0 || minRelation > relation {
			minRelation = relation
			minWaiterId = waiters[i]
		}
		if maxRelation < relation {
			maxRelation = relation
			maxWaiterId = waiters[i]
		}

		response[waiters[i]] = relation
	}
	if maxRelation == minRelation {
		for i := range response {
			response[i] = 2.0
		}

		return response, nil
	}

	point := approximate(1.5 / (maxRelation - minRelation))

	for i := range response {
		if minWaiterId != i && maxWaiterId != i && response[i] != response[maxWaiterId] && response[i] != response[minWaiterId] {
			response[i] = 0.5 + (response[i]-response[minWaiterId])*point
		} else if maxWaiterId != i && response[i] == response[maxWaiterId] {
			response[i] = 2.0
		} else if minWaiterId != i && response[i] == response[minWaiterId] {
			response[i] = 0.5
		}
	}
	response[minWaiterId] = 0.5
	response[maxWaiterId] = 2.0

	return response, nil
}

// calculateKPI : calculates waiters kpi and updates at database
func (r Repository) calculateKPI(ctx context.Context, waiters []int64) error {
	first, err := r.CalculateClientReview(ctx, waiters)
	if err != nil {
		return err
	}
	second, err := r.CalculateOrderScore(ctx, waiters)
	if err != nil {
		return err
	}
	third, err := r.CalculateProfitScore(ctx, waiters)
	if err != nil {
		return err
	}
	fourth, err := r.CalculateOrderServeLikelihood(ctx, waiters)
	if err != nil {
		return err
	}
	fifth, err := r.CalculateProfitAndCountRelation(ctx, waiters)
	if err != nil {
		return err
	}

	for i := range waiters {
		var kpi float64
		if v, ok := first[waiters[i]]; ok {
			kpi += v
		}
		if v, ok := second[waiters[i]]; ok {
			kpi += v
		}
		if v, ok := third[waiters[i]]; ok {
			kpi += v
		}
		if v, ok := fourth[waiters[i]]; ok {
			kpi += v
		}
		if v, ok := fifth[waiters[i]]; ok {
			kpi += v
		}

		update := fmt.Sprintf(`UPDATE users SET rating='%.1f' WHERE id='%d'`, kpi, waiters[i])
		if _, err = r.ExecContext(ctx, update); err != nil {
			return err
		}
	}

	return nil
}

// approximate real number
func approximate(x float64) float64 {
	x, _ = strconv.ParseFloat(fmt.Sprintf("%.5f", x), 64)

	return x
}
