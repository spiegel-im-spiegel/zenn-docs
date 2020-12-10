---
title: "まずはキホンから"
---

## もはや「例外」は Legacy

私は C/C++ や Java などから来た人間なので [Go] を始めたばかりの頃は「例外（Exception）」のないエラーハンドリングに面食らったものだが，今ではすっかり慣れてしまった。

今年（2020年）になって Rust の勉強を少しだけ始めたが，改めて分かった。

💡 **もはや「例外」は Legacy だ！** 💡

たとえば Rust は[列挙型と match 式を組み合わせてエラーの抽出と評価を行う](https://text.baldanders.info/rust-lang/error-handling/ "エラー・ハンドリングのキホン")ことでエラー・ハンドリングを実装できる。

```rust
fn main() {
    let n = match parse_string("-1") {
        Ok(x) => x,
        Err(e) => panic!(e), //Output: thread 'main' panicked at 'Box<Any>', src/main.rs:8:19
    };
    println!("{}", n); //do not reach
}
```

実にスマート！

## 「例外」の問題は “goto” と同じ[^goto1]

[^goto1]: ちなみに [Go] の `goto` や ラベル付きの `break`, `continue` は[飛び先に制約](https://golang.org/test/goto.go)があり，どこにでもジャンプできるわけではない。

「例外」では，あるオブジェクトに関する記述が少なくとも2つ（たとえば try と catch）下手をすると3つ以上のスコープに分割されてしまう。しかもオブジェクトの状態ごと大域脱出するため，その状態（の可能性）の後始末をスコープ間で漏れなく矛盾なく記述しきらなければならない。

この一連に不備があれば，バグやリークやその他の脆弱性のもとになる。考えるだけで面倒である。

[Go] におけるエラーの扱いは，とにかく「シンプル」の一言に尽きる。以降から具体的に見てみよう。

## error 型

まずエラーを扱う組込み interface 型の error は以下のように定義される。

```go
// The error built-in interface type is the conventional interface for
// representing an error condition, with the nil value representing no error.
type error interface {
    Error() string
}
```

つまり，文字列を返す Error() メソッドを持つ型であれば全て error 型として扱うことができる。汎化にもほどがある（笑）

## エラーを含む処理の一連

しかも [Go] ではエラーを普通に関数の返り値として返す。

```go
file, err := os.Open(filename)
```

他に返すべき値があれば組（tuple）にして最後の要素に error 型のインスタンスを配置するのが慣例らしい。

検出したエラーは（投げ出さないでw）その場で評価してしまえばいい。

```go
file, err := os.Open(filename)
if err != nil {
    fmt.Fprintln(os.Stderr, err)
    return
}
```

また [if 構文][if]は内部に構文を含めることもできるので

```go
if err := file.Close(); err != nil {
    fmt.Fprintln(os.Stderr, err)
    return
}
```

てな感じに書くこともできる[^if]。

[^if]: [if 構文][if]内で宣言（:=）された変数は，そのスコープでのみ有効となる。同名変数の shadowing に注意。

Open と Close のように一連の処理が要求される場合は [defer 構文][defer]で後始末を先に書いてしまう。一連の処理をまとめるとこんな感じ。

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

## 一番簡単なエラー型

一番簡単なエラー型は [errors] 標準パッケージで定義されている。

```go:errors/errors.go
// New returns an error that formats as the given text.
// Each call to New returns a distinct error value even if the text is identical.
func New(text string) error {
    return &errorString{text}
}

// errorString is a trivial implementation of error.
type errorString struct {
    s string
}

func (e *errorString) Error() string {
    return e.s
}
```

これを使って，たとえば

```go:internal/oserror/errors.go
var (
    ErrInvalid    = errors.New("invalid argument")
    ErrPermission = errors.New("permission denied")
    ErrExist      = errors.New("file already exists")
    ErrNotExist   = errors.New("file does not exist")
    ErrClosed     = errors.New("file already closed")
)
```

などとエラー・インスタンスを定義できるわけだ。また [fmt].Errorf() 関数を使って

```go
package main

import (
    "fmt"
)

func main() {
    const name, id = "bueller", 17
    err := fmt.Errorf("user %q (id %d) not found", name, id)
    fmt.Println(err.Error())
}
```

のように，その場でエラーメッセージを作ってエラー・インスタンスとして生成することもできる。

[Go]: https://golang.org/ "The Go Programming Language"
[if]: https://golang.org/ref/spec#If_statements "The Go Programming Language Specification - The Go Programming Language"
[defer]: https://golang.org/ref/spec#Defer_statements "The Go Programming Language Specification - The Go Programming Language"
[errors]: https://golang.org/pkg/errors/ "errors - The Go Programming Language"
[os]: https://golang.org/pkg/os/ "os - The Go Programming Language"
[fmt]: https://golang.org/pkg/fmt/ "fmt - The Go Programming Language"
<!-- eof -->
