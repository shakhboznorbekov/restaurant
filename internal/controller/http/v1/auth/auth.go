package auth

import (
	"github.com/restaurant/foundation/web"
	auth2 "github.com/restaurant/internal/auth"
	"github.com/restaurant/internal/service/device"
	"github.com/restaurant/internal/service/user"
	"github.com/restaurant/internal/usecase/auth"
	"net/http"
	"reflect"
)

// build is the git version of this program. It is set using build flags in the makefile.
var build = "develop"

type Controller struct {
	useCase *auth.UseCase
}

func NewController(useCase *auth.UseCase) *Controller {
	return &Controller{useCase}
}

// @client

func (uc Controller) SignInClient(c *web.Context) error {
	var data auth2.SignInClient

	err := c.BindFunc(&data, "Phone", "SMSCode", "Device")
	if err != nil {
		return c.RespondError(err)
	}

	token, isNew, err := uc.useCase.SignInClient(c.Ctx, data)
	if err != nil {
		return c.RespondError(err)
	}

	return c.Respond(map[string]interface{}{
		"new": isNew,
		"data": map[string]string{
			"token": token,
		},
		"error": nil,
	}, http.StatusOK)
}

func (uc Controller) ClientFillUp(c *web.Context) error {
	var data auth2.SignUpRequest

	err := c.BindFunc(&data, "Name", "BirthDate", "Gender")
	if err != nil {
		return c.RespondError(err)
	}

	authStr := c.Request.Header.Get("Authorization")

	data.Token = authStr
	err = uc.useCase.ClientFillUp(c.Ctx, data)
	if err != nil {
		return c.RespondError(err)
	}

	return c.Respond(map[string]interface{}{
		"status": true,
		"error":  nil,
	}, http.StatusOK)
}

func (uc Controller) ClientLogOut(c *web.Context) error {
	id := c.GetParam(reflect.String, "device-id").(string)

	if err := c.ValidParam(); err != nil {
		return c.RespondError(err)
	}

	err := uc.useCase.ClientLogOut(c.Ctx, id)
	if err != nil {
		return c.RespondError(err)
	}

	return c.Respond(map[string]interface{}{
		"data":   "Logged Out!",
		"status": true,
	}, http.StatusOK)
}

func (uc Controller) ClientUpdatePhone(c *web.Context) error {
	var data auth2.UpdatePhone

	err := c.BindFunc(&data, "Phone", "SMSCode", "Gender")
	if err != nil {
		return c.RespondError(err)
	}

	err = uc.useCase.ClientUpdateMePhone(c.Ctx, data)
	if err != nil {
		return c.RespondError(err)
	}

	return c.Respond(map[string]interface{}{
		"status": true,
		"error":  nil,
	}, http.StatusOK)
}

// @super-admin - admin - branch - cashier

func (uc Controller) SignIn(c *web.Context) error {
	var data auth2.SignInRequest

	err := c.BindFunc(&data, "Phone", "SMSCode")
	if err != nil {
		return c.RespondError(err)
	}

	token, err := uc.useCase.SignIn(c.Ctx, data)
	if err != nil {
		return c.RespondError(err)
	}

	return c.Respond(map[string]interface{}{
		"data": map[string]string{
			"token": token,
		},
		"error": nil,
	}, http.StatusOK)
}

// @waiter

func (uc Controller) WaiterSignIn(c *web.Context) error {
	var data auth2.SignInWaiter

	err := c.BindFunc(&data, "Phone", "Password", "Device")
	if err != nil {
		return c.RespondError(err)
	}

	token, err := uc.useCase.SignInWaiter(c.Ctx, data)
	if err != nil {
		return c.RespondError(err)
	}

	return c.Respond(map[string]interface{}{
		"data": map[string]string{
			"token": token,
		},
		"error": nil,
	}, http.StatusOK)
}

func (uc Controller) WaiterUpdatePhone(c *web.Context) error {
	var data auth2.UpdatePhone

	err := c.BindFunc(&data, "Phone", "SMSCode")
	if err != nil {
		return c.RespondError(err)
	}

	err = uc.useCase.WaiterUpdateMePhone(c.Ctx, data)
	if err != nil {
		return c.RespondError(err)
	}

	return c.Respond(map[string]interface{}{
		"status": true,
		"error":  nil,
	}, http.StatusOK)
}

func (uc Controller) WaiterUpdatePassword(c *web.Context) error {

	if err := c.ValidParam(); err != nil {
		return c.RespondError(err)
	}

	var request user.UpdatePasswordRequest

	if err := c.BindFunc(&request, "Password", "SMSCode", "Phone"); err != nil {
		return c.RespondError(err)
	}

	err := uc.useCase.WaiterUpdatePassword(c.Ctx, request)

	if err != nil {
		return c.RespondError(err)
	}

	return c.Respond(map[string]interface{}{
		"data":   "ok!",
		"status": true,
	}, http.StatusOK)
}

// others

func (uc Controller) SendSmsCode(c *web.Context) error {
	var data auth2.SendSms

	err := c.BindFunc(&data, "Phone")
	if err != nil {
		return c.RespondError(err)
	}

	err = uc.useCase.SendSms(c.Ctx, data)
	if err != nil {
		return c.RespondError(err)
	}

	return c.Respond(map[string]interface{}{
		"status": true,
		"error":  nil,
	}, http.StatusOK)
}

func (uc Controller) ChangeDeviceLang(c *web.Context) error {
	var data device.ChangeDeviceLang

	err := c.BindFunc(&data, "DeviceID", "Lang")
	if err != nil {
		return c.RespondError(err)
	}

	err = uc.useCase.ChangeDeviceLang(c.Ctx, data)
	if err != nil {
		return c.RespondError(err)
	}

	return c.Respond(map[string]interface{}{
		"status": true,
		"error":  nil,
	}, http.StatusOK)
}
