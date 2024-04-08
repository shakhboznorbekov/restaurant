package food_recipe_group

import (
	"fmt"
	"github.com/pkg/errors"
	"golang.org/x/net/context"
	"net/http"
	"restu-backend/foundation/web"
	"restu-backend/internal/auth"
	"restu-backend/internal/pkg/repository/postgresql"
	"restu-backend/internal/service/food_recipe_group"
	"time"
)

type Repository struct {
	*postgresql.Database
}

func NewRepository(db *postgresql.Database) *Repository {
	return &Repository{db}
}

// @admin

func (r Repository) AdminGetList(ctx context.Context, filter food_recipe_group.Filter, foodID int64) ([]food_recipe_group.AdminGetListByFoodID, int, error) {
	_, err := r.CheckClaims(ctx, auth.RoleAdmin)
	if err != nil {
		return nil, 0, err
	}

	whereQuery := fmt.Sprintf(`WHERE deleted_at ISNULL AND food_id = %d`, foodID)

	if filter.Name != nil {
		whereQuery += fmt.Sprintf(" AND f.name ilike '%s'", "%"+*filter.Name+"%")
	}

	var limitQuery, offsetQuery string
	if filter.Limit != nil {
		limitQuery = fmt.Sprintf(" LIMIT '%d'", *filter.Limit)
	}
	if filter.Offset != nil {
		offsetQuery = fmt.Sprintf(" OFFSET '%d'", *filter.Offset)
	}

	query := fmt.Sprintf(`SELECT 
    									id,
    									name
								 FROM food_recipe_groups
								 %s %s %s`, whereQuery, limitQuery, offsetQuery)

	list := make([]food_recipe_group.AdminGetListByFoodID, 0)

	rows, err := r.QueryContext(ctx, query)
	if err != nil {
		return nil, 0, web.NewRequestError(errors.Wrap(err, "select food recipe"), http.StatusInternalServerError)
	}

	err = r.ScanRows(ctx, rows, &list)
	if err != nil {
		return nil, 0, web.NewRequestError(errors.Wrap(err, "scanning food recipe"), http.StatusBadRequest)
	}

	var count int
	queryCount := fmt.Sprintf(`SELECT count(id) 
									  FROM food_recipe_groups
									  %s`, whereQuery)
	if err = r.QueryRowContext(ctx, queryCount).Scan(&count); err != nil {
		return nil, 0, err
	}

	return list, count, nil
}

func (r Repository) AdminGetDetail(ctx context.Context, id int64) (*food_recipe_group.AdminGetDetail, error) {
	_, err := r.CheckClaims(ctx, auth.RoleAdmin)
	if err != nil {
		return nil, err
	}

	var detail food_recipe_group.AdminGetDetail

	query := fmt.Sprintf(`SELECT 
    									id,
    									name,
    									food_id,
    								    recipe_ids
								 FROM food_recipe_groups
								 WHERE id='%d'`, id)

	err = r.QueryRowContext(ctx, query).Scan(&detail.ID, &detail.Name, &detail.FoodId, &detail.RecipeIds)
	if err != nil {
		return nil, err
	}

	return &detail, nil
}

func (r Repository) AdminCreate(ctx context.Context, request food_recipe_group.AdminCreateRequest) (*food_recipe_group.AdminCreateResponse, error) {
	claims, err := r.CheckClaims(ctx, auth.RoleAdmin)
	if err != nil {
		return nil, err
	}

	err = r.ValidateStruct(&request, "Name", "RecipeIds", "FoodId")
	if err != nil {
		return nil, err
	}

	response := food_recipe_group.AdminCreateResponse{
		Name:      request.Name,
		FoodID:    request.FoodId,
		RecipeIds: request.RecipeIds,
		CreatedAt: time.Now(),
		CreatedBy: claims.UserId,
	}

	_, err = r.NewInsert().Model(&response).Exec(ctx)
	if err != nil {
		return nil, web.NewRequestError(errors.Wrap(err, "creating food recipe"), http.StatusBadRequest)
	}

	return &response, nil
}

func (r Repository) AdminUpdateAll(ctx context.Context, request food_recipe_group.AdminUpdateRequest) error {
	claims, err := r.CheckClaims(ctx, auth.RoleAdmin)
	if err != nil {
		return err
	}

	if err = r.ValidateStruct(&request, "ID", "Name", "FoodId", "RecipeIds"); err != nil {
		return err
	}

	pqRecipeIdsArray := "{"
	for k, v := range request.RecipeIds {
		if k == 0 {
			pqRecipeIdsArray += fmt.Sprintf(`%d`, v)
		} else {
			pqRecipeIdsArray += fmt.Sprintf(`, %d`, v)
		}
	}
	pqRecipeIdsArray += "}"

	q := r.NewUpdate().Table("food_recipe_groups").Where("deleted_at IS NULL AND id = ?", request.ID)

	q.Set("recipe_ids = array_cat(recipe_ids, ?)", pqRecipeIdsArray)
	q.Set("name = ?", request.Name)
	q.Set("food_id = ?", request.FoodId)
	q.Set("updated_at = ?", time.Now())
	q.Set("updated_by = ?", claims.UserId)

	_, err = q.Exec(ctx)
	if err != nil {
		return web.NewRequestError(errors.Wrap(err, "updating food recipe"), http.StatusBadRequest)
	}

	return nil
}

