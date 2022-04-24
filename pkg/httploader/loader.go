//go:generate mockgen -source=./loader.go -destination=./loader_mock.go -package=httploader

package httploader

import (
	"context"
	"net/http"
	"time"
)

// Report отчёт по нагрузке на сервер
// AvgResponseTime в секундах, если ответ был меньше 0.5 секунд, то в AvgResponseTime будет равен 0
type Report struct {
	Success         int
	Cancelled       int
	Errors          int
	All             int
	AvgResponseTime time.Duration
}

type Loader interface {
	Load(ctx context.Context, host string, headers *http.Header, body []byte) (Report, error)
}

// New создание инстанса объекта, поддерживающего Loader
// аргумент с - количество одновременных запросов к серверу
// если аргумент c будет больше 1, то будет concurrency Loader
func New(timeOut time.Duration, method string, requests, c int) Loader {
	consistentLoader := consistent{
		method:   method,
		requests: requests,
		timeout:  timeOut,
	}

	if c > 1 {
		return &concurrency{consistent: consistentLoader, requestsPerTime: c}
	} else {
		return &consistentLoader
	}
}
