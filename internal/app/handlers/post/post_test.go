package post

import (
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/avGenie/url-shortener/internal/app/entity"
	"github.com/avGenie/url-shortener/internal/app/handlers/errors"
	"github.com/avGenie/url-shortener/internal/app/handlers/post/mock"
	"github.com/avGenie/url-shortener/internal/app/storage/api/model"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	netAddr       = "localhost:8080"
	baseURIPrefix = "http://localhost:8080"
)

func TestPostHandlerURL(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	s := mock.NewMockURLSaver(ctrl)

	type want struct {
		statusCode  int
		contentType string
		message     string
		resp        entity.Response
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
				message:     "42b3e75f",
				resp:        entity.OKResponse(),
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
				message:     errors.WrongURLFormat + "\n",
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
				message:     errors.InternalServerError + "\n",
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			request := httptest.NewRequest(http.MethodPost, test.request, strings.NewReader(test.URL))
			writer := httptest.NewRecorder()

			if test.want.resp.Status == "" {
				s.EXPECT().
					SaveURL(gomock.Any(), gomock.Any(), gomock.Any()).
					Times(0)
			} else {
				s.EXPECT().
					SaveURL(gomock.Any(), gomock.Any(), gomock.Any()).
					Times(1).
					Return(test.want.resp)
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

			if test.isError {
				assert.Equal(t, test.want.message, string(userResult))
				return
			}

			requiredOutput := fmt.Sprintf("http://%s/%s", netAddr, test.want.message)
			assert.Equal(t, requiredOutput, string(userResult))
		})
	}
}

func TestPostHandlerJSON(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	s := mock.NewMockURLSaver(ctrl)

	type want struct {
		statusCode   int
		contentType  string
		expectedBody string
		urlsValue    string
		resp         entity.Response
	}
	tests := []struct {
		name          string
		request       string
		body          string
		baseURIPrefix string
		urlsKey       string
		want          want
	}{
		{
			name:          "correct input data",
			request:       "/",
			body:          `{"url":"https://practicum.yandex.ru/"}`,
			baseURIPrefix: baseURIPrefix,
			urlsKey:       "42b3e75f",

			want: want{
				statusCode:   201,
				contentType:  "application/json",
				expectedBody: `{"result":"http://localhost:8080/42b3e75f"}` + "\n",
				urlsValue:    "https://practicum.yandex.ru/",
				resp:         entity.OKResponse(),
			},
		},
		{
			name:          "empty URL",
			request:       "/",
			body:          `{"url": ""}`,
			baseURIPrefix: baseURIPrefix,

			want: want{
				statusCode:   400,
				contentType:  "text/plain; charset=utf-8",
				expectedBody: errors.WrongJSONFormat + "\n",
			},
		},
		{
			name:    "cannot process JSON",
			request: "/",
			body:    "https://practicum.yandex.ru/",

			want: want{
				statusCode:   400,
				contentType:  "text/plain; charset=utf-8",
				expectedBody: errors.WrongJSONFormat + "\n",
			},
		},
		{
			name:    "empty base URI prefix",
			request: "/",
			body:    `{"url":"https://practicum.yandex.ru"}`,

			want: want{
				statusCode:   500,
				contentType:  "text/plain; charset=utf-8",
				expectedBody: errors.InternalServerError + "\n",
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			request := httptest.NewRequest(http.MethodPost, test.request, strings.NewReader(test.body))
			writer := httptest.NewRecorder()

			if test.want.resp.Status == "" {
				s.EXPECT().
					SaveURL(gomock.Any(), gomock.Any(), gomock.Any()).
					Times(0)
			} else {
				s.EXPECT().
					SaveURL(gomock.Any(), gomock.Any(), gomock.Any()).
					Times(1).
					Return(test.want.resp)
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
			ID: "practicum_id",
			InputURL: "https://practicum.yandex.ru/",
			ShortURL: "http://localhost:8080/42b3e75f",
		},
		model.BatchObject{
			ID: "yandex_id",
			InputURL: "https://yandex.ru/",
			ShortURL: "http://localhost:8080/77fca595",
		},
		model.BatchObject{
			ID: "google_id",
			InputURL: "https://www.google.com",
			ShortURL: "http://localhost:8080/ac6bb669",
		},
	}

	type want struct {
		statusCode   int
		contentType  string
		expectedBody string
		urlsValue    string
		resp         model.BatchResponse
	}
	tests := []struct {
		name          string
		request       string
		body          string
		baseURIPrefix string
		want          want
	}{
		{
			name:          "correct input data",
			request:       "/",
			body:          inputBatch,
			baseURIPrefix: baseURIPrefix,

			want: want{
				statusCode:   201,
				contentType:  "application/json",
				expectedBody: outputBatch,
				urlsValue:    "https://practicum.yandex.ru/",
				resp:         model.OKBatchResponse(batchResponse),
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			request := httptest.NewRequest(http.MethodPost, test.request, strings.NewReader(test.body))
			writer := httptest.NewRecorder()

			if test.want.resp.Status == "" {
				s.EXPECT().
					SaveBatchURL(gomock.Any(), gomock.Any()).
					Times(0)
			} else {
				s.EXPECT().
					SaveBatchURL(gomock.Any(), gomock.Any()).
					Times(1).
					Return(test.want.resp)
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

			assert.JSONEq(t, test.want.expectedBody, string(userResult))
		})
	}
}
