---
title: "DSN を環境変数として取り出す"
---

[PostgreSQL] の DNS (Data Source Name) は以下のような URI で与えられる。

```
postgres://username:password@hostname:port/databasename?option...
```

[Go] 標準の [database/sql] を使う場合でも，この URI をそのまま突っ込めば接続できるようだ。

```go
import (
    "database/sql"
    _ "github.com/lib/pq"
)

db, err := sql.Open("postgres", "postgres://username:password@hostname:port/databasename")
```

でも URI をそのままコード中にハード・コーディングするわけにはいかないので，これを環境変数として渡す方法を考える。

一番簡単なのは [github.com/joho/godotenv] パッケージを使う方法である。

[github.com/joho/godotenv] はプロジェクト・ルートにある .env ファイルの内容を読み込んで環境変数化してくれるパッケージだが，ここはもうひと捻りして [XDG Base Directory](https://specifications.freedesktop.org/basedir-spec/basedir-spec-latest.html) の仕組みを利用しよう。

標準の [os].UserConfigDir() 関数を使えばプラットフォーム毎に適切な設定ファイル置き場（ディレクトリ）を取得できる。

```go:os/file.go
// UserConfigDir returns the default root directory to use for user-specific
// configuration data. Users should create their own application-specific
// subdirectory within this one and use that.
//
// On Unix systems, it returns $XDG_CONFIG_HOME as specified by
// https://specifications.freedesktop.org/basedir-spec/basedir-spec-latest.html if
// non-empty, else $HOME/.config.
// On Darwin, it returns $HOME/Library/Application Support.
// On Windows, it returns %AppData%.
// On Plan 9, it returns $home/lib.
//
// If the location cannot be determined (for example, $HOME is not defined),
// then it will return an error.
func UserConfigDir() (string, error) {
...
```

今回は $XDG_CONFIG_HOME/elephantsql/env ファイルを作って以下の環境変数を設定する（実際の URI は [ElephantSQL] の詳細画面から取得できる）。

```ini:$XDG_CONFIG_HOME/elephantsql/env
ELEPHANTSQL_URL=postgres://username:password@hostname:port/databasename
```

多少なりとセキュリティに気をつけるなら env ファイルのパーミッションをきちんと設定すること。

これを取り出すコードは以下の通り。

```go
//go:build run
// +build run

package main

import (
    "fmt"
    "log"
    "os"

    "github.com/joho/godotenv"
    "github.com/spiegel-im-spiegel/gocli/config"
    "github.com/spiegel-im-spiegel/gocli/exitcode"
)

func Run() exitcode.ExitCode {
    // get PostgreSQL connection URL
    if err := godotenv.Load(config.Path("elephantsql", "env")); err != nil { // load ~/.config/elephantsql/env file
        log.Println(err)
        return exitcode.Abnormal
    }
    fmt.Println(os.Getenv("ELEPHANTSQL_URL")) // Output: postgres://username:password@hostname:port/databasename
    return exitcode.Normal
}

func main() {
    Run().Exit()
}
```

拙作の [github.com/spiegel-im-spiegel/gocli] パッケージは CLI アプリケーションを作る際の細々とした機能を収録している。たとえば

```go
path := config.Path("elephantsql", "env") //get '~/.config/elephantsql/env' path string
```

と書けば ~/.config/elephantsql/env ファイルのパスを取得できる（実際にアクセスするわけではない）。 [github.com/spiegel-im-spiegel/gocli] パッケージの使い方については以下の拙文を参考にどうぞ。

https://text.baldanders.info/release/gocli-package-for-golang/

[godotenv][github.com/joho/godotenv].Load() はアプリケーション起動時に一回だけ呼び出せればいいので init() 関数に含めてしまおう。

```go
func init() {
    //load ~/.config/elephantsql/env file
    if err := godotenv.Load(config.Path("elephantsql", "env")); err != nil {
        panic(err)
    }
}
```

init() 関数内ではエラーを返す先がないので panic() を投げるしかない。ご注意を。

[Go]: https://go.dev/
[PostgreSQL]: https://www.postgresql.org/ "PostgreSQL: The world's most advanced open source database"
[ElephantSQL]: https://www.elephantsql.com/ "ElephantSQL - PostgreSQL as a Service"
[database/sql]: https://pkg.go.dev/database/sql "sql package - database/sql - pkg.go.dev"
[os]: https://pkg.go.dev/os "os package - os - pkg.go.dev"
[github.com/joho/godotenv]: https://github.com/joho/godotenv "joho/godotenv: A Go port of Ruby's dotenv library (Loads environment variables from `.env`.)"
[github.com/spiegel-im-spiegel/gocli]: https://github.com/spiegel-im-spiegel/gocli "spiegel-im-spiegel/gocli: Minimal Packages for Command-Line Interface"
