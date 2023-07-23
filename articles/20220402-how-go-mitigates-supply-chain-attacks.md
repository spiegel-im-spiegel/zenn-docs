---
title: "Go はどのようにしてサプライチェーン攻撃を低減しているか" # 記事のタイトル
emoji: "🤔" # アイキャッチとして使われる絵文字（1文字だけ）
type: "idea" # "tech" : 技術記事 / "idea" : アイデア記事
topics: ["go", "security"] # タグ。["markdown", "rust", "aws"] のように指定する
published: true # 公開設定（true で公開）
---

[Go] 本家ブログが面白い記事を出してたので，かいつまんで紹介してみる。

https://go.dev/blog/supply-chain

## サプライチェーン攻撃とは

知らない方もいるかもしれないので一応説明すると，もともと「サプライチェーン」というのは原料調達から製造，物流，販売を経て顧客に渡るまでの事業の一連の流れ（chain）を指す言葉で，この流れを最適化することで生産性の向上やコストの低減を目指すというのが，いわゆる SCM (Supply Chain Management) と呼ばれるやつである。

これをソフトウェア開発に当てはめて，製品の企画・設計から製造して顧客に渡し，さらにその後の保守・運用を含めた流れもサプライチェーンと呼ぶことがある。さらにさらにソフトウェアのサプライチェーンの場合は複数のソフトウェアを組み合わせた新たなシステムを作って運用することも含まれる。XaaS 全盛の現代ではソフトウェア・サプライチェーンの管理はとても重要である。

で，この「複数のソフトウェアを組み合わせた新たなシステム」というのがセキュリティ管理的には曲者で，攻撃者から見れば組み合わせたソフトウェアのうち一番弱いところを突くことでシステム全体にダメージを与えることができる。これが「サプライチェーン攻撃」の最も大雑把な説明である。

サプライチェーン攻撃はシステムという大きな括りだけではなくソフトウェア製品単体でも起こり得る。昨今のソフトウェアは FOSS を含む多数のライブラリやフレームワークを内包しているため，そのうちのどれかに脆弱性が見つかれば，それは製品そのものの脆弱性となり，ひいてはその製品を使ったシステムの脆弱性となる。最近であれば昨年末に騒ぎになった log4j の脆弱性や最近発覚した Spring4Shell などもソフトウェア・サプライチェーンにダメージを与え得る。だから大きな騒ぎになっているのである。

https://text.baldanders.info/remark/2021/12/log4j-vulnerability/
https://text.baldanders.info/remark/2022/04/spring4shell/

## Go はどのようにしてサプライチェーン攻撃を低減しているか

前置きはこのくらいにして，多数のサードパーティ・パッケージの恩恵を受けているであろう [Go] 製ソフトウェアがどのようにしてサプライチェーン攻撃を低減しているか（しようとしているか）について紹介してみよう。

まず大前提として

