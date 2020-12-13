---
title: "#golang メソッド式とメソッド値" # 記事のタイトル
emoji: "⌨" # アイキャッチとして使われる絵文字（1文字だけ）
type: "tech" # "tech" : 技術記事 / "idea" : アイデア記事
topics: ["go", "programming"] # タグ。["markdown", "rust", "aws"] のように指定する
published: true # 公開設定（true で公開）
---

今回は小ネタ。「[第7回『プログラミング言語Go』オンライン読書会](https://gpl-reading.connpass.com/event/194883/)」で出てきた6.4章の話が面白かったので。

たとえば，型 T に対するメソッド Method() は以下のように表すことができる。

```go
func (t T) Method() { ... }
```

キーワード func 直後にある (t T) は「メソッド・レシーバ」と呼ばれているもので，これが型と関数を繋ぐ役目になっている。で，実はこのメソッドは

```go
func Method(t T) { ... }
```

と等価である。メソッド・レシーバの (t T) が暗黙的な第ゼロ番目の引数になっていると考えれば分かりやすいだろうか。

この性質をよく表しているのが「メソッド式」と「メソッド値」と言える。例として以下のような型とメソッドで考えてみる。

```go
package main

type Number struct {
    n int
}

func (n Number) Add(i int) int {
    return n.n + i
}
```

## メソッド式（Method Expression）

上の定義に対してこんな処理を考えてみる。

```go:sample1.go
func main() {
    add := Number.Add
    fmt.Printf("%T\n", add)
}
```

このコードの実行結果は

```
$ go run sample1.go
func(main.Number, int) int
```

となる。第1引数が Number 型になっているのがお分かりだろうか。このように，型に紐づくメソッドを「メソッド式」に展開するとメソッド・レシーバが第1引数にくる。したがって

```go:sample1b.go
func main() {
    add := Number.Add
    fmt.Println(add(Number{1}, 2))
}
```

このコードの実行結果は

```
$ go run sample1b.go 
3
```

となる。

## メソッド値（Method Value）

さらに面白いのが「メソッド値」だろう。こんな感じのコードを考えてみる。

```go:sample2.go
func main() {
    increment := Number{1}.Add
    fmt.Printf("%T\n", increment)
    fmt.Println(increment(2))
}
```

このコードの実行結果は

```
$ go run sample1.go
func(int) int
3
```

となる。つまり increment にセットされた「メソッド値」には `1` で初期化された Number インスタンスが既に紐付いているのだ。

何かこれってカリー化の部分適用ぽいよね。いや，違うけど（笑）

## 部分適用（partial application）ぽい？

ガチの関数型プログラミング言語 Haskell では，関数定義

```haskell
add x y = x + y
```

はカリー化表現

```haskell
add = \x -> \y -> x + y
```

の糖衣構文となっている[^haskell1]。これを使って部分適用

```haskell
increment = add 1
```

が作れるわけだ。

[^haskell1]: Haskell では関数の引数は1つしかとれないためカリー化は必須の要件となる。

[Go] は第一級関数（first-class function）をサポートしているので，カリー化表現はやろうと思えばできるのだが

```go
package main

import "fmt"

func add(x int) func(int) int {
    return func(y int) int {
        return x + y
    }
}

func main() {
    fmt.Println(add(1)(2)) //Output: 3
    increment := add(1) //partial application
    fmt.Println(increment(2)) //Output: 3
}
```

みたいな感じで余計に面倒くさいコードになってしまう（笑） 単に部分適用がしたいだけならメソッド値を使うほうが簡単かもしれない。

というわけで，どっとはらい。

## リンク

https://text.baldanders.info/remark/2020/03/currying/

## 参考図書

https://www.amazon.co.jp/dp/4621300253

[Go]: https://golang.org/ "The Go Programming Language"
<!-- eof -->
