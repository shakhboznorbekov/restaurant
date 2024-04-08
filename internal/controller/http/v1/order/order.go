package order

import (
	"fmt"
	"github.com/restaurant/foundation/web"
	"net/http"
	"reflect"
)

type Controller struct {
	useCase *order.UseCase
}

func NewController(useCase *order.UseCase) *Controller {
	return &Controller{useCase}
}

// #order-----------------------------------------------------------------------------------------

// @client

func (uc Controller) MobileCreateOrder(c *web.Context) error {
	var (
		body order2.MobileCreateRequest
	)

	if err := c.BindFunc(&body, "TableID", "Menus"); err != nil {
		return c.RespondMobileError(err)
	}

	data, err := uc.useCase.MobileCreateOrder(c.Ctx, body)
	if err != nil {
		return c.RespondMobileError(err)
	}

	return c.Respond(map[string]interface{}{
		"error":  nil,
		"data":   data,
		"status": true,
	}, http.StatusOK)
}

func (uc Controller) MobileUpdateOrder(c *web.Context) error {
	id := c.GetParam(reflect.Int, "id").(int)

	if err := c.ValidParam(); err != nil {
		return c.RespondMobileError(err)
	}

	var request order2.MobileUpdateRequest

	if err := c.BindFunc(&request); err != nil {
		return c.RespondMobileError(err)
	}

	request.Id = int64(id)

	err := uc.useCase.MobileUpdateOrder(c.Ctx, request)
	if err != nil {
		return c.RespondMobileError(err)
	}

	return c.Respond(map[string]interface{}{
		"data":   "updated!",
		"status": true,
		"error":  nil,
	}, http.StatusOK)
}

func (uc Controller) MobileGetOrderList(c *web.Context) error {
	var filter order2.Filter

	if limit, ok := c.GetQueryFunc(reflect.Int, "limit").(*int); ok {
		filter.Limit = limit
	}
	if offset, ok := c.GetQueryFunc(reflect.Int, "offset").(*int); ok {
		filter.Offset = offset
	}
	if page, ok := c.GetQueryFunc(reflect.Int, "page").(*int); ok {
		filter.Page = page
	}

	if err := c.ValidQuery(); err != nil {
		return c.RespondMobileError(err)
	}

	list, count, err := uc.useCase.MobileGetOrderList(c.Ctx, filter)
	if err != nil {
		return c.RespondMobileError(err)
	}

	return c.Respond(map[string]interface{}{
		"data": map[string]interface{}{
			"results": list,
			"count":   count,
		},
		"status": true,
		"error":  nil,
	}, http.StatusOK)
}

func (uc Controller) MobileGetOrderDetail(c *web.Context) error {
	id := c.GetParam(reflect.Int, "id").(int)

	if err := c.ValidParam(); err != nil {
		return c.RespondMobileError(err)
	}

	response, err := uc.useCase.MobileGetOrderDetail(c.Ctx, int64(id))
	if err != nil {
		return c.RespondMobileError(err)
	}

	return c.Respond(map[string]interface{}{
		"data":   response,
		"status": true,
		"error":  nil,
	}, http.StatusOK)
}

func (uc Controller) MobileWaiterCall(c *web.Context) error {
	id := c.GetParam(reflect.Int, "id").(int)

	if err := c.ValidParam(); err != nil {
		return c.RespondMobileError(err)
	}

	err := uc.useCase.MobileWaiterCall(c.Ctx, int64(id))
	if err != nil {
		return c.RespondMobileError(err)
	}

	return c.Respond(map[string]interface{}{
		"data":   "ok!",
		"status": true,
		"error":  nil,
	}, http.StatusOK)
}

func (uc Controller) ClientReviewOrder(c *web.Context) error {
	var (
		body order2.ClientReviewRequest
	)

	if err := c.BindFunc(&body, "OrderId", "Star"); err != nil {
		return c.RespondMobileError(err)
	}

	data, err := uc.useCase.ClientReviewOrder(c.Ctx, body)
	if err != nil {
		return c.RespondMobileError(err)
	}

	return c.Respond(map[string]interface{}{
		"data":   data,
		"status": true,
	}, http.StatusOK)
}

// cashier

func (uc Controller) CashierPaymentOrder(c *web.Context) error {
	id := c.GetParam(reflect.Int, "id").(int)

	if err := c.ValidParam(); err != nil {
		return c.RespondError(err)
	}

	var request order2.CashierUpdateStatusRequest

	if err := c.BindFunc(&request); err != nil {
		return c.RespondError(err)
	}

	request.Id = int64(id)
	status := "PAID"
	request.Status = &status

	err := uc.useCase.CashierUpdateOrder(c.Ctx, request)
	if err != nil {
		return c.RespondError(err)
	}

	return c.Respond(map[string]interface{}{
		"data":   "updated!",
		"status": true,
	}, http.StatusOK)
}

