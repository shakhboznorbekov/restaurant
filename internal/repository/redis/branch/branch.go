package branch

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

func (r Repository) SetBranch(ctx context.Context, key string, value interface{}, expiration time.Duration) error {
	return r.rdb.Set(ctx, key, value, expiration).Err()
}

func (r Repository) GetBranch(ctx context.Context, key string) (string, error) {
	return r.rdb.Get(ctx, key).Result()
}
