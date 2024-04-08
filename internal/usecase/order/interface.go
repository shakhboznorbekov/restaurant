package order

import (
	"context"
	"github.com/restaurant/internal/entity"
	"github.com/restaurant/internal/service/order"
	"github.com/restaurant/internal/service/order_menu"
	"github.com/restaurant/internal/service/order_payment"
)

type Order interface {

	// @admin

	AdminGetListOrder(ctx context.Context, filter order.Filter) ([]order.AdminList, int, error)

	// cashier

	CashierListOrder(ctx context.Context, filter order.Filter) ([]order.CashierGetList, int, error)
	CashierDetailOrder(ctx context.Context, id int64) (order.CashierGetDetail, error)
	CashierUpdateStatus(ctx context.Context, data order.CashierUpdateStatusRequest) error

	// @client

	MobileListOrder(ctx context.Context, filter order.Filter) ([]order.MobileGetList, int, error)
	MobileDetailOrder(ctx context.Context, id int64) (order.MobileGetDetail, error)
	MobileCreateOrder(ctx context.Context, data order.MobileCreateRequest) (*order.MobileCreateResponse, error)
	MobileUpdateOrder(ctx context.Context, data order.MobileUpdateRequest) error
	MobileDeleteOrder(ctx context.Context, id int64) error
	ClientReview(ctx context.Context, request order.ClientReviewRequest) (*order.ClientReviewResponse, error)

	// @waiter

	WaiterGetList(ctx context.Context, filter order.Filter) ([]order.WaiterGetListResponse, int, error)
	WaiterCreate(ctx context.Context, data order.WaiterCreateRequest) (*order.WaiterCreateResponse, error)
	WaiterGetDetail(ctx context.Context, id int64) (*order.WaiterGetDetailResponse, error)
	WaiterUpdate(ctx context.Context, data order.WaiterUpdateRequest) error
	WaiterUpdateStatus(ctx context.Context, id int64, status string) error
	WaiterAccept(ctx context.Context, id int64) (*order.WaiterAcceptOrderResponse, error)
	WaiterGetHistoryActivityList(ctx context.Context, filter order.Filter) ([]order.HistoryActivityListResponse, int, error)
	WaiterGetMyOrderDetail(ctx context.Context, id int64) (*order.WaiterGetOrderDetailResponse, error)

	GetWsMessage(ctx context.Context, orderId int64) (order.GetWsMessageResponse, error)
	GetWsOrderMenus(ctx context.Context, orderId int64, menus []order.Menu) ([]order.GetWsOrderMenusResponse, int64, error)
	CheckOrderIfAccepted(id int64) error
	CancelOrder(id int64) error
	GetWsWaiter(waiterID int64) (order.GetWsWaiterResponse, error)
	OrderChecking(ctx context.Context, Time int) ([]order.GetWsMessageResponse, error)
}

type OrderMenu interface {

	// @client

	ClientGetList(ctx context.Context, filter order_menu.Filter) ([]order_menu.ClientGetList, int, error)
	ClientGetDetail(ctx context.Context, id int64) (entity.OrderMenu, error)
	ClientCreate(ctx context.Context, request order_menu.ClientCreateRequest) (order_menu.ClientCreateResponse, error)
	ClientUpdateAll(ctx context.Context, request order_menu.ClientUpdateRequest) error
	ClientUpdateColumns(ctx context.Context, request order_menu.ClientUpdateRequest) error
	ClientDelete(ctx context.Context, id int64) error
	ClientGetOftenList(ctx context.Context, branchID int) ([]order_menu.ClientGetOftenList, error)

	// @waiter
	WaiterUpdateStatus(ctx context.Context, ids []int64, status string) error

	// @waiter

	CashierUpdateStatus(ctx context.Context, id int64, status string) error
	CashierUpdateStatusByOrderID(ctx context.Context, orderId int64, status string) error
}

type OrderPayment interface {

	// @client

	ClientGetList(ctx context.Context, filter order_payment.Filter) ([]order_payment.ClientGetList, int, error)
	ClientGetDetail(ctx context.Context, id int64) (entity.OrderPayment, error)
	ClientCreate(ctx context.Context, request order_payment.ClientCreateRequest) (order_payment.ClientCreateResponse, error)
	ClientUpdateAll(ctx context.Context, request order_payment.ClientUpdateRequest) error
	ClientUpdateColumns(ctx context.Context, request order_payment.ClientUpdateRequest) error
	ClientDelete(ctx context.Context, id int64) error

	CashierCreate(ctx context.Context, request order_payment.CashierCreateRequest) (order_payment.CashierCreateResponse, error)
}

type OrderReport interface {
	CashierOrderReport(ctx context.Context) error
}
