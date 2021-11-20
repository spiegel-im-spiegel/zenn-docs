---
title: "Interface 型をあらかじめ宣言しなくてもよい" # 記事のタイトル
emoji: "🤔" # アイキャッチとして使われる絵文字（1文字だけ）
type: "idea" # "tech" : 技術記事 / "idea" : アイデア記事
topics: ["go", "programming"] # タグ。["markdown", "rust", "aws"] のように指定する
published: true # 公開設定（true で公開）
---

いつもの小ネタです。起点は以下の tweet から。

https://twitter.com/mattn_jp/status/1461887274744905728

かいつまんで説明すると，[元々の tweet](https://twitter.com/techno_tanoC/status/1461640024253153282) に

> golang、interface Aとinterface Bを満たすものを引数として受け取れる関数を表現するのにinterface ABを宣言しないといけないの？
> 
> rustならtrait使ってT: A +Bでいけるのに。

とあって，それに対して

```go
type A interface {
    DoSomething()
}

type B interface {
    DoAnotherthing()
}

func Do(v interface {A; B}) {
    v.DoSomething()
    v.DoAnotherthing()
}
```

てな感じに書けるよ，という話。もっとも，上の `Do()` 関数を gofmt にかけると

```go
func Do(v interface {
    A
    B
}) {
    v.DoSomething()
    v.DoAnotherthing()
}
```

と整形されてしまうけど（笑）

実はこれ「抽象」と「具象」の間に **継承関係はない** という [Go] のとても重要な機能なの。なので，上の `Do()` 関数のように（仮引数 `v` の実体が何であるかに関係なく）欲しい振る舞いを示す interface 型を即席で作って **制約を課す** ことができる。

たとえば [errors] 標準パッケージに [errors].Unwrap() 関数があるが，これは以下のように実装されている。

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

わざわざ

```go
type Unwrapper interface {
    Unwrap() error
}
```

みたいな interface 型をあらかじめ宣言しなくても，これで必要十分な機能を提供できる。同様に [errors].Is() 関数も

```go:errors/wrap.go
// Is reports whether any error in err's chain matches target.
//
// The chain consists of err itself followed by the sequence of errors obtained by
// repeatedly calling Unwrap.
//
// An error is considered to match a target if it is equal to that target or if
// it implements a method Is(error) bool such that Is(target) returns true.
//
// An error type might provide an Is method so it can be treated as equivalent
// to an existing error. For example, if MyError defines
//
//    func (m MyError) Is(target error) bool { return target == fs.ErrExist }
//
// then Is(MyError{}, fs.ErrExist) returns true. See syscall.Errno.Is for
// an example in the standard library.
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
        // TODO: consider supporting target.Is(err). This would allow
        // user-definable predicates, but also may allow for coping with sloppy
        // APIs, thereby making it easier to get away with them.
        if err = Unwrap(err); err == nil {
            return false
        }
    }
}
```

と書かれている。私はこれを見て目から鱗が落ちた。

私もそうだったが， C++ や Java や Rust のような公称型の部分型付け（nominal subtyping）に慣れていると何となく「抽象型を宣言しなくちゃ」と思ってしまうが， [Go] の場合は抽象型をあらかじめ宣言する必要は微塵もない[^ss1]。むしろ，最初に interface 型を乱発するのは（抽象型に具象型を合わせようという強制力が働くため）開発プロセスの妨げになることさえある。

抽象型で具象型を「囲う」のではなく，必要に応じて最小限の範囲で「接続する」イメージで考えるのがいいのではないだろうか。具象から抽象へ思考（指向）するのが [Go] 流だと思う。

[^ss1]: [Go] のような型付けシステムを「構造型の部分型付け（structural subtyping）」と呼ぶそうな。

[Go] の言語上のメリットのひとつは「継承」という軛（くびき）から自由である，という点だろう。これを実感できるようになれば C++ や Java 上がりのプログラマでももっと自由に [Go] のコードを書けると思う。

[Go]: https://golang.org/ "The Go Programming Language"
[errors]: https://pkg.go.dev/errors "errors package - errors - pkg.go.dev"
<!-- eof -->
