---
title: "uint(1) - uint(2) の評価 または型付けなし定数について" # 記事のタイトル
emoji: "💻" # アイキャッチとして使われる絵文字（1文字だけ）
type: "tech" # "tech" : 技術記事 / "idea" : アイデア記事
topics: ["go", "programming"] # タグ。["markdown", "rust", "aws"] のように指定する
published: true # 公開設定（true で公開）
---

Twitter の TL には脊髄反射で書いてしまったが，この件は [Go] およびそのコンパイラの特徴をよく表していると思うので Zenn 記事としてちょろんと書いておく。

起点となる記事はこれ：

https://qiita.com/xu1718191411/items/c1e2457da384f186ae83

いわゆる符号なし整数の “$1-2$” が計算機的にどう評価されるかという観点では，この記事は概ね合っていると思う。ただし記事のタイトル通り愚直に

```go
package main

import "fmt"

func main() {
    fmt.Println(uint(1) - uint(2))
}
```

と書いたコードを[実行](https://play.golang.org/p/cfPxqIPfnY1)しようとすると

```
./prog.go:6:22: constant -1 overflows uint
```

などとコンパイルエラーになる。何故か。

実は [Go] には「型付けなし定数（untyped constant）」と呼ばれる仕様がある。

> Constants may be [typed](https://golang.org/ref/spec#Types) or untyped. Literal constants, true, false, iota, and certain [constant expressions](https://golang.org/ref/spec#Constant_expressions) containing only untyped constant operands are untyped. 
>
> A constant may be given a type explicitly by a [constant declaration](https://golang.org/ref/spec#Constant_declarations) or [conversion](https://golang.org/ref/spec#Conversions), or implicitly when used in a [variable declaration](https://golang.org/ref/spec#Variable_declarations) or an [assignment](https://golang.org/ref/spec#Assignments) or as an operand in an [expression](https://golang.org/ref/spec#Expressions). It is an error if the constant value cannot be represented as a value of the [respective](https://golang.org/ref/spec#Representability) type.
>
> (via “[The Go Programming Language Specification](https://golang.org/ref/spec)”)

この仕様によりリテラル値はいったん型付けなし定数として評価され型が決まった時点で再評価される。たとえば，符号なし整数は負値を取らないので，リテラル値の `-1` を `uint` にキャストしても

```go
x := uint(-1) // constant -1 overflows uint
```

とコンパイルエラーになる。

`uint(1) - uint(2)` は符号なし整数同士の引き算ぢゃないか！ と思われるかもしれないが，実際のコンパイラの挙動としては，リテラル値の部分を `1 - 2` と型付けなし定数として評価してから `uint` 型にキャストされる。最適化というやつだ。なので

```go
x := uint(1) - uint(2) // constant -1 overflows uint
```

もコンパイルエラーになるのだ。もちろん，それぞれのリテラル値をいったん（型付きの）変数に落とし込んでやれば

```go
package main

import "fmt"

func main() {
    var a, b uint = 1, 2
    fmt.Println(a - b) // 18446744073709551615
}
```

意図通り符号なし整数同士の引き算として機能する。ちなみに

```go
fmt.Println(18446744073709551615) // constant 18446744073709551615 overflows int
```

もコンパイルエラーとなる[^max64]。理由は分かるね（笑）

[^max64]: 18446744073709551615 は 0xffffffffffffffff と同じで，数式で書くなら $2^{64}-1$ と表現できる。 [Go] には階乗を表現する演算子はないが，2の階乗であれば `1<<64 - 1` と記述できる。

型付けなし定数といっても無限サイズの数値を扱えるわけではない。どこまでのサイズを扱えるかは基本的にはコンパイラ依存となるのだが，仕様としては

> Implementation restriction: Although numeric constants have arbitrary precision in the language, a compiler may implement them using an internal representation with limited precision. That said, every implementation must: 
> 
> - Represent integer constants with at least 256 bits.
> - Represent floating-point constants, including the parts of a complex constant, with a mantissa of at least 256 bits and a signed binary exponent of at least 16 bits.
> - Give an error if unable to represent an integer constant precisely.
> - Give an error if unable to represent a floating-point or complex constant due to overflow.
> - Round to the nearest representable constant if unable to represent a floating-point or complex constant due to limits on precision.
>
> (via “[The Go Programming Language Specification](https://golang.org/ref/spec)”)

といったあたりまでは保証されているようだ。

:::message
教訓： プログラムは動くコードだけが正義
:::

## 参考

以下の書籍の3.6.2章で「型付けなし定数」について詳しく解説されている。ご参照あれ！

https://www.amazon.co.jp/dp/B099928SJD

[Go]: https://golang.org/ "The Go Programming Language"
<!-- eof -->
