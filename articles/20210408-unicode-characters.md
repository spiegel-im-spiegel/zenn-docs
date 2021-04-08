---
title: "Unicode 文字列を「文字」単位に分離する" # 記事のタイトル
emoji: "💻" # アイキャッチとして使われる絵文字（1文字だけ）
type: "tech" # "tech" : 技術記事 / "idea" : アイデア記事
topics: ["go", "programming", "unicode"] # タグ。["markdown", "rust", "aws"] のように指定する
published: true # 公開設定（true で公開）
---

前に「[やっかいな日本語](https://zenn.dev/spiegel/articles/20210118-characters)」でも紹介したが Unicode 文字列は「1コードポイント＝1文字」ではない。特にやっかいなのが絵文字で，このあたりの話は自ブログでまとめている。

https://text.baldanders.info/remark/2021/03/terrible-emoji/
https://text.baldanders.info/remark/2021/04/emoji-list/

この記事の中でさらっと紹介しているが， [github.com/rivo/uniseg][rivo/uniseg] という [Go] 言語用パッケージがあって，これを使うと UTF-8 文字列を「文字」単位に切り出してくれるらしい。

早速 [github.com/rivo/uniseg][rivo/uniseg] のサンプルコードを（少しだけアレンジして）動かしてみよう。

```go:sample1.go
// +build run

package main

import (
    "fmt"

    "github.com/rivo/uniseg"
)

func main() {
    text := "👍🏼!"
    fmt.Println("Text:", text)
    gr := uniseg.NewGraphemes(text)
    for gr.Next() {
        rs := gr.Runes()
        fmt.Printf("%v : %U\n", string(rs), rs)
    }
}
```

これを実行すると

```
$ go run sample1.go
Text: 👍🏼!
👍🏼 : [U+1F44D U+1F3FC]
! : [U+0021]
```

となった。じゃあ入力テキストを

```go:sample2.go
text := "ペンギン ﾍﾟﾝｷﾞﾝ"
```

に変えて試してみようか。

```
$ go run sample2.go
Text: ペンギン ﾍﾟﾝｷﾞﾝ
ペ : [U+30D8 U+309A]
ン : [U+30F3]
ギ : [U+30AD U+3099]
ン : [U+30F3]
  : [U+0020]
ﾍﾟ : [U+FF8D U+FF9F]
ﾝ : [U+FF9D]
ｷﾞ : [U+FF77 U+FF9E]
ﾝ : [U+FF9D]
```

ほほう。濁点や半濁点の結合文字もちゃんと認識してくれるんだねぇ。偉い偉い。

ではでは，次は色んなパターンの絵文字

```go:sample3.go
text := "|#️⃣|☝️|☝🏻|🇯🇵|🏴󠁧󠁢󠁥󠁮󠁧󠁿|👩🏻‍❤️‍💋‍👨🏼|"
```

で試してみよう。

```
$ go run sample3.go
Text: |#️⃣|☝️|☝🏻|🇯🇵|🏴󠁧󠁢󠁥󠁮󠁧󠁿|👩🏻‍❤️‍💋‍👨🏼|
| : [U+007C]
#️⃣ : [U+0023 U+FE0F U+20E3]
| : [U+007C]
☝️ : [U+261D U+FE0F]
| : [U+007C]
☝🏻 : [U+261D U+1F3FB]
| : [U+007C]
🇯🇵 : [U+1F1EF U+1F1F5]
| : [U+007C]
🏴󠁧󠁢󠁥󠁮󠁧󠁿 : [U+1F3F4 U+E0067 U+E0062 U+E0065 U+E006E U+E0067 U+E007F]
| : [U+007C]
👩🏻‍❤️‍💋‍👨🏼 : [U+1F469 U+1F3FB U+200D U+2764 U+FE0F U+200D U+1F48B U+200D U+1F468 U+1F3FC]
| : [U+007C]
```

おおっ！ きれいに分離できた。ちなみに各絵文字は

| 絵文字 | シーケンス・タイプ | 名前 |
| :----:| ---- | ----------------- |
| #️⃣ | emoji keycap sequence | keycap: # |
| ☝️ | emoji presentation sequence | index pointing up |
| ☝🏻 | emoji modifier sequence | index pointing up: light skin tone |
| 🇯🇵 | emoji flag sequence | flag: Japan |
| 🏴󠁧󠁢󠁥󠁮󠁧󠁿 | emoji tag sequence | flag: England |
| 👩🏻‍❤️‍💋‍👨🏼 | emoji zwj sequence | kiss: woman, man, light skin tone, medium-light skin tone |

という感じに分類できる（下に行くほど表示できるプラットフォームが限られてくるので悪しからず）。最後のなんか

| 絵文字 | コードポイント | 名前 |
| :----:| ------------- | ---- |
| 👩🏻 | U+1F469 U+1F3FB | woman: light skin ton |
| ❤️ | U+2764 U+FE0F | red heart |
| 💋 | U+1F48B | KISS MARK |
| 👨🏼 | U+1F468 U+1F3FC | man: medium-light skin tone |

の4つの文字を ZWJ (U+200D) で繋げてひとつの絵文字 👩🏻‍❤️‍💋‍👨🏼 とする（全部で10個のコード列）鬼畜仕様である。

でも，まぁ，これで絵文字を含めて Unicode 文字列を「文字」単位に分離できることが確認できた。めでたし。

[Go]: https://golang.org/ "The Go Programming Language"
[rivo/uniseg]: https://github.com/rivo/uniseg "rivo/uniseg: Unicode Text Segmentation for Go (or: How to Count Characters in a String)"
<!-- eof -->
