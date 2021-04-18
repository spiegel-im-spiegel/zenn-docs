---
title: "GOPATH 跡地" # 記事のタイトル
emoji: "🤔" # アイキャッチとして使われる絵文字（1文字だけ）
type: "idea" # "tech" : 技術記事 / "idea" : アイデア記事
topics: ["go", "programming"] # タグ。["markdown", "rust", "aws"] のように指定する
published: false # 公開設定（true で公開）
---

[Go] 1.16 が登場してから2ヶ月経つのだが，どうもネット上の GOPATH モードを前提とした紹介・解説記事を見てハマってるパターンが多いように見受ける。そこで，改めて GOPATH について「こんな時代もあったね」という感じで紹介することにする。

## GOPATH モードに戻したい

2018年8月にリリースされたバージョン 1.11 以降， [Go] ツールーチェーンは以下の2つのモードのどちらかで動作する。

- **GOPATH モード (GOPATH mode)** : バージョン 1.10 までのモード。標準ライブラリを除く全てのパッケージのコード管理とビルドを環境変数 GOPATH で指定されたディレクトリ下で行う。パッケージの管理はリポジトリの最新リビジョンのみが対象となる
- **モジュール対応モード (module-aware mode)** : 標準ライブラリを除く全てのパッケージをモジュールとして管理する。コード管理とビルドは任意のディレクトリで可能で，モジュールはリポジトリのバージョンタグまたはリビジョン毎に管理される

言い方を変えると2018年8月より前の解説記事は GOPATH モードを前提に書かれている可能性が高いので注意が必要である。

2021年2月にリリースされた [Go] 1.16 ではモジュール対応モードが既定のモードとなった。 1.16 へのバージョンアップでよく見かけるトラブルが，この既定モードの変更に起因するもののようだ。 GOPATH モードに戻したいのであれば，環境変数 GO111MODULE の値を `off` または `auto` に変更する[^env1]。

[^env1]: 環境変数 GO111MODULE の既定値は `on`。ちなみに [Go] 1.15 までの既定値は `auto` だった。なお [Go] の環境変数の取り扱いについては，拙文「[Go 言語の環境変数管理](https://text.baldanders.info/golang/go-env/)」をご覧あれ。


```
$ go env -w GO111MODULE=auto
```

変更した環境変数を取り消すには

```
$ go env -u GO111MODULE
```

でOK。各値の意味は以下のとおり。

| 値     | 内容 |
| ------ | ---- |
| `on`   | 常にモジュール対応モードで動作する |
| `off`  | 常に GOPATH モードで動作する |
| `auto` | $GOPATH/src 以下のディレクトリに配置され go.mod ファイルを含まないパッケージは GOPATH モードで，それ以外はモジュール対応モードで動作する |

### モジュール間の循環依存に注意

[Go] では import ディレクティブによる循環参照を禁止している。適切な例が思い浮かばなくて申し訳ないが，たとえば以下の2つのコードがあるとする。

```go:gopath/src/err/say/say.go
package say

import (
	"err/hello"
	"fmt"
)

func Hello() {
	fmt.Println(hello.Hello)
}
```

```go:gopath/src/err/hello/hello.go
package hello

import "err/say"

var Hello = "Hello, World"

func Hello() {
	say.Hello()
}
```

見ての通りお互いがお互いのパッケージを import しているため，これを実行しようとすると

```
$ go run main.go
package command-line-arguments
        imports err/hello
        imports err/say
        imports err/hello: import cycle not allowed
```

という感じにコンパイルエラーになる。

元々 [Go] はコンパイラの高速化のために上のような循環参照を禁じているのだが，コンポーネント設計における「非循環依存関係の原則（Acyclic Dependencies Principle; ADP）」の観点からも理に適ってると言える。

ただし実際には import ディレクティブでは現れてこないが，コンポーネント図等で見ると循環依存の関係になっている場合がある。しかもそれらが別々の [Go] モジュールになっていると管理がとても面倒になる。

最近見かけた例では以下のような話があった。

https://zenn.dev/ikawaha/articles/b840cc185f468997f74c

この例では，最終的に依存関係を見直して循環依存の関係をなくすことで対応したようだ。

https://zenn.dev/ikawaha/articles/1b22e0eae6ce4edc3e90

このようにモジュール対応モードでパッケージ管理がうまくいかない場合に，一時的にでも GOPATH モードに戻すというのはアリだと思う。

## 参考

https://zenn.dev/spiegel/articles/20210223-go-module-aware-mode

[Go]: https://golang.org/ "The Go Programming Language"
<!-- eof -->
