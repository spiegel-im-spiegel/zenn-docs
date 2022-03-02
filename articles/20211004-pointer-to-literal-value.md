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

## リテラル値のポインタ

たとえば Java なら

```java
System.out.println("hello".length()); // Outut: 5
```

みたいな記述ができる。 Java に限らず「オブジェクト指向」を謡っているプログラミング言語はリテラル表現をオブジェクトとして評価するため上述のような芸当ができるのだが， [Go] にはこれができない（[Go] ではリテラルは[型付けなし定数（untyped constant）](https://zenn.dev/spiegel/articles/20210813-untyped-constant "uint(1) - uint(2) の評価 または型付けなし定数について")として扱われる点に注意）。

そもそもリテラル表現で記述できる基本型は，それに紐づくメソッドを持たないので

```go
fmt.Println("Hello".String()) // "Hello".String undefined (type string has no field or method String)
```

とかやっても「そんなメソッドはねーよ！」（←超意訳）と怒られるだけだし，以下のように

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
func (h *Hello) Say() string {
    return "Hello"
}
```

とポインタレシーバで定義した場合は

```go
fmt.Println(Hello{}.Say()) // cannot call pointer method on Hello{}
```

はダメだが（リテラルでは暗黙的にポインタ型への変換が出来ないため）

```go
fmt.Println((&Hello{}).Say()) // Hello
```

と括弧でくくって明示すればインスタンス化されるのでコンパイルエラーにはならない。なお

```go
fmt.Println(&Hello{}.Say())
// cannot take the address of (&Hello{}).Say()
// cannot call pointer method on Hello{}
```

ではコンパイルエラーになるのでご注意を（`&` のスコープが `Hello{}.Say()` までなのが原因）。

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

## 【おまけ】 リテラル値とメソッド

複合型（composite type）だけでなく基本型を基底型（underlying type）とする型の場合でも，たとえば

```go
type Name string

func (n Name) Say() string {
    return strings.Join([]string{"This is", string(n), "speaking!"}, " ")
}
```

という型とメソッドがあるとして

```go
fmt.Println(Name("Hayakawa").Say()) // This is Hayakawa speaking!
```

は問題なく動く（`Name("Hayakawa")` は関数ではなく型変換（type conversion）なので注意）。でも

```go
func (n *Name) Say() string {
    return strings.Join([]string{"This is", string(*n), "speaking!"}, " ")
}
```

と，メソッドレシーバをポインタ型にすると

```go
fmt.Println(Name("Hayakawa").Say())
// cannot call pointer method on Name("Hayakawa")
// cannot take the address of Name("Hayakawa")
```

でも

```go
fmt.Println((&Name("Hayakawa")).Say())
// cannot take the address of Name("Hayakawa")
```

でもコンパイル・エラーになる。前節で述べたように（`&struct{}{}` のような syntax sugar を除き）リテラル表現から直接ポインタ値を得ることは出来ないので，メソッド呼び出し時に暗黙的にポインタ型に変換できないからだ。

もちろん変数に落とし込んでしまえば

````go
n := Name("Hayakawa")
fmt.Println(n.Say())    // This is Hayakawa speaking!
fmt.Println((&n).Say()) // This is Hayakawa speaking!
````

無問題。ややこしい。

## 【2022-03-02 追記】 Slice と Map のアドレッシング

Twitter の「[プログラミング言語Go](https://twitter.com/i/communities/1498095077222400000)」コミュニティで教えてもらった話。

本編で map 型リテラルのポインタ値は取得できるという話をしたが

```go
fmt.Printf("%p\n", &map[string]string{"foo": "bar"}) // print pointer to map
```

角括弧（`[ ]`）を使って取得した要素のポインタ値は取得できずコンパイルエラーになる。

```go
fmt.Printf("%v", map[string]int{"foo": 1, "bar": 2}["foo"])  // 1
fmt.Printf("%p", &map[string]int{"foo": 1, "bar": 2}["foo"]) // cannot take the address of map[string]int{...}["foo"]
```

実はこれ，リテラル云々は関係なく map 型の仕様である。

```go
m := map[string]int{"foo": 1, "bar": 2}
fmt.Printf("%p", &m["foo"]) // cannot take the address of m["foo"]
```

これについて書籍『[プログラミング言語Go](https://www.amazon.co.jp/dp/B099928SJD)』の4.3章では以下のように説明している。

>マップの要素のアドレスが得られない理由の一つは、マップが大きくなる際に既存の要素が再びハッシングされて新たなメモリ位置へ移動するかもしれず、アドレスが無効になる可能性があるからです。
>（『[プログラミング言語Go](https://www.amazon.co.jp/dp/B099928SJD)』4.3章）

何かのきっかけで map の各要素の相対位置がランダムに変わっちゃうから要素のポインタ値は取れないよ，ということらしい。

一方で slice のほうは各要素の相対位置が決まってるので

```go
s := []int{1, 2, 3}
fmt.Printf("%p\n", &s[0])              // print pointer to element in slice
fmt.Printf("%p\n", &[]int{1, 2, 3}[0]) // print pointer to element in slice
```

通常の変数に対してもリテラルに対しても要素へのポインタ値を取ることができる。

[Go]: https://golang.org/ "The Go Programming Language"
<!-- eof -->
