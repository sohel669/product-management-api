package services

import (
	"context"
	"encoding/json"

	"github.com/go-redis/redis/v8"
	"project/models"
)

var rdb *redis.Client
var ctx = context.Background()

func InitRedis(redisURL string) {
	rdb = redis.NewClient(&redis.Options{
		Addr: redisURL,
	})
}

func CacheProduct(product models.Product) {
	data, _ := json.Marshal(product)
	err := rdb.Set(ctx, strconv.Itoa(product.ID), data, 0).Err()
	if err != nil {
		log.Printf("Failed to cache product: %v", err)
	}
}

func GetProductFromCache(id string) *models.Product {
	data, err := rdb.Get(ctx, id).Result()
	if err != nil {
		return nil
	}

	var product models.Product
	if err := json.Unmarshal([]byte(data), &product); err != nil {
		return nil
	}

	return &product
}
