package handlers

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/avGenie/url-shortener/internal/app/entity"
	"github.com/avGenie/url-shortener/internal/app/handlers/errors"
	"github.com/avGenie/url-shortener/internal/app/handlers/post/mock"
	storage_err "github.com/avGenie/url-shortener/internal/app/storage/api/errors"

	"github.com/avGenie/url-shortener/internal/app/storage/api/model"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	baseURIPrefix = "http://localhost:8080"
)

func TestPostHandlerURL(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	s := mock.NewMockURLSaver(ctrl)

	type want struct {
		expectedErr  error
		contentType  string
		expectedBody string
		statusCode   int
		isSaveURL    bool
	}
	tests := []struct {
		name          string
		request       string
		body          string
		baseURIPrefix string
		userIDCtx     entity.UserIDCtx
		want          want
	}{
		{
			name:          "correct input data",
			request:       "/",
			body:          "https://practicum.yandex.ru/",
			baseURIPrefix: baseURIPrefix,
			userIDCtx: entity.UserIDCtx{
				UserID:     "ac2a4811-4f10-487f-bde3-e39a14af7cd8",
				StatusCode: http.StatusOK,
			},

			want: want{
				statusCode:   http.StatusCreated,
				contentType:  "text/plain; charset=utf-8",
				expectedBody: "http://localhost:8080/42b3e75f",
				expectedErr:  nil,
				isSaveURL:    true,
			},
		},
		{
			name:          "url already exists",
			request:       "/",
			body:          "https://practicum.yandex.ru/",
			baseURIPrefix: baseURIPrefix,
			userIDCtx: entity.UserIDCtx{
				UserID:     "ac2a4811-4f10-487f-bde3-e39a14af7cd8",
				StatusCode: http.StatusOK,
			},

			want: want{
				statusCode:   http.StatusConflict,
				contentType:  "text/plain; charset=utf-8",
				expectedBody: "http://localhost:8080/42b3e75f",
				expectedErr:  fmt.Errorf("error: %w", storage_err.ErrURLAlreadyExists),
				isSaveURL:    true,
			},
		},
		{
			name:          "empty URL",
			request:       "/",
			baseURIPrefix: baseURIPrefix,
			userIDCtx: entity.UserIDCtx{
				UserID:     "ac2a4811-4f10-487f-bde3-e39a14af7cd8",
				StatusCode: http.StatusOK,
			},

			want: want{
				statusCode:   http.StatusBadRequest,
				contentType:  "text/plain; charset=utf-8",
				expectedBody: errors.WrongURLFormat + "\n",
			},
		},
		{
			name:    "empty base URI prefix",
			request: "/",
			body:    "https://practicum.yandex.ru/",
			userIDCtx: entity.UserIDCtx{
				UserID:     "ac2a4811-4f10-487f-bde3-e39a14af7cd8",
				StatusCode: http.StatusOK,
			},

			want: want{
				statusCode: http.StatusInternalServerError,
			},
		},
		{
			name:          "unathorized user",
			request:       "/",
			body:          "https://practicum.yandex.ru/",
			baseURIPrefix: baseURIPrefix,
			userIDCtx: entity.UserIDCtx{
				UserID:     "",
				StatusCode: http.StatusUnauthorized,
			},

			want: want{
				statusCode: http.StatusInternalServerError,
			},
		},
		{
			name:          "missing user id",
			request:       "/",
			body:          "https://practicum.yandex.ru/",
			baseURIPrefix: baseURIPrefix,
			userIDCtx: entity.UserIDCtx{
				UserID:     "",
				StatusCode: http.StatusInternalServerError,
			},

			want: want{
				statusCode: http.StatusInternalServerError,
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			request := httptest.NewRequest(http.MethodPost, test.request, strings.NewReader(test.body))
			writer := httptest.NewRecorder()

			request = request.WithContext(context.WithValue(request.Context(), entity.UserIDCtxKey{}, test.userIDCtx))

			if test.want.isSaveURL == false {
				s.EXPECT().
					SaveURL(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
					Times(0)
			} else {
				s.EXPECT().
					SaveURL(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
					Times(1).
					Return(test.want.expectedErr)
			}

			handler := URLHandler(s, test.baseURIPrefix)
			handler(writer, request)

			res := writer.Result()

			assert.Equal(t, test.want.statusCode, res.StatusCode)
			assert.Equal(t, test.want.contentType, res.Header.Get("Content-Type"))

			userResult, err := io.ReadAll(res.Body)
			require.NoError(t, err)

			err = res.Body.Close()
			require.NoError(t, err)

			assert.Equal(t, test.want.expectedBody, string(userResult))
		})
	}
}

func TestPostHandlerJSON(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	s := mock.NewMockURLSaver(ctrl)

	type want struct {
		expectedErr  error
		contentType  string
		expectedBody string
		urlsValue    string
		statusCode   int
		isSaveURL    bool
	}
	tests := []struct {
		name          string
		request       string
		body          string
		baseURIPrefix string
		urlsKey       string
		userIDCtx     entity.UserIDCtx
		want          want
	}{
		{
			name:          "correct input data",
			request:       "/",
			body:          `{"url":"https://practicum.yandex.ru/"}`,
			baseURIPrefix: baseURIPrefix,
			urlsKey:       "42b3e75f",
			userIDCtx: entity.UserIDCtx{
				UserID:     "ac2a4811-4f10-487f-bde3-e39a14af7cd8",
				StatusCode: http.StatusOK,
			},

			want: want{
				statusCode:   http.StatusCreated,
				contentType:  "application/json",
				expectedBody: `{"result":"http://localhost:8080/42b3e75f"}` + "\n",
				urlsValue:    "https://practicum.yandex.ru/",
				expectedErr:  nil,
				isSaveURL:    true,
			},
		},
		{
			name:          "url already exists",
			request:       "/",
			body:          `{"url":"https://practicum.yandex.ru/"}`,
			baseURIPrefix: baseURIPrefix,
			urlsKey:       "42b3e75f",
			userIDCtx: entity.UserIDCtx{
				UserID:     "ac2a4811-4f10-487f-bde3-e39a14af7cd8",
				StatusCode: http.StatusOK,
			},

			want: want{
				statusCode:   http.StatusConflict,
				contentType:  "application/json",
				expectedBody: `{"result":"http://localhost:8080/42b3e75f"}` + "\n",
				urlsValue:    "https://practicum.yandex.ru/",
				expectedErr:  fmt.Errorf("error: %w", storage_err.ErrURLAlreadyExists),
				isSaveURL:    true,
			},
		},
		{
			name:          "empty URL",
			request:       "/",
			body:          `{"url": ""}`,
			baseURIPrefix: baseURIPrefix,
			userIDCtx: entity.UserIDCtx{
				UserID:     "ac2a4811-4f10-487f-bde3-e39a14af7cd8",
				StatusCode: http.StatusOK,
			},

			want: want{
				statusCode:   http.StatusBadRequest,
				contentType:  "text/plain; charset=utf-8",
				expectedBody: errors.WrongJSONFormat + "\n",
			},
		},
		{
			name:          "cannot process JSON",
			request:       "/",
			body:          "https://practicum.yandex.ru/",
			baseURIPrefix: baseURIPrefix,
			userIDCtx: entity.UserIDCtx{
				UserID:     "ac2a4811-4f10-487f-bde3-e39a14af7cd8",
				StatusCode: http.StatusOK,
			},

			want: want{
				statusCode:   http.StatusBadRequest,
				contentType:  "text/plain; charset=utf-8",
				expectedBody: errors.WrongJSONFormat + "\n",
			},
		},
		{
			name:    "empty base URI prefix",
			request: "/",
			body:    `{"url":"https://practicum.yandex.ru"}`,
			userIDCtx: entity.UserIDCtx{
				UserID:     "ac2a4811-4f10-487f-bde3-e39a14af7cd8",
				StatusCode: http.StatusOK,
			},

			want: want{
				statusCode: http.StatusInternalServerError,
			},
		},
		{
			name:          "missing user id",
			request:       "/",
			body:          `{"url":"https://practicum.yandex.ru"}`,
			baseURIPrefix: baseURIPrefix,
			userIDCtx: entity.UserIDCtx{
				UserID:     "",
				StatusCode: http.StatusOK,
			},

			want: want{
				statusCode: http.StatusInternalServerError,
			},
		},
		{
			name:          "error url processing",
			request:       "/",
			body:          `{"url":"https://practicum.yandex.ru/"}`,
			baseURIPrefix: baseURIPrefix,
			urlsKey:       "42b3e75f",
			userIDCtx: entity.UserIDCtx{
				UserID:     "ac2a4811-4f10-487f-bde3-e39a14af7cd8",
				StatusCode: http.StatusOK,
			},

			want: want{
				statusCode:  http.StatusInternalServerError,
				urlsValue:   "https://practicum.yandex.ru/",
				expectedErr: fmt.Errorf("error"),
				isSaveURL:   true,
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			request := httptest.NewRequest(http.MethodPost, test.request, strings.NewReader(test.body))
			writer := httptest.NewRecorder()

			request = request.WithContext(context.WithValue(request.Context(), entity.UserIDCtxKey{}, test.userIDCtx))

			if test.want.isSaveURL == false {
				s.EXPECT().
					SaveURL(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
					Times(0)
			} else {
				s.EXPECT().
					SaveURL(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
					Times(1).
					Return(test.want.expectedErr)
			}

			handler := JSONHandler(s, test.baseURIPrefix)
			handler(writer, request)

			res := writer.Result()

			assert.Equal(t, test.want.statusCode, res.StatusCode)
			assert.Equal(t, test.want.contentType, res.Header.Get("Content-Type"))

			userResult, err := io.ReadAll(res.Body)
			require.NoError(t, err)

			err = res.Body.Close()
			require.NoError(t, err)

			assert.Equal(t, test.want.expectedBody, string(userResult))
		})
	}
}

func TestPostHandlerJSONBatch(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	s := mock.NewMockURLBatchSaver(ctrl)

	inputBatch := `
	[
		{
			"correlation_id": "practicum_id",
			"original_url": "https://practicum.yandex.ru/"
		},
		{
			"correlation_id": "yandex_id",
			"original_url": "https://yandex.ru/"
		},
		{
			"correlation_id": "google_id",
			"original_url": "https://www.google.com"
		}
	]`

	outputBatch := strings.TrimSpace(`
	[
		{
			"correlation_id": "practicum_id",
			"short_url": "http://localhost:8080/42b3e75f"
		},
		{
			"correlation_id": "yandex_id",
			"short_url": "http://localhost:8080/77fca595"
		},
		{
			"correlation_id": "google_id",
			"short_url": "http://localhost:8080/ac6bb669"
		}
	]`)

	batchResponse := model.Batch{
		model.BatchObject{
			ID:       "practicum_id",
			InputURL: "https://practicum.yandex.ru/",
			ShortURL: "42b3e75f",
		},
		model.BatchObject{
			ID:       "yandex_id",
			InputURL: "https://yandex.ru/",
			ShortURL: "77fca595",
		},
		model.BatchObject{
			ID:       "google_id",
			InputURL: "https://www.google.com",
			ShortURL: "ac6bb669",
		},
	}

	type want struct {
		expectedErr   error
		contentType   string
		expectedBody  string
		urlsValue     string
		expectedBatch model.Batch
		statusCode    int
	}
	tests := []struct {
		name          string
		request       string
		body          string
		baseURIPrefix string
		userIDCtx     entity.UserIDCtx
		want          want
		isSaveURL     bool
	}{
		{
			name:          "correct input data",
			request:       "/",
			body:          inputBatch,
			baseURIPrefix: baseURIPrefix,
			userIDCtx: entity.UserIDCtx{
				UserID:     "ac2a4811-4f10-487f-bde3-e39a14af7cd8",
				StatusCode: http.StatusOK,
			},
			isSaveURL: true,

			want: want{
				statusCode:    201,
				contentType:   "application/json",
				expectedBody:  outputBatch,
				urlsValue:     "https://practicum.yandex.ru/",
				expectedBatch: batchResponse,
				expectedErr:   nil,
			},
		},
		{
			name:    "empty base URI prefix",
			request: "/",
			body:    `{"url":"https://practicum.yandex.ru"}`,
			userIDCtx: entity.UserIDCtx{
				UserID:     "ac2a4811-4f10-487f-bde3-e39a14af7cd8",
				StatusCode: http.StatusOK,
			},

			want: want{
				statusCode: http.StatusInternalServerError,
			},
		},
		{
			name:          "missing user id",
			request:       "/",
			body:          `{"url":"https://practicum.yandex.ru"}`,
			baseURIPrefix: baseURIPrefix,
			userIDCtx: entity.UserIDCtx{
				UserID:     "",
				StatusCode: http.StatusOK,
			},

			want: want{
				statusCode: http.StatusInternalServerError,
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			request := httptest.NewRequest(http.MethodPost, test.request, strings.NewReader(test.body))
			writer := httptest.NewRecorder()

			request = request.WithContext(context.WithValue(request.Context(), entity.UserIDCtxKey{}, test.userIDCtx))

			if test.isSaveURL == false {
				s.EXPECT().
					SaveBatchURL(gomock.Any(), gomock.Any(), gomock.Any()).
					Times(0)
			} else {
				s.EXPECT().
					SaveBatchURL(gomock.Any(), gomock.Any(), gomock.Any()).
					Times(1).
					Return(test.want.expectedBatch, test.want.expectedErr)
			}

			handler := JSONBatchHandler(s, test.baseURIPrefix)
			handler(writer, request)

			res := writer.Result()

			assert.Equal(t, test.want.statusCode, res.StatusCode)
			assert.Equal(t, test.want.contentType, res.Header.Get("Content-Type"))

			userResult, err := io.ReadAll(res.Body)
			require.NoError(t, err)

			err = res.Body.Close()
			require.NoError(t, err)

			if test.isSaveURL {
				assert.JSONEq(t, test.want.expectedBody, string(userResult))
			}
		})
	}
}
