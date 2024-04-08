package user

import (
	"context"
	"github.com/restaurant/internal/entity"
	"github.com/restaurant/internal/service/attendance"
	"github.com/restaurant/internal/service/cashier"
	"github.com/restaurant/internal/service/service_percentage"
	"github.com/restaurant/internal/service/sms"
	"github.com/restaurant/internal/service/user"
	"github.com/restaurant/internal/service/waiter"
	wwt "github.com/restaurant/internal/service/waiter_work_time"
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

	//	@cashier

	CashierGetMe(ctx context.Context) (*user.CashierGetMeResponse, error)

	// general get me

	GetMe(ctx context.Context, userID int64) (*user.GetMeResponse, error)
}

type Waiter interface {

	// @admin

	AdminGetList(ctx context.Context, filter waiter.Filter) ([]waiter.AdminGetList, int, error)

	// @branch

	BranchGetList(ctx context.Context, filter waiter.Filter) ([]waiter.BranchGetList, int, error)
	BranchGetDetail(ctx context.Context, id int64) (waiter.BranchGetDetail, error)
	BranchCreate(ctx context.Context, request waiter.BranchCreateRequest) (waiter.BranchCreateResponse, error)
	BranchUpdateAll(ctx context.Context, request waiter.BranchUpdateRequest) error
	BranchUpdateColumns(ctx context.Context, request waiter.BranchUpdateRequest) error
	BranchDelete(ctx context.Context, id int64) error
	BranchUpdateStatus(ctx context.Context, id int64, status string) error

	// others

	UpdatePassword(ctx context.Context, request waiter.BranchUpdatePassword) error
	UpdatePhone(ctx context.Context, request waiter.BranchUpdatePhone) error

	// @waiter

	WaiterGetMe(ctx context.Context) (*waiter.GetMeResponse, error)
	WaiterGetPersonalInfo(ctx context.Context) (*waiter.GetPersonalInfoResponse, error)
	WaiterUpdatePhoto(ctx context.Context, request waiter.WaiterPhotoUpdateRequest) error
	CalculateWaitersKPI(ctx context.Context) error
	WaiterGetActivityStatistics(ctx context.Context) (*waiter.GetActivityStatistics, error)
	WaiterGetWeeklyActivityStatistics(ctx context.Context, filter waiter.EarnedMoneyFilter) (*waiter.GetEarnedMoneyStatistics, error)
	WaiterGetWeeklyAcceptedOrdersStatistics(ctx context.Context, filter waiter.EarnedMoneyFilter) (*waiter.GetAcceptedOrdersStatistics, error)
	WaiterGetWeeklyRatingStatistics(ctx context.Context, filter waiter.Filter) ([]waiter.GetWeeklyRating, error)

	// @cashier

	CashierGetList(ctx context.Context, filter waiter.Filter) ([]waiter.CashierGetList, int, error)

	// @cashier

	CashierGetLists(ctx context.Context, filter waiter.Filter) ([]waiter.CashierGetLists, int, error)
	CashierGetDetails(ctx context.Context, id int64) (waiter.CashierGetDetails, error)
	CashierCreate(ctx context.Context, request waiter.CashierCreateRequest) (waiter.CashierCreateResponse, error)
	CashierUpdateAll(ctx context.Context, request waiter.CashierUpdateRequest) error
	CashierUpdateColumns(ctx context.Context, request waiter.CashierUpdateRequest) error
	CashierDelete(ctx context.Context, id int64) error
	CashierUpdateStatus(ctx context.Context, id int64, status string) error

	// others

	CashierUpdatePassword(ctx context.Context, request waiter.CashierUpdatePassword) error
	CashierUpdatePhone(ctx context.Context, request waiter.CashierUpdatePhone) error
}

type Sms interface {
	SendSMS(ctx context.Context, send sms.Send) error
	WaiterCheckSMSCode(ctx context.Context, check sms.Check) (bool, error)
	CashierCheckSMSCode(ctx context.Context, check sms.Check) (bool, error)
}

type ServicePercentage interface {
	AdminGetDetail(ctx context.Context, id int64) (*service_percentage.AdminGetDetail, error)
}

type AttendanceService interface {
	WaiterCameCreate(ctx context.Context) (attendance.WaiterCreateResponse, error)
	WaiterGoneCreate(ctx context.Context) (attendance.WaiterCreateResponse, error)
}

type WaiterWorkTimeService interface {
	WaiterCreate(ctx context.Context, request attendance.WaiterCreateResponse) error
	GetDetailByWaiterIDAndDate(ctx context.Context, filter wwt.Filter) (wwt.GetDetailByWaiterIDAndDateResponse, error)
	WaiterGetListWorkTime(ctx context.Context, filter wwt.ListFilter) ([]wwt.GetListResponse, int, error)

	BranchGetListWaiterWorkTime(ctx context.Context, filter wwt.BranchFilter) ([]wwt.BranchGetListResponse, int, error)
	CashierGetListWaiterWorkTime(ctx context.Context, filter wwt.BranchFilter) ([]wwt.BranchGetListResponse, int, error)

	BranchGetDetailWaiterWorkTime(ctx context.Context, filter wwt.ListFilter) ([]wwt.GetListResponse, int, error)
	CashierGetDetailWaiterWorkTime(ctx context.Context, filter wwt.ListFilter) ([]wwt.GetListResponse, int, error)
}

type Cashier interface {

	// @admin

	AdminGetList(ctx context.Context, filter cashier.Filter) ([]cashier.AdminGetList, int, error)
	AdminGetDetail(ctx context.Context, id int64) (cashier.AdminGetDetail, error)
	AdminCreate(ctx context.Context, request cashier.AdminCreateRequest) (cashier.AdminCreateResponse, error)
	AdminUpdateAll(ctx context.Context, request cashier.AdminUpdateRequest) error
	AdminUpdateColumns(ctx context.Context, request cashier.AdminUpdateRequest) error
	AdminDelete(ctx context.Context, id int64) error
	AdminUpdateStatus(ctx context.Context, id int64, status string) error

	AdminUpdatePassword(ctx context.Context, request cashier.AdminUpdatePassword) error
	AdminUpdatePhone(ctx context.Context, request cashier.AdminUpdatePhone) error

	// @branch

	BranchGetList(ctx context.Context, filter cashier.Filter) ([]cashier.BranchGetList, int, error)
	BranchGetDetail(ctx context.Context, id int64) (cashier.BranchGetDetail, error)
	BranchCreate(ctx context.Context, request cashier.BranchCreateRequest) (cashier.BranchCreateResponse, error)
	BranchUpdateAll(ctx context.Context, request cashier.BranchUpdateRequest) error
	BranchUpdateColumns(ctx context.Context, request cashier.BranchUpdateRequest) error
	BranchDelete(ctx context.Context, id int64) error
	BranchUpdateStatus(ctx context.Context, id int64, status string) error

	// others

	UpdatePassword(ctx context.Context, request cashier.BranchUpdatePassword) error
	UpdatePhone(ctx context.Context, request cashier.BranchUpdatePhone) error

	// @cashier
}
