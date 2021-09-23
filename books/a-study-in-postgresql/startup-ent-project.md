---
title: "ent プロジェクトの作成"
---

前節で述べたように [ent] は go generate を使ってコードを自動生成する。言い方を変えると自動生成するところまで持っていかないと [PostgreSQL] への接続コードも書けない。というわけで，まずは作業環境を作るところから始めよう。

まずは作業用に適当なディレクトリを掘って以下のコマンドで初期化する。

```
$ go mod init sample
go: creating new go.mod: module sample
```

モジュール名は `sample` で（笑） 次に entc ([ent] CLI ツール) を導入するのだが，不用意に bin ディレクトリを汚されたくないので

```
$ go get -d entgo.io/ent/cmd/ent@latest
```

と -d オプションを付けて実行する。こうすれば必要なモジュールを go.mod ファイルに追記しダウンロードまでは行うがビルド&インストールは行わない。ちなみに [Go] 1.17 で go get を使ってインストールしようとすると

```
go get: installing executables with 'go get' in module mode is deprecated.
    To adjust and download dependencies of the current module, use 'go get -d'.
    To install using requirements of the current module, use 'go install'.
    To install ignoring the current module, use 'go install' with a version,
    like 'go install example.com/cmd@latest'.
    For more information, see https://golang.org/doc/go-get-install-deprecation
    or run 'go help get' or 'go help install'.
```

