package handlers

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/avGenie/url-shortener/internal/app/config"
	"github.com/avGenie/url-shortener/internal/app/entity"
	"github.com/avGenie/url-shortener/internal/app/handlers/errors"
	"github.com/avGenie/url-shortener/internal/app/handlers/get/mock"
	"github.com/avGenie/url-shortener/internal/app/logger"
	"github.com/go-chi/chi/v5"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
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
		userID            entity.UserID
		exitBeforeGetting bool
		want              want
	}{
		{
			name:    "correct input data",
			request: "aHR0cHM6",
			userID:  entity.UserID("ac2a4811-4f10-487f-bde3-e39a14af7cd8"),

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
			userID:  entity.UserID("ac2a4811-4f10-487f-bde3-e39a14af7cd8"),

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
			userID:  entity.UserID("ac2a4811-4f10-487f-bde3-e39a14af7cd8"),

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
			name:    "missing user id",
			request: "/aHR0cHM6",
			userID:  entity.UserID(""),
			exitBeforeGetting: true,

			want: want{
				statusCode:  http.StatusInternalServerError,
				contentType: "",
				location:    "",
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
			if len(test.userID) != 0 {
				request = request.WithContext(context.WithValue(request.Context(), entity.UserIDCtxKey{}, test.userID))
			}

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

			// pingDB(s, writer, request)
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
