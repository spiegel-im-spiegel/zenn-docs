---
title: "インデントおよび行末は EditorConfig で始末する" # 記事のタイトル
emoji: "💮" # アイキャッチとして使われる絵文字（1文字だけ）
type: "idea" # "tech" : 技術記事 / "idea" : アイデア記事
topics: ["editor"] # タグ。["markdown", "rust", "aws"] のように指定する
published: true # 公開設定（true で公開）
---

[EditorConfig] はテキストエディタや IDE (Integrated Development Environment; 統合開発環境) などで文字コードやインデントや改行コードなどの設定を共有するための仕組みで，メジャーなエディタや IDE なら既定で組み込まれているか拡張機能で導入することができる。これによって開発環境や個人設定の差異によるコーデイング・スタイルの混乱を抑えることができる。

## [EditorConfig] による設定

[EditorConfig] を有効にするにはプロジェクトのトップ・フォルダに `.editorconfig` ファイルを置けばよい。たとえば私がよく使う `.editorconfig` ファイルの中身はこんな感じ。

```
root = true

[*]
end_of_line = lf
charset = utf-8
indent_style = tab
indent_size = 4
tab_width = 4
trim_trailing_whitespace = false
insert_final_newline = true

[*.go]
trim_trailing_whitespace = true

[*.md]
indent_style = space
indent_size = 4

[*.yml]
indent_style = space
indent_size = 2
trim_trailing_whitespace = true
```

[EditorConfig] はフォルダを遡って `.editorconfig` ファイルを探し，フォルダの上から順番に評価していく。 `root = true` の記述がないとどこまでも上の階層に遡っていくので，プロジェクトのトップ・フォルダの `.editorconfig` ファイルには必ずこれを記述すること。

- `[...]` は対象となるファイルを指定している。 `[*]` なら全てのファイルが対象， `[*.go]` は拡張子が `go` のファイルが対象となる
- `indent_style` はインデントのスタイルを指定する。 `tab` または `space` を指定する
- `indent_size` はインデントの幅を指定する。 `indent_style` が `tab` の場合は `tab_width` で指定するようだ
- `end_of_line` は改行コードを指定する。 `lf`, `cr`, `crlf` から選択できる
- `chaset` は文字エンコーディングを指定する。 `latin1`, `utf-8`, `utf-8-bom`, `utf-16be` or `utf-16le` から選択できる。残念ながらこれ以外の文字エンコーディングについてはエディタ側の実装に依存する
- `trim_trailing_whitespace` を `true` にすると行末の空白文字を削除してくれる
- `insert_final_newline` を `true` にするとファイルの末尾が改行文字ではない場合に補完してくれる

詳しい仕様については[仕様書ページ](https://editorconfig-specification.readthedocs.io/en/latest/ "EditorConfig Specification — EditorConfig Specification documentation")を参考にどうぞ。ただしエディタや IDE によっては全ての機能を網羅していない場合があるのでご注意を。

## [EditorConfig] によるコーディング規約の統一

こういう仕組みがあればドキュメントで「コーディング規約」を周知しなくても `.editorconfig` ファイルをリポジトリに放りこんでおけば済む[^cr1]。というかリポジトリを作ったらまず `.editorconfig` ファイルをセットするよう習慣づけるべきだろう。

[^cr1]: それでも，曖昧な記述を許容する言語ではコーディング規約がないと困るかもしれないが。

まぁ，最近は Go や Rust みたいに公式の整形ツールが用意されている言語もあるので[^pr1]，昔ほどの需要はないかもしれないけど。

[^pr1]: GitHub で見かける Go パッケージでは pull request を発行する前に gofmt などで整形することを要求しているものもある。

## 参考

- [EditorConfigで文字コード設定を共有して喧嘩しなくなる話。（Frontrend Advent Calendar 2014 – 14日目） | Ginpen.com](http://ginpen.com/2014/12/14/editorconfig/)
- [どんなエディタでもEditorConfigを使ってコードの統一性を高める - Qiita](https://qiita.com/naru0504/items/82f09881abaf3f4dc171)
- [【ATOM Editor】 EditorConfig を使うなら Whitespace は不要 — しっぽのさきっちょ | text.Baldanders.info](https://text.baldanders.info/remark/2016/10/warnig-from-editorconfig-at-atom/)

[EditorConfig]: https://editorconfig.org/
