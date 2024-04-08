package printers

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
	"restu-backend/internal/service/printers"
	"time"
)

type Repository struct {
	*postgresql.Database
}

// @branch

func (r Repository) BranchGetList(ctx context.Context, filter printers.Filter) ([]printers.BranchGetList, int, error) {
	claims, err := r.CheckClaims(ctx, auth.RoleBranch)
	if err != nil {
		return nil, 0, err
	}

	whereQuery := fmt.Sprintf(`WHERE p.deleted_at IS NULL AND p.branch_id = %d`, *claims.BranchID)

	var limitQuery, offsetQuery string
	if filter.Limit != nil {
		limitQuery = fmt.Sprintf(" LIMIT '%d'", *filter.Limit)
	}
	if filter.Offset != nil {
		offsetQuery = fmt.Sprintf(" OFFSET '%d'", *filter.Offset)
	}

	query := fmt.Sprintf(`
					SELECT 
					    p.id, 
					    p.name, 
					    p.ip, 
					    p.port
					FROM 
					    printers as p
					%s %s %s`, whereQuery, limitQuery, offsetQuery)

	list := make([]printers.BranchGetList, 0)

	rows, err := r.QueryContext(ctx, query)
	if err != nil {
		return nil, 0, web.NewRequestError(errors.Wrap(err, "select printers"), http.StatusInternalServerError)
	}

	err = r.ScanRows(ctx, rows, &list)
	if err != nil {
		return nil, 0, web.NewRequestError(errors.Wrap(err, "scanning printers"), http.StatusBadRequest)
	}

	countQuery := fmt.Sprintf(`
		SELECT
			count(p.id)
		FROM
		    printers as p
		%s
	`, whereQuery)
	countRows, err := r.QueryContext(ctx, countQuery)
	if err == sql.ErrNoRows {
		return list, 0, nil
	}
	if err != nil {
		return nil, 0, web.NewRequestError(errors.Wrap(err, "selecting printers"), http.StatusBadRequest)
	}

	count := 0

	for countRows.Next() {
		if err = countRows.Scan(&count); err != nil {
			return nil, 0, web.NewRequestError(errors.Wrap(err, "scanning printers count"), http.StatusBadRequest)
		}
	}

	return list, count, nil
}

func (r Repository) BranchGetDetail(ctx context.Context, id int64) (printers.BranchGetDetail, error) {
	_, err := r.CheckClaims(ctx, auth.RoleBranch)
	if err != nil {
		return printers.BranchGetDetail{}, err
	}

	query := fmt.Sprintf(`
					SELECT 
					    id, 
						name, 
					    ip, 
					    port
					FROM 
					    printers
					WHERE deleted_at IS NULL AND id = %d
					`, id)

	var detail printers.BranchGetDetail

	err = r.QueryRowContext(ctx, query).Scan(
		&detail.ID,
		&detail.Name,
		&detail.IP,
		&detail.Port)

	if err == sql.ErrNoRows {
		return printers.BranchGetDetail{}, web.NewRequestError(postgres.ErrNotFound, http.StatusNotFound)
	}
	if err != nil {
		return printers.BranchGetDetail{}, web.NewRequestError(errors.Wrap(err, "selecting printer detail"), http.StatusBadRequest)
	}

	return detail, nil
}

func (r Repository) BranchCreate(ctx context.Context, request printers.BranchCreateRequest) (printers.BranchCreateResponse, error) {
	claims, err := r.CheckClaims(ctx, auth.RoleBranch)
	if err != nil {
		return printers.BranchCreateResponse{}, err
	}

	err = r.ValidateStruct(&request, "Name", "IP", "Port")
	if err != nil {
		return printers.BranchCreateResponse{}, err
	}

	response := printers.BranchCreateResponse{
		Name:      request.Name,
		IP:        request.IP,
		Port:      request.Port,
		BranchID:  claims.BranchID,
		CreatedAt: time.Now(),
		CreatedBy: claims.UserId,
	}

	_, err = r.NewInsert().Model(&response).Exec(ctx)
	if err != nil {
		return printers.BranchCreateResponse{}, web.NewRequestError(errors.Wrap(err, "creating printers"), http.StatusBadRequest)
	}

	return response, nil
}

func (r Repository) BranchUpdateAll(ctx context.Context, request printers.BranchUpdateRequest) error {
	claims, err := r.CheckClaims(ctx, auth.RoleBranch)
	if err != nil {
		return err
	}

	if err = r.ValidateStruct(&request, "ID", "Name", "IP", "Port"); err != nil {
		return err
	}

	q := r.NewUpdate().Table("printers").Where("deleted_at IS NULL AND id = ?", request.ID)

	q.Set("name = ?", request.Name)
	q.Set("ip = ?", request.IP)
	q.Set("port = ?", request.Port)
	q.Set("updated_at = ?", time.Now())
	q.Set("updated_by = ?", claims.UserId)

	_, err = q.Exec(ctx)
	if err != nil {
		return web.NewRequestError(errors.Wrap(err, "updating printers"), http.StatusBadRequest)
	}

	return nil
}

func (r Repository) BranchUpdateColumns(ctx context.Context, request printers.BranchUpdateRequest) error {
	claims, err := r.CheckClaims(ctx, auth.RoleBranch)
	if err != nil {
		return err
	}

	if err = r.ValidateStruct(&request, "ID"); err != nil {
		return err
	}

	q := r.NewUpdate().Table("printers").Where("deleted_at IS NULL AND id = ?", request.ID)

	if request.Name != nil {
		q.Set("name = ?", request.Name)
	}
	if request.IP != nil {
		q.Set("ip = ?", request.IP)
	}
	if request.Port != nil {
		q.Set("port = ?", request.Port)
	}

	q.Set("updated_at = ?", time.Now())
	q.Set("updated_by = ?", claims.UserId)

	_, err = q.Exec(ctx)
	if err != nil {
		return web.NewRequestError(errors.Wrap(err, "updating printers"), http.StatusBadRequest)
	}

	return nil
}

func (r Repository) BranchDelete(ctx context.Context, id int64) error {
	return r.DeleteRow(ctx, "printers", id, auth.RoleBranch)
}

func NewRepository(DB *postgresql.Database) *Repository {
	return &Repository{DB}
}
