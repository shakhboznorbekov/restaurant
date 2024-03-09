package device

import (
	"context"
	"github.com/restaurant/internal/entity"
)

type Service struct {
	repo Repository
}

func (s Service) Create(ctx context.Context, data Create) (entity.Device, error) {
	return s.repo.Create(ctx, data)
}

func (s Service) Update(ctx context.Context, data Update) error {
	return s.repo.Update(ctx, data)
}

func (s Service) List(ctx context.Context, filter Filter) ([]entity.Device, int, error) {
	return s.repo.List(ctx, filter)
}

func (s Service) Detail(ctx context.Context, id int64) (entity.Device, error) {
	return s.repo.Detail(ctx, id)
}

func (s Service) Delete(ctx context.Context, id string) error {
	return s.repo.Delete(ctx, id)
}

func (s Service) ChangeDeviceLang(ctx context.Context, data ChangeDeviceLang) error {
	return s.repo.ChangeDeviceLang(ctx, data)
}

func NewService(repo Repository) *Service {
	return &Service{repo}
}
