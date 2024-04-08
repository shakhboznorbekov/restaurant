package food_price

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/pkg/errors"
	"net/http"
	"restu-backend/foundation/web"
	"restu-backend/internal/auth"
	"restu-backend/internal/entity"
	"restu-backend/internal/pkg/repository/postgresql"
	"restu-backend/internal/repository/postgres"
	"restu-backend/internal/service/food_price"
	"time"
)

type Repository struct {
	*postgresql.Database
}

// admin

func (r Repository) AdminGetList(ctx context.Context, filter food_price.Filter) ([]food_price.AdminGetList, int, error) {
	_, err := r.CheckClaims(ctx, auth.RoleAdmin)
	if err != nil {
		return nil, 0, err
	}

	table := "food_price"
	whereQuery := fmt.Sprintf(`WHERE %s.deleted_at IS NULL`, table)

	var limitQuery, offsetQuery string
	if filter.Limit != nil {
		limitQuery = fmt.Sprintf(" LIMIT '%d'", *filter.Limit)
	}
	if filter.Offset != nil {
		offsetQuery = fmt.Sprintf(" OFFSET '%d'", *filter.Offset)
	}

	query := fmt.Sprintf(`
			SELECT
				food_price.id,
				food_price.price,
				TO_CHAR(food_price.set_date, 'DD.MM.YYYY'),
				food_price.food_id
			FROM food_price
			%s %s %s
	`, whereQuery, limitQuery, offsetQuery)

	list := make([]food_price.AdminGetList, 0)

	rows, err := r.QueryContext(ctx, query)
	if err != nil {
		return nil, 0, web.NewRequestError(errors.Wrap(err, "select foods"), http.StatusInternalServerError)
	}

	for rows.Next() {
		detail := food_price.AdminGetList{}
		err = rows.Scan(
			&detail.ID,
			&detail.Price,
			&detail.SetDate,
			&detail.FoodID,
		)
		if err != nil {
			return nil, 0, web.NewRequestError(errors.Wrap(err, "scanning foods"), http.StatusBadRequest)
		}
		list = append(list, detail)
	}

	countQuery := fmt.Sprintf(`
		SELECT
			count(%s.id)
		FROM food_price
		%s
	`, table, whereQuery)

	countRows, err := r.QueryContext(ctx, countQuery)
	if err == sql.ErrNoRows {
		return nil, 0, web.NewRequestError(postgres.ErrNotFound, http.StatusNotFound)
	}
	if err != nil {
		return nil, 0, web.NewRequestError(errors.Wrap(err, "selecting foods"), http.StatusBadRequest)
	}

	count := 0

	for countRows.Next() {
		if err = countRows.Scan(&count); err != nil {
			return nil, 0, web.NewRequestError(errors.Wrap(err, "scanning user count"), http.StatusBadRequest)
		}
	}

	return list, count, nil
}

func (r Repository) AdminGetDetail(ctx context.Context, id int64) (entity.FoodPrice, error) {
	_, err := r.CheckClaims(ctx, auth.RoleAdmin)
	if err != nil {
		return entity.FoodPrice{}, err
	}

	var detail entity.FoodPrice

	err = r.NewSelect().Model(&detail).Where("id = ?", id).Scan(ctx)
	if err != nil {
		return entity.FoodPrice{}, err
	}
	return detail, nil
}

func (r Repository) AdminCreate(ctx context.Context, request food_price.AdminCreateRequest) (food_price.AdminCreateResponse, error) {
	claims, err := r.CheckClaims(ctx, auth.RoleAdmin)
	if err != nil {
		return food_price.AdminCreateResponse{}, err
	}

	err = r.ValidateStruct(&request, "Price", "SetDate", "MenuID")
	if err != nil {
		return food_price.AdminCreateResponse{}, err
	}

	setDate, err := time.Parse("02.01.2006", *request.SetDate)
	if err != nil {
		return food_price.AdminCreateResponse{}, errors.Wrap(err, "time parse")
	}

	response := food_price.AdminCreateResponse{
		Price:     request.Price,
		FoodID:    request.FoodID,
		SetDate:   &setDate,
		CreatedAt: time.Now(),
		CreatedBy: claims.UserId,
	}

	_, err = r.NewInsert().Model(&response).Exec(ctx)
	if err != nil {
		return food_price.AdminCreateResponse{}, web.NewRequestError(errors.Wrap(err, "creating food"), http.StatusBadRequest)
	}

	return response, nil
}

func (r Repository) AdminUpdateAll(ctx context.Context, request food_price.AdminUpdateRequest) error {
	claims, err := r.CheckClaims(ctx, auth.RoleAdmin)
	if err != nil {
		return err
	}

	if err = r.ValidateStruct(&request, "ID", "Price", "SetDate", "FoodID"); err != nil {
		return err
	}

	q := r.NewUpdate().Table("food_price").Where("deleted_at IS NULL AND id = ?", request.ID)

	setDate, err := time.Parse("02.01.2006", *request.SetDate)
	if err != nil {
		return errors.Wrap(err, "time parse")
	}

	q.Set("price = ?", request.Price)
	q.Set("set_date = ?", setDate)
	q.Set("food_id", request.FoodID)
	q.Set("updated_at = ?", time.Now())
	q.Set("updated_by = ?", claims.UserId)

	_, err = q.Exec(ctx)
	if err != nil {
		return web.NewRequestError(errors.Wrap(err, "updating food_price"), http.StatusBadRequest)
	}

	return nil
}

func (r Repository) AdminUpdateColumns(ctx context.Context, request food_price.AdminUpdateRequest) error {
	claims, err := r.CheckClaims(ctx, auth.RoleAdmin)
	if err != nil {
		return err
	}

	if err = r.ValidateStruct(&request, "ID"); err != nil {
		return err
	}

	q := r.NewUpdate().Table("food_price").Where("deleted_at IS NULL AND id = ?", request.ID)

	if request.Price != nil {
		q.Set("price = ?", request.Price)
	}
	if request.FoodID != nil {
		q.Set("food_id = ?", request.FoodID)
	}
	if request.SetDate != nil {
		setDate, err := time.Parse("02.01.2006", *request.SetDate)
		if err != nil {
			return errors.Wrap(err, "time parse")
		}
		q.Set("set_date = ?", setDate)
	}
	q.Set("updated_at = ?", time.Now())
	q.Set("updated_by = ?", claims.UserId)

	_, err = q.Exec(ctx)
	if err != nil {
		return web.NewRequestError(errors.Wrap(err, "updating food_price"), http.StatusBadRequest)
	}

	return nil
}

func (r Repository) AdminDelete(ctx context.Context, id int64) error {
	return r.DeleteRow(ctx, "food_price", id, auth.RoleAdmin)
}
