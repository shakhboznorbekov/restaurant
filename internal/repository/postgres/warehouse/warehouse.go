package warehouse

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/pkg/errors"
	"net/http"
	"restu-backend/foundation/web"
	"restu-backend/internal/auth"
	"restu-backend/internal/pkg/repository/postgresql"
	"restu-backend/internal/repository/postgres"
	"restu-backend/internal/service/warehouse"
	"strings"
	"time"
)

type Repository struct {
	*postgresql.Database
}

// #admin

func (r Repository) AdminGetList(ctx context.Context, filter warehouse.Filter) ([]warehouse.AdminGetList, int, error) {
	claims, err := r.CheckClaims(ctx, auth.RoleAdmin)
	if err != nil {
		return nil, 0, err
	}

	whereQuery := fmt.Sprintf(`WHERE w.deleted_at IS NULL AND b.deleted_at IS NULL AND b.restaurant_id = %d`, *claims.RestaurantID)

	var limitQuery, offsetQuery string
	if filter.Limit != nil {
		limitQuery = fmt.Sprintf(" LIMIT '%d'", *filter.Limit)
	}
	if filter.Offset != nil {
		offsetQuery = fmt.Sprintf(" OFFSET '%d'", *filter.Offset)
	}

	if filter.BranchID != nil {
		whereQuery += fmt.Sprintf(` AND b.id = %d`, *filter.BranchID)
	}

	if filter.Type != nil {
		whereQuery += fmt.Sprintf(` AND w.type = '%s'`, *filter.Type)
	}

	query := fmt.Sprintf(`
		SELECT
			w.id,
			w.name,
			w.location,
			w.type,
			b.id,
		    b.name
		FROM
		    warehouses as w
		LEFT OUTER JOIN branches as b ON b.id = w.branch_id
		%s %s %s
	`, whereQuery, limitQuery, offsetQuery)

	rows, err := r.QueryContext(ctx, query)
	if err != nil {
		return nil, 0, web.NewRequestError(errors.Wrap(err, "select user"), http.StatusInternalServerError)
	}

	var list []warehouse.AdminGetList

	for rows.Next() {
		var detail warehouse.AdminGetList
		locationByte := make([]byte, 0)
		if err = rows.Scan(&detail.ID, &detail.Name, &locationByte, &detail.Type, &detail.BranchID, &detail.Branch); err != nil {
			return nil, 0, web.NewRequestError(errors.Wrap(err, "scanning warehouse"), http.StatusBadRequest)
		}
		if len(locationByte) > 0 {
			if err = json.Unmarshal(locationByte, &detail.Location); err != nil {
				return nil, 0, web.NewRequestError(errors.Wrap(err, "parsing location"), http.StatusBadRequest)
			}
		}

		list = append(list, detail)
	}

	countQuery := fmt.Sprintf(`
		SELECT
			count(w.id)
		FROM warehouses as w
		LEFT OUTER JOIN branches as b ON b.id = w.branch_id
		%s
	`, whereQuery)

	countRows, err := r.QueryContext(ctx, countQuery)
	if err == sql.ErrNoRows {
		return nil, 0, web.NewRequestError(postgres.ErrNotFound, http.StatusNotFound)
	}
	if err != nil {
		return nil, 0, web.NewRequestError(errors.Wrap(err, "selecting users"), http.StatusBadRequest)
	}

	count := 0

	for countRows.Next() {
		if err = countRows.Scan(&count); err != nil {
			return nil, 0, web.NewRequestError(errors.Wrap(err, "scanning user count"), http.StatusBadRequest)
		}
	}

	return list, count, nil
}

