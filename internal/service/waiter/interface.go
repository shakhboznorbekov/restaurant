package waiter

import (
	"context"
)

type Repository interface {

	// admin

	AdminGetList(ctx context.Context, filter Filter) ([]AdminGetList, int, error)

	// @branch

	BranchGetList(ctx context.Context, filter Filter) ([]BranchGetList, int, error)
	BranchGetDetail(ctx context.Context, id int64) (BranchGetDetail, error)
	BranchCreate(ctx context.Context, request BranchCreateRequest) (BranchCreateResponse, error)
	BranchUpdateAll(ctx context.Context, request BranchUpdateRequest) error
	BranchUpdateColumns(ctx context.Context, request BranchUpdateRequest) error
	BranchDelete(ctx context.Context, id int64) error
	BranchUpdateStatus(ctx context.Context, id int64, status string) error

	// others

	UpdatePassword(ctx context.Context, request BranchUpdatePassword) error
	UpdatePhone(ctx context.Context, request BranchUpdatePhone) error

	// @waiter

	WaiterGetMe(ctx context.Context) (*GetMeResponse, error)
	WaiterGetPersonalInfo(ctx context.Context) (*GetPersonalInfoResponse, error)
	WaiterUpdatePhoto(ctx context.Context, request WaiterPhotoUpdateRequest) error
	CalculateWaitersKPI(ctx context.Context) error
	WaiterGetActivityStatistics(ctx context.Context) (*GetActivityStatistics, error)
	WaiterGetWeeklyActivityStatistics(ctx context.Context, filter EarnedMoneyFilter) (*GetEarnedMoneyStatistics, error)
	WaiterGetWeeklyAcceptedOrdersStatistics(ctx context.Context, filter EarnedMoneyFilter) (*GetAcceptedOrdersStatistics, error)

	// @cashier

	CashierGetList(ctx context.Context, filter Filter) ([]CashierGetList, int, error)

	CashierGetLists(ctx context.Context, filter Filter) ([]CashierGetLists, int, error)
	CashierGetDetails(ctx context.Context, id int64) (CashierGetDetails, error)
	CashierCreate(ctx context.Context, request CashierCreateRequest) (CashierCreateResponse, error)
	CashierUpdateAll(ctx context.Context, request CashierUpdateRequest) error
	CashierUpdateColumns(ctx context.Context, request CashierUpdateRequest) error
	CashierDelete(ctx context.Context, id int64) error
	CashierUpdateStatus(ctx context.Context, id int64, status string) error

	// others

	CashierUpdatePassword(ctx context.Context, request CashierUpdatePassword) error
	CashierUpdatePhone(ctx context.Context, request CashierUpdatePhone) error
}
