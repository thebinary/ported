package main

import (
	"flag"
	"log"
	"net/http"
	"time"

	"github.com/go-redis/redis/v8"
)

var redisAddr string
var redisDB int
var redisPassword string
var rootKey string
var rootDomain string
var keepAliveTimeout time.Duration

func init() {
	flag.StringVar(&redisAddr, "redisAddr", "localhost:6379", "address of redis being used as dynamic configuration provider")
	flag.IntVar(&redisDB, "redisDB", 0, "redis db number; default: 0")
	flag.StringVar(&redisPassword, "redisPassword", "", "password for redis if used; default: empty (no password)")
	flag.StringVar(&rootKey, "rootKey", "ported", "root key used in redis dynamic configurations")
	flag.StringVar(&rootDomain, "rootDomain", "ported.example.com", "root domain used for porter services")
	flag.DurationVar(&keepAliveTimeout, "keepAliveTimeout", time.Second*60, "keep alive timeout for ported services, services are automatically unregistered if no keep alive messages are recieved from them for this long")

	flag.Parse()

	red = redis.NewClient(&redis.Options{
		Addr:     redisAddr,
		Password: redisPassword,
		DB:       redisDB,
	})

	_, err := red.Ping(ctx).Result()
	if err != nil {
		log.Fatalf("redis unavailable: err=%v", err)
	}
}

func main() {
	http.HandleFunc("/v1/ping", pingHandler)
	http.HandleFunc("/v1/available", availableHandler)
	http.HandleFunc("/v1/service", serviceHandler)
	http.ListenAndServe(":8888", nil)
}
