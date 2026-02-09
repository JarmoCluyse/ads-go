package adserrors_test

import(
	"github.com/jarmocluyse/ads-go/pkg/ads/ads-errors"
	"testing"
	"github.com/stretchr/testify/assert"
)

func TestErrorCodeToString(t *testing.T) {
	tests := []struct {
		name string // description of this test case
		// Named input parameters for target function.
		code uint32
		want string
	}{
		{
			name: "Known error code - No error",
			code: 0,
			want: "No error",
		},
		{
			name: "Known error code - Internal error",
			code: 1,
			want: "Internal error",
		},
		{
			name: "Known error code - General device error",
			code: 1792,
			want: "General device error",
		},
		{
			name: "Unknown error code",
			code: 999999,
			want: "Unknown ads error 999999",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := adserrors.ErrorCodeToString(tt.code)
			assert.Equal(t, tt.want, got, "ErrorCodeToString(%v)", tt.code)
		})
	}
}

