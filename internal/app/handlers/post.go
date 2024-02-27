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
	var shortURL *entity.URL
	exists := true

	encodedURL := base64.StdEncoding.EncodeToString([]byte(inputURL))
	availableURLCount := len(encodedURL) / maxEncodedSize
	for i := 0; i < availableURLCount-1; i++ {
		shortURL = entity.ParseURL(encodedURL[(maxEncodedSize * i):(maxEncodedSize * (i + 1))])
		fmt.Println(shortURL)
		if _, exists = urls.Get(*shortURL); !exists {
			break
		}
	}

	if exists {
		return ""
	}

	userURL := entity.ParseURL(inputURL)
	urls.Add(*shortURL, *userURL)

	return fmt.Sprintf("%s/%s", baseURIPrefix, shortURL)
}
