package food

import (
	"context"
	"fmt"
	"github.com/pkg/errors"
	"github.com/restaurant/internal/pkg/file"
	"github.com/restaurant/internal/service/category"
	"github.com/restaurant/internal/service/food"
	"github.com/restaurant/internal/service/foodCategory"
	"github.com/restaurant/internal/service/food_recipe"
	"github.com/restaurant/internal/service/food_recipe_group"
	"github.com/restaurant/internal/service/food_recipe_group_history"
	"github.com/restaurant/internal/service/menu"
	"github.com/restaurant/internal/service/menu_category"
	"strconv"
)

type UseCase struct {
	foodCategory           FoodCategory
	food                   Food
	menu                   Menu
	foodRecipe             FoodRecipe
	basket                 Basket
	category               Category
	menuCategory           MenuCategory
	foodRecipeGroup        FoodRecipeGroup
	foodRecipeGroupHistory FoodRecipeGroupHistory
}

func NewUseCase(
	foodCategory FoodCategory, food Food,
	menu Menu, foodRecipe FoodRecipe,
	basket Basket, category Category, menuCategory MenuCategory,
	foodRecipeGroup FoodRecipeGroup, foodRecipeGroupHistory FoodRecipeGroupHistory,
) *UseCase {
	return &UseCase{
		foodCategory,
		food,
		menu,
		foodRecipe,
		basket,
		category,
		menuCategory,
		foodRecipeGroup,
		foodRecipeGroupHistory,
	}
}

// #food_category

// @admin

func (uu UseCase) AdminGetFoodCategoryList(ctx context.Context, filter foodCategory.Filter) ([]foodCategory.AdminGetList, int, error) {
	filter.Fields = make(map[string][]string)
	filter.Fields["food_category"] = []string{"id", "name", "logo", "main"}

	list, count, err := uu.foodCategory.AdminGetList(ctx, filter)
	if err != nil {
		return nil, 0, err
	}

	return list, count, err
}

func (uu UseCase) AdminGetFoodCategoryDetail(ctx context.Context, id int64) (foodCategory.AdminGetDetail, error) {
	var detail foodCategory.AdminGetDetail

	data, err := uu.foodCategory.AdminGetDetail(ctx, id)
	if err != nil {
		return foodCategory.AdminGetDetail{}, err
	}

	detail.ID = data.ID
	detail.Name = data.Name
	detail.Logo = data.Logo
	detail.Main = data.Main

	return detail, nil
}

func (uu UseCase) AdminCreateFoodCategory(ctx context.Context, data foodCategory.AdminCreateRequest) (foodCategory.AdminCreateResponse, error) {
	if data.Logo != nil {
		imageLink, _, err := file.UploadSingle(data.Logo, "food_category")
		if err != nil {
			return foodCategory.AdminCreateResponse{}, errors.Wrap(err, "logo upload")
		}
		data.LogoLink = &imageLink
	}

	detail, err := uu.foodCategory.AdminCreate(ctx, data)
	if err != nil {
		return foodCategory.AdminCreateResponse{}, err
	}

	return detail, err
}

func (uu UseCase) AdminUpdateFoodCategory(ctx context.Context, data foodCategory.AdminUpdateRequest) error {
	if data.Logo != nil {
		imageLink, _, err := file.UploadSingle(data.Logo, "food_category")
		if err != nil {
			return errors.Wrap(err, "logo upload")
		}
		data.LogoLink = &imageLink
	}

	return uu.foodCategory.AdminUpdateAll(ctx, data)
}

func (uu UseCase) AdminUpdateFoodCategoryColumn(ctx context.Context, data foodCategory.AdminUpdateRequest) error {
	if data.Logo != nil {
		imageLink, _, err := file.UploadSingle(data.Logo, "food_category")
		if err != nil {
			return errors.Wrap(err, "logo upload")
		}
		data.LogoLink = &imageLink
	}

	return uu.foodCategory.AdminUpdateColumns(ctx, data)
}

func (uu UseCase) AdminDeleteFoodCategory(ctx context.Context, id int64) error {
	return uu.foodCategory.AdminDelete(ctx, id)
}

// @branch

func (uu UseCase) BranchGetFoodCategoryList(ctx context.Context, filter foodCategory.Filter) ([]foodCategory.BranchGetList, int, error) {
	filter.Fields = make(map[string][]string)
	filter.Fields["food_category"] = []string{"id", "name", "logo", "main"}

	list, count, err := uu.foodCategory.BranchGetList(ctx, filter)
	if err != nil {
		return nil, 0, err
	}

	return list, count, err
}

func (uu UseCase) BranchGetFoodCategoryDetail(ctx context.Context, id int64) (foodCategory.BranchGetDetail, error) {
	var detail foodCategory.BranchGetDetail

	data, err := uu.foodCategory.BranchGetDetail(ctx, id)
	if err != nil {
		return foodCategory.BranchGetDetail{}, err
	}

	detail.ID = data.ID
	detail.Name = data.Name
	detail.Logo = data.Logo
	detail.Main = data.Main

	return detail, nil
}

func (uu UseCase) BranchCreateFoodCategory(ctx context.Context, data foodCategory.BranchCreateRequest) (foodCategory.BranchCreateResponse, error) {
	if data.Logo != nil {
		imageLink, _, err := file.UploadSingle(data.Logo, "food_category")
		if err != nil {
			return foodCategory.BranchCreateResponse{}, errors.Wrap(err, "logo upload")
		}
		data.LogoLink = &imageLink
	}

	detail, err := uu.foodCategory.BranchCreate(ctx, data)
	if err != nil {
		return foodCategory.BranchCreateResponse{}, err
	}

	return detail, err
}

