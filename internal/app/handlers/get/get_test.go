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

	"github.com/avGenie/url-shortener/internal/app/handlers/errors"
	"github.com/avGenie/url-shortener/internal/app/handlers/get/mock"
	"github.com/avGenie/url-shortener/internal/app/logger"
	"github.com/avGenie/url-shortener/internal/app/models"
	"github.com/go-chi/chi/v5"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	baseURIPrefix = "http://localhost:8080"
)

func TestGetHandler(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	s := mock.NewMockURLGetter(ctrl)

	type want struct {
		statusCode  int
		contentType string
		location    string
		expectErr   error
		expectURL   *entity.URL
		message     string
	}
	tests := []struct {
		name              string
		request           string
		userIDCtx         entity.UserIDCtx
		exitBeforeGetting bool
		want              want
	}{
		{
			name:    "correct input data",
			request: "aHR0cHM6",
			userIDCtx: entity.UserIDCtx{
				UserID:     "ac2a4811-4f10-487f-bde3-e39a14af7cd8",
				StatusCode: http.StatusOK,
			},

			want: want{
				statusCode:  http.StatusTemporaryRedirect,
				contentType: "text/plain; charset=utf-8",
				location:    "https://practicum.yandex.ru/",
				expectURL:   makeOKURLResponse("https://practicum.yandex.ru/"),
				expectErr:   nil,
				message:     "",
			},
		},
		{
			name:    "request without id",
			request: "",
			userIDCtx: entity.UserIDCtx{
				UserID:     "ac2a4811-4f10-487f-bde3-e39a14af7cd8",
				StatusCode: http.StatusOK,
			},

			want: want{
				statusCode:  http.StatusBadRequest,
				contentType: "text/plain; charset=utf-8",
				location:    "",
				expectURL:   nil,
				expectErr:   fmt.Errorf(""),
				message:     errors.ShortURLNotInDB + "\n",
			},
		},
		{
			name:    "missing URL",
			request: "/fsdfuytu",
			userIDCtx: entity.UserIDCtx{
				UserID:     "ac2a4811-4f10-487f-bde3-e39a14af7cd8",
				StatusCode: http.StatusOK,
			},

			want: want{
				statusCode:  http.StatusBadRequest,
				contentType: "text/plain; charset=utf-8",
				location:    "",
				expectURL:   nil,
				expectErr:   fmt.Errorf(""),
				message:     errors.ShortURLNotInDB + "\n",
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			request := httptest.NewRequest(http.MethodGet, "/{url}", nil)
			writer := httptest.NewRecorder()

			rctx := chi.NewRouteContext()
			rctx.URLParams.Add("url", test.request)

			if test.exitBeforeGetting {
				s.EXPECT().GetURL(gomock.Any(), gomock.Any(), gomock.Any()).Times(0)
			} else {
				s.EXPECT().GetURL(gomock.Any(), gomock.Any(), gomock.Any()).
					Return(test.want.expectURL, test.want.expectErr)
			}

			request = request.WithContext(context.WithValue(request.Context(), chi.RouteCtxKey, rctx))
			request = request.WithContext(context.WithValue(request.Context(), entity.UserIDCtxKey{}, test.userIDCtx))

			handler := URLHandler(s)
			handler(writer, request)

			res := writer.Result()

			assert.Equal(t, test.want.statusCode, res.StatusCode)
			assert.Equal(t, test.want.contentType, res.Header.Get("Content-Type"))
			assert.Equal(t, test.want.location, res.Header.Get("Location"))

			if test.exitBeforeGetting {
				return
			}

			userResult, err := io.ReadAll(res.Body)
			require.NoError(t, err)

			err = res.Body.Close()
			require.NoError(t, err)

			assert.Equal(t, test.want.message, string(userResult))
		})
	}
}