func (r Repository) AdminUpdateColumns(ctx context.Context, request food_recipe_group.AdminUpdateRequest) error {
	claims, err := r.CheckClaims(ctx, auth.RoleAdmin)
	if err != nil {
		return err
	}

	if err = r.ValidateStruct(&request, "ID"); err != nil {
		return err
	}

	q := r.NewUpdate().Table("food_recipe_groups").Where("deleted_at IS NULL AND id = ?", request.ID)

	if request.Name != nil {
		q.Set("Name = ?", request.Name)
	}
	if request.RecipeIds != nil {
		pqRecipeIdsArray := "{"
		for k, v := range request.RecipeIds {
			if k == 0 {
				pqRecipeIdsArray += fmt.Sprintf(`%d`, v)
			} else {
				pqRecipeIdsArray += fmt.Sprintf(`, %d`, v)
			}
		}
		pqRecipeIdsArray += "}"

		q.Set("recipe_ids = array_cat(recipe_ids, ?)", pqRecipeIdsArray)
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

func (r Repository) AdminDeleteRecipe(ctx context.Context, request food_recipe_group.AdminDeleteRecipeRequest) error {
	claims, err := r.CheckClaims(ctx, auth.RoleAdmin)
	if err != nil {
		return err
	}

	if err = r.ValidateStruct(&request, "ID", "RecipeId"); err != nil {
		return err
	}

	q := r.NewUpdate().Table("food_recipe_groups").Where("deleted_at IS NULL AND id = ?", request.ID)

	q.Set("recipe_ids = array_remove(recipe_ids, ?)", request.RecipeId)
	q.Set("updated_at = ?", time.Now())
	q.Set("updated_by = ?", claims.UserId)

	_, err = q.Exec(ctx)
	if err != nil {
		return web.NewRequestError(errors.Wrap(err, "updating food recipe"), http.StatusBadRequest)
	}

	return nil
}

func (r Repository) AdminDelete(ctx context.Context, id int64) error {
	return r.DeleteRow(ctx, "food_recipe_groups", id, auth.RoleAdmin)
}

// @branch

func (r Repository) BranchGetList(ctx context.Context, filter food_recipe_group.Filter, foodID int64) ([]food_recipe_group.BranchGetListByFoodID, int, error) {
	_, err := r.CheckClaims(ctx, auth.RoleBranch)
	if err != nil {
		return nil, 0, err
	}

	whereQuery := fmt.Sprintf(`WHERE deleted_at ISNULL AND food_id = %d`, foodID)

	if filter.Name != nil {
		whereQuery += fmt.Sprintf(" AND f.name ilike '%s'", "%"+*filter.Name+"%")
	}

	var limitQuery, offsetQuery string
	if filter.Limit != nil {
		limitQuery = fmt.Sprintf(" LIMIT '%d'", *filter.Limit)
	}
	if filter.Offset != nil {
		offsetQuery = fmt.Sprintf(" OFFSET '%d'", *filter.Offset)
	}

	query := fmt.Sprintf(`SELECT 
    									id,
    									name
								 FROM food_recipe_groups
								 %s %s %s`, whereQuery, limitQuery, offsetQuery)

	list := make([]food_recipe_group.BranchGetListByFoodID, 0)

	rows, err := r.QueryContext(ctx, query)
	if err != nil {
		return nil, 0, web.NewRequestError(errors.Wrap(err, "select food recipe"), http.StatusInternalServerError)
	}

	err = r.ScanRows(ctx, rows, &list)
	if err != nil {
		return nil, 0, web.NewRequestError(errors.Wrap(err, "scanning food recipe"), http.StatusBadRequest)
	}

	var count int
	queryCount := fmt.Sprintf(`SELECT count(id) 
									  FROM food_recipe_groups
									  %s`, whereQuery)
	if err = r.QueryRowContext(ctx, queryCount).Scan(&count); err != nil {
		return nil, 0, err
	}

	return list, count, nil
}

func (r Repository) BranchGetDetail(ctx context.Context, id int64) (*food_recipe_group.BranchGetDetail, error) {
	_, err := r.CheckClaims(ctx, auth.RoleBranch)
	if err != nil {
		return nil, err
	}

	var detail food_recipe_group.BranchGetDetail

	query := fmt.Sprintf(`SELECT 
    									id,
    									name,
    									food_id,
    								    recipe_ids
								 FROM food_recipe_groups
								 WHERE id='%d'`, id)

	err = r.QueryRowContext(ctx, query).Scan(&detail.ID, &detail.Name, &detail.FoodId, &detail.RecipeIds)
	if err != nil {
		return nil, err
	}

	return &detail, nil
}

func (r Repository) BranchCreate(ctx context.Context, request food_recipe_group.BranchCreateRequest) (*food_recipe_group.BranchCreateResponse, error) {
	claims, err := r.CheckClaims(ctx, auth.RoleBranch)
	if err != nil {
		return nil, err
	}

	err = r.ValidateStruct(&request, "Name", "RecipeIds", "FoodId")
	if err != nil {
		return nil, err
	}

	response := food_recipe_group.BranchCreateResponse{
		Name:      request.Name,
		FoodID:    request.FoodId,
		RecipeIds: request.RecipeIds,
		CreatedAt: time.Now(),
		CreatedBy: claims.UserId,
	}

	_, err = r.NewInsert().Model(&response).Exec(ctx)
	if err != nil {
		return nil, web.NewRequestError(errors.Wrap(err, "creating food recipe"), http.StatusBadRequest)
	}

	return &response, nil
}

func (r Repository) BranchUpdateAll(ctx context.Context, request food_recipe_group.BranchUpdateRequest) error {
	claims, err := r.CheckClaims(ctx, auth.RoleBranch)
	if err != nil {
		return err
	}

	if err = r.ValidateStruct(&request, "ID", "Name", "FoodId", "RecipeIds"); err != nil {
		return err
	}

	pqRecipeIdsArray := "{"
	for k, v := range request.RecipeIds {
		if k == 0 {
			pqRecipeIdsArray += fmt.Sprintf(`%d`, v)
		} else {
			pqRecipeIdsArray += fmt.Sprintf(`, %d`, v)
		}
	}
	pqRecipeIdsArray += "}"

	q := r.NewUpdate().Table("food_recipe_groups").Where("deleted_at IS NULL AND id = ?", request.ID)

	q.Set("recipe_ids = array_cat(recipe_ids, ?)", pqRecipeIdsArray)
	q.Set("name = ?", request.Name)
	q.Set("food_id = ?", request.FoodId)
	q.Set("updated_at = ?", time.Now())
	q.Set("updated_by = ?", claims.UserId)

	_, err = q.Exec(ctx)
	if err != nil {
		return web.NewRequestError(errors.Wrap(err, "updating food recipe"), http.StatusBadRequest)
	}

	return nil
}

func (r Repository) BranchUpdateColumns(ctx context.Context, request food_recipe_group.BranchUpdateRequest) error {
	claims, err := r.CheckClaims(ctx, auth.RoleBranch)
	if err != nil {
		return err
	}

	if err = r.ValidateStruct(&request, "ID"); err != nil {
		return err
	}

	q := r.NewUpdate().Table("food_recipe_groups").Where("deleted_at IS NULL AND id = ?", request.ID)

	if request.Name != nil {
		q.Set("Name = ?", request.Name)
	}
	if request.RecipeIds != nil {
		pqRecipeIdsArray := "{"
		for k, v := range request.RecipeIds {
			if k == 0 {
				pqRecipeIdsArray += fmt.Sprintf(`%d`, v)
			} else {
				pqRecipeIdsArray += fmt.Sprintf(`, %d`, v)
			}
		}
		pqRecipeIdsArray += "}"

		q.Set("recipe_ids = array_cat(recipe_ids, ?)", pqRecipeIdsArray)
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

func (r Repository) BranchDeleteRecipe(ctx context.Context, request food_recipe_group.BranchDeleteRecipeRequest) error {
	claims, err := r.CheckClaims(ctx, auth.RoleBranch)
	if err != nil {
		return err
	}

	if err = r.ValidateStruct(&request, "ID", "RecipeId"); err != nil {
		return err
	}

	q := r.NewUpdate().Table("food_recipe_groups").Where("deleted_at IS NULL AND id = ?", request.ID)

	q.Set("recipe_ids = array_remove(recipe_ids, ?)", request.RecipeId)
	q.Set("updated_at = ?", time.Now())
	q.Set("updated_by = ?", claims.UserId)

	_, err = q.Exec(ctx)
	if err != nil {
		return web.NewRequestError(errors.Wrap(err, "updating food recipe"), http.StatusBadRequest)
	}

	return nil
}

func (r Repository) BranchDelete(ctx context.Context, id int64) error {
	return r.DeleteRow(ctx, "food_recipe_groups", id, auth.RoleBranch)
}
