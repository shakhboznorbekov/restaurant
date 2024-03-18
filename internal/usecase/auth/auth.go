package auth

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/pkg/errors"
	"github.com/restaurant/foundation/web"
	"github.com/restaurant/internal/auth"
	"github.com/restaurant/internal/commands"
	"github.com/restaurant/internal/service/device"
	"github.com/restaurant/internal/service/sms"
	"github.com/restaurant/internal/service/user"
	"golang.org/x/crypto/bcrypt"
	"net/http"
)

type UseCase struct {
	user   User
	sms    Sms
	device Device
	auth   *auth.Auth
}

func NewUseCase(user User, sms Sms, device Device, auth *auth.Auth) *UseCase {
	return &UseCase{user, sms, device, auth}
}

// @client

func (au UseCase) SignInClient(ctx context.Context, request auth.SignInClient) (string, bool, error) {
	exists, err := au.sms.CheckSMSCode(ctx, sms.Check{
		Phone: request.Phone,
		Code:  request.SMSCode,
	})
	if err != nil {
		return "", false, err
	}

	var token string
	var isNew bool
	if exists {
		isNew = false
		detail, err := au.user.GetByPhone(ctx, request.Phone)
		if err != nil {
			return "", isNew, err
		}
		if detail.Phone == nil {
			return "", isNew, errors.New("phone not found: getByPhone")
		}

		if detail.Name == nil {
			isNew = true
		}

		request.Device.UserID = &detail.ID
		_, err = au.device.Create(ctx, request.Device)
		if err != nil {
			return "", isNew, err
		}

		token, err = commands.GenToken(auth.ClaimsAuth{
			ID:    detail.ID,
			Roles: *detail.Role,
		}, "./private.pem")
		if err != nil {
			return "", isNew, err
		}
	} else {
		isNew = true
		detail, err := au.user.ClientCreate(ctx, user.ClientCreateRequest{Phone: &request.Phone})
		if err != nil {
			return "", isNew, err
		}

		request.Device.UserID = &detail.ID
		_, err = au.device.Create(ctx, request.Device)
		if err != nil {
			return "", isNew, err
		}

		token, err = commands.GenToken(auth.ClaimsAuth{
			ID:    detail.ID,
			Roles: *detail.Role,
		}, "./private.pem")
		if err != nil {
			return "", isNew, err
		}
	}

	return token, isNew, nil
}

func (au UseCase) ClientFillUp(ctx context.Context, request auth.SignUpRequest) error {
	return au.user.ClientUpdateAll(ctx, user.ClientUpdateRequest{
		Name:      request.Name,
		BirthDate: request.BirthDate,
		Gender:    request.Gender,
	})
}

func (au UseCase) ClientLogOut(ctx context.Context, deviceID string) error {
	return au.device.Delete(ctx, deviceID)
}

func (au UseCase) ClientUpdateMePhone(ctx context.Context, request auth.UpdatePhone) error {
	exists, err := au.sms.CheckSMSCodeUpdatePhone(ctx, sms.Check{
		Phone: request.Phone,
		Code:  request.SMSCode,
	})
	if err != nil {
		return err
	}

	if exists {
		_, err = au.user.GetByPhone(ctx, request.Phone)
		if err == nil {
			return web.NewRequestError(errors.New("this phone was created before"), http.StatusBadRequest)
		} else if !errors.Is(err, sql.ErrNoRows) {
			return web.NewRequestError(errors.Wrap(err, "updating phone"), http.StatusInternalServerError)
		}
		err = au.user.ClientUpdateMePhone(ctx, request.Phone)
		if err != nil {
			return err
		}
	} else {
		return web.NewRequestError(errors.New("sms code not verified"), http.StatusBadRequest)
	}

	return nil
}

// @super-admin - admin - branch - cashier

