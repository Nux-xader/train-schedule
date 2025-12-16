package main

import (
	"sync"
	"time"

	"github.com/patrickmn/go-cache"
)

var c = cache.New(1*time.Minute, 10*time.Minute)
var mutexes = make(map[string]*sync.Mutex)

func GetMutexForKey(key string) *sync.Mutex {
	if _, exists := mutexes[key]; !exists {
		mutexes[key] = &sync.Mutex{}
	}

	go cleanUpMutexes(key)
	return mutexes[key]
}

func cleanUpMutexes(key string) {
	time.Sleep(10 * time.Minute)

	mu := mutexes[key]
	if mu != nil {
		mu.Lock()
		defer mu.Unlock()

		delete(mutexes, key)
	}
}
