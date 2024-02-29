package user

import (
	"context"
	"github.com/restaurant/internal/entity"
	//"github.com/restaurant/internal/service/attendance"
	//"github.com/restaurant/internal/service/service_percentage"
	//"github.com/restaurant/internal/service/sms"
	"github.com/restaurant/internal/service/user"
	//"github.com/restaurant/internal/service/waiter"
	//wwt "github.com/restaurant/internal/service/waiter_work_time"
)

type User interface {

	// @super-admin

	SuperAdminGetList(ctx context.Context, filter user.Filter) ([]user.SuperAdminGetList, int, error)
	SuperAdminGetDetail(ctx context.Context, id int64) (entity.User, error)
	SuperAdminCreate(ctx context.Context, request user.SuperAdminCreateRequest) (user.SuperAdminCreateResponse, error)
	SuperAdminUpdateAll(ctx context.Context, request user.SuperAdminUpdateRequest) error
	SuperAdminUpdateColumns(ctx context.Context, request user.SuperAdminUpdateRequest) error
	SuperAdminDelete(ctx context.Context, id int64) error

	// client

	ClientGetMe(ctx context.Context, id int64) (entity.User, error)
	ClientUpdateColumn(ctx context.Context, request user.ClientUpdateRequest) error
	ClientDeleteMe(ctx context.Context) error

	// others

	IsPhoneExists(ctx context.Context, phone string) (bool, error)
}

//
//type Waiter interface {
//
//	// @admin
//
//	AdminGetList(ctx context.Context, filter waiter.Filter) ([]waiter.AdminGetList, int, error)
//
//	// @branch
//
//	BranchGetList(ctx context.Context, filter waiter.Filter) ([]waiter.BranchGetList, int, error)
//	BranchGetDetail(ctx context.Context, id int64) (waiter.BranchGetDetail, error)
//	BranchCreate(ctx context.Context, request waiter.BranchCreateRequest) (waiter.BranchCreateResponse, error)
//	BranchUpdateAll(ctx context.Context, request waiter.BranchUpdateRequest) error
//	BranchUpdateColumns(ctx context.Context, request waiter.BranchUpdateRequest) error
//	BranchDelete(ctx context.Context, id int64) error
//	BranchUpdateStatus(ctx context.Context, id int64, status string) error
//
//	// others
//
//	UpdatePassword(ctx context.Context, request waiter.BranchUpdatePassword) error
//	UpdatePhone(ctx context.Context, request waiter.BranchUpdatePhone) error
//
//	// @waiter
//
//	WaiterGetMe(ctx context.Context) (*waiter.GetMeResponse, error)
//	WaiterGetPersonalInfo(ctx context.Context) (*waiter.GetPersonalInfoResponse, error)
//	WaiterUpdatePhoto(ctx context.Context, request waiter.WaiterPhotoUpdateRequest) error
//}
//
//type Sms interface {
//	SendSMS(ctx context.Context, send sms.Send) error
//	WaiterCheckSMSCode(ctx context.Context, check sms.Check) (bool, error)
//}
//
//type ServicePercentage interface {
//	AdminGetDetail(ctx context.Context, id int64) (*service_percentage.AdminGetDetail, error)
//}
//
//type AttendanceService interface {
//	WaiterCameCreate(ctx context.Context) (attendance.WaiterCreateResponse, error)
//	WaiterGoneCreate(ctx context.Context) (attendance.WaiterCreateResponse, error)
//}
//
//type WaiterWorkTimeService interface {
//	WaiterCreate(ctx context.Context, request attendance.WaiterCreateResponse) error
//	GetDetailByWaiterIDAndDate(ctx context.Context, filter wwt.Filter) (wwt.GetDetailByWaiterIDAndDateResponse, error)
//	WaiterGetListWorkTime(ctx context.Context, filter wwt.ListFilter) ([]wwt.GetListResponse, int, error)
//
//	BranchGetListWaiterWorkTime(ctx context.Context, filter wwt.BranchFilter) ([]wwt.BranchGetListResponse, int, error)
//	CashierGetListWaiterWorkTime(ctx context.Context, filter wwt.BranchFilter) ([]wwt.BranchGetListResponse, int, error)
//
//	BranchGetDetailWaiterWorkTime(ctx context.Context, filter wwt.ListFilter) ([]wwt.GetListResponse, int, error)
//	CashierGetDetailWaiterWorkTime(ctx context.Context, filter wwt.ListFilter) ([]wwt.GetListResponse, int, error)
//}
