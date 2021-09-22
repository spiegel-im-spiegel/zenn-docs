---
title: "GORM で CRUD"
---

前節で [PostgreSQL] 上にテーブルを作成するところまでできたので，ひととおり CRUD (Create/Read/Update/Delete) を試してみよう。

## Create

まずはデータの INSERT (Create) から。ファイルの内容をバイナリデータとしてアップロードしたいので，ファイルアクセス用の files サブパッケージとアップロード用の関数を書く。

```go:files/files.go
package files

import (
	"io"
	"os"

	"github.com/spiegel-im-spiegel/errs"
)

func GetBinary(path string) ([]byte, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, errs.Wrap(err, errs.WithContext("path", path))
	}
	defer f.Close()
	b, err := io.ReadAll(f)
	return b, errs.Wrap(err, errs.WithContext("path", path))
}
```

せめて [io].ReadCloser を渡せるといいんだけどねぇ（愚痴）。

これを使って前節の model.User 構造体に値を詰め込んで Create() メソッドをキックする。最初はやっぱり怖いので dry run で試してみる（ファイルは指定のパスに実際にあるものとする）。

```go:sample5.go
func Run() exitcode.ExitCode {
	// create gorm.DB instance for PostgreSQL service
	gormCtx, err := orm.NewGORM()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return exitcode.Abnormal
	}
	defer gormCtx.Close()

	file1 := "files/file1.txt"
	bin1, err := files.GetBinary(file1)
	if err != nil {
		gormCtx.GetLogger().Error().Interface("error", errs.Wrap(err)).Send()
		return exitcode.Abnormal
	}
	file2 := "files/file2.txt"
	bin2, err := files.GetBinary(file2)
	if err != nil {
		gormCtx.GetLogger().Error().Interface("error", errs.Wrap(err)).Send()
		return exitcode.Abnormal
	}

	data := &model.User{
		Username: "Alice",
		BinaryFiles: []model.BinaryFile{
			{Filename: file1, Body: bin1},
			{Filename: file2, Body: bin2},
		},
	}

	// insert data (dry run)
	tx := gormCtx.GetDb().Session(&gorm.Session{DryRun: true}).Create(data)
	if tx.Error != nil {
		gormCtx.GetLogger().Error().Interface("error", errs.Wrap(tx.Error)).Send()
		return exitcode.Abnormal
	}

	return exitcode.Normal
}
```

これの実行結果のログを見ると（[GORM] が出力している部分のみ）

```
[0.000ms] [rows:0] INSERT INTO "binary_files" ("created_at","updated_at","deleted_at","user_id","filename","body") VALUES ('2021-09-20 00:00:00.000','2021-09-20 00:00:00.000',NULL,0,'files/file1.txt','<binary>'),('2021-09-20 00:00:00.000','2021-09-20 00:00:00.000',NULL,0,'files/file2.txt','<binary>') ON CONFLICT ("id") DO UPDATE SET "user_id"="excluded"."user_id" RETURNING "id"
[76.521ms] [rows:0] INSERT INTO "users" ("created_at","updated_at","deleted_at","username") VALUES ('2021-09-20 00:00:00.000','2021-09-20 00:00:00.000',NULL,'Alice') RETURNING "id"
```

という感じで前後関係がおかしい気がするが，とにかくちゃんとした SQL 文を発行しようとしているのが分かる。では，本当に実行してみる。うーやーたー！

```go:sample5b.go
// insert data
tx := gormCtx.GetDb().WithContext(context.TODO()).Create(data)
if tx.Error != nil {
	gormCtx.GetLogger().Error().Interface("error", errs.Wrap(tx.Error)).Send()
	return exitcode.Abnormal
}
```

実行結果は

```
$ go run sample5b.go
0:00AM INF Dialing PostgreSQL server host=hostname module=pgx
0:00AM INF Exec args=[] commandTag=null module=pgx pid=15111 sql=;
0:00AM INF Exec args=[] commandTag=QkVHSU4= module=pgx pid=15111 sql=begin
0:00AM INF Query args=["2021-09-20T00:00:00.0414789+09:00","2021-09-20T00:00:00.0414789+09:00",null,"Alice"] module=pgx pid=15111 rowCount=1 sql="INSERT INTO \"users\" (\"created_at\",\"updated_at\",\"deleted_at\",\"username\") VALUES ($1,$2,$3,$4) RETURNING \"id\""
0:00AM INF Query args=["2021-09-20T00:00:00.118119+09:00","2021-09-20T00:00:00.118119+09:00",null,1,"files/file1.txt","48656c6c6f21204920616d2066696c65206e756d62657220312e0a","2021-09-20T00:00:00.118119+09:00","2021-09-20T00:00:00.118119+09:00",null,1,"files/file2.txt","48656c6c6f21204920616d2066696c65206e756d62657220322e0a"] module=pgx pid=15111 rowCount=2 sql="INSERT INTO \"binary_files\" (\"created_at\",\"updated_at\",\"deleted_at\",\"user_id\",\"filename\",\"body\") VALUES ($1,$2,$3,$4,$5,$6),($7,$8,$9,$10,$11,$12) ON CONFLICT (\"id\") DO UPDATE SET \"user_id\"=\"excluded\".\"user_id\" RETURNING \"id\""
0:00AM INF Exec args=[] commandTag=Q09NTUlU module=pgx pid=15111 sql=commit
0:00AM INF closed connection module=pgx pid=15111
```

