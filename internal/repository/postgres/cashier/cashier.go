package cashier

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
	"restu-backend/internal/service/cashier"
	"restu-backend/internal/service/hashing"
	"strings"
	"time"
)

type Repository struct {
	*postgresql.Database
}

// @admin

func (r Repository) AdminGetList(ctx context.Context, filter cashier.Filter) ([]cashier.AdminGetList, int, error) {
	claims, err := r.CheckClaims(ctx, auth.RoleAdmin)
	if err != nil {
		return nil, 0, err
	}

	whereQuery := fmt.Sprintf(`WHERE c.deleted_at IS NULL and b.restaurant_id = %d and c.role='CASHIER'`, *claims.RestaurantID)

	if filter.BranchID != nil {
		whereQuery += fmt.Sprintf(" AND c.branch_id = %d", *filter.BranchID)
	}

	countWhereQuery := whereQuery

	if filter.Limit != nil {
		whereQuery += fmt.Sprintf(" LIMIT '%d'", *filter.Limit)
	}
	if filter.Offset != nil {
		whereQuery += fmt.Sprintf(" OFFSET '%d'", *filter.Offset)
	}

	query := fmt.Sprintf(`
		SELECT 
		    c.id,
		    c.name,
		    c.phone,
		    c.role,
		    c.status,
		    c.photo,
		    c.address,
		    b.name as branch_name
		FROM 
		    users as c
		LEFT OUTER JOIN branches as b ON b.id = c.branch_id
		LEFT OUTER JOIN restaurants as r ON r.id = b.restaurant_id
		%s
	`, whereQuery)

	list := make([]cashier.AdminGetList, 0)

	rows, err := r.QueryContext(ctx, query)
	if err != nil {
		return nil, 0, web.NewRequestError(errors.Wrap(err, "select user"), http.StatusBadRequest)
	}

	err = r.ScanRows(ctx, rows, &list)
	if err != nil {
		return nil, 0, web.NewRequestError(errors.Wrap(err, "scanning cashiers"), http.StatusInternalServerError)
	}

	countQuery := fmt.Sprintf(`
		SELECT
			count(c.id)
		FROM
		    users as c
		LEFT OUTER JOIN branches as b ON b.id = c.branch_id
		LEFT OUTER JOIN restaurants as r ON r.id = b.restaurant_id
		%s
	`, countWhereQuery)

	countRows, err := r.QueryContext(ctx, countQuery)
	if err == sql.ErrNoRows {
		return nil, 0, web.NewRequestError(postgres.ErrNotFound, http.StatusBadRequest)
	}
	if err != nil {
		return nil, 0, web.NewRequestError(errors.Wrap(err, "selecting cashiers"), http.StatusBadRequest)
	}

	count := 0

	for countRows.Next() {
		if err = countRows.Scan(&count); err != nil {
			return nil, 0, web.NewRequestError(errors.Wrap(err, "scanning user count"), http.StatusBadRequest)
		}
	}

	return list, count, nil
}

