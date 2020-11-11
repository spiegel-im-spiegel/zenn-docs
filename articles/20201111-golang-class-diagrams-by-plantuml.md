---
title: "Go パッケージのクラス図を PlantUML で描く"
emoji: "💮" # アイキャッチとして使われる絵文字（1文字だけ）
type: "idea" # "tech" : 技術記事 / "idea" : アイデア記事
topics: ["golang", "uml"] # タグ。["markdown", "rust", "aws"] のように指定する
published: true # 公開設定（true で公開）
---

今回は軽く小ネタで。

[Go] で作ったパッケージを UML 図で表せたらいいのに，と思ったことはないだろうか。そう思う人は結構いるらしく，いろんなツールが公開されている。それらの中でも今回は [jfeliu007/goplantuml][goplantuml] を紹介する。

[goplantuml] は [Go] パッケージを解析するためのパーサと，それを使った CLI (Command-Line Interface) ツールで構成されている。また，このパッケージを使った [Dumels] という Web サービスもあるらしい。今回は CLI ツールの方を試してみる。

バイナリは用意されてないようなので，おとなしく go get コマンドでダウンロード&ビルドする。モジュール・モードが on になっているなら，以下で無問題。

```
$ go get github.com/jfeliu007/goplantuml/cmd/goplantuml@latest

$ goplantuml -h
Usage of goplantuml:
  -aggregate-private-members
    	Show aggregations for private members. Ignored if -show-aggregations is not used.
  -hide-connections
    	hides all connections in the diagram
  -hide-fields
    	hides fields
  -hide-methods
    	hides methods
  -ignore string
    	comma separated list of folders to ignore
  -notes string
    	Comma separated list of notes to be added to the diagram
  -output string
    	output file path. If omitted, then this will default to standard output
  -recursive
    	walk all directories recursively
  -show-aggregations
    	renders public aggregations even when -hide-connections is used (do not render by default)
  -show-aliases
    	Shows aliases even when -hide-connections is used
  -show-compositions
    	Shows compositions even when -hide-connections is used
  -show-connection-labels
    	Shows labels in the connections to identify the connections types (e.g. extends, implements, aggregates, alias of
  -show-implementations
    	Shows implementations even when -hide-connections is used
  -show-options-as-note
    	Show a note in the diagram with the none evident options ran with this CLI
  -title string
    	Title of the generated diagram
```

よしよし。

解析を行うにはパッケージのあるディレクトリを引数として渡せばよい。

```
$ goplantuml ~/go/src/github.com/spiegel-im-spiegel/pa-api > pa-api.puml
```

解析結果は [PlantUML] の記述形式で標準出力に出力されるので，適当にリダイレクトしておく。あとは [PlantUML] を使って画像データに変換すればよい。

```
java -jar /path/to/plantuml.jar -charset UTF-8 pa-api.puml
```

結果はこんな感じ。

![pa-api.png](https://storage.googleapis.com/zenn-user-upload/lg1kawxhh6ebocxb4sqbfudxk7up)

ちゃんとパッケージ単位でまとめられているのが分かるだろう。なお `-recursive` オプションを付けるとサブディレクトリのパッケージも再帰的に解析してくれる。本来 UML 図を描くなら多重度が必須だが，今回はコードから図を起こしてるのだから重要ではあるまい。

ドキュメンテーションのオトモにどうぞ。

## Windows では dot コマンドに注意

Windows 版 [Graphviz] 2.44 に含まれる dot コマンドを使う場合[^dot1]，コマンドプロンプトで dot.exe コマンドのあるフォルダまで降りて `dot -c` コマンドを打っておく必要があるらしい。

[^dot1]: [PlantUML] は描画に [Graphviz] の dot コマンドを使う。

- [Important note about version](https://plantuml.com/ja/graphviz-dot)

## おまけ

拙作の [spiegel-im-spiegel/depm][depm] を使えばモジュール単位の依存関係を可視化できる。たとえばこんな感じ。

[![depm modules](https://text.baldanders.info/release/dependency-graph-for-golang-modules/output3.png)](https://text.baldanders.info/release/dependency-graph-for-golang-modules/output3.png)

詳しくは以下の紹介ページを参考にどうぞ。

- [Depm: Go 言語用モジュール依存関係可視化ツール](https://text.baldanders.info/release/dependency-graph-for-golang-modules/)

以上，広告でした（笑）

## 参考

- [bykof/go-plantuml: Generate plantuml diagrams from go source files or directories](https://github.com/bykof/go-plantuml)
    - [Generate plantuml diagrams from go source files or directories](https://golangexample.com/generate-plantuml-diagrams-from-go-source-files-or-directories/)
- [真面目に PlantUML (1) : PlantUML のインストール](https://text.baldanders.info/remark/2018/12/plantuml-1/)
- [真面目に PlantUML (3) : クラス図](https://text.baldanders.info/remark/2018/12/plantuml-3-class-diagrams/)

[Go]: https://golang.org/ "The Go Programming Language"
[goplantuml]: https://github.com/jfeliu007/goplantuml "jfeliu007/goplantuml: PlantUML Class Diagram Generator for golang projects"
[Dumels]: https://www.dumels.com/ "Dumels"
[PlantUML]: https://plantuml.com/ "Open-source tool that uses simple textual descriptions to draw beautiful UML diagrams."
[Graphviz]: https://www.graphviz.org/ "Graphviz - Graph Visualization Software"
[depm]: https://github.com/spiegel-im-spiegel/depm "spiegel-im-spiegel/depm: Visualize depndency packages and modules"
<!-- eof -->
