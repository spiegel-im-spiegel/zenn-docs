---
title: "リテラル値のポインタ" # 記事のタイトル
emoji: "💻" # アイキャッチとして使われる絵文字（1文字だけ）
type: "tech" # "tech" : 技術記事 / "idea" : アイデア記事
topics: ["go", "programming"] # タグ。["markdown", "rust", "aws"] のように指定する
published: true # 公開設定（true で公開）
---

今回も

https://gpl-reading.connpass.com/event/224161/

からの小ネタ。

たとえば Java なら

```java
System.out.println("hello".length()); // Outut: 5
```

みたいな記述ができる。 Java に限らず「オブジェクト指向」を謡っているプログラミング言語はリテラル表現をオブジェクトとして評価するため上述のような芸当ができるのだが， [Go] にはこれができない（[Go] では基本型リテラルは[型付けなし定数（untyped constant）](https://zenn.dev/spiegel/articles/20210813-untyped-constant "uint(1) - uint(2) の評価 または型付けなし定数について")として扱われる点に注意）。

そもそもリテラル表現で記述できる基本型は，それに紐づくメソッドを持たないので

```go
fmt.Println("Hello".String()) // "Hello".String undefined (type string has no field or method String)
```

とかやってもコンパイルエラーになるだけだし，以下のように

```go
s := &"Hello" // cannot take the address of "Hello"
```

リテラル表現から直接ポインタ値を得ることもできない。ちなみに

```go
s := &string("Hello") // cannot take the address of string("Hello")
```

と型を明示してもダメ。ただし

```go
s := "Hello"
fmt.Printf("%p\n", &s) // print pointer to variable
```

といった感じにインスタンスとして変数へ落とし込めばポインタ値を得ることは可能である。

ここで皆さん疑問に思わなかっただろうか。リテラル表現から直接ポインタ値が取れないなら，構造体リテラルで

```go
type Hello struct{}

func New() *Hello {
    return &Hello{}
}
```

みたいな記述はなぜコンパイルエラーにならないのか。実は私，今回の読書会で指摘されるまで全く疑問に思わなかった。不覚 orz

この話は『[プログラミング言語Go](https://www.amazon.co.jp/dp/4621300253/)』の「4.4.1 構造体リテラル」にさらりと書かれている。これによると

```go
h := &Hello{}
```

は

```go
h := new(Hello)
*h = Hello{}
```

と等価だと言うのだ[^mem1]。つまり `h := &Hello{}` は一種の syntax sugar として機能しているらしい。ちなみにメソッドを

[^mem1]: 念のために言うと [Go] では new() や make() といった組み込み関数で確保した領域がヒープ上に作られるとは限らない。最適化によってスタック上に積まれる可能性もある。


```go
func (h Hello) Say() string {
    return "Hello"
}
```

と定義すれば

```go
fmt.Println(Hello{}.Say()) // Hello
```

でちゃんと動く。更にメソッドレシーバを

```go
func (h *Hello) Say() string {
    return "Hello"
}
```

とポインタ型にした場合は

```go
fmt.Println((&Hello{}).Say()) // Hello
```

と括弧で明示すれば大丈夫。


```go
fmt.Println(&Hello{}.Say())
// cannot take the address of (&Hello{}).Say()
// cannot call pointer method on Hello{}
```

ではコンパイルエラーになる（`&` のスコープが `Hello{}.Say()` までなのが原因）。

[言語仕様](https://golang.org/ref/spec "The Go Programming Language Specification - The Go Programming Language")をよく読むと

>Calling the built-in function [new](https://golang.org/ref/spec#Allocation) or taking the address of a [composite literal](https://golang.org/ref/spec#Composite_literals) allocates storage for a variable at run time. Such an anonymous variable is referred to via a (possibly implicit) [pointer indirection](https://golang.org/ref/spec#Address_operators).
>(via “[The Go Programming Language Specification](https://golang.org/ref/spec#Variables)”)

と書かれていた。つまり

```go
fmt.Printf("%p\n", &[3]int{1, 2, 3})                 // print pointer to array
fmt.Printf("%p\n", &[]int{4, 5, 6})                  // print pointer to slice
fmt.Printf("%p\n", &map[string]string{"foo": "bar"}) // print pointer to map
```

もアリということか。今回もひとつ賢くなりました（笑）

[Go]: https://golang.org/ "The Go Programming Language"
<!-- eof -->
