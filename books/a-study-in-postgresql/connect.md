---
title: "PostgreSQL に接続する"
---

では実際に PostgreSQL サービスに接続してみよう。

[Go] の標準パッケージである [database/sql] を使って RDBMS に接続するには DB の種別ごとに「ドライバ」と呼ばれるパッケージをインポートする必要がある。 [PostgreSQL] の場合は [github.com/lib/pq] が定番として使われることが多かった。

```go:proto/sample2.go
//go:build run
// +build run

package main

import (
    "database/sql"
    "log"
    "os"

    "github.com/joho/godotenv"
    _ "github.com/lib/pq"
    "github.com/spiegel-im-spiegel/gocli/config"
    "github.com/spiegel-im-spiegel/gocli/exitcode"
)

func init() {
    //load ~/.config/elephantsql/env file
    if err := godotenv.Load(config.Path("elephantsql", "env")); err != nil {
        panic(err)
    }
}

func Run() exitcode.ExitCode {
    // create sql.DB instance for PostgreSQL service
    db, err := sql.Open("postgres", os.Getenv("ELEPHANTSQL_URL"))
    if err != nil {
        log.Println(err)
        return exitcode.Abnormal
    }
    defer db.Close()

    return exitcode.Normal
}

func main() {
    Run().Exit()
}
```

ただ，最近は [github.com/jackc/pgx] パッケージのほうが定番になりつつあるようだ（以下 import と Run() 関数のみ挙げておく）。

```go:proto/sample2b.go
import (
    "database/sql"
    "log"
    "os"

    _ "github.com/jackc/pgx/v4/stdlib"
    "github.com/joho/godotenv"
    "github.com/spiegel-im-spiegel/gocli/config"
    "github.com/spiegel-im-spiegel/gocli/exitcode"
)

func Run() exitcode.ExitCode {
    // create sql.DB instance for PostgreSQL service
    db, err := sql.Open("pgx", os.Getenv("ELEPHANTSQL_URL"))
    if err != nil {
        log.Println(err)
        return exitcode.Abnormal
    }
    defer db.Close()

    return exitcode.Normal
}
```

[database/sql] を使わず [github.com/jackc/pgx] パッケージを直に使って接続することもできる。

```go:proto/sample3.go
import (
    "context"
    "log"
    "os"

    "github.com/jackc/pgx/v4/pgxpool"
    "github.com/joho/godotenv"
    "github.com/spiegel-im-spiegel/gocli/config"
    "github.com/spiegel-im-spiegel/gocli/exitcode"
)

func Run() exitcode.ExitCode {
    // create connection pool for PostgreSQL service
    pool, err := pgxpool.Connect(context.TODO(), os.Getenv("ELEPHANTSQL_URL"))
    if err != nil {
        log.Println(err)
        return exitcode.Abnormal
    }
    defer pool.Close()

    return exitcode.Normal
}
```

pgx.Connect() 関数ではなく pgxpool.Connect() 関数を使って接続している点に注意。 pgx.Connect() 関数で得られる pgx.Conn インスタンスは単一接続を表し並行的に安全（concurrency safe）ではないそうな。

> \`*pgx.Conn\` represents a single connection to the database and is not concurrency safe. Use sub-package pgxpool for a concurrency safe connection pool.
(via “[pgx package - github.com/jackc/pgx/v4 - pkg.go.dev](https://pkg.go.dev/github.com/jackc/pgx/v4#hdr-Connection_Pool)”)

pgxpool.Connect() 関数の引数に [context].Context インスタンスが渡されている点にも注意。 [github.com/jackc/pgx] パッケージではほとんどの処理で [context].Context インスタンスを渡すようになっている。今時ですな。

[github.com/jackc/pgx] パッケージ自体はかなり高機能なのだが，世の中にある [Go] 用の ORM (Object Relational Mapper) フレームワークの大抵は [database/sql] 標準パッケージを前提にしているのが辛いところである。

[Go]: https://go.dev/
[PostgreSQL]: https://www.postgresql.org/ "PostgreSQL: The world's most advanced open source database"
[ElephantSQL]: https://www.elephantsql.com/ "ElephantSQL - PostgreSQL as a Service"
[database/sql]: https://pkg.go.dev/database/sql "sql package - database/sql - pkg.go.dev"
[os]: https://pkg.go.dev/os "os package - os - pkg.go.dev"
[context]: https://pkg.go.dev/context "context package - context - pkg.go.dev"
[github.com/lib/pq]: https://github.com/lib/pq "lib/pq: Pure Go Postgres driver for database/sql"
[github.com/jackc/pgx]: https://github.com/jackc/pgx "jackc/pgx: PostgreSQL driver and toolkit for Go"