func (r Repository) AdminGetDetail(ctx context.Context, id int64) (warehouse.AdminGetDetail, error) {
	claims, err := r.CheckClaims(ctx, auth.RoleAdmin)
	if err != nil {
		return warehouse.AdminGetDetail{}, err
	}

	whereQuery := fmt.Sprintf(`WHERE w.deleted_at IS NULL AND b.deleted_at IS NULL AND b.restaurant_id = %d AND w.id = %d`, *claims.RestaurantID, id)

	query := fmt.Sprintf(`
		SELECT
			w.id,
			w.name,
			w.location,
			w.type,
			w.branch_id
		FROM
		    warehouses as w
		LEFT OUTER JOIN branches as b ON b.id = w.branch_id
		%s 
	`, whereQuery)

	locationByte := make([]byte, 0)

	detail := warehouse.AdminGetDetail{}

	err = r.QueryRowContext(ctx, query).Scan(&detail.ID, &detail.Name, &locationByte, &detail.Type, &detail.BranchID)
	if err != nil {
		return warehouse.AdminGetDetail{}, web.NewRequestError(errors.Wrap(err, "select user"), http.StatusInternalServerError)
	}

	if len(locationByte) > 0 {
		if err = json.Unmarshal(locationByte, &detail.Location); err != nil {
			return warehouse.AdminGetDetail{}, web.NewRequestError(errors.Wrap(err, "parsing location"), http.StatusBadRequest)
		}
	}

	return detail, nil
}

func (r Repository) AdminCreate(ctx context.Context, request warehouse.AdminCreateRequest) (warehouse.AdminCreateResponse, error) {
	claims, err := r.CheckClaims(ctx, auth.RoleAdmin)
	if err != nil {
		return warehouse.AdminCreateResponse{}, err
	}

	err = r.ValidateStruct(&request, "Name", "Type", "BranchID")
	if err != nil {
		return warehouse.AdminCreateResponse{}, err
	}

	if request.Type != nil {
		if strings.ToUpper(*request.Type) != "REAL" && strings.ToUpper(*request.Type) != "VIRTUAL" {
			return warehouse.AdminCreateResponse{}, web.NewRequestError(errors.New("type is invalid"), http.StatusBadRequest)
		}
	} else {
		Type := "REAL"
		request.Type = &Type
	}

	response := warehouse.AdminCreateResponse{
		Name:      request.Name,
		Location:  request.Location,
		Type:      request.Type,
		BranchID:  request.BranchID,
		CreatedAt: time.Now(),
		CreatedBy: claims.UserId,
	}

	_, err = r.NewInsert().Model(&response).Exec(ctx)
	if err != nil {
		return warehouse.AdminCreateResponse{}, web.NewRequestError(errors.Wrap(err, "creating user"), http.StatusBadRequest)
	}

	return response, nil
}

func (r Repository) AdminUpdateAll(ctx context.Context, request warehouse.AdminUpdateRequest) error {
	claims, err := r.CheckClaims(ctx, auth.RoleAdmin)
	if err != nil {
		return err
	}

	err = r.ValidateStruct(&request, "ID", "Location", "Name", "Type", "BranchID")
	if err != nil {
		return err
	}

	if strings.ToUpper(*request.Type) != "REAL" && strings.ToUpper(*request.Type) != "VIRTUAL" {
		return web.NewRequestError(errors.New("type is invalid"), http.StatusBadRequest)
	}

	q := r.NewUpdate().Table("warehouses").Where("deleted_at IS NULL AND id = ?", request.ID)

	q.Set("name = ?", request.Name)
	q.Set("location = ?", request.Location)
	q.Set("type = ?", request.Type)
	q.Set("branch_id = ?", request.BranchID)
	q.Set("updated_at = ?", time.Now())
	q.Set("updated_by = ?", claims.UserId)

	_, err = q.Exec(ctx)
	if err != nil {
		return web.NewRequestError(errors.Wrap(err, "updating warehouse"), http.StatusBadRequest)
	}

	return nil
}

func (r Repository) AdminUpdateColumns(ctx context.Context, request warehouse.AdminUpdateRequest) error {
	claims, err := r.CheckClaims(ctx, auth.RoleAdmin)
	if err != nil {
		return err
	}

	err = r.ValidateStruct(&request)
	if err != nil {
		return err
	}

	if request.Type != nil {
		if strings.ToUpper(*request.Type) != "REAL" && strings.ToUpper(*request.Type) != "VIRTUAL" {
			return web.NewRequestError(errors.New("type is invalid"), http.StatusBadRequest)
		}
	}

	q := r.NewUpdate().Table("warehouses").Where("deleted_at IS NULL AND id = ?", request.ID)

	if request.Name != nil {
		q.Set("name = ?", request.Name)
	}
	if request.Location != nil {
		q.Set("location = ?", request.Location)
	}
	if request.BranchID != nil {
		q.Set("branch_id = ?", request.BranchID)
	}
	if request.Type != nil {
		q.Set("type = ?", request.Type)
	}

	q.Set("updated_at = ?", time.Now())
	q.Set("updated_by = ?", claims.UserId)

	_, err = q.Exec(ctx)
	if err != nil {
		return web.NewRequestError(errors.Wrap(err, "updating warehouse"), http.StatusBadRequest)
	}

	return nil
}

