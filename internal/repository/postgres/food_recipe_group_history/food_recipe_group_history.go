package food_recipe_group_history

import (
	"context"
	"fmt"
	"github.com/pkg/errors"
	"net/http"
	"restu-backend/foundation/web"
	"restu-backend/internal/auth"
	"restu-backend/internal/pkg/repository/postgresql"
	"restu-backend/internal/service/food_recipe_group_history"
	"time"
)

type Repository struct {
	*postgresql.Database
}

// @admin

func (r Repository) AdminGetList(ctx context.Context, filter food_recipe_group_history.Filter) ([]food_recipe_group_history.AdminGetListByFoodID, int, error) {
	_, err := r.CheckClaims(ctx, auth.RoleAdmin)
	if err != nil {
		return nil, 0, err
	}

	whereQuery := fmt.Sprintf(`WHERE fh.deleted_at ISNULL AND fh.food_id = %d`, *filter.FoodID)

	var limitQuery, offsetQuery string
	if filter.Limit != nil {
		limitQuery = fmt.Sprintf(" LIMIT '%d'", *filter.Limit)
	}
	if filter.Offset != nil {
		offsetQuery = fmt.Sprintf(" OFFSET '%d'", *filter.Offset)
	}

	query := fmt.Sprintf(`SELECT 
    									fh.id,
    									TO_CHAR(fh.from, 'DD.MM.YYYY') as "from",
    									TO_CHAR(fh.to, 'DD.MM.YYYY') as "to",
    									f.name as food,
    									fg.name as "group"
								 FROM food_recipe_group_histories as fh
								 LEFT OUTER JOIN foods as f ON f.id = fh.food_id
								 LEFT OUTER JOIN food_recipe_groups as fg ON fg.id = fh.group_id
								 %s %s %s`, whereQuery, limitQuery, offsetQuery)

	list := make([]food_recipe_group_history.AdminGetListByFoodID, 0)

	rows, err := r.QueryContext(ctx, query)
	if err != nil {
		return nil, 0, web.NewRequestError(errors.Wrap(err, "select food recipe"), http.StatusInternalServerError)
	}

	err = r.ScanRows(ctx, rows, &list)
	if err != nil {
		return nil, 0, web.NewRequestError(errors.Wrap(err, "scanning food recipe"), http.StatusBadRequest)
	}

	var count int
	queryCount := fmt.Sprintf(`
							SELECT 
							    count(fh.id) 	  
							FROM food_recipe_group_histories as fh
							LEFT OUTER JOIN foods as f ON f.id = fh.food_id
							LEFT OUTER JOIN food_recipe_groups as fg ON fg.id = fh.group_id
							%s`, whereQuery)
	if err = r.QueryRowContext(ctx, queryCount).Scan(&count); err != nil {
		return nil, 0, err
	}

	return list, count, nil
}

func (r Repository) AdminGetDetail(ctx context.Context, id int64) (*food_recipe_group_history.AdminGetDetail, error) {
	_, err := r.CheckClaims(ctx, auth.RoleAdmin)
	if err != nil {
		return nil, err
	}

	var detail food_recipe_group_history.AdminGetDetail

	query := fmt.Sprintf(`SELECT 
    									fh.id,
    									TO_CHAR(fh.from, 'DD.MM.YYYY') as "from",
    									TO_CHAR(fh.to, 'DD.MM.YYYY') as "to",
    									f.name as food,
    									fg.name as "group",
    									fh.food_id,
    									fh.group_id
								 FROM food_recipe_group_histories as fh
								 LEFT OUTER JOIN foods as f ON f.id = fh.food_id
								 LEFT OUTER JOIN food_recipe_groups as fg ON fg.id = fh.group_id
								 WHERE fh.id='%d'`, id)

	err = r.QueryRowContext(ctx, query).Scan(
		&detail.ID,
		&detail.From,
		&detail.To,
		&detail.Food,
		&detail.Group,
		&detail.FoodID,
		&detail.GroupID)
	if err != nil {
		return nil, err
	}

	return &detail, nil
}

