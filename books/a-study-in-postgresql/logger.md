---
title: "構造化 Logger を組み込む"
---

[github.com/jackc/pgx] には [log] 標準パッケージ以外の logger を組み込むことができる。どうせ組み込むなら構造化 logger を組み込むのがイマドキだよね。ということで，この節では [github.com/rs/zerolog] パッケージを紹介して，これを [github.com/jackc/pgx] パッケージに組み込むところまで行う。

## 構造化ログ zerolog

[github.com/rs/zerolog] パッケージは JSON 形式の構造化ログを出力する。説明するより実際にコードを書いたほうが早いだろう。

```go:sample0.go
//go:build run
// +build run

package main

import (
    "os"

    "github.com/rs/zerolog"
    "github.com/spiegel-im-spiegel/errs"
    "github.com/spiegel-im-spiegel/gocli/exitcode"
)

func Run() exitcode.ExitCode {
    logger := zerolog.New(
        os.Stdout,
    ).Level(zerolog.DebugLevel).With().Timestamp().Logger()
    logger.Err(os.ErrInvalid).Send()
    logger.Error().Interface("error", errs.Wrap(os.ErrInvalid)).Msg("sample error")
    return exitcode.Normal
}

func main() {
    Run().Exit()
}
```

これを実行すると標準出力に

```
$ go run sample0.go
{"level":"error","error":"invalid argument","time":"2021-09-19T20:11:32+09:00"}
{"level":"error","error":{"Type":"*errs.Error","Err":{"Type":"*errors.errorString","Msg":"invalid argument"},"Context":{"function":"main.Run"}},"time":"2021-09-19T20:11:32+09:00","message":"sample error"}
```

などと出力される。ちょっと分かりにくいから [jq] で整形してみるか。

```
$ go run sample0.go | jq .
{
  "level": "error",
  "error": "invalid argument",
  "time": "2021-09-19T20:13:13+09:00"
}
{
  "level": "error",
  "error": {
    "Type": "*errs.Error",
    "Err": {
      "Type": "*errors.errorString",
      "Msg": "invalid argument"
    },
    "Context": {
      "function": "main.Run"
    }
  },
  "time": "2021-09-19T20:13:13+09:00",
  "message": "sample error"
}
```

どやさ！ ちなみに [github.com/spiegel-im-spiegel/errs] は拙作のエラーハンドリング・パッケージで， MarshalJSON() メソッドを持っているため，エラー詳細を JSON 形式で出力することができる。

[github.com/rs/zerolog] は更に標準出力専用のアダプタも持っているので [io].MultiWriter() 関数を使って

```go:sample0b.go
func Run() exitcode.ExitCode {
    file, err := os.Create("access.log")
    if err != nil {
        fmt.Println(err)
        return exitcode.Abnormal
    }
    logger := zerolog.New(
        io.MultiWriter(
            file,
            zerolog.ConsoleWriter{Out: os.Stdout, NoColor: false}),
    ).Level(zerolog.DebugLevel).With().Timestamp().Logger()
    logger.Err(os.ErrInvalid).Send()
    logger.Error().Interface("error", errs.Wrap(os.ErrInvalid)).Msg("sample error")
    return exitcode.Normal
}
```

などとすれば，構造化ログはファイルに，標準出力には

```
$ go run sample0b.go
8:21PM ERR  error="invalid argument"
8:21PM ERR sample error error={"Context":{"function":"main.Run"},"Err":{"Msg":"invalid argument","Type":"*errors.errorString"},"Type":"*errs.Error"}
```

といった内容を色付きで出力してくれる。素晴らしい。

Logger 作成処理も関数として切り出してしまおう。

