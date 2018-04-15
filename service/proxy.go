package service

import (
	"fmt"
    "log"
	"net/http"
	"redisproxy/util"
	"redisproxy/infra"
)

type ProxyService struct {
	Cache *ProxyCache
	RedisHandler *infra.RedisStorager
	Scheduler *util.Scheduler
}

func (s *ProxyService) Get(key string) string {
	found, v := s.Cache.GetIfNotExpired(key)

	if found {
		return v
	} else {
		resp := make(chan string)
		defer close(resp)
		work := util.Job{
			Request: key,
			JobHandler: s.RedisHandler,
			Resp: resp,
		}
		s.Scheduler.JobQueue <- work
		
		v = <-resp
		if v != "" {
			s.Cache.Add(key, v)
		}		
	}
	return v
}

func (s * ProxyService) GetHandler(w http.ResponseWriter, r *http.Request) {
	key := r.URL.Path[1:]
	if key == "" {
		log.Println("[ProxyService.GetHandler]: Key is empty")
		return
	}
	log.Println("[ProxyService.GetHandler]: Searching for key", key)
	v := s.Get(key)
	fmt.Fprintf(w, "%s\n", v)
}