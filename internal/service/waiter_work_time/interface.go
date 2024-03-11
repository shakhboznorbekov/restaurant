package user_work_time

import (
	"context"
	"github.com/restaurant/internal/service/attendance"
)

type Repository interface {

	// @waiter

	WaiterCreate(ctx context.Context, request attendance.WaiterCreateResponse) error
	GetDetailByWaiterIDAndDate(ctx context.Context, filter Filter) (GetDetailByWaiterIDAndDateResponse, error)
	WaiterGetListWorkTime(ctx context.Context, filter ListFilter) ([]GetListResponse, int, error)

	//	 @branch

	BranchGetListWaiterWorkTime(ctx context.Context, filter BranchFilter) ([]BranchGetListResponse, int, error)
	BranchGetDetailWaiterWorkTime(ctx context.Context, filter ListFilter) ([]GetListResponse, int, error)

	CashierGetListWaiterWorkTime(ctx context.Context, filter BranchFilter) ([]BranchGetListResponse, int, error)
	CashierGetDetailWaiterWorkTime(ctx context.Context, filter ListFilter) ([]GetListResponse, int, error)
}
