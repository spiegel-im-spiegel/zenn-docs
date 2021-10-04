---
title: "リテラル値のポインタ" # 記事のタイトル
emoji: "💻" # アイキャッチとして使われる絵文字（1文字だけ）
type: "tech" # "tech" : 技術記事 / "idea" : アイデア記事
topics: ["go", "programming"] # タグ。["markdown", "rust", "aws"] のように指定する
published: false # 公開設定（true で公開）
---

今回も

https://gpl-reading.connpass.com/event/224161/

からの小ネタ。

たとえば Java なら

```java
System.out.println("hello".length()); // Outut: 5
```

みたいな記述ができる。 Java に限らず「オブジェクト指向」を謡っているプログラミング言語はリテラル表現をオブジェクトとして評価するため上述のような芸当ができるのだが， [Go] にはこれができない。

そもそもリテラル表現で記述できる基本型は，それに紐づくメソッドを持たないので

```go
fmt.Println("Hello".String()) // "Hello".String undefined (type string has no field or method String)
```

とかやってもコンパイルエラーになるだけだし，以下のように

```go
s := &"Hello" // cannot take the address of "Hello"
```

リテラル表現から直接ポインタ値を得ることもできない。ただし

```go
s := "Hello"
fmt.Printf("%p\n", &s) // output pointer value
```

という感じに変数へ落とし込めばポインタ値を得ることは可能である。

ここで皆さん疑問に思わなかっただろうか。リテラル表現から直接ポインタ値が取れないなら，構造体リテラルで

```go
type Hello struct{}

func New() *Hello {
    return &Hello{}
}
```

みたいな記述はなぜ通るのか。実は私，今回の読書会で指摘されるまで全く疑問に思わなかった。不覚 orz

この話は『[プログラミング言語Go](https://www.amazon.co.jp/dp/4621300253/)』の「4.4.1 構造体リテラル」にさらりと書かれている。これによると

```go
h := &Hello{}
```

は

```go
h := make(Hello)
*h = Hello{}
```

と等価だと言うのだ。つまり `h := &Hello{}` は一種の syntax sugar として機能しているらしい。

これを踏まえて考えると

```go
type Hello struct{}

func (h *Hello) Say() string {
    return "Hello"
}
```

と定義されているときに

```go
fmt.Println(&Hello{}.Say())
// Output:
// cannot take the address of (&Hello{}).Say()
// cannot call pointer method on Hello{}
```

と書くとコンパイルエラーになるのに

```go
fmt.Println((&Hello{}).Say()) // Hello
```

と括弧でくくるだけでコンパイルが通る理由が分かる。つまり `(&Hello{})` とすることで Say() メソッドを呼ぶ前に変数化されているわけだ。ちなみに

```go
fmt.Printf("%p\n", &Hello{}) // output pointer value
```

も通るけど

```go
fmt.Printf("%p\n", &"Hello")   // cannot take the address of "Hello"
fmt.Printf("%p\n", (&"Hello")) // cannot take the address of "Hello"
fmt.Printf("%p\n", &("Hello")) // cannot take the address of "Hello"
```

はいずれもコンパイルエラーになる。

今回もひとつ賢くなりました（笑）

[Go]: https://golang.org/ "The Go Programming Language"
<!-- eof -->
