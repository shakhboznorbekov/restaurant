package restaurant

import (
	"github.com/restaurant/foundation/web"
	restaurant_service "github.com/restaurant/internal/service/restaurant"
	restaurantCategory_service "github.com/restaurant/internal/service/restaurant_category"
	"github.com/restaurant/internal/usecase/restaurant"
	"net/http"
	"reflect"
)

type Controller struct {
	useCase *restaurant.UseCase
}

func NewController(useCase *restaurant.UseCase) *Controller {
	return &Controller{useCase}
}

// #restaurant ----------------------------------------------------------------------------------------

func (uc Controller) SuperAdminGetRestaurantList(c *web.Context) error {
	var filter restaurant_service.Filter

	if limit, ok := c.GetQueryFunc(reflect.Int, "limit").(*int); ok {
		filter.Limit = limit
	}
	if offset, ok := c.GetQueryFunc(reflect.Int, "offset").(*int); ok {
		filter.Offset = offset
	}
	if name, ok := c.GetQueryFunc(reflect.String, "name").(*string); ok {
		filter.Name = name
	}

	if err := c.ValidQuery(); err != nil {
		return c.RespondError(err)
	}

	list, count, err := uc.useCase.SuperAdminGetRestaurantList(c.Ctx, filter)
	if err != nil {
		return c.RespondError(err)
	}

	return c.Respond(map[string]interface{}{
		"data": map[string]interface{}{
			"results": list,
			"count":   count,
		},
		"status": true,
	}, http.StatusOK)
}

func (uc Controller) SuperAdminGetRestaurantDetail(c *web.Context) error {
	id := c.GetParam(reflect.Int, "id").(int)

	if err := c.ValidParam(); err != nil {
		return c.RespondError(err)
	}

	response, err := uc.useCase.SuperAdminGetRestaurantDetail(c.Ctx, int64(id))
	if err != nil {
		return c.RespondError(err)
	}

	return c.Respond(map[string]interface{}{
		"data":   response,
		"status": true,
	}, http.StatusOK)
}

func (uc Controller) SuperAdminCreateRestaurant(c *web.Context) error {
	var request restaurant_service.SuperAdminCreateRequest

	if err := c.BindFunc(&request, "Name", "CategoryID", "User"); err != nil {
		return c.RespondError(err)
	}

	response, err := uc.useCase.SuperAdminCreateRestaurant(c.Ctx, request)
	if err != nil {
		return c.RespondError(err)
	}

	return c.Respond(map[string]interface{}{
		"data":   response,
		"status": true,
	}, http.StatusOK)
}

func (uc Controller) SuperAdminUpdateRestaurantAll(c *web.Context) error {
	id := c.GetParam(reflect.Int, "id").(int)

	if err := c.ValidParam(); err != nil {
		return c.RespondError(err)
	}

	var request restaurant_service.SuperAdminUpdateRequest

	if err := c.BindFunc(&request, "ID", "Name", "CategoryID", "Logo"); err != nil {
		return c.RespondError(err)
	}

	request.ID = int64(id)

	err := uc.useCase.SuperAdminUpdateRestaurant(c.Ctx, request)
	if err != nil {
		return c.RespondError(err)
	}

	return c.Respond(map[string]interface{}{
		"data":   "ok!",
		"status": true,
	}, http.StatusOK)
}

func (uc Controller) SuperAdminUpdateRestaurantColumns(c *web.Context) error {
	id := c.GetParam(reflect.Int, "id").(int)

	if err := c.ValidParam(); err != nil {
		return c.RespondError(err)
	}

	var request restaurant_service.SuperAdminUpdateRequest

	if err := c.BindFunc(&request, "ID"); err != nil {
		return c.RespondError(err)
	}

	request.ID = int64(id)

	err := uc.useCase.SuperAdminUpdateRestaurantColumn(c.Ctx, request)
	if err != nil {
		return c.RespondError(err)
	}

	return c.Respond(map[string]interface{}{
		"data":   "ok!",
		"status": true,
	}, http.StatusOK)
}

func (uc Controller) SuperAdminDeleteRestaurant(c *web.Context) error {
	id := c.GetParam(reflect.Int, "id").(int)

	if err := c.ValidParam(); err != nil {
		return c.RespondError(err)
	}

	err := uc.useCase.SuperAdminDeleteRestaurant(c.Ctx, int64(id))
	if err != nil {
		return c.RespondError(err)
	}

	return c.Respond(map[string]interface{}{
		"data":   "ok!",
		"status": true,
	}, http.StatusOK)
}

func (uc Controller) SiteGetRestaurantList(c *web.Context) error {
	list, count, err := uc.useCase.SiteGetRestaurantList(c.Ctx)
	if err != nil {
		return c.RespondError(err)
	}

	return c.Respond(map[string]interface{}{
		"data": map[string]interface{}{
			"results": list,
			"count":   count,
		},
		"status": true,
	}, http.StatusOK)
}

// #restaurant_category ----------------------------------------------------------------------------------------

// @super-admin

