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

バージョン 1.11 以降からは [Go] ツールーチェーンは以下の 2 つのモードのどちらかで動作する。

- **GOPATH モード (GOPATH mode)** : バージョン 1.10 までのモード。標準ライブラリを除く全てのパッケージのコード管理とビルドを環境変数 GOPATH で指定されたディレクトリ下で行う。パッケージの管理はリポジトリの最新リビジョンのみが対象となる
- **モジュール対応モード (module-aware mode)** : 標準ライブラリを除く全てのパッケージをモジュールとして管理する。コード管理とビルドは任意のディレクトリで可能で，モジュールはリポジトリのバージョンタグまたはリビジョン毎に管理される

### 「パッケージ」とは

[Go] コンパイラにおける処理単位。具体的にはひとつの物理ディレクトリ内のファイル群をひとつのパッケージとして処理する。また，インポート宣言（import declaration）でパッケージを指定する際には

```go
import "path/to/package"
```

などと，パスで指定する。コンパイラは指定されたパスを見てパッケージを探す。

最優先は標準パッケージで GOPATH モードであれば環境変数 GOPATH で指定されたディレクトリ以下を探す。モジュール対応モードでは[モジュール・キャッシュ](#%E3%83%A2%E3%82%B8%E3%83%A5%E3%83%BC%E3%83%AB%E3%81%AE%E3%82%AD%E3%83%A3%E3%83%83%E3%82%B7%E3%83%A5%E5%85%88)から探す。

go get または go mod tidy コマンドを使ってあらかじめ外部パッケージをダウンロードしておくことで参照可能となる（[Go] 1.16 から[モジュールの自動ダウンロードは禁止](#go-mod-tidy-によるモジュール情報の更新)になった）。

### 「モジュール」とは

モジュール対応モードでは，標準ライブラリを除くパッケージを「モジュール（module）」として管理する。パッケージが単一のディレクトリを指すのに対し，モジュールは go.mod ファイルのあるディレクトリ以下の（go.mod を含まない）全てのパッケージがモジュールの配下となる。

パッケージが [git] 等のバージョン管理ツールで管理されている場合はバージョン（またはリビジョン）ごとに異なるモジュールと見なされる。つまりモジュールの実体は「パッケージ(s)＋バージョン」ということになる。

モジュールのバージョンがバージョン管理ツールのリビジョンと連動している関係上，基本的には「1リポジトリ＝1モジュール」である（ひとつのリポジトリに複数の go.mod ファイルを配置することで複数のモジュールを構成することは可能だが，管理がめちゃめちゃ煩雑になる）。

なお，モジュールのバージョンは go.mod ファイルで管理されるため，外部パッケージへの物理パスとモジュール名が同じであればソース・コードを書き換える必要はない。

## パッケージへの物理パスとモジュール名

たとえば，モジュール対応モードにおいて外部パッケージ github.com/spiegel-im-spiegel/pa-api/entity をインポートするには，ソース・コードにて以下のようにインポート宣言を記述する。

```go
import "github.com/spiegel-im-spiegel/pa-api/entity"
```

GOPATH モードではインポート宣言で指定されたパスがそのままパッケージへの物理パスを指しているが，モジュール対応モードではちょっと複雑な処理を行っている。具体的には

1. 宣言されたパスを解釈して [https://github.com/spiegel-im-spiegel/pa-api][spiegel-im-spiegel/pa-api] にあるリポジトリの指定リビジョンをフェッチする（リビジョンは自パッケージの go.mod ファイルで指定されたバージョンから類推する）
2. フェッチしたリポジトリにある go.mod ファイルからモジュール名 `github.com/spiegel-im-spiegel/pa-api` を取得する（go.mod ファイルがない場合は物理パスがそのままモジュール名となる）
3. 宣言されたパスとモジュール名からサブディレクトリ entity を該当のパッケージと解釈してインポートする

といった感じにパッケージの解釈とインポートを行う。したがって，モジュール名とパッケージへの物理パスはなるべく合わせた方が（外部パッケージを利用する側から見ても）面倒が少ない。ただし，後述する[バージョン管理](#semantic-versioning-%E3%81%AB%E3%82%88%E3%82%8B%E3%83%90%E3%83%BC%E3%82%B8%E3%83%A7%E3%83%B3%E7%AE%A1%E7%90%86)の都合上，どうしても両者が乖離してしまうことがある。

[spiegel-im-spiegel/pa-api]: https://github.com/spiegel-im-spiegel/pa-api "spiegel-im-spiegel/pa-api: APIs for Amazon Product Advertising API v5 by Golang"

## モジュール関連の環境変数

詳しくは以下が参考になる。

https://zenn.dev/tennashi/articles/3b87a8d924bc9c43573e

この記事ではモジュール管理に関連するもののみ挙げていく。

### 環境変数 GO111MODULE によるモードの切り替え

バージョン 1.11 以降では 2 つのモードの切り替えのために環境変数 GO111MODULE が用意されている。

```
$ go env | grep GO111MODULE
GO111MODULE=""
```

GO111MODULE の取りうる値は以下の通り。

| 値     | 内容 |
| ------ | ---- |
| `on`   | 常にモジュール対応モードで動作する |
| `off`  | 常に GOPATH モードで動作する |
| `auto` | $GOPATH/src 以下のディレクトリに配置され go.mod ファイルを含まないパッケージは GOPATH モードで，それ以外はモジュール対応モードで動作する |

バージョン 1.16 から GO111MODULE 未指定時の既定値が `on` になった（1.15 までは `auto`）。 GOPATH モードを使いたいのであれば GO111MODULE の値を `auto` または `off` に設定する[^env1]。

[^env1]: [Go] の環境変数の取り扱いについては，拙文「[Go 言語の環境変数管理](https://text.baldanders.info/golang/go-env/)」をご覧あれ。

```
$ go env -w GO111MODULE=auto
```

ただし [Go] では環境変数 GOPATH への依存を徐々に薄めつつあり，最終的には $GOPATH ディレクトリは使われなくなると思われる。

### パッケージのインストール先

[Go] では go get または go install コマンドでビルドした実行バイナリのインストール先を $GOPATH/bin ディレクトリに配置しているが，これを環境変数 GOBIN で変更することができる。

```
$ go env -w GOBIN=/home/username/bin
```

### モジュールのキャッシュ先

最近の [Go] ではコンパイル結果の中間バイナリを環境変数 GOCACHE で指定したディレクトリにキャッシュし，可能な限り再利用しようとする。

```
$ go env | grep GOCACHE
GOCACHE="/home/usernamee/.cache/go-build"
```

インポートしたモジュールのキャッシュについては $GOPATH/pkg/mod ディレクトリに配置されているが [Go] 1.15 より環境変数 GOMODCACHE で変更できるようになった。

```
$ go env | grep GOMODCACHE
GOMODCACHE="/home/usernamee/go/pkg/mod"

$ go env -w GOBIN=/home/usernamee/.cache/go-mod
```

なお，キャッシュのクリアは

```bash:ビルド・キャッシュのクリア
$ go clean -cache
```

または

```bash:モジュール・キャッシュのクリア
$ go clean -modcache
```

で可能である。

## go.mod と go.sum

既に何度も登場しているが，モジュールの管理は go.mod および go.sum ファイルで行う。

### go.mod の作成

新たに go.mod ファイルを作成するには，以下のコマンドを叩く。

```
$ go mod init github.com/spiegel-im-spiegel/pa-api
go: creating new go.mod: module github.com/spiegel-im-spiegel/pa-api
go: to add module requirements and sums:
	go mod tidy
```

これでモジュール名 github.com/spiegel-im-spiegel/pa-api としてカレント・ディレクトリ直下に go.mod ファイルが作成される。中身はこんな感じ。

```markup:go.mod
module github.com/spiegel-im-spiegel/pa-api

go 1.16
```

`module` や `go` はディレクティブ（directive）と呼ばれるものだ。たとえば `module` ディレクティブはモジュール名を定義する。他に go.mod ファイルで使えるディレクティブは以下の通り。

| ディレクティブ | 記述例                                          | 内容                   |
| -------------- | ----------------------------------------------- | ---------------------- |
| `module`       | `module my/thing`                               | モジュール名       |
| `go`           | `1.16`                                          | 有効な Go バージョン   |
| `require`      | `require other/thing v1.0.2`                    | インポート・モジュール |
| `exclude`      | `exclude old/thing v1.2.3`                      | 除外モジュール         |
| `replace`      | `replace bad/thing v1.4.5 => good/thing v1.4.5` | モジュールの置換       |
| `retract`      | `v1.0.5`                                        | 撤回バージョン         |

`replace` ディレクティブはモジュール名が物理パスと対応してなくて上手くインポートできない等の状況で使える。たとえばこんな感じ。

```markup:go.mod
module sample

require gopkg.in/russross/blackfriday.v2 v2.0.1

replace gopkg.in/russross/blackfriday.v2 v2.0.1 => github.com/russross/blackfriday/v2 v2.0.1
```

また，同一リポジトリ内に複数のモジュールがある場合に `replace` ディレクティブを使って相対パス指定ができるようだ。

```markup:go.mod
module golang.org/x/tools/gopls

replace golang.org/x/tools => ../
```

`retract` ディレクティブはバージョン 1.16 から導入されたもので

```markup:go.mod
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

go.sum ファイルにはインポートするモジュールの SHA-256 チェックサム値が格納されている。たとえば go.mod ファイルで `require` ディレクティブが

```markup:go.mod
require github.com/spiegel-im-spiegel/errs v1.0.2
```

と指定されている場合， go.sum ファイルの内容は

```markup:go.sum
github.com/spiegel-im-spiegel/errs v1.0.2 h1:v4amEwRDqRWjKHOILQnJSovYhZ4ZttEnBBXNXEzS6Sc=
github.com/spiegel-im-spiegel/errs v1.0.2/go.mod h1:UoasJYYujMcdkbT9USv8dfZWoMyaY3btqQxoLJImw0A=
```

といった感じになる。

go.sum ファイルの内容はインポートするモジュールの完全性（integrity）を担保するものだが[^sum1]，手作業で記述できるようなものではないので，次に紹介する go mod tidy コマンドを使って更新する。

[^sum1]: [Go] が go.sum ファイルを使って完全性をどのように管理しているかについては拙文「[Go モジュールのミラーリング・サービス【正式版】](https://text.baldanders.info/golang/mirror-index-and-checksum-database-for-go-module/)」を参考にどうぞ。

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

```markup
$ go mod tidy
```

とするとよい。

## [Semantic Versioning] によるバージョン管理

モジュールのバージョンはリポジトリのリビジョン番号またはバージョンタグによって管理されるが，バージョンタグに関しては [Semantic Versioning] のルールに則ってバージョン番号を設定することが推奨されている。

![research.swtch.com/impver.png](https://research.swtch.com/impver.png)
_[via “Semantic Import Versioning”](https://research.swtch.com/vgo-import "Semantic Import Versioning")_

このように後方互換性のない変更がある場合にはメジャーバージョンを，後方互換性が担保された変更や追加についてはマイナーバージョンを，不具合や脆弱性の修正については第 3 位のパッチバージョンを上げるようにする。またメジャーバージョンを上げる際には，図のようにディレクトリを分離するか， go.mod ファイルの `module` ディレクティブの値を変更するのが簡単である[^ver1]。

[^ver1]: v0 から v1 へのメジャーバージョンアップの場合は例外的に後方互換性は考慮されない。 v0 はベータ版という扱いらしい。それ以外のメジャーバージョンアップの際にモジュール名＋パスで区別せずにバージョンタグを打つと強制的に `v2.0.0+incompatible` みたいなショボい表記にされる。

以下は [github.com/mattn/jvgrep](https://github.com/mattn/jvgrep) の例：

```markup:go.mod
module github.com/mattn/jvgrep/v5
```

このとき，ソース・コード側も

```go
import "github.com/mattn/jvgrep/v5/mmap"
```

などとモジュール名をベースに指定する必要がある[^pasth1]。

[^pasth1]: バージョン 1.16 から import 時の相対パス指定は原則禁止になったので注意。同一リポジトリ内に複数のモジュールがある場合は go.mod ファイルで `replace` ディレクティブを使うとよい。

## 特定バージョンのモジュールをビルド&インストールする

バージョン 1.16 から go install コマンドで モジュール@バージョン を指定できるようになった。

```markup
$ go install golang.org/x/tools/gopls@v0.6.5
```

とにかく最新版が欲しい場合は

```markup
$ go install golang.org/x/tools/gopls@latest
go: downloading golang.org/x/tools/gopls v0.6.5
...
```

などとすれば OK。

バージョン付きで指定する際は，モジュール名で指定する点に注意。たとえば，先ほどの [github.com/mattn/jvgrep](https://github.com/mattn/jvgrep) の場合，パッケージへの物理パスで指定しても

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

なお，バージョン 1.16 では go.mod ファイルに `replace` や `exclude` ディレクティブが含まれていると go install に失敗することがあるみたい。

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

### go get はオワコン？

go get コマンドは元々 $GOPATH ディレクトリ下に指定した外部パッケージを組み込むための仕組みである。

旧来の GOPATH モードでは go get で常に最新リビジョンのパッケージを $GOPATH ディレクトリ下にダウンロードしようとするため「[GOPATH 汚染](https://text.baldanders.info/golang/gopath-pollution/)」の問題がつきまとっていた。現在のモジュール対応モードはこの問題を根本的に解決するためのものと言える。

しかしその結果，モジュール対応モードでは go get コマンドはモジュールをキャッシュするだけのコマンドになってしまった。しかも go.mod & go.sum ファイルを意図せず書き換えてしまう危険がある。

開発中のパッケージにおいて，意図的に（依存モジュールを含む）外部モジュールの読み込みと go.mod & go.sum ファイルの更新を行いたいのであれば go mod tidy で一括処理するほうが便利である。また単にパッケージ（またはモジュール）をビルド&インストールしたいだけであれば [go install を使うほうが不用意に go.mod & go.sum ファイルを汚すこともない](https://zenn.dev/minguu/articles/20210225-difference-of-go-install-and-go-get "go installとgo getをどう使い分ければ良いのか調べた")ので安全と言える。

GOPATH モードが後方互換機能として残されている間は go get コマンド自体もなくならないだろうが，モジュール対応モードが主流となる今後はユーザが手打ちで go get コマンドを叩くことはなくなってくるんじゃないだろうか。

## 参考

https://golang.org/ref/mod
https://golang.org/doc/modules/developing
https://blog.golang.org/go116-module-changes
https://zenn.dev/nobonobo/articles/4fb018a24f9ee9

[go]: https://golang.org/ "The Go Programming Language"
[git]: https://git-scm.com/ "Git"
[semantic versioning]: http://semver.org/ "Semantic Versioning 2.0.0 | Semantic Versioning"
[github]: https://github.com/
