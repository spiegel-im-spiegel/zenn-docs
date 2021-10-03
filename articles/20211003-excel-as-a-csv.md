---
title: "Excel も CSV みたいに扱いたい" # 記事のタイトル
emoji: "💻" # アイキャッチとして使われる絵文字（1文字だけ）
type: "tech" # "tech" : 技術記事 / "idea" : アイデア記事
topics: ["go", "programming"] # タグ。["markdown", "rust", "aws"] のように指定する
published: true # 公開設定（true で公開）
---

半年ほど前に

https://zenn.dev/spiegel/articles/20210302-reading-csv-data

という記事を書いたのだが，その後色々あって Excel や LibreOffice Calc のファイルも同じように扱いたいと思い，拙作の [github.com/spiegel-im-spiegel/csvdata][spiegel-im-spiegel/csvdata] パッケージを改造することにした。取るに足らないパッケージでも残しておくものである（笑）

https://github.com/spiegel-im-spiegel/csvdata/releases/tag/v0.3.0

[前の記事](https://zenn.dev/spiegel/articles/20210302-reading-csv-data "CSV データを読み込むパッケージを書いてみた")と比べてみると，以前は

```go
rc := csvdata.New(strings.NewReader(planets), true)
```

としていたのを

```go
rc := csvdata.NewRows(csvdata.New(strings.NewReader(planets)), true)
```

と New() 関数を2段階に分けた。

外側の [csvdata][spiegel-im-spiegel/csvdata].NewRows() 関数の引数を

```go:rows.go
//RowsReader is interface type for reading columns in a row.
type RowsReader interface {
    Read() ([]string, error)
    Close() error
}

func NewRows(rr RowsReader, headerFlag bool) *Rows { ... }
```

と interface 型にすることによって CSV 形式以外のファイルも受け入れるようにしようというわけ。いわゆる依存の注入（dependency injection）ですな。 [Go] では継承関係とか考えなくても interface 型を間に差し込むことで簡単に依存の注入を設計・実装できる。

まぁ，破壊的変更になるんだけど，私以外使ってる気配はないし，ええじゃろう（笑）

Excel ファイルの場合は

```go:exceldata/example_test.go
package exceldata_test

import (
    "fmt"

    "github.com/spiegel-im-spiegel/csvdata"
    "github.com/spiegel-im-spiegel/csvdata/exceldata"
)

func ExampleNew() {
    xlsx, err := exceldata.OpenFile("testdata/sample.xlsx", "")
    if err != nil {
        fmt.Println(err)
        return
    }
    r, err := exceldata.New(xlsx, "")
    if err != nil {
        fmt.Println(err)
        return
    }
    rc := csvdata.NewRows(r, true)
    defer rc.Close() //dummy

    if err := rc.Next(); err != nil {
        fmt.Println(err)
        return
    }
    fmt.Println(rc.Column("name"))
    // Output:
    // Mercury
}
```

exceldata.OpenFile() 関数の引数で Excel ファイルへのパスとパスワードを指定し（パスワードロックされていなければ空文字列でOK）， exceldata.New() 関数の引数で Excel データインスタンスとシート名を指定する（シート名が空文字列なら最初のシート）。あとは CSV と同じ手順でデータにアクセスできる。

LibreOffice Calc ファイルも同様に

```go:calcdata/example_test.go
package calcdata_test

import (
    "fmt"

    "github.com/spiegel-im-spiegel/csvdata"
    "github.com/spiegel-im-spiegel/csvdata/calcdata"
)

func ExampleNew() {
    ods, err := calcdata.OpenFile("testdata/sample.ods")
    if err != nil {
        fmt.Println(err)
        return
    }
    r, err := calcdata.New(ods, "")
    if err != nil {
        fmt.Println(err)
        return
    }
    rc := csvdata.NewRows(r, true)
    defer rc.Close() //dummy

    if err := rc.Next(); err != nil {
        fmt.Println(err)
        return
    }
    fmt.Println(rc.Column("name"))
    // Output:
    // Mercury
}
```

という感じ。 Calc ファイルの場合はパスワードロックやデータの暗号化には対応していない。ごめんペコン。

更に言うと Excel や Calc ファイルへのアクセスには以下の外部パッケージを利用しているが

https://github.com/qax-os/excelize
https://github.com/knieriem/odf

どちらも中身を全てヒープ上に展開してしまうようなので，数十万行とか大きなファイルは扱えないと思う。こちらもあしからずご了承の程を。

これで [CSV にいちいち変換](https://zenn.dev/spiegel/articles/20210516-excel-to-csv "Go で簡単 Excel → CSV 変換")しなくても直接扱えるようになったよ。めでたし

[Go]: https://golang.org/ "The Go Programming Language"
[spiegel-im-spiegel/csvdata]: https://github.com/spiegel-im-spiegel/csvdata "spiegel-im-spiegel/csvdata: Reading CSV Data"