func (uc Controller) SuperAdminGetRestaurantCategoryList(c *web.Context) error {
	var filter restaurantCategory_service.Filter

	if limit, ok := c.GetQueryFunc(reflect.Int, "limit").(*int); ok {
		filter.Limit = limit
	}
	if offset, ok := c.GetQueryFunc(reflect.Int, "offset").(*int); ok {
		filter.Offset = offset
	}
	if name, ok := c.GetQueryFunc(reflect.String, "name").(*string); ok {
		filter.Name = name
	}

	if err := c.ValidQuery(); err != nil {
		return c.RespondError(err)
	}

	list, count, err := uc.useCase.SuperAdminGetRestaurantCategoryList(c.Ctx, filter)
	if err != nil {
		return c.RespondError(err)
	}

	return c.Respond(map[string]interface{}{
		"data": map[string]interface{}{
			"results": list,
			"count":   count,
		},
		"status": true,
	}, http.StatusOK)
}

func (uc Controller) SuperAdminGetRestaurantCategoryDetail(c *web.Context) error {
	id := c.GetParam(reflect.Int, "id").(int)

	if err := c.ValidParam(); err != nil {
		return c.RespondError(err)
	}

	response, err := uc.useCase.SuperAdminGetRestaurantCategoryDetail(c.Ctx, int64(id))
	if err != nil {
		return c.RespondError(err)
	}

	return c.Respond(map[string]interface{}{
		"data":   response,
		"status": true,
	}, http.StatusOK)
}

func (uc Controller) SuperAdminCreateRestaurantCategory(c *web.Context) error {
	var request restaurantCategory_service.SuperAdminCreateRequest

	if err := c.BindFunc(&request, "Name"); err != nil {
		return c.RespondError(err)
	}

	response, err := uc.useCase.SuperAdminCreateRestaurantCategory(c.Ctx, request)
	if err != nil {
		return c.RespondError(err)
	}

	return c.Respond(map[string]interface{}{
		"data":   response,
		"status": true,
	}, http.StatusOK)
}

func (uc Controller) SuperAdminUpdateRestaurantCategoryAll(c *web.Context) error {
	id := c.GetParam(reflect.Int, "id").(int)

	if err := c.ValidParam(); err != nil {
		return c.RespondError(err)
	}

	var request restaurantCategory_service.SuperAdminUpdateRequest

	if err := c.BindFunc(&request, "ID", "Name"); err != nil {
		return c.RespondError(err)
	}

	request.ID = int64(id)

	err := uc.useCase.SuperAdminUpdateRestaurantCategory(c.Ctx, request)
	if err != nil {
		return c.RespondError(err)
	}

	return c.Respond(map[string]interface{}{
		"data":   "ok!",
		"status": true,
	}, http.StatusOK)
}

func (uc Controller) SuperAdminUpdateRestaurantCategoryColumns(c *web.Context) error {
	id := c.GetParam(reflect.Int, "id").(int)

	if err := c.ValidParam(); err != nil {
		return c.RespondError(err)
	}

	var request restaurantCategory_service.SuperAdminUpdateRequest

	if err := c.BindFunc(&request, "ID"); err != nil {
		return c.RespondError(err)
	}

	request.ID = int64(id)

	err := uc.useCase.SuperAdminUpdateRestaurantCategoryColumn(c.Ctx, request)
	if err != nil {
		return c.RespondError(err)
	}

	return c.Respond(map[string]interface{}{
		"data":   "ok!",
		"status": true,
	}, http.StatusOK)
}

func (uc Controller) SuperAdminDeleteRestaurantCategory(c *web.Context) error {
	id := c.GetParam(reflect.Int, "id").(int)

	if err := c.ValidParam(); err != nil {
		return c.RespondError(err)
	}

	err := uc.useCase.SuperAdminDeleteRestaurantCategory(c.Ctx, int64(id))
	if err != nil {
		return c.RespondError(err)
	}

	return c.Respond(map[string]interface{}{
		"data":   "ok!",
		"status": true,
	}, http.StatusOK)
}

// @admin

func (uc Controller) AdminGetRestaurantCategoryList(c *web.Context) error {
	var filter restaurantCategory_service.Filter

	if limit, ok := c.GetQueryFunc(reflect.Int, "limit").(*int); ok {
		filter.Limit = limit
	}
	if offset, ok := c.GetQueryFunc(reflect.Int, "offset").(*int); ok {
		filter.Offset = offset
	}
	if name, ok := c.GetQueryFunc(reflect.String, "name").(*string); ok {
		filter.Name = name
	}

	if err := c.ValidQuery(); err != nil {
		return c.RespondError(err)
	}

	list, count, err := uc.useCase.AdminGetRestaurantCategoryList(c.Ctx, filter)
	if err != nil {
		return c.RespondError(err)
	}

	return c.Respond(map[string]interface{}{
		"data": map[string]interface{}{
			"results": list,
			"count":   count,
		},
		"status": true,
	}, http.StatusOK)
}

// @site

func (uc Controller) SiteGetRestaurantCategoryList(c *web.Context) error {
	list, count, err := uc.useCase.SiteGetRestaurantCategoryList(c.Ctx)
	if err != nil {
		return c.RespondError(err)
	}

	return c.Respond(map[string]interface{}{
		"data": map[string]interface{}{
			"results": list,
			"count":   count,
		},
		"status": true,
	}, http.StatusOK)
}

// #branch ----------------------------------------------------------------------------------------

// @admin

