---
title: "Zenn の記事を GitHub で管理する" # 記事のタイトル
emoji: "💮" # アイキャッチとして使われる絵文字（1文字だけ）
type: "idea" # "tech" : 技術記事 / "idea" : アイデア記事
topics: ["zenn", "github"] # タグ。["markdown", "rust", "aws"] のように指定する
published: true # 公開設定（false にすると下書き）
---

試しに [Zenn] のアカウントを取ってみた。

[Zenn] では [GitHub] のリポジトリと連携してコンテンツの作成・更新ができるらしい。

- [GitHubリポジトリでZennのコンテンツを管理する | Zenn](https://zenn.dev/zenn/articles/connect-to-github)

ただし [GitHub] → [Zenn] への一方通行の deploy のようで，以下の制限がある。

- リポジトリ上の記事を削除しても [Zenn] に反映されない
- 一度 [Zenn] に deploy された記事の slug は変更できない（別の記事として扱われる）
- 既に [Zenn] でオン書きしたコンテンツは [GitHub] に反映されない

まぁ，この辺はしょうがないだろう。今後に期待しておく。

## 【追記】

自ブログサイトに今回の作業メモを上げておいた。

- [近ごろ流行りらしい “Zenn” のアカウントを作ってみた — しっぽのさきっちょ | text.Baldanders.info](https://text.baldanders.info/remark/2020/09/using-zenn-with-github/)

[Zenn] では記事のライセンスを設定できないが，[リポジトリの方で CC BY-SA ライセンスを付与](https://github.com/spiegel-im-spiegel/zenn-docs/blob/main/LICENSE)している。

ということで再利用は（条件の範囲内で）ご自由にどうぞ。

## 参照記事

- [Zenn CLIをインストールする | Zenn](https://zenn.dev/zenn/articles/install-zenn-cli)
- [Zenn CLIを使ってコンテンツを作成する | Zenn](https://zenn.dev/zenn/articles/zenn-cli-guide)
- [ZennのMarkdown記法 | Zenn](https://zenn.dev/zenn/articles/markdown-guide)
    - [Supported Functions · KaTeX](https://katex.org/docs/supported.html) : [Zenn] で数式として使える表現
- [Ubuntu/Debianに最新のNode.jsをインストールする一番良い方法 | LFI](https://linuxfan.info/install_nodejs_on_ubuntu_debian)
- [Twemoji](https://twemoji.twitter.com/) : [Zenn] は Twitter Emoji を使っているらしい
    - [twitter/twemoji: Emoji for everyone. https://twemoji.twitter.com/](https://github.com/twitter/twemoji)
    - [絵文字一覧 🤣 | Let's EMOJI](https://lets-emoji.com/emojilist/) : よく使いそうなのは以下かな
        - [絵文字一覧（顔文字と感情：Smileys & Emotion）😀 | Let's EMOJI](https://lets-emoji.com/emojilist/emojilist-1/)
        - [絵文字一覧（物：Objects）📌 | Let's EMOJI](https://lets-emoji.com/emojilist/emojilist-7/)
- [🎁 Emoji cheat sheet for GitHub, Basecamp, Slack & more](https://www.webfx.com/tools/emoji-cheat-sheet/)

[Zenn]: https://zenn.dev/ "Zenn｜プログラマーのための情報共有コミュニティ"
[GitHub]: https://github.com/
