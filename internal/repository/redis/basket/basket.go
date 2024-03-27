package basket

import (
	"context"
	"encoding/json"
	"github.com/pkg/errors"
	"github.com/redis/go-redis/v9"
	"github.com/restaurant/internal/service/basket"
	"log"
	"math/rand"
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

func (r Repository) SetBasket(ctx context.Context, data basket.Create) error {
	val, err := json.Marshal(data.Value)
	if err != nil {
		return errors.Wrap(err, "marshal")
	}

	var expiration time.Duration
	if data.Expiration != nil && *data.Expiration != 0 {
		expiration = time.Duration(rand.Int31n(int32(*data.Expiration))) * time.Second
	} else {
		expiration = 0
	}
	return r.rdb.Set(ctx, *data.Key, val, expiration).Err()
}

func (r Repository) GetBasket(ctx context.Context, key string) (basket.OrderStore, error) {
	res, err := r.rdb.Get(ctx, key).Result()
	if err == redis.Nil {
		return basket.OrderStore{}, nil
	} else if err != nil {
		return basket.OrderStore{}, errors.Wrap(err, "get basket")
	}

	var structure basket.OrderStore
	if err = json.Unmarshal([]byte(res), &structure); err != nil {
		return basket.OrderStore{}, errors.Wrap(err, "marshal basket")
	}

	return structure, nil
}

func (r Repository) UpdateBasket(ctx context.Context, key string, value basket.Update) error {

	// Get the current order data from Redis
	val, err := r.rdb.Get(ctx, key).Result()
	if err == redis.Nil {
		log.Println("redis nil - 58")
		foods := make([]basket.Food, 0)
		foods = append(foods, value.Food)
		order := basket.Create{Key: &key, Value: basket.OrderStore{
			Foods:   foods,
			TableID: value.TableID,
			UserID:  value.UserID,
		}}
		return r.SetBasket(ctx, order)
	} else if err != nil {
		return err // Handle the error appropriately
	}

	log.Println("redis not nil - 71")

	// Unmarshal the order data
	var order basket.OrderStore
	if err := json.Unmarshal([]byte(val), &order); err != nil {
		return err // Handle JSON unmarshalling error
	}

	// Check if the food exists and update the count or add a new food
	found := false
	for i, food := range order.Foods {
		if food.ID == value.Food.ID {
			order.Foods[i].Count = value.Food.Count // AdminUpdate the count
			found = true
			break
		}
	}
	if !found {
		order.Foods = append(order.Foods, value.Food) // Add new food
	}

	// Marshal the updated order and set it back in Redis
	return r.SetBasket(ctx, basket.Create{Key: &key, Value: order})
}

func (r Repository) DeleteBasket(ctx context.Context, key string) error {
	return r.rdb.Del(ctx, key).Err()
}