func (uc Controller) AdminGetBranchList(c *web.Context) error {
	var filter branch_service.Filter

	if limit, ok := c.GetQueryFunc(reflect.Int, "limit").(*int); ok {
		filter.Limit = limit
	}
	if offset, ok := c.GetQueryFunc(reflect.Int, "offset").(*int); ok {
		filter.Offset = offset
	}
	if hasMenu, ok := c.GetQueryFunc(reflect.Bool, "has_menu").(*bool); ok {
		filter.HasMenu = hasMenu
	}

	if err := c.ValidQuery(); err != nil {
		return c.RespondError(err)
	}

	list, count, err := uc.useCase.AdminGetBranchList(c.Ctx, filter)
	if err != nil {
		return c.RespondError(err)
	}

	return c.Respond(map[string]interface{}{
		"data": map[string]interface{}{
			"results": list,
			"count":   count,
		},
		"status": true,
	}, http.StatusOK)
}

func (uc Controller) AdminGetBranchDetail(c *web.Context) error {
	id := c.GetParam(reflect.Int, "id").(int)

	if err := c.ValidParam(); err != nil {
		return c.RespondError(err)
	}

	response, err := uc.useCase.AdminGetBranchDetail(c.Ctx, int64(id))
	if err != nil {
		return c.RespondError(err)
	}

	return c.Respond(map[string]interface{}{
		"data":   response,
		"status": true,
	}, http.StatusOK)
}

func (uc Controller) AdminCreateBranch(c *web.Context) error {
	var request branch_service.AdminCreateRequest

	if err := c.BindFunc(&request, "Location", "WorkTime", "User"); err != nil {
		return c.RespondError(err)
	}

	response, err := uc.useCase.AdminCreateBranch(c.Ctx, request)
	if err != nil {
		return c.RespondError(err)
	}

	return c.Respond(map[string]interface{}{
		"data":   response,
		"status": true,
	}, http.StatusOK)
}

func (uc Controller) AdminUpdateBranchAll(c *web.Context) error {
	id := c.GetParam(reflect.Int, "id").(int)

	if err := c.ValidParam(); err != nil {
		return c.RespondError(err)
	}

	var request branch_service.AdminUpdateRequest

	if err := c.BindFunc(&request, "ID", "Location", "Photos", "Status", "WorkTime"); err != nil {
		return c.RespondError(err)
	}

	request.ID = int64(id)

	err := uc.useCase.AdminUpdateBranch(c.Ctx, request)
	if err != nil {
		return c.RespondError(err)
	}

	return c.Respond(map[string]interface{}{
		"data":   "ok!",
		"status": true,
	}, http.StatusOK)
}

func (uc Controller) AdminUpdateBranchColumns(c *web.Context) error {
	id := c.GetParam(reflect.Int, "id").(int)

	if err := c.ValidParam(); err != nil {
		return c.RespondError(err)
	}

	var request branch_service.AdminUpdateRequest

	if err := c.BindFunc(&request); err != nil {
		return c.RespondError(err)
	}

	request.ID = int64(id)

	err := uc.useCase.AdminUpdateBranchColumn(c.Ctx, request)
	if err != nil {
		return c.RespondError(err)
	}

	return c.Respond(map[string]interface{}{
		"data":   "ok!",
		"status": true,
	}, http.StatusOK)
}

func (uc Controller) AdminDeleteBranch(c *web.Context) error {
	id := c.GetParam(reflect.Int, "id").(int)

	if err := c.ValidParam(); err != nil {
		return c.RespondError(err)
	}

	err := uc.useCase.AdminDeleteBranch(c.Ctx, int64(id))
	if err != nil {
		return c.RespondError(err)
	}

	return c.Respond(map[string]interface{}{
		"data":   "ok!",
		"status": true,
	}, http.StatusOK)
}

func (uc Controller) AdminImageDeleteBranch(c *web.Context) error {
	id := c.GetParam(reflect.Int, "id").(int)

	if err := c.ValidParam(); err != nil {
		return c.RespondError(err)
	}

	var request branch_service.AdminDeleteImageRequest

	if err := c.BindFunc(&request); err != nil {
		return c.RespondError(err)
	}

	request.ID = int64(id)

	err := uc.useCase.AdminDeleteImage(c.Ctx, request)
	if err != nil {
		return c.RespondError(err)
	}

	return c.Respond(map[string]interface{}{
		"data":   "ok!",
		"status": true,
	}, http.StatusOK)
}

// @client

func (uc Controller) ClientGetBranchList(c *web.Context) error {
	var filter branch_service.Filter

	if limit, ok := c.GetQueryFunc(reflect.Int, "limit").(*int); ok {
		filter.Limit = limit
	}
	if offset, ok := c.GetQueryFunc(reflect.Int, "offset").(*int); ok {
		filter.Offset = offset
	}
	if page, ok := c.GetQueryFunc(reflect.Int, "page").(*int); ok {
		filter.Page = page
	}
	if isLiked, ok := c.GetQueryFunc(reflect.Bool, "is-liked").(*bool); ok {
		filter.IsLiked = isLiked
	}

	if err := c.ValidQuery(); err != nil {
		return c.RespondMobileError(err)
	}

	list, count, err := uc.useCase.ClientGetBranchList(c.Ctx, filter)
	if err != nil {
		return c.RespondMobileError(err)
	}

	return c.Respond(map[string]interface{}{
		"error": nil,
		"data": map[string]interface{}{
			"results": list,
			"count":   count,
		},
		"status": true,
	}, http.StatusOK)
}

