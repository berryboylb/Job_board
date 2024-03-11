package redis

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
	"github.com/joho/godotenv"
	"log"
	"os"
	"strconv"
	"time"
)

var client *redis.Client
var ctx = context.Background()

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Println("Error loading .env file")
	}

	redisPassword := os.Getenv("REDIS_PASSWORD")
	redisDB := os.Getenv("REDIS_DB")
	redisHost := os.Getenv("REDIS_HOST")

	if redisPassword == "" || redisDB == "" || redisHost == "" {
		panic("Error loading redis password, db or host")
	}

	hostNumber, err := strconv.Atoi(redisHost)
	if err != nil {
		panic("Error converting redis host to int")
	}
	client = redis.NewClient(&redis.Options{
		Addr:     redisHost,
		Password: redisPassword, // no password set
		DB:       hostNumber,    // use default DB
	})
}

func GetClient() *redis.Client {
	return client
}

func Test() {
	ping, err := client.Ping(ctx).Result()
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	fmt.Println(ping)
}

func Store(key string, value interface{}, storageTime int) error {
	var expiration time.Duration
	if storageTime != 0 {
		expiration = time.Duration(storageTime) * time.Minute
	} else {
		expiration = 10 * time.Minute
	}

	err := client.Set(ctx, key, value, expiration).Err()
	if err != nil {
		fmt.Printf("failed to set value in redis instance: %v\n", err)
		return err
	}
	return nil
}

func Retrieve(key string) (interface{}, error) {
	result, err := client.Get(ctx, key).Result()
	if err != nil {
		if err == redis.Nil {
			return nil, fmt.Errorf("%v does not exist", key)
		}
		return nil, fmt.Errorf("failed to get value for key %v from Redis: %v", key, err)
	}
	return result, nil
}
