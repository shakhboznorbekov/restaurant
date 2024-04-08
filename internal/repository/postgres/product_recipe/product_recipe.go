package product_recipe

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
	"restu-backend/internal/service/product_recipe"
	"time"
)

type Repository struct {
	*postgresql.Database
}

func NewRepository(db *postgresql.Database) *Repository {
	return &Repository{db}
}

// @admin

func (r Repository) AdminGetList(ctx context.Context, filter product_recipe.Filter) ([]product_recipe.AdminGetList, int, error) {
	claims, err := r.CheckClaims(ctx, auth.RoleAdmin)
	if err != nil {
		return nil, 0, err
	}

	whereQuery := fmt.Sprintf(`WHERE cr.deleted_at IS NULL AND p.deleted_at IS NULL AND p.restaurant_id='%d'`, *claims.RestaurantID)

	if filter.Name != nil {
		whereQuery += fmt.Sprintf(" AND f.name ilike '%s'", "%"+*filter.Name+"%")
	}
	if filter.ProductId != nil {
		whereQuery += fmt.Sprintf(` AND cr.product_id='%d'`, *filter.ProductId)
	}

	countWhereQuery := whereQuery

	if filter.Limit != nil {
		whereQuery += fmt.Sprintf(" LIMIT '%d'", *filter.Limit)
	}
	if filter.Offset != nil {
		whereQuery += fmt.Sprintf(" OFFSET '%d'", *filter.Offset)
	}

	query := fmt.Sprintf(`SELECT 
    									cr.id as id,
    									f.name as recipe,
    									p.name as product,
    									cr.amount as amount,
    									m.name as measure_unit
								 FROM product_recipe cr 
								     JOIN products f ON cr.recipe_id = f.id 
								     JOIN products p ON cr.product_id = p.id 
								     JOIN measure_unit m ON f.measure_unit_id = m.id
									%s`, whereQuery)

	list := make([]product_recipe.AdminGetList, 0)

	rows, err := r.QueryContext(ctx, query)
	if err != nil {
		return nil, 0, web.NewRequestError(errors.Wrap(err, "select product recipe"), http.StatusInternalServerError)
	}

	err = r.ScanRows(ctx, rows, &list)
	if err != nil {
		return nil, 0, web.NewRequestError(errors.Wrap(err, "scanning product recipe"), http.StatusBadRequest)
	}

	var count int
	queryCount := fmt.Sprintf(`SELECT count(cr.id) 
									  FROM product_recipe cr
									      LEFT JOIN products f
									          ON f.id = cr.recipe_id 
									      LEFT JOIN products p 
									          ON cr.product_id = p.id
									  %s`, countWhereQuery)
	if err = r.QueryRowContext(ctx, queryCount).Scan(&count); err != nil {
		return nil, 0, err
	}

	return list, count, nil
}

func (r Repository) AdminGetDetail(ctx context.Context, id int64) (*product_recipe.AdminGetDetail, error) {
	_, err := r.CheckClaims(ctx, auth.RoleAdmin)
	if err != nil {
		return nil, err
	}

	var detail product_recipe.AdminGetDetail

	query := fmt.Sprintf(`SELECT 
    									cr.id as id,
    									f.name as recipe,
    									cr.recipe_id,
    									p.name as product,
    									cr.product_id,
    									cr.amount as amount,
    									m.name as measure_unit
								 FROM product_recipe cr 
								     JOIN products f 
								         ON cr.recipe_id = f.id 
								     JOIN products p 
								         ON cr.product_id = p.id 
								     JOIN measure_unit m
								         ON p.measure_unit_id = m.id
								 WHERE cr.deleted_at ISNULL AND p.deleted_at IS NULL AND cr.id='%d'`, id)

	err = r.QueryRowContext(ctx, query).Scan(&detail.ID, &detail.Recipe, &detail.RecipeId, &detail.Product, &detail.ProductId, &detail.Amount, &detail.MeasureUnit)
	if err != nil {
		return nil, err
	}

	return &detail, nil
}

func (r Repository) AdminCreate(ctx context.Context, request product_recipe.AdminCreateRequest) (*product_recipe.AdminCreateResponse, error) {
	claims, err := r.CheckClaims(ctx, auth.RoleAdmin)
	if err != nil {
		return nil, err
	}

	err = r.ValidateStruct(&request, "ProductId", "RecipeId")
	if err != nil {
		return nil, err
	}

	response := product_recipe.AdminCreateResponse{
		ProductId: request.ProductId,
		RecipeId:  request.RecipeId,
		Amount:    request.Amount,
		CreatedAt: time.Now(),
		CreatedBy: claims.UserId,
	}

	var id int64

	_, err = r.NewInsert().Model(&response).Returning("id").Exec(ctx, &id)
	if err != nil {
		return nil, web.NewRequestError(errors.Wrap(err, "creating product recipe"), http.StatusBadRequest)
	}

	response.ID = id

	return &response, nil
}

