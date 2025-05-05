package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/redis/go-redis/v9"
)


func GetRedisClient(){
	redisConnString := os.Getenv("REDIS_URL")

	if redisConnString == ""{
		log.Fatal("Redis url not found")
	}

	opt,err:= redis.ParseURL(redisConnString)

	if err != nil{
		log.Fatal("error while parsing url",err)
	}

	client := redis.NewClient(opt)

	ctx := context.Background()

	err = client.Set(ctx,"foo","bar",0).Err()

	if err != nil{
		log.Fatal("Error setting values",err)
	}

	val,err := client.Get(ctx,"foo").Result()

	if err != nil{
		log.Fatal("Value not found",err)
	}

	fmt.Println("foo",val)
}