func (r Repository) AdminGetDetail(ctx context.Context, id int64) (cashier.AdminGetDetail, error) {
	claims, err := r.CheckClaims(ctx, auth.RoleAdmin)
	if err != nil {
		return cashier.AdminGetDetail{}, err
	}
	whereQuery := fmt.Sprintf(`WHERE c.deleted_at IS NULL and b.restaurant_id = %d and c.role='CASHIER' AND c.id = %d`, *claims.RestaurantID, id)

	query := fmt.Sprintf(`
		SELECT
			c.id,
			c.name,
			c.phone,
			to_char(c.birth_date, 'DD.MM.YYYY'),
			c.gender,
			c.role,
			c.photo,
			c.address,
			c.branch_id,
			b.name
		FROM
		    users as c
		LEFT OUTER JOIN branches as b ON b.id = c.branch_id
		LEFT OUTER JOIN restaurants as r ON r.id = b.restaurant_id
	 %s
	`, whereQuery)

	var detail cashier.AdminGetDetail

	err = r.QueryRowContext(ctx, query).Scan(
		&detail.ID,
		&detail.Name,
		&detail.Phone,
		&detail.BirthDate,
		&detail.Gender,
		&detail.Role,
		&detail.Photo,
		&detail.Address,
		&detail.BranchID,
		&detail.BranchName,
	)

	if errors.Is(err, sql.ErrNoRows) {
		return cashier.AdminGetDetail{}, web.NewRequestError(postgres.ErrNotFound, http.StatusBadRequest)
	}

	if err != nil {
		return cashier.AdminGetDetail{}, web.NewRequestError(errors.Wrap(err, "selecting cashier detail"), http.StatusBadRequest)
	}

	if detail.Photo != nil {
		link := hashing.GenerateHash(r.ServerBaseUrl, *detail.Photo)
		detail.Photo = &link

	}

	return detail, nil
}

func (r Repository) AdminCreate(ctx context.Context, request cashier.AdminCreateRequest) (cashier.AdminCreateResponse, error) {
	claims, err := r.CheckClaims(ctx, auth.RoleAdmin)
	if err != nil {
		return cashier.AdminCreateResponse{}, err
	}

	err = r.ValidateStruct(&request, "Phone", "Name", "BirthDate", "Gender", "Password", "BranchID")
	if err != nil {
		return cashier.AdminCreateResponse{}, err
	}

	role := auth.RoleCashier

	var birthDate time.Time
	if request.BirthDate != nil {
		birthDate, err = time.Parse("02.01.2006", *request.BirthDate)

		if err != nil {
			return cashier.AdminCreateResponse{}, web.NewRequestError(fmt.Errorf("incorrect birth-date format: '%v'", err), http.StatusBadRequest)
		}
	}

	var gender string
	if request.Gender != nil {
		if (strings.ToUpper(*request.Gender) == "M") && (strings.ToUpper(*request.Gender) == "F") {
			return cashier.AdminCreateResponse{}, web.NewRequestError(errors.New("incorrect gender. gender should be M (male) or F (female)"), http.StatusBadRequest)
		}
		gender = strings.ToUpper(*request.Gender)
	}
	status := "active"

	response := cashier.AdminCreateResponse{
		Name:      request.Name,
		Password:  request.Password,
		Phone:     request.Phone,
		BirthDate: &birthDate,
		Gender:    &gender,
		Role:      &role,
		CreatedAt: time.Now(),
		BranchID:  claims.BranchID,
		CreatedBy: claims.UserId,
		Photo:     request.PhotoLink,
		Address:   request.Address,
		Status:    &status,
	}

	_, err = r.NewInsert().Model(&response).Exec(ctx)
	if err != nil {
		return cashier.AdminCreateResponse{}, web.NewRequestError(errors.Wrap(err, "creating user"), http.StatusBadRequest)
	}

	if response.Photo != nil {
		photo := hashing.GenerateHash(r.ServerBaseUrl, *response.Photo)
		response.Photo = &photo
	}

	return response, nil
}

