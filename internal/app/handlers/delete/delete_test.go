package handlers

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/avGenie/url-shortener/internal/app/entity"
	"github.com/avGenie/url-shortener/internal/app/handlers/delete/mock"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestDeleteHandler(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	s := mock.NewMockAllURLDeleter(ctrl)

	inputBatch := `["42b3e75f", "77fca595", "ac6bb669"]`
	invalidBatch := `<invalid json>`

	type want struct {
		statusCode int
		expectErr  error
	}
	tests := []struct {
		name              string
		inputBody         string
		userIDCtx         entity.UserIDCtx
		exitBeforeGetting bool
		want              want
	}{
		{
			name:      "correct input data",
			inputBody: inputBatch,
			userIDCtx: entity.UserIDCtx{
				UserID:     "ac2a4811-4f10-487f-bde3-e39a14af7cd8",
				StatusCode: http.StatusOK,
			},

			want: want{
				statusCode: http.StatusAccepted,
			},
		},
		{
			name:      "error json processing",
			inputBody: invalidBatch,
			userIDCtx: entity.UserIDCtx{
				UserID:     "ac2a4811-4f10-487f-bde3-e39a14af7cd8",
				StatusCode: http.StatusOK,
			},
			exitBeforeGetting: true,

			want: want{
				statusCode: http.StatusBadRequest,
			},
		},
		{
			name:      "empty user id",
			inputBody: inputBatch,
			userIDCtx: entity.UserIDCtx{
				UserID:     "",
				StatusCode: http.StatusOK,
			},
			exitBeforeGetting: true,

			want: want{
				statusCode: http.StatusInternalServerError,
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			request := httptest.NewRequest(http.MethodDelete, "/api/user/urls", strings.NewReader(test.inputBody))
			writer := httptest.NewRecorder()

			if test.exitBeforeGetting {
				s.EXPECT().DeleteBatchURL(gomock.Any(), gomock.Any()).Times(0)
			} else {
				s.EXPECT().DeleteBatchURL(gomock.Any(), gomock.Any()).
					Return(test.want.expectErr)
			}

			request = request.WithContext(context.WithValue(request.Context(), entity.UserIDCtxKey{}, test.userIDCtx))

			deleteHandler := NewDeleteHandler(s)
			handler := deleteHandler.DeleteUserURLHandler()
			handler(writer, request)

			res := writer.Result()
			assert.NoError(t, res.Body.Close())

			assert.Equal(t, test.want.statusCode, res.StatusCode)
		})
	}
}
