package server

import (
	"net/http"

	"go.uber.org/zap"
)

// startHTTP Starts HTTP server
func startHTTP(server *http.Server) {
	zap.L().Info("Start HTTP server", zap.String("addr", server.Addr))

	err := server.ListenAndServe()
	if err != nil && err != http.ErrServerClosed {
		zap.L().Fatal("fatal error while starting server", zap.Error(err))
	}
}
