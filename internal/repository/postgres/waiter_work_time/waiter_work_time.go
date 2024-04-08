package waiter_work_time

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/pkg/errors"
	"math/rand"
	"net/http"
	"restu-backend/foundation/web"
	"restu-backend/internal/auth"
	"restu-backend/internal/pkg/repository/postgresql"
	"restu-backend/internal/repository/postgres"
	"restu-backend/internal/service/attendance"
	"restu-backend/internal/service/hashing"
	dto "restu-backend/internal/service/waiter_work_time"
	"time"
)

type Repository struct {
	*postgresql.Database
}

func NewRepository(DB *postgresql.Database) *Repository {
	return &Repository{DB}
}

// @waiter

func (r Repository) WaiterCreate(ctx context.Context, request attendance.WaiterCreateResponse) error {
	claims, err := r.CheckClaims(ctx, auth.RoleWaiter)
	if err != nil {
		return err
	}

	if err := r.ValidateStruct(&request, "ID", "Action"); err != nil {
		return err
	}

	date := request.ActionTime.Format("02.01.2006")
	data, err := r.GetDetailByWaiterIDAndDate(ctx, dto.Filter{Date: &date, WaiterID: &claims.UserId})
	if err != nil {
		if err.Error() == postgres.ErrNotFound.Error() {
			p := make([]dto.Period, 0)
			if *request.Action == "ENTER" {
				p = append(p, dto.Period{
					ComeTime: request.CameAt,
				})
			} else {
				p = append(p, dto.Period{
					GoneTime: request.GoneAt,
				})
			}

			response := dto.CreateResponse{
				WaiterID: request.UserID,
				Date:     request.ActionTime,
				Periods:  p,
			}
			_, err = r.NewInsert().Model(&response).Returning("id").Exec(ctx, &response.ID)
			if err != nil {
				return web.NewRequestError(errors.Wrap(err, "creating user work time"), http.StatusBadRequest)
			}
			return nil
		}
		return err
	}

	response := dto.CreateResponse{}
	if *request.Action == "ENTER" {
		p := dto.Period{
			ComeTime: request.CameAt,
		}
		data.Periods = append(data.Periods, p)
	} else {
		data.Periods[len(data.Periods)-1].GoneTime = request.ActionTime
		if data.Periods[len(data.Periods)-1].ComeTime == nil {
			t, _ := time.Parse("02.01.2006", date)
			d := int(request.ActionTime.Sub(t).Minutes())
			data.Periods[len(data.Periods)-1].ComeTime = &t
			data.Periods[len(data.Periods)-1].Duration = &d
		} else {
			d := int(request.ActionTime.Sub(*data.Periods[len(data.Periods)-1].ComeTime).Minutes())
			data.Periods[len(data.Periods)-1].Duration = &d
		}
	}
	response.WaiterID = &claims.UserId
	response.Date = request.ActionTime
	response.Periods = data.Periods

	_, err = r.NewUpdate().Model(&response).Where("id = ?", data.ID).Returning("id").Exec(ctx, &response.ID)
	if err != nil {
		return web.NewRequestError(errors.Wrap(err, "creating user work time"), http.StatusBadRequest)
	}

	return nil
}

