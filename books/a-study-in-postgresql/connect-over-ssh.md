---
title: "【付録1】 SSH 越しに PostgreSQL に接続する"
---

個人的に訳あって SSH 接続越しに RDBMS に接続する要件があり，その機能を実装したパッケージを公開した。

https://github.com/goark/sshql
https://text.baldanders.info/release/sshql/

とりあえず [github.com/jackc/pgx] と [github.com/goark/sshql] を使って SSH 越しに [PostgreSQL] サービスにアクセスするコードを書いてみる。

```diff go
package main

import (
    "context"
    "fmt"
    "os"

+   "github.com/goark/sshql"
    "github.com/jackc/pgx/v4/pgxpool"
)

func main() {
+   dialer := &sshql.Dialer{
+       Hostname:   "sshserver",
+       Port:       22,
+       Username:   "remoteuser",
+       Password:   "passphraseforauthkey",
+       PrivateKey: "/home/username/.ssh/id_eddsa",
+   }

    cfg, err := pgxpool.ParseConfig("postgres://dbuser:dbpassword@localhost:5432/example?sslmode=disable")
    if err != nil {
        fmt.Fprintln(os.Stderr, err)
        return
    }
+   cfg.ConnConfig.DialFunc = dialer.DialContext
+
+   if err := dialer.Connect(); err != nil {
+       fmt.Fprintln(os.Stderr, err)
+       return
+   }
+   defer dialer.Close()

    pool, err := pgxpool.ConnectConfig(context.TODO(), cfg)
    if err != nil {
        fmt.Fprintln(os.Stderr, err)
        return
    }
    defer pool.Close()

    rows, err := pool.Query(context.TODO(), "SELECT * FROM tablename")
    if err != nil {
        fmt.Fprintln(os.Stderr, err)
        return
    }
    rows.Close()

    for rows.Next() {
        var id int64
        var name string
        if err := rows.Scan(&id, &name); err != nil {
            fmt.Fprintln(os.Stderr, err)
            break
        }
        fmt.Printf("ID: %d  Name: %s\n", id, name)
    }
}
```

色分けされているところが拙作の [github.com/goark/sshql] パッケージに関連する部分。 [sshql][github.com/goark/sshql].Dialer 構造体の使い方については[解説記事](https://text.baldanders.info/release/sshql/ "sshql — SSH 越しに RDBMS にアクセスする")を読んでいただきたい。処理のポイントは [sshql][github.com/goark/sshql].Dialer 構造体の Connect() および Close() メソッドを使って SSH の接続・切断を明示的に行うことと [pgx][github.com/jackc/pgx]/pgxpool.Config 構造体のメンバである ConnConfig.DialFunc に [sshql][github.com/goark/sshql].Dialer 構造体の DialContext() メソッドを登録することである。

[4節]で [pgx][github.com/jackc/pgx]/pgxpool.Config 構造体に [zerolog][github.com/rs/zerolog] を組み込む際に

```go
cfg, err := pgxpool.ParseConfig(os.Getenv("ELEPHANTSQL_URL"))
if err != nil {
    logger.Error().Interface("error", errs.Wrap(err)).Send()
    return exitcode.Abnormal
}
cfg.ConnConfig.Logger = zerologadapter.NewLogger(logger)
cfg.ConnConfig.LogLevel = pgx.LogLevelDebug
```

と書いたが，同じように [sshql][github.com/goark/sshql].Dialer のメソッドも組み込めるわけだ。

[4節]: https://zenn.dev/spiegel/books/a-study-in-postgresql/viewer/logger "構造化 Logger を組み込む"
[Go]: https://go.dev/
[PostgreSQL]: https://www.postgresql.org/ "PostgreSQL: The world's most advanced open source database"
[github.com/jackc/pgx]: https://github.com/jackc/pgx "jackc/pgx: PostgreSQL driver and toolkit for Go"
[github.com/rs/zerolog]: https://github.com/rs/zerolog "rs/zerolog: Zero Allocation JSON Logger"
[github.com/goark/sshql]: https://github.com/goark/sshql "goark/sshql: Go SQL drivers over SSH"
