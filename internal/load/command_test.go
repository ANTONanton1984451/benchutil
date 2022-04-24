package load

import (
	"errors"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestValidateConfig(t *testing.T) {
	type testCase struct {
		name        string
		cfg         config
		expectedErr error
	}

	cases := [...]testCase{
		{
			name:        "empty host error",
			cfg:         config{requestsCount: 123, timeOut: 111},
			expectedErr: errors.New("empty host"),
		},
		{
			name:        "invalid requests count (0)",
			cfg:         config{host: "host", requestsCount: 0},
			expectedErr: errors.New("invalid requests count value - 0"),
		},
		{
			name:        "invalid requests count (negative)",
			cfg:         config{host: "host", requestsCount: -1},
			expectedErr: errors.New("invalid requests count value - -1"),
		},
		{
			name:        "invalid timeout value (0)",
			cfg:         config{host: "host", requestsCount: 1, timeOut: 0},
			expectedErr: errors.New("invalid timeout value - 0"),
		},
		{
			name:        "invalid timeout value (negative)",
			cfg:         config{host: "host", requestsCount: 1, timeOut: -1},
			expectedErr: errors.New("invalid timeout value - -1"),
		},
		{
			name:        "invalid output format",
			cfg:         config{host: "host", requestsCount: 1, timeOut: 1, outputFormat: "xml"},
			expectedErr: errors.New("invalid output format - xml"),
		},
		{
			name:        "invalid concurrency value (negative)",
			cfg:         config{host: "host", requestsCount: 1, timeOut: 1, outputFormat: "json", concurrency: -1},
			expectedErr: errors.New("invalid concurrency value - -1"),
		},
		{
			name: "OK, concurrency == 0",
			cfg:  config{host: "host", requestsCount: 1, timeOut: 1, outputFormat: "json", concurrency: 0},
		},
		{
			name: "OK, concurrency > 0",
			cfg:  config{host: "host", requestsCount: 1, timeOut: 1, outputFormat: "json", concurrency: 100},
		},
		{
			name: "OK, with body file and headers file",
			cfg:  config{host: "host", requestsCount: 1, timeOut: 1, outputFormat: "json", concurrency: 100, bodyPath: "path/to/body", headersPath: "path/to/headers"},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			err := validateConfig(tc.cfg)
			require.Equal(t, tc.expectedErr, err)
		})
	}
}
