package device

import (
	"context"
	"github.com/restaurant/internal/entity"
)

type Repository interface {
	Create(ctx context.Context, data Create) (entity.Device, error)
	Update(ctx context.Context, data Update) error
	List(ctx context.Context, filter Filter) ([]entity.Device, int, error)
	Detail(ctx context.Context, id int64) (entity.Device, error)
	Delete(ctx context.Context, id string) error
	ChangeDeviceLang(ctx context.Context, data ChangeDeviceLang) error
}
