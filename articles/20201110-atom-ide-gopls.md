---
title: "ATOM エディタも gopls に対応していた" # 記事のタイトル
emoji: "🤔" # アイキャッチとして使われる絵文字（1文字だけ）
type: "idea" # "tech" : 技術記事 / "idea" : アイデア記事
topics: ["golang", "editor", "atom"] # タグ。["markdown", "rust", "aws"] のように指定する
published: true # 公開設定（true で公開）
---

Microsoft が誇る Visual Studo Code 最大の功績は「言語サーバ・プロトコル（Language Server Protocol）」を作ったことだと思う。

- [言語サーバープロトコルの概要 - Visual Studio | Microsoft Docs](https://docs.microsoft.com/ja-jp/visualstudio/extensibility/language-server-protocol?view=vs-2019)

これによって，コーディング作業をアシストするための入力補完やコード整形や名前をキーにしたジャンプといった機能を，言語やエディタ実装から独立させることが可能になった。ツール開発者はこのプロトコルに適合する言語サーバ・アプリケーションの作成に専念できる。

- [Langserver.org](https://langserver.org/)

各プログラミング言語コミュニティは自身の言語用サーバを提供している。 [Go] コミュニティでも早い段階からいくつかの実装が存在したが，現在では Google の [gopls][gopls] ( go please と読むらしい 😄) に集約されつつあるようだ。 [gopls] 自体はまだアルファ版だそうだが，それでも VSCode だけでなく Vim や Emacs などでも既に使われている。

この流れに見事に取り残されているのが [ATOM] エディタである。何故かというと [Go] 支援用パッケージの最大手である [go-plus] が今だに [gocode] を使っているのである。つか，実質的に [go-plus] の開発は2年くらい前から止まっているのだが。それでも [go-plus] を越えるパッケージが未だ見当たらないのは（[Go] プログラマから見て）いかに [ATOM] エディタが過疎ってるかを示す傍証と言えるかもしれない[^gcd1]。

[^gcd1]: ちなみに [gocode] は Windows 10 では，普通に go get でダウンロード&コンパイルしてもまともに動かない。動かすためにはコンパイル時に `-ldflags -H=windowsgui` のフラグを付加してやる必要があるが [go-plus] はガン無視して普通に go get して上書きインストールしてくさるので，どうにもならない。

（ホンマ GitHub は [ATOM] をどうするつもりなんだろう）

まぁ，でも，完全に見捨てられているわけではないようで， [gopls] では [ATOM] 用パッケージとして [ide-gopls] を[指定](https://github.com/golang/tools/blob/master/gopls/doc/atom.md)しているっぽい。

[ide-gopls] の機能は以下の通り。

- Auto completion
- Code format
- Diagnostics (errors & warnings)
- Document outline
- Find references
- Go to definition
- Hover
- Reference highlighting

今回ちょろんと試す機会があったのだが，感想は以下のような感じ。

- ユーザインタフェースが [atom-ide-ui] に統合されているためか，入力補完の出来は非常によい
- コード整形は gofmt レベルまでしか対応してない。せめて goimports レベルまではサポートしてほしかった（[gopls] 側の課題かな）
- Lint, test, coverage 等の機能がない。まぁ [x-terminal] からコマンドを叩けばいいんだけど。でもセーブした瞬間に lint や test が走るのって快感だよね 😄

比較対象が [go-plus] なので，評価が多少辛辣なのはご容赦を。まぁ「今後に期待」というところだろうか。 [gocode] はすでにソフトウェアとしての寿命を終えてるし [go-plus] の入力補完に我慢ならないという方は [ide-gopls] を応援し育てていくのもいいだろう。

私はせめて lint & test ができるようになってから，かな。もう少し我慢する。

## 参考

- [Big Sky :: gocode やめます(そして Language Server へ)](https://mattn.kaoriya.net/software/lang/go/20181217000056.htm)
- [Big Sky :: Go 言語の Language Server「gopls」が completeUnimported に対応した。](https://mattn.kaoriya.net/software/lang/c/20191112100330.htm)
- [gopls 0.4.3で構造体を初期化（"fillstruct"）しようとしても、"No code actions found"とだけ表示される - My External Storage](https://budougumi0617.github.io/2020/07/18/use_fillstruct_of_goplus_on_vim/)
- [vim-goを使わず、LSP（gopls）を使ってVimのGo開発環境を構築する - My External Storage](https://budougumi0617.github.io/2020/07/24/make_vimrc_with_lsp/)

[Go]: https://golang.org/ "The Go Programming Language"
[gopls]: https://github.com/golang/tools/tree/master/gopls "tools/gopls at master · golang/tools"
[gocode]: https://github.com/mdempsky/gocode "mdempsky/gocode: An autocompletion daemon for the Go programming language"
[ATOM]: https://atom.io/ "Atom"
[go-plus]: https://atom.io/packages/go-plus
[ide-gopls]: https://atom.io/packages/ide-gopls
[atom-ide-ui]: https://atom.io/packages/atom-ide-ui
[x-terminal]: https://atom.io/packages/x-terminal
<!-- eof -->
