package product

import (
	"context"
	"github.com/restaurant/internal/pkg/utils"
	"github.com/restaurant/internal/service/product"
	"github.com/restaurant/internal/service/product_recipe"
	"github.com/restaurant/internal/service/product_recipe_group"
	"github.com/restaurant/internal/service/product_recipe_group_history"
)

type UseCase struct {
	product                   Product
	productRecipe             ProductRecipe
	productRecipeGroup        ProductRecipeGroup
	productRecipeGroupHistory ProductRecipeGroupHistory
}

func NewUseCase(product Product, productRecipe ProductRecipe, productRecipeGroup ProductRecipeGroup, productRecipeGroupHistory ProductRecipeGroupHistory) *UseCase {
	return &UseCase{product, productRecipe, productRecipeGroup, productRecipeGroupHistory}
}

// product

// @admin

func (uu UseCase) AdminGetProductList(ctx context.Context, filter product.Filter) ([]product.AdminGetList, int, error) {
	filter.Fields = make(map[string][]string)
	filter.Fields["products"] = []string{"id", "name", "measure_unit_id"}

	filter.Joins = make(map[string]utils.Joins)
	joinColumn := "id"
	mainColumn := "measure_unit_id"
	filter.Joins["measure_unit"] = utils.Joins{JoinColumn: &joinColumn, MainColumn: &mainColumn}

	list, count, err := uu.product.AdminGetList(ctx, filter)
	if err != nil {
		return nil, 0, err
	}

	return list, count, err
}

func (uu UseCase) AdminGetProductDetail(ctx context.Context, id int64) (product.AdminGetDetail, error) {
	var detail product.AdminGetDetail

	data, err := uu.product.AdminGetDetail(ctx, id)
	if err != nil {
		return product.AdminGetDetail{}, err
	}

	detail.ID = data.ID
	detail.Name = data.Name
	detail.MeasureUnitID = data.MeasureUnitID
	detail.Barcode = data.Barcode

	return detail, nil
}

func (uu UseCase) AdminCreateProduct(ctx context.Context, data product.AdminCreateRequest) (product.AdminCreateResponse, error) {
	return uu.product.AdminCreate(ctx, data)
}

func (uu UseCase) AdminUpdateProduct(ctx context.Context, data product.AdminUpdateRequest) error {
	return uu.product.AdminUpdateAll(ctx, data)
}

func (uu UseCase) AdminUpdateProductColumn(ctx context.Context, data product.AdminUpdateRequest) error {
	return uu.product.AdminUpdateColumns(ctx, data)
}

func (uu UseCase) AdminDeleteProduct(ctx context.Context, id int64) error {
	return uu.product.AdminDelete(ctx, id)
}

func (uu UseCase) AdminGetProductSpendingByBranch(ctx context.Context, filter product.SpendingFilter) ([]product.AdminGetSpendingByBranchResponse, error) {
	return uu.product.AdminGetSpendingByBranch(ctx, filter)
}

// @branch

func (uu UseCase) BranchGetProductList(ctx context.Context, filter product.Filter) ([]product.BranchGetList, int, error) {
	filter.Fields = make(map[string][]string)
	filter.Fields["products"] = []string{"id", "name", "measure_unit_id"}

	filter.Joins = make(map[string]utils.Joins)
	joinColumn := "id"
	mainColumn := "measure_unit_id"
	filter.Joins["measure_unit"] = utils.Joins{JoinColumn: &joinColumn, MainColumn: &mainColumn}

	list, count, err := uu.product.BranchGetList(ctx, filter)
	if err != nil {
		return nil, 0, err
	}

	return list, count, err
}

func (uu UseCase) BranchGetProductDetail(ctx context.Context, id int64) (product.BranchGetDetail, error) {
	var detail product.BranchGetDetail

	data, err := uu.product.BranchGetDetail(ctx, id)
	if err != nil {
		return product.BranchGetDetail{}, err
	}

	detail.ID = data.ID
	detail.Name = data.Name
	detail.MeasureUnitID = data.MeasureUnitID
	detail.Barcode = data.Barcode

	return detail, nil
}

func (uu UseCase) BranchCreateProduct(ctx context.Context, data product.BranchCreateRequest) (product.BranchCreateResponse, error) {
	return uu.product.BranchCreate(ctx, data)
}

func (uu UseCase) BranchUpdateProduct(ctx context.Context, data product.BranchUpdateRequest) error {
	return uu.product.BranchUpdateAll(ctx, data)
}

func (uu UseCase) BranchUpdateProductColumn(ctx context.Context, data product.BranchUpdateRequest) error {
	return uu.product.BranchUpdateColumns(ctx, data)
}

func (uu UseCase) BranchDeleteProduct(ctx context.Context, id int64) error {
	return uu.product.BranchDelete(ctx, id)
}

// @cashier

