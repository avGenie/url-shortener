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
