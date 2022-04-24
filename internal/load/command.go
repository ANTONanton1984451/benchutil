package load

import (
	"benchutil/pkg/cli"
	"benchutil/pkg/httploader"
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

type config struct {
	requestsCount int
	concurrency   int
	host          string
	method        string
	bodyPath      string
	headersPath   string
	outputFormat  string
	timeOut       int
}

func New() cli.Command {
	var cfg config
	return cli.Command{
		Name:        "load",
		Description: "Нагружает сервер и даёт отчёт по нагрузке",
		Flags: []cli.CmdFlag{
			cli.IntFlag{
				Name:        "n",
				Destination: &cfg.requestsCount,
				Usage:       "Количество запросов к серверу",
			},
			cli.IntFlag{
				Name:        "c",
				Destination: &cfg.concurrency,
				Usage:       "Количество одновременных запросов к серверу в момент времени",
			},
			cli.IntFlag{
				Name:        "t",
				Destination: &cfg.timeOut,
				Default:     1,
				Usage:       "Таймаут для запросов",
			},
			cli.StringFlag{
				Name:        "host",
				Destination: &cfg.host,
				Usage:       "Url адрес для отправки запросов",
			},
			cli.StringFlag{
				Name:        "m",
				Destination: &cfg.method,
				Default:     http.MethodGet,
				Usage:       "Http метод запроса",
			},
			cli.StringFlag{
				Name:        "b",
				Destination: &cfg.bodyPath,
				Usage:       "Путь до файла с телом запроса",
			},
			cli.StringFlag{
				Name:        "h",
				Destination: &cfg.headersPath,
				Usage:       "Путь до файла с заголовками запроса",
			},
			cli.StringFlag{
				Name:        "o",
				Destination: &cfg.outputFormat,
				Default:     outputHuman,
				Usage:       "Формат вывода результатов нагрузки",
			},
		},
		Action: func(ctx context.Context) error {
			return action(ctx, cfg)
		},
	}

}

func action(ctx context.Context, cfg config) error {
	err := validateConfig(cfg)
	if err != nil {
		return err
	}

	ctx = closer(ctx)
	loader := httploader.New(time.Duration(cfg.timeOut)*time.Second, cfg.method, cfg.requestsCount, cfg.concurrency)
	result, err := load(ctx, cfg, loader)
	if err != nil {
		return err
	}

	_, err = os.Stdout.Write(result)
	if err != nil {
		return err
	}

	return nil
}

func closer(ctx context.Context) context.Context {
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)

	newCtx, cancel := context.WithCancel(ctx)

	go func() {
		select {
		case sig := <-sigChan:
			os.Stdout.WriteString(fmt.Sprintf("stop by %s \n", sig))
			cancel()
			fmt.Println("cancel")
		}
	}()

	return newCtx
}

func validateConfig(cfg config) error {
	if cfg.host == "" {
		return errors.New("empty host")
	}

	if cfg.requestsCount <= 0 {
		return fmt.Errorf("invalid requests count value - %d", cfg.requestsCount)
	}

	if cfg.timeOut <= 0 {
		return fmt.Errorf("invalid timeout value - %d", cfg.timeOut)
	}

	if _, ok := outputFormats[cfg.outputFormat]; !ok {
		return fmt.Errorf("invalid output format - %s", cfg.outputFormat)
	}

	if cfg.concurrency < 0 {
		return fmt.Errorf("invalid concurrency value - %d", cfg.concurrency)
	}

	return nil
}
