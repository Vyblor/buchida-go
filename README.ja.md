<div align="center">
  <img src="assets/logo-black.svg" alt="buchida" width="280" />
  <p><strong>buchida Go SDK — AI エージェントのためのメール API</strong></p>

  [English](README.md) | [한국어](README.ko.md) | [**日本語**](README.ja.md) | [中文](README.zh.md)

  [![Go Reference](https://pkg.go.dev/badge/github.com/Vyblor/buchida-go)](https://pkg.go.dev/github.com/Vyblor/buchida-go) [![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](LICENSE)
</div>

---

buchida-go は AI エージェントのために作られたメール API の公式 Go SDK です。buchida は CLI、MCP サーバー、そして 5 言語の SDK (Node、Python、Go、Ruby、Java) を提供しており、すべて同じ REST API 表面を共有しています。`@buchida/email` テンプレートは韓国語、日本語、中国語をネイティブにレンダリングします。

## インストール

```bash
go get github.com/Vyblor/buchida-go
```

## 最初のメールを送信

```go
import "github.com/Vyblor/buchida-go"

client := buchida.NewClient(os.Getenv("BUCHIDA_API_KEY"))

client.Emails.Send(&buchida.SendEmailRequest{
    From:    "hello@yourapp.com",
    To:      []string{"user@example.com"},
    Subject: "こんにちは",
    Html:    "<h1>ようこそ</h1>",
})
```

## ドキュメント

完全なドキュメント: **[buchida.com/docs](https://buchida.com/docs)**

- API リファレンス: https://buchida.com/docs/api-reference
- クイックスタートガイド: https://buchida.com/docs/quickstart
- CJK メールテンプレート: https://buchida.com/docs/templates
- MCP サーバーセットアップ: https://buchida.com/docs/mcp
- CLI リファレンス: https://buchida.com/docs/cli

## リンク

- **ウェブサイト:** [buchida.com](https://buchida.com)
- **ドキュメント:** [buchida.com/docs](https://buchida.com/docs)
- **料金:** [buchida.com/pricing](https://buchida.com/pricing)
- **GitHub:** https://github.com/Vyblor/buchida-go

## ライセンス

MIT