func (uc Controller) ClientGetMapBranchList(c *web.Context) error {
	var filter branch_service.Filter

	if limit, ok := c.GetQueryFunc(reflect.Int, "limit").(*int); ok {
		filter.Limit = limit
	}
	if offset, ok := c.GetQueryFunc(reflect.Int, "offset").(*int); ok {
		filter.Offset = offset
	}
	if page, ok := c.GetQueryFunc(reflect.Int, "page").(*int); ok {
		filter.Page = page
	}

	if err := c.ValidQuery(); err != nil {
		return c.RespondMobileError(err)
	}

	list, count, err := uc.useCase.ClientGetMapList(c.Ctx, filter)
	if err != nil {
		return c.RespondMobileError(err)
	}

	return c.Respond(map[string]interface{}{
		"error": nil,
		"data": map[string]interface{}{
			"results": list,
			"count":   count,
		},
		"status": true,
	}, http.StatusOK)
}

func (uc Controller) ClientGetBranchDetail(c *web.Context) error {
	var branchFilter branch_service.DetailFilter

	if lon, ok := c.GetQueryFunc(reflect.Float64, "lon").(*float64); ok {
		branchFilter.Lon = lon
	}
	if lat, ok := c.GetQueryFunc(reflect.Float64, "lat").(*float64); ok {
		branchFilter.Lat = lat
	}

	id := c.GetParam(reflect.Int, "id").(int)
	branchFilter.ID = int64(id)

	if err := c.ValidParam(); err != nil {
		return c.RespondMobileError(err)
	}

	response, err := uc.useCase.ClientGetBranchDetail(c.Ctx, branchFilter)
	if err != nil {
		return c.RespondMobileError(err)
	}

	return c.Respond(map[string]interface{}{
		"error":  nil,
		"data":   response,
		"status": true,
	}, http.StatusOK)
}

func (uc Controller) ClientGetNearlyBranchList(c *web.Context) error {
	var filter branch_service.Filter

	if limit, ok := c.GetQueryFunc(reflect.Int, "limit").(*int); ok {
		filter.Limit = limit
	}
	if offset, ok := c.GetQueryFunc(reflect.Int, "offset").(*int); ok {
		filter.Offset = offset
	}
	if page, ok := c.GetQueryFunc(reflect.Int, "page").(*int); ok {
		filter.Page = page
	}

	if lon, ok := c.GetQueryFunc(reflect.Float64, "lon").(*float64); ok {
		filter.Lon = lon
	}
	if lat, ok := c.GetQueryFunc(reflect.Float64, "lat").(*float64); ok {
		filter.Lat = lat
	}

	if err := c.ValidQuery(); err != nil {
		return c.RespondMobileError(err)
	}

	if filter.Lon == nil || filter.Lat == nil {
		return c.RespondMobileError(errors.New("lon or lat is required"))
	}

	list, count, err := uc.useCase.ClientGetNearlyBranchList(c.Ctx, filter)
	if err != nil {
		return c.RespondMobileError(err)
	}

	return c.Respond(map[string]interface{}{
		"error": nil,
		"data": map[string]interface{}{
			"results": list,
			"count":   count,
		},
		"status": true,
	}, http.StatusOK)
}

func (uc Controller) ClientUpdateBranchColumns(c *web.Context) error {
	id := c.GetParam(reflect.Int, "id").(int)

	if err := c.ValidParam(); err != nil {
		return c.RespondMobileError(err)
	}

	var request branch_service.ClientUpdateRequest

	if err := c.BindFunc(&request, "ID"); err != nil {
		return c.RespondMobileError(err)
	}

	if request.ID != int64(id) {
		return c.RespondMobileError(errors.New("id param is not equal to body id"))
	}

	request.ID = int64(id)

	err := uc.useCase.ClientUpdateBranchColumn(c.Ctx, request)
	if err != nil {
		return c.RespondMobileError(err)
	}

	return c.Respond(map[string]interface{}{
		"error":  nil,
		"data":   "ok!",
		"status": true,
	}, http.StatusOK)
}

func (uc Controller) ClientAddBranchSearchCount(c *web.Context) error {
	id := c.GetParam(reflect.Int, "id").(int)

	if err := c.ValidParam(); err != nil {
		return c.RespondMobileError(err)
	}

	err := uc.useCase.ClientAddBranchSearchCount(c.Ctx, int64(id))
	if err != nil {
		return c.RespondMobileError(err)
	}

	return c.Respond(map[string]interface{}{
		"error":  nil,
		"data":   "ok!",
		"status": true,
	}, http.StatusOK)
}

func (uc Controller) ClientGetBranchListOrderSearchCount(c *web.Context) error {
	var filter branch_service.Filter

	if limit, ok := c.GetQueryFunc(reflect.Int, "limit").(*int); ok {
		filter.Limit = limit
	}
	if offset, ok := c.GetQueryFunc(reflect.Int, "offset").(*int); ok {
		filter.Offset = offset
	}
	if page, ok := c.GetQueryFunc(reflect.Int, "page").(*int); ok {
		filter.Page = page
	}

	if err := c.ValidQuery(); err != nil {
		return c.RespondMobileError(err)
	}

	list, count, err := uc.useCase.ClientGetBranchListOrderSearchCount(c.Ctx, filter)
	if err != nil {
		return c.RespondMobileError(err)
	}

	return c.Respond(map[string]interface{}{
		"error": nil,
		"data": map[string]interface{}{
			"results": list,
			"count":   count,
		},
		"status": true,
	}, http.StatusOK)
}

