package headers

import (
	"errors"
	"github.com/stretchr/testify/require"
	"net/http"
	"testing"
)

func TestReadFromFile(t *testing.T) {
	type testCase struct {
		name        string
		path        string
		expectedRes *http.Header
		expectedErr error
	}

	cases := [...]testCase{
		{
			name:        "unsupported format",
			path:        "testdata/headers/unsupported.xml",
			expectedErr: errors.New("unsupported format xml"),
		},
		{
			name:        "json format,OK",
			path:        "testdata/headers.json",
			expectedRes: &http.Header{"Test": []string{"header1"}, "Test2": []string{"header2"}},
		},
		{
			name:        "yaml format, OK",
			path:        "testdata/headers.yaml",
			expectedRes: &http.Header{"Testyaml": []string{"header1"}, "Testyaml2": []string{"header2"}},
		},
		{
			name:        "yml format, OK",
			path:        "testdata/headers.yml",
			expectedRes: &http.Header{"Testyml": []string{"header1"}, "Testyml2": []string{"header2"}},
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			h, err := ReadFromFile(c.path)
			require.Equal(t, c.expectedRes, h, "result not equal")
			require.Equal(t, c.expectedErr, err, "error not equal")
		})
	}
}