func (r Repository) AdminUpdateAll(ctx context.Context, request cashier.AdminUpdateRequest) error {
	claims, err := r.CheckClaims(ctx, auth.RoleAdmin)
	if err != nil {
		return err
	}

	if err = r.ValidateStruct(&request, "ID", "Name", "Phone", "BirthDate", "Gender", "Password", "BranchID"); err != nil {
		return err
	}

	if (strings.ToUpper(*request.Gender) == "M") && (strings.ToUpper(*request.Gender) == "F") {
		return web.NewRequestError(errors.New("incorrect gender. gender should be M (male) or F (female)"), http.StatusBadRequest)
	}

	var birthDate time.Time
	if request.BirthDate != nil {
		birthDate, err = time.Parse("02.01.2006", *request.BirthDate)
		if err != nil {
			return web.NewRequestError(errors.New("incorrect birth_date format"), http.StatusBadRequest)
		}
	}

	q := r.NewUpdate().Table("users").Where("deleted_at IS NULL AND id = ? AND role='CASHIER'", request.ID)

	q.Set("name = ?", request.Name)
	q.Set("birth_date = ?", birthDate)
	q.Set("gender = ?", strings.ToUpper(*request.Gender))
	q.Set("updated_at = ?", time.Now())
	q.Set("updated_by = ?", claims.UserId)
	q.Set("photo = ?", request.PhotoLink)
	q.Set("address = ?", request.Address)
	q.Set("branch_id = ?", request.BranchID)

	_, err = q.Exec(ctx)
	if err != nil {
		return web.NewRequestError(errors.Wrap(err, "updating cashier"), http.StatusBadRequest)
	}

	return nil
}

func (r Repository) AdminUpdateColumns(ctx context.Context, request cashier.AdminUpdateRequest) error {
	claims, err := r.CheckClaims(ctx, auth.RoleAdmin)
	if err != nil {
		return err
	}

	if err = r.ValidateStruct(&request, "ID"); err != nil {
		return err
	}

	if (request.Gender != nil) && (strings.ToUpper(*request.Gender) == "M") && (strings.ToUpper(*request.Gender) == "F") {
		return web.NewRequestError(errors.New("incorrect gender. gender should be M (male) or F (female)"), http.StatusBadRequest)
	}

	var birthDate time.Time
	if request.BirthDate != nil {
		birthDate, err = time.Parse("02.01.2006", *request.BirthDate)
		if err != nil {
			return web.NewRequestError(errors.New("incorrect birth_date format"), http.StatusBadRequest)
		}
	}

	q := r.NewUpdate().Table("users").Where("deleted_at IS NULL AND id = ? AND role='CASHIER'", request.ID)

	if request.Name != nil {
		q.Set("name = ?", request.Name)
	}
	if request.BirthDate != nil {
		q.Set("birth_date = ?", birthDate)
	}
	if request.Gender != nil {
		q.Set("gender = ?", strings.ToUpper(*request.Gender))
	}

	if request.PhotoLink != nil {
		q.Set("photo = ?", request.PhotoLink)
	}

	if request.Address != nil {
		q.Set("address = ?", request.Address)
	}

	if request.BranchID != nil {
		q.Set("branch_id = ?", request.BranchID)
	}

	q.Set("updated_at = ?", time.Now())
	q.Set("updated_by = ?", claims.UserId)

	_, err = q.Exec(ctx)
	if err != nil {
		return web.NewRequestError(errors.Wrap(err, "updating cashier"), http.StatusBadRequest)
	}

	return nil
}

func (r Repository) AdminDelete(ctx context.Context, id int64) error {
	return r.DeleteRow(ctx, "users", id, auth.RoleAdmin)
}

func (r Repository) AdminUpdateStatus(ctx context.Context, id int64, status string) error {
	claims, err := r.CheckClaims(ctx, auth.RoleAdmin)
	if err != nil {
		return err
	}

	q := r.NewUpdate().Table("users").Where("deleted_at isnull and role='CASHIER'  and id = ?", id)

	q.Set("updated_at = ?", time.Now())
	q.Set("updated_by = ?", claims.UserId)
	if status == "active" || status == "inactive" {
		q.Set("status = ?", status)
	}

	_, err = q.Exec(ctx)
	if err != nil {
		return web.NewRequestError(errors.Wrap(err, "updating cashier status"), http.StatusBadRequest)
	}

	return nil
}

// others

func (r Repository) AdminUpdatePassword(ctx context.Context, request cashier.AdminUpdatePassword) error {
	q := r.NewUpdate().Table("users").Where("deleted_at IS NULL AND id = ? AND role='CASHIER'", request.ID)

	q.Set("password = ?", request.Password)

	_, err := q.Exec(ctx)
	if err != nil {
		return web.NewRequestError(errors.Wrap(err, "updating cashier password"), http.StatusBadRequest)
	}

	return nil
}

