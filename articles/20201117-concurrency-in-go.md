---
title: "Go における並行処理" # 記事のタイトル
emoji: "🤔" # アイキャッチとして使われる絵文字（1文字だけ）
type: "idea" # "tech" : 技術記事 / "idea" : アイデア記事
topics: ["golang", "concurrency"] # タグ。["markdown", "rust", "aws"] のように指定する
published: false # 公開設定（true で公開）
---

面白い記事を見かける。

https://zenn.dev/it/articles/e975f08392ea846d9d7b

プログラミング言語ごとの差異が分かりやすく紹介されていて，特に Rust における並列処理の最近のトレンドは知らなかったのでとても参考になった。

ただ [Go] に関してはもう少し説明が必要な気がするので補足記事を書いてみる。

## [Go] における「並行処理」と「並列処理」

[Go] の「並行処理」に関してはバイブルと言える本がある。

- [Go言語による並行処理 | Katherine Cox-Buday, 山口 能迪 |本 | 通販 | Amazon][Go言語による並行処理]


- [Coroutine - Wikipedia](https://en.wikipedia.org/wiki/Coroutine)
- [コルーチン - Wikipedia](https://ja.wikipedia.org/wiki/%E3%82%B3%E3%83%AB%E3%83%BC%E3%83%81%E3%83%B3)

> Coroutines are [computer program](https://en.wikipedia.org/wiki/Computer_program) components that generalize [subroutines](https://en.wikipedia.org/wiki/Subroutine) for [non-preemptive multitasking](https://en.wikipedia.org/wiki/Non-preemptive_multitasking), by allowing execution to be suspended and resumed. Coroutines are well-suited for implementing familiar program components such as [cooperative tasks](https://en.wikipedia.org/wiki/Cooperative_multitasking), [exceptions](https://en.wikipedia.org/wiki/Exception_handling), [event loops](https://en.wikipedia.org/wiki/Event_loop), [iterators](https://en.wikipedia.org/wiki/Iterator), [infinite lists](https://en.wikipedia.org/wiki/Lazy_evaluation) and [pipes](https://en.wikipedia.org/wiki/Pipeline_(software)). 



## 参考

- [Go: Goroutine, OS Thread and CPU Management | by Vincent Blanchon | A Journey With Go | Medium](https://medium.com/a-journey-with-go/go-goroutine-os-thread-and-cpu-management-2f5a5eaf518a)
- [Go: GOMAXPROCS & Live Updates. ℹ️ This article is based on Go 1.13. | by Vincent Blanchon | A Journey With Go | Medium](https://medium.com/a-journey-with-go/go-gomaxprocs-live-updates-407ad08624e1)
- [goroutineはなぜ軽量なのか - Carpe Diem](https://christina04.hatenablog.com/entry/why-goroutine-is-good)
- [『Go 言語による並行処理』は Go 言語プログラマ必読書だろう | text.Baldanders.info](https://text.baldanders.info/remark/2018/11/concurrency-in-go/)

[Go]: https://golang.org/ "The Go Programming Language"
[Go言語による並行処理]: https://www.amazon.co.jp/dp/4873118468?tag=baldandersinf-22&linkCode=ogi&th=1&psc=1 "Go言語による並行処理 | Katherine Cox-Buday, 山口 能迪 |本 | 通販 | Amazon"
<!-- eof -->
