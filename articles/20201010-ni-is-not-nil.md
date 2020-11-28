---
title: "nil == nil でないとき（または Go プログラマは息をするように依存を注入する）" # 記事のタイトル
emoji: "⌨" # アイキャッチとして使われる絵文字（1文字だけ）
type: "tech" # "tech" : 技術記事 / "idea" : アイデア記事
topics: ["go", "programming"] # タグ。["markdown", "rust", "aws"] のように指定する
published: true # 公開設定（true で公開）
---

今回も小ネタ。 [Go] でコードを書く人にはお馴染みの話題であるが，たまたま見つけたので。

https://medium.com/@shivi28/go-when-nil-nil-returns-true-a8a014abeffb

お題はこのコード：

```go
package main

import "fmt"

func main() {
    var a *string = nil
    var b interface{} = a
    fmt.Println("a == nil:", a == nil) // true
    fmt.Println("b == nil:", b == nil) // false
    fmt.Println("a == b:  ", a == b)   // true
}
```

実際の実行結果は[こちら](https://play.golang.org/p/wKADfQk3-li)。

[Go] における `nil` はポインタ値のある状態を示すもので，いわゆる「null 参照」を指している。それだけだったら `b == nil` は `true` になりそうなものだが， interface 型が絡むと少し複雑になる。

実は interface 型は，概念的には **型と値への参照を要素として持つ構造体** である。図で描くとこんな感じ（“[Go Data Structures: Interfaces](https://research.swtch.com/interfaces)” より引用）。

[![type interface](https://research.swtch.com/gointer2.png)](https://research.swtch.com/interfaces)

構造体がゼロ値であると言うためには構造体の要素全てがゼロ値である必要がある。今回の文脈で言うと interface 型のインスタンスが `nil` であると言うためには型と値（への参照）がいずれも `nil` でなければならない。

これを踏まえて先程のコードを眺めると

```go
var a *string = nil
var b interface{} = a
```

変数 `a` は「`*string` 型で `nil` 値」である。これを `interface{}` 型の `b` に代入することで `b` は「型は `*string`, 値は `nil`」という内部状態を持つ。つまり `b == nil` は `true` ではないのだ。一方 `a == b` はそれぞれの「値の比較」で，両者比較可能で同じ `nil` だから `true` になる，というわけ。

[Go] では interface 型のこの機能と性質で以って「構造型の部分型付け（structural subtyping）」を実現している。

よく [Go] の interface 型は「C++ や Java の templete や interface のようなもの」と説明されるが， C++ や Java の抽象型は基本的に「公称型の部分型付け（nominal subtyping）」であり根本の設計思想が異なる。ちなみに Rust の trait も公称型である。

構造型部分型付けの何が嬉しいかというと「依存の注入（dependency injection）」がとてつもなく簡単にできるのである。

先程の

```go
var a *string = nil
var b interface{} = a
```

であれば「`b` に `a` を**注入**している」と考えれば理解が容易になるだろう。つまり「`b == nil`」は「`b` に何か注入されているか？」を検査するための式と見なせる。ここで重要なのは `a` と `b` の型にはコード上で事前に明示された関係はない，という点だ。明示された関係はなくても「たまたま構造が同じ」なら注入できてしまうというのが [Go] の素敵で恐ろしいところである（笑）

まぁ，構造型部分型付けによる依存の注入があまりに便利なので総称型の導入が遅れた側面は否めないが（邪推）

💡 **[Go] プログラマは息をするように依存を注入する** 💡

## 参考

- [nil は nil — プログラミング言語 Go | text.Baldanders.info](https://text.baldanders.info/golang/nil-is-nil/)
- [それは Duck Typing ぢゃない（らしい） — しっぽのさきっちょ | text.Baldanders.info](https://text.baldanders.info/remark/2020/04/subtyping/)
- [継承できないなら注入すればいいじゃない！ | slide.Baldanders.info](https://slide.baldanders.info/shimane-go-2020-01-23/)

[Go]: https://golang.org/ "The Go Programming Language"
