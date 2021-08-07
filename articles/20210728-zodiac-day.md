---
title: "土用の丑の日なので..." # 記事のタイトル
emoji: "💻" # アイキャッチとして使われる絵文字（1文字だけ）
type: "tech" # "tech" : 技術記事 / "idea" : アイデア記事
topics: ["go", "programming"] # タグ。["markdown", "rust", "aws"] のように指定する
published: true # 公開設定（true で公開）
---

たまにはこっちでも書かないと（笑） 軽く小ネタです。

2021-07-28 は「土用の丑の日」でした。現行暦では「土用の入り」が雑節のひとつとして定義されていて「太陽黄経が 297°, 27°, 117°, 207° となる日」と決められている。また「土用の明け」は 立春（315°），立夏（45°），立秋（135°），立冬（225°）の前日となる。

したがって，この期間内で「丑の日」を探せばよい。たとえば 2001-01-01 は甲子つまり「子の日」なので，この日を基準に 子，丑，寅，... と数え上げればいいわけだ。

では早速，土用の入り（2021-07-19）から立秋（2021-08-07）の前日までの期間で調べてみる。

```go:sample.go
// +build run

package main

import (
    "fmt"
    "time"
)

var (
    jst         = time.FixedZone("Asia/Tokyo", int((9 * time.Hour).Seconds()))
    zodiacNames = []string{"子", "丑", "寅", "卯", "辰", "巳", "午", "未", "申", "酉", "戌", "亥"}
    baseDay     = time.Date(2001, time.January, 1, 0, 0, 0, 0, jst)
)

func zodiacName(t time.Time) string {
    d := int64(t.Sub(baseDay).Hours()) / 24 % 12
    if d < 0 {
        d += 12
    }
    return zodiacNames[d]
}

func main() {
    day := time.Date(2021, time.July, 18, 0, 0, 0, 0, jst)
    for i := 0; i < 19; i++ {
        day = day.Add(time.Hour * 24)
        fmt.Printf("%v is %v\n", day.Format("2006-01-02"), zodiacName(day))
    }
}
```

やっつけコードでゴメンペコン。実行結果は以下の通り。

```
$ go run sample.go 
2021-07-19 is 辰
2021-07-20 is 巳
2021-07-21 is 午
2021-07-22 is 未
2021-07-23 is 申
2021-07-24 is 酉
2021-07-25 is 戌
2021-07-26 is 亥
2021-07-27 is 子
2021-07-28 is 丑
2021-07-29 is 寅
2021-07-30 is 卯
2021-07-31 is 辰
2021-08-01 is 巳
2021-08-02 is 午
2021-08-03 is 未
2021-08-04 is 申
2021-08-05 is 酉
2021-08-06 is 戌
```

というわけで 2021-07-28 が丑の日であることが分かる。

どっとはらい。

## 【2021-08-07 追記】

調子に乗ってパッケージ化してみた。

https://github.com/spiegel-im-spiegel/jzodiac

[Go]: https://golang.org/ "The Go Programming Language"
<!-- eof -->
