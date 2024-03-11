package story

import (
	"context"
	"github.com/restaurant/internal/entity"
)

type Service struct {
	repo Repository
}

// @admin

func (s Service) AdminGetList(ctx context.Context, filter Filter) ([]AdminGetList, int, error) {
	return s.repo.AdminGetList(ctx, filter)
}

func (s Service) AdminGetDetail(ctx context.Context, id int64) (entity.Story, error) {
	return s.repo.AdminGetDetail(ctx, id)
}

func (s Service) AdminUpdateStatus(ctx context.Context, id int64) error {
	return s.repo.AdminUpdateStatus(ctx, id)
}

func (s Service) AdminCreate(ctx context.Context, request AdminCreateRequest) (AdminCreateResponse, error) {
	return s.repo.AdminCreate(ctx, request)
}

func (s Service) AdminDelete(ctx context.Context, id int64) error {
	return s.repo.AdminDelete(ctx, id)
}

// @client

func (s Service) ClientGetList(ctx context.Context, filter Filter) ([]ClientGetList, int, error) {
	return s.repo.ClientGetList(ctx, filter)
}

func (s Service) ClientSetViewed(ctx context.Context, id int64) error {
	return s.repo.ClientSetViewed(ctx, id)
}

// @super-admin

func (s Service) SuperAdminGetList(ctx context.Context, filter Filter) ([]SuperAdminGetListResponse, int, error) {
	return s.repo.SuperAdminGetList(ctx, filter)
}

func (s Service) SuperAdminUpdateStatus(ctx context.Context, id int64, status string) error {
	return s.repo.SuperAdminUpdateStatus(ctx, id, status)
}

func NewService(repo Repository) *Service {
	return &Service{repo: repo}
}
