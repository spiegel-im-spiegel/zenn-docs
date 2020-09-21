---
title: "Switch 文のナゾ（ってほどでもない）" # 記事のタイトル
emoji: "⌨" # アイキャッチとして使われる絵文字（1文字だけ）
type: "tech" # "tech" : 技術記事 / "idea" : アイデア記事
topics: ["go", "programming"] # タグ。["markdown", "rust", "aws"] のように指定する
published: true # 公開設定（true で公開）
---

いや，[面白い tweet](https://twitter.com/go100and1/status/1307689651998605312) を見かけたので（[解答](https://twitter.com/go100and1/status/1308004478709194753)も出てるからええやろ）。

@[tweet](https://twitter.com/go100and1/status/1307689651998605312)

見事に引っかかっちまったよ（笑）

実際に書いて実行してみれば分かる。

```go
package main

func f() bool {
    return false
}

func main() {
    switch f()
    {
    case true:
        println(1)
    case false:
        println(0)
    }
}
```

[試して](https://play.golang.org/p/qd_XCbpEs6d)みれば分かるが，実は 0 ではなく 1 が出力される。ポイントは `switch` 文の開始ブレスの位置である。

[Go] の[仕様書][Go Spec]によると，式評価の `switch` 文は以下のように定義されている。

```
ExprSwitchStmt = "switch" [ SimpleStmt ";" ] [ Expression ] "{" { ExprCaseClause } "}" .
ExprCaseClause = ExprSwitchCase ":" StatementList .
ExprSwitchCase = "case" ExpressionList | "default" .
```

これを見ると `switch` トークンと開始ブレスの間に文と式を並べて記述できることが分かる。つまり上述のコードは

```go
switch _ = f(); true {
case true:
    println(1)
case false:
    println(0)
}
```

と等価なのである。

なお最初のコードを整形ツールにかけると

```go
package main

func f() bool {
    return false
}

func main() {
    switch f(); {
    case true:
        println(1)
    case false:
        println(0)
    }
}
```

と整形される。これなら分かりやすいよね。

というわけで [Go] のコードを書いたらこまめに整形して確認しましょう。まぁ，最近のエディタは自動で整形してくれるものもあるけど。

[Go]: https://golang.org/ "The Go Programming Language"
[Go Spec]: https://golang.org/ref/spec "The Go Programming Language Specification - The Go Programming Language"
