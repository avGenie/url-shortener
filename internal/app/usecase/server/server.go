package server

import "net/http"

// Start Starts HTTP or HTTPS server
func Start(isHTTPS bool, server *http.Server) {
	if isHTTPS {
		startHTTPS(server)
		return
	}

	startHTTP(server)
}
