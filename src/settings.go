package main

import (
	"sync"
)

type Arguments struct {
	httpPort int
}

var (
	instance *Arguments
	once     sync.Once
	mu       sync.RWMutex
)

func GetSettings() *Arguments {
	once.Do(func() {
		instance = &Arguments{
			httpPort: 8080,
		}
	})
	return instance
}

func UpdateSettings(args Arguments) {
	mu.Lock()
	defer mu.Unlock()

	// Only update non-empty/non-zero values
	if args.httpPort <= 0 {
		instance.httpPort = args.httpPort
	}
}