func (uc Controller) ClientGetBranchListByCategoryID(c *web.Context) error {
	categoryId := c.GetParam(reflect.Int, "category_id").(int)

	if err := c.ValidParam(); err != nil {
		return c.RespondMobileError(err)
	}

	var filter branch_service.Filter

	if limit, ok := c.GetQueryFunc(reflect.Int, "limit").(*int); ok {
		filter.Limit = limit
	}
	if offset, ok := c.GetQueryFunc(reflect.Int, "offset").(*int); ok {
		filter.Offset = offset
	}
	if page, ok := c.GetQueryFunc(reflect.Int, "page").(*int); ok {
		filter.Page = page
	}

	if err := c.ValidQuery(); err != nil {
		return c.RespondMobileError(err)
	}

	list, count, err := uc.useCase.ClientBranchListByCategoryID(c.Ctx, filter, int64(categoryId))
	if err != nil {
		return c.RespondMobileError(err)
	}

	return c.Respond(map[string]interface{}{
		"error": nil,
		"data": map[string]interface{}{
			"results": list,
			"count":   count,
		},
		"status": true,
	}, http.StatusOK)
}

// @branch

func (uc Controller) BranchGetBranchToken(c *web.Context) error {

	response, err := uc.useCase.BranchGetBranchToken(c.Ctx)
	if err != nil {
		return c.RespondMobileError(err)
	}

	return c.Respond(map[string]interface{}{
		"data":   response,
		"status": true,
	}, http.StatusOK)
}

// #table ----------------------------------------------------------------------------------------

//@admin

func (uc Controller) AdminGetTableList(c *web.Context) error {
	var filter table_service.Filter

	if limit, ok := c.GetQueryFunc(reflect.Int, "limit").(*int); ok {
		filter.Limit = limit
	}
	if offset, ok := c.GetQueryFunc(reflect.Int, "offset").(*int); ok {
		filter.Offset = offset
	}

	if err := c.ValidQuery(); err != nil {
		return c.RespondError(err)
	}

	list, count, err := uc.useCase.AdminGetTableList(c.Ctx, filter)
	if err != nil {
		return c.RespondError(err)
	}

	return c.Respond(map[string]interface{}{
		"data": map[string]interface{}{
			"results": list,
			"count":   count,
		},
		"status": true,
	}, http.StatusOK)
}

func (uc Controller) AdminGetTableDetail(c *web.Context) error {
	id := c.GetParam(reflect.Int, "id").(int)

	if err := c.ValidParam(); err != nil {
		return c.RespondError(err)
	}

	response, err := uc.useCase.AdminGetTableDetail(c.Ctx, int64(id))
	if err != nil {
		return c.RespondError(err)
	}

	return c.Respond(map[string]interface{}{
		"data":   response,
		"status": true,
	}, http.StatusOK)
}

func (uc Controller) AdminCreateTable(c *web.Context) error {
	var request table_service.AdminCreateRequest

	if err := c.BindFunc(&request, "Number", "Capacity", "BranchID"); err != nil {
		return c.RespondError(err)
	}

	response, err := uc.useCase.AdminCreateTable(c.Ctx, request)
	if err != nil {
		return c.RespondError(err)
	}

	return c.Respond(map[string]interface{}{
		"data":   response,
		"status": true,
	}, http.StatusOK)
}

func (uc Controller) AdminUpdateTableAll(c *web.Context) error {
	id := c.GetParam(reflect.Int, "id").(int)

	if err := c.ValidParam(); err != nil {
		return c.RespondError(err)
	}

	var request table_service.AdminUpdateRequest

	if err := c.BindFunc(&request, "ID", "Number", "Status", "Capacity", "BranchID"); err != nil {
		return c.RespondError(err)
	}

	request.ID = int64(id)

	err := uc.useCase.AdminUpdateTable(c.Ctx, request)
	if err != nil {
		return c.RespondError(err)
	}

	return c.Respond(map[string]interface{}{
		"data":   "ok!",
		"status": true,
	}, http.StatusOK)
}

func (uc Controller) AdminUpdateTableColumns(c *web.Context) error {
	id := c.GetParam(reflect.Int, "id").(int)

	if err := c.ValidParam(); err != nil {
		return c.RespondError(err)
	}

	var request table_service.AdminUpdateRequest

	if err := c.BindFunc(&request, "ID"); err != nil {
		return c.RespondError(err)
	}

	request.ID = int64(id)

	err := uc.useCase.AdminUpdateTableColumn(c.Ctx, request)
	if err != nil {
		return c.RespondError(err)
	}

	return c.Respond(map[string]interface{}{
		"data":   "ok!",
		"status": true,
	}, http.StatusOK)
}

func (uc Controller) AdminDeleteTable(c *web.Context) error {
	id := c.GetParam(reflect.Int, "id").(int)

	if err := c.ValidParam(); err != nil {
		return c.RespondError(err)
	}

	err := uc.useCase.AdminDeleteTable(c.Ctx, int64(id))
	if err != nil {
		return c.RespondError(err)
	}

	return c.Respond(map[string]interface{}{
		"data":   "ok!",
		"status": true,
	}, http.StatusOK)
}

// @branch

func (uc Controller) BranchGetTableList(c *web.Context) error {
	var filter table_service.Filter

	if limit, ok := c.GetQueryFunc(reflect.Int, "limit").(*int); ok {
		filter.Limit = limit
	}
	if offset, ok := c.GetQueryFunc(reflect.Int, "offset").(*int); ok {
		filter.Offset = offset
	}

	if err := c.ValidQuery(); err != nil {
		return c.RespondError(err)
	}

	list, count, err := uc.useCase.BranchGetTableList(c.Ctx, filter)
	if err != nil {
		return c.RespondError(err)
	}

	return c.Respond(map[string]interface{}{
		"data": map[string]interface{}{
			"results": list,
			"count":   count,
		},
		"status": true,
	}, http.StatusOK)
}

