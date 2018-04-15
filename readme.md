# High-level architecture overview.
A Redis Proxy Service implements as an HTTP web Service.
The system built with thread pool and LRU Cache, and support concurrent processing.

## What the code does.
Service/ProxyService is the HTTP web server to take GET request. 

Each request, service checks if found in cache and not expired.

If yes, update the cache entry's lastUpdated time, move up the cache entry to prevent evict, return the result to client.
If not found, search in Redis, if found, store in cache.

Util contains a implementation of golang like thread pool:

Each job will send to a Job Queue, aiming to process by order.

WorkerPool contains a configurable number of workers to handler jobs.

Worker always listen to incoming job, and run forver until receives shutdown notification

A JobScheduler to listen incoming job, find an avaiable worker from pool, assign the job to the worker. Run forever until receives shutdown notification.

Infra contains the redisHandler low level implementation.

main.go is the entry point.

## Unit tests:
cache_test covered: Cached GET, Gloabl expiry, LRU eviction, Fixed key size.
threadpool_test: Single GET request. Concurrent requests(number of requests are more than the available workers)

## Algorithmic complexity of the cache operations.
O(1). Lru Cache internally uses map to store.


## Instructions for how to run the proxy and tests.
### Build:
For system only contains: make, docker, docker-compose, Bash:
make test
(It download docker images, get all dependencies, build golang project inside the container, run golang unit tests)

### Run:
docker-compose up

If you want just want to run the proxy in docker, and mannualy point to a redis backend:
docker run -p 10000:10000 myproxy -redisIpAndPort=127.0.0.1:6379

### Options

| name | descr |
|---|---|
| `redisIpAndPort` | Redis Ip and Port. Default is localhost:6379 |
| `expiry` | Proxy cache global expiry in second (defaults to 10 sec) |
| `capacity` | Proxy cache capacity (defaults to 100) |
| `port` | Server Port (defaults to 10000 ) |
| `concurrentmax` | Max number of concurrent connections allowed |
| `workers` | Max number of requests can be executed in parallel |


### Test by using curl:
curl -X G http://<IP>:<Port>/<your_word>

### How long you spent on each part of the project.
Coding 2.5h
Unit Testing 0.5h
Docker,makefile, integration testing 1h

### A list of the requirements that you did not implement and the reasons for omitting them.
Redis client protocol - Prefer http way. Easy to use, test, no integration effort.



