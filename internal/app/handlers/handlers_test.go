package handlers

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"github.com/avGenie/url-shortener/internal/app/config"
	"github.com/avGenie/url-shortener/internal/app/entity"
	"github.com/avGenie/url-shortener/internal/app/entity/mock"
	"github.com/go-chi/chi/v5"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	baseURIPrefix = "http://localhost:8080"
	netAddr       = "localhost:8080"
	dbStorage     = "/tmp/short-url-db.json"
)

func initTestConfig() *config.Config {
	return &config.Config{
		NetAddr:           netAddr,
		BaseURIPrefix:     baseURIPrefix,
		DBFileStoragePath: dbStorage,
	}
}

func initTestStorage(t *testing.T, config *config.Config) {
	err := InitStorage(*config)
	if err != nil {
		t.Fatal("cannot initialize config")
	}
}

func deferTestStorage(t *testing.T, config *config.Config) {
	CloseStorage(*config)
	_, err := os.Stat(config.DBFileStoragePath)
	assert.Error(t, err)
}

func TestPostHandlerURL(t *testing.T) {
	config := initTestConfig()
	initTestStorage(t, config)
	defer deferTestStorage(t, config)

	type want struct {
		statusCode  int
		contentType string
		message     string
	}
	tests := []struct {
		name          string
		request       string
		URL           string
		baseURIPrefix string
		isError       bool
		want          want
	}{
		{
			name:          "correct input data",
			request:       "/",
			URL:           "https://practicum.yandex.ru/",
			baseURIPrefix: baseURIPrefix,
			isError:       false,

			want: want{
				statusCode:  201,
				contentType: "text/plain; charset=utf-8",
				message:     "aHR0cHM6",
			},
		},
		{
			name:          "empty URL",
			request:       "/",
			baseURIPrefix: baseURIPrefix,
			isError:       true,

			want: want{
				statusCode:  400,
				contentType: "text/plain; charset=utf-8",
				message:     WrongURLFormat + "\n",
			},
		},
		{
			name:    "empty base URI prefix",
			request: "/",
			URL:     "https://practicum.yandex.ru/",
			isError: true,

			want: want{
				statusCode:  500,
				contentType: "text/plain; charset=utf-8",
				message:     InternalServerError + "\n",
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			request := httptest.NewRequest(http.MethodPost, test.request, strings.NewReader(test.URL))
			writer := httptest.NewRecorder()

			PostHandlerURL(test.baseURIPrefix, writer, request)

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

			requiredOutput := fmt.Sprintf("http://%s/%s", netAddr, test.want.message)
			assert.Equal(t, requiredOutput, string(userResult))

			url, ok := urls.Get(*entity.ParseURL(test.want.message))
			require.True(t, ok)
			assert.Equal(t, url.String(), test.URL)
		})
	}
}

func TestPostHandlerJSON(t *testing.T) {
	config := initTestConfig()
	initTestStorage(t, config)
	defer deferTestStorage(t, config)

	type want struct {
		statusCode   int
		contentType  string
		expectedBody string
		urlsValue    string
	}
	tests := []struct {
		name          string
		request       string
		body          string
		baseURIPrefix string
		urlsKey       string
		isError       bool
		want          want
	}{
		{
			name:          "correct input data",
			request:       "/",
			body:          `{"url":"https://practicum.yandex.ru/"}`,
			baseURIPrefix: baseURIPrefix,
			urlsKey:       "aHR0cHM6",
			isError:       false,

			want: want{
				statusCode:   201,
				contentType:  "application/json",
				expectedBody: `{"result":"http://localhost:8080/aHR0cHM6"}` + "\n",
				urlsValue:    "https://practicum.yandex.ru/",
			},
		},
		{
			name:          "empty URL",
			request:       "/",
			body:          `{"url": ""}`,
			baseURIPrefix: baseURIPrefix,
			isError:       true,

			want: want{
				statusCode:   400,
				contentType:  "text/plain; charset=utf-8",
				expectedBody: WrongURLFormat + "\n",
			},
		},
		{
			name:    "empty base URI prefix",
			request: "/",
			body:    `{"url":"https://practicum.yandex.ru"}`,
			isError: true,

			want: want{
				statusCode:   500,
				contentType:  "text/plain; charset=utf-8",
				expectedBody: InternalServerError + "\n",
			},
		},
		{
			name:    "cannot process JSON",
			request: "/",
			body:    "https://practicum.yandex.ru/",
			isError: true,

			want: want{
				statusCode:   400,
				contentType:  "text/plain; charset=utf-8",
				expectedBody: CannotProcessJSON + "\n",
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			request := httptest.NewRequest(http.MethodPost, test.request, strings.NewReader(test.body))
			writer := httptest.NewRecorder()

			PostHandlerJSON(test.baseURIPrefix, writer, request)

			res := writer.Result()

			assert.Equal(t, test.want.statusCode, res.StatusCode)
			assert.Equal(t, test.want.contentType, res.Header.Get("Content-Type"))

			userResult, err := io.ReadAll(res.Body)
			require.NoError(t, err)

			err = res.Body.Close()
			require.NoError(t, err)

			assert.Equal(t, test.want.expectedBody, string(userResult))

			if test.isError {
				return
			}

			url, ok := urls.Get(*entity.ParseURL(test.urlsKey))
			require.True(t, ok)
			assert.Equal(t, url.String(), test.want.urlsValue)
		})
	}
}

func TestGetHandler(t *testing.T) {
	config := initTestConfig()
	initTestStorage(t, config)
	defer deferTestStorage(t, config)

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

func TestGetPingDBHandler(t *testing.T) {
	ctrl := gomock.NewController(t)
	s := mock.NewMockStorage(ctrl)

	type want struct {
		statusCode int
		err        error
	}
	tests := []struct {
		name string
		want want
	}{
		{
			name: "successfull ping",

			want: want{
				statusCode: http.StatusOK,
				err:        nil,
			},
		},
		{
			name: "fallen ping",

			want: want{
				statusCode: http.StatusInternalServerError,
				err:        context.DeadlineExceeded,
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			request := httptest.NewRequest(http.MethodGet, "/{url}", nil)
			writer := httptest.NewRecorder()

			s.EXPECT().
				PingServer(gomock.Any()).
				Return(test.want.statusCode, test.want.err)

			GetPingDB(s, writer, request)

			res := writer.Result()

			err := res.Body.Close()
			require.NoError(t, err)

			assert.Equal(t, test.want.statusCode, res.StatusCode)
		})
	}
}