func (uu UseCase) CashierGetProductList(ctx context.Context, filter product.Filter) ([]product.CashierGetList, int, error) {
	filter.Fields = make(map[string][]string)
	filter.Fields["products"] = []string{"id", "name", "measure_unit_id"}

	filter.Joins = make(map[string]utils.Joins)
	joinColumn := "id"
	mainColumn := "measure_unit_id"
	filter.Joins["measure_unit"] = utils.Joins{JoinColumn: &joinColumn, MainColumn: &mainColumn}

	list, count, err := uu.product.CashierGetList(ctx, filter)
	if err != nil {
		return nil, 0, err
	}

	return list, count, err
}

func (uu UseCase) CashierGetProductDetail(ctx context.Context, id int64) (product.CashierGetDetail, error) {
	var detail product.CashierGetDetail

	data, err := uu.product.CashierGetDetail(ctx, id)
	if err != nil {
		return product.CashierGetDetail{}, err
	}

	detail.ID = data.ID
	detail.Name = data.Name
	detail.MeasureUnitID = data.MeasureUnitID

	return detail, nil
}

func (uu UseCase) CashierCreateProduct(ctx context.Context, data product.CashierCreateRequest) (product.CashierCreateResponse, error) {
	return uu.product.CashierCreate(ctx, data)
}

func (uu UseCase) CashierUpdateProduct(ctx context.Context, data product.CashierUpdateRequest) error {
	return uu.product.CashierUpdateAll(ctx, data)
}

func (uu UseCase) CashierUpdateProductColumn(ctx context.Context, data product.CashierUpdateRequest) error {
	return uu.product.CashierUpdateColumns(ctx, data)
}

func (uu UseCase) CashierDeleteProduct(ctx context.Context, id int64) error {
	return uu.product.CashierDelete(ctx, id)
}

func (uu UseCase) CashierGetProductSpending(ctx context.Context, filter product.SpendingFilter) ([]product.CashierGetSpendingResponse, error) {
	return uu.product.CashierGetSpending(ctx, filter)
}

// #product_recipe

// @admin

func (uu UseCase) AdminGetProductRecipeList(ctx context.Context, filter product_recipe.Filter) ([]product_recipe.AdminGetList, int, error) {
	return uu.productRecipe.AdminGetList(ctx, filter)
}

func (uu UseCase) AdminGetProductRecipeDetail(ctx context.Context, restaurantID int64) (*product_recipe.AdminGetDetail, error) {
	return uu.productRecipe.AdminGetDetail(ctx, restaurantID)
}

func (uu UseCase) AdminCreateProductRecipe(ctx context.Context, request product_recipe.AdminCreateRequest) (*product_recipe.AdminCreateResponse, error) {
	return uu.productRecipe.AdminCreate(ctx, request)
}

func (uu UseCase) AdminUpdateProductRecipeColumns(ctx context.Context, request product_recipe.AdminUpdateRequest) error {
	return uu.productRecipe.AdminUpdateColumns(ctx, request)
}

func (uu UseCase) AdminUpdateProductRecipeAll(ctx context.Context, request product_recipe.AdminUpdateRequest) error {
	return uu.productRecipe.AdminUpdateAll(ctx, request)
}

func (uu UseCase) AdminDeleteProductRecipe(ctx context.Context, id int64) error {
	return uu.productRecipe.AdminDelete(ctx, id)
}

// @branch

func (uu UseCase) BranchGetProductRecipeList(ctx context.Context, filter product_recipe.Filter) ([]product_recipe.BranchGetList, int, error) {
	return uu.productRecipe.BranchGetList(ctx, filter)
}

func (uu UseCase) BranchGetProductRecipeDetail(ctx context.Context, restaurantID int64) (*product_recipe.BranchGetDetail, error) {
	return uu.productRecipe.BranchGetDetail(ctx, restaurantID)
}

func (uu UseCase) BranchCreateProductRecipe(ctx context.Context, request product_recipe.BranchCreateRequest) (*product_recipe.BranchCreateResponse, error) {
	return uu.productRecipe.BranchCreate(ctx, request)
}

func (uu UseCase) BranchUpdateProductRecipeColumns(ctx context.Context, request product_recipe.BranchUpdateRequest) error {
	return uu.productRecipe.BranchUpdateColumns(ctx, request)
}

func (uu UseCase) BranchUpdateProductRecipeAll(ctx context.Context, request product_recipe.BranchUpdateRequest) error {
	return uu.productRecipe.BranchUpdateAll(ctx, request)
}

func (uu UseCase) BranchDeleteProductRecipe(ctx context.Context, id int64) error {
	return uu.productRecipe.BranchDelete(ctx, id)
}

// #food_recipe_group

// @admin

func (uu UseCase) AdminGetProductRecipeGroupList(ctx context.Context, filter product_recipe_group.Filter) ([]product_recipe_group.AdminGetListByProductID, int, error) {
	return uu.productRecipeGroup.AdminGetList(ctx, filter, int64(*filter.ProductID))
}

func (uu UseCase) AdminGetProductRecipeGroupDetail(ctx context.Context, restaurantID int64) (*product_recipe_group.AdminGetDetail, error) {
	return uu.productRecipeGroup.AdminGetDetail(ctx, restaurantID)
}

