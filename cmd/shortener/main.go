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
	"github.com/avGenie/url-shortener/internal/app/grpc"
	handlers "github.com/avGenie/url-shortener/internal/app/handlers/router"
	"github.com/avGenie/url-shortener/internal/app/logger"
	storage "github.com/avGenie/url-shortener/internal/app/storage/api"
	"github.com/avGenie/url-shortener/internal/app/storage/api/model"
	cidr "github.com/avGenie/url-shortener/internal/app/usecase/CIDR"
	usecase_server "github.com/avGenie/url-shortener/internal/app/usecase/server"
)

// Variables which contains build flag values
var (
	// Version Contains current version of program
	Version = "N/A"
	// BuildTime Contains build time
	BuildTime = "N/A"
	// BuildCommit Contains hash of build commit
	BuildCommit = "N/A"
)

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

	var cidrObj *cidr.CIDR
	if config.TrustedSubnet != "" {
		cidrObj, err = cidr.NewCIDR(config.TrustedSubnet)
		if err != nil {
			sugar.Fatalw(
				err.Error(),
				"event", "cidr creation",
			)
		}
	}

	startHTTPServer(config, storage, cidrObj)
}

func startHTTPServer(config config.Config, storage model.Storage, cidr *cidr.CIDR) {
	ctx, cancel := signal.NotifyContext(
		context.Background(),
		syscall.SIGTERM,
		syscall.SIGINT,
		syscall.SIGQUIT,
		os.Interrupt,
	)
	defer cancel()

	router := handlers.NewRouter(config, storage, cidr)

	server := &http.Server{
		Addr:    config.NetAddr,
		Handler: router.Mux,
	}

	go usecase_server.Start(config.EnableHTTPS, server)

	grpcServer := grpc.NewGRPCServer(config, storage)

	go grpcServer.Start()

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

	grpcServer.Stop()

	zap.L().Info("Got interruption signal. Shutting down HTTP server gracefully...")
	err := server.Shutdown(context.Background())
	if err != nil {
		zap.L().Error("error while shutting down server", zap.Error(err))
	}
	router.Stop()
}

func printProgramInfo() {
	fmt.Printf("Build version: %s\n", Version)
	fmt.Printf("Build date: %s\n", BuildTime)
	fmt.Printf("Build commit: %s\n", BuildCommit)
}