func (uc Controller) CashierGetOrderList(c *web.Context) error {
	var filter order2.Filter

	if limit, ok := c.GetQueryFunc(reflect.Int, "limit").(*int); ok {
		filter.Limit = limit
	}
	if offset, ok := c.GetQueryFunc(reflect.Int, "offset").(*int); ok {
		filter.Offset = offset
	}
	if page, ok := c.GetQueryFunc(reflect.Int, "page").(*int); ok {
		filter.Page = page
	}
	if search, ok := c.GetQueryFunc(reflect.String, "search").(*string); ok {
		filter.Search = search
	}

	if err := c.ValidQuery(); err != nil {
		return c.RespondMobileError(err)
	}

	list, count, err := uc.useCase.CashierGetOrderList(c.Ctx, filter)
	if err != nil {
		return c.RespondMobileError(err)
	}

	return c.Respond(map[string]interface{}{
		"data": map[string]interface{}{
			"results": list,
			"count":   count,
		},
		"status": true,
		"error":  nil,
	}, http.StatusOK)
}

func (uc Controller) CashierGetOrderDetail(c *web.Context) error {
	id := c.GetParam(reflect.Int, "id").(int)

	if err := c.ValidParam(); err != nil {
		return c.RespondMobileError(err)
	}

	response, err := uc.useCase.CashierGetOrderDetail(c.Ctx, int64(id))
	if err != nil {
		return c.RespondMobileError(err)
	}

	return c.Respond(map[string]interface{}{
		"data":   response,
		"status": true,
		"error":  nil,
	}, http.StatusOK)
}

// @admin

func (uc Controller) AdminGetListOrder(c *web.Context) error {
	var filter order2.Filter

	if limit, ok := c.GetQueryFunc(reflect.Int, "limit").(*int); ok {
		filter.Limit = limit
	}
	if offset, ok := c.GetQueryFunc(reflect.Int, "offset").(*int); ok {
		filter.Offset = offset
	}
	if name, ok := c.GetQueryFunc(reflect.String, "name").(*string); ok {
		filter.Name = name
	}

	if err := c.ValidQuery(); err != nil {
		return c.RespondError(err)
	}

	list, count, err := uc.useCase.AdminGetListOrder(c.Ctx, filter)
	if err != nil {
		return c.RespondError(err)
	}

	return c.Respond(map[string]interface{}{
		"data": map[string]interface{}{
			"results": list,
			"count":   count,
		},
		"status": true,
	}, http.StatusOK)
}

// @waiter

func (uc Controller) WaiterGetOrderList(c *web.Context) error {
	var filter order2.Filter

	if limit, ok := c.GetQueryFunc(reflect.Int, "limit").(*int); ok {
		filter.Limit = limit
	}
	if offset, ok := c.GetQueryFunc(reflect.Int, "offset").(*int); ok {
		filter.Offset = offset
	}
	if search, ok := c.GetQueryFunc(reflect.String, "search").(*string); ok {
		filter.Search = search
	}
	if whose, ok := c.GetQueryFunc(reflect.String, "whose").(*string); ok {
		filter.Whose = whose
	}
	if archived, ok := c.GetQueryFunc(reflect.Bool, "archived").(*bool); ok {
		filter.Archived = archived
	}

	if err := c.ValidQuery(); err != nil {
		return c.RespondError(err)
	}

	list, count, err := uc.useCase.WaiterGetOrderList(c.Ctx, filter)
	if err != nil {
		return c.RespondError(err)
	}

	return c.Respond(map[string]interface{}{
		"data": map[string]interface{}{
			"results": list,
			"count":   count,
		},
		"status": true,
	}, http.StatusOK)
}

func (uc Controller) WaiterCreateOrder(c *web.Context) error {
	var (
		body order2.WaiterCreateRequest
	)

	if err := c.BindFunc(&body, "TableID", "Menus"); err != nil {
		return c.RespondMobileError(err)
	}

	data, err := uc.useCase.WaiterCreateOrder(c.Ctx, body)
	if err != nil {
		fmt.Println(err)
		return c.RespondMobileError(err)
	}

	return c.Respond(map[string]interface{}{
		"error":  nil,
		"data":   data,
		"status": true,
	}, http.StatusOK)
}

