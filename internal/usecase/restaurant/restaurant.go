package restaurant

import (
	"bytes"
	"context"
	"fmt"
	"github.com/pkg/errors"
	"github.com/restaurant/internal/pkg/config"
	"github.com/restaurant/internal/pkg/file"
	"github.com/restaurant/internal/pkg/utils"
	"github.com/restaurant/internal/service/branch"
	"github.com/restaurant/internal/service/branchReview"
	"github.com/restaurant/internal/service/printers"
	"github.com/restaurant/internal/service/restaurant"
	"github.com/restaurant/internal/service/restaurant_category"
	"github.com/restaurant/internal/service/service_percentage"
	"github.com/restaurant/internal/service/tables"
	"github.com/restaurant/internal/service/user"
	"github.com/restaurant/internal/utils/location"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
)

type UseCase struct {
	restaurant         Restaurant
	user               User
	restaurantCategory RestaurantCategory
	branch             Branch
	table              Table
	branchReview       BranchReview
	printers           Printers
	servicePercentage  ServicePercentage
}

func NewUseCase(restaurant Restaurant, user User, restaurantCategory RestaurantCategory, branch Branch, table Table, branchReview BranchReview, printers Printers, servicePercentage ServicePercentage) *UseCase {
	return &UseCase{restaurant, user, restaurantCategory, branch, table, branchReview, printers, servicePercentage}
}

// #resaturant

func (uu UseCase) SuperAdminGetRestaurantList(ctx context.Context, filter restaurant.Filter) ([]restaurant.SuperAdminGetList, int, error) {
	fields := make(map[string][]string)
	fields["restaurants"] = []string{"id", "name", "logo", "mini_logo"}

	filter.Fields = fields

	list, count, err := uu.restaurant.SuperAdminGetList(ctx, filter)
	if err != nil {
		return nil, 0, err
	}

	return list, count, err
}

func (uu UseCase) SuperAdminGetRestaurantDetail(ctx context.Context, id int64) (restaurant.SuperAdminGetDetail, error) {
	var detail restaurant.SuperAdminGetDetail

	data, err := uu.restaurant.SuperAdminGetDetail(ctx, id)
	if err != nil {
		return restaurant.SuperAdminGetDetail{}, err
	}

	detail.ID = data.ID
	detail.Name = data.Name
	detail.Logo = data.Logo

	userList, _, err := uu.user.AdminGetList(ctx, user.Filter{RestaurantID: &id})
	if err != nil {
		return restaurant.SuperAdminGetDetail{}, errors.Wrap(err, "user get list")
	}
	if len(userList) > 0 {
		detail.User.ID = userList[0].ID
		detail.User.Name = userList[0].Name
		detail.User.Phone = userList[0].Phone
		detail.User.Role = userList[0].Role
		detail.User.Gender = userList[0].Gender
		detail.User.BirthDate = userList[0].BirthDate
	}

	//user, err := uu.user.AdminGetDetailByRestaurantID(ctx, data.ID)
	//if err != nil && !errors.Is(err, sql.ErrNoRows) {
	//	return restaurant.SuperAdminGetBannerDetail{}, err
	//} else if errors.Is(err, sql.ErrNoRows) {
	//	detail.User = restaurant.User{}
	//} else {
	//	detail.User.Name = user.Name
	//	detail.User.Phone = user.Phone
	//	if user.BirthDate != nil {
	//		birthDate := user.BirthDate.Format("02.01.2006 15:04")
	//		detail.User.BirthDate = &birthDate
	//	}
	//	detail.User.Gender = user.Gender
	//}

	return detail, nil
}

