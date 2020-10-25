---
title: "さようなら SHA-1" # 記事のタイトル
emoji: "🔏" # アイキャッチとして使われる絵文字（1文字だけ）
type: "idea" # "tech" : 技術記事 / "idea" : アイデア記事
topics: ["security", "hash", "sha1"] # タグ。["markdown", "rust", "aws"] のように指定する
published: true # 公開設定（true で公開）
---

先日 [git] v2.29 がリリースされた。

- [[ANNOUNCE] Git v2.29.0](https://lore.kernel.org/git/xmqqy2k2t77l.fsf@gitster.c.googlers.com/)
- [Highlights from Git 2.29 - The GitHub Blog](https://github.blog/2020-10-19-git-2-29-released/)

この中で特筆すべきは SHA-2 コミット・ハッシュの実験的サポートだろう。こんな感じで利用できるらしい。

```
$ git init --object-format=sha256 sample-repo
Initialized empty Git repository in /home/username/sample-repo/.git/

$ cd sample-repo
$ echo 'Hello, SHA-256!' >README.md
$ git add README.md
$ git commit -m "README.md: initial commit"
[main (root-commit) 6d45449] README.md: initial commit
 1 file changed, 1 insertion(+)
 create mode 100644 README.md

 $ git rev-parse HEAD
 6d45449028a8e76500adbfe7330e779d5dc4a3a14fca58ff08ec354c58727b2c
```

この記事では，これに関連すると思われる SHA-1 の危殆化について簡単に紹介する。

## Hash 値の衝突問題

暗号技術における hash 関数とは，以下の機能を持ったアルゴリズムである。

1. 任意のデータ列を一定の長さのデータ列（hash 値）に「要約」する
1. Hash 値から元のデータ列を推測できない
1. ひとつの hash 値に対して複数のデータ列が（実時間で）見つからない

Hash 関数はメッセージ認証符号（Message Authentication Code; MAC）や電子署名（digital signature）の中核技術のひとつであり，データの「完全性（Integrity）」を担保する重要な要素である。

特に3番目の「ひとつの hash 値に対して複数のデータ列が（実時間で）見つからない」という機能が破られると，その hash 関数では完全性を担保できなくなってしまう。これを「Hash 値の衝突問題」という。

## SHA-1 アルゴリズムの危殆化？

そもそもの発端は，2004年に複数の hash 関数において高い確率で hash 値を衝突させる攻略法が発表されたことだった。

- [Collisions for Hash Functions MD4, MD5, HAVAL-128 and RIPEMD](http://eprint.iacr.org/2004/199) : 発端となった論文。この中で SHA-0 も攻略可能であることが示されている

その後の研究で SHA-1 も攻略可能であることが分かってきて暗号技術の周辺は大騒動になった。もともと SHA-1 および SHA-2 の開発には NSA が関与しているとして評判がよくなかった。そこに上の論文の登場で SHA-1, SHA-2 に代わる hash アルゴリズムが求められるようになった。これは最終的に SHA-3 として NIST から勧告されている。

一方， SHA-1 アルゴリズム危殆化については，その後の研究で具体的な SHA-1 および SHA-2 の衝突例が見つからなかったこともあり，議論は縮小し，せっかく作った SHA-3 も SHA-2 のバックアップという位置付けに収まった。

- [『暗号をめぐる最近の話題』 — 旧メイン・ブログ | Baldanders.info](https://baldanders.info/blog/000586/)
- [SHA-3 が正式リリース： あれから10年も... — 旧メイン・ブログ | Baldanders.info](https://baldanders.info/blog/000865/)

## 2010年問題

NIST は SHA-1 の廃止スケジュールを発表し，2010年までに電子署名に使う hash アルゴリズムを SHA-2 (SHA-224/256/364/512) に移行するよう勧告した。

しかし，実際には移行は遅々として進まず NIST による移行スケジュールも2013年に延期された。でも，それから更に2年以上遅延するのだが（笑）


## SHA-1 の性能限界

2010年代に入ると SHA-1 の性能限界が議論されるようになる。その中で SHA-1 を力づくで攻略する論文が発表され始める。これには GPU をふんだんに使った専用ハードウェアやクラウド・サービスの台頭など巨大な計算機パワーの調達コストが下がり始めたことが背景にある。

- [SHA-1 衝突問題： 廃止の前倒し | text.Baldanders.info](https://text.baldanders.info/remark/2015/problem-of-sha1-collision/)
- [最初の SHA-1 衝突例 | text.Baldanders.info](https://text.baldanders.info/remark/2017/02/sha-1-collision/)

この中で決定的になったのは以下の論文である。

- [SHA-1 is a Shambles]
- [（何度目かの）さようなら SHA-1 | text.Baldanders.info](https://text.baldanders.info/remark/2020/01/sayonara-sha-1/)

この論文の注目点は

1. “chosen-prefix collision for SHA-1” なる手法により，衝突可能なデータを用意する際の自由度が高い
2. ハッシュ値を攻略する際の計算機パワーの調達コストが比較的実用的なレベルまで下がった

の2つ。特に2番目が重要で， “Nvidia GTX 1060GPU” × 900 の構成で2ヶ月ほどで攻略できたらしい。コストにして 45k USD だそうだ[^cost1]。

[^cost1]: 単純に 1 USD = 110 JPY とするなら 45k USD = 4.95M JPY ほど。まぁ五百万円以下で攻略できてしまうわけですな。

“[SHA-1 is a Shambles]” を受け， [GnuPG] では 2.2.18 から 2019-01-19 以降に鍵に付与された SHA-1 ベースの電子署名を全て削除することにした。

- [[Announce] GnuPG 2.2.18 released](https://lists.gnupg.org/pipermail/gnupg-announce/2019q4/000442.html)
- [[openpgp] Deprecating SHA1](https://mailarchive.ietf.org/arch/msg/openpgp/Rp-inhYKT8A9H5E34iLTrc9I0gc/)
- [GnuPG 2.2.18 リリース： さようなら SHA-1 | text.Baldanders.info](https://text.baldanders.info/release/2019/11/gnupg-2_2_18-is-released/)

また OpenSSH は SHA-1 を使用する “ssh-rsa” 公開鍵署名アルゴリズムを近い将来に無効にすると言っている。

- [OpenSSH 8.3 released (and ssh-rsa deprecation notice) [LWN.net]](https://lwn.net/Articles/821544/)

将来的には，かつての MD5 と同じく， SHA-1 はレガシー資産の参照のためだけに残されることになるだろう。

## [Git][git] のコミット・ハッシュは安全か？

これまで述べた SHA-1 危殆化の問題は主に電子署名に影響するもので，各種 MAC (Message Authentication Code) アルゴリズムや疑似乱数生成アルゴリズム等には影響しないとされている。 [Git][git] のコミット・ハッシュは単にコミットの IDentity として用いられているだけなので電子署名ほど完全性への要求は厳しくないと思われる。

しかし，今後新たな危殆化問題が浮上する可能性は否定できないし，そうなったときに代替手段がないのは心許ないだろう。そういった意味で，コミット・ハッシュのアルゴリズムとして SHA-2 を確保しておくのは意味があると思われる。欲を言えば将来的に SHA-3 にも対応できる余地も作っておいて欲しいところだが。


[git]: https://git-scm.com/
[SHA-1 is a Shambles]: https://sha-mbles.github.io/
[GnuPG]: https://gnupg.org/ "The GNU Privacy Guard"

## 参考

- [[ANNOUNCE] Git v2.29.1](https://lore.kernel.org/git/xmqq4kmlj9q9.fsf@gitster.c.googlers.com/)
- [Release Git for Windows 2.29.0 · git-for-windows/git · GitHub](https://github.com/git-for-windows/git/releases/tag/v2.29.0.windows.1)
- [Release Git for Windows 2.29.1 · git-for-windows/git · GitHub](https://github.com/git-for-windows/git/releases/tag/v2.29.1.windows.1)
- [「Git for Windows 2.29.0」が公開 ～セットアップ時にデフォルトブランチ名を設定可能 - 窓の杜](https://forest.watch.impress.co.jp/docs/news/1284871.html)
- [Git v2.29 がリリースされた | text.Baldanders.info](https://text.baldanders.info/release/2020/10/git-2_29-is-released/)

<!-- eof -->
