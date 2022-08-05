---
title: "GitHub Actions で Go パッケージの CI 作業を一通り行う" # 記事のタイトル
emoji: "💮" # アイキャッチとして使われる絵文字（1文字だけ）
type: "idea" # "tech" : 技術記事 / "idea" : アイデア記事
topics: ["go", "github"] # タグ。["markdown", "rust", "aws"] のように指定する
published: true # 公開設定（true で公開）
---

この記事では [GitHub] Actions を使って [Go] パッケージのリポジトリに対して以下の作業を自動化する方法を紹介する。

1. 依存パッケージに対する脆弱性検査 (push or pull request 時)
2. lint & test (push or pull request 時)
3. build & deploy (バージョンタグ付加時)

## 依存パッケージに対する脆弱性検査

依存パッケージの検査には [nancy] を使うのがよさげだ。公式の [GitHub] Action も用意されている。

https://github.com/sonatype-nexus-community/nancy
https://github.com/sonatype-nexus-community/nancy-github-action

[GitHub] Action で [nancy] を動かすには `.github/workflows/` ディレクトリに以下の内容の YAML ファイルを設置する。

```yaml:vulns.yml
name: vulns
on:
  push:
    branches:
      - main
  pull_request:
jobs:
  vulns:
    name: Vulnerability scanner
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v2
      - uses: actions/setup-go@v2
        with:
          go-version: ^1.19
      - name: WriteGoList
        run: go list -json -m all > go.list
      - name: Nancy
        uses: sonatype-nexus-community/nancy-github-action@main
```

これで pull request 時， `main` ブランチ[^br1] への push 時に脆弱性検査が走る。