func (uu UseCase) SuperAdminCreateRestaurant(ctx context.Context, data restaurant.SuperAdminCreateRequest) (restaurant.SuperAdminCreateResponse, error) {
	if data.User.Phone == nil {
		return restaurant.SuperAdminCreateResponse{}, errors.New("user.phone required")
	}
	exists, err := uu.user.IsPhoneExists(ctx, *data.User.Phone)
	if err != nil {
		return restaurant.SuperAdminCreateResponse{}, err
	}

	if exists {
		return restaurant.SuperAdminCreateResponse{}, errors.New("user.phone already exists")
	}

	if data.Logo != nil {
		imageLink, _, err := file.UploadSingle(data.Logo, "restaurant")
		if err != nil {
			return restaurant.SuperAdminCreateResponse{}, errors.Wrap(err, "logo upload")
		}
		data.LogoLink = &imageLink
	}

	detail, err := uu.restaurant.SuperAdminCreate(ctx, data)
	if err != nil {
		return restaurant.SuperAdminCreateResponse{}, err
	}

	data.User.RestaurantID = &detail.ID
	data.User.CreatedBy = detail.CreatedBy
	_, err = uu.user.AdminCreate(ctx, data.User)
	if err != nil {
		return restaurant.SuperAdminCreateResponse{}, err
	}

	return detail, err
}

func (uu UseCase) SuperAdminUpdateRestaurant(ctx context.Context, data restaurant.SuperAdminUpdateRequest) error {
	if data.User.Phone != nil {
		exists, err := uu.user.IsPhoneExists(ctx, *data.User.Phone)
		if err != nil {
			return err
		}

		if exists {
			return errors.New("phone already exists")
		}

		data.User.RestaurantID = &data.ID
		err = uu.user.AdminUpdateColumns(ctx, *data.User)
		if err != nil {
			return err
		}
	}

	return uu.restaurant.SuperAdminUpdateAll(ctx, data)
}

func (uu UseCase) SuperAdminUpdateRestaurantColumn(ctx context.Context, data restaurant.SuperAdminUpdateRequest) error {
	if data.User.Phone != nil {
		exists, err := uu.user.IsPhoneExists(ctx, *data.User.Phone)
		if err != nil {
			return err
		}

		if exists {
			return errors.New("phone already exists")
		}
	}

	if data.User != nil {
		data.User.RestaurantID = &data.ID
		err := uu.user.AdminUpdateColumns(ctx, *data.User)
		if err != nil {
			return err
		}
	}

	return uu.restaurant.SuperAdminUpdateColumns(ctx, data)
}

func (uu UseCase) SuperAdminDeleteRestaurant(ctx context.Context, id int64) error {
	return uu.restaurant.SuperAdminDelete(ctx, id)
}

// @site #restaurant

func (uu UseCase) SiteGetRestaurantList(ctx context.Context) ([]restaurant.SiteGetListResponse, int, error) {
	return uu.restaurant.SiteGetList(ctx)
}

// restaurant_category

// -@super-admin

func (uu UseCase) SuperAdminGetRestaurantCategoryList(ctx context.Context, filter restaurant_category.Filter) ([]restaurant_category.SuperAdminGetList, int, error) {
	fields := make(map[string][]string)
	fields["restaurant_category"] = []string{"id", "name"}
	filter.Fields = fields

	list, count, err := uu.restaurantCategory.SuperAdminGetList(ctx, filter)
	if err != nil {
		return nil, 0, err
	}

	return list, count, err
}

func (uu UseCase) SuperAdminGetRestaurantCategoryDetail(ctx context.Context, id int64) (restaurant_category.SuperAdminGetDetail, error) {
	var detail restaurant_category.SuperAdminGetDetail

	data, err := uu.restaurantCategory.SuperAdminGetDetail(ctx, id)
	if err != nil {
		return restaurant_category.SuperAdminGetDetail{}, err
	}

	detail.ID = data.ID
	detail.Name = data.Name

	return detail, nil
}

func (uu UseCase) SuperAdminCreateRestaurantCategory(ctx context.Context, data restaurant_category.SuperAdminCreateRequest) (restaurant_category.SuperAdminCreateResponse, error) {
	detail, err := uu.restaurantCategory.SuperAdminCreate(ctx, data)
	if err != nil {
		return restaurant_category.SuperAdminCreateResponse{}, err
	}

	return detail, err
}

func (uu UseCase) SuperAdminUpdateRestaurantCategory(ctx context.Context, data restaurant_category.SuperAdminUpdateRequest) error {
	return uu.restaurantCategory.SuperAdminUpdateAll(ctx, data)
}

