---
title: "ライセンスファイルからライセンスを推定する" # 記事のタイトル
emoji: "💻" # アイキャッチとして使われる絵文字（1文字だけ）
type: "tech" # "tech" : 技術記事 / "idea" : アイデア記事
topics: ["go", "programming"] # タグ。["markdown", "rust", "aws"] のように指定する
published: true # 公開設定（true で公開）
---

たとえばリポジトリ直下に LICENSE というファイルがあるとして，このファイルが実際に何のライセンスを指しているか機械的に調べる方法はないだろうか。実は Google による [Go] パッケージが公開されている[^g1]。

[^g1]: ただし README.md には “This is not an official Google product” とあり Google 公式パッケージではないことが明記されている。ご注意を。

https://github.com/google/licenseclassifier

私は以前からこのパッケージを利用しているのだが，開発の主力が v2 系に移っているようだ。2022-07-22 に [v2.0.0-pre6](https://github.com/google/licenseclassifier/releases/tag/v2.0.0-pre6) がリリースされていた。さっそく試してみることにする。

今回のサンプルコードはこんな感じ。

```go:sample.go
package main

import (
    "flag"
    "fmt"
    "os"

    "github.com/google/licenseclassifier/v2/assets"
)

func main() {
    flag.Parse()
    args := flag.Args()
    if len(args) < 1 {
        fmt.Fprintln(os.Stderr, os.ErrInvalid)
        return
    }
    file, err := os.Open(args[0])
    if err != nil {
        fmt.Fprintln(os.Stderr, err)
        return
    }
    defer file.Close()

    c, err := assets.DefaultClassifier()
    if err != nil {
        fmt.Fprintln(os.Stderr, err)
        return
    }

    res, err := c.MatchFrom(file)
    if err != nil {
        fmt.Fprintln(os.Stderr, err)
        return
    }
    if len(res.Matches) == 0 {
        fmt.Fprintln(os.Stderr, args[0], "is not license file.")
        return
    }
    for _, m := range res.Matches {
        fmt.Println(m.MatchType, m.Name, )
    }
}
```

手順としては

1. コマンドライン引数で指定したファイルを開く
2. `assets.DefaultClassifier()` で解析のための辞書情報（`*classifier.Classifier` 型）を取得する
3. `MatchFrom()` メソッドでファイルを解析し，結果を表示する

という感じ。では，実際に動かしてみよう。

```
$ go run sample.go ./LICENSE 
License Apache-2.0
```

というわけで，指定した LICENSE ファイルは `License` タイプの `Apache-2.0` ライセンスであることが分かった。よーし，うむうむ，よーし。

まだ正式リリースではないようだが，使えるレベルに達してると思う。上手く利用していただきたい。

## 【付録】 Software Package Data Exchange

（[リクエスト](https://twitter.com/fu7mu4/status/1556141959755886593)にお応えして）

先ほどの `Apache-2.0` だが，これは SPDX (Software Package Data Exchange) のライセンス識別子と呼ばれるものである。

https://spdx.dev/

ちなみに SPDX は [ISO/IEC 5962:2021](https://www.iso.org/standard/81870.html) として[標準化](https://www.linuxfoundation.org/press-release/spdx-becomes-internationally-recognized-standard-for-software-bill-of-materials/ "SPDX Becomes Internationally Recognized Standard for Software Bill of Materials - Linux Foundation")されたそうな。

SPDX ライセンス識別子の一覧は以下のページで確認することができる。

https://spdx.org/licenses/

ソフトウェア・サプライチェーンを構成する際に SPDX ソフトウェア部品表 (software bills of materials; SBOMs) を利用することで情報の共通化を図ることができる。

> Between eighty and ninety percent (80%-90%) of a modern application is assembled from open source software components. An SBOM accounts for the software components contained in an application — open source, proprietary, or third-party — and details their provenance, license, and security attributes. SBOMs are used as a part of a foundational practice to track and trace components across software supply chains. SBOMs also help to proactively identify software issues and risks and establish a starting point for their remediation.
*(via “[SPDX Becomes Internationally Recognized Standard for Software Bill of Materials](https://www.linuxfoundation.org/press-release/spdx-becomes-internationally-recognized-standard-for-software-bill-of-materials/)”)*

SPDX ライセンス識別子はソフトウェア部品表を構成する情報のひとつとして使えるわけだ。 [SPDX のリポジトリ](https://github.com/spdx)に C や [Go] による製品の[部品表サンプル](https://github.com/spdx/spdx-examples)がある。参考になれば幸いである。

[Go]: https://go.dev/ "The Go Programming Language"