と[警告](https://golang.org/doc/go-get-install-deprecation "Deprecation of 'go get' for installing executables - The Go Programming Language")が表示される。

これで準備できたので初期コードを生成してみる。

```
$ go run -mod=mod entgo.io/ent/cmd/ent init User BinaryFile
```

これで ent/ ディレクトリおよびその配下にファイル・ディレクトリが生成される。ちなみに -mod=mod オプションは entgo.io/ent/cmd/ent コマンドを go.mod ファイルに記述されているモジュール・バージョンで起動しろということらしい。こうしておけば entc とパッケージとして組み込まれる [ent] モジュールのバージョンを go.mod ファイルでコントロールすることができる。

なお，上のコマンドを実行した後に

```
$ go mod tidy
```

としておけば go.mod ファイルがきれいに整理される。

さて ent/ ディレクトリ以下を眺めると

```
$ tree ent/
ent/
├── generate.go
└── schema
    ├── binaryfile.go
    └── user.go
```

となっている。 ent/generate.go は go generate で起動する処理が記述されている。こんな感じ。

```go:ent/generate.go
package ent

//go:generate go run -mod=mod entgo.io/ent/cmd/ent generate ./schema
```

ent/schema/user.go ファイルや ent/schema/binaryfile.go ファイルはスキーマを定義するもので， [GORM] などの Model 定義とは異なり，構造体のフィールドではなく，紐付けられたメソッドで定義を行う。スキーマ定義を追加したい場合には

```
$ go run -mod=mod entgo.io/ent/cmd/ent init Foo
```

などとすれば既存のファイルは変更せず，指定した名前のスキーマ定義ファイルのみ追加生成するようだ。お手軽に使えるのはよい。

たとえば ent/schema/user.go ファイルの初期状態はこんな感じ。

```go:ent/schema/user.go
package schema

import "entgo.io/ent"

// User holds the schema definition for the User entity.
type User struct {
    ent.Schema
}

// Fields of the User.
func (User) Fields() []ent.Field {
    return nil
}

// Edges of the User.
func (User) Edges() []ent.Edge {
    return nil
}
```

[ent].Schema 型は以下のように定義されていて， Fields(), Edges() 以外にもいくつかのメソッドを定義できるようになっている（全てを定義する必要はない）。

```go:ent.go
type (
    // The Interface type describes the requirements for an exported type defined in the schema package.
    // It functions as the interface between the user's schema types and codegen loader.
    // Users should use the Schema type for embedding as follows:
    //
    //    type T struct {
    //        ent.Schema
    //    }
    //
    Interface interface {
        // Type is a dummy method, that is used in edge declaration.
        //
        // The Type method should be used as follows:
        //
        //    type S struct { ent.Schema }
        //
        //    type T struct { ent.Schema }
        //
        //    func (T) Edges() []ent.Edge {
        //        return []ent.Edge{
        //            edge.To("S", S.Type),
        //        }
        //    }
        //
        Type()
        // Fields returns the fields of the schema.
        Fields() []Field
        // Edges returns the edges of the schema.
        Edges() []Edge
        // Indexes returns the indexes of the schema.
        Indexes() []Index
        // Config returns an optional config for the schema.
        //
        // Deprecated: the Config method predates the Annotations method and it
        // is planned be removed in v0.5.0. New code should use Annotations instead.
        //
        //    func (T) Annotations() []schema.Annotation {
        //        return []schema.Annotation{
        //            entsql.Annotation{Table: "Name"},
        //        }
        //    }
        //
        Config() Config
        // Mixin returns an optional list of Mixin to extends
        // the schema.
        Mixin() []Mixin
        // Hooks returns an optional list of Hook to apply on
        // mutations.
        Hooks() []Hook
        // Policy returns the privacy policy of the schema.
        Policy() Policy
        // Annotations returns a list of schema annotations to be used by
        // codegen extensions.
        Annotations() []schema.Annotation
    }

    // Schema is the default implementation for the schema Interface.
    // It can be embedded in end-user schemas as follows:
    //
    //    type T struct {
    //        ent.Schema
    //    }
    //
    Schema struct {
        Interface
    }
)
```

まぁ，説明するより書いたほうが早いか。こんな感じでどうだろう。

```go:ent/schema/user.go
package schema

import (
    "time"

    "entgo.io/ent"
    "entgo.io/ent/schema/edge"
    "entgo.io/ent/schema/field"
)

// User holds the schema definition for the User entity.
type User struct {
    ent.Schema
}

// Fields of the User.
func (User) Fields() []ent.Field {
    return []ent.Field{
        field.String("username").
            MaxLen(63).
            NotEmpty().
            Unique(),
        field.Time("created_at").
            Default(time.Now),
        field.Time("updated_at").
            Default(time.Now),
    }
}

// Edges of the User.
func (User) Edges() []ent.Edge {
    return []ent.Edge{
        edge.To("files", BinaryFile.Type),
    }
}
```

```go:ent/schema/binaryfile.go
package schema

import (
    "time"

    "entgo.io/ent"
    "entgo.io/ent/schema/edge"
    "entgo.io/ent/schema/field"
)

// BinaryFile holds the schema definition for the BinaryFile entity.
type BinaryFile struct {
    ent.Schema
}

// Fields of the BinaryFile.
func (BinaryFile) Fields() []ent.Field {
    return []ent.Field{
        field.String("filename").
            NotEmpty().
            Unique(),
        field.Bytes("body").
            Optional().
            Nillable(),
        field.Time("created_at").
            Default(time.Now),
        field.Time("updated_at").
            Default(time.Now),
    }
}

// Edges of the BinaryFile.
func (BinaryFile) Edges() []ent.Edge {
    return []ent.Edge{
        edge.From("owner", User.Type).
            Ref("files").
            Unique(),
    }
}
```

ちょっと微妙だが一応スキーマを定義できたので，コード生成を行う。

```
$ go generate ./ent
```

これで ent/ ディレクトリ以下を見てみると

```
$ tree ent/
ent/
├── binaryfile
│   ├── binaryfile.go
│   └── where.go
├── binaryfile.go
├── binaryfile_create.go
├── binaryfile_delete.go
├── binaryfile_query.go
├── binaryfile_update.go
├── client.go
├── config.go
├── context.go
├── ent.go
├── enttest
│   └── enttest.go
├── generate.go
├── hook
│   └── hook.go
├── migrate
│   ├── migrate.go
│   └── schema.go
├── mutation.go
├── predicate
│   └── predicate.go
├── runtime
│   └── runtime.go
├── runtime.go
├── schema
│   ├── binaryfile.go
│   └── user.go
├── tx.go
├── user
│   ├── user.go
│   └── where.go
├── user.go
├── user_create.go
├── user_delete.go
├── user_query.go
└── user_update.go
```

おうふ。テーブル2つしかないのにエラい量あるな。ここまでできれば，前節の謎の entgo.io/ent/examples/start/ent パッケージを sample/ent に置き換えることができる。

これでようやくスタートラインに立った（笑）

[ent]: https://entgo.io/
[GORM]: https://gorm.io/ "GORM - The fantastic ORM library for Golang, aims to be developer friendly."
[Go]: https://go.dev/
[github.com/ent/ent]: https://github.com/ent/ent "ent/ent: An entity framework for Go"
[PostgreSQL]: https://www.postgresql.org/ "PostgreSQL: The world's most advanced open source database"
[database/sql]: https://pkg.go.dev/database/sql "sql package - database/sql - pkg.go.dev"
[github.com/jackc/pgx]: https://github.com/jackc/pgx "jackc/pgx: PostgreSQL driver and toolkit for Go"
