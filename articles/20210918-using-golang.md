---
title: "コンパイル言語だからといってバイナリ提供しなくていいんだよ，Goならね" # 記事のタイトル
emoji: "💮" # アイキャッチとして使われる絵文字（1文字だけ）
type: "idea" # "tech" : 技術記事 / "idea" : アイデア記事
topics: ["go"] # タグ。["markdown", "rust", "aws"] のように指定する
published: true # 公開設定（true で公開）
---

Zenn 1周年おめでとうございます。今回は記念の小咄をひとつ。いつもと変わらないだろうって？ そりゃすまん。

Twitter で

https://twitter.com/spiegel_2007/status/1438343275237154817

と書いたが，実際にどんな具合か開陳してみる。

私のメインマシンが Ubuntu 機ということもあり，自機の運用ではシェルスクリプトとそれ以外のコードの割合が 1:3 程度になっている。たとえば [Go] のコンパイル環境を /usr/local/go/ にインストールしている場合は /etc/profile.d/ ディレクトリに

```bash:golang-bin-path.sh
# shellcheck shell=sh

# Expand $PATH to include the directory where golang applications go.
golang_bin_path="/usr/local/go/bin"
if [ -d "$golang_bin_path" -a -n "${PATH##*${golang_bin_path}}" -a -n "${PATH##*${golang_bin_path}:*}" ]; then
    export PATH=$PATH:${golang_bin_path}
fi
```

みたいなスクリプトを放り込んでおけばログイン時に（指定ディレクトリが存在すれば）勝手にパスを追加してくれる。私がシェルスクリプトを書くのは，こうした設定と設定に付随する制御を記述したい場合が多い。

一方で「処理」をメインに書きたい場合は [Go] で書くことが多くなった。

