package httploader

import (
	"bytes"
	"context"
	"fmt"
	"math"
	"net/http"
	"net/url"
	"time"
)

type consistent struct {
	timeout  time.Duration
	method   string
	requests int
}

// Load посылает последовательный запрос к host
// перед отправкой каждого запроса проверяет контекст, что делает метод способным поддерживать graceful shutdown
func (l *consistent) Load(ctx context.Context, host string, headers *http.Header, body []byte) (Report, error) {
	req, err := http.NewRequest(l.method, host, bytes.NewBuffer(body))
	if err != nil {
		return Report{}, fmt.Errorf("create request: %w", err)
	}

	if headers != nil {
		req.Header = *headers
	}

	var (
		canceled int
		success  int
		errors   int
		all      int

		responseTime float64
	)

	for all < l.requests {
		select {
		case <-ctx.Done():
			return l.formReport(success, canceled, errors, all, responseTime), nil
		default:
		}
		all++
		cli := http.Client{Timeout: l.timeout}

		now := time.Now()
		resp, err := cli.Do(req)
		if err != nil {
			if urlErr, ok := err.(*url.Error); ok {
				if urlErr.Timeout() {
					canceled++
				} else {
					errors++
				}
			} else {
				errors++
			}
			continue
		}

		if resp.StatusCode != http.StatusOK {
			errors++
			continue
		}
		responseTime += time.Since(now).Seconds()

		success++
	}

	return l.formReport(success, canceled, errors, all, responseTime), nil
}

func (l *consistent) formReport(success, canceled, errors, all int, avgRespTime float64) Report {
	respTime := calcResponseTime(success, avgRespTime)
	return Report{
		Success:         success,
		Cancelled:       canceled,
		Errors:          errors,
		All:             all,
		AvgResponseTime: respTime,
	}
}

func calcResponseTime(success int, avgRespTime float64) time.Duration {
	return time.Duration(math.Round(avgRespTime/float64(success))) * time.Second
}
