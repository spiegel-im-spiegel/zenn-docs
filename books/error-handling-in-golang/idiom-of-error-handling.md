---
title: "【付録3】 エラーハンドリングの慣習"
---

Twitter を眺めてたら，[エラーハンドリングもいくつか慣習（idiom）があるよね](https://twitter.com/nobonobo/status/1659554868431048704)という投稿があった。本編の中でもいくつか挙げているが「慣習」としてまとめるのは意義があると思うので，改めてここでまとめておこう。

## エラーの返し方

[Go] に例外処理はない。エラーを返す場合は返り値として必ず error 型の値を返す。また返り値が複数の組（tuple）になる場合は，最後の要素を error 型とする。

たとえば [os].Open() 関数は以下のように定義される。

```go
func Open(name string) (*File, error) { ... }
```

## エラーがない場合はリテラルの nil を返す

これは返り値を受けた側が error 値を nil と比較するからである。

```go
file, err := os.Open(filename)
if err != nil {
    fmt.Fprintln(os.Stderr, err)
    return
}
```

[2章](./basics "まずはキホンから")で述べた通り， error は interface 型のひとつであるが，そもそも [Go] の interface 型は他言語で言うところの「ボックス化」や「スマートポインタ」に相当するものであり，返ってきた error 値が nil であると言うためには型と値のいずれも nil でなければならない。つまり

```go
func foo() error {
    var err *Error = nil

    ...

    return err
}
```

のような書き方は NG である。

interface 型と nil の関係については以下の拙文を参考にどうぞ。

https://zenn.dev/spiegel/articles/20201010-ni-is-not-nil

## 返ってきたエラーはすぐに評価する

これは俗に “early return” と呼ばれている。喩えるなら迷路パズルを解くのに先に袋小路から塗りつぶしていくようなイメージである。こうして先に異常系の処理を潰してしまうことで本来の処理が見やすくなる。またそういう構成になるような「エラー設計」を考えるべきである。

たとえばエラーハンドリングが複雑なものはなるべく下位レイヤ側で処理して上位レイヤにはいわゆる sentinel error のみを返すようにする。そうすれば上位レイヤでは nil 比較または [errors].Is() 関数による比較のみの単純なハンドリングで済む（みんなが嫌がる単純作業w）。下位レイヤの詳細を知らなくてもエラーハンドリングが書けるわけだ。

なお，関数内のエラー処理をまとめて記述したい場合，以下のように defer で括る方法もある。

```go
func foo(v1 int, v2 string) (_ string, err error) {
    defer func() {
        if err != nil {
            err = fmt.Errorf("error in foo: %w", err)
        }
    }()
    v3, err := bar1(v1)
    if err != nil {
        return "", err
    }
    return bar2(v2, v3)
}
```

例外処理なんて邪法ですよ（笑）

## エラーを返す場合はエラー以外の返り値をゼロ値にする

これも返り値を受けた側のエラーハンドリングを単純にするための慣習と言える。

エラーが発生したのにも関わらず他の返り値に意味のある値がセットされているのであれば，エラーの値と他の値との組み合わせでより複雑なハンドリングを要求されるし，そのハンドリングのためには呼び出した関数の詳細をあらかじめ知っていなければならない。そして当然ながら，呼び出された関数の内部ロジックが変われば呼び出した側のハンドリングも変わる可能性がある。

エラーが発生したにも関わらず他の意味ある値を組で返さないといけないのであれば，それはプログラミング設計に何らかの問題がある可能性が高い。とはいえ，何事にも例外はあるからなぁ...

## Panic はリカバリしない，または Recover が必要な Panic は投げない

Panic については[付録1](./panics "【付録1】 Panic と Recover")を参照のこと。

Panic をいわゆる「例外処理」の代わりに使ってはいけない。 [Go] において panic は続行不可能な致命的エラーである。ブログなどでよく見るサンプルコードで，エラーを返す代わりに panic を投げるコードを見かけるが，絶対に真似してはいけない（同様に安直に [log].fatal() あるいは [os].Exit() するのも悪手）。

## エラーのラッピング

これは [Go] 1.13 以降からの慣習になると思うが，あるエラー値を評価した結果として別のエラー値を返す場合， [fmt].Errorf() などを使って元のエラーをラッピングする。

```go
func checkFileOpen(path string) error {
    file, err := os.Open(path)
    if err != nil {
        return fmt.Errorf("error! : %w", err)
    }
    defer file.Close()
    return nil
}
```

ラッピングされたエラーは [errors].Is() または [errors].As() 関数で捕捉可能である。詳しくは[4章](./layered-error "エラーの階層化")および[付録2](./multi-error "【付録2】 複数のエラーを扱う")を参照のこと。そうそう，[拙作の errs パッケージ]("./error-logging" "ぼくがかんがえたさいきょうのえらーろぐ")もよろしく（笑）



[Go]: https://go.dev/ "The Go Programming Language"
[errors]: https://pkg.go.dev/errors/ "errors - The Go Programming Language"
[fmt]: https://pkg.go.dev/fmt/ "fmt - The Go Programming Language"
[fs]: https://pkg.go.dev/fs/ "fs - The Go Programming Language"
[os]: https://pkg.go.dev/os/ "os - The Go Programming Language"
[io]: https://go.dev/pkg/io/ "io - The Go Programming Language"
[log]: https://go.dev/pkg/log/ "log - The Go Programming Language"
[hashicorp/go-multierror]: https://github.com/hashicorp/go-multierror "hashicorp/go-multierror: A Go (golang) package for representing a list of errors as a single error."
[errs]: https://github.com/goark/errs "goark/errs: Error handling for Golang"
<!-- eof -->
