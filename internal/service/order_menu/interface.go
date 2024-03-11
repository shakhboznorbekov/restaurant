package order_menu

import (
	"context"
	"github.com/restaurant/internal/entity"
)

type Repository interface {
	// @client

	ClientGetList(ctx context.Context, filter Filter) ([]ClientGetList, int, error)
	ClientGetDetail(ctx context.Context, id int64) (entity.OrderMenu, error)
	ClientCreate(ctx context.Context, request ClientCreateRequest) (ClientCreateResponse, error)
	ClientUpdateAll(ctx context.Context, request ClientUpdateRequest) error
	ClientUpdateColumns(ctx context.Context, request ClientUpdateRequest) error
	ClientDelete(ctx context.Context, id int64) error
	ClientGetOftenList(ctx context.Context, branchID int) ([]ClientGetOftenList, error)

	// @waiter
	WaiterUpdateStatus(ctx context.Context, ids []int64, status string) error
}
