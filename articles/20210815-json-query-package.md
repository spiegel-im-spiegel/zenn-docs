---
title: "クエリ型の JSON パーサ" # 記事のタイトル
emoji: "💻" # アイキャッチとして使われる絵文字（1文字だけ）
type: "tech" # "tech" : 技術記事 / "idea" : アイデア記事
topics: ["go", "programming", "json"] # タグ。["markdown", "rust", "aws"] のように指定する
published: true # 公開設定（true で公開）
---

[Go] の標準パッケージには [encoding/json] という JSON パーサがあるが，サードパーティ製のパッケージも色々ある。たとえば [encoding/json] 互換パーサとしては [github.com/goccy/go-json] が速いらしく，これについては以前に紹介している。

https://zenn.dev/spiegel/articles/20210404-another-json-package

[encoding/json] 標準パッケージは JSON データ全体を任意の構造体または map[string]interface{} 型の連想配列に落とし込んで使うが，[jq] のようにクエリを発行して値を取得するタイプもあると便利だろう。

私が2年前に手遊びで作った [gjq](https://github.com/spiegel-im-spiegel/gjq "spiegel-im-spiegel/gjq: Another Implementation of jq by golang") はパーサとして [github.com/savaki/jq](https://github.com/savaki/jq "savaki/jq: A high performance Golang implementation of the incredibly useful jq command line tool.") パッケージを使っているのだが，最後に更新されてから5年ほど経っているようで，モジュールにも未対応で，今となってはあまり使いたくない感じである。

最近知ったのが [github.com/buger/jsonparser] パッケージ。 [jq] とはちょっと違うが，これも要素を指定して JSON データを解析してくれるようだ。こんな感じ。

```go:sample1.go
// +build run

package main

import (
    "fmt"
    "os"

    "github.com/buger/jsonparser"
)

var jsondata = []byte(`{
  "person": {
    "name": {
      "first": "Leonid",
      "last": "Bugaev",
      "fullName": "Leonid Bugaev"
    },
    "github": {
      "handle": "buger",
      "followers": 109
    },
    "avatars": [
      { "url": "https://avatars1.githubusercontent.com/u/14009?v=3&s=460", "type": "thumbnail" }
    ]
  },
  "company": {
    "name": "Acme"
  }
}`)

func main() {
    v, err := jsonparser.GetString(jsondata, "person", "avatars", "[0]", "url")
    if err != nil {
        fmt.Fprintln(os.Stderr, err)
        return
    }
    fmt.Println(v)
    // Output:
    // https://avatars1.githubusercontent.com/u/14009?v=3&s=460
}
```

更に for-each 風の高階関数[^hof1] も用意されていて

[^hof1]: 念のために説明すると「高階関数（higher-order function）」とは，第1級関数（first-class function）をサポートしている言語において (1) 関数を引数に取る (2) 関数を返す の少なくとも1つの機能を満たす関数である。関数型プログラミング言語なんかではお馴染みのやつだが [Go] でも実装できる。ただし総称型を（今のところ）サポートしていない [Go] では，だいぶダサい感じになるのは否めない（笑）

```go:sample2.go
func main() {
    if err := jsonparser.ObjectEach(jsondata, func(key []byte, value []byte, dataType jsonparser.ValueType, offset int) error {
        fmt.Printf("Offset: %d\n\tKey: '%s'\n\tValue: '%s'\n\tType: %s\n", offset, string(key), string(value), dataType)
        return nil
    }, "person", "name"); err != nil {
        fmt.Fprintln(os.Stderr, err)
        return
    }
}
```

てな感じに書くこともできる。ちなみにこれを実行すると

```
$ go run sample2.go 
Offset: 53
    Key: 'first'
    Value: 'Leonid'
    Type: string
Offset: 77
    Key: 'last'
    Value: 'Bugaev'
    Type: string
Offset: 112
    Key: 'fullName'
    Value: 'Leonid Bugaev'
    Type: string
```

と出力される。

[github.com/buger/jsonparser] パッケージは [encoding/json] 標準パッケージより速いと豪語している。公式のベンチマークによると

> Each test processes a 24kb JSON record (based on Discourse API) It should read 2 arrays, and for each item in array get a few fields.  Basically it means processing a full JSON file.
> 
> https://github.com/buger/jsonparser/blob/master/benchmark/benchmark_large_payload_test.go
> 
> | Library | time/op | bytes/op | allocs/op |
> | --- | --- | --- | --- |
> | encoding/json struct | 748336 | 8272 | 307 |
> | encoding/json interface{} | 1224271 | 215425 | 3395 |
> | a8m/djson | 510082 | 213682 | 2845 |
> | pquerna/ffjson | **312271** | **7792** | **298** |
> | mailru/easyjson | **154186** | **6992** | **288** |
> | buger/jsonparser | **85308** | **0** | **0** |
> 
> `jsonparser` now is a winner, but do not forget that it is way more lightweight parser than `ffson` or `easyjson`, and they have to parser all the data, while > `jsonparser` parse only what you need. All `ffjson`, `easysjon` and `jsonparser` have their own parsing code, and does not depend on `encoding/json` or `interface{}`, thats one of the reasons why they are so fast. `easyjson` also use a bit of `unsafe` package to reduce memory consuption (in theory it can lead to some unexpected GC issue, but i did not tested enough)
> 
> (via “[buger/jsonparser: One of the fastest alternative JSON parser for Go that does not require schema][github.com/buger/jsonparser]”)

ということで，（条件付きではあるが）アロケーションを発生させずかなり高速な処理を行っていることが分かる。

[Go]: https://golang.org/ "The Go Programming Language"
[jq]: https://stedolan.github.io/jq/
[encoding/json]: https://pkg.go.dev/encoding/json "json · pkg.go.dev"
[github.com/goccy/go-json]: https://github.com/goccy/go-json "goccy/go-json: Fast JSON encoder/decoder compatible with encoding/json for Go"
[github.com/buger/jsonparser]: https://github.com/buger/jsonparser "buger/jsonparser: One of the fastest alternative JSON parser for Go that does not require schema"










[Go]: https://golang.org/ "The Go Programming Language"
<!-- eof -->
