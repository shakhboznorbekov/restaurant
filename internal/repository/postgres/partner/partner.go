package partner

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
	"restu-backend/internal/service/partner"
	"strings"
	"time"
)

type Repository struct {
	*postgresql.Database
}

// admin

func (r Repository) AdminGetList(ctx context.Context, filter partner.Filter) ([]partner.AdminGetList, int, error) {
	claims, err := r.CheckClaims(ctx, auth.RoleAdmin)
	if err != nil {
		return nil, 0, err
	}

	whereQuery := fmt.Sprintf(`WHERE p.deleted_at IS NULL AND p.restaurant_id = %d`, *claims.RestaurantID)

	var limitQuery, offsetQuery string
	if filter.Limit != nil {
		limitQuery = fmt.Sprintf(" LIMIT '%d'", *filter.Limit)
	}
	if filter.Offset != nil {
		offsetQuery = fmt.Sprintf(" OFFSET '%d'", *filter.Offset)
	}

	orderQuery := fmt.Sprintf(" ORDER BY p.name")

	if filter.Search != nil {
		search := strings.Replace(*filter.Search, " ", "", -1)
		whereQuery += fmt.Sprintf(` AND REPLACE(p.name, ' ', '') ilike '%s'`, "%"+search+"%")
	}

	query := fmt.Sprintf(`
		SELECT
		    p.id,
		    p.name,
		    p.type
		FROM partners AS p
		%s %s %s %s
`, whereQuery, orderQuery, limitQuery, offsetQuery)

	list := make([]partner.AdminGetList, 0)

	rows, err := r.QueryContext(ctx, query)
	if err != nil {
		return nil, 0, web.NewRequestError(errors.Wrap(err, "select partner"), http.StatusInternalServerError)
	}

	err = r.ScanRows(ctx, rows, &list)
	if err != nil {
		return nil, 0, web.NewRequestError(errors.Wrap(err, "scanning partners"), http.StatusBadRequest)
	}

	countQuery := fmt.Sprintf(`
		SELECT
			count(p.id)
		FROM partners AS p
		%s
	`, whereQuery)

	countRows, err := r.QueryContext(ctx, countQuery)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, 0, web.NewRequestError(postgres.ErrNotFound, http.StatusNotFound)
	}
	if err != nil {
		return nil, 0, web.NewRequestError(errors.Wrap(err, "selecting partners"), http.StatusBadRequest)
	}

	count := 0

	for countRows.Next() {
		if err = countRows.Scan(&count); err != nil {
			return nil, 0, web.NewRequestError(errors.Wrap(err, "scanning partners count"), http.StatusBadRequest)
		}
	}

	return list, count, nil
}

func (r Repository) AdminGetDetail(ctx context.Context, id int64) (partner.AdminGetDetail, error) {
	claims, err := r.CheckClaims(ctx, auth.RoleAdmin)
	if err != nil {
		return partner.AdminGetDetail{}, err
	}

	whereQuery := fmt.Sprintf(`WHERE p.deleted_at IS NULL AND p.restaurant_id = %d AND p.id = %d`, *claims.RestaurantID, id)

	query := fmt.Sprintf(`
		SELECT
		    p.id,
		    p.name,
		    p.type
		FROM partners AS p
		%s
`, whereQuery)

	detail := partner.AdminGetDetail{}

	rows, err := r.QueryContext(ctx, query)
	if err != nil {
		return partner.AdminGetDetail{}, web.NewRequestError(errors.Wrap(err, "select partner"), http.StatusInternalServerError)
	}

	err = r.ScanRows(ctx, rows, &detail)
	if err != nil {
		return partner.AdminGetDetail{}, web.NewRequestError(errors.Wrap(err, "scanning partners"), http.StatusBadRequest)
	}

	return detail, nil
}

func (r Repository) AdminCreate(ctx context.Context, request partner.AdminCreateRequest) (partner.AdminCreateResponse, error) {
	claims, err := r.CheckClaims(ctx, auth.RoleAdmin)
	if err != nil {
		return partner.AdminCreateResponse{}, err
	}

	err = r.ValidateStruct(&request, "Name", "Type", "WarehouseID")
	if err != nil {
		return partner.AdminCreateResponse{}, err
	}
	Type := strings.ToUpper(*request.Type)
	if Type != "COUNTER-AGENT" {
		return partner.AdminCreateResponse{}, web.NewRequestError(errors.New("no such type exists"), http.StatusBadRequest)
	}

	response := partner.AdminCreateResponse{
		Name:         request.Name,
		Type:         &Type,
		RestaurantID: claims.RestaurantID,
		CreatedAt:    time.Now(),
		CreatedBy:    claims.UserId,
	}

	_, err = r.NewInsert().Model(&response).Exec(ctx)
	if err != nil {
		return partner.AdminCreateResponse{}, web.NewRequestError(errors.Wrap(err, "creating partner"), http.StatusBadRequest)
	}
	return response, nil
}

