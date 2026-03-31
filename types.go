package buchida

// SendEmailParams are the parameters for sending an email.
type SendEmailParams struct {
	From        string            `json:"from"`
	To          interface{}       `json:"to"`
	Subject     string            `json:"subject"`
	HTML        string            `json:"html,omitempty"`
	Text        string            `json:"text,omitempty"`
	ReplyTo     string            `json:"replyTo,omitempty"`
	CC          interface{}       `json:"cc,omitempty"`
	BCC         interface{}       `json:"bcc,omitempty"`
	Tags        map[string]string `json:"tags,omitempty"`
	ScheduledAt string            `json:"scheduledAt,omitempty"`
}

// SendEmailResponse is the response from sending an email.
type SendEmailResponse struct {
	ID string `json:"id"`
}

// Email represents a sent email.
type Email struct {
	ID        string   `json:"id"`
	From      string   `json:"from"`
	To        []string `json:"to"`
	Subject   string   `json:"subject"`
	HTML      string   `json:"html,omitempty"`
	Text      string   `json:"text,omitempty"`
	Status    string   `json:"status"`
	CreatedAt string   `json:"createdAt"`
}

// ListEmailsParams are the query parameters for listing emails.
type ListEmailsParams struct {
	Cursor string
	Limit  int
	Status string
	From   string
	To     string
}

// ListEmailsResponse is the paginated response for listing emails.
type ListEmailsResponse struct {
	Data   []Email `json:"data"`
	Cursor string  `json:"cursor,omitempty"`
}

// CreateDomainParams are the parameters for creating a domain.
type CreateDomainParams struct {
	Name string `json:"name"`
}

// Domain represents a sending domain.
type Domain struct {
	ID        string      `json:"id"`
	Name      string      `json:"name"`
	Status    string      `json:"status"`
	Records   []DNSRecord `json:"records"`
	CreatedAt string      `json:"createdAt"`
}

// DNSRecord represents a DNS record for domain verification.
type DNSRecord struct {
	Type     string `json:"type"`
	Name     string `json:"name"`
	Value    string `json:"value"`
	Verified bool   `json:"verified"`
}

// CreateApiKeyParams are the parameters for creating an API key.
type CreateApiKeyParams struct {
	Name       string `json:"name"`
	Permission string `json:"permission"`
}

// ApiKey represents an API key.
type ApiKey struct {
	ID         string `json:"id"`
	Name       string `json:"name"`
	Key        string `json:"key,omitempty"`
	Permission string `json:"permission"`
	CreatedAt  string `json:"createdAt"`
}

// CreateWebhookParams are the parameters for creating a webhook.
type CreateWebhookParams struct {
	URL    string   `json:"url"`
	Events []string `json:"events"`
}

// Webhook represents a webhook endpoint.
type Webhook struct {
	ID        string   `json:"id"`
	URL       string   `json:"url"`
	Events    []string `json:"events"`
	CreatedAt string   `json:"createdAt"`
}

// Template represents an email template.
type Template struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	Subject   string `json:"subject,omitempty"`
	HTML      string `json:"html,omitempty"`
	CreatedAt string `json:"createdAt"`
}

// GetMetricsParams are the query parameters for getting metrics.
type GetMetricsParams struct {
	From        string
	To          string
	Granularity string
}

// Metrics represents email sending metrics.
type Metrics struct {
	Sent        int                `json:"sent"`
	Delivered   int                `json:"delivered"`
	Opened      int                `json:"opened"`
	Clicked     int                `json:"clicked"`
	Bounced     int                `json:"bounced"`
	Complained  int                `json:"complained"`
	Timeseries  []MetricsDataPoint `json:"timeseries"`
}

// MetricsDataPoint is a single point in a metrics timeseries.
type MetricsDataPoint struct {
	Timestamp  string `json:"timestamp"`
	Sent       int    `json:"sent"`
	Delivered  int    `json:"delivered"`
	Opened     int    `json:"opened"`
	Clicked    int    `json:"clicked"`
	Bounced    int    `json:"bounced"`
	Complained int    `json:"complained"`
}
