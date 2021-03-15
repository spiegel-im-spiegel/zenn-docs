---
title: "配列とスライス" # 記事のタイトル
emoji: "💻" # アイキャッチとして使われる絵文字（1文字だけ）
type: "tech" # "tech" : 技術記事 / "idea" : アイデア記事
topics: ["go", "programming"] # タグ。["markdown", "rust", "aws"] のように指定する
published: true # 公開設定（true で公開）
---

他所様のブログ記事などを見るに [Go] の学習を始める際に躓きがちなのが interface と nil と slice なのではないかと思う。 [Interface](https://zenn.dev/spiegel/articles/20201129-interface-types-in-golang "Interface 型の使いどころ【Go】") と [nil](https://zenn.dev/spiegel/articles/20201010-ni-is-not-nil "nil == nil でないとき（または Go プログラマは息をするように依存を注入する）") については以前に書いた拙文を見ていただくとして，配列とスライスについては Zenn で書いてなかったな，と思い立ち記事にしてみることにした。なんちうあざとい（笑）

とはいえ，スライスは配列との関係が分かればさほど難しくない。以降からひとつずつ見ていくことにしよう。なお，記事中の図は “[Go Slices: usage and internals](https://blog.golang.org/slices-intro)” から拝借している。つか（英語不得手でないなら）そっちの記事を見た方が早いんだけどね。


## 配列（Array）

まずは配列について。

[Go] における「配列」は基本型のひとつで，単一の型のゼロ個以上の要素で構成されている。コードで書くとこんな感じ。

```go:sample1.go
// +build run

package main

import "fmt"

func main() {
    ary := [4]int{1, 2, 3, 4}
    fmt.Printf("Type: %[1]T , Value: %[1]v\n", ary)
    // Output:
    // Type: [4]int , Value: [1 2 3 4]
}
```

変数 `ary` を図で表すならこんな感じ。

![](https://blog.golang.org/slices-intro/slice-array.png)
*via “[Go Slices: usage and internals - The Go Blog](https://blog.golang.org/slices-intro)”*

ポイントは型名が `[4]int` の固定長データである点。要素の型や数が異なれば異なる型として扱われる。

また，配列は「値」である。つまり，同じ型であれば `==` 演算子で同値性[^eq1]（equality）の評価ができる（ただし，異なる型同士は評価できない。また要素の型が比較可能でない場合も評価できない）。

[^eq1]: 私は演算子における「等価」と「等値」の宗教論争に巻き込まれたくないので，意図的に “equality” を「同値性」と呼んでいる。ごめんペコン。

```go:sample2.go
func main() {
    ary1 := [4]int{1, 2, 3, 4}
    ary2 := [4]int{1, 2, 3, 4}
    ary3 := [4]int{2, 3, 4, 5}
    ary4 := [4]int64{1, 2, 3, 4}

    fmt.Printf("ary1 == ary2: %v\n", ary1 == ary2) // ary1 == ary2: true
    fmt.Printf("ary1 == ary3: %v\n", ary1 == ary3) // ary1 == ary3: false
    fmt.Printf("ary1 == ary4: %v\n", ary1 == ary4) // invalid operation: ary1 == ary4 (mismatched types [4]int and [4]int64)
}
```

さらに，配列が「値」であるということは `=` 等による代入構文[^stmt1] で要素も含めてインスタンスのコピーが発生する。関数の引数に配列を指定した場合も同様にコピーが渡される。

[^stmt1]: [Go] では代入は式（expression）ではなく文（statement）として機能する。式と文の違いは，文は評価結果を値として持たず，式の一部として組み込むことができないことである。

```go:sample3a.go
func displayArray4Int(ary [4]int) {
    fmt.Printf("Pointer: %p , Value: %v\n", &ary, ary)
}

func main() {
    ary1 := [4]int{1, 2, 3, 4}
    ary2 := ary1

    fmt.Printf("Pointer: %p , Value: %v\n", &ary1, ary1)
    fmt.Printf("Pointer: %p , Value: %v\n", &ary2, ary2)
    displayArray4Int(ary1)
    // Output:
    // Pointer: 0xc0000141a0 , Value: [1 2 3 4]
    // Pointer: 0xc0000141c0 , Value: [1 2 3 4]
    // Pointer: 0xc000014240 , Value: [1 2 3 4]
}
```

関数にインスタンス自体を渡したいのであればポインタ値を渡せばよい。

```go:sample3b.go
func referArray4Int(ary *[4]int) {
    fmt.Printf("Pointer: %p , Value: %v\n", ary, ary)
}

func main() {
    ary1 := [4]int{1, 2, 3, 4}

    fmt.Printf("Pointer: %p , Value: %v\n", &ary1, ary1)
    referArray4Int(&ary1)
    // Output:
    // Pointer: 0xc0000141a0 , Value: [1 2 3 4]
    // Pointer: 0xc0000141a0 , Value: &[1 2 3 4]
}
```

ここまでは OK かな。

## スライス（Slice）

スライスをコードで書くとこんな感じになる[^byte1]。

[^byte1]: byte 型は uint8 型の別名定義である。

```go:sample4.go
func main() {
    slc1 := []byte{0, 1, 2, 3, 4}
    fmt.Printf("Type: %[1]T , Value: %[1]v\n", slc1)
    // Output:
    // Type: []uint8 , Value: [0 1 2 3 4]
}
```

配列との記述上の違いは型名の角括弧（bracket）の中に要素数を指定するか否かだ。

配列はスライスに変換することができる。こんな感じ。

```go:sample5.go
func main() {
    ary1 := [5]byte{0, 1, 2, 3, 4}
    slc1 := ary1[:]
    fmt.Printf("Pointer: %p , Refer: %p , Value: %v\n", &ary1, &ary1[0], ary1)
    fmt.Printf("Pointer: %p , Refer: %p , Value: %v\n", &slc1, &slc1[0], slc1)
    // Output:
    // Pointer: 0xc000012088 , Refer: 0xc000012088 , Value: [0 1 2 3 4]
    // Pointer: 0xc000004078 , Refer: 0xc000012088 , Value: [0 1 2 3 4]
}
```

変数 `ary1` と `slc1` について `&x` と `&x[0]` のポインタ値の違いに注目してほしい。異なる変数なのだから変数のポインタ値が異なるのは当然だが，要素のポインタは同じ値になっている。つまりスライスの中身は代入した配列と「同一」なのである。

実はスライスの実体は

- 参照する配列へのポインタ値
- サイズ
- 容量

の3つの状態を属性として持つオブジェクトである。図にするとこんな感じ。

![](https://blog.golang.org/slices-intro/slice-struct.png)
*via “[Go Slices: usage and internals - The Go Blog](https://blog.golang.org/slices-intro)”*

ここで

```go
slc1 := ary1[:]
```

は以下のように図示できる。

![](https://blog.golang.org/slices-intro/slice-1.png)
*via “[Go Slices: usage and internals - The Go Blog](https://blog.golang.org/slices-intro)”*

スライスを使えば配列（またはスライス）の一部を切り出すことができる。たとえば

```go
slc2 := ary1[2:4]
```

とすると

![](https://blog.golang.org/slices-intro/slice-2.png)
*via “[Go Slices: usage and internals - The Go Blog](https://blog.golang.org/slices-intro)”*

という感じに切り出される（元の配列が切り詰められているわけではないので注意）。さらにこの `slc2` に対して

```go
slc3 := sl2[:cap(slc2)]
```

とすると

![](https://blog.golang.org/slices-intro/slice-3.png)
*via “[Go Slices: usage and internals - The Go Blog](https://blog.golang.org/slices-intro)”*

という感じに取り出せる。

```go:sample6.go
func main() {
    ary1 := [5]byte{0, 1, 2, 3, 4}
    slc1 := ary1[:]
    slc2 := ary1[2:4]
    slc3 := slc2[:cap(slc2)]
    fmt.Printf("Refer: %p , Len: %d , Cap: %d , Value: %v\n", &ary1[0], len(ary1), cap(ary1), ary1)
    fmt.Printf("Refer: %p , Len: %d , Cap: %d , Value: %v\n", &slc1[0], len(slc1), cap(slc1), slc1)
    fmt.Printf("Refer: %p , Len: %d , Cap: %d , Value: %v\n", &slc2[0], len(slc2), cap(slc2), slc2)
    fmt.Printf("Refer: %p , Len: %d , Cap: %d , Value: %v\n", &slc3[0], len(slc3), cap(slc3), slc3)
    // Output:
    // Refer: 0xc000012088 , Len: 5 , Cap: 5 , Value: [0 1 2 3 4]
    // Refer: 0xc000012088 , Len: 5 , Cap: 5 , Value: [0 1 2 3 4]
    // Refer: 0xc00001208a , Len: 2 , Cap: 3 , Value: [2 3]
    // Refer: 0xc00001208a , Len: 3 , Cap: 3 , Value: [2 3 4]
}
```

なお `ary[low:high]` とした場合

$$
0 \le \mathrm{low} \le \mathrm{high} \le \mathrm{len(ary)}
$$

となっていなければならない。なお $\mathrm{low} = 0$ または $\mathrm{high} = \mathrm{len(ary)}$ の場合は $\mathrm{low}$ または $\mathrm{high}$ の指定を省略できる。つまり

```go
slc1 := ary1[:]
```

は

```go
slc1 := ary1[0:len(ary1)]
```

と等価である。

あるいは容量の指定も含めて `slc[low:high:max]` と書くこともできる。
この場合 $\mathrm{max}$ は容量を指定するもので

$$
0 \le \mathrm{low} \le \mathrm{high} \le \mathrm{max} \le \mathrm{cap(ary)}
$$

を満たしていればよい。

## スライスは参照であり値である

これまでの説明から分かるようにスライスは配列の「参照」のようにふるまう。「ふるまう」とはどういうことか，もう少し詳しく見てみよう。

```go:sample7.go
func displaySliceByte(slc []byte) {
    fmt.Printf("Pointer: %p , Refer: %p , Value: %v\n", &slc, &slc[0], slc)
}

func main() {
    ary1 := [5]byte{0, 1, 2, 3, 4}
    slc1 := ary1[:]
    fmt.Printf("Pointer: %p , Refer: %p , Value: %v\n", &ary1, &ary1[0], ary1)
    fmt.Printf("Pointer: %p , Refer: %p , Value: %v\n", &slc1, &slc1[0], slc1)
    displaySliceByte(slc1)
    // Output:
    // Pointer: 0xc000102058 , Refer: 0xc000102058 , Value: [0 1 2 3 4]
    // Pointer: 0xc000100048 , Refer: 0xc000102058 , Value: [0 1 2 3 4]
    // Pointer: 0xc000100078 , Refer: 0xc000102058 , Value: [0 1 2 3 4]
}
```

まずは3つの配列・スライスは全て同一の配列を指している点に注目。そして `displaySliceByte()` 関数の引数として渡したスライスと渡す前のスライスは異なるインスタンス（つまり値渡し）であることにも注目してほしい。

このようにスライスは「配列への参照のようにふるまう」だけで（Java 等で言うところの）本当の意味での「参照」ではないということだ。

おそらく Java のような「参照」が言語仕様として組み込まれている言語圏から来た人はここで混乱するんじゃないだろうか。「[Go] に（本当の）参照はない」という点は心に刻み込んだ方がいいだろう。[^ref1]

[^ref1]: 他に [Go] で「参照のようにふるまう」型としてはチャネル，インタフェース，関数，マップがある。スライスも含めてこれらの型はゼロ値が nil になっている。

この参照と値のギャップが最も分かりやすく出るのが `append()` 関数だろう[^cap1]。

[^cap1]: スライスを容量を指定して生成する場合は `slc := make([]int, 0, 5)` などとすればよい。ただし `make()` 関数はインスタンスを必ずヒープ上に生成する。まぁ `append()` 関数でバッファを再取得する場合も結局ヒープになるのだが。

```go:sampe8.go
func main() {
    var slc []int
    fmt.Printf("Pointer: %p , <ZERO value>\n", &slc)
    for i := 0; i < 5; i++ {
        slc = append(slc, i)
        fmt.Printf("Pointer: %p , Refer: %p , Value: %v (%d)\n", &slc, &slc[0], slc, cap(slc))
    }
    // Output:
    // Pointer: 0xc000004078 , <ZERO value>
    // Pointer: 0xc000004078 , Refer: 0xc000012088 , Value: [0] (1)
    // Pointer: 0xc000004078 , Refer: 0xc0000120d0 , Value: [0 1] (2)
    // Pointer: 0xc000004078 , Refer: 0xc0000141c0 , Value: [0 1 2] (4)
    // Pointer: 0xc000004078 , Refer: 0xc0000141c0 , Value: [0 1 2 3] (4)
    // Pointer: 0xc000004078 , Refer: 0xc00000e340 , Value: [0 1 2 3 4] (8)
}
```

`append()` 関数は引数に渡されたスライスに要素を追加する組み込み関数だが，引数として渡される `slc` は単なる「値」なので，関数実行後の〈ポインタ値，サイズ，容量〉の状態をスライスのインスタンスとして返却している。一方 `append()` 関数を呼び出した側は返却値で元のスライスの状態を上書きしているわけだ。

## スライスは複製も比較もできない

配列は値なので，基本的に比較可能だし，代入時にはコピーが作成される。しかしスライスでは `=` 等の代入構文を使っても内容の複製はされない。スライスの複製が必要であれば `copy()` 関数を使う。

```go:sampe9.go
func main() {
    slc1 := []int{0, 1, 2, 3, 4}
    slc2 := slc1
    slc3 := make([]int, len(slc1), cap(slc1))
    copy(slc3, slc1)
    fmt.Printf("Pointer: %p , Refer: %p , Value: %v\n", &slc1, &slc1[0], slc1)
    fmt.Printf("Pointer: %p , Refer: %p , Value: %v\n", &slc2, &slc2[0], slc2)
    fmt.Printf("Pointer: %p , Refer: %p , Value: %v\n", &slc3, &slc3[0], slc3)
    // Output:
    // Pointer: 0xc000004078 , Refer: 0xc00000c2a0 , Value: [0 1 2 3 4]
    // Pointer: 0xc000004090 , Refer: 0xc00000c2a0 , Value: [0 1 2 3 4]
    // Pointer: 0xc0000040a8 , Refer: 0xc00000c2d0 , Value: [0 1 2 3 4]
}
```

スライスを「代入」しても〈ポインタ値，サイズ，容量〉の状態がコピーされるだけなので，まぁ当然だろう。また `copy()` 関数を使う場合はコピー先のインスタンスのサイズと容量をあらかじめ合わせておく必要がある。

さらに，スライスは `==` 演算子による比較もできない（型によらずコンパイルエラーになる。ただし nil との比較は可能）。

```go:sample10a.go
func main() {
    slc1 := []int{0, 1, 2, 3, 4}
    slc2 := []int{0, 1, 2, 3, 4}
    fmt.Printf("slc1 == slc2: %v\n", slc1 == slc2) // invalid operation: slc1 == slc2 (slice can only be compared to nil)
}
```

同じ型のスライス同士で内容の比較がしたいのであれば，たとえば `reflect.DeepEqual()` 関数が使える。

```go:sample10b.go
func main() {
    slc1 := []int{0, 1, 2, 3, 4}
    slc2 := []int{0, 1, 2, 3, 4}
    if reflect.DeepEqual(slc1, slc2) {
        fmt.Println("slc1 == slc2: true")
    } else {
        fmt.Println("slc1 == slc2: false")
    }
    // Output
    // slc1 == slc2: true
}
```

## というわけで

配列とスライスの関係を頭に入れて上手く使い分ければ（C/C++ の配列などに比べれば）簡単に安全にこれらを扱うことができるだろう。色々と試して欲しい。

## 参考

https://text.baldanders.info/golang/array-and-slice/
https://slide.baldanders.info/shimane-go-2020-02-13/

[Go]: https://golang.org/ "The Go Programming Language"
<!-- eof -->
