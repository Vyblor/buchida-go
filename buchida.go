// Package buchida provides a Go client for the buchida email API.
package buchida

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

const (
	defaultBaseURL = "https://api.buchida.com/v1"
	defaultTimeout = 30 * time.Second
	userAgent      = "buchida-go/0.2.0"
)

// Option configures the Client.
type Option func(*Client)

// WithBaseURL sets a custom API base URL.
func WithBaseURL(u string) Option {
	return func(c *Client) { c.baseURL = u }
}

// WithTimeout sets the HTTP client timeout.
func WithTimeout(d time.Duration) Option {
	return func(c *Client) { c.httpClient.Timeout = d }
}

// WithHTTPClient sets a custom HTTP client.
func WithHTTPClient(hc *http.Client) Option {
	return func(c *Client) { c.httpClient = hc }
}

// Client is the buchida API client.
type Client struct {
	apiKey     string
	baseURL    string
	httpClient *http.Client

	Emails    *EmailsService
	Domains   *DomainsService
	ApiKeys   *ApiKeysService
	Webhooks  *WebhooksService
	Templates *TemplatesService
	Metrics   *MetricsService
}

// New creates a new buchida API client.
func New(apiKey string, opts ...Option) (*Client, error) {
	if apiKey == "" {
		return nil, fmt.Errorf("buchida: API key is required")
	}

	c := &Client{
		apiKey:  apiKey,
		baseURL: defaultBaseURL,
		httpClient: &http.Client{
			Timeout: defaultTimeout,
		},
	}

	for _, opt := range opts {
		opt(c)
	}

	c.Emails = &EmailsService{client: c}
	c.Domains = &DomainsService{client: c}
	c.ApiKeys = &ApiKeysService{client: c}
	c.Webhooks = &WebhooksService{client: c}
	c.Templates = &TemplatesService{client: c}
	c.Metrics = &MetricsService{client: c}

	return c, nil
}

func (c *Client) do(ctx context.Context, method, path string, body, result interface{}) error {
	u := c.baseURL + path

	var reqBody io.Reader
	if body != nil {
		data, err := json.Marshal(body)
		if err != nil {
			return fmt.Errorf("buchida: failed to marshal request: %w", err)
		}
		reqBody = bytes.NewReader(data)
	}

	req, err := http.NewRequestWithContext(ctx, method, u, reqBody)
	if err != nil {
		return fmt.Errorf("buchida: failed to create request: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+c.apiKey)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("User-Agent", userAgent)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("buchida: request failed: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("buchida: failed to read response: %w", err)
	}

	if resp.StatusCode >= 400 {
		return parseError(resp.StatusCode, respBody)
	}

	if result != nil && resp.StatusCode != 204 && len(respBody) > 0 {
		if err := json.Unmarshal(respBody, result); err != nil {
			return fmt.Errorf("buchida: failed to decode response: %w", err)
		}
	}

	return nil
}

// doData performs a request and decodes the "data" member of the {"data": ...} success
// envelope into result. A body without a "data" member is decoded as-is.
func (c *Client) doData(ctx context.Context, method, path string, body, result interface{}) error {
	var raw json.RawMessage
	if err := c.do(ctx, method, path, body, &raw); err != nil {
		return err
	}
	if len(raw) == 0 || result == nil {
		return nil
	}
	var env map[string]json.RawMessage
	if err := json.Unmarshal(raw, &env); err == nil {
		if data, ok := env["data"]; ok {
			raw = data
		}
	}
	if err := json.Unmarshal(raw, result); err != nil {
		return fmt.Errorf("buchida: failed to decode response: %w", err)
	}
	return nil
}

func parseError(status int, body []byte) error {
	var errBody struct {
		Message string `json:"message"`
		Code    string `json:"code"`
		Error   *struct {
			Message string `json:"message"`
			Code    string `json:"code"`
		} `json:"error"`
	}
	_ = json.Unmarshal(body, &errBody)
	if errBody.Error != nil {
		if errBody.Error.Message != "" {
			errBody.Message = errBody.Error.Message
		}
		if errBody.Error.Code != "" {
			errBody.Code = errBody.Error.Code
		}
	}
	if errBody.Message == "" {
		errBody.Message = http.StatusText(status)
	}

	switch status {
	case 401:
		return &AuthenticationError{StatusCode: status, Message: errBody.Message}
	case 404:
		return &NotFoundError{StatusCode: status, Message: errBody.Message}
	case 422:
		return &ValidationError{StatusCode: status, Message: errBody.Message}
	case 429:
		return &RateLimitError{StatusCode: status, Message: errBody.Message}
	default:
		return &APIError{StatusCode: status, Message: errBody.Message, Code: errBody.Code}
	}
}
