---
title: "Cobra でテストしやすい CLI を構成する" # 記事のタイトル
emoji: "⌨" # アイキャッチとして使われる絵文字（1文字だけ）
type: "tech" # "tech" : 技術記事 / "idea" : アイデア記事
topics: ["go", "programming", "test"] # タグ。["markdown", "rust", "aws"] のように指定する
published: true # 公開設定（true で公開）
---

ようやく [spf13/cobra] パッケージの v1.1.0 が出たと思ったら，いつのまにか [Cobra.Dev][Cobra] なるサイトも出来てた。折角なので，記念に [Cobra] を使った CLI (Command-Line Interface) を実際に作りながら紹介してみる。

ちなみに [Cobra] は [Go] で CLI を構成するためのフレームワーク・パッケージで，標準の [flag] パッケージと比べて以下の特徴がある（主なもの）。

- 多段階のサブ・コマンドを比較的簡単に組める
- [spf13/pflag] パッケージと組み合わせて POSIX/GNU スタイルのフラグを実装できる。また，フラグのカスケード化も可能
- [spf13/viper] パッケージと組み合わせてフラグと設定ファイルの内容をバインドできる
- Usage を自動で生成・表示してくれる（カスタマイズ可能）

## 今回のお題

今回作る CLI アプリケーションの仕様は以下のものとする。

- コマンド `hash` は `encode` サブコマンドを持つ
- サブコマンド `encode` は入力テキストを SHA256 アルゴリズムでハッシュ値に符号化する

## UNIX Philosophy

実際に CLI アプリケーションを作る前に，設計指針として “UNIX Philosophy” なるものがあるので紹介しておく。曰く

1. Small is beautiful. （小さいものは美しい）
2. Make each program do one thing well. （各プログラムが一つのことをうまくやるようにせよ）
3. Build a prototype as soon as possible. （できる限り早くプロトタイプを作れ）
4. Choose portability over efficiency. （効率よりも移植しやすさを選べ）
5. Store data in flat text files. （単純なテキストファイルにデータを格納せよ）
6. Use software leverage to your advantage. （ソフトウェアの効率を優位さとして利用せよ）
7. Use shell scripts to increase leverage and portability. （効率と移植性を高めるためにシェルスクリプトを利用せよ）
8. Avoid captive user interfaces. （拘束的なユーザーインターフェースは作るな）
9. Make every program a Filter. （全てのプログラムはフィルタとして振る舞うようにせよ）

