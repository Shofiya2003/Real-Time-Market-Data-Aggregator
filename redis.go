package main

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/redis/go-redis/v9"
)

var RedisClient *redis.Client

func ConnectToRedis() {
	client := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	RedisClient = client
}

func StoreCurrency(client *redis.Client, newCurrency Currency) error {
	ctx := context.Background()
	name := newCurrency.Name

	marshalledCurrency, err := json.Marshal(newCurrency)

	if err != nil {
		fmt.Println("marshall error: ", err)
	}

	if RedisClient == nil {
		fmt.Println("redis client not initialized")
		return fmt.Errorf("redis client not initialized")
	}

	return RedisClient.MSet(ctx, name, marshalledCurrency).Err()

}

func GetCurrency(client *redis.Client, name string) (Currency, error) {
	var curr Currency
	ctx := context.Background()
	res, err := client.Get(ctx, name).Result()

	if err != nil {
		fmt.Println("error in retrieving currency from redis")
		return curr, err
	}

	err = json.Unmarshal([]byte(res), &curr)

	if err != nil {
		fmt.Println("error in unmarshalling currency")
		return curr, err
	}

	return curr, nil
}