func (uc Controller) BranchGetTableDetail(c *web.Context) error {
	id := c.GetParam(reflect.Int, "id").(int)

	if err := c.ValidParam(); err != nil {
		return c.RespondError(err)
	}

	response, err := uc.useCase.BranchGetTableDetail(c.Ctx, int64(id))
	if err != nil {
		return c.RespondError(err)
	}

	return c.Respond(map[string]interface{}{
		"data":   response,
		"status": true,
	}, http.StatusOK)
}

func (uc Controller) BranchCreateTable(c *web.Context) error {
	var request table_service.BranchCreateRequest

	if err := c.BindFunc(&request, "From", "Capacity"); err != nil {
		return c.RespondError(err)
	}

	response, err := uc.useCase.BranchCreateTable(c.Ctx, request)
	if err != nil {
		return c.RespondError(err)
	}

	return c.Respond(map[string]interface{}{
		"data":   response,
		"status": true,
	}, http.StatusOK)
}

func (uc Controller) BranchUpdateTableAll(c *web.Context) error {
	id := c.GetParam(reflect.Int, "id").(int)

	if err := c.ValidParam(); err != nil {
		return c.RespondError(err)
	}

	var request table_service.BranchUpdateRequest

	if err := c.BindFunc(&request, "ID", "Number", "Status", "Capacity"); err != nil {
		return c.RespondError(err)
	}

	request.ID = int64(id)

	err := uc.useCase.BranchUpdateTable(c.Ctx, request)
	if err != nil {
		return c.RespondError(err)
	}

	return c.Respond(map[string]interface{}{
		"data":   "ok!",
		"status": true,
	}, http.StatusOK)
}

func (uc Controller) BranchUpdateTableColumns(c *web.Context) error {
	id := c.GetParam(reflect.Int, "id").(int)

	if err := c.ValidParam(); err != nil {
		return c.RespondError(err)
	}

	var request table_service.BranchUpdateRequest

	if err := c.BindFunc(&request, "ID"); err != nil {
		return c.RespondError(err)
	}

	request.ID = int64(id)

	err := uc.useCase.BranchUpdateTableColumn(c.Ctx, request)
	if err != nil {
		return c.RespondError(err)
	}

	return c.Respond(map[string]interface{}{
		"data":   "ok!",
		"status": true,
	}, http.StatusOK)
}

func (uc Controller) BranchDeleteTable(c *web.Context) error {
	id := c.GetParam(reflect.Int, "id").(int)

	if err := c.ValidParam(); err != nil {
		return c.RespondError(err)
	}

	err := uc.useCase.BranchDeleteTable(c.Ctx, int64(id))
	if err != nil {
		return c.RespondError(err)
	}

	return c.Respond(map[string]interface{}{
		"data":   "ok!",
		"status": true,
	}, http.StatusOK)
}

func (uc Controller) BranchGenerateQRTable(c *web.Context) error {
	var request table_service.BranchGenerateQRTable

	if err := c.BindFunc(&request, "Tables"); err != nil {
		return c.RespondError(err)
	}

	response, err := uc.useCase.BranchGenerateQRTable(c.Ctx, request)
	if err != nil {
		return c.RespondError(err)
	}

	return c.Respond(map[string]interface{}{
		"data":   response,
		"status": true,
	}, http.StatusOK)
}

// @branch

func (uc Controller) CashierGetTableList(c *web.Context) error {
	var filter table_service.Filter

	if limit, ok := c.GetQueryFunc(reflect.Int, "limit").(*int); ok {
		filter.Limit = limit
	}
	if offset, ok := c.GetQueryFunc(reflect.Int, "offset").(*int); ok {
		filter.Offset = offset
	}

	if err := c.ValidQuery(); err != nil {
		return c.RespondError(err)
	}

	list, count, err := uc.useCase.CashierGetTableList(c.Ctx, filter)
	if err != nil {
		return c.RespondError(err)
	}

	return c.Respond(map[string]interface{}{
		"data": map[string]interface{}{
			"results": list,
			"count":   count,
		},
		"status": true,
	}, http.StatusOK)
}

func (uc Controller) CashierGetTableDetail(c *web.Context) error {
	id := c.GetParam(reflect.Int, "id").(int)

	if err := c.ValidParam(); err != nil {
		return c.RespondError(err)
	}

	response, err := uc.useCase.CashierGetTableDetail(c.Ctx, int64(id))
	if err != nil {
		return c.RespondError(err)
	}

	return c.Respond(map[string]interface{}{
		"data":   response,
		"status": true,
	}, http.StatusOK)
}

func (uc Controller) CashierCreateTable(c *web.Context) error {
	var request table_service.CashierCreateRequest

	if err := c.BindFunc(&request, "From", "Capacity"); err != nil {
		return c.RespondError(err)
	}

	response, err := uc.useCase.CashierCreateTable(c.Ctx, request)
	if err != nil {
		return c.RespondError(err)
	}

	return c.Respond(map[string]interface{}{
		"data":   response,
		"status": true,
	}, http.StatusOK)
}