func (uu UseCase) SuperAdminUpdateRestaurantCategoryColumn(ctx context.Context, data restaurant_category.SuperAdminUpdateRequest) error {
	return uu.restaurantCategory.SuperAdminUpdateColumns(ctx, data)
}

func (uu UseCase) SuperAdminDeleteRestaurantCategory(ctx context.Context, id int64) error {
	return uu.restaurantCategory.SuperAdminDelete(ctx, id)
}

// -@admin

func (uu UseCase) AdminGetRestaurantCategoryList(ctx context.Context, filter restaurant_category.Filter) ([]restaurant_category.AdminGetList, int, error) {
	fields := make(map[string][]string)
	fields["restaurant_category"] = []string{"id", "name"}
	filter.Fields = fields

	list, count, err := uu.restaurantCategory.AdminGetList(ctx, filter)
	if err != nil {
		return nil, 0, err
	}

	return list, count, err
}

// -@site

func (uu UseCase) SiteGetRestaurantCategoryList(ctx context.Context) ([]restaurant_category.SiteGetListResponse, int, error) {
	return uu.restaurantCategory.SiteGetList(ctx)
}

// #branch

// -@admin

func (uu UseCase) AdminGetBranchList(ctx context.Context, filter branch.Filter) ([]branch.AdminGetList, int, error) {
	list, count, err := uu.branch.AdminGetList(ctx, filter)
	if err != nil {
		return nil, 0, err
	}

	return list, count, err
}

func (uu UseCase) AdminGetBranchDetail(ctx context.Context, id int64) (branch.AdminGetDetail, error) {
	var detail branch.AdminGetDetail

	data, err := uu.branch.AdminGetDetail(ctx, id)
	if err != nil {
		return branch.AdminGetDetail{}, err
	}

	detail.ID = data.ID
	detail.Location = data.Location
	detail.Status = data.Status
	detail.WorkTime = data.WorkTime
	detail.Photos = data.Photos
	detail.Name = data.Name
	detail.CategoryID = data.CategoryID

	role := "BRANCH"
	userList, _, err := uu.user.AdminGetList(ctx, user.Filter{BranchID: &id, Role: &role})
	if err != nil {
		return branch.AdminGetDetail{}, errors.Wrap(err, "user get list")
	}

	if len(userList) > 0 {
		detail.User.ID = userList[0].ID
		detail.User.Name = userList[0].Name
		detail.User.Phone = userList[0].Phone
		detail.User.Role = userList[0].Role
		detail.User.Gender = userList[0].Gender
		detail.User.BirthDate = userList[0].BirthDate
	}

	return detail, nil
}

func (uu UseCase) AdminCreateBranch(ctx context.Context, data branch.AdminCreateRequest) (branch.AdminCreateResponse, error) {
	if data.User.Phone == nil {
		return branch.AdminCreateResponse{}, errors.New("user.phone required")
	}
	exists, err := uu.user.IsPhoneExists(ctx, *data.User.Phone)
	if err != nil {
		return branch.AdminCreateResponse{}, err
	}

	if exists {
		return branch.AdminCreateResponse{}, errors.New("user.phone already exists")
	}

	if data.Photos != nil {
		imageLinks, _, err := file.UploadMultiple(data.Photos, "branch")
		if err != nil {
			return branch.AdminCreateResponse{}, errors.Wrap(err, "logo upload")
		}
		data.PhotosLink = &imageLinks
	}

	var servicePercentage service_percentage.AdminCreateResponse
	if data.DefaultServicePercentage != nil {
		servicePercentage, err = uu.servicePercentage.BranchCreate(ctx, service_percentage.AdminCreateRequest{
			Percent: data.DefaultServicePercentage,
		})
		if err != nil {
			return branch.AdminCreateResponse{}, err
		}
	} else {
		var d float64 = 10
		servicePercentage, err = uu.servicePercentage.BranchCreate(ctx, service_percentage.AdminCreateRequest{
			Percent: &d,
		})
		if err != nil {
			return branch.AdminCreateResponse{}, err
		}
	}

	data.DefaultServicePercentageID = &servicePercentage.ID
	detail, err := uu.branch.AdminCreate(ctx, data)
	if err != nil {
		return branch.AdminCreateResponse{}, err
	}

	err = uu.servicePercentage.AdminUpdateBranchID(ctx, service_percentage.AdminUpdateBranchRequest{
		ID:       servicePercentage.ID,
		BranchID: &detail.ID,
	})
	if err != nil {
		return branch.AdminCreateResponse{}, err
	}

	data.User.BranchID = &detail.ID
	data.User.CreatedBy = detail.CreatedBy
	_, err = uu.user.BranchCreate(ctx, data.User)
	if err != nil {
		return branch.AdminCreateResponse{}, err
	}

	return detail, err
}

