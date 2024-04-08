package food

import (
	"context"
	"github.com/restaurant/internal/entity"
	"github.com/restaurant/internal/service/basket"
	"github.com/restaurant/internal/service/category"
	"github.com/restaurant/internal/service/food"
	"github.com/restaurant/internal/service/foodCategory"
	"github.com/restaurant/internal/service/food_recipe"
	"github.com/restaurant/internal/service/food_recipe_group"
	"github.com/restaurant/internal/service/food_recipe_group_history"
	"github.com/restaurant/internal/service/menu"
	"github.com/restaurant/internal/service/menu_category"
)

type FoodCategory interface {

	// @admin

	AdminGetList(ctx context.Context, filter foodCategory.Filter) ([]foodCategory.AdminGetList, int, error)
	AdminGetDetail(ctx context.Context, id int64) (entity.FoodCategory, error)
	AdminCreate(ctx context.Context, request foodCategory.AdminCreateRequest) (foodCategory.AdminCreateResponse, error)
	AdminUpdateAll(ctx context.Context, request foodCategory.AdminUpdateRequest) error
	AdminUpdateColumns(ctx context.Context, request foodCategory.AdminUpdateRequest) error
	AdminDelete(ctx context.Context, id int64) error

	// @branch

	BranchGetList(ctx context.Context, filter foodCategory.Filter) ([]foodCategory.BranchGetList, int, error)
	BranchGetDetail(ctx context.Context, id int64) (entity.FoodCategory, error)
	BranchCreate(ctx context.Context, request foodCategory.BranchCreateRequest) (foodCategory.BranchCreateResponse, error)
	BranchUpdateAll(ctx context.Context, request foodCategory.BranchUpdateRequest) error
	BranchUpdateColumns(ctx context.Context, request foodCategory.BranchUpdateRequest) error
	BranchDelete(ctx context.Context, id int64) error

	// @cashier

	CashierGetList(ctx context.Context, filter foodCategory.Filter) ([]foodCategory.CashierGetList, int, error)
	CashierGetDetail(ctx context.Context, id int64) (entity.FoodCategory, error)
	CashierCreate(ctx context.Context, request foodCategory.CashierCreateRequest) (foodCategory.CashierCreateResponse, error)
	CashierUpdateAll(ctx context.Context, request foodCategory.CashierUpdateRequest) error
	CashierUpdateColumns(ctx context.Context, request foodCategory.CashierUpdateRequest) error
	CashierDelete(ctx context.Context, id int64) error

	// client

	ClientGetList(ctx context.Context, filter foodCategory.Filter) ([]foodCategory.ClientGetList, int, error)

	// @waiter

	WaiterGetList(ctx context.Context) ([]foodCategory.WaiterGetList, error)
}

type Food interface {

	// @admin

	AdminGetList(ctx context.Context, filter food.Filter) ([]food.AdminGetList, int, error)
	AdminGetDetail(ctx context.Context, id int64) (entity.Foods, error)
	AdminCreate(ctx context.Context, request food.AdminCreateRequest) (food.AdminCreateResponse, error)
	AdminUpdateAll(ctx context.Context, request food.AdminUpdateRequest) error
	AdminUpdateColumns(ctx context.Context, request food.AdminUpdateRequest) error
	AdminDelete(ctx context.Context, id int64) error
	AdminDeleteImage(ctx context.Context, request food.AdminDeleteImageRequest) error

	// @admin

	BranchGetList(ctx context.Context, filter food.Filter) ([]food.BranchGetList, int, error)
	BranchGetDetail(ctx context.Context, id int64) (entity.Foods, error)
	BranchCreate(ctx context.Context, request food.BranchCreateRequest) (food.BranchCreateResponse, error)
	BranchUpdateAll(ctx context.Context, request food.BranchUpdateRequest) error
	BranchUpdateColumns(ctx context.Context, request food.BranchUpdateRequest) error
	BranchDelete(ctx context.Context, id int64) error
	BranchDeleteImage(ctx context.Context, request food.AdminDeleteImageRequest) error

	// @cashier

	CashierGetList(ctx context.Context, filter food.Filter) ([]food.CashierGetList, int, error)
}

