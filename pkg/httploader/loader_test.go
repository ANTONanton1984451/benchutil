package httploader

import (
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func TestNew(t *testing.T) {
	type testCase struct {
		name           string
		expectedLoader Loader

		method   string
		timeOut  time.Duration
		requests int
		c        int
	}

	cases := [...]testCase{
		{
			name:           "get consistent loader",
			c:              0,
			expectedLoader: &consistent{},
		},

		{
			name:           "get concurrency loader",
			c:              10,
			expectedLoader: &concurrency{requestsPerTime: 10},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			res := New(tc.timeOut, tc.method, tc.requests, tc.c)

			require.Equal(t, tc.expectedLoader, res)
		})
	}
}
