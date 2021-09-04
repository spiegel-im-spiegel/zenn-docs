---
title: "iota 出現時の値はゼロとは限らない"
emoji: "💻" # アイキャッチとして使われる絵文字（1文字だけ）
type: "tech" # "tech" : 技術記事 / "idea" : アイデア記事
topics: ["go", "programming"] # タグ。["markdown", "rust", "aws"] のように指定する
published: true # 公開設定（true で公開）
---

今回も小ネタ。

「[第16回『プログラミング言語Go』オンライン読書会](https://gpl-reading.connpass.com/event/221591/)」の『[プログラミング言語Go](https://www.amazon.co.jp/dp/4621300253/)』の 3.6.1章「定数生成器 iota」で出た話で，書籍には

> const 宣言では、 iota の値はゼロから始まり、順番に個々の項目ごとに1増加します。

とあるが，ここで翻訳者であり読書会の主宰である柴田芳樹さんの解説があった。今回はその話。

元々 const では

```go
package main

import "fmt"

const (
    one = 1
    two
    three
    four
)

func main() {
    fmt.Println(one, two, three, four)
    // Output:
    // 1 1 1 1
}
```

と書くと[直前の定数と同じ値がセットされる](https://play.golang.org/p/3SJG2KlZ_iO)という特徴がある。この性質と定数生成器 iota を組み合わせることで


```go
package main

import "fmt"

const (
    one = 1 + iota
    two
    three
    four
)

func main() {
    fmt.Println(one, two, three, four)
    // Output:
    // 1 2 3 4
}
```

[ひとつづつインクリメントした値をセットする](https://play.golang.org/p/_UXJbnK8uyT)ことができる。じゃあ iota の初期値は常にゼロなのかというと，そこは微妙で，たとえば

```go
package main

import "fmt"

const (
    zero = "0"
    one  = 1
    two
    three
    four = iota
)

func main() {
    fmt.Println(zero, one, two, three, four)
    // Output:
    // 0 1 1 1 4
}
```

てな風に書くと [iota 出現時の値は 4 になる](https://play.golang.org/p/3RbtW0-jJis)。つまり iota は出現する前から（見かけ上[^iota1]）カウントしているわけだ。

[^iota1]: 正しくは iota はカウンタではない。この辺の話については拙文「[定数生成器 iota についてちゃんと書く](https://text.baldanders.info/golang/iota-constant-generator/)」で纒めてみた。

iota 出現時の値が常にゼロだと思いこんで，うっかり

```go
package main

import "fmt"

const (
    one = 1 + iota
    two
    three
    four
    zero = iota
)

func main() {
    fmt.Println(zero, one, two, three, four)
    // Output:
    // 4 1 2 3 4
}
```

てなコードを書くと `zero` がゼロにならず「[とひょーん](https://play.golang.org/p/-HvyRN4Doj5)」となってしまう。恥ずかしい話だが，実は昔このパターンでハマったことがあるのだ（テストが通らず，しばらく悩んだ）。

これを回避するには

```go
package main

import "fmt"

const (
    one = 1 + iota
    two
    three
    four
)

const (
    zero = iota
)

func main() {
    fmt.Println(zero, one, two, three, four)
    // Output:
    // 0 1 2 3 4
}
```

という感じに iota 毎に別の const 宣言で括ってやればよい。

:::message
教訓： iota は（別系統の定数と）混ぜるな，危険
:::

[Go]: https://golang.org/ "The Go Programming Language"
<!-- eof -->
