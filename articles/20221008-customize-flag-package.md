---
title: "標準 flag パッケージを pflag パッケージに置き換える" # 記事のタイトル
emoji: "💻" # アイキャッチとして使われる絵文字（1文字だけ）
type: "tech" # "tech" : 技術記事 / "idea" : アイデア記事
topics: ["go", "programming", "test"] # タグ。["markdown", "rust", "aws"] のように指定する
published: true # 公開設定（true で公開）
---

いつもの小ネタです。

## 標準 flag パッケージを pflag パッケージに置き換える

コマンドライン引数を評価する [Go] 標準の [flag] はシンプルながらとてもよく出来ているのだけど [GNU 拡張のシンタックス](https://www.gnu.org/software/libc/manual/html_node/Argument-Syntax.html "Argument Syntax (The GNU C Library)")が使えたらなぁ，と思ったことはありません？

たとえば --foo というパラメータの短縮形として -f が使いたいとか -a -b -c をまとめて -abc と指定したいとか。

これを実装できるのが [github.com/spf13/pflag][pflag] パッケージ。こんな感じに書ける。

```go:sample1.go
package main

import (
    "fmt"

    "github.com/spf13/pflag"
)

func main() {
    f := pflag.BoolP("foo", "f", false, "option foo")
    b := pflag.BoolP("bar", "b", false, "option bar")
    pflag.Parse()

    fmt.Println("foo = ", *f)
    fmt.Println("bar = ", *b)
}
```

これを実行すると

```
$ go run sample1.go
foo =  false
bar =  false

$ go run sample1.go --foo --bar
foo =  true
bar =  true

$ go run sample1.go --foo=true
foo =  true
bar =  false

$ go run sample1.go -fb
foo =  true
bar =  true
```

てな感じにコマンドライン引数を指定できる。

[pflag] は標準 [flag] と互換性があるので


```go:sample2.go
package main

import (
    "fmt"

    flag "github.com/spf13/pflag"
)

func main() {
    f := flag.Bool("foo", false, "option foo")
    b := flag.Bool("bar", false, "option bar")
    flag.Parse()

    fmt.Println("foo = ", *f)
    fmt.Println("bar = ", *b)
}
```

などと書くことで置き換え可能である。ただし挙動は [pflag] の仕様に従うので

```
$ go run sample2.go --foo
foo =  true
bar =  false

$ go run sample2.go -foo
unknown shorthand flag: 'f' in -foo
Usage of /tmp/go-build421334830/b001/exe/sample2:
      --bar   option bar
      --foo   option foo
unknown shorthand flag: 'f' in -foo
```

上のように引数に -foo とかしても「そんなもん知らん」と怒られる（笑）

## 【おまけ】 go test に独自のコマンドライン・パラメータを設定する

恥ずかしながら，[『Go言語による分散サービス』読書会][読書会]で標準 [flag] パッケージを使って go test に独自のコマンドライン・パラメータを設定できることをはじめて知った。こんな感じに書ける。

```go:sample3_test.go
package sample3

import (
    "flag"
    "testing"
)

var foo = flag.Bool("foo", false, "option foo")

func TestMain(m *testing.M) {
    flag.Parse()
    m.Run()
}

func TestFlag(t *testing.T) {
    if !*foo {
        t.Errorf("option foo = %v, want %v.", *foo, true)
    }
}
```

これを実行すると

```
$ go test --shuffle on
-test.shuffle 1665227879001765953
--- FAIL: TestFlag (0.00s)
    sample3_test.go:17: option foo = false, want true.
FAIL
exit status 1
FAIL	pflag-sample/sample3.go	0.001s

$ go test --shuffle on --foo
-test.shuffle 1665227866801533228
PASS
ok  	pflag-sample/sample3.go	0.001s
```

となる。指定したフラグはパッケージ内でのみ有効である点に注意。

まぁ，コマンドライン引数でテスト条件を変えるというのはあまりしないだろうが，『[Go言語による分散サービス](https://www.oreilly.co.jp/books/9784873119977/ "O'Reilly Japan - Go言語による分散サービス")』の6章のサンプルコードでは --debug フラグを設定し，コマンドラインでこれが指定されている場合はトレースログを吐くようにしていて「なるほど」と思った。

これを [pflag] でもできないかなぁ，と思ったのだが上手くいかなかった。残念。 go test では標準 [flag] を使いましょう。

[Go]: https://go.dev/ "The Go Programming Language"
[flag]: https://pkg.go.dev/flag "flag package - flag - Go Packages"
[pflag]: https://github.com/spf13/pflag "spf13/pflag: Drop-in replacement for Go's flag package, implementing POSIX/GNU-style --flags."
[読書会]: https://technical-book-reading-2.connpass.com/event/260183/ "第3回『Go言語による分散サービス』オンライン読書会 - connpass"
