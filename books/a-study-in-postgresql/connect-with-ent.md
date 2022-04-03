---
title: "ent で PostgreSQL に接続できるか"
---

[ent] は Facebook 主管で開発が行われている ORM (Object Relational Mapper) ハンドリング・パッケージである。古いブログ記事だと

```go
import "github.com/facebook/ent"
```

などと書かれていることが多いが，現在のリポジトリは [github.com/ent/ent] に移管されていて，インポートする際も

```go
import "entgo.io/ent"
```

などとするようだ。更に

https://entgo.io/blog/2021/09/01/ent-joins-the-linux-foundation

といったこともあったようで Facebook から中立の立場をとれる [ent] には個人的に注目している。

[GORM] が Model 構造体を interface{} 型で受けて内部で refrect 等を使って頑張っているのに対し [ent] は go generate を使って必要なコードを自動生成するという戦略をとる。コード生成およびそのための定義の記述が若干（いやだいぶ？）手間ではあるが，利用者から見て型安全なコードが書けるのは大きなメリットである。

[ent] はこれまでの SQL ベースの ORM とは一線を画すもので，どちらかというと「データ間のグラフ構造を記述・制御するフレームワーク」といった感じだろうか。近年注目されている GraphQL 等を扱いたい場合には有用だろう。ただしこの本では，あくまでも [ent] と [PostgreSQL] との連携と SQL ベースでのデータアクセスを主軸に記述していく。

[ent] が公式にサポートしているのは MySQL/MariaDB, [PostgreSQL], SQLite, Gremlin (AWS Neptune) といった辺りのようだ（Gremlin は実験的な実装）。 *[sql][database/sql].DB 型のインスタンスを受け付けるので [github.com/jackc/pgx] パッケージも使える。こんな感じ。

```go
//go:build run
// +build run

package main

import (
    "os"
    "sample/loggr"

    "entgo.io/ent/dialect"
    "entgo.io/ent/dialect/sql"
    "entgo.io/ent/examples/start/ent"
    "github.com/jackc/pgx/v4"
    "github.com/jackc/pgx/v4/log/zerologadapter"
    "github.com/jackc/pgx/v4/stdlib"
    "github.com/goark/errs"
    "github.com/goark/gocli/exitcode"
)

func Run() exitcode.ExitCode {
    // get logger
    zlogger := loggr.New()

    // create ent.Client instance for PostgreSQL service
    cfg, err := pgx.ParseConfig(os.Getenv("ELEPHANTSQL_URL"))
    if err != nil {
        zlogger.Error().Interface("error", errs.Wrap(err)).Send()
        return exitcode.Abnormal
    }
    cfg.Logger = zerologadapter.NewLogger(*zlogger)
    cfg.LogLevel = pgx.LogLevelDebug
    client := ent.NewClient(
        ent.Driver(
            sql.OpenDB(dialect.Postgres, stdlib.OpenDB(*cfg)),
        ),
    )
    defer client.Close()

    return exitcode.Normal
}

func main() {
    Run().Exit()
}
```

entgo.io/ent/examples/start/ent パッケージって何？ と思うかもしれないが，これに関しては次節で説明する。とりあえず繋がりますよー，ということで。

[ent]: https://entgo.io/
[GORM]: https://gorm.io/ "GORM - The fantastic ORM library for Golang, aims to be developer friendly."
[Go]: https://go.dev/
[github.com/ent/ent]: https://github.com/ent/ent "ent/ent: An entity framework for Go"
[PostgreSQL]: https://www.postgresql.org/ "PostgreSQL: The world's most advanced open source database"
[database/sql]: https://pkg.go.dev/database/sql "sql package - database/sql - pkg.go.dev"
[github.com/jackc/pgx]: https://github.com/jackc/pgx "jackc/pgx: PostgreSQL driver and toolkit for Go"
