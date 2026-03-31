package buchida

import (
	"context"
	"fmt"
	"net/url"
	"strconv"
)

// EmailsService handles email-related API calls.
type EmailsService struct{ client *Client }

// Send sends an email.
func (s *EmailsService) Send(ctx context.Context, params *SendEmailParams) (*SendEmailResponse, error) {
	var resp SendEmailResponse
	err := s.client.do(ctx, "POST", "/emails", params, &resp)
	return &resp, err
}

// Get retrieves an email by ID.
func (s *EmailsService) Get(ctx context.Context, id string) (*Email, error) {
	var resp Email
	err := s.client.do(ctx, "GET", fmt.Sprintf("/emails/%s", id), nil, &resp)
	return &resp, err
}

// List lists emails with optional filters.
func (s *EmailsService) List(ctx context.Context, params *ListEmailsParams) (*ListEmailsResponse, error) {
	q := url.Values{}
	if params != nil {
		if params.Cursor != "" {
			q.Set("cursor", params.Cursor)
		}
		if params.Limit > 0 {
			q.Set("limit", strconv.Itoa(params.Limit))
		}
		if params.Status != "" {
			q.Set("status", params.Status)
		}
		if params.From != "" {
			q.Set("from", params.From)
		}
		if params.To != "" {
			q.Set("to", params.To)
		}
	}
	path := "/emails"
	if qs := q.Encode(); qs != "" {
		path += "?" + qs
	}

	var resp ListEmailsResponse
	err := s.client.do(ctx, "GET", path, nil, &resp)
	return &resp, err
}

// Cancel cancels a scheduled email.
func (s *EmailsService) Cancel(ctx context.Context, id string) error {
	return s.client.do(ctx, "POST", fmt.Sprintf("/emails/%s/cancel", id), nil, nil)
}

// SendBatch sends multiple emails at once.
func (s *EmailsService) SendBatch(ctx context.Context, emails []SendEmailParams) ([]SendEmailResponse, error) {
	var resp []SendEmailResponse
	err := s.client.do(ctx, "POST", "/emails/batch", emails, &resp)
	return resp, err
}

// DomainsService handles domain-related API calls.
type DomainsService struct{ client *Client }

// Create creates a new sending domain.
func (s *DomainsService) Create(ctx context.Context, params *CreateDomainParams) (*Domain, error) {
	var resp Domain
	err := s.client.do(ctx, "POST", "/domains", params, &resp)
	return &resp, err
}

// List lists all domains.
func (s *DomainsService) List(ctx context.Context) ([]Domain, error) {
	var resp []Domain
	err := s.client.do(ctx, "GET", "/domains", nil, &resp)
	return resp, err
}

// Get retrieves a domain by ID.
func (s *DomainsService) Get(ctx context.Context, id string) (*Domain, error) {
	var resp Domain
	err := s.client.do(ctx, "GET", fmt.Sprintf("/domains/%s", id), nil, &resp)
	return &resp, err
}

// Verify triggers domain verification.
func (s *DomainsService) Verify(ctx context.Context, id string) (*Domain, error) {
	var resp Domain
	err := s.client.do(ctx, "POST", fmt.Sprintf("/domains/%s/verify", id), nil, &resp)
	return &resp, err
}

// ApiKeysService handles API key operations.
type ApiKeysService struct{ client *Client }

// Create creates a new API key.
func (s *ApiKeysService) Create(ctx context.Context, params *CreateApiKeyParams) (*ApiKey, error) {
	var resp ApiKey
	err := s.client.do(ctx, "POST", "/api-keys", params, &resp)
	return &resp, err
}

// List lists all API keys.
func (s *ApiKeysService) List(ctx context.Context) ([]ApiKey, error) {
	var resp []ApiKey
	err := s.client.do(ctx, "GET", "/api-keys", nil, &resp)
	return resp, err
}

// Delete deletes an API key.
func (s *ApiKeysService) Delete(ctx context.Context, id string) error {
	return s.client.do(ctx, "DELETE", fmt.Sprintf("/api-keys/%s", id), nil, nil)
}

// WebhooksService handles webhook operations.
type WebhooksService struct{ client *Client }

// Create creates a new webhook.
func (s *WebhooksService) Create(ctx context.Context, params *CreateWebhookParams) (*Webhook, error) {
	var resp Webhook
	err := s.client.do(ctx, "POST", "/webhooks", params, &resp)
	return &resp, err
}

// List lists all webhooks.
func (s *WebhooksService) List(ctx context.Context) ([]Webhook, error) {
	var resp []Webhook
	err := s.client.do(ctx, "GET", "/webhooks", nil, &resp)
	return resp, err
}

// Delete deletes a webhook.
func (s *WebhooksService) Delete(ctx context.Context, id string) error {
	return s.client.do(ctx, "DELETE", fmt.Sprintf("/webhooks/%s", id), nil, nil)
}

// TemplatesService handles template operations.
type TemplatesService struct{ client *Client }

// List lists all templates.
func (s *TemplatesService) List(ctx context.Context) ([]Template, error) {
	var resp []Template
	err := s.client.do(ctx, "GET", "/templates", nil, &resp)
	return resp, err
}

// Get retrieves a template by ID.
func (s *TemplatesService) Get(ctx context.Context, id string) (*Template, error) {
	var resp Template
	err := s.client.do(ctx, "GET", fmt.Sprintf("/templates/%s", id), nil, &resp)
	return &resp, err
}

// MetricsService handles metrics operations.
type MetricsService struct{ client *Client }

// Get retrieves email metrics.
func (s *MetricsService) Get(ctx context.Context, params *GetMetricsParams) (*Metrics, error) {
	q := url.Values{}
	q.Set("from", params.From)
	q.Set("to", params.To)
	if params.Granularity != "" {
		q.Set("granularity", params.Granularity)
	}

	var resp Metrics
	err := s.client.do(ctx, "GET", "/metrics?"+q.Encode(), nil, &resp)
	return &resp, err
}
