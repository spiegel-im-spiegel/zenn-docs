# Install memo

## Install node.js (for Ubuntu)

```
curl -sL https://deb.nodesource.com/setup_current.x | sudo -E bash -
sudo apt install -y nodejs
```

## Install zenn-cli Package

```
$ npm init --yes
Wrote to /home/username/workspace/zenn-docs/package.json:

{
  "name": "zenn-docs",
  "version": "1.0.0",
  "description": "## Links",
  "main": "index.js",
  "scripts": {
    "test": "echo \"Error: no test specified\" && exit 1"
  },
  "repository": {
    "type": "git",
    "url": "git+https://github.com/spiegel-im-spiegel/zenn-docs.git"
  },
  "keywords": [],
  "author": "",
  "license": "ISC",
  "bugs": {
    "url": "https://github.com/spiegel-im-spiegel/zenn-docs/issues"
  },
  "homepage": "https://github.com/spiegel-im-spiegel/zenn-docs#readme"
}

$ npm install zenn-cli
...
+ zenn-cli@0.1.23
added 900 packages from 393 contributors and audited 903 packages in 66.098s

40 packages are looking for funding
  run `npm fund` for details

found 5 low severity vulnerabilities
  run `npm audit fix` to fix them, or `npm audit` for details

$ npx zenn init

  🎉Done!
  早速コンテンツを作成しましょう

  👇新しい記事を作成する
  $ zenn new:article

  👇新しい本を作成する
  $ zenn new:book

  👇表示をプレビューする
  $ zenn preview
```

## New Article

```
$ npx zenn new:article
📄d309af5057a827deda35.md created.
```

### Default Front Matter

```markdown
---
title: ""
emoji: "🎉"
type: "tech" # tech: 技術記事 / idea: アイデア
topics: []
published: true
---
```

なお，記事 URL のパス名となる slug は [GitHub] のリポジトリでコンテンツ管理していれば任意に指定できる。
Slug の制限は以下の通り。

- 半角英数字（a-z, 0-9）とハイフン（-）の 12〜50 字の組み合わせのみ有効
- `articles` 以下のファイルはディレクトリ階層に出来ない（フラットな構成）
- `books` の場合は「本」ごとに slug を指定できる。本の slug 以下はフラットな構成

[Zenn]: https://zenn.dev/ "Zenn｜プログラマーのための情報共有コミュニティ"
