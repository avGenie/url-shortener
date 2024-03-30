package auth

import (
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
