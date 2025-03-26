package cache

import (
	"sync"
	"time"
)

type LockableCache struct {
	LastUpdated time.Time
	Value       interface{}
}

func CleanCache(c map[string]*LockableCache, mutex *sync.RWMutex, maxAge time.Duration) {
	// Dynamically determine the sleep time based on the maxAge
	sleep := maxAge / 10
	if sleep < time.Minute {
		sleep = time.Minute
	}
	for {
		time.Sleep(sleep)
		mutex.Lock()
		for k, v := range c {
			if time.Since(v.LastUpdated) > maxAge {
				delete(c, k)
			}
		}
		mutex.Unlock()
	}
}