func (uc Controller) WaiterGetOrderDetail(c *web.Context) error {
	id := c.GetParam(reflect.Int, "id").(int)

	if err := c.ValidParam(); err != nil {
		return c.RespondMobileError(err)
	}

	response, err := uc.useCase.WaiterGetOrderDetail(c.Ctx, int64(id))
	if err != nil {
		return c.RespondMobileError(err)
	}

	return c.Respond(map[string]interface{}{
		"data":   response,
		"status": true,
		"error":  nil,
	}, http.StatusOK)
}

func (uc Controller) WaiterHistoryActivityList(c *web.Context) error {
	var filter order2.Filter

	if limit, ok := c.GetQueryFunc(reflect.Int, "limit").(*int); ok {
		filter.Limit = limit
	}
	if page, ok := c.GetQueryFunc(reflect.Int, "page").(*int); ok {
		filter.Page = page
	}

	if err := c.ValidQuery(); err != nil {
		return c.RespondError(err)
	}

	list, count, err := uc.useCase.WaiterGetHistoryActivityList(c.Ctx, filter)
	if err != nil {
		return c.RespondError(err)
	}

	return c.Respond(map[string]interface{}{
		"data": map[string]interface{}{
			"results": list,
			"count":   count,
		},
		"status": true,
	}, http.StatusOK)
}

func (uc Controller) WaiterGetMyOrderDetail(c *web.Context) error {
	id := c.GetParam(reflect.Int, "id").(int)

	if err := c.ValidParam(); err != nil {
		return c.RespondMobileError(err)
	}

	response, err := uc.useCase.WaiterGetMyOrderDetail(c.Ctx, int64(id))
	if err != nil {
		return c.RespondError(err)
	}

	return c.Respond(map[string]interface{}{
		"data":   response,
		"status": true,
		"error":  nil,
	}, http.StatusOK)
}

// #order-payment-----------------------------------------------------------------------------------------

func (uc Controller) ClientGetOrderPaymentList(c *web.Context) error {
	var filter orderPayment.Filter

	if limit, ok := c.GetQueryFunc(reflect.Int, "limit").(*int); ok {
		filter.Limit = limit
	}
	if offset, ok := c.GetQueryFunc(reflect.Int, "offset").(*int); ok {
		filter.Offset = offset
	}

	if err := c.ValidQuery(); err != nil {
		return c.RespondError(err)
	}

	list, count, err := uc.useCase.ClientGetOrderPaymentList(c.Ctx, filter)
	if err != nil {
		return c.RespondError(err)
	}

	return c.Respond(map[string]interface{}{
		"data": map[string]interface{}{
			"results": list,
			"count":   count,
		},
		"status": true,
	}, http.StatusOK)
}

func (uc Controller) ClientGetOrderPaymentDetail(c *web.Context) error {
	id := c.GetParam(reflect.Int, "id").(int)

	if err := c.ValidParam(); err != nil {
		return c.RespondError(err)
	}

	response, err := uc.useCase.ClientGetOrderPaymentDetail(c.Ctx, int64(id))
	if err != nil {
		return c.RespondError(err)
	}

	return c.Respond(map[string]interface{}{
		"data":   response,
		"status": true,
	}, http.StatusOK)
}

func (uc Controller) ClientCreateOrderPayment(c *web.Context) error {
	var request orderPayment.ClientCreateRequest

	if err := c.BindFunc(&request, "Status", "OrderID"); err != nil {
		return c.RespondError(err)
	}

	response, err := uc.useCase.ClientCreateOrderPayment(c.Ctx, request)
	if err != nil {
		return c.RespondError(err)
	}

	return c.Respond(map[string]interface{}{
		"data":   response,
		"status": true,
	}, http.StatusOK)
}

func (uc Controller) ClientUpdateOrderPaymentAll(c *web.Context) error {
	id := c.GetParam(reflect.Int, "id").(int)

	if err := c.ValidParam(); err != nil {
		return c.RespondError(err)
	}

	var request orderPayment.ClientUpdateRequest

	if err := c.BindFunc(&request, "ID", "Status", "OrderID"); err != nil {
		return c.RespondError(err)
	}

	request.ID = int64(id)

	err := uc.useCase.ClientUpdateOrderPayment(c.Ctx, request)
	if err != nil {
		return c.RespondError(err)
	}

	return c.Respond(map[string]interface{}{
		"data":   "ok!",
		"status": true,
	}, http.StatusOK)
}

