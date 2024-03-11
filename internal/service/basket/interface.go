package basket

import (
	"context"
)

type Repository interface {
	SetBasket(ctx context.Context, data Create) error
	GetBasket(ctx context.Context, key string) (OrderStore, error)
	UpdateBasket(ctx context.Context, key string, value Update) error
	DeleteBasket(ctx context.Context, key string) error
}
