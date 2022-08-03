---
title: "Go でサブプロセスを起動する際は LookPath に気をつけろ！" # 記事のタイトル
emoji: "⌨" # アイキャッチとして使われる絵文字（1文字だけ）
type: "tech" # "tech" : 技術記事 / "idea" : アイデア記事
topics: ["go", "security"] # タグ。["markdown", "rust", "aws"] のように指定する
published: true # 公開設定（true で公開）
---

:::message
**【2022-08-03 追記】** [Go] 1.19 で [os/exec] 標準パッケージの挙動が変わり，この記事で述べる脆弱性を標準パッケージで回避できるようになった。詳しくは以下の拙文を参照のこと。

- [Go 1.19 で os/exec パッケージの挙動が変わった話 | text.Baldanders.info](https://text.baldanders.info/golang/exec-package-in-go119/)
:::

## Windows 環境における外部コマンド起動に関する脆弱性

先日 [Git for Windows] 2.29.2 (2) がリリースされたのだが，この中で [Git LFS] の脆弱性の修正が行われている。

- [Release Git for Windows 2.29.2(2) · git-for-windows/git · GitHub](https://github.com/git-for-windows/git/releases/tag/v2.29.2.windows.2)
- [Release v2.12.1 · git-lfs/git-lfs · GitHub](https://github.com/git-lfs/git-lfs/releases/tag/v2.12.1)
- [Git for Windows 2.29.2 (2) のリリース【セキュリテイ・アップデート】 | text.Baldanders.info](https://text.baldanders.info/release/2020/11/git-for-windows-2_29_2-2-is-released/)

この脆弱性は Windows 環境特有のもので

> On Windows, if Git LFS operates on a malicious repository with a git.bat or git.exe file in the current directory, that program is executed, permitting the attacker to execute arbitrary code.

Windows では `PATH` が通ってなくても（パスなしで指定すれば）カレント・フォルダの実行モジュールを起動できるので，そこに malware を紛れ込ませてユーザに起動させることが可能，というわけ。3年前くらいに流行った DLL 読み込みに関する脆弱性のバリエーションと考えると分かりやすいだろう。

- [Windows アプリケーションの DLL 読み込みに関する脆弱性について](https://www.jpcert.or.jp/tips/2017/wr172001.html)

## [os/exec] 標準パッケージ

もの知らずで申し訳ないが，私は今回の件まで [Git LFS] が [Go] 製とは知らなかった（笑） じゃあ [Go] では外部コマンドの呼び出しをどうやっているのかというと [os/exec] 標準パッケージが用意されている。

```go
package main

import (
    "fmt"
    "log"
    "os/exec"
)

func main() {
    out, err := exec.Command("date").Output()
    if err != nil {
        log.Fatal(err)
    }
    fmt.Printf("The date is %s\n", out)
}
```

（実行結果は[こちら](https://play.golang.org/p/XzRbRcDEbvH)）

この `exec.Command()` 関数の中身を見ると

```go
func Command(name string, arg ...string) *Cmd {
    cmd := &Cmd{
        Path: name,
        Args: append([]string{name}, arg...),
    }
    if filepath.Base(name) == name {
        if lp, err := LookPath(name); err != nil {
            cmd.lookPathErr = err
        } else {
            cmd.Path = lp
        }
    }
    return cmd
}
```

てな感じで，直接コマンド名を渡してるわけではなく，いったん `exec.LookPath()` 関数でパスに展開してから渡している。この関数が問題なのだ。

`exec.LookPath()` 関数は OS 毎に別実装になっていて，たとえば Windows では `lp_windows.go` というファイルにこんな感じで実装されている（長めのコードでゴメン）。

```go
func LookPath(file string) (string, error) {
    var exts []string
    x := os.Getenv(`PATHEXT`)
    if x != "" {
        for _, e := range strings.Split(strings.ToLower(x), `;`) {
            if e == "" {
                continue
            }
            if e[0] != '.' {
                e = "." + e
            }
            exts = append(exts, e)
        }
    } else {
        exts = []string{".com", ".exe", ".bat", ".cmd"}
    }

    if strings.ContainsAny(file, `:\/`) {
        if f, err := findExecutable(file, exts); err == nil {
            return f, nil
        } else {
            return "", &Error{file, err}
        }
    }
    if f, err := findExecutable(filepath.Join(".", file), exts); err == nil {
        return f, nil
    }
    path := os.Getenv("path")
    for _, dir := range filepath.SplitList(path) {
        if f, err := findExecutable(filepath.Join(dir, file), exts); err == nil {
            return f, nil
        }
    }
    return "", &Error{file, ErrNotFound}
}
```

注目は

```go
if f, err := findExecutable(filepath.Join(".", file), exts); err == nil {
    return f, nil
}
```

の部分で，パス指定のないコマンド名に対してわざわざカレント・フォルダ `.` を付加して優先的にチェックしてるのだ。なんちうおせっかいな `orz`

ちなみに UNIX 版（`lp_unix.go` ファイル）では，環境変数 `PATH` で明示しない限り，そんなことはしない。

```go
func LookPath(file string) (string, error) {
    if strings.Contains(file, "/") {
        err := findExecutable(file)
        if err == nil {
            return file, nil
        }
        return "", &Error{file, err}
    }
    path := os.Getenv("PATH")
    for _, dir := range filepath.SplitList(path) {
        if dir == "" {
            dir = "."
        }
        path := filepath.Join(dir, file)
        if err := findExecutable(path); err == nil {
            return path, nil
        }
    }
    return "", &Error{file, ErrNotFound}
}
```

拡張子のチェックもしないし，シンプルって素晴らしい！

## 脆弱性への対処

Windows 環境でパスの通っていないカレントのコマンドを安直に実行しないようにするには `exec.Command()` 関数にコマンド名を渡す前にパス名に展開するか， `exec.Cmd`インスタンスの `Path` 要素にパスに展開したコマンド名を直接セットするしかないだろう。 [Git LFS] では `LookPath()` 関数をカスタマイズしたものを実装し，直接 `Path` 要素をセットし直しているようだ。

というわけで [os/exec] パッケージでサブプロセス起動を正確に制御したい場合には `LookPath()` 関数に注意しましょう，ということで。

どっとはらい

## 【2020-12-20 追記】 [github.com/cli/safeexec][safeexec] パッケージを使う

[Hugo v0.79.1 のリリースノート](https://github.com/gohugoio/hugo/releases/tag/v0.79.1 "Release v0.79.1 · gohugoio/hugo")を見て気づいたのだが， GitHub が自身のコマンドライン・ツール用に [github.com/cli/safeexec][safeexec] パッケージを公開している。

これは [os/exec] 標準パッケージ内の `exec.LookPath()` 関数を置き換えるもので

```go
import (
    "os/exec"
    "github.com/cli/safeexec"
)

func gitStatus() error {
    gitBin, err := safeexec.LookPath("git")
    if err != nil {
        return err
    }
    cmd := exec.Command(gitBin, "status")
    return cmd.Run()
}
```

てな感じに使うようだ。

`exec.LookPath()` 関数周りでこれから対策を行うのであれば [github.com/cli/safeexec][safeexec] パッケージを使うことを検討してもいいだろう。

https://text.baldanders.info/golang/safeexec-packge/

## 【2021-01-21 追記】 [golang.org/x/sys/execabs][execabs] パッケージを使う

[Go] 1.15.7 でも今回の脆弱性について改修が行われた。

https://blog.golang.org/path-security

この中で [golang.org/x/sys/execabs][execabs] パッケージの[提案](https://github.com/golang/go/issues/43724 "proposal: os/exec: return error when PATH lookup would use current directory · Issue #43724 · golang/go")について言及されている。このパッケージでは

```go:execabs.go
func fixCmd(name string, cmd *exec.Cmd) {
    if filepath.Base(name) == name && !filepath.IsAbs(cmd.Path) {
        // exec.Command was called with a bare binary name and
        // exec.LookPath returned a path which is not absolute.
        // Set cmd.lookPathErr and clear cmd.Path so that it
        // cannot be run.
        lookPathErr := (*error)(unsafe.Pointer(reflect.ValueOf(cmd).Elem().FieldByName("lookPathErr").Addr().Pointer()))
        if *lookPathErr == nil {
            *lookPathErr = relError(name, cmd.Path)
        }
        cmd.Path = ""
    }
}

// CommandContext is like Command but includes a context.
//
// The provided context is used to kill the process (by calling os.Process.Kill)
// if the context becomes done before the command completes on its own.
func CommandContext(ctx context.Context, name string, arg ...string) *exec.Cmd {
    cmd := exec.CommandContext(ctx, name, arg...)
    fixCmd(name, cmd)
    return cmd

}

// Command returns the Cmd struct to execute the named program with the given arguments.
// See exec.Command for most details.
//
// Command differs from exec.Command in its handling of PATH lookups,
// which are used when the program name contains no slashes.
// If exec.Command would have returned an exec.Cmd configured to run an
// executable from the current directory, Command instead
// returns an exec.Cmd that will return an error from Start or Run.
func Command(name string, arg ...string) *exec.Cmd {
    cmd := exec.Command(name, arg...)
    fixCmd(name, cmd)
    return cmd
}
```

といった感じにパス指定のない外部コマンドについて絶対パスに展開されない場合はエラーとするようだ。

例えばカレントフォルダに

```Batch:hello.cmd
@echo off
echo Hello world!
```

というコマンドがあるとして。

```go:sample.go
package main

import (
    "os"
    "fmt"

    "golang.org/x/sys/execabs"
)

func main() {
    if b, err := execabs.Command("hello").Output(); err != nil {
        fmt.Fprintln(os.Stderr, err)
    } else {
        fmt.Println("Say:", string(b))
    }
}
```

とすれば

```
$ go run sample.go
hello resolves to executable in current directory (.\hello.cmd)
```

てな感じに実行されずにエラーになる。ちなみに外部コマンドを以下のように

```go:sample2.go
if b, err := execabs.Command(".\\hello").Output(); err != nil {
    ...
}
```

明示的にカレント・フォルダを指定すれば問題なく動作する。

```
$ go run sample2.go
Say: Hello world!
```

[execabs] パッケージが標準として組み込まれるかどうかは分からないが，暫定措置として置き換えを検討する価値はあるだろう。

[Go]: https://golang.org/ "The Go Programming Language"
[Git for Windows]: https://gitforwindows.org/ "Git for Windows"
[Git LFS]: https://git-lfs.github.com/ "Git Large File Storage | Git Large File Storage (LFS) replaces large files such as audio samples, videos, datasets, and graphics with text pointers inside Git, while storing the file contents on a remote server like GitHub.com or GitHub Enterprise."
[os/exec]: https://golang.org/pkg/os/exec/ "exec - The Go Programming Language"
[safeexec]: https://github.com/cli/safeexec "cli/safeexec: A safer version of exec.LookPath on Windows"
[execabs]: https://pkg.go.dev/golang.org/x/sys/execabs "execabs · pkg.go.dev"
<!-- eof -->
