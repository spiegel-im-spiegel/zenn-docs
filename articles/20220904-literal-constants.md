---
title: "リテラル定数" # 記事のタイトル
emoji: "💻" # アイキャッチとして使われる絵文字（1文字だけ）
type: "tech" # "tech" : 技術記事 / "idea" : アイデア記事
topics: ["go", "programming"] # タグ。["markdown", "rust", "aws"] のように指定する
published: true # 公開設定（true で公開）
---

Qiita の「[整数型に 'c' が代入できるそれぞれの事情](https://qiita.com/Nabetani/items/0554d5b040d70525ec95)」はなかなか面白かった。

https://qiita.com/Nabetani/items/0554d5b040d70525ec95

特に最後の [zig](https://ziglang.org/ "Home ⚡ Zig Programming Language") は「へー」って感じである。

上の記事はおそらく各言語の比較を重視しているために意図的に解説を端折っていると思うが，折角なので [Go] の「リテラル定数」について少し書いてみる。

## Go の定数

[Go] の言語仕様では「定数」を以下のように説明している。

> There are boolean constants, rune constants, integer constants, floating-point constants, complex constants, and string constants. Rune, integer, floating-point, and complex constants are collectively called numeric constants.
*(via “[The Go Programming Language Specification](https://go.dev/ref/spec#Constants)”)*

さらに数値定数について

> Numeric constants represent exact values of arbitrary precision and do not overflow. Consequently, there are no constants denoting the IEEE-754 negative zero, infinity, and not-a-number values.
>
> Constants may be typed or untyped. Literal constants, true, false, iota, and certain constant expressions containing only untyped constant operands are untyped.
*(via “[The Go Programming Language Specification](https://go.dev/ref/spec#Constants)”)*

とある。たとえば標準の [math](https://pkg.go.dev/math) パッケージでは円周率 $\pi$ は

```go:math/const.go
// Mathematical constants.
const (
    Pi  = 3.14159265358979323846264338327950288419716939937510582097494459 // https://oeis.org/A000796
)
```

と基本型の float64 を大きく越える精度で定義されている。

型付けなし定数（untyped constant）は変数宣言または変数への代入時に型が決定される（実際にはコンパイル時に評価される）。

> An untyped constant has a default type which is the type to which the constant is implicitly converted in contexts where a typed value is required, for instance, in a short variable declaration such as i := 0 where there is no explicit type. The default type of an untyped constant is bool, rune, int, float64, complex128 or string respectively, depending on whether it is a boolean, rune, integer, floating-point, complex, or string constant.
*(via “[The Go Programming Language Specification](https://go.dev/ref/spec#Constants)”)*

あとでまた説明するが，シングルクォートで囲まれた文字 `'c'` は rune リテラルと呼ばれ rune 定数を表している。なので短縮形の変数宣言

```go
r := 'c'
```

で宣言された変数 `r` は rune 型で初期値 `'c' (U+0063)` を与えられる，というわけ。さらに `'c'` は型付けなしの数値定数でもあるので，明示的に

```go
var r int = 'c'
```

で宣言された変数 `r` は int 型で初期値 `0x63` を与えられる。同じ `c` でもダブルクォートで囲まれる `"c"` は文字列リテラルなので

```go
var r int = "c" // cannot use "c" (untyped string constant) as int value in variable declaration
```

はコンパイルエラーになる。ちなみに string と rune 配列は相互変換できるので，明示的に

```go
package main

import "fmt"

func main() {
    var r rune = []rune("c")[0]
    fmt.Printf("%#U", r) // U+0063 'c'
}
```

とすればコンパイルエラーにならない（笑）

## 定数のリテラル表現

定数のリテラル表現には以下の5つがある。

1. [整数リテラル（integer literal）](https://go.dev/ref/spec#Integer_literals)
2. [浮動小数点数リテラル（floating-point literal）](https://go.dev/ref/spec#Floating-point_literals)
3. [虚数リテラル（imaginary literal）](https://go.dev/ref/spec#Imaginary_literals)
4. [rune リテラル（rune literal）](https://go.dev/ref/spec#Rune_literals)
5. [文字列リテラル（string literal）](https://go.dev/ref/spec#String_literals)

以降でひとつずつ見ていこう。

### 整数リテラル

整数リテラルは整数定数を表すリテラル表現で以下の形式で書ける。

```go
42
4_2
0600
0_600
0b0110
0B0110_0011
0o600
0O600       // second character is capital letter 'O'
0xBadFace
0xBad_Face
0x_67_7a_2f_cc_40_c6
```

ちなみに `_` は値として意味を持つものではないが，桁の区切りとして任意の場所に差し込むことができる。

整数リテラルは型付けなしの数値定数としても振る舞うので

```go
package main

import "fmt"

func main() {
    n := 0b0110_0011
    fmt.Printf("%d\n", n) // 99

    var f float64 = 0b0110_0011
    fmt.Printf("%f\n", f) // 99.000000

    var r rune = 0b0110_0011
    fmt.Printf("%#U\n", r) // U+0063 'c'
}
```

などと書くこともできる。

### 浮動小数点数リテラル

浮動小数点数リテラルは浮動小数点数定数を表すリテラル表現で以下の形式で書ける。

```go
0.
72.40
072.40       // == 72.40
2.71828
1.e+0
6.67428e-11
1E6
.25
.12345E+5
1_5.         // == 15.0
0.15e+0_2    // == 15.0

0x1p-2       // == 0.25
0x2.p10      // == 2048.0
0x1.Fp+0     // == 1.9375
0X.8p-0      // == 0.5
0X_1FFFP-16  // == 0.1249847412109375
```

浮動小数点数（IEEE 754）の内部表現で書けるのは凄いと思うが，まず使わないよね（笑）

小数点以下が0であれば整数型（byte や rune を含む）変数宣言時の初期化にも使える。

```go
package main

import "fmt"

func main() {
    var n int = 99.0
    fmt.Printf("%d\n", n) // 99

    var r rune = 99.0
    fmt.Printf("%#U\n", r) // U+0063 'c'
}
```

しかし，小数点以下が0でない場合はコンパイルエラーになる。

```go
var n int = 4.56 // cannot use 4.56 (untyped float constant) as int value in variable declaration (truncated)
```

この場合はリテラル値を一度変数に落とし込めば整数型にキャストできる（小数点以下切り捨て）。

```go
package main

import "fmt"

func main() {
    f := 4.56
    var n int = int(f)
    fmt.Printf("%d\n", n) // 4
}
```

### 虚数リテラル

虚数リテラルは複素数定数の虚数部を表すリテラル表現で以下の形式で書ける。

```go
0i
0123i         // == 123i for backward-compatibility
0o123i        // == 0o123 * 1i == 83i
0xabci        // == 0xabc * 1i == 2748i
0.i
2.71828i
1.e+0i
6.67428e-11i
1E6i
.25i
.12345E+5i
0x1p-2i       // == 0x1p-2 * 1i == 0.25i
```

虚数リテラルを使って複素数定数を以下のように表すことができる。

```go
package main

import "fmt"

func main() {
    c := 1.23 + 4.56i
    fmt.Printf("%f\n", c) // (1.230000+4.560000i)
}
```

当然ながら虚数部が0でない複素数定数を使って整数や浮動小数点数の変数宣言の初期化には使えない。

```go
var f float64 = 1.23 + 4.56i // cannot use 1.23 + 4.56i (untyped complex constant (1.23 + 4.56i)) as float64 value in variable declaration (overflows)
```

もちろん変数に落とし込んでからのキャストもできない（そもそも複素数型から浮動小数点数型への変換ができないので）。

```go
c := 1.23 + 4.56i
var f = float64(c) // cannot convert c (variable of type complex128) to type float64
```

複素数型から実数部・虚数部を取り出すには組み込みの real() または imag() 関数を使う。

```go
package main

import "fmt"

func main() {
    c := 1.23 + 4.56i
    fmt.Printf("real part: %f\n", real(c))      // real part: 1.230000
    fmt.Printf("imaginary part: %f\n", imag(c)) // imaginary part: 4.560000
}
```

複素数型の操作は標準の [math/cmplx](https://pkg.go.dev/math/cmplx "cmplx package - math/cmplx - Go Packages") パッケージを使うといいだろう。

### Rune リテラル

Rune 型は Unicode コードポイントを内部表現として持つ整数型である。 Rune リテラルは rune 定数を表すリテラル表現で以下の形式で書ける。

```go
'a'
'ä'
'本'
'\t'
'\000'
'\007'
'\377'
'\x07'
'\xff'
'\u12e4'
'\U00101234'
```

単一のコードポイントであることがポイント（駄洒落）で，たとえば複数のコードポイントで構成される絵文字を rune リテラルで表そうとしてもコンパイルエラーになる。

```go
const emoji = '👍🏼' // more than one character in rune literal
```

どうやら整数として評価できれば値そのものには頓着しないようで

```go
package main

import "fmt"

func main() {
    var r rune = -1.0
    fmt.Printf("%#U\n", r) // U+FFFFFFFFFFFFFFFF
}
```

みたいな恐ろしい記述も可能らしい（あれ？ Rune って int32 の alias じゃなかったっけ？ まぁいいや）。

### 文字列リテラル

文字列リテラルは文字列定数を表すリテラル表現で以下の形式で書ける。

```go
`abc`                // same as "abc"
`\n
\n`                  // same as "\\n\n\\n"
"\n"
"\""                 // same as `"`
"Hello, world!\n"
"日本語"
"\u65e5本\U00008a9e"
"\xff\u00FF"
```

文字列リテラルは UTF-8 エンコーディングであることを前提としている。また文字列リテラルで表現された文字列定数は string 型にしか使えず不変（immutable）な値として振る舞う。たとえば

```go
var b []byte = "日本語" // cannot use "日本語" (untyped string constant) as []byte value in variable declaration
var r []rune = "日本語" // cannot use "日本語" (untyped string constant) as []rune value in variable declaration
```

のように byte 配列宣言や rune 配列宣言の初期値として使うことはできない。この場合は明示的に

```go
package main

import "fmt"

func main() {
    var b []byte = []byte("日本語")
    fmt.Println(b) // [230 151 165 230 156 172 232 170 158]
    var r []rune = []rune("日本語")
    fmt.Println(r) // [26085 26412 35486]
}
```

などと変換すれば問題なく使える。あるいは for-range 構文を使って

```go
package main

import "fmt"

func main() {
    for _, r := range "日本語" {
        fmt.Printf("%#U\n", r)
    }
    // Output:
    // U+65E5 '日'
    // U+672C '本'
    // U+8A9E '語'
}
```

などと rune 単位で取り出すこともできる[^r1]。

[^r1]: For-range 構文で取り出す値は配列への参照ではなくコピー値なので注意。

## 参考

https://zenn.dev/spiegel/articles/20210813-untyped-constant
https://zenn.dev/spiegel/articles/20211004-pointer-to-literal-value

[Go]: https://go.dev/ "The Go Programming Language"
<!-- eof -->
