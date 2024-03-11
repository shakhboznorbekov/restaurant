package order

import "context"

type Service struct {
	repo Repository
}

func NewService(repo Repository) *Service {
	return &Service{repo: repo}
}

// @admin

func (s *Service) AdminGetListOrder(ctx context.Context, filter Filter) ([]AdminList, int, error) {
	return s.repo.AdminGetListOrder(ctx, filter)
}

// @client

func (s *Service) MobileListOrder(ctx context.Context, filter Filter) ([]MobileGetList, int, error) {
	return s.repo.MobileListOrder(ctx, filter)
}

func (s *Service) MobileDetailOrder(ctx context.Context, id int64) (MobileGetDetail, error) {
	return s.repo.MobileDetailOrder(ctx, id)
}

func (s *Service) MobileDeleteOrder(ctx context.Context, id int64) error {
	return s.repo.MobileDeleteOrder(ctx, id)
}

func (s *Service) MobileCreateOrder(ctx context.Context, request MobileCreateRequest) (*MobileCreateResponse, error) {
	return s.repo.MobileCreateOrder(ctx, request)
}

func (s *Service) MobileUpdateOrder(ctx context.Context, data MobileUpdateRequest) error {
	return s.repo.MobileUpdateOrder(ctx, data)
}

func (s *Service) ClientReview(ctx context.Context, request ClientReviewRequest) (*ClientReviewResponse, error) {
	return s.repo.ClientReview(ctx, request)
}

// cashier

func (s *Service) CashierListOrder(ctx context.Context, filter Filter) ([]CashierGetList, int, error) {
	return s.repo.CashierListOrder(ctx, filter)
}

func (s *Service) CashierDetailOrder(ctx context.Context, id int64) (CashierGetDetail, error) {
	return s.repo.CashierDetailOrder(ctx, id)
}

func (s *Service) CashierUpdateStatus(ctx context.Context, data CashierUpdateStatusRequest) error {
	return s.repo.CashierUpdateStatus(ctx, data)
}

// @waiter

func (s *Service) WaiterGetList(ctx context.Context, filter Filter) ([]WaiterGetListResponse, int, error) {
	return s.repo.WaiterGetList(ctx, filter)
}

func (s *Service) WaiterCreate(ctx context.Context, data WaiterCreateRequest) (*WaiterCreateResponse, error) {
	return s.repo.WaiterCreate(ctx, data)
}

func (s *Service) WaiterGetDetail(ctx context.Context, id int64) (*WaiterGetDetailResponse, error) {
	return s.repo.WaiterGetDetail(ctx, id)
}

func (s *Service) WaiterUpdate(ctx context.Context, request WaiterUpdateRequest) error {
	return s.repo.WaiterUpdate(ctx, request)
}

func (s *Service) WaiterUpdateStatus(ctx context.Context, id int64, status string) error {
	return s.repo.WaiterUpdateStatus(ctx, id, status)
}

func (s *Service) WaiterAccept(ctx context.Context, id int64) (*WaiterAcceptOrderResponse, error) {
	return s.repo.WaiterAccept(ctx, id)
}

// #ws

//func (s *Service) GetWsMessage(ctx context.Context, orderID int64) (GetWsMessageResponse, error) {
//	return s.repo.GetWsMessage(ctx, orderID)
//}
//
//func (s *Service) GetWsOrderMenus(ctx context.Context, orderId int64, menus []Menu) ([]GetWsOrderMenusResponse, int64, error) {
//	return s.repo.GetWsOrderMenus(ctx, orderId, menus)
//}
//
//func (s *Service) GetWsWaiter(waiterID int64) (GetWsWaiterResponse, error) {
//	return s.repo.GetWsWaiter(waiterID)
//}

// others

func (s *Service) CheckOrderIfAccepted(id int64) error {
	return s.repo.CheckOrderIfAccepted(id)
}

func (s *Service) CancelOrder(id int64) error {
	return s.repo.CancelOrder(id)
}