func (uc Controller) ClientUpdateOrderPaymentColumns(c *web.Context) error {
	id := c.GetParam(reflect.Int, "id").(int)

	if err := c.ValidParam(); err != nil {
		return c.RespondError(err)
	}

	var request orderPayment.ClientUpdateRequest

	if err := c.BindFunc(&request, "ID"); err != nil {
		return c.RespondError(err)
	}

	request.ID = int64(id)

	err := uc.useCase.ClientUpdateOrderPaymentColumn(c.Ctx, request)
	if err != nil {
		return c.RespondError(err)
	}

	return c.Respond(map[string]interface{}{
		"data":   "ok!",
		"status": true,
	}, http.StatusOK)
}

func (uc Controller) ClientDeleteOrderPayment(c *web.Context) error {
	id := c.GetParam(reflect.Int, "id").(int)

	if err := c.ValidParam(); err != nil {
		return c.RespondError(err)
	}

	err := uc.useCase.ClientDeleteOrderPayment(c.Ctx, int64(id))
	if err != nil {
		return c.RespondError(err)
	}

	return c.Respond(map[string]interface{}{
		"data":   "ok!",
		"status": true,
	}, http.StatusOK)
}

func (uc Controller) WaiterUpdateOrder(c *web.Context) error {
	id := c.GetParam(reflect.Int, "id").(int)

	if err := c.ValidParam(); err != nil {
		return c.RespondMobileError(err)
	}

	var request order2.WaiterUpdateRequest

	if err := c.BindFunc(&request); err != nil {
		return c.RespondMobileError(err)
	}

	request.Id = int64(id)

	err := uc.useCase.WaiterUpdateOrder(c.Ctx, request)
	if err != nil {
		return c.RespondMobileError(err)
	}

	return c.Respond(map[string]interface{}{
		"data":   "updated!",
		"status": true,
		"error":  nil,
	}, http.StatusOK)
}

func (uc Controller) WaiterUpdateOrderStatus(c *web.Context) error {
	id := c.GetParam(reflect.Int, "id").(int)

	if err := c.ValidParam(); err != nil {
		return c.RespondMobileError(err)
	}

	request := struct {
		Status string `json:"status"`
	}{}

	if err := c.BindFunc(&request); err != nil {
		return c.RespondMobileError(err)
	}

	err := uc.useCase.WaiterUpdateOrderStatus(c.Ctx, int64(id), request.Status)
	if err != nil {
		return c.RespondMobileError(err)
	}

	return c.Respond(map[string]interface{}{
		"data":   "updated!",
		"status": true,
		"error":  nil,
	}, http.StatusOK)
}

func (uc Controller) WaiterAcceptOrder(c *web.Context) error {
	id := c.GetParam(reflect.Int, "id").(int)

	if err := c.ValidParam(); err != nil {
		return c.RespondMobileError(err)
	}

	err := uc.useCase.WaiterAcceptOrder(c.Ctx, int64(id))
	if err != nil {
		return c.RespondMobileError(err)
	}

	return c.Respond(map[string]interface{}{
		"data":   "accepted!",
		"status": true,
		"error":  nil,
	}, http.StatusOK)
}

