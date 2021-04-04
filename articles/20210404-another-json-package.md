---
title: "goccy/go-json パッケージを試す" # 記事のタイトル
emoji: "💻" # アイキャッチとして使われる絵文字（1文字だけ）
type: "tech" # "tech" : 技術記事 / "idea" : アイデア記事
topics: ["go", "programming"] # タグ。["markdown", "rust", "aws"] のように指定する
published: true # 公開設定（true で公開）
---

先日 Twitter で見かけたのだが [github.com/goccy/go-json] という JSON ハンドリング・パッケージがあって，標準の [encoding/json] と置き換え可能かつ [encoding/json] や他の互換パッケージより速いと豪語している。

![](https://user-images.githubusercontent.com/209884/107126757-07ad3480-68f5-11eb-87aa-858cc5eacfcb.png)
![](https://user-images.githubusercontent.com/209884/107979940-bc84d700-7002-11eb-9647-869bbc25c9d9.png)
*via “[github.com/goccy/go-json]”*

まずは本当に [encoding/json] 標準パッケージと置き換えれるか試してみる。

ちょうど最近，絵文字情報を JSON 形式にまとめたパッケージ [spiegel-im-spiegel/emojis] を作ったので，こいつを読み込む処理を書いてみる。 [encoding/json] 標準パッケージを使うとこんな感じ。

```go:sample1.go
// +build run

package main

import (
    "bytes"
    "encoding/json"
    "fmt"
    "os"
    "strings"

    emoji "github.com/spiegel-im-spiegel/emojis/json"
    "github.com/spiegel-im-spiegel/fetch"
)

func getEmojiSequenceJSON() ([]byte, error) {
    u, err := fetch.URL("https://raw.githubusercontent.com/spiegel-im-spiegel/emojis/main/json/emoji-sequences.json")
    if err != nil {
        return nil, err
    }
    resp, err := fetch.New().Get(u)
    if err != nil {
        return nil, err
    }
    return resp.DumpBodyAndClose()
}

func main() {
    b, err := getEmojiSequenceJSON()
    if err != nil {
        fmt.Fprintln(os.Stderr, err)
        return
    }
    list := []emoji.EmojiSequence{}
    if err := json.NewDecoder(bytes.NewReader(b)).Decode(&list); err != nil {
        fmt.Fprintln(os.Stderr, err)
        return
    }

    fmt.Println("| Sequence | Shortcodes |")
    fmt.Println("| :------: | ---------- |")
    for _, ec := range list {
        var bldr strings.Builder
        for _, c := range ec.Shortcodes {
            bldr.WriteString(fmt.Sprintf(" `%s`", c))
        }
        fmt.Printf("| %v |%s |\n", ec.Sequence, bldr.String())
    }
}
```

これを動かしてみると

```
$ go run sample1.go
| Sequence | Shortcodes |
| :------: | ---------- |
| #️⃣ | `:hash:` `:keycap_#:` |
| *️⃣ | `:asterisk:` `:keycap_*:` `:keycap_star:` |
| 0️⃣ | `:zero:` `:keycap_0:` |
| 1️⃣ | `:one:` `:keycap_1:` |
...
```

となる。うんうん。

これを [github.com/goccy/go-json] パッケージに置き換える。

```diff go:sample2.go
import (
    "bytes"
-   "encoding/json"
    "fmt"
    "os"
    "strings"

+   "github.com/goccy/go-json"
    emoji "github.com/spiegel-im-spiegel/emojis/json"
    "github.com/spiegel-im-spiegel/fetch"
)
```

このコードで動かしてみると

```
$ go run sample2.go 
| Sequence | Shortcodes |
| :------: | ---------- |
| #️⃣ | `:hash:` `:keycap_#:` |
| *️⃣ | `:asterisk:` `:keycap_*:` `:keycap_star:` |
| 0️⃣ | `:zero:` `:keycap_0:` |
| 1️⃣ | `:one:` `:keycap_1:` |
...
```

と同じ出力になった。

じゃあ次は，このコードを流用してベンチマークを動かしてみようか。こんな感じでいいかな。

```go:json2_test.go
package jsonbench

import (
    "bytes"
    "encoding/json"
    "testing"

    another "github.com/goccy/go-json"
    emoji "github.com/spiegel-im-spiegel/emojis/json"
    "github.com/spiegel-im-spiegel/fetch"
)

func getMustEmojiSequenceJSON() []byte {
    u, err := fetch.URL("https://raw.githubusercontent.com/spiegel-im-spiegel/emojis/main/json/emoji-sequences.json")
    if err != nil {
        panic(err)
    }
    resp, err := fetch.New().Get(u)
    if err != nil {
        panic(err)
    }
    b, err := resp.DumpBodyAndClose()
    if err != nil {
        panic(err)
    }
    return b
}

var jsonText = getMustEmojiSequenceJSON()

func BenchmarkDecodeOrgPkg(b *testing.B) {
    for i := 0; i < b.N; i++ {
        _ = json.NewDecoder(bytes.NewReader(jsonText)).Decode(&([]emoji.EmojiSequence{}))
    }
}

func BenchmarkDecodeAnotherPkg(b *testing.B) {
    for i := 0; i < b.N; i++ {
        _ = another.NewDecoder(bytes.NewReader(jsonText)).Decode(&([]emoji.EmojiSequence{}))
    }
}
```

これを実行すると，手元の環境では以下の結果になった。

```
$ go test -bench Decode -benchmem
goos: linux
goarch: amd64
pkg: json2
cpu: Intel(R) Core(TM) i5-3470 CPU @ 3.20GHz
BenchmarkDecodeOrgPkg-4                 99      11095017 ns/op     2257059 B/op       17398 allocs/op
BenchmarkDecodeAnotherPkg-4            352       3386700 ns/op      907342 B/op        6403 allocs/op
PASS
ok      json2    3.172s
```

表にするとこんな感じ。

| 使用パッケージ                            |         実行時間 |   Alloc サイズ |       Alloc 回数 |
| ----------------------------------------- | ---------------: | -------------: | ---------------: |
| [encoding/json]                           | 11,095,017 ns/op | 2,257,059 B/op | 17,398 allocs/op |
| [goccy/go-json][github.com/goccy/go-json] |  3,386,700 ns/op |   907,342 B/op |  6,403 allocs/op |

おりょ。アロケーションの回数からまず違うのか。そりゃあ差がでるわな。実行時間全体では標準の3割まで圧縮されている。

こりゃあ，一考の価値ありか？

## 参考

https://zenn.dev/spiegel/articles/20210113-fetch
https://zenn.dev/spiegel/articles/20210322-emoji-shortcode-for-markdown
https://text.baldanders.info/remark/2021/04/emoji-list/

[Go]: https://golang.org/ "The Go Programming Language"
[encoding/json]: https://golang.org/pkg/encoding/json/ "json - The Go Programming Language"
[github.com/goccy/go-json]: https://github.com/goccy/go-json "goccy/go-json: Fast JSON encoder/decoder compatible with encoding/json for Go"
[spiegel-im-spiegel/emojis]: https://github.com/spiegel-im-spiegel/emojis "spiegel-im-spiegel/emojis: List of Emoji-Sequences"
<!-- eof -->
