---
title: "Go 標準のソート・アルゴリズム" # 記事のタイトル
emoji: "⌨" # アイキャッチとして使われる絵文字（1文字だけ）
type: "tech" # "tech" : 技術記事 / "idea" : アイデア記事
topics: ["go", "programming", "sort", "algorithm"] # タグ。["markdown", "rust", "aws"] のように指定する
published: true # 公開設定（true で公開）
---

今回は

https://zenn.dev/satoru_takeuchi/articles/d66b8b69e0218d9f137e

から着想を得て [Go] 標準の [sort] パッケージについて簡単に紹介する。なお今回の記事は，以前 LT 用に書いたスライド「[配列とスライスとソート](https://slide.baldanders.info/shimane-go-2020-02-13/)」からの抜粋となっている。よろしければスライドの方もご笑覧あれ。

## まずはコード

[Go] 標準の [sort] パッケージでは任意の slice ソートについて以下のように記述できる。

```go
package main

import (
    "fmt"
    "sort"
)

func main() {
    ds := []float64{0.055, 0.815, 1.0, 0.107}
    fmt.Println(ds) //before
    sort.Slice(ds, func(i, j int) bool {
        return ds[i] < ds[j]
    })
    fmt.Println(ds) //after
}
```

このコードの[実行結果](https://play.golang.org/p/GbV5WEu5Mic)は以下の通り。

```
[0.055 0.815 1 0.107]
[0.055 0.107 0.815 1]
```

[sort] パッケージでは slice に対して2種類のソート関数を用意していて

- [`func Slice(slice interface{}, less func(i, j int) bool)`](https://pkg.go.dev/sort#Slice)
- [`func SliceStable(slice interface{}, less func(i, j int) bool)`](https://pkg.go.dev/sort#SliceStable)

このうち `sort.SliceStable()` 関数のほうを安定ソート（[挿入ソート]）として記述がされている。

## [Sort][sort] パッケージで使われているアルゴリズム

[Sort][sort] パッケージの  `sort.Slice()` 関数では以下のように複数のアルゴリズムを使い分けている。

1. 基本は[クイックソート]
2. 要素数が12以下なら[シェルソート]（gap 6）へ
3. 要素数が6以下なら[挿入ソート]へ
4. [クイックソート]の再帰レベルが一定を超えたら[ヒープソート]へ（[イントロソート]）
    - $depth = \lceil \mathrm{lb}(n+1)\rceil\times 2$

ちなみにアルゴリズム毎の計算量は以下のとおりである。

|             名前 | 最良                     | 最悪                   | 平均                   |
| ----------------:| ------------------------ | ---------------------- | ---------------------- |
| [クイックソート] | $\mathrm{O}(n \log n)$   | $\mathrm{O}(n^2)$      | $\mathrm{O}(n \log n)$ |
|   [ヒープソート] | $\mathrm{O}(n \log n)$   | $\mathrm{O}(n \log n)$ | $\mathrm{O}(n \log n)$ |
|   [シェルソート] | $\mathrm{O}(n \log^2 n)$ | $\mathrm{O}(n^2)$      | &mdash;                |
|     [挿入ソート] | $\mathrm{O}(n)$          | $\mathrm{O}(n^2)$      | $\mathrm{O}(n^2)$      |

[Go]: https://golang.org/ "The Go Programming Language"
[sort]: https://pkg.go.dev/sort "sort package · go.dev"
[クイックソート]: https://en.wikipedia.org/wiki/Quicksort
[ヒープソート]: https://en.wikipedia.org/wiki/Heapsort
[シェルソート]: https://en.wikipedia.org/wiki/Shellsort
[イントロソート]: https://en.wikipedia.org/wiki/Introsort
[挿入ソート]: https://en.wikipedia.org/wiki/Insertion_sort
