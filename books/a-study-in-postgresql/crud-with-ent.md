---
title: "ent で CRUD"
---

今回の調査もいよいよ大詰め。もうしばらくお付き合いください。

## Create

[ent] って CRUD の dry run ができないのか？ ドキドキするんだが...

CreateBulk() というメソッドを使えば複数レコードを一気に作成できるとあったので，以下のようなコードを組んだわけさ。

```go:sample4.go
import (
    "context"
    "fmt"
    "os"
    "sample/dbconn"
    "sample/ent"
    "sample/files"

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
    client := entCtx.GetClient()

    file1 := "files/file1.txt"
    bin1, err := files.GetBinary(file1)
    if err != nil {
        entCtx.GetLogger().Error().Interface("error", errs.Wrap(err)).Send()
        return exitcode.Abnormal
    }
    file2 := "files/file2.txt"
    bin2, err := files.GetBinary(file2)
    if err != nil {
        entCtx.GetLogger().Error().Interface("error", errs.Wrap(err)).Send()
        return exitcode.Abnormal
    }
    file3 := "files/file3.txt"
    bin3, err := files.GetBinary(file3)
    if err != nil {
        entCtx.GetLogger().Error().Interface("error", errs.Wrap(err)).Send()
        return exitcode.Abnormal
    }

    // create data
    if _, err := client.User.CreateBulk(
        client.User.Create().SetUsername("Alice").AddOwned(
            &ent.BinaryFile{Filename: file1, Body: &bin1},
            &ent.BinaryFile{Filename: file2, Body: &bin2},
        ),
        client.User.Create().SetUsername("Bob").AddOwned(
            &ent.BinaryFile{Filename: file3, Body: &bin3},
        ),
    ).Save(context.TODO()); err != nil {
        entCtx.GetLogger().Error().Interface("error", errs.Wrap(err)).Send()
        return exitcode.Abnormal
    }

    return exitcode.Normal
}
```

これを実行してみたら

```
$ go run sample4.go 
0:00AM INF Dialing PostgreSQL server host=hostname module=pgx
0:00AM INF Exec args=[] commandTag=QkVHSU4= module=pgx pid=4263 sql=begin
0:00AM DBG driver.Tx(6dba389a-f7f2-4ebb-b7cc-b37e0dbc1840): started
0:00AM DBG Tx(6dba389a-f7f2-4ebb-b7cc-b37e0dbc1840).Query: query=INSERT INTO "users" ("created_at", "updated_at", "username") VALUES ($1, $2, $3), ($4, $5, $6) RETURNING "id" args=[2021-09-20 00:00:00.768732103 +0900 JST m=+0.000841017 2021-09-20 00:00:00.768732203 +0900 JST m=+0.000841117 Alice 2021-09-20 00:00:00.768732463 +0900 JST m=+0.000841377 2021-09-20 00:00:00.768732533 +0900 JST m=+0.000841448 Bob]
0:00AM INF Query args=["2021-09-20 00:00:00.768732103+09:00","2021-09-20 00:00:00.768732203+09:00","Alice","2021-09-20 00:00:00.768732463+09:00","2021-09-20 00:00:00.768732533+09:00","Bob"] module=pgx pid=4263 rowCount=2 sql="INSERT INTO \"users\" (\"created_at\", \"updated_at\", \"username\") VALUES ($1, $2, $3), ($4, $5, $6) RETURNING \"id\""
0:00AM DBG Tx(6dba389a-f7f2-4ebb-b7cc-b37e0dbc1840).Exec: query=UPDATE "binary_files" SET "user_owned" = $1 WHERE "id" = $2 AND "user_owned" IS NULL args=[1 0]
0:00AM INF Exec args=[1,0] commandTag=VVBEQVRFIDA= module=pgx pid=4263 sql="UPDATE \"binary_files\" SET \"user_owned\" = $1 WHERE \"id\" = $2 AND \"user_owned\" IS NULL"
0:00AM DBG Tx(6dba389a-f7f2-4ebb-b7cc-b37e0dbc1840): rollbacked
0:00AM INF Exec args=[] commandTag=Uk9MTEJBQ0s= module=pgx pid=4263 sql=rollback
0:00AM ERR  error={"Context":{"function":"main.Run"},"Err":{"Cause":{"Msg":"one of [0] is already connected to a different user_owned","Type":"*sqlgraph.ConstraintError"},"Msg":"ent: constraint failed: one of [0] is already connected to a different user_owned","Type":"*ent.ConstraintError"},"Type":"*errs.Error"}
0:00AM INF closed connection module=pgx pid=4263
```

