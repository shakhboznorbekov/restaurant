package order

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/pkg/errors"
	"github.com/restaurant/internal/service/order"
	"github.com/restaurant/internal/service/order_menu"
	"github.com/restaurant/internal/service/order_payment"
	"github.com/restaurant/internal/socket"
	"log"
	"time"
)

type UseCase struct {
	order        Order
	orderFood    OrderMenu
	orderPayment OrderPayment
	orderReport  OrderReport
	hub          *socket.Hub
}

func NewUseCase(order Order, orderFood OrderMenu, orderPayment OrderPayment, hub *socket.Hub, orderReport OrderReport) *UseCase {
	return &UseCase{order, orderFood, orderPayment, orderReport, hub}
}

// #order

// @waiter

func (uu UseCase) WaiterGetHistoryActivityList(ctx context.Context, filter order.Filter) ([]order.HistoryActivityListResponse, int, error) {
	return uu.order.WaiterGetHistoryActivityList(ctx, filter)
}

// @admin

func (uu UseCase) AdminGetListOrder(ctx context.Context, filter order.Filter) ([]order.AdminList, int, error) {
	return uu.order.AdminGetListOrder(ctx, filter)
}

// @mobile-client

func (uu UseCase) MobileCreateOrder(ctx context.Context, request order.MobileCreateRequest) (*order.MobileCreateResponse, error) {
	response, err := uu.order.MobileCreateOrder(ctx, request)
	if err != nil {
		return nil, err
	}

	message, err := uu.order.GetWsMessage(ctx, response.ID)
	if err != nil {
		return nil, err
	}

	// message broadcast to branch-admin and admin
	wsResponse := socket.Response{
		Action: "new_order",
		Property: socket.NewOrderResponse{
			OrderID:     message.OrderID,
			OrderNumber: message.OrderNumber,
			TableID:     message.TableID,
			TableNumber: message.TableNumber,
			Price:       message.Price,
			CreatedAt:   message.CreatedAt,
			Foods:       response.Menus,
		},
	}
	messageJson, _ := json.Marshal(&wsResponse)
	uu.hub.Broadcast <- socket.Message{
		Action:       "new_order",
		UserID:       message.UserId,
		RestaurantID: message.RestaurantId,
		BranchID:     message.BranchId,
		Message:      messageJson,
	}

	go func() {
		var count, duration int

		for {
			if duration > 600 {
				// message broadcast to user as cancelled
				wsResponse = socket.Response{
					Action: "accepted_order",
					Property: socket.OrderAcceptedResponse{
						OrderID:     &message.OrderID,
						OrderNumber: &message.OrderNumber,
						BranchID:    &message.BranchId,
						TableID:     &message.TableID,
						TableNumber: &message.TableNumber,
						ClientID:    &message.UserId,
						Status:      "CANCELLED",
					},
				}
				messageJson, _ = json.Marshal(&wsResponse)
				uu.hub.Broadcast <- socket.Message{
					UserID:  message.UserId,
					Action:  "client_accepted_order",
					Message: messageJson,
				}

				_ = uu.order.CancelOrder(response.ID)

				break
			}
			wsWaiter := make(map[int64]*socket.Client)
			if v, ok := uu.hub.Clients[fmt.Sprintf("WAITER_%d", message.BranchId)]; ok {
				wsWaiter = v
			}

			Waters := make([]order.GetWsWaiterResponse, 0)
			for _, c := range wsWaiter {
				w, err := uu.order.GetWsWaiter(c.ID)
				if err != nil {
					log.Println(err.Error())
					break
				}
				w.Grade = float32(c.NewOrder) + float32(w.OrderCount)/100
				Waters = append(Waters, w)
			}

			if len(Waters) < 1 {
				fmt.Println(111)
				// sleeps before order accepted...
				time.Sleep(time.Second * 10)
				duration += 10
			}
			Waters = sortWaiter(Waters)
			wsResponse = socket.Response{
				Action: "new_order",
				Property: socket.NewOrderWaiterResponse{
					OrderID:     message.OrderID,
					OrderNumber: message.OrderNumber,
					TableID:     message.TableID,
					TableNumber: message.TableNumber,
					Foods:       response.Menus,
					ClientCount: request.ClientCount,
					Price:       message.Price,
					CreatedAt:   message.CreatedAt,
					Time:        10,
				},
			}
			messageJson, _ = json.Marshal(&wsResponse)
			for _, w := range Waters {
				if v1, ok := uu.hub.Clients[fmt.Sprintf("WAITER_%d", message.BranchId)]; ok {
					if v2, ok := v1[w.ID]; ok {
						v2.Send <- socket.ResMessage{messageJson}
					}
				}

				// sleeps before order accepted...
				time.Sleep(time.Second * 10)
				duration += 10

				err = uu.order.CheckOrderIfAccepted(response.ID)
				if err != nil {
					return
				}
			}
			count++
			err = uu.order.CheckOrderIfAccepted(response.ID)
			if err != nil {
				return
			}
		}
	}()

	return response, err
}

