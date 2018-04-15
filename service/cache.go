package service

import (
	"github.com/hashicorp/golang-lru/simplelru"
    "log"
    "time"
)

type CacheValue struct {
	v string
	lastUpdatedTime int64
}

type ProxyCache struct {
	lru simplelru.LRUCache
	expiry int64
}

func NewProxyCache(cacheSize int, expiredTime int64) (*ProxyCache, error) {
	l, err := simplelru.NewLRU(cacheSize, nil)
	c := &ProxyCache{
		lru : l,
		expiry : expiredTime,
	}
	return c, err
}

func (c *ProxyCache) GetIfNotExpired(key string) (bool, string) {
	// Lru get() will move the entry to front
	if val,ok := c.lru.Get(key); ok {
		// check if expired
		cachedv := val.(CacheValue)
		//log.Println(time.Now().Unix(), cachedv.lastUpdatedTime, c.expiry)
		if (time.Now().Unix() - cachedv.lastUpdatedTime) >= c.expiry {
			// expired
			log.Println("[ProxyCache.GetIfNotExpired]: Key", key, "expired.")
			c.lru.Remove(key)
			return false, ""
		}

		// Update the last updated time
		c.remove(key)
		c.Add(key, cachedv.v)

		log.Println("[ProxyCache.GetIfNotExpired]: Got non-expired Key", key, "in cache. Refreshed lastUpdatedTime.")
		return ok, cachedv.v
	}
	return false, ""
}

func (c *ProxyCache) Add(k, v string) {
	log.Println("[ProxyCache.Add]: Key", k, "value", v, "added in cache")
	t := time.Now().Unix()
	cVal := CacheValue{v, t}
	c.lru.Add(k, cVal)
}

func (c *ProxyCache) remove(key string) {
	c.lru.Remove(key)
}