func (uu UseCase) BranchUpdateFoodCategory(ctx context.Context, data foodCategory.BranchUpdateRequest) error {
	if data.Logo != nil {
		imageLink, _, err := file.UploadSingle(data.Logo, "food_category")
		if err != nil {
			return errors.Wrap(err, "logo upload")
		}
		data.LogoLink = &imageLink
	}

	return uu.foodCategory.BranchUpdateAll(ctx, data)
}

func (uu UseCase) BranchUpdateFoodCategoryColumn(ctx context.Context, data foodCategory.BranchUpdateRequest) error {
	if data.Logo != nil {
		imageLink, _, err := file.UploadSingle(data.Logo, "food_category")
		if err != nil {
			return errors.Wrap(err, "logo upload")
		}
		data.LogoLink = &imageLink
	}

	return uu.foodCategory.BranchUpdateColumns(ctx, data)
}

func (uu UseCase) BranchDeleteFoodCategory(ctx context.Context, id int64) error {
	return uu.foodCategory.BranchDelete(ctx, id)
}

// @cashier

func (uu UseCase) CashierGetFoodCategoryList(ctx context.Context, filter foodCategory.Filter) ([]foodCategory.CashierGetList, int, error) {
	filter.Fields = make(map[string][]string)
	filter.Fields["food_category"] = []string{"id", "name", "logo", "main"}

	list, count, err := uu.foodCategory.CashierGetList(ctx, filter)
	if err != nil {
		return nil, 0, err
	}

	return list, count, err
}

func (uu UseCase) CashierGetFoodCategoryDetail(ctx context.Context, id int64) (foodCategory.CashierGetDetail, error) {
	var detail foodCategory.CashierGetDetail

	data, err := uu.foodCategory.CashierGetDetail(ctx, id)
	if err != nil {
		return foodCategory.CashierGetDetail{}, err
	}

	detail.ID = data.ID
	detail.Name = data.Name
	detail.Logo = data.Logo
	detail.Main = data.Main

	return detail, nil
}

func (uu UseCase) CashierCreateFoodCategory(ctx context.Context, data foodCategory.CashierCreateRequest) (foodCategory.CashierCreateResponse, error) {
	if data.Logo != nil {
		imageLink, _, err := file.UploadSingle(data.Logo, "food_category")
		if err != nil {
			return foodCategory.CashierCreateResponse{}, errors.Wrap(err, "logo upload")
		}
		data.LogoLink = &imageLink
	}

	detail, err := uu.foodCategory.CashierCreate(ctx, data)
	if err != nil {
		return foodCategory.CashierCreateResponse{}, err
	}

	return detail, nil
}

func (uu UseCase) CashierUpdateFoodCategory(ctx context.Context, data foodCategory.CashierUpdateRequest) error {
	if data.Logo != nil {
		imageLink, _, err := file.UploadSingle(data.Logo, "food_category")
		if err != nil {
			return errors.Wrap(err, "logo upload")
		}
		data.LogoLink = &imageLink
	}

	return uu.foodCategory.CashierUpdateAll(ctx, data)
}

func (uu UseCase) CashierUpdateFoodCategoryColumn(ctx context.Context, data foodCategory.CashierUpdateRequest) error {
	if data.Logo != nil {
		imageLink, _, err := file.UploadSingle(data.Logo, "food_category")
		if err != nil {
			return errors.Wrap(err, "logo upload")
		}
		data.LogoLink = &imageLink
	}

	return uu.foodCategory.CashierUpdateColumns(ctx, data)
}

func (uu UseCase) CashierDeleteFoodCategory(ctx context.Context, id int64) error {
	return uu.foodCategory.CashierDelete(ctx, id)
}

// @client

func (uu UseCase) ClientGetFoodCategoryList(ctx context.Context, filter foodCategory.Filter) ([]foodCategory.ClientGetList, int, error) {
	fields := make(map[string][]string)

	fields["food_category"] = []string{"id", "name", "logo"}

	filter.Fields = fields
	list, count, err := uu.foodCategory.ClientGetList(ctx, filter)
	if err != nil {
		return nil, 0, err
	}

	return list, count, err
}

// @waiter

func (uu UseCase) WaiterGetFoodCategoryList(ctx context.Context) ([]foodCategory.WaiterGetList, error) {
	return uu.foodCategory.WaiterGetList(ctx)
}

//----------------------------------------------------------------------------------------------------------

// #food

// @admin

func (uu UseCase) AdminGetFoodList(ctx context.Context, filter food.Filter) ([]food.AdminGetList, int, error) {
	foods := make(map[string][]string)
	foods["foods"] = []string{"id", "name", "photos", "price", "category_id"}
	filter.Fields = foods

	list, count, err := uu.food.AdminGetList(ctx, filter)
	if err != nil {
		return nil, 0, err
	}

	return list, count, err
}

func (uu UseCase) AdminGetFoodDetail(ctx context.Context, id int64) (food.AdminGetDetail, error) {
	var detail food.AdminGetDetail

	data, err := uu.food.AdminGetDetail(ctx, id)
	if err != nil {
		return food.AdminGetDetail{}, err
	}

	detail.ID = data.ID
	detail.Name = data.Name
	detail.Photos = data.Photos
	detail.Price = data.Price
	detail.CategoryID = data.CategoryID

	return detail, nil
}

func (uu UseCase) AdminCreateFood(ctx context.Context, data food.AdminCreateRequest) (food.AdminCreateResponse, error) {
	if data.Photos != nil {
		imageLinks, _, err := file.UploadMultiple(data.Photos, "food")
		if err != nil {
			return food.AdminCreateResponse{}, errors.Wrap(err, "logo upload")
		}
		data.PhotosLink = &imageLinks
	}

	detail, err := uu.food.AdminCreate(ctx, data)
	if err != nil {
		return food.AdminCreateResponse{}, err
	}

	return detail, err
}

