package restaurant

import (
	"context"
	"github.com/restaurant/internal/entity"
)

type Service struct {
	repo Repository
}

func (s Service) SuperAdminUpdateRestaurantAdmin(ctx context.Context, request SuperAdminUpdateRestaurantAdmin) error {
	return s.repo.SuperAdminUpdateRestaurantAdmin(ctx, request)
}

func (s Service) SuperAdminGetList(ctx context.Context, filter Filter) ([]SuperAdminGetList, int, error) {
	return s.repo.SuperAdminGetList(ctx, filter)
}

func (s Service) SuperAdminGetDetail(ctx context.Context, id int64) (entity.Restaurant, error) {
	return s.repo.SuperAdminGetDetail(ctx, id)
}

func (s Service) SuperAdminCreate(ctx context.Context, request SuperAdminCreateRequest) (SuperAdminCreateResponse, error) {
	return s.repo.SuperAdminCreate(ctx, request)
}

func (s Service) SuperAdminUpdateAll(ctx context.Context, request SuperAdminUpdateRequest) error {
	return s.repo.SuperAdminUpdateAll(ctx, request)
}

func (s Service) SuperAdminUpdateColumns(ctx context.Context, request SuperAdminUpdateRequest) error {
	return s.repo.SuperAdminUpdateColumns(ctx, request)
}

func (s Service) SuperAdminDelete(ctx context.Context, id int64) error {
	return s.repo.SuperAdminDelete(ctx, id)
}

func (s Service) SiteGetList(ctx context.Context) ([]SiteGetListResponse, int, error) {
	return s.repo.SiteGetList(ctx)
}

func NewService(repo Repository) *Service {
	return &Service{repo: repo}
}
