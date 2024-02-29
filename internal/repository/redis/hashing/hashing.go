package hashing

import (
	"context"
	"github.com/redis/go-redis/v9"
	"time"
)

type Repository struct {
	rdb *redis.Client
}

func NewRepository(rdb *redis.Client) *Repository {
	return &Repository{
		rdb: rdb,
	}
}

func (r Repository) SetHashing(ctx context.Context, key string, value interface{}, expiration time.Duration) error {
	return r.rdb.Set(ctx, key, value, expiration).Err()
}

func (r Repository) GetHashing(ctx context.Context, key string) (string, error) {
	return r.rdb.Get(ctx, key).Result()
}