func (r Repository) AdminUpdatePhone(ctx context.Context, request cashier.AdminUpdatePhone) error {
	q := r.NewUpdate().Table("users").Where("deleted_at IS NULL AND id = ? AND role='CASHIER'", request.ID)

	q.Set("phone = ?", request.Phone)

	_, err := q.Exec(ctx)
	if err != nil {
		return web.NewRequestError(errors.Wrap(err, "updating cashier phone"), http.StatusBadRequest)
	}

	return nil
}

// @branch

func (r Repository) BranchGetList(ctx context.Context, filter cashier.Filter) ([]cashier.BranchGetList, int, error) {
	claims, err := r.CheckClaims(ctx, auth.RoleBranch)
	if err != nil {
		return nil, 0, err
	}

	whereQuery := fmt.Sprintf(`WHERE c.deleted_at IS NULL AND c.branch_id = '%d' AND c.role='CASHIER'`, *claims.BranchID)

	countWhereQuery := whereQuery

	if filter.Limit != nil {
		whereQuery += fmt.Sprintf(" LIMIT '%d'", *filter.Limit)
	}
	if filter.Offset != nil {
		whereQuery += fmt.Sprintf(" OFFSET '%d'", *filter.Offset)
	}

	query := fmt.Sprintf(`
		SELECT 
		    c.id,
		    c.name,
		    c.phone,
		    c.role,
		    c.status,
		    c.photo,
		    c.address
		FROM 
		    users as c
		%s
	`, whereQuery)

	list := make([]cashier.BranchGetList, 0)

	rows, err := r.QueryContext(ctx, query)
	if err == sql.ErrNoRows {
		return nil, 0, web.NewRequestError(postgres.ErrNotFound, http.StatusBadRequest)
	}
	if err != nil {
		return nil, 0, web.NewRequestError(errors.Wrap(err, "select user"), http.StatusBadRequest)
	}

	err = r.ScanRows(ctx, rows, &list)
	if err != nil {
		return nil, 0, web.NewRequestError(errors.Wrap(err, "scanning cashiers"), http.StatusInternalServerError)
	}

	countQuery := fmt.Sprintf(`
		SELECT
			count(c.id)
		FROM
		    users as c
		%s
	`, countWhereQuery)

	countRows, err := r.QueryContext(ctx, countQuery)
	if err == sql.ErrNoRows {
		return nil, 0, web.NewRequestError(postgres.ErrNotFound, http.StatusBadRequest)
	}
	if err != nil {
		return nil, 0, web.NewRequestError(errors.Wrap(err, "selecting cashiers"), http.StatusInternalServerError)
	}

	for i, v := range list {
		var link string
		if v.Photo != nil {
			link = hashing.GenerateHash(r.ServerBaseUrl, *v.Photo)
		}
		list[i].Photo = &link
	}

	count := 0

	for countRows.Next() {
		if err = countRows.Scan(&count); err != nil {
			return nil, 0, web.NewRequestError(errors.Wrap(err, "scanning user count"), http.StatusBadRequest)
		}
	}

	return list, count, nil
}

func (r Repository) BranchGetDetail(ctx context.Context, id int64) (cashier.BranchGetDetail, error) {
	claims, err := r.CheckClaims(ctx, auth.RoleBranch)
	if err != nil {
		return cashier.BranchGetDetail{}, err
	}

	var data entity.User
	err = r.NewSelect().Model(&data).
		Where("id = ? AND branch_id = ? AND deleted_at IS NULL AND role='CASHIER'", id, claims.BranchID).
		Scan(ctx)
	if err != nil {
		return cashier.BranchGetDetail{}, err
	}
	var detail cashier.BranchGetDetail

	detail.ID = data.ID
	detail.Name = data.Name
	detail.Phone = data.Phone
	detail.Gender = data.Gender
	detail.Role = data.Role
	detail.Address = data.Address

	if data.Photo != nil {
		photo := hashing.GenerateHash(r.ServerBaseUrl, *data.Photo)
		detail.Photo = &photo
	}

	birthDate := data.BirthDate.Format("02.01.2006")
	detail.BirthDate = &birthDate

	return detail, nil
}