func (uc Controller) CashierUpdateTableAll(c *web.Context) error {
	id := c.GetParam(reflect.Int, "id").(int)

	if err := c.ValidParam(); err != nil {
		return c.RespondError(err)
	}

	var request table_service.CashierUpdateRequest

	if err := c.BindFunc(&request, "Number", "Status", "Capacity"); err != nil {
		return c.RespondError(err)
	}

	request.ID = int64(id)

	err := uc.useCase.CashierUpdateTable(c.Ctx, request)
	if err != nil {
		return c.RespondError(err)
	}

	return c.Respond(map[string]interface{}{
		"data":   "ok!",
		"status": true,
	}, http.StatusOK)
}

func (uc Controller) CashierUpdateTableColumns(c *web.Context) error {
	id := c.GetParam(reflect.Int, "id").(int)

	if err := c.ValidParam(); err != nil {
		return c.RespondError(err)
	}

	var request table_service.CashierUpdateRequest

	if err := c.BindFunc(&request); err != nil {
		return c.RespondError(err)
	}

	request.ID = int64(id)

	err := uc.useCase.CashierUpdateTableColumn(c.Ctx, request)
	if err != nil {
		return c.RespondError(err)
	}

	return c.Respond(map[string]interface{}{
		"data":   "ok!",
		"status": true,
	}, http.StatusOK)
}

func (uc Controller) CashierDeleteTable(c *web.Context) error {
	id := c.GetParam(reflect.Int, "id").(int)

	if err := c.ValidParam(); err != nil {
		return c.RespondError(err)
	}

	err := uc.useCase.CashierDeleteTable(c.Ctx, int64(id))
	if err != nil {
		return c.RespondError(err)
	}

	return c.Respond(map[string]interface{}{
		"data":   "ok!",
		"status": true,
	}, http.StatusOK)
}

func (uc Controller) CashierGenerateQRTable(c *web.Context) error {
	var request table_service.CashierGenerateQRTable

	if err := c.BindFunc(&request, "Tables"); err != nil {
		return c.RespondError(err)
	}

	response, err := uc.useCase.CashierGenerateQRTable(c.Ctx, request)
	if err != nil {
		return c.RespondError(err)
	}

	return c.Respond(map[string]interface{}{
		"data":   response,
		"status": true,
	}, http.StatusOK)
}

// @waiter

func (uc Controller) WaiterGetTableList(c *web.Context) error {
	var filter table_service.Filter

	if limit, ok := c.GetQueryFunc(reflect.Int, "limit").(*int); ok {
		filter.Limit = limit
	}
	if offset, ok := c.GetQueryFunc(reflect.Int, "offset").(*int); ok {
		filter.Offset = offset
	}
	if search, ok := c.GetQueryFunc(reflect.String, "search").(*string); ok {
		filter.Search = search
	}

	if err := c.ValidQuery(); err != nil {
		return c.RespondError(err)
	}

	list, count, err := uc.useCase.WaiterGetTableList(c.Ctx, filter)
	if err != nil {
		return c.RespondError(err)
	}

	return c.Respond(map[string]interface{}{
		"data": map[string]interface{}{
			"results": list,
			"count":   count,
		},
		"status": true,
	}, http.StatusOK)
}

// #branchReview ----------------------------------------------------------------------------------------

// @client

func (uc Controller) ClientGetBranchReviewList(c *web.Context) error {
	var filter branchReview_service.Filter

	if limit, ok := c.GetQueryFunc(reflect.Int, "limit").(*int); ok {
		filter.Limit = limit
	}
	if offset, ok := c.GetQueryFunc(reflect.Int, "offset").(*int); ok {
		filter.Offset = offset
	}

	if branchID, ok := c.GetQueryFunc(reflect.Int64, "branch_id").(*int64); ok {
		filter.BranchID = branchID
	}
	if userID, ok := c.GetQueryFunc(reflect.Int64, "user_id").(*int64); ok {
		filter.UserID = userID
	}

	if err := c.ValidQuery(); err != nil {
		return c.RespondError(err)
	}

	list, count, err := uc.useCase.ClientGetBranchReviewList(c.Ctx, filter)
	if err != nil {
		return c.RespondError(err)
	}

	return c.Respond(map[string]interface{}{
		"data": map[string]interface{}{
			"results": list,
			"count":   count,
		},
		"status": true,
	}, http.StatusOK)
}

func (uc Controller) ClientGetBranchReviewDetail(c *web.Context) error {
	id := c.GetParam(reflect.Int, "id").(int)

	if err := c.ValidParam(); err != nil {
		return c.RespondError(err)
	}

	response, err := uc.useCase.ClientGetBranchReviewDetail(c.Ctx, int64(id))
	if err != nil {
		return c.RespondError(err)
	}

	return c.Respond(map[string]interface{}{
		"data":   response,
		"status": true,
	}, http.StatusOK)
}

func (uc Controller) ClientCreateBranchReview(c *web.Context) error {
	var request branchReview_service.ClientCreateRequest

	if err := c.BindFunc(&request, "Point", "Rate", "Comment", "BranchID"); err != nil {
		return c.RespondError(err)
	}

	response, err := uc.useCase.ClientCreateBranchReview(c.Ctx, request)
	if err != nil {
		return c.RespondError(err)
	}

	return c.Respond(map[string]interface{}{
		"data":   response,
		"status": true,
	}, http.StatusOK)
}