[^br1]: 2020年10月から [GitHub の新規リポジトリの既定ブランチ名が `main` になった](https://text.baldanders.info/remark/2020/08/renaming-default-branch-name-in-github-repositries/ "GitHub リポジトリの既定ブランチ名が main になるらしい")。古いリポジトリでは `master` のままなのでご注意を。

私はモジュール情報の収集に拙作の [depm](https://github.com/goark/depm "goark/depm: Visualize depndency packages and modules") を使っている。こんな感じ。

```yaml
      - name: install depm
        run: go install github.com/goark/depm@latest
      - name: WriteGoList
        run: depm list --json > go.list
```

この辺はお好みでどうぞ。

## Lint & Test

[Go] の linter には [golangci-lint] がオススメだ。[golangci-lint] は go vet をはじめ複数の lint を集約して結果を表示してくれる優れものである。こちらも公式の [GitHub] Action が用意されている。

https://github.com/golangci/golangci-lint
https://github.com/golangci/golangci-lint-action

[GitHub] Action で [golangci-lint] を動かすには `.github/workflows/` ディレクトリに以下の内容の YAML ファイルを設置する。

```yaml:lint.yml
name: lint
on:
  push:
    branches:
      - main
  pull_request:

permissions:
  contents: read
  # Optional: allow read access to pull request. Use with `only-new-issues` option.
  # pull-requests: read
jobs:
  golangci:
    name: lint
    runs-on: ubuntu-latest
    steps:
      - uses: actions/setup-go@v3
        with:
          go-version: ^1.19
      - uses: actions/checkout@v3
      - name: golangci-lint
        uses: golangci/golangci-lint-action@v3
        with:
          # Optional: version of golangci-lint to use in form of v1.2 or v1.2.3 or `latest` to use the latest version
          version: latest

          # Optional: working directory, useful for monorepos
          # working-directory: somedir

          # Optional: golangci-lint command line arguments.
          # args: --issues-exit-code=0

          # Optional: show only new issues if it's a pull request. The default value is `false`.
          # only-new-issues: true

          # Optional: if set to true then the all caching functionality will be complete disabled,
          #           takes precedence over all other caching options.
          # skip-cache: true

          # Optional: if set to true then the action don't cache or restore ~/go/pkg.
          # skip-pkg-cache: true

          # Optional: if set to true then the action don't cache or restore ~/.cache/go-build.
          # skip-build-cache: true

      - name: testing
        run: go test -shuffle on ./...
```

これで pull request 時， `main` ブランチ[^br1] への push 時に [golangci-lint] が走る。ちなみに `steps` 項目の

```yaml
      - name: testing
        run: go test -shuffle on ./...
```

は [Go] のテストを実行している部分である。単純な `go test` ではなく `make` コマンド等を使った複雑なテストが必要ならもう少し色々と書く必要がある。

## Build & Deploy

Pure Go であれば [GoReleaser] を使えばクロス・コンパイルと Release ページへのデプロイまで自動でやってくれる。設定は `.goreleaser.yml` に書く[^gr1]。こちらも公式の [GitHub] Action が用意されている。

[^gr1]: [GoReleaser] の使い方等は割愛する。たぶんググったら日本語でも情報が出てくると思う。

https://github.com/goreleaser/goreleaser/
https://github.com/goreleaser/goreleaser-action

[GitHub] Action で [GoReleaser] を動かすには `.github/workflows/` ディレクトリに以下の内容の YAML ファイルを設置する。

```yaml:build.yml
name: build

on:
  push:
    tags:
      - v*
jobs:
  goreleaser:
    runs-on: ubuntu-latest
    steps:
      - name: Checkout
        uses: actions/checkout@v3
        with:
          fetch-depth: 0
      - name: Set up Go
        uses: actions/setup-go@v3
        with:
          go-version: ^1.19
      - name: Run GoReleaser
        uses: goreleaser/goreleaser-action@v3
        with:
          # either 'goreleaser' (default) or 'goreleaser-pro'
          distribution: goreleaser
          version: latest
          args: release --rm-dist
        env:
          GITHUB_TOKEN: ${{ secrets.GITHUB_TOKEN }}
          # Your GoReleaser Pro key, if you are using the 'goreleaser-pro' distribution
          # GORELEASER_KEY: ${{ secrets.GORELEASER_KEY }}
```

これでバージョンタグを打った際に [GoReleaser] によるクロス・コンパイルとデプロイが走る。

## [GitHub] Action バッヂを貼る

`README.md` などのドキュメントに [GitHub] Action の状態を表示するバッヂを貼り付けることができる。バッヂは以下の書式で指定する。

```markdown
[![Actions Status](https://github.com/{user}/{repo}/workflows/{action}/badge.svg)](https://github.com/{user}/{repo}/actions)
```

たとえば リポジトリ [`https://github.com/spiegel-im-spiegel/koyomi`](https://github.com/spiegel-im-spiegel/koyomi) であれば

```markdown
[![lint status](https://github.com/spiegel-im-spiegel/koyomi/workflows/lint/badge.svg)](https://github.com/spiegel-im-spiegel/koyomi/actions)
```

とすれば

[![lint status](https://github.com/spiegel-im-spiegel/koyomi/workflows/lint/badge.svg)](https://github.com/spiegel-im-spiegel/koyomi/actions)

のように表示される。ちなみに `{action}` の名前は YAML のファイル名ではなく先頭行の `name` 項目に対応している。

## 参考ページ

- [reviewdog-golangci-lint を使う](https://zenn.dev/ikawaha/articles/57384e8fc69c7b057f7f)
- [Go の CI を Github Actions に移行した](https://zenn.dev/ikawaha/articles/055cc7070ff0d12c5b10)
- [How to Add a GitHub Actions Badge to Your Project - DEV Community](https://dev.to/robdwaller/how-to-add-a-github-actions-badge-to-your-project-11ci)

- [Go 依存パッケージの脆弱性検査 | text.Baldanders.info](https://text.baldanders.info/golang/check-for-vulns-in-golang-dependencies/)
- [golangci-lint を GitHub Actions で使う | text.Baldanders.info](https://text.baldanders.info/golang/using-golangci-lint-action/)
- [GitHub Actions でクロス・コンパイル（GoReleaser 編） | text.Baldanders.info](https://text.baldanders.info/golang/cross-compiling-in-github-actions-with-goreleaser/)
- [Go のコードでも GitHub Code Scanning が使えるらしい | text.Baldanders.info](https://text.baldanders.info/remark/2020/10/github-code-scanning-with-golang/)
- [CI 用の GitHub Actions が諸々アップデートされていた | text.Baldanders.info](https://text.baldanders.info/golang/update-github-actions/)

[Go]: https://golang.org/ "The Go Programming Language"
[nancy]: https://github.com/sonatype-nexus-community/nancy "sonatype-nexus-community/nancy: A tool to check for vulnerabilities in your Golang dependencies, powered by Sonatype OSS Index"
[golangci-lint]: https://golangci-lint.run/
[GoReleaser]: https://goreleaser.com/ "GoReleaser | Deliver Go binaries as fast and easily as possible"
[GitHub]: https://github.com/ "GitHub"
