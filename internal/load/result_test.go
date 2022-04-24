package load

import (
	"benchutil/pkg/httploader"
	"fmt"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func TestLoaderReportToInternal(t *testing.T) {
	type testCase struct {
		name             string
		expectedInternal report
		loaderReport     httploader.Report
	}

	cases := [...]testCase{
		{
			name: "ok, all field convert",
			loaderReport: httploader.Report{
				All:             1,
				Success:         2,
				Cancelled:       3,
				AvgResponseTime: time.Duration(10) * time.Second,
			},
			expectedInternal: report{
				All:         1,
				Success:     2,
				Canceled:    3,
				AvgRespTime: 10,
			},
		},
		{
			name: "ok, time less then one second",
			loaderReport: httploader.Report{
				All:             1,
				Success:         2,
				Cancelled:       3,
				AvgResponseTime: time.Duration(10),
			},
			expectedInternal: report{
				All:         1,
				Success:     2,
				Canceled:    3,
				AvgRespTime: 0,
			},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			internalRep := loaderReportToInternal(tc.loaderReport)

			require.Equal(t, tc.expectedInternal, internalRep)
		})
	}
}

func TestReportToFormats(t *testing.T) {
	type testCase struct {
		name        string
		rep         report
		format      string
		expectedRes []byte
		expectedErr error
	}

	const (
		jsonOutPut = `{
 "success": 0,
 "canceled": 2,
 "errors": 3,
 "all": 1,
 "avgRespTime": 0
}`
		yamlOutPut = `success: 0
canceled: 123
errors: 9
all: 4
avgRespTime: 0
`

		humanOutPut = `Всего запросов: 5 
Из них 
Успешно: 0 
С ошибкой: 234 
Отменённых: 321 
Среднее время запроса(сек): 0`

		humanAllZeroOutput = `Всего запросов: 0 
Из них 
Успешно: 0 
С ошибкой: 0 
Отменённых: 0 
Среднее время запроса(сек): 0`
	)

	cases := [...]testCase{
		{
			name: "ok, normal json format",
			rep: report{
				All:      1,
				Canceled: 2,
				Errors:   3,
			},
			format:      "json",
			expectedRes: []byte(jsonOutPut),
		},
		{
			name: "ok, normal yaml format",
			rep: report{
				All:      4,
				Canceled: 123,
				Errors:   9,
			},
			format:      "yaml",
			expectedRes: []byte(yamlOutPut),
		},
		{
			name: "ok, normal human format",
			rep: report{
				All:      5,
				Canceled: 321,
				Errors:   234,
			},
			format:      "human",
			expectedRes: []byte(humanOutPut),
		},
		{
			name:        "ok, empty report human format",
			rep:         report{},
			format:      "human",
			expectedRes: []byte(humanAllZeroOutput),
		},
		{
			name:        "error, unknown format",
			rep:         report{},
			format:      "xml",
			expectedErr: fmt.Errorf("unknown format: %s", "xml"),
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			output, err := tc.rep.toBytes(tc.format)

			require.Equal(t, tc.expectedRes, output)
			require.Equal(t, tc.expectedErr, err)
		})
	}
}
