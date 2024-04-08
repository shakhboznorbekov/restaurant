package catalog

import (
	"context"
	"fmt"
	"github.com/dariubs/percent"
	"github.com/pkg/errors"
	"github.com/restaurant/foundation/web"
	"github.com/restaurant/internal/auth"
	"github.com/restaurant/internal/commands"
	"github.com/restaurant/internal/entity"
	"github.com/restaurant/internal/pkg/file"
	"github.com/restaurant/internal/service/banner"
	"github.com/restaurant/internal/service/basket"
	"github.com/restaurant/internal/service/district"
	"github.com/restaurant/internal/service/fcm"
	"github.com/restaurant/internal/service/feedback"
	"github.com/restaurant/internal/service/full_search"
	"github.com/restaurant/internal/service/measureUnit"
	"github.com/restaurant/internal/service/menu"
	"github.com/restaurant/internal/service/notification"
	"github.com/restaurant/internal/service/region"
	"github.com/restaurant/internal/service/service_percentage"
	"github.com/restaurant/internal/service/story"
	"log"
	"net/http"
	"strconv"
	"strings"
)

type UseCase struct {
	measureUnit       MeasureUnit
	region            Region
	district          District
	story             Story
	banner            Banner
	feedback          Feedback
	basket            Basket
	fullSearch        FullSearch
	menus             Menus
	auth              *auth.Auth
	notification      Notification
	firebase          Firebase
	device            Device
	servicePercentage ServicePercentage
}

func NewUseCase(
	measureUnit MeasureUnit, region Region,
	district District, story Story,
	banner Banner, feedback Feedback,
	basket Basket, fullSearch FullSearch,
	menus Menus, auth *auth.Auth,
	notification Notification, firebase Firebase,
	device Device, servicePercentage ServicePercentage,
) *UseCase {
	return &UseCase{
		measureUnit,
		region,
		district,
		story,
		banner,
		feedback,
		basket,
		fullSearch,
		menus,
		auth,
		notification,
		firebase,
		device,
		servicePercentage,
	}
}

// #measure-unit

// @admin

func (uu UseCase) AdminGetMeasureUnitList(ctx context.Context, filter measureUnit.Filter) ([]measureUnit.AdminGetList, int, error) {
	filter.Fields = make(map[string][]string)
	filter.Fields["measure_unit"] = []string{"id", "name"}

	list, count, err := uu.measureUnit.AdminGetList(ctx, filter)
	if err != nil {
		return nil, 0, err
	}

	return list, count, err
}

// @branch

func (uu UseCase) BranchGetMeasureUnitList(ctx context.Context, filter measureUnit.Filter) ([]measureUnit.BranchGetList, int, error) {
	filter.Fields = make(map[string][]string)
	filter.Fields["measure_unit"] = []string{"id", "name"}

	list, count, err := uu.measureUnit.BranchGetList(ctx, filter)
	if err != nil {
		return nil, 0, err
	}

	return list, count, err
}

// @super-admin

func (uu UseCase) SuperAdminGetMeasureUnitList(ctx context.Context, filter measureUnit.Filter) ([]measureUnit.SuperAdminGetList, int, error) {
	filter.Fields = make(map[string][]string)
	filter.Fields["measure_unit"] = []string{"id", "name"}

	list, count, err := uu.measureUnit.SuperAdminGetList(ctx, filter)
	if err != nil {
		return nil, 0, err
	}

	return list, count, err
}

func (uu UseCase) SuperAdminGetMeasureUnitDetail(ctx context.Context, id int64) (measureUnit.SuperAdminGetDetail, error) {
	var detail measureUnit.SuperAdminGetDetail

	data, err := uu.measureUnit.SuperAdminGetDetail(ctx, id)
	if err != nil {
		return measureUnit.SuperAdminGetDetail{}, err
	}

	detail.ID = data.ID
	detail.Name = data.Name

	return detail, nil
}

func (uu UseCase) SuperAdminCreateMeasureUnit(ctx context.Context, data measureUnit.SuperAdminCreateRequest) (measureUnit.SuperAdminCreateResponse, error) {
	return uu.measureUnit.SuperAdminCreate(ctx, data)
}

