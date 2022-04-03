---
title: "ent に zerolog を組み込むよ"
---

自動生成したコードを眺めてみると

```go:ent/client.go
// NewClient creates a new client configured with the given options.
func NewClient(opts ...Option) *Client {
    cfg := config{log: log.Println, hooks: &hooks{}}
    cfg.options(opts...)
    client := &Client{config: cfg}
    client.init()
    return client
}
```

とか書いてある。「なんじゃこらー」って思わなかった？ 思ったよね。更に config 構造体を定義しているところを見ると

```go:ent/config.go
// Config is the configuration for the client and its builder.
type config struct {
    // driver used for executing database requests.
    driver dialect.Driver
    // debug enable a debug logging.
    debug bool
    // log used for logging on debug mode.
    log func(...interface{})
    // hooks to execute on mutations.
    hooks *hooks
}
```

とあって，さらに

```go:ent/config.go
// Options applies the options on the config object.
func (c *config) options(opts ...Option) {
    for _, opt := range opts {
        opt(c)
    }
    if c.debug {
        c.driver = dialect.Debug(c.driver, c.log)
    }
}

// Debug enables debug logging on the ent.Driver.
func Debug() Option {
    return func(c *config) {
        c.debug = true
    }
}

// Log sets the logging function for debug mode.
func Log(fn func(...interface{})) Option {
    return func(c *config) {
        c.log = fn
    }
}
```

と書かれている。ふむむー。つまり func(...interface{}) 関数型にマッチした関数であれば自前で logger を組み込めるわけだ。これなら [zerolog][github.com/rs/zerolog] もなんとか組み込み可能だよね。しかし，凄い雑な構造だな。まっ楽でいいけど（笑）

以上を踏まえて [ent] 用の [PostgreSQL] サービス接続関数を書いてしまおう。こんな感じでどうだろうか。

```go:dbconn/entclient.go
package dbconn

import (
    "fmt"
    "sample/ent"
    "sample/env"

    "entgo.io/ent/dialect"
    "entgo.io/ent/dialect/sql"
    "github.com/rs/zerolog"
    "github.com/goark/errs"
)

type EntContext struct {
    Client *ent.Client
    Logger *zerolog.Logger
}

func NewEnt() (*EntContext, error) {
    pgxCtx, err := NewPgx()
    if err != nil {
        return nil, errs.Wrap(err)
    }
    entCtx := &EntContext{
        Logger: pgxCtx.GetLogger(),
    }
    entCtx.Client = ent.NewClient(
        ent.Driver(
            sql.OpenDB(dialect.Postgres, pgxCtx.GetDb()),
        ),
        ent.Log(func(v ...interface{}) {
            entCtx.Logger.Debug().Msg(fmt.Sprint(v...))
        }),
    )
    if env.LogLevel() >= env.LevelDebug {
        entCtx.Client = entCtx.Client.Debug()
    }
    return entCtx, nil
}

func (entCtx *EntContext) GetClient() *ent.Client {
    if entCtx == nil {
        return nil
    }
    return entCtx.Client
}

func (entCtx *EntContext) GetLogger() *zerolog.Logger {
    if entCtx == nil {
        lggr := zerolog.Nop()
        return &lggr
    }
    return entCtx.Logger
}

func (entCtx *EntContext) Close() error {
    if client := entCtx.GetClient(); client != nil {
        return errs.Wrap(client.Close())
    }
    return nil
}
```

これを使って


```go:sample1.go
//go:build run
// +build run

package main

import (
    "context"
    "fmt"
    "os"
    "sample/dbconn"

    "github.com/goark/errs"
    "github.com/goark/gocli/exitcode"
)

func Run() exitcode.ExitCode {
    // get ent context
    entCtx, err := dbconn.NewEnt()
    if err != nil {
        fmt.Fprintln(os.Stderr, err)
        return exitcode.Abnormal
    }
    defer entCtx.Close()

    // query
    if _, err := entCtx.GetClient().User.Get(context.TODO(), 1); err != nil {
        entCtx.GetLogger().Error().Interface("error", errs.Wrap(err)).Send()
        return exitcode.Abnormal
    }

    return exitcode.Normal
}

func main() {
    Run().Exit()
}
```

と書いて実行すれば

```
$ go run sample/sample1.go 
0:00AM DBG driver.Query: query=SELECT DISTINCT "users"."id", "users"."username", "users"."created_at", "users"."updated_at" FROM "users" WHERE "users"."id" = $1 LIMIT 2 args=[1]
0:00AM INF Dialing PostgreSQL server host=hostname module=pgx
0:00AM ERR Query args=[1] err="ERROR: relation \"users\" does not exist (SQLSTATE 42P01)" module=pgx pid=28859 sql="SELECT DISTINCT \"users\".\"id\", \"users\".\"username\", \"users\".\"created_at\", \"users\".\"updated_at\" FROM \"users\" WHERE \"users\".\"id\" = $1 LIMIT 2"
0:00AM ERR  error={"Context":{"function":"main.Run"},"Err":{"Msg":"ERROR: relation \"users\" does not exist (SQLSTATE 42P01)","Type":"*pgconn.PgError"},"Type":"*errs.Error"}
0:00AM INF closed connection module=pgx pid=28859
```

と表示された（データベースは ent の調査用に初期化したのでエラーになる）。うんうん。上から下まで統一した logger だと気持ちがいいねぇ。

[github.com/rs/zerolog]: https://github.com/rs/zerolog "rs/zerolog: Zero Allocation JSON Logger"
[ent]: https://entgo.io/
[PostgreSQL]: https://www.postgresql.org/ "PostgreSQL: The world's most advanced open source database"
[Go]: https://go.dev/
[github.com/ent/ent]: https://github.com/ent/ent "ent/ent: An entity framework for Go"
[database/sql]: https://pkg.go.dev/database/sql "sql package - database/sql - pkg.go.dev"
[github.com/jackc/pgx]: https://github.com/jackc/pgx "jackc/pgx: PostgreSQL driver and toolkit for Go"
