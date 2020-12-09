---
title: "サードパーティのパッケージ"
---

標準の [errors] パッケージで階層化エラーの基本機能が提供されたが，サードパーティの汎用エラー・パッケージではもう少し高機能なものもある。以下にいくつか紹介してみよう。

## [pkg/errors]

[pkg/errors] は昔から人気の高い汎用エラー・パッケージで，最近のバージョンでは [Go] 1.13 以降の [errors] 標準パッケージと置き換えて使うこともできるようになった。

面白いのはエラーにスタック情報を付加できる点で

```go:sample3.go
package main

import (
    "fmt"
    "os"

    "github.com/pkg/errors"
)

func checkFileOpen(path string) error {
    file, err := os.Open(path)
    if err != nil {
        return errors.WithStack(err)
    }
    defer file.Close()
    return nil
}

func main() {
    if err := checkFileOpen("not-exist.txt"); err != nil {
        fmt.Fprintf(os.Stderr, "%+v\n", err)
        return
    }
}
```

のように errors.WithStack() 関数でラッピングしたエラー・インスタンスを `%+v` 書式で表示すると

```
$ go run sample3.go 
open not-exist.txt: no such file or directory
main.checkFileOpen
    /home/username/path/to/sample3.go:13
main.main
    /home/username/path/to/sample3.go:20
runtime.main
    /usr/local/go/src/runtime/proc.go:204
runtime.goexit
    /usr/local/go/src/runtime/asm_amd64.s:1374
```

てな感じにエラー発生時のスタック情報を吐き出すことができる。さらに

```go
if err != nil {
    return errors.Wrapf(err, "open error (%s)", path)
}
```

のように errors.Wrap() あるいは errors.Wrapf() 関数を使ってエラーメッセージを付加することもできる。

## [hashicorp/go-multierror]

コンテナ操作や goroutine を使った並行処理などで複数のエラーをまとめて処理する場合がある。複数のエラーをまとめて扱えるサードパーティ・パッケージはいくつかあるが，個人的には [hashicorp/go-multierror] がシンプルでお気に入りである。

たとえば，こんな感じに書ける。

```go
func main() {
    paths := []string{"not-exist1.txt", "not-exist2.txt"}
    var result *multierror.Error
    for _, path := range paths {
        if err := checkFileOpen(path); err != nil {
            result = multierror.Append(result, err)
        }
    }
    if err := result.ErrorOrNil(); err != nil {
        fmt.Fprintln(os.Stderr, err)
    }
    // Output:
    // 2 errors occurred:
    //     * error! : open not-exist1.txt: no such file or directory
    //     * error! : open not-exist2.txt: no such file or directory
}
```

また [errors].Is() や [errors].As() を使った評価もできる。

```go
if err := result.ErrorOrNil(); err != nil {
    var perr *os.PathError
    if errors.As(err, &perr) && errors.Is(perr, syscall.ENOENT) {
        fmt.Fprintf(os.Stderr, "\"%v\" ファイルが存在しない\n", perr.Path)
    } else {
        fmt.Fprintln(os.Stderr, "その他のエラー")
    }
}
// Output:
// "not-exist1.txt" ファイルが存在しない
```

簡単・便利！

## [golang.org/x/xerrors]

[errors] 標準パッケージで追加された階層エラー機能の[元ネタ](https://go.googlesource.com/proposal/+/master/design/29934-error-values.md "Proposal: Go 2 Error Inspection")的なパッケージで今でも割と使われているが，可能であれば [errors] 標準パッケージへ移行すべきだろう。

[golang.org/x/xerrors] パッケージにあって [errors] 標準パッケージにない機能として `%+v` 書式を使ったスタック情報の吐き出しがあるが，これについては先に紹介した [pkg/errors] のほうが設計がシンプルでお勧めである。併せて検討していただきたい。

[Go]: https://golang.org/ "The Go Programming Language"
[errors]: https://golang.org/pkg/errors/ "errors - The Go Programming Language"
[pkg/errors]: https://github.com/pkg/errors "pkg/errors: Simple error handling primitives"
[hashicorp/go-multierror]: https://github.com/hashicorp/go-multierror "hashicorp/go-multierror: A Go (golang) package for representing a list of errors as a single error."
[golang.org/x/xerrors]: https://pkg.go.dev/golang.org/x/xerrors "xerrors · pkg.go.dev"
<!-- eof -->
