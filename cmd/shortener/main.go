package main

import (
		"net/http"

	"github.com/avGenie/url-shortener/internal/app/config"
	"github.com/avGenie/url-shortener/internal/app/handlers"
	"github.com/avGenie/url-shortener/internal/app/logger"
	"github.com/avGenie/url-shortener/internal/app/storage"
	"go.uber.org/zap"
)

func main() {
	cnf := config.InitConfig()
	
	err := logger.Initialize(cnf)
	if err != nil {
		panic(err.Error())
	}

	sugar := *zap.S()
	defer sugar.Sync()

	db, err := storage.InitStorage(cnf)
	if err != nil {
		sugar.Fatalw(
			err.Error(),
			"event", "init storage",
		)
	}
	defer db.Close()

	sugar.Infow(
		"Starting server",
		"addr", cnf.NetAddr,
	)

	err = http.ListenAndServe(cnf.NetAddr, handlers.CreateRouter(cnf, db))
	if err != nil && err != http.ErrServerClosed {
		sugar.Fatalw(
			err.Error(),
			"event", "start server",
		)
	}
}