func (r Repository) GetDetailByWaiterIDAndDate(ctx context.Context, filter dto.Filter) (dto.GetDetailByWaiterIDAndDateResponse, error) {

	if err := r.ValidateStruct(&filter, "WaiterID"); err != nil {
		return dto.GetDetailByWaiterIDAndDateResponse{}, err
	}
	rand.Seed(time.Now().UnixNano())
	date := time.Now().Format("2006-01-02")
	if filter.Date != nil {
		Date, err := time.Parse("02.01.2006", *filter.Date)
		if err != nil {
			return dto.GetDetailByWaiterIDAndDateResponse{}, web.NewRequestError(errors.Wrap(err, "date parse"), http.StatusBadRequest)
		}
		date = Date.Format("2006-01-02")
	}

	query := fmt.Sprintf(`
		SELECT
			id,
			TO_CHAR(date,'DD.MM.YYYY'),
			periods
		FROM
		    waiter_work_time
		WHERE waiter_id  = '%d' AND date = '%s'
	`, *filter.WaiterID, date)

	var detail dto.GetDetailByWaiterIDAndDateResponse

	var periods []byte
	err := r.QueryRowContext(ctx, query).Scan(
		&detail.ID,
		&detail.Date,
		&periods,
	)
	if errors.Is(err, sql.ErrNoRows) {
		return dto.GetDetailByWaiterIDAndDateResponse{}, web.NewRequestError(postgres.ErrNotFound, http.StatusBadRequest)
	}

	if err != nil {
		return dto.GetDetailByWaiterIDAndDateResponse{}, web.NewRequestError(errors.Wrap(err, "selecting staffing detail"), http.StatusBadRequest)
	}

	err = json.Unmarshal(periods, &detail.Periods)
	if err != nil {
		return dto.GetDetailByWaiterIDAndDateResponse{}, web.NewRequestError(errors.Wrap(err, "unmarshal period"), http.StatusBadRequest)
	}
	totalDuration := 0
	for _, v := range detail.Periods {
		if v.Duration != nil {
			totalDuration += *v.Duration
		} else {
			if v.ComeTime != nil && v.GoneTime == nil {
				t, _ := time.Parse("02.01.2006 15:04:05", time.Now().Format("02.01.2006 15:04:05"))
				if date != time.Now().Format("2006-01-02") {
					t, _ = time.Parse("2006-01-02", date)
					t = t.Add(24 * time.Hour).Add(-1 * time.Second)
				}
				fmt.Println(t, "\n", *v.ComeTime)
				d := int(t.Sub(*v.ComeTime).Minutes())
				v.GoneTime = &t
				v.Duration = &d
				detail.Now = true
				totalDuration += d
			} else if v.ComeTime == nil && v.GoneTime != nil {
				t, _ := time.Parse("2006-01-02", date)
				d := int(v.GoneTime.Sub(t).Minutes())
				v.ComeTime = &t
				v.Duration = &d
				totalDuration += d
			} else {
				if v.ComeTime != nil {
					d := int(v.GoneTime.Sub(*v.ComeTime).Minutes())
					v.Duration = &d
					totalDuration += d
				}
			}
		}
		//enter := v.ComeTime.Format("15:04:05")
		//exit := v.GoneTime.Format("15:04:05")
		enterInMinute := 0
		if v.ComeTime != nil {
			enterInMinute = v.ComeTime.Hour()*60 + v.ComeTime.Minute()
		}
		exitInMinute := 0
		if v.GoneTime != nil {
			exitInMinute = v.GoneTime.Hour()*60 + v.GoneTime.Minute()
		}
		detail.PeriodsResponse = append(detail.PeriodsResponse, dto.PeriodResponse{
			ComeTime: &enterInMinute,
			GoneTime: &exitInMinute,
			Duration: v.Duration,
		})
	}
	detail.TotalDuration = &totalDuration

	return detail, nil
}

