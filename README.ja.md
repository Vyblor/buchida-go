<div align="center">
  <img src="assets/logo-black.svg" alt="buchida" width="280" />
  <p><strong>CJKサポートを備えた開発者向けメールAPI</strong></p>

  [English](README.md) | [한국어](README.ko.md) | [日本語](README.ja.md) | [中文](README.zh.md)

  [![Go Reference](https://pkg.go.dev/badge/github.com/Vyblor/buchida-go)](https://pkg.go.dev/github.com/Vyblor/buchida-go) [![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](LICENSE)
</div>

---

[buchida](https://buchida.com)メールAPIの公式Go SDKです。

## インストール

```bash
go get github.com/Vyblor/buchida-go
```

## クイックスタート

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
        Subject: "buchidaへようこそ！",
        HTML:    "<h1>こんにちは！</h1><p>ご登録ありがとうございます。</p>",
    })
    if err != nil {
        log.Fatal(err)
    }

    fmt.Printf("メール送信完了: %s\n", resp.ID)
}
```

## 特徴

- `context.Context`サポートによる慣用的なGoコード
- 依存関係ゼロ（標準ライブラリ`net/http`）
- Go 1.21+
- 型付きエラーハンドリング

## ドキュメント

- [クイックスタート](https://buchida.com/ja/docs/quickstart)
- [APIリファレンス](https://buchida.com/ja/docs/sending-email)
- [GitHub](https://github.com/Vyblor/buchida-go)

## ライセンス

MIT
