package buchida

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func setupTestServer(t *testing.T, status int, body interface{}) (*Client, *httptest.Server) {
	t.Helper()
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Verify auth header
		auth := r.Header.Get("Authorization")
		if auth != "Bearer bc_test_xxx" {
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(map[string]string{"message": "Invalid API key"})
			return
		}

		if status == 204 {
			w.WriteHeader(204)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(status)
		if body != nil {
			json.NewEncoder(w).Encode(body)
		}
	}))

	client, err := New("bc_test_xxx", WithBaseURL(server.URL))
	if err != nil {
		t.Fatal(err)
	}

	return client, server
}

// envelope wraps v in the API's {"data": ...} success envelope.
func envelope(v interface{}) map[string]interface{} {
	return map[string]interface{}{"data": v}
}

func TestNew(t *testing.T) {
	_, err := New("")
	if err == nil {
		t.Fatal("expected error for empty API key")
	}

	c, err := New("bc_test_xxx")
	if err != nil {
		t.Fatal(err)
	}
	if c.Emails == nil || c.Domains == nil || c.ApiKeys == nil {
		t.Fatal("services not initialized")
	}
}

func TestEmailsSend(t *testing.T) {
	client, server := setupTestServer(t, 200, envelope(map[string]string{"id": "email_123", "status": "queued"}))
	defer server.Close()

	resp, err := client.Emails.Send(context.Background(), &SendEmailParams{
		From:    "hi@buchida.com",
		To:      "user@example.com",
		Subject: "Hello",
		HTML:    "<p>Hi</p>",
	})
	if err != nil {
		t.Fatal(err)
	}
	if resp.ID != "email_123" {
		t.Fatalf("expected id email_123, got %s", resp.ID)
	}
	if resp.Status != "queued" {
		t.Fatalf("expected status queued, got %s", resp.Status)
	}
}

func TestEmailsGet(t *testing.T) {
	client, server := setupTestServer(t, 200, envelope(map[string]interface{}{
		"id":         "email_123",
		"from":       "hi@buchida.com",
		"to":         []string{"user@example.com"},
		"subject":    "Hello",
		"status":     "delivered",
		"from_name":  "Vyblor",
		"reply_to":   "support@buchida.com",
		"sent_at":    "2026-03-31T00:00:05Z",
		"created_at": "2026-03-31T00:00:00Z",
	}))
	defer server.Close()

	email, err := client.Emails.Get(context.Background(), "email_123")
	if err != nil {
		t.Fatal(err)
	}
	if email.Status != "delivered" {
		t.Fatalf("expected delivered, got %s", email.Status)
	}
	if email.FromName == nil || *email.FromName != "Vyblor" {
		t.Fatalf("expected from_name Vyblor, got %v", email.FromName)
	}
	if email.CreatedAt != "2026-03-31T00:00:00Z" {
		t.Fatalf("created_at not parsed: %q", email.CreatedAt)
	}
	if email.ReplyTo == nil || *email.ReplyTo != "support@buchida.com" {
		t.Fatalf("reply_to not parsed: %v", email.ReplyTo)
	}
	if email.SentAt == nil || *email.SentAt != "2026-03-31T00:00:05Z" {
		t.Fatalf("sent_at not parsed: %v", email.SentAt)
	}
}

func TestEmailsList(t *testing.T) {
	client, server := setupTestServer(t, 200, map[string]interface{}{
		"data":    []interface{}{map[string]interface{}{"id": "e1", "created_at": "2026-03-31T00:00:00Z"}},
		"cursor":  "e1",
		"hasMore": true,
	})
	defer server.Close()

	resp, err := client.Emails.List(context.Background(), &ListEmailsParams{Limit: 10, Status: "delivered"})
	if err != nil {
		t.Fatal(err)
	}
	if len(resp.Data) != 1 || resp.Data[0].ID != "e1" || resp.Data[0].CreatedAt == "" {
		t.Fatalf("unexpected list data: %+v", resp.Data)
	}
	if resp.Cursor != "e1" || !resp.HasMore {
		t.Fatalf("expected cursor e1 / hasMore true, got %q / %v", resp.Cursor, resp.HasMore)
	}
}

