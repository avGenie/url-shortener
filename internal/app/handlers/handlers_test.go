package handlers

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/avGenie/url-shortener/internal/app/config"
	"github.com/avGenie/url-shortener/internal/app/entity"
	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestPostHandlerURL(t *testing.T) {
	config := config.InitConfig()

	type want struct {
		statusCode  int
		contentType string
		message     string
	}
	tests := []struct {
		name    string
		request string
		URL     string
		isError bool
		want    want
	}{
		{
			name:    "correct input data",
			request: "/",
			URL:     "https://practicum.yandex.ru/",
			isError: false,

			want: want{
				statusCode:  201,
				contentType: "text/plain; charset=utf-8",
				message:     "aHR0cHM6",
			},
		},
		{
			name:    "empty URL",
			request: "/",
			URL:     "",
			isError: true,

			want: want{
				statusCode:  400,
				contentType: "text/plain; charset=utf-8",
				message:     WrongURLFormat + "\n",
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			request := httptest.NewRequest(http.MethodPost, test.request, strings.NewReader(test.URL))
			writer := httptest.NewRecorder()

			ctx := context.WithValue(request.Context(), baseURIPrefixCtx, config.BaseURIPrefix)
			request = request.WithContext(ctx)

			PostHandlerURL(writer, request)

			res := writer.Result()

			assert.Equal(t, test.want.statusCode, res.StatusCode)
			assert.Equal(t, test.want.contentType, res.Header.Get("Content-Type"))

			userResult, err := io.ReadAll(res.Body)
			require.NoError(t, err)

			err = res.Body.Close()
			require.NoError(t, err)

			if test.isError {
				assert.Equal(t, test.want.message, string(userResult))
				return
			}

			requiredOutput := fmt.Sprintf("http://%s/%s", config.NetAddr, test.want.message)
			assert.Equal(t, requiredOutput, string(userResult))

			url, ok := urls.Get(*entity.ParseURL(test.want.message))
			require.True(t, ok)
			assert.Equal(t, url.String(), test.URL)
		})
	}
}

func TestGetHandler(t *testing.T) {
	type want struct {
		statusCode  int
		contentType string
		location    string
		message     string
	}
	tests := []struct {
		name    string
		request string
		urls    map[entity.URL]entity.URL
		want    want
	}{
		{
			name:    "correct input data",
			request: "aHR0cHM6",

			urls: map[entity.URL]entity.URL{
				*entity.ParseURL("aHR0cHM6"): *entity.ParseURL("https://practicum.yandex.ru/"),
			},

			want: want{
				statusCode:  http.StatusTemporaryRedirect,
				contentType: "text/plain; charset=utf-8",
				location:    "https://practicum.yandex.ru/",
				message:     "",
			},
		},
		{
			name:    "request without id",
			request: "",

			urls: map[entity.URL]entity.URL{
				*entity.ParseURL("aHR0cHM6"): *entity.ParseURL("https://practicum.yandex.ru/"),
			},

			want: want{
				statusCode:  http.StatusBadRequest,
				contentType: "text/plain; charset=utf-8",
				location:    "",
				message:     ShortURLNotInDB + "\n",
			},
		},
		{
			name:    "missing URL",
			request: "/fsdfuytu",

			urls: map[entity.URL]entity.URL{
				*entity.ParseURL("aHR0cHM6"): *entity.ParseURL("https://practicum.yandex.ru/"),
			},

			want: want{
				statusCode:  http.StatusBadRequest,
				contentType: "text/plain; charset=utf-8",
				location:    "",
				message:     ShortURLNotInDB + "\n",
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			for key, value := range test.urls {
				urls.Add(key, value)
			}

			request := httptest.NewRequest(http.MethodGet, "/{url}", nil)
			writer := httptest.NewRecorder()

			rctx := chi.NewRouteContext()
			rctx.URLParams.Add("url", test.request)

			request = request.WithContext(context.WithValue(request.Context(), chi.RouteCtxKey, rctx))

			GetHandler(writer, request)

			res := writer.Result()

			assert.Equal(t, test.want.statusCode, res.StatusCode)
			assert.Equal(t, test.want.contentType, res.Header.Get("Content-Type"))
			assert.Equal(t, test.want.location, res.Header.Get("Location"))

			userResult, err := io.ReadAll(res.Body)
			require.NoError(t, err)

			err = res.Body.Close()
			require.NoError(t, err)

			assert.Equal(t, test.want.message, string(userResult))
		})
	}
}
