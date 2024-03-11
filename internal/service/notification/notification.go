package notification

import (
	"context"
)

type Service struct {
	repo Repository
}

func NewService(repo Repository) *Service {
	return &Service{repo: repo}
}

// @admin

func (s Service) AdminGetList(ctx context.Context, filter Filter) ([]AdminGetList, int, error) {
	return s.repo.AdminGetList(ctx, filter)
}

func (s Service) AdminGetDetail(ctx context.Context, id int64) (*AdminGetDetail, error) {
	return s.repo.AdminGetDetail(ctx, id)
}

func (s Service) AdminUpdateAll(ctx context.Context, request AdminUpdateRequest) error {
	return s.repo.AdminUpdateAll(ctx, request)
}

func (s Service) AdminUpdateColumn(ctx context.Context, request AdminUpdateRequest) error {
	return s.repo.AdminUpdateColumn(ctx, request)
}

func (s Service) AdminCreate(ctx context.Context, request AdminCreateRequest) (*AdminCreateResponse, error) {
	return s.repo.AdminCreate(ctx, request)
}

func (s Service) AdminDelete(ctx context.Context, id int64) error {
	return s.repo.AdminDelete(ctx, id)
}

// @super-admin

func (s Service) SuperAdminUpdateStatus(ctx context.Context, id int64, status string) error {
	return s.repo.SuperAdminUpdateStatus(ctx, id, status)
}

func (s Service) SuperAdminGetList(ctx context.Context, filter Filter) ([]SuperAdminGetList, int, error) {
	return s.repo.SuperAdminGetList(ctx, filter)
}

func (s Service) SuperAdminGetDetail(ctx context.Context, id int64) (*SuperAdminGetDetail, error) {
	return s.repo.SuperAdminGetDetail(ctx, id)
}

func (s Service) SuperAdminSend(ctx context.Context, request SuperAdminSendRequest) ([]SuperAdminSendResponse, error) {
	return s.repo.SuperAdminSend(ctx, request)
}

// @client

func (s Service) ClientGetList(ctx context.Context, filter Filter) ([]ClientGetListResponse, int, error) {
	return s.repo.ClientGetList(ctx, filter)
}

func (s Service) ClientGetUnseenCount(ctx context.Context) (int, error) {
	return s.repo.ClientGetUnseenCount(ctx)
}

func (s Service) ClientSetAsViewed(ctx context.Context, id int64) error {
	return s.repo.ClientSetAsViewed(ctx, id)
}

func (s Service) ClientSetAllAsViewed(ctx context.Context) error {
	return s.repo.ClientSetAllAsViewed(ctx)
}