func (r Repository) WaiterGetListWorkTime(ctx context.Context, filter dto.ListFilter) ([]dto.GetListResponse, int, error) {
	claims, err := r.CheckClaims(ctx, auth.RoleWaiter)
	if err != nil {
		return nil, 0, err
	}

	whereQuery := fmt.Sprintf(`WHERE wwt.waiter_id = '%d'`, claims.UserId)

	orderQuery := "ORDER BY wwt.date desc"

	var limitQuery, offsetQuery string

	if filter.Page != nil && filter.Limit != nil {
		offset := (*filter.Page - 1) * (*filter.Limit)
		filter.Offset = &offset
	}

	if filter.Limit != nil {
		limitQuery += fmt.Sprintf(" LIMIT %d", *filter.Limit)
	}

	if filter.Offset != nil {
		offsetQuery += fmt.Sprintf(" OFFSET %d", *filter.Offset)
	}

	rand.Seed(time.Now().UnixNano())
	date := time.Now().Format("2006-01-02")

	query := fmt.Sprintf(`
					SELECT 
					  wwt.id,
					  TO_CHAR(wwt.date,'DD.MM.YYYY'),
					  wwt.periods
					FROM 
					    waiter_work_time wwt
					%s %s %s %s
	`, whereQuery, orderQuery, limitQuery, offsetQuery)

	rows, err := r.QueryContext(ctx, query)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, 0, web.NewRequestError(postgres.ErrNotFound, http.StatusBadRequest)
	}
	if err != nil {
		return nil, 0, web.NewRequestError(errors.Wrap(err, "selecting waiter_work_time"), http.StatusInternalServerError)
	}

	var list []dto.GetListResponse

	for rows.Next() {
		var detail dto.GetListResponse
		var periods []byte
		if err = rows.Scan(
			&detail.ID,
			&detail.Date,
			&periods,
		); err != nil {
			return nil, 0, web.NewRequestError(errors.Wrap(err, "scanning waiter_work_time"), http.StatusInternalServerError)
		}

		err = json.Unmarshal(periods, &detail.Periods)
		if err != nil {
			return nil, 0, web.NewRequestError(errors.Wrap(err, "unmarshal periods"), http.StatusBadRequest)
		}
		totalDuration := 0
		for _, v := range detail.Periods {
			if v.Duration != nil {
				totalDuration += *v.Duration
			} else {
				if v.ComeTime != nil && v.GoneTime == nil {
					t := time.Now()
					if date != time.Now().Format("2006-01-02") {
						t, _ = time.Parse("2006-01-02", date)
						t = t.Add(24 * time.Hour).Add(-1 * time.Second)
					}

					d := int(t.Sub(*v.ComeTime).Minutes())

					v.GoneTime = &t
					v.Duration = &d
					detail.Now = true
					totalDuration += d
				} else if v.ComeTime == nil && v.GoneTime != nil {
					t, _ := time.Parse("2006-01-02", date)
					d := int(v.GoneTime.Sub(t).Minutes())
					v.ComeTime = &t
					v.Duration = &d
					totalDuration += d
				} else {
					if v.ComeTime != nil {
						d := int(v.GoneTime.Sub(*v.ComeTime).Minutes())
						v.Duration = &d
						totalDuration += d
					}
				}
			}
			//enter := v.ComeTime.Format("15:04:05")
			//exit := v.GoneTime.Format("15:04:05")
			enterInMinute := 0
			if v.ComeTime != nil {
				enterInMinute = v.ComeTime.Hour()*60 + v.ComeTime.Minute()
			}
			exitInMinute := 0
			if v.GoneTime != nil {
				exitInMinute = v.GoneTime.Hour()*60 + v.GoneTime.Minute()
			}

			detail.PeriodsResponse = append(detail.PeriodsResponse, dto.PeriodResponse{
				ComeTime: &enterInMinute,
				GoneTime: &exitInMinute,
				Duration: v.Duration,
			})
		}

		detail.TotalDuration = &totalDuration
		list = append(list, detail)
	}

	countQuery := fmt.Sprintf(`
		SELECT
			count(wwt.id)
		FROM
		    waiter_work_time wwt
		%s
	`, whereQuery)

	countRows, err := r.QueryContext(ctx, countQuery)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, 0, web.NewRequestError(postgres.ErrNotFound, http.StatusBadRequest)
	}
	if err != nil {
		return nil, 0, web.NewRequestError(errors.Wrap(err, "selecting waiter_work_time"), http.StatusBadRequest)
	}

	count := 0
	for countRows.Next() {
		if err = countRows.Scan(&count); err != nil {
			return nil, 0, web.NewRequestError(errors.Wrap(err, "scanning waiter_work_time count"), http.StatusBadRequest)
		}
	}

	return list, count, nil
}