```go
import (
    "fmt"
    "io"
    "os"
    "time"

    "github.com/rs/zerolog"
    "github.com/spiegel-im-spiegel/errs"
    "github.com/spiegel-im-spiegel/gocli/cache"
    "github.com/spiegel-im-spiegel/gocli/exitcode"
)

func CreateLogger() zerolog.Logger {
    logger := zerolog.Nop()
    logpath := cache.Path("elephantsql", fmt.Sprintf("access.%s.log", time.Now().Local().Format("20060102"))) // logpath = ~/.cache/elephantsql/access.YYYYMMDD.log
    file, err := os.OpenFile(logpath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0600)
    if err != nil {
        logger = zerolog.New(os.Stdout)
    } else {
        logger = zerolog.New(io.MultiWriter(
            file,
            zerolog.ConsoleWriter{Out: os.Stdout, NoColor: false},
        ))
    }
    logger = logger.Level(zerolog.DebugLevel).With().Timestamp().Logger()
    if err != nil {
        logger.Error().Interface("error", errs.Wrap(err, errs.WithContext("logpath", logpath))).Str("logpath", logpath).Msg("error in opening logfile")
    }
    return logger
}
```

この例では，あらかじめ ~/.cache/elephantsql/ ディレクトリを作っておいて，その中の access.YYYYMMDD.log ファイルにログを追記するようにしている。 cache.Path() は拙作の [github.com/spiegel-im-spiegel/gocli]/cache パッケージの関数で，内部で [os].UserCacheDir() を使ってキャッシュ用ディレクトリを取得している。

```go:os/file.go
// UserCacheDir returns the default root directory to use for user-specific
// cached data. Users should create their own application-specific subdirectory
// within this one and use that.
//
// On Unix systems, it returns $XDG_CACHE_HOME as specified by
// https://specifications.freedesktop.org/basedir-spec/basedir-spec-latest.html if
// non-empty, else $HOME/.cache.
// On Darwin, it returns $HOME/Library/Caches.
// On Windows, it returns %LocalAppData%.
// On Plan 9, it returns $home/lib/cache.
//
// If the location cannot be determined (for example, $HOME is not defined),
// then it will return an error.
func UserCacheDir() (string, error) {
...
```


ま，まぁ，最近はログを直接ファイルに吐き出すのは流行りじゃないみたいだし，ちょっとしたツール用ならこれで十分かな。

## zerolog を pgx に組み込む

[github.com/jackc/pgx] パッケージは [github.com/rs/zerolog] 用のアダプタを持っていて logger を [github.com/rs/zerolog] に換装することができる。

```go:sample1.go
import (
    "context"
    "fmt"
    "io"
    "os"
    "time"

    "github.com/jackc/pgx/v4"
    "github.com/jackc/pgx/v4/log/zerologadapter"
    "github.com/jackc/pgx/v4/pgxpool"
    "github.com/joho/godotenv"
    "github.com/rs/zerolog"
    "github.com/spiegel-im-spiegel/errs"
    "github.com/spiegel-im-spiegel/gocli/cache"
    "github.com/spiegel-im-spiegel/gocli/config"
    "github.com/spiegel-im-spiegel/gocli/exitcode"
)

func Run() exitcode.ExitCode {
    // get logger
    logger := CreateLogger()

    // create connection pool for PostgreSQL service
    cfg, err := pgxpool.ParseConfig(os.Getenv("ELEPHANTSQL_URL"))
    if err != nil {
        logger.Error().Interface("error", errs.Wrap(err)).Send()
        return exitcode.Abnormal
    }
    cfg.ConnConfig.Logger = zerologadapter.NewLogger(logger)
    cfg.ConnConfig.LogLevel = pgx.LogLevelDebug
    pool, err := pgxpool.ConnectConfig(context.TODO(), cfg)
    if err != nil {
        logger.Error().Interface("error", errs.Wrap(err)).Send()
        return exitcode.Abnormal
    }
    defer pool.Close()

    // acquire connection
    conn, err := pool.Acquire(context.TODO())
    if err != nil {
        logger.Error().Interface("error", errs.Wrap(err)).Send()
        return exitcode.Abnormal
    }
    defer conn.Release()

    // query
    _, err = conn.Query(context.TODO(), "SELECT * FROM tablename") // "tablename" is not exist
    if err != nil {
        logger.Error().Interface("error", errs.Wrap(err)).Send()
        return exitcode.Abnormal
    }

    return exitcode.Normal
}
```

実はデータベースの中はまだ空っぽなので，これを実行すると

