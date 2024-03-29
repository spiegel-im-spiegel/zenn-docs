---
title: "エンジニアならコードで語れ！"
---

最初は「[Go] で [PostgreSQL] の BLOB (Binary Large OBject) データって簡単に扱えるん？」と軽い動機だったのだが，素の [pq][github.com/lib/pq] または [pgx][github.com/jackc/pgx] ドライバを使うのはちょっとかったるいし [Go] ではメジャーな [GORM] や [ent] を使うのがいいかなぁ，と沼にハマり始めた（笑） さらに個人的にお気に入りの [zerolog][github.com/rs/zerolog] を組み込めないのか，と試行錯誤し始め，気が付いたら「[Go] で [PostgreSQL] の BLOB (Binary Large OBject) データって簡単に扱えるん？」を知りたいために丸2日も費やしてしまったですよ。

で，せっかく2日も使って調べたのだからせめて作業記録をブログ記事として残そうと思ったのだが，途中まで下書きして気が付いた。「これ3回シリーズでも終わらんわ」「じゃあ Zenn 本にすればいいじゃない！」。で，どうせ Zenn 本にするならもう少しちゃんと書かないと，と更に1フィートくらい深堀りした結果が今回のアウトプットである。

公式サイトを含めネットの情報をあちこち彷徨い歩いた印象は「おまーらホンマにこのコードで動かしたのか？」だった。いや，たとえブログ記事でもコードを丸ごと載せたら冗長でピンぼけな内容になりがちだし，最低限のコードに絞るというのは悪くない戦略だけど，せめて import パッケージは明示して欲しかった。

というわけで，この本では言葉の説明は最低限にして，とにかく動くコードを載せまくることにした。目標は来年の私がこの本を起点にして [PostgreSQL] の再調査ができるレベルまで持っていくことである。まぁ，我ながら「本」とは思えない狂った内容だと自覚してるので，あまりいぢめないでやってください（笑）

簡単にまとめっぽいものを書いておくと

1. 全ての Gopher は [PostgreSQL] ドライバを [pgx][github.com/jackc/pgx] に乗り換えるべし
2. もう logger は [zerolog][github.com/rs/zerolog] でいいぢゃん。標準以外の選択肢に [zerolog][github.com/rs/zerolog] を受け付けてくれよ
3. 既にテーブル設計が完了しているなら [GORM] のほうが扱いやすい[^gen1]
4. 設計要件でデータ間のグラフ構造をより重視するのであれば，要件定義の段階から [ent] でコードを書くべし

[^gen1]: [GORM] と組み合わせて使える [smallnest/gen][github.com/smallnest/gen] といった製品もある。 [smallnest/gen][github.com/smallnest/gen] では [GORM] で利用可能な Model 構造体を生成できる。

といったところだろうか。特に [ent] は目から鱗が落ちた。

レガシー・システムではテーブル設計は最重要で真っ先に決定すべき事項とされているが，これはテーブル構造はシステム設計においてもっとも「依存されるもの」であり高安定性が要求されるからだ。でも本当に「依存されるもの」であり高安定性が要求されるのはテーブル構造ではなくグラフ構造なのだ。

[ent] が秀逸なのは未分化だったテーブルとグラフを分離し，テーブルをグラフの内側にカプセル化するという夢の設計を高いレベルで実現している点だろう。その代わり [ent] ではグラフの制御をテーブル側に傾けすぎないよう配慮している印象がある。ゆえにテーブル構造への直接的な干渉が難しく，最初にテーブルありきで設計を始めると [ent] はとてつもなく扱いにくいフレームワークになってしまう。設計の最初の1フィートがこれまでとは違うのだ。

おそらく，実験的実装以外で [ent] を採用するのであれば，要件定義の段階からコードを書いて議論していかないと無理だろう。特にレガシー・システムから入れ替えたい場合などは，かなり斬新な設計変更が求められるんじゃないかな。紙や Excel でスキーマ設計する時代は終わったのである（笑）

まぁ，でも，いまどきはデータ同士の関係が合理的にグラフ化出来ないと DX もへったくれもないので，言語に依らず [ent] のようなアプローチは今後増えていくのかもしれない。

[Go]: https://go.dev/
[PostgreSQL]: https://www.postgresql.org/ "PostgreSQL: The world's most advanced open source database"
[github.com/lib/pq]: https://github.com/lib/pq "lib/pq: Pure Go Postgres driver for database/sql"
[github.com/jackc/pgx]: https://github.com/jackc/pgx "jackc/pgx: PostgreSQL driver and toolkit for Go"
[ent]: https://entgo.io/
[GORM]: https://gorm.io/ "GORM - The fantastic ORM library for Golang, aims to be developer friendly."
[github.com/rs/zerolog]: https://github.com/rs/zerolog "rs/zerolog: Zero Allocation JSON Logger"
[github.com/smallnest/gen]: https://github.com/smallnest/gen "smallnest/gen: Converts a database into gorm structs and RESTful api"
