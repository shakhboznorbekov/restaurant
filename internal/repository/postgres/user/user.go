package user

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/pkg/errors"
	"github.com/restaurant/foundation/web"
	"github.com/restaurant/internal/auth"
	"github.com/restaurant/internal/entity"
	"github.com/restaurant/internal/pkg/repository/postgresql"
	"github.com/restaurant/internal/pkg/utils"
	"github.com/restaurant/internal/repository/postgres"
	"github.com/restaurant/internal/service/hashing"
	"github.com/restaurant/internal/service/user"
	"net/http"
	"strings"
	"time"
)

type Repository struct {
	*postgresql.Database
}

//@super-admin

func (r Repository) SuperAdminCreate(ctx context.Context, request user.SuperAdminCreateRequest) (user.SuperAdminCreateResponse, error) {
	claims, err := r.CheckClaims(ctx, auth.RoleSuperAdmin)
	if err != nil {
		return user.SuperAdminCreateResponse{}, err
	}

	err = r.ValidateStruct(&request, "Phone", "Name", "BirthDate", "Gender")
	if err != nil {
		return user.SuperAdminCreateResponse{}, err
	}
	//if strings.Compare(strings.ToUpper(*request.Role), auth.RoleAdmin) == 0 ||
	//	strings.Compare(strings.ToUpper(*request.Role), auth.RoleBranch) == 0 ||
	//	strings.Compare(strings.ToUpper(*request.Role), auth.RoleCashier) == 0 ||
	//	strings.Compare(strings.ToUpper(*request.Role), auth.RoleWaiter) == 0 ||
	//	strings.Compare(strings.ToUpper(*request.Role), auth.RoleClient) == 0 {
	//	return user.SuperAdminCreateResponse{},
	//		web.NewRequestError(
	//			errors.New(
	//				"you have not permission to create user with role: "+
	//					"ADMIN, BRANCH, CASHIER, WAITER, CLIENT",
	//			),
	//			http.StatusBadRequest,
	//		)
	//}

	role := auth.RoleSuperAdmin
	request.Role = &role

	var birthDate time.Time
	if request.BirthDate != nil {
		birthDate, err = time.Parse("02.01.2006", *request.BirthDate)
		if err != nil {
			return user.SuperAdminCreateResponse{}, web.NewRequestError(errors.New("incorrect birth_date format"), http.StatusBadRequest)
		}
	}

	var gender string
	if request.Gender != nil {
		if (strings.ToUpper(*request.Gender) == "M") && (strings.ToUpper(*request.Gender) == "F") {
			return user.SuperAdminCreateResponse{}, web.NewRequestError(errors.New("incorrect gender. gender should be M (male) or F (female)"), http.StatusBadRequest)
		}
		gender = strings.ToUpper(*request.Gender)
	}

	responce := user.SuperAdminCreateResponse{
		Name:      request.Name,
		Phone:     request.Phone,
		BirthDate: &birthDate,
		Gender:    &gender,
		Role:      request.Role,
		CreatedAt: time.Now(),
		CreatedBy: claims.UserId,
	}

	_, err = r.NewInsert().Model(&responce).Exec(ctx)
	if err != nil {
		return user.SuperAdminCreateResponse{}, web.NewRequestError(errors.Wrap(err, "creating user"), http.StatusBadRequest)
	}

	return responce, nil
}

func (r Repository) SuperAdminGetList(ctx context.Context, filter user.Filter) ([]user.SuperAdminGetList, int, error) {
	_, err := r.CheckClaims(ctx, auth.RoleSuperAdmin)
	if err != nil {
		return nil, 0, err
	}
	table := "users"
	whereQuery := fmt.Sprintf(` WHERE %s.deleted_at IS NULL`, table)

	if filter.Name != nil {
		whereQuery += fmt.Sprintf(" AND %s.name ilike '%s'", table, "%"+*filter.Name+"%")
	}

	if filter.Role != nil {
		whereQuery += fmt.Sprintf(" AND %s.role = '%s'", table, strings.ToUpper(*filter.Role))
	}

	if filter.RestaurantID != nil {
		whereQuery += fmt.Sprintf(" AND %s.restaurant_id = %d", table, *filter.RestaurantID)
	}

	countWhereQuery := whereQuery

	if filter.Limit != nil {
		whereQuery += fmt.Sprintf(" LIMIT '%d'", *filter.Limit)
	}
	if filter.Offset != nil {
		whereQuery += fmt.Sprintf(" OFFSET '%d'", *filter.Offset)
	}

	query, err := utils.SelectQuery(filter.Fields, filter.Joins, &table, &whereQuery)
	if err != nil {
		return nil, 0, errors.Wrap(err, "select query")
	}

	list := make([]user.SuperAdminGetList, 0)

	rows, err := r.QueryContext(ctx, query)
	if err != nil {
		return nil, 0, web.NewRequestError(errors.Wrap(err, "select user"), http.StatusInternalServerError)
	}

	err = r.ScanRow(ctx, rows, &list)
	if err != nil {
		return nil, 0, web.NewRequestError(errors.Wrap(err, "scanning users"), http.StatusBadRequest)
	}

	//countQuery := fmt.Sprintf(`
	//	SELECT
	//		count(id)
	//	FROM
	//	    %s
	//	%s
	//`, table, countWhereQuery)
	//
	//countRows, err := r.QueryContext(ctx, countQuery)
	//if err == sql.ErrNoRows {
	//	return nil, 0, web.NewRequestError(postgres.ErrNotFound, http.StatusNotFound)
	//}
	//if err != nil {
	//	return nil, 0, web.NewRequestError(errors.Wrap(err, "selecting users"), http.StatusBadRequest)
	//}
	//
	//count := 0
	//
	//for countRows.Next() {
	//	if err = countRows.Scan(&count); err != nil {
	//		return nil, 0, web.NewRequestError(errors.Wrap(err, "scanning user count"), http.StatusBadRequest)
	//	}
	//}

	count, err := utils.CountQuery(ctx, r.Database, utils.Count{TableName: &table, WhereQuery: &countWhereQuery})
	if err != nil {
		return nil, 0, errors.Wrap(err, "user count")
	}

	return list, count, nil
}