func (uu UseCase) SuperAdminUpdateMeasureUnit(ctx context.Context, data measureUnit.SuperAdminUpdateRequest) error {
	return uu.measureUnit.SuperAdminUpdateAll(ctx, data)
}

func (uu UseCase) SuperAdminUpdateMeasureUnitColumn(ctx context.Context, data measureUnit.SuperAdminUpdateRequest) error {
	return uu.measureUnit.SuperAdminUpdateColumns(ctx, data)
}

func (uu UseCase) SuperAdminDeleteMeasureUnit(ctx context.Context, id int64) error {
	return uu.measureUnit.SuperAdminDelete(ctx, id)
}

// #region

func (uu UseCase) SuperAdminGetRegionList(ctx context.Context, filter region.Filter) ([]region.SuperAdminGetList, int, error) {
	//filter.Fields = make(map[string][]string)
	//filter.Fields["regions"] = []string{"id", "name", "code"}

	list, count, err := uu.region.SuperAdminGetList(ctx, filter)
	if err != nil {
		return nil, 0, err
	}

	return list, count, err
}

func (uu UseCase) SuperAdminGetRegionDetail(ctx context.Context, id int64) (region.SuperAdminGetDetail, error) {
	var detail region.SuperAdminGetDetail

	data, err := uu.region.SuperAdminGetDetail(ctx, id)
	if err != nil {
		return region.SuperAdminGetDetail{}, err
	}

	detail.ID = data.ID
	detail.Name = data.Name
	detail.Code = data.Code

	return detail, nil
}

func (uu UseCase) SuperAdminCreateRegion(ctx context.Context, data region.SuperAdminCreateRequest) (region.SuperAdminCreateResponse, error) {
	detail, err := uu.region.SuperAdminCreate(ctx, data)
	if err != nil {
		return region.SuperAdminCreateResponse{}, err
	}

	return detail, err
}

func (uu UseCase) SuperAdminUpdateRegion(ctx context.Context, data region.SuperAdminUpdateRequest) error {
	return uu.region.SuperAdminUpdateAll(ctx, data)
}

func (uu UseCase) SuperAdminUpdateRegionColumn(ctx context.Context, data region.SuperAdminUpdateRequest) error {
	return uu.region.SuperAdminUpdateColumns(ctx, data)
}

func (uu UseCase) SuperAdminDeleteRegion(ctx context.Context, id int64) error {
	return uu.region.SuperAdminDelete(ctx, id)
}

// district

func (uu UseCase) SuperAdminGetDistrictList(ctx context.Context, filter district.Filter) ([]district.SuperAdminGetList, int, error) {
	//filter.Fields = make(map[string][]string)
	//filter.Fields["districts"] = []string{"id", "name", "code", "region_id"}
	//filter.Joins = make(map[string]utils.Joins)
	//
	//joinColumn := "id"
	//mainColumn := "region_id"
	//filter.Joins["regions"] = utils.Joins{JoinColumn: &joinColumn, MainColumn: &mainColumn}

	list, count, err := uu.district.SuperAdminGetList(ctx, filter)
	if err != nil {
		return nil, 0, err
	}

	return list, count, err
}

func (uu UseCase) SuperAdminGetDistrictDetail(ctx context.Context, id int64) (district.SuperAdminGetDetail, error) {
	var detail district.SuperAdminGetDetail

	data, err := uu.district.SuperAdminGetDetail(ctx, id)
	if err != nil {
		return district.SuperAdminGetDetail{}, err
	}

	detail.ID = data.ID
	detail.Name = data.Name
	detail.Code = data.Code
	detail.RegionID = data.RegionId

	return detail, nil
}

func (uu UseCase) SuperAdminCreateDistrict(ctx context.Context, data district.SuperAdminCreateRequest) (district.SuperAdminCreateResponse, error) {
	detail, err := uu.district.SuperAdminCreate(ctx, data)
	if err != nil {
		return district.SuperAdminCreateResponse{}, err
	}

	return detail, err
}

