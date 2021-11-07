---
title: "PlantUML で JSON データを簡単視覚化" # 記事のタイトル
emoji: "💮" # アイキャッチとして使われる絵文字（1文字だけ）
type: "idea" # "tech" : 技術記事 / "idea" : アイデア記事
topics: ["json", "yaml", "plantuml"] # タグ。["markdown", "rust", "aws"] のように指定する
published: true # 公開設定（true で公開）
---

最近，仕事で使うことがあってたまたま気がついたのだが， [PlantUML] って JSON や YAML のデータを視覚化できるんだね。

やり方は簡単。たとえば

```json
{
  "firstName": "John",
  "lastName": "Smith",
  "isAlive": true,
  "age": 28,
  "address": {
    "streetAddress": "21 2nd Street",
    "city": "New York",
    "state": "NY",
    "postalCode": "10021-3100"
  },
  "phoneNumbers": [
    {
      "type": "home",
      "number": "212 555-1234"
    },
    {
      "type": "office",
      "number": "646 555-4567"
    }
  ],
  "children": [],
  "spouse": null
}
```

という JSON データあるとすると，これを `@startjson`...`@endjson` で囲ってやればよい。こんな感じ。

```json
@startjson "json"
{
  "firstName": "John",
  "lastName": "Smith",
  "isAlive": true,
  "age": 28,
  "address": {
    "streetAddress": "21 2nd Street",
    "city": "New York",
    "state": "NY",
    "postalCode": "10021-3100"
  },
  "phoneNumbers": [
    {
      "type": "home",
      "number": "212 555-1234"
    },
    {
      "type": "office",
      "number": "646 555-4567"
    }
  ],
  "children": [],
  "spouse": null
}
@endjson
```

これを [PlantUML] で処理すると

![](/images/visualize-json-by-plantuml/json.png)

てな感じになる。

もう少し遊んでみよう。拙作の [github.com/spiegel-im-spiegel/pa-api](https://github.com/spiegel-im-spiegel/pa-api "spiegel-im-spiegel/pa-api: APIs for Amazon Product Advertising API v5 by Golang") パッケージを使って Amazon の書籍情報が取れるのだが，これを使って以下のようなコードを組んでみる。

```go:sample.go
//go:build run
// +build run

package main

import (
    "bytes"
    "fmt"
    "io"
    "net/http"
    "os"

    paapi5 "github.com/spiegel-im-spiegel/pa-api"
    "github.com/spiegel-im-spiegel/pa-api/query"
)

func main() {
    //Create client
    client := paapi5.New(
        paapi5.WithMarketplace(paapi5.LocaleJapan),
    ).CreateClient(
        "mytag-20",
        "AKIAIOSFODNN7EXAMPLE",
        "1234567890",
    )

    //Make query
    q := query.NewGetItems(
        client.Marketplace(),
        client.PartnerTag(),
        client.PartnerType(),
    ).
        ASINs([]string{"B09HK66P5X"}).
        EnableItemInfo().
        EnableImages().
        EnableParentASIN()

    //Requet and response
    body, err := client.Request(q)
    if err != nil {
        fmt.Printf("%+v\n", err)
        return
    }
    fmt.Println("@startjson \"book\"")
    if _, err := io.Copy(os.Stdout, bytes.NewReader(body)); err != nil {
        fmt.Printf("%+v\n", err)
        return
    }
    fmt.Println("\n@endjson")
}
```

（アソシエイト・タグやアクセス・キーはダミーでそのまま使えないので悪しからず）

ちなみに B09HK66P5X は2021年11月に出る『[Java言語で学ぶデザインパターン入門第3版](https://www.amazon.co.jp/dp/B09HK66P5X)』の ASIN コードである。

これを実行すると

```
$ go run sample.go > book.puml
```

以下の出力を得られる（JSON データは分かりやすく整形しています）。

```json: book.puml
@startjson "book"
{
  "ItemsResult": {
    "Items": [
      {
        "ASIN": "B09HK66P5X",
        "DetailPageURL": "https://www.amazon.co.jp/dp/B09HK66P5X?tag=mytag-20&linkCode=ogi&th=1&psc=1",
        "Images": {
          "Primary": {
            "Large": {
              "Height": 500,
              "URL": "https://m.media-amazon.com/images/I/41Fxrb9KFvL._SL500_.jpg",
              "Width": 393
            },
            "Medium": {
              "Height": 160,
              "URL": "https://m.media-amazon.com/images/I/41Fxrb9KFvL._SL160_.jpg",
              "Width": 125
            },
            "Small": {
              "Height": 75,
              "URL": "https://m.media-amazon.com/images/I/41Fxrb9KFvL._SL75_.jpg",
              "Width": 58
            }
          }
        },
        "ItemInfo": {
          "ByLineInfo": {
            "Contributors": [
              {
                "Locale": "ja_JP",
                "Name": "結城 浩",
                "Role": "著",
                "RoleType": "author"
              }
            ],
            "Manufacturer": {
              "DisplayValue": "SBクリエイティブ",
              "Label": "Manufacturer",
              "Locale": "ja_JP"
            }
          },
          "Classifications": {
            "Binding": {
              "DisplayValue": "Kindle版",
              "Label": "Binding",
              "Locale": "ja_JP"
            },
            "ProductGroup": {
              "DisplayValue": "Digital Ebook Purchas",
              "Label": "ProductGroup",
              "Locale": "ja_JP"
            }
          },
          "ContentInfo": {
            "Languages": {
              "DisplayValues": [
                {
                  "DisplayValue": "日本語",
                  "Type": "発行済み"
                }
              ],
              "Label": "Language",
              "Locale": "ja_JP"
            },
            "PublicationDate": {
              "DisplayValue": "2021-11-12T00:00:00.000Z",
              "Label": "PublicationDate",
              "Locale": "en_US"
            }
          },
          "ProductInfo": {
            "IsAdultProduct": {
              "DisplayValue": false,
              "Label": "IsAdultProduct",
              "Locale": "en_US"
            },
            "ReleaseDate": {
              "DisplayValue": "2021-11-13T00:00:00.000Z",
              "Label": "ReleaseDate",
              "Locale": "en_US"
            }
          },
          "Title": {
            "DisplayValue": "Java言語で学ぶデザインパターン入門第3版",
            "Label": "Title",
            "Locale": "ja_JP"
          }
        }
      }
    ]
  }
}
@endjson
```

これを [PlantUML] で処理すると

![](/images/visualize-json-by-plantuml/book.png)

てな感じになる。よーし，うむうむ，よーし。

YAML の場合は `@startyaml`...`@endyaml` で囲む。

```yaml
@startyaml "yaml"
doe: "a deer, a female deer"
ray: "a drop of golden sun"
pi: 3.14159
xmas: true
french-hens: 3
calling-birds: 
    - huey
    - dewey
    - louie
    - fred
xmas-fifth-day: 
    calling-birds: four
    french-hens: 3
    golden-rings: 5
    partridges: 
        count: 1
        location: "a pear tree"
    turtle-doves: two
@endyaml
```

これを [PlantUML] で処理すると

![](/images/visualize-json-by-plantuml/yaml.png)

という感じになる。

生の JSON や YAML のデータをそのまま使えるのがいい感じ。お試しあれ。

## 参考

https://plantuml.com/ja/json
https://plantuml.com/ja/yaml
https://text.baldanders.info/remark/2018/12/plantuml-1/
https://text.baldanders.info/release/pa-api-v5/

[PlantUML]: http://plantuml.com/ "Open-source tool that uses simple textual descriptions to draw UML diagrams."