func (uu UseCase) AdminCreateProductRecipeGroup(ctx context.Context, request product_recipe_group.AdminCreateRequest) (*product_recipe_group.AdminCreateResponse, error) {
	return uu.productRecipeGroup.AdminCreate(ctx, request)
}

func (uu UseCase) AdminUpdateProductRecipeGroupColumns(ctx context.Context, request product_recipe_group.AdminUpdateRequest) error {
	return uu.productRecipeGroup.AdminUpdateColumns(ctx, request)
}

func (uu UseCase) AdminUpdateProductRecipeGroupAll(ctx context.Context, request product_recipe_group.AdminUpdateRequest) error {
	return uu.productRecipeGroup.AdminUpdateAll(ctx, request)
}

func (uu UseCase) AdminDeleteProductRecipeGroupSingleRecipe(ctx context.Context, request product_recipe_group.AdminDeleteRecipeRequest) error {
	return uu.productRecipeGroup.AdminDeleteRecipe(ctx, request)
}

func (uu UseCase) AdminDeleteProductRecipeGroup(ctx context.Context, id int64) error {
	return uu.productRecipeGroup.AdminDelete(ctx, id)
}

// @branch

func (uu UseCase) BranchGetProductRecipeGroupList(ctx context.Context, filter product_recipe_group.Filter) ([]product_recipe_group.BranchGetListByProductID, int, error) {
	return uu.productRecipeGroup.BranchGetList(ctx, filter, int64(*filter.ProductID))
}

func (uu UseCase) BranchGetProductRecipeGroupDetail(ctx context.Context, restaurantID int64) (*product_recipe_group.BranchGetDetail, error) {
	return uu.productRecipeGroup.BranchGetDetail(ctx, restaurantID)
}

func (uu UseCase) BranchCreateProductRecipeGroup(ctx context.Context, request product_recipe_group.BranchCreateRequest) (*product_recipe_group.BranchCreateResponse, error) {
	return uu.productRecipeGroup.BranchCreate(ctx, request)
}

func (uu UseCase) BranchUpdateProductRecipeGroupColumns(ctx context.Context, request product_recipe_group.BranchUpdateRequest) error {
	return uu.productRecipeGroup.BranchUpdateColumns(ctx, request)
}

func (uu UseCase) BranchUpdateProductRecipeGroupAll(ctx context.Context, request product_recipe_group.BranchUpdateRequest) error {
	return uu.productRecipeGroup.BranchUpdateAll(ctx, request)
}

func (uu UseCase) BranchDeleteProductRecipeGroupSingleRecipe(ctx context.Context, request product_recipe_group.BranchDeleteRecipeRequest) error {
	return uu.productRecipeGroup.BranchDeleteRecipe(ctx, request)
}

func (uu UseCase) BranchDeleteProductRecipeGroup(ctx context.Context, id int64) error {
	return uu.productRecipeGroup.BranchDelete(ctx, id)
}

// #food_recipe_group_history

// @admin

func (uu UseCase) AdminGetProductRecipeGroupHistoryList(ctx context.Context, filter product_recipe_group_history.Filter) ([]product_recipe_group_history.AdminGetListByProductID, int, error) {
	return uu.productRecipeGroupHistory.AdminGetList(ctx, filter)
}

func (uu UseCase) AdminGetProductRecipeGroupHistoryDetail(ctx context.Context, restaurantID int64) (*product_recipe_group_history.AdminGetDetail, error) {
	return uu.productRecipeGroupHistory.AdminGetDetail(ctx, restaurantID)
}

func (uu UseCase) AdminCreateProductRecipeGroupHistory(ctx context.Context, request product_recipe_group_history.AdminCreateRequest) (*product_recipe_group_history.AdminCreateResponse, error) {
	return uu.productRecipeGroupHistory.AdminCreate(ctx, request)
}

func (uu UseCase) AdminDeleteProductRecipeGroupHistory(ctx context.Context, id int64) error {
	return uu.productRecipeGroupHistory.AdminDelete(ctx, id)
}

// @branch

func (uu UseCase) BranchGetProductRecipeGroupHistoryList(ctx context.Context, filter product_recipe_group_history.Filter) ([]product_recipe_group_history.BranchGetListByProductID, int, error) {
	return uu.productRecipeGroupHistory.BranchGetList(ctx, filter)
}

func (uu UseCase) BranchGetProductRecipeGroupHistoryDetail(ctx context.Context, restaurantID int64) (*product_recipe_group_history.BranchGetDetail, error) {
	return uu.productRecipeGroupHistory.BranchGetDetail(ctx, restaurantID)
}

func (uu UseCase) BranchCreateProductRecipeGroupHistory(ctx context.Context, request product_recipe_group_history.BranchCreateRequest) (*product_recipe_group_history.BranchCreateResponse, error) {
	return uu.productRecipeGroupHistory.BranchCreate(ctx, request)
}

func (uu UseCase) BranchDeleteProductRecipeGroupHistory(ctx context.Context, id int64) error {
	return uu.productRecipeGroupHistory.BranchDelete(ctx, id)
}
