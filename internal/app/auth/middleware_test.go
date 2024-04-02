package auth

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/avGenie/url-shortener/internal/app/auth/mock"
	"github.com/avGenie/url-shortener/internal/app/entity"
	db_err "github.com/avGenie/url-shortener/internal/app/storage/api/errors"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAuthMiddleware(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	s := mock.NewMockUserAuthorisator(ctrl)

	type want struct {
		statusCode    int
		cookieName    string
		userID        string
		errorAuthUser error
	}
	tests := []struct {
		name         string
		userIDCookie *http.Cookie
		want         want
	}{
		{
			name: "correct cookie",
			userIDCookie: &http.Cookie{
				Name:  entity.UserIDKey,
				Value: "ac2a4811-4f10-487f-bde3-e39a14af7cd8",
			},
			want: want{
				statusCode:    http.StatusOK,
				cookieName:    entity.UserIDKey,
				userID:        "ac2a4811-4f10-487f-bde3-e39a14af7cd8",
				errorAuthUser: nil,
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodPost, "/", nil)
			w := httptest.NewRecorder()

			req.AddCookie(test.userIDCookie)

			s.EXPECT().AddUser(gomock.Any(), gomock.Any()).Times(0)
			s.EXPECT().AuthUser(gomock.Any(), gomock.Any()).Return(entity.UserID(test.want.userID), test.want.errorAuthUser)

			nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				userID, ok := r.Context().Value(entity.UserIDCtxKey{}).(entity.UserID)
				require.True(t, ok)
				assert.Equal(t, userID.String(), test.want.userID)
			})

			authHandler := AuthMiddleware(s)
			handler := authHandler(nextHandler)
			handler.ServeHTTP(w, req)

			res := w.Result()

			err := res.Body.Close()
			require.NoError(t, err)

			assert.Equal(t, test.want.statusCode, res.StatusCode)
		})
	}
}

func TestAuthMiddlewareInvalidAuthUser(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	s := mock.NewMockUserAuthorisator(ctrl)

	nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t.Errorf("reach next handler")
	})

	type want struct {
		statusCode    int
		userID        string
		errorAuthUser error
	}
	tests := []struct {
		name         string
		cookieName   string
		userIDCookie *http.Cookie
		isAddUser    bool
		want         want
	}{
		{
			name:       "unknown user",
			cookieName: entity.UserIDKey,
			userIDCookie: &http.Cookie{
				Name:  entity.UserIDKey,
				Value: "ac2a4811-4f10-487f-bde3-e39a14af7cd8",
			},
			isAddUser: true,
			want: want{
				statusCode:    http.StatusUnauthorized,
				userID:        "ac2a4811-4f10-487f-bde3-e39a14af7cd8",
				errorAuthUser: db_err.ErrUserIDNotFound,
			},
		},
		{
			name:       "storage error",
			cookieName: entity.UserIDKey,
			userIDCookie: &http.Cookie{
				Name:  entity.UserIDKey,
				Value: "ac2a4811-4f10-487f-bde3-e39a14af7cd8",
			},
			isAddUser: false,
			want: want{
				statusCode:    http.StatusInternalServerError,
				userID:        "ac2a4811-4f10-487f-bde3-e39a14af7cd8",
				errorAuthUser: errors.New("error"),
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodPost, "/", nil)
			w := httptest.NewRecorder()

			req.AddCookie(test.userIDCookie)

			s.EXPECT().AuthUser(gomock.Any(), gomock.Any()).Return(entity.UserID(test.want.userID), test.want.errorAuthUser)
			if test.isAddUser {
				s.EXPECT().AddUser(gomock.Any(), gomock.Any()).Return(nil)
			}

			authHandler := AuthMiddleware(s)
			handler := authHandler(nextHandler)
			handler.ServeHTTP(w, req)

			res := w.Result()

			err := res.Body.Close()
			require.NoError(t, err)

			assert.Equal(t, test.want.statusCode, res.StatusCode)

			if !test.isAddUser {
				require.Empty(t, res.Cookies())

				return
			}

			cookies := res.Cookies()
			require.NotEmpty(t, cookies)
			require.Equal(t, cookies[0].Name, test.cookieName)
			require.NotEmpty(t, cookies[0].Value)
		})
	}
}

func TestAuthMiddlewareInvalidAddUser(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	s := mock.NewMockUserAuthorisator(ctrl)

	nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t.Errorf("reach next handler")
	})

	type want struct {
		statusCode   int
		errorAddUser error
	}
	tests := []struct {
		name       string
		cookieName string
		want       want
	}{
		{
			name:       "empty cookie",
			cookieName: entity.UserIDKey,
			want: want{
				statusCode:   http.StatusUnauthorized,
				errorAddUser: nil,
			},
		},
		{
			name:       "invalid cookie",
			cookieName: entity.UserIDKey,
			want: want{
				statusCode:   http.StatusUnauthorized,
				errorAddUser: nil,
			},
		},
		{
			name:       "add user error",
			cookieName: entity.UserIDKey,
			want: want{
				statusCode:   http.StatusInternalServerError,
				errorAddUser: errors.New("error"),
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

			if test.want.errorAddUser != nil {
				require.Empty(t, res.Cookies())

				return
			}

			cookies := res.Cookies()
			require.NotEmpty(t, cookies)
			require.Equal(t, cookies[0].Name, test.cookieName)
			require.NotEmpty(t, cookies[0].Value)
		})
	}
}