func (uu UseCase) AdminUpdateBranch(ctx context.Context, data branch.AdminUpdateRequest) error {
	if data.Photos != nil {
		imageLinks, _, err := file.UploadMultiple(data.Photos, "branch")
		if err != nil {
			return errors.Wrap(err, "logo upload")
		}
		data.PhotosLink = &imageLinks
	}

	return uu.branch.AdminUpdateAll(ctx, data)
}

func (uu UseCase) AdminUpdateBranchColumn(ctx context.Context, data branch.AdminUpdateRequest) error {
	if data.Photos != nil {
		imageLinks, _, err := file.UploadMultiple(data.Photos, "branch")
		if err != nil {
			return errors.Wrap(err, "logo upload")
		}
		data.PhotosLink = &imageLinks
	}
	return uu.branch.AdminUpdateColumns(ctx, data)
}

func (uu UseCase) AdminDeleteBranch(ctx context.Context, id int64) error {
	return uu.branch.AdminDelete(ctx, id)
}

func (uu UseCase) AdminDeleteImage(ctx context.Context, request branch.AdminDeleteImageRequest) error {
	return uu.branch.AdminDeleteImage(ctx, request)
}

// -@client

func (uu UseCase) ClientGetBranchList(ctx context.Context, filter branch.Filter) ([]branch.ClientGetList, int, error) {
	return uu.branch.ClientGetList(ctx, filter)
}

func (uu UseCase) ClientGetMapList(ctx context.Context, filter branch.Filter) ([]branch.ClientGetMapList, int, error) {
	return uu.branch.ClientGetMapList(ctx, filter)
}

func (uu UseCase) ClientGetBranchDetail(ctx context.Context, branchFilter branch.DetailFilter) (branch.ClientGetDetail, error) {
	detail, err := uu.branch.ClientGetDetail(ctx, branchFilter.ID)
	if err != nil {
		return branch.ClientGetDetail{}, err
	}

	if branchFilter.Lon != nil && branchFilter.Lat != nil {
		var unit string
		dist := location.CalculateDistance(
			*branchFilter.Lat, *branchFilter.Lon,
			float64(detail.Location["lat"]), float64(detail.Location["lon"]))

		if dist > 1000 {
			unit = "km"
			dist = dist / 1000
		} else {
			unit = "m"
		}
		distance := fmt.Sprintf("%.1f%s", dist, unit)
		detail.Distance = &distance
	}

	return detail, nil
}

func (uu UseCase) ClientGetNearlyBranchList(ctx context.Context, filter branch.Filter) ([]branch.ClientGetList, int, error) {
	list, count, err := uu.branch.ClientNearlyBranchGetList(ctx, filter)
	if err != nil {
		return nil, 0, err
	}

	for k, v := range list {
		var unit string
		if v.Distance != nil {
			dist, err := strconv.ParseFloat(*v.Distance, 64)
			if err != nil {
				return nil, 0, err
			}
			unit = "km"
			dist = dist / 1000

			// use if needed
			//if dist > 1000 {
			//	unit = "km"
			//	dist = dist / 1000
			//} else {
			//	unit = "m"
			//}

			distance := fmt.Sprintf("%.1f%s", dist, unit)
			list[k].Distance = &distance
		}
	}

	return list, count, err
}

