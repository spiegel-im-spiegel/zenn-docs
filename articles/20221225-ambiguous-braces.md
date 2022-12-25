---
title: "曖昧なブレス" # 記事のタイトル
emoji: "💻" # アイキャッチとして使われる絵文字（1文字だけ）
type: "tech" # "tech" : 技術記事 / "idea" : アイデア記事
topics: ["go", "programming"] # タグ。["markdown", "rust", "aws"] のように指定する
published: true # 公開設定（true で公開）
---

Twitter の「[プログラミング言語Go](https://twitter.com/i/communities/1498095077222400000)」コミュニティで出てきたネタだが，なかなか面白かったので，覚え書きとして残しておく。

コミュニティ内で出てきたのとはちょっと違うが，起点はこんなコード。

```go:prog.go
package main

import "fmt"

type person struct {
    name string
}

func main() {
    if p := person{name: "alice"}; true {
        fmt.Println("hello", p.name)
    }
}
```

if 文に書かれた条件式が常に `true` な点には目をつぶっていただいて，一見無害そうな[コード](https://go.dev/play/p/eVrqInf6-Vc)ではある。しかし，これを実行しようとすると

```
./prog.go:10:7: syntax error: cannot use p := person as value
./prog.go:11:31: syntax error: unexpected newline in composite literal; possibly missing comma or }
```

とコンパイルエラーになる。私も最初は分からなかったのだが，どうやら `if { ... }` 文のブレスと構造体リテラル `person{ ... }` のブレスとが混濁しているようだ。構文解析におけるこの曖昧さについて，どちらか一方に倒して進めるのではなく，安全に「コンパイルエラー」としてしまうのが [Go] らしい（笑）

実はこの辺のことは言語仕様に明記されていた。

> A parsing ambiguity arises when a composite literal using the TypeName form of the LiteralType appears as an operand between the keyword and the opening brace of the block of an "if", "for", or "switch" statement, and the composite literal is not enclosed in parentheses, square brackets, or curly braces. In this rare case, the opening brace of the literal is erroneously parsed as the one introducing the block of statements. To resolve the ambiguity, the composite literal must appear within parentheses.
>
> ```
> if x == (T{a,b,c}[i]) { … }
> if (x == T{a,b,c}[i]) { … }
> ```
*(via “[The Go Programming Language Specification - The Go Programming Language](https://go.dev/ref/spec#Composite_literals)”)*

というわけで，今回の場合は

```go
if p := (person{name: "alice"}); true { ... }
```

という感じに構造体リテラル記述をカッコ `( ... )` で括ってしまえばいいようだ。

勉強になりました。

[Go]: https://go.dev/ "The Go Programming Language"
