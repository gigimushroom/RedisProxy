package main

import (
	"flag"
	"log"
	"net/http"
	"os"
	"redisproxy/infra"
	"redisproxy/service"
	"redisproxy/util"
)

var redisAddr = flag.String("redisIpAndPort", "172.17.0.1:6379", "Redis Ip and Port. Default is localhost:6379")
var expiryTime = flag.Int64("expiry", 10, "Proxy cache global expiry in second")
var capacity = flag.Int("capacity", 100, "Proxy cache capacity")
var port = flag.String("port", "10000", "Server Port.")
var maxConcurrentJob = flag.Int("concurrentmax", 1000, "Max number of concurrent connections allowed")
var numOfWorkers = flag.Int("workers", 100, "Max number of requests can be executed in parallel")


func main() {
	flag.Parse()

	cache, err := service.NewProxyCache(*capacity, *expiryTime)
	if err != nil {
		log.Println("Error when creating proxy cache.")
	}

	log.Println("Proxy Service starts on port", *port, "with cache size", *capacity, "expiry", *expiryTime,
		"in sec. Concurrent clients allowed is", *maxConcurrentJob,
		"Max number of requests can be executed in parallel is", *numOfWorkers)
	
	redisHandler, err := infra.NewRedisHandler(*redisAddr)
	if err != nil {
		log.Println("Failed to connect to redis")
		os.Exit(1)
	}
	log.Println("Successfully connected to Redis", *redisAddr)

	scheduler := util.NewScheduler(*numOfWorkers, *maxConcurrentJob)
	srv := service.ProxyService{cache, redisHandler, scheduler}
	go srv.Scheduler.Run()

	http.HandleFunc("/", srv.GetHandler)

	log.Fatal(http.ListenAndServe(":"+*port, nil))
}
