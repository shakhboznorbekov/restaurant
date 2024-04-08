package food_recipe

import (
	"database/sql"
	"fmt"
	"github.com/pkg/errors"
	"golang.org/x/net/context"
	"net/http"
	"restu-backend/foundation/web"
	"restu-backend/internal/auth"
	"restu-backend/internal/pkg/repository/postgresql"
	"restu-backend/internal/repository/postgres"
	"restu-backend/internal/service/food_recipe"
	"time"
)

type Repository struct {
	*postgresql.Database
}

func NewRepository(db *postgresql.Database) *Repository {
	return &Repository{db}
}

// @admin

func (r Repository) AdminGetList(ctx context.Context, filter food_recipe.Filter) ([]food_recipe.AdminGetList, int, error) {
	_, err := r.CheckClaims(ctx, auth.RoleAdmin)
	if err != nil {
		return nil, 0, err
	}

	whereQuery := fmt.Sprintf(`WHERE fr.deleted_at ISNULL AND f.deleted_at ISNULL AND p.deleted_at ISNULL`)

	if filter.Name != nil {
		whereQuery += fmt.Sprintf(" AND f.name ilike '%s'", "%"+*filter.Name+"%")
	}
	if filter.FoodId != nil {
		whereQuery += fmt.Sprintf(` AND fr.food_id='%d'`, *filter.FoodId)
	}

	countWhereQuery := whereQuery

	if filter.Limit != nil {
		whereQuery += fmt.Sprintf(" LIMIT '%d'", *filter.Limit)
	}
	if filter.Offset != nil {
		whereQuery += fmt.Sprintf(" OFFSET '%d'", *filter.Offset)
	}

	query := fmt.Sprintf(`SELECT 
    									fr.id as id,
    									f.name as food,
    									p.name as product,
    									fr.amount as amount,
    									m.name as measure_unit
								 FROM food_recipe fr 
								     JOIN foods f 
								         ON fr.food_id = f.id 
								     JOIN products p 
								         ON fr.product_id = p.id 
								     JOIN measure_unit m
								         ON p.measure_unit_id = m.id
									%s`, whereQuery)

	list := make([]food_recipe.AdminGetList, 0)

	rows, err := r.QueryContext(ctx, query)
	if err != nil {
		return nil, 0, web.NewRequestError(errors.Wrap(err, "select food recipe"), http.StatusInternalServerError)
	}

	err = r.ScanRows(ctx, rows, &list)
	if err != nil {
		return nil, 0, web.NewRequestError(errors.Wrap(err, "scanning food recipe"), http.StatusBadRequest)
	}

	var count int
	queryCount := fmt.Sprintf(`SELECT count(fr.id) 
									  FROM food_recipe fr
									      LEFT JOIN foods f
									          ON f.id = fr.food_id 
									      LEFT JOIN products p 
									          ON fr.product_id = p.id
									  %s`, countWhereQuery)
	if err = r.QueryRowContext(ctx, queryCount).Scan(&count); err != nil {
		return nil, 0, err
	}

	return list, count, nil
}

func (r Repository) AdminGetDetail(ctx context.Context, id int64) (*food_recipe.AdminGetDetail, error) {
	_, err := r.CheckClaims(ctx, auth.RoleAdmin)
	if err != nil {
		return nil, err
	}

	var detail food_recipe.AdminGetDetail

	query := fmt.Sprintf(`SELECT 
    									fr.id as id,
    									f.name as food,
    									fr.food_id,
    									p.name as product,
    									fr.product_id,
    									fr.amount as amount,
    									m.name as measure_unit
								 FROM food_recipe fr 
								     JOIN foods f 
								         ON fr.food_id = f.id 
								     JOIN products p 
								         ON fr.product_id = p.id 
								     JOIN measure_unit m
								         ON p.measure_unit_id = m.id
								 WHERE fr.deleted_at ISNULL AND f.deleted_at ISNULL AND p.deleted_at ISNULL AND fr.id='%d'`, id)

	err = r.QueryRowContext(ctx, query).Scan(&detail.ID, &detail.Food, &detail.FoodId, &detail.Product, &detail.ProductId, &detail.Amount, &detail.MeasureUnit)
	if err != nil {
		return nil, err
	}

	return &detail, nil
}

func (r Repository) AdminCreate(ctx context.Context, request food_recipe.AdminCreateRequest) (*food_recipe.AdminCreateResponse, error) {
	claims, err := r.CheckClaims(ctx, auth.RoleAdmin)
	if err != nil {
		return nil, err
	}

	err = r.ValidateStruct(&request, "ProductId", "FoodId")
	if err != nil {
		return nil, err
	}

	response := food_recipe.AdminCreateResponse{
		ProductId: request.ProductId,
		FoodId:    request.FoodId,
		Amount:    request.Amount,
		CreatedAt: time.Now(),
		CreatedBy: claims.UserId,
	}

	var id int64

	_, err = r.NewInsert().Model(&response).Returning("id").Exec(ctx, &id)
	if err != nil {
		return nil, web.NewRequestError(errors.Wrap(err, "creating food recipe"), http.StatusBadRequest)
	}

	response.ID = id

	return &response, nil
}