func (uu UseCase) MobileUpdateOrder(ctx context.Context, data order.MobileUpdateRequest) error {
	err := uu.order.MobileUpdateOrder(ctx, data)
	if err != nil {
		return err
	}
	message, err := uu.order.GetWsMessage(ctx, data.Id)
	if err != nil {
		return err
	}

	wsResponse := socket.Response{
		Action: "new_food",
		Property: socket.NewFoodsResponse{
			OrderID:     message.OrderID,
			OrderNumber: message.OrderNumber,
			TableID:     message.TableID,
			TableNumber: message.TableNumber,
			Foods:       data.Menus,
			Price:       message.Price,
		},
	}
	messageJson, err := json.Marshal(&wsResponse)

	uu.hub.Broadcast <- socket.Message{
		Action:       "new_food",
		UserID:       message.UserId,
		RestaurantID: message.RestaurantId,
		BranchID:     message.BranchId,
		Message:      messageJson,
	}

	err = uu.order.CheckOrderIfAccepted(data.Id)
	if err == nil {
		return nil
	}
	printerMessage, branchID, err := uu.order.GetWsOrderMenus(ctx, data.Id, data.Menus)
	if err != nil {
		return err
	}
	for _, v := range printerMessage {
		wsResponse := socket.Response{
			Action:   "new_food",
			Property: v,
		}
		messageJson, _ := json.Marshal(&wsResponse)
		uu.hub.Broadcast <- socket.Message{
			BranchID: branchID,
			Action:   "printer_new_food",
			Message:  messageJson,
		}
	}

	return nil

}

func (uu UseCase) MobileGetOrderList(ctx context.Context, filter order.Filter) ([]order.MobileGetList, int, error) {
	list, count, err := uu.order.MobileListOrder(ctx, filter)
	if err != nil {
		return nil, 0, err
	}

	return list, count, err
}

func (uu UseCase) MobileGetOrderDetail(ctx context.Context, id int64) (order.MobileGetDetail, error) {
	data, err := uu.order.MobileDetailOrder(ctx, id)
	if err != nil {
		return order.MobileGetDetail{}, err
	}

	return data, nil
}

func (uu UseCase) MobileWaiterCall(ctx context.Context, orderID int64) error {
	message, err := uu.order.GetWsMessage(ctx, orderID)
	if err != nil {
		return err
	}

	if message.WaiterID == nil {
		err := errors.New("waiter does not accepted yet!")
		return err
	}

	// message broadcast to user as accepted
	go func() {
		wsResponse := socket.Response{
			Action: "client_call",
			Property: socket.ClientCall{
				OrderID:     &orderID,
				OrderNumber: &message.OrderNumber,
				TableID:     &message.TableID,
				TableNumber: &message.TableNumber,
			},
		}

		messageJson, _ := json.Marshal(&wsResponse)
		uu.hub.Broadcast <- socket.Message{
			UserID:   *message.WaiterID,
			BranchID: message.BranchId,
			Action:   "client_call",
			Message:  messageJson,
		}
	}()

	return nil
}