func (r Repository) BranchGetListWaiterWorkTime(ctx context.Context, filter dto.BranchFilter) ([]dto.BranchGetListResponse, int, error) {
	claims, err := r.CheckClaims(ctx, auth.RoleBranch)
	if err != nil {
		return nil, 0, err
	}

	if err := r.ValidateStruct(&filter, "Date"); err != nil {
		return nil, 0, err
	}

	rand.Seed(time.Now().UnixNano())
	date := time.Now().Format("2006-01-02")
	if filter.Date != nil {
		Date, err := time.Parse("02.01.2006", *filter.Date)
		if err != nil {
			return nil, 0, web.NewRequestError(errors.New("incorrect date format in param"), http.StatusBadRequest)

		}
		date = Date.Format("2006-01-02")
	}

	whereQuery := fmt.Sprintf(`WHERE wwt.date = '%v' and u.branch_id = '%d' `, date, *claims.BranchID)

	orderQuery := "ORDER BY wwt.date desc"

	var limitQuery, offsetQuery string

	if filter.Page != nil && filter.Limit != nil {
		offset := (*filter.Page - 1) * (*filter.Limit)
		filter.Offset = &offset
	}

	if filter.Limit != nil {
		limitQuery += fmt.Sprintf(" LIMIT %d", *filter.Limit)
	}

	if filter.Offset != nil {
		offsetQuery += fmt.Sprintf(" OFFSET %d", *filter.Offset)
	}

	query := fmt.Sprintf(`
					SELECT 
					  wwt.id,
					  wwt.periods,
					  wwt.waiter_id,
					  u.name,
					  u.photo					  
					FROM 
					    waiter_work_time wwt
					LEFT JOIN users u on wwt.waiter_id = u.id and u.branch_id = '%d'
					%s %s %s %s
	`, *claims.BranchID, whereQuery, orderQuery, limitQuery, offsetQuery)

	rows, err := r.QueryContext(ctx, query)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, 0, web.NewRequestError(postgres.ErrNotFound, http.StatusBadRequest)
	}
	if err != nil {
		return nil, 0, web.NewRequestError(errors.Wrap(err, "selecting waiter_work_time"), http.StatusInternalServerError)
	}

	var list []dto.BranchGetListResponse

	for rows.Next() {
		var detail dto.BranchGetListResponse
		var periods []byte
		if err = rows.Scan(
			&detail.ID,
			&periods,
			&detail.WaiterID,
			&detail.Waiter,
			&detail.Avatar); err != nil {
			return nil, 0, web.NewRequestError(errors.Wrap(err, "scanning waiter_work_time"), http.StatusInternalServerError)
		}

		err = json.Unmarshal(periods, &detail.Periods)
		if err != nil {
			return nil, 0, web.NewRequestError(errors.Wrap(err, "unmarshal periods"), http.StatusInternalServerError)
		}
		totalDuration := 0
		for _, v := range detail.Periods {
			if v.Duration != nil {
				totalDuration += *v.Duration
			} else {
				if v.ComeTime != nil && v.GoneTime == nil {
					t := time.Now()
					if date != time.Now().Format("2006-01-02") {
						t, _ = time.Parse("2006-01-02", date)
						t = t.Add(24 * time.Hour).Add(-1 * time.Second)
					}

					d := int(t.Sub(*v.ComeTime).Minutes())

					v.GoneTime = &t
					v.Duration = &d
					detail.Now = true
					totalDuration += d
				} else if v.ComeTime == nil && v.GoneTime != nil {
					t, _ := time.Parse("2006-01-02", date)
					d := int(v.GoneTime.Sub(t).Minutes())
					v.ComeTime = &t
					v.Duration = &d
					totalDuration += d
				} else {
					if v.ComeTime != nil {
						d := int(v.GoneTime.Sub(*v.ComeTime).Minutes())
						v.Duration = &d
						totalDuration += d
					}
				}
			}
			//enter := v.ComeTime.Format("15:04:05")
			//exit := v.GoneTime.Format("15:04:05")
			enterInMinute := 0
			if v.ComeTime != nil {
				enterInMinute = v.ComeTime.Hour()*60 + v.ComeTime.Minute()
			}
			exitInMinute := 0
			if v.GoneTime != nil {
				exitInMinute = v.GoneTime.Hour()*60 + v.GoneTime.Minute()
			}

			detail.PeriodsResponse = append(detail.PeriodsResponse, dto.PeriodResponse{
				ComeTime: &enterInMinute,
				GoneTime: &exitInMinute,
				Duration: v.Duration,
			})
		}

		if detail.Avatar != nil {
			link := hashing.GenerateHash(r.ServerBaseUrl, *detail.Avatar)
			detail.Avatar = &link
		}

		if len(detail.Periods) > 0 {
			lastPeriod := detail.Periods[len(detail.Periods)-1]
			enterInMinute := 0

			if lastPeriod.ComeTime != nil {
				enterInMinute = lastPeriod.ComeTime.Hour()*60 + lastPeriod.ComeTime.Minute()
			}

			detail.StartTime = &enterInMinute
			detail.Duration = lastPeriod.Duration
		}

		detail.TotalDuration = &totalDuration
		list = append(list, detail)
	}

	countQuery := fmt.Sprintf(`
		SELECT
			count(wwt.id)
		FROM
		    waiter_work_time wwt				    
			LEFT JOIN users u on wwt.waiter_id = u.id and u.branch_id = '%d'
		%s
	`, *claims.BranchID, whereQuery)

	countRows, err := r.QueryContext(ctx, countQuery)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, 0, web.NewRequestError(postgres.ErrNotFound, http.StatusBadRequest)
	}
	if err != nil {
		return nil, 0, web.NewRequestError(errors.Wrap(err, "selecting waiter_work_time"), http.StatusInternalServerError)
	}

	count := 0
	for countRows.Next() {
		if err = countRows.Scan(&count); err != nil {
			return nil, 0, web.NewRequestError(errors.Wrap(err, "scanning waiter_work_time count"), http.StatusBadRequest)
		}
	}

	return list, count, nil
}

