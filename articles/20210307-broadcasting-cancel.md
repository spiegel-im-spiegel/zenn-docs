---
title: "context.Context を掴みっぱなしにしてはいけない" # 記事のタイトル
emoji: "🤔" # アイキャッチとして使われる絵文字（1文字だけ）
type: "tech" # "idea" : 技術記事 / "idea" : アイデア記事
topics: ["go", "programming"] # タグ。["markdown", "rust", "aws"] のように指定する
published: true # 公開設定（true で公開）
---

週末に書籍情報にアクセスする以下の自作パッケージをアップデートした。

https://github.com/spiegel-im-spiegel/pa-api/releases/tag/v0.9.0
https://github.com/spiegel-im-spiegel/aozora-api/releases/tag/v0.3.0
https://github.com/spiegel-im-spiegel/openbd-api/releases/tag/v0.3.0

今回はちょっと破壊的変更がある。まぁ，私以外に使う人はあまりいなさそうなので影響は少ないと思うけど。

この3つのパッケージをつくり始めたのは[1,2年くらい前](https://text.baldanders.info/remark/2019/10/pa-api-v5/ "PA-API v5 への移行")なのだが，当時は [context] パッケージの挙動とかあまり考えずに「えいやっ」で組んでいて，最近見直して「拙いなぁ」とは思ってたのよ。

さらに先日の [Go] 1.16 リリース直後に公開された公式ブログ記事

https://blog.golang.org/context-and-structs

を見て「やっぱそうだよね」と納得し，いつ直そうか思案してた。

上の記事って要するに「構造体の要素として [context].Context を掴みっぱなしにするな」という話で，これ自体は [context] パッケージのドキュメントにも明記されている。

>Do not store Contexts inside a struct type; instead, pass a Context explicitly to each function that needs it.

理由は上のブログ記事で詳しく解説されているので，一度は目を通しておくことを強くお勧めする。

で，まぁ，思いっきり言い訳なんだけど [context] パッケージの目的が

>Context provides a means of transmitting deadlines, caller cancellations, and other request-scoped values across API boundaries and between processes.

というものである以上，他レイヤあるいはドメインに渡すコンテキスト情報の一部として [context].Context を含めたくなる気持ちも分かって欲しい orz

奇しくも昨日の[読書会イベント](https://gpl-reading.connpass.com/event/204017/ "第10回『プログラミング言語Go』オンライン読書会 - connpass")でチャネルのクローズを使ったキャンセル・イベント伝搬の節（『[プログラミング言語Go]』8.9章）が出てきて「これは『はよ直せ』という啓示か？」とようやく重い腰を上げることにしたのだった。

なんでチャネルのクローズがキャンセル・イベントに使えるかというと， [Go] のチャネルは，バッファなしまたはバッファが空の状態では待ち受けてるゴルーチンがブロックされるのだが，何も送信しなくてもチャネルがクローズするタイミングでブロックが解除される（または待ちなしで受信処理を抜ける），という性質がある。これはひとつのチャネルを複数のゴルーチンで待ち受けている場合でも並行かつ平等に発生する[^conc1]。つまり「チャネルのクローズ」をイベント・トリガーとして「放送」できるわけだ。

[^conc1]: [Go] のゴルーチン同士には優先順位がない。 RTOS なんかでありがちな priority inversion は起こらないわけだ。言い方を変えるとシビアなリアルタイム処理が要求されるシステムには [Go] は向かないってことなんだけど。という話を[読書会イベント]の雑談でしていた。

んで， [Go] 1.7 から実装された [context] パッケージはこの仕組みを実に上手く使っている。違う見方をすると「チャネル・クローズを使ったキャンセル・イベントの伝搬」の仕組みを頭に入れた上で [context].Context を実装しないと思わぬ副作用が出ててしまうのだ。例えば “[Contexts and structs](https://blog.golang.org/context-and-structs)” で例示されるような。

あと，先日公開された

https://zenn.dev/imamura_sh/articles/retry-46aa586aeb5c3c28244e

という記事に触発されたというのもある。この記事の中に「[context.Contextの終了を確認する](https://zenn.dev/imamura_sh/articles/retry-46aa586aeb5c3c28244e#context.context%E3%81%AE%E7%B5%82%E4%BA%86%E3%82%92%E7%A2%BA%E8%AA%8D%E3%81%99%E3%82%8B)」という節があるが， [context].Context の仕組みが分かっていれば納得の内容だろう。

私の場合「[URI からデータを取ってくるだけの簡単なお仕事](https://zenn.dev/spiegel/articles/20210113-fetch)」をするパッケージを書いたときに構造体の要素として [context].Context を掴んでしまわないよう考慮して書いた（つもりな）ので，他の自作パッケージで HTTP クライアント操作をする部分はこの [github.com/spiegel-im-spiegel/fetch](https://github.com/spiegel-im-spiegel/fetch "spiegel-im-spiegel/fetch: Fetch Data from URL") パッケージで置き換えることで対応できた。

まぁ，そのせいで破壊的変更になってしまったんだけど。こうやって[技術的負債](https://text.baldanders.info/remark/2020/12/technical-debt-and-hacker/ "技術的負債とハッカー")をちまちまと返済していくんですねぇ（笑）

[Go]: https://golang.org/ "The Go Programming Language"
[context]: https://golang.org/pkg/context/ "https://golang.org/pkg/context/"
[プログラミング言語Go]: https://www.amazon.co.jp/dp/4621300253
[読書会イベント]: https://gpl-reading.connpass.com/event/204017/ "第10回『プログラミング言語Go』オンライン読書会 - connpass"
<!-- eof -->