となり [pgx][github.com/jackc/pgx] レベルのログで正しく SQL 文が発行され coomit まで完了していることが分かる。 [ElephantSQL] のブラウザでも確認できたので，大丈夫だろう。

ちなみに同じコマンドをもう一度叩くと

```
$ go run sample5b.go
0:00AM INF Dialing PostgreSQL server host=hostname module=pgx
0:00AM INF Exec args=[] commandTag=null module=pgx pid=16052 sql=;
0:00AM INF Exec args=[] commandTag=QkVHSU4= module=pgx pid=16052 sql=begin
0:00AM ERR Query args=["2021-09-20T00:00:00.7625914+09:00","2021-09-20T00:00:00.7625914+09:00",null,"Alice"] err="ERROR: duplicate key value violates unique constraint \"users_username_key\" (SQLSTATE 23505)" module=pgx pid=16052 sql="INSERT INTO \"users\" (\"created_at\",\"updated_at\",\"deleted_at\",\"username\") VALUES ($1,$2,$3,$4) RETURNING \"id\""
0:00AM INF Exec args=[] commandTag=Uk9MTEJBQ0s= module=pgx pid=16052 sql=rollback
0:00AM ERR  error={"Context":{"function":"main.Run"},"Err":{"Msg":"ERROR: duplicate key value violates unique constraint \"users_username_key\" (SQLSTATE 23505)","Type":"*pgconn.PgError"},"Type":"*errs.Error"}
0:00AM INF closed connection module=pgx pid=16052
```

と `users.username` の一意制約違反でエラーになったのが分かる。よしよし。

## Read

今度は前項で作ったデータを読みだしてみる。

その前に読みだしたデータを JSON 形式で出力する関数を作っておこうか。

```go:files/output.go
package files

import (
	"bytes"
	"encoding/json"
	"io"
	"os"

	"github.com/spiegel-im-spiegel/errs"
)

func Output(dst io.Writer, src interface{}) error {
	buf := &bytes.Buffer{}
	enc := json.NewEncoder(buf)
	if err := enc.Encode(src); err != nil {
		return errs.Wrap(err)
	}
	if _, err := io.Copy(os.Stdout, buf); err != nil {
		return errs.Wrap(err)
	}
	return nil
}
```

準備が整ったところで，まずは全件検索から。最初はやっぱり dry run。

```go:sample6.go
// select all records (dry run)
data := []model.User{}
tx := gormCtx.GetDb().Session(&gorm.Session{DryRun: true}).Find(&data)
if tx.Error != nil {
	gormCtx.GetLogger().Error().Interface("error", errs.Wrap(tx.Error)).Send()
	return exitcode.Abnormal
}
```

[GORM] のログには

```
[1.709ms] [rows:0] SELECT * FROM "users" WHERE "users"."deleted_at" IS NULL
```

と出ている。 `users.deleted_at` は論理削除フラグとして使われているので NULL なら削除されていないということだ。では，本当に実行してみる。

```go:sample6b.go
// select all records
data := []model.User{}
tx := gormCtx.GetDb().WithContext(context.TODO()).Find(&data)
if tx.Error != nil {
	gormCtx.GetLogger().Error().Interface("error", errs.Wrap(tx.Error)).Send()
	return exitcode.Abnormal
}
// output by JSON format
if err := files.Output(os.Stdout, data); err != nil {
	gormCtx.GetLogger().Error().Interface("error", errs.Wrap(tx.Error)).Send()
	return exitcode.Abnormal
}
```

実行結果は以下。

```
$ go run sample6b.go
0:00AM INF Dialing PostgreSQL server host=hostname module=pgx
0:00AM INF Exec args=[] commandTag=null module=pgx pid=12334 sql=;
0:00AM INF Query args=[] module=pgx pid=12334 rowCount=1 sql="SELECT * FROM \"users\" WHERE \"users\".\"deleted_at\" IS NULL"
[{"ID":1,"CreatedAt":"2021-09-20T00:00:00.041478+09:00","UpdatedAt":"2021-09-20T00:00:00.041478+09:00","DeletedAt":null,"Username":"Alice","BinaryFiles":null}]
0:00AM INF closed connection module=pgx pid=12334
```

問題なさそうだね。当然ではあるが User.BinaryFiles のフィールドには何も入らないので null になっている。
















[Go]: https://go.dev/
[PostgreSQL]: https://www.postgresql.org/ "PostgreSQL: The world's most advanced open source database"
[ElephantSQL]: https://www.elephantsql.com/ "ElephantSQL - PostgreSQL as a Service"
[database/sql]: https://pkg.go.dev/database/sql "sql package - database/sql - pkg.go.dev"
[io]: https://pkg.go.dev/io "io package - io - pkg.go.dev"
[encoding/json]: https://pkg.go.dev/encoding/json "json package - encoding/json - pkg.go.dev"
[GORM]: https://gorm.io/ "GORM - The fantastic ORM library for Golang, aims to be developer friendly."
[github.com/jackc/pgx]: https://github.com/jackc/pgx "jackc/pgx: PostgreSQL driver and toolkit for Go"
