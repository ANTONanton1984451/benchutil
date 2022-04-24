package httploader

import (
	"context"
	"github.com/stretchr/testify/require"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestLoadConsistent(t *testing.T) {

	t.Run("happy path all requests are success", func(t *testing.T) {
		okHandler := http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
			writer.WriteHeader(http.StatusOK)
			writer.Write([]byte("OK"))
		})

		serv := httptest.NewServer(okHandler)
		defer serv.Close()

		ctx := context.Background()
		loader := consistent{requests: 10, method: http.MethodGet}
		expectedRep := Report{
			All:     10,
			Success: 10,
		}

		rep, err := loader.Load(ctx, serv.URL, nil, nil)

		require.Equal(t, nil, err)
		require.Equal(t, expectedRep, rep)
	})

	t.Run("50 percents is errors", func(t *testing.T) {
		count := 0

		okHandler := http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
			if count < 5 {
				writer.WriteHeader(http.StatusOK)
			} else {
				writer.WriteHeader(http.StatusInternalServerError)
			}
			count++
		})

		serv := httptest.NewServer(okHandler)
		defer serv.Close()

		ctx := context.Background()
		loader := consistent{requests: 10, method: http.MethodGet}
		expectedRep := Report{
			All:     10,
			Success: 5,
			Errors:  5,
		}

		rep, err := loader.Load(ctx, serv.URL, nil, nil)

		require.Equal(t, nil, err)
		require.Equal(t, expectedRep, rep)
	})

	t.Run("all requests are cancelled", func(t *testing.T) {
		if testing.Short() {
			t.Skip("test execute over 10 sec")
		}
		timeOut := 1

		okHandler := http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
			time.Sleep(time.Second * time.Duration(timeOut) * 2)
			writer.WriteHeader(http.StatusOK)
		})

		serv := httptest.NewServer(okHandler)
		defer serv.Close()

		ctx := context.Background()
		loader := consistent{requests: 10, method: http.MethodGet, timeout: time.Duration(timeOut) * time.Second}
		expectedRep := Report{
			All:       10,
			Cancelled: 10,
		}

		rep, err := loader.Load(ctx, serv.URL, nil, nil)

		require.Equal(t, nil, err)
		require.Equal(t, expectedRep, rep)
	})

	t.Run("cancel context, all requests complete", func(t *testing.T) {
		if testing.Short() {
			t.Skip("test execute over 10 sec")
		}
		timeOut := 4
		serverResponseTime := 2
		loaderTimeOut := 4

		okHandler := http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
			time.Sleep(time.Second * time.Duration(timeOut*serverResponseTime))
			writer.WriteHeader(http.StatusOK)
		})

		serv := httptest.NewServer(okHandler)
		defer serv.Close()

		ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
		loader := consistent{requests: 10, method: http.MethodGet, timeout: time.Duration(timeOut*loaderTimeOut) * time.Second}
		expectedRep := Report{
			All:             2,
			Success:         2,
			AvgResponseTime: time.Second * time.Duration(timeOut*2),
		}

		rep, err := loader.Load(ctx, serv.URL, nil, nil)

		require.Equal(t, nil, err)
		require.Equal(t, expectedRep, rep)
	})

	t.Run("cancel context, all requests cancelled", func(t *testing.T) {
		if testing.Short() {
			t.Skip("test execute over 10 sec")
		}
		timeOut := 4
		serverResponseTime := 4
		loaderTimeOut := 2

		okHandler := http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
			time.Sleep(time.Second * time.Duration(timeOut*serverResponseTime))
			writer.WriteHeader(http.StatusOK)
		})

		serv := httptest.NewServer(okHandler)
		defer serv.Close()

		ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
		loader := consistent{requests: 10, method: http.MethodGet, timeout: time.Duration(timeOut*loaderTimeOut) * time.Second}
		expectedRep := Report{
			All:       2,
			Cancelled: 2,
		}

		rep, err := loader.Load(ctx, serv.URL, nil, nil)

		require.Equal(t, nil, err)
		require.Equal(t, expectedRep, rep)
	})
}
