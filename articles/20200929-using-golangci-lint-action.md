---
title: "golangci-lint を GitHub Actions で使う" # 記事のタイトル
emoji: "💮" # アイキャッチとして使われる絵文字（1文字だけ）
type: "idea" # "tech" : 技術記事 / "idea" : アイデア記事
topics: ["go", "programming", "github", "lint"] # タグ。["markdown", "rust", "aws"] のように指定する
published: true # 公開設定（true で公開）
---

:::message
この[自ブログの記事](https://text.baldanders.info/golang/using-golangci-lint-action/ "golangci-lint を GitHub Actions で使う — プログラミング言語 Go | text.Baldanders.info")とのマルチポストです。要するに宣伝です。
:::

[golangci-lint] は go vet をはじめ複数の lint を集約して結果を表示してくれる優れものである。
かつては GolangCI.com で GitHub と連携できていたのだが，[2020年4月でサービスが停止](https://medium.com/golangci/golangci-com-is-closing-d1fc1bd30e0e "GolangCI.com is closing. Dear customers of GolangCI.com, | by Denis Isaev | golangci | Medium")してしまい，寂しい限り。

と思っていたのだが，いつの間にか公式の GitHub Actions が用意されていた。気付かなんだよ。不覚。

- [golangci/golangci-lint-action: Official GitHub action for golangci-lint from it's authors](https://github.com/golangci/golangci-lint-action)

使い方は簡単。リポジトリの `.github/workflows/` ディレクトリに YAML ファイル（例えば `golangci-lint.yml`）を置き，以下のように記述する。

```yaml
name: golangci-lint
on:
  push:
    tags:
      - v*
    branches:
      - master
  pull_request:
jobs:
  golangci:
    strategy:
      matrix:
        go-version: [1.15.x]
        os: [ubuntu-latest, macos-latest, windows-latest]
    name: lint
    runs-on: ${{ matrix.os }}
    steps:
      - uses: actions/checkout@v2
      - name: golangci-lint
        uses: golangci/golangci-lint-action@v2
        with:
          # Required: the version of golangci-lint is required and must be specified without patch version: we always use the latest patch version.
          version: v1.31

          # Optional: working directory, useful for monorepos
          # working-directory: somedir

          # Optional: golangci-lint command line arguments.
          # args: --issues-exit-code=0

          # Optional: show only new issues if it's a pull request. The default value is `false`.
          # only-new-issues: true
```

これで pull request 時と `master` ブランチ[^br1] にバージョンタグを打った際に [golangci-lint] が走る。
[golangci-lint] は `matrix` の組み合わせで並列処理されるようだ。

[^br1]: 2020年10月から [GitHub の新規リポジトリの既定ブランチ名が `main` になるらしい](https://text.baldanders.info/remark/2020/08/renaming-default-branch-name-in-github-repositries/ "GitHub リポジトリの既定ブランチ名が main になるらしい")。ご注意を。

![Pull Request](https://text.baldanders.info/golang/using-golangci-lint-action/reviews-in-pr.png =500x)

よーし，うむうむ，よーし。

まぁ，プラットフォーム依存のコードでもない限り [Go] 最新バージョンの `ubuntu-latest` だけでいいと思うけどね。

## ブックマーク

- [golangci/golangci-lint: Fast linters Runner for Go](https://github.com/golangci/golangci-lint)
- [golangci-lint に叱られる — プログラミング言語 Go | text.Baldanders.info](https://text.baldanders.info/golang/donot-sleep-through-life/)

[Go]: https://golang.org/ "The Go Programming Language"
[golangci-lint]: https://golangci-lint.run/
