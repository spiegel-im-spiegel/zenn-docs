---
title: "Go のモジュール管理【バージョン 1.16 改訂版】"
emoji: "💻" # アイキャッチとして使われる絵文字（1文字だけ）
type: "tech" # "tech" : 技術記事 / "idea" : アイデア記事
topics: ["go", "programming"] # タグ。["markdown", "rust", "aws"] のように指定する
published: true # 公開設定（true で公開）
---

[Go] のモジュールについては[自ブログ](https://text.baldanders.info/golang/ "プログラミング言語 Go | text.Baldanders.info")でもよく話題にするのだが，差分情報が多く内容が分散しているため，ここの Zenn でまとめておく。

## 用語の整理

まず最初に用語の定義をしておく。

### GOPATH モードとモジュール対応モード

バージョン 1.11 以降からは [Go] ツールーチェーンは以下の2つのモードのどちらかで動作する。

- **GOPATH モード (GOPATH mode)** : バージョン 1.10 までの動作モード。標準ライブラリを除く全てのパッケージのコード管理とビルドを環境変数 GOPATH で指定されたディレクトリ下で行う。パッケージの管理はリポジトリの最新リビジョンのみが対象となる
- **モジュール対応モード (module-aware mode)** : 標準ライブラリを除く全てのパッケージをモジュールとして管理する。コード管理とビルドは任意のディレクトリで可能で，モジュールはリポジトリのバージョンタグまたはリビジョン毎に管理される

### 「モジュール」とは

モジュール対応モードでは，標準ライブラリを除くパッケージを「モジュール（module）」として管理する。パッケージが [git] 等のバージョン管理ツールで管理されている場合はバージョン（またはリビジョン）ごとに異なるモジュールと見なされる。つまりモジュールの実体は「パッケージ＋バージョン」ということになる。
 
ただしコード上ではパッケージとモジュールの間に区別はなく，したがってソースコードを書き換える必要はない。モジュールはソースコードではなく go.mod ファイルで管理される。

## 環境変数 GO111MODULE によるモードの切り替え

バージョン 1.11 以降では2つのモードの切り替えのために環境変数 GO111MODULE が用意されている。

```
$ go env | grep GO111MODULE
GO111MODULE=""
```

GO111MODULE の取りうる値は以下の通り。

| 値     | 内容 |
| ------ | ---- |
| `on`   | 常にモジュール対応モードで動作する |
| `off`  | 常に GOPATH モードで動作する  |
| `auto` | `$GOPATH/src` 以下のディレクトリに配置され `go.mod` を含まないパッケージは GOPATH モードで，それ以外はモジュール対応モードで動作する |

バージョン 1.16 から GO111MODULE 未指定時の既定値が `on` になった（1.15 までは `auto`）。 GOPATH モードを使いたいのであれば GO111MODULE の値を `auto` または `off` に設定する[^env1]。

[^env1]: [Go] の環境変数の取り扱いについては，拙文「[Go 言語の環境変数管理](https://text.baldanders.info/golang/go-env/)」をご覧あれ。

```
$ go env -w GO111MODULE=auto
```

ただし [Go] では環境変数 GOPATH への依存を徐々に薄めつつあり，最終的には $GOPATH/bin ディレクトリ以外は使われなくなると思われる。

https://zenn.dev/tennashi/articles/3b87a8d924bc9c43573e

## go.mod と go.sum

### go.mod の作成

モジュールの管理は go.mod および go.sum ファイルで行う。新たに go.mod を作成する場合は，以下のコマンドを叩く。

```
$ go mod init hello
go: creating new go.mod: module hello
go: to add module requirements and sums:
	go mod tidy
```

これでカレント・ディレクトリ直下に go.mod ファイルが作成される。中身はこんな感じ。

```txt:go.mod
module hello

go 1.16
```

行頭の `module` や `go` はディレクティブ（directive）と呼ばれるものだ。たとえば `module hello` はモジュール名が `hello` であることを示す。

モジュール名は任意に付けられるがソース・コードの `import` で指定するパッケージパスに合わせるのが無難である。たとえば [github.com/spiegel-im-spiegel/fetch](https://github.com/spiegel-im-spiegel/fetch) パッケージであれば

```txt:go.mod
module github.com/spiegel-im-spiegel/fetch
```

とする。

他に go.mod ファイルで使えるディレクティブは以下の通り。

| ディレクティブ | 記述例                                          | 内容                   |
| -------------- | ----------------------------------------------- | ---------------------- |
| `module`       | `module my/thing`                               | モジュール名           |
| `go`           | `1.16`                                          | 有効な Go バージョン   |
| `require`      | `require other/thing v1.0.2`                    | インポート・モジュール |
| `exclude`      | `exclude old/thing v1.2.3`                      | 除外モジュール         |
| `replace`      | `replace bad/thing v1.4.5 => good/thing v1.4.5` | モジュールの置換       |
| `retract`      | `v1.0.5`                                        | 撤回バージョン         |

`replace` ディレクティブはパッケージのあるリポジトリがリダイレクトされていて上手くインポートできない等の状況で使える。たとえばこんな感じ。

```txt:go.mod
module sample

require gopkg.in/russross/blackfriday.v2 v2.0.1

replace gopkg.in/russross/blackfriday.v2 v2.0.1 => github.com/russross/blackfriday/v2 v2.0.1
```

`retract` ディレクティブはバージョン 1.16 から導入されたもので

```txt:go.mod
// Remote-triggered crash in package foo. See CVE-2021-01234.
retract v1.0.5
```

などとコメントと共に指定しておくと

```
$ go list -m -u all
example.com/lib v1.0.0 (retracted)
$ go get .
go: warning: example.com/lib@v1.0.5: retracted by module author:
    Remote-triggered crash in package foo. See CVE-2021-01234.
go: to switch to the latest unretracted version, run:
    go get example.com/lib@latest
```

てな感じにコメントの内容でワーニングを出してくれるらしい。

### go.sum の中身

go.sum ファイルにはインポートするモジュールのハッシュ情報が格納されている。たとえば go.mod ファイルで `require` ディレクティブが

```txt:go.mod
require github.com/spiegel-im-spiegel/errs v1.0.2
```

と指定されている場合， go.sum ファイルの内容は

```txt:go.sum
github.com/spiegel-im-spiegel/errs v1.0.2 h1:v4amEwRDqRWjKHOILQnJSovYhZ4ZttEnBBXNXEzS6Sc=
github.com/spiegel-im-spiegel/errs v1.0.2/go.mod h1:UoasJYYujMcdkbT9USv8dfZWoMyaY3btqQxoLJImw0A=
```

といった感じになる。

go.sum ファイルの内容はインポートするモジュールの完全性（integrity）を保証するものだが，手作業で記述できるようなものではないので，次に紹介する go mod tidy コマンドを使って更新する。

### go mod tidy によるモジュール情報の更新

[Go] バージョン 1.15 までは go build や go test といったコマンドの実行時に go.mod および go.sum ファイルが更新されていたが， 1.16 からはこれができなくなった。たとえば，ソース・コード上で新しい外部パッケージを import しても go.mod と go.sum に記述がないと

```
$ go test ./...
main.go:9:2: no required module provides package github.com/spiegel-im-spiegel/cov19jpn/chart; to add it:
	go get github.com/spiegel-im-spiegel/cov19jpn/chart
```

とか

```
$ go test ./...
go: github.com/spiegel-im-spiegel/cov19jpn@v0.2.0: missing go.sum entry; to add it:
	go mod download github.com/spiegel-im-spiegel/cov19jpn
```

みたいなエラーが出たりする。 go.mod および go.sum ファイルをいい感じに更新したいのであれば

```text
$ go mod tidy
```

とするとよい。

## [Semantic Versioning] によるバージョン管理

モジュールのバージョンははリポジトリのリビジョン番号またはバージョンタグによって管理されるが，バージョンタグに関しては [Semantic Versioning] のルールに則ってバージョン番号を設定することが推奨されている。

[![research.swtch.com/impver.png](https://research.swtch.com/impver.png)](https://research.swtch.com/vgo-import "via “Semantic Import Versioning”")

このように後方互換性のない変更がある場合にはメジャーバージョンを，後方互換性が担保された変更や追加についてはマイナーバージョンを，不具合や脆弱性の修正については第3位のパッチバージョンを上げるようにする。またメジャーバージョンを上げる際には，図のようにディレクトリを分離するか， go.mod ファイルの `module` ディレクティブの値を変更するのが簡単である[^ver1]。

[^ver1]: v0 から v1 へのメジャーバージョンアップの場合は例外的に後方互換性は考慮されない。 v0 はベータ版という扱いらしい。それ以外のメジャーバージョンアップの際にパッケージパスまたはモジュール名で区別せずにバージョンタグを打つと強制的に `v2.0.0+incompatible` みたいなショボい表記にされる。

以下は [github.com/mattn/jvgrep](https://github.com/mattn/jvgrep) の例：

```txt:go.mod
module github.com/mattn/jvgrep/v5
```

## 特定バージョンのモジュールをビルド&インストールする

バージョン 1.16 から go install コマンドで モジュール@バージョン を指定できるようになった。

```text
$ go install golang.org/x/tools/gopls@v0.6.5
```

とにかく最新版が欲しい場合は

```text
$ go install golang.org/x/tools/gopls@latest
go: downloading golang.org/x/tools/gopls v0.6.5
...
```

などとすればOK。

バージョン付きで指定する際は，パッケージ・パスではなくモジュール名で指定する点に注意。たとえば，先ほどの [github.com/mattn/jvgrep](https://github.com/mattn/jvgrep) の場合，パッケージ・パスで指定しても

```
$ go install github.com/mattn/jvgrep@latest
go: downloading github.com/mattn/jvgrep v1.9.0
go: downloading github.com/mattn/jvgrep v5.8.1+incompatible
...
```

てな感じで `latest` を指定しても最新版を取ってくれない。こういうときは

```
$ go install github.com/mattn/jvgrep/v5@latest
go: downloading github.com/mattn/jvgrep/v5 v5.8.9
...
```

とすれば無問題である。

なお，バージョン 1.16 では go.mod ファイルに `replace` や `exclude` ディレクティブが含まれていると go install に失敗する。

```
$ go install github.com/spiegel-im-spiegel/gnkf@latest
go: downloading github.com/spiegel-im-spiegel/gnkf v0.3.0
go install github.com/spiegel-im-spiegel/gnkf@latest: github.com/spiegel-im-spiegel/gnkf@v0.3.0
	The go.mod file for the module providing named packages contains one or
	more replace directives. It must not contain directives that would cause
	it to be interpreted differently than if it were the main module.
```

とほほ orz

こういうときは go install ではなく go get コマンドを使えば取り敢えず大丈夫なようだ。

```
$ go get github.com/spiegel-im-spiegel/gnkf@latest
```

ただし go get コマンドを使ったビルド&インストールは将来バージョンで廃止されるらしいので，何とかしなきゃなぁ...

## 参考

https://blog.golang.org/go116-module-changes

[Go]: https://golang.org/ "The Go Programming Language"
[git]: https://git-scm.com/ "Git"
[Semantic Versioning]: http://semver.org/ "Semantic Versioning 2.0.0 | Semantic Versioning"
[GitHub]: https://github.com/
