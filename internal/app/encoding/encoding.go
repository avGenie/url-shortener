package encoding

import (
	"net/http"
	"strings"

	"github.com/avGenie/url-shortener/internal/app/logger"
	"go.uber.org/zap"
)

const (
	gzipEncodingFormat = "gzip"
)

func GzipMiddleware(h http.Handler) http.Handler {
	gzipFn := func(writer http.ResponseWriter, req *http.Request) {
		ow := writer

		acceptEncoding := req.Header.Get("Accept-Encoding")
		if strings.Contains(acceptEncoding, gzipEncodingFormat) {
			logger.Log.Debug("sending gzip encoded message")
			cw := newCompressWriter(writer)
			cw.writer.Header().Set("Content-Encoding", "gzip")
			ow = cw
			defer cw.Close()
		}

		contentEncoding := req.Header.Get("Content-Encoding")
		if strings.Contains(contentEncoding, gzipEncodingFormat) {
			cr, err := newCompressReader(req.Body)
			if err != nil {
				logger.Log.Error("invalid compress reader creation", zap.Error(err))
				ow.WriteHeader(http.StatusInternalServerError)
				return
			}

			logger.Log.Debug("obtained gzip encoded message")
			req.Body = cr
			defer cr.Close()
		}

		h.ServeHTTP(ow, req)
	}

	return http.HandlerFunc(gzipFn)
}
