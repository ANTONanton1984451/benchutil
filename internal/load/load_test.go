package load

import (
	"benchutil/pkg/httploader"
	"context"
	"errors"
	"fmt"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
	"net/http"
	"testing"
)

func TestMakeLoadFunc(t *testing.T) {

	type testCase struct {
		name           string
		setupFunc      func(mc *gomock.Controller) *httploader.MockLoader
		cfg            config
		expectedResult httploader.Report
		expectedErr    error
	}

	initDefaultMocks := func(mc *gomock.Controller) *httploader.MockLoader {
		return httploader.NewMockLoader(mc)
	}

	var testErr = errors.New("test error")

	ctx := context.Background()

	cases := [...]testCase{
		{
			name: "error http load",
			setupFunc: func(mc *gomock.Controller) *httploader.MockLoader {
				loader := httploader.NewMockLoader(mc)

				loader.EXPECT().Load(ctx, "", nil, nil).Return(httploader.Report{}, testErr)

				return loader
			},
			expectedErr: fmt.Errorf("load: %w", testErr),
		},
		{
			name: "ok, without headers and body",
			setupFunc: func(mc *gomock.Controller) *httploader.MockLoader {
				loader := httploader.NewMockLoader(mc)

				loader.EXPECT().Load(ctx, "hostload", nil, nil).Return(httploader.Report{}, nil)
				return loader
			},
			expectedResult: httploader.Report{},
			cfg:            config{host: "hostload"},
		},
		{
			name: "ok, with headers",
			setupFunc: func(mc *gomock.Controller) *httploader.MockLoader {
				loader := httploader.NewMockLoader(mc)

				loader.EXPECT().Load(ctx, "", &http.Header{"Test": []string{"header1"}}, nil).Return(httploader.Report{}, nil)

				return loader
			},
			cfg: config{headersPath: "testdata/headers.json"},
		},
		{
			name: "ok, with body",
			setupFunc: func(mc *gomock.Controller) *httploader.MockLoader {
				loader := httploader.NewMockLoader(mc)

				loader.EXPECT().Load(ctx, "", nil, []byte(`i am body`)).Return(httploader.Report{}, nil)

				return loader
			},
			cfg: config{bodyPath: "testdata/body.txt"},
		},
		{
			name: "ok with headers and body",
			setupFunc: func(mc *gomock.Controller) *httploader.MockLoader {
				loader := httploader.NewMockLoader(mc)

				loader.EXPECT().Load(ctx, "loadhost", &http.Header{"Test": []string{"header1"}}, []byte(`i am body`)).Return(httploader.Report{}, nil)

				return loader
			},
			cfg: config{host: "loadhost", bodyPath: "testdata/body.txt", headersPath: "testdata/headers.json"},
		},
		{
			name: "ok, return not changed report",
			setupFunc: func(mc *gomock.Controller) *httploader.MockLoader {
				loader := httploader.NewMockLoader(mc)

				loader.EXPECT().Load(ctx, "", nil, nil).Return(httploader.Report{Success: 10, Cancelled: 11, Errors: 12, All: 130}, nil)

				return loader
			},
			expectedResult: httploader.Report{Success: 10, Cancelled: 11, Errors: 12, All: 130},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {

			mc := gomock.NewController(t)
			loaderMock := initDefaultMocks(mc)
			if tc.setupFunc != nil {
				loaderMock = tc.setupFunc(mc)
			}

			rep, err := makeLoad(ctx, tc.cfg, loaderMock)

			require.Equal(t, tc.expectedResult, rep)
			require.Equal(t, tc.expectedErr, err)
		})
	}
}
