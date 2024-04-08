package food

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/lib/pq"
	"github.com/pkg/errors"
	"github.com/uptrace/bun/dialect/pgdialect"
	"net/http"
	"restu-backend/foundation/web"
	"restu-backend/internal/auth"
	"restu-backend/internal/entity"
	"restu-backend/internal/pkg/repository/postgresql"
	"restu-backend/internal/pkg/utils"
	"restu-backend/internal/repository/postgres"
	"restu-backend/internal/service/food"
	"restu-backend/internal/service/hashing"
	"time"
)

type Repository struct {
	*postgresql.Database
}

// admin

func (r Repository) AdminGetList(ctx context.Context, filter food.Filter) ([]food.AdminGetList, int, error) {
	claims, err := r.CheckClaims(ctx, auth.RoleAdmin)
	if err != nil {
		return nil, 0, err
	}

	table := "foods"
	whereQuery := fmt.Sprintf(`WHERE %s.deleted_at IS NULL AND %s.restaurant_id = %d`, table, table, *claims.RestaurantID)

	if filter.Search != nil {
		whereQuery += fmt.Sprintf(` AND (foods.name ilike '%s' OR food_category.name ilike '%s')`,
			"%"+*filter.Search+"%", "%"+*filter.Search+"%")
	}

	if filter.BranchID != nil {
		whereQuery += fmt.Sprintf(` AND %s.branch_id = %d`, table, *filter.BranchID)
	}
	var limitQuery, offsetQuery string
	if filter.Limit != nil {
		limitQuery = fmt.Sprintf(" LIMIT '%d'", *filter.Limit)
	}
	if filter.Offset != nil {
		offsetQuery = fmt.Sprintf(" OFFSET '%d'", *filter.Offset)
	}

	query := fmt.Sprintf(`
			SELECT
				foods.id,
				foods.name,
				foods.photos,
				foods.price,
				foods.category_id,
				food_category.name
			FROM foods
			LEFT JOIN food_category ON foods.category_id = food_category.id
			%s %s %s
`, whereQuery, limitQuery, offsetQuery)

	//whereQuery += fmt.Sprintf(" %s %s", limitQuery, offsetQuery)

	//query, err := utils.SelectQuery(filter.Fields, filter.Joins, &table, &whereQuery)
	//if err != nil {
	//	return nil, 0, errors.Wrap(err, "select query")
	//}

	list := make([]food.AdminGetList, 0)

	rows, err := r.QueryContext(ctx, query)
	if err != nil {
		return nil, 0, web.NewRequestError(errors.Wrap(err, "select foods"), http.StatusInternalServerError)
	}

	for rows.Next() {
		detail := food.AdminGetList{}
		err = rows.Scan(&detail.ID, &detail.Name, &detail.Photos, &detail.Price, &detail.CategoryID, &detail.Category)
		if err != nil {
			return nil, 0, web.NewRequestError(errors.Wrap(err, "scanning foods"), http.StatusBadRequest)
		}
		list = append(list, detail)
	}

	countQuery := fmt.Sprintf(`
		SELECT
			count(%s.id)
		FROM foods
		LEFT JOIN food_category ON foods.category_id = food_category.id
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

	for k, v := range list {
		if v.Photos != nil {
			var photoLink pq.StringArray
			for _, v1 := range *v.Photos {
				baseLink := hashing.GenerateHash(r.ServerBaseUrl, v1)
				photoLink = append(photoLink, baseLink)
			}
			list[k].Photos = &photoLink
		}
	}

	return list, count, nil
}

func (r Repository) AdminGetDetail(ctx context.Context, id int64) (entity.Foods, error) {
	claims, err := r.CheckClaims(ctx, auth.RoleAdmin)
	if err != nil {
		return entity.Foods{}, err
	}

	var detail entity.Foods

	err = r.NewSelect().Model(&detail).Where("id = ? AND restaurant_id = ? AND deleted_at IS NULL", id, claims.RestaurantID).Scan(ctx)
	if err != nil {
		return entity.Foods{}, err
	}

	if detail.Photos != nil {
		var photos pq.StringArray
		for _, v := range *detail.Photos {
			v = hashing.GenerateHash(r.ServerBaseUrl, v)
			photos = append(photos, v)
		}
		detail.Photos = &photos
	}

	return detail, nil
}

func (r Repository) AdminCreate(ctx context.Context, request food.AdminCreateRequest) (food.AdminCreateResponse, error) {
	claims, err := r.CheckClaims(ctx, auth.RoleAdmin)
	if err != nil {
		return food.AdminCreateResponse{}, err
	}

	err = r.ValidateStruct(&request, "Name", "CategoryID")
	if err != nil {
		return food.AdminCreateResponse{}, err
	}

	response := food.AdminCreateResponse{
		Name:         request.Name,
		Photos:       request.PhotosLink,
		Price:        request.Price,
		CategoryID:   request.CategoryID,
		CreatedAt:    time.Now(),
		CreatedBy:    claims.UserId,
		RestaurantID: claims.RestaurantID,
	}

	_, err = r.NewInsert().Model(&response).Exec(ctx)
	if err != nil {
		return food.AdminCreateResponse{}, web.NewRequestError(errors.Wrap(err, "creating food"), http.StatusBadRequest)
	}

	return response, nil
}

func (r Repository) AdminUpdateAll(ctx context.Context, request food.AdminUpdateRequest) error {
	claims, err := r.CheckClaims(ctx, auth.RoleAdmin)
	if err != nil {
		return err
	}

	if err = r.ValidateStruct(&request, "ID", "Name", "Photos", "CategoryID", "Price"); err != nil {
		return err
	}

	q := r.NewUpdate().Table("foods").Where("deleted_at IS NULL AND id = ? AND restaurant_id = ?", request.ID, claims.RestaurantID)
	q.Set("name = ?", request.Name)
	q.Set("photos = array_cat(photos, ?)", request.PhotosLink)
	q.Set("price = ?", request.Price)
	q.Set("category_id = ?", request.CategoryID)
	q.Set("updated_at = ?", time.Now())
	q.Set("updated_by = ?", claims.UserId)

	_, err = q.Exec(ctx)
	if err != nil {
		return web.NewRequestError(errors.Wrap(err, "updating food"), http.StatusBadRequest)
	}

	return nil
}

func (r Repository) AdminUpdateColumns(ctx context.Context, request food.AdminUpdateRequest) error {
	claims, err := r.CheckClaims(ctx, auth.RoleAdmin)
	if err != nil {
		return err
	}

	if err = r.ValidateStruct(&request, "ID"); err != nil {
		return err
	}

	q := r.NewUpdate().Table("foods").Where("deleted_at IS NULL AND id = ? AND restaurant_id = ?",
		request.ID, claims.RestaurantID)

	if request.Name != nil {
		q.Set("name = ?", request.Name)
	}
	if request.PhotosLink != nil {
		q.Set("photos = array_cat(photos, ?)", request.PhotosLink)
	}

	if request.Price != nil {
		q.Set("price = ?", request.Price)
	}

	if request.CategoryID != nil {
		q.Set("category_id = ?", request.CategoryID)
	}
	q.Set("updated_at = ?", time.Now())
	q.Set("updated_by = ?", claims.UserId)

	_, err = q.Exec(ctx)
	if err != nil {
		return web.NewRequestError(errors.Wrap(err, "updating food"), http.StatusBadRequest)
	}

	return nil
}

func (r Repository) AdminDelete(ctx context.Context, id int64) error {
	return r.DeleteRow(ctx, "foods", id, auth.RoleAdmin)
}

func (r Repository) AdminDeleteImage(ctx context.Context, request food.AdminDeleteImageRequest) error {
	_, err := r.CheckClaims(ctx, auth.RoleAdmin)
	if err != nil {
		return err
	}

	if err := r.ValidateStruct(&request,
		"ID",
	); err != nil {
		return err
	}

	var images []string
	err = r.NewSelect().Table("foods").
		Column("photos").
		Where("deleted_at IS NULL AND id = ?", request.ID).
		Scan(ctx,
			pgdialect.Array(&images))

	if err != nil {
		return web.NewRequestError(errors.Wrap(err, "scan photo-stories"), http.StatusBadRequest)
	}

	if request.ImageIndex != nil {
		err = r.DeleteImageIndex(ctx, request.ID, "foods", "photos", images, *request.ImageIndex, auth.RoleAdmin)
		if err != nil {
			return web.NewRequestError(err, http.StatusBadRequest)
		}
	}

	return nil
}

// branch

func (r Repository) BranchGetList(ctx context.Context, filter food.Filter) ([]food.BranchGetList, int, error) {
	claims, err := r.CheckClaims(ctx, auth.RoleBranch)
	if err != nil {
		return nil, 0, err
	}

	whereQuery := fmt.Sprintf(`WHERE f.deleted_at IS NULL AND f.restaurant_id in (select restaurant_id from branches where id = '%d')`, *claims.BranchID)
	if filter.Search != nil {
		whereQuery += fmt.Sprintf(` AND f.name ilike '%s'`, "%"+*filter.Search+"%")
	}

	var limitQuery, offsetQuery string
	if filter.Limit != nil {
		limitQuery = fmt.Sprintf(" LIMIT '%d'", *filter.Limit)
	}
	if filter.Offset != nil {
		offsetQuery = fmt.Sprintf(" OFFSET '%d'", *filter.Offset)
	}

	query := fmt.Sprintf(`
		SELECT 
		    f.id,
		    f.name,
		    f.price,
		    f.photos,
		    f.category_id,
		    fc.name as category
		FROM 
		    foods AS f
		LEFT JOIN food_category AS fc ON f.category_id = fc.id
		%s %s %s
`, whereQuery, limitQuery, offsetQuery)
	if err != nil {
		return nil, 0, errors.Wrap(err, "select query")
	}

	list := make([]food.BranchGetList, 0)

	rows, err := r.QueryContext(ctx, query)
	if err != nil {
		return nil, 0, web.NewRequestError(errors.Wrap(err, "select foods"), http.StatusInternalServerError)
	}

	err = r.ScanRows(ctx, rows, &list)
	if err != nil {
		return nil, 0, web.NewRequestError(errors.Wrap(err, "scanning foods"), http.StatusBadRequest)
	}

	countQuery := fmt.Sprintf(`
		SELECT
			count(f.id)
		FROM 
		    foods AS f
		LEFT JOIN food_category AS fc ON f.category_id = fc.id
		%s
	`, whereQuery)

	countRows, err := r.QueryContext(ctx, countQuery)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, 0, web.NewRequestError(postgres.ErrNotFound, http.StatusBadRequest)
	}
	if err != nil {
		return nil, 0, web.NewRequestError(errors.Wrap(err, "selecting foods"), http.StatusInternalServerError)
	}

	count := 0

	for countRows.Next() {
		if err = countRows.Scan(&count); err != nil {
			return nil, 0, web.NewRequestError(errors.Wrap(err, "scanning user count"), http.StatusInternalServerError)
		}
	}

	for k, v := range list {
		if v.Photos != nil {
			var photoLink pq.StringArray
			for _, v1 := range *v.Photos {
				baseLink := hashing.GenerateHash(r.ServerBaseUrl, v1)
				photoLink = append(photoLink, baseLink)
			}
			list[k].Photos = &photoLink
		}
	}

	return list, count, nil
}

func (r Repository) BranchGetDetail(ctx context.Context, id int64) (entity.Foods, error) {
	claims, err := r.CheckClaims(ctx, auth.RoleBranch)
	if err != nil {
		return entity.Foods{}, err
	}

	var detail entity.Foods

	err = r.NewSelect().Model(&detail).Where("id = ? AND deleted_at IS NULL AND restaurant_id = (select restaurant_id from branches where id = ?)", id, *claims.BranchID).Scan(ctx)

	if err != nil {
		return entity.Foods{}, err
	}

	if detail.Photos != nil {
		var photos pq.StringArray
		for _, v := range *detail.Photos {
			v = hashing.GenerateHash(r.ServerBaseUrl, v)
			photos = append(photos, v)
		}
		detail.Photos = &photos
	}

	return detail, nil
}

func (r Repository) BranchCreate(ctx context.Context, request food.BranchCreateRequest) (food.BranchCreateResponse, error) {
	claims, err := r.CheckClaims(ctx, auth.RoleBranch)
	if err != nil {
		return food.BranchCreateResponse{}, err
	}
	err = r.ValidateStruct(&request, "Name", "CategoryID")
	if err != nil {
		return food.BranchCreateResponse{}, err
	}

	var restaurantID *int64

	err = r.QueryRowContext(ctx, fmt.Sprintf("SELECT restaurant_id FROM branches WHERE id = '%d'", *claims.BranchID)).Scan(&restaurantID)
	if err != nil {
		return food.BranchCreateResponse{}, web.NewRequestError(errors.Wrap(err, "not fount restaurant"), http.StatusBadRequest)
	}

	if restaurantID == nil {
		return food.BranchCreateResponse{}, web.NewRequestError(errors.New("not fount restaurant"), http.StatusBadRequest)
	}

	response := food.BranchCreateResponse{
		Name:         request.Name,
		Photos:       request.PhotosLink,
		Price:        request.Price,
		CategoryID:   request.CategoryID,
		CreatedAt:    time.Now(),
		CreatedBy:    claims.UserId,
		RestaurantID: *restaurantID,
	}

	_, err = r.NewInsert().Model(&response).Exec(ctx)
	if err != nil {
		return food.BranchCreateResponse{}, web.NewRequestError(errors.Wrap(err, "creating food"), http.StatusBadRequest)
	}

	return response, nil
}

func (r Repository) BranchUpdateAll(ctx context.Context, request food.BranchUpdateRequest) error {
	claims, err := r.CheckClaims(ctx, auth.RoleBranch)
	if err != nil {
		return err
	}

	if err = r.ValidateStruct(&request, "ID", "Name", "Photos", "CategoryID", "Price"); err != nil {
		return err
	}

	q := r.NewUpdate().Table("foods").Where("deleted_at IS NULL AND id = ? AND restaurant_id = (select restaurant_id from branches where id = ?)",
		request.ID, *claims.BranchID)

	q.Set("name = ?", request.Name)
	q.Set("photos = array_cat(photos, ?)", request.PhotosLink)
	q.Set("price = ?", request.Price)
	q.Set("category_id", request.CategoryID)
	q.Set("updated_at = ?", time.Now())
	q.Set("updated_by = ?", claims.UserId)

	_, err = q.Exec(ctx)
	if err != nil {
		return web.NewRequestError(errors.Wrap(err, "updating food"), http.StatusBadRequest)
	}

	return nil
}

func (r Repository) BranchUpdateColumns(ctx context.Context, request food.BranchUpdateRequest) error {
	claims, err := r.CheckClaims(ctx, auth.RoleBranch)
	if err != nil {
		return err
	}

	if err = r.ValidateStruct(&request, "ID"); err != nil {
		return err
	}

	q := r.NewUpdate().Table("foods").Where("deleted_at IS NULL AND id = ? AND restaurant_id = (select restaurant_id from branches where id = ?)",
		request.ID, *claims.BranchID)

	if request.Name != nil {
		q.Set("name = ?", request.Name)
	}
	if request.PhotosLink != nil {
		q.Set("photos = array_cat(photos, ?)", request.PhotosLink)
	}

	if request.Price != nil {
		q.Set("price = ?", request.Price)
	}

	if request.CategoryID != nil {
		q.Set("category_id = ?", request.CategoryID)
	}

	q.Set("updated_at = ?", time.Now())
	q.Set("updated_by = ?", claims.UserId)

	_, err = q.Exec(ctx)
	if err != nil {
		return web.NewRequestError(errors.Wrap(err, "updating food"), http.StatusBadRequest)
	}
	if request.Name != nil {
		_, err = r.ExecContext(ctx, fmt.Sprintf(`
			UPDATE branches
    		SET menu_names = (SELECT string_agg( ' {' || text(m.id) || '} ' || m.name, ' ') AS aggregated_names FROM menus m WHERE m.branch_id = branches.id AND m.deleted_at IS NULL AND m.status = 'active')
    		WHERE id in (SELECT branch_id FROM menus WHERE '%d' = any(food_ids) AND status = 'active' AND deleted_at IS NULL)`, request.ID))
		if err != nil {
			return web.NewRequestError(errors.Wrap(err, "updating branch menu_names"), http.StatusBadRequest)
		}
	}
	return nil
}

func (r Repository) BranchDelete(ctx context.Context, id int64) error {
	return r.DeleteRow(ctx, "foods", id, auth.RoleBranch)
}

func (r Repository) BranchDeleteImage(ctx context.Context, request food.AdminDeleteImageRequest) error {
	_, err := r.CheckClaims(ctx, auth.RoleBranch)
	if err != nil {
		return err
	}

	if err := r.ValidateStruct(&request,
		"ID",
	); err != nil {
		return err
	}

	var images []string
	err = r.NewSelect().Table("foods").
		Column("photos").
		Where("deleted_at IS NULL AND id = ?", request.ID).
		Scan(ctx,
			pgdialect.Array(&images))

	if err != nil {
		return web.NewRequestError(errors.Wrap(err, "scan photo-stories"), http.StatusBadRequest)
	}

	if request.ImageIndex != nil {
		err = r.DeleteImageIndex(ctx, request.ID, "foods", "photos", images, *request.ImageIndex, auth.RoleBranch)
		if err != nil {
			return web.NewRequestError(err, http.StatusBadRequest)
		}
	}

	return nil
}

// cashier

func (r Repository) CashierGetList(ctx context.Context, filter food.Filter) ([]food.CashierGetList, int, error) {
	claims, err := r.CheckClaims(ctx, auth.RoleCashier)
	if err != nil {
		return nil, 0, err
	}

	table := "foods"
	whereQuery := fmt.Sprintf(`WHERE %s.deleted_at IS NULL AND %s.restaurant_id in (select restaurant_id from branches where id = '%d')`, table, table, *claims.BranchID)
	countWhereQuery := whereQuery
	if filter.Search != nil {
		whereQuery += fmt.Sprintf(` AND foods.name ilike '%s'`, "%"+*filter.Search+"%")
	}

	var limitQuery, offsetQuery string
	if filter.Limit != nil {
		limitQuery = fmt.Sprintf(" LIMIT '%d'", *filter.Limit)
	}
	if filter.Offset != nil {
		offsetQuery = fmt.Sprintf(" OFFSET '%d'", *filter.Offset)
	}

	whereQuery += fmt.Sprintf(" %s %s", limitQuery, offsetQuery)

	query, err := utils.SelectQuery(filter.Fields, filter.Joins, &table, &whereQuery)
	if err != nil {
		return nil, 0, errors.Wrap(err, "select query")
	}

	list := make([]food.CashierGetList, 0)

	rows, err := r.QueryContext(ctx, query)
	if err != nil {
		return nil, 0, web.NewRequestError(errors.Wrap(err, "select foods"), http.StatusInternalServerError)
	}

	err = r.ScanRows(ctx, rows, &list)
	if err != nil {
		return nil, 0, web.NewRequestError(errors.Wrap(err, "scanning foods"), http.StatusBadRequest)
	}

	countQuery := fmt.Sprintf(`
		SELECT
			count(id)
		FROM
		    %s
		%s
	`, table, countWhereQuery)

	countRows, err := r.QueryContext(ctx, countQuery)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, 0, web.NewRequestError(postgres.ErrNotFound, http.StatusBadRequest)
	}
	if err != nil {
		return nil, 0, web.NewRequestError(errors.Wrap(err, "selecting foods"), http.StatusInternalServerError)
	}

	count := 0

	for countRows.Next() {
		if err = countRows.Scan(&count); err != nil {
			return nil, 0, web.NewRequestError(errors.Wrap(err, "scanning user count"), http.StatusInternalServerError)
		}
	}

	for k, v := range list {
		if v.Photos != nil {
			var photoLink pq.StringArray
			for _, v1 := range *v.Photos {
				baseLink := hashing.GenerateHash(r.ServerBaseUrl, v1)
				photoLink = append(photoLink, baseLink)
			}
			list[k].Photos = &photoLink
		}
	}

	return list, count, nil
}

func NewRepository(DB *postgresql.Database) *Repository {
	return &Repository{DB}
}