func (r Repository) SuperAdminGetDetail(ctx context.Context, id int64) (entity.User, error) {
	_, err := r.CheckClaims(ctx, auth.RoleSuperAdmin)
	if err != nil {
		return entity.User{}, err
	}
	var detail entity.User

	err = r.NewSelect().Model(&detail).Where("id = ? AND deleted_at IS NULL", id).Scan(ctx)
	if err != nil {
		return entity.User{}, err
	}

	return detail, nil
}

func (r Repository) SuperAdminUpdateAll(ctx context.Context, request user.SuperAdminUpdateRequest) error {
	claims, err := r.CheckClaims(ctx, auth.RoleSuperAdmin)
	if err != nil {
		return err
	}

	if err = r.ValidateStruct(&request, "ID", "Name", "Phone", "BirthDate", "Gender"); err != nil {
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

	q := r.NewUpdate().Table("users").Where("deleted_at IS NULL AND id = ?", request.ID)

	q.Set("name = ?", request.Name)
	q.Set("phone = ?", request.Phone)
	q.Set("birth_date = ?", birthDate)
	q.Set("gender = ?", strings.ToUpper(*request.Gender))
	q.Set("updated_at = ?", time.Now())
	q.Set("updated_by = ?", claims.UserId)

	_, err = q.Exec(ctx)
	if err != nil {
		return web.NewRequestError(errors.Wrap(err, "updating user"), http.StatusBadRequest)
	}

	return nil
}

func (r Repository) SuperAdminUpdateColumns(ctx context.Context, request user.SuperAdminUpdateRequest) error {
	claims, err := r.CheckClaims(ctx, auth.RoleSuperAdmin)
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

	q := r.NewUpdate().Table("users").Where("deleted_at IS NULL AND id = ?", request.ID)

	if request.Name != nil {
		q.Set("name = ?", request.Name)
	}
	if request.Phone != nil {
		q.Set("phone = ?", request.Phone)
	}
	if request.BirthDate != nil {
		q.Set("birth_date = ?", birthDate)
	}
	if request.Gender != nil {
		q.Set("gender = ?", strings.ToUpper(*request.Gender))
	}

	q.Set("updated_at = ?", time.Now())
	q.Set("updated_by = ?", claims.UserId)

	_, err = q.Exec(ctx)
	if err != nil {
		return web.NewRequestError(errors.Wrap(err, "updating user"), http.StatusBadRequest)
	}

	return nil
}

func (r Repository) SuperAdminDelete(ctx context.Context, id int64) error {
	return r.DeleteRow(ctx, "users", id, auth.RoleSuperAdmin)
}

// others

func (r Repository) IsPhoneExists(ctx context.Context, phone string) (bool, error) {
	exists, err := r.NewSelect().Table("users").Where("phone = ? and deleted_at IS NULL", phone).Exists(ctx)
	return exists, errors.Wrap(err, "phone exists error")
}

func (r Repository) IsWaiterPhoneExists(ctx context.Context, phone string) (bool, error) {
	exists, err := r.NewSelect().Table("users").Where("phone = ? and deleted_at IS NULL and role = 'WAITER'", phone).Exists(ctx)
	return exists, errors.Wrap(err, "phone exists error")
}

func (r Repository) GetMe(ctx context.Context, userID int64) (*user.GetMeResponse, error) {

	query := fmt.Sprintf(`SELECT
   									u.id,
   									u.name,
   									u.photo,
   									TO_CHAR(u.birth_date, 'DD.MM.YYYY'),
   									u.phone,
   									u.address,
   									u.gender
								 FROM users u
								 WHERE
								     u.id = '%d'
								   AND
								     u.deleted_at IS NULL
								  `, userID)

	var response user.GetMeResponse
	err := r.QueryRowContext(ctx, query).Scan(
		&response.Id,
		&response.Name,
		&response.Photo,
		&response.BirthDate,
		&response.Phone,
		&response.Address,
		&response.Gender)

	if errors.Is(err, sql.ErrNoRows) {
		return nil, web.NewRequestError(postgres.ErrNotFound, http.StatusBadRequest)
	}

	if err != nil {
		return nil, web.NewRequestError(errors.Wrap(err, "selecting user info"), http.StatusBadRequest)
	}

	if response.Photo != nil {
		photo := hashing.GenerateHash(r.ServerBaseUrl, *response.Photo)
		response.Photo = &photo
	}

	return &response, nil
}

func NewRepository(DB *postgresql.Database) *Repository {
	return &Repository{DB}
}
