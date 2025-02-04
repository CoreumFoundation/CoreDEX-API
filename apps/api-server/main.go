package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/CoreumFoundation/CoreDEX-API/apps/api-server/app"
	"github.com/CoreumFoundation/CoreDEX-API/apps/api-server/ports/http"
	"github.com/CoreumFoundation/CoreDEX-API/utils/logger"
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

	application := app.NewApplication()
	server := http.NewHttpServer(application)
	go application.StartUpdater(ctx)
	go func() {
		err := server.Start(ctx)
		if err != nil {
			logger.Errorf(err.Error())
		}
	}()
	<-ctx.Done()
}
