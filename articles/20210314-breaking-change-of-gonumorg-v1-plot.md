---
title: "gonum.org/v1/plot パッケージの破壊的変更" # 記事のタイトル
emoji: "📊" # アイキャッチとして使われる絵文字（1文字だけ）
type: "tech" # "tech" : 技術記事 / "idea" : アイデア記事
topics: ["go", "programming"] # タグ。["markdown", "rust", "aws"] のように指定する
published: true # 公開設定（true で公開）
---

グラフ描画パッケージの一つである [gonum.org/v1/plot][plot] が v0.9.0 にアップデートされたが，破壊的変更を含んでいるようなのでメモしておく。

## plot.New 関数の返り値

[plot].Plot インスタンスを生成する [plot].New() 関数で，以前のバージョンは

```go
p, err := plot.New()
```

という感じに error も一緒に返却されていたが， v0.9.0 では

```go
p := plot.New()
```

と error を返さなくなった。いや，まぁ，エラー・ハンドリングが不要になる分ちょびっとだけ楽になるのでいいのだが。

## 既定フォントとフォント指定の変更

以前は

```go
plot.DefaultFont = "Helvetica"
plotter.DefaultFont = "Helvetica"
```

という感じに文字列でフォントを指定していたが，どうもフォント制御の部分をごっそり作り変えたようで，フォントの指定方法だけでなくフォントそのものも入れ替わった。具体的には

```go
plot.DefaultFont = font.Font{
    Typeface: "Liberation",
    Variant:  "Sans",
}
plotter.DefaultFont = plot.DefaultFont
```

と言った感じに指定する。 [gonum.org/v1/plot][plot] があらかじめキャッシュして指定可能なフォントは以下の通り。

| Typeface       | Variant   | Style              | Weight              |
| -------------- | --------- | ------------------ | ------------------- |
| `"Liberation"` | `"Serif"` | [font].StyleNormal | [font].WeightNormal |
| `"Liberation"` | `"Serif"` | [font].StyleNormal | [font].WeightBold   |
| `"Liberation"` | `"Serif"` | [font].StyleItalic | [font].WeightNormal |
| `"Liberation"` | `"Serif"` | [font].StyleItalic | [font].WeightBold   |
| `"Liberation"` | `"Sans"`  | [font].StyleNormal | [font].WeightNormal |
| `"Liberation"` | `"Sans"`  | [font].StyleNormal | [font].WeightBold   |
| `"Liberation"` | `"Sans"`  | [font].StyleItalic | [font].WeightNormal |
| `"Liberation"` | `"Sans"`  | [font].StyleItalic | [font].WeightBold   |
| `"Liberation"` | `"Mono"`  | [font].StyleNormal | [font].WeightNormal |
| `"Liberation"` | `"Mono"`  | [font].StyleNormal | [font].WeightBold   |
| `"Liberation"` | `"Mono"`  | [font].StyleItalic | [font].WeightNormal |
| `"Liberation"` | `"Mono"`  | [font].StyleItalic | [font].WeightBold   |

このうち [font].StyleNormal と [font].WeightNormal は既定値なので省略可能である。またフォントを指定しない場合は Liberation/Serif が既定値としてセットされている。なお Style と Weight の指定で使用する [font] パッケージは [gonum.org/v1/plot][plot]/font ではなく [golang.org/x/image/font][font] の方なので，パッケージ名の衝突に注意すること。

ちなみに [Liberation フォント][liberation-fonts]は SIL Open Font License 1.1 で提供されていて， [Go] からは [go-fonts/liberation] パッケージでコントロールできる。

まぁ，確かに破壊的変更だが，今の構成なら [golang.org/x/image/font][font] パッケージとの親和性も高そうだし，外部のフォントファイルも扱いやすくなるのかな。面倒くさいのでしないけど。

## 参考

https://text.baldanders.info/golang/chart-with-golang/

[Go]: https://golang.org/ "The Go Programming Language"
[plot]: https://github.com/gonum/plot "gonum/plot: A repository for plotting and visualizing data"
[font]: https://pkg.go.dev/golang.org/x/image/font "font · pkg.go.dev"
[liberation-fonts]: https://github.com/liberationfonts/liberation-fonts/ "liberationfonts/liberation-fonts: The Liberation(tm) Fonts is a font family which aims at metric compatibility with Arial, Times New Roman, and Courier New."
[go-fonts/liberation]: https://github.com/go-fonts/liberation "go-fonts/liberation: Liberation fonts for Go"
<!-- eof -->
