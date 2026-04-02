<div align="center">
  <img src="assets/logo-black.svg" alt="buchida" width="280" />
  <p><strong>CJK 지원을 갖춘 개발자 중심 이메일 API</strong></p>

  [English](README.md) | [한국어](README.ko.md) | [日本語](README.ja.md) | [中文](README.zh.md)

  [![Go Reference](https://pkg.go.dev/badge/github.com/Vyblor/buchida-go)](https://pkg.go.dev/github.com/Vyblor/buchida-go) [![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](LICENSE)
</div>

---

[buchida](https://buchida.com) 이메일 API의 공식 Go SDK입니다.

## 설치

```bash
go get github.com/Vyblor/buchida-go
```

## 빠른 시작

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
        Subject: "buchida에 오신 것을 환영합니다!",
        HTML:    "<h1>안녕하세요!</h1><p>가입을 환영합니다.</p>",
    })
    if err != nil {
        log.Fatal(err)
    }

    fmt.Printf("이메일 발송 완료: %s\n", resp.ID)
}
```

## 주요 기능

- `context.Context` 지원의 관용적 Go 코드
- 의존성 없음 (표준 라이브러리 `net/http`)
- Go 1.21+
- 타입 기반 에러 처리

## 문서

- [빠른 시작 가이드](https://buchida.com/ko/docs/quickstart)
- [API 레퍼런스](https://buchida.com/ko/docs/sending-email)
- [GitHub](https://github.com/Vyblor/buchida-go)

## 라이선스

MIT