func (r Repository) CashierGetListWaiterWorkTime(ctx context.Context, filter dto.BranchFilter) ([]dto.BranchGetListResponse, int, error) {
	claims, err := r.CheckClaims(ctx, auth.RoleCashier)
	if err != nil {
		return nil, 0, err
	}

	if err := r.ValidateStruct(&filter, "Date"); err != nil {
		return nil, 0, err
	}

	rand.Seed(time.Now().UnixNano())
	date := time.Now().Format("2006-01-02")
	if filter.Date != nil {
		Date, err := time.Parse("02.01.2006", *filter.Date)
		if err != nil {
			return nil, 0, web.NewRequestError(errors.New("incorrect date format in param"), http.StatusBadRequest)

		}
		date = Date.Format("2006-01-02")
	}

	whereQuery := fmt.Sprintf(`WHERE wwt.date = '%v' and u.branch_id = '%d' `, date, *claims.BranchID)

	orderQuery := "ORDER BY wwt.date desc"

	var limitQuery, offsetQuery string

	if filter.Page != nil && filter.Limit != nil {
		offset := (*filter.Page - 1) * (*filter.Limit)
		filter.Offset = &offset
	}

	if filter.Limit != nil {
		limitQuery += fmt.Sprintf(" LIMIT %d", *filter.Limit)
	}

	if filter.Offset != nil {
		offsetQuery += fmt.Sprintf(" OFFSET %d", *filter.Offset)
	}

	query := fmt.Sprintf(`
					SELECT 
					  wwt.id,
					  wwt.periods,
					  wwt.waiter_id,
					  u.name,
					  u.photo
					FROM 
					    waiter_work_time wwt
					LEFT JOIN users u on wwt.waiter_id = u.id and u.branch_id = '%d'
					%s %s %s %s
	`, *claims.BranchID, whereQuery, orderQuery, limitQuery, offsetQuery)

	rows, err := r.QueryContext(ctx, query)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, 0, web.NewRequestError(postgres.ErrNotFound, http.StatusBadRequest)
	}
	if err != nil {
		return nil, 0, web.NewRequestError(errors.Wrap(err, "selecting waiter_work_time"), http.StatusInternalServerError)
	}

	var list []dto.BranchGetListResponse

	for rows.Next() {
		var detail dto.BranchGetListResponse
		var periods []byte
		if err = rows.Scan(
			&detail.ID,
			&periods,
			&detail.WaiterID,
			&detail.Waiter,
			&detail.Avatar); err != nil {
			return nil, 0, web.NewRequestError(errors.Wrap(err, "scanning waiter_work_time"), http.StatusInternalServerError)
		}

		err = json.Unmarshal(periods, &detail.Periods)
		if err != nil {
			return nil, 0, web.NewRequestError(errors.Wrap(err, "unmarshal periods"), http.StatusBadRequest)
		}
		totalDuration := 0
		for _, v := range detail.Periods {
			if v.Duration != nil {
				totalDuration += *v.Duration
			} else {
				if v.ComeTime != nil && v.GoneTime == nil {
					t := time.Now()
					if date != time.Now().Format("2006-01-02") {
						t, _ = time.Parse("2006-01-02", date)
						t = t.Add(24 * time.Hour).Add(-1 * time.Second)
					}

					d := int(t.Sub(*v.ComeTime).Minutes())
					v.GoneTime = &t
					v.Duration = &d
					detail.Now = true
					totalDuration += d
				} else if v.ComeTime == nil && v.GoneTime != nil {
					t, _ := time.Parse("2006-01-02", date)
					d := int(v.GoneTime.Sub(t).Minutes())
					v.ComeTime = &t
					v.Duration = &d
					totalDuration += d
				} else {
					if v.ComeTime != nil {
						d := int(v.GoneTime.Sub(*v.ComeTime).Minutes())
						v.Duration = &d
						totalDuration += d
					}
				}
			}
			//enter := v.ComeTime.Format("15:04:05")
			//exit := v.GoneTime.Format("15:04:05")
			enterInMinute := 0
			if v.ComeTime != nil {
				enterInMinute = v.ComeTime.Hour()*60 + v.ComeTime.Minute()
			}
			exitInMinute := 0
			if v.GoneTime != nil {
				exitInMinute = v.GoneTime.Hour()*60 + v.GoneTime.Minute()
			}
			detail.PeriodsResponse = append(detail.PeriodsResponse, dto.PeriodResponse{
				ComeTime: &enterInMinute,
				GoneTime: &exitInMinute,
				Duration: v.Duration,
			})
		}

		if detail.Avatar != nil {
			link := hashing.GenerateHash(r.ServerBaseUrl, *detail.Avatar)
			detail.Avatar = &link
		}

		if len(detail.Periods) > 0 {
			lastPeriod := detail.Periods[len(detail.Periods)-1]
			enterInMinute := 0

			if lastPeriod.ComeTime != nil {
				enterInMinute = lastPeriod.ComeTime.Hour()*60 + lastPeriod.ComeTime.Minute()
			}

			detail.StartTime = &enterInMinute
			detail.Duration = lastPeriod.Duration
		}

		detail.TotalDuration = &totalDuration
		list = append(list, detail)
	}

	countQuery := fmt.Sprintf(`
		SELECT
			count(wwt.id)
		FROM
		    waiter_work_time wwt				    
			LEFT JOIN users u on wwt.waiter_id = u.id and u.branch_id = '%d'
		%s
	`, *claims.BranchID, whereQuery)

	countRows, err := r.QueryContext(ctx, countQuery)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, 0, web.NewRequestError(postgres.ErrNotFound, http.StatusBadRequest)
	}
	if err != nil {
		return nil, 0, web.NewRequestError(errors.Wrap(err, "selecting waiter_work_time"), http.StatusBadRequest)
	}

	count := 0
	for countRows.Next() {
		if err = countRows.Scan(&count); err != nil {
			return nil, 0, web.NewRequestError(errors.Wrap(err, "scanning waiter_work_time count"), http.StatusBadRequest)
		}
	}

	return list, count, nil
}

