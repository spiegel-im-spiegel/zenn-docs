---
title: "GORM で PostgreSQL に接続する"
---

[GORM] は [Go] 向けの ORM (Object Relational Mapper) ハンドリング・パッケージとしては定番のひとつになっている。

https://github.com/go-gorm/gorm

ORM としての基本機能はきちんと押さえられている。その上でテーブル設計や SQL 文構築の自由度が高いのが特徴と言えるだろう。その代わり型の解釈で挙動が怪しくなる場合があるため，取り回しでは若干の注意を要する（まぁ interface{} 型を使って無理やり汎化してるので仕方ない面はあるのだが）。

SQL 文の方言があるためどんな DB ドライバでも受け付けるというわけではないのだが MySQL, [PostgreSQL], SQLite, SQL Server といったメジャーどころは対応している。内部では *[sql][database/sql].DB 型のインスタンスとして保持しているようなので，上述の製品と互換性のある RDBMS であれば対応できる可能性がある。

今回であれば stdlib.OpenDB() 関数で生成した [sql][database/sql].DB インスタンスを渡してしまえばよい。こんな感じ。

```go
import (
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
    "gorm.io/driver/postgres"
    "gorm.io/gorm"
    "gorm.io/gorm/logger"
)

func Run() exitcode.ExitCode {
    // get logger
    zlogger := CreateLogger()

    // create gorm.DB instance for PostgreSQL service
    cfg, err := pgx.ParseConfig(os.Getenv("ELEPHANTSQL_URL"))
    if err != nil {
        zlogger.Error().Interface("error", errs.Wrap(err)).Send()
        return exitcode.Abnormal
    }
    cfg.Logger = zerologadapter.NewLogger(zlogger)
    cfg.LogLevel = pgx.LogLevelDebug
    db, err := gorm.Open(postgres.New(postgres.Config{
        Conn: stdlib.OpenDB(*cfg),
    }), &gorm.Config{
        Logger: logger.Discard,
    })
    defer func() {
        if sqlDb, err := db.DB(); err == nil {
            sqlDb.Close()
        }
    }()

    // query
    var results []map[string]interface{}
    tx := db.Table("tablename").Find(&results) // "tablename" is not exist
    if tx.Error != nil {
        zlogger.Error().Interface("error", errs.Wrap(tx.Error)).Send()
        return exitcode.Abnormal
    }

    return exitcode.Normal
}
```

この中の

```go
db, err := gorm.Open(postgres.New(postgres.Config{
    Conn: stdlib.OpenDB(*cfg),
}), &gorm.Config{
    Logger: logger.Discard,
})
```

が核心部分だね。なお，このコードでは [GORM] の logger は潰してある（nil をセットしても潰せないので注意）。

[GORM] は専用の logger を持っているのだがショージキ微妙。これなら [github.com/simukti/sqldb-logger] パッケージみたいなのを使ってサードパーティの logger を受け付けるようにして欲しかった。まぁ，でも，デバッグ中に [GORM] の内部状態を知りたいこともあるだろうから，その場合は logger.Default などをセットすればいいだろう。

[Go]: https://go.dev/
[PostgreSQL]: https://www.postgresql.org/ "PostgreSQL: The world's most advanced open source database"
[database/sql]: https://pkg.go.dev/database/sql "sql package - database/sql - pkg.go.dev"
[GORM]: https://gorm.io/ "GORM - The fantastic ORM library for Golang, aims to be developer friendly."
[github.com/simukti/sqldb-logger]: https://github.com/simukti/sqldb-logger "simukti/sqldb-logger: A logger for Go SQL database driver without modifying existing *sql.DB stdlib usage."
