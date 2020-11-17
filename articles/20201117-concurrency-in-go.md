---
title: "Go における「並行処理」の抽象化" # 記事のタイトル
emoji: "🤔" # アイキャッチとして使われる絵文字（1文字だけ）
type: "idea" # "tech" : 技術記事 / "idea" : アイデア記事
topics: ["golang", "concurrency"] # タグ。["markdown", "rust", "aws"] のように指定する
published: true # 公開設定（true で公開）
---

面白い記事を見かける。

https://zenn.dev/it/articles/e975f08392ea846d9d7b

プログラミング言語ごとの差異が分かりやすく紹介されていて，特に [Rust] における並列処理の最近のトレンドは知らなかったのでとても参考になった。ただ [Go] に関してはもう少し説明が必要な気がするので補足記事を書いてみる。

## [Go] における「並行処理」と「並列処理」

まずは「並行処理」と「並列処理」の違いについて書こうと思ったが，以下の素晴らしい記事にて言及しつくされていた orz

https://zenn.dev/koron/articles/3ddcaaeae37f9befdf70

...気を取り直して， [Go] の「並行処理」に関してはバイブルと言える本がある。

https://www.amazon.co.jp/dp/4873118468

この本の中で「並行処理」と「並列処理」の違いについて以下のように書かれている（2.1章）。

> 並行性はコードの性質を指し，並列性は動作しているプログラムの性質を指します。

何故このように言い切れるのか。それはこの一文に集約されるだろう。

> 並行性と並列性の違いはコードの設計をする際に非常に強力な抽象化になることがわかり，そして Go はこの違いを最大限に活かしています。

## Coroutine としての Goroutine

抽象化された並行処理として思い浮かぶキーワードは coroutine だろう。 Goroutine は coroutine の一種と考えられる。ちなみに [Wikipedia で coroutine を引く](https://en.wikipedia.org/wiki/Coroutine "Coroutine - Wikipedia")と

> Coroutines are [computer program](https://en.wikipedia.org/wiki/Computer_program) components that generalize [subroutines](https://en.wikipedia.org/wiki/Subroutine) for [non-preemptive multitasking](https://en.wikipedia.org/wiki/Non-preemptive_multitasking), by allowing execution to be suspended and resumed. Coroutines are well-suited for implementing familiar program components such as [cooperative tasks](https://en.wikipedia.org/wiki/Cooperative_multitasking), [exceptions](https://en.wikipedia.org/wiki/Exception_handling), [event loops](https://en.wikipedia.org/wiki/Event_loop), [iterators](https://en.wikipedia.org/wiki/Iterator), [infinite lists](https://en.wikipedia.org/wiki/Lazy_evaluation) and [pipes](https://en.wikipedia.org/wiki/Pipeline_(software)). 

などと書かれている。ここで coroutine は非プリエンプティブであると書かれているが， [Go] の goroutine は [1.14](https://golang.org/doc/go1.14 "Go 1.14 Release Notes - The Go Programming Language") から（一部のプラットフォームを除いて）プリエンプティブな非同期処理を獲得している。

> Goroutines are now asynchronously preemptible. As a result, loops without function calls no longer potentially deadlock the scheduler or significantly delay garbage collection. This is supported on all platforms except `windows/arm`, `darwin/arm`, `js/wasm`, and `plan9/*`.

これによってプログラマは「並列処理」の実装詳細を考えることなく「並行処理」のコード化に専念できるわけだ。

## 並行処理のトレードオフ

私の個人的な好みとしては [Go] と [Rust] が好きなので[^rst1] 両者の比較を[よく考える](https://text.baldanders.info/remark/2020/04/subtyping/ "それは Duck Typing ぢゃない（らしい） | text.Baldanders.info")が，最初に紹介した「[多言語からみるマルチコアの活かし方](https://zenn.dev/it/articles/e975f08392ea846d9d7b)」を見るに，並行処理に関しても両者は対象的なんだなぁ，と感じる。

[^rst1]: [Rust] については今のところ停滞しているが。個人が余暇でちょっとしたコードを書くには [Rust] はちょっとヘヴィなんだよなぁ。

大抵のプログラミング言語は「並列処理」をライブラリまたはフレームワークの一部として装備している。なので大方において「並行処理」と「並列処理」は未分化のまま設計せざるを得ない。だから（並列処理に対する）並行処理を「複数の処理を順番に実行すること」みたいな勘違いが発生するのだと思う。しかし見方を変えると「並列処理」と「並行処理」が密結合しているからこそ得られるパフォーマンスもあるわけだ， [Rust] のように。

一方 [Go] は「並列処理」をランタイム内に標準装備して「並行処理」と明確に分離し，両者を疎結合とすることで「並行処理」そのものの自由度を高めている。その意味で [Go] と [Rust] はトレードオフのような関係になっているわけだ。まぁ「[Go] か [Rust] かどっちか選べ」みたいな究極の選択はないだろうが（笑）

## というわけで...

[Go] で並行処理を勉強するなら，まずは『[Go言語による並行処理]』を読みなはれ，ということで。

## 参考リンク

- [Go: Goroutine, OS Thread and CPU Management | by Vincent Blanchon | A Journey With Go | Medium](https://medium.com/a-journey-with-go/go-goroutine-os-thread-and-cpu-management-2f5a5eaf518a)
- [Go: GOMAXPROCS & Live Updates. ℹ️ This article is based on Go 1.13. | by Vincent Blanchon | A Journey With Go | Medium](https://medium.com/a-journey-with-go/go-gomaxprocs-live-updates-407ad08624e1)
- [In-process caching in Go: scaling lakeFS to 100k requests/second](https://lakefs.io/2020/09/23/in-process-caching-in-go-scaling-lakefs-to-100k-requests-second/)
- [goroutineはなぜ軽量なのか - Carpe Diem](https://christina04.hatenablog.com/entry/why-goroutine-is-good)
- [『Go 言語による並行処理』は Go 言語プログラマ必読書だろう | text.Baldanders.info](https://text.baldanders.info/remark/2018/11/concurrency-in-go/)

## 参考図書

https://www.amazon.co.jp/dp/4621300253

[Go]: https://golang.org/ "The Go Programming Language"
[Go言語による並行処理]: https://www.amazon.co.jp/dp/4873118468?tag=baldandersinf-22&linkCode=ogi&th=1&psc=1 "Go言語による並行処理 | Katherine Cox-Buday, 山口 能迪 |本 | 通販 | Amazon"
[Rust]: https://www.rust-lang.org/ "Rust Programming Language"
<!-- eof -->
