# URI からデータを取ってくるだけの簡単なお仕事

今回も仕事の合間の息抜きを！

自分で作って使ってるツール群について RESTful API なデータを取ってくるコードをあちこちで個別に書いていて，そろそろコピペにも飽きてきたので，「取ってくる」機能だけに特化した機能を独立パッケージとして切り出すことにした。

https://github.com/spiegel-im-spiegel/fetch

まずは GET にのみ注力する。

このパッケージを使って，たとえばこんな感じに書ける。

```go
package main

import (
    "context"
    "fmt"
    "io"
    "net/http"
    "os"

    "github.com/spiegel-im-spiegel/fetch"
)

func main() {
    githubUser := "spiegel-im-spiegel"
    u, err := fetch.URL("https://api.github.com/users/" + githubUser + "/gpg_keys")
    if err != nil {
        fmt.Fprintln(os.Stderr, err)
        return
    }
    resp, err := fetch.New(
        fetch.WithHTTPClient(&http.Client{}),
    ).Get(
        u,
        fetch.WithContext(context.Background()),
        fetch.WithRequestHeaderSet("Accept", "application/vnd.github.v3+json"),
    )
    if err != nil {
        fmt.Fprintln(os.Stderr, err)
        return
    }
    defer func() {
        // _, _ = io.Copy(ioutil.Discard, resp.Body)
        resp.Body.Close()
    }()
    if _, err := io.Copy(os.Stdout, resp.Body); err != nil {
        fmt.Fprintln(os.Stderr, err)
    }
}
```

（何をやってるかは以前に書いた「[GitHub に登録した OpenPGP 公開鍵を取り出す](https://zenn.dev/spiegel/articles/20201014-openpgp-pubkey-in-github)」を参照のこと）

OAuth みたいな認証機能があるわけじゃないし，需要がニッチ過ぎて汎用パッケージとしては使えないと思うけど，[net/http](https://golang.org/pkg/net/http/ "http - The Go Programming Language") パッケージを使う上でありがちなバグや脆弱性が出にくいように書いてるつもりではある。

## [cURL as DSL]

上のサンプル・コードは以下の curl コマンドラインとほぼ同等である。

```
$ curl "https://api.github.com/users/spiegel-im-spiegel/gpg_keys" -H "Accept: application/vnd.github.v3+json"
```

というより，まず「[API 仕様](https://docs.github.com/en/free-pro-team@latest/rest/reference/users#list-gpg-keys-for-a-user)」としての curl コマンド例があって，それに沿うようにコードを書くことが多い。

実は，その名もずばり “[cURL as DSL]” というツールというかサービスがあって，私は昔からお世話になっている。
“[cURL as DSL]” は curl のコマンドラインから Go, Python, node.js, Java, PHP, Vim script などのコードに変換してくれる。

もっとも “[cURL as DSL]” 自体はもうメンテナンスされてないようで，同様のサービスとして

https://curl.trillworks.com/

を勧めている。
こちらは Python, Ansible URI, MATLAB, Node.js, R, PHP, Strest, Go, Dart, JSON, Elixir, Rust と本当に多岐に及んでいる。

この手のツールを使うようになったきっかけは2015年に “[cURL as DSL]” を公開された渋川よしきさんの以下の記事で

http://blog.shibu.jp/article/115602749.html

特に

> このツールを発想したきっかけが、Google Chromeの開発者ツールの「Copy as cURL」というメニューです。cURL形式でコピーできるなら、cURLコマンドは多くのユーザが気軽に作れます。コミュニケーション手段になると思いました。例え冗長でも、自動生成できるというのは、非プログラマにとってはとてもありがたい選択肢になります。ExcelのVBAだって、FlashのJSFLだって、プログラマじゃない人がたくさん使っていますからね。
>
> また、PythonやらRubyやらnode.jsを使ったことがある人は、インタラクティブモードとかirbとかREPLといった環境の便利さはよく理解されていると思います。cURLが広まって、みんながウェブのドキュメントにcURLの擬似コードを書いてくれるようになれば、cURLはHTTP界におけるREPLになる可能性もあると思っています。

という部分は当時「激しく同意」してしまった。

ぶっちゃけ curl コマンドライン例と利用可能なパラメータの一覧表があればコード化できるもんね。
[PA-APIv5 のクライアント側パッケージを作った](https://text.baldanders.info/remark/2019/10/pa-api-v5/ "PA-API v5 への移行")ときも Amazon さんが curl のサンプルコードを公開してくれたおかげで理解が容易だった面は否めない（あと PHP のコード。 Java はコード詳細が隠蔽されていてワケが分からんかったw）。

というわけで，これから RESTful API 仕様を一般公開される場合は，是非とも curl で書いてください 🙇 いや，時代は GraphQL だろうけど（笑）

## 関連リンク

- [Goのnet/httpのkeep-aliveで気をつけること - Carpe Diem](https://christina04.hatenablog.com/entry/go-keep-alive)

[cURL as DSL]: https://shibukawa.github.io/curl_as_dsl/ "cURL as DSL — cURL as DSL 1.0 documentation"

## 参考図書

https://www.amazon.co.jp/dp/4908686033
<!-- eof -->
