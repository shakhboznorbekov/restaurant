package catalog

import (
	"context"
	"github.com/restaurant/internal/entity"
	"github.com/restaurant/internal/service/banner"
	"github.com/restaurant/internal/service/basket"
	"github.com/restaurant/internal/service/device"
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
)

type MeasureUnit interface {

	// @super-admin

	SuperAdminGetList(ctx context.Context, filter measureUnit.Filter) ([]measureUnit.SuperAdminGetList, int, error)
	SuperAdminGetDetail(ctx context.Context, id int64) (entity.MeasureUnit, error)
	SuperAdminCreate(ctx context.Context, request measureUnit.SuperAdminCreateRequest) (measureUnit.SuperAdminCreateResponse, error)
	SuperAdminUpdateAll(ctx context.Context, request measureUnit.SuperAdminUpdateRequest) error
	SuperAdminUpdateColumns(ctx context.Context, request measureUnit.SuperAdminUpdateRequest) error
	SuperAdminDelete(ctx context.Context, id int64) error

	// @admin

	AdminGetList(ctx context.Context, filter measureUnit.Filter) ([]measureUnit.AdminGetList, int, error)

	// @branch

	BranchGetList(ctx context.Context, filter measureUnit.Filter) ([]measureUnit.BranchGetList, int, error)
}

type Region interface {

	// @super-admin

	SuperAdminGetList(ctx context.Context, filter region.Filter) ([]region.SuperAdminGetList, int, error)
	SuperAdminGetDetail(ctx context.Context, id int64) (entity.Region, error)
	SuperAdminCreate(ctx context.Context, request region.SuperAdminCreateRequest) (region.SuperAdminCreateResponse, error)
	SuperAdminUpdateAll(ctx context.Context, request region.SuperAdminUpdateRequest) error
	SuperAdminUpdateColumns(ctx context.Context, request region.SuperAdminUpdateRequest) error
	SuperAdminDelete(ctx context.Context, id int64) error
}

type District interface {

	// @super-admin

	SuperAdminGetList(ctx context.Context, filter district.Filter) ([]district.SuperAdminGetList, int, error)
	SuperAdminGetDetail(ctx context.Context, id int64) (entity.District, error)
	SuperAdminCreate(ctx context.Context, request district.SuperAdminCreateRequest) (district.SuperAdminCreateResponse, error)
	SuperAdminUpdateAll(ctx context.Context, request district.SuperAdminUpdateRequest) error
	SuperAdminUpdateColumns(ctx context.Context, request district.SuperAdminUpdateRequest) error
	SuperAdminDelete(ctx context.Context, id int64) error
}

type Story interface {

	// @admin

	AdminGetList(ctx context.Context, filter story.Filter) ([]story.AdminGetList, int, error)
	AdminGetDetail(ctx context.Context, id int64) (entity.Story, error)
	AdminCreate(ctx context.Context, request story.AdminCreateRequest) (story.AdminCreateResponse, error)
	AdminDelete(ctx context.Context, id int64) error
	AdminUpdateStatus(ctx context.Context, id int64) error

	// @client

	ClientGetList(ctx context.Context, filter story.Filter) ([]story.ClientGetList, int, error)
	ClientSetViewed(ctx context.Context, id int64) error

	// @super-admin
	SuperAdminGetList(ctx context.Context, filter story.Filter) ([]story.SuperAdminGetListResponse, int, error)
	SuperAdminUpdateStatus(ctx context.Context, id int64, status string) error
}

type Banner interface {

	// @admin

	BranchGetList(ctx context.Context, filter banner.Filter) ([]banner.BranchGetList, int, error)
	BranchGetDetail(ctx context.Context, id int64) (*banner.BranchGetDetail, error)
	BranchCreate(ctx context.Context, request banner.BranchCreateRequest) (*banner.BranchCreateResponse, error)
	BranchUpdateAll(ctx context.Context, request banner.BranchUpdateRequest) error
	BranchUpdateColumn(ctx context.Context, request banner.BranchUpdateRequest) error
	BranchDelete(ctx context.Context, id int64) error
	BranchUpdateStatus(ctx context.Context, id int64, expireAt string) error

	// @client

	ClientGetList(ctx context.Context, filter banner.Filter) ([]banner.ClientGetList, int, error)
	ClientGetDetail(ctx context.Context, id int64) (*banner.ClientGetDetail, error)

	// @super-admin

	SuperAdminGetList(ctx context.Context, filter banner.Filter) ([]banner.SuperAdminGetListResponse, int, error)
	SuperAdminGetDetail(ctx context.Context, id int64) (*banner.SuperAdminGetDetailResponse, error)
	SuperAdminUpdateStatus(ctx context.Context, id int64, status string) error
}

