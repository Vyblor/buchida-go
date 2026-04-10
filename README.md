# buchida-go

**buchida Go SDK — Email API for AI agents**

buchida-go is the official Go SDK for **buchida** — an email API built for AI agents. buchida ships a CLI, an MCP server, and SDKs in 5 languages (Node, Python, Go, Ruby, Java), all sharing the same REST API surface. `@buchida/email` templates render Korean, Japanese, and Chinese natively.

## Install

```bash
go get github.com/Vyblor/buchida-go
```

## Send your first email

```go
import "github.com/Vyblor/buchida-go"

client := buchida.NewClient(os.Getenv("BUCHIDA_API_KEY"))

client.Emails.Send(&buchida.SendEmailRequest{
    From:    "hello@yourapp.com",
    To:      []string{"user@example.com"},
    Subject: "Hello",
    Html:    "<h1>Welcome</h1>",
})
```

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