func (uu UseCase) SuperAdminUpdateDistrict(ctx context.Context, data district.SuperAdminUpdateRequest) error {
	return uu.district.SuperAdminUpdateAll(ctx, data)
}

func (uu UseCase) SuperAdminUpdateDistrictColumn(ctx context.Context, data district.SuperAdminUpdateRequest) error {
	return uu.district.SuperAdminUpdateColumns(ctx, data)
}

func (uu UseCase) SuperAdminDeleteDistrict(ctx context.Context, id int64) error {
	return uu.district.SuperAdminDelete(ctx, id)
}

// story

func (uu UseCase) AdminGetStoryList(ctx context.Context, filter story.Filter) ([]story.AdminGetList, int, error) {

	list, count, err := uu.story.AdminGetList(ctx, filter)
	if err != nil {
		return nil, 0, err
	}

	return list, count, err
}

func (uu UseCase) AdminCreateStory(ctx context.Context, data story.AdminCreateRequest) (story.AdminCreateResponse, error) {
	if data.File != nil {
		if data.File.Size > 31457280 {
			return story.AdminCreateResponse{}, web.NewRequestError(errors.New("the file size is large"), http.StatusBadRequest)
		}
		if data.Type != nil {
			if strings.ToUpper(*data.Type) == "VIDEO" {
				if ok := commands.CheckFileType(ctx, data.File, "video"); !ok {
					return story.AdminCreateResponse{}, web.NewRequestError(errors.New("type is invalid"), http.StatusBadRequest)
				}
				Type := "VIDEO"
				data.Type = &Type
			} else if strings.ToUpper(*data.Type) == "IMAGE" {
				if ok := commands.CheckFileType(ctx, data.File, "image"); !ok {
					return story.AdminCreateResponse{}, web.NewRequestError(errors.New("type is invalid"), http.StatusBadRequest)
				}
				Type := "IMAGE"
				data.Type = &Type
			} else {
				return story.AdminCreateResponse{}, web.NewRequestError(errors.New("type is invalid"), http.StatusBadRequest)
			}
		} else {
			return story.AdminCreateResponse{}, web.NewRequestError(errors.New("not fount type"), http.StatusBadRequest)
		}
	} else {
		return story.AdminCreateResponse{}, web.NewRequestError(errors.New("not fount file"), http.StatusBadRequest)
	}

	if data.File != nil {
		imageLinks, _, err := file.UploadSingle(data.File, "story")
		if err != nil {
			return story.AdminCreateResponse{}, errors.Wrap(err, "logo upload")
		}
		data.FileLink = &imageLinks
	}

	detail, err := uu.story.AdminCreate(ctx, data)
	if err != nil {
		return story.AdminCreateResponse{}, err
	}

	return detail, err
}

func (uu UseCase) AdminDeleteStory(ctx context.Context, id int64) error {
	storyM, err := uu.story.AdminGetDetail(ctx, id)
	if err != nil {
		return err
	}

	err = uu.story.AdminDelete(ctx, id)
	if err != nil {
		return err
	}

	if storyM.File != nil {
		err = file.DeleteFiles(*storyM.File)
		if err != nil {
			return err
		}
	}

	return nil
}

func (uu UseCase) AdminUpdateStatusStory(ctx context.Context, id int64) error {
	return uu.story.AdminUpdateStatus(ctx, id)
}

func (uu UseCase) ClientGetStoryList(ctx context.Context, filter story.Filter) ([]story.ClientGetList, int, error) {

	list, count, err := uu.story.ClientGetList(ctx, filter)
	if err != nil {
		return nil, 0, err
	}

	return list, count, err
}

func (uu UseCase) ClientSetViewed(ctx context.Context, id int64) error {
	return uu.story.ClientSetViewed(ctx, id)
}

// @super-admin

func (uu UseCase) SuperAdminGetStoryList(ctx context.Context, filter story.Filter) ([]story.SuperAdminGetListResponse, int, error) {
	return uu.story.SuperAdminGetList(ctx, filter)
}

