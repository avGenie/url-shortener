package handlers

import (
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestPostHandler(t *testing.T) {
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
				contentType: "text/plain",
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
				contentType: "text/plain",
				message:     EmptyURL,
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			request := httptest.NewRequest(http.MethodPost, test.request, strings.NewReader(test.URL))
			writer := httptest.NewRecorder()

			PostHandler(writer, request)

			res := writer.Result()

			assert.Equal(t, test.want.statusCode, res.StatusCode)
			assert.Equal(t, test.want.contentType, res.Header.Get("Content-Type"))

			userResult, err := io.ReadAll(res.Body)
			require.NoError(t, err)

			err = res.Body.Close()
			require.NoError(t, err)

			if test.isError {
				assert.Equal(t, test.want.message, string(userResult))
			} else {
				requiredOutput := fmt.Sprintf("http://%s/%s", request.Host, test.want.message)
				assert.Equal(t, requiredOutput, string(userResult))

				url, ok := urls[test.want.message]
				require.True(t, ok)
				assert.Equal(t, url, test.URL)
			}
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
		urls    map[string]string
		want    want
	}{
		{
			name:    "correct input data",
			request: "/aHR0cHM6",

			urls: map[string]string{
				"aHR0cHM6": "https://practicum.yandex.ru/",
			},

			want: want{
				statusCode:  http.StatusTemporaryRedirect,
				contentType: "text/plain",
				location:    "https://practicum.yandex.ru/",
				message:     "",
			},
		},
		{
			name:    "request without id",
			request: "/",

			urls: map[string]string{
				"aHR0cHM6": "https://practicum.yandex.ru/",
			},

			want: want{
				statusCode:  http.StatusBadRequest,
				contentType: "text/plain",
				location:    "",
				message:     ErrInvalidGivenURL.Error(),
			},
		},
		{
			name:    "missing URL",
			request: "/fsdfuytu",

			urls: map[string]string{
				"aHR0cHM6": "https://practicum.yandex.ru/",
			},

			want: want{
				statusCode:  http.StatusBadRequest,
				contentType: "text/plain",
				location:    "",
				message:     ShortURLNotInDB,
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			urls = test.urls

			request := httptest.NewRequest(http.MethodGet, test.request, nil)
			writer := httptest.NewRecorder()

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
