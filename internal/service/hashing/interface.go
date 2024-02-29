package hashing

import (
	"context"
	"time"
)

type Repository interface {
	SetHashing(ctx context.Context, key string, value interface{}, expiration time.Duration) error
	GetHashing(ctx context.Context, key string) (string, error)
}
