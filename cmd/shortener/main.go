package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"runtime"
	"runtime/pprof"
	"syscall"

	_ "net/http/pprof"

	"github.com/avGenie/url-shortener/internal/app/config"
	handlers "github.com/avGenie/url-shortener/internal/app/handlers/router"
	"github.com/avGenie/url-shortener/internal/app/logger"
	storage "github.com/avGenie/url-shortener/internal/app/storage/api"
	"github.com/avGenie/url-shortener/internal/app/storage/api/model"
	"go.uber.org/zap"
)

func main() {
	config := config.InitConfig()

	err := logger.Initialize(config)
	if err != nil {
		panic(err.Error())
	}

	sugar := *zap.S()
	defer sugar.Sync()

	storage, err := storage.InitStorage(config)
	if err != nil {
		sugar.Fatalw(
			err.Error(),
			"event", "init storage",
		)
	}
	defer storage.Close()

	sugar.Infow(
		"Starting server",
		"addr", config.NetAddr,
	)

	startHTTPServer(config, storage)
}

func startHTTPServer(config config.Config, storage model.Storage) {
	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGTERM, os.Interrupt)
	defer cancel()

	router := handlers.NewRouter(config, storage)

	server := &http.Server{
		Addr:    config.NetAddr,
		Handler: router.Mux,
	}

	go func() {
		err := http.ListenAndServe(config.NetAddr, router.Mux)
		if err != nil && err != http.ErrServerClosed {
			zap.L().Fatal("fatal error while starting server", zap.Error(err))
		}
	}()

	<-ctx.Done()

	fmem, err := os.Create(`result.pprof`)
	if err != nil {
		panic(err)
	}
	defer fmem.Close()
	runtime.GC()
	if err := pprof.WriteHeapProfile(fmem); err != nil {
		panic(err)
	}

	zap.L().Info("Got interruption signal. Shutting down HTTP server gracefully...")
	err = server.Shutdown(context.Background())
	if err != nil {
		zap.L().Error("error while shutting down server", zap.Error(err))
	}
	router.Stop()
}
