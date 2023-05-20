---
title: "ぼくがかんがえたさいきょうのえらーろぐ"
---

## 私の欲しいエラーと貴方の欲しいエラーは違う

エラーハンドリングで最も考慮すべきことは

💡 **利用者が欲しいエラー情報と提供者が欲しいエラー情報は異なる** 💡

という点だろう。

エラーが発生した際に利用者が最も欲しい情報は「どうすればいいのか？」である。そのためのヒントとして「何故エラーが起こったのか？」も欲しいわけだ。

たとえば，コマンドライン・ツールのフレームワークを提供する [spf13/cobra] は利用者がコマンド入力を間違えた際に正しいコマンドを推測して教えてくれる。[私が公開しているコマンドライン・ツール](https://github.com/goark/gpgpdump "goark/gpgpdump: OpenPGP packet visualizer")だとこんな感じ。

```
$ gpgpdump http
Error: unknown command "http" for "gpgpdump"

Did you mean this?
    hkp

Run 'gpgpdump --help' for usage.
```

これで本当に `gpgpdump hkp` と打ち間違えたのなら打ち直せばいいし，全然違うというのなら `gpgpdump --help` で使い方を表示してみればいい，と分かる。

一方，提供者側にとって（利用者にエラー情報を提供するにせよ）最も欲しい情報は「どうやって起きたか？」である。これを知るためには，何故（＝原因）も含めて，エラー発生時の「文脈」をできるだけかき集めることが重要である。

エラー発生時にスタック情報を欲しがるエンジニアが多いのは，この情報が「文脈」の一部となりうるからだ。でも，これは個人的な見解だが，スタック情報は9割以上がノイズである（実行時のプログラム構造解析がしたいなら別だが）。喩えるなら藁束の中から金の針を探すようなものだ。

じゃあ，エラーハンドリングはどういう戦術をとるのがいいのだろう。

...というわけで，そろそろ「ぼくがかんがえたさいきょうのえらーろぐ」の出番だ（笑）

## [goark/errs][errs]

[goark/errs][errs] は自作のエラーハンドリンク・パッケージで，他で公開している自作のコマンドライン・ツールの中で主に使っているが，一応は汎用で使えるよう構成している。主な特徴は以下の通り。

- [errors] 標準パッケージと置き換え可能（[errs].Is(), [errs].As() 等の関数が用意されている）
- [errs].WithContext() 関数を使って任意のコンテキスト情報を付加できる。付加した情報は map[string]interface{} 型の連想配列で保持される 
    - 既定でエラーが発生した関数名を格納している
- `%+v` 書式を使ってエラー情報の詳細を JSON 形式で出力できる。また MarshalJSON() メソッドを備えているので [encoding/json][json] 標準パッケージを使って JSON 形式にエンコーディングできる（デコード機能はない）

簡単な使い方は以下の記事を参照のこと。

https://text.baldanders.info/release/errs-package-for-golang/

たとえば，こんな感じに書ける。

```go:sample5.go
package main

import (
    "fmt"
    "os"

    "github.com/goark/errs"
)

func checkFileOpen(path string) error {
    file, err := os.Open(path)
    if err != nil {
        return errs.Wrap(
            err,
            errs.WithContext("path", path),
        )
    }
    defer file.Close()
    return nil
}

func main() {
    if err := checkFileOpen("not-exist.txt"); err != nil {
        fmt.Printf("%+v\n", err)
        return
    }
}
```

これを実行すると

```
$ go run sample5.go
{"Type":"*errs.Error","Err":{"Type":"*os.PathError","Msg":"open not-exist.txt: no such file or directory","Cause":{"Type":"syscall.Errno","Msg":"no such file or directory"}},"Context":{"function":"main.checkFileOpen","path":"not-exist.txt"}}
```

てな感じにエラーが JSON 形式で表示される。このままだと見にくいので [jq] コマンド等を使って

```
$ go run sample5.go | jq .
{
  "Type": "*errs.Error",
  "Err": {
    "Type": "*os.PathError",
    "Msg": "open not-exist.txt: no such file or directory",
    "Cause": {
      "Type": "syscall.Errno",
      "Msg": "no such file or directory"
    }
  },
  "Context": {
    "function": "main.checkFileOpen",
    "path": "not-exist.txt"
  }
}
```

などとすれば分かりやすいかな。自作コマンドライン・ツールでは `--debug` オプションをつけると JSON 形式のエラーを吐くようにしている。

## [rs/zerolog][zerolog] を使って構造化ログを出力する

[rs/zerolog][zerolog] ロギング・パッケージはパフォーマンスがよく，しかも JSON 形式でログを出力する優れものである。これを拙作の [goark/errs][errs] と組み合わせることを考える。

こんな感じでどうだろう。

```go:sample6.go
func main() {
    logger := zerolog.New(os.Stdout).Level(zerolog.DebugLevel).With().
        Timestamp().
        Str("role", "logger-sample").
        Logger()

    if err := checkFileOpen("not-exist.txt"); err != nil {
        logger.Error().Interface("error", err).Send()
        return
    }
}
```

これを実行すると

```
$ go run sample6.go | jq -s .
[
  {
    "level": "error",
    "role": "logger-sample",
    "error": {
      "Type": "*errs.Error",
      "Err": {
        "Type": "*os.PathError",
        "Msg": "open not-exist.txt: no such file or directory",
        "Cause": {
          "Type": "syscall.Errno",
          "Msg": "no such file or directory"
        }
      },
      "Context": {
        "function": "main.checkFileOpen",
        "path": "not-exist.txt"
      }
    },
    "time": "2020-12-10T19:16:29+09:00"
  }
]
```

という感じにエラー情報を JSON 形式で埋め込むことができる[^elog1]。

[^elog1]: 実際には [zerolog] パッケージと [errs] パッケージの結合では少し試行錯誤している。その辺の様子は「[構造化エラーをログ出力する](https://text.baldanders.info/golang/logging-error/)」をご覧あれ。

ちなみに，上のように [jq] に `-s` オプションをつけると複数の JSON オブジェクトを配列に組みなおして出力してくれる。 [encoding/json][json] 標準パッケージを使うなら [json].NewDecoder() 関数でデコーダを作ればオブジェクト単位でデコードしてくれるので， [for 文][for]で EOF まで回せばよい。

[rs/zerolog][zerolog] を使えばログを再利用しやすくなるので，是非とも活用していきたいところである。

## [zap] を使って構造化ログを出力する【2023-05-20 追記】

[Zap][zap] は gRPC 関連サービスや分散システムなどで人気の高い logger で，柔軟なカスタマイズができ，かつ高速で JSON 形式の構造化ログを出力できる。この logger に拙作のパッケージを食わせてみる。ソースコードはこんな感じ。

```go:sample7.go
package main

import (
    "os"

    "github.com/goark/errs"
    "go.uber.org/zap"
)

func checkFileOpen(path string) error {
    file, err := os.Open(path)
    if err != nil {
        return errs.New(
            "file open error",
            errs.WithCause(err),
            errs.WithContext("path", path),
        )
    }
    defer file.Close()

    return nil
}

func main() {
    logger := zap.NewExample()
    defer logger.Sync()

    path := "not-exist.txt"
    if err := checkFileOpen("not-exist.txt"); err != nil {
        logger.Error("error in checkFileOpen function", zap.Error(err), zap.String("file", path))
    }
}
```

これを実行すると

```
$ go run sample7.go | jq .
{
  "level": "error",
  "msg": "error in checkFileOpen function",
  "error": "file open error: open not-exist.txt: no such file or directory",
  "errorVerbose": "{\"Type\":\"*errs.Error\",\"Err\":{\"Type\":\"*errors.errorString\",\"Msg\":\"file open error\"},\"Context\":{\"function\":\"main.checkFileOpen\",\"path\":\"not-exist.txt\"},\"Cause\":{\"Type\":\"*fs.PathError\",\"Msg\":\"open not-exist.txt: no such file or directory\",\"Cause\":{\"Type\":\"syscall.Errno\",\"Msg\":\"no such file or directory\"}}}",
  "file": "not-exist.txt"
}
```

となる。
`"error"` 項目も `"errorVerbose"` 項目も文字列として出力されてしまうため構造化されているとは言えない。

[Zap][zap] には [zap].Object() 関数があって，これを使えば内部構造を出力することができるのだが，この関数を使うためには対象のオブジェクトが [zapcore].ObjectMarshaler 型の interface を満たす必要がある。

```go
type ObjectMarshaler interface {
    MarshalLogObject(ObjectEncoder) error
}
```

この要件を満たすために [goark/errs/zapobject][zapobject] モジュールを作った。こんな感じに error をラッピングして使う。

```go:sample8.go
package main

import (
    "os"

    "github.com/goark/errs"
    "github.com/goark/errs/zapobject"
    "go.uber.org/zap"
)

func checkFileOpen(path string) error {
    ...
}

func main() {
    logger := zap.NewExample()
    defer logger.Sync()

    path := "not-exist.txt"
    if err := checkFileOpen("not-exist.txt"); err != nil {
        logger.Error("error in checkFileOpen function", zap.Object("error", zapobject.New(err)), zap.String("file", path))
    }
}
```

これを実行すると

```
$ go run sample8.go | jq .
{
  "level": "error",
  "msg": "error in checkFileOpen function",
  "error": {
    "type": "*errs.Error",
    "msg": "file open error: open not-exist.txt: no such file or directory",
    "error": {
      "type": "*errors.errorString",
      "msg": "file open error"
    },
    "cause": {
      "type": "*fs.PathError",
      "msg": "open not-exist.txt: no such file or directory",
      "cause": {
        "type": "syscall.Errno",
        "msg": "no such file or directory"
      }
    },
    "context": {
      "function": "main.checkFileOpen",
      "path": "not-exist.txt"
    }
  },
  "file": "not-exist.txt"
}
```

と，いい感じに構造化されて出力される。

[Go]: https://go.dev/ "The Go Programming Language"
[for]: https://go.dev/ref/spec#For_statements "The Go Programming Language Specification - The Go Programming Language"
[spf13/cobra]: https://github.com/spf13/cobra "spf13/cobra: A Commander for modern Go CLI interactions"
[jq]: https://stedolan.github.io/jq/
[errs]: https://github.com/goark/errs "goark/errs: Error handling for Golang"
[errors]: https://pkg.go.dev/errors/ "errors - The Go Programming Language"
[json]: https://pkg.go.dev/encoding/json/ "json - The Go Programming Language"
[zerolog]: https://github.com/rs/zerolog "rs/zerolog: Zero Allocation JSON Logger"
[zapobject]: https://pkg.go.dev/github.com/goark/errs/zapobject "zapobject package - github.com/goark/errs/zapobject - Go Packages"
[zap]: https://pkg.go.dev/go.uber.org/zap "zap package - go.uber.org/zap - Go Packages"
[zapcore]: https://pkg.go.dev/go.uber.org/zap@v1.24.0/zapcore "zapcore package - go.uber.org/zap/zapcore - Go Packages"