> every dependency is unavoidably a trust relationship
*(via [How Go Mitigates Supply Chain Attacks](https://go.dev/blog/supply-chain))*

という点は押さえておくべきだろう。 log4j の脆弱性が発覚したときに件のライブラリのメンテナに批難を向ける声があったと聞くが，はっきり言って無茶な話である。私なら「なら使うな！」って言ってしまいそうである[^log4j]（笑） 提供する側も利用する側も互いの信頼の下に「協働」することで物事は上手く回るようになる。

[^log4j]: もっとも log4j のメンテナの方は「相応のお金くれるなら面倒見るヨ」ってスタンスのようなので（私から見れば）随分と優しい対応だと思うが（笑） 個人的には FOSS のフリーライドは構わないと思っているが，フリーライドするなら脆弱性に文句を言う筋合いはないし，利用するソフトウェアに対して継続的な保証や確証が欲しいならプロジェクトに積極的にコミットしていくしかない。何も支払わないし責任も丸投げだけど保証はしろ，などというもの言いが通用するのは小学生までである（笑）

### go.mod によるバージョンのロック

[Go] 1.11 以降でモジュール・モードが導入されたが，これによりインポートする外部パッケージのバージョンをコントロールできるようになった。この「外部パッケージのバージョン」を記述しているのが go.mod ファイルである。

[Go] のパッケージとモジュールについては以下の拙文を参照あれ。

https://zenn.dev/spiegel/articles/20210223-go-module-aware-mode

さらに [Go] 1.16 からはモジュール（パッケージ＋バージョン）の自動ダウンロードが禁止になった。外部パッケージをインポートする場合には `go get` または `go mod tidy` コマンドで明示的にモジュールをダウンロードする必要がある。これにより不用意に未知のバージョンを取り込むリスクは減るわけだ。

### go.sum によるモジュールの完全性（Integrity）の担保

[Go] のモジュール・モードではビルドに関わる依存モジュールのハッシュ値（SHA-256）を go.sum ファイルに記録している。これにより同一バージョンのパッケージのコードが密かに書き換えられた場合でも検知しやすくなる。このため go.sum ファイルは go.mod ファイルと共にリポジトリのコミットに含めておく必要がある。

さらに [Go] 1.13 リリースのタイミングでモジュールのミラーリングとチェックサム・データベースのサービスが正式に開始された。

https://go.dev/blog/module-mirror-launch

ダウンロードされるモジュールはこれらのサービスを経由しハッシュ値も記録・照会されるため，真っさらの状態からパッケージをインポートする場合でも（以前に誰かがインポートしたものであれば）ある程度の完全性を担保してくれる[^disable]。

![](https://go.dev/blog/module-mirror-launch/sumdb-protocol.png)
*via [Module Mirror and Checksum Database Launched - The Go Programming Language](https://go.dev/blog/module-mirror-launch)*

[^disable]: ミラーリングおよびチェックサム・データベースのサービス利用は一部制限したり無効にすることも可能。

これのおかげで同じバージョンでリリースのやり直しができなくなるが，致し方ないところである。

### 公式リポジトリなんて飾りです。偉い人には分からんのです

これは [Go] に特徴的なことだと思うが [Go] のエコシステムにはサードパーティのライブラリやフレームワークを一元管理する公式リポジトリは用意されていない。サードパーティのパッケージをインポートするのであれば提供者の VCS (git など) リポジトリから直接インポートする。

前節で紹介したミラーリング・サービスはユーザからはただのプロキシとして透過的にふるまう。これによりユーザが実際にモジュールをダウンロードする際に参照するのはミラーリング・サービスに限定されるため，攻撃ポイントを減らすことが期待できる。またミラーリング・サービス内部の動作はサンドボックスで隔離されているそうだ。

ユーザから見て公式リポジトリは便利な一方，攻撃者から見れば各パッケージのリポジトリに加えて公式リポジトリも攻撃対象にできるためリスク管理の観点からは微妙なのは確かである。まぁ，この辺の考え方は人によって賛否があるだろう。いかにも [Go] らしいシンプルさだとは思うが（笑）

### モノリスは砕けない

大事な点なので原文をそのまま抜き出しておく。ごめんペコン。

>It is an explicit security design goal of the Go toolchain that neither fetching nor building code will let that code execute, even if it is untrusted and malicious. This is different from most other ecosystems, many of which have first-class support for running code at package fetch time. These “post-install” hooks have been used in the past as the most convenient way to turn a compromised dependency into compromised developer machines, and to worm through module authors.
*(via [How Go Mitigates Supply Chain Attacks](https://go.dev/blog/supply-chain))*

[Go] の実行バイナリがダイナミックリンク等を必要とせずモノリシックな構造になっているのはセキュリティ上の意味もあるということやね[^embed]。

[^embed]: ダイナミックリンクのような仕組みが全く使えないというわけではない。たとえば「[Goで多数実行ファイルでディスク使用量を削減する](https://zenn.dev/nobonobo/articles/a8c07284247b64)」といったやり方もある。 [Go] のモノリシックな構造は計算機リソースが潤沢な現在では有効な割り切りだと思うが，組込みシステムでは計算機リソースの制約が大きい場合も多く，その辺のリスク・トレードオフについて頭を悩ませることになるだろう。

### “A little copying is better than a little dependency”

これって [Go] のことわざなんだそうな。何かあちこちにケンカ売ってる気がするが（笑）

https://www.youtube.com/watch?v=PAAkCSZUG1c

（この公演の9分30秒のあたり）

最重要ポイントはこれかな。

>the strongest mitigation will always be a small dependency tree
*(via [How Go Mitigates Supply Chain Attacks](https://go.dev/blog/supply-chain))*

個人の感想だけど，こういうのってリファクタリングを繰り返してチューニングしていくしかないんじゃないかな。その際の指針として “A little copying is better than a little dependency” を考慮しましょうねってことでいいだろう。

## 参考

https://security.googleblog.com/2023/04/supply-chain-security-for-go-part-1.html
https://security.googleblog.com/2023/06/supply-chain-security-for-go-part-2.html
https://security.googleblog.com/2023/07/supply-chain-security-for-go-part-3.html
https://google.github.io/osv.dev/
https://github.com/google/osv-scanner
https://www.amazon.co.jp/dp/B07FSBHS2V

[Go]: https://go.dev/ "The Go Programming Language"
