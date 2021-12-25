---
title: "Zenn のダッシュボードに統計情報が表示されるようになるらしい" # 記事のタイトル
emoji: "💮" # アイキャッチとして使われる絵文字（1文字だけ）
type: "idea" # "tech" : 技術記事 / "idea" : アイデア記事
topics: ["zenn"] # タグ。["markdown", "rust", "aws"] のように指定する
published: true # 公開設定（true で公開）
---

[Zenn] から「振り返りレポート」のメールが届いてました。

![](/images/zenn-analytics/zenn-report.png)

うちみたいな小ネタや与太話ばっかりのページでも読んでいただけるとは有り難い話である。
来年もボチボチ頑張ります（笑）

「振り返りレポート」に PV の項目があって「ユーザごとに集計してて偉いなぁ」と思っていたが， [zenn-dev/zenn-community](https://github.com/zenn-dev/zenn-community "zenn-dev/zenn-community: zenn.dev roadmap") の「[ダッシュボードでPVなどの統計データを見れるように](https://github.com/zenn-dev/zenn-community/issues/98)」が更新されていて


> - 以前はAnalyticsのReporting APIを使ってページビュー等の統計情報を取得していたため、ユーザーごとに統計情報を見れるようにするとLimitを超えてしまう可能性が高いという問題がありました。
> - 数ヶ月前にBigQueryを導入したことにより、ユーザーごとに細かなデータを取得することが可能になりました。
> - 2022年にはzennのダッシュボードに統計ページを追加し、以下のような情報を見れるようにする予定です。
>   - 投稿ごとの合計ページビュー / Like
>   - 日別や月別のページビューの推移
>   - 投稿アクティビティ（GitHubのContributionsの草のような形で投稿量を振り返れるように）
>
>（「[ダッシュボードでPVなどの統計データを見れるように](https://github.com/zenn-dev/zenn-community/issues/98#issuecomment-1000624531)」より）

と書かれていた。なんと！

ひょっとして今回の「振り返りレポート」は来年からの機能追加の前ふりということなのだろうか。楽しみである。

[Zenn]: https://zenn.dev/ "Zenn｜エンジニアのための情報共有コミュニティ"