func (uu UseCase) ClientUpdateBranchColumn(ctx context.Context, data branch.ClientUpdateRequest) error {
	return uu.branch.ClientUpdateColumns(ctx, data)
}

func (uu UseCase) ClientAddBranchSearchCount(ctx context.Context, branchID int64) error {
	return uu.branch.ClientAddSearchCount(ctx, branchID)
}

func (uu UseCase) ClientGetBranchListOrderSearchCount(ctx context.Context, filter branch.Filter) ([]branch.ClientGetList, int, error) {
	return uu.branch.ClientGetListOrderSearchCount(ctx, filter)
}

func (uu UseCase) ClientBranchListByCategoryID(ctx context.Context, filter branch.Filter, CategoryID int64) ([]branch.ClientGetList, int, error) {
	return uu.branch.ClientGetListByCategoryID(ctx, filter, CategoryID)
}

// -@branch

func (uu UseCase) BranchGetBranchToken(ctx context.Context) (branch.BranchGetToken, error) {
	return uu.branch.BranchGetToken(ctx)
}

// -@ws

func (uu UseCase) WsGetBranchByToken(ctx context.Context, token string) (branch.WsGetByTokenResponse, error) {
	return uu.branch.WsGetByToken(ctx, token)
}

func (uu UseCase) WsBranchUpdateTokenExpiredAt(ctx context.Context, id int64) (string, error) {
	return uu.branch.WsUpdateTokenExpiredAt(ctx, id)
}

// #table

// @admin

func (uu UseCase) AdminGetTableList(ctx context.Context, filter tables.Filter) ([]tables.AdminGetList, int, error) {
	m := make(map[string][]string)

	m["tables"] = []string{"id", "number", "status", "capacity", "branch_id"}
	filter.Fields = m

	joinColumn := "id"
	mainColumn := "branch_id"
	joins := make(map[string]utils.Joins)
	joins["branches"] = utils.Joins{JoinColumn: &joinColumn, MainColumn: &mainColumn}
	filter.Joins = joins

	list, count, err := uu.table.AdminGetList(ctx, filter)
	if err != nil {
		return nil, 0, err
	}

	return list, count, err
}

func (uu UseCase) AdminGetTableDetail(ctx context.Context, id int64) (tables.AdminGetDetail, error) {
	var detail tables.AdminGetDetail

	data, err := uu.table.AdminGetDetail(ctx, id)
	if err != nil {
		return tables.AdminGetDetail{}, err
	}

	detail.ID = data.ID
	detail.Number = data.Number
	detail.Status = data.Status
	detail.Capacity = data.Capacity
	detail.BranchID = data.BranchID

	return detail, nil
}

func (uu UseCase) AdminCreateTable(ctx context.Context, data tables.AdminCreateRequest) ([]tables.AdminCreateResponse, error) {
	if data.To == nil {
		data.To = data.From
	}

	response := make([]tables.AdminCreateResponse, 0)

	for i := *data.From; i <= *data.To; i++ {
		data.Number = &i
		detail, err := uu.table.AdminCreate(ctx, data)
		if err != nil {
			return nil, err
		}
		response = append(response, detail)
	}

	return response, nil
}

func (uu UseCase) AdminUpdateTable(ctx context.Context, data tables.AdminUpdateRequest) error {
	return uu.table.AdminUpdateAll(ctx, data)
}

func (uu UseCase) AdminUpdateTableColumn(ctx context.Context, data tables.AdminUpdateRequest) error {
	return uu.table.AdminUpdateColumns(ctx, data)
}

func (uu UseCase) AdminDeleteTable(ctx context.Context, id int64) error {
	return uu.table.AdminDelete(ctx, id)
}

// @waiter

func (uu UseCase) WaiterGetTableList(ctx context.Context, filter tables.Filter) ([]tables.WaiterGetListResponse, int, error) {
	return uu.table.WaiterGetList(ctx, filter)
}

// @branch