func (uu UseCase) ClientReviewOrder(ctx context.Context, request order.ClientReviewRequest) (*order.ClientReviewResponse, error) {
	return uu.order.ClientReview(ctx, request)
}

// @cashier

func (uu UseCase) CashierUpdateOrder(ctx context.Context, data order.CashierUpdateStatusRequest) error {
	err := uu.order.CashierUpdateStatus(ctx, data)
	if err != nil {
		return err
	}

	_, err = uu.orderPayment.CashierCreate(ctx, order_payment.CashierCreateRequest{
		OrderID: &data.Id,
	})
	if err != nil {
		return err
	}

	err = uu.orderFood.CashierUpdateStatusByOrderID(ctx, data.Id, "PAID")
	if err != nil {
		return err
	}

	message, err := uu.order.GetWsMessage(ctx, data.Id)
	if err != nil {
		return err
	}

	wsResponse := socket.Response{
		Action: "order_payment",
		Property: socket.OrderPaymentResponse{
			OrderID:     message.OrderID,
			OrderNumber: message.OrderNumber,
		},
	}

	messageJson, err := json.Marshal(&wsResponse)

	uu.hub.Broadcast <- socket.Message{
		Action:       "order_payment",
		UserID:       message.UserId,
		RestaurantID: message.RestaurantId,
		BranchID:     message.BranchId,
		Message:      messageJson,
	}

	return nil

}

// @others

func (uu UseCase) OrderChecking(ctx context.Context, Time int) error {
	messages, err := uu.order.OrderChecking(ctx, Time)
	if err != nil {
		return err
	}

	for _, m := range messages {
		go func(message order.GetWsMessageResponse) {
			fmt.Println(message.OrderID)
			var count, duration int

			for {
				if duration > 600 {
					// message broadcast to user as cancelled
					wsResponse := socket.Response{
						Action: "accepted_order",
						Property: socket.OrderAcceptedResponse{
							OrderID:     &message.OrderID,
							OrderNumber: &message.OrderNumber,
							BranchID:    &message.BranchId,
							TableID:     &message.TableID,
							TableNumber: &message.TableNumber,
							ClientID:    &message.UserId,
							Status:      "CANCELLED",
						},
					}
					messageJson, _ := json.Marshal(&wsResponse)
					uu.hub.Broadcast <- socket.Message{
						UserID:  message.UserId,
						Action:  "client_accepted_order",
						Message: messageJson,
					}

					_ = uu.order.CancelOrder(message.OrderID)

					break
				}
				wsWaiter := make(map[int64]*socket.Client)
				if v, ok := uu.hub.Clients[fmt.Sprintf("WAITER_%d", message.BranchId)]; ok {
					wsWaiter = v
				}

				Waters := make([]order.GetWsWaiterResponse, 0)
				for _, c := range wsWaiter {
					w, err := uu.order.GetWsWaiter(c.ID)
					if err != nil {
						log.Println(err.Error())
						break
					}
					w.Grade = float32(c.NewOrder) + float32(w.OrderCount)/100
					Waters = append(Waters, w)
				}

				if len(Waters) < 1 {
					// sleeps before order accepted...
					time.Sleep(time.Second * 10)
					duration += 10
				}
				Waters = sortWaiter(Waters)
				wsResponse := socket.Response{
					Action: "new_order",
					Property: socket.NewOrderWaiterResponse{
						OrderID:     message.OrderID,
						OrderNumber: message.OrderNumber,
						TableID:     message.TableID,
						TableNumber: message.TableNumber,
						Price:       message.Price,
						CreatedAt:   message.CreatedAt,
						Time:        10,
					},
				}
				messageJson, _ := json.Marshal(&wsResponse)
				for _, w := range Waters {
					if v1, ok := uu.hub.Clients[fmt.Sprintf("WAITER_%d", message.BranchId)]; ok {
						if v2, ok := v1[w.ID]; ok {
							v2.Send <- socket.ResMessage{messageJson}
						}
					}

					// sleeps before order accepted...
					time.Sleep(time.Second * 10)
					duration += 10

					err = uu.order.CheckOrderIfAccepted(message.OrderID)
					if err != nil {
						return
					}
				}
				count++
				err = uu.order.CheckOrderIfAccepted(message.OrderID)
				if err != nil {
					return
				}
			}
		}(m)
	}

	return err
}

