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

前後関係がおかしい気はするが，とにかくちゃんとした SQL 文を発行しようとしているのが分かる。では，本当に実行してみる。うーやーたー！

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

という感じに [pgx][github.com/jackc/pgx] レベルのログで正しく SQL 文が発行され coomit まで完了していることが分かる。 [ElephantSQL] の SQL Browser でも確認できたので，大丈夫だろう。 `binary_files` テーブルのデータが（slice を読み取って）複数ちゃんと作成できてる点に注目してほしい。

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

`users` テーブルのレコードが1つしかないのは寂しいのでもう一つ追加しておくか。

```go:sample5c.go
file3 := "files/file3.txt"
bin3, err := files.GetBinary(file3)
if err != nil {
    gormCtx.GetLogger().Error().Interface("error", errs.Wrap(err)).Send()
    return exitcode.Abnormal
}

data := &model.User{
    Username: "Bob",
    BinaryFiles: []model.BinaryFile{
        {Filename: file3, Body: bin3},
    },
}
```

## Read (Query)

今度は前項で作ったデータを読み出してみる。

その前に読み出したデータを JSON 形式で出力する関数を作っておこうか。

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
0:00AM INF Exec args=[] commandTag=null module=pgx pid=16676 sql=;
0:00AM INF Query args=[] module=pgx pid=16676 rowCount=2 sql="SELECT * FROM \"users\" WHERE \"users\".\"deleted_at\" IS NULL"
[{"ID":1,"CreatedAt":"2021-09-20T00:00:00.041478+09:00","UpdatedAt":"2021-09-20T00:00:00.041478+09:00","DeletedAt":null,"Username":"Alice","BinaryFiles":null},{"ID":3,"CreatedAt":"2021-09-20T00:00:00.282635+09:00","UpdatedAt":"2021-09-20T00:00:00.282635+09:00","DeletedAt":null,"Username":"Bob","BinaryFiles":null}]
0:00AM INF closed connection module=pgx pid=16676
```

問題なさそうだね。当然ではあるが User.BinaryFiles のフィールドには何も入らないので null になっている。

では User.BinaryFiles にも値を詰めたい場合はどうするかというと Preload() オプションを使う。

```go:sample7.go
// select all records (with preload)
data := []model.User{}
tx := gormCtx.GetDb().WithContext(context.TODO()).Preload(clause.Associations).Find(&data)
if tx.Error != nil {
    gormCtx.GetLogger().Error().Interface("error", errs.Wrap(tx.Error)).Send()
    return exitcode.Abnormal
}
```

（clause パッケージのパスは `"gorm.io/gorm/clause"`）

これは dry run が上手く働かないみたいなので，いきなり実行した。実行結果は以下の通り。

```
$ go run sample7.go
0:00AM INF Dialing PostgreSQL server host=hostname module=pgx
0:00AM INF Exec args=[] commandTag=null module=pgx pid=17949 sql=;
0:00AM INF Query args=[] module=pgx pid=17949 rowCount=2 sql="SELECT * FROM \"users\" WHERE \"users\".\"deleted_at\" IS NULL"
0:00AM INF Query args=[1,3] module=pgx pid=17949 rowCount=3 sql="SELECT * FROM \"binary_files\" WHERE \"binary_files\".\"user_id\" IN ($1,$2) AND \"binary_files\".\"deleted_at\" IS NULL"
[{"ID":1,"CreatedAt":"2021-09-20T00:00:00.041478+09:00","UpdatedAt":"2021-09-20T00:00:00.041478+09:00","DeletedAt":null,"Username":"Alice","BinaryFiles":[{"ID":1,"CreatedAt":"2021-09-20T00:00:00.118119+09:00","UpdatedAt":"2021-09-20T00:00:00.118119+09:00","DeletedAt":null,"UserId":1,"Filename":"files/file1.txt","Body":"SGVsbG8hIEkgYW0gZmlsZSBudW1iZXIgMS4K"},{"ID":2,"CreatedAt":"2021-09-20T00:00:00.118119+09:00","UpdatedAt":"2021-09-20T00:00:00.118119+09:00","DeletedAt":null,"UserId":1,"Filename":"files/file2.txt","Body":"SGVsbG8hIEkgYW0gZmlsZSBudW1iZXIgMi4K"}]},{"ID":3,"CreatedAt":"2021-09-20T00:00:00.282635+09:00","UpdatedAt":"2021-09-20T00:00:00.282635+09:00","DeletedAt":null,"Username":"Bob","BinaryFiles":[{"ID":3,"CreatedAt":"2021-09-20T00:00:00.361274+09:00","UpdatedAt":"2021-09-20T00:00:00.361274+09:00","DeletedAt":null,"UserId":3,"Filename":"files/file3.txt","Body":"SGVsbG8hIEkgYW0gZmlsZSBudW1iZXIgMy4K"}]}]
0:00AM INF closed connection module=pgx pid=17949
```

ん？ バイナリデータは，符号化されてるのか？ いや [Go] の [encoding/json] パッケージの仕様か[^jsn1]。

[^jsn1]: JSON の公式な仕様にはバイナリデータの書式については言及されてないらしい。

>Array and slice values encode as JSON arrays, except that []byte encodes as a base64-encoded string, and a nil slice encodes as the null JSON value.
(via “[json package - encoding/json - pkg.go.dev](https://pkg.go.dev/encoding/json#Marshal)”)

上の結果を更に "Alice" で絞り込みたいなら Where() メソッドが使える。

```go:sample7b.go
// select data for 'Alice' (with preload)
data := []model.User{}
tx := gormCtx.GetDb().WithContext(context.TODO()).Preload(clause.Associations).Where(&model.User{Username: "Alice"}).Find(&data)
if tx.Error != nil {
    gormCtx.GetLogger().Error().Interface("error", errs.Wrap(tx.Error)).Send()
    return exitcode.Abnormal
}
```

これの実行結果はこんな感じ。

```
$ go run sample7b.go 
0:00AM INF Dialing PostgreSQL server host=hostname module=pgx
0:00AM INF Exec args=[] commandTag=null module=pgx pid=26015 sql=;
0:00AM INF Query args=["Alice"] module=pgx pid=26015 rowCount=1 sql="SELECT * FROM \"users\" WHERE \"users\".\"username\" = $1 AND \"users\".\"deleted_at\" IS NULL"
0:00AM INF Query args=[1] module=pgx pid=26015 rowCount=2 sql="SELECT * FROM \"binary_files\" WHERE \"binary_files\".\"user_id\" = $1 AND \"binary_files\".\"deleted_at\" IS NULL"
[{"ID":1,"CreatedAt":"2021-09-20T00:00:00.041478+09:00","UpdatedAt":"2021-09-20T00:00:00.041478+09:00","DeletedAt":null,"Username":"Alice","BinaryFiles":[{"ID":1,"CreatedAt":"2021-09-20T00:00:00.118119+09:00","UpdatedAt":"2021-09-20T00:00:00.118119+09:00","DeletedAt":null,"UserId":1,"Filename":"files/file1.txt","Body":"SGVsbG8hIEkgYW0gZmlsZSBudW1iZXIgMS4K"},{"ID":2,"CreatedAt":"2021-09-20T00:00:00.118119+09:00","UpdatedAt":"2021-09-20T00:00:00.118119+09:00","DeletedAt":null,"UserId":1,"Filename":"files/file2.txt","Body":"SGVsbG8hIEkgYW0gZmlsZSBudW1iZXIgMi4K"}]}]
0:00AM INF closed connection module=pgx pid=26015
```

他にも Where() メソッドと Or() メソッドを組み合わせたり Select() メソッドでカラム名を指定したり Table() メソッドでテーブル名を指定したり，更に Order(), Limit(), Offset(), Group(), Having(), Distinct(), Joins() といったメソッドも提供されていて SQL の組み立てに関してはかなり自由度が高い。何なら Raw() メソッドを使って

```go
var data []model.User
tx := gormCtx.GetDb().WithContext(context.TODO()).Raw("SELECT id, username FROM users WHERE username = ?", "Alice").Scan(&data)
```

なんてなこともできる[^cud1]。SQLスキー な人には天国だろう（笑）

[^cud1]: DROP/INSERT/UPDATE/DELETE 用に Exec() メソッドも用意されている。

## Update

Save() メソッドを使えばレコード全体をまるっと更新してくれる。

```go:sample8b.go
// edit and uodate
var data model.User
tx := gormCtx.GetDb().WithContext(context.TODO()).Where(&model.User{Username: "Bob"}).First(&data)
if tx.Error != nil {
    gormCtx.GetLogger().Error().Interface("error", errs.Wrap(tx.Error)).Send()
    return exitcode.Abnormal
}
data.Username = "Bob 2nd"
tx = gormCtx.GetDb().WithContext(context.TODO()).Save(&data)
if tx.Error != nil {
    gormCtx.GetLogger().Error().Interface("error", errs.Wrap(tx.Error)).Send()
    return exitcode.Abnormal
}
```

これの実行結果はこんな感じ。

```
$ go run sample8b.go 
8:10PM INF Dialing PostgreSQL server host=hostname module=pgx
8:10PM INF Exec args=[] commandTag=null module=pgx pid=11183 sql=;
8:10PM INF Query args=["Bob"] module=pgx pid=11183 rowCount=1 sql="SELECT * FROM \"users\" WHERE \"users\".\"username\" = $1 AND \"users\".\"deleted_at\" IS NULL ORDER BY \"users\".\"id\" LIMIT 1"
8:10PM INF Exec args=[] commandTag=QkVHSU4= module=pgx pid=11183 sql=begin
8:10PM INF Exec args=["2021-09-20T00:00:00.282635+09:00","2021-09-20T20:10:55.779497619+09:00",null,"Bob 2nd",3] commandTag=VVBEQVRFIDE= module=pgx pid=11183 sql="UPDATE \"users\" SET \"created_at\"=$1,\"updated_at\"=$2,\"deleted_at\"=$3,\"username\"=$4 WHERE \"id\" = $5"
8:10PM INF Exec args=[] commandTag=Q09NTUlU module=pgx pid=11183 sql=commit
8:10PM INF closed connection module=pgx pid=11183
```

Save() メソッドでは引数で与えられた構造体の各フィールド（primary key 以外）を全て更新しようとする。なお `users.deleted_at` については「現在時刻」で更新されている点に注目。

Updates() メソッドを使うと指定したカラムのみを更新できる。

```go:sample8c.go
file4 := "files/file4.txt"
bin4, err := files.GetBinary(file4)
if err != nil {
    gormCtx.GetLogger().Error().Interface("error", errs.Wrap(err)).Send()
    return exitcode.Abnormal
}
// select data for 'Bob 2nd' (with preload)
data := []model.User{}
tx := gormCtx.GetDb().WithContext(context.TODO()).Preload(clause.Associations).Where(&model.User{Username: "Bob 2nd"}).Find(&data)
if tx.Error != nil {
    gormCtx.GetLogger().Error().Interface("error", errs.Wrap(tx.Error)).Send()
    return exitcode.Abnormal
}
// update data in binary_files table
tx = gormCtx.GetDb().WithContext(context.TODO()).Model(&data[0].BinaryFiles[0]).Updates(model.BinaryFile{Filename: file4, Body: bin4})
if tx.Error != nil {
    gormCtx.GetLogger().Error().Interface("error", errs.Wrap(tx.Error)).Send()
    return exitcode.Abnormal
}
```

この例では `Model(&data[0].BinaryFiles[0])` の引数構造体データの primary key フィールドを条件にしている。また，上のコードのように Updates() の引数に構造体データをセットする場合，ゼロ値のフィールドは更新対象にならない（つまりゼロ値への更新は出来ない）ので注意。この場合は map[string]interface{} 型の連想配列を使うとよいだろう。

更に Model() メソッドの引数の構造体が空だと全件が対象になってしまうので注意。この場合は Where() メソッドで条件を指定してやればよい。

実行結果を見てみよう。

```
$ go run sample8c.go 
9:00PM INF Dialing PostgreSQL server host=hostname module=pgx
9:00PM INF Exec args=[] commandTag=null module=pgx pid=10480 sql=;
9:00PM INF Query args=["Bob 2nd"] module=pgx pid=10480 rowCount=1 sql="SELECT * FROM \"users\" WHERE \"users\".\"username\" = $1 AND \"users\".\"deleted_at\" IS NULL"
9:00PM INF Query args=[3] module=pgx pid=10480 rowCount=1 sql="SELECT * FROM \"binary_files\" WHERE \"binary_files\".\"user_id\" = $1 AND \"binary_files\".\"deleted_at\" IS NULL"
9:00PM INF Exec args=[] commandTag=QkVHSU4= module=pgx pid=10480 sql=begin
9:00PM INF Exec args=["2021-09-20T21:00:21.811307485+09:00","files/file4.txt","48656c6c6f21204920616d2066696c65206e756d62657220342e0a",3] commandTag=VVBEQVRFIDE= module=pgx pid=10480 sql="UPDATE \"binary_files\" SET \"updated_at\"=$1,\"filename\"=$2,\"body\"=$3 WHERE \"id\" = $4"
9:00PM INF Exec args=[] commandTag=Q09NTUlU module=pgx pid=10480 sql=commit
9:00PM INF closed connection module=pgx pid=10480
```

うんうん。上手に出来ました。 `binary_files.updated_at` カラムは更新対象として指定されてないけど，現在時刻で更新されている点に注目。

## Delete

Delete は怖いので dry run で。 Delete() メソッドを使えば引数によって指定された条件でレコードを削除できるのだが...

```go:sample9.go
// select data for 'Bob 2nd'
data := model.User{}
tx := gormCtx.GetDb().WithContext(context.TODO()).Where(&model.User{Username: "Bob 2nd"}).First(&data)
if tx.Error != nil {
    gormCtx.GetLogger().Error().Interface("error", errs.Wrap(tx.Error)).Send()
    return exitcode.Abnormal
}

