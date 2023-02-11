---
title: "【付録2】 複数のエラーを扱う"
---

コンテナ操作や並行処理など，単一のエラーではなく複数のエラーをまとめて処理したい場合がある。サードパーティでは [hashicorp/go-multierror] など複数エラーを扱うパッケージがあるが， 2023年2月にリリースされた [Go] 1.20 から標準で複数エラーを扱えるようになった。

## [errors].Is() および [errors].As() 関数の拡張

[Go] 1.13 から「[エラーの階層化](./layered-error)」をサポートするためにエラーインスタンスに対して Unwrap() error メソッドがの有無を検査するようになったが， 1.20 ではこれに加えて Unwrap() []error メソッドも対象となった。たとえば [errors].Is() 関数は

```go:errors/wrap.go
func Is(err, target error) bool {
    if target == nil {
        return err == target
    }

    isComparable := reflectlite.TypeOf(target).Comparable()
    for {
        if isComparable && err == target {
            return true
        }
        if x, ok := err.(interface{ Is(error) bool }); ok && x.Is(target) {
            return true
        }
        switch x := err.(type) {
        case interface{ Unwrap() error }:
            err = x.Unwrap()
            if err == nil {
                return false
            }
        case interface{ Unwrap() []error }:
            for _, err := range x.Unwrap() {
                if Is(err, target) {
                    return true
                }
            }
            return false
        default:
            return false
        }
    }
}
```

と拡張されている。 Unwrap() error と Unwrap() []error を型 switch を使って検査しているのがおわかりだろうか。 [errors].As() 関数関数についても同様の拡張がされている。

ここで Unwrap() error と Unwrap() []error は，メソッドとして同時に定義できない点に注意。さらに既存の [errors].Unwrap() 関数は原因エラーが単一の場合（Unwrap() error メソッドを備えている）にのみ値を返し Unwrap() []error メソッドに対応した関数は用意されていないらしい。

したがって，原因エラーが複数ある（つまり Unwrap() []error を備えている）場合に原因エラーのリストを取得したければ

```go
func Unwraps(err error) []error {
    if err != nil {
        if es, ok := err.(interface {
            Unwrap() []error
        }); ok {
            return es.Unwrap()
        }
    }
    e := errors.Unwrap(err)
    if e != nil {
        return []error{e}
    }
    return nil
}
```

みたいな関数を自前で用意する必要があるだろう。

:::message
拙作の [goark/errs][errs] パッケージでは [v1.2.1](https://text.baldanders.info/release/2023/02/errs-package-v1_2_1-is-released/ "goark/errs パッケージ v1.2.1 をリリースした") から上述の [errs].Unwraps() 関数を用意した。
:::

ともかく，これで複数エラーを持つエラーインスタンスについても [errors].Is() および [errors].As() 関数を使ってハンドリングできるようになった。

## [errors].Join() 関数の追加と [fmt].Errorf 関数の拡張

[Go] 1.20 から複数エラーを持つインスタンスを作るため [errors].Join() 関数が追加された。こんな感じに使える。

```go
package main

import (
    "errors"
    "fmt"
    "io"
    "os"
)

func main() {
    err := errors.Join(os.ErrInvalid, io.EOF)
    fmt.Println(err)
    fmt.Println("Error is EOF ? >", errors.Is(err, io.EOF))
}
```

これを実行するとこんな風に出力される。

```
invalid argument
EOF
Error is EOF ? > true
```

改行で区切るのか。  

さらに [fmt].Errorf 関数について `%w` を複数使えるようになった。こんな感じ。

```go
import (
    "errors"
    "fmt"
    "io"
    "os"
)

func main() {
    err := fmt.Errorf(`multiple errors: "%w" and "%w"`, os.ErrInvalid, io.EOF)
    fmt.Println(err)
    fmt.Println("Error is EOF ? >", errors.Is(err, io.EOF))
}
```

これを実行するとこんな風に出力される。

```
multiple errors: "invalid argument" and "EOF"
Error is EOF ? > true
```

`%w` を複数使えるようになったのは嬉しい。これでより柔軟なエラーハンドリングができるようになるだろう。

[Go]: https://golang.org/ "The Go Programming Language"
[errors]: https://pkg.go.dev/errors/ "errors - The Go Programming Language"
[fmt]: https://pkg.go.dev/fmt/ "fmt - The Go Programming Language"
[os]: https://pkg.go.dev/os/ "os - The Go Programming Language"
[io]: https://golang.org/pkg/io/ "io - The Go Programming Language"
[hashicorp/go-multierror]: https://github.com/hashicorp/go-multierror "hashicorp/go-multierror: A Go (golang) package for representing a list of errors as a single error."
[errs]: https://github.com/goark/errs "goark/errs: Error handling for Golang"
<!-- eof -->