func (r Repository) BranchCreate(ctx context.Context, request cashier.BranchCreateRequest) (cashier.BranchCreateResponse, error) {
	claims, err := r.CheckClaims(ctx, auth.RoleBranch)
	if err != nil {
		return cashier.BranchCreateResponse{}, err
	}

	err = r.ValidateStruct(&request, "Phone", "Name", "BirthDate", "Gender", "Password")
	if err != nil {
		return cashier.BranchCreateResponse{}, err
	}

	role := auth.RoleCashier

	var birthDate time.Time
	if request.BirthDate != nil {
		birthDate, err = time.Parse("02.01.2006", *request.BirthDate)

		if err != nil {
			return cashier.BranchCreateResponse{}, web.NewRequestError(fmt.Errorf("incorrect birth-date format: '%v'", err), http.StatusBadRequest)
		}
	}

	var gender string
	if request.Gender != nil {
		if (strings.ToUpper(*request.Gender) == "M") && (strings.ToUpper(*request.Gender) == "F") {
			return cashier.BranchCreateResponse{}, web.NewRequestError(errors.New("incorrect gender. gender should be M (male) or F (female)"), http.StatusBadRequest)
		}
		gender = strings.ToUpper(*request.Gender)
	}

	response := cashier.BranchCreateResponse{
		Name:      request.Name,
		Password:  request.Password,
		Phone:     request.Phone,
		BirthDate: &birthDate,
		Gender:    &gender,
		Role:      &role,
		CreatedAt: time.Now(),
		BranchID:  claims.BranchID,
		CreatedBy: claims.UserId,
		Photo:     request.PhotoLink,
		Address:   request.Address,
	}

	_, err = r.NewInsert().Model(&response).Exec(ctx)
	if err != nil {
		return cashier.BranchCreateResponse{}, web.NewRequestError(errors.Wrap(err, "creating user"), http.StatusBadRequest)
	}

	if response.Photo != nil {
		photo := hashing.GenerateHash(r.ServerBaseUrl, *response.Photo)
		response.Photo = &photo
	}

	return response, nil
}

func (r Repository) BranchUpdateAll(ctx context.Context, request cashier.BranchUpdateRequest) error {
	claims, err := r.CheckClaims(ctx, auth.RoleBranch)
	if err != nil {
		return err
	}

	if err = r.ValidateStruct(&request, "ID", "Name", "Phone", "BirthDate", "Gender", "Password"); err != nil {
		return err
	}

	if (strings.ToUpper(*request.Gender) == "M") && (strings.ToUpper(*request.Gender) == "F") {
		return web.NewRequestError(errors.New("incorrect gender. gender should be M (male) or F (female)"), http.StatusBadRequest)
	}

	var birthDate time.Time
	if request.BirthDate != nil {
		birthDate, err = time.Parse("02.01.2006", *request.BirthDate)
		if err != nil {
			return web.NewRequestError(errors.New("incorrect birth_date format"), http.StatusBadRequest)
		}
	}

	q := r.NewUpdate().Table("users").Where("deleted_at IS NULL AND id = ? AND role='CASHIER'", request.ID)

	q.Set("name = ?", request.Name)
	q.Set("birth_date = ?", birthDate)
	q.Set("gender = ?", strings.ToUpper(*request.Gender))
	q.Set("updated_at = ?", time.Now())
	q.Set("updated_by = ?", claims.UserId)
	q.Set("photo = ?", request.PhotoLink)
	q.Set("address = ?", request.Address)

	_, err = q.Exec(ctx)
	if err != nil {
		return web.NewRequestError(errors.Wrap(err, "updating cashier"), http.StatusBadRequest)
	}

	return nil
}