エラーになっちゃったよ。ログを眺めてもよく分からないが，ヘンテコな SQL 文を投げてるように見える。どうも one-to-many で入れ子になっているデータ構造をそのまま突っ込んでもダメらしい。しょうがない，地道にやるか。

```go:sample4-1.go
import (
    "context"
    "fmt"
    "os"
    "sample/dbconn"
    "sample/files"

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
    client := entCtx.GetClient()

    file1 := "files/file1.txt"
    bin1, err := files.GetBinary(file1)
    if err != nil {
        entCtx.GetLogger().Error().Interface("error", errs.Wrap(err)).Send()
        return exitcode.Abnormal
    }
    file2 := "files/file2.txt"
    bin2, err := files.GetBinary(file2)
    if err != nil {
        entCtx.GetLogger().Error().Interface("error", errs.Wrap(err)).Send()
        return exitcode.Abnormal
    }

    // create data
    user, err := client.User.Create().SetUsername("Alice").Save(context.TODO())
    if err != nil {
        entCtx.GetLogger().Error().Interface("error", errs.Wrap(err)).Send()
        return exitcode.Abnormal
    }
    if _, err := client.BinaryFile.CreateBulk(
        client.BinaryFile.Create().SetFilename(file1).SetBody(bin1).SetOwner(user),
        client.BinaryFile.Create().SetFilename(file2).SetBody(bin2).SetOwner(user),
    ).Save(context.TODO()); err != nil {
        entCtx.GetLogger().Error().Interface("error", errs.Wrap(err)).Send()
        return exitcode.Abnormal
    }

    return exitcode.Normal
}
```

これを実行すると

```
$ go run sample4-1.go 
0:00AM INF Dialing PostgreSQL server host=hostname module=pgx
0:00AM INF Exec args=[] commandTag=QkVHSU4= module=pgx pid=24148 sql=begin
0:00AM DBG driver.Tx(e7f78224-0b97-4c79-b3d0-74d8421cdfd8): started
0:00AM DBG Tx(e7f78224-0b97-4c79-b3d0-74d8421cdfd8).Query: query=INSERT INTO "users" ("username", "created_at", "updated_at") VALUES ($1, $2, $3) RETURNING "id" args=[Alice 2021-09-20 00:00:00.138394439 +0900 JST m=+0.001185546 2021-09-20 00:00:00.13839455 +0900 JST m=+0.001185656]
0:00AM INF Query args=["Alice","2021-09-20 00:00:00.138394439+09:00","2021-09-20 00:00:00.13839455+09:00"] module=pgx pid=24148 rowCount=1 sql="INSERT INTO \"users\" (\"username\", \"created_at\", \"updated_at\") VALUES ($1, $2, $3) RETURNING \"id\""
0:00AM DBG Tx(e7f78224-0b97-4c79-b3d0-74d8421cdfd8): committed
0:00AM INF Exec args=[] commandTag=Q09NTUlU module=pgx pid=24148 sql=commit
0:00AM INF Exec args=[] commandTag=QkVHSU4= module=pgx pid=24148 sql=begin
0:00AM DBG driver.Tx(afdf601a-6f84-4b0b-88b7-1b6dec244063): started
0:00AM DBG Tx(afdf601a-6f84-4b0b-88b7-1b6dec244063).Query: query=INSERT INTO "binary_files" ("body", "created_at", "filename", "updated_at", "user_owned") VALUES ($1, $2, $3, $4, $5), ($6, $7, $8, $9, $10) RETURNING "id" args=[[72 101 108 108 111 33 32 73 32 97 109 32 102 105 108 101 32 110 117 109 98 101 114 32 49 46 10] 2021-09-20 00:00:00.981711959 +0900 JST m=+0.844503075 files/file1.txt 2021-09-20 00:00:00.981712269 +0900 JST m=+0.844503385 4 [72 101 108 108 111 33 32 73 32 97 109 32 102 105 108 101 32 110 117 109 98 101 114 32 50 46 10] 2021-09-20 00:00:00.98171266 +0900 JST m=+0.844503776 files/file2.txt 2021-09-20 00:00:00.98171277 +0900 JST m=+0.844503876 4]
0:00AM INF Query args=["48656c6c6f21204920616d2066696c65206e756d62657220312e0a","2021-09-20 00:00:00.981711959+09:00","files/file1.txt","2021-09-20 00:00:00.981712269+09:00",4,"48656c6c6f21204920616d2066696c65206e756d62657220322e0a","2021-09-20 00:00:00.98171266+09:00","files/file2.txt","2021-09-20 00:00:00.98171277+09:00",4] module=pgx pid=24148 rowCount=2 sql="INSERT INTO \"binary_files\" (\"body\", \"created_at\", \"filename\", \"updated_at\", \"user_owned\") VALUES ($1, $2, $3, $4, $5), ($6, $7, $8, $9, $10) RETURNING \"id\""
0:00AM DBG Tx(afdf601a-6f84-4b0b-88b7-1b6dec244063): committed
0:00AM INF Exec args=[] commandTag=Q09NTUlU module=pgx pid=24148 sql=commit
0:00AM INF closed connection module=pgx pid=24148
```

