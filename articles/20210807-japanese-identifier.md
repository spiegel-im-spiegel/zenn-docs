---
title: "日本語は公開できない #golang" # 記事のタイトル
emoji: "💻" # アイキャッチとして使われる絵文字（1文字だけ）
type: "tech" # "tech" : 技術記事 / "idea" : アイデア記事
topics: ["go", "programming"] # タグ。["markdown", "rust", "aws"] のように指定する
published: true # 公開設定（true で公開）
---

[第15回『プログラミング言語Go』オンライン読書会](https://gpl-reading.connpass.com/event/218308/) で思い出した小ネタをひとつ。

[前回](https://zenn.dev/spiegel/articles/20210728-zodiac-day "土用の丑の日なので...")書いた記事から調子に乗って十干十二支を数え上げるパッケージを作ってみた。実用性は考えない（笑）

https://github.com/spiegel-im-spiegel/jzodiac
https://text.baldanders.info/release/2021/07/japanese-zodiac/

この中で十干十二支をシンボル化するのに

```go:zodiac.go
type Kan10 uint

const (
    Kinoe     Kan10 = iota // 甲（木の兄）
    Kinoto                 // 乙（木の弟）
    Hinoe                  // 丙（火の兄）
    Hinoto                 // 丁（火の弟）
    Tsutinoe               // 戊（土の兄）
    Tsutinoto              // 己（土の弟）
    Kanoe                  // 庚（金の兄）
    Kanoto                 // 辛（金の弟）
    Mizunoe                // 壬（水の兄）
    Mizunoto               // 癸（水の弟）
    KanMax
)

type Shi12 uint

const (
    Rat     Shi12 = iota // 子
    Ox                   // 丑
    Tiger                // 寅
    Rabbit               // 卯
    Dragon               // 辰
    Snake                // 巳
    Horse                // 午
    Sheep                // 未
    Monkey               // 申
    Rooster              // 酉
    Dog                  // 戌
    Boar                 // 亥
    ShiMax
)
```

という感じに書いたんだけど，本当は

```go
type Kan10 uint

const (
    甲 Kan10 = iota
    乙
    丙
    丁
    戊
    己
    庚
    辛
    壬
    癸
)

const (
    子 Shi12 = iota
    丑
    寅
    卯
    辰
    巳
    午
    未
    申
    酉
    戌
    亥
 )
```

と書きたかったのよ。でも lint で

```
zodiac.go:8:2: `甲` is unused (deadcode)
    甲 Kan10 = iota
    ^
```

という感じに，ものごっつ怒られてしまった。

考えてみたら日本語の漢字や仮名は「大文字ではない」から，これらの文字で始まる識別子はパッケージ外部に公開できないんだよね。ちなみに全角の英字なら大文字があるので公開できる。わざわざ全角にする意味はないけど。

というわけで [Go] の仕様の意外な落とし穴にハマってしまったよ orz

## 参考

https://text.baldanders.info/golang/go-with-japanese/

[Go]: https://golang.org/ "The Go Programming Language"
<!-- eof -->
