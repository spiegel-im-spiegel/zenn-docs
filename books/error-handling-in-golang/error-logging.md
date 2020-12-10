---
title: "ぼくがかんがえたさいきょうのえらーろぐ"
---

というわけで，そろそろ「ぼくがかんがえたさいきょうのえらーろぐ」の出番だ（笑）

## [spiegel-im-spiegel/errs][errs]

[spiegel-im-spiegel/errs][errs] は自作のエラーハンドリンク・パッケージで，他で公開している自作のコマンドライン・ツールで主に使っているが，一応汎用で使えるよう構成している。主な特徴は以下の通り。

- [errors] 標準パッケージと置き換え可能（[errs].Is(), [errs].As() 等の関数が用意されている）
- [errs].WithContext() 関数を使ってコンテキスト情報を付加できる。付加した情報は map[string]interface{} 型の連想配列で保持される 
    - 既定でエラーが発生した関数名を格納している
- `%+v` 書式を使ってエラー情報の詳細を JSON 形式で出力できる。また MarshalJSON() メソッドを備えているので [encoding/json][json] 標準パッケージを使って JSON 形式にエンコーディングできる

簡単な使い方は以下の記事を参照のこと。

https://text.baldanders.info/release/errs-package-for-golang/

たとえば [spiegel-im-spiegel/errs][errs] を使ってこんな感じに書ける。

```go:sample5.go
package main

import (
    "fmt"
    "os"

    "github.com/spiegel-im-spiegel/errs"
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
{"Type":"*errs.Error","Err":{"Type":"*os.PathError","Msg":"open not-exist.txt: The system cannot find the file specified.","Cause":{"Type":"syscall.Errno","Msg":"The system cannot find the file specified."}},"Context":{"function":"main.checkFileOpen","path":"not-exist.txt"}}
```

てな感じにエラーが JSON 形式で表示される。
このままだと見にくいので [jq] コマンド等を使って

```
$ go run sample5.go | jq .
{
  "Type": "*errs.Error",
  "Err": {
    "Type": "*os.PathError",
    "Msg": "open not-exist.txt: The system cannot find the file specified.",
    "Cause": {
      "Type": "syscall.Errno",
      "Msg": "The system cannot find the file specified."
    }
  },
  "Context": {
    "function": "main.checkFileOpen",
    "path": "not-exist.txt"
  }
}
```

などとすれば分かりやすいかな。自作コマンドライン・ツールでは `--debug` オプションをつけると JSON 形式のエラーを吐くようにしている。

どや，かっこええやろ（CV は久川綾さんでw）

## [rs/zerolog][zerolog] を使って構造化ログを出力する

[rs/zerolog][zerolog] はパフォーマンスがよく，しかも JSON 形式でログを出力する優れものである。これを拙作の [spiegel-im-spiegel/errs][errs] と組み合わせることを考える。

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
        "Msg": "open not-exist.txt: The system cannot find the file specified.",
        "Cause": {
          "Type": "syscall.Errno",
          "Msg": "The system cannot find the file specified."
        }
      },
      "Context": {
        "function": "main.checkFileOpen",
        "path": "not-exist.txt"
      }
    },
    "time": "2020-12-10T14:39:48+09:00"
  }
]
```

という感じにエラー情報を JSON 形式で埋め込むことができる[^elog1]。

[^elog1]: 実際には [zerolog] パッケージと [errs] パッケージの結合では少し試行錯誤している。その辺の様子は「[構造化エラーをログ出力する](https://text.baldanders.info/golang/logging-error/)」をご覧あれ。

ちなみに，上のように [jq] に `-s` オプションをつけると複数の JSON オブジェクトを配列に組みなおして出力してくれる。 [encoding/json][json] 標準パッケージを使うなら [json].NewDecoder() 関数でデコーダを作ればオブジェクト単位でデコードしてくれるので， [for 文][for]で EOF まで回せばよい。

[rs/zerolog][zerolog] を使えばログを再利用しやすくなるので，是非とも活用していきたいところである。

[Go]: https://golang.org/ "The Go Programming Language"
[for]: https://golang.org/ref/spec#For_statements "The Go Programming Language Specification - The Go Programming Language"
[jq]: https://stedolan.github.io/jq/
[errs]: https://github.com/spiegel-im-spiegel/errs "spiegel-im-spiegel/errs: Error handling for Golang"
[errors]: https://golang.org/pkg/errors/ "errors - The Go Programming Language"
[json]: https://golang.org/pkg/encoding/json/ "json - The Go Programming Language"
[zerolog]: https://github.com/rs/zerolog "rs/zerolog: Zero Allocation JSON Logger"
<!-- eof -->