//// #order-food
//
//func (uc Controller) ClientGetOrderMenuList(c *web.Context) error {
//	var filter orderFood.Filter
//
//	if limit, ok := c.GetQueryFunc(reflect.Int, "limit").(*int); ok {
//		filter.Limit = limit
//	}
//	if offset, ok := c.GetQueryFunc(reflect.Int, "offset").(*int); ok {
//		filter.Offset = offset
//	}
//
//	if err := c.ValidQuery(); err != nil {
//		return c.RespondError(err)
//	}
//
//	list, count, err := uc.useCase.ClientGetOrderMenuList(c.Ctx, filter)
//	if err != nil {
//		return c.RespondError(err)
//	}
//
//	return c.Respond(map[string]interface{}{
//		"data": map[string]interface{}{
//			"results": list,
//			"count":   count,
//		},
//		"status": true,
//	}, http.StatusOK)
//}
//
//func (uc Controller) ClientGetOrderMenuDetail(c *web.Context) error {
//	id := c.GetParam(reflect.Int, "id").(int)
//
//	if err := c.ValidParam(); err != nil {
//		return c.RespondError(err)
//	}
//
//	response, err := uc.useCase.ClientGetOrderMenuDetail(c.Ctx, int64(id))
//	if err != nil {
//		return c.RespondError(err)
//	}
//
//	return c.Respond(map[string]interface{}{
//		"data":   response,
//		"status": true,
//	}, http.StatusOK)
//}
//
//func (uc Controller) ClientCreateOrderFood(c *web.Context) error {
//	var request orderFood.ClientCreateRequest
//
//	if err := c.BindFunc(&request, "Count", "MenuID", "OrderID"); err != nil {
//		return c.RespondError(err)
//	}
//
//	response, err := uc.useCase.ClientCreateOrderFood(c.Ctx, request)
//	if err != nil {
//		return c.RespondError(err)
//	}
//
//	return c.Respond(map[string]interface{}{
//		"data":   response,
//		"status": true,
//	}, http.StatusOK)
//}
//
//func (uc Controller) ClientUpdateOrderFoodAll(c *web.Context) error {
//	id := c.GetParam(reflect.Int, "id").(int)
//
//	if err := c.ValidParam(); err != nil {
//		return c.RespondError(err)
//	}
//
//	var request orderFood.ClientUpdateRequest
//
//	if err := c.BindFunc(&request, "ID", "MenuID", "OrderID", "Count"); err != nil {
//		return c.RespondError(err)
//	}
//
//	request.ID = int64(id)
//
//	err := uc.useCase.ClientUpdateOrderFood(c.Ctx, request)
//	if err != nil {
//		return c.RespondError(err)
//	}
//
//	return c.Respond(map[string]interface{}{
//		"data":   "ok!",
//		"status": true,
//	}, http.StatusOK)
//}
//
//func (uc Controller) ClientUpdateOrderFoodColumns(c *web.Context) error {
//	id := c.GetParam(reflect.Int, "id").(int)
//
//	if err := c.ValidParam(); err != nil {
//		return c.RespondError(err)
//	}
//
//	var request orderFood.ClientUpdateRequest
//
//	if err := c.BindFunc(&request, "ID"); err != nil {
//		return c.RespondError(err)
//	}
//
//	request.ID = int64(id)
//
//	err := uc.useCase.ClientUpdateOrderFoodColumn(c.Ctx, request)
//	if err != nil {
//		return c.RespondError(err)
//	}
//
//	return c.Respond(map[string]interface{}{
//		"data":   "ok!",
//		"status": true,
//	}, http.StatusOK)
//}
//
//func (uc Controller) ClientDeleteOrderFood(c *web.Context) error {
//	id := c.GetParam(reflect.Int, "id").(int)
//
//	if err := c.ValidParam(); err != nil {
//		return c.RespondError(err)
//	}
//
//	err := uc.useCase.ClientDeleteOrderFood(c.Ctx, int64(id))
//	if err != nil {
//		return c.RespondError(err)
//	}
//
//	return c.Respond(map[string]interface{}{
//		"data":   "ok!",
//		"status": true,
//	}, http.StatusOK)
//}

// order_menu-----------------------------------------------------------------------------------------------

// @waiter

func (uc Controller) WaiterUpdateOrderMenuStatus(c *web.Context) error {

	request := struct {
		Status string  `json:"status"`
		Ids    []int64 `json:"ids"`
	}{}

	if err := c.BindFunc(&request); err != nil {
		return c.RespondMobileError(err)
	}

	err := uc.useCase.WaiterUpdateOrderMenuStatus(c.Ctx, request.Ids, request.Status)
	if err != nil {
		return c.RespondMobileError(err)
	}

	return c.Respond(map[string]interface{}{
		"data":   "updated!",
		"status": true,
	}, http.StatusOK)
}

// @cashier

func (uc Controller) CashierUpdateOrderMenuStatus(c *web.Context) error {
	request := struct {
		Status string `json:"status"`
		Id     int64  `json:"id"`
	}{}

	if err := c.BindFunc(&request, "Id", "Status"); err != nil {
		return c.RespondMobileError(err)
	}

	err := uc.useCase.CashierUpdateOrderMenuStatus(c.Ctx, request.Id, request.Status)
	if err != nil {
		return c.RespondMobileError(err)
	}

	return c.Respond(map[string]interface{}{
		"data":   "updated!",
		"status": true,
	}, http.StatusOK)
}

// @client

func (uc Controller) MobileClientGetOrderMenuOftenByBranchID(c *web.Context) error {
	id := c.GetParam(reflect.Int, "branch_id").(int)

	if err := c.ValidParam(); err != nil {
		return c.RespondMobileError(err)
	}

	response, err := uc.useCase.ClientGetOrderMenuOftenList(c.Ctx, id)
	if err != nil {
		return c.RespondMobileError(err)
	}

	return c.Respond(map[string]interface{}{
		"data":   response,
		"status": true,
		"error":  nil,
	}, http.StatusOK)
}

// order_report----------------------------------------------------------------------------------------------

// @cashier

func (uc Controller) CashierReportOrder(c *web.Context) error {

	err := uc.useCase.CashierOrderReport(c.Ctx)
	if err != nil {
		return c.RespondMobileError(err)
	}

	return c.Respond(map[string]interface{}{
		"error":  nil,
		"status": true,
	}, http.StatusOK)
}
