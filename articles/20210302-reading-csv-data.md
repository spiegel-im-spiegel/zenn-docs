---
title: "CSV データを読み込むパッケージを書いてみた" # 記事のタイトル
emoji: "💻" # アイキャッチとして使われる絵文字（1文字だけ）
type: "tech" # "tech" : 技術記事 / "idea" : アイデア記事
topics: ["go", "programming"] # タグ。["markdown", "rust", "aws"] のように指定する
published: true # 公開設定（true で公開）
---

## [spiegel-im-spiegel/csvdata][csvdata] パッケージ

標準パッケージに [encoding/csv][csv] というのがあって [RFC 4180] に従って処理してくれるのだが， [encoding/csv][csv] 自体は基本的な機能しか用意されてないため，毎回ゴチャゴチャと周辺コード（とテスト）を書いていくのが面倒くさくなってきたんだよね。

ちうわけで [encoding/csv][csv] 標準パッケージに機能をちょい足しした [spiegel-im-spiegel/csvdata][csvdata] という小さいパッケージを書いてみた。

たとえば，こんな感じの CSV ファイルがあるとして

```markup:sample.csv
"order", name ,"mass","distance","habitable"
1, Mercury, 0.055, 0.4,false
2, Venus, 0.815, 0.7,false
3, Earth, 1.0, 1.0,true
4, Mars, 0.107, 1.5,false
```

以下のように読み込み処理を書く。

```go:sample.go
// +build run

package main

import (
    _ "embed"
    "errors"
    "fmt"
    "io"
    "os"
    "strings"

    "github.com/spiegel-im-spiegel/csvdata"
)

//go:embed sample.csv
var planets string

func main() {
    rc := csvdata.New(strings.NewReader(planets), true)
    for {
        if err := rc.Next(); err != nil {
            if errors.Is(err, io.EOF) {
                break
            }
            fmt.Fprintln(os.Stderr, err)
            return
        }
        order, err := rc.ColumnInt64("order", 10)
        if err != nil {
            fmt.Fprintln(os.Stderr, err)
            return
        }
        fmt.Println("    Order =", order)
        fmt.Println("     Name =", rc.Column("name"))
        mass, err := rc.ColumnFloat64("mass")
        if err != nil {
            fmt.Fprintln(os.Stderr, err)
            return
        }
        fmt.Println("     Mass =", mass)
        habitable, err := rc.ColumnBool("habitable")
        if err != nil {
            fmt.Fprintln(os.Stderr, err)
            return
        }
        fmt.Println("Habitable =", habitable)
    }
}
```

これを実行すると

```
$ go run sample.go
    Order = 1
     Name = Mercury
     Mass = 0.055
Habitable = false
    Order = 2
     Name = Venus
     Mass = 0.815
Habitable = false
    Order = 3
     Name = Earth
     Mass = 1
Habitable = true
    Order = 4
     Name = Mars
     Mass = 0.107
Habitable = false
```

てな感じに出力される。

ちなみに

```go
rt := csvdata.New(tsvReader, true).WithComma('\t')
```

とか WithComma() メソッドでセパレータを指定すれば TSV 等にも対応可能である。

[Go] 1.16 で登場した [embed] 標準パッケージと `//go:embed` ディレクティブは本当に素晴らしくて，これを使えばテストデータを用意するのが格段に楽になる。テスト準備データとして CSV や JSON ファイルを用意し，今回作ったようなパッケージでさくっと読んでテストに食わせるなんてケースがこれから増えるんじゃないかと夢想する。

とりあえず COVID-2019 関連の CSV データ読み込み処理を [spiegel-im-spiegel/csvdata][csvdata] パッケージで置き換えていくことにしよう。

## 【付録】 Shift-JIS エンコーディングの CSV データを読み込む

Excel 等でエクスポートした CSV ファイルの場合，文字エンコーディングが Shift-JIS だったりする場合がある。この場合は [golang.org/x/text/encoding/japanese](https://pkg.go.dev/golang.org/x/text/encoding/japanese) パッケージを使って UTF-8 エンコーディングに変換しつつ読み込むとよい。

つまり先程の sample.go のコードの [csvdata].New() 関数をこんな感じに書き換える。

```go
rc := csvdata.New(japanese.ShiftJIS.NewDecoder().Reader(os.Stdin), true)
```

こうすれば CSV データを必要なだけ読み込みつつ処理できる。

## 参考

https://zenn.dev/koya_iwamura/articles/53a4469271022e
https://text.baldanders.info/golang/embeded-filesystem/

[Go]: https://golang.org/ "The Go Programming Language"
[csv]: https://golang.org/pkg/encoding/csv/ "csv - The Go Programming Language"
[embed]: https://golang.org/pkg/embed/ "embed - The Go Programming Language"
[RFC 4180]: https://tools.ietf.org/html/rfc4180 "RFC 4180 - Common Format and MIME Type for Comma-Separated Values (CSV) Files"
[csvdata]: https://github.com/spiegel-im-spiegel/csvdata "spiegel-im-spiegel/csvdata: Reading CSV Data]"
<!-- eof -->
