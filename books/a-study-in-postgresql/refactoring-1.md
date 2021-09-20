---
title: "リファクタリング（その1）"
---

さて，コードがかなりカオスになってきたので，ここらで一発，リファクタリングをかましておこう。この節は長いよ。コードに興味のない方は適当に読み飛ばしてよい（笑）

そうそう。今回書いたコードは

```
$ go mod init sample
```

で初期化したディレクトリ内にあるので，サブパッケージをインポートする際には

```go
import "sample/env"
```

という記述になっている。あしからずご了承の程を。

## env サブパッケージ

まず env ファイルの中身を以下のようにセットする。

```ini:$XDG_CONFIG_HOME/elephantsql/env
ELEPHANTSQL_URL=postgres://username:password@hostname:port/databasename
ENABLE_LOGFILE=true
LOGLEVEL=info
```

`ENABLE_LOGFILE` が `true` ならファイルに構造化ログを出力し，標準出力へは [zerolog][github.com/rs/zerolog].ConsoleWriter を使う。また `LOGLEVEL` でログの最低出力レベルを指定する。

`LOGLEVEL` を解釈するために以下の列挙型を導入する。

```go:env/loglevel.go
package env

import (
    "strings"

    "github.com/jackc/pgx/v4"
    "github.com/rs/zerolog"
)

type LoggerLevel int

const (
    LevelNop LoggerLevel = iota
    LevelError
    LevelWarn
    LevelInfo
    LevelDebug
)

var levelMap = map[LoggerLevel]string{
    LevelNop:   "nop",
    LevelError: "error",
    LevelWarn:  "warn",
    LevelInfo:  "info",
    LevelDebug: "debug",
}

var zerologLevelMap = map[LoggerLevel]zerolog.Level{
    LevelNop:   zerolog.NoLevel,
    LevelError: zerolog.ErrorLevel,
    LevelWarn:  zerolog.WarnLevel,
    LevelInfo:  zerolog.InfoLevel,
    LevelDebug: zerolog.DebugLevel,
}

var pgxlogLevelMap = map[LoggerLevel]pgx.LogLevel{
    LevelNop:   pgx.LogLevelNone,
    LevelError: pgx.LogLevelError,
    LevelWarn:  pgx.LogLevelWarn,
    LevelInfo:  pgx.LogLevelInfo,
    LevelDebug: pgx.LogLevelDebug,
}

func getLogLevel(s string) LoggerLevel {
    for k, v := range levelMap {
        if strings.EqualFold(v, s) {
            return k
        }
    }
    return LevelInfo
}

func (lvl LoggerLevel) ZerlogLevel() zerolog.Level {
    if l, ok := zerologLevelMap[lvl]; ok {
        return l
    }
    return zerolog.InfoLevel
}

func (lvl LoggerLevel) PgxLogLevel() pgx.LogLevel {
    if l, ok := pgxlogLevelMap[lvl]; ok {
        return l
    }
    return pgx.LogLevelInfo
}
```

これを踏まえて環境変数の取り出し処理を以下のようにする。

```go:env/env.go
package env

import (
    "os"
    "strings"

    "github.com/jackc/pgx/v4"
    "github.com/joho/godotenv"
    "github.com/rs/zerolog"
    "github.com/spiegel-im-spiegel/gocli/config"
)

const (
    ServiceName = "elephantsql"
)

func init() {
    //load ${XDG_CONFIG_HOME}/${ServiceName}/env file
    if err := godotenv.Load(config.Path(ServiceName, "env")); err != nil {
        panic(err)
    }
}

func PostgresDSN() string {
    return os.Getenv("ELEPHANTSQL_URL")
}

func LogLevel() LoggerLevel {
    return getLogLevel(os.Getenv("LOGLEVEL"))
}

func ZerologLevel() zerolog.Level {
    return LogLevel().ZerlogLevel()
}

func PgxlogLevel() pgx.LogLevel {
    return LogLevel().PgxLogLevel()
}

func EnableLogFile() bool {
    return strings.EqualFold(os.Getenv("ENABLE_LOGFILE"), "true")
}
```

## loggr サブパッケージ

環境変数が読めるようになったので logger を生成するパッケージを書こう。パッケージ名が loggr となっているのは，世の中に logger という名前のパッケージがあまりに多いので回避のため（笑）

