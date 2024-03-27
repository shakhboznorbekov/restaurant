package user

import (
	"fmt"
	"github.com/restaurant/foundation/web"
	"github.com/restaurant/internal/pkg/config"
	user_service "github.com/restaurant/internal/service/user"
	"github.com/restaurant/internal/service/waiter"
	"github.com/restaurant/internal/usecase/user"
	"net/http"
	"reflect"
	"strconv"
)

type Controller struct {
	useCase *user.UseCase
}

func NewController(useCase *user.UseCase) *Controller {
	return &Controller{useCase}
}

// #user

// @super-Admin

func (uc Controller) SuperAdminGetUserList(c *web.Context) error {
	var filter user_service.Filter

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

	list, count, err := uc.useCase.SuperAdminGetUserList(c.Ctx, filter)
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

func (uc Controller) SuperAdminGetUserDetail(c *web.Context) error {
	id := c.GetParam(reflect.Int, "id").(int)

	if err := c.ValidParam(); err != nil {
		return c.RespondError(err)
	}

	response, err := uc.useCase.SuperAdminGetUserDetail(c.Ctx, int64(id))
	if err != nil {
		return c.RespondError(err)
	}

	return c.Respond(map[string]interface{}{
		"data":   response,
		"status": true,
	}, http.StatusOK)
}

func (uc Controller) SuperAdminCreateUser(c *web.Context) error {
	var request user_service.SuperAdminCreateRequest

	if err := c.BindFunc(&request, "Name", "Phone", "Gender", "BirthDate"); err != nil {
		return c.RespondError(err)
	}

	response, err := uc.useCase.SuperAdminCreateUser(c.Ctx, request)
	if err != nil {
		return c.RespondError(err)
	}

	return c.Respond(map[string]interface{}{
		"data":   response,
		"status": true,
	}, http.StatusOK)
}

func (uc Controller) SuperAdminUpdateUserAll(c *web.Context) error {
	id := c.GetParam(reflect.Int, "id").(int)

	if err := c.ValidParam(); err != nil {
		return c.RespondError(err)
	}

	var request user_service.SuperAdminUpdateRequest

	if err := c.BindFunc(&request, "ID", "Name", "Phone", "Gender", "BirthDate"); err != nil {
		return c.RespondError(err)
	}

	request.ID = int64(id)

	err := uc.useCase.SuperAdminUpdateUser(c.Ctx, request)
	if err != nil {
		return c.RespondError(err)
	}

	return c.Respond(map[string]interface{}{
		"data":   "ok!",
		"status": true,
	}, http.StatusOK)
}

func (uc Controller) SuperAdminUpdateUserColumns(c *web.Context) error {
	id := c.GetParam(reflect.Int, "id").(int)

	if err := c.ValidParam(); err != nil {
		return c.RespondError(err)
	}

	var request user_service.SuperAdminUpdateRequest

	if err := c.BindFunc(&request); err != nil {
		return c.RespondError(err)
	}

	request.ID = int64(id)

	err := uc.useCase.SuperAdminUpdateUserColumn(c.Ctx, request)
	if err != nil {
		return c.RespondError(err)
	}

	return c.Respond(map[string]interface{}{
		"data":   "ok!",
		"status": true,
	}, http.StatusOK)
}

func (uc Controller) SuperAdminDeleteUser(c *web.Context) error {
	id := c.GetParam(reflect.Int, "id").(int)

	if err := c.ValidParam(); err != nil {
		return c.RespondError(err)
	}

	err := uc.useCase.SuperAdminDeleteUser(c.Ctx, int64(id))
	if err != nil {
		return c.RespondError(err)
	}

	return c.Respond(map[string]interface{}{
		"data":   "ok!",
		"status": true,
	}, http.StatusOK)
}

// @client

func (uc Controller) ClientGetUserMe(c *web.Context) error {
	authStr := c.Request.Header.Get("Authorization")

	response, err := uc.useCase.ClientGetUserMe(c.Ctx, authStr)
	if err != nil {
		return c.RespondMobileError(err)
	}

	return c.Respond(map[string]interface{}{
		"data":   response,
		"status": true,
		"error":  nil,
	}, http.StatusOK)
}

func (uc Controller) ClientUpdateUserColumns(c *web.Context) error {
	var request user_service.ClientUpdateRequest

	if err := c.BindFunc(&request); err != nil {
		return c.RespondMobileError(err)
	}

	err := uc.useCase.ClientUpdateUserColumn(c.Ctx, request)
	if err != nil {
		return c.RespondMobileError(err)
	}

	return c.Respond(map[string]interface{}{
		"data":   "ok!",
		"status": true,
		"error":  nil,
	}, http.StatusOK)
}

func (uc Controller) ClientDeleteUserMe(c *web.Context) error {
	err := uc.useCase.ClientDeleteMe(c.Ctx)

	if err != nil {
		return c.RespondMobileError(err)
	}

	return c.Respond(map[string]interface{}{
		"data":   "ok!",
		"status": true,
		"error":  nil,
	}, http.StatusOK)
}

func (uc Controller) SendFeedback(c *web.Context) error {
	var data user_service.FeedBack

	if err := c.BindFunc(&data,
		"Title",
		"Description",
		"FullName",
		"PhoneNumber",
	); err != nil {
		return c.RespondMobileError(err)
	}

	if data.Email == nil {
		email := "-"
		data.Email = &email
	}

	breakLine := "%0A"
	botData := fmt.Sprintf(""+
		"FullName: %s%s"+
		"Email: %s%s"+
		"Title: %s%s"+
		"Description: %s%s"+
		"PhoneNumber: %s",
		*data.FullName, breakLine,
		*data.Email, breakLine,
		*data.Title, breakLine,
		*data.Description, breakLine,
		*data.PhoneNumber)
	url := fmt.Sprintf("https://api.telegram.org/bot%s/sendMessage?chat_id=%s&text=%s",
		config.NewConfig().BOTToken, config.NewConfig().ChatID, botData)

	res, err := http.Get(url)
	if res.StatusCode == 200 {
		return c.Respond(map[string]interface{}{
			"error":  err,
			"status": res.Status,
		}, http.StatusOK)
	} else {
		return c.RespondMobileError(err)
	}
}

// @admin
// #waiter

func (uc Controller) AdminGetWaiterList(c *web.Context) error {
	var filter waiter.Filter

	if limit, ok := c.GetQueryFunc(reflect.Int, "limit").(*int); ok {
		filter.Limit = limit
	}
	if offset, ok := c.GetQueryFunc(reflect.Int, "offset").(*int); ok {
		filter.Offset = offset
	}
	if branchId, ok := c.GetQueryFunc(reflect.Int, "branchId").(int); ok {
		branchID := int64(branchId)
		filter.BranchID = &branchID
	}

	if err := c.ValidQuery(); err != nil {
		return c.RespondError(err)
	}

	list, count, err := uc.useCase.AdminGetWaiterList(c.Ctx, filter)
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

// #cashier

func (uc Controller) AdminGetCashierList(c *web.Context) error {
	var filter cashier.Filter

	if limit, ok := c.GetQueryFunc(reflect.Int, "limit").(*int); ok {
		filter.Limit = limit
	}
	if offset, ok := c.GetQueryFunc(reflect.Int, "offset").(*int); ok {
		filter.Offset = offset
	}

	if err := c.ValidQuery(); err != nil {
		return c.RespondError(err)
	}

	list, count, err := uc.useCase.AdminGetCashierList(c.Ctx, filter)
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

func (uc Controller) AdminGetCashierDetail(c *web.Context) error {
	id := c.GetParam(reflect.Int, "id").(int)

	if err := c.ValidParam(); err != nil {
		return c.RespondError(err)
	}

	response, err := uc.useCase.AdminGetCashierDetail(c.Ctx, int64(id))
	if err != nil {
		return c.RespondError(err)
	}

	return c.Respond(map[string]interface{}{
		"data":   response,
		"status": true,
	}, http.StatusOK)
}

func (uc Controller) AdminCreateCashier(c *web.Context) error {
	var request cashier.AdminCreateRequest

	if err := c.BindFunc(&request, "Name", "Phone", "Gender", "BirthDate", "Password"); err != nil {
		return c.RespondError(err)
	}

	response, err := uc.useCase.AdminCreateCashier(c.Ctx, request)
	if err != nil {
		return c.RespondError(err)
	}

	return c.Respond(map[string]interface{}{
		"data":   response,
		"status": true,
	}, http.StatusOK)
}

func (uc Controller) AdminUpdateCashierAll(c *web.Context) error {
	id := c.GetParam(reflect.Int, "id").(int)

	if err := c.ValidParam(); err != nil {
		return c.RespondError(err)
	}

	var request cashier.AdminUpdateRequest

	if err := c.BindFunc(&request, "Name", "Gender", "BirthDate"); err != nil {
		return c.RespondError(err)
	}

	request.ID = int64(id)

	err := uc.useCase.AdminUpdateCashier(c.Ctx, request)
	if err != nil {
		return c.RespondError(err)
	}

	return c.Respond(map[string]interface{}{
		"data":   "ok!",
		"status": true,
	}, http.StatusOK)
}

func (uc Controller) AdminUpdateCashierColumns(c *web.Context) error {
	id := c.GetParam(reflect.Int, "id").(int)

	if err := c.ValidParam(); err != nil {
		return c.RespondError(err)
	}

	var request cashier.AdminUpdateRequest

	if err := c.BindFunc(&request); err != nil {
		return c.RespondError(err)
	}

	request.ID = int64(id)

	err := uc.useCase.AdminUpdateCashierColumn(c.Ctx, request)
	if err != nil {
		return c.RespondError(err)
	}

	return c.Respond(map[string]interface{}{
		"data":   "ok!",
		"status": true,
	}, http.StatusOK)
}

func (uc Controller) AdminDeleteCashier(c *web.Context) error {
	id := c.GetParam(reflect.Int, "id").(int)

	if err := c.ValidParam(); err != nil {
		return c.RespondError(err)
	}

	err := uc.useCase.AdminDeleteCashier(c.Ctx, int64(id))
	if err != nil {
		return c.RespondError(err)
	}

	return c.Respond(map[string]interface{}{
		"data":   "ok!",
		"status": true,
	}, http.StatusOK)
}

func (uc Controller) AdminUpdateCashierStatus(c *web.Context) error {
	id := c.GetParam(reflect.Int, "id").(int)

	if err := c.ValidParam(); err != nil {
		return c.RespondError(err)
	}

	request := struct {
		Status string `json:"status"`
	}{}

	if err := c.BindFunc(&request); err != nil {
		return c.RespondError(err)
	}

	err := uc.useCase.AdminUpdateCashierStatus(c.Ctx, int64(id), request.Status)
	if err != nil {
		return c.RespondError(err)
	}

	return c.Respond(map[string]interface{}{
		"data":   "ok!",
		"status": true,
	}, http.StatusOK)
}

func (uc Controller) AdminUpdateCashierPassword(c *web.Context) error {
	id := c.GetParam(reflect.Int, "id").(int)

	if err := c.ValidParam(); err != nil {
		return c.RespondError(err)
	}

	var request cashier.AdminUpdatePassword

	if err := c.BindFunc(&request); err != nil {
		return c.RespondError(err)
	}

	request.ID = int64(id)

	err := uc.useCase.AdminUpdateCashierPassword(c.Ctx, request)
	if err != nil {
		return c.RespondError(err)
	}

	return c.Respond(map[string]interface{}{
		"data":   "ok!",
		"status": true,
	}, http.StatusOK)
}

func (uc Controller) AdminUpdateCashierPhone(c *web.Context) error {
	id := c.GetParam(reflect.Int, "id").(int)

	if err := c.ValidParam(); err != nil {
		return c.RespondError(err)
	}

	var request cashier.AdminUpdatePhone

	if err := c.BindFunc(&request); err != nil {
		return c.RespondError(err)
	}

	request.ID = int64(id)

	err := uc.useCase.AdminUpdateCashierPhone(c.Ctx, request)
	if err != nil {
		return c.RespondError(err)
	}

	return c.Respond(map[string]interface{}{
		"data":   "ok!",
		"status": true,
	}, http.StatusOK)
}

// @branch
// #waiter

func (uc Controller) BranchGetListWaiterWorkTime(c *web.Context) error {
	var filter waiter_work_time2.BranchFilter

	if limit, ok := c.GetQueryFunc(reflect.Int, "limit").(*int); ok {
		filter.Limit = limit
	}
	if offset, ok := c.GetQueryFunc(reflect.Int, "offset").(*int); ok {
		filter.Offset = offset
	}
	if page, ok := c.GetQueryFunc(reflect.Int, "page").(*int); ok {
		filter.Page = page
	}
	if date, ok := c.GetQueryFunc(reflect.String, "date").(*string); ok {
		filter.Date = date
	}

	if err := c.ValidQuery(); err != nil {
		return c.RespondError(err)
	}

	list, count, err := uc.useCase.BranchGetListWaiterWorkTime(c.Ctx, filter)
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

func (uc Controller) BranchGetWaiterList(c *web.Context) error {
	var filter waiter.Filter

	if limit, ok := c.GetQueryFunc(reflect.Int, "limit").(*int); ok {
		filter.Limit = limit
	}
	if offset, ok := c.GetQueryFunc(reflect.Int, "offset").(*int); ok {
		filter.Offset = offset
	}
	if branchId, ok := c.GetQueryFunc(reflect.Int, "branchId").(int); ok {
		branchID := int64(branchId)
		filter.BranchID = &branchID
	}

	if err := c.ValidQuery(); err != nil {
		return c.RespondError(err)
	}

	list, count, err := uc.useCase.BranchGetWaiterList(c.Ctx, filter)
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

func (uc Controller) BranchGetWaiterDetail(c *web.Context) error {
	id := c.GetParam(reflect.Int, "id").(int)

	if err := c.ValidParam(); err != nil {
		return c.RespondError(err)
	}

	response, err := uc.useCase.BranchGetWaiterDetail(c.Ctx, int64(id))
	if err != nil {
		return c.RespondError(err)
	}

	return c.Respond(map[string]interface{}{
		"data":   response,
		"status": true,
	}, http.StatusOK)
}

func (uc Controller) BranchCreateWaiter(c *web.Context) error {
	var request waiter.BranchCreateRequest

	if err := c.BindFunc(&request, "Name", "Phone", "Gender", "BirthDate"); err != nil {
		return c.RespondError(err)
	}

	response, err := uc.useCase.BranchCreateWaiter(c.Ctx, request)
	if err != nil {
		return c.RespondError(err)
	}

	return c.Respond(map[string]interface{}{
		"data":   response,
		"status": true,
	}, http.StatusOK)
}

func (uc Controller) BranchUpdateWaiterAll(c *web.Context) error {
	id := c.GetParam(reflect.Int, "id").(int)

	if err := c.ValidParam(); err != nil {
		return c.RespondError(err)
	}

	var request waiter.BranchUpdateRequest

	if err := c.BindFunc(&request, "Name", "Gender", "BirthDate"); err != nil {
		return c.RespondError(err)
	}

	request.ID = int64(id)

	err := uc.useCase.BranchUpdateWaiter(c.Ctx, request)
	if err != nil {
		return c.RespondError(err)
	}

	return c.Respond(map[string]interface{}{
		"data":   "ok!",
		"status": true,
	}, http.StatusOK)
}

func (uc Controller) BranchUpdateWaiterColumns(c *web.Context) error {
	id := c.GetParam(reflect.Int, "id").(int)

	if err := c.ValidParam(); err != nil {
		return c.RespondError(err)
	}

	var request waiter.BranchUpdateRequest

	if err := c.BindFunc(&request); err != nil {
		return c.RespondError(err)
	}

	request.ID = int64(id)

	err := uc.useCase.BranchUpdateWaiterColumn(c.Ctx, request)
	if err != nil {
		return c.RespondError(err)
	}

	return c.Respond(map[string]interface{}{
		"data":   "ok!",
		"status": true,
	}, http.StatusOK)
}

func (uc Controller) BranchDeleteWaiter(c *web.Context) error {
	id := c.GetParam(reflect.Int, "id").(int)

	if err := c.ValidParam(); err != nil {
		return c.RespondError(err)
	}

	err := uc.useCase.BranchDeleteWaiter(c.Ctx, int64(id))
	if err != nil {
		return c.RespondError(err)
	}

	return c.Respond(map[string]interface{}{
		"data":   "ok!",
		"status": true,
	}, http.StatusOK)
}

func (uc Controller) BranchUpdateWaiterPassword(c *web.Context) error {
	id := c.GetParam(reflect.Int, "id").(int)

	if err := c.ValidParam(); err != nil {
		return c.RespondError(err)
	}

	var request waiter.BranchUpdatePassword

	if err := c.BindFunc(&request); err != nil {
		return c.RespondError(err)
	}

	request.ID = int64(id)

	err := uc.useCase.BranchUpdateWaiterPassword(c.Ctx, request)
	if err != nil {
		return c.RespondError(err)
	}

	return c.Respond(map[string]interface{}{
		"data":   "ok!",
		"status": true,
	}, http.StatusOK)
}

func (uc Controller) BranchUpdateWaiterPhone(c *web.Context) error {
	id := c.GetParam(reflect.Int, "id").(int)

	if err := c.ValidParam(); err != nil {
		return c.RespondError(err)
	}

	var request waiter.BranchUpdatePhone

	if err := c.BindFunc(&request); err != nil {
		return c.RespondError(err)
	}

	request.ID = int64(id)

	err := uc.useCase.BranchUpdateWaiterPhone(c.Ctx, request)
	if err != nil {
		return c.RespondError(err)
	}

	return c.Respond(map[string]interface{}{
		"data":   "ok!",
		"status": true,
	}, http.StatusOK)
}

func (uc Controller) BranchUpdateWaiterStatus(c *web.Context) error {
	id := c.GetParam(reflect.Int, "id").(int)

	if err := c.ValidParam(); err != nil {
		return c.RespondError(err)
	}

	request := struct {
		Status string `json:"status"`
	}{}

	if err := c.BindFunc(&request); err != nil {
		return c.RespondError(err)
	}

	err := uc.useCase.BranchUpdateWaiterStatus(c.Ctx, int64(id), request.Status)
	if err != nil {
		return c.RespondError(err)
	}

	return c.Respond(map[string]interface{}{
		"data":   "ok!",
		"status": true,
	}, http.StatusOK)
}

// #cashier

func (uc Controller) BranchGetCashierList(c *web.Context) error {
	var filter cashier.Filter

	if limit, ok := c.GetQueryFunc(reflect.Int, "limit").(*int); ok {
		filter.Limit = limit
	}
	if offset, ok := c.GetQueryFunc(reflect.Int, "offset").(*int); ok {
		filter.Offset = offset
	}

	if err := c.ValidQuery(); err != nil {
		return c.RespondError(err)
	}

	list, count, err := uc.useCase.BranchGetCashierList(c.Ctx, filter)
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

func (uc Controller) BranchGetCashierDetail(c *web.Context) error {
	id := c.GetParam(reflect.Int, "id").(int)

	if err := c.ValidParam(); err != nil {
		return c.RespondError(err)
	}

	response, err := uc.useCase.BranchGetCashierDetail(c.Ctx, int64(id))
	if err != nil {
		return c.RespondError(err)
	}

	return c.Respond(map[string]interface{}{
		"data":   response,
		"status": true,
	}, http.StatusOK)
}

func (uc Controller) BranchCreateCashier(c *web.Context) error {
	var request cashier.BranchCreateRequest

	if err := c.BindFunc(&request, "Password", "Name", "Phone", "Gender", "BirthDate"); err != nil {
		return c.RespondError(err)
	}

	response, err := uc.useCase.BranchCreateCashier(c.Ctx, request)
	if err != nil {
		return c.RespondError(err)
	}

	return c.Respond(map[string]interface{}{
		"data":   response,
		"status": true,
	}, http.StatusOK)
}

func (uc Controller) BranchUpdateCashierAll(c *web.Context) error {
	id := c.GetParam(reflect.Int, "id").(int)

	if err := c.ValidParam(); err != nil {
		return c.RespondError(err)
	}

	var request cashier.BranchUpdateRequest

	if err := c.BindFunc(&request, "Name", "Gender", "BirthDate"); err != nil {
		return c.RespondError(err)
	}

	request.ID = int64(id)

	err := uc.useCase.BranchUpdateCashier(c.Ctx, request)
	if err != nil {
		return c.RespondError(err)
	}

	return c.Respond(map[string]interface{}{
		"data":   "ok!",
		"status": true,
	}, http.StatusOK)
}

func (uc Controller) BranchUpdateCashierColumns(c *web.Context) error {
	id := c.GetParam(reflect.Int, "id").(int)

	if err := c.ValidParam(); err != nil {
		return c.RespondError(err)
	}

	var request cashier.BranchUpdateRequest

	if err := c.BindFunc(&request); err != nil {
		return c.RespondError(err)
	}

	request.ID = int64(id)

	err := uc.useCase.BranchUpdateCashierColumn(c.Ctx, request)
	if err != nil {
		return c.RespondError(err)
	}

	return c.Respond(map[string]interface{}{
		"data":   "ok!",
		"status": true,
	}, http.StatusOK)
}

func (uc Controller) BranchDeleteCashier(c *web.Context) error {
	id := c.GetParam(reflect.Int, "id").(int)

	if err := c.ValidParam(); err != nil {
		return c.RespondError(err)
	}

	err := uc.useCase.BranchDeleteCashier(c.Ctx, int64(id))
	if err != nil {
		return c.RespondError(err)
	}

	return c.Respond(map[string]interface{}{
		"data":   "ok!",
		"status": true,
	}, http.StatusOK)
}

func (uc Controller) BranchUpdateCashierStatus(c *web.Context) error {
	id := c.GetParam(reflect.Int, "id").(int)

	if err := c.ValidParam(); err != nil {
		return c.RespondError(err)
	}

	request := struct {
		Status string `json:"status"`
	}{}

	if err := c.BindFunc(&request); err != nil {
		return c.RespondError(err)
	}

	err := uc.useCase.BranchUpdateCashierStatus(c.Ctx, int64(id), request.Status)
	if err != nil {
		return c.RespondError(err)
	}

	return c.Respond(map[string]interface{}{
		"data":   "ok!",
		"status": true,
	}, http.StatusOK)
}

func (uc Controller) BranchUpdateCashierPassword(c *web.Context) error {
	id := c.GetParam(reflect.Int, "id").(int)

	if err := c.ValidParam(); err != nil {
		return c.RespondError(err)
	}

	var request cashier.BranchUpdatePassword

	if err := c.BindFunc(&request, "Password"); err != nil {
		return c.RespondError(err)
	}

	request.ID = int64(id)

	err := uc.useCase.BranchUpdateCashierPassword(c.Ctx, request)
	if err != nil {
		return c.RespondError(err)
	}

	return c.Respond(map[string]interface{}{
		"data":   "ok!",
		"status": true,
	}, http.StatusOK)
}

func (uc Controller) BranchUpdateCashierPhone(c *web.Context) error {
	id := c.GetParam(reflect.Int, "id").(int)

	if err := c.ValidParam(); err != nil {
		return c.RespondError(err)
	}

	var request cashier.BranchUpdatePhone

	if err := c.BindFunc(&request); err != nil {
		return c.RespondError(err)
	}

	request.ID = int64(id)

	err := uc.useCase.BranchUpdateCashierPhone(c.Ctx, request)
	if err != nil {
		return c.RespondError(err)
	}

	return c.Respond(map[string]interface{}{
		"data":   "ok!",
		"status": true,
	}, http.StatusOK)
}

// @cashier
// #cashier

func (uc Controller) CashierGetMe(c *web.Context) error {
	response, err := uc.useCase.CashierGetMe(c.Ctx)
	if err != nil {
		return c.RespondError(err)
	}

	return c.Respond(map[string]interface{}{
		"data":   response,
		"status": true,
	}, http.StatusOK)
}

// waiter

func (uc Controller) CashierGetWaiterLists(c *web.Context) error {
	var filter waiter.Filter

	if limit, ok := c.GetQueryFunc(reflect.Int, "limit").(*int); ok {
		filter.Limit = limit
	}
	if offset, ok := c.GetQueryFunc(reflect.Int, "offset").(*int); ok {
		filter.Offset = offset
	}
	if branchId, ok := c.GetQueryFunc(reflect.Int, "branch_id").(int); ok {
		branchID := int64(branchId)
		filter.BranchID = &branchID
	}

	if err := c.ValidQuery(); err != nil {
		return c.RespondError(err)
	}

	list, count, err := uc.useCase.CashierGetWaiterLists(c.Ctx, filter)
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

func (uc Controller) CashierGetWaiterDetails(c *web.Context) error {
	id := c.GetParam(reflect.Int, "id").(int)

	if err := c.ValidParam(); err != nil {
		return c.RespondError(err)
	}

	response, err := uc.useCase.CashierGetWaiterDetails(c.Ctx, int64(id))
	if err != nil {
		return c.RespondError(err)
	}

	return c.Respond(map[string]interface{}{
		"data":   response,
		"status": true,
	}, http.StatusOK)
}

func (uc Controller) CashierCreateWaiter(c *web.Context) error {
	var request waiter.CashierCreateRequest

	if err := c.BindFunc(&request, "Name", "Phone", "Gender", "BirthDate"); err != nil {
		return c.RespondError(err)
	}

	response, err := uc.useCase.CashierCreateWaiter(c.Ctx, request)
	if err != nil {
		return c.RespondError(err)
	}

	return c.Respond(map[string]interface{}{
		"data":   response,
		"status": true,
	}, http.StatusOK)
}

func (uc Controller) CashierUpdateWaiterAll(c *web.Context) error {
	id := c.GetParam(reflect.Int, "id").(int)

	if err := c.ValidParam(); err != nil {
		return c.RespondError(err)
	}

	var request waiter.CashierUpdateRequest

	if err := c.BindFunc(&request, "Name", "Gender", "BirthDate"); err != nil {
		return c.RespondError(err)
	}

	request.ID = int64(id)

	err := uc.useCase.CashierUpdateWaiter(c.Ctx, request)
	if err != nil {
		return c.RespondError(err)
	}

	return c.Respond(map[string]interface{}{
		"data":   "ok!",
		"status": true,
	}, http.StatusOK)
}

func (uc Controller) CashierUpdateWaiterColumns(c *web.Context) error {
	id := c.GetParam(reflect.Int, "id").(int)

	if err := c.ValidParam(); err != nil {
		return c.RespondError(err)
	}

	var request waiter.CashierUpdateRequest

	if err := c.BindFunc(&request); err != nil {
		return c.RespondError(err)
	}

	request.ID = int64(id)

	err := uc.useCase.CashierUpdateWaiterColumn(c.Ctx, request)
	if err != nil {
		return c.RespondError(err)
	}

	return c.Respond(map[string]interface{}{
		"data":   "ok!",
		"status": true,
	}, http.StatusOK)
}

func (uc Controller) CashierDeleteWaiter(c *web.Context) error {
	id := c.GetParam(reflect.Int, "id").(int)

	if err := c.ValidParam(); err != nil {
		return c.RespondError(err)
	}

	err := uc.useCase.CashierDeleteWaiter(c.Ctx, int64(id))
	if err != nil {
		return c.RespondError(err)
	}

	return c.Respond(map[string]interface{}{
		"data":   "ok!",
		"status": true,
	}, http.StatusOK)
}

func (uc Controller) CashierUpdateWaiterPassword(c *web.Context) error {
	id := c.GetParam(reflect.Int, "id").(int)

	if err := c.ValidParam(); err != nil {
		return c.RespondError(err)
	}

	var request waiter.CashierUpdatePassword

	if err := c.BindFunc(&request); err != nil {
		return c.RespondError(err)
	}

	request.ID = int64(id)

	err := uc.useCase.CashierUpdateWaiterPassword(c.Ctx, request)
	if err != nil {
		return c.RespondError(err)
	}

	return c.Respond(map[string]interface{}{
		"data":   "ok!",
		"status": true,
	}, http.StatusOK)
}

func (uc Controller) CashierUpdateWaiterPhone(c *web.Context) error {
	id := c.GetParam(reflect.Int, "id").(int)

	if err := c.ValidParam(); err != nil {
		return c.RespondError(err)
	}

	var request waiter.CashierUpdatePhone

	if err := c.BindFunc(&request); err != nil {
		return c.RespondError(err)
	}

	request.ID = int64(id)

	err := uc.useCase.CashierUpdateWaiterPhone(c.Ctx, request)
	if err != nil {
		return c.RespondError(err)
	}

	return c.Respond(map[string]interface{}{
		"data":   "ok!",
		"status": true,
	}, http.StatusOK)
}

func (uc Controller) CashierUpdateWaiterStatus(c *web.Context) error {
	id := c.GetParam(reflect.Int, "id").(int)

	if err := c.ValidParam(); err != nil {
		return c.RespondError(err)
	}

	request := struct {
		Status string `json:"status"`
	}{}

	if err := c.BindFunc(&request); err != nil {
		return c.RespondError(err)
	}

	err := uc.useCase.CashierUpdateWaiterStatus(c.Ctx, int64(id), request.Status)
	if err != nil {
		return c.RespondError(err)
	}

	return c.Respond(map[string]interface{}{
		"data":   "ok!",
		"status": true,
	}, http.StatusOK)
}

func (uc Controller) CashierGetListWaiterWorkTime(c *web.Context) error {
	var filter waiter_work_time2.BranchFilter

	if limit, ok := c.GetQueryFunc(reflect.Int, "limit").(*int); ok {
		filter.Limit = limit
	}
	if offset, ok := c.GetQueryFunc(reflect.Int, "offset").(*int); ok {
		filter.Offset = offset
	}
	if page, ok := c.GetQueryFunc(reflect.Int, "page").(*int); ok {
		filter.Page = page
	}
	if date, ok := c.GetQueryFunc(reflect.String, "date").(*string); ok {
		filter.Date = date
	}

	if err := c.ValidQuery(); err != nil {
		return c.RespondError(err)
	}

	list, count, err := uc.useCase.CashierGetListWaiterWorkTime(c.Ctx, filter)
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

func (uc Controller) CashierGetDetailWaiterWorkTime(c *web.Context) error {
	var filter waiter_work_time2.ListFilter

	if limit, ok := c.GetQueryFunc(reflect.Int, "limit").(*int); ok {
		filter.Limit = limit
	}
	if offset, ok := c.GetQueryFunc(reflect.Int, "offset").(*int); ok {
		filter.Offset = offset
	}
	if page, ok := c.GetQueryFunc(reflect.Int, "page").(*int); ok {
		filter.Page = page
	}

	waiterID := c.GetParam(reflect.String, "waiter_id").(string)

	if err := c.ValidQuery(); err != nil {
		return c.RespondError(err)
	}

	waiterInt, _ := strconv.Atoi(waiterID)

	filter.WaiterID = &waiterInt

	list, count, err := uc.useCase.CashierGetDetailWaiterWorkTime(c.Ctx, filter)
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

func (uc Controller) CashierGetWaiterList(c *web.Context) error {
	var filter waiter.Filter

	if limit, ok := c.GetQueryFunc(reflect.Int, "limit").(*int); ok {
		filter.Limit = limit
	}
	if offset, ok := c.GetQueryFunc(reflect.Int, "offset").(*int); ok {
		filter.Offset = offset
	}

	if err := c.ValidQuery(); err != nil {
		return c.RespondError(err)
	}

	list, count, err := uc.useCase.CashierGetWaiterList(c.Ctx, filter)
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

func (uc Controller) WaiterGetMe(c *web.Context) error {
	response, err := uc.useCase.WaiterGetMe(c.Ctx)
	if err != nil {
		return c.RespondError(err)
	}

	return c.Respond(map[string]interface{}{
		"data":   response,
		"status": true,
	}, http.StatusOK)
}

func (uc Controller) WaiterGetPersonalInfo(c *web.Context) error {
	response, err := uc.useCase.WaiterGetPersonalInfo(c.Ctx)
	if err != nil {
		return c.RespondError(err)
	}

	return c.Respond(map[string]interface{}{
		"data":   response,
		"status": true,
	}, http.StatusOK)
}

func (uc Controller) WaiterUpdatePhoto(c *web.Context) error {

	if err := c.ValidParam(); err != nil {
		return c.RespondError(err)
	}

	var request waiter.WaiterPhotoUpdateRequest

	if err := c.BindFunc(&request, "Photo"); err != nil {
		return c.RespondError(err)
	}

	err := uc.useCase.WaiterUpdatePhoto(c.Ctx, request)

	if err != nil {
		return c.RespondError(err)
	}

	return c.Respond(map[string]interface{}{
		"data":   "ok!",
		"status": true,
	}, http.StatusOK)
}

func (uc Controller) BranchGetDetailWaiterWorkTime(c *web.Context) error {
	var filter waiter_work_time2.ListFilter

	if limit, ok := c.GetQueryFunc(reflect.Int, "limit").(*int); ok {
		filter.Limit = limit
	}
	if offset, ok := c.GetQueryFunc(reflect.Int, "offset").(*int); ok {
		filter.Offset = offset
	}
	if page, ok := c.GetQueryFunc(reflect.Int, "page").(*int); ok {
		filter.Page = page
	}

	waiterID := c.GetParam(reflect.String, "waiter_id").(string)

	if err := c.ValidQuery(); err != nil {
		return c.RespondError(err)
	}

	waiterInt, _ := strconv.Atoi(waiterID)

	filter.WaiterID = &waiterInt

	list, count, err := uc.useCase.BranchGetDetailWaiterWorkTime(c.Ctx, filter)
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

func (uc Controller) WaiterGetActivityStatistics(c *web.Context) error {
	response, err := uc.useCase.WaiterGetActivityStatistics(c.Ctx)
	if err != nil {
		return c.RespondError(err)
	}

	return c.Respond(map[string]interface{}{
		"data":   response,
		"status": true,
	}, http.StatusOK)
}

func (uc Controller) WaiterGetWeeklyActivityStatistics(c *web.Context) error {
	var filter waiter.EarnedMoneyFilter

	if date, ok := c.GetQueryFunc(reflect.String, "date").(*string); ok {
		filter.Date = date
	}
	if err := c.ValidQuery(); err != nil {
		return c.RespondError(err)
	}

	response, err := uc.useCase.WaiterGetWeeklyActivityStatistics(c.Ctx, filter)
	if err != nil {
		return c.RespondError(err)
	}

	return c.Respond(map[string]interface{}{
		"data":   response,
		"status": true,
	}, http.StatusOK)
}

func (uc Controller) WaiterGetWeeklyAcceptedOrdersStatistics(c *web.Context) error {
	var filter waiter.EarnedMoneyFilter

	if date, ok := c.GetQueryFunc(reflect.String, "date").(*string); ok {
		filter.Date = date
	}
	if err := c.ValidQuery(); err != nil {
		return c.RespondError(err)
	}

	response, err := uc.useCase.WaiterGetWeeklyAcceptedOrdersStatistics(c.Ctx, filter)
	if err != nil {
		return c.RespondError(err)
	}

	return c.Respond(map[string]interface{}{
		"data":   response,
		"status": true,
	}, http.StatusOK)
}

func (uc Controller) WaiterCreateComeAttendance(c *web.Context) error {
	if err := c.ValidParam(); err != nil {
		return c.RespondError(err)
	}

	err := uc.useCase.WaiterCreateComeAttendance(c.Ctx)

	if err != nil {
		return c.RespondError(err)
	}

	return c.Respond(map[string]interface{}{
		"data":   "ok!",
		"status": true,
		"error":  nil,
	}, http.StatusOK)
}

func (uc Controller) WaiterCreateGoneAttendance(c *web.Context) error {
	if err := c.ValidParam(); err != nil {
		return c.RespondError(err)
	}

	err := uc.useCase.WaiterCreateGoneAttendance(c.Ctx)

	if err != nil {
		return c.RespondError(err)
	}

	return c.Respond(map[string]interface{}{
		"data":   "ok!",
		"status": true,
		"error":  nil,
	}, http.StatusOK)
}

func (uc Controller) WaiterGetListWorkTime(c *web.Context) error {
	var filter waiter_work_time2.ListFilter

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
		return c.RespondError(err)
	}

	list, count, err := uc.useCase.WaiterGetListWorkTime(c.Ctx, filter)
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

func (uc Controller) WaiterGetWeeklyRatingStatistics(c *web.Context) error {
	var filter waiter.Filter

	if err := c.ValidQuery(); err != nil {
		return c.RespondError(err)
	}

	response, err := uc.useCase.WaiterGetWeeklyRatingStatistics(c.Ctx, filter)
	if err != nil {
		return c.RespondError(err)
	}

	return c.Respond(map[string]interface{}{
		"data":   response,
		"status": true,
	}, http.StatusOK)
}

func (uc Controller) GetMe(c *web.Context) error {
	authStr := c.Request.Header.Get("Authorization")

	response, err := uc.useCase.GetMe(c.Ctx, authStr)
	if err != nil {
		return c.RespondError(err)
	}

	return c.Respond(map[string]interface{}{
		"data":   response,
		"status": true,
	}, http.StatusOK)
}
