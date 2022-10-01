---
title: "【付録2】 kra で sql をラップする"
---

[2022年春の Go Conference](https://gocon.connpass.com/event/212162/ "Go Conference 2022 Spring Online - connpass") でとても感銘を受けたのが [github.com/taichi/kra] パッケージである。

https://github.com/taichi/kra
https://gocon.jp/2022spring/ja/sessions/b7-l/

[Kra][github.com/taichi/kra] は ORM ではなく，標準 [database/sql] をラップするヘルパー・パッケージという位置付けになるだろう。個人的には [context].Context の強制やプレースホルダに対応するパラメータや Scan() メソッドの受け皿に（構造体などの）構造化されたデータを使う点が気に入っている。

[github.com/jackc/pgx] に対応する専用ドライバも用意されていて

```go
package main

import (
    "context"
    "fmt"
    "os"

    "github.com/taichi/kra/pgx"
)

func main() {
    db, err := pgx.Open(context.TODO(), "postgres://dbuser:dbpassword@dbserver:5432/example?sslmode=require")
    if err != nil {
        fmt.Fprintln(os.Stderr, err)
        return
    }
    defer db.Close()

    ...
}
```

という感じに [database/sql] から置き換えて使うことができる。

ただし，専用 logger を組み込むとかいった要件がある場合には，ちょっと記述が複雑になる。たとえば [zerolog][github.com/rs/zerolog] および[前節]の [sshql][github.com/goark/sshql] パッケージを組み合わせるとこんな感じ。

```diff go
package main

import (
    "context"
    "fmt"
    "os"

    "github.com/goark/errs"
    "github.com/goark/sshql"
    "github.com/jackc/pgx/v4"
    "github.com/jackc/pgx/v4/log/zerologadapter"
    "github.com/jackc/pgx/v4/pgxpool"
    "github.com/rs/zerolog"
    "github.com/taichi/kra"
    krapgx "github.com/taichi/kra/pgx"
)

func main() {
    dialer := &sshql.Dialer{
        Hostname:   "sshserver",
        Port:       22,
        Username:   "remoteuser",
        Password:   "passphraseforauthkey",
        PrivateKey: "/home/username/.ssh/id_eddsa",
    }
    logger := zerolog.New(os.Stderr)

    cfg, err := pgxpool.ParseConfig("postgres://dbuser:dbpassword@localhost:5432/example?sslmode=disable")
    if err != nil {
        logger.Error().Interface("error", errs.Wrap(err)).Send()
        return
    }
    cfg.ConnConfig.Logger = zerologadapter.NewLogger(logger)
    cfg.ConnConfig.LogLevel = pgx.LogLevelDebug
    cfg.ConnConfig.DialFunc = dialer.DialContext

    if err := dialer.Connect(); err != nil {
        logger.Error().Interface("error", errs.Wrap(err)).Send()
        return
    }
    defer dialer.Close()

    pool, err := pgxpool.ConnectConfig(context.TODO(), cfg)
    if err != nil {
        logger.Error().Interface("error", errs.Wrap(err)).Send()
        return
    }
+   db := krapgx.NewDB(pool, krapgx.NewCore(kra.NewCore(kra.PostgreSQL)))
+   defer db.Close()

    ...
}
```

まぁ，複雑なのは [pgx][github.com/jackc/pgx]/pgxpool.Config 構造体を組み立てる部分と SSH への接続部分が余計にくっ付いてる部分だけで，あとは [pgx][github.com/jackc/pgx] のコネクションプールを [kra][github.com/taichi/kra]/pgx.NewDB() 関数でラップしてる（色付きの行）だけなのだが。

## 適材適所

突然だが，データベースの設計要件には大きく2種類ある。

1. データ構成とその関係をいちから設計できる場合
2. 既にあるデータ構造を使って機能を追加・変更あるいはリプレースする場合

前者の場合は [ent] のようなフレームワークが有効だろう。特に最初の頃はデータベース周りはサービスごと変わることも多く安定性が低い。であれば，フレームワークでデータベース周りを抽象化・隠蔽できるのは有用である。

一方で，後者の場合はデータベース周りを大きく変更できないので，ビジネスロジック側でデータベースの都合に合わせることになる。データベースのほうが安定性が高くなってしまうのだ。しかも，こういう状況では汎用のフレームワークやクエリビルダでは設計要件に上手くマッチしないことが多く，結局は「自分で書いたほうが速い」になってしまう。

またバッチ処理などで大量のデータを動かす場合は [PostgreSQL] の COPY FROM 文や MySQL の LOAD DATA INFILE 文などの特殊構文を使うことが多く[^mysql1]，他の部分で ORM やクエリビルダを使っていてもそこで抽象化や隠蔽が崩れてしまうわけで，それなら下手な隠蔽などせず SQL 文も最初からベタ書きで設計し，クエリプランを使って事前に評価・最適化すればいい，となりがちである。

[^mysql1]: MySQL の LOAD DATA INFILE 文を使った事例については自ブログの「[SSH, MySQL, Zerolog, そして Kra](https://text.baldanders.info/golang/ssh-mysql-zerolog-and-kra/)」で紹介している。

かなり愚痴になってしまったが，私達はエンジニアであって書道家ではない。道具は状況に合わせて最適なものを選ぶべきである（選ぶ自由が欲しい）。[kra][github.com/taichi/kra] はそうした道具のひとつとして優れていると私は思う。

[前節]: https://zenn.dev/spiegel/books/a-study-in-postgresql/viewer/connect-over-ssh "SSH 越しに PostgreSQL に接続する"
[Go]: https://go.dev/
[PostgreSQL]: https://www.postgresql.org/ "PostgreSQL: The world's most advanced open source database"
[ent]: https://entgo.io/
[database/sql]: https://pkg.go.dev/database/sql "sql package - database/sql - pkg.go.dev"
[context]: https://pkg.go.dev/context "context package - context - pkg.go.dev"
[github.com/jackc/pgx]: https://github.com/jackc/pgx "jackc/pgx: PostgreSQL driver and toolkit for Go"
[github.com/rs/zerolog]: https://github.com/rs/zerolog "rs/zerolog: Zero Allocation JSON Logger"
[github.com/goark/sshql]: https://github.com/goark/sshql "goark/sshql: Go SQL drivers over SSH"
[github.com/taichi/kra]: https://github.com/taichi/kra "taichi/kra: relational database access helper library"
