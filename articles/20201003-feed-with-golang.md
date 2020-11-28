---
title: "フィードを取得する Go 言語パッケージ" # 記事のタイトル
emoji: "⌨" # アイキャッチとして使われる絵文字（1文字だけ）
type: "tech" # "tech" : 技術記事 / "idea" : アイデア記事
topics: ["go", "programming"] # タグ。["markdown", "rust", "aws"] のように指定する
published: true # 公開設定（true で公開）
---

「[シェルスクリプトで作る Twitter bot 作成入門](https://zenn.dev/mattn/books/bb181f3f4731920f29a5 "シェルスクリプトで作る Twitter bot 作成入門 | Zenn")」を見て簡単なボットでも作ろうかと色々と調べているのだが[^ifttt1]，ブログ等が公開している RSS/Atom フィードを取得する構造が簡単な [Go] パッケージがないかとググってみたら丁度いいのがあった。

[^ifttt1]: いや Twitter にブログ更新情報等を流し込むのに IFTTT を使ってるんだけど，ほとんど使ってないのに[アップグレードして金払え](https://forest.watch.impress.co.jp/docs/news/1278901.html "無償版“IFTTT”で利用可能なアプレットは3つまでに ～超過分は10月8日にアーカイブ - 窓の杜")って五月蝿くってさぁ。ほぼ毎日メールを投げてくるの。スパムなサービスはお断りだよ。

- [mmcdole/gofeed: Parse RSS, Atom and JSON feeds in Go][mmcdole/gofeed]

[mmcdole/gofeed] が優れているのは，フィードの種別に関わらず食べてくれて，統一された構造体に落とし込んでくれるところ。たとえばこんな感じ。

```go:sample.go
package main

import (
    "fmt"
    "os"
    "time"

    "github.com/mmcdole/gofeed"
)

func main() {
    feed, err := gofeed.NewParser().ParseURL("https://zenn.dev/spiegel/feed")
    if err != nil {
        fmt.Fprintln(os.Stderr, err)
        return
    }
    fmt.Println(feed.Title)
    fmt.Println(feed.FeedType, feed.FeedVersion)
    for _, item := range feed.Items {
        if item == nil {
            break
        }
        fmt.Println(item.Title)
        fmt.Println("\t->", item.Link)
        fmt.Println("\t->", item.PublishedParsed.Format(time.RFC3339))
    }
}
```

このコードの実行結果はこんな感じになる。

```
$ go run sample.go 
Spiegelさんのフィード
rss 2.0
GitHub Actions で Go パッケージの CI 作業を一通り行う
        -> https://zenn.dev/spiegel/articles/20200929-using-golangci-lint-action
        -> 2020-09-29T10:27:34Z
errors.Is, errors.As は（単なる）比較関数ではない
        -> https://zenn.dev/spiegel/articles/20200926-error-handling-with-golang
        -> 2020-09-26T06:54:49Z
...
```

標準の [context] パッケージにも対応してて，たとえば並行処理（goroutine）下で

```go
ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
defer cancel()
feed, err := gofeed.NewParser().ParseURLWithContext("https://zenn.dev/spiegel/feed", ctx)
if err != nil {
    return
}
...
```

みたいにキャンセル・イベントを絡めて書くこともできるらしい。

よしよし。使えそうだな。

## 参考

https://text.baldanders.info/golang/ticker/

[Go]: https://golang.org/ "The Go Programming Language"
[context]: https://golang.org/pkg/context/ "context - The Go Programming Language"
[mmcdole/gofeed]: https://github.com/mmcdole/gofeed "mmcdole/gofeed: Parse RSS, Atom and JSON feeds in Go"
