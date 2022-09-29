---
title: "Go で簡単なメール送信" # 記事のタイトル
emoji: "💻" # アイキャッチとして使われる絵文字（1文字だけ）
type: "tech" # "tech" : 技術記事 / "idea" : アイデア記事
topics: ["go", "programming"] # タグ。["markdown", "rust", "aws"] のように指定する
published: true # 公開設定（true で公開）
---

今回も小ネタで。

[Go] 標準の [net/smtp][smtp] パッケージを使って簡単なメール送信を組むことができる。本当に簡単なのでいきなりコードから。

```go:sample1.go
package main

import (
    "fmt"
    "net/smtp"
    "os"
    "strings"
)

var (
    hostname = "mail.example.com"
    port     = 587
    username = "user@example.com"
    password = "password"
)

func main() {
    from := "gopher@example.net"
    recipients := []string{"foo@example.com", "bar@example.com"}
    subject := "hello"
    body := "Hello World!\nHello Gopher!"

    auth := smtp.CRAMMD5Auth(username, password)
    msg := []byte(strings.ReplaceAll(fmt.Sprintf("To: %s\nSubject: %s\n\n%s", strings.Join(recipients, ","), subject, body), "\n", "\r\n"))
    if err := smtp.SendMail(fmt.Sprintf("%s:%d", hostname, port), auth, from, recipients, msg); err != nil {
        fmt.Fprintln(os.Stderr, err)
    }
}
```

ね。簡単でしょ。 SMTP 認証部分を平文認証にするのであれば smtp.CRAMMD5Auth() 関数を

```go
auth := smtp.PlainAuth("", username, password, hostname)
```

で置き換えればいい。ね。簡単でしょ。

ただし [net/smtp][smtp] パッケージは本当にメール送信（プロトコル）に特化しているため

- IRV (旧 US-ASCII) 以外の文字集合を含む場合は Content-Type フィールドを自前で追加し charset を指定する必要あり
- From, To, Cc , Bcc および Subject といったフィールドに IRV 以外の文字集合を使う場合は [RFC 2047] に従って符号化する必要あり
- マルチパートのメールを送る場合は Content-Type フィールドを自前で追加して boundary を指定し，さらに [mime/multipart][multipart] パッケージ等を使って本文の組み立てを行う必要あり（HTML メールやファイル等を添付して送る場合はマルチパートの制御が必要）

という感じに，ちょっと凝ったことをしようとすると途端に面倒くさくなる。

言い方を変えると IRV で単純なメッセージを送るなら [net/smtp][smtp] パッケージのみで全然 OK ということだ。たとえば [Go] でバッチ処理を組む際に上のようなコードを組み込んで，何らかの理由で正常終了しなかったらバッチ処理を閉じる前にちょろんとメールを送ったりできるわけだ。 Linux だと [msmtp](https://marlam.de/msmtp/) などのツールもあるが [Go] で認証情報ごとシングルバイナリに固めてしまうなら，取り扱いも簡単になるだろう。

ちなみに最初のコードは同じ [net/smtp][smtp] パッケージを使って以下のように書くこともできる。

```go:sample2.go
package main

import (
    "fmt"
    "io"
    "net/smtp"
    "os"
    "strings"
)

var (
    hostname = "mail.example.com"
    port     = 587
    username = "user@example.com"
    password = "password"
)

func main() {
    from := "gopher@example.net"
    recipients := []string{"foo@example.com", "bar@example.com"}
    subject := "hello"
    body := "Hello World!\nHello Gopher!"

    client, err := smtp.Dial(fmt.Sprintf("%s:%d", hostname, port))
    if err != nil {
        fmt.Fprintln(os.Stderr, err)
        return
    }
    defer client.Close()

    if err := client.Auth(smtp.CRAMMD5Auth(username, password)); err != nil {
        fmt.Fprintln(os.Stderr, err)
        return
    }

    if err := client.Mail(from); err != nil {
        fmt.Fprintln(os.Stderr, err)
        return
    }

    for _, addr := range recipients {
        if err := client.Rcpt(addr); err != nil {
            fmt.Fprintln(os.Stderr, err)
            return
        }
    }

    if err := func() error {
        w, err := client.Data()
        if err != nil {
            return err
        }
        defer w.Close()
        msg := strings.ReplaceAll(fmt.Sprintf("To: %s\nSubject: %s\n\n%s", strings.Join(recipients, ","), subject, body), "\n", "\r\n")
        if _, err := io.WriteString(w, msg); err != nil {
            return err
        }
        return nil
    }(); err != nil {
        fmt.Fprintln(os.Stderr, err)
        return
    }

    if err := client.Quit(); err != nil {
        fmt.Fprintln(os.Stderr, err)
    }
}
```

このコードではいきなり Quit() して一気に処理を終わらせているが，プロトコルをもっと細かく制御することもできるようだ。

# 参考

https://github.com/go-mail/mail
https://github.com/ungerik/go-mail
https://twinbird-htn.hatenablog.com/entry/2017/08/02/233000


[Go]: https://go.dev/ "The Go Programming Language"
[smtp]: https://pkg.go.dev/net/smtp "smtp package - net/smtp - Go Packages"
[multipart]: https://pkg.go.dev/mime/multipart "multipart package - mime/multipart - Go Packages"
[RFC 2047]: https://www.rfc-editor.org/rfc/rfc2047.html "RFC 2047: MIME (Multipurpose Internet Mail Extensions) Part Three: Message Header Extensions for Non-ASCII Text"