```go:loggr/loggr.go
package loggr

import (
    "fmt"
    "io"
    "os"
    "sample/env"
    "time"

    "github.com/rs/zerolog"
    "github.com/spiegel-im-spiegel/errs"
    "github.com/spiegel-im-spiegel/gocli/cache"
)

func New() *zerolog.Logger {
    logger := zerolog.Nop()
    if env.ZerologLevel() == zerolog.NoLevel {
        return &logger
    }
    if env.EnableLogFile() {
        // make path to ${XDG_CACHE_HOME}/${ServiceName}/access.YYYYMMDD.log file and create logger
        if file, err := os.OpenFile(cache.Path(env.ServiceName, fmt.Sprintf("access.%s.log", time.Now().Local().Format("20060102"))), os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0600); err == nil {
            logger = zerolog.New(io.MultiWriter(
                file,
                zerolog.ConsoleWriter{Out: os.Stdout, NoColor: false},
            )).Level(env.ZerologLevel()).With().Timestamp().Logger()
        } else {
            logger = zerolog.New(os.Stdout).Level(env.ZerologLevel()).With().Timestamp().Logger()
            logger.Error().Interface("error", errs.Wrap(err)).Msg("error in opening logfile")
        }
    } else {
        logger = zerolog.New(os.Stdout).Level(env.ZerologLevel()).With().Timestamp().Logger()
    }
    return &logger
}
```

前節までとはログファイルの出力先が異なるのでご注意。

## dbconn サブパッケージ

次は [github.com/jackc/pgx] パッケージを使って接続インスタンスを作るパッケージ。

```go:dbconn/dbconn.go
package dbconn

import (
    "database/sql"
    "sample/env"
    "sample/loggr"

    "github.com/jackc/pgx/v4"
    "github.com/jackc/pgx/v4/log/zerologadapter"
    "github.com/jackc/pgx/v4/stdlib"
    "github.com/rs/zerolog"
    "github.com/spiegel-im-spiegel/errs"
)

type PgxContext struct {
    Db     *sql.DB
    Logger *zerolog.Logger
}

func NewPgx() (*PgxContext, error) {
    dbctx := &PgxContext{
        Logger: loggr.New(),
    }
    cfg, err := pgx.ParseConfig(env.PostgresDSN())
    if err != nil {
        dbctx.Logger.Error().Interface("error", errs.Wrap(err)).Msg("error in pgx.ParseConfig() method")
        return nil, errs.Wrap(err, errs.WithContext("dsn", env.PostgresDSN()))
    }
    cfg.Logger = zerologadapter.NewLogger(*dbctx.Logger)
    cfg.LogLevel = env.PgxlogLevel()
    dbctx.Db = stdlib.OpenDB(*cfg)

    return dbctx, nil
}

func (dbctx *PgxContext) GetDb() *sql.DB {
    if dbctx == nil {
        return nil
    }
    return dbctx.Db
}

func (dbctx *PgxContext) GetLogger() *zerolog.Logger {
    if dbctx == nil {
        lggr := zerolog.Nop()
        return &lggr
    }
    return dbctx.Logger
}

func (dbctx *PgxContext) Acquire() (*pgx.Conn, error) {
    if db := dbctx.GetDb(); db != nil {
        conn, err := stdlib.AcquireConn(db)
        return conn, errs.Wrap(err)
    }
    return nil, errs.New("*sql.DB instance is nil.")
}

func (dbctx *PgxContext) Close() error {
    if db := dbctx.GetDb(); db != nil {
        return errs.Wrap(db.Close())
    }
    return nil
}
```

dbconn.PgxContext というコンテキストを作って，これを返している。これは logger を *[sql][database/sql].DB インスタンスの外側で使えるようにするため。また Acquire() メソッドを使って *[pgx][github.com/jackc/pgx].Conn インスタンスを取り出せるようにしてみた。

## orm サブパッケージ

最後に [GORM] のインスタンスを作るパッケージ。

