package domain

import (
	"sync"
	"time"
)

type LockableCache struct {
	LastUpdated time.Time
	Value       interface{}
}

func CleanCache(c map[string]*LockableCache, mutex *sync.RWMutex, maxAge time.Duration) {
	for {
		time.Sleep(5 * time.Minute)
		mutex.Lock()
		for k, v := range c {
			if time.Since(v.LastUpdated) > maxAge {
				delete(c, k)
			}
		}
		mutex.Unlock()
	}
}
