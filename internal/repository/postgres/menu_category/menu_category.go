package menu_category

import (
	"context"
	"fmt"
	"restu-backend/internal/auth"
	"restu-backend/internal/entity"
	"restu-backend/internal/pkg/repository/postgresql"
	"restu-backend/internal/service/hashing"
	"restu-backend/internal/service/menu_category"
	"time"
)

type Repository struct {
	*postgresql.Database
}

// @admin

func (r Repository) AdminCreate(ctx context.Context, data menu_category.AdminCreateRequest) (*entity.MenuCategory, error) {
	claims, err := r.CheckClaims(ctx, auth.RoleAdmin)
	if err != nil {
		return nil, err
	}

	timeNow := time.Now()

	detail := entity.MenuCategory{
		Name:         data.Name,
		RestaurantId: claims.RestaurantID,
		Logo:         data.LogoLink,
		CreatedAt:    &timeNow,
	}

	_, err = r.NewInsert().Model(&detail).Exec(ctx)
	if err != nil {
		return nil, err
	}

	return &detail, nil
}

func (r Repository) AdminUpdate(ctx context.Context, data menu_category.AdminUpdateRequest) error {
	claims, err := r.CheckClaims(ctx, auth.RoleAdmin)
	if err != nil {
		return err
	}

	q := r.NewUpdate().Table("menu_categories").Where("id = ? AND deleted_at isnull AND restaurant_id = ?", data.ID, *claims.RestaurantID)

	if data.Name != nil {
		q.Set("name = ?", data.Name)
	}
	if data.Logo != nil {
		q.Set("logo = ?", data.LogoLink)
	}

	_, err = q.Exec(ctx)
	if err != nil {
		return err
	}

	return nil
}

func (r Repository) AdminGetList(ctx context.Context, filter menu_category.Filter) ([]entity.MenuCategory, int, error) {
	claims, err := r.CheckClaims(ctx, auth.RoleAdmin)
	if err != nil {
		return nil, 0, err
	}

	var list []entity.MenuCategory

	q := r.NewSelect().Model(&list)

	q.Where(" deleted_at isnull AND restaurant_id = ?", *claims.RestaurantID)
	if filter.Limit != nil {
		q.Limit(*filter.Limit)
	}
	if filter.Offset != nil {
		q.Offset(*filter.Offset)
	}

	count, err := q.ScanAndCount(ctx)

	for i := range list {
		if list[i].Logo != nil {
			logo := hashing.GenerateHash(r.ServerBaseUrl, *list[i].Logo)
			list[i].Logo = &logo
		}
	}

	return list, count, err
}

func (r Repository) AdminGetDetail(ctx context.Context, id int64) (*entity.MenuCategory, error) {
	claims, err := r.CheckClaims(ctx, auth.RoleAdmin)
	if err != nil {
		return nil, err
	}

	var detail entity.MenuCategory

	err = r.NewSelect().Model(&detail).Where("id = ? AND deleted_at isnull AND restaurant_id = ?", id, *claims.RestaurantID).Scan(ctx)

	if detail.Logo != nil {
		logo := hashing.GenerateHash(r.ServerBaseUrl, *detail.Logo)
		detail.Logo = &logo
	}

	return &detail, err
}

func (r Repository) AdminDelete(ctx context.Context, id int64) error {
	claims, err := r.CheckClaims(ctx, auth.RoleAdmin)
	if err != nil {
		return err
	}

	_, err = r.NewUpdate().Table("menu_categories").
		Set("deleted_at = ?", time.Now()).
		Where("restaurant_id = ? and id = ?", *claims.RestaurantID, id).Exec(ctx)
	if err != nil {
		return err
	}

	return nil
}

// @branch

func (r Repository) BranchCreate(ctx context.Context, data menu_category.BranchCreateRequest) (*entity.MenuCategory, error) {
	claims, err := r.CheckClaims(ctx, auth.RoleBranch)
	if err != nil {
		return nil, err
	}

	timeNow := time.Now()

	var restaurantID *int64
	query := fmt.Sprintf(`SELECT restaurant_id FROM branches WHERE id='%d'`, *claims.BranchID)
	if err = r.QueryRowContext(ctx, query).Scan(&restaurantID); err != nil || restaurantID == nil {
		return nil, err
	}

	detail := entity.MenuCategory{
		Name:         data.Name,
		RestaurantId: restaurantID,
		Logo:         data.LogoLink,
		CreatedAt:    &timeNow,
	}

	_, err = r.NewInsert().Model(&detail).Exec(ctx)
	if err != nil {
		return nil, err
	}

	return &detail, nil
}

func (r Repository) BranchUpdate(ctx context.Context, data menu_category.BranchUpdateRequest) error {
	claims, err := r.CheckClaims(ctx, auth.RoleBranch)
	if err != nil {
		return err
	}

	q := r.NewUpdate().Table("menu_categories").Where("id = ? AND deleted_at isnull AND restaurant_id = (select restaurant_id from branches where id = ?)", data.ID, *claims.BranchID)

	if data.Name != nil {
		q.Set("name = ?", data.Name)
	}
	if data.Logo != nil {
		q.Set("logo = ?", data.LogoLink)
	}

	_, err = q.Exec(ctx)
	if err != nil {
		return err
	}

	return nil
}