func TestEmailsCancel(t *testing.T) {
	client, server := setupTestServer(t, 204, nil)
	defer server.Close()

	err := client.Emails.Cancel(context.Background(), "email_123")
	if err != nil {
		t.Fatal(err)
	}
}

func TestEmailsSendBatch(t *testing.T) {
	client, server := setupTestServer(t, 200, envelope([]interface{}{
		map[string]interface{}{"data": map[string]interface{}{"id": "email_1", "status": "queued"}},
		map[string]interface{}{"error": map[string]string{"code": "invalid_email_address", "message": "bad to"}},
	}))
	defer server.Close()

	resp, err := client.Emails.SendBatch(context.Background(), []SendEmailParams{
		{From: "hi@buchida.com", To: "a@example.com", Subject: "A"},
		{From: "hi@buchida.com", To: "b@example.com", Subject: "B"},
	})
	if err != nil {
		t.Fatal(err)
	}
	if len(resp) != 2 {
		t.Fatalf("expected 2 responses, got %d", len(resp))
	}
	if resp[0].Data == nil || resp[0].Data.ID != "email_1" || resp[0].Error != nil {
		t.Fatalf("unexpected first result: %+v", resp[0])
	}
	if resp[1].Data != nil || resp[1].Error == nil || resp[1].Error.Code != "invalid_email_address" {
		t.Fatalf("unexpected second result: %+v", resp[1])
	}
}

func TestDomainsCreate(t *testing.T) {
	client, server := setupTestServer(t, 200, envelope(map[string]interface{}{
		"id":          "dom_1",
		"name":        "example.com",
		"status":      "pending",
		"dns_records": []interface{}{map[string]string{"type": "TXT", "purpose": "spf"}},
		"created_at":  "2026-03-31T00:00:00Z",
	}))
	defer server.Close()

	domain, err := client.Domains.Create(context.Background(), &CreateDomainParams{Name: "example.com"})
	if err != nil {
		t.Fatal(err)
	}
	if domain.Name != "example.com" {
		t.Fatalf("expected example.com, got %s", domain.Name)
	}
	if domain.CreatedAt != "2026-03-31T00:00:00Z" {
		t.Fatalf("created_at not parsed: %q", domain.CreatedAt)
	}
	if len(domain.Records) != 1 || domain.Records[0].Purpose != "spf" {
		t.Fatalf("dns_records not parsed: %+v", domain.Records)
	}
}

