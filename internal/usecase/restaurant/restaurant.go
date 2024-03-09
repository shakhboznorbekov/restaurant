package restaurant

import (
	"context"
	"github.com/pkg/errors"
	"github.com/restaurant/internal/pkg/file"
	"github.com/restaurant/internal/service/restaurant"
	"github.com/restaurant/internal/service/restaurant_category"
	"github.com/restaurant/internal/service/user"
)

type UseCase struct {
	restaurant         Restaurant
	user               User
	restaurantCategory RestaurantCategory
}

func NewUseCase(restaurant Restaurant, user User, restaurantCategory RestaurantCategory) *UseCase {
	return &UseCase{restaurant, user, restaurantCategory}
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