type Feedback interface {

	// @admin

	AdminGetList(ctx context.Context, filter feedback.Filter) ([]feedback.AdminGetList, int, error)
	AdminGetDetail(ctx context.Context, id int64) (entity.Feedback, error)
	AdminCreate(ctx context.Context, request feedback.AdminCreate) (entity.Feedback, error)
	AdminUpdateColumns(ctx context.Context, request feedback.AdminUpdate) error
	AdminDelete(ctx context.Context, id int64) error

	// @client

	ClientGetList(ctx context.Context, filter feedback.Filter) ([]feedback.ClientGetList, int, error)
}

type Basket interface {
	SetBasket(ctx context.Context, data basket.Create) error
	GetBasket(ctx context.Context, key string) (basket.OrderStore, error)
	UpdateBasket(ctx context.Context, key string, value basket.Update) error
	DeleteBasket(ctx context.Context, key string) error
}

type FullSearch interface {
	ClientGetList(ctx context.Context, filter full_search.Filter) ([]full_search.ClientGetList, error)
}

type Menus interface {
	ClientGetDetail(ctx context.Context, id int64) (menu.ClientGetDetail, error)
}

type Notification interface {
	// @admin

	AdminGetList(ctx context.Context, filter notification.Filter) ([]notification.AdminGetList, int, error)
	AdminGetDetail(ctx context.Context, id int64) (*notification.AdminGetDetail, error)
	AdminUpdateAll(ctx context.Context, request notification.AdminUpdateRequest) error
	AdminUpdateColumn(ctx context.Context, request notification.AdminUpdateRequest) error
	AdminCreate(ctx context.Context, request notification.AdminCreateRequest) (*notification.AdminCreateResponse, error)
	AdminDelete(ctx context.Context, id int64) error

	// @super-admin

	SuperAdminUpdateStatus(ctx context.Context, id int64, status string) error
	SuperAdminGetList(ctx context.Context, filter notification.Filter) ([]notification.SuperAdminGetList, int, error)
	SuperAdminGetDetail(ctx context.Context, id int64) (*notification.SuperAdminGetDetail, error)
	SuperAdminSend(ctx context.Context, request notification.SuperAdminSendRequest) ([]notification.SuperAdminSendResponse, error)

	// @client

	ClientGetList(ctx context.Context, filter notification.Filter) ([]notification.ClientGetListResponse, int, error)
	ClientGetUnseenCount(ctx context.Context) (int, error)
	ClientSetAsViewed(ctx context.Context, id int64) error
	ClientSetAllAsViewed(ctx context.Context) error
}

type Firebase interface {
	SendCloudMessage(model fcm.CloudMessage) (string, error)
}

type Device interface {
	List(ctx context.Context, filter device.Filter) ([]entity.Device, int, error)
	Detail(ctx context.Context, id int64) (entity.Device, error)
}

type ServicePercentage interface {
	// @branch-admin

	AdminGetList(ctx context.Context, filter service_percentage.Filter) ([]service_percentage.AdminGetList, int, error)
	AdminGetDetail(ctx context.Context, id int64) (*service_percentage.AdminGetDetail, error)
	AdminCreate(ctx context.Context, request service_percentage.AdminCreateRequest) (service_percentage.AdminCreateResponse, error)
	AdminUpdateAll(ctx context.Context, request service_percentage.AdminUpdateRequest) error
	AdminDelete(ctx context.Context, id int64) error

	BranchCreate(ctx context.Context, request service_percentage.AdminCreateRequest) (service_percentage.AdminCreateResponse, error)
	AdminUpdateBranchID(ctx context.Context, request service_percentage.AdminUpdateBranchRequest) error
}