func (r Repository) AdminCreate(ctx context.Context, request food_recipe_group_history.AdminCreateRequest) (*food_recipe_group_history.AdminCreateResponse, error) {
	claims, err := r.CheckClaims(ctx, auth.RoleAdmin)
	if err != nil {
		return nil, nil
	}

	err = r.ValidateStruct(&request, "Date", "FoodID", "GroupID")
	if err != nil {
		return nil, err
	}

	var detail food_recipe_group_history.AdminCreateResponse
	var dateT time.Time
	if request.Date != nil {
		dateT, err = time.Parse("02.01.2006", *request.Date)
		if err != nil {
			return nil, err
		}
	}

	query := fmt.Sprintf(`WITH new_record AS (
										INSERT INTO food_recipe_group_histories ("from", "to", food_id, group_id, created_at, created_by)
											VALUES ('%s', (
												SELECT "from"
												FROM food_recipe_group_histories
												WHERE "from" > '%s' AND food_id = %d
												ORDER BY "from"
												LIMIT 1
											), '%d', '%d', '%v', '%d')
											ON CONFLICT("from", food_id)
												DO UPDATE SET group_id='%d'
											RETURNING id, "from",food_id)
								UPDATE food_recipe_group_histories
								SET "to" = (SELECT "from" FROM new_record LIMIT 1)
								WHERE id = (
									SELECT id
									FROM food_recipe_group_histories
									WHERE
										"from" < (SELECT "from" FROM new_record LIMIT 1) AND
										food_id = (SELECT food_id FROM new_record LIMIT 1)
									ORDER BY "from" DESC
									LIMIT 1
								);`,
		dateT,
		dateT,
		*request.FoodID,
		*request.FoodID,
		*request.GroupID,
		time.Now(),
		claims.UserId,
		*request.GroupID,
	)

	if _, err = r.ExecContext(ctx, query, &detail); err != nil {
		return nil, err
	}

	return &detail, nil
}

func (r Repository) AdminDelete(ctx context.Context, id int64) error {
	return r.DeleteRow(ctx, "food_recipe_group_histories", id, auth.RoleAdmin)

}

// @branch

func (r Repository) BranchGetList(ctx context.Context, filter food_recipe_group_history.Filter) ([]food_recipe_group_history.BranchGetListByFoodID, int, error) {
	_, err := r.CheckClaims(ctx, auth.RoleBranch)
	if err != nil {
		return nil, 0, err
	}

	whereQuery := fmt.Sprintf(`WHERE fh.deleted_at ISNULL AND fh.food_id = %d`, *filter.FoodID)

	var limitQuery, offsetQuery string
	if filter.Limit != nil {
		limitQuery = fmt.Sprintf(" LIMIT '%d'", *filter.Limit)
	}
	if filter.Offset != nil {
		offsetQuery = fmt.Sprintf(" OFFSET '%d'", *filter.Offset)
	}

	query := fmt.Sprintf(`SELECT 
    									fh.id,
    									TO_CHAR(fh.from, 'DD.MM.YYYY') as "from",
    									TO_CHAR(fh.to, 'DD.MM.YYYY') as "to",
    									f.name as food,
    									fg.name as "group"
								 FROM food_recipe_group_histories as fh
								 LEFT OUTER JOIN foods as f ON f.id = fh.food_id
								 LEFT OUTER JOIN food_recipe_groups as fg ON fg.id = fh.group_id
								 %s %s %s`, whereQuery, limitQuery, offsetQuery)

	list := make([]food_recipe_group_history.BranchGetListByFoodID, 0)

	rows, err := r.QueryContext(ctx, query)
	if err != nil {
		return nil, 0, web.NewRequestError(errors.Wrap(err, "select food recipe"), http.StatusInternalServerError)
	}

	err = r.ScanRows(ctx, rows, &list)
	if err != nil {
		return nil, 0, web.NewRequestError(errors.Wrap(err, "scanning food recipe"), http.StatusBadRequest)
	}

	var count int
	queryCount := fmt.Sprintf(`
							SELECT 
							    count(fh.id) 	  
							FROM food_recipe_group_histories as fh
							LEFT OUTER JOIN foods as f ON f.id = fh.food_id
							LEFT OUTER JOIN food_recipe_groups as fg ON fg.id = fh.group_id
							%s`, whereQuery)
	if err = r.QueryRowContext(ctx, queryCount).Scan(&count); err != nil {
		return nil, 0, err
	}

	return list, count, nil
}

