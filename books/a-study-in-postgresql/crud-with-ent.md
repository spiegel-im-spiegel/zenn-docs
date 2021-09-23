---
title: "ent で CRUD"
---

前節でちょっぴりモチベーションが下がったが CRUD (Create/Read/Update/Delete) の検証を終えてしまおう。もうしばらくお付き合いください。

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

https://entgo.io/ja/docs/transactions/#%E3%83%99%E3%82%B9%E3%83%88%E3%83%97%E3%83%A9%E3%82%AF%E3%83%86%E3%82%A3%E3%82%B9







[ent]: https://entgo.io/
[PostgreSQL]: https://www.postgresql.org/ "PostgreSQL: The world's most advanced open source database"
[Go]: https://go.dev/
[github.com/ent/ent]: https://github.com/ent/ent "ent/ent: An entity framework for Go"
[database/sql]: https://pkg.go.dev/database/sql "sql package - database/sql - pkg.go.dev"
[github.com/jackc/pgx]: https://github.com/jackc/pgx "jackc/pgx: PostgreSQL driver and toolkit for Go"

