package auth

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/avGenie/url-shortener/internal/app/entity"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAuthMiddleware(t *testing.T) {
	type want struct {
		nextHandler http.HandlerFunc
	}
	tests := []struct {
		userIDCookie *http.Cookie
		want         want
		name         string
	}{
		{
			name: "correct cookie",
			userIDCookie: &http.Cookie{
				Name:  entity.UserIDKey,
				Value: "ac2a4811-4f10-487f-bde3-e39a14af7cd8",
			},
			want: want{
				nextHandler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					userIDCtx, ok := r.Context().Value(entity.UserIDCtxKey{}).(entity.UserIDCtx)

					require.True(t, ok)
					assert.Equal(t, userIDCtx.UserID.String(), "ac2a4811-4f10-487f-bde3-e39a14af7cd8")
					assert.Equal(t, userIDCtx.StatusCode, http.StatusOK)
				}),
			},
		},
		{
			name: "empty cookie",
			want: want{
				nextHandler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					userIDCtx, ok := r.Context().Value(entity.UserIDCtxKey{}).(entity.UserIDCtx)

					require.True(t, ok)
					assert.NotEmpty(t, userIDCtx.UserID.String())
					assert.Equal(t, userIDCtx.StatusCode, http.StatusUnauthorized)
				}),
			},
		},
		{
			name: "invalid cookie",
			userIDCookie: &http.Cookie{
				Name:  entity.UserIDKey,
				Value: "",
			},
			want: want{
				nextHandler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					userIDCtx, ok := r.Context().Value(entity.UserIDCtxKey{}).(entity.UserIDCtx)

					require.True(t, ok)
					assert.Empty(t, userIDCtx.UserID.String())
					assert.Equal(t, userIDCtx.StatusCode, http.StatusUnauthorized)
				}),
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodPost, "/", nil)
			w := httptest.NewRecorder()

			if test.userIDCookie != nil {
				req.AddCookie(test.userIDCookie)
			}

			handler := AuthMiddleware(test.want.nextHandler)
			handler.ServeHTTP(w, req)
		})
	}
}