// delete data in users table (dry run)
tx = gormCtx.GetDb().Session(&gorm.Session{DryRun: true}).Delete(&data)
if tx.Error != nil {
    gormCtx.GetLogger().Error().Interface("error", errs.Wrap(tx.Error)).Send()
    return exitcode.Abnormal
}
```

[GORM] のログが

```
[99.438ms] [rows:0] UPDATE "users" SET "deleted_at"='2021-09-20 00:00:00.00' WHERE "users"."id" = 3 AND "users"."deleted_at" IS NULL
```

となっている。どうやら `deleted_at` カラムがあると論理削除になるようだ。この構成で物理削除がしたいなら。

```go:sample9b.go
tx = gormCtx.GetDb().Session(&gorm.Session{DryRun: true}).Unscoped().Delete(&data)
```

と Unscoped() メソッドを噛ませればいいようだ。これで

```
[119.670ms] [rows:0] DELETE FROM "users" WHERE "users"."id" = 3
```

となった。ちなみに，検索（SELECT）のときも Unscoped() メソッドを噛ませれば `deleted_at` カラムを無視してくれるらしい。もっとも本当に

```go:sample9c.go
// select data for 'Bob 2nd'
data := model.User{}
tx := gormCtx.GetDb().WithContext(context.TODO()).Unscoped().Where(&model.User{Username: "Bob 2nd"}).First(&data)
if tx.Error != nil {
    gormCtx.GetLogger().Error().Interface("error", errs.Wrap(tx.Error)).Send()
    return exitcode.Abnormal
}
// delete data in users table (dry run)
tx = gormCtx.GetDb().WithContext(context.TODO()).Unscoped().Delete(&data)
if tx.Error != nil {
    gormCtx.GetLogger().Error().Interface("error", errs.Wrap(tx.Error)).Send()
    return exitcode.Abnormal
}
```

で物理削除しようとしたら

```
$ go run sample9c.go 
0:00AM INF Dialing PostgreSQL server host=hostname module=pgx
0:00AM INF Exec args=[] commandTag=null module=pgx pid=7578 sql=;
0:00AM INF Query args=["Bob 2nd"] module=pgx pid=7578 rowCount=1 sql="SELECT * FROM \"users\" WHERE \"users\".\"username\" = $1 ORDER BY \"users\".\"id\" LIMIT 1"
0:00AM INF Exec args=[] commandTag=QkVHSU4= module=pgx pid=7578 sql=begin
0:00AM ERR Exec args=[3] err="ERROR: update or delete on table \"users\" violates foreign key constraint \"fk_users_binary_files\" on table \"binary_files\" (SQLSTATE 23503)" module=pgx pid=7578 sql="DELETE FROM \"users\" WHERE \"users\".\"id\" = $1"
0:00AM INF Exec args=[] commandTag=Uk9MTEJBQ0s= module=pgx pid=7578 sql=rollback
0:00AM ERR  error={"Context":{"function":"main.Run"},"Err":{"Msg":"ERROR: update or delete on table \"users\" violates foreign key constraint \"fk_users_binary_files\" on table \"binary_files\" (SQLSTATE 23503)","Type":"*pgconn.PgError"},"Type":"*errs.Error"}
0:00AM INF closed connection module=pgx pid=7578
```

などと foreign key のせいでエラーになったけどね（笑）

ところで，全件削除しようとして

```go:sample9d.go
tx := gormCtx.GetDb().Session(&gorm.Session{DryRun: true}).Delete(&model.User{})
```

と書いたら [GORM] に “WHERE conditions required” と怒られてエラーになった。安易な全件削除はアカンらしい。でも，明示的に生の SQL 文で

```go
tx := gormCtx.GetDb().Session(&gorm.Session{DryRun: true}).Exec("UPDATE users SET deleted_at=now() WHERE deleted_at IS NULL") // soft delete
```

または

```go
tx := gormCtx.GetDb().Session(&gorm.Session{DryRun: true}).Exec("DELETE FROM users") // delete permanently
```

と書けば行けるようだ。なんだかなぁ。

## トランザクション処理

ひとつのトランザクションの中で複数の CRUD 処理を行いたいことも当然ある。 [GORM] では Begin(), Commit(), Rollback() といった伝統的なメソッドも用意されている。

```go
tx := gormCtx.GetDb().WithContext(context.TODO()).Begin()
result := tx.Create(data)
if result.Error != nil {
    tx.Rollback()
    return exitcode.Abnormal
}
tx.Commit()
return exitcode.Normal
```

しかし [GORM] には Transaction() メソッドが用意されていて，これがかなり秀逸である。実際にコードを書いたほうが早いだろう。

```go:sample10.go
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
file3 := "files/file3.txt"
bin3, err := files.GetBinary(file3)
if err != nil {
    gormCtx.GetLogger().Error().Interface("error", errs.Wrap(err)).Send()
    return exitcode.Abnormal
}
data1 := &model.User{
    Username: "Alice",
    BinaryFiles: []model.BinaryFile{
        {Filename: file1, Body: bin1},
        {Filename: file2, Body: bin2},
    },
}
data2 := &model.User{
    Username: "Bob",
    BinaryFiles: []model.BinaryFile{
        {Filename: file3, Body: bin3},
    },
}