func (r Repository) BranchUpdateColumns(ctx context.Context, request cashier.BranchUpdateRequest) error {
	claims, err := r.CheckClaims(ctx, auth.RoleBranch)
	if err != nil {
		return err
	}

	if err = r.ValidateStruct(&request, "ID"); err != nil {
		return err
	}

	if (request.Gender != nil) && (strings.ToUpper(*request.Gender) == "M") && (strings.ToUpper(*request.Gender) == "F") {
		return web.NewRequestError(errors.New("incorrect gender. gender should be M (male) or F (female)"), http.StatusBadRequest)
	}

	var birthDate time.Time
	if request.BirthDate != nil {
		birthDate, err = time.Parse("02.01.2006", *request.BirthDate)
		if err != nil {
			return web.NewRequestError(errors.New("incorrect birth_date format"), http.StatusBadRequest)
		}
	}

	q := r.NewUpdate().Table("users").Where("deleted_at IS NULL AND id = ? AND role='CASHIER'", request.ID)

	if request.Name != nil {
		q.Set("name = ?", request.Name)
	}
	if request.BirthDate != nil {
		q.Set("birth_date = ?", birthDate)
	}
	if request.Gender != nil {
		q.Set("gender = ?", strings.ToUpper(*request.Gender))
	}

	if request.PhotoLink != nil {
		q.Set("photo = ?", request.PhotoLink)
	}

	if request.Address != nil {
		q.Set("address = ?", request.Address)
	}

	q.Set("updated_at = ?", time.Now())
	q.Set("updated_by = ?", claims.UserId)

	_, err = q.Exec(ctx)
	if err != nil {
		return web.NewRequestError(errors.Wrap(err, "updating cashier"), http.StatusBadRequest)
	}

	return nil
}

func (r Repository) BranchDelete(ctx context.Context, id int64) error {
	return r.DeleteRow(ctx, "users", id, auth.RoleBranch)
}

func (r Repository) BranchUpdateStatus(ctx context.Context, id int64, status string) error {
	claims, err := r.CheckClaims(ctx, auth.RoleBranch)
	if err != nil {
		return err
	}

	q := r.NewUpdate().Table("users").Where("deleted_at isnull and role='CASHIER' and branch_id = ? and id = ?", *claims.BranchID, id)

	q.Set("updated_at = ?", time.Now())
	q.Set("updated_by = ?", claims.UserId)
	if status == "active" || status == "inactive" {
		q.Set("status = ?", status)
	}

	_, err = q.Exec(ctx)
	if err != nil {
		return web.NewRequestError(errors.Wrap(err, "updating cashier status"), http.StatusBadRequest)
	}

	return nil
}

// others

func (r Repository) UpdatePassword(ctx context.Context, request cashier.BranchUpdatePassword) error {
	q := r.NewUpdate().Table("users").Where("deleted_at IS NULL AND id = ? AND role='CASHIER'", request.ID)

	q.Set("password = ?", request.Password)

	_, err := q.Exec(ctx)
	if err != nil {
		return web.NewRequestError(errors.Wrap(err, "updating cashier password"), http.StatusBadRequest)
	}

	return nil
}

func (r Repository) UpdatePhone(ctx context.Context, request cashier.BranchUpdatePhone) error {
	q := r.NewUpdate().Table("users").Where("deleted_at IS NULL AND id = ? AND role='CASHIER'", request.ID)

	q.Set("phone = ?", request.Phone)

	_, err := q.Exec(ctx)
	if err != nil {
		return web.NewRequestError(errors.Wrap(err, "updating cashier phone"), http.StatusBadRequest)
	}

	return nil
}

func NewRepository(DB *postgresql.Database) *Repository {
	return &Repository{DB}
}
