package auth

import (
	"fmt"

	"github.com/avGenie/url-shortener/internal/app/entity"
)

// ExampleDecodeUserID Example of decoding user ID
func ExampleDecodeUserID() {
	uuid := entity.UserID("ac2a4811-4f10-487f-bde3-e39a14af7cd8")

	encoded, _ := EncodeUserID(uuid)

	decoded, _ := DecodeUserID(encoded)
	fmt.Println(decoded)

	// Output:
	// ac2a4811-4f10-487f-bde3-e39a14af7cd8
}