```go:orm/gorm.go
package orm

import (
    "sample/dbconn"
    "sample/env"

    "github.com/rs/zerolog"
    "github.com/spiegel-im-spiegel/errs"
    "gorm.io/driver/postgres"
    "gorm.io/gorm"
    "gorm.io/gorm/logger"
)

type GormContext struct {
    Db     *gorm.DB
    Logger *zerolog.Logger
}

func NewGORM() (*GormContext, error) {
    pgxCtx, err := dbconn.NewPgx()
    if err != nil {
        return nil, errs.Wrap(err)
    }
    gormCtx := &GormContext{
        Logger: pgxCtx.GetLogger(),
    }
    loggr := logger.Discard
    if env.LogLevel() == env.LevelDebug {
        loggr = logger.Default
    }
    gormCtx.Db, err = gorm.Open(postgres.New(postgres.Config{
        Conn: pgxCtx.GetDb(),
    }), &gorm.Config{
        Logger: loggr,
    })
    if err != nil {
        pgxCtx.GetLogger().Error().Interface("error", errs.Wrap(err)).Msg("error in gorm.Open() method")
        pgxCtx.Close()
        return nil, errs.Wrap(err)
    }
    return gormCtx, nil
}

func (gormCtx *GormContext) GetDb() *gorm.DB {
    if gormCtx == nil {
        return nil
    }
    return gormCtx.Db
}

func (gormCtx *GormContext) GetLogger() *zerolog.Logger {
    if gormCtx == nil {
        lggr := zerolog.Nop()
        return &lggr
    }
    return gormCtx.Logger
}

func (gormCtx *GormContext) Close() error {
    if db := gormCtx.GetDb(); db != nil {
        if sqlDb, err := db.DB(); err == nil {
            sqlDb.Close()
        }
    }
    return nil
}
```

これも同じく logger を含めた orm.GormContext を返している。あと env ファイルで `LOGLEVEL` が `debug` と指定されている場合は [GORM] の logger も有効とするようにした。

## 起動サンプル

以上のサブパッケージ群を使ってサンプルコードを書き直すとこうなる。

```go:sampele2.go
//go:build run
// +build run

package main

import (
    "fmt"
    "os"
    "sample/orm"

    "github.com/spiegel-im-spiegel/errs"
    "github.com/spiegel-im-spiegel/gocli/exitcode"
)

func Run() exitcode.ExitCode {
    // create gorm.DB instance for PostgreSQL service
    gormCtx, err := orm.NewGORM()
    if err != nil {
        fmt.Fprintln(os.Stderr, err)
        return exitcode.Abnormal
    }
    defer gormCtx.Close()

    // query
    var results []map[string]interface{}
    tx := gormCtx.GetDb().Table("tablename").Find(&results) // "tablename" is not exist
    if tx.Error != nil {
        gormCtx.GetLogger().Error().Interface("error", errs.Wrap(tx.Error)).Send()
        return exitcode.Abnormal
    }

    return exitcode.Normal
}

func main() {
    Run().Exit()
}
```

うんうん。だいぶスッキリしたねぇ。これの実行結果は以下の通り。

```
$ go run sample2.go 
12:25PM INF Dialing PostgreSQL server host=hostname module=pgx
12:25PM INF Exec args=[] commandTag=null module=pgx pid=11556 sql=;
12:25PM ERR Query args=[] err="ERROR: relation \"tablename\" does not exist (SQLSTATE 42P01)" module=pgx pid=11556 sql="SELECT * FROM \"tablename\""
12:25PM ERR  error={"Context":{"function":"main.Run"},"Err":{"Msg":"ERROR: relation \"tablename\" does not exist (SQLSTATE 42P01)","Type":"*pgconn.PgError"},"Type":"*errs.Error"}
12:25PM INF closed connection module=pgx pid=11556
```

ん，問題なく SELECT 文でエラーになるね（笑）

まぁ，パッケージ間の関係が密すぎてテストが書きにくい（つか書けんな，これ）のはご容赦。もう少し弄りたい気持ちはあるが，今回はこの辺で勘弁してやろう（笑）


[Go]: https://go.dev/
[PostgreSQL]: https://www.postgresql.org/ "PostgreSQL: The world's most advanced open source database"
[os]: https://pkg.go.dev/os "os package - os - pkg.go.dev"
[database/sql]: https://pkg.go.dev/database/sql "sql package - database/sql - pkg.go.dev"
[GORM]: https://gorm.io/ "GORM - The fantastic ORM library for Golang, aims to be developer friendly."
[github.com/rs/zerolog]: https://github.com/rs/zerolog "rs/zerolog: Zero Allocation JSON Logger"
[github.com/jackc/pgx]: https://github.com/jackc/pgx "jackc/pgx: PostgreSQL driver and toolkit for Go"
