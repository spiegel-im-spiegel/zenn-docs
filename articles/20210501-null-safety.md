---
title: "Null の始末は「誰」がするのか" # 記事のタイトル
emoji: "🤔" # アイキャッチとして使われる絵文字（1文字だけ）
type: "idea" # "tech" : 技術記事 / "idea" : アイデア記事
topics: ["go", "programming"] # タグ。["markdown", "rust", "aws"] のように指定する
published: true # 公開設定（true で公開）
---

みなさま GW はいかがお過ごしでしょうか。どこぞの為政者どもが下手こいた所為で自粛や制限を余儀なくされてる状態なので，私は引きこもって[ガンブラ組んで](https://www.flickr.com/photos/spiegel/51147596917/ "GWの工作 | Flickr")遊んでたりします。

ところで

https://zenn.dev/zetamatta/scraps/07cb6835e172a6

の後半で C# の [null 許容参照型][nullable-reference-types]について言及されている[^nrt]。

[^nrt]: C# の [null 許容参照型][nullable-reference-types]は2019年にリリースされた 8.0 から導入されたらしい。ただし null の静的検査および警告はオプションで，既定では無効になっている。

そういえば，縁あって昨年末に職業プログラマに復帰できたのだが，ン年ぶりに書いた Java コードで null 参照がいわゆる「ヌルポ」になることを失念していて思い出すまでしばらく悩んでいたよ[^dbg1]。

[^dbg1]: 仕事で [Go] コードを書いたことがないというのもあるのだが [Go] のデバッガを使ったことがない。必要なかったし。んで Java コードの null 参照のバグに気がつくまでこれまた久しぶりにデバッガを使わざるをえなかった。組み込み開発だとオシロスコープや ICE + デバッガは必須の道具だけど，ただのアプリケーションでデバッガを使わざるをえない状況になるってのは設計がヘタレだからだよなぁ，と自己嫌悪したり。まぁ，他人が書いたコードは把握できなくてもしょうがないのでデバッガのお世話になることも多いのだが。

## Null 安全（Null Safety）

仕様として値の「参照」ができるプログラミング言語では常に null 参照の問題がつきまとう。中には [null 参照による損失を10億ドルと見積もっている人もいる](https://en.wikipedia.org/wiki/Tony_Hoare "Tony Hoare - Wikipedia")。

ある変数が「null 安全」かどうかをコンパイラで制御するというのはメンタル・モデルの観点からも合理的である。 C# の [null 許容参照型][nullable-reference-types]もそうした制御のひとつと言える。「null 安全」については5年前にブログ記事を書いたので参考にどうぞ。

https://text.baldanders.info/remark/2016/11/null-safety/

Swift, Kotlin, Dart や上述の C# のように「原則として null 参照を許容しない（許容する際は明示的に型宣言する）」のが最も素直な戦略だが，他の戦略をとるプログラム言語もある。

## [Rust] の所有権（Ownership）

たとえば，現時点で最もトレンディなコンパイル言語である [Rust] には「所有権」と呼ばれる仕組みがある。所有権のルールは以下の3つ。

1. [Rust] の各値は、所有者と呼ばれる変数と対応している
2. いかなる時も所有者は一つである
3. 所有者がスコープから外れたら、値は破棄される

また明示的に参照型として宣言される変数を使って（所有権の移動のない）借用を行うこともできる。値を所有または借用していない変数を使おうとするとコンパイル・エラーになるため，これを利用して「null 安全」をある程度達成できる。

## [Go] に「参照」はないが...

個人的な好みで [Go] の話をしてしまうが， [Go] には仕様上の「参照」は存在しない。その代わりに値へのアドレッシングを指すポインタ型が存在する。ただし以下の型は「参照のように振る舞う」ので注意が必要である。

| 型名 | 記述例 |
| ---- | ------ |
| チャネル | `chan int` |
| インタフェース | `interface{}` |
| 関数 | `func(int) int` |
| スライス | `[]int` |
| マップ | `map[int]int` |

ポインタ型および上記の型は初期値として nil 値をとる。これが事実上の null 参照と言えるが，残念ながら **[Go] は「null 安全」とは言えない**。

面白いのは nil 値をとる変数でもメソッドを呼び出し可能な点だ。以下のコードはコンパイル・エラーにも実行時 panic にもならない。

```go
package main

import "fmt"

type Hello struct{}

func (h *Hello) String() string {
    return "Hello, World!"
}

func main() {
    var hello *Hello = nil
    fmt.Println(hello.String())
    // Output:
    // Hello, World!
}
```

もちろん nil 変数のメソッド内で内部状態にアクセスすれば[実行時 panic](https://play.golang.org/p/gBl0Imdobmu) になる。

```go
package main

import "fmt"

type Hello struct {
    Name string
}

func (h *Hello) String() string {
    return "Hello, " + h.Name
}

func main() {
    var hello *Hello = nil
    fmt.Println(hello.String()) // panic: runtime error: invalid memory address or nil pointer dereference
}
```

このように [Go] ではメソッドの設計によって null 参照の始末をメソッド側の責務とすることができる。上手く利用して欲しい。

## 参考

https://mattn.kaoriya.net/software/lang/go/20190516095124.htm

[Go]: https://golang.org/ "The Go Programming Language"
[Rust]: https://www.rust-lang.org/ "Rust Programming Language"
[nullable-reference-types]: https://docs.microsoft.com/ja-jp/dotnet/csharp/language-reference/builtin-types/nullable-reference-types "Null 許容参照型 - C# リファレンス | Microsoft Docs"
