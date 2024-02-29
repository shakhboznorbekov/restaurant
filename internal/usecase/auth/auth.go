package auth

import (
	"context"
	"fmt"
	"github.com/pkg/errors"
	"github.com/restaurant/internal/auth"
	"github.com/restaurant/internal/commands"
	"github.com/restaurant/internal/service/sms"
	"github.com/restaurant/internal/service/user"
	"golang.org/x/crypto/bcrypt"
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

func (au UseCase) SignInClient(ctx context.Context, request auth.SignInRequest) (string, bool, error) {
	exist, err := au.sms.CheckSMSCode(ctx, sms.Check{
		Phone: request.Phone,
		Code:  request.SMSCode,
	})
	if err != nil {
		return "", false, err
	}

	var token string
	var isNew bool
	if exist {
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
		_, err := au.user.Create(ctx, request.Device)
		if err != nil {
			return "", isNew, err
		}

		token, err := commands.GenToken(auth.ClaimsAuth{
			ID:    detail.ID,
			Roles: *detail.Role,
		}, "private.pem")
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
		_, err := au.device.Create(ctx, request.Device)
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

func (au UseCase) SignInWaiter(ctx context.Context, request auth.SignInWaiter) (string, error) {
	exist, err := au.user.IsWaiterPhoneExists(ctx, request.Phone)
	if err != nil {
		return "", err
	}

	var token string
	if !exist {
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

		request.Device, UserID = &detail.ID
		_, err := au.device.Create(ctx, request.Device)
		if err != nil {
			return "", err
		}

		token, err = commands.GenToken(auth.ClaimsAuth{
			ID:       detail.ID,
			Roles:    *detail.Role,
			BranchID: detail.BranchID,
		}, "private.pem")
		if err != nil {
			return "", err
		}
	}
	return token, nil
}

func (au UseCase) ClientFillUp(ctx context.Context, request auth.SignUpRequest) error {
	claims, err := au.auth.GetTokenData(request.Token)
	if err != nil {
		return err
	}

	return au.user.ClientUpdateAll(ctx, user.ClientUpdateRequest{
		ID:        claims.UserId,
		Name:      request.Name,
		BirthDate: request.BirthDate,
		Gender:    request.Gender,
	})
}

func (au UseCase) SendSms(ctx context.Context, request auth.SendSms) error {
	return au.sms.SendSMS(ctx, sms.Send{
		Phone:   request.Phone,
		SmsType: 1,
	})
}

func (au UseCase) SignInAdmin(ctx context.Context, request auth.SignInRequest) (string, error) {
	exist, err := au.sms.CheckSMSCode(ctx, sms.Check{
		Phone: request.Phone,
		Code:  request.SMSCode,
	})
	if err != nil {
		return "", err
	}

	var token string
	if exist {

	}
}