func (r Repository) AdminDelete(ctx context.Context, id int64) error {
	return r.DeleteRow(ctx, "warehouses", id, auth.RoleAdmin)
}

func NewRepository(DB *postgresql.Database) *Repository {
	return &Repository{DB}
}

// #branch

func (r Repository) BranchGetList(ctx context.Context, filter warehouse.Filter) ([]warehouse.BranchGetList, int, error) {
	claims, err := r.CheckClaims(ctx, auth.RoleBranch)
	if err != nil {
		return nil, 0, err
	}

	whereQuery := fmt.Sprintf(`WHERE w.deleted_at IS NULL AND b.deleted_at IS NULL AND b.id = %d`, *claims.BranchID)

	var limitQuery, offsetQuery string
	if filter.Limit != nil {
		limitQuery = fmt.Sprintf(" LIMIT '%d'", *filter.Limit)
	}
	if filter.Offset != nil {
		offsetQuery = fmt.Sprintf(" OFFSET '%d'", *filter.Offset)
	}

	if filter.Type != nil {
		whereQuery += fmt.Sprintf(` AND w.type = '%s'`, *filter.Type)
	}

	query := fmt.Sprintf(`
		SELECT
			w.id,
			w.name,
			w.location,
			w.type
		FROM
		    warehouses as w
		LEFT OUTER JOIN branches as b ON b.id = w.branch_id
		%s %s %s
	`, whereQuery, limitQuery, offsetQuery)

	rows, err := r.QueryContext(ctx, query)
	if err != nil {
		return nil, 0, web.NewRequestError(errors.Wrap(err, "select user"), http.StatusInternalServerError)
	}

	var list []warehouse.BranchGetList

	for rows.Next() {
		var detail warehouse.BranchGetList
		locationByte := make([]byte, 0)
		if err = rows.Scan(&detail.ID, &detail.Name, &locationByte, &detail.Type); err != nil {
			return nil, 0, web.NewRequestError(errors.Wrap(err, "scanning warehouse"), http.StatusBadRequest)
		}
		if len(locationByte) > 0 {
			if err = json.Unmarshal(locationByte, &detail.Location); err != nil {
				return nil, 0, web.NewRequestError(errors.Wrap(err, "parsing location"), http.StatusBadRequest)
			}
		}

		list = append(list, detail)
	}

	countQuery := fmt.Sprintf(`
		SELECT
			count(w.id)
		FROM warehouses as w
		LEFT OUTER JOIN branches as b ON b.id = w.branch_id
		%s
	`, whereQuery)

	countRows, err := r.QueryContext(ctx, countQuery)
	if err == sql.ErrNoRows {
		return nil, 0, web.NewRequestError(postgres.ErrNotFound, http.StatusNotFound)
	}
	if err != nil {
		return nil, 0, web.NewRequestError(errors.Wrap(err, "selecting users"), http.StatusBadRequest)
	}

	count := 0

	for countRows.Next() {
		if err = countRows.Scan(&count); err != nil {
			return nil, 0, web.NewRequestError(errors.Wrap(err, "scanning user count"), http.StatusBadRequest)
		}
	}

	return list, count, nil
}

func (r Repository) BranchGetDetail(ctx context.Context, id int64) (warehouse.BranchGetDetail, error) {
	claims, err := r.CheckClaims(ctx, auth.RoleBranch)
	if err != nil {
		return warehouse.BranchGetDetail{}, err
	}

	whereQuery := fmt.Sprintf(`WHERE w.deleted_at IS NULL AND b.deleted_at IS NULL AND b.id = %d AND w.id = %d`, *claims.BranchID, id)

	query := fmt.Sprintf(`
		SELECT
			w.id,
			w.name,
			w.location,
			w.type
		FROM
		    warehouses as w
		LEFT OUTER JOIN branches as b ON b.id = w.branch_id
		%s 
	`, whereQuery)

	locationByte := make([]byte, 0)

	detail := warehouse.BranchGetDetail{}

	err = r.QueryRowContext(ctx, query).Scan(&detail.ID, &detail.Name, &locationByte, &detail.Type)
	if err != nil {
		return warehouse.BranchGetDetail{}, web.NewRequestError(errors.Wrap(err, "select user"), http.StatusInternalServerError)
	}

	if len(locationByte) > 0 {
		if err = json.Unmarshal(locationByte, &detail.Location); err != nil {
			return warehouse.BranchGetDetail{}, web.NewRequestError(errors.Wrap(err, "parsing location"), http.StatusBadRequest)
		}
	}

	return detail, nil
}

