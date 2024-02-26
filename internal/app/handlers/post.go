package handlers

import (
	"context"
	"encoding/base64"
	"fmt"
	"net/http"

	"github.com/avGenie/url-shortener/internal/app/config"
	"github.com/avGenie/url-shortener/internal/app/entity"
)

func PostMiddleware(config config.Config, h http.HandlerFunc) http.HandlerFunc {
	logFn := func(w http.ResponseWriter, r *http.Request) {
		r = r.WithContext(context.WithValue(r.Context(), baseURIPrefixCtx, config.BaseURIPrefix))

		h.ServeHTTP(w, r)
	}

	return http.HandlerFunc(logFn)
}

func postURLProcessing(inputURL, baseURIPrefix string) string {
	userURL := entity.ParseURL(inputURL)

	encodedURL := base64.StdEncoding.EncodeToString([]byte(inputURL))
	shortURL := entity.ParseURL(encodedURL[:maxEncodedSize])

	urls.Add(*shortURL, *userURL)

	return fmt.Sprintf("%s/%s", baseURIPrefix, shortURL)
}
