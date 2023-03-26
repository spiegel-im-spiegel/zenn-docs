---
title: "エラー評価のいろいろ"
---

エラーハンドリングを行うためにはまず何らかの形でエラーを評価する必要がある。 [Go] におけるエラー評価は大雑把に以下の4つに分けられるだろう。

1. nil との比較
2. インスタンスの同値性
3. インスタンスのボックス化解除
4. Error() メソッドの返り値を解析する

以降，ひとつずつ見ていこう。

## nil との比較（エラーの有無の判定）

おそらく [Go] のコードで最もよく見かけるパターンは

```go
if err != nil {
    ...
}
```

だろう。最近の（[Go] の支援機能を備えた）高機能エディタなら `iferr` と入力して [Tab] キーを押すと

```go
if err != nil {
    return
}
```

まで展開してくれたりする。重宝してます，ホンマ（笑）

[前節](./basics)でも述べたように error は interface 型のひとつだが，そもそも interface 型の機能はボックス化（boxing）の一種と見なせる[^boxing1]。つまり

[^boxing1]: 念のために解説すると「ボックス化」とは，あるインスタンスを型と値をセットにして（大抵はヒープ上の）特定領域に格納することを言う。言い方を変えるとボックス化インスタンスは内部属性としてインスタンスの型と値を持っているわけだ。スマートポインタや依存の注入などは，このボックス化の仕組みと密接な関係がある。

```go
if err != nil {
    ...
}
```

は「err の中にボックス化されたエラー・インスタンスが入っているか」という評価と言えるだろう。

interface 型と nil の関係については以下の拙文を参照のこと。

https://zenn.dev/spiegel/articles/20201010-ni-is-not-nil

## インスタンスの同値性（equality）

次によく見るのは

```go
if err != io.EOF {
    ...
}
```

のようなパターン。 [io].EOF は [io] 標準パッケージにおいて

```go:io/io.go
// EOF is the error returned by Read when no more input is available.
// Functions should return EOF only to signal a graceful end of input.
// If the EOF occurs unexpectedly in a structured data stream,
// the appropriate error is either ErrUnexpectedEOF or some other error
// giving more detail.
var EOF = errors.New("EOF")
```

などと定義されていて，ストリームの終端を示す EOF エラーとして広く使われている。なので，エラー・インスタンスが [io].EOF と同値[^equality1] であれば EOF エラーであると評価できるわけだ。

[^equality1]: IT 用語としての “equality” は日本語では何故か「等価性」と訳されることが多いが，等価性ならむしろ “equivalency” だよなぁ。というわけで，この辺の「用語」は混乱していて宗教論争に発展することも多い。私はそういうものに巻き込まれたくないので，この本では「equality ＝ 同値性」と定義する。ちなみに [Go] では `==` や `!=` は「値」の比較しかしない。ポインタ値の比較で同じ値であれば結果的に2つのインスタンスは「同一」であると見なせるが，やっていることはあくまで「値」の比較である。この辺も [Go] ならではのシンプルさと言えよう。

このように，あらかじめエラー・インスタンスを定義しておいて，それらとの比較を行うことで簡単にエラーの評価を行うことができる。

なお [Go] 1.13 からは [errors].Is() 関数が正式に用意されていて，上のコードは

```go
if !errors.Is(err, io.EOF) {
    ...
}
```

と置き換えることができる。むしろ今後は [errors].Is() 関数を使うことを強くお勧めする（理由は[次節](./layered-error)にて）。

## インスタンスのボックス化解除（unboxing）

たとえば [os].Open() 関数の返り値のエラー型は以下の内部状態を持っている。

```go
// PathError records an error and the operation and file path that caused it.
type PathError struct {
    Op   string
    Path string
    Err  error
}
```

しかし error 型でボックス化している状態では Error() メソッドしか使えないため [os].PathError 型の要素を使うことが出来ない。使うためにはボックス化の解除が必要である。 [Go] では[型アサーション][type assertion]を使ってボックス化解除ができる。

こんな感じ。

```go
switch e := err.(type) {
case *os.PathError:
    if errno, ok := e.Err.(syscall.Errno); ok {
        switch errno {
        case syscall.ENOENT:
            fmt.Fprintln(os.Stderr, "ファイルが存在しない")
        case syscall.ENOTDIR:
            fmt.Fprintln(os.Stderr, "ディレクトリが存在しない")
        default:
            fmt.Fprintln(os.Stderr, "Errno =", errno)
        }
    } else {
        fmt.Fprintln(os.Stderr, "その他の PathError")
    }
default:
    fmt.Fprintln(os.Stderr, "その他のエラー")
}
```

[Go] 1.13 からは [errors].As() 関数が正式に用意され，ボックス化解除が少し楽になった。

```go
var perr *os.PathError
if errors.As(err, &perr) {
    fmt.Fprintf(os.Stderr, "file is \"%v\"\n", perr.Path)
}
```

## Error() メソッドの返り値（文字列）を解析する

[バッドノウハウ](http://0xcc.net/misc/bad-knowhow.html "バッドノウハウと「奥が深い症候群」")。

と切り捨てたいところだが，これまで述べた評価方法が使えない場合は Error() メソッドの返り値である文字列を解析するしかない。特に [fmt].Errorf() 関数でエラー・インスタンスを作ると他の評価手段が封じられてしまう。

なお [fmt].Errorf() 関数については [errors].Is() や [errors].As() などと組み合わせてもう少し構造的に評価できるようになった。これについては[次節](./layered-error)で改めて紹介する。

[Go]: https://go.dev/ "The Go Programming Language"
[io]: https://pkg.go.dev/io/ "io - The Go Programming Language"
[errors]: https://pkg.go.dev/errors/ "errors - The Go Programming Language"
[os]: https://pkg.go.dev/os/ "os - The Go Programming Language"
[fmt]: https://pkg.go.dev/fmt/ "fmt - The Go Programming Language"
[conversion]: https://go.dev/ref/spec#Conversions "The Go Programming Language Specification - The Go Programming Language"
[type assertion]: https://go.dev/ref/spec#Type_assertions "The Go Programming Language Specification - The Go Programming Language"
<!-- eof -->
