package load

import (
	"benchutil/pkg/headers"
	"benchutil/pkg/httploader"
	"context"
	"fmt"
	"net/http"
	"os"
)

func load(ctx context.Context, cfg config, loader httploader.Loader) (output []byte, err error) {
	loadRep, err := makeLoad(ctx, cfg, loader)
	if err != nil {
		return nil, err
	}

	return readResult(loadRep, cfg.outputFormat)
}

func makeLoad(ctx context.Context, cfg config, loader httploader.Loader) (rep httploader.Report, err error) {
	var h *http.Header
	if cfg.headersPath != "" {
		if h, err = headers.ReadFromFile(cfg.headersPath); err != nil {
			return httploader.Report{}, fmt.Errorf("read headers:%w", err)
		}
	}

	var body []byte
	if cfg.bodyPath != "" {
		if body, err = readBody(cfg.bodyPath); err != nil {
			return httploader.Report{}, fmt.Errorf("read body: %w", err)
		}
	}

	rep, err = loader.Load(ctx, cfg.host, h, body)
	if err != nil {
		return httploader.Report{}, fmt.Errorf("load: %w", err)
	}

	return rep, nil
}

func readBody(path string) ([]byte, error) {
	body, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	return body, nil
}
