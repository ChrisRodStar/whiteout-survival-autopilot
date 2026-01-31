// internal/ocrclient/client.go
package ocrclient

import (
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"github.com/avast/retry-go"
	"github.com/spf13/viper"
)

// RetryTransport — wrapper over RoundTripper with automatic retry.
type RetryTransport struct {
	Base     http.RoundTripper
	Attempts uint
	Delay    time.Duration
	Logger   *slog.Logger
}

func (t *RetryTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	var resp *http.Response
	err := retry.Do(
		func() error {
			r, err := t.Base.RoundTrip(req)
			if err != nil {
				t.Logger.Warn("HTTP error, retrying", "err", err)
				return err
			}
			if r.StatusCode >= 500 {
				// server error — close body and retry
				r.Body.Close()
				t.Logger.Warn("Server 5xx, retrying", "status", r.StatusCode)
				return fmt.Errorf("server status %d", r.StatusCode)
			}
			resp = r
			return nil
		},
		retry.Attempts(t.Attempts),
		retry.Delay(t.Delay),
		retry.LastErrorOnly(true),
	)
	return resp, err
}

// Client — HTTP client for communicating with OCR service.
type Client struct {
	ServiceURL string
	DeviceID   string
	Logger     *slog.Logger
	HTTP       *http.Client
}

// NewClient creates an OCR client with retry middleware.
// All requests through c.HTTP will be automatically
// retried 3 times with a 500ms delay.
func NewClient(deviceID string, logger *slog.Logger) *Client {
	viper.SetDefault("OCR_SERVICE_URL", "http://localhost:8000")

	transport := &RetryTransport{
		Base:     http.DefaultTransport,
		Attempts: 3,
		Delay:    500 * time.Millisecond,
		Logger:   logger,
	}
	httpClient := &http.Client{
		Transport: transport,
		Timeout:   40 * time.Second,
	}

	return &Client{
		ServiceURL: viper.GetString("OCR_SERVICE_URL"),
		DeviceID:   deviceID,
		Logger:     logger,
		HTTP:       httpClient,
	}
}