type Menu interface {

	// @admin

	AdminGetList(ctx context.Context, filter menu.Filter) ([]menu.AdminGetList, int, error)
	AdminGetDetail(ctx context.Context, id int64) (entity.Menu, error)
	AdminCreate(ctx context.Context, request menu.AdminCreateRequest) ([]menu.AdminCreateResponse, error)
	AdminUpdateAll(ctx context.Context, request menu.AdminUpdateRequest) error
	AdminUpdateColumns(ctx context.Context, request menu.AdminUpdateRequest) error
	AdminDelete(ctx context.Context, id int64) error
	AdminRemovePhoto(ctx context.Context, id int64, index *int) (*string, error)

	// @client

	ClientGetList(ctx context.Context, filter menu.Filter) ([]menu.ClientGetList, error)
	ClientGetListByCategoryID(ctx context.Context, foodCategoryID int, filter menu.Filter) ([]menu.ClientGetListByCategoryID, error)

	// @branch

	BranchGetList(ctx context.Context, filter menu.Filter) ([]menu.BranchGetList, int, error)
	BranchGetDetail(ctx context.Context, id int64) (entity.Menu, error)
	BranchCreate(ctx context.Context, request menu.BranchCreateRequest) (menu.BranchCreateResponse, error)
	BranchUpdateAll(ctx context.Context, request menu.BranchUpdateRequest) error
	BranchUpdateColumns(ctx context.Context, request menu.BranchUpdateRequest) error
	BranchDelete(ctx context.Context, id int64) error
	BranchUpdatePrinterID(ctx context.Context, request menu.BranchUpdatePrinterIDRequest) error
	BranchDeletePrinterID(ctx context.Context, menuID int64) error
	BranchRemovePhoto(ctx context.Context, id int64, index int) (*string, error)

	// @cashier

	CashierUpdateColumns(ctx context.Context, request menu.CashierUpdateMenuStatus) error
	CashierGetList(ctx context.Context, filter menu.Filter) ([]menu.CashierGetList, int, error)
	CashierGetDetail(ctx context.Context, id int64) (entity.Menu, error)
	CashierCreate(ctx context.Context, request menu.CashierCreateRequest) (menu.CashierCreateResponse, error)
	CashierUpdateAll(ctx context.Context, request menu.CashierUpdateRequest) error
	CashierUpdateColumn(ctx context.Context, request menu.CashierUpdateRequest) error
	CashierDelete(ctx context.Context, id int64) error
	CashierUpdatePrinterID(ctx context.Context, request menu.CashierUpdatePrinterIDRequest) error
	CashierDeletePrinterID(ctx context.Context, menuID int64) error
	CashierRemovePhoto(ctx context.Context, id int64, index int) (*string, error)

	// @waiter

	WaiterGetMenuList(ctx context.Context, filter menu.Filter) ([]menu.WaiterGetMenuListResponse, error)
}

type FoodRecipe interface {
	AdminGetList(ctx context.Context, filter food_recipe.Filter) ([]food_recipe.AdminGetList, int, error)
	AdminGetDetail(ctx context.Context, id int64) (*food_recipe.AdminGetDetail, error)
	AdminCreate(ctx context.Context, request food_recipe.AdminCreateRequest) (*food_recipe.AdminCreateResponse, error)
	AdminUpdateAll(ctx context.Context, request food_recipe.AdminUpdateRequest) error
	AdminUpdateColumns(ctx context.Context, request food_recipe.AdminUpdateRequest) error
	AdminDelete(ctx context.Context, id int64) error

	// @branch

	BranchGetList(ctx context.Context, filter food_recipe.Filter) ([]food_recipe.BranchGetList, int, error)
	BranchGetDetail(ctx context.Context, id int64) (*food_recipe.BranchGetDetail, error)
	BranchCreate(ctx context.Context, request food_recipe.BranchCreateRequest) (*food_recipe.BranchCreateResponse, error)
	BranchUpdateAll(ctx context.Context, request food_recipe.BranchUpdateRequest) error
	BranchUpdateColumns(ctx context.Context, request food_recipe.BranchUpdateRequest) error
	BranchDelete(ctx context.Context, id int64) error
}

type Basket interface {
	GetBasket(ctx context.Context, key string) (basket.OrderStore, error)
}

type Category interface {

	// @super-admin

	SuperAdminGetList(ctx context.Context, filter category.Filter) ([]category.SuperAdminGetList, int, error)
	SuperAdminGetDetail(ctx context.Context, id int64) (entity.Category, error)
	SuperAdminCreate(ctx context.Context, request category.SuperAdminCreateRequest) (category.SuperAdminCreateResponse, error)
	SuperAdminUpdateAll(ctx context.Context, request category.SuperAdminUpdateRequest) error
	SuperAdminUpdateColumns(ctx context.Context, request category.SuperAdminUpdateRequest) error
	SuperAdminDelete(ctx context.Context, id int64) error

	// @client

	ClientGetList(ctx context.Context, filter category.Filter) ([]category.ClientGetList, int, error)

	// @branch

	BranchGetList(ctx context.Context, filter category.Filter) ([]category.BranchGetList, int, error)

	// @cashier

	CashierGetList(ctx context.Context, filter category.Filter) ([]category.CashierGetList, int, error)

	// @admin

	AdminGetList(ctx context.Context, filter category.Filter) ([]category.AdminGetList, int, error)

	// @waiter

	WaiterGetList(ctx context.Context) ([]category.WaiterGetList, error)
}

