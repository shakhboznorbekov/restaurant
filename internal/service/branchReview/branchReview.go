package branchReview

import (
	"context"
)

type Service struct {
	repo Repository
}

// @client

func (s Service) ClientGetList(ctx context.Context, filter Filter) ([]ClientGetList, int, error) {
	return s.repo.ClientGetList(ctx, filter)
}

func (s Service) ClientGetDetail(ctx context.Context, id int64) (ClientGetDetail, error) {
	return s.repo.ClientGetDetail(ctx, id)
}

func (s Service) ClientCreate(ctx context.Context, request ClientCreateRequest) (ClientCreateResponse, error) {
	return s.repo.ClientCreate(ctx, request)
}

func (s Service) ClientUpdateAll(ctx context.Context, request ClientUpdateRequest) error {
	return s.repo.ClientUpdateAll(ctx, request)
}

func (s Service) ClientUpdateColumns(ctx context.Context, request ClientUpdateRequest) error {
	return s.repo.ClientUpdateColumns(ctx, request)
}

func (s Service) ClientDelete(ctx context.Context, id int64) error {
	return s.repo.ClientDelete(ctx, id)
}

func NewService(repo Repository) *Service {
	return &Service{repo: repo}
}