今度は正常終了した。 CreateBulk() もうまく動いてるみたいだし [ElephantSQL] の SQL Browser で見てもちゃんと値が格納されているようだ。でも Create() や CreateBulk() が走るたびに commit されるのは面白くないよねぇ。

### トランザクション処理

[ent] では User とか BinaryFile とかのノード単位で CRUD 処理を行うようで，トランザクションもその処理単位で閉じている。もちろん Commit() や Rollback() といったメソッドは用意されているが， [GORM] の CRUD 紹介の節でも述べた通り，トランザクション処理は成功時と失敗時で後始末が異なるのですこぶる鬱陶しい。そこで [GORM] の Transaction() メソッドのようなものがあると便利である。 [ent] にはそういうものは見当たらないが，ないなら書けばいいのである（笑）

幸い[サンプルコード](https://entgo.io/ja/docs/transactions/#%E3%83%99%E3%82%B9%E3%83%88%E3%83%97%E3%83%A9%E3%82%AF%E3%83%86%E3%82%A3%E3%82%B9 "トランザクション | ent")があるので，これを利用させてもらおう。こんな感じでどうだろう。

```go:dbconn/entclient.go
package dbconn

import (
    "context"
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

func (entCtx *EntContext) Transaction(ctx context.Context, fn func(tx *ent.Tx) error) error {
    client := entCtx.GetClient()
    if client == nil {
        return errs.New("null reference instance")
    }
    logger := entCtx.GetLogger()

    logger.Info().Msg("begining transaction")
    tx, err := client.Tx(ctx)
    if err != nil {
        return errs.Wrap(err)
    }
    defer func() {
        if v := recover(); v != nil {
            _ = tx.Rollback()
            panic(v)
        }
    }()

    if err := fn(tx); err != nil {
        txErr := errs.Wrap(err)
        if err := tx.Rollback(); err != nil {
            return errs.Wrap(err, errs.WithCause(txErr))
        }
        return txErr
    }

    logger.Info().Msg("committing transaction")
    if err := tx.Commit(); err != nil {
        return errs.Wrap(err)
    }
    return nil
}
```

処理の中で panic を拾っているのは SaveX() とかフザけたメソッドに対応するためではなく（メモリ不足など）本当に不測の事態でも最低限 rollback を走らせるようにするため（rollback がまともに動くとは限らないので返り値は無視する）。 Defer が使えないというのは事程左様に面倒くさい。

これを踏まえて，先程の sample4-1.go を以下のように書き直す。

```go:sample4-2.go
import (
    "context"
    "fmt"
    "os"
    "sample/dbconn"
    "sample/ent"
    "sample/files"

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

    file1 := "files/file1.txt"
    bin1, err := files.GetBinary(file1)
    if err != nil {
        entCtx.GetLogger().Error().Interface("error", errs.Wrap(err)).Send()
        return exitcode.Abnormal
    }
    file2 := "files/file2.txt"
    bin2, err := files.GetBinary(file2)
    if err != nil {
        entCtx.GetLogger().Error().Interface("error", errs.Wrap(err)).Send()
        return exitcode.Abnormal
    }

    // create data
    if err := entCtx.Transaction(context.TODO(), func(tx *ent.Tx) error {
        user, err := tx.User.Create().SetUsername("Alice").Save(context.TODO())
        if err != nil {
            return errs.Wrap(err, errs.WithContext("username", "Alice"))
        }
        if _, err := tx.BinaryFile.CreateBulk(
            tx.BinaryFile.Create().SetFilename(file1).SetBody(bin1).SetOwner(user),
            tx.BinaryFile.Create().SetFilename(file2).SetBody(bin2).SetOwner(user),
        ).Save(context.TODO()); err != nil {
            return errs.Wrap(err, errs.WithContext("files", []string{file1, file2}))
        }
        return nil
    }); err != nil {
        entCtx.GetLogger().Error().Interface("error", errs.Wrap(err)).Send()
        return exitcode.Abnormal
    }

    return exitcode.Normal
}
```

テーブルを（手動で）真っさらにしてからコレを動かすと

```
$ go run sample4-2.go 
0:00AM INF begining transaction
0:00AM INF Dialing PostgreSQL server host=hostname module=pgx
0:00AM INF Exec args=[] commandTag=QkVHSU4= module=pgx pid=12126 sql=begin
0:00AM DBG driver.Tx(1208c8df-310e-4121-b8fe-37ca2a4f3d19): started
0:00AM DBG Tx(1208c8df-310e-4121-b8fe-37ca2a4f3d19).Query: query=INSERT INTO "users" ("username", "created_at", "updated_at") VALUES ($1, $2, $3) RETURNING "id" args=[Alice 2021-09-20 00:00:00.929532828 +0900 JST m=+0.513447613 2021-09-20 00:00:00.929533219 +0900 JST m=+0.513448003]
0:00AM INF Query args=["Alice","2021-09-20 00:00:00.929532828+09:00","2021-09-20 00:00:00.929533219+09:00"] module=pgx pid=12126 rowCount=1 sql="INSERT INTO \"users\" (\"username\", \"created_at\", \"updated_at\") VALUES ($1, $2, $3) RETURNING \"id\""
0:00AM DBG Tx(1208c8df-310e-4121-b8fe-37ca2a4f3d19).Query: query=INSERT INTO "binary_files" ("body", "created_at", "filename", "updated_at", "user_owned") VALUES ($1, $2, $3, $4, $5), ($6, $7, $8, $9, $10) RETURNING "id" args=[[72 101 108 108 111 33 32 73 32 97 109 32 102 105 108 101 32 110 117 109 98 101 114 32 49 46 10] 2021-09-20 00:00:00.04947576 +0900 JST m=+0.633390554 files/file1.txt 2021-09-20 00:00:00.049476091 +0900 JST m=+0.633390875 1 [72 101 108 108 111 33 32 73 32 97 109 32 102 105 108 101 32 110 117 109 98 101 114 32 50 46 10] 2021-09-20 00:00:00.049476832 +0900 JST m=+0.633391616 files/file2.txt 2021-09-20 00:00:00.049477082 +0900 JST m=+0.633391867 1]
0:00AM INF Query args=["48656c6c6f21204920616d2066696c65206e756d62657220312e0a","2021-09-20 00:00:00.04947576+09:00","files/file1.txt","2021-09-20 00:00:00.049476091+09:00",1,"48656c6c6f21204920616d2066696c65206e756d62657220322e0a","2021-09-20 00:00:00.049476832+09:00","files/file2.txt","2021-09-20 00:00:00.049477082+09:00",1] module=pgx pid=12126 rowCount=2 sql="INSERT INTO \"binary_files\" (\"body\", \"created_at\", \"filename\", \"updated_at\", \"user_owned\") VALUES ($1, $2, $3, $4, $5), ($6, $7, $8, $9, $10) RETURNING \"id\""
0:00AM INF committing transaction
0:00AM DBG Tx(1208c8df-310e-4121-b8fe-37ca2a4f3d19): committed
0:00AM INF Exec args=[] commandTag=Q09NTUlU module=pgx pid=12126 sql=commit
0:00AM INF closed connection module=pgx pid=12126
```

よーし，うむうむ，よーし。

## Read (Query)

サクッと全件検索するコードを書いてみた。

```go:sample5.go
import (
    "context"
    "fmt"
    "os"
    "sample/dbconn"
    "sample/files"

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

    // query all data
    users, err := entCtx.GetClient().User.Query().WithOwned().All(context.TODO())
    if err != nil {
        entCtx.GetLogger().Error().Interface("error", errs.Wrap(err)).Send()
        return exitcode.Abnormal
    }

    if err := files.Output(os.Stdout, users); err != nil {
        entCtx.GetLogger().Error().Interface("error", errs.Wrap(err)).Send()
        return exitcode.Abnormal
    }

    return exitcode.Normal
}
```

WithOwned() オプションを付けることで `owned` エッヂに紐づくノード（BinaryFile）のデータも一緒に取得できる。これを実行すると

```json
$ go run sample5.go | jq .
[
  {
    "id": 1,
    "username": "Alice",
    "created_at": "2021-09-20 00:00:00.034896+09:00",
    "updated_at": "2021-09-20 00:00:00.034897+09:00",
    "edges": {
      "owned": [
        {
          "id": 1,
          "filename": "files/file1.txt",
          "body": "SGVsbG8hIEkgYW0gZmlsZSBudW1iZXIgMS4K",
          "created_at": "2021-09-20 00:00:00.157279+09:00",
          "updated_at": "2021-09-20 00:00:00.15728+09:00",
          "edges": {}
        },
        {
          "id": 2,
          "filename": "files/file2.txt",
          "body": "SGVsbG8hIEkgYW0gZmlsZSBudW1iZXIgMi4K",
          "created_at": "2021-09-20 00:00:00.157281+09:00",
          "updated_at": "2021-09-20 00:00:00.157281+09:00",
          "edges": {}
        }
      ]
    }
  },
  {
    "id": 2,
    "username": "Bob",
    "created_at": "2021-09-20 00:00:00.285065+09:00",
    "updated_at": "2021-09-20 00:00:00.285066+09:00",
    "edges": {
      "owned": [
        {
          "id": 3,
          "filename": "files/file3.txt",
          "body": "SGVsbG8hIEkgYW0gZmlsZSBudW1iZXIgMy4K",
          "created_at": "2021-09-20 00:00:00.343913+09:00",
          "updated_at": "2021-09-20 00:00:00.343913+09:00",
          "edges": {}
        }
      ]
    }
  }
]
```

となる（[jq] で整形済）。 実体テーブルの `binary_files.user_owned` カラムは隠蔽されている点に注意。ログを見るとこんな感じになっていて JOIN で連結させて取ってきているわけではないようだ。

```
$ go run sample5.go
0:00AM INF Dialing PostgreSQL server host=hostname module=pgx
0:00AM INF Query args=[] module=pgx pid=28651 rowCount=2 sql="SELECT DISTINCT \"users\".\"id\", \"users\".\"username\", \"users\".\"created_at\", \"users\".\"updated_at\" FROM \"users\""
0:00AM INF Query args=[1,2] module=pgx pid=28651 rowCount=3 sql="SELECT DISTINCT \"binary_files\".\"id\", \"binary_files\".\"filename\", \"binary_files\".\"body\", \"binary_files\".\"created_at\", \"binary_files\".\"updated_at\", \"binary_files\".\"user_owned\" FROM \"binary_files\" WHERE \"user_owned\" IN ($1, $2)"
0:00AM INF closed connection module=pgx pid=28651
```

まぁ [GORM] の Preload() オプションも JOIN して取ってきてるわけじゃないし，こんなもんかな。オンライン・ドキュメントによると

>Since a query-builder can load more than one association, it's not possible to load them using one `JOIN` operation. Therefore, `ent` executes additional queries for loading associations. One query for `M2O/O2M` and `O2O` edges, and 2 queries for loading `M2M` edges.
>
>Note that, we expect to improve this in the next versions of `ent`.
(via “[Eager Loading | ent](https://entgo.io/docs/eager-load/)”)

とのこと。 [ent] 先生の次回作にご期待ください，といったところだろうか（笑）

## Update

Bob が所有する BinaryFile の内容を file4.txt に入れ替える処理。

```go:sample6.go
import (
    "context"
    "fmt"
    "os"
    "sample/dbconn"
    "sample/ent"
    "sample/ent/binaryfile"
    "sample/ent/user"
    "sample/files"

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

    file4 := "files/file4.txt"
    bin4, err := files.GetBinary(file4)
    if err != nil {
        entCtx.GetLogger().Error().Interface("error", errs.Wrap(err)).Send()
        return exitcode.Abnormal
    }

    // search data and update
    if err := entCtx.Transaction(context.TODO(), func(tx *ent.Tx) error {
        ct, err := tx.BinaryFile.Update().Where(
            binaryfile.HasOwnerWith(user.Username("Bob")),
        ).SetFilename(file4).SetBody(bin4).Save(context.TODO())
        if err != nil {
            return errs.Wrap(err)
        }
        if ct <= 0 {
            return errs.New("not change record", errs.WithContext("username", "Bob"))
        }
        return nil
    }); err != nil {
        entCtx.GetLogger().Error().Interface("error", errs.Wrap(err)).Send()
        return exitcode.Abnormal
    }

    return exitcode.Normal
}
```

HasOwnerWith() 関数を使って `owner` エッヂの関係を利用しているのがポイント。実行結果は以下の通り。

```
$ go run sample6.go
0:00AM INF begining transaction
0:00AM INF Dialing PostgreSQL server host=hostname module=pgx
0:00AM INF Exec args=[] commandTag=QkVHSU4= module=pgx pid=31502 sql=begin
0:00AM INF Exec args=["files/file4.txt","48656c6c6f21204920616d2066696c65206e756d62657220342e0a","Bob"] commandTag=VVBEQVRFIDE= module=pgx pid=31502 sql="UPDATE \"binary_files\" SET \"filename\" = $1, \"body\" = $2 WHERE \"binary_files\".\"user_owned\" IN (SELECT \"users\".\"id\" FROM \"users\" WHERE \"users\".\"username\" = $3)"
0:00AM INF committing transaction
0:00AM INF Exec args=[] commandTag=Q09NTUlU module=pgx pid=31502 sql=commit
0:00AM INF closed connection module=pgx pid=31502
```

おっ。ちゃんと副問合せを使ってるのか。えらいえらい。

## Delete

sample6.go の Update() メソッドを Delete() に変えれば同じ条件でレコードの削除ができる。

```go:sample7.go
// search data and delete
if err := entCtx.Transaction(context.TODO(), func(tx *ent.Tx) error {
    ct, err := tx.BinaryFile.Delete().Where(
        binaryfile.HasOwnerWith(user.Username("Bob")),
    ).Exec(context.TODO())
    if err != nil {
        return errs.Wrap(err)
    }
    if ct <= 0 {
        return errs.New("not delete record", errs.WithContext("username", "Bob"))
    }
    return nil
}); err != nil {
    entCtx.GetLogger().Error().Interface("error", errs.Wrap(err)).Send()
    return exitcode.Abnormal
}
```

実行結果は以下の通り。

```
$ go run sample7.go
0:00AM INF begining transaction
0:00AM INF Dialing PostgreSQL server host=hostname module=pgx
0:00AM INF Exec args=[] commandTag=QkVHSU4= module=pgx pid=13871 sql=begin
0:00AM INF Exec args=["Bob"] commandTag=REVMRVRFIDE= module=pgx pid=13871 sql="DELETE FROM \"binary_files\" WHERE \"binary_files\".\"user_owned\" IN (SELECT \"users\".\"id\" FROM \"users\" WHERE \"users\".\"username\" = $1)"
0:00AM INF committing transaction
0:00AM INF Exec args=[] commandTag=Q09NTUlU module=pgx pid=13871 sql=commit
0:00AM INF closed connection module=pgx pid=13871
```

問題なし。ちなみに

```go:sample7-1.go
ct, err := tx.BinaryFile.Delete().Exec(context.TODO())
```

と無条件で Delete() メソッドを実行したらどうなるのかと思ったら。

```
$ go run sample7-1.go
0:00AM INF begining transaction
0:00AM INF Dialing PostgreSQL server host=hostname module=pgx
0:00AM INF Exec args=[] commandTag=QkVHSU4= module=pgx pid=17631 sql=begin
0:00AM INF Exec args=[] commandTag=REVMRVRFIDI= module=pgx pid=17631 sql="DELETE FROM \"binary_files\""
0:00AM INF committing transaction
0:00AM INF Exec args=[] commandTag=Q09NTUlU module=pgx pid=17631 sql=commit
0:00AM INF closed connection module=pgx pid=17631
```

普通に全件削除された。怖い怖い（笑）

[ent]: https://entgo.io/
[GORM]: https://gorm.io/ "GORM - The fantastic ORM library for Golang, aims to be developer friendly."
[PostgreSQL]: https://www.postgresql.org/ "PostgreSQL: The world's most advanced open source database"
[ElephantSQL]: https://www.elephantsql.com/ "ElephantSQL - PostgreSQL as a Service"
[Go]: https://go.dev/
[github.com/ent/ent]: https://github.com/ent/ent "ent/ent: An entity framework for Go"
[database/sql]: https://pkg.go.dev/database/sql "sql package - database/sql - pkg.go.dev"
[github.com/jackc/pgx]: https://github.com/jackc/pgx "jackc/pgx: PostgreSQL driver and toolkit for Go"
[jq]: https://stedolan.github.io/jq/