func (uu UseCase) SuperAdminUpdateStoryStatus(ctx context.Context, id int64, status string) error {
	return uu.story.SuperAdminUpdateStatus(ctx, id, status)
}

// banner

func (uu UseCase) BranchGetBannerList(ctx context.Context, filter banner.Filter) ([]banner.BranchGetList, int, error) {

	list, count, err := uu.banner.BranchGetList(ctx, filter)
	if err != nil {
		return nil, 0, err
	}

	return list, count, err
}

func (uu UseCase) BranchGetBannerByID(ctx context.Context, id int64) (*banner.BranchGetDetail, error) {
	return uu.banner.BranchGetDetail(ctx, id)
}

func (uu UseCase) BranchCreateBanner(ctx context.Context, data banner.BranchCreateRequest) (*banner.BranchCreateResponse, error) {
	if data.Photo != nil {
		if data.Photo.Size > 31457280 {
			return nil, web.NewRequestError(errors.New("the file size is large"), http.StatusBadRequest)
		}
		if ok := commands.CheckFileType(ctx, data.Photo, "image"); !ok {
			return nil, web.NewRequestError(errors.New("type is invalid"), http.StatusBadRequest)
		}

	} else {
		return nil, web.NewRequestError(errors.New("not fount file"), http.StatusBadRequest)
	}

	if data.Photo != nil {
		imageLinks, _, err := file.UploadSingle(data.Photo, "story")
		if err != nil {
			return nil, errors.Wrap(err, "logo upload")
		}
		data.PhotoLink = &imageLinks
	}

	detail, err := uu.banner.BranchCreate(ctx, data)
	if err != nil {
		return nil, err
	}

	return detail, err
}

func (uu UseCase) BranchUpdateAll(ctx context.Context, data banner.BranchUpdateRequest) error {
	if data.Photo != nil {
		if data.Photo.Size > 31457280 {
			return web.NewRequestError(errors.New("the file size is large"), http.StatusBadRequest)
		}
		if ok := commands.CheckFileType(ctx, data.Photo, "image"); !ok {
			return web.NewRequestError(errors.New("type is invalid"), http.StatusBadRequest)
		}
	}

	if data.Photo != nil {
		imageLinks, _, err := file.UploadSingle(data.Photo, "story")
		if err != nil {
			return errors.Wrap(err, "logo upload")
		}
		data.PhotoLink = &imageLinks
	}

	return uu.banner.BranchUpdateAll(ctx, data)
}

func (uu UseCase) BranchUpdateColumn(ctx context.Context, data banner.BranchUpdateRequest) error {
	if data.Photo != nil {
		if data.Photo.Size > 31457280 {
			return web.NewRequestError(errors.New("the file size is large"), http.StatusBadRequest)
		}
		if ok := commands.CheckFileType(ctx, data.Photo, "image"); !ok {
			return web.NewRequestError(errors.New("type is invalid"), http.StatusBadRequest)
		}
	}

	if data.Photo != nil {
		imageLinks, _, err := file.UploadSingle(data.Photo, "story")
		if err != nil {
			return errors.Wrap(err, "logo upload")
		}
		data.PhotoLink = &imageLinks
	}

	return uu.banner.BranchUpdateColumn(ctx, data)
}

func (uu UseCase) BranchDeleteBanner(ctx context.Context, id int64) error {
	return uu.banner.BranchDelete(ctx, id)
}

func (uu UseCase) BranchUpdateStatusBanner(ctx context.Context, id int64, expireAt string) error {
	return uu.banner.BranchUpdateStatus(ctx, id, expireAt)
}

