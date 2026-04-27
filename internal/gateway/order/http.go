package order

import (
	"context"
	"course-go-avito-SitnikovArtem06/internal/observability"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

const (
	TimeOut       = 5
	RetryAttempts = 5
	RetryDelay    = 100 * time.Millisecond
)

type HttpGateway struct {
	baseURL string
	client  *http.Client
}

func NewHttpGateway(url string, client *http.Client) *HttpGateway {
	return &HttpGateway{
		baseURL: url,
		client:  client,
	}
}

func (g *HttpGateway) GetOrder(ctx context.Context, orderId string) (*OrderDto, error) {

	url := g.baseURL + "/public/api/v1/order/" + orderId

	var lastErr error

	for i := 1; i <= RetryAttempts; i++ {
		if i != 1 {
			observability.GatewayRetriesTotal.WithLabelValues().Inc()
		}
		req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
		if err != nil {
			return nil, err
		}

		req.Header.Set("Accept", "application/json")

		resp, err := g.client.Do(req)
		if err != nil {
			return nil, err
		}

		switch resp.StatusCode {
		case http.StatusOK:
			var order OrderDto
			err = json.NewDecoder(resp.Body).Decode(&order)
			_ = resp.Body.Close()
			if err != nil {
				return nil, err
			}
			return &order, nil
		case http.StatusTooManyRequests,
			http.StatusBadGateway,
			http.StatusServiceUnavailable,
			http.StatusGatewayTimeout:
			_ = resp.Body.Close()
			lastErr = fmt.Errorf("order gateway: status=%d", resp.StatusCode)
		default:
			_ = resp.Body.Close()
			return nil, fmt.Errorf("order gateway: status=%d", resp.StatusCode)
		}
		if i < RetryAttempts {
			select {
			case <-ctx.Done():
				return nil, ctx.Err()
			case <-time.After(RetryDelay):
			}
		}

	}

	return nil, fmt.Errorf("order gateway: failed after %d attempts: %w", RetryAttempts, lastErr)

}
