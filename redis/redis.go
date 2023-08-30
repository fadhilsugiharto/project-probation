// pkg/redis/redis.go

package redis

import (
	"context"
	"encoding/json"
	"project-probation/consts"

	"project-probation/model"

	"github.com/go-redis/redis/v8"
)

func GetProductsFromRedis(ctx context.Context, redisKey string) ([]model.Product, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     consts.RedisHost + ":" + consts.RedisPort,
		Password: "",
		DB:       0,
	})

	val, err := client.Get(ctx, redisKey).Result()
	if err != nil {
		return nil, err
	}

	var products []model.Product
	if err := json.Unmarshal([]byte(val), &products); err != nil {
		return nil, err
	}

	return products, nil
}

func StoreProductsInRedis(ctx context.Context, redisKey string, products []model.Product) error {
	client := redis.NewClient(&redis.Options{
		Addr:     consts.RedisHost + ":" + consts.RedisPort,
		Password: "",
		DB:       0,
	})

	productsJSON, err := json.Marshal(products)
	if err != nil {
		return err
	}

	err = client.Set(ctx, redisKey, productsJSON, 0).Err()
	if err != nil {
		return err
	}

	return nil
}