の9つである（[翻訳は Wikipedia](https://ja.wikipedia.org/wiki/UNIX%E5%93%B2%E5%AD%A6 "UNIX哲学 - Wikipedia") より）。これらを踏まえて [Go] で CLI を構成する際には以下の点に気をつけて設計していけばいいだろう。

- 他のツールと shell を介して連携できるよう標準入出力を活かすフィルタプログラムとする
- 結果の出力には JSON, YAML といったフォーマットも有効にし，外部ツールとの連携を取りやすくする
- 可能なら入出力を UTF-8 エンコーディングで統一する
- コードの可搬性（または移植性）を考慮し，プラットフォーム依存を避けるようにする

なお， CLI でもいわゆる「対話モード」ではこの指針は当てはまらないので悪しからず。

### サブコマンドとファサード・パターン

サブコマンド方式は一見 “UNIX Philosophy” に反しているように見えるが， [Go] の場合は全てのパッケージをひとつのバイナリに結合するため，関連する機能をサブコマンドとして組み込むのは悪くないやりかたである。

サブコマンドを構成する場合は「ファサード・パターン（facade pattern）」で考えるとよい。

![facade pattern](https://text.baldanders.info/golang/cli-and-facade-pattern/facade-pattern.svg)

「ファサード」は「建物の正面」という意味だそうで，システム内の各サブシステムの窓口のように機能する。ファサード自身はサブシステムの詳細を知らずコンテキスト情報を渡してキックするのみ。サブシステム側はファサードに依存せず，コンテキスト情報さえあれば処理可能とするのがコツである。

CLI なので，サブシステム側の処理結果は string, []byte, io.Reader またはそれらを出力可能な型でファサードに返せばいいだろう。

## Cobra コマンドを使ってひな型を生成する

前置きが長くなったが，さっそく [Cobra] を使って簡単な CLI を組んでみよう。

### [Cobra] コマンドのインストール

[Cobra] パッケージにはコードのひな型を出力するツールが用意されている。バイナリでの提供はないので `go get` コマンドでビルド&インストールする。

```
$ go get -u github.com/spf13/cobra/cobra
```

### main および cmd パッケージの生成

まずは環境を作る。こんな感じでどうだろう。

```
$ mkdir hash & cd hash
$ go mod init sample/hash
```

便宜上，パッケージというかモジュールのパスを `sample/hash` としておく。実際には `github.com/user/hash` みたいなパスになる筈である。

次に，前節で作成した `cobra` コマンドを使って `main` パッケージのひな型を作る。こんな感じ

```
$ cobra init --pkg-name sample/hash --viper=false
$ tree .
.
├── LICENSE
├── cmd
│   └── root.go
├── go.mod
└── main.go
```

今回は [spf13/viper] を使った設定ファイルの制御は行わないので `--viper=false` オプションを付けている。

`LICENSE` ファイルは削除するか妥当なものに差し替えてもらって構わない。`main.go` の中身はこんな感じ（コメント部分を除く）。

```go:main.go
package main

import "sample/hash/cmd"

func main() {
    cmd.Execute()
}
```

また `cmd/root.go` の中身はこんな感じ（コメント部分を除く）。

```go:cmd/root.go
package cmd

import (
    "fmt"
    "github.com/spf13/cobra"
    "os"
)

var rootCmd = &cobra.Command{
    Use:   "hash",
    Short: "A brief description of your application",
    Long: `A longer description that spans multiple lines and likely contains
examples and usage of using your application. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
}

func Execute() {
    if err := rootCmd.Execute(); err != nil {
        fmt.Println(err)
        os.Exit(1)
    }
}

func init() {
    rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
```

一応，この状態でも起動はできる。

```
$ go run main.go -h
A longer description that spans multiple lines and likely contains
examples and usage of using your application. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.
```

手を入れたいところだけど，もう少し我慢して次へ進もう。

### サブコマンドの追加

サブコマンド `encode` を追加してみよう。こんな感じ。

```
$ cobra add encode
$ tree .
.
├── LICENSE
├── cmd
│   ├── encode.go
│   └── root.go
├── go.mod
└── main.go

```

`cmd/encode.go` が追加されたのがお分かりだろうか。なお，他の `main.go` や `cmd/root.go` には手が加えられていない。

`cmd/encode.go` の中身はこんな感じ（コメント部分を除く）。

```go:cmd/encode.go
package cmd

import (
    "fmt"

    "github.com/spf13/cobra"
)

var encodeCmd = &cobra.Command{
    Use:   "encode",
    Short: "A brief description of your command",
    Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
    Run: func(cmd *cobra.Command, args []string) {
        fmt.Println("encode called")
    },
}

func init() {
    rootCmd.AddCommand(encodeCmd)
}
```

この状態で動かすとこんな感じになる。

```
$ go run main.go
A longer description that spans multiple lines and likely contains
examples and usage of using your application. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.

Usage:
  hash [command]

Available Commands:
  encode      A brief description of your command
  help        Help about any command

Flags:
  -h, --help     help for hash
  -t, --toggle   Help message for toggle

Use "hash [command] --help" for more information about a command.

$ go run main.go encode -h
A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.

Usage:
  hash encode [flags]

Flags:
  -h, --help   help for encode

$ go run main.go encode
encode called
```

というわけで，サブコマンド `encode` が組み込まれていることが分かる。

## ひな型コードを書き直す

`cobra` コマンドで生成したひな型コードにそのまま機能を組み込んでいってもいいのだが，ひな型コードには以下の問題がある。

- ひな型の `cmd` パッケージ内部で直接標準入出力やコマンドライン引数を使っている
- `cobra.Command` インスタンスを `cmd` パッケージ内で静的変数として定義している
- `cmd.Execute()` 関数内で `os.Exit()` 関数により強制終了させている箇所がある

このままだと `cmd` パッケージ自体のテストが難しいので，テストしやすいよう弄ってみる。

### spiegel-im-spiegel/gocli パッケージの導入

手前味噌で恐縮だが，ここで [spiegel-im-spiegel/gocli] パッケージを導入する。

[spiegel-im-spiegel/gocli] パッケージは入出力をコンテキスト情報として受け渡しできるようにしたもので，こんな風に使える。

```go
package main

import (
    "os"

    "github.com/spiegel-im-spiegel/gocli/exitcode"
    "github.com/spiegel-im-spiegel/gocli/rwi"
)

func run(ui *rwi.RWI) exitcode.ExitCode {
    ui.Outputln("Hello world")
    return exitcode.Normal
}

func main() {
    run(rwi.New(
        rwi.WithReader(os.Stdin),
        rwi.WithWriter(os.Stdout),
        rwi.WithErrorWriter(os.Stderr),
    )).Exit()
}
```

これを使って `cmd` パッケージのひな型パッケージを書き直していこう。

### サブコマンドを定義し直す

まずはサブコマンド `encode` から。こんな感じにしてみた。

```go:cmd/encode.go
func newEncodeCmd(ui *rwi.RWI) *cobra.Command {
    encodeCmd := &cobra.Command{
        Use:     "encode",
        Aliases: []string{"enc", "e"},
        Short:   "hash input data",
        Long:    "hash input data (detail)",
        RunE: func(cmd *cobra.Command, args []string) error {
            if err := ui.Outputln("encode called"); err != nil {
                return err
            }
            return nil
        },
    }
    return encodeCmd
}
```

これでサブコマンド `encode` 用の `cobra.Command` インスタンスを内部関数を使って動的に生成することができる。

### コマンドを定義し直す

同じように `root.go` の中身も書き直す。

```go:cmd/root.go
func newRootCmd(ui *rwi.RWI, args []string) *cobra.Command {
    rootCmd := &cobra.Command{
        Use:   "hash",
        Short: "Hash functions",
        Long:  "Hash functions (detail)",
    }
    rootCmd.SilenceUsage = true
    rootCmd.SetArgs(args)            //arguments of command-line
    rootCmd.SetIn(ui.Reader())       //Stdin
    rootCmd.SetOut(ui.ErrorWriter()) //Stdout -> Stderr
    rootCmd.SetErr(ui.ErrorWriter()) //Stderr
    rootCmd.AddCommand(
        newEncodeCmd(ui),
    )
    return rootCmd
}

func Execute(ui *rwi.RWI, args []string) exitcode.ExitCode {
    if err := newRootCmd(ui, args).Execute(); err != nil {
        return exitcode.Abnormal
    }
    return exitcode.Normal
}
```

ちなみに `newRootCmd()` 関数の

```go
rootCmd.SetArgs(args)            //arguments of command-line
rootCmd.SetIn(ui.Reader())       //Stdin
rootCmd.SetOut(ui.ErrorWriter()) //Stdout -> Stderr
rootCmd.SetErr(ui.ErrorWriter()) //Stderr
```

の部分でコマンドライン引数と入出力をセットしている。また `cmd.Execute()` 関数内部で返り値の error インスタンスの評価をしていないように見えるが，実は `cobra.Command.Execute()` メソッド内部でエラー内容を評価済みなので， `cmd.Execute()` 関数では `err != nil` の真偽のみ見ている。

最後に `main()` 関数も上記に合わせて直しておこう。

```go:main.go
func main() {
    cmd.Execute(
        rwi.New(
            rwi.WithReader(os.Stdin),
            rwi.WithWriter(os.Stdout),
            rwi.WithErrorWriter(os.Stderr),
        ),
        os.Args[1:],
    ).Exit()
}
```

これで OK。

## コマンドをテストする

書いたコードのテストを書かなきゃね。こんな感じ？

```go:cmd/encode_test.go
func TestEncode(t *testing.T) {
    testCases := []struct {
        inp  string
        outp string
        ext  exitcode.ExitCode
    }{
        {inp: "hello world\n", outp: "a948904f2f0f479b8f8197694b30184b0d2ed1c1cd2a1ec0fb85d299a192a447\n", ext: exitcode.Normal},
    }

    for _, tc := range testCases {
        r := strings.NewReader(tc.inp)
        wbuf := &bytes.Buffer{}
        ebuf := &bytes.Buffer{}
        ext := Execute(
            rwi.New(
                rwi.WithReader(r),
                rwi.WithWriter(wbuf),
                rwi.WithErrorWriter(ebuf),
            ),
            []string{"encode"},
        )
        if ext != tc.ext {
            t.Errorf("Execute() is \"%v\", want \"%v\".", ext, tc.ext)
            fmt.Println(ebuf.String())
        }
        str := wbuf.String()
        if str != tc.outp {
            t.Errorf("Execute() -> \"%v\", want \"%v\".", str, tc.outp)
        }
    }
}
```

`Execute()` 関数の引数に注意。実行結果は

```
$ go test ./...
?       sample/hash    [no test files]
--- FAIL: TestEncode (0.00s)
    encode_test.go:36: Execute() -> "encode called
        ", want "a948904f2f0f479b8f8197694b30184b0d2ed1c1cd2a1ec0fb85d299a192a447
        ".
FAIL
FAIL    sample/hash/cmd    0.002s
FAIL
```

おおう。肝心の中身を書くのを忘れてたぜ（笑） じゃあ `encode/encode.go` ファイルを作って

```go:encode/encode.go
package encode

import (
    "crypto"
    "errors"
    "io"
)

var (
    ErrNoImplement = errors.New("no implementation")
)

//Value returns hash value string from io.Reader
func Value(r io.Reader, alg crypto.Hash) ([]byte, error) {
    if !alg.Available() {
        return nil, ErrNoImplement
    }
    h := alg.New()
    if _, err := io.Copy(h, r); err != nil {
        return nil, err
    }
    return h.Sum(nil), nil
}
```

って感じでいいかな。これを `cmd.newEncodeCmd()` 関数に組み込んで

```go:cmd/encode.go
func newEncodeCmd(ui *rwi.RWI) *cobra.Command {
    encodeCmd := &cobra.Command{
        Use:     "encode",
        Aliases: []string{"enc", "e"},
        Short:   "hash input data",
        Long:    "hash input data (detail)",
        RunE: func(cmd *cobra.Command, args []string) error {
            v, err := encode.Value(ui.Reader(), crypto.SHA256)
            if err != nil {
                return err
            }
            fmt.Fprintf(ui.Writer(), "%x\n", v)
            return nil
        },
    }
    return encodeCmd
}
```

テスト再開。

```
$ go test ./...
?       sample/hash    [no test files]
ok      sample/hash/cmd    0.003s
?       sample/hash/encode    [no test files]
```

よーし，うむうむ，よーし。

最後に実際にコマンドを実行してみる。

```
$ go run main.go encode -h
hash input data (detail)

Usage:
  hash encode [flags]

Aliases:
  encode, enc, e

Flags:
  -h, --help   help for encode

$ echo hello world | go run main.go encode
a948904f2f0f479b8f8197694b30184b0d2ed1c1cd2a1ec0fb85d299a192a447
```

無事に実行できた。

なお，今回書いたコードは [`github.com/spiegel-im-spiegel/zenn-docs/code/cobra/hash`](https://github.com/spiegel-im-spiegel/zenn-docs/tree/main/code/cobra/hash) に置いている。参考にどうぞ。

[Go]: https://golang.org/ "The Go Programming Language"
[flag]: https://golang.org/pkg/flag/ "flag - The Go Programming Language"
[Cobra]: https://cobra.dev/ "Cobra.Dev"
[spf13/cobra]: https://github.com/spf13/cobra "spf13/cobra: A Commander for modern Go CLI interactions"
[spf13/pflag]: https://github.com/spf13/pflag "spf13/pflag: Drop-in replacement for Go's flag package, implementing POSIX/GNU-style --flags."
[spf13/viper]: https://github.com/spf13/viper "spf13/viper: Go configuration with fangs"
[spiegel-im-spiegel/gocli]: https://github.com/spiegel-im-spiegel/gocli "spiegel-im-spiegel/gocli: Minimal Packages for Command-Line Interface"
<!-- eof -->
