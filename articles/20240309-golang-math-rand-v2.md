---
title: "Go 1.22 の math/rand/v2 を使ってみる" # 記事のタイトル
emoji: "💻" # アイキャッチとして使われる絵文字（1文字だけ）
type: "tech" # "tech" : 技術記事 / "idea" : アイデア記事
topics: ["go", "programming", "random"] # タグ。["markdown", "rust", "aws"] のように指定する
published: true # 公開設定（true で公開）
---

[Go] 1.22 の変更点はたくさんあるが，この記事では 疑似乱数パッケージ [math/rand] の新版である [math/rand/v2] を触ってみる。

## ChaCha8 疑似乱数生成器が標準になった

疑似乱数パッケージについて最大の変更は，ストリーム暗号 ChaCha8 の疑似乱数生成器が実装されランタイムに組み込まれたことだろう。 ChaCha20/8 については以下の記事を参考にするとよい。

- [The Salsa20 family of stream ciphers](https://cr.yp.to/snuffle/salsafamily-20071225.pdf)
- [C2SP/chacha8rand.md at main · C2SP/C2SP · GitHub](https://github.com/C2SP/C2SP/blob/main/chacha8rand.md)

ちなみに 20 とか 8 とかってのは疑似乱数を生成する際のラウンド数を指すものらしい。

で， [math/rand] および [math/rand/v2] パッケージのトップレベル関数群（[rand][math/rand/v2].IntN(), [rand][math/rand/v2].Float64() など）ではランタイムに組み込まれた ChaCha8 疑似乱数生成器を使っている。どういう実装になっているかは，拙文を参考にしていただければ。

https://text.baldanders.info/golang/pseudo-random-number-generator-v2/

疑似乱数を生成する際には seed を与える必要があるが，ランタイムに組み込んだ ChaCha8 の疑似乱数生成器ではアプリケーション初期化時に乱数デバイスを使って seed を与えるため，コード上で明示的に seed を与える必要はなくなった。

さっそく簡単なコードを書いてみよう。

```go:sample1.go
package main

import (
    "fmt"
    "math/rand/v2"
    "time"
)

func main() {
    for i := range 10 {
        fmt.Println(i+1, ":", rand.N(10*time.Minute))
    }
}
```

これを[実行](https://go.dev/play/p/Db4PYRoG6nM)すると以下のような出力になる。

```
1 : 2m26.946461444s
2 : 5m29.459083773s
3 : 4m24.754220634s
4 : 8m40.996019235s
5 : 2m33.695282473s
6 : 7m13.504302116s
7 : 2m11.664224654s
8 : 3m4.584656523s
9 : 6.659580705s
10 : 5m58.737834834s
```

明示的に seed を与えなくても毎回異なる値になることを確かめてほしい。以前に

https://zenn.dev/spiegel/articles/20211016-crypt-rand-as-a-math-rand

という記事を書いたが， ChaCha8 疑似乱数生成器は暗号技術的にセキュア（「予測困難性」要件を満たす）と言えるので，上の記事のような変換をしなくてもよくなった。なお [math/rand] パッケージのトップレベル関数群を使う場合は [rand][math/rand].Seed() メソッドを呼んでしまうと従来の（セキュアでない）疑似乱数生成器を使うので要注意である。つか，可能であれば [math/rand] は使わないほうがいいと思う。

[rand][math/rand/v2].N() は [math/rand/v2] パッケージで新たに追加された関数で Generics になっている。定義は以下の通り。

```go:math/rand/v2/rand.go
// N returns a pseudo-random number in the half-open interval [0,n) from the default Source.
// The type parameter Int can be any integer type.
// It panics if n <= 0.
func N[Int intType](n Int) Int {
    if n <= 0 {
        panic("invalid argument to N")
    }
    return Int(globalRand.uint64n(uint64(n)))
}

type intType interface {
    ~int | ~int8 | ~int16 | ~int32 | ~int64 |
        ~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64 | ~uintptr
}
```

この定義により int など整数の基本型を基底型（underlying type）とする任意の型を扱うことができるようになった。

## ChaCha8 疑似乱数生成器を Source として構成する

ChaCha8 疑似乱数生成器を [rand][math/rand/v2].Source として構成することもできる。こんな感じ。

```go:sample2.go
package main

import (
    crnd "crypto/rand"
    "fmt"
    "math/rand/v2"
)

func main() {
    var seed [32]byte
    _, _ = crnd.Read(seed[:]) //エラー処理をサボってます 🙇
    rnd := rand.New(rand.NewChaCha8(seed))
    for i := range 10 {
        fmt.Println(i+1, ":", rnd.IntN(1000))
    }
}
```

なお，この方法で作成した疑似乱数生成器（[rand][math/rand/v2].ChaCha8）は並行的に安全（concurrency safe）ではないため，並行処理下で使用する場合は注意が必要である。まぁ，流石に [panic を吐いてコケる](https://text.baldanders.info/golang/pseudo-random-number-generator/ "Go の疑似乱数生成器は並行的に安全ではないらしい")まではないみたいだけど。

ランタイムに組み込まれた ChaCha8 疑似乱数生成器はちゃんと排他処理を行っているので並行処理下で使いまくっても大丈夫である。 [math/rand] および [math/rand/v2] パッケージのトップレベル関数群も同様。

## PCG 疑似乱数生成器を Source として構成する

[math/rand/v2] では，もうひとつ疑似乱数生成器が用意されている。こんな感じ。

```go:sample3.go
package main

import (
    "fmt"
    "math/rand/v2"
)

func main() {
    rnd := rand.New(rand.NewPCG(rand.Uint64(), rand.Uint64()))
    for i := range 10 {
        fmt.Println(i+1, ":", rnd.IntN(1000))
    }
}
```

PCG (Permuted Congruential Generator) は線形合同法（LCG）のバリエーションだそうで LCG の統計学上の欠点を改善したものらしい。当然ながら暗号技術的にセキュアではない。また [rand][math/rand/v2].PCG は内部状態を持ち並行的に安全ではないため，こちらも並行処理下で使用する場合は注意が必要である。

## [math/rand] パッケージの疑似乱数生成器を [math/rand/v2] パッケージで使う

ところで [math/rand] パッケージの疑似乱数生成器は [math/rand/v2] パッケージでは使えないのだろうか。試してみよう。

```go:sample3.go
package main

import (
    "fmt"
    oldrand "math/rand"
    "math/rand/v2"
)

func main() {

    rnd := rand.New(oldrand.NewSource(rand.Int64()).(rand.Source))
    for i := range 10 {
        fmt.Println(i+1, ":", rnd.IntN(1000))
    }
}
```

おー。こいつ，[動くぞ](https://go.dev/play/p/b6Jf74YQNPU)。[math/rand] パッケージで定義されている疑似乱数生成器は[ラグ付フィボナッチ法（Lagged Fibonacci Generator）の一種](https://text.baldanders.info/golang/estimate-of-pi-4-prng/ "モンテカルロ法による円周率の推定（その4 PRNG）")だそうだ。これも LCG の改良版と言える。

## ベンチマーク

標準で用意されているアルゴリズムについては網羅できたと思うので，ベンチマークを取ってみよう。ベンチマーク用にこんな感じのコードを書いてみた。

```go:random_test.go
package randoms

import (
    crnd "crypto/rand"
    oldrand "math/rand"
    "math/rand/v2"
    "testing"
)

func makeSeedChaCha8() [32]byte {
    var seed [32]byte
    _, _ = crnd.Read(seed[:]) //エラー処理をサボってます 🙇
    return seed
}

func makeSeedPPCG() (uint64, uint64) {
    return rand.Uint64(), rand.Uint64()
}

func makeSeedLaggedFibonacci() int64 {
    return rand.Int64()
}

var seedChaCha8 = makeSeedChaCha8()
var seedPCG1, seedPCG2 = makeSeedPPCG()
var seedLaggedFibonacci = makeSeedLaggedFibonacci()
var count = 1000000

func BenchmarkRandomChaCha8(b *testing.B) {
    rnd := rand.New(rand.NewChaCha8(seedChaCha8))
    for i := 0; i < b.N; i++ {
        for range count {
            _ = rnd.IntN(1000)
        }
    }
}

func BenchmarkRandomChaCha8runtime(b *testing.B) {
    for i := 0; i < b.N; i++ {
        for range count {
            _ = rand.IntN(1000)
        }
    }
}

func BenchmarkRandomPCG(b *testing.B) {
    rnd := rand.New(rand.NewPCG(seedPCG1, seedPCG2))
    for i := 0; i < b.N; i++ {
        for range count {
            _ = rnd.IntN(1000)
        }
    }
}

func BenchmarkRandomLaggedFibonacci(b *testing.B) {
    rnd := rand.New(oldrand.NewSource(seedLaggedFibonacci).(rand.Source))
    for i := 0; i < b.N; i++ {
        for range count {
            _ = rnd.IntN(1000)
        }
    }
}
```

各関数の内容は以下の通り。

| 関数名 | 内容 |
| --- | --- |
| `BenchmarkRandomChaCha8` | ChaCha8  |
| `BenchmarkRandomChaCha8runtime` | ChaCha8 (ランタイム組込版) |
| `BenchmarkRandomPCG` | PCG |
| `BenchmarkRandomLaggedFibonacci` | Lagged Fibonacci |

これを手元で実行するとこんな結果になった。

```
$ go test -bench Random -benchmem
goos: linux
goarch: amd64
pkg: randoms
cpu: AMD Ryzen 5 PRO 4650G with Radeon Graphics
BenchmarkRandomChaCha8-12                     184       6389731 ns/op           1 B/op           0 allocs/op
BenchmarkRandomChaCha8runtime-12              151       7874495 ns/op           0 B/op           0 allocs/op
BenchmarkRandomPCG-12                         297       4062647 ns/op           0 B/op           0 allocs/op
BenchmarkRandomLaggedFibonacci-12             355       3376871 ns/op          15 B/op           0 allocs/op
PASS
ok      randoms    6.974s
```

分かりにくいので表にまとめる。

| 関数名 | 実行時間 (ns) | Alloc サイズ(Byte) | Alloc 回数 |
| --- | ---: | ---: | ---: |
| `BenchmarkRandomChaCha8` | 6,389,731 | 1 | 0 |
| `BenchmarkRandomChaCha8runtime` | 7,874,495 | 0 | 0 |
| `BenchmarkRandomPCG` | 4,062,647 | 0 | 0 |
| `BenchmarkRandomLaggedFibonacci` | 3,376,871 | 15 | 0 |

実際には乱数取得処理をループで100万回ずつ回してるので，本当の1回あたりの実行時間は上の表の100万分の1である。

ChaCha8 疑似乱数生成器が（相対的に）遅いのは予想通りだけど PCG より [math/rand] の疑似乱数生成器のほうが速いんだな。科学技術シミュレーションなどでは疑似乱数生成器の速さも求められる。上手く使っていきたいものである。

## 参考

https://go.dev/doc/go1.22

[Go]: https://go.dev/ "The Go Programming Language"
[internal/chacha8rand]: https://pkg.go.dev/internal/chacha8rand "chacha8rand package - internal/chacha8rand - Go Packages"
[runtime]: https://pkg.go.dev/runtime "runtime package - runtime - Go Packages"
[unsafe]: https://pkg.go.dev/unsafe "unsafe package - unsafe - Go Packages"
[math/rand]: https://pkg.go.dev/math/rand "rand package - math/rand - Go Packages"
[math/rand/v2]: https://pkg.go.dev/math/rand/v2 "rand package - math/rand/v2 - Go Packages"
[crypto/rand]: https://pkg.go.dev/crypto/rand "rand package - crypto/rand - Go Packages"
<!-- eof -->
