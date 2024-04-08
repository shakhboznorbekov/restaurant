package service_percentage

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/pkg/errors"
	"net/http"
	"restu-backend/foundation/web"
	"restu-backend/internal/auth"
	"restu-backend/internal/pkg/repository/postgresql"
	"restu-backend/internal/repository/postgres"
	"restu-backend/internal/service/service_percentage"
	"time"
)

type Repository struct {
	*postgresql.Database
}

func (r Repository) AdminGetList(ctx context.Context, filter service_percentage.Filter) ([]service_percentage.AdminGetList, int, error) {
	claims, err := r.CheckClaims(ctx, auth.RoleBranch)
	if err != nil {
		return nil, 0, err
	}

	table := "service_percentage"
	whereQuery := fmt.Sprintf(`WHERE %s.deleted_at IS NULL AND branch_id = '%d'`, table, *claims.BranchID)

	var limitQuery, offsetQuery string

	if filter.Limit != nil {
		limitQuery = fmt.Sprintf(" LIMIT '%d'", *filter.Limit)
	}
	if filter.Offset != nil {
		offsetQuery = fmt.Sprintf(" OFFSET '%d'", *filter.Offset)
	}

	query := fmt.Sprintf(`
						SELECT 
						    id, 
						    percent 
						FROM service_percentage
						%s %s %s`, whereQuery, limitQuery, offsetQuery)

	list := make([]service_percentage.AdminGetList, 0)

	rows, err := r.QueryContext(ctx, query)
	if err != nil {
		return nil, 0, web.NewRequestError(errors.Wrap(err, "select service_percentage"), http.StatusInternalServerError)
	}

	err = r.ScanRows(ctx, rows, &list)
	if err != nil {
		return nil, 0, web.NewRequestError(errors.Wrap(err, "scanning service_percentages"), http.StatusBadRequest)
	}

	countQuery := fmt.Sprintf(`
		SELECT
			count(id)
		FROM
		    service_percentage
		%s
	`, whereQuery)

	countRows, err := r.QueryContext(ctx, countQuery)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, 0, web.NewRequestError(postgres.ErrNotFound, http.StatusNotFound)
	}
	if err != nil {
		return nil, 0, web.NewRequestError(errors.Wrap(err, "selecting service_percentage"), http.StatusBadRequest)
	}

	count := 0

	for countRows.Next() {
		if err = countRows.Scan(&count); err != nil {
			return nil, 0, web.NewRequestError(errors.Wrap(err, "scanning service_percentage count"), http.StatusBadRequest)
		}
	}

	return list, count, nil
}

func (r Repository) AdminGetDetail(ctx context.Context, id int64) (*service_percentage.AdminGetDetail, error) {
	claims, err := r.CheckClaims(ctx, auth.RoleBranch)
	if err != nil {
		return nil, err
	}

	var detail service_percentage.AdminGetDetail

	query := fmt.Sprintf(`SELECT 
    									id,
    									percent
								 FROM service_percentage
								 WHERE deleted_at IS NULL AND id='%d' AND branch_id = '%d'`, id, *claims.BranchID)

	err = r.QueryRowContext(ctx, query).Scan(&detail.ID, &detail.Percent)
	if err != nil {
		return nil, err
	}

	return &detail, nil
}

func (r Repository) AdminCreate(ctx context.Context, request service_percentage.AdminCreateRequest) (service_percentage.AdminCreateResponse, error) {
	claims, err := r.CheckClaims(ctx, auth.RoleBranch)
	if err != nil {
		return service_percentage.AdminCreateResponse{}, err
	}

	err = r.ValidateStruct(&request, "Percent")
	if err != nil {
		return service_percentage.AdminCreateResponse{}, err
	}

	response := service_percentage.AdminCreateResponse{
		Percent:   request.Percent,
		CreatedAt: time.Now(),
		CreatedBy: claims.UserId,
		BranchID:  claims.BranchID,
	}

	_, err = r.NewInsert().Model(&response).Returning("id").Exec(ctx, &response.ID)
	if err != nil {
		return service_percentage.AdminCreateResponse{}, web.NewRequestError(errors.Wrap(err, "creating service_percentage"), http.StatusBadRequest)
	}

	return response, nil
}

func (r Repository) AdminUpdateAll(ctx context.Context, request service_percentage.AdminUpdateRequest) error {
	claims, err := r.CheckClaims(ctx, auth.RoleBranch)
	if err != nil {
		return err
	}

	if err = r.ValidateStruct(&request, "ID", "Percent"); err != nil {
		return err
	}

	q := r.NewUpdate().Table("service_percentage").Where("deleted_at IS NULL AND id = ? AND branch_id = ?", request.ID, *claims.BranchID)

	q.Set("percent = ?", request.Percent)

	_, err = q.Exec(ctx)
	if err != nil {
		return web.NewRequestError(errors.Wrap(err, "updating service_percentage"), http.StatusBadRequest)
	}

	return nil
}

func (r Repository) AdminUpdateBranchID(ctx context.Context, request service_percentage.AdminUpdateBranchRequest) error {

	err := r.ValidateStruct(&request, "ID", "BranchID")
	if err != nil {
		return err
	}

	q := r.NewUpdate().Table("service_percentage").Where("deleted_at IS NULL AND id = ? ", request.ID)

	q.Set("branch_id = ?", request.BranchID)

	_, err = q.Exec(ctx)
	if err != nil {
		return web.NewRequestError(errors.Wrap(err, "updating service_percentage"), http.StatusBadRequest)
	}

	return nil
}

func (r Repository) AdminDelete(ctx context.Context, id int64) error {
	return r.DeleteRow(ctx, "service_percentage", id, auth.RoleBranch)
}

// branch

func (r Repository) BranchCreate(ctx context.Context, request service_percentage.AdminCreateRequest) (service_percentage.AdminCreateResponse, error) {
	claims, err := r.CheckClaims(ctx, auth.RoleAdmin)
	if err != nil {
		return service_percentage.AdminCreateResponse{}, err
	}

	err = r.ValidateStruct(&request, "Percent")
	if err != nil {
		return service_percentage.AdminCreateResponse{}, err
	}

	response := service_percentage.AdminCreateResponse{
		Percent:   request.Percent,
		CreatedAt: time.Now(),
		CreatedBy: claims.UserId,
		BranchID:  claims.BranchID,
	}

	_, err = r.NewInsert().Model(&response).Returning("id").Exec(ctx, &response.ID)
	if err != nil {
		return service_percentage.AdminCreateResponse{}, web.NewRequestError(errors.Wrap(err, "creating service_percentage"), http.StatusBadRequest)
	}

	return response, nil
}

func NewRepository(DB *postgresql.Database) *Repository {
	return &Repository{DB}
}
