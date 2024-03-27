package food

import (
	"github.com/restaurant/foundation/web"
	"github.com/restaurant/internal/service/food_recipe"
	"net/http"
	"reflect"
)

type Controller struct {
	useCase *food.UseCase
}

func NewController(useCase *food.UseCase) *Controller {
	return &Controller{useCase}
}

// #food_category ---------------------------------------------------------------

// @admin

func (uc Controller) AdminGetFoodCategoryList(c *web.Context) error {
	var filter foodCategory_service.Filter

	if limit, ok := c.GetQueryFunc(reflect.Int, "limit").(*int); ok {
		filter.Limit = limit
	}
	if offset, ok := c.GetQueryFunc(reflect.Int, "offset").(*int); ok {
		filter.Offset = offset
	}
	if page, ok := c.GetQueryFunc(reflect.Int, "page").(*int); ok {
		filter.Page = page
	}
	if name, ok := c.GetQueryFunc(reflect.String, "name").(*string); ok {
		filter.Name = name
	}

	if err := c.ValidQuery(); err != nil {
		return c.RespondError(err)
	}

	list, count, err := uc.useCase.AdminGetFoodCategoryList(c.Ctx, filter)
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

func (uc Controller) AdminGetFoodCategoryDetail(c *web.Context) error {
	id := c.GetParam(reflect.Int, "id").(int)

	if err := c.ValidParam(); err != nil {
		return c.RespondError(err)
	}

	response, err := uc.useCase.AdminGetFoodCategoryDetail(c.Ctx, int64(id))
	if err != nil {
		return c.RespondError(err)
	}

	return c.Respond(map[string]interface{}{
		"data":   response,
		"status": true,
	}, http.StatusOK)
}

func (uc Controller) AdminCreateFoodCategory(c *web.Context) error {
	var request foodCategory_service.AdminCreateRequest

	if err := c.BindFunc(&request, "Name", "Logo"); err != nil {
		return c.RespondError(err)
	}

	response, err := uc.useCase.AdminCreateFoodCategory(c.Ctx, request)
	if err != nil {
		return c.RespondError(err)
	}

	return c.Respond(map[string]interface{}{
		"data":   response,
		"status": true,
	}, http.StatusOK)
}

func (uc Controller) AdminUpdateFoodCategoryAll(c *web.Context) error {
	id := c.GetParam(reflect.Int, "id").(int)

	if err := c.ValidParam(); err != nil {
		return c.RespondError(err)
	}

	var request foodCategory_service.AdminUpdateRequest

	if err := c.BindFunc(&request, "Name", "Logo", "Main"); err != nil {
		return c.RespondError(err)
	}

	request.ID = int64(id)

	err := uc.useCase.AdminUpdateFoodCategory(c.Ctx, request)
	if err != nil {
		return c.RespondError(err)
	}

	return c.Respond(map[string]interface{}{
		"data":   "ok!",
		"status": true,
	}, http.StatusOK)
}

func (uc Controller) AdminUpdateFoodCategoryColumns(c *web.Context) error {
	id := c.GetParam(reflect.Int, "id").(int)

	if err := c.ValidParam(); err != nil {
		return c.RespondError(err)
	}

	var request foodCategory_service.AdminUpdateRequest

	if err := c.BindFunc(&request); err != nil {
		return c.RespondError(err)
	}

	request.ID = int64(id)

	err := uc.useCase.AdminUpdateFoodCategoryColumn(c.Ctx, request)
	if err != nil {
		return c.RespondError(err)
	}

	return c.Respond(map[string]interface{}{
		"data":   "ok!",
		"status": true,
	}, http.StatusOK)
}

func (uc Controller) AdminDeleteFoodCategory(c *web.Context) error {
	id := c.GetParam(reflect.Int, "id").(int)

	if err := c.ValidParam(); err != nil {
		return c.RespondError(err)
	}

	err := uc.useCase.AdminDeleteFoodCategory(c.Ctx, int64(id))
	if err != nil {
		return c.RespondError(err)
	}

	return c.Respond(map[string]interface{}{
		"data":   "ok!",
		"status": true,
	}, http.StatusOK)
}

// @branch

func (uc Controller) BranchGetFoodCategoryList(c *web.Context) error {
	var filter foodCategory_service.Filter

	if limit, ok := c.GetQueryFunc(reflect.Int, "limit").(*int); ok {
		filter.Limit = limit
	}
	if offset, ok := c.GetQueryFunc(reflect.Int, "offset").(*int); ok {
		filter.Offset = offset
	}
	if page, ok := c.GetQueryFunc(reflect.Int, "page").(*int); ok {
		filter.Page = page
	}
	if name, ok := c.GetQueryFunc(reflect.String, "name").(*string); ok {
		filter.Name = name
	}

	if err := c.ValidQuery(); err != nil {
		return c.RespondError(err)
	}

	list, count, err := uc.useCase.BranchGetFoodCategoryList(c.Ctx, filter)
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

func (uc Controller) BranchGetFoodCategoryDetail(c *web.Context) error {
	id := c.GetParam(reflect.Int, "id").(int)

	if err := c.ValidParam(); err != nil {
		return c.RespondError(err)
	}

	response, err := uc.useCase.BranchGetFoodCategoryDetail(c.Ctx, int64(id))
	if err != nil {
		return c.RespondError(err)
	}

	return c.Respond(map[string]interface{}{
		"data":   response,
		"status": true,
	}, http.StatusOK)
}

func (uc Controller) BranchCreateFoodCategory(c *web.Context) error {
	var request foodCategory_service.BranchCreateRequest

	if err := c.BindFunc(&request, "Name", "Logo"); err != nil {
		return c.RespondError(err)
	}

	response, err := uc.useCase.BranchCreateFoodCategory(c.Ctx, request)
	if err != nil {
		return c.RespondError(err)
	}

	return c.Respond(map[string]interface{}{
		"data":   response,
		"status": true,
	}, http.StatusOK)
}

func (uc Controller) BranchUpdateFoodCategoryAll(c *web.Context) error {
	id := c.GetParam(reflect.Int, "id").(int)

	if err := c.ValidParam(); err != nil {
		return c.RespondError(err)
	}

	var request foodCategory_service.BranchUpdateRequest

	if err := c.BindFunc(&request, "Name", "Logo", "Main"); err != nil {
		return c.RespondError(err)
	}

	request.ID = int64(id)

	err := uc.useCase.BranchUpdateFoodCategory(c.Ctx, request)
	if err != nil {
		return c.RespondError(err)
	}

	return c.Respond(map[string]interface{}{
		"data":   "ok!",
		"status": true,
	}, http.StatusOK)
}

func (uc Controller) BranchUpdateFoodCategoryColumns(c *web.Context) error {
	id := c.GetParam(reflect.Int, "id").(int)

	if err := c.ValidParam(); err != nil {
		return c.RespondError(err)
	}

	var request foodCategory_service.BranchUpdateRequest

	if err := c.BindFunc(&request); err != nil {
		return c.RespondError(err)
	}

	request.ID = int64(id)

	err := uc.useCase.BranchUpdateFoodCategoryColumn(c.Ctx, request)
	if err != nil {
		return c.RespondError(err)
	}

	return c.Respond(map[string]interface{}{
		"data":   "ok!",
		"status": true,
	}, http.StatusOK)
}

func (uc Controller) BranchDeleteFoodCategory(c *web.Context) error {
	id := c.GetParam(reflect.Int, "id").(int)

	if err := c.ValidParam(); err != nil {
		return c.RespondError(err)
	}

	err := uc.useCase.BranchDeleteFoodCategory(c.Ctx, int64(id))
	if err != nil {
		return c.RespondError(err)
	}

	return c.Respond(map[string]interface{}{
		"data":   "ok!",
		"status": true,
	}, http.StatusOK)
}

// @cashier

func (uc Controller) CashierGetFoodCategoryList(c *web.Context) error {
	var filter foodCategory_service.Filter

	if limit, ok := c.GetQueryFunc(reflect.Int, "limit").(*int); ok {
		filter.Limit = limit
	}
	if offset, ok := c.GetQueryFunc(reflect.Int, "offset").(*int); ok {
		filter.Offset = offset
	}
	if page, ok := c.GetQueryFunc(reflect.Int, "page").(*int); ok {
		filter.Page = page
	}
	if name, ok := c.GetQueryFunc(reflect.String, "name").(*string); ok {
		filter.Name = name
	}

	if err := c.ValidQuery(); err != nil {
		return c.RespondError(err)
	}

	list, count, err := uc.useCase.CashierGetFoodCategoryList(c.Ctx, filter)
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

func (uc Controller) CashierGetFoodCategoryDetail(c *web.Context) error {
	id := c.GetParam(reflect.Int, "id").(int)

	if err := c.ValidParam(); err != nil {
		return c.RespondError(err)
	}

	response, err := uc.useCase.CashierGetFoodCategoryDetail(c.Ctx, int64(id))
	if err != nil {
		return c.RespondError(err)
	}

	return c.Respond(map[string]interface{}{
		"data":   response,
		"status": true,
	}, http.StatusOK)
}

func (uc Controller) CashierCreateFoodCategory(c *web.Context) error {
	var request foodCategory_service.CashierCreateRequest

	if err := c.BindFunc(&request, "Name", "Logo"); err != nil {
		return c.RespondError(err)
	}

	response, err := uc.useCase.CashierCreateFoodCategory(c.Ctx, request)
	if err != nil {
		return c.RespondError(err)
	}

	return c.Respond(map[string]interface{}{
		"data":   response,
		"status": true,
	}, http.StatusOK)
}

func (uc Controller) CashierUpdateFoodCategoryAll(c *web.Context) error {
	id := c.GetParam(reflect.Int, "id").(int)

	if err := c.ValidParam(); err != nil {
		return c.RespondError(err)
	}

	var request foodCategory_service.CashierUpdateRequest

	if err := c.BindFunc(&request, "Name", "Logo", "Main"); err != nil {
		return c.RespondError(err)
	}

	request.ID = int64(id)

	err := uc.useCase.CashierUpdateFoodCategory(c.Ctx, request)
	if err != nil {
		return c.RespondError(err)
	}

	return c.Respond(map[string]interface{}{
		"data":   "ok!",
		"status": true,
	}, http.StatusOK)
}

func (uc Controller) CashierUpdateFoodCategoryColumns(c *web.Context) error {
	id := c.GetParam(reflect.Int, "id").(int)

	if err := c.ValidParam(); err != nil {
		return c.RespondError(err)
	}

	var request foodCategory_service.CashierUpdateRequest

	if err := c.BindFunc(&request); err != nil {
		return c.RespondError(err)
	}

	request.ID = int64(id)

	err := uc.useCase.CashierUpdateFoodCategoryColumn(c.Ctx, request)
	if err != nil {
		return c.RespondError(err)
	}

	return c.Respond(map[string]interface{}{
		"data":   "ok!",
		"status": true,
	}, http.StatusOK)
}

func (uc Controller) CashierDeleteFoodCategory(c *web.Context) error {
	id := c.GetParam(reflect.Int, "id").(int)

	if err := c.ValidParam(); err != nil {
		return c.RespondError(err)
	}

	err := uc.useCase.CashierDeleteFoodCategory(c.Ctx, int64(id))
	if err != nil {
		return c.RespondError(err)
	}

	return c.Respond(map[string]interface{}{
		"data":   "ok!",
		"status": true,
	}, http.StatusOK)
}

// @client

func (uc Controller) ClientGetFoodCategoryList(c *web.Context) error {
	var filter foodCategory_service.Filter

	if limit, ok := c.GetQueryFunc(reflect.Int, "limit").(*int); ok {
		filter.Limit = limit
	}
	if offset, ok := c.GetQueryFunc(reflect.Int, "offset").(*int); ok {
		filter.Offset = offset
	}
	if page, ok := c.GetQueryFunc(reflect.Int, "page").(*int); ok {
		filter.Page = page
	}
	if name, ok := c.GetQueryFunc(reflect.String, "name").(*string); ok {
		filter.Name = name
	}

	if err := c.ValidQuery(); err != nil {
		return c.RespondMobileError(err)
	}

	list, count, err := uc.useCase.ClientGetFoodCategoryList(c.Ctx, filter)
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

// @waiter

func (uc Controller) WaiterGetFoodCategoryList(c *web.Context) error {
	list, err := uc.useCase.WaiterGetFoodCategoryList(c.Ctx)
	if err != nil {
		return c.RespondError(err)
	}

	return c.Respond(map[string]interface{}{
		"data":   list,
		"status": true,
	}, http.StatusOK)
}

// #food ---------------------------------------------------------------

// @admin

func (uc Controller) AdminGetFoodList(c *web.Context) error {
	var filter food_service.Filter

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

	list, count, err := uc.useCase.AdminGetFoodList(c.Ctx, filter)
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

func (uc Controller) AdminGetFoodDetail(c *web.Context) error {
	id := c.GetParam(reflect.Int, "id").(int)

	if err := c.ValidParam(); err != nil {
		return c.RespondError(err)
	}

	response, err := uc.useCase.AdminGetFoodDetail(c.Ctx, int64(id))
	if err != nil {
		return c.RespondError(err)
	}

	return c.Respond(map[string]interface{}{
		"data":   response,
		"status": true,
	}, http.StatusOK)
}

func (uc Controller) AdminCreateFood(c *web.Context) error {
	var request food_service.AdminCreateRequest

	if err := c.BindFunc(&request, "Name", "CategoryID"); err != nil {
		return c.RespondError(err)
	}

	response, err := uc.useCase.AdminCreateFood(c.Ctx, request)
	if err != nil {
		return c.RespondError(err)
	}

	return c.Respond(map[string]interface{}{
		"data":   response,
		"status": true,
	}, http.StatusOK)
}

func (uc Controller) AdminUpdateFoodAll(c *web.Context) error {
	id := c.GetParam(reflect.Int, "id").(int)

	if err := c.ValidParam(); err != nil {
		return c.RespondError(err)
	}

	var request food_service.AdminUpdateRequest

	if err := c.BindFunc(&request, "ID", "Name", "Photos", "CategoryID", "Price"); err != nil {
		return c.RespondError(err)
	}

	request.ID = int64(id)

	err := uc.useCase.AdminUpdateFood(c.Ctx, request)
	if err != nil {
		return c.RespondError(err)
	}

	return c.Respond(map[string]interface{}{
		"data":   "ok!",
		"status": true,
	}, http.StatusOK)
}

func (uc Controller) AdminUpdateFoodColumns(c *web.Context) error {
	id := c.GetParam(reflect.Int, "id").(int)

	if err := c.ValidParam(); err != nil {
		return c.RespondError(err)
	}

	var request food_service.AdminUpdateRequest

	if err := c.BindFunc(&request); err != nil {
		return c.RespondError(err)
	}

	request.ID = int64(id)

	err := uc.useCase.AdminUpdateFoodColumn(c.Ctx, request)
	if err != nil {
		return c.RespondError(err)
	}

	return c.Respond(map[string]interface{}{
		"data":   "ok!",
		"status": true,
	}, http.StatusOK)
}

func (uc Controller) AdminDeleteFood(c *web.Context) error {
	id := c.GetParam(reflect.Int, "id").(int)

	if err := c.ValidParam(); err != nil {
		return c.RespondError(err)
	}

	err := uc.useCase.AdminDeleteFood(c.Ctx, int64(id))
	if err != nil {
		return c.RespondError(err)
	}

	return c.Respond(map[string]interface{}{
		"data":   "ok!",
		"status": true,
	}, http.StatusOK)
}

func (uc Controller) AdminImageDelete(c *web.Context) error {
	id := c.GetParam(reflect.Int, "id").(int)

	if err := c.ValidParam(); err != nil {
		return c.RespondError(err)
	}

	var request food_service.AdminDeleteImageRequest

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

func (uc Controller) ClientGetMenuList(c *web.Context) error {
	var filter menu_service.Filter
	if search, ok := c.GetQueryFunc(reflect.String, "search").(*string); ok {
		filter.Search = search
	}
	if branchID, ok := c.GetQueryFunc(reflect.Int, "branch-id").(*int); ok {
		filter.BranchID = branchID
	}

	if err := c.ValidQuery(); err != nil {
		return c.RespondMobileError(err)
	}

	list, err := uc.useCase.ClientGetMenuList(c.Ctx, filter)
	if err != nil {
		return c.RespondMobileError(err)
	}

	return c.Respond(map[string]interface{}{
		"error": nil,
		"data": map[string]interface{}{
			"results": list,
		},
		"status": true,
	}, http.StatusOK)
}

func (uc Controller) ClientGetMenuListByCategoryID(c *web.Context) error {
	id := c.GetParam(reflect.Int, "category_id").(int)

	if err := c.ValidParam(); err != nil {
		return c.RespondError(err)
	}

	var filter menu_service.Filter

	if lat, ok := c.GetQueryFunc(reflect.Float64, "lat").(*float64); ok {
		filter.Lat = lat
	}
	if lon, ok := c.GetQueryFunc(reflect.Float64, "lon").(*float64); ok {
		filter.Lon = lon
	}

	if err := c.ValidQuery(); err != nil {
		return c.RespondMobileError(err)
	}

	list, err := uc.useCase.ClientGetMenuListByCategoryID(c.Ctx, id, filter)
	if err != nil {
		return c.RespondMobileError(err)
	}

	return c.Respond(map[string]interface{}{
		"error":  nil,
		"data":   list,
		"status": true,
	}, http.StatusOK)
}

// @branch

func (uc Controller) BranchGetFoodList(c *web.Context) error {
	var filter food_service.Filter

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

	list, count, err := uc.useCase.BranchGetFoodList(c.Ctx, filter)
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

func (uc Controller) BranchGetFoodDetail(c *web.Context) error {
	id := c.GetParam(reflect.Int, "id").(int)

	if err := c.ValidParam(); err != nil {
		return c.RespondError(err)
	}

	response, err := uc.useCase.BranchGetFoodDetail(c.Ctx, int64(id))
	if err != nil {
		return c.RespondError(err)
	}

	return c.Respond(map[string]interface{}{
		"data":   response,
		"status": true,
	}, http.StatusOK)
}

func (uc Controller) BranchCreateFood(c *web.Context) error {
	var request food_service.BranchCreateRequest

	if err := c.BindFunc(&request, "Name", "CategoryID"); err != nil {
		return c.RespondError(err)
	}

	response, err := uc.useCase.BranchCreateFood(c.Ctx, request)
	if err != nil {
		return c.RespondError(err)
	}

	return c.Respond(map[string]interface{}{
		"data":   response,
		"status": true,
	}, http.StatusOK)
}

func (uc Controller) BranchUpdateFoodAll(c *web.Context) error {
	id := c.GetParam(reflect.Int, "id").(int)

	if err := c.ValidParam(); err != nil {
		return c.RespondError(err)
	}

	var request food_service.BranchUpdateRequest

	if err := c.BindFunc(&request, "ID", "Name", "Photos", "Price", "CategoryID"); err != nil {
		return c.RespondError(err)
	}

	id64 := int64(id)
	request.ID = &id64

	err := uc.useCase.BranchUpdateFood(c.Ctx, request)
	if err != nil {
		return c.RespondError(err)
	}

	return c.Respond(map[string]interface{}{
		"data":   "ok!",
		"status": true,
	}, http.StatusOK)
}

func (uc Controller) BranchUpdateFoodColumns(c *web.Context) error {
	id := c.GetParam(reflect.Int, "id").(int)

	if err := c.ValidParam(); err != nil {
		return c.RespondError(err)
	}

	var request food_service.BranchUpdateRequest

	if err := c.BindFunc(&request); err != nil {
		return c.RespondError(err)
	}

	id64 := int64(id)
	request.ID = &id64

	err := uc.useCase.BranchUpdateFoodColumn(c.Ctx, request)
	if err != nil {
		return c.RespondError(err)
	}

	return c.Respond(map[string]interface{}{
		"data":   "ok!",
		"status": true,
	}, http.StatusOK)
}

func (uc Controller) BranchDeleteFood(c *web.Context) error {
	id := c.GetParam(reflect.Int, "id").(int)

	if err := c.ValidParam(); err != nil {
		return c.RespondError(err)
	}

	err := uc.useCase.BranchDeleteFood(c.Ctx, int64(id))
	if err != nil {
		return c.RespondError(err)
	}

	return c.Respond(map[string]interface{}{
		"data":   "ok!",
		"status": true,
	}, http.StatusOK)
}

func (uc Controller) BranchImageDelete(c *web.Context) error {
	id := c.GetParam(reflect.Int, "id").(int)

	if err := c.ValidParam(); err != nil {
		return c.RespondError(err)
	}

	var request food_service.AdminDeleteImageRequest

	if err := c.BindFunc(&request); err != nil {
		return c.RespondError(err)
	}

	request.ID = int64(id)

	err := uc.useCase.BranchDeleteImage(c.Ctx, request)
	if err != nil {
		return c.RespondError(err)
	}

	return c.Respond(map[string]interface{}{
		"data":   "ok!",
		"status": true,
	}, http.StatusOK)
}

// @cashier

func (uc Controller) CashierGetFoodList(c *web.Context) error {
	var filter food_service.Filter

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

	list, count, err := uc.useCase.CashierGetFoodList(c.Ctx, filter)
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

// #menu ---------------------------------------------------------------

// @admin

func (uc Controller) AdminGetMenuList(c *web.Context) error {
	var filter menu_service.Filter

	if limit, ok := c.GetQueryFunc(reflect.Int, "limit").(*int); ok {
		filter.Limit = limit
	}
	if offset, ok := c.GetQueryFunc(reflect.Int, "offset").(*int); ok {
		filter.Offset = offset
	}
	if branchID, ok := c.GetQueryFunc(reflect.Int, "branch_id").(*int); ok {
		filter.BranchID = branchID
	}

	if err := c.ValidQuery(); err != nil {
		return c.RespondError(err)
	}

	list, count, err := uc.useCase.AdminGetMenuList(c.Ctx, filter)
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

func (uc Controller) AdminGetMenuDetail(c *web.Context) error {
	id := c.GetParam(reflect.Int, "id").(int)

	if err := c.ValidParam(); err != nil {
		return c.RespondError(err)
	}

	response, err := uc.useCase.AdminGetMenuDetail(c.Ctx, int64(id))
	if err != nil {
		return c.RespondError(err)
	}

	return c.Respond(map[string]interface{}{
		"data":   response,
		"status": true,
	}, http.StatusOK)
}

func (uc Controller) AdminCreateMenu(c *web.Context) error {
	var request menu_service.AdminCreateRequest

	if err := c.BindFunc(&request, "MenuID", "BranchID"); err != nil {
		return c.RespondError(err)
	}

	response, err := uc.useCase.AdminCreateMenu(c.Ctx, request)
	if err != nil {
		return c.RespondError(err)
	}

	return c.Respond(map[string]interface{}{
		"data":   response,
		"status": true,
	}, http.StatusOK)
}

func (uc Controller) AdminUpdateMenuAll(c *web.Context) error {
	id := c.GetParam(reflect.Int, "id").(int)

	if err := c.ValidParam(); err != nil {
		return c.RespondError(err)
	}

	var request menu_service.AdminUpdateRequest

	if err := c.BindFunc(&request, "MenuID", "BranchID", "Status"); err != nil {
		return c.RespondError(err)
	}

	request.ID = int64(id)

	err := uc.useCase.AdminUpdateMenu(c.Ctx, request)
	if err != nil {
		return c.RespondError(err)
	}

	return c.Respond(map[string]interface{}{
		"data":   "ok!",
		"status": true,
	}, http.StatusOK)
}

func (uc Controller) AdminUpdateMenuColumns(c *web.Context) error {
	id := c.GetParam(reflect.Int, "id").(int)

	if err := c.ValidParam(); err != nil {
		return c.RespondError(err)
	}

	var request menu_service.AdminUpdateRequest

	if err := c.BindFunc(&request); err != nil {
		return c.RespondError(err)
	}

	request.ID = int64(id)

	err := uc.useCase.AdminUpdateMenuColumn(c.Ctx, request)
	if err != nil {
		return c.RespondError(err)
	}

	return c.Respond(map[string]interface{}{
		"data":   "ok!",
		"status": true,
	}, http.StatusOK)
}

func (uc Controller) AdminDeleteMenu(c *web.Context) error {
	id := c.GetParam(reflect.Int, "id").(int)

	if err := c.ValidParam(); err != nil {
		return c.RespondError(err)
	}

	err := uc.useCase.AdminDeleteMenu(c.Ctx, int64(id))
	if err != nil {
		return c.RespondError(err)
	}

	return c.Respond(map[string]interface{}{
		"data":   "ok!",
		"status": true,
	}, http.StatusOK)
}

func (uc Controller) AdminRemoveMenuPhoto(c *web.Context) error {
	id := c.GetParam(reflect.Int, "id").(int)

	if err := c.ValidParam(); err != nil {
		return c.RespondError(err)
	}

	var imageIndex *int
	if index, ok := c.GetQueryFunc(reflect.Int, "image_index").(*int); ok {
		imageIndex = index
	}

	err := uc.useCase.AdminRemoveMenuPhoto(c.Ctx, int64(id), imageIndex)
	if err != nil {
		return c.RespondError(err)
	}

	return c.Respond(map[string]interface{}{
		"data":   "ok!",
		"status": true,
	}, http.StatusOK)
}

// @branch

func (uc Controller) BranchGetMenuList(c *web.Context) error {
	var filter menu_service.Filter

	if limit, ok := c.GetQueryFunc(reflect.Int, "limit").(*int); ok {
		filter.Limit = limit
	}
	if offset, ok := c.GetQueryFunc(reflect.Int, "offset").(*int); ok {
		filter.Offset = offset
	}
	if printerId, ok := c.GetQueryFunc(reflect.Int, "printer_id").(*int); ok {
		filter.PrinterID = printerId
	}
	if printer, ok := c.GetQueryFunc(reflect.Bool, "printer").(*bool); ok {
		filter.Printer = printer
	}

	if err := c.ValidQuery(); err != nil {
		return c.RespondError(err)
	}

	list, count, err := uc.useCase.BranchGetMenuList(c.Ctx, filter)
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

func (uc Controller) BranchGetMenuDetail(c *web.Context) error {
	id := c.GetParam(reflect.Int, "id").(int)

	if err := c.ValidParam(); err != nil {
		return c.RespondError(err)
	}

	response, err := uc.useCase.BranchGetMenuDetail(c.Ctx, int64(id))
	if err != nil {
		return c.RespondError(err)
	}

	return c.Respond(map[string]interface{}{
		"data":   response,
		"status": true,
	}, http.StatusOK)
}

func (uc Controller) BranchCreateMenu(c *web.Context) error {
	var request menu_service.BranchCreateRequest

	if err := c.BindFunc(&request, "MenuID", "BranchID"); err != nil {
		return c.RespondError(err)
	}

	response, err := uc.useCase.BranchCreateMenu(c.Ctx, request)
	if err != nil {
		return c.RespondError(err)
	}

	return c.Respond(map[string]interface{}{
		"data":   response,
		"status": true,
	}, http.StatusOK)
}

func (uc Controller) BranchUpdateMenuAll(c *web.Context) error {
	id := c.GetParam(reflect.Int, "id").(int)

	if err := c.ValidParam(); err != nil {
		return c.RespondError(err)
	}

	var request menu_service.BranchUpdateRequest

	if err := c.BindFunc(&request, "MenuID", "BranchID", "Status"); err != nil {
		return c.RespondError(err)
	}

	request.ID = int64(id)

	err := uc.useCase.BranchUpdateMenu(c.Ctx, request)
	if err != nil {
		return c.RespondError(err)
	}

	return c.Respond(map[string]interface{}{
		"data":   "ok!",
		"status": true,
	}, http.StatusOK)
}

func (uc Controller) BranchUpdateMenuColumns(c *web.Context) error {
	id := c.GetParam(reflect.Int, "id").(int)

	if err := c.ValidParam(); err != nil {
		return c.RespondError(err)
	}

	var request menu_service.BranchUpdateRequest

	if err := c.BindFunc(&request); err != nil {
		return c.RespondError(err)
	}

	request.ID = int64(id)

	err := uc.useCase.BranchUpdateMenuColumn(c.Ctx, request)
	if err != nil {
		return c.RespondError(err)
	}

	return c.Respond(map[string]interface{}{
		"data":   "ok!",
		"status": true,
	}, http.StatusOK)
}

func (uc Controller) BranchDeleteMenu(c *web.Context) error {
	id := c.GetParam(reflect.Int, "id").(int)

	if err := c.ValidParam(); err != nil {
		return c.RespondError(err)
	}

	err := uc.useCase.BranchDeleteMenu(c.Ctx, int64(id))
	if err != nil {
		return c.RespondError(err)
	}

	return c.Respond(map[string]interface{}{
		"data":   "ok!",
		"status": true,
	}, http.StatusOK)
}

func (uc Controller) BranchUpdateMenuPrinterID(c *web.Context) error {
	var request menu_service.BranchUpdatePrinterIDRequest

	if err := c.BindFunc(&request, "PrinterID", "MenuIds"); err != nil {
		return c.RespondError(err)
	}

	err := uc.useCase.BranchUpdatePrinterID(c.Ctx, request)
	if err != nil {
		return c.RespondError(err)
	}

	return c.Respond(map[string]interface{}{
		"data":   "ok!",
		"status": true,
	}, http.StatusOK)
}

func (uc Controller) BranchDeleteMenuPrinterID(c *web.Context) error {
	id := c.GetParam(reflect.Int, "id").(int)

	if err := c.ValidParam(); err != nil {
		return c.RespondError(err)
	}

	err := uc.useCase.BranchDeletePrinterID(c.Ctx, int64(id))
	if err != nil {
		return c.RespondError(err)
	}

	return c.Respond(map[string]interface{}{
		"data":   "ok!",
		"status": true,
	}, http.StatusOK)
}

func (uc Controller) BranchRemoveMenuPhoto(c *web.Context) error {
	id := c.GetParam(reflect.Int, "id").(int)

	if err := c.ValidParam(); err != nil {
		return c.RespondError(err)
	}

	var request struct {
		Index int `json:"index"`
	}

	if err := c.BindFunc(&request, "Index"); err != nil {
		return c.RespondError(err)
	}

	err := uc.useCase.BranchRemoveMenuPhoto(c.Ctx, int64(id), request.Index)
	if err != nil {
		return c.RespondError(err)
	}

	return c.Respond(map[string]interface{}{
		"data":   "ok!",
		"status": true,
	}, http.StatusOK)
}

// @cashier

func (uc Controller) CashierGetMenuList(c *web.Context) error {
	var filter menu_service.Filter

	if limit, ok := c.GetQueryFunc(reflect.Int, "limit").(*int); ok {
		filter.Limit = limit
	}
	if offset, ok := c.GetQueryFunc(reflect.Int, "offset").(*int); ok {
		filter.Offset = offset
	}
	if printerId, ok := c.GetQueryFunc(reflect.Int, "printer_id").(*int); ok {
		filter.PrinterID = printerId
	}
	if printer, ok := c.GetQueryFunc(reflect.Bool, "printer").(*bool); ok {
		filter.Printer = printer
	}

	if err := c.ValidQuery(); err != nil {
		return c.RespondError(err)
	}

	list, count, err := uc.useCase.CashierGetMenuList(c.Ctx, filter)
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

func (uc Controller) CashierGetMenuDetail(c *web.Context) error {
	id := c.GetParam(reflect.Int, "id").(int)

	if err := c.ValidParam(); err != nil {
		return c.RespondError(err)
	}

	response, err := uc.useCase.CashierGetMenuDetail(c.Ctx, int64(id))
	if err != nil {
		return c.RespondError(err)
	}

	return c.Respond(map[string]interface{}{
		"data":   response,
		"status": true,
	}, http.StatusOK)
}

func (uc Controller) CashierCreateMenu(c *web.Context) error {
	var request menu_service.CashierCreateRequest

	if err := c.BindFunc(&request, "MenuID", "CashierID"); err != nil {
		return c.RespondError(err)
	}

	response, err := uc.useCase.CashierCreateMenu(c.Ctx, request)
	if err != nil {
		return c.RespondError(err)
	}

	return c.Respond(map[string]interface{}{
		"data":   response,
		"status": true,
	}, http.StatusOK)
}

func (uc Controller) CashierUpdateMenuAll(c *web.Context) error {
	id := c.GetParam(reflect.Int, "id").(int)

	if err := c.ValidParam(); err != nil {
		return c.RespondError(err)
	}

	var request menu_service.CashierUpdateRequest

	if err := c.BindFunc(&request, "MenuID", "CashierID", "Status"); err != nil {
		return c.RespondError(err)
	}

	request.ID = int64(id)

	err := uc.useCase.CashierUpdateMenu(c.Ctx, request)
	if err != nil {
		return c.RespondError(err)
	}

	return c.Respond(map[string]interface{}{
		"data":   "ok!",
		"status": true,
	}, http.StatusOK)
}

func (uc Controller) CashierUpdateMenuColumn(c *web.Context) error {
	id := c.GetParam(reflect.Int, "id").(int)

	if err := c.ValidParam(); err != nil {
		return c.RespondError(err)
	}

	var request menu_service.CashierUpdateRequest

	if err := c.BindFunc(&request); err != nil {
		return c.RespondError(err)
	}

	request.ID = int64(id)

	err := uc.useCase.CashierUpdateMenuColumn(c.Ctx, request)
	if err != nil {
		return c.RespondError(err)
	}

	return c.Respond(map[string]interface{}{
		"data":   "ok!",
		"status": true,
	}, http.StatusOK)
}

func (uc Controller) CashierDeleteMenu(c *web.Context) error {
	id := c.GetParam(reflect.Int, "id").(int)

	if err := c.ValidParam(); err != nil {
		return c.RespondError(err)
	}

	err := uc.useCase.CashierDeleteMenu(c.Ctx, int64(id))
	if err != nil {
		return c.RespondError(err)
	}

	return c.Respond(map[string]interface{}{
		"data":   "ok!",
		"status": true,
	}, http.StatusOK)
}

func (uc Controller) CashierUpdateMenuPrinterID(c *web.Context) error {
	var request menu_service.CashierUpdatePrinterIDRequest

	if err := c.BindFunc(&request, "PrinterID", "MenuIds"); err != nil {
		return c.RespondError(err)
	}

	err := uc.useCase.CashierUpdatePrinterID(c.Ctx, request)
	if err != nil {
		return c.RespondError(err)
	}

	return c.Respond(map[string]interface{}{
		"data":   "ok!",
		"status": true,
	}, http.StatusOK)
}

func (uc Controller) CashierDeleteMenuPrinterID(c *web.Context) error {
	id := c.GetParam(reflect.Int, "id").(int)

	if err := c.ValidParam(); err != nil {
		return c.RespondError(err)
	}

	err := uc.useCase.BranchDeletePrinterID(c.Ctx, int64(id))
	if err != nil {
		return c.RespondError(err)
	}

	return c.Respond(map[string]interface{}{
		"data":   "ok!",
		"status": true,
	}, http.StatusOK)
}

func (uc Controller) CashierRemoveMenuPhoto(c *web.Context) error {
	id := c.GetParam(reflect.Int, "id").(int)

	if err := c.ValidParam(); err != nil {
		return c.RespondError(err)
	}

	var request struct {
		Index int `json:"index"`
	}

	if err := c.BindFunc(&request, "Index"); err != nil {
		return c.RespondError(err)
	}

	err := uc.useCase.CashierRemoveMenuPhoto(c.Ctx, int64(id), request.Index)
	if err != nil {
		return c.RespondError(err)
	}

	return c.Respond(map[string]interface{}{
		"data":   "ok!",
		"status": true,
	}, http.StatusOK)
}

// @waiter

func (uc Controller) WaiterGetMenuList(c *web.Context) error {
	var filter menu_service.Filter

	if limit, ok := c.GetQueryFunc(reflect.Int, "limit").(*int); ok {
		filter.Limit = limit
	}
	if offset, ok := c.GetQueryFunc(reflect.Int, "offset").(*int); ok {
		filter.Offset = offset
	}
	if search, ok := c.GetQueryFunc(reflect.String, "search").(*string); ok {
		filter.Search = search
	}
	if categoryId, ok := c.GetQueryFunc(reflect.Int, "category_id").(*int); ok {
		filter.CategoryId = categoryId
	}

	if err := c.ValidQuery(); err != nil {
		return c.RespondError(err)
	}

	list, err := uc.useCase.WaiterGetMenuList(c.Ctx, filter)
	if err != nil {
		return c.RespondError(err)
	}

	return c.Respond(map[string]interface{}{
		"result": list,
		"status": true,
	}, http.StatusOK)
}

// #food_recipe

// @admin

func (uc Controller) AdminGetFoodRecipeList(c *web.Context) error {
	var filter food_recipe.Filter

	if limit, ok := c.GetQueryFunc(reflect.Int, "limit").(*int); ok {
		filter.Limit = limit
	}
	if offset, ok := c.GetQueryFunc(reflect.Int, "offset").(*int); ok {
		filter.Offset = offset
	}
	if name, ok := c.GetQueryFunc(reflect.String, "name").(*string); ok {
		filter.Name = name
	}
	if foodId, ok := c.GetQueryFunc(reflect.Int, "food_id").(*int); ok {
		filter.FoodId = foodId
	}

	if err := c.ValidQuery(); err != nil {
		return c.RespondError(err)
	}

	list, count, err := uc.useCase.AdminGetFoodRecipeList(c.Ctx, filter)
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

func (uc Controller) AdminGetFoodRecipeDetail(c *web.Context) error {
	id := c.GetParam(reflect.Int, "id").(int)

	if err := c.ValidParam(); err != nil {
		return c.RespondError(err)
	}

	response, err := uc.useCase.AdminGetFoodRecipeDetail(c.Ctx, int64(id))
	if err != nil {
		return c.RespondError(err)
	}

	return c.Respond(map[string]interface{}{
		"data":   response,
		"status": true,
	}, http.StatusOK)
}

func (uc Controller) AdminCreateFoodRecipe(c *web.Context) error {
	var request food_recipe.AdminCreateRequest

	if err := c.BindFunc(&request, "FoodId", "ProductId"); err != nil {
		return c.RespondError(err)
	}

	response, err := uc.useCase.AdminCreateFoodRecipe(c.Ctx, request)
	if err != nil {
		return c.RespondError(err)
	}

	return c.Respond(map[string]interface{}{
		"data":   response,
		"status": true,
	}, http.StatusOK)
}

func (uc Controller) AdminUpdateFoodRecipeAll(c *web.Context) error {
	id := c.GetParam(reflect.Int, "id").(int)

	if err := c.ValidParam(); err != nil {
		return c.RespondError(err)
	}

	var request food_recipe.AdminUpdateRequest

	if err := c.BindFunc(&request); err != nil {
		return c.RespondError(err)
	}

	request.ID = int64(id)

	err := uc.useCase.AdminUpdateFoodRecipeAll(c.Ctx, request)
	if err != nil {
		return c.RespondError(err)
	}

	return c.Respond(map[string]interface{}{
		"data":   "ok!",
		"status": true,
	}, http.StatusOK)
}

func (uc Controller) AdminUpdateFoodRecipeColumns(c *web.Context) error {
	id := c.GetParam(reflect.Int, "id").(int)

	if err := c.ValidParam(); err != nil {
		return c.RespondError(err)
	}

	var request food_recipe.AdminUpdateRequest

	if err := c.BindFunc(&request); err != nil {
		return c.RespondError(err)
	}

	request.ID = int64(id)

	err := uc.useCase.AdminUpdateFoodRecipeColumns(c.Ctx, request)
	if err != nil {
		return c.RespondError(err)
	}

	return c.Respond(map[string]interface{}{
		"data":   "ok!",
		"status": true,
	}, http.StatusOK)
}

func (uc Controller) AdminDeleteFoodRecipe(c *web.Context) error {
	id := c.GetParam(reflect.Int, "id").(int)

	if err := c.ValidParam(); err != nil {
		return c.RespondError(err)
	}

	err := uc.useCase.AdminDeleteFoodRecipe(c.Ctx, int64(id))
	if err != nil {
		return c.RespondError(err)
	}

	return c.Respond(map[string]interface{}{
		"data":   "ok!",
		"status": true,
	}, http.StatusOK)
}

// @branch

func (uc Controller) BranchGetFoodRecipeList(c *web.Context) error {
	var filter food_recipe.Filter

	if limit, ok := c.GetQueryFunc(reflect.Int, "limit").(*int); ok {
		filter.Limit = limit
	}
	if offset, ok := c.GetQueryFunc(reflect.Int, "offset").(*int); ok {
		filter.Offset = offset
	}
	if name, ok := c.GetQueryFunc(reflect.String, "name").(*string); ok {
		filter.Name = name
	}
	if foodId, ok := c.GetQueryFunc(reflect.Int, "food_id").(*int); ok {
		filter.FoodId = foodId
	}

	if err := c.ValidQuery(); err != nil {
		return c.RespondError(err)
	}

	list, count, err := uc.useCase.BranchGetFoodRecipeList(c.Ctx, filter)
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

func (uc Controller) BranchGetFoodRecipeDetail(c *web.Context) error {
	id := c.GetParam(reflect.Int, "id").(int)

	if err := c.ValidParam(); err != nil {
		return c.RespondError(err)
	}

	response, err := uc.useCase.BranchGetFoodRecipeDetail(c.Ctx, int64(id))
	if err != nil {
		return c.RespondError(err)
	}

	return c.Respond(map[string]interface{}{
		"data":   response,
		"status": true,
	}, http.StatusOK)
}

func (uc Controller) BranchCreateFoodRecipe(c *web.Context) error {
	var request food_recipe.BranchCreateRequest

	if err := c.BindFunc(&request, "FoodId", "ProductId"); err != nil {
		return c.RespondError(err)
	}

	response, err := uc.useCase.BranchCreateFoodRecipe(c.Ctx, request)
	if err != nil {
		return c.RespondError(err)
	}

	return c.Respond(map[string]interface{}{
		"data":   response,
		"status": true,
	}, http.StatusOK)
}

func (uc Controller) BranchUpdateFoodRecipeAll(c *web.Context) error {
	id := c.GetParam(reflect.Int, "id").(int)

	if err := c.ValidParam(); err != nil {
		return c.RespondError(err)
	}

	var request food_recipe.BranchUpdateRequest

	if err := c.BindFunc(&request); err != nil {
		return c.RespondError(err)
	}

	request.ID = int64(id)

	err := uc.useCase.BranchUpdateFoodRecipeAll(c.Ctx, request)
	if err != nil {
		return c.RespondError(err)
	}

	return c.Respond(map[string]interface{}{
		"data":   "ok!",
		"status": true,
	}, http.StatusOK)
}

func (uc Controller) BranchUpdateFoodRecipeColumns(c *web.Context) error {
	id := c.GetParam(reflect.Int, "id").(int)

	if err := c.ValidParam(); err != nil {
		return c.RespondError(err)
	}

	var request food_recipe.BranchUpdateRequest

	if err := c.BindFunc(&request); err != nil {
		return c.RespondError(err)
	}

	request.ID = int64(id)

	err := uc.useCase.BranchUpdateFoodRecipeColumns(c.Ctx, request)
	if err != nil {
		return c.RespondError(err)
	}

	return c.Respond(map[string]interface{}{
		"data":   "ok!",
		"status": true,
	}, http.StatusOK)
}

func (uc Controller) BranchDeleteFoodRecipe(c *web.Context) error {
	id := c.GetParam(reflect.Int, "id").(int)

	if err := c.ValidParam(); err != nil {
		return c.RespondError(err)
	}

	err := uc.useCase.BranchDeleteFoodRecipe(c.Ctx, int64(id))
	if err != nil {
		return c.RespondError(err)
	}

	return c.Respond(map[string]interface{}{
		"data":   "ok!",
		"status": true,
	}, http.StatusOK)
}

// #category ---------------------------------------------------------------

// @admin

func (uc Controller) SuperAdminGetCategoryList(c *web.Context) error {
	var filter category.Filter

	if limit, ok := c.GetQueryFunc(reflect.Int, "limit").(*int); ok {
		filter.Limit = limit
	}
	if offset, ok := c.GetQueryFunc(reflect.Int, "offset").(*int); ok {
		filter.Offset = offset
	}
	if page, ok := c.GetQueryFunc(reflect.Int, "page").(*int); ok {
		filter.Page = page
	}
	if name, ok := c.GetQueryFunc(reflect.String, "name").(*string); ok {
		filter.Name = name
	}

	if err := c.ValidQuery(); err != nil {
		return c.RespondError(err)
	}

	list, count, err := uc.useCase.SuperAdminGetCategoryList(c.Ctx, filter)
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

func (uc Controller) SuperAdminGetCategoryDetail(c *web.Context) error {
	id := c.GetParam(reflect.Int, "id").(int)

	if err := c.ValidParam(); err != nil {
		return c.RespondError(err)
	}

	response, err := uc.useCase.SuperAdminGetCategoryDetail(c.Ctx, int64(id))
	if err != nil {
		return c.RespondError(err)
	}

	return c.Respond(map[string]interface{}{
		"data":   response,
		"status": true,
	}, http.StatusOK)
}

func (uc Controller) SuperAdminCreateCategory(c *web.Context) error {
	var request category.SuperAdminCreateRequest

	if err := c.BindFunc(&request, "Name", "Logo"); err != nil {
		return c.RespondError(err)
	}

	response, err := uc.useCase.SuperAdminCreateCategory(c.Ctx, request)
	if err != nil {
		return c.RespondError(err)
	}

	return c.Respond(map[string]interface{}{
		"data":   response,
		"status": true,
	}, http.StatusOK)
}

func (uc Controller) SuperAdminUpdateCategoryAll(c *web.Context) error {
	id := c.GetParam(reflect.Int, "id").(int)

	if err := c.ValidParam(); err != nil {
		return c.RespondError(err)
	}

	var request category.SuperAdminUpdateRequest

	if err := c.BindFunc(&request, "Name", "Logo", "Main"); err != nil {
		return c.RespondError(err)
	}

	request.ID = int64(id)

	err := uc.useCase.SuperAdminUpdateCategory(c.Ctx, request)
	if err != nil {
		return c.RespondError(err)
	}

	return c.Respond(map[string]interface{}{
		"data":   "ok!",
		"status": true,
	}, http.StatusOK)
}

func (uc Controller) SuperAdminUpdateCategoryColumns(c *web.Context) error {
	id := c.GetParam(reflect.Int, "id").(int)

	if err := c.ValidParam(); err != nil {
		return c.RespondError(err)
	}

	var request category.SuperAdminUpdateRequest

	if err := c.BindFunc(&request); err != nil {
		return c.RespondError(err)
	}

	request.ID = int64(id)

	err := uc.useCase.SuperAdminUpdateCategoryColumn(c.Ctx, request)
	if err != nil {
		return c.RespondError(err)
	}

	return c.Respond(map[string]interface{}{
		"data":   "ok!",
		"status": true,
	}, http.StatusOK)
}

func (uc Controller) SuperAdminDeleteCategory(c *web.Context) error {
	id := c.GetParam(reflect.Int, "id").(int)

	if err := c.ValidParam(); err != nil {
		return c.RespondError(err)
	}

	err := uc.useCase.SuperAdminDeleteCategory(c.Ctx, int64(id))
	if err != nil {
		return c.RespondError(err)
	}

	return c.Respond(map[string]interface{}{
		"data":   "ok!",
		"status": true,
	}, http.StatusOK)
}

// @client

func (uc Controller) ClientGetCategoryList(c *web.Context) error {
	var filter category.Filter

	if limit, ok := c.GetQueryFunc(reflect.Int, "limit").(*int); ok {
		filter.Limit = limit
	}
	if offset, ok := c.GetQueryFunc(reflect.Int, "offset").(*int); ok {
		filter.Offset = offset
	}
	if page, ok := c.GetQueryFunc(reflect.Int, "page").(*int); ok {
		filter.Page = page
	}
	if name, ok := c.GetQueryFunc(reflect.String, "name").(*string); ok {
		filter.Name = name
	}

	if err := c.ValidQuery(); err != nil {
		return c.RespondMobileError(err)
	}

	list, count, err := uc.useCase.ClientGetCategoryList(c.Ctx, filter)
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

func (uc Controller) BranchGetCategoryList(c *web.Context) error {
	var filter category.Filter

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

	list, count, err := uc.useCase.BranchGetCategoryList(c.Ctx, filter)
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

// @cashier

func (uc Controller) CashierGetCategoryList(c *web.Context) error {
	var filter category.Filter

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

	list, count, err := uc.useCase.CashierGetCategoryList(c.Ctx, filter)
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

// @admin

func (uc Controller) AdminGetCategoryList(c *web.Context) error {
	var filter category.Filter

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

	list, count, err := uc.useCase.AdminGetCategoryList(c.Ctx, filter)
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

// @waiter

func (uc Controller) WaiterGetCategoryList(c *web.Context) error {
	list, err := uc.useCase.WaiterGetCategoryList(c.Ctx)
	if err != nil {
		return c.RespondError(err)
	}

	return c.Respond(map[string]interface{}{
		"data":   list,
		"status": true,
	}, http.StatusOK)
}

// #menu_category

// @admin

func (uc Controller) AdminGetMenuCategoryList(c *web.Context) error {
	var filter menu_category.Filter

	if limit, ok := c.GetQueryFunc(reflect.Int, "limit").(*int); ok {
		filter.Limit = limit
	}
	if offset, ok := c.GetQueryFunc(reflect.Int, "offset").(*int); ok {
		filter.Offset = offset
	}

	if err := c.ValidQuery(); err != nil {
		return c.RespondError(err)
	}

	list, count, err := uc.useCase.AdminGetListMenuCategory(c.Ctx, filter)
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

func (uc Controller) AdminGetMenuCategoryDetail(c *web.Context) error {
	id := c.GetParam(reflect.Int, "id").(int)

	if err := c.ValidParam(); err != nil {
		return c.RespondError(err)
	}
	response, err := uc.useCase.AdminGetDetailMenuCategory(c.Ctx, int64(id))
	if err != nil {
		return c.RespondError(err)
	}

	return c.Respond(map[string]interface{}{
		"data":   response,
		"status": true,
	}, http.StatusOK)
}

func (uc Controller) AdminCreateMenuCategory(c *web.Context) error {
	var request menu_category.AdminCreateRequest

	if err := c.BindFunc(&request, "FoodId", "ProductId"); err != nil {
		return c.RespondError(err)
	}

	response, err := uc.useCase.AdminCreateMenuCategory(c.Ctx, request)
	if err != nil {
		return c.RespondError(err)
	}

	return c.Respond(map[string]interface{}{
		"data":   response,
		"status": true,
	}, http.StatusOK)
}

func (uc Controller) AdminUpdateMenuCategory(c *web.Context) error {
	id := c.GetParam(reflect.Int, "id").(int)

	if err := c.ValidParam(); err != nil {
		return c.RespondError(err)
	}

	var request menu_category.AdminUpdateRequest

	if err := c.BindFunc(&request); err != nil {
		return c.RespondError(err)
	}

	request.ID = int64(id)

	err := uc.useCase.AdminUpdateMenuCategory(c.Ctx, request)
	if err != nil {
		return c.RespondError(err)
	}

	return c.Respond(map[string]interface{}{
		"data":   "ok!",
		"status": true,
	}, http.StatusOK)
}

func (uc Controller) AdminDeleteMenuCategory(c *web.Context) error {
	id := c.GetParam(reflect.Int, "id").(int)

	if err := c.ValidParam(); err != nil {
		return c.RespondError(err)
	}

	err := uc.useCase.AdminDeleteMenuCategory(c.Ctx, int64(id))
	if err != nil {
		return c.RespondError(err)
	}

	return c.Respond(map[string]interface{}{
		"data":   "ok!",
		"status": true,
	}, http.StatusOK)
}

// @branch

func (uc Controller) BranchGetMenuCategoryList(c *web.Context) error {
	var filter menu_category.Filter

	if limit, ok := c.GetQueryFunc(reflect.Int, "limit").(*int); ok {
		filter.Limit = limit
	}
	if offset, ok := c.GetQueryFunc(reflect.Int, "offset").(*int); ok {
		filter.Offset = offset
	}

	if err := c.ValidQuery(); err != nil {
		return c.RespondError(err)
	}

	list, count, err := uc.useCase.BranchGetListMenuCategory(c.Ctx, filter)
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

func (uc Controller) BranchGetMenuCategoryDetail(c *web.Context) error {
	id := c.GetParam(reflect.Int, "id").(int)

	if err := c.ValidParam(); err != nil {
		return c.RespondError(err)
	}

	response, err := uc.useCase.BranchGetDetailMenuCategory(c.Ctx, int64(id))
	if err != nil {
		return c.RespondError(err)
	}

	return c.Respond(map[string]interface{}{
		"data":   response,
		"status": true,
	}, http.StatusOK)
}

func (uc Controller) BranchCreateMenuCategory(c *web.Context) error {
	var request menu_category.BranchCreateRequest

	if err := c.BindFunc(&request); err != nil {
		return c.RespondError(err)
	}

	response, err := uc.useCase.BranchCreateMenuCategory(c.Ctx, request)
	if err != nil {
		return c.RespondError(err)
	}

	return c.Respond(map[string]interface{}{
		"data":   response,
		"status": true,
	}, http.StatusOK)
}

func (uc Controller) BranchUpdateMenuCategory(c *web.Context) error {
	id := c.GetParam(reflect.Int, "id").(int)

	if err := c.ValidParam(); err != nil {
		return c.RespondError(err)
	}

	var request menu_category.BranchUpdateRequest

	if err := c.BindFunc(&request); err != nil {
		return c.RespondError(err)
	}

	request.ID = int64(id)

	err := uc.useCase.BranchUpdateMenuCategory(c.Ctx, request)
	if err != nil {
		return c.RespondError(err)
	}

	return c.Respond(map[string]interface{}{
		"data":   "ok!",
		"status": true,
	}, http.StatusOK)
}

func (uc Controller) BranchDeleteMenuCategory(c *web.Context) error {
	id := c.GetParam(reflect.Int, "id").(int)

	if err := c.ValidParam(); err != nil {
		return c.RespondError(err)
	}

	err := uc.useCase.BranchDeleteMenuCategory(c.Ctx, int64(id))
	if err != nil {
		return c.RespondError(err)
	}

	return c.Respond(map[string]interface{}{
		"data":   "ok!",
		"status": true,
	}, http.StatusOK)
}

// @cashier

func (uc Controller) CashierGetMenuCategoryList(c *web.Context) error {
	var filter menu_category.Filter

	if limit, ok := c.GetQueryFunc(reflect.Int, "limit").(*int); ok {
		filter.Limit = limit
	}
	if offset, ok := c.GetQueryFunc(reflect.Int, "offset").(*int); ok {
		filter.Offset = offset
	}

	if err := c.ValidQuery(); err != nil {
		return c.RespondError(err)
	}

	list, count, err := uc.useCase.CashierGetListMenuCategory(c.Ctx, filter)
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

func (uc Controller) CashierGetMenuCategoryDetail(c *web.Context) error {
	id := c.GetParam(reflect.Int, "id").(int)

	if err := c.ValidParam(); err != nil {
		return c.RespondError(err)
	}

	response, err := uc.useCase.CashierGetDetailMenuCategory(c.Ctx, int64(id))
	if err != nil {
		return c.RespondError(err)
	}

	return c.Respond(map[string]interface{}{
		"data":   response,
		"status": true,
	}, http.StatusOK)
}

func (uc Controller) CashierCreateMenuCategory(c *web.Context) error {
	var request menu_category.CashierCreateRequest

	if err := c.BindFunc(&request, "FoodId", "ProductId"); err != nil {
		return c.RespondError(err)
	}

	response, err := uc.useCase.CashierCreateMenuCategory(c.Ctx, request)
	if err != nil {
		return c.RespondError(err)
	}

	return c.Respond(map[string]interface{}{
		"data":   response,
		"status": true,
	}, http.StatusOK)
}

func (uc Controller) CashierUpdateMenuCategory(c *web.Context) error {
	id := c.GetParam(reflect.Int, "id").(int)

	if err := c.ValidParam(); err != nil {
		return c.RespondError(err)
	}

	var request menu_category.CashierUpdateRequest

	if err := c.BindFunc(&request); err != nil {
		return c.RespondError(err)
	}

	request.ID = int64(id)

	err := uc.useCase.CashierUpdateMenuCategory(c.Ctx, request)
	if err != nil {
		return c.RespondError(err)
	}

	return c.Respond(map[string]interface{}{
		"data":   "ok!",
		"status": true,
	}, http.StatusOK)
}

func (uc Controller) CashierDeleteMenuCategory(c *web.Context) error {
	id := c.GetParam(reflect.Int, "id").(int)

	if err := c.ValidParam(); err != nil {
		return c.RespondError(err)
	}

	err := uc.useCase.CashierDeleteMenuCategory(c.Ctx, int64(id))
	if err != nil {
		return c.RespondError(err)
	}

	return c.Respond(map[string]interface{}{
		"data":   "ok!",
		"status": true,
	}, http.StatusOK)
}

// #food_recipe_group

// @admin

func (uc Controller) AdminGetFoodRecipeGroupList(c *web.Context) error {
	var filter food_recipe_group.Filter

	if limit, ok := c.GetQueryFunc(reflect.Int, "limit").(*int); ok {
		filter.Limit = limit
	}
	if offset, ok := c.GetQueryFunc(reflect.Int, "offset").(*int); ok {
		filter.Offset = offset
	}
	if name, ok := c.GetQueryFunc(reflect.String, "name").(*string); ok {
		filter.Name = name
	}
	if foodId, ok := c.GetQueryFunc(reflect.Int, "food_id").(*int); ok {
		filter.FoodId = foodId
	}

	if err := c.ValidQuery(); err != nil {
		return c.RespondError(err)
	}

	if filter.FoodId == nil {
		return c.RespondError(errors.New("food_id is required"))
	}

	list, count, err := uc.useCase.AdminGetFoodRecipeGroupList(c.Ctx, filter)
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

func (uc Controller) AdminGetFoodRecipeGroupDetail(c *web.Context) error {
	id := c.GetParam(reflect.Int, "id").(int)
	if err := c.ValidParam(); err != nil {
		return c.RespondError(err)
	}

	response, err := uc.useCase.AdminGetFoodRecipeGroupDetail(c.Ctx, int64(id))
	if err != nil {
		return c.RespondError(err)
	}

	return c.Respond(map[string]interface{}{
		"data":   response,
		"status": true,
	}, http.StatusOK)
}

func (uc Controller) AdminCreateFoodRecipeGroup(c *web.Context) error {
	var request food_recipe_group.AdminCreateRequest

	if err := c.BindFunc(&request, "FoodId", "ProductId"); err != nil {
		return c.RespondError(err)
	}

	response, err := uc.useCase.AdminCreateFoodRecipeGroup(c.Ctx, request)
	if err != nil {
		return c.RespondError(err)
	}

	return c.Respond(map[string]interface{}{
		"data":   response,
		"status": true,
	}, http.StatusOK)
}

func (uc Controller) AdminUpdateFoodRecipeGroupAll(c *web.Context) error {
	id := c.GetParam(reflect.Int, "id").(int)
	if err := c.ValidParam(); err != nil {
		return c.RespondError(err)
	}

	var request food_recipe_group.AdminUpdateRequest

	if err := c.BindFunc(&request); err != nil {
		return c.RespondError(err)
	}

	request.ID = int64(id)

	err := uc.useCase.AdminUpdateFoodRecipeGroupAll(c.Ctx, request)
	if err != nil {
		return c.RespondError(err)
	}

	return c.Respond(map[string]interface{}{
		"data":   "ok!",
		"status": true,
	}, http.StatusOK)
}

func (uc Controller) AdminUpdateFoodRecipeGroupColumns(c *web.Context) error {
	id := c.GetParam(reflect.Int, "id").(int)
	if err := c.ValidParam(); err != nil {
		return c.RespondError(err)
	}

	var request food_recipe_group.AdminUpdateRequest

	if err := c.BindFunc(&request); err != nil {
		return c.RespondError(err)
	}

	request.ID = int64(id)

	err := uc.useCase.AdminUpdateFoodRecipeGroupColumns(c.Ctx, request)
	if err != nil {
		return c.RespondError(err)
	}

	return c.Respond(map[string]interface{}{
		"data":   "ok!",
		"status": true,
	}, http.StatusOK)
}

func (uc Controller) AdminDeleteFoodRecipeGroupSingleRecipe(c *web.Context) error {
	id := c.GetParam(reflect.Int, "id").(int)

	if err := c.ValidParam(); err != nil {
		return c.RespondError(err)
	}

	var request food_recipe_group.AdminDeleteRecipeRequest

	if err := c.BindFunc(&request, "RecipeId"); err != nil {
		return c.RespondError(err)
	}

	request.ID = int64(id)

	err := uc.useCase.AdminDeleteFoodRecipeGroupSingleRecipe(c.Ctx, request)
	if err != nil {
		return c.RespondError(err)
	}

	return c.Respond(map[string]interface{}{
		"data":   "ok!",
		"status": true,
	}, http.StatusOK)
}

func (uc Controller) AdminDeleteFoodRecipeGroup(c *web.Context) error {
	id := c.GetParam(reflect.Int, "id").(int)

	if err := c.ValidParam(); err != nil {
		return c.RespondError(err)
	}

	err := uc.useCase.AdminDeleteFoodRecipeGroup(c.Ctx, int64(id))
	if err != nil {
		return c.RespondError(err)
	}

	return c.Respond(map[string]interface{}{
		"data":   "ok!",
		"status": true,
	}, http.StatusOK)
}

// @branch

func (uc Controller) BranchGetFoodRecipeGroupList(c *web.Context) error {
	var filter food_recipe_group.Filter

	if limit, ok := c.GetQueryFunc(reflect.Int, "limit").(*int); ok {
		filter.Limit = limit
	}
	if offset, ok := c.GetQueryFunc(reflect.Int, "offset").(*int); ok {
		filter.Offset = offset
	}
	if name, ok := c.GetQueryFunc(reflect.String, "name").(*string); ok {
		filter.Name = name
	}
	if foodId, ok := c.GetQueryFunc(reflect.Int, "food_id").(*int); ok {
		filter.FoodId = foodId
	}

	if err := c.ValidQuery(); err != nil {
		return c.RespondError(err)
	}

	if filter.FoodId == nil {
		return c.RespondError(errors.New("food_id is required"))
	}

	list, count, err := uc.useCase.BranchGetFoodRecipeGroupList(c.Ctx, filter)
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

func (uc Controller) BranchGetFoodRecipeGroupDetail(c *web.Context) error {
	id := c.GetParam(reflect.Int, "id").(int)

	if err := c.ValidParam(); err != nil {
		return c.RespondError(err)
	}

	response, err := uc.useCase.BranchGetFoodRecipeGroupDetail(c.Ctx, int64(id))
	if err != nil {
		return c.RespondError(err)
	}

	return c.Respond(map[string]interface{}{
		"data":   response,
		"status": true,
	}, http.StatusOK)
}

func (uc Controller) BranchCreateFoodRecipeGroup(c *web.Context) error {
	var request food_recipe_group.BranchCreateRequest

	if err := c.BindFunc(&request, "FoodId", "ProductId"); err != nil {
		return c.RespondError(err)
	}

	response, err := uc.useCase.BranchCreateFoodRecipeGroup(c.Ctx, request)
	if err != nil {
		return c.RespondError(err)
	}

	return c.Respond(map[string]interface{}{
		"data":   response,
		"status": true,
	}, http.StatusOK)
}

func (uc Controller) BranchUpdateFoodRecipeGroupAll(c *web.Context) error {
	id := c.GetParam(reflect.Int, "id").(int)
	if err := c.ValidParam(); err != nil {
		return c.RespondError(err)
	}

	var request food_recipe_group.BranchUpdateRequest

	if err := c.BindFunc(&request); err != nil {
		return c.RespondError(err)
	}

	request.ID = int64(id)

	err := uc.useCase.BranchUpdateFoodRecipeGroupAll(c.Ctx, request)
	if err != nil {
		return c.RespondError(err)
	}

	return c.Respond(map[string]interface{}{
		"data":   "ok!",
		"status": true,
	}, http.StatusOK)
}

func (uc Controller) BranchUpdateFoodRecipeGroupColumns(c *web.Context) error {
	id := c.GetParam(reflect.Int, "id").(int)
	if err := c.ValidParam(); err != nil {
		return c.RespondError(err)
	}

	var request food_recipe_group.BranchUpdateRequest

	if err := c.BindFunc(&request); err != nil {
		return c.RespondError(err)
	}

	request.ID = int64(id)

	err := uc.useCase.BranchUpdateFoodRecipeGroupColumns(c.Ctx, request)
	if err != nil {
		return c.RespondError(err)
	}

	return c.Respond(map[string]interface{}{
		"data":   "ok!",
		"status": true,
	}, http.StatusOK)
}

func (uc Controller) BranchDeleteFoodRecipeGroupSingleRecipe(c *web.Context) error {
	id := c.GetParam(reflect.Int, "id").(int)
	if err := c.ValidParam(); err != nil {
		return c.RespondError(err)
	}

	var request food_recipe_group.BranchDeleteRecipeRequest
	if err := c.BindFunc(&request, "RecipeId"); err != nil {
		return c.RespondError(err)
	}

	request.ID = int64(id)

	err := uc.useCase.BranchDeleteFoodRecipeGroupSingleRecipe(c.Ctx, request)
	if err != nil {
		return c.RespondError(err)
	}

	return c.Respond(map[string]interface{}{
		"data":   "ok!",
		"status": true,
	}, http.StatusOK)
}

func (uc Controller) BranchDeleteFoodRecipeGroup(c *web.Context) error {
	id := c.GetParam(reflect.Int, "id").(int)

	if err := c.ValidParam(); err != nil {
		return c.RespondError(err)
	}

	err := uc.useCase.BranchDeleteFoodRecipeGroup(c.Ctx, int64(id))
	if err != nil {
		return c.RespondError(err)
	}

	return c.Respond(map[string]interface{}{
		"data":   "ok!",
		"status": true,
	}, http.StatusOK)
}
