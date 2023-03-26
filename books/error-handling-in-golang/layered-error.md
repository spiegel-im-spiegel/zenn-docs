---
title: "エラーの階層化"
---

:::message
[Go] 1.20 から複数エラーをもつインスタンスを扱えるようになった。詳しくは「[複数のエラーを扱う](./multi-error)」を参照のこと。
:::


## 階層化エラーの導入

[Go] 1.13 から [errors] 等の標準パッケージでエラーの階層化が取り入れられた。具体的には Unwrap() 関数の導入である。

```go:errors/wrap.go
// Unwrap returns the result of calling the Unwrap method on err, if err's
// type contains an Unwrap method returning error.
// Otherwise, Unwrap returns nil.
func Unwrap(err error) error {
    u, ok := err.(interface {
        Unwrap() error
    })
    if !ok {
        return nil
    }
    return u.Unwrap()
}
```

たとえば[前節](./evaluations)で紹介した [os].PathError 型のメソッドは以下のように定義されている。

```go:os/error.go
// PathError records an error and the operation and file path that caused it.
type PathError struct {
    Op   string
    Path string
    Err  error
}

func (e *PathError) Error() string { return e.Op + " " + e.Path + ": " + e.Err.Error() }

func (e *PathError) Unwrap() error { return e.Err }
```

つまり [errors].Unwrap() の引数に [os].PathError 型のエラー・インスタンスを指定すると Err 要素を返してくれるわけだ。これによってエラーの原因を遡ることが可能になる。

実は [errors].Is() および [errors].As() 関数も内部で Unwrap() 関数を呼んでいて，たとえば以下のように原因を遡って評価を行うことができる。

```go
file, err := os.Open(filename)
switch {
case errors.Is(err, syscall.ENOENT):
    fmt.Fprintln(os.Stderr, "ファイルが存在しない")
default:
    fmt.Fprintln(os.Stderr, "その他のエラー")
}
```

## [fmt].Errorf() 関数によるエラーのラッピング

もうひとつ [Go] 1.13 からエラーに関する重要な仕様拡張があり， [fmt].Errorf() 関数で書式 `%w` を指定することでエラーのラッピングを行うことができるようになった。

例として以下のようなコードを考えてみる。

```go:sample.go
package main

import (
    "fmt"
    "os"
)

func checkFileOpen(path string) error {
    file, err := os.Open(path)
    if err != nil {
        return fmt.Errorf("error! : %w", err)
    }
    defer file.Close()
    return nil
}

func main() {
    if err := checkFileOpen("not-exist.txt"); err != nil {
        fmt.Fprintf(os.Stderr, "%#v\n", err)
        return
    }
}
```

これを実行すると

```
$ go run sample.go 
&fmt.wrapError{msg:"error! : open not-exist.txt: no such file or directory", err:(*os.PathError)(0xc00005a150)}
```

となり [os].PathError 型のインスタンスがラッピングされていることが分かる。したがって

```go:sample2.go
func main() {
    if err := checkFileOpen("not-exist.txt"); err != nil {
        var perr *os.PathError
        if errors.As(err, &perr) {
            switch {
            case errors.Is(err, syscall.ENOENT):
                fmt.Fprintf(os.Stderr, "\"%v\" ファイルが存在しない\n", perr.Path)
            default:
                fmt.Fprintln(os.Stderr, "その他の PathError")
            }
        } else {
            fmt.Fprintln(os.Stderr, "その他のエラー")
        }
        return
    }
}
```

などとすれば

```
$ go run sample2.go 
"not-exist.txt" ファイルが存在しない
```

といった評価も可能になる。

これで，エラーメッセージを解析するしかなかった [fmt].Errorf() 関数もだいぶ「使える」ようになっただろう（笑）

[Go]: https://go.dev/ "The Go Programming Language"
[io]: https://pkg.go.dev/io/ "io - The Go Programming Language"
[errors]: https://pkg.go.dev/errors/ "errors - The Go Programming Language"
[os]: https://pkg.go.dev/os/ "os - The Go Programming Language"
[fmt]: https://pkg.go.dev/fmt/ "fmt - The Go Programming Language"
[conversion]: https://go.dev/ref/spec#Conversions "The Go Programming Language Specification - The Go Programming Language"
[type assertion]: https://go.dev/ref/spec#Type_assertions "The Go Programming Language Specification - The Go Programming Language"
<!-- eof -->
