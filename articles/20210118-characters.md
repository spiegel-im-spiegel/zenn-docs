---
title: "やっかいな日本語" # 記事のタイトル
emoji: "🤔" # アイキャッチとして使われる絵文字（1文字だけ）
type: "idea" # "tech" : 技術記事 / "idea" : アイデア記事
topics: ["engineering", "unicode"] # タグ。["markdown", "rust", "aws"] のように指定する
published: true # 公開設定（true で公開）
---

Go 実装関連で個人的に注目している [ikawaha](https://zenn.dev/ikawaha) さんが面白いパッケージを公開されている。

https://zenn.dev/ikawaha/articles/20210116-ab1ac4a692ae8bb4d9cf
https://github.com/ikawaha/encoding

これを使えば

```go
package main

import (
    "fmt"
    "unicode"

    "github.com/ikawaha/encoding/jisx0208"
)

func main() {
    for _, c := range "１二③Ⅳ" {
        fmt.Printf("%#U %v a JIS X 0208 character\n", c, func() string {
            if unicode.Is(jisx0208.RangeTable, c) {
                return "is"
            }
            return "is not"
        }())
    }
    // Outpu:
    // U+FF11 '１' is a JIS X 0208 character
    // U+4E8C '二' is a JIS X 0208 character
    // U+2462 '③' is not a JIS X 0208 character
    // U+2163 'Ⅳ' is not a JIS X 0208 character
}
```

という感じに 指定した文字が JIS X 0208 文字集合に含まれるか否かの判定ができる。素晴らしい！

このパッケージに敬意を表して，おぢさんがとりとめない昔話をしようじゃないか（笑）

## 2つの「文字コード」

いわゆる「文字コード」と呼ばれるものは以下の2つの意味が混濁していることが多い。

- 符号化文字集合
- 文字エンコーディング

7ビット空間の IRV[^irv1] や8ビット空間の JIS X 0201 などは両者が一体化しているので区別することにあまり意味はないが， JIS X 0213 や Unicode のような大きな文字集合や ISO 2022 のように複数の文字集合を出し入れできる文字エンコーディングでは区別して議論する必要がある。

[^irv1]: 昔で言うところの US-ASCII のこと。現在では ISO 646 の IRV (International Reference Version; 国際基準版) として定義し直されている。

文字集合については次節から紹介するとして，文字エンコードディングでよく知られているのは Shift-JIS や EUC (Extended UNIX Code) や UTF (Unicode Transformation Format) といったあたりか。

### Shift-JIS

Shift-JIS は JIS 文字集合用の文字エンコーディング。仕組みが単純で CUI 端末全盛時代には重宝されたが拡張性に乏しく，もはや時代遅れといっていい。

### EUC/EUC-JP

EUC は ISO 2022 実装のひとつで複数の文字集合に対応可能である。更に JIS 文字集合を使った EUC サブセットを EUC-JP と呼ぶ。 Shift-JIS と EUC-JP は元となる文字集合が同じなのでお互いに変換可能である。

### ISO-2022-JP

[RFC 1468] で規格化されている文字エンコーディング。いわゆる「JIS コード」と呼ぶ場合は大抵この文字エンコードディングを指す。もともとは7ビット伝送の電子メール配信で日本語を通すために考えられたようだ。名前に反して，厳密には ISO 2022 に準拠していない。

JIS X 0201 は許容するが，いわゆる半角カナは許容しない。また漢字についても JIS X 0208 までしか許容しない。ぶっちゃけ廃れた規格である（笑）

### UTF/UTF-8

UTF は Unicode 用の文字エンコーディングである。特に Unicode をオクテット単位のデータ・ストリームに変換する UTF-8 はインターネットとも相性がよく，情報交換用の文字エンコーディングとしては事実上の世界標準と言っていいだろう。

Unicode は JIS X 0201 や JIS X 0213 を取り込んでいる形になっているので，少なくとも JIS 系文字エンコーディングから Unicode 系エンコーディングへの変換は比較的容易である。



## JIS X 0201 と IRV

JIS X 0201 は8ビットサイズの文字集合で，以下の文字種を含んでいる。

- ラテン文字
- 片仮名

あるいは「半角英数」「半角カナ」と言った方が通りがいいだろうか。

厄介なことに JIS ラテン文字は IRV と2文字違う。

| コード | IRV                   | ラテン文字       |
|:------:| --------------------- | ---------------- |
|  0x5C  | `\` (REVERSE SOLIDUS) | `¥` (YEN SIGN)   |
|  0x7E  | `~` (TILDA)           | `￣` (OVER LINE) |

しかも 0x7E は環境によって実装が異なっていて，たとえば Shift-JIS ベースの日本語 Windows では TILDA なのに昔の UNIX 系システムでは OVER LINE だったりするんだよねぇ。まぁ，どちらも制御記号として使われることが多いので文字が異なっていても大したインパクトはないけどね。

それに対して 0x5C は JIS ラテン文字では通貨記号を割り当てている一方で IRV の REVERSE SOLIDUS はエスケープシーケンスの制御文字として多用されるのでインパクトが大きい。ちなみに Unicode で YEN SIGN のコードポイントは U+00A5 である。


## JIS 漢字の変遷

JIS X 0201 を除くいわゆる「JIS 漢字」には大きく分けて以下の3つのバージョンがある。

### JIS C 6226

私のようなロートル世代は JIS C 6226 (または JIS X 0208-1978) を「旧 JIS」と呼ぶ。 Windows の台頭により，1990年代にはパソコン等においてはほぼ駆逐されたが，汎用機などではしぶとく生き残っていることがあり，ビックリする。

### JIS X 0208 および JIS X 0212

同じく古い世代の間では JIS X 0208 を（「旧 JIS」に対する）「新 JIS」と呼ぶことがある。「JIS 漢字第一水準」「JIS 漢字第二水準」を含む文字集合で，94区×94点のコード空間を持っている。

JIS X 0212 は「補助漢字」と呼ばれ JIS X 0208 にない文字を補うものだったが，実装例が少なくあまり普及しなかった。

### JIS X 0213

JIS X 0213 は Unicode を意識した文字集合で，JIS X 0208 および JIS X 0212 を含んでいる。更にこれまで「システム外字」とか「機種依存文字」とか呼ばれていた文字も取り込んでいる。「JIS 漢字第一水準」から「JIS 漢字第四水準」までを網羅し，2面×94区×94点のコード空間を持っている。

### やっかいな JIS 漢字

JIS 漢字は，厄介なことに，漢字だけでなくラテン文字や片仮名も含んでいる。こちらは「全角英数」とか「全角カナ」でおなじみだろう。同じ文字で違うコードが割り当てられているわけだ。

また JIS X 0208 までは異体字をできるだけ排除してひとつの字形に包摂する方針だったのに対し， JIS X 0213 では一転して多くの異体字を許容し取り込むようになった。機種依存文字にかなりの異体字が含まれているので，機種依存文字を取り込むなら異体字もある程度許容せざるを得ないのだろう。

更に更に JIS X 0213 は2004年版で字形に関する大改訂を行ったため，OS や使用するフォントによって同じ UTF コードでも表示される字形が異なってしまい，これによる混乱がしばらくあった。最近のフォントは JIS X 0213:2004 に対応している筈だが，逆に古い文書は，フォントの埋め込み等をしない限り「文字化け」することになる。

## Unicode でええぢゃん？

Unicode は元々米国企業が中心になって作った文字集合で[^ucs]，簡単に言うと「各国の文字を全部ひとつのコード体系に組み込んでしまえばいいじゃない？」という雑な設計で成り立っている。それでもあらゆる文字をひとつのコード体系で表せるのは相当なメリットだし，拡張の余地が大きく残されているのも魅力的である。

[^ucs]: Unicode とは別に ISO 10646 という国際規格もある。両者の関係は微妙だが，現在では Unicode で決めた規格を ISO 10646 が追認するという形をとっているらしい。

しかし，各国における符号化文字集合の問題（たとえば日本なら全角・半角文字や異体字など）もそのまま引き継いでしまっているため，かなり厄介な状況になっているのも確かである。

以下にいくつか挙げてみよう。

### 後方互換性

たとえば Go では UTF-8 と Shift-JIS や EUC-JP との間で相互変換する[パッケージが存在する](https://pkg.go.dev/golang.org/x/text/encoding/japanese "japanese · pkg.go.dev")が，変換したコードが有効なコードか否かは別に検証する必要がある。

そこで最初に紹介した [ikawaha/encoding] パッケージの出番となるわけだ。また `¥` (YEN SIGN) のように単純変換できない場合もあるだろう。

おそらく仕事で使うなら，必要に応じて JIS C 6226 や JIS X 0213 の各文字集合についても検証や（場合によっては）変換の仕組みを用意する必要があるかもしれない。でも，そのためには結構デカい unicode.RangeTable を用意する（そしてテストする）必要がある。これが大変なんだよなぁ。

### CJK 統合問題

これは昔から有名な話だが，初期の Unicode を決めた人たちは同じ漢字でも国や地域によって字形や意味（ニュアンス）さえも異なるということが分からなかったらしい。というわけで，東アジアで使われる漢字を勝手に包摂（統合）してしまったのだ。今でも Unicode で書いたドキュメントを表示する際に国や言語ごとにフォントを正しく切り替えないと微妙に「文字化け」してしまう。

文字コードの問題に限らないが，こういう「国際規格」は国家としてきちんとプレゼンスを確保・維持していかないと，こんな風に「あとの祭り」になってしまう。

### Unicode 正規化

たとえば Unicode で「ペンギン」という文字列の表現には

- ペ：U+30DA
- ン：U+30F3
- ギ：U+30AE
- ン：U+30F3



または

- ペ：U+30D8 + U+309A
- ン：U+30F3
- ギ：U+30AD + U+3099
- ン：U+30F3

または

- ﾍﾟ：U+FF8D + U+FF9F
- ﾝ：U+FF9D
- ｷﾞ：U+FF77 + U+FF9E
- ﾝ：U+FF9D

の3種類（またはその組み合わせ）がある。ちなみに3番目はいわゆる「半角カナ」である。

これらの「3羽のペンギン」が等価であることを保証するために「Unicode 正規化」のルールが決められている。しかしこの正規化がまた微妙で，たとえば CJK 互換文字「神 (U+FA19)」を含む文字列を正規等価で評価しようとすると勝手に「神 (U+795E)」に変換され，しかも逆変換できなかったりする。詳しくは以下を参照のこと。

https://text.baldanders.info/golang/unicode-normalization/

こういう副作用もあるので， Unicode テキストを安直に丸ごと正規化するのは危険である。

### 文字の結合または異体字セレクタ

上述の

- ペ：U+30D8 + U+309A

の U+309A に相当する半濁点は「結合文字」と呼ばれていて，直前の基底文字と組み合わせてひとつの文字を構成している。

絵文字でも，たとえば 🙇‍♂️ のコードを拙作の [gnkf] で見ると

```
$ echo 🙇‍♂️ | gnkf dump --unicode
0x0001f647, 0x200d, 0x2642, 0xfe0f, 0x000a
```

のように複数のコードポイントで構成されているのが分かる。具体的には

| Unicode Point | 字形 | Unicode 名称          |
| ------------- | ---- | --------------------- |
| U+1F647       | 🙇   | PERSON BOWING DEEPLY  |
| U+200D        |      | ZERO WIDTH JOINER     |
| U+2642        | ♂    | MALE SIGN             |
| U+FE0F        |      | VARIATION SELECTOR-16 |

という内容になっていて，これで「土下座する男性」を表す絵文字（の異体字）になるらしい。

ちなみに，汎用の異体字セレクタは以下の256個が定義されているようだ。

- U+FE00 〜 U+FE0F: VARIATION SELECTOR-1〜16 (Standardized Variation Sequence; SVS)
- U+E0100 〜 U+E01EF: VARIATION SELECTOR-17〜256 (Ideographic Variation Sequence; IVS)

このうち IVS 用の異体字セレクタは漢字の異体字として使える。

厳密に Unicode 文字の操作を行うのであれば「ひとつの文字がひとつのコードポイントで表されるとは限らない」ことを前提に結合文字（結合子）や異体字セレクタを解釈しながら文字列を分析していくしかない。

## 最後に愚痴

UTF-8 の BOM はそろそろ滅びてくれんかねぇ...

## 参考リンク

- [特別編24 JIS X 0213の改正は、文字コードにどんな未来をもたらすか（7）　番外編：改正JIS X 0213とUnicodeの等価属性／正規化について（上）](https://internet.watch.impress.co.jp/www/column/ogata/sp24.htm)
    - [JIS X 0213の改正は、文字コードにどんな未来をもたらすか（8）　番外編：改正JIS X 0213とUnicodeの等価属性／正規化について（下）](https://internet.watch.impress.co.jp/www/column/ogata/sp25.htm)
- [「ユニコード」で予期せぬ目に遭った話 - moriyoshiの日記](http://moriyoshi.hatenablog.com/entry/2017/03/13/011005)
- [絵文字をひとつだけを取り出す技術](https://zenn.dev/catnose99/scraps/79012f2617ffd9)
- [Unicodeとの異体字バトルがはじまったぜ](https://zenn.dev/zetamatta/scraps/7737ec9a9a426f)
- [GitHub - spiegel-im-spiegel/charset_document: 「文字コードとその実装」 upLaTeX ドキュメント](https://github.com/spiegel-im-spiegel/charset_document) : 既に内容が古びているが，前半はまだ役に立つはず（笑）
- [かなカナ変換 | text.Baldanders.info](https://text.baldanders.info/golang/kana-conversion/)
- [こんな埼「玉」修正してやるぅ | text.Baldanders.info](https://text.baldanders.info/golang/unicode-kangxi-radical/)
- [絵文字と異体字と Markdown | text.Baldanders.info](https://text.baldanders.info/remark/2020/10/emoji-variation-and-markdown/)

[RFC 1468]: https://tools.ietf.org/html/rfc1468 "RFC 1468 - Japanese Character Encoding for Internet Messages"
[ikawaha/encoding]: https://github.com/ikawaha/encoding
[gnkf]: https://github.com/spiegel-im-spiegel/gnkf "spiegel-im-spiegel/gnkf: Network Kanji Filter by Golang"

## 参考図書

https://zenn.dev/zetamatta/books/b820d588f4856bcf836c