func (uu UseCase) BranchGetTableList(ctx context.Context, filter tables.Filter) ([]tables.BranchGetList, int, error) {
	m := make(map[string][]string)
	m["tables"] = []string{"id", "number", "status", "capacity", "branch_id"}
	filter.Fields = m

	joinColumn := "id"
	mainColumn := "branch_id"
	joins := make(map[string]utils.Joins)
	joins["branches"] = utils.Joins{JoinColumn: &joinColumn, MainColumn: &mainColumn}
	filter.Joins = joins

	list, count, err := uu.table.BranchGetList(ctx, filter)
	if err != nil {
		return nil, 0, err
	}

	return list, count, err
}

func (uu UseCase) BranchGetTableDetail(ctx context.Context, id int64) (tables.BranchGetDetail, error) {
	var detail tables.BranchGetDetail

	data, err := uu.table.BranchGetDetail(ctx, id)
	if err != nil {
		return tables.BranchGetDetail{}, err
	}

	detail.ID = data.ID
	detail.Number = data.Number
	detail.Status = data.Status
	detail.Capacity = data.Capacity
	detail.BranchID = data.BranchID

	return detail, nil
}

func (uu UseCase) BranchCreateTable(ctx context.Context, data tables.BranchCreateRequest) ([]tables.BranchCreateResponse, error) {
	if data.To == nil {
		data.To = data.From
	}

	response := make([]tables.BranchCreateResponse, 0)

	for i := *data.From; i <= *data.To; i++ {
		data.Number = &i
		detail, err := uu.table.BranchCreate(ctx, data)
		if err != nil {
			return nil, err
		}
		response = append(response, detail)
	}

	return response, nil
}

func (uu UseCase) BranchUpdateTable(ctx context.Context, data tables.BranchUpdateRequest) error {
	return uu.table.BranchUpdateAll(ctx, data)
}

func (uu UseCase) BranchUpdateTableColumn(ctx context.Context, data tables.BranchUpdateRequest) error {
	return uu.table.BranchUpdateColumns(ctx, data)
}

func (uu UseCase) BranchDeleteTable(ctx context.Context, id int64) error {
	return uu.table.BranchDelete(ctx, id)
}

func (uu UseCase) BranchGenerateQRTable(ctx context.Context, data tables.BranchGenerateQRTable) (QRResponse, error) {
	var (
		branchBody []struct {
			Logo     *string `json:"logo"`
			Number   *int    `json:"number"`
			BranchID *int64  `json:"branch_id"`
			TableID  *int64  `json:"table_id"`
		}
		branchDetail branch.BranchGetDetail
	)

	if data.Tables != nil && len(data.Tables) >= 1 {
		for _, v := range data.Tables {
			var body struct {
				Logo     *string `json:"logo"`
				Number   *int    `json:"number"`
				BranchID *int64  `json:"branch_id"`
				TableID  *int64  `json:"table_id"`
			}
			tableDetail, err := uu.table.BranchGetDetail(ctx, v)
			if err != nil {
				return QRResponse{}, err
			}

			if tableDetail.BranchID != nil {
				branchDetail, err = uu.branch.BranchGetDetail(ctx, *tableDetail.BranchID)
				if err != nil {
					return QRResponse{}, err
				}

				//if branchDetail.Logo == nil {
				//	return QRResponse{}, errors.New("branch logo is empty")
				//}
			} else {
				continue
			}

			body.Number = tableDetail.Number
			body.Logo = branchDetail.Logo
			body.BranchID = tableDetail.BranchID
			body.TableID = &tableDetail.ID
			branchBody = append(branchBody, body)
		}

		response, err := sendRequestQRService(request.QRBody{
			BaseUrl: config.NewConfig().PythonBaseURL + "/generate/restu/table",
			Model:   branchBody,
		})
		if err != nil {
			return QRResponse{}, err
		}

		return response, nil
	} else {
		return QRResponse{}, errors.New("tables list is empty")
	}
}

// @branch