type MenuCategory interface {

	// @admin

	AdminGetList(ctx context.Context, filter menu_category.Filter) ([]entity.MenuCategory, int, error)
	AdminGetDetail(ctx context.Context, id int64) (*entity.MenuCategory, error)
	AdminCreate(ctx context.Context, request menu_category.AdminCreateRequest) (*entity.MenuCategory, error)
	AdminUpdate(ctx context.Context, request menu_category.AdminUpdateRequest) error
	AdminDelete(ctx context.Context, id int64) error

	// @branch

	BranchGetList(ctx context.Context, filter menu_category.Filter) ([]entity.MenuCategory, int, error)
	BranchGetDetail(ctx context.Context, id int64) (*entity.MenuCategory, error)
	BranchCreate(ctx context.Context, request menu_category.BranchCreateRequest) (*entity.MenuCategory, error)
	BranchUpdate(ctx context.Context, request menu_category.BranchUpdateRequest) error
	BranchDelete(ctx context.Context, id int64) error

	// cahier

	CashierGetList(ctx context.Context, filter menu_category.Filter) ([]entity.MenuCategory, int, error)
	CashierGetDetail(ctx context.Context, id int64) (*entity.MenuCategory, error)
	CashierCreate(ctx context.Context, request menu_category.CashierCreateRequest) (*entity.MenuCategory, error)
	CashierUpdate(ctx context.Context, request menu_category.CashierUpdateRequest) error
	CashierDelete(ctx context.Context, id int64) error
}

type FoodRecipeGroup interface {

	// @admin

	AdminGetList(ctx context.Context, filter food_recipe_group.Filter, foodID int64) ([]food_recipe_group.AdminGetListByFoodID, int, error)
	AdminGetDetail(ctx context.Context, id int64) (*food_recipe_group.AdminGetDetail, error)
	AdminCreate(ctx context.Context, request food_recipe_group.AdminCreateRequest) (*food_recipe_group.AdminCreateResponse, error)
	AdminUpdateAll(ctx context.Context, request food_recipe_group.AdminUpdateRequest) error
	AdminUpdateColumns(ctx context.Context, request food_recipe_group.AdminUpdateRequest) error
	AdminDeleteRecipe(ctx context.Context, request food_recipe_group.AdminDeleteRecipeRequest) error
	AdminDelete(ctx context.Context, id int64) error

	// @branch

	BranchGetList(ctx context.Context, filter food_recipe_group.Filter, foodID int64) ([]food_recipe_group.BranchGetListByFoodID, int, error)
	BranchGetDetail(ctx context.Context, id int64) (*food_recipe_group.BranchGetDetail, error)
	BranchCreate(ctx context.Context, request food_recipe_group.BranchCreateRequest) (*food_recipe_group.BranchCreateResponse, error)
	BranchUpdateAll(ctx context.Context, request food_recipe_group.BranchUpdateRequest) error
	BranchUpdateColumns(ctx context.Context, request food_recipe_group.BranchUpdateRequest) error
	BranchDeleteRecipe(ctx context.Context, request food_recipe_group.BranchDeleteRecipeRequest) error
	BranchDelete(ctx context.Context, id int64) error
}

type FoodRecipeGroupHistory interface {

	// @admin

	AdminGetList(ctx context.Context, filter food_recipe_group_history.Filter) ([]food_recipe_group_history.AdminGetListByFoodID, int, error)
	AdminGetDetail(ctx context.Context, id int64) (*food_recipe_group_history.AdminGetDetail, error)
	AdminCreate(ctx context.Context, request food_recipe_group_history.AdminCreateRequest) (*food_recipe_group_history.AdminCreateResponse, error)
	AdminDelete(ctx context.Context, id int64) error

	// @branch

	BranchGetList(ctx context.Context, filter food_recipe_group_history.Filter) ([]food_recipe_group_history.BranchGetListByFoodID, int, error)
	BranchGetDetail(ctx context.Context, id int64) (*food_recipe_group_history.BranchGetDetail, error)
	BranchCreate(ctx context.Context, request food_recipe_group_history.BranchCreateRequest) (*food_recipe_group_history.BranchCreateResponse, error)
	BranchDelete(ctx context.Context, id int64) error
}
