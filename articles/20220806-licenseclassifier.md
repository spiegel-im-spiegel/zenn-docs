---
title: "ライセンスファイルからライセンスを推定する" # 記事のタイトル
emoji: "💻" # アイキャッチとして使われる絵文字（1文字だけ）
type: "tech" # "tech" : 技術記事 / "idea" : アイデア記事
topics: ["go", "programming"] # タグ。["markdown", "rust", "aws"] のように指定する
published: true # 公開設定（true で公開）
---

たとえばリポジトリ直下に LICENSE というファイルがあるとして，このファイルが実際に何のライセンスを指しているか機械的に調べる方法はないだろうか。実は Google による [Go] パッケージが公開されている[^g1]。

[^g1]: ただし README.md には “This is not an official Google product” とあり Google 公式パッケージではないことが明記されている。ご注意を。

https://github.com/google/licenseclassifier

私は以前からこのパッケージを利用しているのだが，開発の主力が v2 系に移っているようだ。2022-07-22 に [v2.0.0-pre6](https://github.com/google/licenseclassifier/releases/tag/v2.0.0-pre6) がリリースされていた。さっそく試してみることにする。

今回のサンプルコードはこんな感じ。

```go:sample.go
package main

import (
    "flag"
    "fmt"
    "os"

    "github.com/google/licenseclassifier/v2/assets"
)

func main() {
    flag.Parse()
    args := flag.Args()
    if len(args) < 1 {
        fmt.Fprintln(os.Stderr, os.ErrInvalid)
        return
    }
    file, err := os.Open(args[0])
    if err != nil {
        fmt.Fprintln(os.Stderr, err)
        return
    }
    defer file.Close()

    c, err := assets.DefaultClassifier()
    if err != nil {
        fmt.Fprintln(os.Stderr, err)
        return
    }

    res, err := c.MatchFrom(file)
    if err != nil {
        fmt.Fprintln(os.Stderr, err)
        return
    }
    if len(res.Matches) == 0 {
        fmt.Fprintln(os.Stderr, args[0], "is not license file.")
        return
    }
    for _, m := range res.Matches {
        fmt.Println(m.MatchType, m.Name, )
    }
}
```

手順としては

1. コマンドライン引数で指定したファイルを開く
2. `assets.DefaultClassifier()` で解析のための辞書情報（`*classifier.Classifier` 型）を取得する
3. `MatchFrom()` メソッドでファイルを解析し，結果を表示する

という感じ。では，実際に動かしてみよう。

```
$ go run sample.go ./LICENSE 
License Apache-2.0
```

というわけで，指定した LICENSE ファイルは `License` タイプの `Apache-2.0` ライセンスであることが分かった。よーし，うむうむ，よーし。

まだ正式リリースではないようだが，使えるレベルに達してると思う。上手く利用していただきたい。

[Go]: https://go.dev/ "The Go Programming Language"
