package cidr

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCIDR(t *testing.T) {
	tests := []struct {
		name           string
		subnet         string
		inputIP        string
		expectedOutput bool
	}{
		{
			name:           "passed IP",
			subnet:         "192.168.1.0/24",
			inputIP:        "192.168.1.14",
			expectedOutput: true,
		},
		{
			name:           "incorrect IP",
			subnet:         "192.168.1.0/24",
			inputIP:        "192.168.2.14",
			expectedOutput: false,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			cidr, err := NewCIDR(test.subnet)
			require.NoError(t, err)

			assert.Equal(t, test.expectedOutput, cidr.Contains(test.inputIP))
		})
	}
}