func (au UseCase) SignIn(ctx context.Context, request auth.SignInRequest) (string, error) {
	exists, err := au.user.IsSABCPhoneExists(ctx, request.Phone)
	if err != nil {
		return "", err
	}

	var token string
	if exists {
		detail, err := au.user.GetByPhone(ctx, request.Phone)
		if err != nil {
			return "", err
		}

		if (*detail.Role != auth.RoleSuperAdmin) &&
			(*detail.Role != auth.RoleAdmin) &&
			(*detail.Role != auth.RoleBranch) &&
			(*detail.Role != auth.RoleCashier) {
			errStr := fmt.Sprintf("you have not permission as %s", *detail.Role)
			return "", errors.New(errStr)
		}

		if *detail.Role == auth.RoleAdmin && detail.RestaurantID == nil {
			return "", errors.New(fmt.Sprintf("role %s doesn't contain restaurant_id", *detail.Role))
		}

		if (*detail.Role == auth.RoleBranch ||
			*detail.Role == auth.RoleCashier) &&
			detail.BranchID == nil {
			return "", errors.New(fmt.Sprintf("role %s doesn't contain branch_id", *detail.Role))
		}

		if detail.Password == nil {
			return "", errors.New("user doesn't contain password")
		}

		if err = bcrypt.CompareHashAndPassword([]byte(*detail.Password), []byte(request.Password)); err != nil {
			return "", errors.New(fmt.Sprintf("incorrect password. error: %v", err))
		}

		token, err = commands.GenToken(auth.ClaimsAuth{
			ID:           detail.ID,
			Roles:        *detail.Role,
			RestaurantID: detail.RestaurantID,
			BranchID:     detail.BranchID,
		}, "./private.pem")
		if err != nil {
			return "", err
		}
	} else {
		return "", errors.New("user does not exists")
	}

	return token, nil
}

// @waiter

func (au UseCase) SignInWaiter(ctx context.Context, request auth.SignInWaiter) (string, error) {
	exists, err := au.user.IsWaiterPhoneExists(ctx, request.Phone)
	if err != nil {
		return "", err
	}

	var token string
	if !exists {
		return "", errors.New("phone not found")
	} else {
		detail, err := au.user.GetByPhone(ctx, request.Phone)
		if err != nil {
			return "", err
		}

		if detail.Password == nil {
			return "", errors.New("user doesn't contain password")
		}

		if err = bcrypt.CompareHashAndPassword([]byte(*detail.Password), []byte(request.Password)); err != nil {
			return "", errors.New(fmt.Sprintf("incorrect password. error: %v", err))
		}

		request.Device.UserID = &detail.ID
		_, err = au.device.Create(ctx, request.Device)
		if err != nil {
			return "", err
		}

		token, err = commands.GenToken(auth.ClaimsAuth{
			ID:       detail.ID,
			Roles:    *detail.Role,
			BranchID: detail.BranchID,
		}, "./private.pem")
		if err != nil {
			return "", err
		}
	}

	return token, nil
}

func (au UseCase) WaiterUpdateMePhone(ctx context.Context, request auth.UpdatePhone) error {
	exists, err := au.sms.CheckSMSCodeUpdatePhone(ctx, sms.Check{
		Phone: request.Phone,
		Code:  request.SMSCode,
	})
	if err != nil {
		return err
	}

	if exists {
		_, err = au.user.GetByPhone(ctx, request.Phone)
		if err == nil {
			return web.NewRequestError(errors.New("this phone was created before"), http.StatusBadRequest)
		} else if !errors.Is(err, sql.ErrNoRows) {
			return web.NewRequestError(errors.Wrap(err, "updating phone"), http.StatusInternalServerError)
		}
		err = au.user.WaiterUpdateMePhone(ctx, request.Phone)
		if err != nil {
			return err
		}
	} else {
		return web.NewRequestError(errors.New("sms code not verified"), http.StatusBadRequest)
	}

	return nil
}

func (au UseCase) WaiterUpdatePassword(ctx context.Context, request user.UpdatePasswordRequest) error {
	exists, err := au.sms.CheckSMSCodeUpdatePhone(ctx, sms.Check{
		Phone: request.Phone,
		Code:  request.SMSCode,
	})
	if err != nil {
		return err
	}

	if exists {
		detail, err := au.user.GetByPhone(ctx, request.Phone)
		if err != nil {
			return web.NewRequestError(err, http.StatusBadRequest)
		}

		password, err := bcrypt.GenerateFromPassword([]byte(request.Password), bcrypt.DefaultCost)
		err = au.user.WaiterUpdatePassword(ctx, string(password), detail.ID)
		if err != nil {
			return err
		}
	} else {
		return web.NewRequestError(errors.New("sms code not verified"), http.StatusBadRequest)
	}

	return nil
}

// @others

func (au UseCase) ChangeDeviceLang(ctx context.Context, request device.ChangeDeviceLang) error {
	return au.device.ChangeDeviceLang(ctx, request)
}

func (au UseCase) SendSms(ctx context.Context, request auth.SendSms) error {
	return au.sms.SendSMS(ctx, sms.Send{
		Phone:   request.Phone,
		SmsType: 1,
	})
}