func (uc Controller) ClientUpdateBranchReviewAll(c *web.Context) error {
	id := c.GetParam(reflect.Int, "id").(int)

	if err := c.ValidParam(); err != nil {
		return c.RespondError(err)
	}

	var request branchReview_service.ClientUpdateRequest

	if err := c.BindFunc(&request, "ID", "Point", "Rate", "Comment", "UserID", "BranchID"); err != nil {
		return c.RespondError(err)
	}

	request.ID = int64(id)

	err := uc.useCase.ClientUpdateBranchReview(c.Ctx, request)
	if err != nil {
		return c.RespondError(err)
	}

	return c.Respond(map[string]interface{}{
		"data":   "ok!",
		"status": true,
	}, http.StatusOK)
}

func (uc Controller) ClientUpdateBranchReviewColumns(c *web.Context) error {
	id := c.GetParam(reflect.Int, "id").(int)

	if err := c.ValidParam(); err != nil {
		return c.RespondError(err)
	}

	var request branchReview_service.ClientUpdateRequest

	if err := c.BindFunc(&request, "ID"); err != nil {
		return c.RespondError(err)
	}

	request.ID = int64(id)

	err := uc.useCase.ClientUpdateBranchReviewColumn(c.Ctx, request)
	if err != nil {
		return c.RespondError(err)
	}

	return c.Respond(map[string]interface{}{
		"data":   "ok!",
		"status": true,
	}, http.StatusOK)
}

func (uc Controller) ClientDeleteBranchReview(c *web.Context) error {
	id := c.GetParam(reflect.Int, "id").(int)

	if err := c.ValidParam(); err != nil {
		return c.RespondError(err)
	}

	err := uc.useCase.ClientDeleteBranchReview(c.Ctx, int64(id))
	if err != nil {
		return c.RespondError(err)
	}

	return c.Respond(map[string]interface{}{
		"data":   "ok!",
		"status": true,
	}, http.StatusOK)
}

// #devices ----------------------------------------------------------------------------------------

// @branch

func (uc Controller) BranchGetPrintersList(c *web.Context) error {
	var filter printers_service.Filter

	if limit, ok := c.GetQueryFunc(reflect.Int, "limit").(*int); ok {
		filter.Limit = limit
	}
	if offset, ok := c.GetQueryFunc(reflect.Int, "offset").(*int); ok {
		filter.Offset = offset
	}

	if warehouseID, ok := c.GetQueryFunc(reflect.Int, "warehouse_id").(*int); ok {
		filter.WarehouseID = warehouseID
	}

	if err := c.ValidQuery(); err != nil {
		return c.RespondError(err)
	}

	list, count, err := uc.useCase.BranchGetPrintersList(c.Ctx, filter)
	if err != nil {
		return c.RespondError(err)
	}

	return c.Respond(map[string]interface{}{
		"data": map[string]interface{}{
			"results": list,
			"count":   count,
		},
		"status": true,
	}, http.StatusOK)
}

func (uc Controller) BranchGetPrintersDetail(c *web.Context) error {
	id := c.GetParam(reflect.Int, "id").(int)

	if err := c.ValidParam(); err != nil {
		return c.RespondError(err)
	}

	response, err := uc.useCase.BranchGetPrintersDetail(c.Ctx, int64(id))
	if err != nil {
		return c.RespondError(err)
	}

	return c.Respond(map[string]interface{}{
		"data":   response,
		"status": true,
	}, http.StatusOK)
}

func (uc Controller) BranchCreatePrinters(c *web.Context) error {
	var request printers_service.BranchCreateRequest

	if err := c.BindFunc(&request, "Name", "IP", "Port"); err != nil {
		return c.RespondError(err)
	}

	response, err := uc.useCase.BranchCreatePrinters(c.Ctx, request)
	if err != nil {
		return c.RespondError(err)
	}

	return c.Respond(map[string]interface{}{
		"data":   response,
		"status": true,
	}, http.StatusOK)
}

func (uc Controller) BranchUpdatePrintersAll(c *web.Context) error {
	id := c.GetParam(reflect.Int, "id").(int)

	if err := c.ValidParam(); err != nil {
		return c.RespondError(err)
	}

	var request printers_service.BranchUpdateRequest

	if err := c.BindFunc(&request, "ID", "Name", "IP", "Port"); err != nil {
		return c.RespondError(err)
	}

	request.ID = int64(id)

	err := uc.useCase.BranchUpdatePrinters(c.Ctx, request)
	if err != nil {
		return c.RespondError(err)
	}

	return c.Respond(map[string]interface{}{
		"data":   "ok!",
		"status": true,
	}, http.StatusOK)
}

func (uc Controller) BranchUpdatePrintersColumns(c *web.Context) error {
	id := c.GetParam(reflect.Int, "id").(int)

	if err := c.ValidParam(); err != nil {
		return c.RespondError(err)
	}

	var request printers_service.BranchUpdateRequest

	if err := c.BindFunc(&request); err != nil {
		return c.RespondError(err)
	}

	request.ID = int64(id)

	err := uc.useCase.BranchUpdatePrintersColumn(c.Ctx, request)
	if err != nil {
		return c.RespondError(err)
	}

	return c.Respond(map[string]interface{}{
		"data":   "ok!",
		"status": true,
	}, http.StatusOK)
}

func (uc Controller) BranchDeletePrinters(c *web.Context) error {
	id := c.GetParam(reflect.Int, "id").(int)

	if err := c.ValidParam(); err != nil {
		return c.RespondError(err)
	}

	err := uc.useCase.BranchDeletePrinters(c.Ctx, int64(id))
	if err != nil {
		return c.RespondError(err)
	}

	return c.Respond(map[string]interface{}{
		"data":   "ok!",
		"status": true,
	}, http.StatusOK)
}
