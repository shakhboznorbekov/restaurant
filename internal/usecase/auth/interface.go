package auth

import (
	"context"
	"github.com/restaurant/internal/entity"
	"github.com/restaurant/internal/service/device"
	"github.com/restaurant/internal/service/sms"
	"github.com/restaurant/internal/service/user"
)

type User interface {
	ClientCreate(ctx context.Context, request user.ClientCreateRequest) (user.ClientCreateResponse, error)
	ClientUpdateAll(ctx context.Context, request user.ClientUpdateRequest) error
	GetByPhone(ctx context.Context, phone string) (entity.User, error)
	ClientUpdateMePhone(ctx context.Context, newPhone string) error
	IsPhoneExists(ctx context.Context, phone string) (bool, error)
	IsWaiterPhoneExists(ctx context.Context, phone string) (bool, error)

	//watier

	WaiterUpdateMePhone(ctx context.Context, newPhone string) error
}

type Sms interface {
	SendSMS(ctx context.Context, send sms.Send) error
	CheckSMSCode(ctx context.Context, check sms.Check) (bool, error)
	CheckSMSCodeUpdatePhone(ctx context.Context, check sms.Check) (bool, error)
}

type Device interface {
	Create(ctx context.Context, data device.Create) (entity.Device, error)
	Update(ctx context.Context, data device.Update) error
	List(ctx context.Context, filter device.Filter) ([]entity.Device, int, error)
	Detail(ctx context.Context, id int64) (entity.Device, error)
	Delete(ctx context.Context, id string) error
	ChangeDeviceLang(ctx context.Context, data device.ChangeDeviceLang) error
}
