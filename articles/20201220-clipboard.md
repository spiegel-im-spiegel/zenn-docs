---
title: "クリップボードのテキストを取り出す" # 記事のタイトル
emoji: "⌨" # アイキャッチとして使われる絵文字（1文字だけ）
type: "tech" # "tech" : 技術記事 / "idea" : アイデア記事
topics: ["go", "programming"] # タグ。["markdown", "rust", "aws"] のように指定する
published: true # 公開設定（true で公開）
---

今回も小ネタ。

たとえばクリップボードに「Hello, world!」と入っているときに，これを取り出して標準出力に出力することを考える。

[Go] でクリップボードへ読み書きできるパッケージとしては [github.com/atotto/clipboard][clipboard] が有名で，こんな感じに書ける。

```go:sample.go
package main

import (
    "fmt"
    "os"

    "github.com/atotto/clipboard"
)

func main() {
    s, err := clipboard.ReadAll()
    if err != nil {
        fmt.Fprintln(os.Stderr, err)
        return
    }
    fmt.Print(s)
}
```

これを実行すると

```
$ go run sample.go
Hello, world!
```

という具合にクリップボードの内容を取り出せるわけだ。

元々 Linux では xsel (または xclip) コマンドってのがあって，クリップボードからの取り出しも

```
$ xsel
Hello, world!
```

てな感じで取り出すことができる。じゃあ [github.com/atotto/clipboard][clipboard] パッケージの何が嬉しいかというと，マルチプラットフォーム対応になっていて，どの環境でも ReadAll/WriteAll 関数を呼び出すだけで機能するのよ[^xsel1]。

[^xsel1]: Linux では xsel または xclip コマンドが標準で入ってない。 Linux 等の UNIX 系環境で [github.com/atotto/clipboard][clipboard] パッケージを使う場合には，あらかじめ xsel または xclip コマンドをインストールしておく必要がある。

というわけで，拙作の [gpgpdump] に [github.com/atotto/clipboard][clipboard] パッケージを組み込んで，クリップボードから直接 ASCII armor テキストを読み込むようにした。

https://text.baldanders.info/release/2020/12/gpgpdump-v0_11_0-is-released/

クリップボード操作が [Go] コード内で気軽に使えるようになると何かと便利である。

[Go]: https://golang.org/ "The Go Programming Language"
[clipboard]: https://github.com/atotto/clipboard "atotto/clipboard: clipboard for golang"
[gpgpdump]: https://github.com/spiegel-im-spiegel/gpgpdump "spiegel-im-spiegel/gpgpdump: OpenPGP packet visualizer"
<!-- eof -->