func (r Repository) AdminUpdateAll(ctx context.Context, request food_recipe.AdminUpdateRequest) error {
	claims, err := r.CheckClaims(ctx, auth.RoleAdmin)
	if err != nil {
		return err
	}

	if err = r.ValidateStruct(&request, "ID"); err != nil {
		return err
	}

	q := r.NewUpdate().Table("food_recipe").Where("deleted_at IS NULL AND id = ?", request.ID)

	q.Set("updated_at = ?", time.Now())
	q.Set("updated_by = ?", claims.UserId)
	q.Set("amount = ?", request.Amount)
	q.Set("product_id = ?", request.ProductId)
	q.Set("food_id = ?", request.FoodId)

	_, err = q.Exec(ctx)
	if err != nil {
		return web.NewRequestError(errors.Wrap(err, "updating food recipe"), http.StatusBadRequest)
	}

	return nil
}

func (r Repository) AdminUpdateColumns(ctx context.Context, request food_recipe.AdminUpdateRequest) error {
	claims, err := r.CheckClaims(ctx, auth.RoleAdmin)
	if err != nil {
		return err
	}

	if err = r.ValidateStruct(&request, "ID"); err != nil {
		return err
	}

	q := r.NewUpdate().Table("food_recipe").Where("deleted_at IS NULL AND id = ?", request.ID)

	if request.Amount != nil {
		q.Set("amount = ?", request.Amount)
	}
	if request.ProductId != nil {
		q.Set("product_id = ?", request.ProductId)
	}
	if request.FoodId != nil {
		q.Set("food_id = ?", request.FoodId)
	}

	q.Set("updated_at = ?", time.Now())
	q.Set("updated_by = ?", claims.UserId)

	_, err = q.Exec(ctx)
	if err != nil {
		return web.NewRequestError(errors.Wrap(err, "updating food recipe"), http.StatusBadRequest)
	}

	return nil
}

func (r Repository) AdminDelete(ctx context.Context, id int64) error {
	return r.DeleteRow(ctx, "food_recipe", id, auth.RoleAdmin)
}

// @branch

func (r Repository) BranchGetList(ctx context.Context, filter food_recipe.Filter) ([]food_recipe.BranchGetList, int, error) {
	claims, err := r.CheckClaims(ctx, auth.RoleBranch)
	if err != nil {
		return nil, 0, err
	}

	whereQuery := fmt.Sprintf(`WHERE fr.deleted_at ISNULL 
		AND f.deleted_at ISNULL AND p.deleted_at ISNULL
		AND f.restaurant_id in (select restaurant_id from branches where id = '%d') `, *claims.BranchID)

	if filter.Name != nil {
		whereQuery += fmt.Sprintf(" AND f.name ilike '%s'", "%"+*filter.Name+"%")
	}
	if filter.FoodId != nil {
		whereQuery += fmt.Sprintf(` AND fr.food_id='%d'`, *filter.FoodId)
	}

	countWhereQuery := whereQuery

	if filter.Limit != nil {
		whereQuery += fmt.Sprintf(" LIMIT '%d'", *filter.Limit)
	}
	if filter.Offset != nil {
		whereQuery += fmt.Sprintf(" OFFSET '%d'", *filter.Offset)
	}

	query := fmt.Sprintf(`SELECT 
    									fr.id as id,
    									f.name as food,
    									p.name as product,
    									fr.amount as amount,
    									m.name as measure_unit
								 FROM food_recipe fr 
								     JOIN foods f 
								         ON fr.food_id = f.id 
								     JOIN products p 
								         ON fr.product_id = p.id 
								     JOIN measure_unit m
								         ON p.measure_unit_id = m.id
									%s`, whereQuery)

	list := make([]food_recipe.BranchGetList, 0)

	rows, err := r.QueryContext(ctx, query)
	if err != nil {
		return nil, 0, web.NewRequestError(errors.Wrap(err, "select food recipe"), http.StatusInternalServerError)
	}

	err = r.ScanRows(ctx, rows, &list)
	if err != nil {
		return nil, 0, web.NewRequestError(errors.Wrap(err, "scanning food recipe"), http.StatusBadRequest)
	}

	var count int
	queryCount := fmt.Sprintf(`SELECT count(fr.id) 
									  FROM food_recipe fr
									      LEFT JOIN foods f
									          ON f.id = fr.food_id 
									      LEFT JOIN products p 
									          ON fr.product_id = p.id
									  %s`, countWhereQuery)
	if err = r.QueryRowContext(ctx, queryCount).Scan(&count); err != nil {
		return nil, 0, err
	}

	return list, count, nil
}

