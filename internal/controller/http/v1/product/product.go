package product

import (
	"github.com/restaurant/foundation/web"
	"github.com/restaurant/internal/service/product_recipe"
	"net/http"
	"reflect"
	"time"
)

type Controller struct {
	useCase *product.UseCase
}

func NewController(useCase *product.UseCase) *Controller {
	return &Controller{useCase}
}

// product

// @admin

func (uc Controller) AdminGetProductList(c *web.Context) error {
	var filter product_service.Filter

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

	list, count, err := uc.useCase.AdminGetProductList(c.Ctx, filter)
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

func (uc Controller) AdminGetProductDetail(c *web.Context) error {
	id := c.GetParam(reflect.Int, "id").(int)

	if err := c.ValidParam(); err != nil {
		return c.RespondError(err)
	}

	response, err := uc.useCase.AdminGetProductDetail(c.Ctx, int64(id))
	if err != nil {
		return c.RespondError(err)
	}

	return c.Respond(map[string]interface{}{
		"data":   response,
		"status": true,
	}, http.StatusOK)
}

func (uc Controller) AdminCreateProduct(c *web.Context) error {
	var request product_service.AdminCreateRequest

	if err := c.BindFunc(&request, "Name", "MeasureUnitID"); err != nil {
		return c.RespondError(err)
	}

	response, err := uc.useCase.AdminCreateProduct(c.Ctx, request)
	if err != nil {
		return c.RespondError(err)
	}

	return c.Respond(map[string]interface{}{
		"data":   response,
		"status": true,
	}, http.StatusOK)
}

func (uc Controller) AdminUpdateProductAll(c *web.Context) error {
	id := c.GetParam(reflect.Int, "id").(int)

	if err := c.ValidParam(); err != nil {
		return c.RespondError(err)
	}

	var request product_service.AdminUpdateRequest

	if err := c.BindFunc(&request, "ID", "Name", "MeasureUnitID"); err != nil {
		return c.RespondError(err)
	}

	request.ID = int64(id)

	err := uc.useCase.AdminUpdateProduct(c.Ctx, request)
	if err != nil {
		return c.RespondError(err)
	}

	return c.Respond(map[string]interface{}{
		"data":   "ok!",
		"status": true,
	}, http.StatusOK)
}

func (uc Controller) AdminUpdateProductColumns(c *web.Context) error {
	id := c.GetParam(reflect.Int, "id").(int)

	if err := c.ValidParam(); err != nil {
		return c.RespondError(err)
	}

	var request product_service.AdminUpdateRequest

	if err := c.BindFunc(&request); err != nil {
		return c.RespondError(err)
	}

	request.ID = int64(id)

	err := uc.useCase.AdminUpdateProductColumn(c.Ctx, request)
	if err != nil {
		return c.RespondError(err)
	}

	return c.Respond(map[string]interface{}{
		"data":   "ok!",
		"status": true,
	}, http.StatusOK)
}

func (uc Controller) AdminDeleteProduct(c *web.Context) error {
	id := c.GetParam(reflect.Int, "id").(int)

	if err := c.ValidParam(); err != nil {
		return c.RespondError(err)
	}

	err := uc.useCase.AdminDeleteProduct(c.Ctx, int64(id))
	if err != nil {
		return c.RespondError(err)
	}

	return c.Respond(map[string]interface{}{
		"data":   "ok!",
		"status": true,
	}, http.StatusOK)
}

func (uc Controller) AdminGetProductSpendingByBranch(c *web.Context) error {
	var filter product_service.SpendingFilter

	if limit, ok := c.GetQueryFunc(reflect.Int, "limit").(*int); ok {
		filter.Limit = limit
	}
	if offset, ok := c.GetQueryFunc(reflect.Int, "offset").(*int); ok {
		filter.Offset = offset
	}
	if fromDate, ok := c.GetQueryFunc(reflect.String, "from_date").(*time.Time); ok {
		filter.FromDate = fromDate
	}
	if toDate, ok := c.GetQueryFunc(reflect.String, "to_date").(*time.Time); ok {
		filter.ToDate = toDate
	}
	if branchId, ok := c.GetQueryFunc(reflect.Int, "branch_id").(*int); ok {
		filter.BranchId = branchId
	}

	if err := c.ValidQuery(); err != nil {
		return c.RespondError(err)
	}

	list, err := uc.useCase.AdminGetProductSpendingByBranch(c.Ctx, filter)
	if err != nil {
		return c.RespondError(err)
	}

	return c.Respond(map[string]interface{}{
		"data":   list,
		"status": true,
	}, http.StatusOK)
}

// @branch

func (uc Controller) BranchGetProductList(c *web.Context) error {
	var filter product_service.Filter

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

	list, count, err := uc.useCase.BranchGetProductList(c.Ctx, filter)
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

func (uc Controller) BranchGetProductDetail(c *web.Context) error {
	id := c.GetParam(reflect.Int, "id").(int)

	if err := c.ValidParam(); err != nil {
		return c.RespondError(err)
	}

	response, err := uc.useCase.BranchGetProductDetail(c.Ctx, int64(id))
	if err != nil {
		return c.RespondError(err)
	}

	return c.Respond(map[string]interface{}{
		"data":   response,
		"status": true,
	}, http.StatusOK)
}

func (uc Controller) BranchCreateProduct(c *web.Context) error {
	var request product_service.BranchCreateRequest

	if err := c.BindFunc(&request, "Name", "MeasureUnitID"); err != nil {
		return c.RespondError(err)
	}

	response, err := uc.useCase.BranchCreateProduct(c.Ctx, request)
	if err != nil {
		return c.RespondError(err)
	}

	return c.Respond(map[string]interface{}{
		"data":   response,
		"status": true,
	}, http.StatusOK)
}

func (uc Controller) BranchUpdateProductAll(c *web.Context) error {
	id := c.GetParam(reflect.Int, "id").(int)

	if err := c.ValidParam(); err != nil {
		return c.RespondError(err)
	}

	var request product_service.BranchUpdateRequest

	if err := c.BindFunc(&request, "Name", "MeasureUnitID"); err != nil {
		return c.RespondError(err)
	}

	request.ID = int64(id)

	err := uc.useCase.BranchUpdateProduct(c.Ctx, request)
	if err != nil {
		return c.RespondError(err)
	}

	return c.Respond(map[string]interface{}{
		"data":   "ok!",
		"status": true,
	}, http.StatusOK)
}

func (uc Controller) BranchUpdateProductColumns(c *web.Context) error {
	id := c.GetParam(reflect.Int, "id").(int)

	if err := c.ValidParam(); err != nil {
		return c.RespondError(err)
	}

	var request product_service.BranchUpdateRequest

	if err := c.BindFunc(&request); err != nil {
		return c.RespondError(err)
	}

	request.ID = int64(id)

	err := uc.useCase.BranchUpdateProductColumn(c.Ctx, request)
	if err != nil {
		return c.RespondError(err)
	}

	return c.Respond(map[string]interface{}{
		"data":   "ok!",
		"status": true,
	}, http.StatusOK)
}

func (uc Controller) BranchDeleteProduct(c *web.Context) error {
	id := c.GetParam(reflect.Int, "id").(int)

	if err := c.ValidParam(); err != nil {
		return c.RespondError(err)
	}

	err := uc.useCase.BranchDeleteProduct(c.Ctx, int64(id))
	if err != nil {
		return c.RespondError(err)
	}

	return c.Respond(map[string]interface{}{
		"data":   "ok!",
		"status": true,
	}, http.StatusOK)
}

// @cashier

func (uc Controller) CashierGetProductList(c *web.Context) error {
	var filter product_service.Filter

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

	list, count, err := uc.useCase.CashierGetProductList(c.Ctx, filter)
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

func (uc Controller) CashierGetProductDetail(c *web.Context) error {
	id := c.GetParam(reflect.Int, "id").(int)

	if err := c.ValidParam(); err != nil {
		return c.RespondError(err)
	}

	response, err := uc.useCase.CashierGetProductDetail(c.Ctx, int64(id))
	if err != nil {
		return c.RespondError(err)
	}

	return c.Respond(map[string]interface{}{
		"data":   response,
		"status": true,
	}, http.StatusOK)
}

func (uc Controller) CashierCreateProduct(c *web.Context) error {
	var request product_service.CashierCreateRequest

	if err := c.BindFunc(&request, "Name", "MeasureUnitID"); err != nil {
		return c.RespondError(err)
	}

	response, err := uc.useCase.CashierCreateProduct(c.Ctx, request)
	if err != nil {
		return c.RespondError(err)
	}

	return c.Respond(map[string]interface{}{
		"data":   response,
		"status": true,
	}, http.StatusOK)
}

func (uc Controller) CashierUpdateProductAll(c *web.Context) error {
	id := c.GetParam(reflect.Int, "id").(int)

	if err := c.ValidParam(); err != nil {
		return c.RespondError(err)
	}

	var request product_service.CashierUpdateRequest

	if err := c.BindFunc(&request, "Name", "MeasureUnitID"); err != nil {
		return c.RespondError(err)
	}

	request.ID = int64(id)

	err := uc.useCase.CashierUpdateProduct(c.Ctx, request)
	if err != nil {
		return c.RespondError(err)
	}

	return c.Respond(map[string]interface{}{
		"data":   "ok!",
		"status": true,
	}, http.StatusOK)
}

func (uc Controller) CashierUpdateProductColumns(c *web.Context) error {
	id := c.GetParam(reflect.Int, "id").(int)

	if err := c.ValidParam(); err != nil {
		return c.RespondError(err)
	}

	var request product_service.CashierUpdateRequest

	if err := c.BindFunc(&request); err != nil {
		return c.RespondError(err)
	}

	request.ID = int64(id)

	err := uc.useCase.CashierUpdateProductColumn(c.Ctx, request)
	if err != nil {
		return c.RespondError(err)
	}

	return c.Respond(map[string]interface{}{
		"data":   "ok!",
		"status": true,
	}, http.StatusOK)
}

func (uc Controller) CashierDeleteProduct(c *web.Context) error {
	id := c.GetParam(reflect.Int, "id").(int)

	if err := c.ValidParam(); err != nil {
		return c.RespondError(err)
	}

	err := uc.useCase.CashierDeleteProduct(c.Ctx, int64(id))
	if err != nil {
		return c.RespondError(err)
	}

	return c.Respond(map[string]interface{}{
		"data":   "ok!",
		"status": true,
	}, http.StatusOK)
}

func (uc Controller) CashierGetProductSpending(c *web.Context) error {
	var filter product_service.SpendingFilter

	if limit, ok := c.GetQueryFunc(reflect.Int, "limit").(*int); ok {
		filter.Limit = limit
	}
	if offset, ok := c.GetQueryFunc(reflect.Int, "offset").(*int); ok {
		filter.Offset = offset
	}
	if fromDate, ok := c.GetQueryFunc(reflect.String, "from_date").(*time.Time); ok {
		filter.FromDate = fromDate
	}
	if toDate, ok := c.GetQueryFunc(reflect.String, "to_date").(*time.Time); ok {
		filter.ToDate = toDate
	}

	if err := c.ValidQuery(); err != nil {
		return c.RespondError(err)
	}

	list, err := uc.useCase.CashierGetProductSpending(c.Ctx, filter)
	if err != nil {
		return c.RespondError(err)
	}

	return c.Respond(map[string]interface{}{
		"data":   list,
		"status": true,
	}, http.StatusOK)
}

// #product-recipe

// @admin

func (uc Controller) AdminGetProductRecipeList(c *web.Context) error {
	var filter product_recipe.Filter

	if limit, ok := c.GetQueryFunc(reflect.Int, "limit").(*int); ok {
		filter.Limit = limit
	}
	if offset, ok := c.GetQueryFunc(reflect.Int, "offset").(*int); ok {
		filter.Offset = offset
	}
	if name, ok := c.GetQueryFunc(reflect.String, "name").(*string); ok {
		filter.Name = name
	}
	if productId, ok := c.GetQueryFunc(reflect.Int, "product_id").(*int); ok {
		filter.ProductId = productId
	}

	if err := c.ValidQuery(); err != nil {
		return c.RespondError(err)
	}

	list, count, err := uc.useCase.AdminGetProductRecipeList(c.Ctx, filter)
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

func (uc Controller) AdminGetProductRecipeDetail(c *web.Context) error {
	id := c.GetParam(reflect.Int, "id").(int)

	if err := c.ValidParam(); err != nil {
		return c.RespondError(err)
	}

	response, err := uc.useCase.AdminGetProductRecipeDetail(c.Ctx, int64(id))
	if err != nil {
		return c.RespondError(err)
	}

	return c.Respond(map[string]interface{}{
		"data":   response,
		"status": true,
	}, http.StatusOK)
}

func (uc Controller) AdminCreateProductRecipe(c *web.Context) error {
	var request product_recipe.AdminCreateRequest

	if err := c.BindFunc(&request, "RecipeId", "ProductId"); err != nil {
		return c.RespondError(err)
	}

	response, err := uc.useCase.AdminCreateProductRecipe(c.Ctx, request)
	if err != nil {
		return c.RespondError(err)
	}

	return c.Respond(map[string]interface{}{
		"data":   response,
		"status": true,
	}, http.StatusOK)
}

func (uc Controller) AdminUpdateProductRecipeAll(c *web.Context) error {
	id := c.GetParam(reflect.Int, "id").(int)

	if err := c.ValidParam(); err != nil {
		return c.RespondError(err)
	}

	var request product_recipe.AdminUpdateRequest

	if err := c.BindFunc(&request); err != nil {
		return c.RespondError(err)
	}

	request.ID = int64(id)

	err := uc.useCase.AdminUpdateProductRecipeAll(c.Ctx, request)
	if err != nil {
		return c.RespondError(err)
	}

	return c.Respond(map[string]interface{}{
		"data":   "ok!",
		"status": true,
	}, http.StatusOK)
}

func (uc Controller) AdminUpdateProductRecipeColumns(c *web.Context) error {
	id := c.GetParam(reflect.Int, "id").(int)

	if err := c.ValidParam(); err != nil {
		return c.RespondError(err)
	}

	var request product_recipe.AdminUpdateRequest

	if err := c.BindFunc(&request); err != nil {
		return c.RespondError(err)
	}

	request.ID = int64(id)

	err := uc.useCase.AdminUpdateProductRecipeColumns(c.Ctx, request)
	if err != nil {
		return c.RespondError(err)
	}

	return c.Respond(map[string]interface{}{
		"data":   "ok!",
		"status": true,
	}, http.StatusOK)
}

func (uc Controller) AdminDeleteProductRecipe(c *web.Context) error {
	id := c.GetParam(reflect.Int, "id").(int)

	if err := c.ValidParam(); err != nil {
		return c.RespondError(err)
	}

	err := uc.useCase.AdminDeleteProductRecipe(c.Ctx, int64(id))
	if err != nil {
		return c.RespondError(err)
	}

	return c.Respond(map[string]interface{}{
		"data":   "ok!",
		"status": true,
	}, http.StatusOK)
}

// @branch

func (uc Controller) BranchGetProductRecipeList(c *web.Context) error {
	var filter product_recipe.Filter

	if limit, ok := c.GetQueryFunc(reflect.Int, "limit").(*int); ok {
		filter.Limit = limit
	}
	if offset, ok := c.GetQueryFunc(reflect.Int, "offset").(*int); ok {
		filter.Offset = offset
	}
	if name, ok := c.GetQueryFunc(reflect.String, "name").(*string); ok {
		filter.Name = name
	}
	if productId, ok := c.GetQueryFunc(reflect.Int, "product_id").(*int); ok {
		filter.ProductId = productId
	}

	if err := c.ValidQuery(); err != nil {
		return c.RespondError(err)
	}

	list, count, err := uc.useCase.BranchGetProductRecipeList(c.Ctx, filter)
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

func (uc Controller) BranchGetProductRecipeDetail(c *web.Context) error {
	id := c.GetParam(reflect.Int, "id").(int)

	if err := c.ValidParam(); err != nil {
		return c.RespondError(err)
	}

	response, err := uc.useCase.BranchGetProductRecipeDetail(c.Ctx, int64(id))
	if err != nil {
		return c.RespondError(err)
	}

	return c.Respond(map[string]interface{}{
		"data":   response,
		"status": true,
	}, http.StatusOK)
}

func (uc Controller) BranchCreateProductRecipe(c *web.Context) error {
	var request product_recipe.BranchCreateRequest

	if err := c.BindFunc(&request, "RecipeId", "ProductId"); err != nil {
		return c.RespondError(err)
	}

	response, err := uc.useCase.BranchCreateProductRecipe(c.Ctx, request)
	if err != nil {
		return c.RespondError(err)
	}

	return c.Respond(map[string]interface{}{
		"data":   response,
		"status": true,
	}, http.StatusOK)
}

func (uc Controller) BranchUpdateProductRecipeAll(c *web.Context) error {
	id := c.GetParam(reflect.Int, "id").(int)

	if err := c.ValidParam(); err != nil {
		return c.RespondError(err)
	}

	var request product_recipe.BranchUpdateRequest

	if err := c.BindFunc(&request); err != nil {
		return c.RespondError(err)
	}

	request.ID = int64(id)

	err := uc.useCase.BranchUpdateProductRecipeAll(c.Ctx, request)
	if err != nil {
		return c.RespondError(err)
	}

	return c.Respond(map[string]interface{}{
		"data":   "ok!",
		"status": true,
	}, http.StatusOK)
}

func (uc Controller) BranchUpdateProductRecipeColumns(c *web.Context) error {
	id := c.GetParam(reflect.Int, "id").(int)

	if err := c.ValidParam(); err != nil {
		return c.RespondError(err)
	}

	var request product_recipe.BranchUpdateRequest

	if err := c.BindFunc(&request); err != nil {
		return c.RespondError(err)
	}

	request.ID = int64(id)

	err := uc.useCase.BranchUpdateProductRecipeColumns(c.Ctx, request)
	if err != nil {
		return c.RespondError(err)
	}

	return c.Respond(map[string]interface{}{
		"data":   "ok!",
		"status": true,
	}, http.StatusOK)
}

func (uc Controller) BranchDeleteProductRecipe(c *web.Context) error {
	id := c.GetParam(reflect.Int, "id").(int)

	if err := c.ValidParam(); err != nil {
		return c.RespondError(err)
	}

	err := uc.useCase.BranchDeleteProductRecipe(c.Ctx, int64(id))
	if err != nil {
		return c.RespondError(err)
	}

	return c.Respond(map[string]interface{}{
		"data":   "ok!",
		"status": true,
	}, http.StatusOK)
}

// #food_recipe_group---------------------------------------------------------

// @admin

func (uc Controller) AdminGetProductRecipeGroupList(c *web.Context) error {
	var filter product_recipe_group.Filter

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
		filter.ProductID = foodId
	}

	if err := c.ValidQuery(); err != nil {
		return c.RespondError(err)
	}

	if filter.ProductID == nil {
		return c.RespondError(errors.New("food_id is required"))
	}

	list, count, err := uc.useCase.AdminGetProductRecipeGroupList(c.Ctx, filter)
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

func (uc Controller) AdminGetProductRecipeGroupDetail(c *web.Context) error {
	id := c.GetParam(reflect.Int, "id").(int)
	if err := c.ValidParam(); err != nil {
		return c.RespondError(err)
	}

	response, err := uc.useCase.AdminGetProductRecipeGroupDetail(c.Ctx, int64(id))
	if err != nil {
		return c.RespondError(err)
	}

	return c.Respond(map[string]interface{}{
		"data":   response,
		"status": true,
	}, http.StatusOK)
}

func (uc Controller) AdminCreateProductRecipeGroup(c *web.Context) error {
	var request product_recipe_group.AdminCreateRequest

	if err := c.BindFunc(&request, "ProductID", "ProductId"); err != nil {
		return c.RespondError(err)
	}

	response, err := uc.useCase.AdminCreateProductRecipeGroup(c.Ctx, request)
	if err != nil {
		return c.RespondError(err)
	}

	return c.Respond(map[string]interface{}{
		"data":   response,
		"status": true,
	}, http.StatusOK)
}

func (uc Controller) AdminUpdateProductRecipeGroupAll(c *web.Context) error {
	id := c.GetParam(reflect.Int, "id").(int)
	if err := c.ValidParam(); err != nil {
		return c.RespondError(err)
	}

	var request product_recipe_group.AdminUpdateRequest

	if err := c.BindFunc(&request); err != nil {
		return c.RespondError(err)
	}

	request.ID = int64(id)

	err := uc.useCase.AdminUpdateProductRecipeGroupAll(c.Ctx, request)
	if err != nil {
		return c.RespondError(err)
	}

	return c.Respond(map[string]interface{}{
		"data":   "ok!",
		"status": true,
	}, http.StatusOK)
}

func (uc Controller) AdminUpdateProductRecipeGroupColumns(c *web.Context) error {
	id := c.GetParam(reflect.Int, "id").(int)
	if err := c.ValidParam(); err != nil {
		return c.RespondError(err)
	}

	var request product_recipe_group.AdminUpdateRequest

	if err := c.BindFunc(&request); err != nil {
		return c.RespondError(err)
	}

	request.ID = int64(id)

	err := uc.useCase.AdminUpdateProductRecipeGroupColumns(c.Ctx, request)
	if err != nil {
		return c.RespondError(err)
	}

	return c.Respond(map[string]interface{}{
		"data":   "ok!",
		"status": true,
	}, http.StatusOK)
}

func (uc Controller) AdminDeleteProductRecipeGroupSingleRecipe(c *web.Context) error {
	id := c.GetParam(reflect.Int, "id").(int)

	if err := c.ValidParam(); err != nil {
		return c.RespondError(err)
	}

	var request product_recipe_group.AdminDeleteRecipeRequest

	if err := c.BindFunc(&request, "RecipeId"); err != nil {
		return c.RespondError(err)
	}

	request.ID = int64(id)

	err := uc.useCase.AdminDeleteProductRecipeGroupSingleRecipe(c.Ctx, request)
	if err != nil {
		return c.RespondError(err)
	}

	return c.Respond(map[string]interface{}{
		"data":   "ok!",
		"status": true,
	}, http.StatusOK)
}

func (uc Controller) AdminDeleteProductRecipeGroup(c *web.Context) error {
	id := c.GetParam(reflect.Int, "id").(int)

	if err := c.ValidParam(); err != nil {
		return c.RespondError(err)
	}

	err := uc.useCase.AdminDeleteProductRecipeGroup(c.Ctx, int64(id))
	if err != nil {
		return c.RespondError(err)
	}

	return c.Respond(map[string]interface{}{
		"data":   "ok!",
		"status": true,
	}, http.StatusOK)
}

// @branch

func (uc Controller) BranchGetProductRecipeGroupList(c *web.Context) error {
	var filter product_recipe_group.Filter

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
		filter.ProductID = foodId
	}

	if err := c.ValidQuery(); err != nil {
		return c.RespondError(err)
	}

	if filter.ProductID == nil {
		return c.RespondError(errors.New("food_id is required"))
	}

	list, count, err := uc.useCase.BranchGetProductRecipeGroupList(c.Ctx, filter)
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

func (uc Controller) BranchGetProductRecipeGroupDetail(c *web.Context) error {
	id := c.GetParam(reflect.Int, "id").(int)

	if err := c.ValidParam(); err != nil {
		return c.RespondError(err)
	}

	response, err := uc.useCase.BranchGetProductRecipeGroupDetail(c.Ctx, int64(id))
	if err != nil {
		return c.RespondError(err)
	}

	return c.Respond(map[string]interface{}{
		"data":   response,
		"status": true,
	}, http.StatusOK)
}

func (uc Controller) BranchCreateProductRecipeGroup(c *web.Context) error {
	var request product_recipe_group.BranchCreateRequest

	if err := c.BindFunc(&request, "ProductID", "ProductId"); err != nil {
		return c.RespondError(err)
	}

	response, err := uc.useCase.BranchCreateProductRecipeGroup(c.Ctx, request)
	if err != nil {
		return c.RespondError(err)
	}

	return c.Respond(map[string]interface{}{
		"data":   response,
		"status": true,
	}, http.StatusOK)
}

func (uc Controller) BranchUpdateProductRecipeGroupAll(c *web.Context) error {
	id := c.GetParam(reflect.Int, "id").(int)
	if err := c.ValidParam(); err != nil {
		return c.RespondError(err)
	}

	var request product_recipe_group.BranchUpdateRequest

	if err := c.BindFunc(&request); err != nil {
		return c.RespondError(err)
	}

	request.ID = int64(id)

	err := uc.useCase.BranchUpdateProductRecipeGroupAll(c.Ctx, request)
	if err != nil {
		return c.RespondError(err)
	}

	return c.Respond(map[string]interface{}{
		"data":   "ok!",
		"status": true,
	}, http.StatusOK)
}

func (uc Controller) BranchUpdateProductRecipeGroupColumns(c *web.Context) error {
	id := c.GetParam(reflect.Int, "id").(int)
	if err := c.ValidParam(); err != nil {
		return c.RespondError(err)
	}

	var request product_recipe_group.BranchUpdateRequest

	if err := c.BindFunc(&request); err != nil {
		return c.RespondError(err)
	}

	request.ID = int64(id)

	err := uc.useCase.BranchUpdateProductRecipeGroupColumns(c.Ctx, request)
	if err != nil {
		return c.RespondError(err)
	}

	return c.Respond(map[string]interface{}{
		"data":   "ok!",
		"status": true,
	}, http.StatusOK)
}

func (uc Controller) BranchDeleteProductRecipeGroupSingleRecipe(c *web.Context) error {
	id := c.GetParam(reflect.Int, "id").(int)
	if err := c.ValidParam(); err != nil {
		return c.RespondError(err)
	}

	var request product_recipe_group.BranchDeleteRecipeRequest
	if err := c.BindFunc(&request, "RecipeId"); err != nil {
		return c.RespondError(err)
	}

	request.ID = int64(id)

	err := uc.useCase.BranchDeleteProductRecipeGroupSingleRecipe(c.Ctx, request)
	if err != nil {
		return c.RespondError(err)
	}

	return c.Respond(map[string]interface{}{
		"data":   "ok!",
		"status": true,
	}, http.StatusOK)
}

func (uc Controller) BranchDeleteProductRecipeGroup(c *web.Context) error {
	id := c.GetParam(reflect.Int, "id").(int)

	if err := c.ValidParam(); err != nil {
		return c.RespondError(err)
	}

	err := uc.useCase.BranchDeleteProductRecipeGroup(c.Ctx, int64(id))
	if err != nil {
		return c.RespondError(err)
	}

	return c.Respond(map[string]interface{}{
		"data":   "ok!",
		"status": true,
	}, http.StatusOK)
}

// #food_recipe_group_history---------------------------------------------------------

// @admin

func (uc Controller) AdminGetProductRecipeGroupHistoryList(c *web.Context) error {
	var filter product_recipe_group_history.Filter

	if limit, ok := c.GetQueryFunc(reflect.Int, "limit").(*int); ok {
		filter.Limit = limit
	}
	if offset, ok := c.GetQueryFunc(reflect.Int, "offset").(*int); ok {
		filter.Offset = offset
	}
	if foodId, ok := c.GetQueryFunc(reflect.Int, "food_id").(*int); ok {
		filter.ProductID = foodId
	}

	if err := c.ValidQuery(); err != nil {
		return c.RespondError(err)
	}

	if filter.ProductID == nil {
		return c.RespondError(errors.New("food_id is required"))
	}

	list, count, err := uc.useCase.AdminGetProductRecipeGroupHistoryList(c.Ctx, filter)
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

func (uc Controller) AdminGetProductRecipeGroupHistoryDetail(c *web.Context) error {
	id := c.GetParam(reflect.Int, "id").(int)
	if err := c.ValidParam(); err != nil {
		return c.RespondError(err)
	}

	response, err := uc.useCase.AdminGetProductRecipeGroupHistoryDetail(c.Ctx, int64(id))
	if err != nil {
		return c.RespondError(err)
	}

	return c.Respond(map[string]interface{}{
		"data":   response,
		"status": true,
	}, http.StatusOK)
}

func (uc Controller) AdminCreateProductRecipeGroupHistory(c *web.Context) error {
	var request product_recipe_group_history.AdminCreateRequest

	if err := c.BindFunc(&request, "FoodId", "ProductId"); err != nil {
		return c.RespondError(err)
	}

	response, err := uc.useCase.AdminCreateProductRecipeGroupHistory(c.Ctx, request)
	if err != nil {
		return c.RespondError(err)
	}

	return c.Respond(map[string]interface{}{
		"data":   response,
		"status": true,
	}, http.StatusOK)
}

func (uc Controller) AdminDeleteProductRecipeGroupHistory(c *web.Context) error {
	id := c.GetParam(reflect.Int, "id").(int)

	if err := c.ValidParam(); err != nil {
		return c.RespondError(err)
	}

	err := uc.useCase.AdminDeleteProductRecipeGroupHistory(c.Ctx, int64(id))
	if err != nil {
		return c.RespondError(err)
	}

	return c.Respond(map[string]interface{}{
		"data":   "ok!",
		"status": true,
	}, http.StatusOK)
}

// @branch

func (uc Controller) BranchGetProductRecipeGroupHistoryList(c *web.Context) error {
	var filter product_recipe_group_history.Filter

	if limit, ok := c.GetQueryFunc(reflect.Int, "limit").(*int); ok {
		filter.Limit = limit
	}
	if offset, ok := c.GetQueryFunc(reflect.Int, "offset").(*int); ok {
		filter.Offset = offset
	}
	if foodId, ok := c.GetQueryFunc(reflect.Int, "food_id").(*int); ok {
		filter.ProductID = foodId
	}

	if err := c.ValidQuery(); err != nil {
		return c.RespondError(err)
	}

	if filter.ProductID == nil {
		return c.RespondError(errors.New("food_id is required"))
	}

	list, count, err := uc.useCase.BranchGetProductRecipeGroupHistoryList(c.Ctx, filter)
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

func (uc Controller) BranchGetProductRecipeGroupHistoryDetail(c *web.Context) error {
	id := c.GetParam(reflect.Int, "id").(int)

	if err := c.ValidParam(); err != nil {
		return c.RespondError(err)
	}

	response, err := uc.useCase.BranchGetProductRecipeGroupHistoryDetail(c.Ctx, int64(id))
	if err != nil {
		return c.RespondError(err)
	}

	return c.Respond(map[string]interface{}{
		"data":   response,
		"status": true,
	}, http.StatusOK)
}

func (uc Controller) BranchCreateProductRecipeGroupHistory(c *web.Context) error {
	var request product_recipe_group_history.BranchCreateRequest

	if err := c.BindFunc(&request, "FoodId", "ProductId"); err != nil {
		return c.RespondError(err)
	}

	response, err := uc.useCase.BranchCreateProductRecipeGroupHistory(c.Ctx, request)
	if err != nil {
		return c.RespondError(err)
	}

	return c.Respond(map[string]interface{}{
		"data":   response,
		"status": true,
	}, http.StatusOK)
}

func (uc Controller) BranchDeleteProductRecipeGroupHistory(c *web.Context) error {
	id := c.GetParam(reflect.Int, "id").(int)

	if err := c.ValidParam(); err != nil {
		return c.RespondError(err)
	}

	err := uc.useCase.BranchDeleteProductRecipeGroupHistory(c.Ctx, int64(id))
	if err != nil {
		return c.RespondError(err)
	}

	return c.Respond(map[string]interface{}{
		"data":   "ok!",
		"status": true,
	}, http.StatusOK)
}
