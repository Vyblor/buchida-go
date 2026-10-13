<div align="center">
  <img src="assets/logo-black.svg" alt="buchida" width="280" />
  <p><strong>buchida Go SDK — Email API for AI agents</strong></p>

  [**English**](README.md) | [한국어](README.ko.md) | [日本語](README.ja.md) | [中文](README.zh.md)

  [![Go Reference](https://pkg.go.dev/badge/github.com/Vyblor/buchida-go.svg)](https://pkg.go.dev/github.com/Vyblor/buchida-go) [![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](LICENSE)
</div>

---

buchida-go is the official Go SDK for **buchida** — an email API built for AI agents. buchida ships a CLI, an MCP server, and SDKs in 5 languages (Node, Python, Go, Ruby, Java), all sharing the same REST API surface. `@buchida/email` templates render Korean, Japanese, and Chinese natively.

## Install

```bash
go get github.com/Vyblor/buchida-go
```

## Send your first email

```go
package main

import (
	"context"
	"fmt"
	"log"
	"os"

	buchida "github.com/Vyblor/buchida-go"
)

func main() {
	client, err := buchida.New(os.Getenv("BUCHIDA_API_KEY"))
	if err != nil {
		log.Fatal(err)
	}

	email, err := client.Emails.Send(context.Background(), &buchida.SendEmailParams{
		From:     "hello@yourapp.com",
		FromName: "Your App",
		To:       []string{"user@example.com"},
		Subject:  "Hello",
		HTML:     "<h1>Welcome</h1>",
	})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(email.ID, email.Status)
}
```

Lists and batches:

```go
page, err := client.Emails.List(ctx, &buchida.ListEmailsParams{Limit: 20})
// page.Data []Email, page.Cursor string, page.HasMore bool

results, err := client.Emails.SendBatch(ctx, []buchida.SendEmailParams{
	{From: "hello@yourapp.com", To: "a@example.com", Subject: "A", Text: "hi"},
})
for _, r := range results {
	if r.Error != nil {
		fmt.Println(r.Error.Code, r.Error.Message)
	} else {
		fmt.Println(r.Data.ID)
	}
}
```

## Upgrading to 0.2

The API wraps every success response in `{"data": ...}`. Earlier versions decoded the
envelope straight into the result structs, so most fields (for example `Emails.Send` -> `ID`,
`Emails.Get`, `Domains.Get`) came back empty, and `Domains.List`, `ApiKeys.List`,
`Webhooks.List`, `Templates.List` and `Emails.SendBatch` returned a decode error. v0.2 unwraps
`data` and reads the API's snake_case fields (`created_at`, `from_name`, `reply_to`, `dns_records`, ...).

- `SendEmailResponse` is now an alias of `Email` (it only had `ID` before).
- `Emails.SendBatch` now returns `[]BatchResult` (`Data *Email` or `Error *BatchError` per item)
  instead of `[]SendEmailResponse`.
- `ListEmailsResponse` gained `HasMore`; `Domain.Records` now reads `dns_records`.
- New: `SendEmailParams.FromName`, `Email.FromName`; request bodies use snake_case
  (`reply_to`, `scheduled_at`); nested API errors (`{"error": {"code", "message"}}`) are parsed.

## Documentation

Full docs: **[buchida.com/docs](https://buchida.com/docs)**

- API reference: https://buchida.com/docs/api-reference
- Quickstart guide: https://buchida.com/docs/quickstart
- CJK email templates: https://buchida.com/docs/templates
- MCP server setup: https://buchida.com/docs/mcp
- CLI reference: https://buchida.com/docs/cli

## Links

- **Website:** [buchida.com](https://buchida.com)
- **Documentation:** [buchida.com/docs](https://buchida.com/docs)
- **Pricing:** [buchida.com/pricing](https://buchida.com/pricing)
- **GitHub:** https://github.com/Vyblor/buchida-go

## License

MIT
