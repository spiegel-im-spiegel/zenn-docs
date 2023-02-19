---
title: "GoReleaser の設定のいくつかが DEPRECATED になっていた" # 記事のタイトル
emoji: "🤔" # アイキャッチとして使われる絵文字（1文字だけ）
type: "idea" # "tech" : 技術記事 / "idea" : アイデア記事
topics: ["go"] # タグ。["markdown", "rust", "aws"] のように指定する
published: true # 公開設定（true で公開）
---

私は [GoReleaser] を使って [GitHub] に対してビルドとデプロイを行っているのだが[^ga1]，ビルド用の設定ファイル .goreleaser.yml については最初の頃に作ったのを（ほとんど設定を変えずに）使いまわしていた。で， [Go] 1.20 が[出た](https://text.baldanders.info/release/2023/02/go-1_20-is-released/ "Go 1.20 がリリースされた")のをきっかけに「RICS-V 用のバイナリも追加しよう」と色気を出したのだが，実際に手元で動かしてみたら，いくつかの設定が DEPRECATED になってるらしい。

[^ga1]: 詳しくは拙文「[GitHub Actions で Go パッケージの CI 作業を一通り行う](https://zenn.dev/spiegel/articles/20200929-using-golangci-lint-action)」を参照のこと。

この記事では非推奨になる（なった）設定のいくつかを覚え書きの形で残しておく。なお .goreleaser.yml ファイルの書きかた等についてはここでは言及しない。適当にググったら出てくると思う。 ChatGPT が吐くポエムに騙されないように（笑）

## --rm-dist オプションは非推奨

手元の環境で，たとえば

```
$ goreleaser release --snapshot --skip-publish --rm-dist
```

などと起動すると，いきなり

```
DEPRECATED: --rm-dist was deprecated in favor of --clean, check https://goreleaser.com/deprecations#-rm-dist for more details
```

と叱られる。おぅふ orz

指示された Web ページを見ると

> `--rm-dist` has been deprecated in favor of `--clean`.
（via “[Deprecation notices - GoReleaser](https://goreleaser.com/deprecations#-rm-dist)”）

と書かれていた。指示のとおりコマンドラインを

```
$ goreleaser release --snapshot --skip-publish --clean
```

などとすればOK。 GitHub Action の設定も変更する必要があるので，お忘れなく。

## rlcp を有効にしておく

さらにワーニングは続く。

```
DEPRECATED: `archives.rlcp` will be the default soon, check https://goreleaser.com/deprecations#archivesrlcp for more info
```

これも指示された Web ページを見ると以下のように書かれていた。

> This is not so much a deprecation property (yet), as it is a default behavior change.
> 
> The usage of relative longest common path (`rlcp`) on the destination side of archive files will be enabled by default by June 2023. Then, this option will be deprecated, and you will have another 6 months (until December 2023) to remove it.
（via “[Deprecation notices - GoReleaser](https://goreleaser.com/deprecations#archivesrlcp)”）

とりあえず .goreleaser.yml ファイルで

```yaml:.goreleaser.yml
archives:
-
  rlcp: true
```

と設定しておけば今までどおりの動作でワーニングは出なくなるし，鼻の先はOKのようだ。問題の先送りとも言う（笑）

## 名前の置き換えはテンプレートで行う

つぎ行ってみよう。

```
DEPRECATED: `archives.replacements` should not be used anymore, check https://goreleaser.com/deprecations#archivesreplacements for more info
```

たとえば出力するアーカイブ・ファイル名のテンプレートを

```yaml:.goreleaser.yml
archives:
  - rlcp: true
    format: tar.gz
    format_overrides:
      - goos: windows
        format: zip
    name_template: '{{ .ProjectName }}_{{ .Version }}_{{ .Os }}_{{ .Arch }}{{ if .Arm }}v{{ .Arm }}{{ end }}'
    replacements:
      darwin: Darwin
      linux: Linux
      windows: Windows
      freebsd: FreeBSD
      386: 32bit
      amd64: 64bit
      arm64: ARM64
      riscv64: RISCV
```

などと定義していたとする。これらのうち `replacements` の項目が DEPRECATED だと言っているようだ[^ga2]。じゃあ，どうやって文字列を置き換えるかというと [Go] のテンプレート機能[^ga3] を使えってことらしい。

[^ga2]: `replacements` の項目の非推奨化は `archives` だけでなく `nfpms`, `snapcrafts` でも同様だそうな。
[^ga3]: [Go] のテンプレート機能については[公式ドキュメント](https://pkg.go.dev/text/template "template package - text/template - Go Packages")を参照のこと。クセが強いので慣れるまでは大変かも。

[Go] のテンプレート機能を使って上の記述を書き直すと

```yaml:.goreleaser.yml
archives:
  - rlcp: true
    format: tar.gz
    format_overrides:
      - goos: windows
        format: zip
    name_template: >-
      {{ .ProjectName }}_
      {{- .Version }}_
      {{- if eq .Os "freebsd" }}FreeBSD
      {{- else }}{{ title .Os }}{{ end }}_
      {{- if eq .Arch "amd64" }}64bit
      {{- else if eq .Arch "386" }}32bit
      {{- else if eq .Arch "arm64" }}ARM64
      {{- else if eq .Arch "riscv64" }}RISCV
      {{- else }}{{ .Arch }}{{ if .Arm }}v{{ .Arm }}{{ end }}{{ end }}
```

という感じになる。この設定の下に [GoReleaser] で処理すると

```
  • archives
    • creating      archive=dist/product_1.0.0_Windows_armv6.zip
    • creating      archive=dist/product_1.0.0_Linux_64bit.tar.gz
    • creating      archive=dist/product_1.0.0_FreeBSD_armv6.tar.gz
    • creating      archive=dist/product_1.0.0_Darwin_ARM64.tar.gz
    • creating      archive=dist/product_1.0.0_Darwin_64bit.tar.gz
    • creating      archive=dist/product_1.0.0_Windows_ARM64.zip
    • creating      archive=dist/product_1.0.0_Linux_32bit.tar.gz
    • creating      archive=dist/product_1.0.0_Linux_armv6.tar.gz
    • creating      archive=dist/product_1.0.0_Linux_ARM64.tar.gz
    • creating      archive=dist/product_1.0.0_Windows_32bit.zip
    • creating      archive=dist/product_1.0.0_Windows_64bit.zip
    • creating      archive=dist/product_1.0.0_Linux_RISCV.tar.gz
    • creating      archive=dist/product_1.0.0_FreeBSD_64bit.tar.gz
    • creating      archive=dist/product_1.0.0_FreeBSD_ARM64.tar.gz
    • creating      archive=dist/product_1.0.0_FreeBSD_32bit.tar.gz
```

という感じにファイルが生成された。よーし，うむうむ，よーし。

[Go]: https://golang.org/ "The Go Programming Language"
[GoReleaser]: https://goreleaser.com/ "GoReleaser | Deliver Go binaries as fast and easily as possible"
[GitHub]: https://github.com/ "GitHub"
