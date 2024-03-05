package handlers

import (
	"encoding/base64"
	"fmt"
	"net/http"

	"github.com/avGenie/url-shortener/internal/app/config"
	"github.com/avGenie/url-shortener/internal/app/entity"
)

type PostContext struct {
	baseURIPrefix string
	handle        func(string, http.ResponseWriter, *http.Request)
}

func (ctx *PostContext) Handle() http.HandlerFunc {
	return func(writer http.ResponseWriter, req *http.Request) {
		ctx.handle(ctx.baseURIPrefix, writer, req)
	}
}

func NewPostContextURL(config config.Config) *PostContext {
	return &PostContext{
		baseURIPrefix: config.BaseURIPrefix,
		handle: PostHandlerURL,
	}
}

func NewPostContextJSON(config config.Config) *PostContext {
	return &PostContext{
		baseURIPrefix: config.BaseURIPrefix,
		handle: PostHandlerJSON,
	}
}

func postURLProcessing(inputURL, baseURIPrefix string) (string, error) {
	var shortURL *entity.URL

	userURL := entity.ParseURL(inputURL)
	added := true

	encodedURL := base64.StdEncoding.EncodeToString([]byte(inputURL))
	availableURLCount := len(encodedURL) / maxEncodedSize
	for i := 0; i < availableURLCount-1; i++ {
		shortURL = entity.ParseURL(encodedURL[(maxEncodedSize * i):(maxEncodedSize * (i + 1))])
		isAdded, err := urls.Add(*shortURL, *userURL)
		if err != nil {
			return "", err
		}
		added = isAdded
		if isAdded {
			break
		}
	}

	if !added {
		return "", nil
	}

	return fmt.Sprintf("%s/%s", baseURIPrefix, shortURL), nil
}
