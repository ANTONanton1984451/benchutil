package httploader

import (
	"context"
	"github.com/stretchr/testify/require"
	"net/http"
	"net/http/httptest"
	"sync/atomic"
	"testing"
	"time"
)

func TestConcurrencyLoader(t *testing.T) {

	t.Run("happy path all requests are success", func(t *testing.T) {
		okHandler := http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
			writer.WriteHeader(http.StatusOK)
			writer.Write([]byte("OK"))
		})

		serv := httptest.NewServer(okHandler)
		defer serv.Close()

		ctx := context.Background()
		loader := concurrency{consistent: consistent{timeout: 5 * time.Second, requests: 10, method: http.MethodGet}, requestsPerTime: 5}
		expectedRep := Report{
			All:     10,
			Success: 10,
		}

		rep, err := loader.Load(ctx, serv.URL, nil, nil)

		require.Equal(t, nil, err)
		require.Equal(t, expectedRep, rep)
	})

	t.Run("all requests are cancelled", func(t *testing.T) {
		timeOut := 1

		okHandler := http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
			time.Sleep(time.Second * time.Duration(timeOut) * 2)
			writer.WriteHeader(http.StatusOK)
		})

		serv := httptest.NewServer(okHandler)
		defer serv.Close()

		ctx := context.Background()
		loader := concurrency{consistent: consistent{requests: 10, method: http.MethodGet, timeout: time.Duration(timeOut) * time.Second}, requestsPerTime: 10}
		expectedRep := Report{
			All:       10,
			Cancelled: 10,
		}

		rep, err := loader.Load(ctx, serv.URL, nil, nil)

		require.Equal(t, nil, err)
		require.Equal(t, expectedRep, rep)
	})

	t.Run("50 percents is error", func(t *testing.T) {
		errorGenerator := responser{
			eachErrorResponse: 2,
		}

		okHandler := http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
			if errorGenerator.responseError() {
				writer.WriteHeader(http.StatusInternalServerError)
			} else {
				writer.WriteHeader(http.StatusOK)
			}
		})

		serv := httptest.NewServer(okHandler)
		defer serv.Close()

		ctx := context.Background()
		loader := concurrency{consistent: consistent{requests: 20, method: http.MethodGet, timeout: time.Second}, requestsPerTime: 10}
		expectedRep := Report{
			All:     20,
			Success: 10,
			Errors:  10,
		}

		rep, err := loader.Load(ctx, serv.URL, nil, nil)

		require.Equal(t, nil, err)
		require.Equal(t, expectedRep, rep)
	})

	t.Run("cancel context", func(t *testing.T) {
		if testing.Short() {
			t.Skipf("test execute over 4 seconds")
		}

		timeOut := 1
		serverResponseTime := 4
		reqTimeOut := 5
		reqPerTime := 5
		ctxTimeOut := 3

		okHandler := http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
			time.Sleep(time.Second * time.Duration(timeOut*serverResponseTime))
			writer.WriteHeader(http.StatusOK)
		})

		serv := httptest.NewServer(okHandler)
		defer serv.Close()

		ctx, _ := context.WithTimeout(context.Background(), time.Duration(ctxTimeOut)*time.Second)
		loader := concurrency{consistent: consistent{requests: 10, method: http.MethodGet, timeout: time.Duration(timeOut*reqTimeOut) * time.Second}, requestsPerTime: reqPerTime}
		expectedRep := Report{
			All:             6,
			Success:         6,
			AvgResponseTime: 4 * time.Second,
		}

		rep, err := loader.Load(ctx, serv.URL, nil, nil)

		require.Equal(t, nil, err)
		require.Equal(t, expectedRep, rep)
	})

}

type responser struct {
	eachErrorResponse int
	counter           uint64
}

func (r *responser) responseError() bool {
	atomic.AddUint64(&r.counter, 1)
	return r.counter%uint64(r.eachErrorResponse) == 0
}
