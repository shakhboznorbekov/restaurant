package order

import "context"

type Repository interface {

	// mobile-client

	MobileListOrder(ctx context.Context, filter Filter) ([]MobileGetList, int, error)
	MobileDetailOrder(ctx context.Context, id int64) (MobileGetDetail, error)
	MobileCreateOrder(ctx context.Context, data MobileCreateRequest) (*MobileCreateResponse, error)
	MobileUpdateOrder(ctx context.Context, data MobileUpdateRequest) error
	MobileDeleteOrder(ctx context.Context, id int64) error
	ClientReview(ctx context.Context, request ClientReviewRequest) (*ClientReviewResponse, error)

	// admin

	AdminGetListOrder(ctx context.Context, filter Filter) ([]AdminList, int, error)

	// cashier

	CashierListOrder(ctx context.Context, filter Filter) ([]CashierGetList, int, error)
	CashierDetailOrder(ctx context.Context, id int64) (CashierGetDetail, error)
	CashierUpdateStatus(ctx context.Context, data CashierUpdateStatusRequest) error

	// @waiter

	WaiterGetList(ctx context.Context, filter Filter) ([]WaiterGetListResponse, int, error)
	WaiterCreate(ctx context.Context, data WaiterCreateRequest) (*WaiterCreateResponse, error)
	WaiterGetDetail(ctx context.Context, id int64) (*WaiterGetDetailResponse, error)
	WaiterUpdate(ctx context.Context, data WaiterUpdateRequest) error
	WaiterUpdateStatus(ctx context.Context, id int64, status string) error
	WaiterAccept(ctx context.Context, id int64) (*WaiterAcceptOrderResponse, error)

	// others

	GetWsOrderMenus(ctx context.Context, orderId int64, menus []Menu) ([]GetWsOrderMenusResponse, int64, error)
	GetWsMessage(ctx context.Context, orderId int64) (GetWsMessageResponse, error)
	CheckOrderIfAccepted(id int64) error
	CancelOrder(id int64) error
	GetWsWaiter(waiterID int64) (GetWsWaiterResponse, error)
}
