---
title: "Go 言語で改行コードを変換する（正規表現以外の解）"
emoji: "😀"
type: "tech"
topics: [Go]
published: false
---
改行コード（LF, CR, CRLF）を変換する [Go 言語]のコードを考える。

真っ先に浮かぶのは [regexp] パケージを使って，たとえば

```go
package main

import (
	"fmt"
	"regexp"
)

var regxNewline = regexp.MustCompile(`\r\n|\r|\n`) //throw panic if fail

func convNewline(str, nlcode string) string {
	return regxNewline.Copy().ReplaceAllString(str, nlcode)
}

func main() {
	before := "あ\nい\rう\r\nえ"
	fmt.Printf("%U\n", []rune(before))

	after := convNewline(before, "\n")

	fmt.Printf("%U\n", []rune(after))
}
```

と書く[^r1]。これの実行結果は以下の通り。

[^r1]: `regexp.Regexp.ReplaceAllString()` メソッドを教えていただいた。感謝。なお Go 1.12 では複数の goroutine 下で `regexp.Regexp` インスタンスを使う際に `Copy()` メソッドでコピー・インスタンスを作らなくても処理がブロックされることはなくなった。

```
[U+3042 U+000A U+3044 U+000D U+3046 U+000D U+000A U+3048]
[U+3042 U+000A U+3044 U+000A U+3046 U+000A U+3048]
```

私も最初はこんな感じで書いていたのだが，**「[Go 言語]で正規表現を使ったら負け」**な気がして，何か方法はないかと [strings] パッケージをつらつら眺めてたら `strings.Replacer` 型が使えそうである。

たとえば，こんな感じで書ける。

```go
package main

import (
	"fmt"
	"strings"
)

func convNewline(str, nlcode string) string {
	return strings.NewReplacer(
		"\r\n", nlcode,
		"\r", nlcode,
		"\n", nlcode,
	).Replace(str)
}

func main() {
	before := "あ\nい\rう\r\nえ"
	fmt.Printf("%U\n", []rune(before))

	after := convNewline(before, "\n")

	fmt.Printf("%U\n", []rune(after))
}
```

これの実行結果は以下の通りで同じ結果が得られた。

```
[U+3042 U+000A U+3044 U+000D U+3046 U+000D U+000A U+3048]
[U+3042 U+000A U+3044 U+000A U+3046 U+000A U+3048]
```

あぁ，これで気持ちよく週末を過ごせる（笑）

[strings] パッケージと（今回は使わなかったが）[unicode] パッケージを組み合わせるとかなり色々できるので正規表現に手を出す前に検討してみるのもいいかもしれない。

[Go 言語]: https://golang.org/ "The Go Programming Language"
[regexp]: https://golang.org/pkg/regexp/ "regexp - The Go Programming Language"
[strings]: https://golang.org/pkg/strings/ "strings - The Go Programming Language"
[unicode]: https://golang.org/pkg/unicode/ "unicode - The Go Programming Language"