func TestDomainsList(t *testing.T) {
	client, server := setupTestServer(t, 200, envelope([]interface{}{}))
	defer server.Close()

	domains, err := client.Domains.List(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if domains == nil {
		t.Fatal("expected non-nil")
	}
}

func TestDomainsVerify(t *testing.T) {
	client, server := setupTestServer(t, 200, envelope(map[string]interface{}{
		"id":          "dom_1",
		"name":        "example.com",
		"status":      "verified",
		"dns_records": []interface{}{map[string]string{"type": "TXT", "purpose": "spf"}},
		"created_at":  "2026-03-31T00:00:00Z",
	}))
	defer server.Close()

	domain, err := client.Domains.Verify(context.Background(), "dom_1")
	if err != nil {
		t.Fatal(err)
	}
	if domain.Status != "verified" {
		t.Fatalf("expected verified, got %s", domain.Status)
	}
}

func TestApiKeysCreate(t *testing.T) {
	client, server := setupTestServer(t, 200, envelope(map[string]interface{}{
		"id":         "key_1",
		"name":       "test",
		"key":        "bc_live_newkey",
		"permission": "full_access",
		"created_at": "2026-03-31T00:00:00Z",
	}))
	defer server.Close()

	key, err := client.ApiKeys.Create(context.Background(), &CreateApiKeyParams{
		Name:       "test",
		Permission: "full_access",
	})
	if err != nil {
		t.Fatal(err)
	}
	if key.Key != "bc_live_newkey" {
		t.Fatalf("expected bc_live_newkey, got %s", key.Key)
	}
}

func TestApiKeysDelete(t *testing.T) {
	client, server := setupTestServer(t, 204, nil)
	defer server.Close()

	err := client.ApiKeys.Delete(context.Background(), "key_1")
	if err != nil {
		t.Fatal(err)
	}
}

func TestWebhooksCreate(t *testing.T) {
	client, server := setupTestServer(t, 200, envelope(map[string]interface{}{
		"id":         "wh_1",
		"url":        "https://example.com/wh",
		"events":     []string{"email.delivered"},
		"created_at": "2026-03-31T00:00:00Z",
	}))
	defer server.Close()

	wh, err := client.Webhooks.Create(context.Background(), &CreateWebhookParams{
		URL:    "https://example.com/wh",
		Events: []string{"email.delivered"},
	})
	if err != nil {
		t.Fatal(err)
	}
	if wh.ID != "wh_1" {
		t.Fatalf("expected wh_1, got %s", wh.ID)
	}
}

func TestTemplatesList(t *testing.T) {
	client, server := setupTestServer(t, 200, envelope([]interface{}{}))
	defer server.Close()

	templates, err := client.Templates.List(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if templates == nil {
		t.Fatal("expected non-nil")
	}
}

func TestTemplatesGet(t *testing.T) {
	client, server := setupTestServer(t, 200, envelope(map[string]interface{}{
		"id":         "tpl_1",
		"name":       "Welcome",
		"created_at": "2026-03-31T00:00:00Z",
	}))
	defer server.Close()

	tpl, err := client.Templates.Get(context.Background(), "tpl_1")
	if err != nil {
		t.Fatal(err)
	}
	if tpl.Name != "Welcome" {
		t.Fatalf("expected Welcome, got %s", tpl.Name)
	}
}

func TestMetricsGet(t *testing.T) {
	client, server := setupTestServer(t, 200, envelope(map[string]interface{}{
		"sent":       100,
		"delivered":  95,
		"opened":     50,
		"clicked":    10,
		"bounced":    3,
		"complained": 1,
		"timeseries": []interface{}{},
	}))
	defer server.Close()

	metrics, err := client.Metrics.Get(context.Background(), &GetMetricsParams{
		From:        "2026-03-01",
		To:          "2026-03-31",
		Granularity: "day",
	})
	if err != nil {
		t.Fatal(err)
	}
	if metrics.Sent != 100 {
		t.Fatalf("expected 100, got %d", metrics.Sent)
	}
}

func TestError401(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(401)
		json.NewEncoder(w).Encode(map[string]string{"message": "Invalid API key"})
	}))
	defer server.Close()

	client, _ := New("bc_test_xxx", WithBaseURL(server.URL))
	_, err := client.Emails.List(context.Background(), nil)

	if err == nil {
		t.Fatal("expected error")
	}

	if _, ok := err.(*AuthenticationError); !ok {
		t.Fatalf("expected AuthenticationError, got %T", err)
	}
}

func TestError429(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(429)
		json.NewEncoder(w).Encode(map[string]string{"message": "Rate limit exceeded"})
	}))
	defer server.Close()

	client, _ := New("bc_test_xxx", WithBaseURL(server.URL))
	_, err := client.Emails.List(context.Background(), nil)

	if _, ok := err.(*RateLimitError); !ok {
		t.Fatalf("expected RateLimitError, got %T", err)
	}
}

func TestError500(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(500)
		json.NewEncoder(w).Encode(map[string]string{"message": "Internal server error"})
	}))
	defer server.Close()

	client, _ := New("bc_test_xxx", WithBaseURL(server.URL))
	_, err := client.Emails.List(context.Background(), nil)

	if _, ok := err.(*APIError); !ok {
		t.Fatalf("expected APIError, got %T", err)
	}
}

