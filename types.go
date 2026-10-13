package buchida

// SendEmailParams are the parameters for sending an email.
type SendEmailParams struct {
	From        string            `json:"from"`
	FromName    string            `json:"from_name,omitempty"`
	To          interface{}       `json:"to"`
	Subject     string            `json:"subject"`
	HTML        string            `json:"html,omitempty"`
	Text        string            `json:"text,omitempty"`
	ReplyTo     string            `json:"reply_to,omitempty"`
	CC          interface{}       `json:"cc,omitempty"`
	BCC         interface{}       `json:"bcc,omitempty"`
	Tags        map[string]string `json:"tags,omitempty"`
	ScheduledAt string            `json:"scheduled_at,omitempty"`
}

// SendEmailResponse is the response from sending an email: the created email.
// (Before v0.2 this struct only declared ID, and it was never populated because the
// {"data": ...} envelope was not unwrapped.)
type SendEmailResponse = Email

// Email represents an email.
type Email struct {
	ID          string       `json:"id"`
	From        string       `json:"from"`
	FromName    *string      `json:"from_name,omitempty"`
	To          []string     `json:"to"`
	Subject     string       `json:"subject"`
	HTML        string       `json:"html,omitempty"`
	Text        string       `json:"text,omitempty"`
	ReplyTo     *string      `json:"reply_to,omitempty"`
	CC          []string     `json:"cc,omitempty"`
	BCC         []string     `json:"bcc,omitempty"`
	Status      string       `json:"status"`
	Direction   string       `json:"direction,omitempty"`
	TemplateID  *string      `json:"template_id,omitempty"`
	Tags        interface{}  `json:"tags,omitempty"`
	ScheduledAt *string      `json:"scheduled_at,omitempty"`
	SentAt      *string      `json:"sent_at,omitempty"`
	DeliveredAt *string      `json:"delivered_at,omitempty"`
	CreatedAt   string       `json:"created_at"`
	Events      []EmailEvent `json:"events,omitempty"`
}

// EmailEvent is an entry of an email's event timeline (returned by Emails.Get).
type EmailEvent struct {
	ID        string      `json:"id"`
	Type      string      `json:"type"`
	Timestamp string      `json:"timestamp"`
	Data      interface{} `json:"data,omitempty"`
}

// BatchError is the per-item error of a batch send.
type BatchError struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

// BatchResult is the result for one email of a batch send: either Data or Error is set.
type BatchResult struct {
	Data  *Email      `json:"data,omitempty"`
	Error *BatchError `json:"error,omitempty"`
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
	Data    []Email `json:"data"`
	Cursor  string  `json:"cursor,omitempty"`
	HasMore bool    `json:"hasMore"`
}

// CreateDomainParams are the parameters for creating a domain.
type CreateDomainParams struct {
	Name string `json:"name"`
}

// Domain represents a sending domain.
type Domain struct {
	ID         string      `json:"id"`
	Name       string      `json:"name"`
	Status     string      `json:"status"`
	Region     string      `json:"region,omitempty"`
	Records    []DNSRecord `json:"dns_records"`
	VerifiedAt *string     `json:"verified_at,omitempty"`
	CreatedAt  string      `json:"created_at"`
}

// DNSRecord represents a DNS record for domain verification.
type DNSRecord struct {
	Type        string `json:"type"`
	Name        string `json:"name"`
	Value       string `json:"value"`
	Purpose     string `json:"purpose,omitempty"`
	Label       string `json:"label,omitempty"`
	Description string `json:"description,omitempty"`
	Verified    bool   `json:"verified"`
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
	KeyPrefix  string `json:"key_prefix,omitempty"`
	Permission string `json:"permission"`
	CreatedAt  string `json:"created_at"`
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
	Secret    string   `json:"secret,omitempty"`
	Active    bool     `json:"active"`
	CreatedAt string   `json:"created_at"`
}

// Template represents an email template.
type Template struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	Subject   string `json:"subject,omitempty"`
	HTML      string `json:"html,omitempty"`
	CreatedAt string `json:"created_at"`
}

// GetMetricsParams are the query parameters for getting metrics.
type GetMetricsParams struct {
	From        string
	To          string
	Granularity string
}

// Metrics represents email sending metrics.
type Metrics struct {
	Sent       int                `json:"sent"`
	Delivered  int                `json:"delivered"`
	Opened     int                `json:"opened"`
	Clicked    int                `json:"clicked"`
	Bounced    int                `json:"bounced"`
	Complained int                `json:"complained"`
	Timeseries []MetricsDataPoint `json:"timeseries"`

	// Fields returned by the current API.
	Total         int            `json:"total"`
	ByStatus      map[string]int `json:"by_status,omitempty"`
	DeliveryRate  float64        `json:"delivery_rate"`
	BounceRate    float64        `json:"bounce_rate"`
	ComplaintRate float64        `json:"complaint_rate"`
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