たとえば [CVE-2021-33560](https://nvd.nist.gov/vuln/detail/CVE-2021-33560) の CVSSv3 評価ベクタ `CVSS:3.1/AV:N/AC:L/PR:N/UI:N/S:U/C:H/I:N/A:N` の内容を知りたいとする。
CVSS 評価ベクタを解釈するパッケージは自作しているので

https://github.com/spiegel-im-spiegel/go-cvss

これを使って

```go:main.go
package main

import (
    "flag"
    "fmt"
    "io"
    "os"

    "github.com/spiegel-im-spiegel/go-cvss/v3/metric"
    "github.com/spiegel-im-spiegel/go-cvss/v3/report"
    "golang.org/x/text/language"
)

var template = "- `{{ .Vector }}`" + `
- {{ .SeverityName }}: {{ .SeverityValue }} (Score: {{ .BaseScore }})

| {{ .BaseMetrics }} | {{ .BaseMetricValue }} |
|--------|-------|
| {{ .AVName }} | {{ .AVValue }} |
| {{ .ACName }} | {{ .ACValue }} |
| {{ .PRName }} | {{ .PRValue }} |
| {{ .UIName }} | {{ .UIValue }} |
| {{ .SName }} | {{ .SValue }} |
| {{ .CName }} | {{ .CValue }} |
| {{ .IName }} | {{ .IValue }} |
| {{ .AName }} | {{ .AValue }} |
`

func main() {
    flag.Parse()
    if flag.NArg() < 1 {
        fmt.Fprintln(os.Stderr, "Set CVSS vector")
        return
    }
    bm, err := metric.NewBase().Decode(flag.Arg(0))
    if err != nil {
        fmt.Fprintf(os.Stderr, "%+v\n", err)
        return
    }
    r, err := report.NewBase(bm, report.WithOptionsLanguage(language.Japanese)).ExportWithString(template)
    if err != nil {
        fmt.Fprintf(os.Stderr, "%+v\n", err)
        return
    }
    if _, err := io.Copy(os.Stdout, r); err != nil {
        fmt.Fprintf(os.Stderr, "%+v\n", err)
    }
}
```

と書いておけば

```
$ go run main.go CVSS:3.1/AV:N/AC:L/PR:N/UI:N/S:U/C:H/I:N/A:N
- `CVSS:3.1/AV:N/AC:L/PR:N/UI:N/S:U/C:H/I:N/A:N`
- 深刻度: 重要 (Score: 7.5)

| 基本評価基準 | 評価値 |
|--------|-------|
| 攻撃元区分 | ネットワーク |
| 攻撃条件の複雑さ | 低 |
| 必要な特権レベル | 不要 |
| ユーザ関与レベル | 不要 |
| スコープ | 変更なし |
| 機密性への影響 | 高 |
| 完全性への影響 | なし |
| 可用性への影響 | なし |
```

などと出力してくれる。 [Go] のコンパイルは殆ど瞬時と言っていいくらい速いし，わざわざバイナリをビルドするほどの内容でもないちょっとした処理ではこういう運用をすることが多い。

こんな運用をするようになったのも（2015年あたりから始めた新参者とはいえ）手元に [Go] のコード資源が貯まってきたからであるとは言える。コード資源云々については他の言語でも言えると思うが [Go] はリファクタリングに厚い言語ゆえに再利用可能な処理を別パッケージに切り出したり，何ならまるっと他のパッケージに置き換えたりが比較的容易である点が[個人的に気に入っている](https://text.baldanders.info/remark/2021/03/awesome-golang/ "Go を褒め殺ししてみる | text.Baldanders.info")。

もっと言うとコードの置き場所はローカルでなくてもいい。たとえばシェルスクリプトを

```
$ curl https://sh.rustup.rs -sSf | sh -s -- --help
rustup-init 1.24.3 (c1c769109 2021-05-31)
The installer for rustup

USAGE:
    rustup-init [FLAGS] [OPTIONS]

FLAGS:
    -v, --verbose           Enable verbose output
    -q, --quiet             Disable progress output
    -y                      Disable confirmation prompt.
        --no-modify-path    Don't configure the PATH environment variable
    -h, --help              Prints help information
    -V, --version           Prints version information

OPTIONS:
        --default-host <default-host>              Choose a default host triple
        --default-toolchain <default-toolchain>    Choose a default toolchain to install
        --default-toolchain none                   Do not install any toolchains
        --profile [minimal|default|complete]       Choose a profile
    -c, --component <components>...                Component name to also install
    -t, --target <targets>...                      Target name to also install
```

てな感じにネットから取ってきてそのまま sh に流し込むというのはよくある手法だが，似たようなことを [Go] でもできる。

たとえば [Go] 用の ORM (Object Relational Mapper) のフレームワークとして有名な [entgo.io/ent](https://entgo.io/) では，割と最近のブログ記事でも github.com/facebook/ent/cmd/entc を go get しろみたいな記述が多いが，実際にはローカルにインストールしなくても

```
$ go run entgo.io/ent/cmd/ent@latest init TableName
```

で事足りる（[Go] 1.17 の場合）。吐き出される ent/generate.go ファイルの内容も

```go:generate.go
package ent

//go:generate go run -mod=mod entgo.io/ent/cmd/ent generate ./schema
```

という感じに遠隔リポジトリの ent モジュールを実行するようになっている。ちなみに拙作の [depm](https://github.com/spiegel-im-spiegel/depm "spiegel-im-spiegel/depm: Visualize depndency packages and modules") も同じように

```
go run github.com/spiegel-im-spiegel/depm@latest m --dot
```

とすればカレントディレクトリの [Go] コードの依存関係を [DOT 言語](https://www.graphviz.org/doc/info/lang.html "DOT Language | Graphviz")形式で出力してくれる。お試しあれ（笑）

まぁ（[Go] の処理系を持っていない）不特定のユーザに対しては実行バイナリを提供したほうがいいに決まっているし，見知らぬ第3者のコードをいきなり go run path/to/package で実行するのはセキュリティ上のリスクもあるが，スクリプト言語のようにコードベースで運用できる点は覚えておいて損はないだろう。

[Go]: https://golang.org/ "The Go Programming Language"