func TestEmailsSendWireKeys(t *testing.T) {
	var got map[string]interface{}
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		json.NewDecoder(r.Body).Decode(&got)
		w.WriteHeader(200)
		json.NewEncoder(w).Encode(map[string]string{"id": "email_123"})
	}))
	defer server.Close()

	client, _ := New("bc_test_xxx", WithBaseURL(server.URL))
	_, err := client.Emails.Send(context.Background(), &SendEmailParams{
		From:        "hi@buchida.com",
		FromName:    "Vyblor",
		To:          "user@example.com",
		Subject:     "Hello",
		ReplyTo:     "support@buchida.com",
		ScheduledAt: "2026-12-01T00:00:00Z",
	})
	if err != nil {
		t.Fatal(err)
	}
	if got["from_name"] != "Vyblor" || got["reply_to"] != "support@buchida.com" ||
		got["scheduled_at"] != "2026-12-01T00:00:00Z" {
		t.Fatalf("unexpected wire keys: %v", got)
	}
	if _, ok := got["replyTo"]; ok {
		t.Fatal("camelCase replyTo must not be sent")
	}
}

func TestErrorNestedBody(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(422)
		w.Write([]byte(`{"error":{"code":"validation_error","message":"from_name is invalid"}}`))
	}))
	defer server.Close()

	client, _ := New("bc_test_xxx", WithBaseURL(server.URL))
	_, err := client.Emails.List(context.Background(), nil)
	if _, ok := err.(*ValidationError); !ok {
		t.Fatalf("expected ValidationError, got %T", err)
	}
	if !strings.Contains(err.Error(), "from_name is invalid") {
		t.Fatalf("nested message not surfaced: %v", err)
	}
}

func TestListsUnwrapEnvelope(t *testing.T) {
	client, server := setupTestServer(t, 200, envelope([]interface{}{
		map[string]interface{}{"id": "dom_1", "name": "a.com", "created_at": "2026-03-31T00:00:00Z"},
		map[string]interface{}{"id": "dom_2", "name": "b.com", "created_at": "2026-04-01T00:00:00Z"},
	}))
	defer server.Close()

	domains, err := client.Domains.List(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if len(domains) != 2 || domains[1].Name != "b.com" || domains[1].CreatedAt != "2026-04-01T00:00:00Z" {
		t.Fatalf("unexpected domains: %+v", domains)
	}
}

func TestSnakeCaseModelFields(t *testing.T) {
	client, server := setupTestServer(t, 200, envelope(map[string]interface{}{
		"id": "key_1", "name": "k", "key_prefix": "bc_live_ab", "permission": "full_access",
		"created_at": "2026-03-31T00:00:00Z",
	}))
	defer server.Close()

	key, err := client.ApiKeys.Create(context.Background(), &CreateApiKeyParams{Name: "k", Permission: "full_access"})
	if err != nil {
		t.Fatal(err)
	}
	if key.ID != "key_1" || key.KeyPrefix != "bc_live_ab" || key.CreatedAt != "2026-03-31T00:00:00Z" {
		t.Fatalf("unexpected key: %+v", key)
	}
}

func TestMetricsCurrentShape(t *testing.T) {
	client, server := setupTestServer(t, 200, envelope(map[string]interface{}{
		"total": 10, "by_status": map[string]int{"delivered": 9}, "delivery_rate": 0.9,
	}))
	defer server.Close()

	m, err := client.Metrics.Get(context.Background(), &GetMetricsParams{From: "2026-03-01", To: "2026-03-31"})
	if err != nil {
		t.Fatal(err)
	}
	if m.Total != 10 || m.ByStatus["delivered"] != 9 || m.DeliveryRate != 0.9 {
		t.Fatalf("unexpected metrics: %+v", m)
	}
}

// The API is mounted under /v1; without it every call returned 404.
func TestDefaultBaseURLIncludesV1(t *testing.T) {
	c, err := New("bc_test_xxx")
	if err != nil {
		t.Fatal(err)
	}
	if c.baseURL != "https://api.buchida.com/v1" {
		t.Fatalf("unexpected default base URL %q", c.baseURL)
	}
}
