package branchReview

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/pkg/errors"
	"github.com/restaurant/foundation/web"
	"github.com/restaurant/internal/auth"
	"github.com/restaurant/internal/pkg/repository/postgresql"
	"github.com/restaurant/internal/repository/postgres"
	"github.com/restaurant/internal/service/branchReview"
	"net/http"
	"time"
)

type Repository struct {
	*postgresql.Database
}

// @client

func (r Repository) ClientGetList(ctx context.Context, filter branchReview.Filter) ([]branchReview.ClientGetList, int, error) {
	_, err := r.CheckClaims(ctx, auth.RoleClient)
	if err != nil {
		return nil, 0, err
	}

	whereQuery := fmt.Sprintf(`WHERE br.deleted_at IS NULL`)
	countWhereQuery := whereQuery

	var limitQuery, offsetQuery string
	if filter.Limit != nil {
		limitQuery = fmt.Sprintf(" LIMIT '%d'", *filter.Limit)
	}
	if filter.Offset != nil {
		offsetQuery = fmt.Sprintf(" OFFSET '%d'", *filter.Offset)
	}

	if filter.UserID != nil {
		whereQuery += fmt.Sprintf(" br.user_id '%d'", *filter.UserID)
	}

	if filter.BranchID != nil {
		whereQuery += fmt.Sprintf(" br.branch_id '%d'", *filter.BranchID)
	}

	whereQuery += fmt.Sprintf(" %s %s", limitQuery, offsetQuery)

	query := fmt.Sprintf(`
					SELECT 
					    br.id, 
					    br.point, 
					    br.comment, 
					    br.rate,
					    br.user_id,
					    u.name as user_name,
					    br.branch_id,
					    b.name as branch_name
					FROM 
					    branch_reviews as br
					LEFT OUTER JOIN users as u ON u.id = br.user_id
					LEFT OUTER JOIN branches as b ON b.id = br.branch_id
					%s`, whereQuery)

	list := make([]branchReview.ClientGetList, 0)

	rows, err := r.QueryContext(ctx, query)
	if err != nil {
		return nil, 0, web.NewRequestError(errors.Wrap(err, "select branchReviews"), http.StatusInternalServerError)
	}

	err = r.ScanRows(ctx, rows, &list)
	if err != nil {
		return nil, 0, web.NewRequestError(errors.Wrap(err, "scanning branchReviews"), http.StatusBadRequest)
	}

	countQuery := fmt.Sprintf(`
		SELECT
			count(id)
		FROM
		    branch_reviews as br
		%s
	`, countWhereQuery)
	countRows, err := r.QueryContext(ctx, countQuery)
	if err == sql.ErrNoRows {
		return nil, 0, web.NewRequestError(postgres.ErrNotFound, http.StatusNotFound)
	}
	if err != nil {
		return nil, 0, web.NewRequestError(errors.Wrap(err, "selecting branchReview"), http.StatusBadRequest)
	}

	count := 0

	for countRows.Next() {
		if err = countRows.Scan(&count); err != nil {
			return nil, 0, web.NewRequestError(errors.Wrap(err, "scanning branchReview count"), http.StatusBadRequest)
		}
	}

	return list, count, nil
}

func (r Repository) ClientGetDetail(ctx context.Context, id int64) (branchReview.ClientGetDetail, error) {
	_, err := r.CheckClaims(ctx, auth.RoleClient)
	if err != nil {
		return branchReview.ClientGetDetail{}, err
	}

	query := fmt.Sprintf(`
					SELECT 
					    br.id, 
					    br.point, 
					    br.comment, 
					    br.rate,
					    br.user_id,
					    u.name as user_name,
					    br.branch_id,
					    b.name as branch_name
					FROM 
					    branch_reviews as br
					LEFT OUTER JOIN users as u ON u.id = br.user_id
					LEFT OUTER JOIN branches as b ON b.id = br.branch_id
					WHERE br.deleted_at IS NULL AND br.id = %d
					`, id)

	var detail branchReview.ClientGetDetail

	err = r.QueryRowContext(ctx, query).Scan(
		&detail.ID,
		&detail.Point,
		&detail.Comment,
		&detail.Rate,
		&detail.UserID,
		&detail.UserName,
		&detail.BranchID,
		&detail.BranchName)

	if err == sql.ErrNoRows {
		return branchReview.ClientGetDetail{}, web.NewRequestError(postgres.ErrNotFound, http.StatusNotFound)
	}
	if err != nil {
		return branchReview.ClientGetDetail{}, web.NewRequestError(errors.Wrap(err, "selecting area detail"), http.StatusBadRequest)
	}

	return detail, nil
}