func (r Repository) BranchGetDetail(ctx context.Context, id int64) (*food_recipe.BranchGetDetail, error) {
	claims, err := r.CheckClaims(ctx, auth.RoleBranch)
	if err != nil {
		return nil, err
	}

	var restaurantID *int64

	err = r.QueryRowContext(ctx, fmt.Sprintf("SELECT restaurant_id FROM branches WHERE id = '%d'", *claims.BranchID)).Scan(&restaurantID)
	if err != nil {
		return nil, web.NewRequestError(errors.Wrap(err, "restaurant not found in "), http.StatusBadRequest)
	}

	if restaurantID == nil {
		return nil, web.NewRequestError(errors.New("not found restaurant"), http.StatusBadRequest)
	}

	var detail food_recipe.BranchGetDetail

	query := fmt.Sprintf(`SELECT 
    									fr.id as id,
    									f.name as food,
    									fr.food_id,
    									p.name as product,
    									fr.product_id,
    									fr.amount as amount,
    									m.name as measure_unit
								 FROM food_recipe fr 
								     JOIN foods f 
								         ON fr.food_id = f.id 
								     JOIN products p 
								         ON fr.product_id = p.id 
								     JOIN measure_unit m
								         ON p.measure_unit_id = m.id
								 WHERE fr.deleted_at ISNULL AND f.deleted_at ISNULL AND f.restaurant_id = '%d' AND p.deleted_at ISNULL AND fr.id='%d' `, *restaurantID, id)

	err = r.QueryRowContext(ctx, query).Scan(&detail.ID, &detail.Food, &detail.FoodId, &detail.Product, &detail.ProductId, &detail.Amount, &detail.MeasureUnit)

	if errors.Is(err, sql.ErrNoRows) {
		return nil, web.NewRequestError(postgres.ErrNotFound, http.StatusBadRequest)
	}
	if err != nil {
		return nil, err
	}

	return &detail, nil
}

func (r Repository) BranchCreate(ctx context.Context, request food_recipe.BranchCreateRequest) (*food_recipe.BranchCreateResponse, error) {
	claims, err := r.CheckClaims(ctx, auth.RoleBranch)
	if err != nil {
		return nil, err
	}

	err = r.ValidateStruct(&request, "ProductId", "FoodId")
	if err != nil {
		return nil, err
	}

	response := food_recipe.BranchCreateResponse{
		ProductId: request.ProductId,
		FoodId:    request.FoodId,
		Amount:    request.Amount,
		CreatedAt: time.Now(),
		CreatedBy: claims.UserId,
	}

	var id int64

	_, err = r.NewInsert().Model(&response).Returning("id").Exec(ctx, &id)
	if err != nil {
		return nil, web.NewRequestError(errors.Wrap(err, "creating food recipe"), http.StatusBadRequest)
	}

	response.ID = id

	return &response, nil
}

func (r Repository) BranchUpdateAll(ctx context.Context, request food_recipe.BranchUpdateRequest) error {
	claims, err := r.CheckClaims(ctx, auth.RoleBranch)
	if err != nil {
		return err
	}

	if err = r.ValidateStruct(&request, "ID"); err != nil {
		return err
	}

	q := r.NewUpdate().Table("food_recipe").Where("deleted_at IS NULL AND id = ?", request.ID)

	q.Set("updated_at = ?", time.Now())
	q.Set("updated_by = ?", claims.UserId)
	q.Set("amount = ?", request.Amount)
	q.Set("product_id = ?", request.ProductId)
	q.Set("food_id = ?", request.FoodId)

	_, err = q.Exec(ctx)
	if err != nil {
		return web.NewRequestError(errors.Wrap(err, "updating food recipe"), http.StatusBadRequest)
	}

	return nil
}

func (r Repository) BranchUpdateColumns(ctx context.Context, request food_recipe.BranchUpdateRequest) error {
	claims, err := r.CheckClaims(ctx, auth.RoleBranch)
	if err != nil {
		return err
	}

	if err = r.ValidateStruct(&request, "ID"); err != nil {
		return err
	}

	q := r.NewUpdate().Table("food_recipe").Where("deleted_at IS NULL AND id = ?", request.ID)

	if request.Amount != nil {
		q.Set("amount = ?", request.Amount)
	}
	if request.ProductId != nil {
		q.Set("product_id = ?", request.ProductId)
	}
	if request.FoodId != nil {
		q.Set("food_id = ?", request.FoodId)
	}

	q.Set("updated_at = ?", time.Now())
	q.Set("updated_by = ?", claims.UserId)

	_, err = q.Exec(ctx)
	if err != nil {
		return web.NewRequestError(errors.Wrap(err, "updating food recipe"), http.StatusBadRequest)
	}

	return nil
}

func (r Repository) BranchDelete(ctx context.Context, id int64) error {
	return r.DeleteRow(ctx, "food_recipe", id, auth.RoleBranch)
}
