---
title: "GitHub の Profile Readme に Feed を表示する" # 記事のタイトル
emoji: "💮" # アイキャッチとして使われる絵文字（1文字だけ）
type: "idea" # "tech" : 技術記事 / "idea" : アイデア記事
topics: ["github"] # タグ。["markdown", "rust", "aws"] のように指定する
published: true # 公開設定（true で公開）
---

[GitHub] ではユーザ名（私なら `spiegel-im-spiegel`）と同名のリポジトリにある `README.md` を使ってプロファイル・ページの Overview に追加の記述を載せることができる。

今回はその `README.md` にブログ等の feed の内容を表示する方法を紹介する。つっても [gautamkrishnar/blog-post-workflow] の [GitHub] Action を利用するだけの簡単なお仕事（笑）

まず，リポジトリにある `README.md` に以下の記述を追加する。

```
<!-- BLOG-POST-LIST:START -->
<!-- BLOG-POST-LIST:END -->
```

次にリポジトリ直下の `.github/workflows` ディレクトリ（なければ作成する）に `blog-post-workflow.yml` ファイルを作成する。中身はこんな感じ。

```yaml
name: Latest blog post workflow
on:
  schedule: # Run workflow automatically
    - cron: '0 * * * *' # Runs every hour, on the hour
  workflow_dispatch: # Run workflow manually (without waiting for the cron to be called), through the Github Actions Workflow page directly
jobs:
  update-readme-with-blog:
    name: Update this repo's README with latest blog posts
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v2
      - uses: gautamkrishnar/blog-post-workflow@master
        with:
          feed_list: "https://text.baldanders.info/index.xml, https://baldanders.info/index.xml"
```

この中の `feed_list` 項目を書き換えて参照したい feed の URL を列挙していく。他はとりあえず弄らなくてよい。

これで commit & push し， Action を起動すれば


```markdown
<!-- BLOG-POST-LIST:START -->
- [個人番号と個人番号カード](https://text.baldanders.info/remark/2020/09/my-number-and-my-number-card/)
- [2020-09-20 のブックマーク](https://text.baldanders.info/bookmarks/2020/09/20-bookmarks/)
- [近ごろ流行りらしい “Zenn” のアカウントを作ってみた](https://text.baldanders.info/remark/2020/09/using-zenn-with-github/)
- [NIST SP 800-207: “Zero Trust Architecture”](https://text.baldanders.info/remark/2020/09/nist-sp-800-207-zero-trust-architecture/)
- [Java 15 がリリースされた](https://text.baldanders.info/release/2020/09/java-15-is-released/)
<!-- BLOG-POST-LIST:END -->
```

てな感じに一覧を挿入してくれる。ちなみに，上述の YAML 設定だと cron で1時間毎に Action が起動する設定になっている。

cron のタイミングを変えたり，複数の feed を別々に取得して `README.md` の異なる位置に挿入することもできる。詳しくは [gautamkrishnar/blog-post-workflow] にカスタマイズ方法が載っているので参考になるだろう。

## 参考

- [GitHub プロファイルを（ちょっとだけ）カッコよくしてみる — しっぽのさきっちょ | text.Baldanders.info](https://text.baldanders.info/remark/2020/09/using-github-profile-readme/)

[GitHub]: https://github.com/
[gautamkrishnar/blog-post-workflow]: https://github.com/gautamkrishnar/blog-post-workflow "gautamkrishnar/blog-post-workflow: Show your latest blog posts from any sources or StackOverflow activity or Youtube Videos on your GitHub profile/project readme automatically using the RSS feed"