func TestGetUserURLHandler(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	s := mock.NewMockAllURLGetter(ctrl)

	outputStorageBatch := models.AllUrlsBatch{
		{
			ShortURL:    "42b3e75f",
			OriginalURL: "https://practicum.yandex.ru/",
		},
		{
			ShortURL:    "77fca595",
			OriginalURL: "https://yandex.ru/",
		},
		{
			ShortURL:    "ac6bb669",
			OriginalURL: "https://www.google.com",
		},
	}

	outputBatch := strings.TrimSpace(`
	[
		{
			"short_url": "http://localhost:8080/42b3e75f",
			"original_url": "https://practicum.yandex.ru/"
		},
		{
			"short_url": "http://localhost:8080/77fca595",
			"original_url": "https://yandex.ru/"
		},
		{
			"short_url": "http://localhost:8080/ac6bb669",
			"original_url": "https://www.google.com"
		}
	]`)

	type want struct {
		statusCode  int
		contentType string
		expectErr   error
		message     string
	}
	tests := []struct {
		name               string
		baseURIPrefix      string
		outputStorageBatch models.AllUrlsBatch
		userIDCtx          entity.UserIDCtx
		invalidOutput      bool
		exitBeforeGetting  bool
		want               want
	}{
		{
			name:               "correct input data",
			baseURIPrefix:      baseURIPrefix,
			outputStorageBatch: outputStorageBatch,
			userIDCtx: entity.UserIDCtx{
				UserID:     "ac2a4811-4f10-487f-bde3-e39a14af7cd8",
				StatusCode: http.StatusOK,
			},

			want: want{
				statusCode:  http.StatusCreated,
				contentType: "application/json",
				expectErr:   nil,
				message:     outputBatch,
			},
		},
		{
			name:               "empty output",
			baseURIPrefix:      baseURIPrefix,
			outputStorageBatch: models.AllUrlsBatch{},
			userIDCtx: entity.UserIDCtx{
				UserID:     "ac2a4811-4f10-487f-bde3-e39a14af7cd8",
				StatusCode: http.StatusOK,
			},
			invalidOutput: true,

			want: want{
				statusCode: http.StatusNoContent,
				expectErr:  nil,
			},
		},
		{
			name:               "error while getting from storage",
			baseURIPrefix:      baseURIPrefix,
			outputStorageBatch: models.AllUrlsBatch{},
			userIDCtx: entity.UserIDCtx{
				UserID:     "ac2a4811-4f10-487f-bde3-e39a14af7cd8",
				StatusCode: http.StatusOK,
			},
			invalidOutput: true,

			want: want{
				statusCode: http.StatusInternalServerError,
				expectErr:  fmt.Errorf("error"),
			},
		},
		{
			name: "empty base URI prefix",
			userIDCtx: entity.UserIDCtx{
				UserID:     "ac2a4811-4f10-487f-bde3-e39a14af7cd8",
				StatusCode: http.StatusOK,
			},
			invalidOutput:     true,
			exitBeforeGetting: true,

			want: want{
				statusCode: http.StatusInternalServerError,
			},
		},
		{
			name:    "unathorized user",
			baseURIPrefix: baseURIPrefix,
			userIDCtx: entity.UserIDCtx{
				UserID:     "",
				StatusCode: http.StatusUnauthorized,
			},
			invalidOutput:     true,
			exitBeforeGetting: true,

			want: want{
				statusCode:  http.StatusUnauthorized,
				contentType: "",
			},
		},
		{
			name:    "missing user id",
			baseURIPrefix: baseURIPrefix,
			userIDCtx: entity.UserIDCtx{
				UserID:     "",
				StatusCode: http.StatusOK,
			},
			invalidOutput:     true,
			exitBeforeGetting: true,

			want: want{
				statusCode:  http.StatusInternalServerError,
				contentType: "",
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			request := httptest.NewRequest(http.MethodGet, "/api/user/urls", nil)
			writer := httptest.NewRecorder()

			request = request.WithContext(context.WithValue(request.Context(), entity.UserIDCtxKey{}, test.userIDCtx))

			if test.exitBeforeGetting {
				s.EXPECT().GetAllURLByUserID(gomock.Any(), gomock.Any()).Times(0)
			} else {
				s.EXPECT().GetAllURLByUserID(gomock.Any(), gomock.Any()).Return(test.outputStorageBatch, test.want.expectErr)
			}

			handler := UserURLsHandler(s, test.baseURIPrefix)
			handler(writer, request)

			res := writer.Result()

			assert.Equal(t, test.want.statusCode, res.StatusCode)
			assert.Equal(t, test.want.contentType, res.Header.Get("Content-Type"))

			userResult, err := io.ReadAll(res.Body)
			require.NoError(t, err)

			err = res.Body.Close()
			require.NoError(t, err)

			if !test.invalidOutput {
				assert.JSONEq(t, test.want.message, string(userResult))
			}
		})
	}
}

func TestGetPingDBHandler(t *testing.T) {
	cnf := config.InitConfig()

	logger.Initialize(cnf)

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	s := mock.NewMockStoragePinger(ctrl)

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
				Return(test.want.err)

			handler := PingDBHandler(s)
			handler(writer, request)

			res := writer.Result()

			err := res.Body.Close()
			require.NoError(t, err)

			assert.Equal(t, test.want.statusCode, res.StatusCode)
		})
	}
}

func makeOKURLResponse(URL string) *entity.URL {
	outURL, _ := entity.NewURL(URL)

	return outURL
}