func (uu UseCase) ClientGetBannerList(ctx context.Context, filter banner.Filter) ([]banner.ClientGetList, int, error) {

	list, count, err := uu.banner.ClientGetList(ctx, filter)
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

func (uu UseCase) ClientGetBannerDetail(ctx context.Context, id int64) (*banner.ClientGetDetail, error) {
	return uu.banner.ClientGetDetail(ctx, id)
}

// @super-admin

func (uu UseCase) SuperAdminGetBannerList(ctx context.Context, filter banner.Filter) ([]banner.SuperAdminGetListResponse, int, error) {
	return uu.banner.SuperAdminGetList(ctx, filter)
}

func (uu UseCase) SuperAdminGetBannerDetail(ctx context.Context, id int64) (*banner.SuperAdminGetDetailResponse, error) {
	return uu.banner.SuperAdminGetDetail(ctx, id)
}

func (uu UseCase) SuperAdminUpdateBannerStatus(ctx context.Context, id int64, status string) error {
	return uu.banner.SuperAdminUpdateStatus(ctx, id, status)
}

// feedback

// @admin

func (uu UseCase) AdminGetFeedbackList(ctx context.Context, filter feedback.Filter) ([]feedback.AdminGetList, int, error) {
	list, count, err := uu.feedback.AdminGetList(ctx, filter)
	if err != nil {
		return nil, 0, err
	}

	return list, count, err
}

func (uu UseCase) AdminGetFeedbackByID(ctx context.Context, id int64) (entity.Feedback, error) {
	res, err := uu.feedback.AdminGetDetail(ctx, id)
	if err != nil {
		return entity.Feedback{}, err
	}

	return res, err
}

func (uu UseCase) AdminCreateFeedback(ctx context.Context, data feedback.AdminCreate) (entity.Feedback, error) {
	detail, err := uu.feedback.AdminCreate(ctx, data)
	if err != nil {
		return entity.Feedback{}, err
	}

	return detail, err
}

func (uu UseCase) AdminUpdateFeedbackColumn(ctx context.Context, data feedback.AdminUpdate) error {
	return uu.feedback.AdminUpdateColumns(ctx, data)
}

func (uu UseCase) AdminDeleteFeedback(ctx context.Context, id int64) error {
	return uu.feedback.AdminDelete(ctx, id)
}

// @client

func (uu UseCase) ClientGetFeedbackList(ctx context.Context, filter feedback.Filter) ([]feedback.ClientGetList, int, error) {
	list, count, err := uu.feedback.ClientGetList(ctx, filter)
	if err != nil {
		return nil, 0, err
	}

	return list, count, err
}

// #basket

// @client

func (uu UseCase) ClientGetBasket(ctx context.Context, branchID int64, token string) (basket.Detail, error) {
	claims, err := uu.auth.GetTokenData(token)
	if err != nil {
		return basket.Detail{}, err
	}

	var key string
	key = fmt.Sprintf("%d_%d", claims.UserId, branchID)

	detail, err := uu.basket.GetBasket(ctx, key)
	if err != nil {
		return basket.Detail{}, err
	}

	var det basket.Detail
	det.UserID = detail.UserID
	det.TableID = detail.TableID

	foodsList := make([]menu.ClientGetDetail, 0)
	prices := make([]float32, 0)
	if detail.Foods != nil {
		for _, v := range detail.Foods {

			foodDetail, err := uu.menus.ClientGetDetail(ctx, v.ID)
			if err != nil {
				return basket.Detail{}, err
			}

			foodDetail.Count = v.Count

			if foodDetail.Price != nil {
				allPrice := *foodDetail.Price * float32(*v.Count)
				prices = append(prices, allPrice)
			}

			if foodDetail.Count != nil && *foodDetail.Count != 0 {
				foodsList = append(foodsList, foodDetail)
			}
		}
		det.Foods = foodsList
	}

	percentage := 10
	var sum float32
	for i := 0; i < len(prices); i++ {
		sum = sum + prices[i]
	}
	service := float32(percent.PercentFloat(float64(percentage), float64(sum)))
	overAll := sum + service

	det.Sum = &sum
	det.Service = &service
	det.OverAll = &overAll
	det.ServicePercentage = &percentage

	return det, err
}

func (uu UseCase) ClientUpdateBasket(ctx context.Context, data basket.Update, token string) error {
	claims, err := uu.auth.GetTokenData(token)
	if err != nil {
		return err
	}

	var key string
	key = fmt.Sprintf("%d_%d", claims.UserId, *data.BranchID)
	log.Println("541 -> ", key)

	data.UserID = &claims.UserId

	return uu.basket.UpdateBasket(ctx, key, data)
}

func (uu UseCase) ClientDeleteBasket(ctx context.Context, branchID int64, token string) error {
	claims, err := uu.auth.GetTokenData(token)
	if err != nil {
		return err
	}

	var key string
	key = fmt.Sprintf("%d_%d", claims.UserId, branchID)

	return uu.basket.DeleteBasket(ctx, key)
}

//func (uu UseCase) ClientCreateBasket(ctx context.Context, data basket.AdminCreate, token string) (string, error) {
//	claims, err := uu.auth.GetTokenData(token)
//	if err != nil {
//		return "", err
//	}
//
//	var key string
//	key = fmt.Sprintf("%d_%d", claims.UserId, *data.BranchID)
//	data.Key = &key
//
//	err = uu.basket.SetBasket(ctx, data)
//	if err != nil {
//		return "", errors.Wrap(err, "redis create")
//	}
//
//	return key, err
//}

// #full_search

// @client

func (uu UseCase) ClientGetFullSearchList(ctx context.Context, filter full_search.Filter) ([]full_search.ClientGetList, error) {
	list, err := uu.fullSearch.ClientGetList(ctx, filter)
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

// #notification

// @admin

func (uu UseCase) AdminGetNotificationList(ctx context.Context, filter notification.Filter) ([]notification.AdminGetList, int, error) {
	filter.Fields = make(map[string][]string)
	filter.Fields["notifications"] = []string{"id", "title", "description", "photo", "status"}

	list, count, err := uu.notification.AdminGetList(ctx, filter)
	if err != nil {
		return nil, 0, err
	}

	return list, count, err
}

func (uu UseCase) AdminGetNotificationDetail(ctx context.Context, id int64) (*notification.AdminGetDetail, error) {
	data, err := uu.notification.AdminGetDetail(ctx, id)
	if err != nil {
		return nil, err
	}

	return data, nil
}

func (uu UseCase) AdminCreateNotification(ctx context.Context, data notification.AdminCreateRequest) (*notification.AdminCreateResponse, error) {
	if data.Photo != nil {
		imageLink, _, err := file.UploadSingle(data.Photo, "notification")
		if err != nil {
			return nil, errors.Wrap(err, "photo upload")
		}
		data.PhotoLink = &imageLink
	}

	return uu.notification.AdminCreate(ctx, data)
}

func (uu UseCase) AdminUpdateNotification(ctx context.Context, data notification.AdminUpdateRequest) error {
	if data.Photo != nil {
		imageLink, _, err := file.UploadSingle(data.Photo, "notification")
		if err != nil {
			return errors.Wrap(err, "photo upload")
		}
		data.PhotoLink = &imageLink
	}

	return uu.notification.AdminUpdateAll(ctx, data)
}

func (uu UseCase) AdminUpdateColumnNotification(ctx context.Context, data notification.AdminUpdateRequest) error {
	if data.Photo != nil {
		imageLink, _, err := file.UploadSingle(data.Photo, "notification")
		if err != nil {
			return errors.Wrap(err, "photo upload")
		}
		data.PhotoLink = &imageLink
	}

	return uu.notification.AdminUpdateColumn(ctx, data)
}

func (uu UseCase) AdminDeleteNotification(ctx context.Context, id int64) error {
	return uu.notification.AdminDelete(ctx, id)
}

// @super-admin

func (uu UseCase) SuperAdminUpdateNotificationStatus(ctx context.Context, id int64, status string) error {
	if err := uu.notification.SuperAdminUpdateStatus(ctx, id, status); err != nil {
		return err
	}

	if status == "SENT" {
		detail, err := uu.notification.SuperAdminGetDetail(ctx, id)
		if err != nil {
			return err
		}

		if detail.DeviceTokens != nil && len(*detail.DeviceTokens) > 0 {
			for _, v := range *detail.DeviceTokens {
				// send cloud message
				notificate := fcm.CloudMessage{
					To: v,
					Notification: fcm.Notification{
						Title:          *detail.Title,
						Body:           *detail.Description,
						Image:          *detail.Photo,
						MutableContent: true,
						Sound:          "True-Tone",
					},
					Data: nil,
				}
				_, err = uu.firebase.SendCloudMessage(notificate)
				if err != nil {
					return err
				}
			}
		}
	}

	return nil
}

func (uu UseCase) SuperAdminGetNotificationList(ctx context.Context, filter notification.Filter) ([]notification.SuperAdminGetList, int, error) {
	return uu.notification.SuperAdminGetList(ctx, filter)
}

func (uu UseCase) SuperAdminSendNotification(ctx context.Context, request notification.SuperAdminSendRequest) error {
	if request.Photo != nil {
		photoLink, _, err := file.UploadSingle(request.Photo, "notification")
		if err != nil {
			return err
		}

		request.PhotoLink = &photoLink
	}

	response, err := uu.notification.SuperAdminSend(ctx, request)
	if err != nil {
		return err
	}

	if *request.Status == "SENT" {
		note := fcm.CloudMessage{
			Notification: fcm.Notification{
				MutableContent:   true,
				Sound:            "True-Tone",
				ContentAvailable: true,
			},
			Data: nil,
		}
		for _, v := range response {
			note.To = v.DeviceToken
			title, ok := request.Title[v.DeviceLang].(string)
			if !ok {
				continue
			}
			description, ok := request.Description[v.DeviceLang].(string)
			if !ok {
				continue
			}
			note.Notification.Title = title
			note.Notification.Body = description

			if request.PhotoLink != nil {
				note.Notification.Image = *request.PhotoLink
			}

			_, err = uu.firebase.SendCloudMessage(note)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

// @client

func (uu UseCase) ClientGetNotificationList(ctx context.Context, filter notification.Filter) ([]notification.ClientGetListResponse, int, error) {
	return uu.notification.ClientGetList(ctx, filter)
}

func (uu UseCase) ClientGetCountUnseenNotification(ctx context.Context) (int, error) {
	return uu.notification.ClientGetUnseenCount(ctx)
}

func (uu UseCase) ClientSetNotificationAsViewed(ctx context.Context, id int64) error {
	return uu.notification.ClientSetAsViewed(ctx, id)
}

func (uu UseCase) ClientSetAllNotificationsAsViewed(ctx context.Context) error {
	return uu.notification.ClientSetAllAsViewed(ctx)
}

// @branch-admin

func (uu UseCase) AdminGetServicePercentageList(ctx context.Context, filter service_percentage.Filter) ([]service_percentage.AdminGetList, int, error) {
	return uu.servicePercentage.AdminGetList(ctx, filter)
}

func (uu UseCase) AdminGetServicePercentageDetail(ctx context.Context, id int64) (*service_percentage.AdminGetDetail, error) {
	return uu.servicePercentage.AdminGetDetail(ctx, id)
}

func (uu UseCase) AdminCreateServicePercentage(ctx context.Context, request service_percentage.AdminCreateRequest) (service_percentage.AdminCreateResponse, error) {
	detail, err := uu.servicePercentage.AdminCreate(ctx, request)
	if err != nil {
		return service_percentage.AdminCreateResponse{}, err
	}

	return detail, err
}
func (uu UseCase) AdminUpdateServicePercentageAll(ctx context.Context, request service_percentage.AdminUpdateRequest) error {
	return uu.servicePercentage.AdminUpdateAll(ctx, request)
}

func (uu UseCase) AdminDeleteServicePercentage(ctx context.Context, id int64) error {
	return uu.servicePercentage.AdminDelete(ctx, id)
}

// @branch

func (uu UseCase) BranchCreateServicePercentage(ctx context.Context, request service_percentage.AdminCreateRequest) (service_percentage.AdminCreateResponse, error) {
	detail, err := uu.servicePercentage.BranchCreate(ctx, request)
	if err != nil {
		return service_percentage.AdminCreateResponse{}, err
	}

	return detail, err
}
