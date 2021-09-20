---
title: "GORM による Entity の設計"
---

今回は以下の関係を持つ2つのテーブルを作って基本的な CRUD (Create/Read/Update/Delete) を試してみることにする。

```mermaid
erDiagram

users ||--o{ binary_files : ""

users {
  integer id
  varchar username
  timestamp created_at
  timestamp updated_at
  timestamp deleted_at
}

binary_files {
  integer id
  integer user_id
  varchar filename
  bytea body
  timestamp created_at
  timestamp updated_at
  timestamp deleted_at
}

```
















[Go]: https://go.dev/
[PostgreSQL]: https://www.postgresql.org/ "PostgreSQL: The world's most advanced open source database"
[database/sql]: https://pkg.go.dev/database/sql "sql package - database/sql - pkg.go.dev"
[GORM]: https://gorm.io/ "GORM - The fantastic ORM library for Golang, aims to be developer friendly."
[github.com/simukti/sqldb-logger]: https://github.com/simukti/sqldb-logger "simukti/sqldb-logger: A logger for Go SQL database driver without modifying existing *sql.DB stdlib usage."