func (r Repository) BranchGetDetailWaiterWorkTime(ctx context.Context, filter dto.ListFilter) ([]dto.GetListResponse, int, error) {
	claims, err := r.CheckClaims(ctx, auth.RoleBranch)
	if err != nil {
		return nil, 0, err
	}

	if err := r.ValidateStruct(&filter, "WaiterID"); err != nil {
		return nil, 0, err
	}

	whereQuery := fmt.Sprintf(`WHERE u.branch_id = '%d' and wwt.waiter_id = '%d'`, *claims.BranchID, *filter.WaiterID)

	orderQuery := "ORDER BY wwt.date desc"

	var limitQuery, offsetQuery string

	if filter.Page != nil && filter.Limit != nil {
		offset := (*filter.Page - 1) * (*filter.Limit)
		filter.Offset = &offset
	}

	if filter.Limit != nil {
		limitQuery += fmt.Sprintf(" LIMIT %d", *filter.Limit)
	}

	if filter.Offset != nil {
		offsetQuery += fmt.Sprintf(" OFFSET %d", *filter.Offset)
	}

	rand.Seed(time.Now().UnixNano())
	date := time.Now().Format("2006-01-02")

	query := fmt.Sprintf(`
					SELECT 
					  wwt.id,
					  TO_CHAR(wwt.date,'DD.MM.YYYY'),
					  wwt.periods
					FROM 
					    waiter_work_time wwt
					LEFT JOIN users u on wwt.waiter_id = u.id and u.branch_id = '%d'
					%s %s %s %s
	`, *claims.BranchID, whereQuery, orderQuery, limitQuery, offsetQuery)

	rows, err := r.QueryContext(ctx, query)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, 0, web.NewRequestError(postgres.ErrNotFound, http.StatusBadRequest)
	}
	if err != nil {
		return nil, 0, web.NewRequestError(errors.Wrap(err, "selecting waiter_work_time"), http.StatusInternalServerError)
	}

	var list []dto.GetListResponse

	for rows.Next() {
		var detail dto.GetListResponse
		var periods []byte
		if err = rows.Scan(
			&detail.ID,
			&detail.Date,
			&periods); err != nil {
			return nil, 0, web.NewRequestError(errors.Wrap(err, "scanning waiter_work_time"), http.StatusInternalServerError)
		}

		err = json.Unmarshal(periods, &detail.Periods)
		if err != nil {
			return nil, 0, web.NewRequestError(errors.Wrap(err, "unmarshal periods"), http.StatusBadRequest)
		}
		totalDuration := 0
		for _, v := range detail.Periods {
			if v.Duration != nil {
				totalDuration += *v.Duration
			} else {
				if v.ComeTime != nil && v.GoneTime == nil {
					t := time.Now()
					if date != time.Now().Format("2006-01-02") {
						t, _ = time.Parse("2006-01-02", date)
						t = t.Add(24 * time.Hour).Add(-1 * time.Second)
					}

					d := int(t.Sub(*v.ComeTime).Minutes())

					v.GoneTime = &t
					v.Duration = &d
					detail.Now = true
					totalDuration += d
				} else if v.ComeTime == nil && v.GoneTime != nil {
					t, _ := time.Parse("2006-01-02", date)
					d := int(v.GoneTime.Sub(t).Minutes())
					v.ComeTime = &t
					v.Duration = &d
					totalDuration += d
				} else {
					if v.ComeTime != nil {
						d := int(v.GoneTime.Sub(*v.ComeTime).Minutes())
						v.Duration = &d
						totalDuration += d
					}
				}
			}
			//enter := v.ComeTime.Format("15:04:05")
			//exit := v.GoneTime.Format("15:04:05")
			enterInMinute := 0
			if v.ComeTime != nil {
				enterInMinute = v.ComeTime.Hour()*60 + v.ComeTime.Minute()
			}
			exitInMinute := 0
			if v.GoneTime != nil {
				exitInMinute = v.GoneTime.Hour()*60 + v.GoneTime.Minute()
			}
			detail.PeriodsResponse = append(detail.PeriodsResponse, dto.PeriodResponse{
				ComeTime: &enterInMinute,
				GoneTime: &exitInMinute,
				Duration: v.Duration,
			})
		}

		detail.TotalDuration = &totalDuration
		list = append(list, detail)
	}

	countQuery := fmt.Sprintf(`
		SELECT
			count(wwt.id)
		FROM
		    waiter_work_time wwt					
			LEFT JOIN users u on wwt.waiter_id = u.id and u.branch_id = '%d'
		%s
	`, *claims.BranchID, whereQuery)

	countRows, err := r.QueryContext(ctx, countQuery)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, 0, web.NewRequestError(postgres.ErrNotFound, http.StatusBadRequest)
	}
	if err != nil {
		return nil, 0, web.NewRequestError(errors.Wrap(err, "selecting waiter_work_time"), http.StatusBadRequest)
	}

	count := 0
	for countRows.Next() {
		if err = countRows.Scan(&count); err != nil {
			return nil, 0, web.NewRequestError(errors.Wrap(err, "scanning waiter_work_time count"), http.StatusBadRequest)
		}
	}

	return list, count, nil
}