func (r Repository) BranchCreate(ctx context.Context, request warehouse.BranchCreateRequest) (warehouse.BranchCreateResponse, error) {
	claims, err := r.CheckClaims(ctx, auth.RoleBranch)
	if err != nil {
		return warehouse.BranchCreateResponse{}, err
	}

	err = r.ValidateStruct(&request, "Name", "Type", "BranchID")
	if err != nil {
		return warehouse.BranchCreateResponse{}, err
	}

	if request.Type != nil {
		if strings.ToUpper(*request.Type) != "REAL" && strings.ToUpper(*request.Type) != "VIRTUAL" {
			return warehouse.BranchCreateResponse{}, web.NewRequestError(errors.New("type is invalid"), http.StatusBadRequest)
		}
	} else {
		Type := "REAL"
		request.Type = &Type
	}

	response := warehouse.BranchCreateResponse{
		Name:      request.Name,
		Location:  request.Location,
		Type:      request.Type,
		BranchID:  claims.BranchID,
		CreatedAt: time.Now(),
		CreatedBy: claims.UserId,
	}

	_, err = r.NewInsert().Model(&response).Exec(ctx)
	if err != nil {
		return warehouse.BranchCreateResponse{}, web.NewRequestError(errors.Wrap(err, "creating user"), http.StatusBadRequest)
	}

	return response, nil
}

func (r Repository) BranchUpdateAll(ctx context.Context, request warehouse.BranchUpdateRequest) error {
	claims, err := r.CheckClaims(ctx, auth.RoleBranch)
	if err != nil {
		return err
	}

	err = r.ValidateStruct(&request, "ID", "Location", "Name", "Type")
	if err != nil {
		return err
	}

	if strings.ToUpper(*request.Type) != "REAL" && strings.ToUpper(*request.Type) != "VIRTUAL" {
		return web.NewRequestError(errors.New("type is invalid"), http.StatusBadRequest)
	}

	q := r.NewUpdate().Table("warehouses").Where("deleted_at IS NULL AND id = ?", request.ID)

	q.Set("name = ?", request.Name)
	q.Set("location = ?", request.Location)
	q.Set("type = ?", request.Type)
	q.Set("updated_at = ?", time.Now())
	q.Set("updated_by = ?", claims.UserId)

	_, err = q.Exec(ctx)
	if err != nil {
		return web.NewRequestError(errors.Wrap(err, "updating warehouse"), http.StatusBadRequest)
	}

	return nil
}

func (r Repository) BranchUpdateColumns(ctx context.Context, request warehouse.BranchUpdateRequest) error {
	claims, err := r.CheckClaims(ctx, auth.RoleBranch)
	if err != nil {
		return err
	}

	err = r.ValidateStruct(&request)
	if err != nil {
		return err
	}

	if request.Type != nil {
		if strings.ToUpper(*request.Type) != "REAL" && strings.ToUpper(*request.Type) != "VIRTUAL" {
			return web.NewRequestError(errors.New("type is invalid"), http.StatusBadRequest)
		}
	}

	q := r.NewUpdate().Table("warehouses").Where("deleted_at IS NULL AND id = ?", request.ID)

	if request.Name != nil {
		q.Set("name = ?", request.Name)
	}
	if request.Location != nil {
		q.Set("location = ?", request.Location)
	}
	if request.Type != nil {
		q.Set("type = ?", request.Type)
	}

	q.Set("updated_at = ?", time.Now())
	q.Set("updated_by = ?", claims.UserId)

	_, err = q.Exec(ctx)
	if err != nil {
		return web.NewRequestError(errors.Wrap(err, "updating warehouse"), http.StatusBadRequest)
	}

	return nil
}

func (r Repository) BranchDelete(ctx context.Context, id int64) error {
	return r.DeleteRow(ctx, "warehouses", id, auth.RoleBranch)
}