func (uu UseCase) CashierGetTableList(ctx context.Context, filter tables.Filter) ([]tables.CashierGetList, int, error) {
	m := make(map[string][]string)
	m["tables"] = []string{"id", "number", "status", "capacity", "branch_id"}
	filter.Fields = m

	joinColumn := "id"
	mainColumn := "branch_id"
	joins := make(map[string]utils.Joins)
	joins["branches"] = utils.Joins{JoinColumn: &joinColumn, MainColumn: &mainColumn}
	filter.Joins = joins

	list, count, err := uu.table.CashierGetList(ctx, filter)
	if err != nil {
		return nil, 0, err
	}

	return list, count, err
}

func (uu UseCase) CashierGetTableDetail(ctx context.Context, id int64) (tables.CashierGetDetail, error) {
	var detail tables.CashierGetDetail

	data, err := uu.table.CashierGetDetail(ctx, id)
	if err != nil {
		return tables.CashierGetDetail{}, err
	}

	detail.ID = data.ID
	detail.Number = data.Number
	detail.Status = data.Status
	detail.Capacity = data.Capacity
	detail.BranchID = data.BranchID

	return detail, nil
}

func (uu UseCase) CashierCreateTable(ctx context.Context, data tables.CashierCreateRequest) ([]tables.CashierCreateResponse, error) {
	if data.To == nil {
		data.To = data.From
	}

	response := make([]tables.CashierCreateResponse, 0)

	for i := *data.From; i <= *data.To; i++ {
		data.Number = &i
		detail, err := uu.table.CashierCreate(ctx, data)
		if err != nil {
			return nil, err
		}
		response = append(response, detail)
	}

	return response, nil
}

func (uu UseCase) CashierUpdateTable(ctx context.Context, data tables.CashierUpdateRequest) error {
	return uu.table.CashierUpdateAll(ctx, data)
}

func (uu UseCase) CashierUpdateTableColumn(ctx context.Context, data tables.CashierUpdateRequest) error {
	return uu.table.CashierUpdateColumns(ctx, data)
}

func (uu UseCase) CashierDeleteTable(ctx context.Context, id int64) error {
	return uu.table.CashierDelete(ctx, id)
}

func (uu UseCase) CashierGenerateQRTable(ctx context.Context, data tables.CashierGenerateQRTable) (QRResponse, error) {
	var (
		branchBody []struct {
			Logo     *string `json:"logo"`
			Number   *int    `json:"number"`
			BranchID *int64  `json:"branch_id"`
			TableID  *int64  `json:"table_id"`
		}
		branchDetail branch.CashierGetDetail
	)

	if data.Tables != nil && len(data.Tables) >= 1 {
		for _, v := range data.Tables {
			var body struct {
				Logo     *string `json:"logo"`
				Number   *int    `json:"number"`
				BranchID *int64  `json:"branch_id"`
				TableID  *int64  `json:"table_id"`
			}
			tableDetail, err := uu.table.CashierGetDetail(ctx, v)
			if err != nil {
				return QRResponse{}, err
			}

			if tableDetail.BranchID != nil {
				branchDetail, err = uu.branch.CashierGetDetail(ctx, *tableDetail.BranchID)
				if err != nil {
					return QRResponse{}, err
				}

				//if branchDetail.Logo == nil {
				//	return QRResponse{}, errors.New("branch logo is empty")
				//}
			} else {
				continue
			}

			body.Number = tableDetail.Number
			body.Logo = branchDetail.Logo
			body.BranchID = tableDetail.BranchID
			body.TableID = &tableDetail.ID
			branchBody = append(branchBody, body)
		}

		response, err := sendRequestQRService(request.QRBody{
			BaseUrl: config.NewConfig().PythonBaseURL + "/generate/restu/table",
			Model:   branchBody,
		})
		if err != nil {
			return QRResponse{}, err
		}

		return response, nil
	} else {
		return QRResponse{}, errors.New("tables list is empty")
	}
}

// #branch_review

// -@client

func (uu UseCase) ClientGetBranchReviewList(ctx context.Context, filter branchReview.Filter) ([]branchReview.ClientGetList, int, error) {
	list, count, err := uu.branchReview.ClientGetList(ctx, filter)
	if err != nil {
		return nil, 0, err
	}

	return list, count, err
}

