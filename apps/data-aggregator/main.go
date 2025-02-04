package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/CoreumFoundation/CoreDEX-API/apps/data-aggregator/app"
	"github.com/CoreumFoundation/CoreDEX-API/apps/data-aggregator/ports/http"
)

func main() {
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go func() {
		<-c
		cancel()
	}()

	l := app.NewApplication(ctx)
	go l.StartOHLCProcessor(ctx)
	go l.StartScanners(ctx)
	go http.NewListener() // Provide a liveness probe
	<-ctx.Done()
}
