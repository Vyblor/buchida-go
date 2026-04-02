<div align="center">
  <img src="assets/logo-black.svg" alt="buchida" width="280" />
  <p><strong>支持CJK的开发者优先邮件API</strong></p>

  [English](README.md) | [한국어](README.ko.md) | [日本語](README.ja.md) | [中文](README.zh.md)

  [![Go Reference](https://pkg.go.dev/badge/github.com/Vyblor/buchida-go)](https://pkg.go.dev/github.com/Vyblor/buchida-go) [![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](LICENSE)
</div>

---

[buchida](https://buchida.com)邮件API的官方Go SDK。

## 安装

```bash
go get github.com/Vyblor/buchida-go
```

## 快速开始

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
        Subject: "欢迎使用buchida！",
        HTML:    "<h1>你好！</h1><p>欢迎加入。</p>",
    })
    if err != nil {
        log.Fatal(err)
    }

    fmt.Printf("邮件发送成功: %s\n", resp.ID)
}
```

## 特性

- 支持`context.Context`的惯用Go代码
- 零依赖（标准库`net/http`）
- Go 1.21+
- 类型化错误处理

## 文档

- [快速开始](https://buchida.com/zh/docs/quickstart)
- [API参考](https://buchida.com/zh/docs/sending-email)
- [GitHub](https://github.com/Vyblor/buchida-go)

## 许可证

MIT