if err := gormCtx.GetDb().WithContext(context.TODO()).Transaction(func(tx *gorm.DB) error {
    if err := tx.Create(data1).Error; err != nil {
        return errs.Wrap(err) // return any error will rollback
    }
    if err := tx.Create(data2).Error; err != nil {
        return errs.Wrap(err) // return any error will rollback
    }
    return nil // return nil will commit the whole transaction
}); err != nil {
    gormCtx.GetLogger().Error().Interface("error", errs.Wrap(err)).Send()
    return exitcode.Abnormal
}
```

上のコードはバイナリデータを全てメモリ上に保持ってしまってるが，良い子はマネしないように（笑） このコードの実行結果は以下の通り。

```
$ go run sample10.go 
0:00AM INF Dialing PostgreSQL server host=hostname module=pgx
0:00AM INF Exec args=[] commandTag=null module=pgx pid=23629 sql=;
0:00AM INF Exec args=[] commandTag=QkVHSU4= module=pgx pid=23629 sql=begin
0:00AM INF Query args=["2021-09-20T00:00:00.57261701+09:00","2021-09-20T00:00:00.57261701+09:00",null,"Alice"] module=pgx pid=23629 rowCount=1 sql="INSERT INTO \"users\" (\"created_at\",\"updated_at\",\"deleted_at\",\"username\") VALUES ($1,$2,$3,$4) RETURNING \"id\""
0:00AM INF Query args=["2021-09-20T00:00:00.677252677+09:00","2021-09-20T00:00:00.677252677+09:00",null,1,"files/file1.txt","48656c6c6f21204920616d2066696c65206e756d62657220312e0a","2021-09-20T00:00:00.677252677+09:00","2021-09-20T00:00:00.677252677+09:00",null,1,"files/file2.txt","48656c6c6f21204920616d2066696c65206e756d62657220322e0a"] module=pgx pid=23629 rowCount=2 sql="INSERT INTO \"binary_files\" (\"created_at\",\"updated_at\",\"deleted_at\",\"user_id\",\"filename\",\"body\") VALUES ($1,$2,$3,$4,$5,$6),($7,$8,$9,$10,$11,$12) ON CONFLICT (\"id\") DO UPDATE SET \"user_id\"=\"excluded\".\"user_id\" RETURNING \"id\""
0:00AM INF Query args=["2021-09-20T00:00:00.820487738+09:00","2021-09-20T00:00:00.820487738+09:00",null,"Bob"] module=pgx pid=23629 rowCount=1 sql="INSERT INTO \"users\" (\"created_at\",\"updated_at\",\"deleted_at\",\"username\") VALUES ($1,$2,$3,$4) RETURNING \"id\""
0:00AM INF Query args=["2021-09-20T00:00:00.915438304+09:00","2021-09-20T00:00:00.915438304+09:00",null,2,"files/file3.txt","48656c6c6f21204920616d2066696c65206e756d62657220332e0a"] module=pgx pid=23629 rowCount=1 sql="INSERT INTO \"binary_files\" (\"created_at\",\"updated_at\",\"deleted_at\",\"user_id\",\"filename\",\"body\") VALUES ($1,$2,$3,$4,$5,$6) ON CONFLICT (\"id\") DO UPDATE SET \"user_id\"=\"excluded\".\"user_id\" RETURNING \"id\""
0:00AM INF Exec args=[] commandTag=Q09NTUlU module=pgx pid=23629 sql=commit
0:00AM INF closed connection module=pgx pid=23629
```

このように commit や rollback に関するハンドリングは Transaction() 側に丸投げでき Transaction() を呼び出した側は error ハンドリングに専念できる。素晴らしい！

トランザクション処理では成功時と失敗時で後始末が異なるので（defer が使えないため）すこぶる鬱陶しいのだが，こうやって一連の処理をリテラル関数で括ってしまえばいいのか。どこぞの try-catch よりはだいぶマシなアイデアかな。これは覚えておこう。

[GORM]: https://gorm.io/ "GORM - The fantastic ORM library for Golang, aims to be developer friendly."
[PostgreSQL]: https://www.postgresql.org/ "PostgreSQL: The world's most advanced open source database"
[ElephantSQL]: https://www.elephantsql.com/ "ElephantSQL - PostgreSQL as a Service"
[database/sql]: https://pkg.go.dev/database/sql "sql package - database/sql - pkg.go.dev"
[io]: https://pkg.go.dev/io "io package - io - pkg.go.dev"
[encoding/json]: https://pkg.go.dev/encoding/json "json package - encoding/json - pkg.go.dev"
[github.com/jackc/pgx]: https://github.com/jackc/pgx "jackc/pgx: PostgreSQL driver and toolkit for Go"
[Go]: https://go.dev/
