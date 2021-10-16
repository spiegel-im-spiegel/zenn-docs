---
title: "crypt/rand を math/rand として使う" # 記事のタイトル
emoji: "💻" # アイキャッチとして使われる絵文字（1文字だけ）
type: "tech" # "tech" : 技術記事 / "idea" : アイデア記事
topics: ["go", "programming", "random"] # タグ。["markdown", "rust", "aws"] のように指定する
published: true # 公開設定（true で公開）
---

みなさん一度は思いませんでした？ なんで [math/rand] と [crypto/rand] は内部構成が全く違うんだろう，と。まぁ，目的が異なるので構成が違っててもおかしくないけど，せめて [crypto/rand] を [math/rand] のソースとして使えればいいのに，と。

実は [crypto/rand] を [math/rand] のソースにするのはそんなに難しくない。 [math/rand] のソースの定義は

```go:math/rand/rand.go
// A Source represents a source of uniformly-distributed
// pseudo-random int64 values in the range [0, 1<<63).
type Source interface {
	Int63() int64
	Seed(seed int64)
}

// A Source64 is a Source that can also generate
// uniformly-distributed pseudo-random uint64 values in
// the range [0, 1<<64) directly.
// If a Rand r's underlying Source s implements Source64,
// then r.Uint64 returns the result of one call to s.Uint64
// instead of making two calls to s.Int63.
type Source64 interface {
	Source
	Uint64() uint64
}
```

となっているので，これにマッチするラッパーを作ってやればいいわけ。たとえばこんな感じ。

```go:sample.go
type Source struct{}

// Seed method is dummy function for rand.Source interface.
func (s Source) Seed(seed int64) {}

// Uint64 method generates a random number in the range [0, 1<<64).
func (s Source) Uint64() uint64 {
	b := [8]byte{}
	ct, _ := rand.Read(b[:])
	return binary.BigEndian.Uint64(b[:ct])
}

// Int63 method generates a random number in the range [0, 1<<63).
func (s Source) Int63() int64 {
	return (int64)(s.Uint64() >> 1)
}
```

こうすれば

```go:sample.go
fmt.Println(rand.New(Source{}).Float64()) // 0.9581627789424901
```

という感じに [rand][math/rand].Rand 型が提供するメソッドを利用することができる[^cs1]。コード全体ではこんな感じ。

[^cs1]: [rand][math/rand].Rand 型が提供するメソッドは並行的に安全（concurrency safe）ではないのでご注意を。

```go:sample.go
//go:build run
// +build run

package main

import (
	"crypto/rand"
	"encoding/binary"
	"fmt"
	mrand "math/rand"
)

type Source struct{}

// Seed method is dummy function for rand.Source interface.
func (s Source) Seed(seed int64) {}

// Uint64 method generates a random number in the range [0, 1<<64).
func (s Source) Uint64() uint64 {
	b := [8]byte{}
	ct, _ := rand.Read(b[:])
	return binary.BigEndian.Uint64(b[:ct])
}

// Int63 method generates a random number in the range [0, 1<<63).
func (s Source) Int63() int64 {
	return (int64)(s.Uint64() >> 1)
}

func main() {
	fmt.Println(mrand.New(Source{}).Float64())
}
```

というわけで，これをパッケージ化することにした。といっても，たったこれだけの機能のためにリポジトリを作るのはもったいないので，拙作の疑似乱数パッケージ [github.com/spiegel-im-spiegel/mt] のオマケ機能として組み込んでみた。こんな感じに使える。

```go
//go:build run
// +build run

package main

import (
	"fmt"
	"math/rand"

	"github.com/spiegel-im-spiegel/mt/secure"
)

func main() {
	fmt.Println(rand.New(secure.Source{}).Uint64())
}
```

よーし，うむうむ，よーし。では，作業の続きをするか。

https://text.baldanders.info/release/mersenne-twister-by-golang/

[Go]: https://golang.org/ "The Go Programming Language"
[crypto/rand]: https://pkg.go.dev/crypto/rand "rand package - crypto/rand - pkg.go.dev"
[math/rand]: https://pkg.go.dev/math/rand "rand package - math/rand - pkg.go.dev"
[github.com/spiegel-im-spiegel/mt]: https://github.com/spiegel-im-spiegel/mt "spiegel-im-spiegel/mt: Mersenne Twister; Pseudo Random Number Generator, Implemented by Golang"
<!-- eof -->
