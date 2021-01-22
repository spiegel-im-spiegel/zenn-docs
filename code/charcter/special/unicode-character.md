# Unicode 文字種を判別する

いつもの小ネタ。

別記事で「[やっかいな日本語](https://zenn.dev/spiegel/articles/20210118-characters)」という記事を書いたが，今回はもう少し [Go] 寄りに Unicode 文字の判別について紹介してみる。

Unicode 文字の種類を判別するには [unicode] 標準パッケージが使える。判別用の [unicode].RangeTable を用意し，これを参照することで文字種を判別する。

[unicode] 標準パッケージの中身を見ると分かるが，かなりの数の定義済み [unicode].RangeTable テーブルを取り揃えている。今回はこの定義済みテーブルのみ使っていく。

## 図形文字と制御文字

まずは大雑把に「図形文字」と「制御文字」を判別してみよう。

図形文字の判別には [unicode].IsGraphic() 関数が，制御文字の判別には [unicode].IsControl() 関数が使える。ただし [unicode].IsControl() 関数では U+00FF 以下の ISO 8859 で定義されている制御文字領域しか判別してくれないようで BOM (U+FEFF) などの Unicode 独自の制御文字も含めて判別したいなら [unicode].C テーブルを使う必要がある。

そこで，こんな関数を考えてみる。

```go:sample1.go
import "unicode"

func check(r rune) string {
    switch {
    case unicode.IsGraphic(r):
        return "Graphic"
    case unicode.IsControl(r):
        return "Latin1 Control"
    case unicode.Is(unicode.C, r):
        return "Unicode Control"
    }
    return "Unknown"
}
```

これを使って実際に文字列をチェックしてみよう。

```go:sample1.go
func main() {
    text := string([]byte{0xef, 0xbb, 0xbf, 0xe3, 0x82, 0x84, 0x09, 0xe3, 0x81, 0x82})
    fmt.Println(text)
    for _, c := range text {
        fmt.Printf("%#U (%v)\n", c, check(c))
    }
}
```

これを実行すると

```
$ go run sample1.go
﻿や     あ
U+FEFF (Unicode Control)
U+3084 'や' (Graphic)
U+0009 (Latin1 Control)
U+3042 'あ' (Graphic)
```

となった。うんうん。

## 結合子と異体字セレクタ

上述の check() 関数を使って，今度は絵文字の中身を見てみる。

```go:sample2.go
func main() {
    text := "🙇‍♂️"
    for _, c := range text {
        fmt.Printf("%#U (%v)\n", c, check(c))
    }
}
```

これを実行すると

```
$ go run sample2.go
U+1F647 '🙇' (Graphic)
U+200D (Unicode Control)
U+2642 '♂' (Graphic)
U+FE0F '️' (Graphic)
```

となった。ありゃ。 ZWJ はともかく異体字セレクタって図形文字あつかいなんだ。しかし，これでは大雑把すぎるので check() 関数を少し弄って...

```go:sample2.go
func check(r rune) string {
    switch {
    case unicode.Is(unicode.Sc, r):
        return "Symbol/currency"
    case unicode.Is(unicode.Sk, r):
        return "Symbol/modifier"
    case unicode.Is(unicode.Sm, r):
        return "Symbol/math"
    case unicode.Is(unicode.So, r):
        return "Symbol/other"
    case unicode.Is(unicode.Variation_Selector, r):
        return "Variation Selector"
    case unicode.Is(unicode.Join_Control, r):
        return "Join Control"
    case unicode.IsGraphic(r):
        return "Graphic"
    case unicode.IsControl(r):
        return "Latin1 Control"
    case unicode.Is(unicode.C, r):
        return "Unicode Control"
    }
    return "Unknown"
}
```

これを使ってもう一度実行してみると

```
$ go run sample2.go
U+1F647 '🙇' (Symbol/other)
U+200D (Join Control)
U+2642 '♂' (Symbol/other)
U+FE0F '️' (Variation Selector)
```

となった。なお，シンボルを区別しなくていいのなら [unicode].IsSymbol() 関数を使う手もある。

## 漢字と部首

Unicode って漢字の部首にもコードポイントが割り当てられているのよ。幸いなことに [unicode] 標準パッケージで部首を判別可能だ。先ほどの check() 関数に

```go:sample3.go
switch {
case unicode.Is(unicode.Radical, r):
    return "Radical"
case unicode.Is(unicode.Ideographic, r):
    return "Ideographic"
}
```

を加えればよい。これで

```go:sample4.go
func main() {
    text := "⽟玉"
    for _, c := range text {
        fmt.Printf("%#U (%v)\n", c, check(c))
    }
}
```

を実行すると

```
$ go run sample3.go
U+2F5F '⽟' (Radical)
U+7389 '玉' (Ideographic)
```

となった。なお，部首はシンボル扱いなので [unicode].IsSymbol() 関数でも一応は区別できる。

## 濁点とか

次は check() 関数を，以下のように，カナ文字を判別するよう書き換える。

```go:sample4.go
func check(r rune) string {
    switch {
    case unicode.Is(unicode.Katakana, r):
        return "Katakana"
    case unicode.Is(unicode.Hiragana, r):
        return "Hiragana"
    case unicode.Is(unicode.Lm, r):
        return "Letter/modifier"
    case unicode.Is(unicode.Lo, r):
        return "Letter"
    case unicode.Is(unicode.Mc, r):
        return "Mark/spacing combining"
    case unicode.Is(unicode.Me, r):
        return "Mark/enclosing"
    case unicode.Is(unicode.Mn, r):
        return "Mark/nonspacing"
    case unicode.IsSymbol(r):
        return "Symbol"
    case unicode.IsGraphic(r):
        return "Graphic"
    case unicode.IsControl(r):
        return "Latin1 Control"
    case unicode.Is(unicode.C, r):
        return "Unicode Control"
    }
    return "Unknown"
}
```

これで以下の文字列を調べてみる。

```go:sample4.go
func main() {
    text := "ペンギンペンギンﾍﾟﾝｷﾞﾝ"
    for _, c := range text {
        fmt.Printf("%#U (%v)\n", c, check(c))
    }
}
```

実行結果は以下の通り

```
$ go run sample4.go
U+30DA 'ペ' (Katakana)
U+30F3 'ン' (Katakana)
U+30AE 'ギ' (Katakana)
U+30F3 'ン' (Katakana)
U+30D8 'ヘ' (Katakana)
U+309A '゚' (Mark/nonspacing)
U+30F3 'ン' (Katakana)
U+30AD 'キ' (Katakana)
U+3099 '゙' (Mark/nonspacing)
U+30F3 'ン' (Katakana)
U+FF8D 'ﾍ' (Katakana)
U+FF9F 'ﾟ' (Letter/modifier)
U+FF9D 'ﾝ' (Katakana)
U+FF77 'ｷ' (Katakana)
U+FF9E 'ﾞ' (Letter/modifier)
U+FF9D 'ﾝ' (Katakana)
```

濁点や半濁点の文字種が全角と半角で異なっている点に注意。おそらく濁点等の判別に関しては専用の [unicode].RangeTable のテーブルを用意した方がいいと思う。

## Unicode はやっかい

ね。普通の日本語文字でこれだもの。ホンマやっかいだよ。

[unicode] 標準パッケージの定義済み [unicode].RangeTable テーブルはよくできてるし，ある程度日本語も考慮されているけど，細かい制御を行うのであれば用途に応じて専用の [unicode].RangeTable テーブルを用意したほうがいいだろう。量が多くて面倒くさいけどね。

## 参考リンク

- [その文字が JIS X 0208 に含まれるか？ あるいは unicode.RangeTable の使い方](https://zenn.dev/ikawaha/articles/20210116-ab1ac4a692ae8bb4d9cf)
- [かなカナ変換 | text.Baldanders.info](https://text.baldanders.info/golang/kana-conversion/)
- [こんな埼「玉」修正してやるぅ | text.Baldanders.info](https://text.baldanders.info/golang/unicode-kangxi-radical/)

[Go]: https://golang.org/ "The Go Programming Language"
[unicode]: https://golang.org/pkg/unicode/ "unicode - The Go Programming Language"

## 参考図書

https://www.amazon.co.jp/dp/4621300253
