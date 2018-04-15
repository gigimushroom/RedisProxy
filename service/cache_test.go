package service

import (
	"time"
	"testing"
)

func TestLRU_Evict(t *testing.T) {

	cache, err := NewProxyCache(1, 10) // cachesize, timeout in sec
	if err != nil {
		t.Fatalf("err: %v", err)
	}

	cache.Add("testK1", "testV1")
	cache.Add("testK2", "testV2") // should evict


	result1, _:= cache.GetIfNotExpired("testK1")
	if result1 == true {
		t.Errorf("should have an eviction")
	}

	result2, _:= cache.GetIfNotExpired("testK2")
	if result2 == false {
		t.Errorf("should not have an eviction")
	}
}

func TestLRU_Expiry(t *testing.T) {

	cache, err := NewProxyCache(3, 1) // cachesize, timeout in sec
	if err != nil {
		t.Fatalf("err: %v", err)
	}

	cache.Add("testK1", "testV1")
	cache.Add("testK2", "testV2") // should evict


	result1, _:= cache.GetIfNotExpired("testK1")
	if result1 == false {
		t.Errorf("should not have an eviction")
	}

	result2, _:= cache.GetIfNotExpired("testK2")
	if result2 == false {
		t.Errorf("should not have an eviction")
	}

	time.Sleep(time.Second)
	result2, _= cache.GetIfNotExpired("testK2")
	if result2 == true {
		t.Errorf("After sleep, testK2 should have an eviction")
	}
}
