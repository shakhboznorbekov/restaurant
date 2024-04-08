package user

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
	"restu-backend/internal/pkg/utils"
	"restu-backend/internal/repository/postgres"
	"restu-backend/internal/service/hashing"
	"restu-backend/internal/service/user"
	"strings"
	"time"
)

type Repository struct {
	*postgresql.Database
}

// @super-admin

func (r Repository) SuperAdminGetList(ctx context.Context, filter user.Filter) ([]user.SuperAdminGetList, int, error) {
	_, err := r.CheckClaims(ctx, auth.RoleSuperAdmin)
	if err != nil {
		return nil, 0, err
	}

	table := "users"
	whereQuery := fmt.Sprintf(`WHERE %s.deleted_at IS NULL`, table)

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

	err = r.ScanRows(ctx, rows, &list)
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

func (r Repository) SuperAdminCreate(ctx context.Context, request user.SuperAdminCreateRequest) (user.SuperAdminCreateResponse, error) {
	claims, err := r.CheckClaims(ctx, auth.RoleSuperAdmin)
	if err != nil {
		return user.SuperAdminCreateResponse{}, err
	}

	err = r.ValidateStruct(&request, "Phone", "Name", "BirthDate", "Gender", "Password")
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

	response := user.SuperAdminCreateResponse{
		Name:      request.Name,
		Phone:     request.Phone,
		Password:  request.Password,
		BirthDate: &birthDate,
		Gender:    &gender,
		Role:      request.Role,
		CreatedAt: time.Now(),
		CreatedBy: claims.UserId,
	}

	_, err = r.NewInsert().Model(&response).Exec(ctx)
	if err != nil {
		return user.SuperAdminCreateResponse{}, web.NewRequestError(errors.Wrap(err, "creating user"), http.StatusBadRequest)
	}

	return response, nil
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

// @client

func (r Repository) ClientCreate(ctx context.Context, request user.ClientCreateRequest) (user.ClientCreateResponse, error) {
	err := r.ValidateStruct(&request, "Phone")
	if err != nil {
		return user.ClientCreateResponse{}, err
	}

	role := auth.RoleClient
	request.Role = &role

	status := "active"
	response := user.ClientCreateResponse{
		Phone:     request.Phone,
		Role:      request.Role,
		Status:    &status,
		CreatedAt: time.Now(),
	}

	_, err = r.NewInsert().Model(&response).Exec(ctx)
	if err != nil {
		return user.ClientCreateResponse{}, web.NewRequestError(errors.Wrap(err, "creating user"), http.StatusBadRequest)
	}

	return response, nil
}

func (r Repository) ClientUpdateAll(ctx context.Context, request user.ClientUpdateRequest) error {
	claims, err := r.CheckClaims(ctx, auth.RoleClient)
	if err != nil {
		return err
	}

	var userDetail entity.User

	err = r.ValidateStruct(&request, "Name", "BirthDate", "Gender")
	if err != nil {
		return err
	}

	err = r.NewSelect().Model(&userDetail).Where("id = ? AND deleted_at IS NULL", claims.UserId).Scan(ctx)
	if err != nil && errors.Is(err, sql.ErrNoRows) {
		return postgres.ErrInvalidID
	} else if err != nil {
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

	q := r.NewUpdate().Table("users").Where("id = ? AND deleted_at IS NULL", claims.UserId)

	q.Set("name = ?", request.Name)
	q.Set("birth_date = ?", birthDate)
	q.Set("gender = ?", strings.ToUpper(*request.Gender))
	q.Set("updated_at = ?", time.Now())

	_, err = q.Exec(ctx)
	if err != nil {
		return web.NewRequestError(errors.Wrap(err, "updating client"), http.StatusBadRequest)
	}

	return nil
}

func (r Repository) ClientUpdateColumn(ctx context.Context, request user.ClientUpdateRequest) error {
	claims, err := r.CheckClaims(ctx, auth.RoleClient)
	if err != nil {
		return err
	}

	var userDetail entity.User

	err = r.NewSelect().Model(&userDetail).Where("id = ? AND deleted_at IS NULL", claims.UserId).Scan(ctx)
	if err != nil && errors.Is(err, sql.ErrNoRows) {
		return postgres.ErrInvalidID
	} else if err != nil {
		return err
	}

	q := r.NewUpdate().Table("users").Where("id = ? AND deleted_at IS NULL", claims.UserId)

	if request.Name != nil {
		q.Set("name = ?", request.Name)
	}
	if request.BirthDate != nil {
		var birthDate time.Time
		if request.BirthDate != nil {
			birthDate, err = time.Parse("02.01.2006", *request.BirthDate)
			if err != nil {
				return web.NewRequestError(errors.New("incorrect birth_date format"), http.StatusBadRequest)
			}
		}
		q.Set("birth_date = ?", birthDate)
	}
	if request.Gender != nil {
		if (strings.ToUpper(*request.Gender) == "M") && (strings.ToUpper(*request.Gender) == "F") {
			return web.NewRequestError(errors.New("incorrect gender. gender should be M (male) or F (female)"), http.StatusBadRequest)
		}
		q.Set("gender = ?", strings.ToUpper(*request.Gender))
	}
	q.Set("updated_at = ?", time.Now())

	_, err = q.Exec(ctx)
	if err != nil {
		return web.NewRequestError(errors.Wrap(err, "updating client"), http.StatusBadRequest)
	}

	return nil
}

func (r Repository) GetByPhone(ctx context.Context, phone string) (entity.User, error) {
	var detail entity.User

	err := r.NewSelect().Model(&detail).Where("phone = ? AND deleted_at IS NULL AND (status is null OR status = 'active')", phone).Scan(ctx)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return entity.User{}, web.NewRequestError(errors.New("user not found"), http.StatusBadRequest)
		}
		return entity.User{}, err
	}

	return detail, nil
}

func (r Repository) ClientGetMe(ctx context.Context, id int64) (entity.User, error) {
	_, err := r.CheckClaims(ctx, auth.RoleClient)
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

func (r Repository) ClientDeleteMe(ctx context.Context) error {
	claims, err := r.CheckClaims(ctx, auth.RoleClient)
	if err != nil {
		return err
	}

	return r.DeleteRow(ctx, "users", claims.UserId, auth.RoleClient)
}

func (r Repository) ClientUpdateMePhone(ctx context.Context, newPhone string) error {
	claims, err := r.CheckClaims(ctx, auth.RoleClient)
	if err != nil {
		return err
	}

	q := r.NewUpdate().Table("users").Where("id = ? AND deleted_at IS NULL", claims.UserId)

	q.Set("phone = ?", newPhone)
	q.Set("updated_at = ?", time.Now())

	_, err = q.Exec(ctx)
	if err != nil {
		return web.NewRequestError(errors.Wrap(err, "updating phone"), http.StatusBadRequest)
	}

	return nil
}

// @admin

func (r Repository) AdminGetList(ctx context.Context, filter user.Filter) ([]user.AdminGetList, int, error) {
	//_, err := r.CheckClaims(ctx, auth.RoleAdmin)
	//if err != nil {
	//	return nil, 0, err
	//}

	table := "users"
	whereQuery := fmt.Sprintf(`WHERE %s.deleted_at IS NULL`, table)

	if filter.Name != nil {
		whereQuery += fmt.Sprintf(" AND %s.name ilike '%s'", table, "%"+*filter.Name+"%")
	}

	if filter.Role != nil {
		whereQuery += fmt.Sprintf(" AND %s.role = '%s'", table, strings.ToUpper(*filter.Role))
	}

	if filter.RestaurantID != nil {
		whereQuery += fmt.Sprintf(" AND %s.restaurant_id = %d", table, *filter.RestaurantID)
	}

	if filter.BranchID != nil {
		whereQuery += fmt.Sprintf(" AND %s.branch_id = %d", table, *filter.BranchID)
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
					    id, 
					    name,
					    phone,
					    role,
					    TO_CHAR(birth_date, 'DD.MM.YYYY') as birth_date,
					    gender
					FROM 
					    %s
					%s`, table, whereQuery)

	//query, err := utils.SelectQuery(filter.Fields, filter.Joins, &table, &whereQuery)
	//if err != nil {
	//	return nil, 0, errors.Wrap(err, "select query")
	//}

	list := make([]user.AdminGetList, 0)

	rows, err := r.QueryContext(ctx, query)
	if err != nil {
		return nil, 0, web.NewRequestError(errors.Wrap(err, "select user"), http.StatusInternalServerError)
	}

	err = r.ScanRows(ctx, rows, &list)
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

func (r Repository) AdminGetDetailByRestaurantID(ctx context.Context, restaurantID int64) (entity.User, error) {
	_, err := r.CheckClaims(ctx, auth.RoleSuperAdmin)
	if err != nil {
		return entity.User{}, err
	}

	var detail entity.User

	err = r.NewSelect().Model(&detail).Where("restaurant_id = ? AND role='ADMIN' AND deleted_at IS NULL", restaurantID).Scan(ctx)
	if err != nil {
		return entity.User{}, err
	}

	return detail, nil
}

func (r Repository) AdminUpdateColumnsByRestaurantID(ctx context.Context, request user.AdminUpdateByRestaurantIDRequest) error {
	claims, err := r.CheckClaims(ctx, auth.RoleSuperAdmin)
	if err != nil {
		return err
	}

	if err = r.ValidateStruct(&request, "RestaurantID"); err != nil {
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

	q := r.NewUpdate().Table("users").Where("deleted_at IS NULL AND restaurant_id = ? AND role = 'ADMIN'", request.RestaurantID)

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

func (r Repository) AdminCreate(ctx context.Context, request user.AdminCreateRequest) (user.AdminCreateResponse, error) {
	err := r.ValidateStruct(&request,
		"Name", "Phone", "BirthDate", "Gender", "RestaurantID", "Password")
	if err != nil {
		return user.AdminCreateResponse{}, err
	}

	role := auth.RoleAdmin
	request.Role = &role

	var birthDate time.Time
	if request.BirthDate != nil {
		birthDate, err = time.Parse("02.01.2006", *request.BirthDate)
		if err != nil {
			return user.AdminCreateResponse{}, web.NewRequestError(errors.New("incorrect birth_date format"), http.StatusBadRequest)
		}
	}

	var gender string
	if request.Gender != nil {
		if (strings.ToUpper(*request.Gender) == "M") && (strings.ToUpper(*request.Gender) == "F") {
			return user.AdminCreateResponse{}, web.NewRequestError(errors.New("incorrect gender. gender should be M (male) or F (female)"), http.StatusBadRequest)
		}
		gender = strings.ToUpper(*request.Gender)
	}

	status := "active"
	response := user.AdminCreateResponse{
		Name:         request.Name,
		Phone:        request.Phone,
		Password:     request.Password,
		BirthDate:    &birthDate,
		Gender:       &gender,
		Role:         request.Role,
		RestaurantID: request.RestaurantID,
		Status:       &status,
		CreatedAt:    time.Now(),
		CreatedBy:    request.CreatedBy,
	}

	_, err = r.NewInsert().Model(&response).Exec(ctx)
	if err != nil {
		return user.AdminCreateResponse{}, web.NewRequestError(errors.Wrap(err, "creating user"), http.StatusBadRequest)
	}

	return response, nil
}

func (r Repository) AdminUpdateColumns(ctx context.Context, request user.AdminUpdateRequest) error {
	//claims, err := r.CheckClaims(ctx, auth.RoleAdmin)
	//if err != nil {
	//	return err
	//}

	err := r.ValidateStruct(&request, "ID")
	if err != nil {
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

	q := r.NewUpdate().Table("users").Where("deleted_at IS NULL AND id = ? AND restaurant_id = ?",
		request.ID, request.RestaurantID)

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
	//q.Set("updated_by = ?", request.UpdatedBy)

	_, err = q.Exec(ctx)
	if err != nil {
		return web.NewRequestError(errors.Wrap(err, "updating user"), http.StatusBadRequest)
	}

	return nil
}

// @branch

func (r Repository) BranchCreate(ctx context.Context, request user.BranchCreateRequest) (user.BranchCreateResponse, error) {
	err := r.ValidateStruct(&request,
		"Name", "Phone", "BirthDate", "Gender", "BranchID", "Password")
	if err != nil {
		return user.BranchCreateResponse{}, err
	}

	role := auth.RoleBranch
	request.Role = &role

	var birthDate time.Time
	if request.BirthDate != nil {
		birthDate, err = time.Parse("02.01.2006", *request.BirthDate)
		if err != nil {
			return user.BranchCreateResponse{}, web.NewRequestError(errors.New("incorrect birth_date format"), http.StatusBadRequest)
		}
	}

	var gender string
	if request.Gender != nil {
		if (strings.ToUpper(*request.Gender) == "M") && (strings.ToUpper(*request.Gender) == "F") {
			return user.BranchCreateResponse{}, web.NewRequestError(errors.New("incorrect gender. gender should be M (male) or F (female)"), http.StatusBadRequest)
		}
		gender = strings.ToUpper(*request.Gender)
	}

	status := "active"
	response := user.BranchCreateResponse{
		Name:      request.Name,
		Phone:     request.Phone,
		Password:  request.Password,
		BirthDate: &birthDate,
		Gender:    &gender,
		Role:      request.Role,
		BranchID:  request.BranchID,
		Status:    &status,
		CreatedAt: time.Now(),
		CreatedBy: request.CreatedBy,
	}

	_, err = r.NewInsert().Model(&response).Exec(ctx)
	if err != nil {
		return user.BranchCreateResponse{}, web.NewRequestError(errors.Wrap(err, "creating user"), http.StatusBadRequest)
	}

	return response, nil
}

// @waiter

func (r Repository) WaiterUpdateMePhone(ctx context.Context, newPhone string) error {
	claims, err := r.CheckClaims(ctx, auth.RoleWaiter)
	if err != nil {
		return err
	}

	q := r.NewUpdate().Table("users").Where("id = ? AND deleted_at IS NULL AND role='WAITER'", claims.UserId)

	q.Set("phone = ?", newPhone)
	q.Set("updated_at = ?", time.Now())

	_, err = q.Exec(ctx)
	if err != nil {
		return web.NewRequestError(errors.Wrap(err, "updating phone"), http.StatusBadRequest)
	}

	return nil
}

func (r Repository) WaiterUpdatePassword(ctx context.Context, password string, waiterId int64) error {
	_, err := r.CheckClaims(ctx, auth.RoleWaiter)
	if err != nil {
		return err
	}

	q := r.NewUpdate().Table("users").Where("deleted_at is null and role='WAITER' and id = ? and status='active'", waiterId)

	q.Set("password = ?", password)
	q.Set("updated_at = ?", time.Now())
	q.Set("updated_by = ?", waiterId)

	_, err = q.Exec(ctx)
	if err != nil {
		return web.NewRequestError(errors.Wrap(err, "updating waiter password"), http.StatusBadRequest)
	}

	return nil
}

// others

func (r Repository) IsPhoneExists(ctx context.Context, phone string) (bool, error) {
	exists, err := r.NewSelect().
		Table("users").
		Where("phone = ? and deleted_at IS NULL", phone).
		Exists(ctx)
	return exists, errors.Wrap(err, "phone exists error")
}

func (r Repository) IsWaiterPhoneExists(ctx context.Context, phone string) (bool, error) {
	exists, err := r.NewSelect().
		Table("users").
		Where("phone = ? and deleted_at IS NULL and role = 'WAITER'", phone).
		Exists(ctx)
	return exists, errors.Wrap(err, "phone exists error")
}

func (r Repository) IsSABCPhoneExists(ctx context.Context, phone string) (bool, error) {
	exists, err := r.NewSelect().
		Table("users").
		Where("phone = ? and deleted_at IS NULL and "+
			"role in ('SUPER-ADMIN', 'ADMIN', 'BRANCH', 'CASHIER')", phone).
		Exists(ctx)
	return exists, errors.Wrap(err, "phone exists error")
}

// @cashier

func (r Repository) CashierGetMe(ctx context.Context) (*user.CashierGetMeResponse, error) {
	claims, err := r.CheckClaims(ctx, auth.RoleCashier)
	if err != nil {
		return nil, err
	}

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
								   AND
								     u.role='CASHIER'`, claims.UserId)

	var response user.CashierGetMeResponse
	if err = r.QueryRowContext(ctx, query).Scan(
		&response.Id,
		&response.Name,
		&response.Photo,
		&response.BirthDate,
		&response.Phone,
		&response.Address,
		&response.Gender); err != nil {
		return nil, err
	}

	if response.Photo != nil {
		photo := hashing.GenerateHash(r.ServerBaseUrl, *response.Photo)
		response.Photo = &photo
	}

	return &response, nil
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
