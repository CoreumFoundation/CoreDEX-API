package main

import (
	"context"
	"encoding/json"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()
	// Load the app config from the env APP_CONFIG
	ac := os.Getenv("APP_CONFIG")
	if ac == "" {
		panic("APP_CONFIG env is not set")
	}
	appConfig := &AppConfig{}
	// Parse the app config
	err := json.Unmarshal([]byte(ac), appConfig)
	if err != nil {
		panic("Can not parse APP_CONFIG. Error: " + err.Error())
	}
	application, err := NewApp(ctx, *appConfig)
	if err != nil {
		panic(err)
	}

	t := time.NewTicker(10 * time.Second)
	for {
		select {
		case <-ctx.Done():
			t.Stop()
			return
		case <-t.C:
			application.CreateOrder(ctx, application.GetAccounts())
		}
	}
}
