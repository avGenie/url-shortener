package auth

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/avGenie/url-shortener/internal/app/auth/mock"
	"github.com/avGenie/url-shortener/internal/app/entity"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAuthMiddlewareInvalidUserID(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	s := mock.NewMockUserAuthorisator(ctrl)

	nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t.Errorf("reach next handler")
	})

	type want struct {
		statusCode   int
		cookieExact  bool //opt??
		cookieName   string
		errorAddUser error
	}
	tests := []struct {
		name         string
		userIDCookie *http.Cookie
		want         want
	}{
		{
			name:         "empty cookie",
			userIDCookie: nil,
			want: want{
				statusCode:   http.StatusUnauthorized,
				cookieExact:  true,
				cookieName:   entity.UserIDKey,
				errorAddUser: nil,
			},
		},
		{
			name: "invalid cookie",
			userIDCookie: &http.Cookie{
				Name:  entity.UserIDKey,
				Value: "",
			},
			want: want{
				statusCode:   http.StatusUnauthorized,
				cookieExact:  true,
				cookieName:   entity.UserIDKey,
				errorAddUser: nil,
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodPost, "/", nil)
			w := httptest.NewRecorder()

			s.EXPECT().AddUser(gomock.Any(), gomock.Any()).Return(test.want.errorAddUser)

			authHandler := AuthMiddleware(s)
			handler := authHandler(nextHandler)
			handler.ServeHTTP(w, req)

			res := w.Result()

			err := res.Body.Close()
			require.NoError(t, err)

			assert.Equal(t, test.want.statusCode, res.StatusCode)

			cookies := req.Cookies()
			require.NotEmpty(t, cookies)

			require.Equal(t, test.want.cookieName, cookies[0].Name)
			require.NotEmpty(t, cookies[0].Value)
		})
	}
}
