---
title: "ent で CRUD"
---

今回の調査もいよいよ大詰め。もうしばらくお付き合いください。

## トランザクション処理

いきなりトランザクション処理の話だが [ent] では User とか BinaryFile とかのノード単位で CRUD 処理を行うようで，トランザクションもその処理単位で閉じている。もちろん commit や rollback といったメソッドは用意されているが， [GORM] の CRUD 紹介の節でも述べた通り，トランザクション処理は成功時と失敗時で後始末処理が異なるのがすこぶる鬱陶しい。そこで [GORM] の Transaction() メソッドのようなものがあると便利である。 [ent] にはそういうものは見当たらないが，ないなら書けばいいのである（笑）

幸い[サンプルコード](https://entgo.io/ja/docs/transactions/#%E3%83%99%E3%82%B9%E3%83%88%E3%83%97%E3%83%A9%E3%82%AF%E3%83%86%E3%82%A3%E3%82%B9 "トランザクション | ent")があるので，これを利用させてもらおう。


## Create

[ent] って CRUD の dry run ができないのか？ ドキドキするんだが...

それはともかく [ent] の create は

```go
// create data
user, err := entCtx.GetClient().User.Create().SetUsername("Alice").Save(context.TODO())
if err != nil {
	entCtx.GetLogger().Error().Interface("error", errs.Wrap(err)).Send()
	return exitcode.Abnormal
}
...
```

みたいな感じでデータを作って Save() するのだが（SaveX() を使うと失敗時に panic を投げる。なにそれこわい）， Save() ごとにトランザクションが完結してしまうのが面白くない感じである。









[ent]: https://entgo.io/
[GORM]: https://gorm.io/ "GORM - The fantastic ORM library for Golang, aims to be developer friendly."
[PostgreSQL]: https://www.postgresql.org/ "PostgreSQL: The world's most advanced open source database"
[Go]: https://go.dev/
[github.com/ent/ent]: https://github.com/ent/ent "ent/ent: An entity framework for Go"
[database/sql]: https://pkg.go.dev/database/sql "sql package - database/sql - pkg.go.dev"
[github.com/jackc/pgx]: https://github.com/jackc/pgx "jackc/pgx: PostgreSQL driver and toolkit for Go"

