package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/buaazp/fasthttprouter"
	"github.com/go-redis/redis"
	"github.com/joho/godotenv"
	"github.com/valyala/fasthttp"
)

func main() {
	log.Println("Loading .env")
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	log.Println("Connecting to redis")
	client := redis.NewClient(&redis.Options{
		Addr:     os.Getenv("REDIS_ADDR"),
		Password: os.Getenv("REDIS_PASS"),
		DB:       0,
	})

	_, err = client.Ping().Result()
	if err != nil {
		fmt.Printf("Couldn't connect to redis: '%s'", err)
	}

	log.Println("Starting cron thread")
	go redisTick(client)

	log.Println("Opening web service")
	router := fasthttprouter.New()

	log.Fatal(fasthttp.ListenAndServe(":8080", router.Handler))
}

func redisTick(redis *redis.Client) {
	for true {
		time.Sleep(100)
	}
}
