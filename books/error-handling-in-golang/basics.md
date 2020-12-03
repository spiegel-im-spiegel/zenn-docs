---
title: "まずはキホンから"
---

[Go] におけるエラーの扱いは，とにかく「シンプル」の一言に尽きる。

まずエラーを扱う組込み interface 型の error は以下のように定義される。

```go
// The error built-in interface type is the conventional interface for
// representing an error condition, with the nil value representing no error.
type error interface {
    Error() string
}
```

つまり，文字列を返す Error() メソッドを持つ型であれば全て error 型として扱うことができる。

[Go] ではエラーを普通に関数の返値として返す。

```go
file, err := os.Open(filename)
```

他に返すべき値があれば組（tuple）にして最後の要素に error 型のインスタンスを配置するのが慣例らしい。

検出したエラーは（どこにも投げないで）その場で評価してしまえばよい。

```go
file, err := os.Open(filename)
if err != nil {
    fmt.Fprintln(os.Stderr, err)
    return
}
```

また if 構文は内部に構文を含めることができるので

```go
if err := file.Close(); err != nil {
    fmt.Fprintln(os.Stderr, err)
    return
}
```

なんてな感じに書くこともできる。

Open と Close のように一連の処理が要求される場合は defer 構文で後始末を先に書いてしまう。

```go
defer func() {
    if err := file.Close(); err != nil {
        fmt.Fprintln(os.Stderr, err)
    }
}()
```

まとめるとこんな感じ。

```go
file, err := os.Open(filename)
if err != nil {
    fmt.Fprintln(os.Stderr, err)
    return
}
defer func() {
    if err := file.Close(); err != nil {
        fmt.Fprintln(os.Stderr, err)
    }
}()
```

これが [Go] の基本的な書き方。特徴的なのは，あるオブジェクトに纏わる処理をセットで記述できる点である。とても文芸的なコードであるとも言える[^bb]。

[^bb]: これからのコードは「文芸的」であることが必要条件だと思う。何故ならエンジニアにとって最も信頼できる「設計書」は（動いている）コードだからだ。コードをひとりで考えてひとりで書いてひとりで使ってひとりでメンテナンスするなら（本人さえ理解していれば）文芸的である必要はないかもしれない。が，実用的なコードでそんな状況はもはやありえない。コードにおいても暗黙知をできるだけ排除していくことが重要である。










[Go]: https://golang.org/ "The Go Programming Language"
<!-- eof -->
