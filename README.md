# buchida-go

Official Go SDK for the [buchida](https://buchida.com) email API.

## Installation

```bash
go get github.com/Vyblor/buchida-go
```

## Quick Start

```go
package main

import (
    "context"
    "fmt"
    "log"

    buchida "github.com/Vyblor/buchida-go"
)

func main() {
    client, err := buchida.New("bc_live_xxxxxxxxxxxxxxxxxxxxx")
    if err != nil {
        log.Fatal(err)
    }

    resp, err := client.Emails.Send(context.Background(), &buchida.SendEmailParams{
        From:    "hello@yourdomain.com",
        To:      "user@example.com",
        Subject: "Welcome to buchida!",
        HTML:    "<h1>Hello!</h1><p>Welcome aboard.</p>",
    })
    if err != nil {
        log.Fatal(err)
    }

    fmt.Printf("Email sent: %s\n", resp.ID)
}
```

## Features

- Idiomatic Go with `context.Context` support
- Zero dependencies (stdlib `net/http`)
- Go 1.21+
- Typed error handling

## Error Handling

```go
resp, err := client.Emails.Send(ctx, params)
if err != nil {
    switch err.(type) {
    case *buchida.AuthenticationError:
        // 401 - invalid API key
    case *buchida.RateLimitError:
        // 429 - too many requests
    case *buchida.NotFoundError:
        // 404 - resource not found
    case *buchida.Error:
        // Other API errors
        apiErr := err.(*buchida.Error)
        fmt.Println(apiErr.StatusCode, apiErr.Message)
    }
}
```

## License

MIT
