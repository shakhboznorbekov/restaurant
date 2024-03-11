package attendance

import (
	"context"
)

type Repository interface {

	// @waiter

	WaiterCameCreate(ctx context.Context) (WaiterCreateResponse, error)
	WaiterGoneCreate(ctx context.Context) (WaiterCreateResponse, error)
}
