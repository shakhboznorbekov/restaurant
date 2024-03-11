package branchReview

import (
	"context"
)

type Repository interface {
	// @client
	ClientGetList(ctx context.Context, filter Filter) ([]ClientGetList, int, error)
	ClientGetDetail(ctx context.Context, id int64) (ClientGetDetail, error)
	ClientCreate(ctx context.Context, request ClientCreateRequest) (ClientCreateResponse, error)
	ClientUpdateAll(ctx context.Context, request ClientUpdateRequest) error
	ClientUpdateColumns(ctx context.Context, request ClientUpdateRequest) error
	ClientDelete(ctx context.Context, id int64) error
}