func (r Repository) AdminUpdateAll(ctx context.Context, request product_recipe.AdminUpdateRequest) error {
	claims, err := r.CheckClaims(ctx, auth.RoleAdmin)
	if err != nil {
		return err
	}

	if err = r.ValidateStruct(&request, "ID"); err != nil {
		return err
	}

	q := r.NewUpdate().Table("product_recipe").Where("deleted_at IS NULL AND id = ?", request.ID)

	q.Set("updated_at = ?", time.Now())
	q.Set("updated_by = ?", claims.UserId)
	q.Set("amount = ?", request.Amount)
	q.Set("product_id = ?", request.ProductId)
	q.Set("recipe_id = ?", request.RecipeId)

	_, err = q.Exec(ctx)
	if err != nil {
		return web.NewRequestError(errors.Wrap(err, "updating product recipe"), http.StatusBadRequest)
	}

	return nil
}

func (r Repository) AdminUpdateColumns(ctx context.Context, request product_recipe.AdminUpdateRequest) error {
	claims, err := r.CheckClaims(ctx, auth.RoleAdmin)
	if err != nil {
		return err
	}

	if err = r.ValidateStruct(&request, "ID"); err != nil {
		return err
	}

	q := r.NewUpdate().Table("product_recipe").Where("deleted_at IS NULL AND id = ?", request.ID)

	if request.Amount != nil {
		q.Set("amount = ?", request.Amount)
	}
	if request.ProductId != nil {
		q.Set("product_id = ?", request.ProductId)
	}
	if request.RecipeId != nil {
		q.Set("recipe_id = ?", request.RecipeId)
	}

	q.Set("updated_at = ?", time.Now())
	q.Set("updated_by = ?", claims.UserId)

	_, err = q.Exec(ctx)
	if err != nil {
		return web.NewRequestError(errors.Wrap(err, "updating product recipe"), http.StatusBadRequest)
	}

	return nil
}

func (r Repository) AdminDelete(ctx context.Context, id int64) error {
	return r.DeleteRow(ctx, "product_recipe", id, auth.RoleAdmin)
}

// @branch

func (r Repository) BranchGetList(ctx context.Context, filter product_recipe.Filter) ([]product_recipe.BranchGetList, int, error) {
	claims, err := r.CheckClaims(ctx, auth.RoleBranch)
	if err != nil {
		return nil, 0, err
	}

	whereQuery := fmt.Sprintf(`WHERE cr.deleted_at IS NULL AND p.deleted_at IS NULL AND p.restaurant_id in (select restaurant_id from branches where id = '%d')`, *claims.BranchID)

	if filter.Name != nil {
		whereQuery += fmt.Sprintf(" AND f.name ilike '%s'", "%"+*filter.Name+"%")
	}
	if filter.ProductId != nil {
		whereQuery += fmt.Sprintf(` AND cr.product_id='%d'`, *filter.ProductId)
	}

	countWhereQuery := whereQuery

	if filter.Limit != nil {
		whereQuery += fmt.Sprintf(" LIMIT '%d'", *filter.Limit)
	}
	if filter.Offset != nil {
		whereQuery += fmt.Sprintf(" OFFSET '%d'", *filter.Offset)
	}

	query := fmt.Sprintf(`SELECT 
    									cr.id as id,
    									f.name as recipe,
    									p.name as product,
    									cr.amount as amount,
    									m.name as measure_unit
								 FROM product_recipe cr 
								     JOIN products f ON cr.recipe_id = f.id 
								     JOIN products p ON cr.product_id = p.id 
								     JOIN measure_unit m ON f.measure_unit_id = m.id
									%s`, whereQuery)

	list := make([]product_recipe.BranchGetList, 0)

	rows, err := r.QueryContext(ctx, query)
	if err != nil {
		return nil, 0, web.NewRequestError(errors.Wrap(err, "select product recipe"), http.StatusInternalServerError)
	}

	err = r.ScanRows(ctx, rows, &list)
	if err != nil {
		return nil, 0, web.NewRequestError(errors.Wrap(err, "scanning product recipe"), http.StatusBadRequest)
	}

	var count int
	queryCount := fmt.Sprintf(`SELECT count(cr.id) 
									  FROM product_recipe cr
									      LEFT JOIN products f
									          ON f.id = cr.recipe_id 
									      LEFT JOIN products p 
									          ON cr.product_id = p.id
									  %s`, countWhereQuery)
	if err = r.QueryRowContext(ctx, queryCount).Scan(&count); err != nil {
		return nil, 0, err
	}

	return list, count, nil
}