func (uu UseCase) AdminUpdateFood(ctx context.Context, data food.AdminUpdateRequest) error {
	if data.Photos != nil {
		imageLinks, _, err := file.UploadMultiple(data.Photos, "food")
		if err != nil {
			return errors.Wrap(err, "food upload")
		}
		data.PhotosLink = &imageLinks
	}
	return uu.food.AdminUpdateAll(ctx, data)
}

func (uu UseCase) AdminUpdateFoodColumn(ctx context.Context, data food.AdminUpdateRequest) error {
	if data.Photos != nil {
		imageLinks, _, err := file.UploadMultiple(data.Photos, "food")
		if err != nil {
			return errors.Wrap(err, "food upload")
		}
		data.PhotosLink = &imageLinks
	}
	return uu.food.AdminUpdateColumns(ctx, data)
}

func (uu UseCase) AdminDeleteFood(ctx context.Context, id int64) error {
	return uu.food.AdminDelete(ctx, id)
}

func (uu UseCase) AdminDeleteImage(ctx context.Context, request food.AdminDeleteImageRequest) error {
	return uu.food.AdminDeleteImage(ctx, request)
}

// @branch

func (uu UseCase) BranchGetFoodList(ctx context.Context, filter food.Filter) ([]food.BranchGetList, int, error) {
	foods := make(map[string][]string)
	foods["foods"] = []string{"id", "name", "photos", "price", "category_id"}
	filter.Fields = foods

	list, count, err := uu.food.BranchGetList(ctx, filter)
	if err != nil {
		return nil, 0, err
	}

	return list, count, err
}

func (uu UseCase) BranchGetFoodDetail(ctx context.Context, id int64) (food.BranchGetDetail, error) {
	var detail food.BranchGetDetail

	data, err := uu.food.BranchGetDetail(ctx, id)
	if err != nil {
		return food.BranchGetDetail{}, err
	}

	detail.ID = data.ID
	detail.Name = data.Name
	detail.Photos = data.Photos
	detail.Price = data.Price
	detail.CategoryID = data.CategoryID

	return detail, nil
}

func (uu UseCase) BranchCreateFood(ctx context.Context, data food.BranchCreateRequest) (food.BranchCreateResponse, error) {
	if data.Photos != nil {
		imageLinks, _, err := file.UploadMultiple(data.Photos, "food")
		if err != nil {
			return food.BranchCreateResponse{}, errors.Wrap(err, "logo upload")
		}
		data.PhotosLink = &imageLinks
	}

	detail, err := uu.food.BranchCreate(ctx, data)
	if err != nil {
		return food.BranchCreateResponse{}, err
	}

	return detail, err
}

func (uu UseCase) BranchUpdateFood(ctx context.Context, data food.BranchUpdateRequest) error {
	if data.Photos != nil {
		imageLinks, _, err := file.UploadMultiple(data.Photos, "food")
		if err != nil {
			return errors.Wrap(err, "food upload")
		}
		data.PhotosLink = &imageLinks
	}

	return uu.food.BranchUpdateAll(ctx, data)
}

func (uu UseCase) BranchUpdateFoodColumn(ctx context.Context, data food.BranchUpdateRequest) error {
	if data.Photos != nil {
		imageLinks, _, err := file.UploadMultiple(data.Photos, "food")
		if err != nil {
			return errors.Wrap(err, "food upload")
		}
		data.PhotosLink = &imageLinks
	}
	return uu.food.BranchUpdateColumns(ctx, data)
}

func (uu UseCase) BranchDeleteFood(ctx context.Context, id int64) error {
	return uu.food.BranchDelete(ctx, id)
}

func (uu UseCase) BranchDeleteImage(ctx context.Context, request food.AdminDeleteImageRequest) error {
	return uu.food.BranchDeleteImage(ctx, request)
}

// @cashier

func (uu UseCase) CashierGetFoodList(ctx context.Context, filter food.Filter) ([]food.CashierGetList, int, error) {
	foods := make(map[string][]string)
	foods["foods"] = []string{"id", "name", "photos", "price", "category_id"}
	filter.Fields = foods

	list, count, err := uu.food.CashierGetList(ctx, filter)
	if err != nil {
		return nil, 0, err
	}

	return list, count, err
}

// #menu

// @client

func (uu UseCase) ClientGetMenuList(ctx context.Context, filter menu.Filter) ([]menu.ClientGetList, error) {
	list, err := uu.menu.ClientGetList(ctx, filter)
	if err != nil {
		return nil, err
	}

	return list, err
}

