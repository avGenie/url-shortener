package server

import (
	"crypto/tls"
	"net/http"

	"github.com/avGenie/url-shortener/internal/app/usecase/https"
	"go.uber.org/zap"
)

// startHTTPS Starts HTTPS server
func startHTTPS(server *http.Server) {
	zap.L().Info("Start HTTPS server", zap.String("addr", server.Addr))

	cert, err := https.GenerateTLSCert()
	if err != nil && err != http.ErrServerClosed {
		zap.L().Fatal("fatal error while starting https server", zap.Error(err))
	}

	server.TLSConfig = &tls.Config{
		Certificates: []tls.Certificate{cert},
	}

	err = server.ListenAndServeTLS("", "")
	if err != nil && err != http.ErrServerClosed {
		zap.L().Fatal("fatal error while starting https server", zap.Error(err))
	}
}
