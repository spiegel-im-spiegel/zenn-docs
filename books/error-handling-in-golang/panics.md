---
title: "【付録1】 Panic と Recover"
---

ゼロ除算を行ったり配列などで領域外を参照・設定しようとしたりヒープメモリが不足したり... 等々，致命的なエラーが発生する場合がある。

```go:sample7.go
package main

import "fmt"

func main() {
    foo()
}

func foo() {
    numbers := []int{0, 1, 2}
    fmt.Println(numbers[3])
}
```

これを実行すると

```
$ go run sample7.go
panic: runtime error: index out of range [3] with length 3

goroutine 1 [running]:
main.foo()
    /home/username/path/to/sample7/sample7.go:11 +0x1b
main.main()
    /home/username/path/to/sample7/sample7.go:6 +0x25
exit status 2
```

といった感じになり，大域脱出させてアプリケーションを強制終了させているのが分かる。この仕組みを panic と呼ぶ。

panic は意図的に発生させることもできる。

```go:sample8.go
package main

func main() {
    foo()
}

func foo() {
    panic("Panic!")
}
```

これを実行すると

```
$ go run sample8.go
panic: Panic!

goroutine 1 [running]:
main.foo(...)
    /home/username/path/to/sample8/sample8.go:8
main.main()
    /home/username/path/to/sample8/sample8.go:4 +0x39
exit status 2
```

てな感じになる。

また panic を recover することもできる[^recover1]。

[^recover1]: recover は [defer 構文][defer]とともに使用する。つまり panic 発生時でも [defer 構文][defer]で予約された処理は実行される。

```go:sample9.go
package main

import (
    "errors"
    "fmt"
)

func main() {
    err := bar()
    if err != nil {
        fmt.Printf("%#v\n", err)
        return
    }
    fmt.Println("Normal End.")
}

func bar() (err error) {
    defer func() {
        if rec := recover(); rec != nil {
            err = fmt.Errorf("Recovered from: %w", rec)
        }
    }()

    foo()
    return
}

func foo() {
    panic(errors.New("Panic!"))
}
```

これを実行すると

```
$ go run sample9.go
&fmt.wrapError{msg:"Recovered from: Panic!", err:(*errors.errorString)(0xc00010a040)}
```

といった感じになる。 panic を recover() 関数で捕まえて通常の error として返しているのがお分かりだろうか。

一般的に panic はアプリケーション内で続行不可能な致命的エラーが発生した場合に投げられる。

まぁ，ゼロ除算や領域外アクセスのようなエラーは panic が発生する前に回避するコードにすべきだが，ヒープメモリ不足のような回避不能なエラーの場合は panic が投げられるのもやむを得ないだろう。しかし，その場合でも recover して処理を継続させることに殆ど意味はない。

例外的な使い方として [bytes].Buffer では，メモリ確保で panic が発生した際に recover で捕まえ，定義済みの error インスタンスに入れ替えて panic を投げ直している。

```go
// makeSlice allocates a slice of size n. If the allocation fails, it panics
// with ErrTooLarge.
func makeSlice(n int) []byte {
    // If the make fails, give a known error.
    defer func() {
        if recover() != nil {
            panic(ErrTooLarge)
        }
    }()
    return make([]byte, n)
}
```

このような用途で recover を使うことはあり得る。

また再帰処理中に続行不能なエラーが発生した場合に panic を投げてトップレベルの関数に一気に復帰するような使い方をすることもあるらしい。この場合，トップレベルの関数は（続行不可能なら）改めて panic を投げるか（処理続行できる根拠があるのなら）通常の error を返すことになる[^recover2]。

[^recover2]: サーバ用途などでプロセスを落とせない場合に recover で回避することもあるそうだが，既に続行不可能な状態で無理やりプロセスを続行するのが正しい動きなのかどうかは疑問が残る。

いずれにしろ，外部パッケージが（何らかの理由で）投げた panic を安易に拾って「例外処理」すべきではないし， panic を投げる側も本当にそれが正しいハンドリングなのかよくよく考えるべきだろう。

なお，ビルド時（`go run` コマンド時を含む）に `-trimpath` オプションを付けるとスタック情報吐き出し時にフルパスでの表示を抑制できる。

```
$ go run -trimpath sample8.go
panic: Panic!

goroutine 1 [running]:
main.foo(...)
    command-line-arguments/sample8.go:8
main.main()
    command-line-arguments/sample8.go:4 +0x39
exit status 2
```

開発中はともかく，バイナリを公にリリースする際に（たとえ Docker 上でビルドするにしても）開発環境のパスが丸見えなのはどうかと思うので，リリース用ビルドのスクリプトに `-trimpath` オプションを付けてビルドするよう手を加えておくといいだろう。

[Go]: https://go.dev/ "The Go Programming Language"
[defer]: https://go.dev/ref/spec#Defer_statements "The Go Programming Language Specification - The Go Programming Language"
[bytes]: https://pkg.go.dev/bytes/ "bytes - The Go Programming Language"
<!-- eof -->