func (r Repository) BranchGetDetail(ctx context.Context, id int64) (*product_recipe.BranchGetDetail, error) {
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

	var detail product_recipe.BranchGetDetail

	query := fmt.Sprintf(`SELECT 
    									cr.id as id,
    									f.name as recipe,
    									cr.recipe_id,
    									p.name as product,
    									cr.product_id,
    									cr.amount as amount,
    									m.name as measure_unit
								 FROM product_recipe cr 
								     JOIN products f 
								         ON cr.recipe_id = f.id 
								     JOIN products p 
								         ON cr.product_id = p.id 
								     JOIN measure_unit m
								         ON p.measure_unit_id = m.id
								 WHERE cr.deleted_at ISNULL AND p.restaurant_id = '%d' AND p.deleted_at IS NULL AND cr.id='%d'`, *restaurantID, id)

	err = r.QueryRowContext(ctx, query).Scan(&detail.ID, &detail.Recipe, &detail.RecipeId, &detail.Product, &detail.ProductId, &detail.Amount, &detail.MeasureUnit)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, web.NewRequestError(postgres.ErrNotFound, http.StatusBadRequest)
	}
	if err != nil {
		return nil, web.NewRequestError(errors.Wrap(err, "selecting products"), http.StatusBadRequest)
	}

	return &detail, nil
}

func (r Repository) BranchCreate(ctx context.Context, request product_recipe.BranchCreateRequest) (*product_recipe.BranchCreateResponse, error) {
	claims, err := r.CheckClaims(ctx, auth.RoleBranch)
	if err != nil {
		return nil, err
	}

	err = r.ValidateStruct(&request, "ProductId", "RecipeId")
	if err != nil {
		return nil, err
	}

	response := product_recipe.BranchCreateResponse{
		ProductId: request.ProductId,
		RecipeId:  request.RecipeId,
		Amount:    request.Amount,
		CreatedAt: time.Now(),
		CreatedBy: claims.UserId,
	}

	var id int64

	_, err = r.NewInsert().Model(&response).Returning("id").Exec(ctx, &id)
	if err != nil {
		return nil, web.NewRequestError(errors.Wrap(err, "creating product recipe"), http.StatusBadRequest)
	}

	response.ID = id

	return &response, nil
}

func (r Repository) BranchUpdateAll(ctx context.Context, request product_recipe.BranchUpdateRequest) error {
	claims, err := r.CheckClaims(ctx, auth.RoleBranch)
	if err != nil {
		return err
	}

	if err = r.ValidateStruct(&request, "ID"); err != nil {
		return err
	}

	q := r.NewUpdate().Table("product_recipe").Where("deleted_at IS NULL AND id = ?", request.ID)

	q.Set("updated_at = ?", time.Now())
	q.Set("updated_by = ?", claims.UserId)
	q.Set("amount = ?", request.Amount)
	q.Set("product_id = ?", request.ProductId)
	q.Set("recipe_id = ?", request.RecipeId)

	_, err = q.Exec(ctx)
	if err != nil {
		return web.NewRequestError(errors.Wrap(err, "updating product recipe"), http.StatusBadRequest)
	}

	return nil
}

func (r Repository) BranchUpdateColumns(ctx context.Context, request product_recipe.BranchUpdateRequest) error {
	claims, err := r.CheckClaims(ctx, auth.RoleBranch)
	if err != nil {
		return err
	}

	if err = r.ValidateStruct(&request, "ID"); err != nil {
		return err
	}

	q := r.NewUpdate().Table("product_recipe").Where("deleted_at IS NULL AND id = ?", request.ID)

	if request.Amount != nil {
		q.Set("amount = ?", request.Amount)
	}
	if request.ProductId != nil {
		q.Set("product_id = ?", request.ProductId)
	}
	if request.RecipeId != nil {
		q.Set("recipe_id = ?", request.RecipeId)
	}

	q.Set("updated_at = ?", time.Now())
	q.Set("updated_by = ?", claims.UserId)

	_, err = q.Exec(ctx)
	if err != nil {
		return web.NewRequestError(errors.Wrap(err, "updating product recipe"), http.StatusBadRequest)
	}

	return nil
}

func (r Repository) BranchDelete(ctx context.Context, id int64) error {
	return r.DeleteRow(ctx, "product_recipe", id, auth.RoleBranch)
}