func (uu UseCase) ClientGetMenuListByCategoryID(ctx context.Context, foodCategoryId int, filter menu.Filter) ([]menu.ClientGetListByCategoryID, error) {

	list, err := uu.menu.ClientGetListByCategoryID(ctx, foodCategoryId, filter)
	if err != nil {
		return nil, err
	}

	for k, v := range list {
		var unit string
		if v.Distance != nil {
			dist, err := strconv.ParseFloat(*v.Distance, 64)
			if err != nil {
				return nil, err
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

	return list, err
}

// @admin

func (uu UseCase) AdminGetMenuList(ctx context.Context, filter menu.Filter) ([]menu.AdminGetList, int, error) {
	list, count, err := uu.menu.AdminGetList(ctx, filter)
	if err != nil {
		return nil, 0, err
	}

	return list, count, err
}

func (uu UseCase) AdminGetMenuDetail(ctx context.Context, id int64) (menu.AdminGetDetail, error) {
	var detail menu.AdminGetDetail

	data, err := uu.menu.AdminGetDetail(ctx, id)
	if err != nil {
		return menu.AdminGetDetail{}, err
	}

	detail.ID = data.ID
	detail.Name = data.Name
	detail.Photos = data.Photos
	detail.FoodID = data.FoodID
	detail.BranchID = data.BranchID
	detail.Status = data.Status
	detail.NewPrice = data.NewPrice
	detail.OldPrice = data.OldPrice
	detail.Description = data.Description
	detail.MenuCategoryID = data.MenuCategoryID
	detail.CategoryID = data.CategoryID

	return detail, nil
}

func (uu UseCase) AdminCreateMenu(ctx context.Context, data menu.AdminCreateRequest) ([]menu.AdminCreateResponse, error) {
	link, photos, err := file.UploadMultiple(data.Photos, "menus")
	if err != nil {
		return nil, err
	}
	data.PhotosLink = &link
	detail, err := uu.menu.AdminCreate(ctx, data)
	if err != nil {
		if err := file.DeleteFiles(photos...); err != nil {
			return nil, err
		}
		return nil, err
	}

	return detail, err
}

func (uu UseCase) AdminUpdateMenu(ctx context.Context, data menu.AdminUpdateRequest) error {
	link, photos, err := file.UploadMultiple(data.Photos, "menus")
	if err != nil {
		return err
	}
	data.PhotosLink = &link
	if err = uu.menu.AdminUpdateAll(ctx, data); err != nil {
		if err := file.DeleteFiles(photos...); err != nil {
			return err
		}
		return err
	}

	return nil
}

func (uu UseCase) AdminUpdateMenuColumn(ctx context.Context, data menu.AdminUpdateRequest) error {
	link, photos, err := file.UploadMultiple(data.Photos, "menus")
	if err != nil {
		return err
	}
	data.PhotosLink = &link
	if err = uu.menu.AdminUpdateColumns(ctx, data); err != nil {
		if err := file.DeleteFiles(photos...); err != nil {
			return err
		}
		return err
	}

	return nil
}

func (uu UseCase) AdminDeleteMenu(ctx context.Context, id int64) error {
	return uu.menu.AdminDelete(ctx, id)
}

func (uu UseCase) AdminRemoveMenuPhoto(ctx context.Context, id int64, index *int) error {
	photo, err := uu.menu.AdminRemovePhoto(ctx, id, index)
	if err != nil {
		return err
	}

	if photo != nil {
		_ = file.DeleteFiles(*photo)
	}

	return nil
}

// #menu

// @branch

func (uu UseCase) BranchGetMenuList(ctx context.Context, filter menu.Filter) ([]menu.BranchGetList, int, error) {
	list, count, err := uu.menu.BranchGetList(ctx, filter)
	if err != nil {
		return nil, 0, err
	}

	return list, count, err
}

func (uu UseCase) BranchGetMenuDetail(ctx context.Context, id int64) (menu.BranchGetDetail, error) {
	var detail menu.BranchGetDetail

	data, err := uu.menu.BranchGetDetail(ctx, id)
	if err != nil {
		return menu.BranchGetDetail{}, err
	}

	detail.ID = data.ID
	detail.Name = data.Name
	detail.Photos = data.Photos
	detail.FoodID = data.FoodID
	detail.BranchID = data.BranchID
	detail.Status = data.Status
	detail.OldPrice = data.OldPrice
	detail.NewPrice = data.NewPrice
	detail.Description = data.Description
	detail.MenuCategoryID = data.MenuCategoryID
	detail.CategoryID = data.CategoryID

	return detail, nil
}

func (uu UseCase) BranchCreateMenu(ctx context.Context, data menu.BranchCreateRequest) (menu.BranchCreateResponse, error) {
	link, photos, err := file.UploadMultiple(data.Photos, "menus")
	if err != nil {
		return menu.BranchCreateResponse{}, err
	}
	data.PhotosLink = &link
	detail, err := uu.menu.BranchCreate(ctx, data)
	if err != nil {
		if err := file.DeleteFiles(photos...); err != nil {
			return menu.BranchCreateResponse{}, err
		}
		return menu.BranchCreateResponse{}, err
	}

	return detail, err
}

func (uu UseCase) BranchUpdateMenu(ctx context.Context, data menu.BranchUpdateRequest) error {
	link, photos, err := file.UploadMultiple(data.Photos, "menus")
	if err != nil {
		return err
	}
	data.PhotosLink = &link
	if err = uu.menu.BranchUpdateAll(ctx, data); err != nil {
		if err := file.DeleteFiles(photos...); err != nil {
			return err
		}
		return err
	}

	return nil
}

func (uu UseCase) BranchUpdateMenuColumn(ctx context.Context, data menu.BranchUpdateRequest) error {
	link, photos, err := file.UploadMultiple(data.Photos, "menus")
	if err != nil {
		return err
	}
	data.PhotosLink = &link
	if err = uu.menu.BranchUpdateColumns(ctx, data); err != nil {
		if err := file.DeleteFiles(photos...); err != nil {
			return err
		}
		return err
	}

	return nil
}

func (uu UseCase) BranchDeleteMenu(ctx context.Context, id int64) error {
	return uu.menu.BranchDelete(ctx, id)
}

func (uu UseCase) BranchUpdatePrinterID(ctx context.Context, request menu.BranchUpdatePrinterIDRequest) error {
	return uu.menu.BranchUpdatePrinterID(ctx, request)
}

func (uu UseCase) BranchDeletePrinterID(ctx context.Context, menuID int64) error {
	return uu.menu.BranchDeletePrinterID(ctx, menuID)
}

func (uu UseCase) BranchRemoveMenuPhoto(ctx context.Context, id int64, index int) error {
	photo, err := uu.menu.BranchRemovePhoto(ctx, id, index)
	if err != nil {
		return err
	}

	if photo != nil {
		_ = file.DeleteFiles(*photo)
	}

	return nil
}

// @cashier

func (uu UseCase) CashierUpdateColumns(ctx context.Context, request menu.CashierUpdateMenuStatus) error {
	return uu.menu.CashierUpdateColumns(ctx, request)
}

func (uu UseCase) CashierGetMenuList(ctx context.Context, filter menu.Filter) ([]menu.CashierGetList, int, error) {
	list, count, err := uu.menu.CashierGetList(ctx, filter)
	if err != nil {
		return nil, 0, err
	}

	return list, count, err
}

func (uu UseCase) CashierGetMenuDetail(ctx context.Context, id int64) (menu.CashierGetDetail, error) {
	var detail menu.CashierGetDetail

	data, err := uu.menu.CashierGetDetail(ctx, id)
	if err != nil {
		return menu.CashierGetDetail{}, err
	}

	detail.ID = data.ID
	detail.Name = data.Name
	detail.Photos = data.Photos
	detail.FoodID = data.FoodID
	detail.BranchID = data.BranchID
	detail.Status = data.Status
	detail.OldPrice = data.OldPrice
	detail.NewPrice = data.NewPrice
	detail.Description = data.Description

	return detail, nil
}

func (uu UseCase) CashierCreateMenu(ctx context.Context, data menu.CashierCreateRequest) (menu.CashierCreateResponse, error) {
	link, photos, err := file.UploadMultiple(data.Photos, "menus")
	if err != nil {
		return menu.CashierCreateResponse{}, err
	}
	data.PhotosLink = &link
	detail, err := uu.menu.CashierCreate(ctx, data)
	if err != nil {
		if err := file.DeleteFiles(photos...); err != nil {
			return menu.CashierCreateResponse{}, err
		}
		return menu.CashierCreateResponse{}, err
	}

	return detail, err
}

func (uu UseCase) CashierUpdateMenu(ctx context.Context, data menu.CashierUpdateRequest) error {
	link, photos, err := file.UploadMultiple(data.Photos, "menus")
	if err != nil {
		return err
	}
	data.PhotosLink = &link
	if err = uu.menu.CashierUpdateAll(ctx, data); err != nil {
		if err := file.DeleteFiles(photos...); err != nil {
			return err
		}
		return err
	}

	return nil
}

func (uu UseCase) CashierUpdateMenuColumn(ctx context.Context, data menu.CashierUpdateRequest) error {
	link, photos, err := file.UploadMultiple(data.Photos, "menus")
	if err != nil {
		return err
	}
	data.PhotosLink = &link
	if err = uu.menu.CashierUpdateColumn(ctx, data); err != nil {
		if err = file.DeleteFiles(photos...); err != nil {
			return err
		}
		return err
	}

	return nil
}

func (uu UseCase) CashierDeleteMenu(ctx context.Context, id int64) error {
	return uu.menu.CashierDelete(ctx, id)
}

func (uu UseCase) CashierUpdatePrinterID(ctx context.Context, request menu.CashierUpdatePrinterIDRequest) error {
	return uu.menu.CashierUpdatePrinterID(ctx, request)
}

func (uu UseCase) CashierDeletePrinterID(ctx context.Context, menuID int64) error {
	return uu.menu.CashierDeletePrinterID(ctx, menuID)
}

func (uu UseCase) CashierRemoveMenuPhoto(ctx context.Context, id int64, index int) error {
	photo, err := uu.menu.CashierRemovePhoto(ctx, id, index)
	if err != nil {
		return err
	}

	if photo != nil {
		_ = file.DeleteFiles(*photo)
	}

	return nil
}

// @waiter

func (uu UseCase) WaiterGetMenuList(ctx context.Context, filter menu.Filter) ([]menu.WaiterGetMenuListResponse, error) {
	return uu.menu.WaiterGetMenuList(ctx, filter)
}

// #food_recipe

// @admin

func (uu UseCase) AdminGetFoodRecipeList(ctx context.Context, filter food_recipe.Filter) ([]food_recipe.AdminGetList, int, error) {
	return uu.foodRecipe.AdminGetList(ctx, filter)
}

func (uu UseCase) AdminGetFoodRecipeDetail(ctx context.Context, restaurantID int64) (*food_recipe.AdminGetDetail, error) {
	return uu.foodRecipe.AdminGetDetail(ctx, restaurantID)
}

func (uu UseCase) AdminCreateFoodRecipe(ctx context.Context, request food_recipe.AdminCreateRequest) (*food_recipe.AdminCreateResponse, error) {
	return uu.foodRecipe.AdminCreate(ctx, request)
}

func (uu UseCase) AdminUpdateFoodRecipeColumns(ctx context.Context, request food_recipe.AdminUpdateRequest) error {
	return uu.foodRecipe.AdminUpdateColumns(ctx, request)
}

func (uu UseCase) AdminUpdateFoodRecipeAll(ctx context.Context, request food_recipe.AdminUpdateRequest) error {
	return uu.foodRecipe.AdminUpdateAll(ctx, request)
}

func (uu UseCase) AdminDeleteFoodRecipe(ctx context.Context, id int64) error {
	return uu.foodRecipe.AdminDelete(ctx, id)
}

// @branch

func (uu UseCase) BranchGetFoodRecipeList(ctx context.Context, filter food_recipe.Filter) ([]food_recipe.BranchGetList, int, error) {
	return uu.foodRecipe.BranchGetList(ctx, filter)
}

func (uu UseCase) BranchGetFoodRecipeDetail(ctx context.Context, restaurantID int64) (*food_recipe.BranchGetDetail, error) {
	return uu.foodRecipe.BranchGetDetail(ctx, restaurantID)
}

func (uu UseCase) BranchCreateFoodRecipe(ctx context.Context, request food_recipe.BranchCreateRequest) (*food_recipe.BranchCreateResponse, error) {
	return uu.foodRecipe.BranchCreate(ctx, request)
}

func (uu UseCase) BranchUpdateFoodRecipeColumns(ctx context.Context, request food_recipe.BranchUpdateRequest) error {
	return uu.foodRecipe.BranchUpdateColumns(ctx, request)
}

func (uu UseCase) BranchUpdateFoodRecipeAll(ctx context.Context, request food_recipe.BranchUpdateRequest) error {
	return uu.foodRecipe.BranchUpdateAll(ctx, request)
}

func (uu UseCase) BranchDeleteFoodRecipe(ctx context.Context, id int64) error {
	return uu.foodRecipe.BranchDelete(ctx, id)
}

// #category

// @super-admin

func (uu UseCase) SuperAdminGetCategoryList(ctx context.Context, filter category.Filter) ([]category.SuperAdminGetList, int, error) {
	filter.Fields = make(map[string][]string)
	filter.Fields["categories"] = []string{"id", "name", "logo", "status"}

	list, count, err := uu.category.SuperAdminGetList(ctx, filter)
	if err != nil {
		return nil, 0, err
	}

	return list, count, err
}

func (uu UseCase) SuperAdminGetCategoryDetail(ctx context.Context, id int64) (category.SuperAdminGetDetail, error) {
	var detail category.SuperAdminGetDetail

	data, err := uu.category.SuperAdminGetDetail(ctx, id)
	if err != nil {
		return category.SuperAdminGetDetail{}, err
	}

	detail.ID = data.ID
	detail.Name = data.Name
	detail.Logo = data.Logo
	detail.Status = data.Status

	return detail, nil
}

func (uu UseCase) SuperAdminCreateCategory(ctx context.Context, data category.SuperAdminCreateRequest) (category.SuperAdminCreateResponse, error) {
	if data.Logo != nil {
		imageLink, _, err := file.UploadSingle(data.Logo, "menu_categories")
		if err != nil {
			return category.SuperAdminCreateResponse{}, errors.Wrap(err, "logo upload")
		}
		data.LogoLink = &imageLink
	}

	detail, err := uu.category.SuperAdminCreate(ctx, data)
	if err != nil {
		return category.SuperAdminCreateResponse{}, err
	}

	return detail, err
}

func (uu UseCase) SuperAdminUpdateCategory(ctx context.Context, data category.SuperAdminUpdateRequest) error {
	if data.Logo != nil {
		imageLink, _, err := file.UploadSingle(data.Logo, "menu_categories")
		if err != nil {
			return errors.Wrap(err, "logo upload")
		}
		data.LogoLink = &imageLink
	}

	return uu.category.SuperAdminUpdateAll(ctx, data)
}

func (uu UseCase) SuperAdminUpdateCategoryColumn(ctx context.Context, data category.SuperAdminUpdateRequest) error {
	if data.Logo != nil {
		imageLink, _, err := file.UploadSingle(data.Logo, "menu_categories")
		if err != nil {
			return errors.Wrap(err, "logo upload")
		}
		data.LogoLink = &imageLink
	}

	return uu.category.SuperAdminUpdateColumns(ctx, data)
}

func (uu UseCase) SuperAdminDeleteCategory(ctx context.Context, id int64) error {
	return uu.category.SuperAdminDelete(ctx, id)
}

// @client

func (uu UseCase) ClientGetCategoryList(ctx context.Context, filter category.Filter) ([]category.ClientGetList, int, error) {
	list, count, err := uu.category.ClientGetList(ctx, filter)
	if err != nil {
		return nil, 0, err
	}

	return list, count, err
}

// @branch

func (uu UseCase) BranchGetCategoryList(ctx context.Context, filter category.Filter) ([]category.BranchGetList, int, error) {
	filter.Fields = make(map[string][]string)
	filter.Fields["categories"] = []string{"id", "name", "logo"}

	list, count, err := uu.category.BranchGetList(ctx, filter)
	if err != nil {
		return nil, 0, err
	}

	return list, count, err
}

// @cashier

func (uu UseCase) CashierGetCategoryList(ctx context.Context, filter category.Filter) ([]category.CashierGetList, int, error) {
	filter.Fields = make(map[string][]string)
	filter.Fields["categories"] = []string{"id", "name", "logo"}

	list, count, err := uu.category.CashierGetList(ctx, filter)
	if err != nil {
		return nil, 0, err
	}

	return list, count, err
}

// @admin

func (uu UseCase) AdminGetCategoryList(ctx context.Context, filter category.Filter) ([]category.AdminGetList, int, error) {
	filter.Fields = make(map[string][]string)
	filter.Fields["categories"] = []string{"id", "name", "logo"}

	list, count, err := uu.category.AdminGetList(ctx, filter)
	if err != nil {
		return nil, 0, err
	}

	return list, count, err
}

// @waiter

func (uu UseCase) WaiterGetCategoryList(ctx context.Context) ([]category.WaiterGetList, error) {
	return uu.category.WaiterGetList(ctx)
}

// #menu_category

// @admin

func (uu UseCase) AdminCreateMenuCategory(ctx context.Context, request menu_category.AdminCreateRequest) (*menu_category.AdminCreateResponse, error) {
	logo, _, err := file.UploadSingle(request.Logo, "menu_category")
	if err != nil {
		return nil, err
	}
	request.LogoLink = &logo

	data, err := uu.menuCategory.AdminCreate(ctx, request)
	if err != nil {
		return nil, err
	}

	createdAt := data.CreatedAt.String()

	response := menu_category.AdminCreateResponse{
		ID:        data.ID,
		Name:      data.Name,
		CreatedAt: &createdAt,
	}

	return &response, nil
}

func (uu UseCase) AdminUpdateMenuCategory(ctx context.Context, request menu_category.AdminUpdateRequest) error {
	logo, _, err := file.UploadSingle(request.Logo, "menu_category")
	if err != nil {
		return err
	}
	request.LogoLink = &logo

	return uu.menuCategory.AdminUpdate(ctx, request)
}

func (uu UseCase) AdminGetListMenuCategory(ctx context.Context, filter menu_category.Filter) ([]menu_category.AdminGetListResponse, int, error) {
	list, count, err := uu.menuCategory.AdminGetList(ctx, filter)
	if err != nil {
		return nil, 0, err
	}

	response := make([]menu_category.AdminGetListResponse, 0)

	var detail menu_category.AdminGetListResponse
	for i := range list {
		detail.ID = list[i].ID
		detail.Name = list[i].Name
		detail.Logo = list[i].Logo

		response = append(response, detail)
	}

	return response, count, nil
}

func (uu UseCase) AdminGetDetailMenuCategory(ctx context.Context, id int64) (*menu_category.AdminGetDetailResponse, error) {
	data, err := uu.menuCategory.AdminGetDetail(ctx, id)
	if err != nil {
		return nil, err
	}

	createdAt := data.CreatedAt.String()

	response := menu_category.AdminGetDetailResponse{
		ID:        data.ID,
		Name:      data.Name,
		Logo:      data.Logo,
		CreatedAt: &createdAt,
	}

	return &response, nil
}

func (uu UseCase) AdminDeleteMenuCategory(ctx context.Context, id int64) error {
	return uu.menuCategory.AdminDelete(ctx, id)
}

// @branch

func (uu UseCase) BranchCreateMenuCategory(ctx context.Context, request menu_category.BranchCreateRequest) (*menu_category.BranchCreateResponse, error) {
	if request.Logo != nil {
		logo, _, err := file.UploadSingle(request.Logo, "menu_category")
		if err != nil {
			return nil, err
		}
		request.LogoLink = &logo
	}

	data, err := uu.menuCategory.BranchCreate(ctx, request)
	if err != nil {
		return nil, err
	}

	createdAt := data.CreatedAt.String()

	response := menu_category.BranchCreateResponse{
		ID:        data.ID,
		Name:      data.Name,
		Logo:      data.Logo,
		CreatedAt: &createdAt,
	}

	return &response, nil
}

func (uu UseCase) BranchUpdateMenuCategory(ctx context.Context, request menu_category.BranchUpdateRequest) error {
	logo, _, err := file.UploadSingle(request.Logo, "menu_category")
	if err != nil {
		return err
	}
	request.LogoLink = &logo

	return uu.menuCategory.BranchUpdate(ctx, request)
}

func (uu UseCase) BranchGetListMenuCategory(ctx context.Context, filter menu_category.Filter) ([]menu_category.BranchGetListResponse, int, error) {
	list, count, err := uu.menuCategory.BranchGetList(ctx, filter)
	if err != nil {
		return nil, 0, err
	}

	response := make([]menu_category.BranchGetListResponse, 0)

	var detail menu_category.BranchGetListResponse
	for i := range list {
		detail.ID = list[i].ID
		detail.Name = list[i].Name
		detail.Logo = list[i].Logo

		response = append(response, detail)
	}

	return response, count, nil
}

func (uu UseCase) BranchGetDetailMenuCategory(ctx context.Context, id int64) (*menu_category.BranchGetDetailResponse, error) {
	data, err := uu.menuCategory.BranchGetDetail(ctx, id)
	if err != nil {
		return nil, err
	}

	createdAt := data.CreatedAt.String()

	response := menu_category.BranchGetDetailResponse{
		ID:        data.ID,
		Name:      data.Name,
		Logo:      data.Logo,
		CreatedAt: &createdAt,
	}

	return &response, nil
}

func (uu UseCase) BranchDeleteMenuCategory(ctx context.Context, id int64) error {
	return uu.menuCategory.BranchDelete(ctx, id)
}

// @cashier

func (uu UseCase) CashierCreateMenuCategory(ctx context.Context, request menu_category.CashierCreateRequest) (*menu_category.CashierCreateResponse, error) {
	logo, _, err := file.UploadSingle(request.Logo, "menu_category")
	if err != nil {
		return nil, err
	}
	request.LogoLink = &logo

	data, err := uu.menuCategory.CashierCreate(ctx, request)
	if err != nil {
		return nil, err
	}

	createdAt := data.CreatedAt.String()

	response := menu_category.CashierCreateResponse{
		ID:        data.ID,
		Name:      data.Name,
		CreatedAt: &createdAt,
	}

	return &response, nil
}

func (uu UseCase) CashierUpdateMenuCategory(ctx context.Context, request menu_category.CashierUpdateRequest) error {
	logo, _, err := file.UploadSingle(request.Logo, "menu_category")
	if err != nil {
		return err
	}
	request.LogoLink = &logo

	return uu.menuCategory.CashierUpdate(ctx, request)
}

func (uu UseCase) CashierGetListMenuCategory(ctx context.Context, filter menu_category.Filter) ([]menu_category.CashierGetListResponse, int, error) {
	list, count, err := uu.menuCategory.CashierGetList(ctx, filter)
	if err != nil {
		return nil, 0, err
	}

	response := make([]menu_category.CashierGetListResponse, 0)

	var detail menu_category.CashierGetListResponse
	for i := range list {
		detail.ID = list[i].ID
		detail.Name = list[i].Name
		detail.Logo = list[i].Logo

		response = append(response, detail)
	}

	return response, count, nil
}

func (uu UseCase) CashierGetDetailMenuCategory(ctx context.Context, id int64) (*menu_category.CashierGetDetailResponse, error) {
	data, err := uu.menuCategory.CashierGetDetail(ctx, id)
	if err != nil {
		return nil, err
	}

	createdAt := data.CreatedAt.String()

	response := menu_category.CashierGetDetailResponse{
		ID:        data.ID,
		Name:      data.Name,
		Logo:      data.Logo,
		CreatedAt: &createdAt,
	}

	return &response, nil
}

func (uu UseCase) CashierDeleteMenuCategory(ctx context.Context, id int64) error {
	return uu.menuCategory.CashierDelete(ctx, id)
}

// #food_recipe_group

// @admin

func (uu UseCase) AdminGetFoodRecipeGroupList(ctx context.Context, filter food_recipe_group.Filter) ([]food_recipe_group.AdminGetListByFoodID, int, error) {
	return uu.foodRecipeGroup.AdminGetList(ctx, filter, int64(*filter.FoodId))
}

func (uu UseCase) AdminGetFoodRecipeGroupDetail(ctx context.Context, restaurantID int64) (*food_recipe_group.AdminGetDetail, error) {
	return uu.foodRecipeGroup.AdminGetDetail(ctx, restaurantID)
}

func (uu UseCase) AdminCreateFoodRecipeGroup(ctx context.Context, request food_recipe_group.AdminCreateRequest) (*food_recipe_group.AdminCreateResponse, error) {
	return uu.foodRecipeGroup.AdminCreate(ctx, request)
}

func (uu UseCase) AdminUpdateFoodRecipeGroupColumns(ctx context.Context, request food_recipe_group.AdminUpdateRequest) error {
	return uu.foodRecipeGroup.AdminUpdateColumns(ctx, request)
}

func (uu UseCase) AdminUpdateFoodRecipeGroupAll(ctx context.Context, request food_recipe_group.AdminUpdateRequest) error {
	return uu.foodRecipeGroup.AdminUpdateAll(ctx, request)
}

func (uu UseCase) AdminDeleteFoodRecipeGroupSingleRecipe(ctx context.Context, request food_recipe_group.AdminDeleteRecipeRequest) error {
	return uu.foodRecipeGroup.AdminDeleteRecipe(ctx, request)
}

func (uu UseCase) AdminDeleteFoodRecipeGroup(ctx context.Context, id int64) error {
	return uu.foodRecipeGroup.AdminDelete(ctx, id)
}

// @branch

func (uu UseCase) BranchGetFoodRecipeGroupList(ctx context.Context, filter food_recipe_group.Filter) ([]food_recipe_group.BranchGetListByFoodID, int, error) {
	return uu.foodRecipeGroup.BranchGetList(ctx, filter, int64(*filter.FoodId))
}

func (uu UseCase) BranchGetFoodRecipeGroupDetail(ctx context.Context, restaurantID int64) (*food_recipe_group.BranchGetDetail, error) {
	return uu.foodRecipeGroup.BranchGetDetail(ctx, restaurantID)
}

func (uu UseCase) BranchCreateFoodRecipeGroup(ctx context.Context, request food_recipe_group.BranchCreateRequest) (*food_recipe_group.BranchCreateResponse, error) {
	return uu.foodRecipeGroup.BranchCreate(ctx, request)
}

func (uu UseCase) BranchUpdateFoodRecipeGroupColumns(ctx context.Context, request food_recipe_group.BranchUpdateRequest) error {
	return uu.foodRecipeGroup.BranchUpdateColumns(ctx, request)
}

func (uu UseCase) BranchUpdateFoodRecipeGroupAll(ctx context.Context, request food_recipe_group.BranchUpdateRequest) error {
	return uu.foodRecipeGroup.BranchUpdateAll(ctx, request)
}

func (uu UseCase) BranchDeleteFoodRecipeGroupSingleRecipe(ctx context.Context, request food_recipe_group.BranchDeleteRecipeRequest) error {
	return uu.foodRecipeGroup.BranchDeleteRecipe(ctx, request)
}

func (uu UseCase) BranchDeleteFoodRecipeGroup(ctx context.Context, id int64) error {
	return uu.foodRecipeGroup.BranchDelete(ctx, id)
}

// #food_recipe_group_history

// @admin

func (uu UseCase) AdminGetFoodRecipeGroupHistoryList(ctx context.Context, filter food_recipe_group_history.Filter) ([]food_recipe_group_history.AdminGetListByFoodID, int, error) {
	return uu.foodRecipeGroupHistory.AdminGetList(ctx, filter)
}

func (uu UseCase) AdminGetFoodRecipeGroupHistoryDetail(ctx context.Context, restaurantID int64) (*food_recipe_group_history.AdminGetDetail, error) {
	return uu.foodRecipeGroupHistory.AdminGetDetail(ctx, restaurantID)
}

func (uu UseCase) AdminCreateFoodRecipeGroupHistory(ctx context.Context, request food_recipe_group_history.AdminCreateRequest) (*food_recipe_group_history.AdminCreateResponse, error) {
	return uu.foodRecipeGroupHistory.AdminCreate(ctx, request)
}

func (uu UseCase) AdminDeleteFoodRecipeGroupHistory(ctx context.Context, id int64) error {
	return uu.foodRecipeGroupHistory.AdminDelete(ctx, id)
}

// @branch

func (uu UseCase) BranchGetFoodRecipeGroupHistoryList(ctx context.Context, filter food_recipe_group_history.Filter) ([]food_recipe_group_history.BranchGetListByFoodID, int, error) {
	return uu.foodRecipeGroupHistory.BranchGetList(ctx, filter)
}

func (uu UseCase) BranchGetFoodRecipeGroupHistoryDetail(ctx context.Context, restaurantID int64) (*food_recipe_group_history.BranchGetDetail, error) {
	return uu.foodRecipeGroupHistory.BranchGetDetail(ctx, restaurantID)
}

func (uu UseCase) BranchCreateFoodRecipeGroupHistory(ctx context.Context, request food_recipe_group_history.BranchCreateRequest) (*food_recipe_group_history.BranchCreateResponse, error) {
	return uu.foodRecipeGroupHistory.BranchCreate(ctx, request)
}

func (uu UseCase) BranchDeleteFoodRecipeGroupHistory(ctx context.Context, id int64) error {
	return uu.foodRecipeGroupHistory.BranchDelete(ctx, id)
}
