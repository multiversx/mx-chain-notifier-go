package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/ElrondNetwork/notifier-go/config"
	"github.com/ElrondNetwork/notifier-go/proxy"
)

var (
	backgroundContextTimeout = 5 * time.Second
)

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		panic(err)
	}

	api, err := proxy.NewWebServer(cfg)
	if err != nil {
		panic(err)
	}
	server := api.Run()
	waitForGracefulShutdown(server)
}

func waitForGracefulShutdown(server *http.Server) {
	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt, os.Kill)
	<-quit

	ctx, cancel := context.WithTimeout(context.Background(), backgroundContextTimeout)
	defer cancel()
	if err := server.Shutdown(ctx); err != nil {
		panic(err)
	}
	_ = server.Close()
}
