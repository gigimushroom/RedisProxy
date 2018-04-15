package util

import (
	"time"
	"testing"
	"sync"
)

type fakeHandler struct{}

func (s fakeHandler) Lookup(key string) string {
	return "foundKey"
}

func TestThreadPool_FindKey(t *testing.T) {
	handler := fakeHandler{}
	resp := make(chan string)
	work := Job{
		Request: "findkey",
		JobHandler: handler,
		Resp: resp,
	}
	scheduler := NewScheduler(10, 100)
	go scheduler.Run()
	scheduler.JobQueue <- work
	
	val := <-resp
	if val != "foundKey" {
		t.Errorf("should found the key!")
	}
}

func TestThreadPool_FindKeyConcurrent(t *testing.T) {
	handler := fakeHandler{}
	numberOfJobs := 10
	resp := make(chan string, numberOfJobs)
	work := Job{
		Request: "findkey",
		JobHandler: handler,
		Resp: resp,
	}
	scheduler := NewScheduler(3, 100)
	go scheduler.Run()

	wg := new(sync.WaitGroup)
	for i:=0; i<numberOfJobs; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			scheduler.JobQueue <- work
			time.Sleep(time.Second)
		}()	
	}

	// after group is done, close response channel to avoid deadlock
	go func() {
		wg.Wait()
		close(resp)
	}()

	count := 0
	for _= range resp {
		count++
	}

	scheduler.Stop()
	if count != 10 {
		t.Errorf("should found 10 responses!")
	}
}