package attendance

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/pkg/errors"
	"github.com/restaurant/foundation/web"
	"github.com/restaurant/internal/auth"
	"github.com/restaurant/internal/pkg/repository/postgresql"
	"github.com/restaurant/internal/service/attendance"
	"net/http"
	"time"
)

type Repository struct {
	*postgresql.Database
}

func NewRepository(DB *postgresql.Database) *Repository {
	return &Repository{DB}
}

// @waiter

func (r Repository) WaiterCameCreate(ctx context.Context) (attendance.WaiterCreateResponse, error) {
	claims, err := r.CheckClaims(ctx, auth.RoleWaiter)
	if err != nil {
		return attendance.WaiterCreateResponse{}, err
	}

	var existingCameTime *time.Time
	err = r.QueryRowContext(ctx,
		fmt.Sprintf("SELECT came_at FROM attendances WHERE gone_at is null and user_id = '%d' order by came_at desc", claims.UserId)).Scan(&existingCameTime)

	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return attendance.WaiterCreateResponse{}, web.NewRequestError(errors.Wrap(err, "querying attendance"), http.StatusInternalServerError)
	}

	if existingCameTime != nil {
		return attendance.WaiterCreateResponse{}, web.NewRequestError(errors.New("you are already at work"), http.StatusBadRequest)
	}

	cameTime := time.Now()
	var response attendance.WaiterCreateResponse
	response.UserID = &claims.UserId
	response.CameAt = &cameTime

	_, err = r.NewInsert().Model(&response).Returning("id").Exec(ctx, &response.ID)
	if err != nil {
		return attendance.WaiterCreateResponse{}, web.NewRequestError(errors.Wrap(err, "creating waiter come attendance"), http.StatusInternalServerError)
	}

	q := r.NewUpdate().Table("users").Where("deleted_at is null and role='WAITER' and branch_id = ? and id = ?", *claims.BranchID, claims.UserId)

	q.Set("attendance_status = ?", true)

	q.Set("updated_at = ?", time.Now())
	q.Set("updated_by = ?", claims.UserId)

	_, err = q.Exec(ctx)
	if err != nil {
		return response, web.NewRequestError(errors.Wrap(err, "updating waiter attendance status"), http.StatusBadRequest)
	}

	return response, nil
}

func (r Repository) WaiterGoneCreate(ctx context.Context) (attendance.WaiterCreateResponse, error) {
	claims, err := r.CheckClaims(ctx, auth.RoleWaiter)
	if err != nil {
		return attendance.WaiterCreateResponse{}, err
	}

	var existingComeTime *time.Time
	var id int64
	err = r.QueryRowContext(ctx,
		fmt.Sprintf("SELECT id, came_at FROM attendances WHERE gone_at is null and user_id = '%d'  order by came_at desc", claims.UserId)).Scan(&id, &existingComeTime)

	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return attendance.WaiterCreateResponse{}, web.NewRequestError(errors.Wrap(err, "querying attendance"), http.StatusInternalServerError)
	}

	if existingComeTime == nil {
		return attendance.WaiterCreateResponse{}, web.NewRequestError(errors.New("you are not at work"), http.StatusBadRequest)
	}
	goneTime := time.Now()

	var response attendance.WaiterCreateResponse
	response.UserID = &claims.UserId
	response.GoneAt = &goneTime
	response.ID = id
	response.CameAt = existingComeTime

	_, err = r.NewUpdate().Table("attendances").
		Set("gone_at = ?", response.GoneAt).
		Where("user_id = ? and came_at is not null and id = ?", claims.UserId, id).
		Exec(ctx)

	if err != nil {
		return attendance.WaiterCreateResponse{}, web.NewRequestError(errors.Wrap(err, "updating waiter gone attendance"), http.StatusInternalServerError)
	}

	q := r.NewUpdate().Table("users").Where("deleted_at is null and role='WAITER' and branch_id = ? and id = ?", *claims.BranchID, claims.UserId)

	q.Set("attendance_status = ?", false)

	q.Set("updated_at = ?", time.Now())
	q.Set("updated_by = ?", claims.UserId)

	_, err = q.Exec(ctx)
	if err != nil {
		return response, web.NewRequestError(errors.Wrap(err, "updating waiter attendance status"), http.StatusBadRequest)
	}

	return response, nil
}