func (uu UseCase) ClientGetBranchReviewDetail(ctx context.Context, id int64) (branchReview.ClientGetDetail, error) {
	var detail branchReview.ClientGetDetail

	data, err := uu.branchReview.ClientGetDetail(ctx, id)
	if err != nil {
		return branchReview.ClientGetDetail{}, err
	}
	detail.ID = data.ID
	detail.Point = data.Point
	detail.Comment = data.Comment
	detail.Rate = data.Rate
	detail.UserID = data.UserID
	detail.UserName = data.UserName
	detail.BranchID = data.BranchID
	detail.BranchName = data.BranchName

	return detail, nil
}

func (uu UseCase) ClientCreateBranchReview(ctx context.Context, data branchReview.ClientCreateRequest) (branchReview.ClientCreateResponse, error) {
	detail, err := uu.branchReview.ClientCreate(ctx, data)
	if err != nil {
		return branchReview.ClientCreateResponse{}, err
	}

	return detail, err
}

func (uu UseCase) ClientUpdateBranchReview(ctx context.Context, data branchReview.ClientUpdateRequest) error {
	return uu.branchReview.ClientUpdateAll(ctx, data)
}

func (uu UseCase) ClientUpdateBranchReviewColumn(ctx context.Context, data branchReview.ClientUpdateRequest) error {
	return uu.branchReview.ClientUpdateColumns(ctx, data)
}

func (uu UseCase) ClientDeleteBranchReview(ctx context.Context, id int64) error {
	return uu.branchReview.ClientDelete(ctx, id)
}

// #devices

// -@branch

func (uu UseCase) BranchGetPrintersList(ctx context.Context, filter printers.Filter) ([]printers.BranchGetList, int, error) {
	list, count, err := uu.printers.BranchGetList(ctx, filter)
	if err != nil {
		return nil, 0, err
	}

	return list, count, err
}

func (uu UseCase) BranchGetPrintersDetail(ctx context.Context, id int64) (printers.BranchGetDetail, error) {
	return uu.printers.BranchGetDetail(ctx, id)
}

func (uu UseCase) BranchCreatePrinters(ctx context.Context, data printers.BranchCreateRequest) (printers.BranchCreateResponse, error) {
	detail, err := uu.printers.BranchCreate(ctx, data)
	if err != nil {
		return printers.BranchCreateResponse{}, err
	}

	return detail, err
}

func (uu UseCase) BranchUpdatePrinters(ctx context.Context, data printers.BranchUpdateRequest) error {
	return uu.printers.BranchUpdateAll(ctx, data)
}

func (uu UseCase) BranchUpdatePrintersColumn(ctx context.Context, data printers.BranchUpdateRequest) error {
	return uu.printers.BranchUpdateColumns(ctx, data)
}

func (uu UseCase) BranchDeletePrinters(ctx context.Context, id int64) error {
	return uu.printers.BranchDelete(ctx, id)
}

// other

type QRResponse struct {
	Error  string `json:"error"`
	Status bool   `json:"status"`
	Data   string `json:"data"`
}

func sendRequestQRService(body request.QRBody) (QRResponse, error) {
	client := http.Client{}

	data, err := json.Marshal(body.Model)
	if err != nil {
		return QRResponse{}, err
	}

	req, err := http.NewRequest(http.MethodPost, body.BaseUrl, bytes.NewBuffer(data))
	if err != nil {
		return QRResponse{}, errors.Wrap(err, "new request")
	}

	req.Header.Set("Content-Type", "application/json")
	if body.Token != nil {
		req.Header.Set("Authorization", *body.Token)
	}

	res, err := client.Do(req)
	if err != nil {
		return QRResponse{}, errors.Wrap(err, "doing request")
	}

	defer res.Body.Close()

	responseBody, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Fatalf("error reading response body: %v\n", err)
	}

	var response QRResponse
	err = json.Unmarshal(responseBody, &response)
	if err != nil {
		return QRResponse{}, err
	}

	log.Println(res.StatusCode, " -> ", response)

	if res.StatusCode < 200 || res.StatusCode > 299 {
		return QRResponse{}, errors.New("something went wrong")
	}

	return response, nil
}