func (uu UseCase) CashierGetOrderList(ctx context.Context, filter order.Filter) ([]order.CashierGetList, int, error) {
	list, count, err := uu.order.CashierListOrder(ctx, filter)
	if err != nil {
		return nil, 0, err
	}

	return list, count, err
}

func (uu UseCase) CashierGetOrderDetail(ctx context.Context, id int64) (order.CashierGetDetail, error) {
	data, err := uu.order.CashierDetailOrder(ctx, id)
	if err != nil {
		return order.CashierGetDetail{}, err
	}

	return data, nil
}

// #order_menu

// @mobile-client

func (uu UseCase) ClientGetOrderMenuList(ctx context.Context, filter order_menu.Filter) ([]order_menu.ClientGetList, int, error) {
	list, count, err := uu.orderFood.ClientGetList(ctx, filter)
	if err != nil {
		return nil, 0, err
	}

	return list, count, err
}

func (uu UseCase) ClientGetOrderMenuDetail(ctx context.Context, id int64) (order_menu.ClientGetDetail, error) {
	var detail order_menu.ClientGetDetail

	data, err := uu.orderFood.ClientGetDetail(ctx, id)
	if err != nil {
		return order_menu.ClientGetDetail{}, err
	}
	detail.ID = data.ID
	detail.Count = data.Count
	detail.MenuID = data.MenuID
	detail.OrderID = data.OrderID

	return detail, nil
}

func (uu UseCase) ClientGetOrderMenuOftenList(ctx context.Context, branchID int) ([]order_menu.ClientGetOftenList, error) {
	return uu.orderFood.ClientGetOftenList(ctx, branchID)
}

// @waiter

func (uu UseCase) WaiterUpdateOrderMenuStatus(ctx context.Context, ids []int64, status string) error {
	return uu.orderFood.WaiterUpdateStatus(ctx, ids, status)
}

// @cashier

func (uu UseCase) CashierUpdateOrderMenuStatus(ctx context.Context, id int64, status string) error {
	return uu.orderFood.CashierUpdateStatus(ctx, id, status)
}

//// @client
//
//func (uu UseCase) ClientGetOrderMenuList(ctx context.Context, filter order_food.Filter) ([]order_food.ClientGetList, int, error) {
//	list, count, err := uu.orderFood.ClientGetList(ctx, filter)
//	if err != nil {
//		return nil, 0, err
//	}
//
//	return list, count, err
//}
//
//func (uu UseCase) ClientGetOrderMenuDetail(ctx context.Context, id int64) (order_food.ClientGetBannerDetail, error) {
//	var detail order_food.ClientGetBannerDetail
//
//	data, err := uu.orderFood.ClientGetBannerDetail(ctx, id)
//	if err != nil {
//		return order_food.ClientGetBannerDetail{}, err
//	}
//	detail.ID = data.ID
//	detail.Count = data.Count
//	detail.MenuID = data.MenuID
//	detail.OrderID = data.OrderID
//
//	return detail, nil
//}
//
//func (uu UseCase) ClientCreateOrderFood(ctx context.Context, data order_food.ClientCreateRequest) (order_food.ClientCreateResponse, error) {
//	detail, err := uu.orderFood.ClientCreate(ctx, data)
//	if err != nil {
//		return order_food.ClientCreateResponse{}, err
//	}
//
//	return detail, err
//}
//
//func (uu UseCase) ClientUpdateOrderFood(ctx context.Context, data order_food.ClientUpdateRequest) error {
//	return uu.orderFood.ClientUpdateAll(ctx, data)
//}
//
//func (uu UseCase) ClientUpdateOrderFoodColumn(ctx context.Context, data order_food.ClientUpdateRequest) error {
//	return uu.orderFood.ClientUpdateColumns(ctx, data)
//}
//
//func (uu UseCase) ClientDeleteOrderFood(ctx context.Context, id int64) error {
//	return uu.orderFood.ClientDelete(ctx, id)
//}