func (r Repository) BranchGetDetail(ctx context.Context, id int64) (*food_recipe_group_history.BranchGetDetail, error) {
	_, err := r.CheckClaims(ctx, auth.RoleBranch)
	if err != nil {
		return nil, err
	}

	var detail food_recipe_group_history.BranchGetDetail

	query := fmt.Sprintf(`SELECT 
    									fh.id,
    									TO_CHAR(fh.from, 'DD.MM.YYYY') as "from",
    									TO_CHAR(fh.to, 'DD.MM.YYYY') as "to",
    									f.name as food,
    									fg.name as "group",
    									fh.food_id,
    									fh.group_id
								 FROM food_recipe_group_histories as fh
								 LEFT OUTER JOIN foods as f ON f.id = fh.food_id
								 LEFT OUTER JOIN food_recipe_groups as fg ON fg.id = fh.group_id
								 WHERE fh.id='%d'`, id)

	err = r.QueryRowContext(ctx, query).Scan(
		&detail.ID,
		&detail.From,
		&detail.To,
		&detail.Food,
		&detail.Group,
		&detail.FoodID,
		&detail.GroupID)
	if err != nil {
		return nil, err
	}

	return &detail, nil
}

func (r Repository) BranchCreate(ctx context.Context, request food_recipe_group_history.BranchCreateRequest) (*food_recipe_group_history.BranchCreateResponse, error) {
	claims, err := r.CheckClaims(ctx, auth.RoleBranch)
	if err != nil {
		return nil, nil
	}

	err = r.ValidateStruct(&request, "Date", "FoodID", "GroupID")
	if err != nil {
		return nil, err
	}

	var detail food_recipe_group_history.BranchCreateResponse
	var dateT time.Time
	if request.Date != nil {
		dateT, err = time.Parse("02.01.2006", *request.Date)
		if err != nil {
			return nil, err
		}
	}

	query := fmt.Sprintf(`WITH new_record AS (
										INSERT INTO food_recipe_group_histories ("from", "to", food_id, group_id, created_at, created_by)
											VALUES ('%s', (
												SELECT "from"
												FROM food_recipe_group_histories
												WHERE "from" > '%s' AND food_id = %d
												ORDER BY "from"
												LIMIT 1
											), '%d', '%d', '%v', '%d')
											ON CONFLICT("from", food_id)
												DO UPDATE SET group_id='%d'
											RETURNING id, "from",food_id)
								UPDATE food_recipe_group_histories
								SET "to" = (SELECT "from" FROM new_record LIMIT 1)
								WHERE id = (
									SELECT id
									FROM food_recipe_group_histories
									WHERE
										"from" < (SELECT "from" FROM new_record LIMIT 1) AND
										food_id = (SELECT food_id FROM new_record LIMIT 1)
									ORDER BY "from" DESC
									LIMIT 1
								);`,
		dateT,
		dateT,
		*request.FoodID,
		*request.FoodID,
		*request.GroupID,
		time.Now(),
		claims.UserId,
		*request.GroupID,
	)

	if _, err = r.ExecContext(ctx, query, &detail); err != nil {
		return nil, err
	}

	return &detail, nil
}

func (r Repository) BranchDelete(ctx context.Context, id int64) error {
	return r.DeleteRow(ctx, "food_recipe_group_histories", id, auth.RoleBranch)

}

func NewRepository(db *postgresql.Database) *Repository {
	return &Repository{db}
}