```
$ go run sample1.go
9:39PM INF Dialing PostgreSQL server host=hostname module=pgx
9:39PM ERR Query args=[] err="ERROR: relation \"tablename\" does not exist (SQLSTATE 42P01)" module=pgx pid=27036 sql="SELECT * FROM tablename"
9:39PM ERR  error={"Context":{"function":"main.Run"},"Err":{"Msg":"ERROR: relation \"tablename\" does not exist (SQLSTATE 42P01)","Type":"*pgconn.PgError"},"Type":"*errs.Error"}
9:39PM INF closed connection module=pgx pid=27036
```

という感じにエラーになる。まぁ，ログ出力の確認は出来たということで（笑）

[database/sql] を使いたい場合は stdlib.OpenDB() 関数から [sql][database/sql].DB インスタンスを取得する。

```go:sample2.go
import (
    "context"
    "fmt"
    "io"
    "os"
    "time"

    "github.com/jackc/pgx/v4"
    "github.com/jackc/pgx/v4/log/zerologadapter"
    "github.com/jackc/pgx/v4/stdlib"
    "github.com/joho/godotenv"
    "github.com/rs/zerolog"
    "github.com/spiegel-im-spiegel/errs"
    "github.com/spiegel-im-spiegel/gocli/cache"
    "github.com/spiegel-im-spiegel/gocli/config"
    "github.com/spiegel-im-spiegel/gocli/exitcode"
)

func Run() exitcode.ExitCode {
    // get logger
    logger := CreateLogger()

    // create sql.DB instance for PostgreSQL service
    cfg, err := pgx.ParseConfig(os.Getenv("ELEPHANTSQL_URL"))
    if err != nil {
        logger.Error().Interface("error", errs.Wrap(err)).Send()
        return exitcode.Abnormal
    }
    cfg.Logger = zerologadapter.NewLogger(logger)
    cfg.LogLevel = pgx.LogLevelDebug
    db := stdlib.OpenDB(*cfg)
    defer db.Close()

    // get connection from connection pool
    conn, err := db.Conn(context.TODO())
    if err != nil {
        logger.Error().Interface("error", errs.Wrap(err)).Send()
        return exitcode.Abnormal
    }
    defer conn.Close()

    // query
    _, err = conn.QueryContext(context.TODO(), "SELECT * FROM tablename") // "tablename" is not exist
    if err != nil {
        logger.Error().Interface("error", errs.Wrap(err)).Send()
        return exitcode.Abnormal
    }

    return exitcode.Normal
}
```

なお [github.com/jackc/pgx] 以外のパッケージで [github.com/rs/zerolog] を直接組み込めない場合は [github.com/simukti/sqldb-logger] パッケージ経由で組み込めるようだ。

https://pod.hatenablog.com/entry/2020/09/30/073034

[Go]: https://go.dev/
[PostgreSQL]: https://www.postgresql.org/ "PostgreSQL: The world's most advanced open source database"
[jq]: https://stedolan.github.io/jq/
[database/sql]: https://pkg.go.dev/database/sql "sql package - database/sql - pkg.go.dev"
[log]: https://pkg.go.dev/log "log package - log - pkg.go.dev"
[io]: https://pkg.go.dev/io "io package - io - pkg.go.dev"
[os]: https://pkg.go.dev/os "os package - os - pkg.go.dev"
[github.com/jackc/pgx]: https://github.com/jackc/pgx "jackc/pgx: PostgreSQL driver and toolkit for Go"
[github.com/rs/zerolog]: https://github.com/rs/zerolog "rs/zerolog: Zero Allocation JSON Logger"
[github.com/spiegel-im-spiegel/errs]: https://github.com/spiegel-im-spiegel/errs "spiegel-im-spiegel/errs: Error handling for Golang"
[github.com/spiegel-im-spiegel/gocli]: https://github.com/spiegel-im-spiegel/gocli "spiegel-im-spiegel/gocli: Minimal Packages for Command-Line Interface"
[github.com/simukti/sqldb-logger]: https://github.com/simukti/sqldb-logger "simukti/sqldb-logger: A logger for Go SQL database driver without modifying existing *sql.DB stdlib usage."
