package handlers

import (
	"fmt"
	"io"
	"net"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gopkg.in/h2non/gentleman.v2"
)

const (
	testURL = "127.0.0.1:8080"
)

func testCreateServer(t *testing.T) *httptest.Server {
	ts := httptest.NewUnstartedServer(CreateRouter())
	l, err := net.Listen("tcp", testURL)
	require.NoError(t, err)

	ts.Listener.Close()
	ts.Listener = l

	return ts
}

func testRequest(t *testing.T, request string) *gentleman.Response {
	cli := gentleman.New()
	cli.URL(testURL)
	req := cli.Request()
	req.Method("GET")
	req.Path(request)

	// Perform the request
	res, err := req.Send()
	require.NoError(t, err)

	return res
}

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
	ts := testCreateServer(t)
	defer ts.Close()
	ts.Start()

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
				statusCode:  http.StatusMethodNotAllowed,
				contentType: "",
				location:    "",
				message:     "",
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
		urls = test.urls

		res := testRequest(t, test.request)

		assert.Equal(t, test.want.statusCode, res.StatusCode)
		assert.Equal(t, test.want.contentType, res.Header.Get("Content-Type"))
		assert.Equal(t, test.want.location, res.Header.Get("Location"))
		assert.Equal(t, test.want.message, res.String())
	}
}