func (r Repository) CashierGetDetailWaiterWorkTime(ctx context.Context, filter dto.ListFilter) ([]dto.GetListResponse, int, error) {
	claims, err := r.CheckClaims(ctx, auth.RoleCashier)
	if err != nil {
		return nil, 0, err
	}

	whereQuery := fmt.Sprintf(`WHERE u.branch_id = '%d' and wwt.waiter_id = '%d'`, *claims.BranchID, *filter.WaiterID)

	orderQuery := "ORDER BY wwt.date desc"

	var limitQuery, offsetQuery string

	if filter.Page != nil && filter.Limit != nil {
		offset := (*filter.Page - 1) * (*filter.Limit)
		filter.Offset = &offset
	}

	if filter.Limit != nil {
		limitQuery += fmt.Sprintf(" LIMIT %d", *filter.Limit)
	}

	if filter.Offset != nil {
		offsetQuery += fmt.Sprintf(" OFFSET %d", *filter.Offset)
	}

	rand.Seed(time.Now().UnixNano())
	date := time.Now().Format("2006-01-02")

	query := fmt.Sprintf(`
					SELECT 
					  wwt.id,
					  TO_CHAR(wwt.date,'DD.MM.YYYY'),
					  wwt.periods
					FROM 
					    waiter_work_time wwt
					LEFT JOIN users u on wwt.waiter_id = u.id and u.branch_id = '%d'
					%s %s %s %s
	`, *claims.BranchID, whereQuery, orderQuery, limitQuery, offsetQuery)

	rows, err := r.QueryContext(ctx, query)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, 0, web.NewRequestError(postgres.ErrNotFound, http.StatusBadRequest)
	}
	if err != nil {
		return nil, 0, web.NewRequestError(errors.Wrap(err, "selecting waiter_work_time"), http.StatusInternalServerError)
	}

	var list []dto.GetListResponse

	for rows.Next() {
		var detail dto.GetListResponse
		var periods []byte
		if err = rows.Scan(
			&detail.ID,
			&detail.Date,
			&periods); err != nil {
			return nil, 0, web.NewRequestError(errors.Wrap(err, "scanning waiter_work_time"), http.StatusInternalServerError)
		}

		err = json.Unmarshal(periods, &detail.Periods)
		if err != nil {
			return nil, 0, web.NewRequestError(errors.Wrap(err, "unmarshal periods"), http.StatusBadRequest)
		}
		totalDuration := 0
		for _, v := range detail.Periods {
			if v.Duration != nil {
				totalDuration += *v.Duration
			} else {
				if v.ComeTime != nil && v.GoneTime == nil {
					t := time.Now()
					if date != time.Now().Format("2006-01-02") {
						t, _ = time.Parse("2006-01-02", date)
						t = t.Add(24 * time.Hour).Add(-1 * time.Second)
					}

					d := int(t.Sub(*v.ComeTime).Minutes())

					v.GoneTime = &t
					v.Duration = &d
					detail.Now = true
					totalDuration += d
				} else if v.ComeTime == nil && v.GoneTime != nil {
					t, _ := time.Parse("2006-01-02", date)
					d := int(v.GoneTime.Sub(t).Minutes())
					v.ComeTime = &t
					v.Duration = &d
					totalDuration += d
				} else {
					if v.ComeTime != nil {
						d := int(v.GoneTime.Sub(*v.ComeTime).Minutes())
						v.Duration = &d
						totalDuration += d
					}
				}
			}
			//enter := v.ComeTime.Format("15:04:05")
			//exit := v.GoneTime.Format("15:04:05")
			enterInMinute := 0
			if v.ComeTime != nil {
				enterInMinute = v.ComeTime.Hour()*60 + v.ComeTime.Minute()
			}
			exitInMinute := 0
			if v.GoneTime != nil {
				exitInMinute = v.GoneTime.Hour()*60 + v.GoneTime.Minute()
			}
			detail.PeriodsResponse = append(detail.PeriodsResponse, dto.PeriodResponse{
				ComeTime: &enterInMinute,
				GoneTime: &exitInMinute,
				Duration: v.Duration,
			})
		}

		detail.TotalDuration = &totalDuration
		list = append(list, detail)
	}

	countQuery := fmt.Sprintf(`
		SELECT
			count(wwt.id)
		FROM
		 waiter_work_time wwt
		LEFT JOIN users u on wwt.waiter_id = u.id and u.branch_id = '%d'
		%s
	`, *claims.BranchID, whereQuery)

	countRows, err := r.QueryContext(ctx, countQuery)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, 0, web.NewRequestError(postgres.ErrNotFound, http.StatusBadRequest)
	}
	if err != nil {
		return nil, 0, web.NewRequestError(errors.Wrap(err, "selecting waiter_work_time"), http.StatusBadRequest)
	}

	count := 0
	for countRows.Next() {
		if err = countRows.Scan(&count); err != nil {
			return nil, 0, web.NewRequestError(errors.Wrap(err, "scanning waiter_work_time count"), http.StatusBadRequest)
		}
	}

	return list, count, nil
}
