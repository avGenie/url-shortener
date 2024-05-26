package auth

import (
	"github.com/avGenie/url-shortener/internal/app/usecase/user"
	"go.uber.org/zap"
	"testing"

	"github.com/avGenie/url-shortener/internal/app/entity"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestEncodeDecode(t *testing.T) {
	const testCount = 10

	for i := 0; i < testCount; i++ {
		rawUUID := uuid.New()

		encodedUUID, err := EncodeUserID(entity.UserID(rawUUID.String()))
		require.NoError(t, err)

		decodedUUID, err := DecodeUserID(encodedUUID)
		require.NoError(t, err)

		assert.Equal(t, rawUUID.String(), decodedUUID.String())
	}
}

func BenchmarkEncodeDecode(b *testing.B) {
	elems := make([]string, 0, b.N)

	b.Run("encode", func(b *testing.B) {
		for n := 0; n < b.N; n++ {
			userID := user.CreateUserID()
			encodedUUID, err := EncodeUserID(userID)
			if err != nil {
				zap.L().Error("benchmark encode error", zap.Error(err))
				continue
			}

			elems = append(elems, encodedUUID)
		}
	})

	b.Run("decode", func(b *testing.B) {
		for _, elem := range elems {
			_, err := DecodeUserID(elem)
			if err != nil {
				zap.L().Error("benchmark decode error", zap.Error(err))
			}
		}
	})
}
