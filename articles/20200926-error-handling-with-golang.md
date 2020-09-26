---
title: "errors.Is, errors.As は（単なる）比較関数ではない" # 記事のタイトル
emoji: "⌨" # アイキャッチとして使われる絵文字（1文字だけ）
type: "tech" # "tech" : 技術記事 / "idea" : アイデア記事
topics: ["go", "programming"] # タグ。["markdown", "rust", "aws"] のように指定する
published: true # 公開設定（true で公開）
---

- [［Golang］errors.Is() errors.As() 完全ガイド〜使い方と違いをしっかり調査しました〜 | Zenn](https://zenn.dev/kskumgk63/articles/550dc9d42078d968beac)

という記事を見かけたが微妙に「？？？」な印象だったので，私なりに書き直してみる。

## [Go] におけるエラー・ハンドリングの戦術

[Go] におけるエラー・ハンドリングの戦術は大まかに以下の3つのいずれか，または組み合わせである。

1. `error` インスタンス同士の同値性[^eq1] を調べる（ポインタ値を含む）
2. `error` インスタンスから具体的な型で括り出す
3. `error.Error()` メソッドで出力される文字列を解釈する

[^eq1]: 等値とか等価とかの言葉を使うと絶対に混乱が起きるし，この手の宗教論争に巻き込まれるのは御免なので，インスタンスの値またはポインタ値が単純に同じという意味で「同値」とした。

まぁ，3番めは[バッドノウハウ](http://0xcc.net/misc/bad-knowhow.html "バッドノウハウと「奥が深い症候群」")なので華麗にスルーするとして，1番目に相当するのが [`errors.Is()`]，2番目に相当するのが [`errors.As()`] の各関数である。

## 昔々...

[`errors.Is()`]，[`errors.As()`] 各関数がなかった頃はどうしてたか。

たとえば，以下のようなファイルをオープンしてみるだけの関数があったとする。


```go
func checkFileOpen(path string) error {
    file, err := os.Open(path)
    if err != nil {
        return err
    }
    defer file.Close()
    return nil
}
```

この関数の返り値は以下のように評価できる。

```go
func main() {
    if err := checkFileOpen("not-exist.txt"); err != nil {
        switch e := err.(type) {
        case *os.PathError:
            if errno, ok := e.Err.(syscall.Errno); ok {
                switch errno {
                case syscall.ENOENT:
                    fmt.Fprintf(os.Stderr, "%v ファイルが存在しない\n", e.Path)
                default:
                    fmt.Fprintln(os.Stderr, "Errno =", errno)
                }
            } else {
                fmt.Fprintln(os.Stderr, "その他の PathError")
            }
        default:
            fmt.Fprintln(os.Stderr, "その他のエラー")
        }
        return
    }
    fmt.Println("正常終了")
}
```

まず，返ってきた `error` インスタンスから [`*os.PathError`] 型で括り出し，更にその属性 `Err` を [`syscall.Errno`] 型で括りだしている。その上で [`syscall.Errno`] 型の値を定義済みインスタンスと比較してエラーを判定しているのだ。

このように [Go] でエラーハンドリングを行う際はエラーの内部構造をあらかじめ知っておく必要があるため，どうしても煩雑になる。

## 改訂版エラーハンドリング

上の評価を [`errors.Is()`]，[`errors.As()`] 各関数を使って行うとこんな感じにできる。

```go
func main() {
    if err := checkFileOpen("not-exist.txt"); err != nil {
        var errPath *os.PathError
        if errors.As(err, &errPath) {
            switch {
            case errors.Is(errPath.Err, syscall.ENOENT):
                fmt.Fprintf(os.Stderr, "%v ファイルが存在しない\n", errPath.Path)
            default:
                fmt.Fprintln(os.Stderr, "その他の PathError")
            }
        } else {
            fmt.Fprintln(os.Stderr, "その他のエラー")
        }
        return
    }
    fmt.Println("正常終了")
}
```

もっと言えば [`syscall.Errno`] 型の値を定義済みインスタンスと比較するだけでいいのなら

```go
func main() {
    if err := checkFileOpen("not-exist.txt"); err != nil {
        switch {
        case errors.Is(err, syscall.ENOENT):
            fmt.Fprintln(os.Stderr, "ファイルが存在しない")
        default:
            fmt.Fprintln(os.Stderr, "その他のエラー")
        }
        return
    }
    fmt.Println("正常終了")
}
```

で済む。カンタン！

### `Unwrap()` メソッドで垂直方向にエラーを構造化する

[Go] 1.13 からは，エラーハンドリングにおいて `Unwrap()` メソッドの有無が考慮される。

```go
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

これによって標準パッケージだけで垂直方向の構造化エラーを扱えるようになった。たとえば，先程の [`*os.PathError`] 型であれば

```go
func (e *PathError) Unwrap() error { return e.Err }
```

と定義されていて `Unwrap()` メソッドで原因エラーを返すようになっている。

これにより，内部構造をすっ飛ばして

```go
if errors.Is(err, syscall.ENOENT) {
    ...
}
```

のように原因エラーを直接評価できるのである。

## `errors.As()` は恥だが役に立つ

[`errors.Is()`] 関数 はともかく [`errors.As()`] 関数はちょっと... いや，だいぶカッコ悪い。

```go
var errPath *os.PathError
if errors.As(err, &errPath) {
    ...
}
```

本来なら変換した型を返り値として返すべきなのに引数としてポインタ渡ししている。C言語か（笑）

実は，元々の proposal では [`errors.As()`] 関数は総称型（generics）の実装と抱き合わせだったのだ。たとえばこんな感じ。

```go
func As(type E)(err error) (e E, ok bool) {
    for {
        if e, ok := err.(E); ok {
            return e, true
        }
        err = Unwrap(err)
        if err == nil {
            return nil, false
        }
    }
}
```

でも総称型が投入されるのはしばらく先だし，とりあえず現行の仕様の範囲で実装するとあんな感じにダサくなってしまうのだ。まぁ，内部構造を気にせず指定した型で括り出せるのは便利だし，総称型の登場を楽しみにしつつ現状でなんとかやりくりしよう。

## というわけで，宣伝

自作パッケージで使っているエラーハンドリングを切り出して，独立したパッケージとして公開している。そのまま使うなりコピペしてアレンジして使うなり，ご自由にどうぞ。

- [Go 言語用エラーハンドリング・パッケージ — リリース情報 | text.Baldanders.info](https://text.baldanders.info/release/errs-package-for-golang/)

## その他，参考

- [エラー・ハンドリングについて（追記あり） — プログラミング言語 Go | text.Baldanders.info](https://text.baldanders.info/golang/error-handling/)
- [Go 1.13 のエラー・ハンドリング — プログラミング言語 Go | text.Baldanders.info](https://text.baldanders.info/golang/error-handling-in-go-1_3/)
- [構造化エラーをログ出力する — プログラミング言語 Go | text.Baldanders.info](https://text.baldanders.info/golang/logging-error/)


[Go]: https://golang.org/ "The Go Programming Language"
[errors]: https://pkg.go.dev/errors "errors package · go.dev"
[`errors.Is()`]: https://pkg.go.dev/errors#Is
[`errors.As()`]: https://pkg.go.dev/errors#As
[`*os.PathError`]: https://pkg.go.dev/os#PathError
[`syscall.Errno`]: https://pkg.go.dev/syscall#Errno