// #order_payment

// @client

func (uu UseCase) ClientGetOrderPaymentList(ctx context.Context, filter order_payment.Filter) ([]order_payment.ClientGetList, int, error) {
	list, count, err := uu.orderPayment.ClientGetList(ctx, filter)
	if err != nil {
		return nil, 0, err
	}

	return list, count, err
}

func (uu UseCase) ClientGetOrderPaymentDetail(ctx context.Context, id int64) (order_payment.ClientGetDetail, error) {
	var detail order_payment.ClientGetDetail

	data, err := uu.orderPayment.ClientGetDetail(ctx, id)
	if err != nil {
		return order_payment.ClientGetDetail{}, err
	}
	detail.ID = data.ID
	detail.Status = data.Status
	detail.OrderID = data.OrderID

	return detail, nil
}

func (uu UseCase) ClientCreateOrderPayment(ctx context.Context, data order_payment.ClientCreateRequest) (order_payment.ClientCreateResponse, error) {
	detail, err := uu.orderPayment.ClientCreate(ctx, data)
	if err != nil {
		return order_payment.ClientCreateResponse{}, err
	}

	return detail, err
}

func (uu UseCase) ClientUpdateOrderPayment(ctx context.Context, data order_payment.ClientUpdateRequest) error {
	return uu.orderPayment.ClientUpdateAll(ctx, data)
}

func (uu UseCase) ClientUpdateOrderPaymentColumn(ctx context.Context, data order_payment.ClientUpdateRequest) error {
	return uu.orderPayment.ClientUpdateColumns(ctx, data)
}

func (uu UseCase) ClientDeleteOrderPayment(ctx context.Context, id int64) error {
	return uu.orderPayment.ClientDelete(ctx, id)
}

// @waiter

func (uu UseCase) WaiterGetOrderList(ctx context.Context, filter order.Filter) ([]order.WaiterGetListResponse, int, error) {
	return uu.order.WaiterGetList(ctx, filter)
}

func (uu UseCase) WaiterCreateOrder(ctx context.Context, data order.WaiterCreateRequest) (*order.WaiterCreateResponse, error) {
	response, err := uu.order.WaiterCreate(ctx, data)
	if err != nil {
		return nil, err
	}
	printerMessage, branchID, err := uu.order.GetWsOrderMenus(ctx, response.ID, nil)
	if err != nil {
		return nil, err
	}
	for _, v := range printerMessage {
		wsResponse := socket.Response{
			Action:   "new_food",
			Property: v,
		}
		messageJson, _ := json.Marshal(&wsResponse)
		uu.hub.Broadcast <- socket.Message{
			BranchID: branchID,
			Action:   "printer_new_food",
			Message:  messageJson,
		}
	}

	return response, err
}

func (uu UseCase) WaiterGetOrderDetail(ctx context.Context, id int64) (*order.WaiterGetDetailResponse, error) {
	return uu.order.WaiterGetDetail(ctx, id)
}

func (uu UseCase) WaiterUpdateOrder(ctx context.Context, request order.WaiterUpdateRequest) error {
	err := uu.order.WaiterUpdate(ctx, request)
	if err != nil {
		return err
	}

	printerMessage, branchID, err := uu.order.GetWsOrderMenus(ctx, request.Id, request.Menus)
	if err != nil {
		return err
	}
	for _, v := range printerMessage {
		wsResponse := socket.Response{
			Action:   "new_food",
			Property: v,
		}
		messageJson, _ := json.Marshal(&wsResponse)
		uu.hub.Broadcast <- socket.Message{
			BranchID: branchID,
			Action:   "printer_new_food",
			Message:  messageJson,
		}
	}

	return nil
}