func (r Repository) AdminUpdateAll(ctx context.Context, request partner.AdminUpdateRequest) error {
	claims, err := r.CheckClaims(ctx, auth.RoleAdmin)
	if err != nil {
		return err
	}

	if err = r.ValidateStruct(&request, "ID", "Name", "Type"); err != nil {
		return err
	}

	Type := strings.ToUpper(*request.Type)
	if Type != "COUNTER-AGENT" {
		return web.NewRequestError(errors.New("no such type exists"), http.StatusBadRequest)
	}

	q := r.NewUpdate().Table("partners").Where("deleted_at IS NULL AND id = ? AND restaurant_id = ?", request.ID, claims.RestaurantID)

	q.Set("name = ?", request.Name)
	q.Set("type =?", Type)
	q.Set("updated_at = ?", time.Now())
	q.Set("updated_by = ?", claims.UserId)

	_, err = q.Exec(ctx)
	if err != nil {
		return web.NewRequestError(errors.Wrap(err, "updating partner"), http.StatusBadRequest)
	}

	return nil
}

func (r Repository) AdminUpdateColumns(ctx context.Context, request partner.AdminUpdateRequest) error {
	claims, err := r.CheckClaims(ctx, auth.RoleAdmin)
	if err != nil {
		return err
	}

	if err = r.ValidateStruct(&request, "ID"); err != nil {
		return err
	}

	q := r.NewUpdate().Table("partners").Where("deleted_at IS NULL AND id = ? AND restaurant_id = ?", request.ID, claims.RestaurantID)

	if request.Name != nil {
		q.Set("name = ?", request.Name)
	}

	if request.Type != nil {
		Type := strings.ToUpper(*request.Type)
		if Type != "COUNTER-AGENT" {
			return web.NewRequestError(errors.New("no such type exists"), http.StatusBadRequest)
		}

		q.Set("type = ?", Type)
	}

	q.Set("updated_at = ?", time.Now())
	q.Set("updated_by = ?", claims.UserId)

	_, err = q.Exec(ctx)
	if err != nil {
		return web.NewRequestError(errors.Wrap(err, "updating partner"), http.StatusBadRequest)
	}

	return nil
}

func (r Repository) AdminDelete(ctx context.Context, id int64) error {
	return r.DeleteRow(ctx, "partners", id, auth.RoleAdmin)
}

// branch

func (r Repository) BranchGetList(ctx context.Context, filter partner.Filter) ([]partner.BranchGetList, int, error) {
	claims, err := r.CheckClaims(ctx, auth.RoleBranch)
	if err != nil {
		return nil, 0, err
	}

	whereQuery := fmt.Sprintf(`WHERE p.deleted_at IS NULL AND p.restaurant_id = b.restaurant_id`)

	var limitQuery, offsetQuery string
	if filter.Limit != nil {
		limitQuery = fmt.Sprintf(" LIMIT '%d'", *filter.Limit)
	}
	if filter.Offset != nil {
		offsetQuery = fmt.Sprintf(" OFFSET '%d'", *filter.Offset)
	}

	orderQuery := fmt.Sprintf(" ORDER BY p.name")

	if filter.Search != nil {
		search := strings.Replace(*filter.Search, " ", "", -1)
		whereQuery += fmt.Sprintf(` AND REPLACE(p.name, ' ', '') ilike '%s'`, "%"+search+"%")
	}

	query := fmt.Sprintf(`
		SELECT
		    p.id,
		    p.name,
		    p.type
		FROM partners AS p
		LEFT JOIN branches AS b ON b.id = %d
		%s %s %s %s
`, *claims.BranchID, whereQuery, orderQuery, limitQuery, offsetQuery)

	list := make([]partner.BranchGetList, 0)

	rows, err := r.QueryContext(ctx, query)
	if err != nil {
		return nil, 0, web.NewRequestError(errors.Wrap(err, "select partner"), http.StatusInternalServerError)
	}

	err = r.ScanRows(ctx, rows, &list)
	if err != nil {
		return nil, 0, web.NewRequestError(errors.Wrap(err, "scanning partners"), http.StatusBadRequest)
	}

	countQuery := fmt.Sprintf(`
		SELECT
			count(p.id)
		FROM partners AS p
		LEFT JOIN branches AS b ON b.id = %d
		%s
	`, *claims.BranchID, whereQuery)

	countRows, err := r.QueryContext(ctx, countQuery)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, 0, web.NewRequestError(postgres.ErrNotFound, http.StatusNotFound)
	}
	if err != nil {
		return nil, 0, web.NewRequestError(errors.Wrap(err, "selecting partners"), http.StatusBadRequest)
	}

	count := 0

	for countRows.Next() {
		if err = countRows.Scan(&count); err != nil {
			return nil, 0, web.NewRequestError(errors.Wrap(err, "scanning partners count"), http.StatusBadRequest)
		}
	}

	return list, count, nil
}
func NewRepository(DB *postgresql.Database) *Repository {
	return &Repository{DB}
}
