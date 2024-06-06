package main

import (
	"context"
	"fmt"
	"net/http"
	_ "net/http/pprof"
	"os"
	"os/signal"
	"runtime"
	"runtime/pprof"
	"syscall"

	"go.uber.org/zap"

	"github.com/avGenie/url-shortener/internal/app/config"
	handlers "github.com/avGenie/url-shortener/internal/app/handlers/router"
	"github.com/avGenie/url-shortener/internal/app/logger"
	storage "github.com/avGenie/url-shortener/internal/app/storage/api"
	"github.com/avGenie/url-shortener/internal/app/storage/api/model"
)

var (
	// Version Contains current version of program
	Version string
	// BuildTime Contains build time
	BuildTime string
	// BuildCommit Contains hash of build commit
	BuildCommit string
)

const naValue = "N/A"

func main() {
	printProgramInfo()

	config, err := config.InitConfig()
	if err != nil {
		zap.L().Fatal("Failed to initialize config", zap.Error(err))
	}

	err = logger.Initialize(config)
	if err != nil {
		zap.L().Fatal("Failed to initialize logger", zap.Error(err))
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

	if len(config.ProfilerFile) != 0 {
		fmem, err := os.Create(config.ProfilerFile)
		if err != nil {
			panic(err)
		}
		defer fmem.Close()
		runtime.GC()
		if err := pprof.WriteHeapProfile(fmem); err != nil {
			panic(err)
		}
	}

	zap.L().Info("Got interruption signal. Shutting down HTTP server gracefully...")
	err := server.Shutdown(context.Background())
	if err != nil {
		zap.L().Error("error while shutting down server", zap.Error(err))
	}
	router.Stop()
}

func printProgramInfo() {
	version := naValue
	if Version != "" {
		version = Version
	}
	fmt.Printf("Build version: %s\n", version)

	date := naValue
	if BuildTime != "" {
		date = BuildTime
	}
	fmt.Printf("Build date: %s\n", date)

	commit := naValue
	if BuildCommit != "" {
		commit = BuildCommit
	}
	fmt.Printf("Build commit: %s\n", commit)
}
