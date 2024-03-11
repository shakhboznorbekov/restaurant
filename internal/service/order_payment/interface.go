package order_payment

import (
	"context"
	"github.com/restaurant/internal/entity"
)

type Repository interface {
	// @client
	ClientGetList(ctx context.Context, filter Filter) ([]ClientGetList, int, error)
	ClientGetDetail(ctx context.Context, id int64) (entity.OrderPayment, error)
	ClientCreate(ctx context.Context, request ClientCreateRequest) (ClientCreateResponse, error)
	ClientUpdateAll(ctx context.Context, request ClientUpdateRequest) error
	ClientUpdateColumns(ctx context.Context, request ClientUpdateRequest) error
	ClientDelete(ctx context.Context, id int64) error

	// @cashier
	CashierGetList(ctx context.Context, filter Filter) ([]CashierGetList, int, error)
	CashierGetDetail(ctx context.Context, id int64) (entity.OrderPayment, error)
	CashierCreate(ctx context.Context, request CashierCreateRequest) (CashierCreateResponse, error)
	CashierUpdateAll(ctx context.Context, request CashierUpdateRequest) error
	CashierUpdateColumns(ctx context.Context, request CashierUpdateRequest) error
	CashierDelete(ctx context.Context, id int64) error
}