func (r Repository) ClientCreate(ctx context.Context, request branchReview.ClientCreateRequest) (branchReview.ClientCreateResponse, error) {
	claims, err := r.CheckClaims(ctx, auth.RoleClient)
	if err != nil {
		return branchReview.ClientCreateResponse{}, err
	}

	err = r.ValidateStruct(&request, "Point", "Rate", "Comment", "BranchID")
	if err != nil {
		return branchReview.ClientCreateResponse{}, err
	}

	response := branchReview.ClientCreateResponse{
		Point:     request.Point,
		Comment:   request.Comment,
		Rate:      request.Rate,
		UserID:    &claims.UserId,
		BranchID:  request.BranchID,
		CreatedAt: time.Now(),
		CreatedBy: claims.UserId,
	}

	_, err = r.NewInsert().Model(&response).Exec(ctx)
	if err != nil {
		return branchReview.ClientCreateResponse{}, web.NewRequestError(errors.Wrap(err, "creating branchReview"), http.StatusBadRequest)
	}

	return response, nil
}

func (r Repository) ClientUpdateAll(ctx context.Context, request branchReview.ClientUpdateRequest) error {
	claims, err := r.CheckClaims(ctx, auth.RoleClient)
	if err != nil {
		return err
	}

	if err = r.ValidateStruct(&request, "ID", "Point", "Rate", "Comment", "BranchID", "UserID"); err != nil {
		return err
	}

	q := r.NewUpdate().Table("branch_reviews").Where("deleted_at IS NULL AND id = ?", request.ID)

	q.Set("point = ?", request.Point)
	q.Set("comment = ?", request.Comment)
	q.Set("rate = ?", request.Rate)
	q.Set("user_id = ?", request.UserID)
	q.Set("branch_id = ?", request.BranchID)
	q.Set("updated_at = ?", time.Now())
	q.Set("updated_by = ?", claims.UserId)

	_, err = q.Exec(ctx)
	if err != nil {
		return web.NewRequestError(errors.Wrap(err, "updating branchReview"), http.StatusBadRequest)
	}

	return nil
}

func (r Repository) ClientUpdateColumns(ctx context.Context, request branchReview.ClientUpdateRequest) error {
	claims, err := r.CheckClaims(ctx, auth.RoleClient)
	if err != nil {
		return err
	}

	if err = r.ValidateStruct(&request, "ID"); err != nil {
		return err
	}

	q := r.NewUpdate().Table("branch_reviews").Where("deleted_at IS NULL AND id = ?", request.ID)

	if request.Point != nil {
		q.Set("point = ?", request.Point)
	}
	if request.Comment != nil {
		q.Set("comment = ?", request.Comment)
	}
	if request.Rate != nil {
		q.Set("rate = ?", request.Rate)
	}
	if request.UserID != nil {
		q.Set("user_id = ?", request.UserID)
	}
	if request.BranchID != nil {
		q.Set("branch_id = ?", request.BranchID)
	}

	q.Set("updated_at = ?", time.Now())
	q.Set("updated_by = ?", claims.UserId)

	_, err = q.Exec(ctx)
	if err != nil {
		return web.NewRequestError(errors.Wrap(err, "updating branchReview"), http.StatusBadRequest)
	}

	return nil
}

func (r Repository) ClientDelete(ctx context.Context, id int64) error {
	return r.DeleteRow(ctx, "branch_reviews", id, auth.RoleClient)
}

func NewRepository(DB *postgresql.Database) *Repository {
	return &Repository{DB}
}
