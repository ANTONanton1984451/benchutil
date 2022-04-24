package httploader

import (
	"bytes"
	"context"
	"fmt"
	"net/http"
	"net/url"
	"sync"
	"time"
)

type concurrency struct {
	consistent
	requestsPerTime int
}

type concurrencyResp struct {
	cancelled, success, error bool
	respTime                  time.Duration
}

// Load отсылает параллельные запросы к host
// при прерывании контекстом перестаёт слать запросы и дождидается выполнения всех, уже запущенных запросов
// поддерживает graceful shutdown
func (l *concurrency) Load(ctx context.Context, host string, headers *http.Header, body []byte) (Report, error) {
	req, err := http.NewRequest(l.method, host, bytes.NewBuffer(body))
	if err != nil {
		return Report{}, fmt.Errorf("create request: %w", err)
	}

	if headers != nil {
		req.Header = *headers
	}
	throttle := make(chan struct{}, l.requestsPerTime)
	wg := sync.WaitGroup{}

	responseList := make([]*concurrencyResp, l.requests)
	for i := 0; i < l.requests; i++ {

		select {
		case <-ctx.Done():
			goto wait
		default:
		}

		throttle <- struct{}{}
		wg.Add(1)
		go func(index int) {
			defer func() {
				wg.Done()
				<-throttle
			}()

			cli := http.Client{Timeout: l.timeout}
			now := time.Now()

			reqResult := concurrencyResp{}
			resp, err := cli.Do(req)
			if err != nil {
				if urlErr, ok := err.(*url.Error); ok {
					if urlErr.Timeout() {
						reqResult.cancelled = true
					} else {
						reqResult.error = true
					}
				} else {
					reqResult.error = true
				}
				responseList[index] = &reqResult
				return
			}
			resp.Body.Close()

			if resp.StatusCode != http.StatusOK {
				reqResult.error = true
			}

			reqResult.success = true
			reqResult.respTime = time.Since(now)
			responseList[index] = &reqResult
		}(i)
	}

wait:
	wg.Wait()
	close(throttle)

	return l.calcReport(responseList), nil
}

func (l *concurrency) calcReport(respList []*concurrencyResp) Report {
	var (
		all       int
		success   int
		cancelled int
		errored   int

		avgRespTime float64
	)

	for _, res := range respList {
		if res != nil {
			all++
			if res.cancelled {
				cancelled++
				continue
			}
			if res.error {
				errored++
				continue
			}

			if res.success {
				success++

				avgRespTime += res.respTime.Seconds()
			}
		}
	}

	return Report{
		Success:   success,
		Cancelled: cancelled,
		Errors:    errored,
		All:       all,

		AvgResponseTime: calcResponseTime(success, avgRespTime),
	}
}
