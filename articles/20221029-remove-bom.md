---
title: "Decorator Pattern で BOM を除去する" # 記事のタイトル
emoji: "💻" # アイキャッチとして使われる絵文字（1文字だけ）
type: "tech" # "tech" : 技術記事 / "idea" : アイデア記事
topics: ["go", "programming", "unicode"] # タグ。["markdown", "rust", "aws"] のように指定する
published: true # 公開設定（true で公開）
---

最近 CSV ファイルを扱う必要がありまして。 Windows では有名な [A5:SQL Mk-2](https://a5m2.mmatsubara.com/ "A5:SQL Mk-2 - フリーのSQLクライアント/ER図作成ソフト (松原正和)") のエクスポート機能を使って吸い上げたデータを再利用するのですが，例によって BOM (Byte Order Mark) が付いてるのですよ。

BOM は忘れた頃にやってくる（遠い目）

で，最近読んだ『[実用 Go言語](https://www.oreilly.co.jp/books/9784873119694/ "O'Reilly Japan - 実用 Go言語")』に Decorator Pattern で BOM を除去する方法が載っていた（8章）ので，早速試してみることにした。こんな感じ。

```go:sample1.go
package main

import (
    "fmt"
    "io"
    "strings"

    "github.com/spkg/bom"
)

const text = "\xEF\xBB\xBFhello"

func main() {
    fmt.Println([]byte(text))
    r := bom.NewReader(strings.NewReader(text))
    b, err := io.ReadAll(r)
    if err != nil {
        fmt.Println(err)
        return
    }
    fmt.Println(b)
}
```

これを実行すると

```
$ go run sample1.go 
[239 187 191 104 101 108 108 111]
[104 101 108 108 111]
```

と出力される。先頭の BOM が除去されているのがお分かりだろうか。

これとは別に [github.com/dimchansky/utfbom] パッケージってのがあって，同じように

```go:sample2
package main

import (
    "fmt"
    "io"
    "strings"

    "github.com/dimchansky/utfbom"
)

const text = "\xEF\xBB\xBFhello"

func main() {
    fmt.Println([]byte(text))
    r := utfbom.SkipOnly(strings.NewReader(text))
    b, err := io.ReadAll(r)
    if err != nil {
        fmt.Println(err)
        return
    }
    fmt.Println(b)
}
```

と書けば全く同じ出力を得られた。

[github.com/spkg/bom] パッケージは，ソースコードを見ると分かるが，とてもシンプルな作りになっていて， UTF-8 エンコーディングに限るならお手軽に使えるのがよい。もう一方の [github.com/dimchansky/utfbom] パッケージは UTF-8 以外に UTF-16 や UTF-32 にも対応していて，先程のコードを

```go
r, enc := utfbom.Skip(strings.NewReader(text))
```

と置き換えれば UTF テキストのエンコーディングも取得できる。必要に応じて使い分けるのがいいだろう。

ただし，いずれのパッケージも先頭の BOM しか取り除いてくれない。何らかの理由（BOM 付きテキストを安直に結合した場合とか）でテキストの先頭以外に紛れ込んでいる BOM があっても素通ししてしまう。まぁ，今時そういうケースは殆どないだろうが。

CSV ファイルは巨大になりがちで，システムの規模によってはすぐに十万レコードとか百万レコードとかになってしまう。この点で [csv].Reader 型はとてもよく出来ていて， Read() メソッドを使って順次アクセスで1レコードづつ切り出して返してくれる。

```go
func (r *Reader) Read() (record []string, err error)
```

これを活かすのであれば Decorator Pattern で入力をラッピングするのが最善だろう。たとえばこんな感じ。

```go:sample3.go
package main

import (
    "encoding/csv"
    "errors"
    "fmt"
    "io"
    "os"

    "github.com/spkg/bom"
)

func main() {
    file, err := os.Open("./sample3.csv")
    if err != nil {
        fmt.Fprintln(os.Stderr, err)
        return
    }
    defer file.Close()

    r := csv.NewReader(bom.NewReader(file))
    for {
        row, err := r.Read()
        if err != nil {
            if !errors.Is(err, io.EOF) {
                fmt.Fprintln(os.Stderr, err)
                return
            }
            break
        }
        fmt.Println(row)
    }
}
```

そもそも [csv].Reader 型自体が入力のラッパーである点に注目。 [io].Reader interface 型をベースにした Decorator Pattern で，動的な機能追加が簡単に出来てしまうのが嬉しい。

なお，文字エンコーディング変換も Decorator Pattern で実装できる。この記事では詳細は割愛するが，ググればあちこちに見つかると思うので探してみていただきたい。

[Go]: https://go.dev/ "The Go Programming Language"
[github.com/spkg/bom]: https://github.com/spkg/bom "spkg/bom: Strip UTF-8 byte order marks"
[github.com/dimchansky/utfbom]: https://github.com/dimchansky/utfbom "dimchansky/utfbom: Detection of the BOM and removing as necessary"
[csv]: https://pkg.go.dev/encoding/csv "csv package - encoding/csv - Go Packages"
[io]: https://pkg.go.dev/io "io package - io - Go Packages"
