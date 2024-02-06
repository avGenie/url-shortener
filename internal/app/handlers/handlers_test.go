package handlers

import (
	"fmt"
	"io"
	"net"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/avGenie/url-shortener/internal/app/config"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gopkg.in/h2non/gentleman.v2"
)

func testCreateServer(t *testing.T) *httptest.Server {
	ts := httptest.NewUnstartedServer(CreateRouter())
	l, err := net.Listen("tcp", config.NetAddr)
	require.NoError(t, err)

	ts.Listener.Close()
	ts.Listener = l

	return ts
}

func testGetResponse(t *testing.T, request string) *gentleman.Response {
	cli := gentleman.New()
	cli.URL(config.NetAddr)
	req := cli.Request()
	req.Method("GET")
	req.Path(request)

	// Perform the request
	res, err := req.Send()
	require.NoError(t, err)

	return res
}

func testPostResponse(t *testing.T, ts *httptest.Server, method, path, data string) *http.Response {
	req, err := http.NewRequest(method, ts.URL+path, strings.NewReader(data))
	require.NoError(t, err)

	resp, err := ts.Client().Do(req)
	require.NoError(t, err)

	return resp
}

func TestPostHandler(t *testing.T) {
	ts := testCreateServer(t)
	defer ts.Close()
	ts.Start()

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
				message:     EmptyURL + "\n",
			},
		},
	}

	for _, test := range tests {
		res := testPostResponse(t, ts, "POST", test.request, test.URL)

		respBody, err := io.ReadAll(res.Body)
		require.NoError(t, err)

		err = res.Body.Close()
		require.NoError(t, err)

		assert.Equal(t, test.want.statusCode, res.StatusCode)
		assert.Equal(t, test.want.contentType, res.Header.Get("Content-Type"))

		if test.isError {
			assert.Equal(t, test.want.message, string(respBody))
		} else {
			requiredOutput := fmt.Sprintf("%s/%s", config.BaseURIPrefix, test.want.message)
			assert.Equal(t, requiredOutput, string(respBody))

			url, ok := urls[test.want.message]
			require.True(t, ok)
			assert.Equal(t, url, test.URL)
		}
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
				contentType: "text/plain; charset=utf-8",
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
				contentType: "text/plain; charset=utf-8",
				location:    "",
				message:     ShortURLNotInDB + "\n",
			},
		},
	}

	for _, test := range tests {
		urls = test.urls

		res := testGetResponse(t, test.request)

		assert.Equal(t, test.want.statusCode, res.StatusCode)
		assert.Equal(t, test.want.contentType, res.Header.Get("Content-Type"))
		assert.Equal(t, test.want.location, res.Header.Get("Location"))
		assert.Equal(t, test.want.message, res.String())
	}
}
