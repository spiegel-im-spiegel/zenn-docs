---
title: "単一リポジトリで Workspace モードを試してみる" # 記事のタイトル
emoji: "💻" # アイキャッチとして使われる絵文字（1文字だけ）
type: "tech" # "tech" : 技術記事 / "idea" : アイデア記事
topics: ["go", "programming"] # タグ。["markdown", "rust", "aws"] のように指定する
published: true # 公開設定（true で公開）
---

## 単一リポジトリで複数モジュールを扱いたい

週末に SSH 越しに RDBMS サービスにアクセスする [Go] パッケージをリリースしたのだが

https://github.com/goark/sshql
https://text.baldanders.info/release/2022/09/sql-over-ssh/

このパッケージは3つのモジュールで構成されている。ディレクトリ構造はこんな感じ。

```
$ tree sshql
sshql
├── go.mod
├── go.sum
├── mysqldrv
│   ├── go.mod
│   └── go.sum
└── pgdrv
    ├── go.mod
    └── go.sum
```

ちなみに mysqldrv および pgdrv パッケージは親の sshql パッケージに依存している。これら3つのパッケージを別々のモジュールにしたのは，ルートの go.mod ファイルに各ドライバの外部パッケージを混ぜたくなかったから。モジュールを分けることで

https://github.com/goark/sshql/blob/main/go.mod
https://github.com/goark/sshql/blob/main/mysqldrv/go.mod
https://github.com/goark/sshql/blob/main/pgdrv/go.mod

という感じにドライバ用のパッケージをきれいに分離することができる。

問題はこの3つのモジュールが（見かけ上）独立しているためモジュール間で import する際にローカル・ファイル群ではなく GitHub リポジトリを参照してしまうことにある。これは開発環境としては具合が悪い。

## Workspace モードを試す

そこで [Go] 1.18 から登場した workspace モードを試すこととにした。ただ，以下のドキュメントを読む限り workspace 管理用の親ディレクトリを作ってその下にモジュールを配置していく，というのがセオリーのようだ。

https://go.dev/doc/tutorial/workspaces
https://go.dev/blog/get-familiar-with-workspaces
https://future-architect.github.io/articles/20220216a/

これでは今回のケースでは使えない。でも，まぁ，色々試してみよう。

まずリポジトリのルート直下に go.work ファイルを作ってみる。こんな感じでどうだろう。

```
$ go work init . mysqldrv pgdrv
```

これで

https://github.com/goark/sshql/blob/main/go.work

というファイルができた。どうやらこれでカレントの sshql モジュールと配下の mysqldrv, pgdrv モジュールを認識することができるっぽい。

これで例えば

```
$ go test -shuffle on ./pgdrv/...
```

とかいった感じにモジュール毎にテストを実行できる。

### もし go.work ファイルがないと...

もし go.work ファイルがない状態でテストを実行しようとしても

```
$ go test -shuffle on ./pgdrv/...
pattern ./pgdrv/...: directory prefix pgdrv does not contain main module or its selected dependencies
```

てな感じに怒られる。まぁ

```
$ pushd pgdrv
$ go test -shuffle on ./...
```

とかすれば一応通るのだが，今度は親モジュールを GitHub リポジトリからインポートしようとする。ままならぬ。この辺は本当に色々と試行錯誤した。



## サブモジュール毎にバージョンを振る

気を取り直して。

go.work を設置することでローカル環境で閉じた開発ができるようになった。またリポジトリに go.work を含めることで GitHub Actions で複数モジュールを扱うこともできるようになる。

一応それっぽくできたのでバージョンタグを打つのだが，たとえば `v0.1.0` などとタグを打ってもこのバージョンタグはリポジトリ直下のモジュールにしか有効にならないようだ。今回のケースではサブモジュールの mysqldrv, pgdrv には適用されず `v0.0.0-20220910000000-abcde012345` みたいなダサいバージョンにされてしまう。

サブモジュールにバージョンタグを打つには `mysqldrv/v0.1.0` という感じに頭にパッケージ名を付ければいいようだ。これで [Go] はサブモジュールのバージョンと認識してくれるらしい。別プロジェクトで `go mod tidy` とかすると

```:go.mod
module sample

go 1.19

require (
    github.com/goark/sshql v0.1.0
    github.com/goark/sshql/mysqldrv v0.1.0
)

require (
    github.com/go-sql-driver/mysql v1.6.0 // indirect
    github.com/goark/errs v1.1.0 // indirect
    golang.org/x/crypto v0.0.0-20220829220503-c86fa9a7ed90 // indirect
    golang.org/x/sys v0.0.0-20210615035016-665e8c7367d1 // indirect
)
```

という感じに mysqldrv モジュールにもバージョンが付くようになる。またモジュール毎に独立にバージョン管理ができるのもありがたい。

[Go]: https://go.dev/ "The Go Programming Language"
<!-- eof -->
