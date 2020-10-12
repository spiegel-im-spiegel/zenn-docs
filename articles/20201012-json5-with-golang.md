---
title: "Go で JSON5 データを読み込む" # 記事のタイトル
emoji: "⌨" # アイキャッチとして使われる絵文字（1文字だけ）
type: "tech" # "tech" : 技術記事 / "idea" : アイデア記事
topics: ["go", "programming", "json"] # タグ。["markdown", "rust", "aws"] のように指定する
published: true # 公開設定（true で公開）
---

いつものように小ネタです。

JSON (JavaScript Object Notation) はデータ交換用のテキストフォーマットとしてはシンプルで非常に優れているのだが，人間側から見て読み／書きやすいかというと微妙だったりする。

その辺の不満を解消しようというのが [JSON5] である。

- [JSON5 | JSON for Humans][JSON5]

[JSON5] の特徴は以下の通り（「[知らないうちにJSON5 in Babel - Qiita](https://qiita.com/jkr_2255/items/026e0fdb4570c88c4f51)」からの引用）

> * オブジェクト
>     * キーは一重引用符でも、識別子として使えるものならOK
>     * ケツカンマOK
> * 配列
>     * ケツカンマOK
> * 文字列
>     * 一重引用符もOK
>     * 文字列の途中で（バックスラッシュでエスケープすれば）改行OK
>     * `\x0f`のようなエスケープシーケンスを使える
> * 数値
>     * `0xabcd`のような16進法でもOK
>     * `16.`や`.23`のような、小数点前後の片方だけ書くのもOK
>     * `+12`のようなプラス記号も使える
>     * `Infinity`や`NaN`も値として受け付ける 
> * コメント
>     * `/* ... */`あるいは`// (このあと行末まで) `といったコメントを使える

例えばこんな感じの記述ができる（[オフィシャル・サイト][JSON5]より）。

```json5:json5-sample.json5
{
  // comments
  unquoted: 'and you can quote me on that',
  singleQuotes: 'I can use "double quotes" here',
  lineBreaks: "Look, Mom! \
No \\n's!",
  hexadecimal: 0xdecaf,
  leadingDecimalPoint: .8675309, andTrailing: 8675309.,
  positiveSign: +1,
  trailingComma: 'in objects', andIn: ['arrays',],
  "backwardsCompatible": "with JSON",
}
```

しかし，逆に既存のツール等で [JSON5] テキストを読み込ませるのは難しかったりする。たとえば事実上の標準である [jq] に読み込ませようとしても

```
$ cat json5-sample.json5 | jq .
parse error: Invalid numeric literal at line 2, column 5
```

てな具合に最初の1フィートでぶっコケる。

[オフィシャルのリポジトリ](https://github.com/json5 "JSON5")を見ると少なくとも node.js 用のパーサはあるようだが，残念なことに [Go] 用のパッケージは Archive 化されていて使いものにならないようだ。

で，どなたか [Go] 用の [JSON5] パーサ作ってないかなぁ，と軽くググってみたら，以下のパッケージがよさげである。

- [flynn/json5: Go JSON5 decoder package based on encoding/json](https://github.com/flynn/json5)

[JSON5] テキストの Decode (または Unmarshal) が標準の [encoding/json] と互換になっていて，パッケージ名を `json5` に置き換えれば済みそうだ。

こんな感じ。

```go:sample.go
package main

import (
    "encoding/json"
    "fmt"
    "os"

    "github.com/flynn/json5"
)

func main() {
    elms := make(map[string]interface{})
    if err := json5.NewDecoder(os.Stdin).Decode(&elms); err != nil {
        fmt.Fprintln(os.Stderr, err)
        return
    }
    if err := json.NewEncoder(os.Stdout).Encode(elms); err != nil {
        fmt.Fprintln(os.Stderr, err)
    }
}
```

これを使えば

```
$ cat json5-sample.json5 | go run sample.go | jq .
{
  "andIn": [
    "arrays"
  ],
  "andTrailing": 8675309,
  "backwardsCompatible": "with JSON",
  "hexadecimal": 912559,
  "leadingDecimalPoint": 0.8675309,
  "lineBreaks": "Look, Mom! No \\n's!",
  "positiveSign": 1,
  "singleQuotes": "I can use \"double quotes\" here",
  "trailingComma": "in objects",
  "unquoted": "and you can quote me on that"
}
```

といい感じに JSON に変換できる。

4年前に更新されたきりでモジュールにも対応していないのだが，まぁ鼻の先は問題ないだろう。

[Go]: https://golang.org/ "The Go Programming Language"
[JSON5]: https://json5.org/ "JSON5 | JSON for Humans"
[jq]: https://stedolan.github.io/jq/
[encoding/json]: https://golang.org/pkg/encoding/json/ "json - The Go Programming Language"
<!-- eof -->
