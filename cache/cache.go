package cache

import (
	"github.com/go-redis/redis/v8"
	"context"
	"encoding/json"
	"log"
	"your_project/models"
)

var rdb *redis.Client
var ctx = context.Background()

func InitRedis() {
	rdb = redis.NewClient(&redis.Options{
		Addr: "localhost:6379", 
	})
}

func GetProductCache(productID int) (models.Product, bool) {
	var product models.Product
	cacheKey := fmt.Sprintf("product:%d", productID)

	val, err := rdb.Get(ctx, cacheKey).Result()
	if err != nil {
		return product, false
	}

	err = json.Unmarshal([]byte(val), &product)
	if err != nil {
		log.Println("Error unmarshalling cached product:", err)
		return product, false
	}

	return product, true
}

func SetProductCache(productID int, product models.Product) {
	cacheKey := fmt.Sprintf("product:%d", productID)
	productBytes, err := json.Marshal(product)
	if err != nil {
		log.Println("Error marshalling product:", err)
		return
	}

	err = rdb.Set(ctx, cacheKey, productBytes, 0).Err()
	if err != nil {
		log.Println("Error setting product in cache:", err)
	}
}


