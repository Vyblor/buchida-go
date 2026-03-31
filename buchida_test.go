package buchida

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
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
	client, server := setupTestServer(t, 200, map[string]string{"id": "email_123"})
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
}

func TestEmailsGet(t *testing.T) {
	client, server := setupTestServer(t, 200, map[string]interface{}{
		"id":        "email_123",
		"from":      "hi@buchida.com",
		"to":        []string{"user@example.com"},
		"subject":   "Hello",
		"status":    "delivered",
		"createdAt": "2026-03-31T00:00:00Z",
	})
	defer server.Close()

	email, err := client.Emails.Get(context.Background(), "email_123")
	if err != nil {
		t.Fatal(err)
	}
	if email.Status != "delivered" {
		t.Fatalf("expected delivered, got %s", email.Status)
	}
}

func TestEmailsList(t *testing.T) {
	client, server := setupTestServer(t, 200, map[string]interface{}{
		"data":   []interface{}{},
		"cursor": "",
	})
	defer server.Close()

	resp, err := client.Emails.List(context.Background(), &ListEmailsParams{Limit: 10, Status: "delivered"})
	if err != nil {
		t.Fatal(err)
	}
	if resp == nil {
		t.Fatal("expected response")
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
	client, server := setupTestServer(t, 200, []map[string]string{
		{"id": "email_1"},
		{"id": "email_2"},
	})
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
}

func TestDomainsCreate(t *testing.T) {
	client, server := setupTestServer(t, 200, map[string]interface{}{
		"id":        "dom_1",
		"name":      "example.com",
		"status":    "pending",
		"records":   []interface{}{},
		"createdAt": "2026-03-31T00:00:00Z",
	})
	defer server.Close()

	domain, err := client.Domains.Create(context.Background(), &CreateDomainParams{Name: "example.com"})
	if err != nil {
		t.Fatal(err)
	}
	if domain.Name != "example.com" {
		t.Fatalf("expected example.com, got %s", domain.Name)
	}
}

func TestDomainsList(t *testing.T) {
	client, server := setupTestServer(t, 200, []interface{}{})
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
	client, server := setupTestServer(t, 200, map[string]interface{}{
		"id":        "dom_1",
		"name":      "example.com",
		"status":    "verified",
		"records":   []interface{}{},
		"createdAt": "2026-03-31T00:00:00Z",
	})
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
	client, server := setupTestServer(t, 200, map[string]interface{}{
		"id":         "key_1",
		"name":       "test",
		"key":        "bc_live_newkey",
		"permission": "full_access",
		"createdAt":  "2026-03-31T00:00:00Z",
	})
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
	client, server := setupTestServer(t, 200, map[string]interface{}{
		"id":        "wh_1",
		"url":       "https://example.com/wh",
		"events":    []string{"email.delivered"},
		"createdAt": "2026-03-31T00:00:00Z",
	})
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
	client, server := setupTestServer(t, 200, []interface{}{})
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
	client, server := setupTestServer(t, 200, map[string]interface{}{
		"id":        "tpl_1",
		"name":      "Welcome",
		"createdAt": "2026-03-31T00:00:00Z",
	})
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
	client, server := setupTestServer(t, 200, map[string]interface{}{
		"sent":       100,
		"delivered":  95,
		"opened":     50,
		"clicked":    10,
		"bounced":    3,
		"complained": 1,
		"timeseries": []interface{}{},
	})
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

	if _, ok := err.(*Error); !ok {
		t.Fatalf("expected Error, got %T", err)
	}
}