func (uu UseCase) WaiterUpdateOrderStatus(ctx context.Context, id int64, status string) error {
	return uu.order.WaiterUpdateStatus(ctx, id, status)
}

func (uu UseCase) WaiterAcceptOrder(ctx context.Context, id int64) error {
	message, err := uu.order.WaiterAccept(ctx, id)
	if err != nil {
		return err
	}

	printerMessage, branchID, err := uu.order.GetWsOrderMenus(ctx, id, nil)
	if err != nil {
		return err
	}

	for _, v := range printerMessage {
		wsResponse := socket.Response{
			Action:   "new_food",
			Property: v,
		}
		messageJson, _ := json.Marshal(&wsResponse)
		uu.hub.Broadcast <- socket.Message{
			BranchID: branchID,
			Action:   "printer_new_food",
			Message:  messageJson,
		}
	}

	// message broadcast to user as accepted
	go func() {
		wsResponse := socket.Response{
			Action: "accepted_order",
			Property: socket.OrderAcceptedResponse{
				OrderID:     &id,
				OrderNumber: message.OrderNumber,
				TableID:     message.TableID,
				TableNumber: message.TableNumber,
				BranchID:    &branchID,
				ClientID:    message.ClientID,
				WaiterID:    message.WaiterID,
				WaiterName:  message.WaiterName,
				WaiterPhoto: message.WaiterPhoto,
				AcceptedAt:  message.AcceptedAt,
				Status:      "ACCEPTED",
			},
		}

		if message.ClientID != nil {
			messageJson, _ := json.Marshal(&wsResponse)
			uu.hub.Broadcast <- socket.Message{
				UserID:  *message.ClientID,
				Action:  "client_accepted_order",
				Message: messageJson,
			}
		}
	}()

	//message broadcast to waiters that order accepted by other waiter
	go func() {
		wsResponse := socket.Response{
			Action: "accepted",
			Property: socket.OrderAcceptedResponse{
				OrderID:     &id,
				OrderNumber: message.OrderNumber,
				TableID:     message.TableID,
				TableNumber: message.TableNumber,
				ClientID:    message.ClientID,
				WaiterID:    message.WaiterID,
				WaiterName:  message.WaiterName,
				WaiterPhoto: message.WaiterPhoto,
				AcceptedAt:  message.AcceptedAt,
				Status:      "ACCEPTED",
			},
		}

		if message.WaiterID != nil {
			messageJson, _ := json.Marshal(&wsResponse)
			uu.hub.Broadcast <- socket.Message{
				UserID:   *message.WaiterID,
				BranchID: branchID,
				Action:   "waiter_accepted_order",
				Message:  messageJson,
			}
		}
	}()

	return nil
}

func (uu UseCase) WaiterGetMyOrderDetail(ctx context.Context, id int64) (*order.WaiterGetOrderDetailResponse, error) {
	return uu.order.WaiterGetMyOrderDetail(ctx, id)
}

func sortWaiter(request []order.GetWsWaiterResponse) []order.GetWsWaiterResponse {
	response := make([]order.GetWsWaiterResponse, 0)
	for k, req := range request {
		for k2, res := range response {
			if req.Grade < res.Grade {
				response = append(append(response[:k2], req), response[k2:]...)
				break
			}
			if k2 == len(response)-1 {
				response = append(response, req)
			}
		}
		if k == 0 {
			response = append(response, req)
		}
	}
	return response
}

// #order_report

// @cashier

func (uu UseCase) CashierOrderReport(ctx context.Context) error {
	return uu.orderReport.CashierOrderReport(ctx)
}
