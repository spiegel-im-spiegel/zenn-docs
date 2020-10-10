---
title: "ゼロから始める Hugo — インストールから Hello World まで"
emoji: "😀"
type: "tech"
topics: [Hugo]
published: false
---
現在，[私の本家サイトのブログ機能を外出しにする作業](http://www.baldanders.info/spiegel/log2/000870.shtml)を行っていて，ツールとしては [Hugo] を使うことにしたのだが，初めて使うツールなので，いろいろ試しながら作業している。ここではその時のメモを公開する。

## Hugo のインストール

[Hugo] の実行モジュールは公式サイトから取得できるのだが Go 言語環境のある方は是非 `go get` から取得することをお勧めする。 [Hugo の repository](https://github.com/spf13/hugo) はサイトのドキュメントのビルド環境も同梱されているため色々と参考になる。

なお [Hugo] では非常に多くのパッケージを利用していて，取得の際には git の他に Mercurial が必要である（参照： 「[はじめての Go 言語 (on Windows) その3](http://qiita.com/spiegel-im-spiegel/items/a52a47942fd3946bb583)」）。

```shell
C:\workspace>mkdir hugo

C:\workspace>cd hugo

C:\workspace\hugo>SET GOPATH=C:\workspace\hugo

C:\workspace\hugo>go get -v github.com/spf13/hugo
github.com/spf13/hugo (download)
github.com/kardianos/osext (download)
github.com/spf13/afero (download)
github.com/spf13/cast (download)
github.com/spf13/jwalterweatherman (download)
github.com/spf13/cobra (download)
github.com/cpuguy83/go-md2man (download)
github.com/russross/blackfriday (download)
github.com/shurcooL/sanitized_anchor_name (download)
github.com/inconshreveable/mousetrap (download)
github.com/spf13/pflag (download)
github.com/spf13/fsync (download)
github.com/PuerkitoBio/purell (download)
github.com/opennota/urlesc (download)
github.com/miekg/mmark (download)
github.com/BurntSushi/toml (download)
Fetching https://gopkg.in/yaml.v2?go-get=1
Parsing meta tags from https://gopkg.in/yaml.v2?go-get=1 (status code 200)
get "gopkg.in/yaml.v2": found meta tag main.metaImport{Prefix:"gopkg.in/yaml.v2", VCS:"git", RepoRoot:"https://gopkg.in/yaml.v2"} at https://gopkg.in/yaml.v2?go-get=1
gopkg.in/yaml.v2 (download)
github.com/spf13/viper (download)
github.com/kr/pretty (download)
github.com/kr/text (download)
github.com/magiconair/properties (download)
github.com/mitchellh/mapstructure (download)
Fetching https://golang.org/x/text/transform?go-get=1
Parsing meta tags from https://golang.org/x/text/transform?go-get=1 (status code 200)
get "golang.org/x/text/transform": found meta tag main.metaImport{Prefix:"golang.org/x/text", VCS:"git", RepoRoot:"https://go.googlesource.com/text"} at https://golang.org/x/text/transform?go-get=1
get "golang.org/x/text/transform": verifying non-authoritative meta tag
Fetching https://golang.org/x/text?go-get=1
Parsing meta tags from https://golang.org/x/text?go-get=1 (status code 200)
golang.org/x/text (download)
Fetching https://golang.org/x/text/unicode/norm?go-get=1
Parsing meta tags from https://golang.org/x/text/unicode/norm?go-get=1 (status code 200)
get "golang.org/x/text/unicode/norm": found meta tag main.metaImport{Prefix:"golang.org/x/text", VCS:"git", RepoRoot:"https://go.googlesource.com/text"} at https://golang.org/x/text/unicode/norm?go-get=1
get "golang.org/x/text/unicode/norm": verifying non-authoritative meta tag
bitbucket.org/pkg/inflect (download)
github.com/dchest/cssmin (download)
github.com/eknkc/amber (download)
github.com/yosssi/ace (download)
github.com/spf13/nitro (download)
github.com/gorilla/websocket (download)
Fetching https://gopkg.in/fsnotify.v1?go-get=1
Parsing meta tags from https://gopkg.in/fsnotify.v1?go-get=1 (status code 200)
get "gopkg.in/fsnotify.v1": found meta tag main.metaImport{Prefix:"gopkg.in/fsnotify.v1", VCS:"git", RepoRoot:"https://gopkg.in/fsnotify.v1"} at https://gopkg.in/fsnotify.v1?go-get=1
gopkg.in/fsnotify.v1 (download)
github.com/shurcooL/sanitized_anchor_name
github.com/spf13/afero
github.com/inconshreveable/mousetrap
github.com/spf13/hugo/bufferpool
github.com/kr/text
github.com/kardianos/osext
github.com/spf13/jwalterweatherman
github.com/spf13/pflag
github.com/russross/blackfriday
github.com/opennota/urlesc
github.com/BurntSushi/toml
github.com/PuerkitoBio/purell
gopkg.in/yaml.v2
github.com/spf13/cast
github.com/kr/pretty
github.com/magiconair/properties
github.com/spf13/fsync
github.com/cpuguy83/go-md2man/md2man
github.com/spf13/hugo/hugofs
github.com/mitchellh/mapstructure
golang.org/x/text/transform
bitbucket.org/pkg/inflect
github.com/dchest/cssmin
github.com/miekg/mmark
github.com/eknkc/amber/parser
github.com/spf13/cobra
github.com/yosssi/ace
golang.org/x/text/unicode/norm
github.com/spf13/nitro
github.com/spf13/hugo/parser
github.com/spf13/viper
github.com/eknkc/amber
github.com/gorilla/websocket
github.com/spf13/hugo/utils
gopkg.in/fsnotify.v1
github.com/spf13/hugo/transform
github.com/spf13/hugo/watcher
github.com/spf13/hugo/livereload
github.com/spf13/hugo/helpers
github.com/spf13/hugo/source
github.com/spf13/hugo/target
github.com/spf13/hugo/tpl
github.com/spf13/hugo/hugolib
github.com/spf13/hugo/create
github.com/spf13/hugo/commands
github.com/spf13/hugo

C:\workspace\hugo>bin\hugo.exe version
Hugo Static Site Generator v0.15-DEV BuildDate: 2015-09-05T13:45:44+09:00
```

作成した実行モジュールは，作成された場所に `PATH` を通すか， `PATH` の通った場所にコピーすればよい。

## 作業環境の作成

[Hugo] 用の作業環境を作るには `hugo new site` コマンドを起動する。

```shell
C:>hugo new site C:\hugo-env\www

C:>cd C:\hugo-env\www

C:\hugo-env\www>tree /f .
C:\HUGO-ENV\WWW
│  config.toml
│
├─archetypes
├─content
├─data
├─layouts
└─static
```

ここで `config.toml` はサイト設定用のファイルである。初期値は以下のようになっている。

```toml:config.toml
baseurl = "http://replace-this-with-your-hugo-site.com/"
languageCode = "en-us"
title = "My New Hugo Site"
```

なお， TOML 形式のファイルの読み書きは以下が参考になる。

- [設定ファイル記述言語 TOML - Qiita](http://qiita.com/b4b4r07/items/77c327742fc2256d6cbe)

この中で必須なのは baseurl のみで，これがないとビルドエラーになる（ただし `--baseUrl` オプションを付ければ回避できる）。今回は `config.toml` を以下のようにセットする。

```toml:config.toml
baseurl = "http://hello.example.com/"
languageCode = "jp"
title = "Hello World!"
```

では，いきなりビルドしてみよう（えー）

```shell
C:\hugo-env\www>hugo
0 draft content
0 future content
0 pages created
0 paginator pages created
0 categories created
0 tags created
in 7 ms

C:\hugo-env\www>tree /f .
C:\HUGO-ENV\WWW
│  config.toml
│
├─archetypes
├─content
├─data
├─layouts
├─public
│      404.html
│      index.html
│      index.xml
│      sitemap.xml
│
└─static
```

`public` フォルダとその下に4つのファイルが作成されている（ちなみに出力先は `-d` または `--destination` オプションで変更できる）。現時点では `404.html` と `index.html` は空。

`index.xml` は feed 用のファイルで RSS 2.0 形式で書かれている（何故か atom の語彙を使っている）。

```xml:index.xml
<?xml version="1.0" encoding="utf-8" standalone="yes" ?>
<rss version="2.0" xmlns:atom="http://www.w3.org/2005/Atom">
  <channel>
    <title>Hello World!</title>
    <link>http://hello.example.com/</link>
    <description>Recent content on Hello World!</description>
    <generator>Hugo -- gohugo.io</generator>
    <language>jp</language>
    <atom:link href="http://hello.example.com/index.xml" rel="self" type="application/rss+xml" />

  </channel>
</rss>
```

`sitemap.xml` は文字通りサイトマップ記述ファイル（SEO 用？）。

```xml:sitemap.xml
<?xml version="1.0" encoding="utf-8" standalone="yes" ?>
<urlset xmlns="http://www.sitemaps.org/schemas/sitemap/0.9">

  <url>
    <loc>http://hello.example.com/</loc>
    <priority>0</priority>
  </url>

</urlset>
```

### サーバモード

[Hugo] では `hugo server` コマンドによりサーバモードで起動する。

```shell
C:\hugo-env\www>hugo server
0 draft content
0 future content
0 pages created
0 paginator pages created
0 categories created
0 tags created
in 7 ms
Serving pages from C:\hugo-env\www\public
Web Server is available at http://127.0.0.1:1313/
Press Ctrl+C to stop
```

この状態で `http://127.0.0.1:1313/` をブラウザで開けばいいのだが，現時点では `index.html` が空なので何も表示されない。ちなみに，この状態の `index.xml` は

```xml:index.xml
<?xml version="1.0" encoding="utf-8" standalone="yes" ?>
<rss version="2.0" xmlns:atom="http://www.w3.org/2005/Atom">
  <channel>
    <title>Hello World!</title>
    <link>http://localhost:1313/</link>
    <description>Recent content on Hello World!</description>
    <generator>Hugo -- gohugo.io</generator>
    <language>jp</language>
    <atom:link href="http://localhost:1313/index.xml" rel="self" type="application/rss+xml" />

  </channel>
</rss>
```

となっている。 URL が変更されていることにお気づきだろうか。 [Hugo] は動作モードによって URL を自動的に書き換えてくれるのでデバッグが容易である。ただしサーバモードの状態で本番環境に deploy しようとすると大変なことになるのでご注意を。 deploy する前には必ずビルドしよう（またはビルド結果の出力先を分ける）。

個人で作業する場合はこれで問題ないが，デバッグ用サーバを共有する場合には `http://localhost:1313/` では都合が悪いので `--bind` および `--port` オプションを付加する。

```shell
C:\hugo-env\www>hugo server --bind="hostname" --port=8080 --watch
0 draft content
0 future content
0 pages created
0 paginator pages created
0 tags created
0 categories created
in 8 ms
Watching for changes in C:\hugo-env\www/{data,content,layouts,static}
Serving pages from C:\hugo-env\www\public
Web Server is available at http://hostname:8080/
Press Ctrl+C to stop
```

これで `http://hostname:8080/` でアクセスできるようになる。ちなみに `--watch` オプションを付加すると，入力フォルダを監視して変更があればリコンパイルする。便利！

## はじめてのテンプレート

`index.html` を生成するにはテンプレートを用意する必要がある。 `index.html` のテンプレートは `layouts` フォルダに `index.html` というファイル名で配置する（安直！）

それじゃあ，早速テンプレートを組んでみる。

```html:layouts/index.html
<!DOCTYPE html>
{{ with .Site.LanguageCode }}<html lang="{{ . }}">{{ else }}<html>{{ end }}
<head>
<meta charset="utf-8">
<title>{{ .Title }}</title>
</head>
<body>
<h1>{{ .Title }}</h1>
<p>Hello Hugo World!</p>
</body>
<html>
```

`{{ }}` で囲まれている部分が埋め込み構文だ。たとえば `{{ .Title }}` にはタイトルが入る。今回はトップページなので， `config.toml` で `title` に設定した文字列が入る。

`{{ with .Site.LanguageCode }} ... {{ else }} ... {{ end }}` はいわゆる with 構文ってやつ。これ便利だよね。ちなみに `with` で指定する変数が存在しない場合は `{{ else }}` 以降が走る。似た書式で `{{ if .Site.LanguageCode }} ... {{ else }} ... {{ end }}` もあるが， `with` とは変数のスコープが異なるので注意（慣れるまでは結構ハマった）。

```html
{{ if .Site.LanguageCode }}<html lang="{{ .Site.LanguageCode }}">{{ else }}<html>{{ end }}
```

これをビルドした結果はこうなる。

```html:public/index.html
<!DOCTYPE html>
<html lang="jp">
<head>
<meta charset="utf-8">
<title>Hello World!</title>
</head>
<body>
<h1>Hello World!</h1>
<p>Hello Hugo World!</p>
</body>
<html>
```

ちなみに `config.toml` で `baseurl` 以外を削除してビルドするとこうなる。

```html:public/index.html
<!DOCTYPE html>
<html>
<head>
<meta charset="utf-8">
<title></title>
</head>
<body>
<h1></h1>
<p>Hello Hugo World!</p>
</body>
<html>
```

埋め込み構文がどのように機能しているか分かると思う。

### 404.html をどうしよう

`404.html` はページが存在しない場合（404）に代わりに表示するページである。普通 404 では Web サーバで決められたページを表示するのだが， GitHub Pages のようにサイトのトップページにある `404.html` を表示してくれる場合もある。なら `404.html` を空っぽにしておくわけにはいかないよねぇ。

固定の内容なら `static` フォルダ直下に `404.html` を置く手もある（[Hugo] では `static` フォルダ以下のファイルはそのままコピーされるが [Hugo] のコントロール外となる）。また， `layouts` フォルダ直下にテンプレートを置く手もある。たとえばこんな感じ。

```html:layouts/404.html
<!DOCTYPE html>
{{ with .Site.LanguageCode }}<html lang="{{ . }}">{{ else }}<html>{{ end }}
<head>
<meta charset="utf-8">
{{ .Hugo.Generator }}
<title>{{ .Title }}</title>
</head>
<body>{{ with .Site.Params.error404 }}<p>{{ . }}</p>{{ end }}</body>
</html>
```

これをビルドするとこんな感じになる。

```html:public/404.html
<!DOCTYPE html>
<html lang="jp">
<head>
<meta charset="utf-8">
<meta name="generator" content="Hugo 0.15-DEV" />
<title>404 Page not found</title>
</head>
<body><p>あると思った？ 残念！ 404 でした。</p></body>
</html>
```

## 記事を書く

さて，いよいよ記事を書いてみる。 Hugo では今どきの流行にのっとって Markdown 形式のファイルを入力ファイルにしている。自前でファイルを用意してもよいが，新規に作成するなら `hugo new` コマンドで作成するのがよいだろう。

```shell
C:\hugo-env\www>hugo new hello.md
C:\hugo-env\www\content\hello.md created
```

作成されたファイルを見てみると

```markdown:content\hello.md
+++
date = "2015-09-05T16:40:41+09:00"
draft = true
title = "hello"

+++
```

などとなっている。

`+++` で囲まれている部分は “front matter” と呼ばれている領域で，記事ごとの設定情報を格納する。 `+++` の場合は TOML， `---` の場合は YAML， `{ }` で挟んである場合は JSON として解釈されるようだ。

このファイルに記事を入力する。

```markdown:content\hello.md
+++
date = "2015-09-05T16:40:41+09:00"
draft = false
title = "Hello!"
+++

ようこそ， [Hugo](http://gohugo.io/) の世界へ！
```

草稿記事の場合は `draft` は `true` のままでよいが，そうでない場合は `false` にすること。ちなみに `-D` または `--buildDrafts` オプションを付けてビルドすると草稿版も含めて出力される。先ほどのサーバモードと組み合わせてデバッグに使うとよいだろう。

ではビルドしてみる（えー）

```shell
C:\hugo-env\www>hugo
0 draft content
0 future content
1 pages created
0 paginator pages created
0 tags created
0 categories created
in 42 ms

C:\hugo-env\www>tree /f .
C:\HUGO-ENV\WWW
│  config.toml
│
├─archetypes
├─content
│      hello.md
│
├─data
├─layouts
│      404.html
│      index.html
│
├─public
│  │  404.html
│  │  index.html
│  │  index.xml
│  │  sitemap.xml
│  │
│  └─hello
│          index.html
│
└─static
```

[Hugo] の標準では `hello.md` は `hello/index.html` に展開される。これを deploy すると Web ブラウザからは `http://hello.example.com/hello/` でアクセスできることになる。この path の展開のさせ方はいろいろ設定できるようなのだが，今回は割愛する。

現時点では `hello/index.html` は空。まぁテンプレートがないから当たり前なのだが。記事用のテンプレートは `layouts` フォルダ以下に `_default` フォルダを作成し，さらにその下に `single.html` を配置する。

それじゃあ，早速テンプレートを組んでみる。

```html:layouts/_default/single.html
<!DOCTYPE html>
{{ with .Site.LanguageCode }}<html lang="{{ . }}">{{ else }}<html>{{ end }}
<head>
<meta charset="utf-8">
<title>{{ .Title }} -- {{ .Site.Title }}</title>
</head>
<body>
<h1>{{ .Title }}</h1>
<div>{{ .Content }}</div>
</body>
<html>
```

`{{ .Title }}` には（今回は）記事のタイトルが入る。 `{{ .Content }}` は実際の記事の中身である。これでビルドすると

```html:public/hello/index.html
<!DOCTYPE html>
<html lang="jp">
<head>
<meta charset="utf-8">
<title>Hello! -- Hello World!</title>
</head>
<body>
<h1>Hello!</h1>
<div><p>ようこそ， <a href="http://gohugo.io/">Hugo</a> の世界へ！</p>
</div>
</body>
<html>
```

となる。 Feed はこんな感じ。

```xml:index.xml
<?xml version="1.0" encoding="utf-8" standalone="yes" ?>
<rss version="2.0" xmlns:atom="http://www.w3.org/2005/Atom">
  <channel>
    <title>Hello World!</title>
    <link>http://hello.example.com/</link>
    <description>Recent content on Hello World!</description>
    <generator>Hugo -- gohugo.io</generator>
    <language>jp</language>
    <lastBuildDate>Sat, 05 Sep 2015 16:40:41 +0900</lastBuildDate>
    <atom:link href="http://hello.example.com/index.xml" rel="self" type="application/rss+xml" />

    <item>
      <title>Hello!</title>
      <link>http://hello.example.com/hello/</link>
      <pubDate>Sat, 05 Sep 2015 16:40:41 +0900</pubDate>

      <guid>http://hello.example.com/hello/</guid>
      <description>&lt;p&gt;ようこそ， &lt;a href=&#34;http://gohugo.io/&#34;&gt;Hugo&lt;/a&gt; の世界へ！&lt;/p&gt;
</description>
    </item>

  </channel>
</rss>
```

ちなみに front matter をまるっと削除してビルドすると

```html:public/hello/index.html
<!DOCTYPE html>
<html lang="jp">
<head>
<meta charset="utf-8">
<title> -- Hello World!</title>
</head>
<body>
<h1></h1>
<div><p>ようこそ， <a href="http://gohugo.io/">Hugo</a> の世界へ！</p>
</div>
</body>
<html>
```

```xml:index.xml
<?xml version="1.0" encoding="utf-8" standalone="yes" ?>
<rss version="2.0" xmlns:atom="http://www.w3.org/2005/Atom">
  <channel>
    <title>Hello World!</title>
    <link>http://hello.example.com/</link>
    <description>Recent content on Hello World!</description>
    <generator>Hugo -- gohugo.io</generator>
    <language>jp</language>
    <atom:link href="http://hello.example.com/index.xml" rel="self" type="application/rss+xml" />

    <item>
      <title></title>
      <link>http://hello.example.com/hello/</link>
      <pubDate>Mon, 01 Jan 0001 00:00:00 +0000</pubDate>

      <guid>http://hello.example.com/hello/</guid>
      <description>&lt;p&gt;ようこそ， &lt;a href=&#34;http://gohugo.io/&#34;&gt;Hugo&lt;/a&gt; の世界へ！&lt;/p&gt;
</description>
    </item>

  </channel>
</rss>
```

のような感じになる。タイトルや日付の情報が欠落するがエラーにはならないようだ。

### 記事の一覧をつくる

ついでなので `index.html` に記事の一覧を表示するようにしよう。テンプレートはこんな感じで書く。

```html:layouts/index.html
<!DOCTYPE html>
{{ with .Site.LanguageCode }}<html lang="{{ . }}">{{ else }}<html>{{ end }}
<head>
<meta charset="utf-8">
<title>{{ .Title }}</title>
</head>
<body>
<h1>{{ .Title }}</h1>

<h2>What's New</h2>
<ul style="list-style:none;">
{{ range first 15 .Site.Pages }}
	<li><a href="{{ .Permalink }}">{{ .Title }}</a> (<time pubdate="{{ .Date.Format "2006-01-02" }}">{{ .Date.Format "2006-01-02" }}</time>){{ if .Draft }} #Draft{{ end }}</li>
{{ end }}
</ul>

</body>
<html>
```

これで最新の15記事を列挙できる。で，ビルド結果はこんな感じ。

```html:public/index.html
<!DOCTYPE html>
<html lang="jp">
<head>
<meta charset="utf-8">
<title>Hello World!</title>
</head>
<body>
<h1>Hello World!</h1>

<h2>What's New</h2>
<ul style="list-style:none;">

	<li><a href="http://hello.example.com/hello/">Hello!</a> (<time pubdate="2015-09-05">2015-09-05</time>)</li>

</ul>

</body>
<html>
```

これで大まかなテンプレートのイメージがつかめただろうか。

## 【おまけ】 HTML ファイルをそのまま突っ込んだらどうなるの？

というわけで，上で作成した `hello/index.html` を `hello2.html` として `content` フォルダに突っ込んでみた。

```shell
C:\hugo-env\www>hugo
0 draft content
0 future content
2 pages created
0 paginator pages created
0 tags created
0 categories created
in 16 ms

C:\hugo-env\www>tree /f .
C:\HUGO-ENV\WWW
│  config.toml
│
├─archetypes
├─content
│      hello.md
│      hello2.html
│
├─data
├─layouts
│  │  404.html
│  │  index.html
│  │
│  └─_default
│          single.html
│
├─public
│  │  404.html
│  │  index.html
│  │  index.xml
│  │  sitemap.xml
│  │
│  ├─hello
│  │      index.html
│  │
│  └─hello2
│          index.html
│
└─static
```

おおう。 `hello2/index.html` に変換されている。ルートの `index.html` ファイルはどうなっているだろう。

```html:public/index.html
<!DOCTYPE html>
<html lang="jp">
<head>
<meta charset="utf-8">
<title>Hello World!</title>
</head>
<body>
<h1>Hello World!</h1>

<h2>What's New</h2>
<ul style="list-style:none;">

	<li><a href="http://hello.example.com/hello/">Hello!</a> (<time pubdate="2015-09-05">2015-09-05</time>)</li>

	<li><a href="http://hello.example.com/hello2/"></a> (<time pubdate="0001-01-01">0001-01-01</time>)</li>

</ul>

</body>
<html>
```

おおう。残念な結果に。でもこれってもしかして front matter つけりゃいいのか。

```markdown:content\hello2.html
+++
date = "2015-09-05T18:00:00+09:00"
draft = false
title = "Hello! Part 2"
+++

<!DOCTYPE html>
<html lang="jp">
<head>
<meta charset="utf-8">
<title>Hello! Part 2 -- Hello World!</title>
</head>
<body>
<h1>Hello! Part 2</h1>
<div><p>再びようこそ， <a href="http://gohugo.io/">Hugo</a> の世界へ！</p>
</div>
</body>
<html>
```

んで，ビルドしてみる。

```html:public/index.html
<!DOCTYPE html>
<html lang="jp">
<head>
<meta charset="utf-8">
<title>Hello World!</title>
</head>
<body>
<h1>Hello World!</h1>

<h2>What's New</h2>
<ul style="list-style:none;">

	<li><a href="http://hello.example.com/hello2/">Hello! Part 2</a> (<time pubdate="2015-09-05">2015-09-05</time>)</li>

	<li><a href="http://hello.example.com/hello/">Hello!</a> (<time pubdate="2015-09-05">2015-09-05</time>)</li>

</ul>

</body>
<html>
```

おおお！

```html:public/hello2/index.html
<!DOCTYPE html>
<html lang="jp">
<head>
<meta charset="utf-8">
<title>Hello! Part 2 -- Hello World!</title>
</head>
<body>
<h1>Hello! Part 2</h1>
<div><!DOCTYPE html>
<html lang="jp">
<head>
<meta charset="utf-8">
<title>Hello! Part 2 -- Hello World!</title>
</head>
<body>
<h1>Hello! Part 2</h1>
<div><p>再びようこそ， <a href="http://gohugo.io/">Hugo</a> の世界へ！</p>
</div>
</body>
<html>
</div>
</body>
<html>
```

おおう orz なんてこったい。つまり，元の HTML ファイルからヘッダ等の要素を抜いて  front matter を付加すれば見た目も保持できて [Hugo] からも Controllable な状態になるってことか？

```markdown:content\hello2.html
+++
date = "2015-09-05T18:00:00+09:00"
draft = false
title = "Hello! Part 2"
+++

<p>再びようこそ， <a href="http://gohugo.io/">Hugo</a> の世界へ！</p>
```

これでビルドしてみると

```html:public/hello2/index.html
<!DOCTYPE html>
<html lang="jp">
<head>
<meta charset="utf-8">
<title>Hello! Part 2 -- Hello World!</title>
</head>
<body>
<h1>Hello! Part 2</h1>
<div><p>再びようこそ， <a href="http://gohugo.io/">Hugo</a> の世界へ！</p>
</div>
</body>
<html>
```

ふむむむむ。これは古い資産を [Hugo] に組み入れるのは結構骨かもなぁ。

## ブックマーク

- [Hugoサイト構築 | Watanabe-DENKI Inc. 渡辺電気株式会社](http://wdkk.co.jp/lab/hugo/) : お勧め！

[Hugo]: http://gohugo.io/ "Hugo :: A fast and modern static website engine"