func (r Repository) BranchGetList(ctx context.Context, filter menu_category.Filter) ([]entity.MenuCategory, int, error) {
	claims, err := r.CheckClaims(ctx, auth.RoleBranch)
	if err != nil {
		return nil, 0, err
	}

	var list []entity.MenuCategory

	q := r.NewSelect().Model(&list)

	q.Where(" deleted_at isnull AND restaurant_id = (select restaurant_id from branches where id = ?)", *claims.BranchID)
	if filter.Limit != nil {
		q.Limit(*filter.Limit)
	}
	if filter.Offset != nil {
		q.Offset(*filter.Offset)
	}

	count, err := q.ScanAndCount(ctx)

	for i := range list {
		if list[i].Logo != nil {
			logo := hashing.GenerateHash(r.ServerBaseUrl, *list[i].Logo)
			list[i].Logo = &logo
		}
	}

	return list, count, err
}

func (r Repository) BranchGetDetail(ctx context.Context, id int64) (*entity.MenuCategory, error) {
	claims, err := r.CheckClaims(ctx, auth.RoleBranch)
	if err != nil {
		return nil, err
	}

	var detail entity.MenuCategory

	err = r.NewSelect().Model(&detail).Where("id = ? AND deleted_at isnull AND restaurant_id = (select restaurant_id from branches where id = ?)", id, *claims.BranchID).Scan(ctx)

	if detail.Logo != nil {
		logo := hashing.GenerateHash(r.ServerBaseUrl, *detail.Logo)
		detail.Logo = &logo
	}

	return &detail, err
}

func (r Repository) BranchDelete(ctx context.Context, id int64) error {
	claims, err := r.CheckClaims(ctx, auth.RoleBranch)
	if err != nil {
		return err
	}

	_, err = r.NewUpdate().Table("menu_categories").
		Set("deleted_at = ?", time.Now()).
		Where("restaurant_id = (select restaurant_id from branches where id = ?) and id = ?", *claims.BranchID, id).Exec(ctx)
	if err != nil {
		return err
	}

	return nil
}

// @cashier

func (r Repository) CashierCreate(ctx context.Context, data menu_category.CashierCreateRequest) (*entity.MenuCategory, error) {
	claims, err := r.CheckClaims(ctx, auth.RoleCashier)
	if err != nil {
		return nil, err
	}

	timeNow := time.Now()

	var restaurantID *int64
	query := fmt.Sprintf(`SELECT restaurant_id FROM branches WHERE id='%d'`, *claims.BranchID)
	if err = r.QueryRowContext(ctx, query).Scan(&restaurantID); err != nil || restaurantID == nil {
		return nil, err
	}

	detail := entity.MenuCategory{
		Name:         data.Name,
		RestaurantId: restaurantID,
		Logo:         data.LogoLink,
		CreatedAt:    &timeNow,
	}

	_, err = r.NewInsert().Model(&detail).Exec(ctx)
	if err != nil {
		return nil, err
	}

	return &detail, nil
}

func (r Repository) CashierUpdate(ctx context.Context, data menu_category.CashierUpdateRequest) error {
	claims, err := r.CheckClaims(ctx, auth.RoleCashier)
	if err != nil {
		return err
	}

	q := r.NewUpdate().Table("menu_categories").Where("id = ? AND deleted_at isnull AND restaurant_id = (select restaurant_id from branches where id = ?)", data.ID, *claims.BranchID)

	if data.Name != nil {
		q.Set("name = ?", data.Name)
	}
	if data.Logo != nil {
		q.Set("logo = ?", data.LogoLink)
	}

	_, err = q.Exec(ctx)
	if err != nil {
		return err
	}

	return nil
}

func (r Repository) CashierGetList(ctx context.Context, filter menu_category.Filter) ([]entity.MenuCategory, int, error) {
	claims, err := r.CheckClaims(ctx, auth.RoleCashier)
	if err != nil {
		return nil, 0, err
	}

	var list []entity.MenuCategory

	q := r.NewSelect().Model(&list)

	q.Where(" deleted_at isnull AND restaurant_id = (select restaurant_id from branches where id = ?)", *claims.BranchID)
	if filter.Limit != nil {
		q.Limit(*filter.Limit)
	}
	if filter.Offset != nil {
		q.Offset(*filter.Offset)
	}

	count, err := q.ScanAndCount(ctx)

	for i := range list {
		if list[i].Logo != nil {
			logo := hashing.GenerateHash(r.ServerBaseUrl, *list[i].Logo)
			list[i].Logo = &logo
		}
	}

	return list, count, err
}

func (r Repository) CashierGetDetail(ctx context.Context, id int64) (*entity.MenuCategory, error) {
	claims, err := r.CheckClaims(ctx, auth.RoleCashier)
	if err != nil {
		return nil, err
	}

	var detail entity.MenuCategory

	err = r.NewSelect().Model(&detail).Where("id = ? AND deleted_at isnull AND restaurant_id = (select restaurant_id from branches where id = ?)", id, *claims.BranchID).Scan(ctx)

	if detail.Logo != nil {
		logo := hashing.GenerateHash(r.ServerBaseUrl, *detail.Logo)
		detail.Logo = &logo
	}

	return &detail, err
}

func (r Repository) CashierDelete(ctx context.Context, id int64) error {
	claims, err := r.CheckClaims(ctx, auth.RoleCashier)
	if err != nil {
		return err
	}

	_, err = r.NewUpdate().Table("menu_categories").
		Set("deleted_at = ?", time.Now()).
		Where("restaurant_id = (select restaurant_id from branches where id = ?) and id = ?", *claims.BranchID, id).Exec(ctx)
	if err != nil {
		return err
	}

	return nil
}

func NewRepository(DB *postgresql.Database) *Repository {
	return &Repository{DB}
}
